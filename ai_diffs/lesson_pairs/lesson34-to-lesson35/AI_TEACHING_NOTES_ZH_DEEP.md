# `lesson34 -> lesson35` 独立讲义（详细注释版）

这一组是一个很典型的“课程进入下一阶段”的分水岭。

因为它不再只是补一个小工具函数、修一个命名细节，而是同时推进了两条真正会改变系统行为的主线：

1. `RabbitMQ retry` 不再只是 broker 里有个工具函数，而是正式接入 consumer
2. `order` 仓储从内存实现，第一次切到了 MongoDB 持久化实现

再加上这一组还顺手把 `Viper` 配置初始化方式修了一遍，所以整组看上去会比较杂。

你如果直接机械看 diff，很容易觉得它“什么都改了一点”。

但如果按主线去看，这组其实非常清楚：

`系统正在从“课程演示级样例”走向“更像真实服务”的形态。`

字幕上也能看出两条线：
- `36.TXT`：Viper 配置初始化 / 路径问题
- `37.TXT`：Mongo 接入

而 `retry` 这条线，则是上一组 `HandleRetry(...)` 写出来之后，终于在这组真正接到 consumer 里。

所以你这次要把它看成一个“三件事同时收口”的 lesson：
- 配置初始化修正
- 消息失败重试接线
- 订单存储改成 MongoDB

## 1. 这组你应该怎么读

建议顺序：

1. [internal/order/infrastructure/consumer/consumer.go](/g:/shi/go_shop_second/internal/order/infrastructure/consumer/consumer.go)
2. [internal/payment/infrastructure/consumer/consumer.go](/g:/shi/go_shop_second/internal/payment/infrastructure/consumer/consumer.go)
3. [internal/common/broker/rabbitmq.go](/g:/shi/go_shop_second/internal/common/broker/rabbitmq.go)
4. [internal/common/config/viper.go](/g:/shi/go_shop_second/internal/common/config/viper.go)
5. [internal/common/config/global.yaml](/g:/shi/go_shop_second/internal/common/config/global.yaml)
6. [internal/order/adapters/order_mongo_repository.go](/g:/shi/go_shop_second/internal/order/adapters/order_mongo_repository.go)
7. [internal/order/service/application.go](/g:/shi/go_shop_second/internal/order/service/application.go)
8. [docker-compose.yml](/g:/shi/go_shop_second/docker-compose.yml)
9. [internal/order/domain/order/order.go](/g:/shi/go_shop_second/internal/order/domain/order/order.go)
10. [internal/order/app/command/create_order.go](/g:/shi/go_shop_second/internal/order/app/command/create_order.go)
11. [internal/order/main.go](/g:/shi/go_shop_second/internal/order/main.go)
12. [internal/payment/main.go](/g:/shi/go_shop_second/internal/payment/main.go)
13. [internal/stock/main.go](/g:/shi/go_shop_second/internal/stock/main.go)

为什么这样看：
- 先看 consumer，是因为 retry 终于真正接线了，这是最直接的行为变化。
- 再看 `rabbitmq.go`，你就知道 consumer 接的到底是什么能力。
- 然后看 `viper.go + global.yaml`，理解为什么这组开始需要认真解决“配置从哪读”的问题。
- 再看 `order_mongo_repository.go + application.go + docker-compose.yml`，理解持久化切换这条主线。
- 最后看 `NewPendingOrder` 和 `create_order.go`，理解 Mongo 接入时顺手做的领域清理。

## 2. 这组到底在解决什么问题

先说 lesson34 之前的状态。

那时系统有几个关键特点：

### 2.1 retry 逻辑已经存在，但没正式接进 consumer

也就是说：
- `broker.HandleRetry(...)` 有了
- 但是 `order consumer` 和 `payment consumer` 里很多地方还是 `// TODO: retry`
- 真正消费失败时，系统并没有统一走那套 retry 策略

### 2.2 配置初始化方式很脆弱

前面课程里很多地方想读配置时，都依赖你显式调用：
- `NewViperConfig()`

而且原先 `viper.AddConfigPath("../common/config")` 这种写法，本质上依赖的是“当前 working directory 恰好对”。

这在多模块、多入口、多服务场景下是很脆弱的。

### 2.3 `order` 的存储还是内存版

也就是说：
- 服务一重启，订单没了
- 并不是真正持久化
- 事务、查询、一致性这些现实问题还没进场

所以 lesson35 的推进方向很明确：

1. 把 retry 真正接上线
2. 把配置初始化方式从“手动 + 脆弱相对路径”修成更稳的方式
3. 把 `order` 仓储切到 MongoDB

这三件事组合起来，系统就开始从“课堂 demo”走向“能真正留住状态、能容忍失败、能多服务运行”的阶段。

## 3. 第一条主线：retry 终于正式接进 consumer

### 3.1 `internal/order/infrastructure/consumer/consumer.go`

#### 3.1.1 这个文件自己的原始 diff

```diff
diff --git a/internal/order/infrastructure/consumer/consumer.go b/internal/order/infrastructure/consumer/consumer.go
index 623bb50..511e0e8 100644
--- a/internal/order/infrastructure/consumer/consumer.go
+++ b/internal/order/infrastructure/consumer/consumer.go
@@ -40,25 +40,33 @@ func (c *Consumer) Listen(ch *amqp.Channel) {
 	var forever chan struct{}
 	go func() {
 		for msg := range msgs {
-			c.handleMessage(msg, q)
+			c.handleMessage(ch, msg, q)
 		}
 	}()
 	<-forever
 }
 
-func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
+func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
 	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
 	t := otel.Tracer("rabbitmq")
 	_, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
 	defer span.End()
 
+	var err error
+	defer func() {
+		if err != nil {
+			_ = msg.Nack(false, false)
+		} else {
+			_ = msg.Ack(false)
+		}
+	}()
+
 	o := &domain.Order{}
 	if err := json.Unmarshal(msg.Body, o); err != nil {
 		logrus.Infof("error unmarshal msg.body into domain.order, err = %v", err)
-		_ = msg.Nack(false, false)
 		return
 	}
-	_, err := c.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
+	_, err = c.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
@@ -69,11 +77,12 @@ func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
 	})
 	if err != nil {
 		logrus.Infof("error updating order, orderID = %s, err = %v", o.ID, err)
-		// TODO: retry
+		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
+			logrus.Warnf("retry_error, error handling retry, messageID=%s, err=%v", msg.MessageId, err)
+		}
 		return
 	}
 
 	span.AddEvent("order.updated")
-	_ = msg.Ack(false)
 	logrus.Info("order consume paid event success!")
 }
```

#### 3.1.2 旧代码做什么，新代码做什么

旧代码：
- `handleMessage` 没拿 `ch`
- 反序列化失败时立即 `Nack`
- 更新订单失败时只是打日志，`// TODO: retry`
- 成功时手动 `Ack`

新代码：
- `handleMessage` 现在把 `ch` 也传进来，因为 retry 需要拿 channel 重新 publish
- 统一用 `defer` 管理 `Ack/Nack`
- 更新订单失败时正式调用 `broker.HandleRetry(ctx, ch, &msg)`
- 成功就 `Ack`，错误就 `Nack`

一句话：

`order consumer 终于从“知道以后要重试”变成“真的开始走统一重试逻辑”。`

#### 3.1.3 关键代码和详细注释

```go
func (c *Consumer) Listen(ch *amqp.Channel) {
    ...
    go func() {
        for msg := range msgs {
            // 这次把 ch 也传进 handleMessage，
            // 因为 retry 逻辑不是本地处理一下就完，它需要重新 publish 消息。
            c.handleMessage(ch, msg, q)
        }
    }()
}

func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
    ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
    t := otel.Tracer("rabbitmq")
    _, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
    defer span.End()

    var err error
    defer func() {
        // 这里把 Ack/Nack 收口到 defer，很像前面 HTTP handler 把 Response 收口到 defer。
        // 成功就 Ack，失败就 Nack。
        if err != nil {
            _ = msg.Nack(false, false)
        } else {
            _ = msg.Ack(false)
        }
    }()

    o := &domain.Order{}
    if err := json.Unmarshal(msg.Body, o); err != nil {
        logrus.Infof("error unmarshal msg.body into domain.order, err = %v", err)
        return
    }

    _, err = c.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{...})
    if err != nil {
        logrus.Infof("error updating order, orderID = %s, err = %v", o.ID, err)

        // lesson35 的真正变化：
        // 失败时不再只是 TODO，而是开始正式走统一重试逻辑。
        if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
            logrus.Warnf("retry_error, error handling retry, messageID=%s, err=%v", msg.MessageId, err)
        }
        return
    }

    span.AddEvent("order.updated")
}
```

#### 3.1.4 为什么这里用 `defer` 统一 `Ack/Nack`

这和前面统一 HTTP response 的思想很像。

好处是：
- 成功和失败的收尾动作集中管理
- 减少到处散落的 `Ack()` / `Nack()`
- 逻辑更容易推理

但这里也有一个微妙点：
- `err` 变量现在同时承担“业务错误”和“retry 处理错误”
- 所以你得小心它在 defer 执行前的最终值到底是什么

这就是 Go 里“闭包 + defer + 可变变量”的典型思考点。

### 3.2 `internal/payment/infrastructure/consumer/consumer.go`

#### 3.2.1 这个文件自己的原始 diff

```diff
diff --git a/internal/payment/infrastructure/consumer/consumer.go b/internal/payment/infrastructure/consumer/consumer.go
index 99046de..611c5ae 100644
--- a/internal/payment/infrastructure/consumer/consumer.go
+++ b/internal/payment/infrastructure/consumer/consumer.go
@@ -38,33 +38,41 @@ func (c *Consumer) Listen(ch *amqp.Channel) {
 	var forever chan struct{}
 	go func() {
 		for msg := range msgs {
-			c.handleMessage(msg, q)
+			c.handleMessage(ch, msg, q)
 		}
 	}()
 	<-forever
 }
 
-func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
+func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
@@
+	var err error
+	defer func() {
+		if err != nil {
+			_ = msg.Nack(false, false)
+		} else {
+			_ = msg.Ack(false)
+		}
+	}()
@@
 	if _, err := c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{Order: o}); err != nil {
-		// TODO: retry
-		logrus.Infof("failed to create order, err=%v", err)
-		_ = msg.Nack(false, false)
+		logrus.Infof("failed to create payment, err=%v", err)
+		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
+			logrus.Warnf("retry_error, error handling retry, messageID=%s, err=%v", msg.MessageId, err)
+		}
 		return
 	}
@@
-	_ = msg.Ack(false)
 	logrus.Info("consume success")
 }
```

#### 3.2.2 这段和 order consumer 的关系

它本质上是同一类变化：
- 失败处理从 `TODO` 变成真实调用 `HandleRetry(...)`
- `Ack/Nack` 收口到 `defer`
- `handleMessage` 需要拿到 `ch`

也就是说，这组不是只修一个 consumer，而是把两条主要异步链的消费端都正式接进 retry 逻辑。

#### 3.2.3 关键代码和详细注释

```go
func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
    logrus.Infof("Payment receive a message from %s, msg=%v", q.Name, string(msg.Body))

    ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
    tr := otel.Tracer("rabbitmq")
    _, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
    defer span.End()

    var err error
    defer func() {
        if err != nil {
            _ = msg.Nack(false, false)
        } else {
            _ = msg.Ack(false)
        }
    }()

    o := &orderpb.Order{}
    if err := json.Unmarshal(msg.Body, o); err != nil {
        logrus.Infof("failed to unmarshall msg to order, err=%v", err)
        return
    }

    if _, err := c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{Order: o}); err != nil {
        logrus.Infof("failed to create payment, err=%v", err)

        // 和 order consumer 一样，lesson35 开始正式接线。
        if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
            logrus.Warnf("retry_error, error handling retry, messageID=%s, err=%v", msg.MessageId, err)
        }
        return
    }

    span.AddEvent("payment.created")
    logrus.Info("consume success")
}
```

#### 3.2.4 为什么这一步很重要

因为只有当 consumer 真正接上 retry 逻辑，上一组写的 `HandleRetry(...)` 才不再只是“工具函数躺在那儿”。

这一组的实际系统行为变化就是：

`消费失败时，系统终于开始真正重试，而不是只记录“以后再做”。`

## 4. 第二条主线：把 retry 配置化，而不是写死

### 4.1 `internal/common/config/global.yaml`

#### 4.1.1 这个文件自己的原始 diff

```diff
diff --git a/internal/common/config/global.yaml b/internal/common/config/global.yaml
index 6b3d5eb..099b0f5 100644
--- a/internal/common/config/global.yaml
+++ b/internal/common/config/global.yaml
@@ -30,6 +30,16 @@ rabbitmq:
   password: guest
   host: 127.0.0.1
   port: 5672
+  max-retry: 3
+
+mongo:
+  user: root
+  password: password
+  host: 127.0.0.1
+  port: 27017
+  db-name: "order"
+  coll-name: "order"
```

#### 4.1.2 这段在做什么

它一次性补了两块配置：

1. `rabbitmq.max-retry`
- 把上一组写死的 retry 次数正式变成配置

2. `mongo.*`
- 为 order 服务切 Mongo 做准备

这也说明这一组的主题确实不止一条线。

### 4.2 `internal/common/broker/rabbitmq.go`

#### 4.2.1 这个文件自己的原始 diff

```diff
diff --git a/internal/common/broker/rabbitmq.go b/internal/common/broker/rabbitmq.go
index d8a3e84..ef74e7d 100644
--- a/internal/common/broker/rabbitmq.go
+++ b/internal/common/broker/rabbitmq.go
@@ -5,8 +5,10 @@ import (
 	"fmt"
 	"time"
 
+	_ "github.com/ghost-yu/go_shop_second/common/config"
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
+	"github.com/spf13/viper"
 	"go.opentelemetry.io/otel"
 )
@@ -17,8 +19,7 @@ const (
 )
 
 var (
-	//maxRetryCount = viper.GetInt64("rabbitmq.max-retry")
-	maxRetryCount int64 = 3
+	maxRetryCount = viper.GetInt64("rabbitmq.max-retry")
 )
@@ -63,6 +64,7 @@ func createDLX(ch *amqp.Channel) error {
 }
 
 func HandleRetry(ctx context.Context, ch *amqp.Channel, d *amqp.Delivery) error {
+	logrus.Info("handleretry_max-retry-count", maxRetryCount)
```

#### 4.2.2 旧代码做什么，新代码做什么

旧代码：
- `maxRetryCount` 写死成 `3`

新代码：
- 通过 `viper.GetInt64("rabbitmq.max-retry")` 从配置里读
- 还加了 `_ ".../common/config"`，通过副作用导入保证配置先初始化

#### 4.2.3 关键代码和详细注释

```go
import (
    _ "github.com/ghost-yu/go_shop_second/common/config"
    ...
    "github.com/spf13/viper"
)

var (
    maxRetryCount = viper.GetInt64("rabbitmq.max-retry")
)
```

这里有两个重点。

##### 重点 1：为什么是副作用导入

```go
_ "github.com/.../common/config"
```

这个写法表示：
- 我不直接用包里的导出符号
- 但我希望这个包的 `init()` 被执行

而 `common/config` 在这一组里已经把 `init()` 变成了自动初始化 Viper。

所以这里的含义是：

`只要 broker 包被加载，config 包就会先把 Viper 配好。`

##### 重点 2：为什么要配置化

因为 retry 次数本来就不该硬编码死在业务代码里。

不同环境可能想要：
- 本地调试：`1` 或 `2`
- 测试环境：`3`
- 生产环境：可能更保守或更激进

所以 lesson35 把它拉到 `global.yaml`，这是对的。

#### 4.2.4 这一行 debug 日志说明什么

```go
logrus.Info("handleretry_max-retry-count", maxRetryCount)
```

这行明显带一点“调试味”。

它说明作者在课程推进时，想先验证：
- 配置到底有没有被正确读到

你要能分辨这种代码：
- 它在课程演示阶段很正常
- 但在更成熟版本里通常会删掉或者换成更规范的 structured log

## 5. 第三条主线：Viper 初始化方式终于被修正

### 5.1 `internal/common/config/viper.go`

#### 5.1.1 这个文件自己的原始 diff

```diff
diff --git a/internal/common/config/viper.go b/internal/common/config/viper.go
index a0c91d2..4ba14c7 100644
--- a/internal/common/config/viper.go
+++ b/internal/common/config/viper.go
@@ -1,17 +1,52 @@
 package config
 
 import (
+	"fmt"
+	"os"
+	"path/filepath"
+	"runtime"
 	"strings"
+	"sync"
 
 	"github.com/spf13/viper"
 )
 
-func NewViperConfig() error {
+func init() {
+	if err := NewViperConfig(); err != nil {
+		panic(err)
+	}
+}
+
+var once sync.Once
+
+func NewViperConfig() (err error) {
+	once.Do(func() {
+		err = newViperConfig()
+	})
+	return
+}
+
+func newViperConfig() error {
+	relPath, err := getRelativePathFromCaller()
+	if err != nil {
+		return err
+	}
 	viper.SetConfigName("global")
 	viper.SetConfigType("yaml")
-	viper.AddConfigPath("../common/config")
+	viper.AddConfigPath(relPath)
 	viper.EnvKeyReplacer(strings.NewReplacer("_", "-"))
 	viper.AutomaticEnv()
 	_ = viper.BindEnv("stripe-key", "STRIPE_KEY", "endpoint-stripe-secret", "ENDPOINT_STRIPE_SECRET")
 	return viper.ReadInConfig()
 }
+
+func getRelativePathFromCaller() (relPath string, err error) {
+	callerPwd, err := os.Getwd()
+	if err != nil {
+		return
+	}
+	_, here, _, _ := runtime.Caller(0)
+	relPath, err = filepath.Rel(callerPwd, filepath.Dir(here))
+	fmt.Printf("caller from: %s, here: %s, relpath: %s", callerPwd, here, relPath)
+	return
+}
```

#### 5.1.2 旧代码的问题是什么

旧代码最大的问题是这行：

```go
viper.AddConfigPath("../common/config")
```

它隐含假设：

`当前工作目录相对于 common/config 的位置永远固定。`

这在单一入口、固定启动方式下也许还能侥幸工作。

但当你有：
- order
- payment
- stock
- 不同模块目录
- 不同启动入口

这个假设很容易失效。

#### 5.1.3 新代码在做什么

新代码做了三件事：

1. 用 `init()` 自动初始化 Viper
2. 用 `sync.Once` 确保只初始化一次
3. 通过 `os.Getwd()` + `runtime.Caller(0)` + `filepath.Rel(...)` 动态算相对路径

### 5.2 关键代码和详细注释

```go
func init() {
    if err := NewViperConfig(); err != nil {
        panic(err)
    }
}

var once sync.Once

func NewViperConfig() (err error) {
    once.Do(func() {
        err = newViperConfig()
    })
    return
}
```

这里你要看懂两个点。

#### `init()`
- 包一被导入就自动执行
- 这里的意思是：配置初始化变成包级副作用

#### `sync.Once`
- 确保即使多个地方调用 `NewViperConfig()`，真正初始化逻辑也只跑一次

这在配置初始化、logger 初始化这种“全局单次初始化”场景里非常常见。

继续看路径部分：

```go
func getRelativePathFromCaller() (relPath string, err error) {
    callerPwd, err := os.Getwd()
    if err != nil {
        return
    }

    _, here, _, _ := runtime.Caller(0)

    // callerPwd: 当前工作目录
    // here: 当前 viper.go 这个文件所在绝对路径
    // 然后计算从当前工作目录到配置目录的相对路径
    relPath, err = filepath.Rel(callerPwd, filepath.Dir(here))

    fmt.Printf("caller from: %s, here: %s, relpath: %s", callerPwd, here, relPath)
    return
}
```

这里你要真正理解的不是 API 名，而是思路：

`既然固定写死相对路径不稳，那就运行时动态计算。`

### 5.3 为什么这个改动很重要

因为 lesson35 以后，不只是 main 会读配置了：
- broker 包想读 `rabbitmq.max-retry`
- mongo repository 想读 `mongo.db-name`
- 各种包都可能通过 `viper.GetString(...)` 拿配置

所以“配置初始化必须稳”这件事开始变成基础设施问题，不再只是 main 里的小工具。

### 5.4 这里有哪些局限

#### 局限 1：`fmt.Printf(...)` 明显是调试代码

这行打印调用路径，在课程演示阶段很有帮助，但长期应该删掉或换日志。

#### 局限 2：`init()` 虽方便，但也会带来隐式副作用

课程里这么写是为了省心。
但更大型项目里，很多人会更偏向：
- 显式初始化顺序
- 在 main 中统一 orchestrate

因为 `init()` 太多时，顺序和副作用不透明。

所以你要理解：
- 课程这里是一个实用折中
- 不是所有项目都该滥用 `init()`

## 6. 第四条主线：Order 仓储第一次切到 MongoDB

### 6.1 `docker-compose.yml`

#### 6.1.1 这个文件自己的原始 diff

```diff
diff --git a/docker-compose.yml b/docker-compose.yml
index 12f1e17..ac2b22f 100644
--- a/docker-compose.yml
+++ b/docker-compose.yml
@@ -22,4 +22,24 @@ services:
       COLLECTOR_OTLP_ENABLED: true
+
+  order-mongo:
+    image: "mongo:7.0.8"
+    restart: always
+    environment:
+      MONGO_INITDB_ROOT_USERNAME: root
+      MONGO_INITDB_ROOT_PASSWORD: password
+    ports:
+      - "27017:27017"
+
+  mongo-express:
+    image: "mongo-express"
+    restart: always
+    ports:
+      - "8082:8081"
+    environment:
+      ME_CONFIG_MONGODB_ADMINUSERNAME: root
+      ME_CONFIG_MONGODB_ADMINPASSWORD: password
+      ME_CONFIG_MONGODB_URL: mongodb://root:password@order-mongo:27017/
+      ME_CONFIG_BASICAUTH: false
```

#### 6.1.2 这段在做什么

很直接：
- 增加 MongoDB 服务
- 增加 `mongo-express` 可视化管理页面

这说明从这一组开始：

`order` 的数据不再只是存在进程内内存，而是准备落到真正数据库里。`

### 6.2 `internal/order/adapters/order_mongo_repository.go`

#### 6.2.1 这个文件自己的原始 diff

这是一整个新文件，我只摘核心段落：

```diff
+type OrderRepositoryMongo struct {
+	db *mongo.Client
+}
+
+func NewOrderRepositoryMongo(db *mongo.Client) *OrderRepositoryMongo {
+	return &OrderRepositoryMongo{db: db}
+}
+
+type orderModel struct {
+	MongoID     primitive.ObjectID `bson:"_id"`
+	ID          string             `bson:"id"`
+	CustomerID  string             `bson:"customer_id"`
+	Status      string             `bson:"status"`
+	PaymentLink string             `bson:"payment_link"`
+	Items       []*entity.Item     `bson:"items"`
+}
+
+func (r *OrderRepositoryMongo) Create(ctx context.Context, order *domain.Order) (created *domain.Order, err error) {
+    ...
+}
+
+func (r *OrderRepositoryMongo) Get(ctx context.Context, id, customerID string) (got *domain.Order, err error) {
+    ...
+}
+
+func (r *OrderRepositoryMongo) Update(...)
```

#### 6.2.2 这段的意义是什么

这是 order 仓储第一次真正落地到数据库。

也就是说，前面那些：
- `CreateOrder`
- `UpdateOrder`
- `GetCustomerOrder`

现在终于不是只操作一个进程内 map 了。

这会带来本质变化：
- 数据重启后还能保留
- 多实例共享数据开始有可能
- 事务、一致性、查询条件这些真实问题开始进入系统

#### 6.2.3 关键代码和详细注释

```go
type OrderRepositoryMongo struct {
    db *mongo.Client
}

func NewOrderRepositoryMongo(db *mongo.Client) *OrderRepositoryMongo {
    return &OrderRepositoryMongo{db: db}
}

func (r *OrderRepositoryMongo) collection() *mongo.Collection {
    return r.db.Database(dbName).Collection(collName)
}
```

这三段先做的是最基本的仓储适配器包装：
- 持有 `mongo.Client`
- 根据配置拿到目标 database + collection

再看模型：

```go
type orderModel struct {
    MongoID     primitive.ObjectID `bson:"_id"`
    ID          string             `bson:"id"`
    CustomerID  string             `bson:"customer_id"`
    Status      string             `bson:"status"`
    PaymentLink string             `bson:"payment_link"`
    Items       []*entity.Item     `bson:"items"`
}
```

这里值得你注意：
- Mongo 真正主键是 `_id`，类型是 `ObjectID`
- 同时作者还保留了一个业务层的 `ID` 字段

但后面 `Create` / `unmarshal` 里又会看到，这里的处理还带一点过渡态味道。

继续看 `Create`：

```go
func (r *OrderRepositoryMongo) Create(ctx context.Context, order *domain.Order) (created *domain.Order, err error) {
    defer r.logWithTag("create", err, order, created)

    write := r.marshalToModel(order)
    res, err := r.collection().InsertOne(ctx, write)
    if err != nil {
        return nil, err
    }

    created = order
    created.ID = res.InsertedID.(primitive.ObjectID).Hex()
    return created, nil
}
```

这里的逻辑是：
- 先把 domain.Order 转成 Mongo 模型
- `InsertOne`
- 然后把 Mongo 返回的 `_id` 取出来，转成 hex string 回填到业务 `ID`

也就是说，这一版课程决定：

`业务里的订单 ID 直接复用 Mongo 的 ObjectID 字符串形式。`

这不是唯一方案，但在课程阶段很省事。

再看 `Get`：

```go
func (r *OrderRepositoryMongo) Get(ctx context.Context, id, customerID string) (got *domain.Order, err error) {
    defer r.logWithTag("get", err, nil, got)

    read := &orderModel{}
    mongoID, _ := primitive.ObjectIDFromHex(id)
    cond := bson.M{"_id": mongoID}
    if err = r.collection().FindOne(ctx, cond).Decode(read); err != nil {
        return
    }
    if read == nil {
        return nil, domain.NotFoundError{OrderID: id}
    }
    got = r.unmarshal(read)
    return got, nil
}
```

这里你要看懂两个点：
- 查询时是拿传进来的字符串 `id` 转回 `ObjectID`
- 实际查询条件只按 `_id` 查，没有真正用到 `customerID`

这说明这一版仓储是“能跑起来”的第一版，不是已经把所有条件和索引设计得很严密。

再看 `Update`：

```go
func (r *OrderRepositoryMongo) Update(...)(err error) {
    ...
    session, err := r.db.StartSession()
    ...
    if err = session.StartTransaction(); err != nil {
        return err
    }
    defer func() {
        if err == nil {
            _ = session.CommitTransaction(ctx)
        } else {
            _ = session.AbortTransaction(ctx)
        }
    }()

    oldOrder, err := r.Get(ctx, order.ID, order.CustomerID)
    if err != nil {
        return
    }
    updated, err := updateFn(ctx, oldOrder)
    if err != nil {
        return
    }

    mongoID, _ := primitive.ObjectIDFromHex(oldOrder.ID)
    _, err = r.collection().UpdateOne(
        ctx,
        bson.M{"_id": mongoID, "customer_id": oldOrder.CustomerID},
        bson.M{"$set": bson.M{
            "status":       updated.Status,
            "payment_link": updated.PaymentLink,
        }},
    )
    ...
}
```

这段你至少要学到三件事：

##### 1. 它试图做事务保护

虽然这里的事务用法还比较课程风格，但作者已经开始意识到：
- 先查旧值
- 再执行更新函数
- 再写回去

这类流程如果要更严肃地处理一致性，是会想到事务的。

##### 2. `updateFn` 思路被保留了

这和前面内存仓储版本是同一个抽象：
- 仓储负责取旧值
- 领域更新逻辑由 `updateFn` 决定
- 最后仓储负责持久化

这说明课程在切 Mongo 时，并没有推翻之前的更新抽象。

##### 3. 当前只更新了部分字段

这里只更新：
- `status`
- `payment_link`

没有更新 `items`。

这说明这版仓储实现是按当前业务最需要的字段先补的，不是一次性把所有更新场景都做完。

### 6.3 `internal/order/service/application.go`

#### 6.3.1 这个文件自己的原始 diff

```diff
diff --git a/internal/order/service/application.go b/internal/order/service/application.go
index 8de0e53..59bbe5a 100644
--- a/internal/order/service/application.go
+++ b/internal/order/service/application.go
@@ -36,7 +41,9 @@ func NewApplication(ctx context.Context) (app.Application, func()) {
 }
 
 func newApplication(_ context.Context, stockGRPC query.StockService, ch *amqp.Channel) app.Application {
-	orderRepo := adapters.NewMemoryOrderRepository()
+	//orderRepo := adapters.NewMemoryOrderRepository()
+	mongoClient := newMongoClient()
+	orderRepo := adapters.NewOrderRepositoryMongo(mongoClient)
@@
+func newMongoClient() *mongo.Client {
+    ...
+}
```

#### 6.3.2 这段在做什么

它把 order 服务真正从：
- 内存仓储

切到了：
- Mongo 仓储

这是整个组里对运行行为影响最大的改动之一。

#### 6.3.3 关键代码和详细注释

```go
func newApplication(_ context.Context, stockGRPC query.StockService, ch *amqp.Channel) app.Application {
    // 旧版是内存仓储
    // orderRepo := adapters.NewMemoryOrderRepository()

    mongoClient := newMongoClient()
    orderRepo := adapters.NewOrderRepositoryMongo(mongoClient)

    logger := logrus.NewEntry(logrus.StandardLogger())
    metricClient := metrics.TodoMetrics{}
    return app.Application{
        Commands: app.Commands{
            CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, ch, logger, metricClient),
            UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
        },
        Queries: app.Queries{
            GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
        },
    }
}
```

这段要学会看出一个架构点：

`应用层其实并不在乎底下是内存仓储还是 Mongo 仓储，只要 Repository 接口不变，application 组装时换实现就行。`

这就是接口抽象真正带来的收益。

继续看 Mongo client：

```go
func newMongoClient() *mongo.Client {
    uri := fmt.Sprintf(
        "mongodb://%s:%s@%s:%s",
        viper.GetString("mongo.user"),
        viper.GetString("mongo.password"),
        viper.GetString("mongo.host"),
        viper.GetString("mongo.port"),
    )

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        panic(err)
    }
    if err = c.Ping(ctx, readpref.Primary()); err != nil {
        panic(err)
    }
    return c
}
```

这里又补了几个基础设施知识点：

##### `mongo.Connect(...)` 不等于“连接一定可用”

所以后面还要 `Ping(...)`。

##### `context.WithTimeout(...)` 的意义

表示：
- 连接 Mongo 最多等 10 秒
- 超过就失败
- 避免程序卡死在无限等待上

这是一种典型的启动期外部依赖连接保护。

## 7. 领域层顺手补了一个真正的“待支付订单构造器”

### 7.1 `internal/order/domain/order/order.go`

#### 7.1.1 这个文件自己的原始 diff

```diff
diff --git a/internal/order/domain/order/order.go b/internal/order/domain/order/order.go
index 08e7e96..1afa1cb 100644
--- a/internal/order/domain/order/order.go
+++ b/internal/order/domain/order/order.go
@@ -38,6 +38,20 @@ func NewOrder(id, customerID, status, paymentLink string, items []*entity.Item)
 	}, nil
 }
 
+func NewPendingOrder(customerId string, items []*entity.Item) (*Order, error) {
+	if customerId == "" {
+		return nil, errors.New("empty customerID")
+	}
+	if items == nil {
+		return nil, errors.New("empty items")
+	}
+	return &Order{
+		CustomerID: customerId,
+		Status:     "pending",
+		Items:      items,
+	}, nil
+	}
```

#### 7.1.2 为什么这段值得讲

前面 `CreateOrder` 里是直接手写：
- `&domain.Order{CustomerID:..., Items:...}`

现在改成 `domain.NewPendingOrder(...)`，这是个很好的趋势：

`把“订单刚创建时应该长什么样”收拢到领域层。`

也就是说，“新订单默认状态是 pending”不再散落在应用层里，而是由领域构造器统一规定。

### 7.2 `internal/order/app/command/create_order.go`

#### 7.2.1 这个文件自己的原始 diff

```diff
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
index 64f9070..6384a96 100644
--- a/internal/order/app/command/create_order.go
+++ b/internal/order/app/command/create_order.go
@@ -75,10 +75,11 @@ func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*Creat
 	if err != nil {
 		return nil, err
 	}
-	o, err := c.orderRepo.Create(ctx, &domain.Order{
-		CustomerID: cmd.CustomerID,
-		Items:      validItems,
-	})
+	pendingOrder, err := domain.NewPendingOrder(cmd.CustomerID, validItems)
+	if err != nil {
+		return nil, err
+	}
+	o, err := c.orderRepo.Create(ctx, pendingOrder)
```

#### 7.2.2 这段在做什么

它配合上面的 `NewPendingOrder(...)`，把订单创建流程从“裸 struct 初始化”改成“通过领域构造器拿合法对象”。

这在 DDD 或普通干净架构里都是很合理的：
- 应用层负责 orchestration
- 领域层负责对象合法性和默认业务状态

## 8. 第五条主线：各服务入口继续收口到统一配置初始化

### 8.1 `internal/order/main.go` / `payment/main.go` / `stock/main.go`

这三个文件的变化本身不大，但它们背后的意义是：
- 现在很多地方都只靠 `_ ".../common/config"` 做副作用导入
- 不再要求 main 手工写一遍配置初始化

这和 `viper.go` 这组的改造是联动的。

你可以理解成：

`配置初始化开始从“显式调用某个函数”转向“导入配置包即初始化”。`

这能让多个服务入口更整齐。

但你也要记住前面讲过的边界：
- 这种方式省事
- 但会增加隐式副作用
- 大项目里需要谨慎控制

## 9. 这组真正的运行行为变化是什么

如果你只抓最重要的行为变化，我会这样总结：

### 9.1 订单终于持久化到 Mongo

意味着：
- 订单不再只存在内存里
- 服务重启后数据还能留住
- 订单 ID 现在直接来源于 Mongo `_id`

### 9.2 MQ 失败处理终于开始真正 retry

意味着：
- `payment consumer` 和 `order consumer` 出错时不再只是 TODO
- 系统开始按 `x-retry-count` + `max-retry` + `DLQ` 这套逻辑处理失败消息

### 9.3 配置读取终于更稳一点了

意味着：
- broker、mongo repository、main 等地方开始更容易稳定地拿到配置

## 10. 这组的局限和你必须警惕的点

### 10.1 Mongo 仓储还不是成熟版

有几个明显点：

#### `Get(...)` 没真正用 `customerID`

签名里收了 `customerID`，但查询条件主要还是 `_id`。

#### `Update(...)` 只更新部分字段

只改了：
- `status`
- `payment_link`

没有全面更新 `items` 等其他字段。

#### 事务写法是课程式的

有事务意识是好事，但这版实现还比较粗。真实生产里还要更仔细地确认：
- session context 用法
- mongo deployment 是否支持事务
- 错误处理和重试边界

### 10.2 `HandleRetry` 接进 consumer 后，`Ack/Nack` 语义需要你仔细想

现在逻辑是：
- 如果业务报错，先尝试 `HandleRetry(...)`
- 之后 defer 看 `err != nil` 就 `Nack`

这说明：
- retry publish 成功与否
- 原消息是否还需要立刻 Nack
- 是否会导致重复投递语义

这些在真实系统里都要仔细推敲。

课程阶段先这么写，是为了把主流程串起来。

### 10.3 `viper.go` 里的 `fmt.Printf(...)` 和隐式 init 还是偏教学版

它们能帮助课程演示“路径到底算出来没有”，
但如果你以后自己写项目，最好还是把这种调试输出和初始化副作用控制得更谨慎。

### 10.4 `mongo-express` 裸开在本地 compose 里很方便，但不是生产做法

它是教学/调试工具，不是生产环境默认组件。

## 11. 这组最该记住的话

1. 这一组最大的行为变化有两个：retry 正式接线，order 持久化切到 Mongo。
2. `HandleRetry(...)` 只有接进 consumer，才真正开始改变系统行为。
3. `NewPendingOrder(...)` 说明领域层开始负责“新订单默认状态是什么”。
4. `Viper` 初始化从手工调用转向副作用导入 + once，是为了让多入口多模块更容易稳定拿配置。
5. Mongo 仓储这版是“第一版能跑通”，不是已经成熟到可以直接照搬生产。

## 12. 你现在应该怎么复习这组

建议顺序：

1. [internal/order/infrastructure/consumer/consumer.go](/g:/shi/go_shop_second/internal/order/infrastructure/consumer/consumer.go)
2. [internal/payment/infrastructure/consumer/consumer.go](/g:/shi/go_shop_second/internal/payment/infrastructure/consumer/consumer.go)
   - 先确认 retry 是怎么真正接进去的
3. [internal/common/broker/rabbitmq.go](/g:/shi/go_shop_second/internal/common/broker/rabbitmq.go)
   - 再回看 retry 底层策略怎么被配置化
4. [internal/common/config/viper.go](/g:/shi/go_shop_second/internal/common/config/viper.go)
   - 理解这组为什么必须修配置初始化
5. [internal/order/adapters/order_mongo_repository.go](/g:/shi/go_shop_second/internal/order/adapters/order_mongo_repository.go)
   - 理解 order 仓储第一次落库
6. [internal/order/service/application.go](/g:/shi/go_shop_second/internal/order/service/application.go)
   - 看接口抽象怎么让“内存仓储 -> Mongo 仓储”切换变得很自然

如果你继续，我下一组就写 `lesson35 -> lesson36`。
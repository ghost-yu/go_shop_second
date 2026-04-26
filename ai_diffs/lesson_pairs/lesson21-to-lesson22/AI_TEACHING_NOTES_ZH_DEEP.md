# `lesson21 -> lesson22` 独立讲义（详细注释版）

这一组是上一组 webhook 的直接后续。

上一组 `lesson20 -> lesson21` 做完之后，系统已经具备这个能力：

- 用户在 Stripe 支付成功
- Stripe 会回调 `payment` 服务的 `/api/webhook`
- `payment` 服务会把 `order.paid` 事件发到 RabbitMQ

但是当时整条链还没有真正闭环。

为什么？

因为：

`payment` 虽然把 “订单已支付” 这条消息发出来了，但 `order` 服务还没有真正去消费这条消息。`

消息已经在 MQ 里了，订单状态却还没有真的推进到 `paid`。

这就是这一组要补的核心：

`让 order 服务开始消费 order.paid 事件，把消息里的“已支付”状态真正落到订单自身。`

所以这组你要把它理解成：

`支付链第一次从“能发 paid 消息”推进到“收到 paid 消息后，订单自己真的进入 paid 状态”。`

这组对项目整体非常关键，因为这是第一次把“外部支付成功”真正变成“我们自己系统内部订单状态变更”的事实。

## 1. 原始 diff 去哪里看

原始差异在这里：
[diff.md](/g:/shi/go_shop_second/ai_diffs/lesson_pairs/lesson21-to-lesson22/diff.md)

这次我继续严格按固定要求写：
- 每个文件先贴自己的 diff
- 再贴带中文注释的代码/关键代码
- 讲清楚旧代码做什么、新代码做什么、为什么要这样改
- 并结合对应字幕 `23.txt`

## 2. 正确阅读顺序

这组建议按下面顺序读：

1. [internal/order/infrastructure/consumer/consumer.go](/g:/shi/go_shop_second/internal/order/infrastructure/consumer/consumer.go)
2. [internal/order/main.go](/g:/shi/go_shop_second/internal/order/main.go)
3. [internal/order/domain/order/order.go](/g:/shi/go_shop_second/internal/order/domain/order/order.go)

原因：
- `consumer.go` 是这组真正的业务核心，负责消费 `order.paid`
- `main.go` 解释 order 服务为什么现在也要连 MQ，并启动自己的 consumer
- `order.go` 解释订单领域对象为什么要补“已支付”校验/状态能力

## 3. 总调用链

这组最重要的脑图是：

```text
Stripe 支付成功
-> payment 服务收到 webhook
-> payment 服务发布 order.paid 事件到 RabbitMQ
-> order 服务的 consumer 订阅 order.paid
-> consumer 把消息反序列化成 Order
-> 调用 order 应用层的 UpdateOrder 命令
-> UpdateFn 检查/更新订单状态
-> order 仓储把这次状态变更写回去
-> 前端 success 页面下一次轮询时，就能看到 status=paid
```

你一定要看懂这件事：

`这组不是“又加了个 MQ consumer”，而是“让 order 服务第一次真正接手 payment 发回来的支付结果”。`

也就是说，系统开始出现一个很完整的事件闭环：
- payment 负责感知支付成功
- order 负责把自己的订单状态更新掉

这就符合微服务里“谁拥有数据，谁负责改自己的数据”这个基本原则。

## 4. 关键文件详细讲解

### 4.1 [internal/order/infrastructure/consumer/consumer.go](/g:/shi/go_shop_second/internal/order/infrastructure/consumer/consumer.go)

这份文件是这组最核心的新文件。

在旧版本里，`order` 服务并不会消费 `order.paid`。也就是说：
- payment 发出消息
- order 并没有接手
- 消息层面是成功的，但订单状态层面还没有真正闭环

这组新增的 `consumer.go`，就是让 `order` 服务正式进入“消息消费者”的角色。

先看这个文件自己的原始 diff。

```diff
diff --git a/internal/order/infrastructure/consumer/consumer.go b/internal/order/infrastructure/consumer/consumer.go
new file mode 100644
index 0000000..7e2b8a6
--- /dev/null
+++ b/internal/order/infrastructure/consumer/consumer.go
@@ -0,0 +1,70 @@
+package consumer
+
+import (
+	"context"
+	"encoding/json"
+
+	"github.com/ghost-yu/go_shop_second/common/broker"
+	"github.com/ghost-yu/go_shop_second/order/app"
+	"github.com/ghost-yu/go_shop_second/order/app/command"
+	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
+	amqp "github.com/rabbitmq/amqp091-go"
+	"github.com/sirupsen/logrus"
+)
+
+type Consumer struct {
+	app app.Application
+}
+
+func NewConsumer(app app.Application) *Consumer {
+	return &Consumer{
+		app: app,
+	}
+}
+
+func (c *Consumer) Listen(ch *amqp.Channel) {
+	q, err := ch.QueueDeclare(broker.EventOrderPaid, true, false, true, false, nil)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	err = ch.QueueBind(q.Name, "", broker.EventOrderPaid, false, nil)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	var forever chan struct{}
+	go func() {
+		for msg := range msgs {
+			c.handleMessage(msg, q, ch)
+		}
+	}()
+	<-forever
+}
+
+func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue, ch *amqp.Channel) {
+	o := &domain.Order{}
+	if err := json.Unmarshal(msg.Body, o); err != nil {
+		logrus.Infof("error unmarshal msg.body into domain.order, err = %v", err)
+		_ = msg.Nack(false, false)
+		return
+	}
+	_, err := c.app.Commands.UpdateOrder.Handle(context.Background(), command.UpdateOrder{
+		Order: o,
+		UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
+			if err := order.IsPaid(); err != nil {
+				return nil, err
+			}
+			return order, nil
+		},
+	})
+	if err != nil {
+		logrus.Infof("error updating order, orderID = %s, err = %v", o.ID, err)
+		// TODO: retry
+		return
+	}
+	_ = msg.Ack(false)
+	logrus.Info("order consume paid event success!")
+}
```

这个 diff 很长，但本质上只做了三件事：

1. 声明并绑定一个消费 `order.paid` 的队列
2. 从 RabbitMQ 里拿到消息
3. 把消息转成订单对象，再交给 `UpdateOrder` 命令更新订单状态

下面我贴当前代码里这个文件的关键部分，并写详细中文注释。

```go
package consumer

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/ghost-yu/go_shop_second/common/broker"
    "github.com/ghost-yu/go_shop_second/common/logging"
    "github.com/ghost-yu/go_shop_second/order/app"
    "github.com/ghost-yu/go_shop_second/order/app/command"
    domain "github.com/ghost-yu/go_shop_second/order/domain/order"
    "github.com/pkg/errors"
    amqp "github.com/rabbitmq/amqp091-go"
    "github.com/sirupsen/logrus"
    "go.opentelemetry.io/otel"
)

type Consumer struct {
    // 持有应用层入口。consumer 自己不直接改仓储，
    // 它还是通过应用层命令来更新订单。
    app app.Application
}

func NewConsumer(app app.Application) *Consumer {
    return &Consumer{
        app: app,
    }
}

func (c *Consumer) Listen(ch *amqp.Channel) {
    // 声明一个队列，用来接收 order.paid 交换机发出来的消息。
    q, err := ch.QueueDeclare(broker.EventOrderPaid, true, false, true, false, nil)
    if err != nil {
        logrus.Fatal(err)
    }

    // 把这个队列绑定到 order.paid 这个 exchange 上。
    // 这样 payment 服务发布到 order.paid 的事件，order 才能收到。
    err = ch.QueueBind(q.Name, "", broker.EventOrderPaid, false, nil)
    if err != nil {
        logrus.Fatal(err)
    }

    // 开始消费消息。autoAck=false，表示不是一收到消息就自动确认，
    // 而是等真正处理成功后，再手动 Ack。
    msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
    if err != nil {
        logrus.Fatal(err)
    }

    var forever chan struct{}
    go func() {
        for msg := range msgs {
            c.handleMessage(ch, msg, q)
        }
    }()

    // 用一个永远不写入的 channel 把 goroutine 挂住，让 consumer 持续运行。
    <-forever
}
```

#### 这段代码到底在做什么

简单说：

`order 服务现在开始盯着 order.paid 这个消息源，一旦 payment 发来“这张订单已支付”，它就立刻接手处理。`

#### 这里补一个 RabbitMQ 基础知识：`QueueDeclare`、`QueueBind`、`Consume`

你可以把这三个动作理解成：

1. `QueueDeclare`
   - 先准备一个收信箱

2. `QueueBind`
   - 告诉 RabbitMQ：某个 exchange 发出来的消息，应该投递到这个收信箱

3. `Consume`
   - 真正开始从这个收信箱里一条条拿消息

这三个动作缺一不可。

尤其是 `QueueBind`，新手很容易漏掉。

如果你只声明队列，不绑定 exchange，那么消息可能根本不会路由到这里来。

#### `autoAck=false` 为什么重要

这里 `Consume(..., false, ...)` 的第一个 `false` 表示：

`不要一收到消息就自动确认。`

为什么要这样？

因为这条消息不是“收到了就算处理成功”，而是：
- 要先反序列化
- 要再调用应用层更新订单
- 真成功后，才能 Ack

如果你一上来自动 Ack：
- 后面处理报错了
- RabbitMQ 也会以为这条消息已经成功消费完了
- 那你就丢消息了

所以这条设计非常关键。

---

继续看 `handleMessage`：

```go
func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
    tr := otel.Tracer("rabbitmq")
    ctx, span := tr.Start(
        broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers),
        fmt.Sprintf("rabbitmq.%s.consume", q.Name),
    )
    defer span.End()

    var err error
    defer func() {
        if err != nil {
            logging.Warnf(ctx, nil, "consume failed||from=%s||msg=%+v||err=%v", q.Name, msg, err)
            _ = msg.Nack(false, false)
        } else {
            logging.Infof(ctx, nil, "%s", "consume success")
            _ = msg.Ack(false)
        }
    }()

    o := &domain.Order{}
    if err = json.Unmarshal(msg.Body, o); err != nil {
        err = errors.Wrap(err, "error unmarshal msg.body into domain.order")
        return
    }

    _, err = c.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
        Order: o,
        UpdateFn: func(ctx context.Context, oldOrder *domain.Order) (*domain.Order, error) {
            if err := oldOrder.UpdateStatus(o.Status); err != nil {
                return nil, err
            }
            return oldOrder, nil
        },
    })
    if err != nil {
        logging.Errorf(ctx, nil, "error updating order||orderID=%s||err=%v", o.ID, err)
        if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
            err = errors.Wrapf(err, "retry_error, error handling retry, messageID=%s||err=%v", msg.MessageId, err)
        }
        return
    }

    span.AddEvent("order.updated")
}
```

这段代码非常值得你细读，因为它把“MQ 消息”真正翻译成了“订单状态更新命令”。

#### 第一步：从消息里还原出订单对象

```go
o := &domain.Order{}
if err = json.Unmarshal(msg.Body, o); err != nil {
    ...
}
```

这里的意思是：

`payment 服务发出来的消息体，本质上是一段 JSON；order consumer 先把它反序列化成自己的领域订单对象。`

你要注意“自己的领域订单对象”这几个字。

不是说：
- MQ 里有什么结构
- 就直接原样往数据库里塞

而是：
- 先解析成 `domain.Order`
- 再通过应用层命令去更新

这说明这组虽然用了 MQ，但仍然尽量保持应用层和领域层的边界。

#### 第二步：不是自己改仓储，而是调 `UpdateOrder` 命令

```go
_, err = c.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
    Order: o,
    UpdateFn: func(ctx context.Context, oldOrder *domain.Order) (*domain.Order, error) {
        if err := oldOrder.UpdateStatus(o.Status); err != nil {
            return nil, err
        }
        return oldOrder, nil
    },
})
```

这是这组设计里非常值得你记住的一点：

`consumer 没有越过应用层去直接操作仓储。`

它做的是：
- 拿到消息里的新状态 `o.Status`
- 把它交给 `UpdateOrder` 命令
- 让命令内部通过 `UpdateFn` 去更新旧订单

为什么这样设计更好？

因为应用层命令是整个系统“允许怎么改订单”的统一入口。

如果 consumer 自己直接改仓储，后面会有两个问题：
- 业务规则分散
- 很多状态校验逻辑绕开了命令层

所以这里虽然看起来多绕了一层，但其实是在守住应用层边界。

#### 第三步：`UpdateFn` 为什么拿的是 `oldOrder`

这里很关键：

```go
UpdateFn: func(ctx context.Context, oldOrder *domain.Order) (*domain.Order, error) {
    if err := oldOrder.UpdateStatus(o.Status); err != nil {
        return nil, err
    }
    return oldOrder, nil
}
```

这个函数的语义是：

- `o` 是消息里带来的“新状态提示”
- `oldOrder` 是仓储里原本已经存在的订单
- 真正被更新的是 `oldOrder`

也就是说，这里不是拿 MQ 里的对象直接覆盖数据库，而是：

`以数据库里的旧订单为准，在它身上尝试做一次合法状态迁移。`

这很重要，因为状态流转通常不能瞎跳。

比如：
- 不能从 `pending` 直接跳到 `ready`
- 不能乱改成一个完全未知状态

所以“先取旧订单，再在旧订单上更新”是更合理的。

#### 第四步：成功才 Ack，失败就 Nack / Retry

这组当前代码已经比当时 diff 更进一步：
- 成功时 `Ack`
- 失败时 `Nack`
- 更新失败时还会走 `broker.HandleRetry(...)`

这说明这一条链已经开始考虑“消息不要轻易丢”的问题了。

你可以把它理解成：

`MQ 消费不是拿到消息就结束，而是要为失败负责。`

### 这里要补的一个误区：字幕里的 `IsPaid()` 和当前代码的 `UpdateStatus(...)`

你会发现 diff 里当时的实现是：
- `order.IsPaid()`
- 只是校验“消息里的订单状态是不是 paid”

而你当前工作区里的代码已经进一步演进成：
- `oldOrder.UpdateStatus(o.Status)`
- 真正按领域规则去推进状态

这两者有本质差别：

旧实现更像：
- 验证“消息是不是 paid”

新实现更像：
- 验证“当前订单能不能从旧状态合法地推进到这个新状态”

后者显然更成熟。

所以你读这一组时要记住：
- diff 代表课程当时的落地形态
- 当前工作区代表这条链后来又被继续优化过

### 这一文件你最该带走的结论

`consumer.go` 这组真正完成的是：让 order 服务第一次具备“从 MQ 收到 paid 事件后，正式把支付结果落回自己订单状态”的能力。`

---

### 4.2 [internal/order/main.go](/g:/shi/go_shop_second/internal/order/main.go)

这个文件的变化不复杂，但它解释了 order 服务为什么现在也要接 RabbitMQ。

先看这个文件自己的 diff。

```diff
diff --git a/internal/order/main.go b/internal/order/main.go
index b35cdc8..4453aaf 100644
--- a/internal/order/main.go
+++ b/internal/order/main.go
@@ -3,11 +3,13 @@ package main
 import (
 	"context"
 
+	"github.com/ghost-yu/go_shop_second/common/broker"
 	"github.com/ghost-yu/go_shop_second/common/config"
 	"github.com/ghost-yu/go_shop_second/common/discovery"
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/ghost-yu/go_shop_second/common/server"
+	"github.com/ghost-yu/go_shop_second/order/infrastructure/consumer"
@@ -40,6 +42,18 @@ func main() {
 		_ = deregisterFunc()
 	}()
 
+	ch, closeCh := broker.Connect(
+		viper.GetString("rabbitmq.user"),
+		viper.GetString("rabbitmq.password"),
+		viper.GetString("rabbitmq.host"),
+		viper.GetString("rabbitmq.port"),
+	)
+	defer func() {
+		_ = ch.Close()
+		_ = closeCh()
+	}()
+	go consumer.NewConsumer(application).Listen(ch)
+
 	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
```

再看当前代码里的关键部分：

```go
func main() {
    serviceName := viper.GetString("order.service-name")

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    application, cleanup := service.NewApplication(ctx)
    defer cleanup()

    deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
    if err != nil {
        logrus.Fatal(err)
    }
    defer func() {
        _ = deregisterFunc()
    }()

    // 这次新加：order 服务也要连 RabbitMQ。
    ch, closeCh := broker.Connect(
        viper.GetString("rabbitmq.user"),
        viper.GetString("rabbitmq.password"),
        viper.GetString("rabbitmq.host"),
        viper.GetString("rabbitmq.port"),
    )
    defer func() {
        _ = ch.Close()
        _ = closeCh()
    }()

    // 启动 order 自己的消息消费者，去监听 order.paid。
    go consumer.NewConsumer(application).Listen(ch)

    go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
        ...
    })

    server.RunHTTPServer(serviceName, func(router *gin.Engine) {
        ...
    })
}
```

旧代码做的事情是：
- order 服务只提供 HTTP / gRPC
- 它自己不消费 MQ

新代码做的事情是：
- order 服务现在也连接 RabbitMQ
- 并且自己启动 consumer 去监听 `order.paid`

这代表什么？

代表 `order` 这个服务的职责开始扩展了：
- 它不再只是“被动等 API 调用”
- 它开始主动监听系统里的异步事件

这也是事件驱动架构的典型样子。

你以后看这种 `main.go` 初始化代码，要有个习惯：

`谁在 main 里多了一条连接、多起了一个 goroutine，通常就说明这个服务的新职责增加了。`

### 这一文件你最该带走的结论

`main.go` 这组真正表达的是：order 服务现在不只是 HTTP/gRPC 提供者，它也成了 MQ 里的正式事件消费者。`

---

### 4.3 [internal/order/domain/order/order.go](/g:/shi/go_shop_second/internal/order/domain/order/order.go)

这份文件在 diff 里改动很小，但它代表了“订单领域如何看待 paid 状态”这件事。

先看这个文件自己的 diff。

```diff
diff --git a/internal/order/domain/order/order.go b/internal/order/domain/order/order.go
index b3a9fa9..d87e406 100644
--- a/internal/order/domain/order/order.go
+++ b/internal/order/domain/order/order.go
@@ -1,8 +1,11 @@
 package order
 
 import (
+	"fmt"
+
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	"github.com/pkg/errors"
+	"github.com/stripe/stripe-go/v80"
 )
@@ -44,3 +47,10 @@ func (o *Order) ToProto() *orderpb.Order {
 		PaymentLink: o.PaymentLink,
 	}
 }
+
+func (o *Order) IsPaid() error {
+	if o.Status == string(stripe.CheckoutSessionPaymentStatusPaid) {
+		return nil
+	}
+	return fmt.Errorf("order status not paid, order id = %s, status = %s", o.ID, o.Status)
+}
```

这段 diff 当时做的事情很直接：

`给订单领域对象补一个 IsPaid()，专门用来判断 MQ 消息里携带的订单状态是不是 paid。`

对应到当时 consumer 里的写法，就是：
- consumer 收到消息
- 调 `order.IsPaid()`
- 如果不是 paid，就报错，不继续更新

再看当前工作区里，这个文件已经继续演进成更成熟的状态迁移模型：

```go
func (o *Order) UpdateStatus(to string) error {
    if !o.isValidStatusTransition(to) {
        return fmt.Errorf("cannot transit from '%s' to '%s'", o.Status, to)
    }
    o.Status = to
    return nil
}

func (o *Order) isValidStatusTransition(to string) bool {
    switch o.Status {
    default:
        return false
    case consts.OrderStatusPending:
        return slices.Contains([]string{consts.OrderStatusWaitingForPayment}, to)
    case consts.OrderStatusWaitingForPayment:
        return slices.Contains([]string{consts.OrderStatusPaid}, to)
    case consts.OrderStatusPaid:
        return slices.Contains([]string{consts.OrderStatusReady}, to)
    }
}
```

也就是说，这条链后来从：
- 只会判断“是不是 paid”

进一步升级到了：
- 按完整订单状态机去校验“能不能从旧状态推进到新状态”

这说明什么？

说明这一组当时解决的是第一个关键问题：

`让 order 服务具备对 paid 状态做领域判断的意识。`

虽然当时判断还比较简单，但方向是对的：
- 订单状态不是随便一个字符串
- 它应该受领域规则约束

#### 为什么当时要 import Stripe 的状态常量

你会看到 diff 里用了：

```go
stripe.CheckoutSessionPaymentStatusPaid
```

作者当时这么写，是因为 payment webhook 里发过来的状态值，直接沿用了 Stripe 的 `paid` 字面量。

这样做的好处是：
- 不用自己再映射一次
- 当时上手快

坏处是：
- order 域直接依赖了 Stripe 的概念
- 域模型被外部支付平台语义“污染”了一点

所以你看当前工作区后面已经在往 `common/consts.OrderStatusPaid` 这种内部统一状态收敛了。

这就是教学项目很常见的演进路径：
- 先借用外部平台状态跑通链路
- 后面再把内部状态收口成自己的常量

### 这一文件你最该带走的结论

`order.go` 这组变化虽然小，但它代表订单领域对象第一次开始正式承认“paid”是一种需要被校验的业务状态，而不是随便来一条字符串就能改。`

## 5. 第三方库和易错点

### RabbitMQ / amqp

这组里你要特别记住：
- `QueueDeclare` 是准备队列
- `QueueBind` 是把队列挂到 exchange 上
- `Consume` 是开始拿消息
- `Ack` 是处理成功后确认消费
- `Nack` 是处理失败后拒绝消费

最容易忘的坑：
- 忘了 `Bind`，消息根本到不了队列
- 自动 Ack 开太早，导致失败消息丢失
- consumer 启动了，但没有一直阻塞住，进程很快退出

### Stripe 状态常量

这组 diff 里 order 域直接用了 Stripe 的 `PaymentStatusPaid`。

你要知道：
- 这样写短期方便
- 但长期会让领域层依赖外部支付平台语义

所以这只是教学过程中的过渡态，不是最优状态。

### 应用层命令

这组一个非常好的点是：
- consumer 没直接写仓储
- 还是走 `UpdateOrder.Handle(...)`

这说明：

`消息消费不是业务逻辑例外，它仍然要走应用层入口。`

这点你以后做项目一定要记住，不然系统很容易出现“HTTP 一套规则、MQ 一套规则”的分裂。

## 6. 为什么这么设计

这一组设计选择的是：

`payment 负责感知支付成功，order 负责消费 paid 事件并更新自己的订单状态。`

这非常符合“谁拥有数据，谁维护数据”的原则。

为什么不让 payment 直接改 order 仓储？

因为那样会有几个问题：
- payment 越权
- 服务边界变模糊
- order 自己的状态规则被绕开

而当前设计是：
- payment 只发事实：这张订单 paid 了
- order 自己决定：这个事实如何落进自己的领域模型和仓储

这是一种更干净的职责划分。

## 7. 当前还不完美的地方

这组虽然关键，但也明显还在过渡态：

1. diff 里的 `IsPaid()` 规则太单薄，只会看“是不是 paid”
2. 当时版本里 retry 只是 `TODO`，失败处理不完整
3. order 域对 Stripe 状态常量存在耦合
4. MQ 消息结构和领域对象结构仍然绑得比较近
5. 当前系统虽然已经能推进到 `paid`，但 `paid -> ready` 还得后续服务继续接力

也就是说，这组完成的是“支付闭环的中段”，不是最终整条业务终态。

## 8. 这组最该带走的知识点

1. webhook 发出 `order.paid` 只是第一步，真正闭环要看 order 有没有消费并落库
2. MQ consumer 最合理的做法，不是直接改仓储，而是继续走应用层命令
3. 成功后 Ack、失败时 Nack/Retry，是消息系统里非常重要的基本纪律
4. 订单状态不应该随便改，至少要经过领域规则验证
5. 微服务里“谁拥有订单，谁维护订单状态”是很重要的边界原则
6. 这一组让前端页面轮询看到 `paid` 成为可能，因为 order 终于真正接住了 payment 发回来的支付成功消息

## 9. 一句话收住这组

`lesson21 -> lesson22` 的本质，不是“order 也接了 RabbitMQ”，而是“支付成功这件外部事实，第一次真正被 order 服务接住，并转化成了订单自身的 paid 状态”。`
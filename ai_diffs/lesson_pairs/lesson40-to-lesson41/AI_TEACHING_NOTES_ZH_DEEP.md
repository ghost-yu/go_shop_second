# `lesson40 -> lesson41` 独立讲义（详细注释版）

这一组的核心不是修一个小 bug，而是：

`系统里正式多出一个新的微服务：kitchen（厨房服务）。`

如果你把前面课程理解成：

- `order` 负责创建订单
- `payment` 负责收款
- `stock` 负责库存和商品信息

那么这一组就是把“支付完成以后，订单进入制作阶段”这条链补出来。

你可以把它理解成：

1. 用户支付成功
2. `payment` 发布 `order.paid` 事件
3. 新增的 `kitchen` 服务消费这条消息
4. `kitchen` 模拟做餐
5. 做完以后再回调 `order` 服务，把订单状态改成 `ready`

这一组还有一条很重要的副线：

`stock` 不再只是“随便回点假数据”，它开始真正把商品 ID 转成 Stripe 里的价格 ID，并且把 proto 类型和内部 entity 类型分开。`

## 1. 正确阅读顺序

建议按这个顺序看：

1. [internal/kitchen/main.go](/g:/shi/go_shop_second/internal/kitchen/main.go)
2. [internal/kitchen/infrastructure/consumer/consumer.go](/g:/shi/go_shop_second/internal/kitchen/infrastructure/consumer/consumer.go)
3. [internal/kitchen/adapters/order_grpc_client.go](/g:/shi/go_shop_second/internal/kitchen/adapters/order_grpc_client.go)
4. [internal/stock/app/query/check_if_items_in_stock.go](/g:/shi/go_shop_second/internal/stock/app/query/check_if_items_in_stock.go)
5. [internal/stock/infrastructure/integration/stripe.go](/g:/shi/go_shop_second/internal/stock/infrastructure/integration/stripe.go)
6. [internal/stock/ports/grpc.go](/g:/shi/go_shop_second/internal/stock/ports/grpc.go)
7. [internal/stock/convertor/convertor.go](/g:/shi/go_shop_second/internal/stock/convertor/convertor.go)
8. [internal/stock/convertor/facade.go](/g:/shi/go_shop_second/internal/stock/convertor/facade.go)
9. [internal/stock/entity/entity.go](/g:/shi/go_shop_second/internal/stock/entity/entity.go)
10. [internal/stock/service/application.go](/g:/shi/go_shop_second/internal/stock/service/application.go)
11. [internal/order/adapters/grpc/stock_grpc.go](/g:/shi/go_shop_second/internal/order/adapters/grpc/stock_grpc.go)
12. [internal/common/decorator/logging.go](/g:/shi/go_shop_second/internal/common/decorator/logging.go)
13. [internal/order/tests/create_order_test.go](/g:/shi/go_shop_second/internal/order/tests/create_order_test.go)

## 2. 总调用链

这一组之后，支付完成到出餐完成的大致链路变成：

`payment 收到支付成功 -> payment 发布 order.paid 到 RabbitMQ -> kitchen consumer 订阅到消息 -> kitchen 模拟做餐 -> kitchen 通过 gRPC 调用 order.UpdateOrder -> order 状态被改成 ready`

与此同时，创建订单时的库存 / 商品检查链也变得更严谨：

`order 调 stock.CheckIfItemsInStock -> stock gRPC 层把 proto 请求转成 entity -> query handler 调 Stripe API 根据商品 ID 取 price ID -> 再把 entity 转回 proto 返回给 order`

## 3. 这组到底解决什么问题

### 3.1 新问题：系统里缺“厨房阶段”

前面课程里，订单链路大致是：

- 创建订单
- 支付
- 回写已支付

但还缺少“支付之后进入制作”这一步。现实里支付成功不等于订单完成。

### 3.2 老问题：stock 层边界仍然不干净

之前 `stock` 某些地方还直接拿 proto 类型做业务数据结构，而且 `priceID` 依赖 stub map。这样写在课程初期没问题，但随着服务越来越多，会带来两个问题：

1. 业务层和传输层耦合得太紧
2. Stripe 价格信息不是实时拿的，而是写死的假数据

## 4. 关键文件一：`internal/kitchen/main.go`

### 4.1 这个文件自己的原始 diff

```diff
diff --git a/internal/kitchen/main.go b/internal/kitchen/main.go
new file mode 100644
index 0000000..35b4159
--- /dev/null
+++ b/internal/kitchen/main.go
@@ -0,0 +1,66 @@
+package main
+
+import (
+	"context"
+	"os"
+	"os/signal"
+	"syscall"
+
+	"github.com/ghost-yu/go_shop_second/common/broker"
+	grpcClient "github.com/ghost-yu/go_shop_second/common/client"
+	_ "github.com/ghost-yu/go_shop_second/common/config"
+	"github.com/ghost-yu/go_shop_second/common/logging"
+	"github.com/ghost-yu/go_shop_second/common/tracing"
+	"github.com/ghost-yu/go_shop_second/kitchen/adapters"
+	"github.com/ghost-yu/go_shop_second/kitchen/infrastructure/consumer"
+	"github.com/sirupsen/logrus"
+	"github.com/spf13/viper"
+)
+
+func init() {
+	logging.Init()
+}
+
+func main() {
+	serviceName := viper.GetString("kitchen.service-name")
+
+	ctx, cancel := context.WithCancel(context.Background())
+	defer cancel()
+
+	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	defer shutdown(ctx)
+
+	orderClient, closeFunc, err := grpcClient.NewOrderGRPCClient(ctx)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	defer closeFunc()
+
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
+
+	orderGRPC := adapters.NewOrderGRPC(orderClient)
+	go consumer.NewConsumer(orderGRPC).Listen(ch)
+
+	sigs := make(chan os.Signal, 1)
+	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
+
+	go func() {
+		<-sigs
+		logrus.Infof("receive signal, exiting...")
+		os.Exit(0)
+	}()
+	logrus.Println("to exit, press ctrl+c")
+	select {}
+}
+```

### 4.2 这段代码在项目里做什么

这就是 `kitchen` 服务的启动入口。它负责：

- 初始化日志
- 初始化 tracing
- 连上 order 服务的 gRPC client
- 连上 RabbitMQ
- 启动 kitchen 的消息消费者
- 阻塞主 goroutine，让进程一直活着

### 4.3 带中文注释的关键代码

```go
func main() {
    // viper.GetString(key) 从配置里读取字符串值。
    // 如果 key 没配到，通常返回空串。
    serviceName := viper.GetString("kitchen.service-name")

    // context.Background() 是根 context。
    // context.WithCancel 会返回一个可取消的子 context 和 cancel 函数。
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // 初始化 Jaeger tracing。shutdown 是退出时的收尾函数。
    shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
    if err != nil {
        // logrus.Fatal = 打日志 + os.Exit(1)
        logrus.Fatal(err)
    }
    defer shutdown(ctx)

    // 建立到 order 服务的 gRPC client。
    orderClient, closeFunc, err := grpcClient.NewOrderGRPCClient(ctx)
    if err != nil {
        logrus.Fatal(err)
    }
    defer closeFunc()

    // 连接 RabbitMQ。返回 channel 和关闭函数。
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

    orderGRPC := adapters.NewOrderGRPC(orderClient)

    // go 关键字会启动新的 goroutine。
    go consumer.NewConsumer(orderGRPC).Listen(ch)

    // make(chan os.Signal, 1) 创建一个缓冲为 1 的 channel。
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

    go func() {
        <-sigs
        logrus.Infof("receive signal, exiting...")
        os.Exit(0)
    }()

    logrus.Println("to exit, press ctrl+c")

    // 空 select 会永久阻塞当前 goroutine。
    select {}
}
```

### 4.4 这里最需要你记住的库函数

#### `context.WithCancel()`

来源：标准库 `context`

作用：

- 基于父 context 生成一个可取消的子 context
- 返回 `(ctx, cancel)`

#### `signal.Notify()`

来源：标准库 `os/signal`

作用：

- 把系统信号转发到 channel

#### `select {}`

这是 Go 语法，不是库函数。效果是永久阻塞。

## 5. 关键文件二：`internal/kitchen/infrastructure/consumer/consumer.go`

### 5.1 这个文件自己的原始 diff

```diff
diff --git a/internal/kitchen/infrastructure/consumer/consumer.go b/internal/kitchen/infrastructure/consumer/consumer.go
new file mode 100644
index 0000000..f17b444
--- /dev/null
+++ b/internal/kitchen/infrastructure/consumer/consumer.go
@@ -0,0 +1,107 @@
+package consumer
+...
+func (c *Consumer) Listen(ch *amqp.Channel) {
+	q, err := ch.QueueDeclare("", true, false, true, false, nil)
+	...
+	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
+	...
+}
+...
+func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
+	...
+	o := &Order{}
+	if err := json.Unmarshal(msg.Body, o); err != nil {
+		...
+	}
+	if o.Status != "paid" {
+		err = errors.New("order not paid, cannot cook")
+		return
+	}
+	cook(o)
+	...
+	if err := c.orderGRPC.UpdateOrder(mqCtx, &orderpb.Order{...}); err != nil {
+		if err = broker.HandleRetry(mqCtx, ch, &msg); err != nil {
+			...
+		}
+		return
+	}
+}
+```

### 5.2 这段代码在项目里做什么

这是 kitchen 的核心业务文件。它负责：

- 创建 RabbitMQ 队列并绑定 `order.paid`
- 持续消费消息
- 反序列化成订单对象
- 校验状态是不是 `paid`
- 模拟做餐
- 调 order gRPC 把状态改成 `ready`
- 成功就 `Ack`，失败就 `Nack` 或走重试

### 5.3 带中文注释的关键代码

```go
type OrderService interface {
    // consumer 依赖接口，不依赖具体 gRPC 实现。
    UpdateOrder(ctx context.Context, request *orderpb.Order) error
}

func (c *Consumer) Listen(ch *amqp.Channel) {
    // QueueDeclare 用于声明队列。
    // 第一个参数传空串，表示让 RabbitMQ 自动生成队列名。
    q, err := ch.QueueDeclare("", true, false, true, false, nil)
    if err != nil {
        logrus.Fatal(err)
    }

    // QueueBind 把这个队列绑定到 order.paid 事件源。
    if err = ch.QueueBind(q.Name, "", broker.EventOrderPaid, false, nil); err != nil {
        logrus.Fatal(err)
    }

    // Consume 返回一个 Go channel，后续消息会从这里不断流出来。
    msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
    if err != nil {
        logrus.Warnf("fail to consume: queue=%s, err=%v", q.Name, err)
    }

    var forever chan struct{}
    go func() {
        for msg := range msgs {
            c.handleMessage(ch, msg, q)
        }
    }()

    // nil channel 上接收会永久阻塞。
    <-forever
}

func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
    var err error

    logrus.Infof("kitchen receive a message from %s, msg=%v", q.Name, string(msg.Body))

    // 从 RabbitMQ header 里提取 trace 上下文，避免异步链路 trace 断掉。
    ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)

    // otel.Tracer(name) 可以理解成 span 工厂。
    tr := otel.Tracer("rabbitmq")
    mqCtx, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))

    defer func() {
        span.End()

        // Ack = 成功消费
        // Nack = 失败消费
        if err != nil {
            _ = msg.Nack(false, false)
        } else {
            _ = msg.Ack(false)
        }
    }()

    o := &Order{}
    if err := json.Unmarshal(msg.Body, o); err != nil {
        logrus.Infof("failed to unmarshall msg to order, err=%v", err)
        return
    }

    if o.Status != "paid" {
        err = errors.New("order not paid, cannot cook")
        return
    }

    cook(o)
    span.AddEvent(fmt.Sprintf("order_cook: %v", o))

    // 通过 gRPC 回调 order 服务，把订单推进到 ready。
    if err := c.orderGRPC.UpdateOrder(mqCtx, &orderpb.Order{
        ID:          o.ID,
        CustomerID:  o.CustomerID,
        Status:      "ready",
        Items:       o.Items,
        PaymentLink: o.PaymentLink,
    }); err != nil {
        if err = broker.HandleRetry(mqCtx, ch, &msg); err != nil {
            logrus.Warnf("kitchen: error handling retry: err=%v", err)
        }
        return
    }

    span.AddEvent("kitchen.order.finished.updated")
    logrus.Info("consume success")
}

func cook(o *Order) {
    logrus.Printf("cooking order: %s", o.ID)

    // time.Sleep 只是模拟做餐耗时。
    time.Sleep(5 * time.Second)

    logrus.Printf("order %s done!", o.ID)
}
```

### 5.4 这里最需要拆开的库函数 / 组件概念

#### `QueueDeclare()`

来源：RabbitMQ Go client `amqp091-go`

作用：声明一个队列。

注意：声明队列不等于收到消息，你还需要绑定 exchange，再开始 consume。

#### `QueueBind()`

作用：把队列绑定到某个 exchange / 路由规则。

#### `Consume()`

作用：开启消费，返回一个 Go channel。

#### `Ack()` / `Nack()`

作用：

- `Ack`：确认消费成功，broker 删除消息
- `Nack`：确认消费失败

#### `json.Unmarshal()`

来源：标准库 `encoding/json`

作用：把 JSON 字节反序列化为 Go 结构体。

坑：目标对象必须是指针。

#### `otel.Tracer().Start()`

来源：OpenTelemetry

作用：创建 span，并返回带 span 的新 context。

### 5.5 这份实现还不完美的地方

1. `Order` 结构直接复用了 `orderpb.Item`，边界还不够干净。
2. `cook()` 只是 `Sleep`，纯模拟。
3. `Listen()` 用 nil channel 永久阻塞，生命周期管理粗糙。
4. kitchen 还没有自己的持久化。

## 6. `internal/kitchen/adapters/order_grpc_client.go`

### 6.1 这个文件自己的原始 diff

```diff
diff --git a/internal/kitchen/adapters/order_grpc_client.go b/internal/kitchen/adapters/order_grpc_client.go
new file mode 100644
index 0000000..92402ab
--- /dev/null
+++ b/internal/kitchen/adapters/order_grpc_client.go
@@ -0,0 +1,20 @@
+package adapters
+...
+type OrderGRPC struct {
+	client orderpb.OrderServiceClient
+}
+...
+func (g *OrderGRPC) UpdateOrder(ctx context.Context, request *orderpb.Order) error {
+	_, err := g.client.UpdateOrder(ctx, request)
+	return err
+}
+```

### 6.2 带中文注释的代码

```go
type OrderGRPC struct {
    // 这是 protobuf 生成出来的 gRPC client 接口。
    client orderpb.OrderServiceClient
}

func (g *OrderGRPC) UpdateOrder(ctx context.Context, request *orderpb.Order) error {
    // 真正发起 gRPC 调用。
    // 返回值一般是 response + error，这里只关心 error。
    _, err := g.client.UpdateOrder(ctx, request)
    return err
}
```

### 6.3 为什么要包一层

因为 consumer 不应该直接到处写 proto client 调用细节。包一层 adapter 的好处是：

- consumer 只依赖接口
- 测试更容易 mock
- 以后换实现时影响更小

## 7. `stock` 这一条副线：边界开始变干净

### 7.1 `internal/stock/app/query/check_if_items_in_stock.go`

#### 这个文件自己的原始 diff

```diff
diff --git a/internal/stock/app/query/check_if_items_in_stock.go b/internal/stock/app/query/check_if_items_in_stock.go
index 3d3114e..00b5c4c 100644
--- a/internal/stock/app/query/check_if_items_in_stock.go
+++ b/internal/stock/app/query/check_if_items_in_stock.go
@@
-type CheckIfItemsInStock struct {
-	Items []*orderpb.ItemWithQuantity
-}
+type CheckIfItemsInStock struct {
+	Items []*entity.ItemWithQuantity
+}
@@
-type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*orderpb.Item]
+type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*entity.Item]
@@
+	stripeAPI *integration.StripeAPI
@@
-	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*orderpb.Item](
+	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*entity.Item](
@@
-		priceID, ok := stub[i.ID]
-		if !ok {
-			priceID = stub["1"]
-		}
-		res = append(res, &orderpb.Item{
+		priceID, err := h.stripeAPI.GetPriceByProductID(ctx, i.ID)
+		if err != nil || priceID == "" {
+			return nil, err
+		}
+		res = append(res, &entity.Item{
```

#### 这段变化的真正意义

它同时做了两件事：

1. 业务输入输出不再直接用 proto 类型，而是换成 `entity`
2. `priceID` 不再用本地 stub map 硬编码，而是去查 Stripe

#### 带中文注释的关键代码

```go
type CheckIfItemsInStock struct {
    // 从 proto 类型改成 entity 类型，说明应用层开始和传输层解耦。
    Items []*entity.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*entity.Item]

type checkIfItemsInStockHandler struct {
    stockRepo domain.Repository
    stripeAPI *integration.StripeAPI
}

func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*entity.Item, error) {
    var res []*entity.Item

    for _, i := range query.Items {
        // 商品 ID -> Stripe Product -> DefaultPrice ID
        priceID, err := h.stripeAPI.GetPriceByProductID(ctx, i.ID)
        if err != nil || priceID == "" {
            return nil, err
        }

        res = append(res, &entity.Item{
            ID:       i.ID,
            Quantity: i.Quantity,
            PriceID:  priceID,
        })
    }
    return res, nil
}
```

### 7.2 `internal/stock/infrastructure/integration/stripe.go`

#### 这个文件自己的原始 diff

```diff
diff --git a/internal/stock/infrastructure/integration/stripe.go b/internal/stock/infrastructure/integration/stripe.go
new file mode 100644
index 0000000..c24cee9
--- /dev/null
+++ b/internal/stock/infrastructure/integration/stripe.go
@@ -0,0 +1,32 @@
+package integration
+...
+func NewStripeAPI() *StripeAPI {
+	key := viper.GetString("stripe-key")
+	if key == "" {
+		logrus.Fatal("empty key")
+	}
+	return &StripeAPI{apiKey: key}
+}
+
+func (s *StripeAPI) GetPriceByProductID(ctx context.Context, pid string) (string, error) {
+	stripe.Key = s.apiKey
+	result, err := product.Get(pid, &stripe.ProductParams{})
+	if err != nil {
+		return "", err
+	}
+	return result.DefaultPrice.ID, err
+}
+```

#### 带中文注释的代码

```go
type StripeAPI struct {
    apiKey string
}

func NewStripeAPI() *StripeAPI {
    key := viper.GetString("stripe-key")
    if key == "" {
        logrus.Fatal("empty key")
    }
    return &StripeAPI{apiKey: key}
}

func (s *StripeAPI) GetPriceByProductID(ctx context.Context, pid string) (string, error) {
    // stripe.Key 是 Stripe SDK 的全局认证配置。
    stripe.Key = s.apiKey

    // product.Get 会到 Stripe API 拉商品对象。
    result, err := product.Get(pid, &stripe.ProductParams{})
    if err != nil {
        return "", err
    }

    // 取商品默认价格 ID。
    return result.DefaultPrice.ID, err
}
```

#### 这里最该知道的坑

1. `stripe.Key` 是全局状态。
2. 这版代码虽然收了 `ctx`，但没有真正把它传给 Stripe SDK。
3. `stock` 开始依赖外部平台可用性，可靠性风险会上升。

### 7.3 `internal/stock/ports/grpc.go`

#### 这个文件自己的原始 diff

```diff
diff --git a/internal/stock/ports/grpc.go b/internal/stock/ports/grpc.go
index 53b5eaf..b2c659d 100644
--- a/internal/stock/ports/grpc.go
+++ b/internal/stock/ports/grpc.go
@@
-	items, err := G.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{Items: request.Items})
+	items, err := G.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{
+		Items: convertor.NewItemWithQuantityConvertor().ProtosToEntities(request.Items),
+	})
@@
-		Items:   items,
+		Items:   convertor.NewItemConvertor().EntitiesToProtos(items),
```

#### 这段改动的意义

它明确了：

- gRPC 入口层负责做 proto <-> entity 转换
- 应用层 handler 不再直接吃 proto

#### 带中文注释的代码

```go
func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
    _, span := tracing.Start(ctx, "CheckIfItemsInStock")
    defer span.End()

    items, err := G.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{
        // 先把 protobuf 请求转成内部 entity。
        Items: convertor.NewItemWithQuantityConvertor().ProtosToEntities(request.Items),
    })
    if err != nil {
        return nil, err
    }

    return &stockpb.CheckIfItemsInStockResponse{
        InStock: 1,
        // 出口再把 entity 转回 proto。
        Items: convertor.NewItemConvertor().EntitiesToProtos(items),
    }, nil
}
```

### 7.4 `internal/stock/convertor/convertor.go`

#### 这段代码的本质

它就是一个翻译层。职责非常纯：

- proto 转 entity
- entity 转 proto

#### 带中文注释的关键代码

```go
func (c *ItemWithQuantityConvertor) ProtosToEntities(items []*orderpb.ItemWithQuantity) (res []*entity.ItemWithQuantity) {
    for _, i := range items {
        res = append(res, c.ProtoToEntity(i))
    }
    return
}

func (c *ItemWithQuantityConvertor) ProtoToEntity(i *orderpb.ItemWithQuantity) *entity.ItemWithQuantity {
    return &entity.ItemWithQuantity{
        ID:       i.ID,
        Quantity: i.Quantity,
    }
}

func (c *ItemConvertor) EntitiesToProtos(items []*entity.Item) (res []*orderpb.Item) {
    for _, i := range items {
        res = append(res, c.EntityToProto(i))
    }
    return
}

func (c *ItemConvertor) EntityToProto(i *entity.Item) *orderpb.Item {
    return &orderpb.Item{
        ID:       i.ID,
        Name:     i.Name,
        Quantity: i.Quantity,
        PriceID:  i.PriceID,
    }
}
```

### 7.5 `internal/stock/convertor/facade.go`

#### `sync.Once` 是什么

来源：标准库 `sync`

作用：保证某段初始化逻辑只执行一次。

这里为什么用它：作者想让 convertor 作为单例式可复用对象存在。不过这些 convertor 基本无状态，所以这更像写法偏好，不是强需求。

### 7.6 `internal/stock/entity/entity.go`

这段代码为什么重要：

它给 `stock` 明确建立了自己的内部数据模型，意味着：

- `stock` 不再完全借用别人的 proto 类型活着
- 自己开始有真正内部语义对象

## 8. `internal/stock/service/application.go`

### 8.1 这个文件自己的原始 diff

```diff
diff --git a/internal/stock/service/application.go b/internal/stock/service/application.go
index 0ccb6a8..b51f2ad 100644
--- a/internal/stock/service/application.go
+++ b/internal/stock/service/application.go
@@
+	stripeAPI := integration.NewStripeAPI()
@@
-			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, logger, metricsClient),
+			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, stripeAPI, logger, metricsClient),
```

### 8.2 这段改动的意义

这是装配层变化。说明 query handler 不只是“理论上能用 StripeAPI”，而是运行时真的被注入了 StripeAPI。

## 9. `internal/order/adapters/grpc/stock_grpc.go`

### 9.1 这个文件自己的原始 diff

```diff
diff --git a/internal/order/adapters/grpc/stock_grpc.go b/internal/order/adapters/grpc/stock_grpc.go
index c2397d6..549bdd1 100644
--- a/internal/order/adapters/grpc/stock_grpc.go
+++ b/internal/order/adapters/grpc/stock_grpc.go
@@
+	if items == nil {
+		return nil, errors.New("grpc items cannot be nil")
+	}
```

### 9.2 为什么这点改动值得讲

它补了一个很典型的边界保护：gRPC 调用前先检查输入不能为空。

### 9.3 `errors.New()` 是什么

来源：标准库 `errors`

作用：快速创建一个普通 error 值。

## 10. `internal/common/decorator/logging.go`

### 10.1 这个文件自己的原始 diff

```diff
diff --git a/internal/common/decorator/logging.go b/internal/common/decorator/logging.go
index dbe402c..32f2f93 100644
--- a/internal/common/decorator/logging.go
+++ b/internal/common/decorator/logging.go
@@
-	return q.base.Handle(ctx, cmd)
+	result, err = q.base.Handle(ctx, cmd)
+	return result, err
@@
-	return q.base.Handle(ctx, cmd)
+	result, err = q.base.Handle(ctx, cmd)
+	return result, err
```

### 10.2 这到底修了什么

它修的是 `defer` 和命名返回值的配合问题。

原来直接：

```go
return q.base.Handle(ctx, cmd)
```

现在改成：

```go
result, err = q.base.Handle(ctx, cmd)
return result, err
```

这样 `defer` 里的日志判断才能读到正确的 `err`。

## 11. `internal/order/tests/create_order_test.go`

### 11.1 这个文件自己的原始 diff

```diff
diff --git a/internal/order/tests/create_order_test.go b/internal/order/tests/create_order_test.go
index 2af7d2b..d57d2d2 100644
--- a/internal/order/tests/create_order_test.go
+++ b/internal/order/tests/create_order_test.go
@@
-				Id:       "test-item-1",
+				Id:       "prod_R3g7MikGYsXKzr",
@@
+			{
+				Id:       "prod_R285C3Wb7FDprc",
+				Quantity: 10,
+			},
@@
+	t.Logf("getResponse body=%+v", body)
```

### 11.2 这组测试改动在说明什么

它说明测试数据已经开始依赖“真实的 Stripe Product ID 风格数据”，而不是随便写个 `test-item-1`。

## 12. 第三方库和基础设施，这组必须带走的理解

### 12.1 RabbitMQ 在这里扮演什么角色

不是数据库，不是 RPC。它在这里的作用是：

`作为异步事件总线，让 payment 和 kitchen 解耦。`

### 12.2 gRPC 在这里扮演什么角色

它负责 kitchen -> order 的同步回写。

### 12.3 Stripe 在这里扮演什么角色

它现在不只是 payment 侧在用，`stock` 也开始依赖 Stripe 的商品/价格信息。

## 13. 这组还不完美的地方

1. kitchen 没有真正“制作状态机”，只是 `paid -> ready` 一步到位。
2. `cook()` 只是 `Sleep`，没有业务持久化。
3. `stock` 直接查 Stripe product，性能和可靠性都还没优化。
4. `StripeAPI.GetPriceByProductID` 没真正使用传入的 `ctx`。
5. `convertor` 单例化有点过度设计。
6. MQ 消息体里的 `Order` 结构仍然混用了 proto item，边界还没完全清爽。

## 14. 这组最该带走的知识点

1. 新增一个微服务，不只是加个目录，而是要把启动、依赖、消费入口、退出逻辑都补齐。
2. RabbitMQ consumer 的核心循环就是：声明队列、绑定、消费、Ack/Nack。
3. 异步事件链路和同步 gRPC 回写可以同时存在，各自负责不同阶段。
4. 应用层不要长期直接依赖 proto 类型，边界层做 convertor 才更稳。
5. `defer` 和命名返回值的配合在 Go 里非常重要，写错会导致日志判断失真。
6. 教学代码能跑通不等于生产可用，你要学会主动识别课程里的过渡态实现。

## 15. 补充：这一组里的 `channel` 到底在干什么

你刚刚提到这个点是对的。`channel` 如果忘了，后面看 Go 并发代码会很容易断上下文。

这一组里虽然业务主线是 kitchen 和 stock，但真正让 kitchen 服务“活起来”的并发骨架，靠的就是：

- `goroutine`
- `channel`
- 阻塞接收
- 永久阻塞

所以这里单独补一节，把这组里出现的 channel 全部讲透。

### 15.1 先给你一个最短定义

`channel` 是 Go 里用来在不同 goroutine 之间传值的通道。

你可以先把它粗略理解成：

- 队列
- 管道
- 收发口

但这三个类比都不完全准确。Go 里 channel 最大的特点不是“能存数据”，而是：

`它同时也是一种同步机制。`

也就是说，很多时候 channel 不只是拿来传一个值，而是拿来协调“谁先走、谁后走、谁要等谁”。

### 15.2 这一组里出现了哪些 channel

这组里你实际看到了这几个：

1. `sigs := make(chan os.Signal, 1)`
2. `msgs, err := ch.Consume(...)`
3. `var forever chan struct{}`
4. `<-sigs`
5. `<-forever`
6. `for msg := range msgs`

这几个已经覆盖了 Go 初学者最容易混的几类 channel 用法。

### 15.3 `make(chan os.Signal, 1)` 是什么意思

代码：

```go
sigs := make(chan os.Signal, 1)
```

拆开讲：

- `make(...)` 是 Go 用来创建 `slice`、`map`、`channel` 的内建函数
- `chan os.Signal` 表示“这个 channel 里流动的值，类型是 `os.Signal`”
- 最后的 `1` 表示这个 channel 的缓冲区大小是 1

也就是说，这一行的意思是：

`创建一个能传递 os.Signal 值、缓冲区大小为 1 的 channel。`

#### 为什么这里元素类型是 `os.Signal`

因为 `signal.Notify(...)` 会往这个 channel 里塞“系统信号”这个类型的值。

比如：

- `SIGINT`
- `SIGTERM`

#### 为什么缓冲区是 `1`

这就涉及 `channel` 的一个重要概念：

- 无缓冲 channel
- 有缓冲 channel

##### 无缓冲 channel

例如：

```go
ch := make(chan int)
```

没有写第二个参数，默认就是无缓冲。

它的特点是：

- 发送方 `ch <- x` 会阻塞，直到有人接收
- 接收方 `<-ch` 会阻塞，直到有人发送

所以无缓冲 channel 不只是传值，更像“你必须和对方同步握手”。

##### 有缓冲 channel

例如：

```go
ch := make(chan int, 1)
```

它的特点是：

- 只要缓冲区没满，发送可以先放进去，不必立刻有人接收
- 只要缓冲区里有值，接收就能直接拿出来

这里把 `sigs` 设成缓冲 1，意思是：

- 即使接收 goroutine 那一瞬间还没来得及 `<-sigs`
- 系统信号也可以先临时放进去
- 不至于因为“没人马上接”而出问题

这是一种很常见的信号监听写法。

### 15.4 `signal.Notify(sigs, ...)` 和 channel 的关系

代码：

```go
signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
```

这行的意思不是“监听后自动退出”，而是：

`当进程收到 SIGTERM 或 SIGINT 时，把这个信号值发进 sigs 这个 channel。`

也就是说，`signal.Notify` 做的是“转发”动作。

后面这个 goroutine：

```go
go func() {
    <-sigs
    logrus.Infof("receive signal, exiting...")
    os.Exit(0)
}()
```

本质上就是：

- 一直阻塞等 `sigs` 里来一个值
- 来了之后就执行退出逻辑

这里的 `<-sigs` 不是“语法糖”，而是一个真正的阻塞点。

#### `<-sigs` 到底在做什么

`<-ch` 这种写法表示：

- 从 channel `ch` 接收一个值

如果你不接收赋值给变量，而是直接单独写：

```go
<-sigs
```

意思就是：

- 我只关心“等到一个值来”
- 我不关心这个值具体叫什么

在这段代码里，这是一种很常见的写法，因为它只是想等“收到退出信号”这个事件发生。

### 15.5 `msgs, err := ch.Consume(...)` 为什么也是 channel

这行最容易让初学者忽略，因为它看起来不像 `make(chan ...)`。

但本质上它返回的就是一个 channel。

你后面看到：

```go
for msg := range msgs {
    c.handleMessage(ch, msg, q)
}
```

只要能 `range`，而且每次迭代拿一个消息，那你就可以判断：

`msgs` 是一个 channel。

这里它的真实含义是：

- RabbitMQ 客户端库在底层不断从 broker 拉消息
- 每来一条消息，就往 `msgs` 这个 Go channel 里塞一个 `amqp.Delivery`
- 你的代码只需要像读普通 channel 一样去消费它

所以 `Consume()` 其实帮你把“外部 MQ 消息流”适配成了“Go channel 流”。

这也是 Go 写并发 I/O 时很常见的风格。

### 15.6 `for msg := range msgs` 到底是什么意思

代码：

```go
for msg := range msgs {
    c.handleMessage(ch, msg, q)
}
```

这行你可以翻译成：

`只要 msgs 这个 channel 还在持续产出值，我就不断从里面取一个 msg 出来处理。`

这里的执行模型是：

1. `msgs` 里来一条 RabbitMQ 消息
2. `range msgs` 读到它
3. 赋值给变量 `msg`
4. 调 `handleMessage(...)`
5. 回来继续等下一条

#### 为什么 `range channel` 很常见

因为它天然适合处理“持续流入”的事件流：

- 消息队列
- worker 任务流
- 网络事件
- 定时事件

#### `range channel` 什么时候会结束

只有一种情况：

`channel 被 close，并且里面剩余数据也读完了。`

这也是为什么很多 Go 代码会用 `for range ch` 做 worker loop。

### 15.7 `var forever chan struct{}` 为什么会永久阻塞

代码：

```go
var forever chan struct{}
...
<-forever
```

这段是这组里最值得单独讲的 channel 细节之一。

#### 第一步：`var forever chan struct{}` 到底创建了什么

这行只是“声明”了一个 channel 变量，**并没有真正创建 channel**。

也就是说，这时候：

- `forever` 的类型是 `chan struct{}`
- 但它的值是 `nil`

这和下面这种写法不一样：

```go
forever := make(chan struct{})
```

后者才是真的创建了一个可用 channel。

#### 第二步：`chan struct{}` 又是什么

`struct{}` 是空结构体。

它的特点是：

- 不占实际数据空间
- 常被 Go 程序员用来表示“这里只需要一个信号，不需要实际负载数据”

所以：

```go
chan struct{}
```

通常表示：

- 这是一个只用于同步/通知的 channel
- 不是为了传业务数据

#### 第三步：为什么 `<-forever` 会永久阻塞

因为：

- `forever` 是 `nil channel`
- 对 `nil channel` 做发送和接收，都会永久阻塞

这是 Go channel 很容易忘记、但非常重要的一条规则：

1. 向 `nil channel` 发送：永久阻塞
2. 从 `nil channel` 接收：永久阻塞
3. `close(nil channel)`：直接 panic

所以这里：

```go
<-forever
```

不是作者写错了语法，而是故意利用了这个特性来“卡住当前 goroutine 不退出”。

#### 这种写法好不好

能工作，但不优雅。

因为它的可读性一般。对初学者来说，看到这里很容易以为：

- 是不是忘了初始化
- 是不是 bug

更直观的阻塞写法通常是：

```go
select {}
```

或者用 `context` / `WaitGroup` 做生命周期管理。

所以这里你要学到的是：

`它不是错误，而是一种利用 nil channel 阻塞特性的写法。`

### 15.8 `chan struct{}` 为什么经常被拿来做“通知信号”

这是 Go 里一个很常见的约定俗成用法。

比如：

```go
done := make(chan struct{})
```

然后某处：

```go
close(done)
```

另一处：

```go
<-done
```

意思往往是：

- 我不需要传“内容”
- 我只需要传“事件发生了”

因为空结构体没有负载，所以它非常适合做：

- done 信号
- stop 信号
- ready 信号

### 15.9 `channel` 和 `goroutine` 的关系是什么

这两个概念经常被一起提，但你要区分清楚：

- `goroutine`：执行单元，谁来跑代码
- `channel`：通信手段，goroutine 之间怎么传值和同步

在这一组里：

#### 场景 1：监听系统信号

```go
go func() {
    <-sigs
    os.Exit(0)
}()
```

这里表示：

- 新起一个 goroutine 专门等退出信号
- 主 goroutine 继续干别的事
- 两者之间靠 `sigs` 这个 channel 协调

#### 场景 2：消费 RabbitMQ 消息

```go
go func() {
    for msg := range msgs {
        c.handleMessage(ch, msg, q)
    }
}()
```

这里表示：

- 新起一个 goroutine 专门消费消息
- `msgs` channel 持续给它喂消息
- 它收到一条处理一条

所以你可以把 channel 想成 goroutine 之间的“交通管道”。

### 15.10 阻塞到底是什么意思

Go 里经常说“阻塞”，你一定要把这个词真正理解。

阻塞不是报错，也不是卡死整个程序。

阻塞的准确意思是：

`当前 goroutine 暂时停在这里，等条件满足再继续。`

比如：

```go
x := <-ch
```

如果 `ch` 里现在没有值，那么：

- 当前 goroutine 会停在这里等待
- 其他 goroutine 仍然可以继续运行

所以在这组里：

- `<-sigs` 是“等一个系统信号”
- `<-forever` 是“永远等不到，所以一直卡着”
- `for msg := range msgs` 是“不断等下一条消息”

### 15.11 缓冲 channel 和非缓冲 channel，再用这组代码重新理解一次

#### `sigs := make(chan os.Signal, 1)`

这是有缓冲 channel。

含义：

- 最多可以先缓存 1 个信号
- 发送方不必马上等接收方准备好

#### `msgs`

这是库返回的 channel，你可以先把它理解成“客户端帮你维护的消息流 channel”。

它具体是缓冲还是非缓冲，不需要你在这层代码里直接控制，但你要知道：

- 你对它的读取速度会影响消息处理节奏

#### `forever`

这里根本不是缓冲/非缓冲的问题，因为它是 `nil channel`，不是“可用 channel”。

### 15.12 为什么 channel 容易死锁

最常见的几种情况：

1. 没人发，你却一直收
2. 没人收，你却一直发
3. 所有 goroutine 都在互相等
4. 对 `nil channel` 做收发

在这组里，`<-forever` 之所以没变成“程序 bug”，是因为作者就是故意想永久卡住它。

但如果你在别的业务代码里无意中写出了类似模式，就会非常危险。

### 15.13 这组里你应该带走的 channel 直觉

看完这组，你可以先记住这几条：

1. `channel` 不只是传值工具，也是同步工具。
2. `make(chan T)` 才是真的创建 channel，`var ch chan T` 默认是 `nil`。
3. `<-ch` 表示接收，`ch <- x` 表示发送。
4. `for v := range ch` 适合处理持续流入的数据流。
5. `nil channel` 上收发会永久阻塞。
6. `chan struct{}` 经常不是传数据，而是传“一个事件信号”。
7. `goroutine` 决定“谁在跑”，`channel` 决定“他们怎么通信”。

### 15.14 回到这一组代码，重新用一句话总结

- `sigs`：负责把操作系统的退出信号送给退出 goroutine
- `msgs`：负责把 RabbitMQ 的消息送给消费 goroutine
- `forever`：负责把当前 goroutine 永久卡住，不让 `Listen()` 提前退出

如果你愿意，我下一步可以继续把这份讲义里所有跟 `goroutine` 有关的地方也补成同样粒度，比如：

- `go func()` 究竟什么时候并发
- 为什么这里不会自动开新线程
- goroutine 和线程到底什么关系
- 为什么 `time.Sleep` 只是阻塞当前 goroutine，不是卡死整个程序

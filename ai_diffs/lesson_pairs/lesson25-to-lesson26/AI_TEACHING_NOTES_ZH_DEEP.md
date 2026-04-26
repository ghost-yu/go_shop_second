# `lesson25 -> lesson26` 独立讲义（详细注释版）

这一组不是在“新增一个业务功能”，而是在把上一组 `lesson24 -> lesson25` 只做到一半的链路追踪继续补完。

上一组你已经把这些东西接上了：
- HTTP 入口可以自动起 span
- gRPC server / client 可以自动起 span
- trace 数据可以送到 Jaeger

但是还有一个明显断点：

`order 创建订单 -> 发 MQ -> payment 消费消息 -> 生成支付链接 -> webhook -> 再发 MQ -> order 消费 paid 事件`

这里面有两段是 RabbitMQ 异步消息。

而 OpenTelemetry 对 HTTP/gRPC 这种“标准协议”可以靠现成中间件自动帮你传递 trace 上下文，
但对 RabbitMQ 这种消息队列，这个项目当前并没有直接用一个现成的自动接入方案，所以这组课要做的事情是：

`手动把 trace 上下文塞进 MQ header，再在消费端从 header 里取出来，这样异步链路上的 span 才不会断。`

这也是字幕 `26补充.txt` 的主线：
- 之前 MQ 后面的链路消失了
- 需要手动 instrument 一个 span
- 需要理解 propagator / carrier / inject / extract
- 需要把 context 放进 header，再从 header 取出来

你可以把这一组理解成：

`把“同步链路可观测”继续推进到“异步消息链路也能接上 trace”。`

## 1. 这一组你应该先看什么

建议阅读顺序：

1. [internal/common/broker/rabbitmq.go](/g:/shi/go_shop_second/internal/common/broker/rabbitmq.go)
2. [internal/order/app/command/create_order.go](/g:/shi/go_shop_second/internal/order/app/command/create_order.go)
3. [internal/payment/infrastructure/consumer/consumer.go](/g:/shi/go_shop_second/internal/payment/infrastructure/consumer/consumer.go)
4. [internal/payment/infrastructure/processor/stripe.go](/g:/shi/go_shop_second/internal/payment/infrastructure/processor/stripe.go)
5. [internal/payment/adapters/order_grpc.go](/g:/shi/go_shop_second/internal/payment/adapters/order_grpc.go)
6. [internal/payment/http.go](/g:/shi/go_shop_second/internal/payment/http.go)
7. [internal/order/infrastructure/consumer/consumer.go](/g:/shi/go_shop_second/internal/order/infrastructure/consumer/consumer.go)
8. [internal/stock/ports/grpc.go](/g:/shi/go_shop_second/internal/stock/ports/grpc.go)

为什么这样看：
- `rabbitmq.go` 先解释基础设施：MQ header 怎么承载 trace 上下文。
- `create_order.go` 是第一处“发送 MQ 时注入 header”的生产者。
- `payment consumer` 是第一处“消费 MQ 时提取 header”的消费者。
- 然后再看 `stripe`、`order_grpc`，你会明白为什么 trace 能继续串到 Stripe 创建支付链接、再串到回写 order 的 gRPC 调用。
- 最后再看 `payment/http.go` 和 `order consumer`，把第二条异步链也补齐。
- `stock/ports/grpc.go` 是同步 gRPC 入口补 span，和这一组 “more otel” 属于同一主题。

## 2. 先把整条链路想清楚

这一组完成后，链路大致变成这样：

```text
客户端请求 order 创建订单
-> order HTTP span
-> order 内部业务 span
-> order 发布 order.created 消息时，手动起一个 rabbitmq publish span
-> 把 trace 上下文注入到 MQ headers
-> payment consumer 收到消息
-> 从 MQ headers 提取 trace 上下文
-> payment consumer 手动起一个 rabbitmq consume span
-> payment 创建 Stripe 支付链接
-> payment 通过 gRPC 回写 order
-> 后面 Stripe webhook 再触发另一条 MQ publish
-> order consumer 再 extract header，继续接上 trace
```

你要注意：

`异步消息本质上没有像普通函数调用那样天然共享同一个 context。`

所以如果你什么都不做：
- 生产者这边有 trace
- 消费者这边也许自己能起一个新 trace
- 但这两边不会自动连成一条链

这就是“链路断了”的本质。

## 3. 先补最底层基础设施：`internal/common/broker/rabbitmq.go`

### 3.1 这个文件自己的原始 diff

```diff
diff --git a/internal/common/broker/rabbitmq.go b/internal/common/broker/rabbitmq.go
index 599af4e..dd6f8b9 100644
--- a/internal/common/broker/rabbitmq.go
+++ b/internal/common/broker/rabbitmq.go
@@ -1,10 +1,12 @@
 package broker
 
 import (
+	"context"
 	"fmt"
 
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
+	"go.opentelemetry.io/otel"
 )
@@ -27,3 +29,37 @@ func Connect(user, password, host, port string) (*amqp.Channel, func() error) {
 	}
 	return ch, conn.Close
 }
+
+type RabbitMQHeaderCarrier map[string]interface{}
+
+func (r RabbitMQHeaderCarrier) Get(key string) string {
+	value, ok := r[key]
+	if !ok {
+		return ""
+	}
+	return value.(string)
+}
+
+func (r RabbitMQHeaderCarrier) Set(key string, value string) {
+	r[key] = value
+}
+
+func (r RabbitMQHeaderCarrier) Keys() []string {
+	keys := make([]string, len(r))
+	i := 0
+	for key := range r {
+		keys[i] = key
+		i++
+	}
+	return keys
+}
+
+func InjectRabbitMQHeaders(ctx context.Context) map[string]interface{} {
+	carrier := make(RabbitMQHeaderCarrier)
+	otel.GetTextMapPropagator().Inject(ctx, carrier)
+	return carrier
+}
+
+func ExtractRabbitMQHeaders(ctx context.Context, headers map[string]interface{}) context.Context {
+	return otel.GetTextMapPropagator().Extract(ctx, RabbitMQHeaderCarrier(headers))
+}
```

### 3.2 旧代码做什么，新代码做什么

旧代码：
- 这个文件主要负责连 RabbitMQ，声明 exchange，做一些 broker 级公共能力。
- 它还不知道 trace 上下文怎么在 MQ 里传递。

新代码：
- 增加了 `RabbitMQHeaderCarrier`。
- 让 RabbitMQ 的 `headers` 结构可以伪装成 OpenTelemetry 需要的 `TextMapCarrier`。
- 增加 `InjectRabbitMQHeaders`：把 `context` 里的 trace 信息写进 MQ headers。
- 增加 `ExtractRabbitMQHeaders`：从 MQ headers 读回 trace 信息，恢复到新的 `context`。

一句话：

`这段代码是在做“协议适配”。它把 RabbitMQ 的 header 结构适配成了 OpenTelemetry propagator 能理解的载体。`

### 3.3 关键代码和详细注释

```go
package broker

import (
    "context" // 这里需要传入/返回 context，因为 trace 上下文都挂在 context 里
    "fmt"

    amqp "github.com/rabbitmq/amqp091-go"
    "github.com/sirupsen/logrus"
    "go.opentelemetry.io/otel" // 这里提供全局 propagator
)

// RabbitMQHeaderCarrier 本质上就是对 map[string]interface{} 起了一个新类型名。
// 这样做不是为了“好看”，而是为了给这个类型挂方法，让它满足 otel 的 TextMapCarrier 接口。
type RabbitMQHeaderCarrier map[string]interface{}

func (r RabbitMQHeaderCarrier) Get(key string) string {
    value, ok := r[key]
    if !ok {
        // 没有这个 key，按接口约定返回空字符串。
        return ""
    }

    // 这里默认把 header 的值断言成 string。
    // 这在当前 trace 传播场景里通常成立，因为 traceparent / tracestate 都是字符串。
    // 但这也是一个容易踩坑的点：
    // 如果某个 header 实际不是 string，这里会 panic。
    return value.(string)
}

func (r RabbitMQHeaderCarrier) Set(key string, value string) {
    // inject 时会走到这里，把 trace 信息写入 map。
    r[key] = value
}

func (r RabbitMQHeaderCarrier) Keys() []string {
    // otel 可能需要遍历已有 keys，这里把 map 的 key 全部吐出来。
    keys := make([]string, len(r))
    i := 0
    for key := range r {
        keys[i] = key
        i++
    }
    return keys
}

func InjectRabbitMQHeaders(ctx context.Context) map[string]interface{} {
    // 先准备一个 carrier，它本质就是一个 map。
    carrier := make(RabbitMQHeaderCarrier)

    // 这一步非常关键：
    // 把当前 context 里的 trace 信息，通过当前全局 propagator 注入进 carrier。
    // 注入后，carrier 里通常会出现 traceparent 之类的键值。
    otel.GetTextMapPropagator().Inject(ctx, carrier)

    // 最后返回给 MQ Publish 使用。
    return carrier
}

func ExtractRabbitMQHeaders(ctx context.Context, headers map[string]interface{}) context.Context {
    // 这一步和 Inject 相反：
    // 把 MQ header 里的 trace 信息读出来，恢复进一个新的 context。
    // 这样 consumer 继续往后调用时，就还能接在原 trace 上。
    return otel.GetTextMapPropagator().Extract(ctx, RabbitMQHeaderCarrier(headers))
}
```

### 3.4 这里你必须懂的几个基础概念

#### `propagator` 是什么

你可以把它理解成：

`专门负责“把 trace 上下文写出去 / 读回来”的编解码器。`

它不负责起 span，也不负责上报 Jaeger。
它只负责“传播上下文”。

对应两个动作：
- `Inject`：把上下文写到外部载体里
- `Extract`：从外部载体里恢复上下文

#### `carrier` 是什么

carrier 不是特指某个库，它只是一个“承载上下文的容器”。

在不同协议里 carrier 不一样：
- HTTP 里通常是 header
- gRPC 里通常是 metadata
- RabbitMQ 里这里就是 `Publishing.Headers` / `Delivery.Headers`

#### 为什么这里要自己实现 `Get/Set/Keys`

因为 OpenTelemetry 的 propagator 不认识 RabbitMQ 原生 header 类型。
它只认识一个抽象接口：`TextMapCarrier`。

所以你要自己告诉它：
- 怎么取值
- 怎么写值
- 怎么列出 key

这就是适配器模式，非常典型。

#### 容易错的点

1. `value.(string)` 有 panic 风险
- 如果 header 里某个值不是字符串，这里会直接崩。
- 更稳妥的写法通常会做类型判断。

2. `headers == nil` 的情况
- 当前 `ExtractRabbitMQHeaders` 里如果传入 `nil`，一般问题不大，但你要知道有些消息未必带 header。
- 真实生产代码常会更防御式。

3. 这只是“传播上下文”，不是“自动起 span”
- inject/extract 本身不会创建 span。
- 你还是得在 publish / consume 两端自己 `Start(...)`。

## 4. 生产者第一站：`internal/order/app/command/create_order.go`

### 4.1 这个文件自己的原始 diff

```diff
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
index 7185157..282f66c 100644
--- a/internal/order/app/command/create_order.go
+++ b/internal/order/app/command/create_order.go
@@ -4,6 +4,7 @@ import (
 	"context"
 	"encoding/json"
 	"errors"
+	"fmt"
@@ -12,6 +13,7 @@ import (
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
+	"go.opentelemetry.io/otel"
 )
@@ -59,6 +61,15 @@ func NewCreateOrderHandler(
 }
 
 func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
+	q, err := c.channel.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
+	if err != nil {
+		return nil, err
+	}
+
+	t := otel.Tracer("rabbitmq")
+	ctx, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", q.Name))
+	defer span.End()
+
 	validItems, err := c.validate(ctx, cmd.Items)
 	if err != nil {
 		return nil, err
@@ -71,19 +82,16 @@ func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*Creat
 		return nil, err
 	}
 
-	q, err := c.channel.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
-	if err != nil {
-		return nil, err
-	}
-
 	marshalledOrder, err := json.Marshal(o)
 	if err != nil {
 		return nil, err
 	}
+	header := broker.InjectRabbitMQHeaders(ctx)
 	err = c.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
 		ContentType:  "application/json",
 		DeliveryMode: amqp.Persistent,
 		Body:         marshalledOrder,
+		Headers:      header,
 	})
 	if err != nil {
 		return nil, err
```

### 4.2 旧代码做什么，新代码做什么

旧代码：
- 校验库存
- 创建订单
- `json.Marshal(o)`
- 直接 `PublishWithContext(...)` 发到 RabbitMQ

旧代码的问题不是“消息发不出去”，而是：
- 消息发出去之后，trace 没有跟着过去
- 后面的 payment consumer 虽然在执行，但 Jaeger 里链路断开了

新代码多做了两件关键事：
1. 发布消息之前，手动起一个 `rabbitmq.xxx.publish` span
2. 把当前 `ctx` 里的 trace 上下文注入到 MQ `Headers`

这两步缺一不可：
- 只起 span 不注入 header：消费者接不上
- 只注入 header 不起 publish span：你看不到“发消息”这一段动作

### 4.3 关键代码和详细注释

```go
func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
    // 先声明/确保队列存在。
    // 这里把 QueueDeclare 提到前面，是因为后面起 span 时需要拿到 q.Name 作为 span 名称的一部分。
    q, err := c.channel.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
    if err != nil {
        return nil, err
    }

    // 这里不是用 HTTP/gRPC 自动中间件，而是手动拿一个 tracer。
    // 名字写成 rabbitmq，表示这类 span 属于 MQ 这一层。
    t := otel.Tracer("rabbitmq")

    // Start 之后会返回一个新的 ctx。
    // 这个新 ctx 非常重要，因为它携带了“当前 publish span”这个上下文。
    // 后面 inject 的必须是这个新的 ctx，而不是旧的 ctx。
    ctx, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", q.Name))
    defer span.End()

    // 这里继续做业务校验。由于已经换成新的 ctx，
    // 后续校验、调用 stock 等如果内部也打点，就会挂在这个链路下面。
    validItems, err := c.validate(ctx, cmd.Items)
    if err != nil {
        return nil, err
    }

    pendingOrder, err := domain.NewPendingOrder(cmd.CustomerID, validItems)
    if err != nil {
        return nil, err
    }

    o, err := service.NewOrderDomainService(c.orderRepo, c.eventPublisher).CreateOrder(ctx, *pendingOrder)
    if err != nil {
        return nil, err
    }

    marshalledOrder, err := json.Marshal(o)
    if err != nil {
        return nil, err
    }

    // 这一步是整组课最关键的一步之一。
    // 它会把 ctx 里的 trace 信息编码到一个 map 里。
    header := broker.InjectRabbitMQHeaders(ctx)

    err = c.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
        ContentType:  "application/json",
        DeliveryMode: amqp.Persistent,
        Body:         marshalledOrder,
        Headers:      header, // 把 trace 上下文塞进 MQ headers
    })
    if err != nil {
        return nil, err
    }

    return &CreateOrderResult{OrderID: o.ID}, nil
}
```

### 4.4 为什么 `QueueDeclare` 要被提前

这点你很容易忽略。

旧代码里 `QueueDeclare` 在后面也能工作。
新代码把它提前，不是业务逻辑变化，而是为了：
- 拿到 `q.Name`
- 用这个名字生成 span 名
- 让 span 名更贴近真实队列

这是观测性代码经常会做的事情：

`为了让 trace 更可读，会把资源名（queue、topic、service、method）带进 span 名或 attributes。`

### 4.5 为什么必须把 `ctx` 接回来

这也是 Go 小白最容易丢的点。

`ctx, span := t.Start(ctx, ...)`

不是随便写的。

因为 `Start` 返回的新 `ctx` 里挂着“当前 span 是谁”的信息。
你如果这样写：

```go
_, span := t.Start(ctx, ...)
header := broker.InjectRabbitMQHeaders(ctx)
```

那你 inject 的还是旧 `ctx`，新 span 可能不会被正确传播。

所以：
- `Start` 之后要接住新 `ctx`
- inject / 后续调用都尽量用这个新 `ctx`

## 5. 第一段消费链：`internal/payment/infrastructure/consumer/consumer.go`

### 5.1 这个文件自己的原始 diff

```diff
diff --git a/internal/payment/infrastructure/consumer/consumer.go b/internal/payment/infrastructure/consumer/consumer.go
index 7f01e7a..99046de 100644
--- a/internal/payment/infrastructure/consumer/consumer.go
+++ b/internal/payment/infrastructure/consumer/consumer.go
@@ -3,6 +3,7 @@ package consumer
 import (
 	"context"
 	"encoding/json"
+	"fmt"
@@ -10,6 +11,7 @@ import (
 	"github.com/ghost-yu/go_shop_second/payment/app/command"
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
+	"go.opentelemetry.io/otel"
 )
@@ -44,6 +46,10 @@ func (c *Consumer) Listen(ch *amqp.Channel) {
 
 func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
 	logrus.Infof("Payment receive a message from %s, msg=%v", q.Name, string(msg.Body))
+	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
+	tr := otel.Tracer("rabbitmq")
+	_, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
+	defer span.End()
 
 	o := &orderpb.Order{}
 	if err := json.Unmarshal(msg.Body, o); err != nil {
@@ -51,13 +57,14 @@ func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
 		_ = msg.Nack(false, false)
 		return
 	}
-	if _, err := c.app.Commands.CreatePayment.Handle(context.TODO(), command.CreatePayment{Order: o}); err != nil {
+	if _, err := c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{Order: o}); err != nil {
 		// TODO: retry
 		logrus.Infof("failed to create order, err=%v", err)
 		_ = msg.Nack(false, false)
 		return
 	}
 
+	span.AddEvent("payment.created")
 	_ = msg.Ack(false)
 	logrus.Info("consume success")
 }
```

### 5.2 旧代码做什么，新代码做什么

旧代码：
- 收到 MQ 消息
- 反序列化为订单
- 用 `context.TODO()` 调 `CreatePayment`
- 成功就 Ack

问题：
- `context.TODO()` 等于明确告诉别人“这里上下文我还没想好”
- 这会直接丢掉生产者带来的 trace 上下文
- payment 这边的链路会变成一个新的孤岛

新代码：
- 先从 `msg.Headers` 提取 trace 上下文
- 再手动起一个 consume span
- 再把这个 `ctx` 传给 `CreatePayment`
- 成功后给 span 增加事件 `payment.created`

这相当于把“消费消息 -> 创建支付”这一段挂回原链路上。

### 5.3 关键代码和详细注释

```go
func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
    logrus.Infof("Payment receive a message from %s, msg=%v", q.Name, string(msg.Body))

    // 第一步：从 MQ headers 里恢复 trace 上下文。
    // 注意这里不是继续沿用某个旧 ctx，因为消息消费是异步发生的，
    // 当前 goroutine 本来并没有那个请求上下文。
    ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)

    tr := otel.Tracer("rabbitmq")

    // 第二步：基于刚恢复出来的 ctx，再起一个“consume” span。
    // 这样在 Jaeger 里你能清楚看到：生产者 publish -> 消费者 consume
    _, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
    defer span.End()

    o := &orderpb.Order{}
    if err := json.Unmarshal(msg.Body, o); err != nil {
        logrus.Infof("error unmarshal msg.body into orderpb.order, err = %v", err)
        _ = msg.Nack(false, false)
        return
    }

    // 第三步：不能再用 context.TODO()。
    // 必须把刚 extract 回来的 ctx 继续往下传。
    // 否则 CreatePayment 里面的 trace 会重新断掉。
    if _, err := c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{Order: o}); err != nil {
        logrus.Infof("failed to create order, err=%v", err)
        _ = msg.Nack(false, false)
        return
    }

    // 额外补一个事件，方便在 trace 上看出“支付已创建”这个业务节点。
    span.AddEvent("payment.created")

    _ = msg.Ack(false)
    logrus.Info("consume success")
}
```

### 5.4 为什么这里 `Extract` 之后还要 `Start`

这是另一个很容易误解的点。

`Extract` 的作用只是“恢复父上下文”，不是“自动创建一个新的 consume span”。

所以完整步骤是：
1. `Extract`：先把父 trace 找回来
2. `Start`：在这个父上下文下面新建当前消费动作的 span

如果只 `Extract` 不 `Start`：
- 下游可能还能接上 trace
- 但你看不到“消费者处理消息”这一层本身

如果只 `Start` 不 `Extract`：
- 你会有 consume span
- 但它会变成一条新的 trace，不会接在生产者后面

## 6. 支付处理器继续往下透传：`internal/payment/infrastructure/processor/stripe.go`

### 6.1 这个文件自己的原始 diff

```diff
diff --git a/internal/payment/infrastructure/processor/stripe.go b/internal/payment/infrastructure/processor/stripe.go
index 37fcb1d..e597fe5 100644
--- a/internal/payment/infrastructure/processor/stripe.go
+++ b/internal/payment/infrastructure/processor/stripe.go
@@ -6,6 +6,7 @@ import (
 	"fmt"
 
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/common/tracing"
 	"github.com/stripe/stripe-go/v79"
 	"github.com/stripe/stripe-go/v79/checkout/session"
 )
@@ -27,6 +28,9 @@ const (
 )
 
 func (s StripeProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
+	_, span := tracing.Start(ctx, "stripe_processor.create_payment_link")
+	defer span.End()
+
 	var items []*stripe.CheckoutSessionLineItemParams
 	for _, item := range order.Items {
 		items = append(items, &stripe.CheckoutSessionLineItemParams{
```

### 6.2 这段改动为什么重要

这段看起来只加了 2 行 tracing，但意义很大。

因为 payment consumer 把 MQ trace 接回来之后，如果这里不继续打 span，
你在 Jaeger 里仍然只会看到：
- 收到消息
- 然后下一步没有细节

新代码等于告诉你：

`现在“真正调用 Stripe 创建支付链接”这一步，也被显式标出来了。`

### 6.3 关键代码和详细注释

```go
func (s StripeProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
    // 这里不是自动中间件能覆盖到的地方，因为这是普通业务函数。
    // 所以要手动起 span，标记“调用 Stripe 创建支付链接”这个动作。
    _, span := tracing.Start(ctx, "stripe_processor.create_payment_link")
    defer span.End()

    var items []*stripe.CheckoutSessionLineItemParams
    for _, item := range order.Items {
        items = append(items, &stripe.CheckoutSessionLineItemParams{
            Price:    stripe.String(item.PriceID),
            Quantity: stripe.Int64(int64(item.Quantity)),
        })
    }

    // 后面省略的是拼 metadata、组装 CheckoutSessionParams、调用 session.New(params)。
    // 因为这整段逻辑现在都跑在这个 span 下面，所以在 Jaeger 里能看到它属于 payment consumer 链路的一部分。
}
```

### 6.4 为什么业务函数里也要手动打 span

因为“自动 tracing”只会覆盖它知道的边界：
- HTTP 请求入口
- gRPC 服务端入口
- gRPC 客户端调用

但普通函数，比如：
- `CreatePaymentLink`
- `publish event`
- `validate order`

这些都不会自动有 span。

如果你想在可观测性上看到更细粒度的业务节点，就得自己打。

## 7. payment 回写 order：`internal/payment/adapters/order_grpc.go`

### 7.1 这个文件自己的原始 diff

```diff
diff --git a/internal/payment/adapters/order_grpc.go b/internal/payment/adapters/order_grpc.go
index 4ee8d88..3bd171d 100644
--- a/internal/payment/adapters/order_grpc.go
+++ b/internal/payment/adapters/order_grpc.go
@@ -4,6 +4,7 @@ import (
 	"context"
 
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/common/tracing"
 	"github.com/sirupsen/logrus"
 )
@@ -16,6 +17,9 @@ func NewOrderGRPC(client orderpb.OrderServiceClient) *OrderGRPC {
 }
 
 func (o OrderGRPC) UpdateOrder(ctx context.Context, order *orderpb.Order) error {
+	ctx, span := tracing.Start(ctx, "order_grpc.update_order")
+	defer span.End()
+
 	_, err := o.client.UpdateOrder(ctx, order)
 	logrus.Infof("payment_adapter||update_order,err=%v", err)
 	return err
```

### 7.2 为什么这里还要再包一层 span

你可能会问：
- gRPC client 不是上一组已经加了 otel handler 吗
- 那这里为什么还要手动 `tracing.Start(...)`

答案是：

`这两个层次不冲突。`

- gRPC client handler 负责网络层、协议层的调用追踪
- 这里的手动 span 负责业务语义层：我现在要“回写订单”

所以你最终会看到两层信息：
- 一个更业务化的 `order_grpc.update_order`
- 一个更底层的 gRPC client 调用 span

这对排查问题很有帮助。

### 7.3 关键代码和详细注释

```go
func (o OrderGRPC) UpdateOrder(ctx context.Context, order *orderpb.Order) error {
    // 业务语义层的 span：表示 payment 服务准备回写 order。
    ctx, span := tracing.Start(ctx, "order_grpc.update_order")
    defer span.End()

    // 注意这里继续把新的 ctx 往下传，
    // 这样底层 gRPC client instrumentation 才能挂在这个 span 下面。
    _, err := o.client.UpdateOrder(ctx, order)
    logrus.Infof("payment_adapter||update_order,err=%v", err)
    return err
}
```

## 8. 第二段发布链：`internal/payment/http.go`

### 8.1 这个文件自己的原始 diff

```diff
diff --git a/internal/payment/http.go b/internal/payment/http.go
index 9017d2f..884e2c9 100644
--- a/internal/payment/http.go
+++ b/internal/payment/http.go
@@ -3,6 +3,7 @@ package main
 import (
 	"context"
 	"encoding/json"
+	"fmt"
 	"io"
 	"net/http"
@@ -15,6 +16,7 @@ import (
 	"github.com/spf13/viper"
 	"github.com/stripe/stripe-go/v79"
 	"github.com/stripe/stripe-go/v79/webhook"
+	"go.opentelemetry.io/otel"
 )
@@ -80,10 +82,16 @@ func (h *PaymentHandler) handleWebhook(c *gin.Context) {
 				return
 			}
 
-			_ = h.channel.PublishWithContext(ctx, broker.EventOrderPaid, "", false, false, amqp.Publishing{
+			tr := otel.Tracer("rabbitmq")
+			mqCtx, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", broker.EventOrderPaid))
+			defer span.End()
+
+			headers := broker.InjectRabbitMQHeaders(mqCtx)
+			_ = h.channel.PublishWithContext(mqCtx, broker.EventOrderPaid, "", false, false, amqp.Publishing{
 				ContentType:  "application/json",
 				DeliveryMode: amqp.Persistent,
 				Body:         marshalledOrder,
+				Headers:      headers,
 			})
 			logrus.Infof("message published to %s, body: %s", broker.EventOrderPaid, string(marshalledOrder))
 		}
```

### 8.2 这段在做什么

这段和 `create_order.go` 本质上是同一个模式，只是它发生在另一条业务链上：

`Stripe webhook -> payment 服务 -> 发布 order.paid 事件`

旧代码里：
- 能发消息
- 但 trace 不会被带到 order consumer

新代码里：
- 给 publish 动作手动起 span
- 注入 MQ headers
- 把这个链路继续传给 order consumer

### 8.3 关键代码和详细注释

```go
tr := otel.Tracer("rabbitmq")

// 这里使用 ctx 作为父上下文。
// 这个 ctx 来自当前 HTTP webhook 请求链，所以 MQ publish 会挂在 webhook 请求下。
mqCtx, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", broker.EventOrderPaid))
defer span.End()

// 把刚刚这个 mqCtx 的 trace 信息塞进 headers
headers := broker.InjectRabbitMQHeaders(mqCtx)

_ = h.channel.PublishWithContext(mqCtx, broker.EventOrderPaid, "", false, false, amqp.Publishing{
    ContentType:  "application/json",
    DeliveryMode: amqp.Persistent,
    Body:         marshalledOrder,
    Headers:      headers,
})
```

### 8.4 为什么这里变量名特意叫 `mqCtx`

这是个很好的小细节。

它在语义上提醒你：
- 这是基于当前 webhook `ctx` 派生出来的
- 但它特别服务于当前这段 MQ publish 动作

虽然技术上你也可以继续复用 `ctx` 变量名，但拆成 `mqCtx` 更不容易读混。

## 9. order 消费 paid 事件：`internal/order/infrastructure/consumer/consumer.go`

### 9.1 这个文件自己的原始 diff

```diff
diff --git a/internal/order/infrastructure/consumer/consumer.go b/internal/order/infrastructure/consumer/consumer.go
index e7ce252..623bb50 100644
--- a/internal/order/infrastructure/consumer/consumer.go
+++ b/internal/order/infrastructure/consumer/consumer.go
@@ -3,6 +3,7 @@ package consumer
 import (
 	"context"
 	"encoding/json"
+	"fmt"
@@ -10,6 +11,7 @@ import (
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
+	"go.opentelemetry.io/otel"
 )
@@ -38,20 +40,25 @@ func (c *Consumer) Listen(ch *amqp.Channel) {
 	var forever chan struct{}
 	go func() {
 		for msg := range msgs {
-			c.handleMessage(msg)
+			c.handleMessage(msg, q)
 		}
 	}()
 	<-forever
 }
 
-func (c *Consumer) handleMessage(msg amqp.Delivery) {
+func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
+	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
+	t := otel.Tracer("rabbitmq")
+	_, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
+	defer span.End()
+
 	o := &domain.Order{}
 	if err := json.Unmarshal(msg.Body, o); err != nil {
 		logrus.Infof("error unmarshal msg.body into domain.order, err = %v", err)
 		_ = msg.Nack(false, false)
 		return
 	}
-	_, err := c.app.Commands.UpdateOrder.Handle(context.Background(), command.UpdateOrder{
+	_, err := c.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
 		Order: o,
 		UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
 			if err := order.IsPaid(); err != nil {
@@ -65,6 +72,8 @@ func (c *Consumer) handleMessage(msg amqp.Delivery) {
 		// TODO: retry
 		return
 	}
+
+	span.AddEvent("order.updated")
 	_ = msg.Ack(false)
 	logrus.Info("order consume paid event success!")
 }
```

### 9.2 它和 payment consumer 的关系

这段其实就是上一节 payment consumer 的镜像版本。

payment consumer 负责：
- 消费 `order.created`
- 继续 payment 链

order consumer 负责：
- 消费 `order.paid`
- 把支付成功状态回写到 order 领域里

它们都遵循同一个 tracing 模式：
1. 从 `msg.Headers` extract 上下文
2. 手动起 consume span
3. 用这个 `ctx` 调真正业务命令
4. 成功后给 span 加事件

### 9.3 关键代码和详细注释

```go
func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
    // 先从 header 里恢复出 trace 上下文。
    ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)

    t := otel.Tracer("rabbitmq")
    _, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
    defer span.End()

    o := &domain.Order{}
    if err := json.Unmarshal(msg.Body, o); err != nil {
        logrus.Infof("error unmarshal msg.body into domain.order, err = %v", err)
        _ = msg.Nack(false, false)
        return
    }

    // 关键点：不再用 context.Background()。
    // 否则 order 这边会和 payment webhook 那条链断开。
    _, err := c.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
        Order: o,
        UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
            if err := order.IsPaid(); err != nil {
                return nil, err
            }
            return order, nil
        },
    })
    if err != nil {
        return
    }

    span.AddEvent("order.updated")
    _ = msg.Ack(false)
}
```

### 9.4 为什么这里把 `q` 也传给 `handleMessage`

旧代码 `handleMessage(msg)` 已经够处理业务了。
新代码变成 `handleMessage(msg, q)`，主要是为了 tracing：
- 需要 `q.Name`
- 方便 span 名包含真实队列名

这类改动常见于“为了观测性改函数签名”。
业务没变，但调试体验会明显变好。

## 10. stock gRPC 入口补 span：`internal/stock/ports/grpc.go`

### 10.1 这个文件自己的原始 diff

```diff
diff --git a/internal/stock/ports/grpc.go b/internal/stock/ports/grpc.go
index 91e34e8..53b5eaf 100644
--- a/internal/stock/ports/grpc.go
+++ b/internal/stock/ports/grpc.go
@@ -4,6 +4,7 @@ import (
 	context "context"
 
 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
+	"github.com/ghost-yu/go_shop_second/common/tracing"
 	"github.com/ghost-yu/go_shop_second/stock/app"
 	"github.com/ghost-yu/go_shop_second/stock/app/query"
 )
@@ -17,6 +18,9 @@ func NewGRPCServer(app app.Application) *GRPCServer {
 }
 
 func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
+	_, span := tracing.Start(ctx, "GetItems")
+	defer span.End()
+
 	items, err := G.app.Queries.GetItems.Handle(ctx, query.GetItems{ItemIDs: request.ItemIDs})
@@ -25,6 +29,9 @@ func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsReque
 }
 
 func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
+	_, span := tracing.Start(ctx, "CheckIfItemsInStock")
+	defer span.End()
+
 	items, err := G.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{Items: request.Items})
```

### 10.2 这段和 MQ tracing 有什么关系

它不是 MQ 代码，但属于同一个 “more otel” 主题。

因为这组不是只补 MQ，还在补整条链路里那些原来没有显式业务 span 的地方。

#### 自动 gRPC tracing 和这里手动 span 的区别

自动 gRPC tracing：
- 负责“有一次 gRPC 请求进入了服务”
- 更偏协议层

这里手动 span：
- 负责“当前这个业务方法是 `GetItems` / `CheckIfItemsInStock`”
- 更偏应用语义层

所以加上之后，trace 会更好读。

### 10.3 关键代码和详细注释

```go
func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
    // 给 stock 的业务入口补一个更语义化的 span。
    _, span := tracing.Start(ctx, "CheckIfItemsInStock")
    defer span.End()

    items, err := G.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{Items: request.Items})
    if err != nil {
        return nil, err
    }
    return &stockpb.CheckIfItemsInStockResponse{...}, nil
}
```

## 11. 这组最关键的基础设施认知

### 11.1 `context` 为什么可以跨服务传播

严格说，不是 `context` 对象本身跨服务传过去了。

真正发生的是：
1. 本地 `context` 里有 trace 上下文
2. `Inject` 把它编码到 headers 里
3. headers 跟着消息/请求发出去
4. 对端 `Extract` 再从 headers 解码回新的 `context`

所以跨边界传播的是“上下文字段”，不是内存里的那个 Go 对象。

### 11.2 为什么 MQ 必须手动处理，HTTP/gRPC 看起来却比较自动

因为这个项目里：
- HTTP 已经用了 `otelgin`
- gRPC 已经用了 `otelgrpc`

这些库已经替你做了：
- 自动起 span
- 自动读写 header/metadata

但 RabbitMQ 这里没有同样直接接好的自动化封装，所以要自己做：
- carrier
- inject/extract
- publish span
- consume span

### 11.3 `span.AddEvent(...)` 是干什么的

它不是新 span。
它是在当前 span 里追加一个“事件标记”。

比如：
- `payment.created`
- `order.updated`

这类 event 能帮助你在时间线上看到：
- 这个 span 执行期间发生了什么关键业务节点

适合做“轻量标记”。
不适合拿来替代真正的子 span。

### 11.4 `context.Background()` 和 `context.TODO()` 为什么在这里不能乱用

- `context.Background()`：通常表示“全新的根上下文”
- `context.TODO()`：表示“这里以后再想上下文怎么传”

在 tracing 场景里，这两者如果乱用，常常就意味着：

`你主动把原链路切断了。`

所以只要你是从上游拿到了一个有意义的 `ctx`，就尽量继续传，不要随便换成新的根 context。

## 12. 这组为什么要这么做，而不是别的做法

### 12.1 为什么不只依赖日志

因为日志只能告诉你“这里打印过什么”。
但当链路跨：
- HTTP
- gRPC
- MQ
- 第三方 webhook

你会遇到这些问题：
- 同一个订单的日志散落在不同服务
- 先后顺序不好拼
- 并发时容易串台
- 只能靠人工搜索关键词

trace 的优势是：
- 它天然按请求链组织
- 可以看到父子关系
- 可以看到每段耗时
- 可以跳服务查看

### 12.2 为什么不用消息体传 trace，而要用 header

因为 trace 上下文属于“传输元信息”，不是业务数据。

把它放 header 的好处：
- 不污染业务 body 结构
- 业务消费者不需要把 trace 字段建模到业务对象里
- 更符合传播语义

这和 HTTP 把认证、trace、用户代理放 header 是同一类思想。

### 12.3 为什么 span 名里带队列名

因为以后系统复杂起来后，可能会有：
- 多个 exchange
- 多个 queue
- 多个 consumer

如果 span 名只是 `rabbitmq.publish` / `rabbitmq.consume`，信息太少。
带上队列名后，你在 Jaeger 里一眼就知道这段操作针对谁。

## 13. 这组还没完全解决的地方

虽然这组已经把 MQ trace 打通了第一版，但还不是终点。

### 13.1 `RabbitMQHeaderCarrier.Get` 里强转 string 有风险

```go
return value.(string)
```

如果未来 header 里出现非字符串值，这里会 panic。
更稳的写法会做类型判断。

### 13.2 发布失败/消费失败时 span 属性还不够丰富

现在主要是：
- 起 span
- 加 event

但还没系统性地加：
- error status
- queue name attributes
- message id attributes
- retry count attributes

这在真实生产观测里通常会继续补。

### 13.3 业务语义的 span 层级还可以更清晰

比如：
- “publish order.created”
- “consume order.created”
- “create stripe checkout session”
- “update order via grpc”

现在已经开始有这个方向，但还不算特别完整。

## 14. 这一组你最该记住的结论

1. `HTTP/gRPC` 能自动传 trace，是因为现成中间件/handler 帮你做了 inject/extract。
2. `RabbitMQ` 这组课里没有自动方案，所以要手写 `carrier + inject + extract`。
3. `Extract` 只是恢复父上下文，不会自动替你创建 consume span。
4. `Start` 之后一定要继续用新的 `ctx`，否则 trace 很容易断。
5. `context.Background()` / `context.TODO()` 在 tracing 链里乱用，通常就是把链路切断。
6. 这一组真正补上的不是“业务功能”，而是“异步消息链路的可观测性”。

## 15. 你现在应该怎么复习这组

建议你按这个顺序回看：

1. 先重读 [internal/common/broker/rabbitmq.go](/g:/shi/go_shop_second/internal/common/broker/rabbitmq.go)
   - 只盯住 `carrier / inject / extract`
2. 再看 [internal/order/app/command/create_order.go](/g:/shi/go_shop_second/internal/order/app/command/create_order.go)
   - 理解生产者怎么把 trace 塞进 header
3. 再看 [internal/payment/infrastructure/consumer/consumer.go](/g:/shi/go_shop_second/internal/payment/infrastructure/consumer/consumer.go)
   - 理解消费者怎么把 trace 取回来
4. 然后看 [internal/payment/http.go](/g:/shi/go_shop_second/internal/payment/http.go) 和 [internal/order/infrastructure/consumer/consumer.go](/g:/shi/go_shop_second/internal/order/infrastructure/consumer/consumer.go)
   - 理解第二条 MQ 链也采用同样模式
5. 最后打开 Jaeger 去对应想象：
   - publish span 在哪
   - consume span 在哪
   - Stripe span 在哪
   - order update gRPC span 在哪

如果你下一步继续，我就按同样标准写 `lesson26 -> lesson27`。那一组大概率会继续围绕 MQ 重试、死信队列或观测性收尾去推进。
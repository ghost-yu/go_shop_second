# `lesson20 -> lesson21` 独立讲义（详细注释版）

这一组是支付链里的一个关键分水岭。

上一组 `lesson19 -> lesson20` 已经做到：
- 订单创建成功后，页面能轮询订单状态
- 如果订单进入 `waiting_for_payment`，页面能拿到 Stripe 支付链接并跳过去

但这里还有一个根本问题没有解决：

`用户在 Stripe 页面完成支付之后，我们自己的系统怎么知道这笔钱真的付成功了？`

这就是这组要解决的核心问题。

所以这组你不要只看成“加了个 webhook 接口”。

它真正推进的是：

1. 让 Stripe 在支付完成后主动回调我们的 `payment` 服务
2. 让 `payment` 服务验证这次回调真的是 Stripe 发来的，不是别人伪造的
3. 让 `payment` 服务把“订单已支付”这件事重新发成 MQ 事件，交给系统后面继续处理

也就是说，这组完成的是：

`从“我们把支付链接发出去”推进到“外部支付平台把支付结果通知回来”。`

这是一条方向完全相反的调用链：
- 以前是我们调 Stripe
- 现在是 Stripe 调我们

这也是为什么这组对新手来说更容易绕。

## 1. 原始 diff 去哪里看

原始差异在这里：
[diff.md](/g:/shi/go_shop_second/ai_diffs/lesson_pairs/lesson20-to-lesson21/diff.md)

这次我还是按固定规则来：
- 每个文件先贴自己的 diff
- 再贴带中文注释的代码/关键代码
- 再讲旧代码做什么、新代码做什么、为什么这样改

## 2. 正确阅读顺序

这组建议你按下面顺序读：

1. [internal/payment/http.go](/g:/shi/go_shop_second/internal/payment/http.go)
2. [internal/payment/main.go](/g:/shi/go_shop_second/internal/payment/main.go)
3. [internal/payment/infrastructure/processor/stripe.go](/g:/shi/go_shop_second/internal/payment/infrastructure/processor/stripe.go)
4. [internal/common/config/global.yaml](/g:/shi/go_shop_second/internal/common/config/global.yaml)
5. [internal/common/config/viper.go](/g:/shi/go_shop_second/internal/common/config/viper.go)
6. [internal/payment/domain/payment.go](/g:/shi/go_shop_second/internal/payment/domain/payment.go)

原因：
- `payment/http.go` 是真正的 webhook 处理入口，也是这组主角
- `payment/main.go` 解释 MQ channel 为什么要注入进 handler
- `stripe.go` 解释为什么 success URL 和 metadata 要一起补
- 配置文件解释 endpoint secret 从哪来
- `domain/payment.go` 解释为什么要补一个 payment 自己的 Order 结构

## 3. 总调用链

这组最关键的脑图你先记住：

```text
用户点开 Stripe 支付页面
-> 用户在 Stripe 完成支付
-> Stripe 服务器主动 POST /api/webhook 到 payment 服务
-> payment/http.go 读取请求体并校验签名
-> webhook 事件被解析成 Stripe event
-> 如果是 checkout.session.completed 且 payment_status=paid
-> 从 metadata 里还原订单上下文
-> payment 服务把 order.paid 事件重新发到 RabbitMQ
-> 后面的 order 服务再去消费这个 paid 事件
```

这里有一个你必须先转过来的思维：

`success 页面跳转不等于支付成功的最终确认。真正可靠的确认来源，是 Stripe 发来的 webhook。`

也就是说：
- 前端页面跳转更像“用户体验”
- webhook 才是“系统事实”

这组就是在补这一层系统事实。

## 4. 关键文件详细讲解

### 4.1 [internal/payment/http.go](/g:/shi/go_shop_second/internal/payment/http.go)

这是这组最重要的文件，没有之一。

旧版本里，这个文件几乎什么都没做，只是简单打一句日志：
- 收到了 Stripe webhook
- 但并没有真的验证、解析、处理它

新版本开始，它正式承担 webhook 消费入口的职责。

先看这个文件自己的 diff。

```diff
diff --git a/internal/payment/http.go b/internal/payment/http.go
index b98bbe3..9017d2f 100644
--- a/internal/payment/http.go
+++ b/internal/payment/http.go
@@ -1,15 +1,28 @@
 package main
 
 import (
+	"context"
+	"encoding/json"
+	"io"
+	"net/http"
+
+	"github.com/ghost-yu/go_shop_second/common/broker"
+	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/payment/domain"
 	"github.com/gin-gonic/gin"
+	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
+	"github.com/spf13/viper"
+	"github.com/stripe/stripe-go/v79"
+	"github.com/stripe/stripe-go/v79/webhook"
 )
 
 type PaymentHandler struct {
+	channel *amqp.Channel
 }
 
-func NewPaymentHandler() *PaymentHandler {
-	return &PaymentHandler{}
+func NewPaymentHandler(ch *amqp.Channel) *PaymentHandler {
+	return &PaymentHandler{channel: ch}
 }
@@
 func (h *PaymentHandler) handleWebhook(c *gin.Context) {
 	logrus.Info("receive webhook from stripe")
+	const MaxBodyBytes = int64(65536)
+	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
+	payload, err := io.ReadAll(c.Request.Body)
+	if err != nil {
+		...
+	}
+
+	event, err := webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"),
+		viper.GetString("ENDPOINT_STRIPE_SECRET"))
+
+	if err != nil {
+		...
+	}
+
+	switch event.Type {
+	case stripe.EventTypeCheckoutSessionCompleted:
+		var session stripe.CheckoutSession
+		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
+			...
+		}
+
+		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
+			...
+			_ = h.channel.PublishWithContext(ctx, broker.EventOrderPaid, "", false, false, amqp.Publishing{...})
+		}
+	}
+	c.JSON(http.StatusOK, nil)
 }
```

这个 diff 很大，但它可以拆成 5 个动作。

#### 动作 1：把 webhook handler 从“空壳”变成真正的处理入口

下面我先贴当前代码里这份文件最关键的部分，并直接写上中文注释。

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"

    "github.com/ghost-yu/go_shop_second/common/broker"
    "github.com/ghost-yu/go_shop_second/common/consts"
    "github.com/ghost-yu/go_shop_second/common/entity"
    "github.com/ghost-yu/go_shop_second/common/logging"
    "github.com/gin-gonic/gin"
    "github.com/pkg/errors"
    amqp "github.com/rabbitmq/amqp091-go"
    "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
    "github.com/stripe/stripe-go/v79"
    "github.com/stripe/stripe-go/v79/webhook"
    "go.opentelemetry.io/otel"
)

type PaymentHandler struct {
    // 这次最关键的新增字段：handler 里要持有 RabbitMQ channel。
    // 因为 webhook 收到“支付成功”后，它要立刻把 order.paid 事件发回 MQ。
    channel *amqp.Channel
}

func NewPaymentHandler(ch *amqp.Channel) *PaymentHandler {
    // 旧代码是 NewPaymentHandler()，什么也不带。
    // 新代码必须把 channel 注入进来，因为 handler 已经不只是“接请求”，
    // 它还要负责对外发消息了。
    return &PaymentHandler{channel: ch}
}

// stripe listen --forward-to localhost:8284/api/webhook
func (h *PaymentHandler) RegisterRoutes(c *gin.Engine) {
    // 这条路由就是 Stripe webhook 的入口。
    // Stripe 后面会主动 POST 到 /api/webhook。
    c.POST("/api/webhook", h.handleWebhook)
}
```

这里你首先要理解的是：

旧代码的 `PaymentHandler` 只是一个 HTTP 壳。

新代码的 `PaymentHandler` 已经开始承担“支付结果入口 + 事件发布者”双重角色。

也就是说：
- 它不是单纯收一个请求然后返回 200
- 它收完请求还要把消息发进 RabbitMQ

这也是为什么 `channel *amqp.Channel` 会加进来。

#### 这里补一个新手容易断掉的点：`*amqp.Channel` 是什么

`amqp091-go` 是 Go 里操作 RabbitMQ 的库。

里面的 `Channel` 你可以先粗暴理解成：

`跟 RabbitMQ 交互的一条工作通道。`

你发消息、声明队列、消费消息，通常都要通过这个 channel 来做。

所以这里把 channel 注入进 `PaymentHandler`，意思就是：

`当 webhook 确认支付成功后，这个 HTTP handler 自己就有能力把消息继续发到 MQ。`

---

#### 动作 2：读取 webhook 请求体，并限制 body 大小

```go
func (h *PaymentHandler) handleWebhook(c *gin.Context) {
    logrus.WithContext(c.Request.Context()).Info("receive webhook from stripe")
    var err error
    defer func() {
        if err != nil {
            logging.Warnf(c.Request.Context(), nil, "handleWebhook err=%v", err)
        } else {
            logging.Infof(c.Request.Context(), nil, "%s", "handleWebhook success")
        }
    }()

    // 先限制 body 大小，防止请求体无限大。
    const MaxBodyBytes = int64(65536)
    c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

    // 再把整个 body 读出来。
    payload, err := io.ReadAll(c.Request.Body)
    if err != nil {
        err = errors.Wrap(err, "Error reading request body")
        c.JSON(http.StatusServiceUnavailable, err.Error())
        return
    }
```

这段看起来像样板代码，但不能跳过。

它在做两件事：

1. 限制 webhook 请求体大小
2. 把原始 body 全部读出来，后面拿去做签名校验

#### 为什么要保留原始 body

因为 Stripe webhook 的签名校验不是“看 JSON 解析后的对象”来做的，而是：

`拿原始 payload + 请求头里的 Stripe-Signature + endpoint secret 一起算。`

也就是说，这里必须先把原始字节流读出来。

如果你先把它改得面目全非，再去校验，签名就不一定对得上了。

#### `http.MaxBytesReader` 是干什么的

这是 Go 标准库用来限制请求体大小的一个工具。

你可以把它理解成：

`最多只允许这个请求读 64KB，再大就不接受。`

它的意义主要是安全和稳健性：
- 避免别人随便发一个特别大的 body 来搞你
- 避免无上限读取请求体

对小白来说，你先记住：

`凡是接外部平台回调，限制 body 大小是个很正常的防守动作。`

---

#### 动作 3：校验 webhook 真的是 Stripe 发来的

```go
    event, err := webhook.ConstructEvent(
        payload,
        c.Request.Header.Get("Stripe-Signature"),
        viper.GetString("ENDPOINT_STRIPE_SECRET"),
    )

    if err != nil {
        err = errors.Wrap(err, "error verifying webhook signature")
        c.JSON(http.StatusBadRequest, err.Error())
        return
    }
```

这一段是整条 webhook 链里最关键的安全点。

为什么？

因为：
- 任何人理论上都能给你的 `/api/webhook` 发 HTTP 请求
- 你不能因为别人 POST 了一段 JSON，就相信“这一定是 Stripe”

所以这里要用 Stripe 官方 SDK 提供的 `webhook.ConstructEvent(...)` 来验签。

它会用三个东西：

1. `payload`
   - 原始请求体
2. `Stripe-Signature`
   - Stripe 放在请求头里的签名
3. `ENDPOINT_STRIPE_SECRET`
   - 你自己在 Stripe CLI 或 dashboard 里配置好的 secret

只有三者能对上，才说明：

`这次 webhook 是 Stripe 发来的，而且中途没被篡改。`

#### 这里补一个云服务/外部平台基础知识：`endpoint secret` 是什么

你可以把它理解成：

`Stripe 和你这个 webhook 接口之间共享的一把暗号。`

Stripe 发请求时会带签名。
你这边拿这个 secret 去验。

如果验不过：
- 说明要么请求不是 Stripe 发的
- 要么 secret 配错了
- 要么 body 被改了

这就是为什么这组要专门加 `endpoint-stripe-secret` 配置。

---

#### 动作 4：只处理我们真正关心的事件类型

```go
    switch event.Type {
    case stripe.EventTypeCheckoutSessionCompleted:
        var session stripe.CheckoutSession
        if err = json.Unmarshal(event.Data.Raw, &session); err != nil {
            err = errors.Wrap(err, "error unmarshal event.data.raw into session")
            c.JSON(http.StatusBadRequest, err.Error())
            return
        }
```

这一步的意思是：

Stripe 会发很多种事件，但当前这组课程只关心一种：

`checkout.session.completed`

为什么只关心它？

因为前一组我们创建的就是 Checkout Session。

所以现在最自然的确认点就是：
- 当这个 session 完成了
- 我们再继续判断它是不是已经 paid

这里还有一个重要动作：

```go
json.Unmarshal(event.Data.Raw, &session)
```

意思是把 Stripe event 里面原始的 `data` 字段，解析成 SDK 里的 `stripe.CheckoutSession` 结构体。

这一步之后，你才能用：
- `session.PaymentStatus`
- `session.Metadata`
- `session.ID`

这些字段。

#### 这里小白最容易绕的点

你会看到这里不是直接 `json.Unmarshal(payload, &session)`。

因为：
- 整个 webhook payload 是一个“事件对象”
- 真正的 checkout session 数据包在 `event.Data.Raw` 里面

也就是说：
- 先验签，拿到 `event`
- 再从 `event` 里拆出真正的 session

这是两层结构，不是一层。

---

#### 动作 5：如果确认已支付，就发 `order.paid` 事件

```go
        if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
            var items []*entity.Item
            _ = json.Unmarshal([]byte(session.Metadata["items"]), &items)

            tr := otel.Tracer("rabbitmq")
            ctx, span := tr.Start(c.Request.Context(), fmt.Sprintf("rabbitmq.%s.publish", broker.EventOrderPaid))
            defer span.End()

            _ = broker.PublishEvent(ctx, broker.PublishEventReq{
                Channel:  h.channel,
                Routing:  broker.FanOut,
                Queue:    "",
                Exchange: broker.EventOrderPaid,
                Body: entity.NewOrder(
                    session.Metadata["orderID"],
                    session.Metadata["customerID"],
                    consts.OrderStatusPaid,
                    session.Metadata["paymentLink"],
                    items,
                ),
            })
        }
    }
    c.JSON(http.StatusOK, nil)
}
```

这一段是整条业务逻辑真正落地的地方。

它做了 4 件事：

1. 确认 `PaymentStatus == paid`
2. 从 `session.Metadata` 里把订单上下文拿回来
3. 重建一个订单对象
4. 发布 `order.paid` 事件到 RabbitMQ

这里最关键的一句话是：

`payment 服务并没有在 webhook 里直接改 order 服务数据库，而是继续走“发事件”的异步模式。`

这和前面整个项目的设计主线是一致的。

#### 为什么还要靠 `metadata`

因为 Stripe 只知道它自己的 session。

如果 webhook 回来时你想知道：
- 这是哪张订单
- 这是哪个用户
- 当时的 items 是什么

你就必须在创建 Checkout Session 时，提前把这些信息塞进 `metadata`。

所以这组和上一组是紧密配合的：
- 上一组在创建支付链接时把 `orderID/customerID/items/paymentLink` 放进 metadata
- 这一组 webhook 回来时再把这些值拿出来

这就是一次非常典型的“外部平台上下文透传”。

#### 这里为什么发 MQ，而不是直接调 order gRPC

这是个设计问题，也是你以后面试里很容易被问的问题。

当前选择发 MQ 的好处是：
- payment 服务不用同步等 order 服务响应
- webhook 处理速度更快
- 后续可以多个消费者都感知到 `order.paid`
- 事件驱动风格更统一

坏处是：
- 时序更复杂
- 你得处理异步一致性
- 如果后面没有消费者绑定，这条消息就只是进了交换机，不会真的落地处理

而字幕最后也正好指出了当前系统的缺口：

`payment 已经把 order.paid 发出来了，但 order 这边还没开始消费它。`

这正是下一组要补的内容。

#### 这里当前代码还有哪些不完美的地方

1. `json.Unmarshal([]byte(session.Metadata["items"]), &items)` 错误被忽略了
2. `PublishEvent(...)` 返回值也被忽略了
3. `c.JSON(http.StatusOK, nil)` 很简陋，只能说明“我收到了 webhook”，不能表达更多业务结果
4. webhook handler 里逻辑已经有点厚了，后面更成熟时通常会抽 service/usecase

也就是说，这组是“链路打通版”，不是“最终架构版”。

### 这一文件你最该带走的结论

`payment/http.go` 这组真正完成的是：把 Stripe 的外部支付结果，正式接入到我们自己的系统事件流里。`

---

### 4.2 [internal/payment/main.go](/g:/shi/go_shop_second/internal/payment/main.go)

这个文件改动不多，但它解释了 webhook 为什么现在能发 MQ。

先看它自己的 diff。

```diff
diff --git a/internal/payment/main.go b/internal/payment/main.go
index 55b131a..b80b31c 100644
--- a/internal/payment/main.go
+++ b/internal/payment/main.go
@@ -42,7 +42,7 @@ func main() {
 
 	go consumer.NewConsumer(application).Listen(ch)
 
-	paymentHandler := NewPaymentHandler()
+	paymentHandler := NewPaymentHandler(ch)
 	switch serverType {
 	case "http":
 		server.RunHTTPServer(viper.GetString("payment.service-name"), paymentHandler.RegisterRoutes)
```

再看当前代码对应片段：

```go
func main() {
    serviceName := viper.GetString("payment.service-name")
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    application, cleanup := service.NewApplication(ctx)
    defer cleanup()

    // 先建立到 RabbitMQ 的连接，拿到 channel。
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

    // 这边还是原来那条链：payment 自己也会消费 order.created。
    go consumer.NewConsumer(application).Listen(ch)

    // 这次最重要的变化：把 MQ channel 注入给 PaymentHandler。
    // 这样 webhook 收到支付成功后，也能继续往 MQ 发 order.paid 事件。
    paymentHandler := NewPaymentHandler(ch)

    switch serverType {
    case "http":
        server.RunHTTPServer(viper.GetString("payment.service-name"), paymentHandler.RegisterRoutes)
    }
}
```

旧代码做的事情是：
- `paymentHandler` 只是个纯 HTTP handler
- 不需要 MQ channel

新代码做的事情是：
- `paymentHandler` 已经变成“HTTP 入口 + MQ 发布者”
- 所以启动时必须把 `ch` 注进去

这就是典型的依赖注入：

`谁需要什么能力，就在构造它的时候把依赖传进去。`

对 Go 初学者来说，这里不用把“依赖注入”想得太框架化。

你先把它理解成一句简单的话：

`PaymentHandler 想发 MQ，就得先拿到 MQ 的 channel。`

---

### 4.3 [internal/payment/infrastructure/processor/stripe.go](/g:/shi/go_shop_second/internal/payment/infrastructure/processor/stripe.go)

这一组里这个文件不是主角，但它是 webhook 能顺利配合起来的重要补丁。

先看这个文件自己的 diff。

```diff
diff --git a/internal/payment/infrastructure/processor/stripe.go b/internal/payment/infrastructure/processor/stripe.go
index 7111272..37fcb1d 100644
--- a/internal/payment/infrastructure/processor/stripe.go
+++ b/internal/payment/infrastructure/processor/stripe.go
@@ -3,6 +3,7 @@ package processor
 import (
 	"context"
 	"encoding/json"
+	"fmt"
@@ -21,7 +22,7 @@ func NewStripeProcessor(apiKey string) *StripeProcessor {
 	return &StripeProcessor{apiKey: apiKey}
 }
 
-var (
+const (
 	successURL = "http://localhost:8282/success"
 )
@@ -36,16 +37,17 @@ func (s StripeProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.O
 
 	marshalledItems, _ := json.Marshal(order.Items)
 	metadata := map[string]string{
-		"orderID":    order.ID,
-		"customerID": order.CustomerID,
-		"status":     order.Status,
-		"items":      string(marshalledItems),
+		"orderID":     order.ID,
+		"customerID":  order.CustomerID,
+		"status":      order.Status,
+		"items":       string(marshalledItems),
+		"paymentLink": order.PaymentLink,
 	}
 	params := &stripe.CheckoutSessionParams{
 		Metadata:   metadata,
 		LineItems:  items,
 		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
-		SuccessURL: stripe.String(successURL),
+		SuccessURL: stripe.String(fmt.Sprintf("%s?customerID=%s&orderID=%s", successURL, order.CustomerID, order.ID)),
 	}
```

这组真正重要的变化有两个。

#### 变化 1：`SuccessURL` 补上了查询参数

当前代码相关部分：

```go
const (
    successURL = "http://localhost:8282/success"
)

func (s StripeProcessor) CreatePaymentLink(ctx context.Context, order *entity.Order) (string, error) {
    ...
    params := &stripe.CheckoutSessionParams{
        Metadata:   metadata,
        LineItems:  items,
        Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),

        // 这次关键改动：支付成功后跳回 success 页面时，
        // 把 customerID 和 orderID 一起带回去。
        SuccessURL: stripe.String(fmt.Sprintf(
            "%s?customerID=%s&orderID=%s",
            successURL,
            order.CustomerID,
            order.ID,
        )),
    }
```

为什么这行很重要？

因为上一组页面 `success.html` 是靠 URL 参数来知道自己查哪张订单的。

如果 Stripe 支付成功后只跳回：

```text
/success
```

而不带：
- `customerID`
- `orderID`

那页面就不知道自己该轮询谁。

所以这一组是在把“支付成功回跳页面”跟“前端页面状态查询”真正接起来。

#### 变化 2：metadata 里补了 `paymentLink`

```go
metadata := map[string]string{
    "orderID":     order.ID,
    "customerID":  order.CustomerID,
    "status":      order.Status,
    "items":       string(marshalledItems),
    "paymentLink": order.PaymentLink,
}
```

为什么要多带一个 `paymentLink`？

因为 webhook 回来以后，payment 服务会根据 metadata 重建订单上下文，再发 `order.paid` 事件。

如果你希望那条事件里把订单数据尽量带完整，那就要把之前的 `paymentLink` 也传回来。

这就是典型的“前面创建 session 时多存一点上下文，后面 webhook 回来时就少一点重建痛苦”。

#### 这一文件你最该带走的结论

`stripe.go` 这组不是主逻辑入口，但它给 webhook 链补了两个关键上下文：回跳页面参数，以及回调后继续发事件要用到的 metadata。`

---

### 4.4 [internal/common/config/global.yaml](/g:/shi/go_shop_second/internal/common/config/global.yaml)

先看 diff。

```diff
diff --git a/internal/common/config/global.yaml b/internal/common/config/global.yaml
index 17a4f41..7dd93c9 100644
--- a/internal/common/config/global.yaml
+++ b/internal/common/config/global.yaml
@@ -27,4 +27,5 @@ rabbitmq:
   host: 127.0.0.1
   port: 5672
 
-stripe-key: "${STRIPE_KEY}"
+stripe-key: "${STRIPE_KEY}"
+endpoint-stripe-secret: "${ENDPOINT_STRIPE_SECRET}"
```

当前代码对应部分：

```yaml
stripe-key: "${STRIPE_KEY}"
endpoint-stripe-secret: "${ENDPOINT_STRIPE_SECRET}"
```

这一组配置文件的关键点不是 YAML 语法，而是：

`系统第一次同时依赖两种 Stripe 侧敏感配置。`

1. `stripe-key`
   - 我们主动调用 Stripe API 创建支付链接时用

2. `endpoint-stripe-secret`
   - Stripe 主动回调我们 webhook 时验签用

这两个东西不是一回事。

新手很容易混：
- 以为都是 Stripe 的 secret，所以随便一个都能通用

其实不行。

你可以这样记：
- `stripe-key`：我们调用 Stripe 的“客户端凭证”
- `endpoint-stripe-secret`：Stripe 回调我们时，用来验证身份的“共享暗号”

这两个用途完全不同。

---

### 4.5 [internal/common/config/viper.go](/g:/shi/go_shop_second/internal/common/config/viper.go)

先看 diff。

```diff
diff --git a/internal/common/config/viper.go b/internal/common/config/viper.go
index 0246761..a0c91d2 100644
--- a/internal/common/config/viper.go
+++ b/internal/common/config/viper.go
@@ -12,6 +12,6 @@ func NewViperConfig() error {
 	viper.AddConfigPath("../common/config")
 	viper.EnvKeyReplacer(strings.NewReplacer("_", "-"))
 	viper.AutomaticEnv()
-	_ = viper.BindEnv("stripe-key", "STRIPE_KEY")
+	_ = viper.BindEnv("stripe-key", "STRIPE_KEY", "endpoint-stripe-secret", "ENDPOINT_STRIPE_SECRET")
 	return viper.ReadInConfig()
 }
```

当前代码里这一段已经进一步演进，但核心意思还是一样：

```go
viper.EnvKeyReplacer(strings.NewReplacer("_", "-"))
viper.AutomaticEnv()
_ = viper.BindEnv("stripe-key", "STRIPE_KEY", "endpoint-stripe-secret", "ENDPOINT_STRIPE_SECRET")
```

这段代码在做什么？

就是告诉 Viper：

- 配置项 `stripe-key` 对应环境变量 `STRIPE_KEY`
- 配置项 `endpoint-stripe-secret` 对应环境变量 `ENDPOINT_STRIPE_SECRET`

#### 为什么这里要显式绑定

因为 webhook 这条链很依赖 secret 配对正确。

如果 secret 没读出来：
- 代码不一定立刻崩
- 但 webhook 验签一定会失败
- 你就会看到“怎么老提示签名不对”

所以这组配置绑定虽然改动小，但它是 webhook 能跑起来的必要条件。

#### 这里小白最容易忘的点

配置里叫：
- `endpoint-stripe-secret`

环境变量里叫：
- `ENDPOINT_STRIPE_SECRET`

中间这个名字映射不是 Go 自动魔法，而是 Viper 帮你做的。

---

### 4.6 [internal/payment/domain/payment.go](/g:/shi/go_shop_second/internal/payment/domain/payment.go)

先看 diff。

```diff
diff --git a/internal/payment/domain/payment.go b/internal/payment/domain/payment.go
index 503bf26..07d56aa 100644
--- a/internal/payment/domain/payment.go
+++ b/internal/payment/domain/payment.go
@@ -9,3 +9,11 @@ import (
 type Processor interface {
 	CreatePaymentLink(context.Context, *orderpb.Order) (string, error)
 }
+
+type Order struct {
+	ID          string
+	CustomerID  string
+	Status      string
+	PaymentLink string
+	Items       []*orderpb.Item
+}
```

当前代码已经进一步演进到 `common/entity.Order`，但你还是要理解当时为什么加这份结构。

这里的核心原因是：

`payment 服务在 webhook 里需要自己表达一份“订单数据”，但它又不应该直接依赖 order 服务的内部 domain。`

所以当时作者选择在 payment 自己的 domain 里加一份 `Order`：

```go
type Order struct {
    ID          string
    CustomerID  string
    Status      string
    PaymentLink string
    Items       []*orderpb.Item
}
```

虽然这份结构后面还会继续演进，但这一步体现了一个很重要的微服务意识：

`不同服务哪怕都在说“订单”，也不一定共用同一个内部模型。`

这也是为什么字幕里作者会说“你不能直接 import order 那边的东西，因为虽然写在一个仓库里，本质上还是分开的服务”。

## 5. 第三方库和易错点

### Stripe webhook SDK

这组最关键的第三方能力是：

```go
webhook.ConstructEvent(payload, signature, secret)
```

你要记住它的作用不是“解析 JSON”，而是：

`验证这次 webhook 是否可信，并顺便构造 Stripe event 对象。`

新手最容易忽略：
- `payload` 要用原始 body
- 不能随便改动 body 后再验签
- `endpoint secret` 和 `stripe-key` 不是一个东西

### RabbitMQ / amqp.Channel

这组 `payment` handler 之所以要拿 `*amqp.Channel`，就是因为 webhook 一旦确认支付成功，它要把 `order.paid` 发布出去。

这里你要理解：
- webhook 是 HTTP 入口
- MQ 是系统内部异步传播机制

这组代码把两个世界接起来了。

### Viper

这组 Viper 的坑在于配置名和环境变量名不完全一样。

如果你没配好 `ENDPOINT_STRIPE_SECRET`，你会看到 webhook 验签失败，但不一定第一时间想到是配置问题。

## 6. 为什么这么设计

这组设计选择的是：

`支付成功后，不直接在 webhook 里同步改 order 服务，而是仍然通过 MQ 继续传播一个 order.paid 事件。`

这跟整个项目的风格是一致的：
- 服务之间尽量事件驱动
- payment 只负责感知支付成功
- 后面的 order、kitchen 等服务再根据 paid 事件各自处理

好处：
- 解耦
- 扩展性更好
- 多个下游都能感知支付完成

坏处：
- 链路更绕
- 排错更难
- 少一个消费者，消息就没人处理

而这组最后正处在这种“中间态”：

`payment 端已经会发 order.paid，但 order 端还没开始消费它。`

## 7. 当前还不完美的地方

这组你一定要看到这些过渡态问题：

1. webhook 里一些错误仍然被忽略了，比如反序列化 `items`
2. 发布 MQ 的返回值没有认真处理
3. handler 业务逻辑已经开始变厚，后面更适合继续下沉到应用层
4. success 页面虽然会轮询，但此时系统还没把 paid 事件接到 order 里，所以页面状态不会真正推进
5. 当前实现强依赖 metadata 完整性，如果 metadata 漏字段，后续链路就会断

## 8. 这组最该带走的知识点

1. success 页面跳转不等于最终支付成功确认
2. 真正可靠的支付确认，通常来自第三方平台的 webhook
3. webhook 一定要做签名校验，不能只看请求路径和 JSON 内容
4. 创建 Stripe Checkout Session 时，提前放进 metadata 的上下文，后面 webhook 才拿得回来
5. webhook 收到支付成功后，当前项目仍然选择继续发 MQ，而不是直接同步改别的服务
6. 这组最核心的变化方向是“Stripe 调我们”，而不是“我们调 Stripe”

## 9. 一句话收住这组

`lesson20 -> lesson21` 的本质，不是“补了个 webhook 接口”，而是“支付系统第一次能从 Stripe 反向接收成功信号，并把这个外部事实重新转成我们自己系统里的 order.paid 事件”。`
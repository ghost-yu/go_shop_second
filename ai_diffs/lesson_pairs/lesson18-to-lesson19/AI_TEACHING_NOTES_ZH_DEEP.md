# `lesson18 -> lesson19` 独立讲义（重写版，结合字幕）

这一组差异你不要把它看成“接了个 Stripe SDK”。

它真正做的是三件彼此依赖的事：

1. 把 `payment` 服务里的假支付处理器，切换成真正会向 Stripe 创建 Checkout Session 的处理器
2. 修掉 `order` 内存仓储里一个会导致 `payment link` 和 `status` 回写失败的 bug
3. 让 `stock` 在校验库存时顺手把 `PriceID` 带回来，不然 Stripe 根本不知道你到底要卖哪一个价格对象

也就是说，这一组不是单纯“接 SDK”，而是：

`把上一节已经跑通的 payment 主链，从“假的演示版”推进到“第一次能打到真实第三方支付平台的版本”。`

这也是为什么视频字幕里 `19.txt` 先在修 bug，`20.TXT` 又在补 `PriceID`。这两个看起来像小修，但其实都在给 Stripe 真正跑通铺路。

## 1. 先看原始 diff

原始差异文件在这里：
[diff.md](/g:/shi/go_shop_second/ai_diffs/lesson_pairs/lesson18-to-lesson19/diff.md)

为了让你先建立整体印象，我先把这一组的完整 diff 全文放在这里。你先整体扫一遍，知道到底改了哪些文件，再往下看讲解。

```diff
diff --git a/internal/order/adapters/order_inmem_repository.go b/internal/order/adapters/order_inmem_repository.go
index 818b51f..f910521 100644
--- a/internal/order/adapters/order_inmem_repository.go
+++ b/internal/order/adapters/order_inmem_repository.go
@@ -70,7 +70,7 @@ func (m *MemoryOrderRepository) Update(ctx context.Context, order *domain.Order,
 	for i, o := range m.store {
 		if o.ID == order.ID && o.CustomerID == order.CustomerID {
 			found = true
-			updatedOrder, err := updateFn(ctx, o)
+			updatedOrder, err := updateFn(ctx, order)
 			if err != nil {
 				return err
 			}
 diff --git a/internal/payment/go.mod b/internal/payment/go.mod
 index cceef38..7b26a13 100644
--- a/internal/payment/go.mod
+++ b/internal/payment/go.mod
@@ -6,14 +6,15 @@ replace github.com/ghost-yu/go_shop_second/common => ../common
 
 require (
 	github.com/ghost-yu/go_shop_second/common v0.0.0-00010101000000-000000000000
-	github.com/armon/go-metrics v0.4.1
 	github.com/gin-gonic/gin v1.9.1
 	github.com/rabbitmq/amqp091-go v1.10.0
 	github.com/sirupsen/logrus v1.8.1
 	github.com/spf13/viper v1.19.0
+	github.com/stripe/stripe-go/v79 v79.12.0
 )
 
 require (
+	github.com/armon/go-metrics v0.4.1 // indirect
 	github.com/bytedance/sonic v1.10.0-rc3 // indirect
 	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
 	github.com/chenzhuoyu/iasm v0.9.0 // indirect
 diff --git a/internal/payment/go.sum b/internal/payment/go.sum
 index 53fad56..05d1402 100644
--- a/internal/payment/go.sum
+++ b/internal/payment/go.sum
@@ -278,6 +278,8 @@ github.com/stretchr/testify v1.8.2/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o
 github.com/stretchr/testify v1.8.4/go.mod h1:sz/lmYIOXD/1dqDmKjjqLyZ2RngseejIcXlSw2iwfAo=
 github.com/stretchr/testify v1.9.0 h1:HtqpIVDClZ4nwg75+f6Lvsy/wHu+3BoSGCbBAcpTsTg=
 github.com/stretchr/testify v1.9.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8C91i36aY=
+github.com/stripe/stripe-go/v79 v79.12.0 h1:HQs/kxNEB3gYA7FnkSFkp0kSOeez0fsmCWev6SxftYs=
+github.com/stripe/stripe-go/v79 v79.12.0/go.mod h1:cuH6X0zC8peY6f1AubHwgJ/fJSn2dh5pfiCr6CjyKVU=
 github.com/subosito/gotenv v1.6.0 h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=
 github.com/subosito/gotenv v1.6.0/go.mod h1:Dk4QP5c2W3ibzajGcXpNraDfq2IrhjMIvMSWPKKo0FU=
 github.com/tv42/httpunix v0.0.0-20150427012821-b75d8614f926/go.mod h1:9ESjWnEqriFuLhtthL60Sar/7RFoluCcXsuvEwTV5KM=
@@ -331,6 +333,7 @@ golang.org/x/net v0.0.0-20200226121028-0de0cce0169b/go.mod h1:z5CRVTTTmAJ677TzLL
 golang.org/x/net v0.0.0-20201021035429-f5854403a974/go.mod h1:sp8m0HH+o8qH0wwXwYZr8TS3Oi6o0r6Gce1SSxlDquU=
 golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
 golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1/go.mod h1:9tjilg8BloeKEkVJvy7fQ90B1CfIiPueXVOjqfkSzI8=
+golang.org/x/net v0.0.0-20210520170846-37e1c6afe023/go.mod h1:9nx3DQGgdP8bBQD5qxJ1jj9UTztislL4KSBs9R2vV5Y=
 golang.org/x/net v0.28.0 h1:a9JDOJc5GMUJ0+UDqmLT86WiEy7iWyIhz8gz8E4e5hE=
 golang.org/x/net v0.28.0/go.mod h1:yqtgsTWOOnlGLG9GFRrK3++bGOUEkNBoHZc8MEDWPNg=
 golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be/go.mod h1:N/0e6XlmueqKjAGxoOufVs8QHGRruUQn6yWY3a++T0U=
@@ -359,6 +362,7 @@ golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f/go.mod h1:h1NjWce9XRLGQEsW7w
 golang.org/x/sys v0.0.0-20201119102817-f84b799fce68/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210303074136-134d130e1a04/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210330210617-4fbd30eecc44/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20210423082822-04245dca01da/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 diff --git a/internal/payment/infrastructure/processor/stripe.go b/internal/payment/infrastructure/processor/stripe.go
 new file mode 100644
 index 0000000..7111272
--- /dev/null
+++ b/internal/payment/infrastructure/processor/stripe.go
@@ -0,0 +1,55 @@
+package processor
+
+import (
+	"context"
+	"encoding/json"
+
+	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/stripe/stripe-go/v79"
+	"github.com/stripe/stripe-go/v79/checkout/session"
+)
+
+type StripeProcessor struct {
+	apiKey string
+}
+
+func NewStripeProcessor(apiKey string) *StripeProcessor {
+	if apiKey == "" {
+		panic("empty api key")
+	}
+	stripe.Key = apiKey
+	return &StripeProcessor{apiKey: apiKey}
+}
+
+var (
+	successURL = "http://localhost:8282/success"
+)
+
+func (s StripeProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
+	var items []*stripe.CheckoutSessionLineItemParams
+	for _, item := range order.Items {
+		items = append(items, &stripe.CheckoutSessionLineItemParams{
+			Price:    stripe.String(item.PriceID),
+			Quantity: stripe.Int64(int64(item.Quantity)),
+		})
+	}
+
+	marshalledItems, _ := json.Marshal(order.Items)
+	metadata := map[string]string{
+		"orderID":    order.ID,
+		"customerID": order.CustomerID,
+		"status":     order.Status,
+		"items":      string(marshalledItems),
+	}
+	params := &stripe.CheckoutSessionParams{
+		Metadata:   metadata,
+		LineItems:  items,
+		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
+		SuccessURL: stripe.String(successURL),
+	}
+	result, err := session.New(params)
+	if err != nil {
+		return "", err
+	}
+	return result.URL, nil
+}
 diff --git a/internal/payment/service/application.go b/internal/payment/service/application.go
 index b265265..31b7d69 100644
--- a/internal/payment/service/application.go
+++ b/internal/payment/service/application.go
@@ -11,6 +11,7 @@ import (
 	"github.com/ghost-yu/go_shop_second/payment/domain"
 	"github.com/ghost-yu/go_shop_second/payment/infrastructure/processor"
 	"github.com/sirupsen/logrus"
+	"github.com/spf13/viper"
 )
 
 func NewApplication(ctx context.Context) (app.Application, func()) {
@@ -19,8 +20,9 @@ func NewApplication(ctx context.Context) (app.Application, func()) {
 		panic(err)
 	}
 	orderGRPC := adapters.NewOrderGRPC(orderClient)
-	memoryProcessor := processor.NewInmemProcessor()
-	return newApplication(ctx, orderGRPC, memoryProcessor), func() {
+	//memoryProcessor := processor.NewInmemProcessor()
+	stripeProcessor := processor.NewStripeProcessor(viper.GetString("stripe-key"))
+	return newApplication(ctx, orderGRPC, stripeProcessor), func() {
 		_ = closeOrderClient()
 	}
 }
 diff --git a/internal/stock/app/query/check_if_items_in_stock.go b/internal/stock/app/query/check_if_items_in_stock.go
 index d1078f0..4d07201 100644
--- a/internal/stock/app/query/check_if_items_in_stock.go
+++ b/internal/stock/app/query/check_if_items_in_stock.go
@@ -34,12 +34,24 @@ func NewCheckIfItemsInStockHandler(
 	)
 }
 
+// TODO: 删掉
+var stub = map[string]string{
+	"1": "price_1QBYvXRuyMJmUCSsEyQm2oP7",
+	"2": "price_1QBYl4RuyMJmUCSsWt2tgh6d",
+}
+
 func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*orderpb.Item, error) {
 	var res []*orderpb.Item
 	for _, i := range query.Items {
+		// TODO: 改成从数据库 or stripe 获取
+		priceId, ok := stub[i.ID]
+		if !ok {
+			priceId = stub["1"]
+		}
 		res = append(res, &orderpb.Item{
 			ID:       i.ID,
 			Quantity: i.Quantity,
+			PriceID:  priceId,
 		})
 	}
 	return res, nil
~~~
```

## 2. 这组差异到底在解决什么问题

先把问题说透：

上一组 `lesson17 -> lesson18` 已经做出了 `payment` 服务的主链骨架：

`MQ 收到订单 -> CreatePayment 命令 -> 处理器生成链接 -> gRPC 回写 order`

但那时用的是假的 `InmemProcessor`。假的处理器最大的问题不是“它不够真实”，而是：

1. 它不会跟真实支付平台对接
2. 它不能帮你验证订单里的 `PriceID`、支付模式、跳转 URL 这些真实世界约束
3. 你看起来像“链路通了”，但其实很多关键字段根本没被真实消费过

所以这组代码就是在把这条链从演示态推进到第一版真实态。

结合字幕看，作者这组实际上按下面顺序在推进：

1. 先发现回写订单时 `status` 和 `payment link` 是空的
2. 定位到 `order` 内存仓储 `Update` 的入参传错了
3. 接入 `stripe-go` SDK
4. 发现 Stripe 创建 Checkout Session 时 `Price` 不能为空
5. 于是又去 `stock` 侧先临时补一个 `PriceID`

你会发现这里很像真实开发：

`不是先设计出完美模型再一次性交付，而是先把真实链路打通，再沿着报错一层层补齐缺失字段。`

## 3. 运行行为有没有变化

答案：`有，而且是本质变化。`

这组之前：
- `payment` 生成的是本地假链接
- `stock` 返回的 item 只有 `ID` 和 `Quantity`
- `order` 回写可能成功，但也可能因为内存仓储 bug 导致你以为更新了，实际没更新到正确对象

这组之后：
- `payment` 开始直接调用 Stripe 创建真实 Checkout Session
- `stock` 返回的 item 里开始带 `PriceID`
- `payment` 回写 `order` 的 `status/payment link` 终于能真正落进去

所以这不是“代码整理”，而是：

`支付链第一次真正打到了外部支付平台。`

## 4. 正确阅读顺序

这组别从 `stripe.go` 直接硬啃。正确顺序是：

1. [application.go](/g:/shi/go_shop_second/internal/payment/service/application.go)
2. [create_payment.go](/g:/shi/go_shop_second/internal/payment/app/command/create_payment.go)
3. [payment.go](/g:/shi/go_shop_second/internal/payment/domain/payment.go)
4. [stripe.go](/g:/shi/go_shop_second/internal/payment/infrastructure/processor/stripe.go)
5. [check_if_items_in_stock.go](/g:/shi/go_shop_second/internal/stock/app/query/check_if_items_in_stock.go)
6. [order_inmem_repository.go](/g:/shi/go_shop_second/internal/order/adapters/order_inmem_repository.go)

原因是：

- `application.go` 告诉你“系统最终注入的是谁”
- `create_payment.go` 告诉你“处理器接口在业务里怎么被调用”
- `payment.go` 告诉你“处理器抽象长什么样”
- `stripe.go` 才是这个接口的新实现
- `stock` 文件解释 `PriceID` 从哪里来
- `order` 内存仓储解释为什么之前 payment link 没写进去

## 5. 总调用链

这组最重要的脑图你要反复看：

```text
用户创建订单
-> order service 调 stock service 校验库存
-> stock service 返回 item 列表，并补上 PriceID
-> order service 把订单消息发到 MQ
-> payment service 消费订单消息
-> CreatePaymentHandler 调用 Processor.CreatePaymentLink(...)
-> StripeProcessor 向 Stripe 创建 Checkout Session
-> Stripe 返回 result.URL
-> payment service 再通过 gRPC 回写 order.status / order.paymentLink
```

这里有两个你作为小白最容易忽略的事实：

1. `PriceID` 不是 payment 侧现算出来的，而是更早在 `stock` 那边准备好的
2. `payment link` 不是 order 自己生成的，而是 payment 服务拿着订单信息向 Stripe 要来的

## 6. 关键文件详细讲解

### 6.1 `internal/payment/service/application.go`

这文件是装配层。它不写业务细节，但它决定系统最终到底跑哪一个实现。

这组最关键的变化就在这里：以前注入的是 `InmemProcessor`，现在注入的是 `StripeProcessor`。

下面我把这段关键代码按“伪注释”方式写给你看：

```go
func NewApplication(ctx context.Context) (app.Application, func()) {
    // 先建立到 order 服务的 gRPC 客户端。
    // 后面 payment 创建完支付链接，还要回头调用 order.UpdateOrder。
    orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
    if err != nil {
        panic(err)
    }

    // 这里把底层 gRPC client 包一层 adapter，
    // 这样上层命令处理器不用直接依赖生成出来的 gRPC 细节。
    orderGRPC := adapters.NewOrderGRPC(orderClient)

    // 旧实现：只返回假链接，属于教学演示版。
    // memoryProcessor := processor.NewInmemProcessor()

    // 新实现：真正读取配置里的 Stripe API Key，构造 StripeProcessor。
    stripeProcessor := processor.NewStripeProcessor(viper.GetString("stripe-key"))

    // 把真实处理器注入应用层。
    return newApplication(ctx, orderGRPC, stripeProcessor), func() {
        _ = closeOrderClient()
    }
}
```

这段代码为什么重要？

因为它说明：

`这组的真实变化，不是多了一个 stripe.go 文件，而是系统启动时已经决定以后所有 CreatePayment 都走 Stripe 实现了。`

### 这里要补的基础知识

#### `viper.GetString("stripe-key")` 是什么

`viper` 是一个配置库。你可以把它理解成“统一从配置文件、环境变量里读配置”的工具。

这里它的作用是：
- 不把 API Key 写死在代码里
- 启动时从配置中取出 `stripe-key`
- 再交给 `StripeProcessor`

新手常见误区：
- 以为 `viper.GetString(...)` 会自动保证配置存在。其实不会。
- 这里真正做“兜底”的不是 Viper，而是 `NewStripeProcessor` 里的 `panic("empty api key")`

也就是说，这组代码的策略是：

`如果关键配置没给，就让服务启动时直接失败，而不是运行到一半再莫名其妙报错。`

这在基础设施层很常见。

---

### 6.2 `internal/payment/domain/payment.go`

这个文件虽然这组没改，但你必须看，因为 `StripeProcessor` 不是凭空出现的，它是在实现这里定义的接口。

```go
type Processor interface {
    CreatePaymentLink(context.Context, *orderpb.Order) (string, error)
}
```

你先别怕 `interface`。

在这里你可以把它理解成：

`只要谁能根据订单生成支付链接，谁就可以当 Processor。`

上一组是：
- `InmemProcessor` 实现了它

这一组是：
- `StripeProcessor` 也实现了它

所以应用层命令处理器根本不需要知道底层到底是“内存假实现”还是“Stripe 真实现”。

这就是接口抽象的价值：

`上层只关心能力，不关心底层具体怎么做。`

对 Go 小白来说，这里要记住一句话：

`Go 里接口通常不是拿来“炫技”的，而是拿来隔离变化。`

这组变化里，“会变化”的就是支付链接生成方式，所以这里抽象成接口很合理。

---

### 6.3 `internal/payment/app/command/create_payment.go`

这是真正的业务调用点。理解这段后，你才知道 `StripeProcessor` 在整条链里处于什么位置。

```go
func (c createPaymentHandler) Handle(ctx context.Context, cmd CreatePayment) (string, error) {
    // 第一步：让 Processor 根据订单生成支付链接。
    link, err := c.processor.CreatePaymentLink(ctx, cmd.Order)
    if err != nil {
        return "", err
    }

    // 第二步：构造一个新的订单对象，把状态改成 waiting_for_payment，
    // 同时把 Stripe 返回的支付链接写进去。
    newOrder := &orderpb.Order{
        ID:          cmd.Order.ID,
        CustomerID:  cmd.Order.CustomerID,
        Status:      "waiting_for_payment",
        Items:       cmd.Order.Items,
        PaymentLink: link,
    }

    // 第三步：调用 order 服务的 gRPC 接口，把订单状态回写回去。
    err = c.orderGRPC.UpdateOrder(ctx, newOrder)
    return link, err
}
```

这段代码解释了为什么字幕 `19.txt` 会先去查 `Update` bug。

因为这里的业务意图非常明确：

- `processor.CreatePaymentLink(...)` 成功后
- 一定要再调用 `UpdateOrder(...)`
- 否则你在 order 服务里查到的订单，就还是旧状态，`payment link` 也还是空的

所以那不是一个“无关紧要的小 bug”，而是直接卡住主链的 bug。

### 这里要补的基础知识

#### `context.Context` 是干什么的

你会看到这里一路都在传 `ctx`。

在这个项目里，你先把 `context` 理解成“请求上下文传送带”就够了。它至少承担三种作用：

1. 传递超时和取消信号
2. 让日志/链路追踪拿到同一条请求的信息
3. 在跨服务调用时，把同一条请求继续往下传

这一组里虽然还没把 tracing 补完整，但这个 `ctx` 已经是未来做超时控制和链路追踪的基础了。

---

### 6.4 `internal/payment/infrastructure/processor/stripe.go`

这是本组最核心的新文件。你应该按“整文件”去看它，而不是盯某一两行。

下面是这份文件的教学版注释写法，我会把关键行都解释出来。

```go
package processor

import (
    "context"
    "encoding/json"

    // 这里依赖的是 orderpb.Order，也就是 protobuf 生成的订单类型。
    // 这说明此时 payment 层仍然直接使用外层 DTO，
    // 还没有进一步做领域对象隔离。
    "github.com/ghost-yu/go_shop_second/common/genproto/orderpb"

    // stripe-go 是 Stripe 官方 Go SDK。
    // 这里用它来创建 Checkout Session。
    "github.com/stripe/stripe-go/v79"
    "github.com/stripe/stripe-go/v79/checkout/session"
)

type StripeProcessor struct {
    // 当前结构体只保存 apiKey。
    // 但要注意：真正生效的不只是这个字段，
    // 下面还会把 stripe.Key 这个全局变量赋值掉。
    apiKey string
}

func NewStripeProcessor(apiKey string) *StripeProcessor {
    // 如果 key 为空，启动直接失败。
    // 这里用 panic 很激进，但在教学项目里能尽早暴露配置问题。
    if apiKey == "" {
        panic("empty api key")
    }

    // 这是 Stripe SDK 的一个重要点：
    // stripe.Key 是全局变量，不是某个 client 实例字段。
    // 赋值一次后，后面 session.New(...) 就会默认拿这个 key。
    stripe.Key = apiKey

    return &StripeProcessor{apiKey: apiKey}
}

var (
    // 支付成功后的跳转地址先写死。
    // 这很明显是教学版/过渡态实现，后面大概率要改成配置项。
    successURL = "http://localhost:8282/success"
)

func (s StripeProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
    // Stripe Checkout Session 里的商品行。
    // 你可以把它理解成支付页面上展示的购买项列表。
    var items []*stripe.CheckoutSessionLineItemParams

    // 遍历订单里的每个商品，把它转换成 Stripe 认识的 line item。
    for _, item := range order.Items {
        items = append(items, &stripe.CheckoutSessionLineItemParams{
            // Stripe 这里要的不是你的商品 ID，而是 Stripe 那边的 Price ID。
            Price: stripe.String(item.PriceID),

            // SDK 这里要求 int64，所以项目里的数量要做一次类型转换。
            Quantity: stripe.Int64(int64(item.Quantity)),
        })
    }

    // 把订单商品序列化成 JSON，后面塞进 metadata。
    // 这里故意忽略了 Marshal 的错误，这是一个不够严谨的地方。
    marshalledItems, _ := json.Marshal(order.Items)

    // metadata 是 Stripe 提供的一块“自定义透传信息”。
    // 你可以理解成“这次支付顺手绑在 session 上的一包业务上下文”。
    metadata := map[string]string{
        "orderID":    order.ID,
        "customerID": order.CustomerID,
        "status":     order.Status,
        "items":      string(marshalledItems),
    }

    // 组装创建 Checkout Session 所需参数。
    params := &stripe.CheckoutSessionParams{
        Metadata:   metadata,
        LineItems:  items,

        // 这里明确声明这是一次性支付，不是订阅。
        Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),

        // 支付成功后跳回本地页面。
        SuccessURL: stripe.String(successURL),
    }

    // 真正发请求给 Stripe。
    result, err := session.New(params)
    if err != nil {
        return "", err
    }

    // Stripe 返回的 URL 就是给前端跳转的支付链接。
    return result.URL, nil
}
```

### 这一整份文件真正重要的点

#### 第一，`PriceID` 才是 Stripe 真正认识的价格对象

新手最容易把这件事想错成：

`我不是已经有商品 ID 了吗，为什么还要 PriceID？`

因为在 Stripe 的模型里：
- 商品是商品
- 价格是价格
- 一个商品甚至可能挂多个价格

所以真正创建支付单时，要传的是 `PriceID`，不是你业务里的 `Item.ID`。

这也是为什么 `20.TXT` 字幕会单独讲这一件事。

#### 第二，`metadata` 不是“必需字段”，但非常有用

这段代码把 `orderID/customerID/status/items` 都塞进了 `metadata`。

为什么这么做？

因为后面 Stripe 回调 webhook 给你时，你需要知道：
- 这次支付到底对应哪个订单
- 属于哪个用户
- 你当时传进去的商品是什么

`metadata` 就像给远程支付平台绑了一张“业务便签”。

#### 第三，`stripe.Key = apiKey` 是全局状态

这是 `stripe-go` 的典型用法，但也是新手容易踩的点。

坑在这里：
- 它不是“某个对象自己的配置”
- 它是整个进程级别的全局变量

这意味着：
- 如果你同一个进程里要连多个 Stripe 账户，会很麻烦
- 测试时也容易出现互相污染

但在当前这个教学项目里，只连一个 Stripe 账户，所以这样写能先快速跑通。

#### 第四，`successURL` 写死是明显的过渡态

这行能跑，但不成熟。为什么？

因为真实系统里这个 URL 往往跟环境有关：
- 本地环境一个地址
- 测试环境一个地址
- 生产环境一个地址

而且有时候还会带订单参数、用户参数、前端路由参数。

所以这里是典型的：

`先把链路打通，再谈配置化。`

---

### 6.5 `internal/stock/app/query/check_if_items_in_stock.go`

这一组里这份文件看起来像“小修”，但实际上它是 Stripe 能不能工作的前提之一。

因为 `StripeProcessor` 会读 `item.PriceID`，如果这里不补，前面 `Price` 就是空字符串，Stripe 会直接报错。

这和字幕 `20.TXT` 完全对应。

下面把这组变更对应的历史代码按注释方式写出来：

```go
// TODO: 删掉
var stub = map[string]string{
    // 商品 1 先硬编码映射到一个 Stripe Price ID。
    "1": "price_1QBYvXRuyMJmUCSsEyQm2oP7",

    // 商品 2 映射到另一个 Stripe Price ID。
    "2": "price_1QBYl4RuyMJmUCSsWt2tgh6d",
}

func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*orderpb.Item, error) {
    var res []*orderpb.Item

    for _, i := range query.Items {
        // 这里先临时从 stub 里拿 price id。
        // 如果没有对应项，就退回到默认值 "1"。
        // 这非常粗糙，但能保证 Stripe 那边不再收到空 price。
        priceId, ok := stub[i.ID]
        if !ok {
            priceId = stub["1"]
        }

        res = append(res, &orderpb.Item{
            ID:       i.ID,
            Quantity: i.Quantity,
            PriceID:  priceId,
        })
    }

    return res, nil
}
```

### 这里的设计意图是什么

这段不是在认真建库存模型，而是在做一件非常务实的事：

`让上游 order 拿到带 PriceID 的 item，让下游 payment 真能调用 Stripe。`

这正是教学项目里很典型的一种过渡写法。

### 这里的问题也要看清楚

#### 1. `stub` 是硬编码

这说明：
- 当前商品系统和 Stripe 商品体系还没真正打通
- 只是先手工映射了两个商品

#### 2. 未命中时直接回退到 `stub["1"]`

这是很危险的默认行为。

因为这会导致：
- 你请求商品 `999`
- 系统却可能给你商品 `1` 的价格
- 这样虽然“不报错”，但业务是错的

教学里这样写是为了快速跑通，但生产里很容易造成错单。

#### 3. `TODO: 改成从数据库 or stripe 获取`

这句注释其实是在明确告诉你：

`作者自己也知道这不是最终方案。`

这点你以后读代码要养成习惯：

不是所有“能跑”的代码都代表“设计完成”。

---

### 6.6 `internal/order/adapters/order_inmem_repository.go`

这一行改动最少，但不能忽视。

原来：

```go
updatedOrder, err := updateFn(ctx, o)
```

改成：

```go
updatedOrder, err := updateFn(ctx, order)
```

你第一次看可能会疑惑：

`这不就一个参数从 o 换成 order 吗，为什么字幕里专门讲它？`

因为这里直接决定 `UpdateFn` 到底收到哪个对象。

### 先理解这两个变量分别是谁

- `o`：仓储里已经存着的旧订单对象
- `order`：调用 `Update(...)` 时传进来的新订单对象

在当前这条链里，payment 服务在生成支付链接后，会构造一个新订单对象：
- `Status = waiting_for_payment`
- `PaymentLink = Stripe 返回的链接`

然后把这个新订单交给 `UpdateOrder(...)`。

如果这里把旧对象 `o` 传给 `updateFn`，那么 `updateFn` 里看到的还是旧数据，
就可能把新的 `status/payment link` 丢掉。

这就是为什么字幕里会出现这种现象：
- 你明明在 payment 侧收到 link 了
- 但是再去查 order，发现 `status` 和 `payment link` 还是空的

### 这说明了什么

说明这组不只是“第三方支付接入”，还顺手修了一个会掩盖主链正确性的 bug。

也就是说，如果你只接 Stripe，不修这行，效果可能是：

`Stripe 真创建成功了，但 order 服务查出来仍像没成功一样。`

这会让你排查方向完全跑偏。

## 7. 第三方库和容易忘的点

### 7.1 `stripe-go`

这组最关键的外部库就是 `github.com/stripe/stripe-go/v79`。

这里你至少要知道四件事：

1. 它是 Stripe 官方 Go SDK
2. 当前代码只用了 Checkout Session 这一小块功能
3. `stripe.Key` 是全局变量，不是实例字段
4. `session.New(params)` 会真实发请求到 Stripe

新手容易忽略的坑：

- `Price` 不能为空
- `Mode` 要跟价格类型匹配，不然 Stripe 会报错
- `SuccessURL` 只是支付成功后的浏览器跳转，不等于支付完成的最终确认机制

最后这一条特别重要。

因为很多新手会误以为：

`浏览器跳到 success 页面 = 订单一定支付成功`

其实不是。

真正更稳的做法通常是：
- 前端跳 success 页面只是用户体验
- 后端再靠 webhook 接收 Stripe 的异步通知，确认最终支付状态

而下一组课程正是在补 webhook。

### 7.2 `json.Marshal`

这组 `metadata` 里用到了 `json.Marshal(order.Items)`。

这是什么意思？

就是把 Go 里的结构体切片，转成 JSON 字符串，方便塞进 `map[string]string`。

这里最容易忽略的坑是：

```go
marshalledItems, _ := json.Marshal(order.Items)
```

错误被直接丢掉了。

这在教学项目里常见，但不是好习惯。

更稳的写法应该是：
- 检查 `err`
- 如果序列化失败，就明确返回错误

### 7.3 `viper`

这组对 Viper 的使用非常轻，但作用很关键：读 Stripe 的 API Key。

要记住：
- `viper` 只是“取配置”
- 真正决定“不允许空值”的，是你自己的业务代码

所以配置系统和业务校验是两层职责，不要混为一谈。

## 8. 为什么这么设计

这组设计如果一句话总结，就是：

`先最小代价接上真实支付，再沿真实报错把缺的字段和 bug 补齐。`

这样做的好处：
- 推进快
- 反馈真实
- 哪条链没通，马上会暴露

坏处也很明显：
- 过程里会出现很多硬编码
- 代码里会混着教学态和真实态
- 结构上还不够干净

比如这组里就有几个非常明显的过渡态信号：

1. `successURL` 写死
2. `stub` 写死 `PriceID`
3. `json.Marshal` 错误忽略
4. `payment` 直接依赖 `orderpb.Order`
5. `stripe.Key` 用的是全局变量配置方式

所以这组的正确评价不是“设计完成了”，而是：

`第一版真实集成已经成立，但远没有到可生产的程度。`

## 9. 当前还不完美的地方

这部分你必须会看，因为以后面试或者自己复盘项目，这里才是拉开水平的地方。

### 1. `PriceID` 获取方式非常粗糙

现在是 `stub`。这意味着商品和 Stripe 价格体系还没真正建模。

更成熟的方案通常是：
- 商品表里存 Stripe 对应 PriceID
- 或者由 stock / product 服务统一维护映射
- 绝不会未命中时偷偷退回默认价格

### 2. `successURL` 没带业务参数

这一版只是跳到固定页面，没把订单信息显式带回去。

真实项目里通常会：
- 带 `orderID`
- 带 `customerID`
- 或者带一个前端能拿去查询状态的 token

### 3. 只靠 success 页面还不够

这一版能让你“看起来像支付成功了”，但最终一致性还没完整建立。

真正的支付确认应该依赖 webhook，而不是浏览器跳转。

### 4. `payment` 层直接吃 `orderpb.Order`

这说明 protobuf 生成的外层传输对象，已经渗进了应用层和基础设施层。

教学里先这么写没问题，但后面课程也确实会继续清理这类 pb 泄漏问题。

### 5. 内存仓储 bug 暴露了测试覆盖不足

`Update(ctx, order, updateFn)` 这种错误其实很适合被单元测试提前抓出来。

这也说明：
- 链路打通了不代表行为完全对
- 没测试时，这类错误通常要靠手工联调才会暴露

## 10. 这组最该带走的知识点

你把这组学完，最该记住的不是某一行 SDK 调用，而是下面这几件事：

1. 真正接第三方支付时，`PriceID` 往往比商品 ID 更重要
2. 抽象成 `Processor` 接口，是为了隔离“支付链接怎么生成”这件会变化的事
3. 配置库只负责取配置，关键配置的合法性仍要你自己兜底
4. 支付成功页不等于最终支付确认，最终还是要靠 webhook
5. 一个很小的仓储更新 bug，就能把整条业务链的结果伪装成“没成功”
6. 教学项目常常会先写硬编码 stub，这不是最佳实践，但有助于先跑通主链

## 11. 最后用一句话收住这组

`lesson18 -> lesson19` 的本质，不是“接了 Stripe SDK”，而是“支付链第一次开始面对真实支付平台的约束，于是项目被迫补上 PriceID、修正回写 bug，并把假支付链接替换成真实 Checkout Session URL”。`

你后面再看下一组 `lesson19 -> lesson20`、`lesson20 -> lesson21` 时，会更容易理解为什么马上就会出现：
- 前端 success 页面
- Stripe webhook
- 支付完成后的订单状态同步

因为这组已经把“真实支付入口”打开了。
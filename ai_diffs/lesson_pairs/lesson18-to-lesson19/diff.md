# Lesson Pair Diff Report

- FromBranch: lesson18
- ToBranch: lesson19

## Short Summary

~~~text
 6 files changed, 78 insertions(+), 4 deletions(-)
~~~

## File Stats

~~~text
 internal/order/adapters/order_inmem_repository.go  |  2 +-
 internal/payment/go.mod                            |  3 +-
 internal/payment/go.sum                            |  4 ++
 .../payment/infrastructure/processor/stripe.go     | 55 ++++++++++++++++++++++
 internal/payment/service/application.go            |  6 ++-
 .../stock/app/query/check_if_items_in_stock.go     | 12 +++++
 6 files changed, 78 insertions(+), 4 deletions(-)
~~~

## Commit Comparison

~~~text
> 1cb5423 small fix
> 9520cda stripe processor fix
> 14602b2 stripe processor
~~~

## Changed Files

~~~text
internal/order/adapters/order_inmem_repository.go
internal/payment/go.mod
internal/payment/go.sum
internal/payment/infrastructure/processor/stripe.go
internal/payment/service/application.go
internal/stock/app/query/check_if_items_in_stock.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
internal/order/adapters/order_inmem_repository.go
internal/payment/infrastructure/processor/stripe.go
internal/payment/service/application.go
internal/stock/app/query/check_if_items_in_stock.go
~~~

## Full Diff

~~~diff
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

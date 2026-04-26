# Lesson Pair Diff Report

- FromBranch: lesson20
- ToBranch: lesson21

## Short Summary

~~~text
 6 files changed, 93 insertions(+), 11 deletions(-)
~~~

## File Stats

~~~text
 internal/common/config/global.yaml                 |  3 +-
 internal/common/config/viper.go                    |  2 +-
 internal/payment/domain/payment.go                 |  8 +++
 internal/payment/http.go                           | 75 +++++++++++++++++++++-
 .../payment/infrastructure/processor/stripe.go     | 14 ++--
 internal/payment/main.go                           |  2 +-
 6 files changed, 93 insertions(+), 11 deletions(-)
~~~

## Commit Comparison

~~~text
> 55bc3a6 stripe webhook
~~~

## Changed Files

~~~text
internal/common/config/global.yaml
internal/common/config/viper.go
internal/payment/domain/payment.go
internal/payment/http.go
internal/payment/infrastructure/processor/stripe.go
internal/payment/main.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
internal/common/config/global.yaml
internal/common/config/viper.go
internal/payment/domain/payment.go
internal/payment/http.go
internal/payment/infrastructure/processor/stripe.go
internal/payment/main.go
~~~

## Full Diff

~~~diff
diff --git a/internal/common/config/global.yaml b/internal/common/config/global.yaml
index 17a4f41..7dd93c9 100644
--- a/internal/common/config/global.yaml
+++ b/internal/common/config/global.yaml
@@ -27,4 +27,5 @@ rabbitmq:
   host: 127.0.0.1
   port: 5672
 
-stripe-key: "${STRIPE_KEY}"
\ No newline at end of file
+stripe-key: "${STRIPE_KEY}"
+endpoint-stripe-secret: "${ENDPOINT_STRIPE_SECRET}"
\ No newline at end of file
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
 
 func (h *PaymentHandler) RegisterRoutes(c *gin.Engine) {
@@ -18,4 +31,62 @@ func (h *PaymentHandler) RegisterRoutes(c *gin.Engine) {
 
 func (h *PaymentHandler) handleWebhook(c *gin.Context) {
 	logrus.Info("receive webhook from stripe")
+	const MaxBodyBytes = int64(65536)
+	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
+	payload, err := io.ReadAll(c.Request.Body)
+	if err != nil {
+		logrus.Infof("Error reading request body: %v\n", err)
+		c.JSON(http.StatusServiceUnavailable, err.Error())
+		return
+	}
+
+	event, err := webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"),
+		viper.GetString("ENDPOINT_STRIPE_SECRET"))
+
+	if err != nil {
+		logrus.Infof("Error verifying webhook signature: %v\n", err)
+		c.JSON(http.StatusBadRequest, err.Error())
+		return
+	}
+
+	switch event.Type {
+	case stripe.EventTypeCheckoutSessionCompleted:
+		var session stripe.CheckoutSession
+		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
+			logrus.Infof("error unmarshal event.data.raw into session, err = %v", err)
+			c.JSON(http.StatusBadRequest, err.Error())
+			return
+		}
+
+		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
+			logrus.Infof("payment for checkout session %v success!", session.ID)
+
+			ctx, cancel := context.WithCancel(context.TODO())
+			defer cancel()
+
+			var items []*orderpb.Item
+			_ = json.Unmarshal([]byte(session.Metadata["items"]), &items)
+
+			marshalledOrder, err := json.Marshal(&domain.Order{
+				ID:          session.Metadata["orderID"],
+				CustomerID:  session.Metadata["customerID"],
+				Status:      string(stripe.CheckoutSessionPaymentStatusPaid),
+				PaymentLink: session.Metadata["paymentLink"],
+				Items:       items,
+			})
+			if err != nil {
+				logrus.Infof("error marshal domain.order, err = %v", err)
+				c.JSON(http.StatusBadRequest, err.Error())
+				return
+			}
+
+			_ = h.channel.PublishWithContext(ctx, broker.EventOrderPaid, "", false, false, amqp.Publishing{
+				ContentType:  "application/json",
+				DeliveryMode: amqp.Persistent,
+				Body:         marshalledOrder,
+			})
+			logrus.Infof("message published to %s, body: %s", broker.EventOrderPaid, string(marshalledOrder))
+		}
+	}
+	c.JSON(http.StatusOK, nil)
 }
diff --git a/internal/payment/infrastructure/processor/stripe.go b/internal/payment/infrastructure/processor/stripe.go
index 7111272..37fcb1d 100644
--- a/internal/payment/infrastructure/processor/stripe.go
+++ b/internal/payment/infrastructure/processor/stripe.go
@@ -3,6 +3,7 @@ package processor
 import (
 	"context"
 	"encoding/json"
+	"fmt"
 
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	"github.com/stripe/stripe-go/v79"
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
 	result, err := session.New(params)
 	if err != nil {
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
~~~

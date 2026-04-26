# Lesson Pair Diff Report

- FromBranch: lesson25
- ToBranch: lesson26

## Short Summary

~~~text
 8 files changed, 93 insertions(+), 10 deletions(-)
~~~

## File Stats

~~~text
 internal/common/broker/rabbitmq.go                 | 36 ++++++++++++++++++++++
 internal/order/app/command/create_order.go         | 18 ++++++++---
 internal/order/infrastructure/consumer/consumer.go | 15 +++++++--
 internal/payment/adapters/order_grpc.go            |  4 +++
 internal/payment/http.go                           | 10 +++++-
 .../payment/infrastructure/consumer/consumer.go    |  9 +++++-
 .../payment/infrastructure/processor/stripe.go     |  4 +++
 internal/stock/ports/grpc.go                       |  7 +++++
 8 files changed, 93 insertions(+), 10 deletions(-)
~~~

## Commit Comparison

~~~text
> acdb857 more otel
~~~

## Changed Files

~~~text
internal/common/broker/rabbitmq.go
internal/order/app/command/create_order.go
internal/order/infrastructure/consumer/consumer.go
internal/payment/adapters/order_grpc.go
internal/payment/http.go
internal/payment/infrastructure/consumer/consumer.go
internal/payment/infrastructure/processor/stripe.go
internal/stock/ports/grpc.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
internal/common/broker/rabbitmq.go
internal/order/app/command/create_order.go
internal/order/infrastructure/consumer/consumer.go
internal/payment/adapters/order_grpc.go
internal/payment/http.go
internal/payment/infrastructure/consumer/consumer.go
internal/payment/infrastructure/processor/stripe.go
internal/stock/ports/grpc.go
~~~

## Full Diff

~~~diff
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
 
 func Connect(user, password, host, port string) (*amqp.Channel, func() error) {
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
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
index 7185157..282f66c 100644
--- a/internal/order/app/command/create_order.go
+++ b/internal/order/app/command/create_order.go
@@ -4,6 +4,7 @@ import (
 	"context"
 	"encoding/json"
 	"errors"
+	"fmt"
 
 	"github.com/ghost-yu/go_shop_second/common/broker"
 	"github.com/ghost-yu/go_shop_second/common/decorator"
@@ -12,6 +13,7 @@ import (
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
+	"go.opentelemetry.io/otel"
 )
 
 type CreateOrder struct {
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
diff --git a/internal/order/infrastructure/consumer/consumer.go b/internal/order/infrastructure/consumer/consumer.go
index e7ce252..623bb50 100644
--- a/internal/order/infrastructure/consumer/consumer.go
+++ b/internal/order/infrastructure/consumer/consumer.go
@@ -3,6 +3,7 @@ package consumer
 import (
 	"context"
 	"encoding/json"
+	"fmt"
 
 	"github.com/ghost-yu/go_shop_second/common/broker"
 	"github.com/ghost-yu/go_shop_second/order/app"
@@ -10,6 +11,7 @@ import (
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
+	"go.opentelemetry.io/otel"
 )
 
 type Consumer struct {
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
 
 type PaymentHandler struct {
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
diff --git a/internal/payment/infrastructure/consumer/consumer.go b/internal/payment/infrastructure/consumer/consumer.go
index 7f01e7a..99046de 100644
--- a/internal/payment/infrastructure/consumer/consumer.go
+++ b/internal/payment/infrastructure/consumer/consumer.go
@@ -3,6 +3,7 @@ package consumer
 import (
 	"context"
 	"encoding/json"
+	"fmt"
 
 	"github.com/ghost-yu/go_shop_second/common/broker"
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
@@ -10,6 +11,7 @@ import (
 	"github.com/ghost-yu/go_shop_second/payment/app/command"
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
+	"go.opentelemetry.io/otel"
 )
 
 type Consumer struct {
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
 	if err != nil {
 		return nil, err
@@ -25,6 +29,9 @@ func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsReque
 }
 
 func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
+	_, span := tracing.Start(ctx, "CheckIfItemsInStock")
+	defer span.End()
+
 	items, err := G.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{Items: request.Items})
 	if err != nil {
 		return nil, err
~~~

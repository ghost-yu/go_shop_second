# Lesson Pair Diff Report

- FromBranch: lesson62
- ToBranch: lesson63

## Short Summary

~~~text
 25 files changed, 65 insertions(+), 210 deletions(-)
~~~

## File Stats

~~~text
 internal/{order => common}/convertor/convertor.go  | 15 ++--
 internal/{order => common}/convertor/facade.go     |  0
 internal/{stock => common}/entity/entity.go        | 16 ++--
 .../kitchen/infrastructure/consumer/consumer.go    | 16 ++--
 internal/order/adapters/order_mongo_repository.go  |  2 +-
 internal/order/app/command/create_order.go         |  4 +-
 internal/order/domain/order/order.go               |  2 +-
 internal/order/entity/entity.go                    | 13 ---
 internal/order/http.go                             | 10 ++-
 internal/order/ports/grpc.go                       | 11 ++-
 internal/payment/app/command/create_payment.go     | 11 ++-
 internal/payment/domain/payment.go                 |  6 +-
 internal/payment/http.go                           |  7 +-
 .../payment/infrastructure/consumer/consumer.go    |  4 +-
 internal/payment/infrastructure/processor/inmem.go |  4 +-
 .../payment/infrastructure/processor/stripe.go     |  4 +-
 internal/stock/adapters/stock_inmem_repository.go  |  2 +-
 internal/stock/adapters/stock_mysql_repository.go  |  2 +-
 .../stock/adapters/stock_mysql_repository_test.go  |  2 +-
 .../stock/app/query/check_if_items_in_stock.go     |  2 +-
 internal/stock/app/query/get_items.go              |  2 +-
 internal/stock/convertor/convertor.go              | 97 ----------------------
 internal/stock/convertor/facade.go                 | 39 ---------
 internal/stock/domain/stock/repository.go          |  2 +-
 internal/stock/ports/grpc.go                       |  2 +-
 25 files changed, 65 insertions(+), 210 deletions(-)
~~~

## Commit Comparison

~~~text
> d767896 acl cleanup
~~~

## Changed Files

~~~text
internal/common/convertor/convertor.go
internal/common/convertor/facade.go
internal/common/entity/entity.go
internal/kitchen/infrastructure/consumer/consumer.go
internal/order/adapters/order_mongo_repository.go
internal/order/app/command/create_order.go
internal/order/domain/order/order.go
internal/order/entity/entity.go
internal/order/http.go
internal/order/ports/grpc.go
internal/payment/app/command/create_payment.go
internal/payment/domain/payment.go
internal/payment/http.go
internal/payment/infrastructure/consumer/consumer.go
internal/payment/infrastructure/processor/inmem.go
internal/payment/infrastructure/processor/stripe.go
internal/stock/adapters/stock_inmem_repository.go
internal/stock/adapters/stock_mysql_repository.go
internal/stock/adapters/stock_mysql_repository_test.go
internal/stock/app/query/check_if_items_in_stock.go
internal/stock/app/query/get_items.go
internal/stock/convertor/convertor.go
internal/stock/convertor/facade.go
internal/stock/domain/stock/repository.go
internal/stock/ports/grpc.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
internal/common/convertor/convertor.go
internal/common/convertor/facade.go
internal/common/entity/entity.go
internal/kitchen/infrastructure/consumer/consumer.go
internal/order/adapters/order_mongo_repository.go
internal/order/app/command/create_order.go
internal/order/domain/order/order.go
internal/order/entity/entity.go
internal/order/http.go
internal/order/ports/grpc.go
internal/payment/app/command/create_payment.go
internal/payment/domain/payment.go
internal/payment/http.go
internal/payment/infrastructure/consumer/consumer.go
internal/payment/infrastructure/processor/inmem.go
internal/payment/infrastructure/processor/stripe.go
internal/stock/adapters/stock_inmem_repository.go
internal/stock/adapters/stock_mysql_repository.go
internal/stock/adapters/stock_mysql_repository_test.go
internal/stock/app/query/check_if_items_in_stock.go
internal/stock/app/query/get_items.go
internal/stock/convertor/convertor.go
internal/stock/convertor/facade.go
internal/stock/domain/stock/repository.go
internal/stock/ports/grpc.go
~~~

## Full Diff

~~~diff
diff --git a/internal/order/convertor/convertor.go b/internal/common/convertor/convertor.go
similarity index 90%
rename from internal/order/convertor/convertor.go
rename to internal/common/convertor/convertor.go
index dc39c0b..fc5c2fd 100644
--- a/internal/order/convertor/convertor.go
+++ b/internal/common/convertor/convertor.go
@@ -2,9 +2,8 @@ package convertor
 
 import (
 	client "github.com/ghost-yu/go_shop_second/common/client/order"
+	"github.com/ghost-yu/go_shop_second/common/entity"
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
-	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
-	"github.com/ghost-yu/go_shop_second/order/entity"
 )
 
 type OrderConvertor struct{}
@@ -53,7 +52,7 @@ func (c *ItemWithQuantityConvertor) ClientToEntity(i client.ItemWithQuantity) *e
 	}
 }
 
-func (c *OrderConvertor) EntityToProto(o *domain.Order) *orderpb.Order {
+func (c *OrderConvertor) EntityToProto(o *entity.Order) *orderpb.Order {
 	c.check(o)
 	return &orderpb.Order{
 		ID:          o.ID,
@@ -64,9 +63,9 @@ func (c *OrderConvertor) EntityToProto(o *domain.Order) *orderpb.Order {
 	}
 }
 
-func (c *OrderConvertor) ProtoToEntity(o *orderpb.Order) *domain.Order {
+func (c *OrderConvertor) ProtoToEntity(o *orderpb.Order) *entity.Order {
 	c.check(o)
-	return &domain.Order{
+	return &entity.Order{
 		ID:          o.ID,
 		CustomerID:  o.CustomerID,
 		Status:      o.Status,
@@ -75,9 +74,9 @@ func (c *OrderConvertor) ProtoToEntity(o *orderpb.Order) *domain.Order {
 	}
 }
 
-func (c *OrderConvertor) ClientToEntity(o *client.Order) *domain.Order {
+func (c *OrderConvertor) ClientToEntity(o *client.Order) *entity.Order {
 	c.check(o)
-	return &domain.Order{
+	return &entity.Order{
 		ID:          o.Id,
 		CustomerID:  o.CustomerId,
 		Status:      o.Status,
@@ -86,7 +85,7 @@ func (c *OrderConvertor) ClientToEntity(o *client.Order) *domain.Order {
 	}
 }
 
-func (c *OrderConvertor) EntityToClient(o *domain.Order) *client.Order {
+func (c *OrderConvertor) EntityToClient(o *entity.Order) *client.Order {
 	c.check(o)
 	return &client.Order{
 		Id:          o.ID,
diff --git a/internal/order/convertor/facade.go b/internal/common/convertor/facade.go
similarity index 100%
rename from internal/order/convertor/facade.go
rename to internal/common/convertor/facade.go
diff --git a/internal/stock/entity/entity.go b/internal/common/entity/entity.go
similarity index 100%
rename from internal/stock/entity/entity.go
rename to internal/common/entity/entity.go
index 671b85b..7ad56e7 100644
--- a/internal/stock/entity/entity.go
+++ b/internal/common/entity/entity.go
@@ -1,13 +1,5 @@
 package entity
 
-type Order struct {
-	ID          string
-	CustomerID  string
-	Status      string
-	PaymentLink string
-	Items       []*Item
-}
-
 type Item struct {
 	ID       string
 	Name     string
@@ -19,3 +11,11 @@ type ItemWithQuantity struct {
 	ID       string
 	Quantity int32
 }
+
+type Order struct {
+	ID          string
+	CustomerID  string
+	Status      string
+	PaymentLink string
+	Items       []*Item
+}
diff --git a/internal/kitchen/infrastructure/consumer/consumer.go b/internal/kitchen/infrastructure/consumer/consumer.go
index 9e17a9e..0376646 100644
--- a/internal/kitchen/infrastructure/consumer/consumer.go
+++ b/internal/kitchen/infrastructure/consumer/consumer.go
@@ -7,6 +7,8 @@ import (
 	"time"
 
 	"github.com/ghost-yu/go_shop_second/common/broker"
+	"github.com/ghost-yu/go_shop_second/common/convertor"
+	"github.com/ghost-yu/go_shop_second/common/entity"
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/pkg/errors"
@@ -23,14 +25,6 @@ type Consumer struct {
 	orderGRPC OrderService
 }
 
-type Order struct {
-	ID          string
-	CustomerID  string
-	Status      string
-	PaymentLink string
-	Items       []*orderpb.Item
-}
-
 func NewConsumer(orderGRPC OrderService) *Consumer {
 	return &Consumer{
 		orderGRPC: orderGRPC,
@@ -75,7 +69,7 @@ func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Que
 		}
 	}()
 
-	o := &Order{}
+	o := &entity.Order{}
 	if err = json.Unmarshal(msg.Body, o); err != nil {
 		err = errors.Wrap(err, "error unmarshal msg.body into order")
 		return
@@ -90,7 +84,7 @@ func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Que
 		ID:          o.ID,
 		CustomerID:  o.CustomerID,
 		Status:      "ready",
-		Items:       o.Items,
+		Items:       convertor.NewItemConvertor().EntitiesToProtos(o.Items),
 		PaymentLink: o.PaymentLink,
 	}); err != nil {
 		logging.Errorf(ctx, nil, "error updating order||orderID=%s||err=%v", o.ID, err)
@@ -102,7 +96,7 @@ func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Que
 	span.AddEvent("kitchen.order.finished.updated")
 }
 
-func cook(ctx context.Context, o *Order) {
+func cook(ctx context.Context, o *entity.Order) {
 	logrus.WithContext(ctx).Printf("cooking order: %s", o.ID)
 	time.Sleep(5 * time.Second)
 	logrus.WithContext(ctx).Printf("order %s done!", o.ID)
diff --git a/internal/order/adapters/order_mongo_repository.go b/internal/order/adapters/order_mongo_repository.go
index fa7cf40..31bac38 100644
--- a/internal/order/adapters/order_mongo_repository.go
+++ b/internal/order/adapters/order_mongo_repository.go
@@ -4,9 +4,9 @@ import (
 	"context"
 
 	_ "github.com/ghost-yu/go_shop_second/common/config"
+	"github.com/ghost-yu/go_shop_second/common/entity"
 	"github.com/ghost-yu/go_shop_second/common/logging"
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
-	"github.com/ghost-yu/go_shop_second/order/entity"
 	"github.com/spf13/viper"
 	"go.mongodb.org/mongo-driver/bson"
 	"go.mongodb.org/mongo-driver/bson/primitive"
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
index 74a0319..d98a379 100644
--- a/internal/order/app/command/create_order.go
+++ b/internal/order/app/command/create_order.go
@@ -5,12 +5,12 @@ import (
 	"fmt"
 
 	"github.com/ghost-yu/go_shop_second/common/broker"
+	"github.com/ghost-yu/go_shop_second/common/convertor"
 	"github.com/ghost-yu/go_shop_second/common/decorator"
+	"github.com/ghost-yu/go_shop_second/common/entity"
 	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/ghost-yu/go_shop_second/order/app/query"
-	"github.com/ghost-yu/go_shop_second/order/convertor"
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
-	"github.com/ghost-yu/go_shop_second/order/entity"
 	"github.com/pkg/errors"
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
diff --git a/internal/order/domain/order/order.go b/internal/order/domain/order/order.go
index 1afa1cb..5776b33 100644
--- a/internal/order/domain/order/order.go
+++ b/internal/order/domain/order/order.go
@@ -3,7 +3,7 @@ package order
 import (
 	"fmt"
 
-	"github.com/ghost-yu/go_shop_second/order/entity"
+	"github.com/ghost-yu/go_shop_second/common/entity"
 	"github.com/pkg/errors"
 	"github.com/stripe/stripe-go/v80"
 )
diff --git a/internal/order/entity/entity.go b/internal/order/entity/entity.go
deleted file mode 100644
index 51016cc..0000000
--- a/internal/order/entity/entity.go
+++ /dev/null
@@ -1,13 +0,0 @@
-package entity
-
-type Item struct {
-	ID       string
-	Name     string
-	Quantity int32
-	PriceID  string
-}
-
-type ItemWithQuantity struct {
-	ID       string
-	Quantity int32
-}
diff --git a/internal/order/http.go b/internal/order/http.go
index 2129717..a4878df 100644
--- a/internal/order/http.go
+++ b/internal/order/http.go
@@ -6,12 +6,12 @@ import (
 	"github.com/ghost-yu/go_shop_second/common"
 	client "github.com/ghost-yu/go_shop_second/common/client/order"
 	"github.com/ghost-yu/go_shop_second/common/consts"
+	"github.com/ghost-yu/go_shop_second/common/convertor"
 	"github.com/ghost-yu/go_shop_second/common/handler/errors"
 	"github.com/ghost-yu/go_shop_second/order/app"
 	"github.com/ghost-yu/go_shop_second/order/app/command"
 	"github.com/ghost-yu/go_shop_second/order/app/dto"
 	"github.com/ghost-yu/go_shop_second/order/app/query"
-	"github.com/ghost-yu/go_shop_second/order/convertor"
 	"github.com/gin-gonic/gin"
 )
 
@@ -70,7 +70,13 @@ func (H HTTPServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerI
 		return
 	}
 
-	resp = convertor.NewOrderConvertor().EntityToClient(o)
+	resp = client.Order{
+		CustomerId:  o.CustomerID,
+		Id:          o.ID,
+		Items:       convertor.NewItemConvertor().EntitiesToClients(o.Items),
+		PaymentLink: o.PaymentLink,
+		Status:      o.Status,
+	}
 }
 
 func (H HTTPServer) validate(req client.CreateOrderRequest) error {
diff --git a/internal/order/ports/grpc.go b/internal/order/ports/grpc.go
index 7863a7f..a98acf4 100644
--- a/internal/order/ports/grpc.go
+++ b/internal/order/ports/grpc.go
@@ -3,11 +3,11 @@ package ports
 import (
 	"context"
 
+	"github.com/ghost-yu/go_shop_second/common/convertor"
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	"github.com/ghost-yu/go_shop_second/order/app"
 	"github.com/ghost-yu/go_shop_second/order/app/command"
 	"github.com/ghost-yu/go_shop_second/order/app/query"
-	"github.com/ghost-yu/go_shop_second/order/convertor"
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
 	"github.com/golang/protobuf/ptypes/empty"
 	"google.golang.org/grpc/codes"
@@ -42,7 +42,14 @@ func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderReque
 	if err != nil {
 		return nil, status.Error(codes.NotFound, err.Error())
 	}
-	return convertor.NewOrderConvertor().EntityToProto(o), nil
+
+	return &orderpb.Order{
+		ID:          o.ID,
+		CustomerID:  o.CustomerID,
+		Status:      o.Status,
+		Items:       convertor.NewItemConvertor().EntitiesToProtos(o.Items),
+		PaymentLink: o.PaymentLink,
+	}, nil
 }
 
 func (G GRPCServer) UpdateOrder(ctx context.Context, request *orderpb.Order) (_ *emptypb.Empty, err error) {
diff --git a/internal/payment/app/command/create_payment.go b/internal/payment/app/command/create_payment.go
index d168300..3474861 100644
--- a/internal/payment/app/command/create_payment.go
+++ b/internal/payment/app/command/create_payment.go
@@ -3,17 +3,16 @@ package command
 import (
 	"context"
 
+	"github.com/ghost-yu/go_shop_second/common/convertor"
 	"github.com/ghost-yu/go_shop_second/common/decorator"
-	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/common/entity"
 	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/ghost-yu/go_shop_second/payment/domain"
 	"github.com/sirupsen/logrus"
 )
 
-// TODO: ACL 清理
-
 type CreatePayment struct {
-	Order *orderpb.Order
+	Order *entity.Order
 }
 
 type CreatePaymentHandler decorator.CommandHandler[CreatePayment, string]
@@ -31,14 +30,14 @@ func (c createPaymentHandler) Handle(ctx context.Context, cmd CreatePayment) (st
 	if err != nil {
 		return "", err
 	}
-	newOrder := &orderpb.Order{
+	newOrder := &entity.Order{
 		ID:          cmd.Order.ID,
 		CustomerID:  cmd.Order.CustomerID,
 		Status:      "waiting_for_payment",
 		Items:       cmd.Order.Items,
 		PaymentLink: link,
 	}
-	err = c.orderGRPC.UpdateOrder(ctx, newOrder)
+	err = c.orderGRPC.UpdateOrder(ctx, convertor.NewOrderConvertor().EntityToProto(newOrder))
 	return link, err
 }
 
diff --git a/internal/payment/domain/payment.go b/internal/payment/domain/payment.go
index 07d56aa..599619c 100644
--- a/internal/payment/domain/payment.go
+++ b/internal/payment/domain/payment.go
@@ -3,11 +3,11 @@ package domain
 import (
 	"context"
 
-	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/common/entity"
 )
 
 type Processor interface {
-	CreatePaymentLink(context.Context, *orderpb.Order) (string, error)
+	CreatePaymentLink(context.Context, *entity.Order) (string, error)
 }
 
 type Order struct {
@@ -15,5 +15,5 @@ type Order struct {
 	CustomerID  string
 	Status      string
 	PaymentLink string
-	Items       []*orderpb.Item
+	Items       []*entity.Item
 }
diff --git a/internal/payment/http.go b/internal/payment/http.go
index 0d4eddf..2cd5cb3 100644
--- a/internal/payment/http.go
+++ b/internal/payment/http.go
@@ -7,9 +7,8 @@ import (
 	"net/http"
 
 	"github.com/ghost-yu/go_shop_second/common/broker"
-	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/common/entity"
 	"github.com/ghost-yu/go_shop_second/common/logging"
-	"github.com/ghost-yu/go_shop_second/payment/domain"
 	"github.com/gin-gonic/gin"
 	"github.com/pkg/errors"
 	amqp "github.com/rabbitmq/amqp091-go"
@@ -72,7 +71,7 @@ func (h *PaymentHandler) handleWebhook(c *gin.Context) {
 		}
 
 		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
-			var items []*orderpb.Item
+			var items []*entity.Item
 			_ = json.Unmarshal([]byte(session.Metadata["items"]), &items)
 
 			tr := otel.Tracer("rabbitmq")
@@ -84,7 +83,7 @@ func (h *PaymentHandler) handleWebhook(c *gin.Context) {
 				Routing:  broker.FanOut,
 				Queue:    "",
 				Exchange: broker.EventOrderPaid,
-				Body: &domain.Order{
+				Body: &entity.Order{
 					ID:          session.Metadata["orderID"],
 					CustomerID:  session.Metadata["customerID"],
 					Status:      string(stripe.CheckoutSessionPaymentStatusPaid),
diff --git a/internal/payment/infrastructure/consumer/consumer.go b/internal/payment/infrastructure/consumer/consumer.go
index 60f9c19..d708551 100644
--- a/internal/payment/infrastructure/consumer/consumer.go
+++ b/internal/payment/infrastructure/consumer/consumer.go
@@ -6,7 +6,7 @@ import (
 	"fmt"
 
 	"github.com/ghost-yu/go_shop_second/common/broker"
-	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/common/entity"
 	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/ghost-yu/go_shop_second/payment/app"
 	"github.com/ghost-yu/go_shop_second/payment/app/command"
@@ -63,7 +63,7 @@ func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Que
 		}
 	}()
 
-	o := &orderpb.Order{}
+	o := &entity.Order{}
 	if err = json.Unmarshal(msg.Body, o); err != nil {
 		err = errors.Wrap(err, "failed to unmarshall msg to order")
 		return
diff --git a/internal/payment/infrastructure/processor/inmem.go b/internal/payment/infrastructure/processor/inmem.go
index 7e8a6a3..7bd44c9 100644
--- a/internal/payment/infrastructure/processor/inmem.go
+++ b/internal/payment/infrastructure/processor/inmem.go
@@ -3,7 +3,7 @@ package processor
 import (
 	"context"
 
-	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/common/entity"
 )
 
 type InmemProcessor struct {
@@ -13,6 +13,6 @@ func NewInmemProcessor() *InmemProcessor {
 	return &InmemProcessor{}
 }
 
-func (i InmemProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
+func (i InmemProcessor) CreatePaymentLink(ctx context.Context, order *entity.Order) (string, error) {
 	return "inmem-payment-link", nil
 }
diff --git a/internal/payment/infrastructure/processor/stripe.go b/internal/payment/infrastructure/processor/stripe.go
index e597fe5..b04c2c5 100644
--- a/internal/payment/infrastructure/processor/stripe.go
+++ b/internal/payment/infrastructure/processor/stripe.go
@@ -5,7 +5,7 @@ import (
 	"encoding/json"
 	"fmt"
 
-	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/common/entity"
 	"github.com/ghost-yu/go_shop_second/common/tracing"
 	"github.com/stripe/stripe-go/v79"
 	"github.com/stripe/stripe-go/v79/checkout/session"
@@ -27,7 +27,7 @@ const (
 	successURL = "http://localhost:8282/success"
 )
 
-func (s StripeProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
+func (s StripeProcessor) CreatePaymentLink(ctx context.Context, order *entity.Order) (string, error) {
 	_, span := tracing.Start(ctx, "stripe_processor.create_payment_link")
 	defer span.End()
 
diff --git a/internal/stock/adapters/stock_inmem_repository.go b/internal/stock/adapters/stock_inmem_repository.go
index 16f48bd..550aab2 100644
--- a/internal/stock/adapters/stock_inmem_repository.go
+++ b/internal/stock/adapters/stock_inmem_repository.go
@@ -4,8 +4,8 @@ import (
 	"context"
 	"sync"
 
+	"github.com/ghost-yu/go_shop_second/common/entity"
 	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
-	"github.com/ghost-yu/go_shop_second/stock/entity"
 )
 
 type MemoryStockRepository struct {
diff --git a/internal/stock/adapters/stock_mysql_repository.go b/internal/stock/adapters/stock_mysql_repository.go
index 867911b..bf5f54d 100644
--- a/internal/stock/adapters/stock_mysql_repository.go
+++ b/internal/stock/adapters/stock_mysql_repository.go
@@ -3,7 +3,7 @@ package adapters
 import (
 	"context"
 
-	"github.com/ghost-yu/go_shop_second/stock/entity"
+	"github.com/ghost-yu/go_shop_second/common/entity"
 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent"
 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent/builder"
 	"github.com/pkg/errors"
diff --git a/internal/stock/adapters/stock_mysql_repository_test.go b/internal/stock/adapters/stock_mysql_repository_test.go
index fa77ed0..8991b9e 100644
--- a/internal/stock/adapters/stock_mysql_repository_test.go
+++ b/internal/stock/adapters/stock_mysql_repository_test.go
@@ -8,7 +8,7 @@ import (
 	"time"
 
 	_ "github.com/ghost-yu/go_shop_second/common/config"
-	"github.com/ghost-yu/go_shop_second/stock/entity"
+	"github.com/ghost-yu/go_shop_second/common/entity"
 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent"
 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent/builder"
 	"github.com/spf13/viper"
diff --git a/internal/stock/app/query/check_if_items_in_stock.go b/internal/stock/app/query/check_if_items_in_stock.go
index c61aeed..ea898bc 100644
--- a/internal/stock/app/query/check_if_items_in_stock.go
+++ b/internal/stock/app/query/check_if_items_in_stock.go
@@ -6,10 +6,10 @@ import (
 	"time"
 
 	"github.com/ghost-yu/go_shop_second/common/decorator"
+	"github.com/ghost-yu/go_shop_second/common/entity"
 	"github.com/ghost-yu/go_shop_second/common/handler/redis"
 	"github.com/ghost-yu/go_shop_second/common/logging"
 	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
-	"github.com/ghost-yu/go_shop_second/stock/entity"
 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/integration"
 	"github.com/pkg/errors"
 	"github.com/sirupsen/logrus"
diff --git a/internal/stock/app/query/get_items.go b/internal/stock/app/query/get_items.go
index 4de7bd5..67f9ae2 100644
--- a/internal/stock/app/query/get_items.go
+++ b/internal/stock/app/query/get_items.go
@@ -4,8 +4,8 @@ import (
 	"context"
 
 	"github.com/ghost-yu/go_shop_second/common/decorator"
+	"github.com/ghost-yu/go_shop_second/common/entity"
 	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
-	"github.com/ghost-yu/go_shop_second/stock/entity"
 	"github.com/sirupsen/logrus"
 )
 
diff --git a/internal/stock/convertor/convertor.go b/internal/stock/convertor/convertor.go
deleted file mode 100644
index be35f65..0000000
--- a/internal/stock/convertor/convertor.go
+++ /dev/null
@@ -1,97 +0,0 @@
-package convertor
-
-import (
-	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
-	"github.com/ghost-yu/go_shop_second/stock/entity"
-)
-
-type OrderConvertor struct{}
-type ItemConvertor struct{}
-type ItemWithQuantityConvertor struct{}
-
-func (c *ItemWithQuantityConvertor) EntitiesToProtos(items []*entity.ItemWithQuantity) (res []*orderpb.ItemWithQuantity) {
-	for _, i := range items {
-		res = append(res, c.EntityToProto(i))
-	}
-	return
-}
-
-func (c *ItemWithQuantityConvertor) EntityToProto(i *entity.ItemWithQuantity) *orderpb.ItemWithQuantity {
-	return &orderpb.ItemWithQuantity{
-		ID:       i.ID,
-		Quantity: i.Quantity,
-	}
-}
-
-func (c *ItemWithQuantityConvertor) ProtosToEntities(items []*orderpb.ItemWithQuantity) (res []*entity.ItemWithQuantity) {
-	for _, i := range items {
-		res = append(res, c.ProtoToEntity(i))
-	}
-	return
-}
-
-func (c *ItemWithQuantityConvertor) ProtoToEntity(i *orderpb.ItemWithQuantity) *entity.ItemWithQuantity {
-	return &entity.ItemWithQuantity{
-		ID:       i.ID,
-		Quantity: i.Quantity,
-	}
-}
-
-func (c *OrderConvertor) EntityToProto(o *entity.Order) *orderpb.Order {
-	c.check(o)
-	return &orderpb.Order{
-		ID:          o.ID,
-		CustomerID:  o.CustomerID,
-		Status:      o.Status,
-		Items:       NewItemConvertor().EntitiesToProtos(o.Items),
-		PaymentLink: o.PaymentLink,
-	}
-}
-
-func (c *OrderConvertor) ProtoToEntity(o *orderpb.Order) *entity.Order {
-	c.check(o)
-	return &entity.Order{
-		ID:          o.ID,
-		CustomerID:  o.CustomerID,
-		Status:      o.Status,
-		PaymentLink: o.PaymentLink,
-		Items:       NewItemConvertor().ProtosToEntities(o.Items),
-	}
-}
-func (c *OrderConvertor) check(o interface{}) {
-	if o == nil {
-		panic("connot convert nil order")
-	}
-}
-
-func (c *ItemConvertor) EntitiesToProtos(items []*entity.Item) (res []*orderpb.Item) {
-	for _, i := range items {
-		res = append(res, c.EntityToProto(i))
-	}
-	return
-}
-
-func (c *ItemConvertor) ProtosToEntities(items []*orderpb.Item) (res []*entity.Item) {
-	for _, i := range items {
-		res = append(res, c.ProtoToEntity(i))
-	}
-	return
-}
-
-func (c *ItemConvertor) EntityToProto(i *entity.Item) *orderpb.Item {
-	return &orderpb.Item{
-		ID:       i.ID,
-		Name:     i.Name,
-		Quantity: i.Quantity,
-		PriceID:  i.PriceID,
-	}
-}
-
-func (c *ItemConvertor) ProtoToEntity(i *orderpb.Item) *entity.Item {
-	return &entity.Item{
-		ID:       i.ID,
-		Name:     i.Name,
-		Quantity: i.Quantity,
-		PriceID:  i.PriceID,
-	}
-}
diff --git a/internal/stock/convertor/facade.go b/internal/stock/convertor/facade.go
deleted file mode 100644
index 63ad8d0..0000000
--- a/internal/stock/convertor/facade.go
+++ /dev/null
@@ -1,39 +0,0 @@
-package convertor
-
-import "sync"
-
-var (
-	orderConvertor *OrderConvertor
-	orderOnce      sync.Once
-)
-
-var (
-	itemConvertor *ItemConvertor
-	itemOnce      sync.Once
-)
-
-var (
-	itemWithQuantityConvertor *ItemWithQuantityConvertor
-	itemWithQuantityOnce      sync.Once
-)
-
-func NewOrderConvertor() *OrderConvertor {
-	orderOnce.Do(func() {
-		orderConvertor = new(OrderConvertor)
-	})
-	return orderConvertor
-}
-
-func NewItemConvertor() *ItemConvertor {
-	itemOnce.Do(func() {
-		itemConvertor = new(ItemConvertor)
-	})
-	return itemConvertor
-}
-
-func NewItemWithQuantityConvertor() *ItemWithQuantityConvertor {
-	itemWithQuantityOnce.Do(func() {
-		itemWithQuantityConvertor = new(ItemWithQuantityConvertor)
-	})
-	return itemWithQuantityConvertor
-}
diff --git a/internal/stock/domain/stock/repository.go b/internal/stock/domain/stock/repository.go
index ea9d7c3..1165b6b 100644
--- a/internal/stock/domain/stock/repository.go
+++ b/internal/stock/domain/stock/repository.go
@@ -5,7 +5,7 @@ import (
 	"fmt"
 	"strings"
 
-	"github.com/ghost-yu/go_shop_second/stock/entity"
+	"github.com/ghost-yu/go_shop_second/common/entity"
 )
 
 type Repository interface {
diff --git a/internal/stock/ports/grpc.go b/internal/stock/ports/grpc.go
index 13f212e..5333b80 100644
--- a/internal/stock/ports/grpc.go
+++ b/internal/stock/ports/grpc.go
@@ -3,11 +3,11 @@ package ports
 import (
 	"context"
 
+	"github.com/ghost-yu/go_shop_second/common/convertor"
 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
 	"github.com/ghost-yu/go_shop_second/common/tracing"
 	"github.com/ghost-yu/go_shop_second/stock/app"
 	"github.com/ghost-yu/go_shop_second/stock/app/query"
-	"github.com/ghost-yu/go_shop_second/stock/convertor"
 	"google.golang.org/grpc/codes"
 	"google.golang.org/grpc/status"
 )
~~~

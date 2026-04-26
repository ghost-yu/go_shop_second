# Lesson Pair Diff Report

- FromBranch: lesson21
- ToBranch: lesson22

## Short Summary

~~~text
 5 files changed, 99 insertions(+)
~~~

## File Stats

~~~text
 internal/order/domain/order/order.go               | 10 ++++
 internal/order/go.mod                              |  1 +
 internal/order/go.sum                              |  4 ++
 internal/order/infrastructure/consumer/consumer.go | 70 ++++++++++++++++++++++
 internal/order/main.go                             | 14 +++++
 5 files changed, 99 insertions(+)
~~~

## Commit Comparison

~~~text
> c49c498 order consume event.paid
~~~

## Changed Files

~~~text
internal/order/domain/order/order.go
internal/order/go.mod
internal/order/go.sum
internal/order/infrastructure/consumer/consumer.go
internal/order/main.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
internal/order/domain/order/order.go
internal/order/infrastructure/consumer/consumer.go
internal/order/main.go
~~~

## Full Diff

~~~diff
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
 
 type Order struct {
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
diff --git a/internal/order/go.mod b/internal/order/go.mod
index 73b9119..99f2d66 100644
--- a/internal/order/go.mod
+++ b/internal/order/go.mod
@@ -13,6 +13,7 @@ require (
 	github.com/rabbitmq/amqp091-go v1.10.0
 	github.com/sirupsen/logrus v1.9.3
 	github.com/spf13/viper v1.19.0
+	github.com/stripe/stripe-go/v80 v80.2.0
 	google.golang.org/grpc v1.67.1
 	google.golang.org/protobuf v1.35.1
 )
diff --git a/internal/order/go.sum b/internal/order/go.sum
index 213b049..9f6846a 100644
--- a/internal/order/go.sum
+++ b/internal/order/go.sum
@@ -282,6 +282,8 @@ github.com/stretchr/testify v1.8.0/go.mod h1:yNjHg4UonilssWZ8iaSj1OCr/vHnekPRkoO
 github.com/stretchr/testify v1.8.1/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o6fzry7u4=
 github.com/stretchr/testify v1.9.0 h1:HtqpIVDClZ4nwg75+f6Lvsy/wHu+3BoSGCbBAcpTsTg=
 github.com/stretchr/testify v1.9.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8C91i36aY=
+github.com/stripe/stripe-go/v80 v80.2.0 h1:rCl1PyIAG+gi7tj9prOuWt6XNjKK0BjMoZSvtdiQwUc=
+github.com/stripe/stripe-go/v80 v80.2.0/go.mod h1:n7tsDvdltYlzOLGXlseMSJM6ik5uv3guptqtae/VSak=
 github.com/subosito/gotenv v1.6.0 h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=
 github.com/subosito/gotenv v1.6.0/go.mod h1:Dk4QP5c2W3ibzajGcXpNraDfq2IrhjMIvMSWPKKo0FU=
 github.com/tv42/httpunix v0.0.0-20150427012821-b75d8614f926/go.mod h1:9ESjWnEqriFuLhtthL60Sar/7RFoluCcXsuvEwTV5KM=
@@ -332,6 +334,7 @@ golang.org/x/net v0.0.0-20200226121028-0de0cce0169b/go.mod h1:z5CRVTTTmAJ677TzLL
 golang.org/x/net v0.0.0-20201021035429-f5854403a974/go.mod h1:sp8m0HH+o8qH0wwXwYZr8TS3Oi6o0r6Gce1SSxlDquU=
 golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
 golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1/go.mod h1:9tjilg8BloeKEkVJvy7fQ90B1CfIiPueXVOjqfkSzI8=
+golang.org/x/net v0.0.0-20210520170846-37e1c6afe023/go.mod h1:9nx3DQGgdP8bBQD5qxJ1jj9UTztislL4KSBs9R2vV5Y=
 golang.org/x/net v0.30.0 h1:AcW1SDZMkb8IpzCdQUaIq2sP4sZ4zw+55h6ynffypl4=
 golang.org/x/net v0.30.0/go.mod h1:2wGyMJ5iFasEhkwi13ChkO/t1ECNC4X4eBKkVFyYFlU=
 golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be/go.mod h1:N/0e6XlmueqKjAGxoOufVs8QHGRruUQn6yWY3a++T0U=
@@ -360,6 +363,7 @@ golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f/go.mod h1:h1NjWce9XRLGQEsW7w
 golang.org/x/sys v0.0.0-20201119102817-f84b799fce68/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210303074136-134d130e1a04/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210330210617-4fbd30eecc44/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20210423082822-04245dca01da/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
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
 	"github.com/ghost-yu/go_shop_second/order/ports"
 	"github.com/ghost-yu/go_shop_second/order/service"
 	"github.com/gin-gonic/gin"
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
 		svc := ports.NewGRPCServer(application)
 		orderpb.RegisterOrderServiceServer(server, svc)
~~~

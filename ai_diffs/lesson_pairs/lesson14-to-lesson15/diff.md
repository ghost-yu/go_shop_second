# Lesson Pair Diff Report

- FromBranch: lesson14
- ToBranch: lesson15

## Short Summary

~~~text
 16 files changed, 154 insertions(+), 7 deletions(-)
~~~

## File Stats

~~~text
 docker-compose.yml                         |  8 ++++++-
 internal/common/broker/event.go            |  6 +++++
 internal/common/broker/rabbitmq.go         | 29 ++++++++++++++++++++++++
 internal/common/client/grpc.go             | 10 ++++++++-
 internal/common/config/global.yaml         |  9 ++++++++
 internal/common/discovery/grpc.go          | 19 ++++++++++++++++
 internal/common/go.mod                     |  1 +
 internal/common/go.sum                     |  2 ++
 internal/order/app/command/create_order.go | 36 +++++++++++++++++++++++++++++-
 internal/order/go.mod                      |  1 +
 internal/order/go.sum                      |  4 ++++
 internal/order/service/application.go      | 17 +++++++++++---
 internal/payment/go.mod                    |  1 +
 internal/payment/go.sum                    |  4 ++++
 internal/payment/http.go                   |  2 +-
 internal/payment/main.go                   | 12 ++++++++++
 16 files changed, 154 insertions(+), 7 deletions(-)
~~~

## Commit Comparison

~~~text
> a0deb6f publish order create event
~~~

## Changed Files

~~~text
docker-compose.yml
internal/common/broker/event.go
internal/common/broker/rabbitmq.go
internal/common/client/grpc.go
internal/common/config/global.yaml
internal/common/discovery/grpc.go
internal/common/go.mod
internal/common/go.sum
internal/order/app/command/create_order.go
internal/order/go.mod
internal/order/go.sum
internal/order/service/application.go
internal/payment/go.mod
internal/payment/go.sum
internal/payment/http.go
internal/payment/main.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
docker-compose.yml
internal/common/broker/event.go
internal/common/broker/rabbitmq.go
internal/common/client/grpc.go
internal/common/config/global.yaml
internal/common/discovery/grpc.go
internal/order/app/command/create_order.go
internal/order/service/application.go
internal/payment/http.go
internal/payment/main.go
~~~

## Full Diff

~~~diff
diff --git a/docker-compose.yml b/docker-compose.yml
index 079748a..06beeb4 100644
--- a/docker-compose.yml
+++ b/docker-compose.yml
@@ -5,4 +5,10 @@ services:
     command: agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0
     ports:
       - 8500:8500
-      - 8600:8600/udp
\ No newline at end of file
+      - 8600:8600/udp
+
+  rabbit-mq:
+    image: "rabbitmq:3-management"
+    ports:
+      - "15672:15672"
+      - "5672:5672"
\ No newline at end of file
diff --git a/internal/common/broker/event.go b/internal/common/broker/event.go
new file mode 100644
index 0000000..e53531c
--- /dev/null
+++ b/internal/common/broker/event.go
@@ -0,0 +1,6 @@
+package broker
+
+const (
+	EventOrderCreated = "order.created"
+	EventOrderPaid    = "order.paid"
+)
diff --git a/internal/common/broker/rabbitmq.go b/internal/common/broker/rabbitmq.go
new file mode 100644
index 0000000..599af4e
--- /dev/null
+++ b/internal/common/broker/rabbitmq.go
@@ -0,0 +1,29 @@
+package broker
+
+import (
+	"fmt"
+
+	amqp "github.com/rabbitmq/amqp091-go"
+	"github.com/sirupsen/logrus"
+)
+
+func Connect(user, password, host, port string) (*amqp.Channel, func() error) {
+	address := fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port)
+	conn, err := amqp.Dial(address)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	ch, err := conn.Channel()
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	err = ch.ExchangeDeclare(EventOrderCreated, "direct", true, false, false, false, nil)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	err = ch.ExchangeDeclare(EventOrderPaid, "fanout", true, false, false, false, nil)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	return ch, conn.Close
+}
diff --git a/internal/common/client/grpc.go b/internal/common/client/grpc.go
index f9f514a..323329f 100644
--- a/internal/common/client/grpc.go
+++ b/internal/common/client/grpc.go
@@ -3,14 +3,22 @@ package client
 import (
 	"context"
 
+	"github.com/ghost-yu/go_shop_second/common/discovery"
 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
+	"github.com/sirupsen/logrus"
 	"github.com/spf13/viper"
 	"google.golang.org/grpc"
 	"google.golang.org/grpc/credentials/insecure"
 )
 
 func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
-	grpcAddr := viper.GetString("stock.grpc-addr")
+	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("stock.service-name"))
+	if err != nil {
+		return nil, func() error { return nil }, err
+	}
+	if grpcAddr == "" {
+		logrus.Warn("empty grpc addr for stock grpc")
+	}
 	opts, err := grpcDialOpts(grpcAddr)
 	if err != nil {
 		return nil, func() error { return nil }, err
diff --git a/internal/common/config/global.yaml b/internal/common/config/global.yaml
index ad5fe69..17a4f41 100644
--- a/internal/common/config/global.yaml
+++ b/internal/common/config/global.yaml
@@ -1,5 +1,8 @@
 fallback-grpc-addr: 127.0.0.1:3030
 
+consul:
+  addr: 127.0.0.1:8500
+
 order:
   service-name: order
   server-to-run: http
@@ -18,4 +21,10 @@ payment:
   http-addr: 127.0.0.1:8284
   grpc-addr: 127.0.0.1:5004
 
+rabbitmq:
+  user: guest
+  password: guest
+  host: 127.0.0.1
+  port: 5672
+
 stripe-key: "${STRIPE_KEY}"
\ No newline at end of file
diff --git a/internal/common/discovery/grpc.go b/internal/common/discovery/grpc.go
index cda522a..f404a0b 100644
--- a/internal/common/discovery/grpc.go
+++ b/internal/common/discovery/grpc.go
@@ -2,6 +2,8 @@ package discovery
 
 import (
 	"context"
+	"fmt"
+	"math/rand"
 	"time"
 
 	"github.com/ghost-yu/go_shop_second/common/discovery/consul"
@@ -35,3 +37,20 @@ func RegisterToConsul(ctx context.Context, serviceName string) (func() error, er
 		return registry.Deregister(ctx, instanceID, serviceName)
 	}, nil
 }
+
+func GetServiceAddr(ctx context.Context, serviceName string) (string, error) {
+	registry, err := consul.New(viper.GetString("consul.addr"))
+	if err != nil {
+		return "", err
+	}
+	addrs, err := registry.Discover(ctx, serviceName)
+	if err != nil {
+		return "", err
+	}
+	if len(addrs) == 0 {
+		return "", fmt.Errorf("got empty %s addrs from consul", serviceName)
+	}
+	i := rand.Intn(len(addrs))
+	logrus.Infof("Discovered %d instance of %s, addrs=%v", len(addrs), serviceName, addrs)
+	return addrs[i], nil
+}
diff --git a/internal/common/go.mod b/internal/common/go.mod
index 2df82af..af1aa77 100644
--- a/internal/common/go.mod
+++ b/internal/common/go.mod
@@ -53,6 +53,7 @@ require (
 	github.com/onsi/ginkgo v1.16.5 // indirect
 	github.com/onsi/gomega v1.34.2 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
+	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
 	github.com/sourcegraph/conc v0.3.0 // indirect
diff --git a/internal/common/go.sum b/internal/common/go.sum
index 395888c..05ee4a2 100644
--- a/internal/common/go.sum
+++ b/internal/common/go.sum
@@ -259,6 +259,8 @@ github.com/prometheus/common v0.9.1/go.mod h1:yhUN8i9wzaXS3w1O07YhxHEBxD+W35wd8b
 github.com/prometheus/procfs v0.0.0-20181005140218-185b4288413d/go.mod h1:c3At6R/oaqEKCNdg8wHV1ftS6bRYblBhIjjI8uT2IGk=
 github.com/prometheus/procfs v0.0.2/go.mod h1:TjEm7ze935MbeOT/UhFTIMYKhuLP4wbCsTZCD3I8kEA=
 github.com/prometheus/procfs v0.0.8/go.mod h1:7Qr8sr6344vo1JqZ6HhLceV9o3AJ1Ff+GxbHq6oeK9A=
+github.com/rabbitmq/amqp091-go v1.10.0 h1:STpn5XsHlHGcecLmMFCtg7mqq0RnD+zFr4uzukfVhBw=
+github.com/rabbitmq/amqp091-go v1.10.0/go.mod h1:Hy4jKW5kQART1u+JkDTF9YYOQUHXqMuhrgxOEeS7G4o=
 github.com/rogpeppe/go-internal v1.9.0 h1:73kH8U+JUqXU8lRuOHeVHaa/SZPifC7BkcraZVejAe8=
 github.com/rogpeppe/go-internal v1.9.0/go.mod h1:WtVeX8xhTBvf0smdhujwtBcq4Qrzq/fJaraNFVN+nFs=
 github.com/ryanuber/columnize v0.0.0-20160712163229-9b3edd62028f/go.mod h1:sm1tb6uqfes/u+d4ooFouqFdy9/2g9QGwK3SQygK0Ts=
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
index 1c72b04..7185157 100644
--- a/internal/order/app/command/create_order.go
+++ b/internal/order/app/command/create_order.go
@@ -2,12 +2,15 @@ package command
 
 import (
 	"context"
+	"encoding/json"
 	"errors"
 
+	"github.com/ghost-yu/go_shop_second/common/broker"
 	"github.com/ghost-yu/go_shop_second/common/decorator"
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	"github.com/ghost-yu/go_shop_second/order/app/query"
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
+	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
 )
 
@@ -25,19 +28,31 @@ type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult
 type createOrderHandler struct {
 	orderRepo domain.Repository
 	stockGRPC query.StockService
+	channel   *amqp.Channel
 }
 
 func NewCreateOrderHandler(
 	orderRepo domain.Repository,
 	stockGRPC query.StockService,
+	channel *amqp.Channel,
 	logger *logrus.Entry,
 	metricClient decorator.MetricsClient,
 ) CreateOrderHandler {
 	if orderRepo == nil {
 		panic("nil orderRepo")
 	}
+	if stockGRPC == nil {
+		panic("nil stockGRPC")
+	}
+	if channel == nil {
+		panic("nil channel ")
+	}
 	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
-		createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC},
+		createOrderHandler{
+			orderRepo: orderRepo,
+			stockGRPC: stockGRPC,
+			channel:   channel,
+		},
 		logger,
 		metricClient,
 	)
@@ -55,6 +70,25 @@ func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*Creat
 	if err != nil {
 		return nil, err
 	}
+
+	q, err := c.channel.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
+	if err != nil {
+		return nil, err
+	}
+
+	marshalledOrder, err := json.Marshal(o)
+	if err != nil {
+		return nil, err
+	}
+	err = c.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
+		ContentType:  "application/json",
+		DeliveryMode: amqp.Persistent,
+		Body:         marshalledOrder,
+	})
+	if err != nil {
+		return nil, err
+	}
+
 	return &CreateOrderResult{OrderID: o.ID}, nil
 }
 
diff --git a/internal/order/go.mod b/internal/order/go.mod
index e5c9d5c..5478ff8 100644
--- a/internal/order/go.mod
+++ b/internal/order/go.mod
@@ -55,6 +55,7 @@ require (
 	github.com/modern-go/reflect2 v1.0.2 // indirect
 	github.com/nxadm/tail v1.4.11 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
+	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
 	github.com/sagikazarmark/locafero v0.6.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
 	github.com/sourcegraph/conc v0.3.0 // indirect
diff --git a/internal/order/go.sum b/internal/order/go.sum
index 42065e9..213b049 100644
--- a/internal/order/go.sum
+++ b/internal/order/go.sum
@@ -241,6 +241,8 @@ github.com/prometheus/common v0.9.1/go.mod h1:yhUN8i9wzaXS3w1O07YhxHEBxD+W35wd8b
 github.com/prometheus/procfs v0.0.0-20181005140218-185b4288413d/go.mod h1:c3At6R/oaqEKCNdg8wHV1ftS6bRYblBhIjjI8uT2IGk=
 github.com/prometheus/procfs v0.0.2/go.mod h1:TjEm7ze935MbeOT/UhFTIMYKhuLP4wbCsTZCD3I8kEA=
 github.com/prometheus/procfs v0.0.8/go.mod h1:7Qr8sr6344vo1JqZ6HhLceV9o3AJ1Ff+GxbHq6oeK9A=
+github.com/rabbitmq/amqp091-go v1.10.0 h1:STpn5XsHlHGcecLmMFCtg7mqq0RnD+zFr4uzukfVhBw=
+github.com/rabbitmq/amqp091-go v1.10.0/go.mod h1:Hy4jKW5kQART1u+JkDTF9YYOQUHXqMuhrgxOEeS7G4o=
 github.com/rogpeppe/go-internal v1.9.0 h1:73kH8U+JUqXU8lRuOHeVHaa/SZPifC7BkcraZVejAe8=
 github.com/rogpeppe/go-internal v1.9.0/go.mod h1:WtVeX8xhTBvf0smdhujwtBcq4Qrzq/fJaraNFVN+nFs=
 github.com/ryanuber/columnize v0.0.0-20160712163229-9b3edd62028f/go.mod h1:sm1tb6uqfes/u+d4ooFouqFdy9/2g9QGwK3SQygK0Ts=
@@ -293,6 +295,8 @@ github.com/yuin/goldmark v1.1.27/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9de
 github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 go.uber.org/atomic v1.7.0/go.mod h1:fEN4uk6kAWBTFdckzkM89CLk9XfWZrxpCo0nPH17wJc=
 go.uber.org/goleak v1.1.10/go.mod h1:8a7PlsEVH3e/a/GLqe5IIrQx6GzcnRmZEufDUTk4A7A=
+go.uber.org/goleak v1.3.0 h1:2K3zAYmnTNqV73imy9J1T3WC+gmCePx2hEGkimedGto=
+go.uber.org/goleak v1.3.0/go.mod h1:CoHD4mav9JJNrW/WLlf7HGZPjdw8EucARQHekz1X6bE=
 go.uber.org/multierr v1.6.0/go.mod h1:cdWPpRnG4AhwMwsgIHip0KRBQjJy5kYEpYjJxpXp9iU=
 go.uber.org/multierr v1.11.0 h1:blXXJkSxSSfBVBlC76pxqeO+LN3aDfLQo+309xJstO0=
 go.uber.org/multierr v1.11.0/go.mod h1:20+QtiLqy0Nd6FdQB9TLXag12DsQkrbs3htMFfDN80Y=
diff --git a/internal/order/service/application.go b/internal/order/service/application.go
index 648e438..8de0e53 100644
--- a/internal/order/service/application.go
+++ b/internal/order/service/application.go
@@ -3,6 +3,7 @@ package service
 import (
 	"context"
 
+	"github.com/ghost-yu/go_shop_second/common/broker"
 	grpcClient "github.com/ghost-yu/go_shop_second/common/client"
 	"github.com/ghost-yu/go_shop_second/common/metrics"
 	"github.com/ghost-yu/go_shop_second/order/adapters"
@@ -10,7 +11,9 @@ import (
 	"github.com/ghost-yu/go_shop_second/order/app"
 	"github.com/ghost-yu/go_shop_second/order/app/command"
 	"github.com/ghost-yu/go_shop_second/order/app/query"
+	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
+	"github.com/spf13/viper"
 )
 
 func NewApplication(ctx context.Context) (app.Application, func()) {
@@ -18,19 +21,27 @@ func NewApplication(ctx context.Context) (app.Application, func()) {
 	if err != nil {
 		panic(err)
 	}
+	ch, closeCh := broker.Connect(
+		viper.GetString("rabbitmq.user"),
+		viper.GetString("rabbitmq.password"),
+		viper.GetString("rabbitmq.host"),
+		viper.GetString("rabbitmq.port"),
+	)
 	stockGRPC := grpc.NewStockGRPC(stockClient)
-	return newApplication(ctx, stockGRPC), func() {
+	return newApplication(ctx, stockGRPC, ch), func() {
 		_ = closeStockClient()
+		_ = closeCh()
+		_ = ch.Close()
 	}
 }
 
-func newApplication(_ context.Context, stockGRPC query.StockService) app.Application {
+func newApplication(_ context.Context, stockGRPC query.StockService, ch *amqp.Channel) app.Application {
 	orderRepo := adapters.NewMemoryOrderRepository()
 	logger := logrus.NewEntry(logrus.StandardLogger())
 	metricClient := metrics.TodoMetrics{}
 	return app.Application{
 		Commands: app.Commands{
-			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, logger, metricClient),
+			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, ch, logger, metricClient),
 			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
 		},
 		Queries: app.Queries{
diff --git a/internal/payment/go.mod b/internal/payment/go.mod
index 320311e..3d977b2 100644
--- a/internal/payment/go.mod
+++ b/internal/payment/go.mod
@@ -37,6 +37,7 @@ require (
 	github.com/modern-go/reflect2 v1.0.2 // indirect
 	github.com/nxadm/tail v1.4.11 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
+	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
 	github.com/sourcegraph/conc v0.3.0 // indirect
diff --git a/internal/payment/go.sum b/internal/payment/go.sum
index dbe94ae..d754f6e 100644
--- a/internal/payment/go.sum
+++ b/internal/payment/go.sum
@@ -113,6 +113,8 @@ github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZN
 github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 h1:Jamvg5psRIccs7FGNTlIRMkT8wgtp5eCXdBlqhYGL6U=
 github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
 github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
+github.com/rabbitmq/amqp091-go v1.10.0 h1:STpn5XsHlHGcecLmMFCtg7mqq0RnD+zFr4uzukfVhBw=
+github.com/rabbitmq/amqp091-go v1.10.0/go.mod h1:Hy4jKW5kQART1u+JkDTF9YYOQUHXqMuhrgxOEeS7G4o=
 github.com/rogpeppe/go-internal v1.9.0 h1:73kH8U+JUqXU8lRuOHeVHaa/SZPifC7BkcraZVejAe8=
 github.com/rogpeppe/go-internal v1.9.0/go.mod h1:WtVeX8xhTBvf0smdhujwtBcq4Qrzq/fJaraNFVN+nFs=
 github.com/sagikazarmark/locafero v0.4.0 h1:HApY1R9zGo4DBgr7dqsTH/JJxLTTsOt7u6keLGt6kNQ=
@@ -162,6 +164,8 @@ go.uber.org/atomic v1.7.0/go.mod h1:fEN4uk6kAWBTFdckzkM89CLk9XfWZrxpCo0nPH17wJc=
 go.uber.org/atomic v1.9.0 h1:ECmE8Bn/WFTYwEW/bpKD3M8VtR/zQVbavAoalC1PYyE=
 go.uber.org/atomic v1.9.0/go.mod h1:fEN4uk6kAWBTFdckzkM89CLk9XfWZrxpCo0nPH17wJc=
 go.uber.org/goleak v1.1.10/go.mod h1:8a7PlsEVH3e/a/GLqe5IIrQx6GzcnRmZEufDUTk4A7A=
+go.uber.org/goleak v1.3.0 h1:2K3zAYmnTNqV73imy9J1T3WC+gmCePx2hEGkimedGto=
+go.uber.org/goleak v1.3.0/go.mod h1:CoHD4mav9JJNrW/WLlf7HGZPjdw8EucARQHekz1X6bE=
 go.uber.org/multierr v1.6.0/go.mod h1:cdWPpRnG4AhwMwsgIHip0KRBQjJy5kYEpYjJxpXp9iU=
 go.uber.org/multierr v1.9.0 h1:7fIwc/ZtS0q++VgcfqFDxSBZVv/Xo49/SYnDFupUwlI=
 go.uber.org/multierr v1.9.0/go.mod h1:X2jQV1h+kxSjClGpnseKVIxpmcjrj7MNnI0bnlfKTVQ=
diff --git a/internal/payment/http.go b/internal/payment/http.go
index 6ea6825..b98bbe3 100644
--- a/internal/payment/http.go
+++ b/internal/payment/http.go
@@ -17,5 +17,5 @@ func (h *PaymentHandler) RegisterRoutes(c *gin.Engine) {
 }
 
 func (h *PaymentHandler) handleWebhook(c *gin.Context) {
-	logrus.Info("Got webhook from stripe")
+	logrus.Info("receive webhook from stripe")
 }
diff --git a/internal/payment/main.go b/internal/payment/main.go
index 23758d4..b8d159b 100644
--- a/internal/payment/main.go
+++ b/internal/payment/main.go
@@ -1,6 +1,7 @@
 package main
 
 import (
+	"github.com/ghost-yu/go_shop_second/common/broker"
 	"github.com/ghost-yu/go_shop_second/common/config"
 	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/ghost-yu/go_shop_second/common/server"
@@ -18,6 +19,17 @@ func init() {
 func main() {
 	serverType := viper.GetString("payment.server-to-run")
 
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
 	paymentHandler := NewPaymentHandler()
 	switch serverType {
 	case "http":
~~~

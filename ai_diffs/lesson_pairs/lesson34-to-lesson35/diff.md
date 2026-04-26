# Lesson Pair Diff Report

- FromBranch: lesson34
- ToBranch: lesson35

## Short Summary

~~~text
 15 files changed, 349 insertions(+), 38 deletions(-)
~~~

## File Stats

~~~text
 docker-compose.yml                                 |  22 ++-
 internal/common/broker/rabbitmq.go                 |   6 +-
 internal/common/config/global.yaml                 |  10 ++
 internal/common/config/viper.go                    |  39 ++++-
 internal/order/adapters/order_mongo_repository.go  | 158 +++++++++++++++++++++
 internal/order/app/command/create_order.go         |   9 +-
 internal/order/domain/order/order.go               |  14 ++
 internal/order/go.mod                              |  10 +-
 internal/order/go.sum                              |  32 ++++-
 internal/order/infrastructure/consumer/consumer.go |  21 ++-
 internal/order/main.go                             |   5 +-
 internal/order/service/application.go              |  29 +++-
 .../payment/infrastructure/consumer/consumer.go    |  22 ++-
 internal/payment/main.go                           |   5 +-
 internal/stock/main.go                             |   5 +-
 15 files changed, 349 insertions(+), 38 deletions(-)
~~~

## Commit Comparison

~~~text
> 95d04c4 mongo
> ede4e45 viper fix
~~~

## Changed Files

~~~text
docker-compose.yml
internal/common/broker/rabbitmq.go
internal/common/config/global.yaml
internal/common/config/viper.go
internal/order/adapters/order_mongo_repository.go
internal/order/app/command/create_order.go
internal/order/domain/order/order.go
internal/order/go.mod
internal/order/go.sum
internal/order/infrastructure/consumer/consumer.go
internal/order/main.go
internal/order/service/application.go
internal/payment/infrastructure/consumer/consumer.go
internal/payment/main.go
internal/stock/main.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
docker-compose.yml
internal/common/broker/rabbitmq.go
internal/common/config/global.yaml
internal/common/config/viper.go
internal/order/adapters/order_mongo_repository.go
internal/order/app/command/create_order.go
internal/order/domain/order/order.go
internal/order/infrastructure/consumer/consumer.go
internal/order/main.go
internal/order/service/application.go
internal/payment/infrastructure/consumer/consumer.go
internal/payment/main.go
internal/stock/main.go
~~~

## Full Diff

~~~diff
diff --git a/docker-compose.yml b/docker-compose.yml
index 12f1e17..ac2b22f 100644
--- a/docker-compose.yml
+++ b/docker-compose.yml
@@ -22,4 +22,24 @@ services:
       - "4318:4318"
       - "4317:4317"
     environment:
-      COLLECTOR_OTLP_ENABLED: true
\ No newline at end of file
+      COLLECTOR_OTLP_ENABLED: true
+
+  order-mongo:
+    image: "mongo:7.0.8"
+    restart: always
+    environment:
+      MONGO_INITDB_ROOT_USERNAME: root
+      MONGO_INITDB_ROOT_PASSWORD: password
+    ports:
+      - "27017:27017"
+
+  mongo-express:
+    image: "mongo-express"
+    restart: always
+    ports:
+      - "8082:8081"
+    environment:
+      ME_CONFIG_MONGODB_ADMINUSERNAME: root
+      ME_CONFIG_MONGODB_ADMINPASSWORD: password
+      ME_CONFIG_MONGODB_URL: mongodb://root:password@order-mongo:27017/
+      ME_CONFIG_BASICAUTH: false
\ No newline at end of file
diff --git a/internal/common/broker/rabbitmq.go b/internal/common/broker/rabbitmq.go
index d8a3e84..ef74e7d 100644
--- a/internal/common/broker/rabbitmq.go
+++ b/internal/common/broker/rabbitmq.go
@@ -5,8 +5,10 @@ import (
 	"fmt"
 	"time"
 
+	_ "github.com/ghost-yu/go_shop_second/common/config"
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
+	"github.com/spf13/viper"
 	"go.opentelemetry.io/otel"
 )
 
@@ -17,8 +19,7 @@ const (
 )
 
 var (
-	//maxRetryCount = viper.GetInt64("rabbitmq.max-retry")
-	maxRetryCount int64 = 3
+	maxRetryCount = viper.GetInt64("rabbitmq.max-retry")
 )
 
 func Connect(user, password, host, port string) (*amqp.Channel, func() error) {
@@ -63,6 +64,7 @@ func createDLX(ch *amqp.Channel) error {
 }
 
 func HandleRetry(ctx context.Context, ch *amqp.Channel, d *amqp.Delivery) error {
+	logrus.Info("handleretry_max-retry-count", maxRetryCount)
 	if d.Headers == nil {
 		d.Headers = amqp.Table{}
 	}
diff --git a/internal/common/config/global.yaml b/internal/common/config/global.yaml
index 6b3d5eb..099b0f5 100644
--- a/internal/common/config/global.yaml
+++ b/internal/common/config/global.yaml
@@ -30,6 +30,16 @@ rabbitmq:
   password: guest
   host: 127.0.0.1
   port: 5672
+  max-retry: 3
+
+mongo:
+  user: root
+  password: password
+  host: 127.0.0.1
+  port: 27017
+  db-name: "order"
+  coll-name: "order"
+
 
 stripe-key: "${STRIPE_KEY}"
 endpoint-stripe-secret: "${ENDPOINT_STRIPE_SECRET}"
\ No newline at end of file
diff --git a/internal/common/config/viper.go b/internal/common/config/viper.go
index a0c91d2..4ba14c7 100644
--- a/internal/common/config/viper.go
+++ b/internal/common/config/viper.go
@@ -1,17 +1,52 @@
 package config
 
 import (
+	"fmt"
+	"os"
+	"path/filepath"
+	"runtime"
 	"strings"
+	"sync"
 
 	"github.com/spf13/viper"
 )
 
-func NewViperConfig() error {
+func init() {
+	if err := NewViperConfig(); err != nil {
+		panic(err)
+	}
+}
+
+var once sync.Once
+
+func NewViperConfig() (err error) {
+	once.Do(func() {
+		err = newViperConfig()
+	})
+	return
+}
+
+func newViperConfig() error {
+	relPath, err := getRelativePathFromCaller()
+	if err != nil {
+		return err
+	}
 	viper.SetConfigName("global")
 	viper.SetConfigType("yaml")
-	viper.AddConfigPath("../common/config")
+	viper.AddConfigPath(relPath)
 	viper.EnvKeyReplacer(strings.NewReplacer("_", "-"))
 	viper.AutomaticEnv()
 	_ = viper.BindEnv("stripe-key", "STRIPE_KEY", "endpoint-stripe-secret", "ENDPOINT_STRIPE_SECRET")
 	return viper.ReadInConfig()
 }
+
+func getRelativePathFromCaller() (relPath string, err error) {
+	callerPwd, err := os.Getwd()
+	if err != nil {
+		return
+	}
+	_, here, _, _ := runtime.Caller(0)
+	relPath, err = filepath.Rel(callerPwd, filepath.Dir(here))
+	fmt.Printf("caller from: %s, here: %s, relpath: %s", callerPwd, here, relPath)
+	return
+}
diff --git a/internal/order/adapters/order_mongo_repository.go b/internal/order/adapters/order_mongo_repository.go
new file mode 100644
index 0000000..fc4527d
--- /dev/null
+++ b/internal/order/adapters/order_mongo_repository.go
@@ -0,0 +1,158 @@
+package adapters
+
+import (
+	"context"
+	"time"
+
+	_ "github.com/ghost-yu/go_shop_second/common/config"
+	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
+	"github.com/ghost-yu/go_shop_second/order/entity"
+	"github.com/sirupsen/logrus"
+	"github.com/spf13/viper"
+	"go.mongodb.org/mongo-driver/bson"
+	"go.mongodb.org/mongo-driver/bson/primitive"
+	"go.mongodb.org/mongo-driver/mongo"
+)
+
+var (
+	dbName   = viper.GetString("mongo.db-name")
+	collName = viper.GetString("mongo.coll-name")
+)
+
+type OrderRepositoryMongo struct {
+	db *mongo.Client
+}
+
+func NewOrderRepositoryMongo(db *mongo.Client) *OrderRepositoryMongo {
+	return &OrderRepositoryMongo{db: db}
+}
+
+func (r *OrderRepositoryMongo) collection() *mongo.Collection {
+	return r.db.Database(dbName).Collection(collName)
+}
+
+type orderModel struct {
+	MongoID     primitive.ObjectID `bson:"_id"`
+	ID          string             `bson:"id"`
+	CustomerID  string             `bson:"customer_id"`
+	Status      string             `bson:"status"`
+	PaymentLink string             `bson:"payment_link"`
+	Items       []*entity.Item     `bson:"items"`
+}
+
+func (r *OrderRepositoryMongo) Create(ctx context.Context, order *domain.Order) (created *domain.Order, err error) {
+	defer r.logWithTag("create", err, order, created)
+	write := r.marshalToModel(order)
+	res, err := r.collection().InsertOne(ctx, write)
+	if err != nil {
+		return nil, err
+	}
+	created = order
+	created.ID = res.InsertedID.(primitive.ObjectID).Hex()
+	return created, nil
+}
+
+func (r *OrderRepositoryMongo) logWithTag(tag string, err error, input *domain.Order, result interface{}) {
+	l := logrus.WithFields(logrus.Fields{
+		"tag":            "order_repository_mongo",
+		"input_order":    input,
+		"performed_time": time.Now().Unix(),
+		"err":            err,
+		"result":         result,
+	})
+	if err != nil {
+		l.Infof("%s_fail", tag)
+	} else {
+		l.Infof("%s_success", tag)
+	}
+}
+
+func (r *OrderRepositoryMongo) Get(ctx context.Context, id, customerID string) (got *domain.Order, err error) {
+	defer r.logWithTag("get", err, nil, got)
+	read := &orderModel{}
+	mongoID, _ := primitive.ObjectIDFromHex(id)
+	cond := bson.M{"_id": mongoID}
+	if err = r.collection().FindOne(ctx, cond).Decode(read); err != nil {
+		return
+	}
+	if read == nil {
+		return nil, domain.NotFoundError{OrderID: id}
+	}
+	got = r.unmarshal(read)
+	return got, nil
+}
+
+// Update 先查找对应的order，然后apply updateFn，再写入回去
+func (r *OrderRepositoryMongo) Update(
+	ctx context.Context,
+	order *domain.Order,
+	updateFn func(context.Context, *domain.Order,
+	) (*domain.Order, error)) (err error) {
+	defer r.logWithTag("update", err, order, nil)
+	if order == nil {
+		panic("got nil order")
+	}
+	// 事务
+	session, err := r.db.StartSession()
+	if err != nil {
+		return
+	}
+	defer session.EndSession(ctx)
+
+	if err = session.StartTransaction(); err != nil {
+		return err
+	}
+	defer func() {
+		if err == nil {
+			_ = session.CommitTransaction(ctx)
+		} else {
+			_ = session.AbortTransaction(ctx)
+		}
+	}()
+
+	// inside transaction:
+	oldOrder, err := r.Get(ctx, order.ID, order.CustomerID)
+	if err != nil {
+		return
+	}
+	updated, err := updateFn(ctx, oldOrder)
+	if err != nil {
+		return
+	}
+	logrus.Infof("update||oldOrder=%+v||updated=%+v", oldOrder, updated)
+	mongoID, _ := primitive.ObjectIDFromHex(oldOrder.ID)
+	res, err := r.collection().UpdateOne(
+		ctx,
+		bson.M{"_id": mongoID, "customer_id": oldOrder.CustomerID},
+		bson.M{"$set": bson.M{
+			"status":       updated.Status,
+			"payment_link": updated.PaymentLink,
+		}},
+	)
+	if err != nil {
+		return
+	}
+	r.logWithTag("finish_update", err, order, res)
+	return
+}
+
+func (r *OrderRepositoryMongo) marshalToModel(order *domain.Order) *orderModel {
+	return &orderModel{
+		MongoID:     primitive.NewObjectID(),
+		ID:          order.ID,
+		CustomerID:  order.CustomerID,
+		Status:      order.Status,
+		PaymentLink: order.PaymentLink,
+		Items:       order.Items,
+	}
+}
+
+func (r *OrderRepositoryMongo) unmarshal(m *orderModel) *domain.Order {
+	return &domain.Order{
+		ID:          m.MongoID.Hex(),
+		CustomerID:  m.CustomerID,
+		Status:      m.Status,
+		PaymentLink: m.PaymentLink,
+		Items:       m.Items,
+	}
+}
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
index 64f9070..6384a96 100644
--- a/internal/order/app/command/create_order.go
+++ b/internal/order/app/command/create_order.go
@@ -75,10 +75,11 @@ func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*Creat
 	if err != nil {
 		return nil, err
 	}
-	o, err := c.orderRepo.Create(ctx, &domain.Order{
-		CustomerID: cmd.CustomerID,
-		Items:      validItems,
-	})
+	pendingOrder, err := domain.NewPendingOrder(cmd.CustomerID, validItems)
+	if err != nil {
+		return nil, err
+	}
+	o, err := c.orderRepo.Create(ctx, pendingOrder)
 	if err != nil {
 		return nil, err
 	}
diff --git a/internal/order/domain/order/order.go b/internal/order/domain/order/order.go
index 08e7e96..1afa1cb 100644
--- a/internal/order/domain/order/order.go
+++ b/internal/order/domain/order/order.go
@@ -38,6 +38,20 @@ func NewOrder(id, customerID, status, paymentLink string, items []*entity.Item)
 	}, nil
 }
 
+func NewPendingOrder(customerId string, items []*entity.Item) (*Order, error) {
+	if customerId == "" {
+		return nil, errors.New("empty customerID")
+	}
+	if items == nil {
+		return nil, errors.New("empty items")
+	}
+	return &Order{
+		CustomerID: customerId,
+		Status:     "pending",
+		Items:      items,
+	}, nil
+}
+
 func (o *Order) IsPaid() error {
 	if o.Status == string(stripe.CheckoutSessionPaymentStatusPaid) {
 		return nil
diff --git a/internal/order/go.mod b/internal/order/go.mod
index 9c2d3ea..9cd0692 100644
--- a/internal/order/go.mod
+++ b/internal/order/go.mod
@@ -14,6 +14,7 @@ require (
 	github.com/sirupsen/logrus v1.9.3
 	github.com/spf13/viper v1.19.0
 	github.com/stripe/stripe-go/v80 v80.2.0
+	go.mongodb.org/mongo-driver v1.17.1
 	go.opentelemetry.io/otel v1.31.0
 	google.golang.org/grpc v1.67.1
 	google.golang.org/protobuf v1.35.1
@@ -36,6 +37,7 @@ require (
 	github.com/go-playground/universal-translator v0.18.1 // indirect
 	github.com/go-playground/validator/v10 v10.22.1 // indirect
 	github.com/goccy/go-json v0.10.3 // indirect
+	github.com/golang/snappy v0.0.4 // indirect
 	github.com/google/uuid v1.6.0 // indirect
 	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
 	github.com/hashicorp/consul/api v1.28.2 // indirect
@@ -49,6 +51,7 @@ require (
 	github.com/hashicorp/hcl v1.0.0 // indirect
 	github.com/hashicorp/serf v0.10.1 // indirect
 	github.com/json-iterator/go v1.1.12 // indirect
+	github.com/klauspost/compress v1.17.2 // indirect
 	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
 	github.com/leodido/go-urn v1.4.0 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
@@ -58,6 +61,7 @@ require (
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
 	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
 	github.com/modern-go/reflect2 v1.0.2 // indirect
+	github.com/montanaflynn/stats v0.7.1 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
 	github.com/sagikazarmark/locafero v0.6.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
@@ -68,7 +72,10 @@ require (
 	github.com/subosito/gotenv v1.6.0 // indirect
 	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
 	github.com/ugorji/go/codec v1.2.12 // indirect
-	github.com/zenazn/goji v1.0.1 // indirect
+	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
+	github.com/xdg-go/scram v1.1.2 // indirect
+	github.com/xdg-go/stringprep v1.0.4 // indirect
+	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
 	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 // indirect
 	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
 	go.opentelemetry.io/contrib/propagators/b3 v1.31.0 // indirect
@@ -81,6 +88,7 @@ require (
 	golang.org/x/crypto v0.28.0 // indirect
 	golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c // indirect
 	golang.org/x/net v0.30.0 // indirect
+	golang.org/x/sync v0.8.0 // indirect
 	golang.org/x/sys v0.26.0 // indirect
 	golang.org/x/text v0.19.0 // indirect
 	google.golang.org/genproto/googleapis/rpc v0.0.0-20241007155032-5fefd90f89a9 // indirect
diff --git a/internal/order/go.sum b/internal/order/go.sum
index 176a808..aacb2f0 100644
--- a/internal/order/go.sum
+++ b/internal/order/go.sum
@@ -91,6 +91,8 @@ github.com/golang/protobuf v1.3.2/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5y
 github.com/golang/protobuf v1.3.3/go.mod h1:vzj43D7+SQXF/4pzW/hwtAqwc6iTitCiVSaWz5lYuqw=
 github.com/golang/protobuf v1.5.4 h1:i7eJL8qZTpSEXOPTxNKhASYpMn+8e5Q6AdndVa1dWek=
 github.com/golang/protobuf v1.5.4/go.mod h1:lnTiLA8Wa4RWRcIUkrtSVa5nRhsEGBg48fD6rSs7xps=
+github.com/golang/snappy v0.0.4 h1:yAGX7huGHXlcLOEtBnF4w7FQwA26wojNCwOYAEhLjQM=
+github.com/golang/snappy v0.0.4/go.mod h1:/XxbfmMg8lxefKM7IXC3fBNl/7bRcc72aCRzEWrmP2Q=
 github.com/google/btree v0.0.0-20180813153112-4030bb1f1f0c/go.mod h1:lNA+9X1NB3Zf8V7Ke586lFgjr2dZNuvo3lPJSGZ5JPQ=
 github.com/google/btree v1.0.1 h1:gK4Kx5IaGY9CD5sPJ36FHiBJ6ZXl0kilRiiCj+jdYp4=
 github.com/google/btree v1.0.1/go.mod h1:xXMiIv4Fb/0kKde4SpL7qlzvu5cMJDRkFDxJfI9uaxA=
@@ -158,6 +160,8 @@ github.com/juju/gnuflag v0.0.0-20171113085948-2ce1bb71843d/go.mod h1:2PavIy+JPci
 github.com/julienschmidt/httprouter v1.2.0/go.mod h1:SYymIcj16QtmaHHD7aYtjjsJG7VTCxuUUipMqKk8s4w=
 github.com/kisielk/errcheck v1.5.0/go.mod h1:pFxgyoBC7bSaBwPgfKdkLd5X25qrDl4LWUI2bnpBCr8=
 github.com/kisielk/gotool v1.0.0/go.mod h1:XhKaO+MFFWcvkIS/tQcRk01m1F5IRFswLeQ+oQHNcck=
+github.com/klauspost/compress v1.17.2 h1:RlWWUY/Dr4fL8qk9YG7DTZ7PDgME2V4csBXA8L/ixi4=
+github.com/klauspost/compress v1.17.2/go.mod h1:ntbaceVETuRiXiv4DpjP66DpAtAGkEQskQzEyD//IeE=
 github.com/klauspost/cpuid/v2 v2.0.9/go.mod h1:FInQzS24/EEf25PyTYn52gqo7WaD8xa0213Md/qVLRg=
 github.com/klauspost/cpuid/v2 v2.2.8 h1:+StwCXwm9PdpiEkPyzBXIy+M9KUb4ODm0Zarf1kS5BM=
 github.com/klauspost/cpuid/v2 v2.2.8/go.mod h1:Lcz8mBdAVJIBVzewtcLocK12l3Y+JytZYpaMropDUws=
@@ -207,6 +211,8 @@ github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742/go.mod h1:bx2lN
 github.com/modern-go/reflect2 v1.0.1/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3RllmbCylyMrvgv0=
 github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
 github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
+github.com/montanaflynn/stats v0.7.1 h1:etflOAAHORrCC44V+aR6Ftzort912ZU+YLiSTuV8eaE=
+github.com/montanaflynn/stats v0.7.1/go.mod h1:etXPPgVO6n31NxCd9KQUMvCM+ve0ruNzt6R8Bnaayow=
 github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
 github.com/oapi-codegen/runtime v1.1.1 h1:EXLHh0DXIJnWhdRPN2w4MXAzFyE4CskzhNLUmtpMYro=
 github.com/oapi-codegen/runtime v1.1.1/go.mod h1:SK9X900oXmPWilYR5/WKPzt3Kqxn/uS/+lbpREv+eCg=
@@ -287,10 +293,19 @@ github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS
 github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
 github.com/ugorji/go/codec v1.2.12 h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65EE=
 github.com/ugorji/go/codec v1.2.12/go.mod h1:UNopzCgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=
+github.com/xdg-go/pbkdf2 v1.0.0 h1:Su7DPu48wXMwC3bs7MCNG+z4FhcyEuz5dlvchbq0B0c=
+github.com/xdg-go/pbkdf2 v1.0.0/go.mod h1:jrpuAogTd400dnrH08LKmI/xc1MbPOebTwRqcT5RDeI=
+github.com/xdg-go/scram v1.1.2 h1:FHX5I5B4i4hKRVRBCFRxq1iQRej7WO3hhBuJf+UUySY=
+github.com/xdg-go/scram v1.1.2/go.mod h1:RT/sEzTbU5y00aCK8UOx6R7YryM0iF1N2MOmC3kKLN4=
+github.com/xdg-go/stringprep v1.0.4 h1:XLI/Ng3O1Atzq0oBs3TWm+5ZVgkq2aqdlvP9JtoZ6c8=
+github.com/xdg-go/stringprep v1.0.4/go.mod h1:mPGuuIYwz7CmR2bT9j4GbQqutWS1zV24gijq1dTyGkM=
+github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 h1:ilQV1hzziu+LLM3zUTJ0trRztfwgjqKnBWNtSRkbmwM=
+github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78/go.mod h1:aL8wCCfTfSfmXjznFBSZNN13rSJjlIOI1fUNAtF7rmI=
 github.com/yuin/goldmark v1.1.27/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
-github.com/zenazn/goji v1.0.1 h1:4lbD8Mx2h7IvloP7r2C0D6ltZP6Ufip8Hn0wmSK5LR8=
-github.com/zenazn/goji v1.0.1/go.mod h1:7S9M489iMyHBNxwZnk9/EHS098H4/F6TATF2mIxtB1Q=
+github.com/yuin/goldmark v1.4.13/go.mod h1:6yULJ656Px+3vBD8DxQVa3kxgyrAnzto9xy5taEt/CY=
+go.mongodb.org/mongo-driver v1.17.1 h1:Wic5cJIwJgSpBhe3lx3+/RybR5PiYRMpVFgO7cOHyIM=
+go.mongodb.org/mongo-driver v1.17.1/go.mod h1:wwWm/+BuOddhcq3n68LKRmgk2wXzmF6s0SFOa0GINL4=
 go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 h1:0nTRpaCaILLdooXAQnfktlL6Zw1ECKEW9DZGH2byi2c=
 go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0/go.mod h1:A7aFlp4WSLmeOnFRZwf2dMU+40THPc+rsr6KOwZLOcg=
 go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 h1:4Pp6oUg3+e/6M4C0A/3kJ2VYa++dsWVTtGgLVj5xtHg=
@@ -322,6 +337,7 @@ golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACk
 golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392/go.mod h1:/lpIB1dKB+9EgE3H3cr1v9wB50oz8l4C4h62xy7jSTY=
 golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
 golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9/go.mod h1:LzIPMQfyMNhhGPhUkYOs5KpL4U8rLKemX1yGLhDgUto=
+golang.org/x/crypto v0.0.0-20210921155107-089bfa567519/go.mod h1:GvvjBRRGRdwPK5ydBHafDWAxML/pGHZbMvKqRZ5+Abc=
 golang.org/x/crypto v0.28.0 h1:GBDwsMXVQi34v5CCYUm2jkJvu4cbtru2U4TN2PSyQnw=
 golang.org/x/crypto v0.28.0/go.mod h1:rmgy+3RHxRZMyY0jjAJShp2zgEdOqj2AO7U0pYmeQ7U=
 golang.org/x/exp v0.0.0-20190121172915-509febef88a4/go.mod h1:CJ0aWSM057203Lf6IL+f9T1iT9GByDxfZKAQTCR3kQA=
@@ -333,6 +349,7 @@ golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3/go.mod h1:6SW0HCj/g11FgYtHl
 golang.org/x/lint v0.0.0-20190930215403-16217165b5de/go.mod h1:6SW0HCj/g11FgYtHlgUYUwCkIfeOF89ocIRzGO/8vkc=
 golang.org/x/mod v0.2.0/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
 golang.org/x/mod v0.3.0/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
+golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4/go.mod h1:jJ57K6gSWd91VN4djpZkiMVwK6gcyfeH4XE8wZrZaV4=
 golang.org/x/net v0.0.0-20180724234803-3673e40ba225/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20180826012351-8a410e7b638d/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20181114220301-adae6a3d119a/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
@@ -347,6 +364,7 @@ golang.org/x/net v0.0.0-20201021035429-f5854403a974/go.mod h1:sp8m0HH+o8qH0wwXwY
 golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
 golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1/go.mod h1:9tjilg8BloeKEkVJvy7fQ90B1CfIiPueXVOjqfkSzI8=
 golang.org/x/net v0.0.0-20210520170846-37e1c6afe023/go.mod h1:9nx3DQGgdP8bBQD5qxJ1jj9UTztislL4KSBs9R2vV5Y=
+golang.org/x/net v0.0.0-20220722155237-a158d28d115b/go.mod h1:XRhObCWvk6IyKnWLug+ECip1KBveYUHfp+8e9klMJ9c=
 golang.org/x/net v0.30.0 h1:AcW1SDZMkb8IpzCdQUaIq2sP4sZ4zw+55h6ynffypl4=
 golang.org/x/net v0.30.0/go.mod h1:2wGyMJ5iFasEhkwi13ChkO/t1ECNC4X4eBKkVFyYFlU=
 golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be/go.mod h1:N/0e6XlmueqKjAGxoOufVs8QHGRruUQn6yWY3a++T0U=
@@ -357,6 +375,9 @@ golang.org/x/sync v0.0.0-20190423024810-112230192c58/go.mod h1:RxMgew5VJxzue5/jJ
 golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20210220032951-036812b2e83c/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
+golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
+golang.org/x/sync v0.8.0 h1:3NFvSEYkUoMifnESzZl15y791HH1qU2xm6eCJU5ZPXQ=
+golang.org/x/sync v0.8.0/go.mod h1:Czt+wKu1gCyEFDUtn0jG5QVvpJ6rzVqr5aXyt9drQfk=
 golang.org/x/sys v0.0.0-20180823144017-11551d06cbcc/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20180830151530-49385e6e1522/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20180905080454-ebe1bf3edb33/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
@@ -376,11 +397,14 @@ golang.org/x/sys v0.0.0-20201119102817-f84b799fce68/go.mod h1:h1NjWce9XRLGQEsW7w
 golang.org/x/sys v0.0.0-20210303074136-134d130e1a04/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210330210617-4fbd30eecc44/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210423082822-04245dca01da/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
@@ -388,10 +412,13 @@ golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.26.0 h1:KHjCJyddX0LoSTb3J+vWpupP9p0oznkqVk/IfjymZbo=
 golang.org/x/sys v0.26.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
 golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
+golang.org/x/term v0.0.0-20210927222741-03fcf44c2211/go.mod h1:jbD1KX2456YbFQfuXm/mYQcufACuNUgVhRMnK/tPxf8=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
 golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
 golang.org/x/text v0.3.6/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
+golang.org/x/text v0.3.7/go.mod h1:u+2+/6zg+i71rQMx5EYifcz6MCKuco9NR6JIITiCfzQ=
+golang.org/x/text v0.3.8/go.mod h1:E6s5w1FMmriuDzIBO73fBruAKo1PCIq6d2Q6DHfQ8WQ=
 golang.org/x/text v0.19.0 h1:kTxAhCbGbxhK0IwgSKiMO5awPoDQ0RpfiVYBfK860YM=
 golang.org/x/text v0.19.0/go.mod h1:BuEKDfySbSR4drPmRPG/7iBdf8hvFMuRexcpahXilzY=
 golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
@@ -404,6 +431,7 @@ golang.org/x/tools v0.0.0-20191108193012-7d206e10da11/go.mod h1:b+2E5dAYhXwXZwtn
 golang.org/x/tools v0.0.0-20191119224855-298f0cb1881e/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
 golang.org/x/tools v0.0.0-20200619180055-7c47624df98f/go.mod h1:EkVYQZoAsY45+roYkvgYkIh4xh/qjgUK9TdY2XT94GE=
 golang.org/x/tools v0.0.0-20210106214847-113979e3529a/go.mod h1:emZCQorbCU4vsT4fOWvOPXz4eW1wZW4PmDk9uLelYpA=
+golang.org/x/tools v0.1.12/go.mod h1:hNGJHUnrk76NpqgfD5Aqm5Crs+Hm0VOH/i9J2+nxYbc=
 golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
 golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
 golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
diff --git a/internal/order/infrastructure/consumer/consumer.go b/internal/order/infrastructure/consumer/consumer.go
index 623bb50..511e0e8 100644
--- a/internal/order/infrastructure/consumer/consumer.go
+++ b/internal/order/infrastructure/consumer/consumer.go
@@ -40,25 +40,33 @@ func (c *Consumer) Listen(ch *amqp.Channel) {
 	var forever chan struct{}
 	go func() {
 		for msg := range msgs {
-			c.handleMessage(msg, q)
+			c.handleMessage(ch, msg, q)
 		}
 	}()
 	<-forever
 }
 
-func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
+func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
 	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
 	t := otel.Tracer("rabbitmq")
 	_, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
 	defer span.End()
 
+	var err error
+	defer func() {
+		if err != nil {
+			_ = msg.Nack(false, false)
+		} else {
+			_ = msg.Ack(false)
+		}
+	}()
+
 	o := &domain.Order{}
 	if err := json.Unmarshal(msg.Body, o); err != nil {
 		logrus.Infof("error unmarshal msg.body into domain.order, err = %v", err)
-		_ = msg.Nack(false, false)
 		return
 	}
-	_, err := c.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
+	_, err = c.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
 		Order: o,
 		UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
 			if err := order.IsPaid(); err != nil {
@@ -69,11 +77,12 @@ func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
 	})
 	if err != nil {
 		logrus.Infof("error updating order, orderID = %s, err = %v", o.ID, err)
-		// TODO: retry
+		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
+			logrus.Warnf("retry_error, error handling retry, messageID=%s, err=%v", msg.MessageId, err)
+		}
 		return
 	}
 
 	span.AddEvent("order.updated")
-	_ = msg.Ack(false)
 	logrus.Info("order consume paid event success!")
 }
diff --git a/internal/order/main.go b/internal/order/main.go
index 04fb795..67dc3b0 100644
--- a/internal/order/main.go
+++ b/internal/order/main.go
@@ -4,7 +4,7 @@ import (
 	"context"
 
 	"github.com/ghost-yu/go_shop_second/common/broker"
-	"github.com/ghost-yu/go_shop_second/common/config"
+	_ "github.com/ghost-yu/go_shop_second/common/config"
 	"github.com/ghost-yu/go_shop_second/common/discovery"
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	"github.com/ghost-yu/go_shop_second/common/logging"
@@ -21,9 +21,6 @@ import (
 
 func init() {
 	logging.Init()
-	if err := config.NewViperConfig(); err != nil {
-		logrus.Fatal(err)
-	}
 }
 
 func main() {
diff --git a/internal/order/service/application.go b/internal/order/service/application.go
index 8de0e53..59bbe5a 100644
--- a/internal/order/service/application.go
+++ b/internal/order/service/application.go
@@ -2,6 +2,8 @@ package service
 
 import (
 	"context"
+	"fmt"
+	"time"
 
 	"github.com/ghost-yu/go_shop_second/common/broker"
 	grpcClient "github.com/ghost-yu/go_shop_second/common/client"
@@ -14,6 +16,9 @@ import (
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
 	"github.com/spf13/viper"
+	"go.mongodb.org/mongo-driver/mongo"
+	"go.mongodb.org/mongo-driver/mongo/options"
+	"go.mongodb.org/mongo-driver/mongo/readpref"
 )
 
 func NewApplication(ctx context.Context) (app.Application, func()) {
@@ -36,7 +41,9 @@ func NewApplication(ctx context.Context) (app.Application, func()) {
 }
 
 func newApplication(_ context.Context, stockGRPC query.StockService, ch *amqp.Channel) app.Application {
-	orderRepo := adapters.NewMemoryOrderRepository()
+	//orderRepo := adapters.NewMemoryOrderRepository()
+	mongoClient := newMongoClient()
+	orderRepo := adapters.NewOrderRepositoryMongo(mongoClient)
 	logger := logrus.NewEntry(logrus.StandardLogger())
 	metricClient := metrics.TodoMetrics{}
 	return app.Application{
@@ -49,3 +56,23 @@ func newApplication(_ context.Context, stockGRPC query.StockService, ch *amqp.Ch
 		},
 	}
 }
+
+func newMongoClient() *mongo.Client {
+	uri := fmt.Sprintf(
+		"mongodb://%s:%s@%s:%s",
+		viper.GetString("mongo.user"),
+		viper.GetString("mongo.password"),
+		viper.GetString("mongo.host"),
+		viper.GetString("mongo.port"),
+	)
+	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
+	defer cancel()
+	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
+	if err != nil {
+		panic(err)
+	}
+	if err = c.Ping(ctx, readpref.Primary()); err != nil {
+		panic(err)
+	}
+	return c
+}
diff --git a/internal/payment/infrastructure/consumer/consumer.go b/internal/payment/infrastructure/consumer/consumer.go
index 99046de..611c5ae 100644
--- a/internal/payment/infrastructure/consumer/consumer.go
+++ b/internal/payment/infrastructure/consumer/consumer.go
@@ -38,33 +38,41 @@ func (c *Consumer) Listen(ch *amqp.Channel) {
 	var forever chan struct{}
 	go func() {
 		for msg := range msgs {
-			c.handleMessage(msg, q)
+			c.handleMessage(ch, msg, q)
 		}
 	}()
 	<-forever
 }
 
-func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
+func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
 	logrus.Infof("Payment receive a message from %s, msg=%v", q.Name, string(msg.Body))
 	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
 	tr := otel.Tracer("rabbitmq")
 	_, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
 	defer span.End()
 
+	var err error
+	defer func() {
+		if err != nil {
+			_ = msg.Nack(false, false)
+		} else {
+			_ = msg.Ack(false)
+		}
+	}()
+
 	o := &orderpb.Order{}
 	if err := json.Unmarshal(msg.Body, o); err != nil {
 		logrus.Infof("failed to unmarshall msg to order, err=%v", err)
-		_ = msg.Nack(false, false)
 		return
 	}
 	if _, err := c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{Order: o}); err != nil {
-		// TODO: retry
-		logrus.Infof("failed to create order, err=%v", err)
-		_ = msg.Nack(false, false)
+		logrus.Infof("failed to create payment, err=%v", err)
+		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
+			logrus.Warnf("retry_error, error handling retry, messageID=%s, err=%v", msg.MessageId, err)
+		}
 		return
 	}
 
 	span.AddEvent("payment.created")
-	_ = msg.Ack(false)
 	logrus.Info("consume success")
 }
diff --git a/internal/payment/main.go b/internal/payment/main.go
index d0e4165..f810389 100644
--- a/internal/payment/main.go
+++ b/internal/payment/main.go
@@ -4,7 +4,7 @@ import (
 	"context"
 
 	"github.com/ghost-yu/go_shop_second/common/broker"
-	"github.com/ghost-yu/go_shop_second/common/config"
+	_ "github.com/ghost-yu/go_shop_second/common/config"
 	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/ghost-yu/go_shop_second/common/server"
 	"github.com/ghost-yu/go_shop_second/common/tracing"
@@ -16,9 +16,6 @@ import (
 
 func init() {
 	logging.Init()
-	if err := config.NewViperConfig(); err != nil {
-		logrus.Fatal(err)
-	}
 }
 
 func main() {
diff --git a/internal/stock/main.go b/internal/stock/main.go
index 548e6f1..18fa2e3 100644
--- a/internal/stock/main.go
+++ b/internal/stock/main.go
@@ -3,7 +3,7 @@ package main
 import (
 	"context"
 
-	"github.com/ghost-yu/go_shop_second/common/config"
+	_ "github.com/ghost-yu/go_shop_second/common/config"
 	"github.com/ghost-yu/go_shop_second/common/discovery"
 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
 	"github.com/ghost-yu/go_shop_second/common/logging"
@@ -18,9 +18,6 @@ import (
 
 func init() {
 	logging.Init()
-	if err := config.NewViperConfig(); err != nil {
-		logrus.Fatal(err)
-	}
 }
 
 func main() {
~~~

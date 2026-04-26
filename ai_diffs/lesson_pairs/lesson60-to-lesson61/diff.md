# Lesson Pair Diff Report

- FromBranch: lesson60
- ToBranch: lesson61

## Short Summary

~~~text
 10 files changed, 168 insertions(+), 85 deletions(-)
~~~

## File Stats

~~~text
 internal/common/broker/event.go            | 89 ++++++++++++++++++++++++++++++
 internal/common/broker/rabbitmq.go         | 19 +++++--
 internal/common/config/global.yaml         |  3 +
 internal/common/go.mod                     |  1 +
 internal/common/logging/when.go            |  5 +-
 internal/kitchen/go.mod                    | 24 ++++----
 internal/kitchen/go.sum                    | 52 ++++++++---------
 internal/order/app/command/create_order.go | 26 +++------
 internal/payment/go.mod                    |  1 +
 internal/payment/http.go                   | 33 ++++-------
 10 files changed, 168 insertions(+), 85 deletions(-)
~~~

## Commit Comparison

~~~text
> 10be7fe 全链路可观测性建设-下(mq异步链路)
~~~

## Changed Files

~~~text
internal/common/broker/event.go
internal/common/broker/rabbitmq.go
internal/common/config/global.yaml
internal/common/go.mod
internal/common/logging/when.go
internal/kitchen/go.mod
internal/kitchen/go.sum
internal/order/app/command/create_order.go
internal/payment/go.mod
internal/payment/http.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
internal/common/broker/event.go
internal/common/broker/rabbitmq.go
internal/common/config/global.yaml
internal/common/logging/when.go
internal/order/app/command/create_order.go
internal/payment/http.go
~~~

## Full Diff

~~~diff
diff --git a/internal/common/broker/event.go b/internal/common/broker/event.go
index e53531c..2e93939 100644
--- a/internal/common/broker/event.go
+++ b/internal/common/broker/event.go
@@ -1,6 +1,95 @@
 package broker
 
+import (
+	"context"
+	"encoding/json"
+
+	"github.com/ghost-yu/go_shop_second/common/logging"
+	"github.com/pkg/errors"
+	amqp "github.com/rabbitmq/amqp091-go"
+	"github.com/sirupsen/logrus"
+)
+
 const (
 	EventOrderCreated = "order.created"
 	EventOrderPaid    = "order.paid"
 )
+
+type RoutingType string
+
+const (
+	FanOut RoutingType = "fan-out"
+	Direct RoutingType = "direct"
+)
+
+type PublishEventReq struct {
+	Channel  *amqp.Channel
+	Routing  RoutingType
+	Queue    string
+	Exchange string
+	Body     any
+}
+
+func PublishEvent(ctx context.Context, p PublishEventReq) (err error) {
+	_, dLog := logging.WhenEventPublish(ctx, p)
+	defer dLog(nil, &err)
+
+	if err = checkParam(p); err != nil {
+		return err
+	}
+
+	switch p.Routing {
+	default:
+		logrus.WithContext(ctx).Panicf("unsupported routing type: %s", string(p.Routing))
+	case FanOut:
+		return fanOut(ctx, p)
+	case Direct:
+		return directQueue(ctx, p)
+	}
+	return nil
+}
+
+func checkParam(p PublishEventReq) error {
+	if p.Channel == nil {
+		return errors.New("nil channel")
+	}
+	return nil
+}
+
+func directQueue(ctx context.Context, p PublishEventReq) (err error) {
+	_, err = p.Channel.QueueDeclare(p.Queue, true, false, false, false, nil)
+	if err != nil {
+		return err
+	}
+	jsonBody, err := json.Marshal(p.Body)
+	if err != nil {
+		return err
+	}
+	return doPublish(ctx, p.Channel, p.Exchange, p.Queue, false, false, amqp.Publishing{
+		ContentType:  "application/json",
+		DeliveryMode: amqp.Persistent,
+		Body:         jsonBody,
+		Headers:      InjectRabbitMQHeaders(ctx),
+	})
+}
+
+func doPublish(ctx context.Context, ch *amqp.Channel, exchange, key string, mandatory bool, immediate bool, msg amqp.Publishing) error {
+	if err := ch.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg); err != nil {
+		logging.Warnf(ctx, nil, "_publish_event_failed||exchange=%s||key=%s||msg=%v", exchange, key, msg)
+		return errors.Wrap(err, "publish event error")
+	}
+	return nil
+}
+
+func fanOut(ctx context.Context, p PublishEventReq) (err error) {
+	jsonBody, err := json.Marshal(p.Body)
+	if err != nil {
+		return err
+	}
+	return doPublish(ctx, p.Channel, p.Exchange, "", false, false, amqp.Publishing{
+		ContentType:  "application/json",
+		DeliveryMode: amqp.Persistent,
+		Body:         jsonBody,
+		Headers:      InjectRabbitMQHeaders(ctx),
+	})
+}
diff --git a/internal/common/broker/rabbitmq.go b/internal/common/broker/rabbitmq.go
index ef74e7d..1a835d2 100644
--- a/internal/common/broker/rabbitmq.go
+++ b/internal/common/broker/rabbitmq.go
@@ -6,6 +6,7 @@ import (
 	"time"
 
 	_ "github.com/ghost-yu/go_shop_second/common/config"
+	"github.com/ghost-yu/go_shop_second/common/logging"
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
 	"github.com/spf13/viper"
@@ -63,8 +64,13 @@ func createDLX(ch *amqp.Channel) error {
 	return err
 }
 
-func HandleRetry(ctx context.Context, ch *amqp.Channel, d *amqp.Delivery) error {
-	logrus.Info("handleretry_max-retry-count", maxRetryCount)
+func HandleRetry(ctx context.Context, ch *amqp.Channel, d *amqp.Delivery) (err error) {
+	fields, dLog := logging.WhenRequest(ctx, "HandleRetry", map[string]any{
+		"delivery":        d,
+		"max_retry_count": maxRetryCount,
+	})
+	defer dLog(nil, &err)
+
 	if d.Headers == nil {
 		d.Headers = amqp.Table{}
 	}
@@ -74,19 +80,20 @@ func HandleRetry(ctx context.Context, ch *amqp.Channel, d *amqp.Delivery) error
 	}
 	retryCount++
 	d.Headers[amqpRetryHeaderKey] = retryCount
+	fields["retry_count"] = retryCount
 
 	if retryCount >= maxRetryCount {
-		logrus.Infof("moving message %s to dlq", d.MessageId)
-		return ch.PublishWithContext(ctx, "", DLQ, false, false, amqp.Publishing{
+		logrus.WithContext(ctx).Infof("moving message %s to dlq", d.MessageId)
+		return doPublish(ctx, ch, "", DLQ, false, false, amqp.Publishing{
 			Headers:      d.Headers,
 			ContentType:  "application/json",
 			Body:         d.Body,
 			DeliveryMode: amqp.Persistent,
 		})
 	}
-	logrus.Infof("retring message %s, count=%d", d.MessageId, retryCount)
+	logrus.WithContext(ctx).Debugf("retring message %s, count=%d", d.MessageId, retryCount)
 	time.Sleep(time.Second * time.Duration(retryCount))
-	return ch.PublishWithContext(ctx, d.Exchange, d.RoutingKey, false, false, amqp.Publishing{
+	return doPublish(ctx, ch, "", DLQ, false, false, amqp.Publishing{
 		Headers:      d.Headers,
 		ContentType:  "application/json",
 		Body:         d.Body,
diff --git a/internal/common/config/global.yaml b/internal/common/config/global.yaml
index 4b0b19c..b65dfe4 100644
--- a/internal/common/config/global.yaml
+++ b/internal/common/config/global.yaml
@@ -25,6 +25,9 @@ payment:
   http-addr: 127.0.0.1:8284
   grpc-addr: 127.0.0.1:5004
 
+kitchen:
+  service-name: kitchen
+
 rabbitmq:
   user: guest
   password: guest
diff --git a/internal/common/go.mod b/internal/common/go.mod
index 6ec52f4..c208169 100644
--- a/internal/common/go.mod
+++ b/internal/common/go.mod
@@ -7,6 +7,7 @@ require (
 	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
 	github.com/hashicorp/consul/api v1.28.2
 	github.com/oapi-codegen/runtime v1.1.1
+	github.com/pkg/errors v0.9.1
 	github.com/rabbitmq/amqp091-go v1.10.0
 	github.com/redis/go-redis/v9 v9.7.0
 	github.com/sirupsen/logrus v1.8.1
diff --git a/internal/common/logging/when.go b/internal/common/logging/when.go
index f123008..8698a55 100644
--- a/internal/common/logging/when.go
+++ b/internal/common/logging/when.go
@@ -38,10 +38,9 @@ func WhenRequest(ctx context.Context, method string, args ...any) (logrus.Fields
 	}
 }
 
-func WhenEventPublish(ctx context.Context, method string, args ...any) (logrus.Fields, func(any, *error)) {
+func WhenEventPublish(ctx context.Context, args ...any) (logrus.Fields, func(any, *error)) {
 	fields := logrus.Fields{
-		Method: method,
-		Args:   formatArgs(args),
+		Args: formatArgs(args),
 	}
 	start := time.Now()
 	return fields, func(resp any, err *error) {
diff --git a/internal/kitchen/go.mod b/internal/kitchen/go.mod
index 2a29c39..22b3e73 100644
--- a/internal/kitchen/go.mod
+++ b/internal/kitchen/go.mod
@@ -6,7 +6,7 @@ replace github.com/ghost-yu/go_shop_second/common => ../common
 
 require (
 	github.com/ghost-yu/go_shop_second/common v0.0.0-00010101000000-000000000000
-	github.com/prometheus/client_golang v1.20.5
+	github.com/pkg/errors v0.9.1
 	github.com/rabbitmq/amqp091-go v1.10.0
 	github.com/sirupsen/logrus v1.9.3
 	github.com/spf13/viper v1.19.0
@@ -15,8 +15,6 @@ require (
 
 require (
 	github.com/armon/go-metrics v0.4.1 // indirect
-	github.com/beorn7/perks v1.0.1 // indirect
-	github.com/cespare/xxhash/v2 v2.3.0 // indirect
 	github.com/fatih/color v1.14.1 // indirect
 	github.com/fsnotify/fsnotify v1.7.0 // indirect
 	github.com/go-logr/logr v1.4.2 // indirect
@@ -32,17 +30,15 @@ require (
 	github.com/hashicorp/golang-lru v0.5.4 // indirect
 	github.com/hashicorp/hcl v1.0.0 // indirect
 	github.com/hashicorp/serf v0.10.1 // indirect
-	github.com/klauspost/compress v1.17.9 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
 	github.com/mattn/go-colorable v0.1.13 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
+	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
 	github.com/mitchellh/go-homedir v1.1.0 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
-	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
+	github.com/nxadm/tail v1.4.11 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
-	github.com/prometheus/client_model v0.6.1 // indirect
-	github.com/prometheus/common v0.55.0 // indirect
-	github.com/prometheus/procfs v0.15.1 // indirect
+	github.com/rogpeppe/go-internal v1.10.0 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
 	github.com/sourcegraph/conc v0.3.0 // indirect
@@ -50,6 +46,7 @@ require (
 	github.com/spf13/cast v1.6.0 // indirect
 	github.com/spf13/pflag v1.0.5 // indirect
 	github.com/subosito/gotenv v1.6.0 // indirect
+	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
 	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
 	go.opentelemetry.io/contrib/propagators/b3 v1.31.0 // indirect
 	go.opentelemetry.io/otel/exporters/jaeger v1.17.0 // indirect
@@ -58,13 +55,16 @@ require (
 	go.opentelemetry.io/otel/trace v1.31.0 // indirect
 	go.uber.org/atomic v1.9.0 // indirect
 	go.uber.org/multierr v1.9.0 // indirect
+	golang.org/x/crypto v0.31.0 // indirect
 	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
-	golang.org/x/net v0.30.0 // indirect
-	golang.org/x/sys v0.26.0 // indirect
-	golang.org/x/text v0.19.0 // indirect
+	golang.org/x/net v0.33.0 // indirect
+	golang.org/x/sys v0.28.0 // indirect
+	golang.org/x/term v0.27.0 // indirect
+	golang.org/x/text v0.21.0 // indirect
 	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
 	google.golang.org/grpc v1.67.1 // indirect
-	google.golang.org/protobuf v1.35.1 // indirect
+	google.golang.org/protobuf v1.36.1 // indirect
+	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
 	gopkg.in/ini.v1 v1.67.0 // indirect
 	gopkg.in/yaml.v3 v3.0.1 // indirect
 )
diff --git a/internal/kitchen/go.sum b/internal/kitchen/go.sum
index 1ab1071..1528c76 100644
--- a/internal/kitchen/go.sum
+++ b/internal/kitchen/go.sum
@@ -11,12 +11,9 @@ github.com/armon/go-radix v0.0.0-20180808171621-7fddfc383310/go.mod h1:ufUuZ+zHj
 github.com/armon/go-radix v1.0.0/go.mod h1:ufUuZ+zHj4x4TnLV4JWEpy2hxWSpsRywHrMgIH9cCH8=
 github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973/go.mod h1:Dwedo/Wpr24TaqPxmxbtue+5NUziq4I4S80YR8gNf3Q=
 github.com/beorn7/perks v1.0.0/go.mod h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+CedLV8=
-github.com/beorn7/perks v1.0.1 h1:VlbKKnNfV8bJzeqoa4cOKqO6bYr3WgKZxO8Z16+hsOM=
 github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
 github.com/bgentry/speakeasy v0.1.0/go.mod h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=
 github.com/cespare/xxhash/v2 v2.1.1/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
-github.com/cespare/xxhash/v2 v2.3.0 h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=
-github.com/cespare/xxhash/v2 v2.3.0/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
 github.com/circonus-labs/circonus-gometrics v2.3.1+incompatible/go.mod h1:nmEj6Dob7S7YxXgwXpfOuvO54S+tGdZdw9fuRZt25Ag=
 github.com/circonus-labs/circonusllhist v0.1.3/go.mod h1:kMXHVDlOchFAehlya5ePtbp5jckzBHf4XRpQvBOLI+I=
 github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
@@ -30,6 +27,7 @@ github.com/fatih/color v1.14.1 h1:qfhVLaG5s+nCROl1zJsZRxFeYrHLqWroPOQ8BWiNb4w=
 github.com/fatih/color v1.14.1/go.mod h1:2oHN61fhTpgcxD3TSWCgKDiH1+x4OiDVVGH8WlgGZGg=
 github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHkI4W8=
 github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
+github.com/fsnotify/fsnotify v1.6.0/go.mod h1:sl3t1tCWJFWoRz9R8WJCbQihKKwmorjAbSClcnxKAGw=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
 github.com/go-kit/kit v0.8.0/go.mod h1:xBxKIO96dXMWWy0MnWVtmwkA9/13aqxPnvrjFYMA2as=
@@ -105,19 +103,16 @@ github.com/hashicorp/serf v0.10.1/go.mod h1:yL2t6BqATOLGc5HF7qbFkTfXoPIY0WZdWHfE
 github.com/json-iterator/go v1.1.6/go.mod h1:+SdeFBvtyEkXs7REEP0seUULqWtbJapLOCVDaaPEHmU=
 github.com/json-iterator/go v1.1.9/go.mod h1:KdQUCv79m/52Kvf8AW2vK1V8akMuk1QjK/uOdHXbAo4=
 github.com/julienschmidt/httprouter v1.2.0/go.mod h1:SYymIcj16QtmaHHD7aYtjjsJG7VTCxuUUipMqKk8s4w=
-github.com/klauspost/compress v1.17.9 h1:6KIumPrER1LHsvBVuDa0r5xaG0Es51mhhB9BQB2qeMA=
-github.com/klauspost/compress v1.17.9/go.mod h1:Di0epgTjJY877eYKx5yC51cX2A2Vl2ibi7bDH9ttBbw=
 github.com/konsorten/go-windows-terminal-sequences v1.0.1/go.mod h1:T0+1ngSBFLxvqU3pZ+m/2kptfBszLMUkC4ZK/EgS/cQ=
 github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515/go.mod h1:+0opPa2QZZtGFBFZlji/RkVcI2GknAs/DXo4wKdlNEc=
 github.com/kr/pretty v0.1.0/go.mod h1:dAy3ld7l9f0ibDNOQOHHMYYIIbhfbHSm3C4ZsoJORNo=
+github.com/kr/pretty v0.2.1/go.mod h1:ipq/a2n7PKx3OHsz4KJII5eveXtPO4qwEXGdVfWzfnI=
 github.com/kr/pretty v0.3.1 h1:flRD4NNwYAUpkphVc1HcthR4KEIFJ65n8Mw5qdRn3LE=
 github.com/kr/pretty v0.3.1/go.mod h1:hoEshYVHaxMs3cyo3Yncou5ZscifuDolrwPKZanG3xk=
 github.com/kr/pty v1.1.1/go.mod h1:pFQYn66WHrOpPYNljwOMqo10TkYh1fy3cYio2l3bCsQ=
 github.com/kr/text v0.1.0/go.mod h1:4Jbv+DJW3UT/LiOwJeYQe1efqtUx/iVham/4vfdArNI=
 github.com/kr/text v0.2.0 h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=
 github.com/kr/text v0.2.0/go.mod h1:eLer722TekiGuMkidMxC/pM04lWEeraHUUmBw8l2grE=
-github.com/kylelemons/godebug v1.1.0 h1:RPNrshWIDI6G2gRW9EHilWtl7Z6Sb1BR0xunSBf0SNc=
-github.com/kylelemons/godebug v1.1.0/go.mod h1:9/0rRGxNHcop5bhtWyNeEfOS8JIWk580+fNqagV/RAw=
 github.com/magiconair/properties v1.8.7 h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=
 github.com/magiconair/properties v1.8.7/go.mod h1:Dhd985XPs7jluiymwWYZ0G4Z61jb3vdS329zhj2hYo0=
 github.com/mattn/go-colorable v0.0.9/go.mod h1:9vuHe8Xs5qXnSaW/c/ABM9alt+Vo+STaOChaDxuIBZU=
@@ -136,6 +131,8 @@ github.com/mattn/go-isatty v0.0.16/go.mod h1:kYGgaQfpe5nmfYZH+SKPsOc2e4SrIfOl2e/
 github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
 github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
 github.com/matttproud/golang_protobuf_extensions v1.0.1/go.mod h1:D8He9yQNgCq6Z5Ld7szi9bcBfOoFv/3dc6xSMkL2PC0=
+github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d h1:5PJl274Y63IEHC+7izoQE9x6ikvDFZS2mDVS3drnohI=
+github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d/go.mod h1:01TrycV0kFyexm33Z7vhZRXopbI8J3TDReVlkTgMUxE=
 github.com/miekg/dns v1.1.26/go.mod h1:bPDLeHnStXmXAq1m/Ch/hvfNHr14JKNPMBo3VZKjuso=
 github.com/miekg/dns v1.1.41 h1:WMszZWJG0XmzbK9FEmzH2TVcqYzFesusSIB41b8KHxY=
 github.com/miekg/dns v1.1.41/go.mod h1:p6aan82bvRIyn+zDIv9xYNUpwa73JcSh9BKwknJysuI=
@@ -149,9 +146,13 @@ github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421/go.mod h1:6dJ
 github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
 github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3RllmbCylyMrvgv0=
 github.com/modern-go/reflect2 v1.0.1/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3RllmbCylyMrvgv0=
-github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 h1:C3w9PqII01/Oq1c1nUAm88MOHcQC9l5mIlSMApZMrHA=
-github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822/go.mod h1:+n7T8mK8HuQTcFwEeznm/DIxMOiR9yIdICNftLE1DvQ=
 github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
+github.com/nxadm/tail v1.4.11 h1:8feyoE3OzPrcshW5/MJ4sGESc5cqmGkGCWlco4l0bqY=
+github.com/nxadm/tail v1.4.11/go.mod h1:OTaG3NK980DZzxbRq6lEuzgU+mug70nY11sMd4JXXHc=
+github.com/onsi/ginkgo v1.16.5 h1:8xi0RTUf59SOSfEtZMvwTvXYMzG4gV23XVHOZiXNtnE=
+github.com/onsi/ginkgo v1.16.5/go.mod h1:+E8gABHa3K6zRBolWtd+ROzc/U5bkGt0FwiG042wbpU=
+github.com/onsi/gomega v1.36.2 h1:koNYke6TVk6ZmnyHrCXba/T/MoLBXFjeC1PtvYgw0A8=
+github.com/onsi/gomega v1.36.2/go.mod h1:DdwyADRjrc825LhMEkD76cHR5+pUnjhUN8GlHlRPHzY=
 github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
 github.com/pascaldekloe/goe v0.1.0/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
@@ -169,22 +170,14 @@ github.com/posener/complete v1.2.3/go.mod h1:WZIdtGGp+qx0sLrYKtIRAruyNpv6hFCicSg
 github.com/prometheus/client_golang v0.9.1/go.mod h1:7SWBe2y4D6OKWSNQJUaRYU/AaXPKyh/dDVn+NZz0KFw=
 github.com/prometheus/client_golang v1.0.0/go.mod h1:db9x61etRT2tGnBNRi70OPL5FsnadC4Ky3P0J6CfImo=
 github.com/prometheus/client_golang v1.4.0/go.mod h1:e9GMxYsXl05ICDXkRhurwBS4Q3OK1iX/F2sw+iXX5zU=
-github.com/prometheus/client_golang v1.20.5 h1:cxppBPuYhUnsO6yo/aoRol4L7q7UFfdm+bR9r+8l63Y=
-github.com/prometheus/client_golang v1.20.5/go.mod h1:PIEt8X02hGcP8JWbeHyeZ53Y/jReSnHgO035n//V5WE=
 github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910/go.mod h1:MbSGuTsp3dbXC40dX6PRTWyKYBIrTGTE9sqQNg2J8bo=
 github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
 github.com/prometheus/client_model v0.2.0/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
-github.com/prometheus/client_model v0.6.1 h1:ZKSh/rekM+n3CeS952MLRAdFwIKqeY8b62p8ais2e9E=
-github.com/prometheus/client_model v0.6.1/go.mod h1:OrxVMOVHjw3lKMa8+x6HeMGkHMQyHDk9E3jmP2AmGiY=
 github.com/prometheus/common v0.4.1/go.mod h1:TNfzLD0ON7rHzMJeJkieUDPYmFC7Snx/y86RQel1bk4=
 github.com/prometheus/common v0.9.1/go.mod h1:yhUN8i9wzaXS3w1O07YhxHEBxD+W35wd8bs7vj7HSQ4=
-github.com/prometheus/common v0.55.0 h1:KEi6DK7lXW/m7Ig5i47x0vRzuBsHuvJdi5ee6Y3G1dc=
-github.com/prometheus/common v0.55.0/go.mod h1:2SECS4xJG1kd8XF9IcM1gMX6510RAEL65zxzNImwdc8=
 github.com/prometheus/procfs v0.0.0-20181005140218-185b4288413d/go.mod h1:c3At6R/oaqEKCNdg8wHV1ftS6bRYblBhIjjI8uT2IGk=
 github.com/prometheus/procfs v0.0.2/go.mod h1:TjEm7ze935MbeOT/UhFTIMYKhuLP4wbCsTZCD3I8kEA=
 github.com/prometheus/procfs v0.0.8/go.mod h1:7Qr8sr6344vo1JqZ6HhLceV9o3AJ1Ff+GxbHq6oeK9A=
-github.com/prometheus/procfs v0.15.1 h1:YagwOFzUgYfKKHX6Dr+sHT7km/hxC76UB0learggepc=
-github.com/prometheus/procfs v0.15.1/go.mod h1:fB45yRUv8NstnjriLhBQLuOUt+WW4BsoGhij/e3PBqk=
 github.com/rabbitmq/amqp091-go v1.10.0 h1:STpn5XsHlHGcecLmMFCtg7mqq0RnD+zFr4uzukfVhBw=
 github.com/rabbitmq/amqp091-go v1.10.0/go.mod h1:Hy4jKW5kQART1u+JkDTF9YYOQUHXqMuhrgxOEeS7G4o=
 github.com/rogpeppe/go-internal v1.10.0 h1:TMyTOH3F/DB16zRVcYyreMH6GnZZrwQVAoYjRBZyWFQ=
@@ -224,6 +217,8 @@ github.com/stretchr/testify v1.9.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8
 github.com/subosito/gotenv v1.6.0 h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=
 github.com/subosito/gotenv v1.6.0/go.mod h1:Dk4QP5c2W3ibzajGcXpNraDfq2IrhjMIvMSWPKKo0FU=
 github.com/tv42/httpunix v0.0.0-20150427012821-b75d8614f926/go.mod h1:9ESjWnEqriFuLhtthL60Sar/7RFoluCcXsuvEwTV5KM=
+github.com/x-cray/logrus-prefixed-formatter v0.5.2 h1:00txxvfBM9muc0jiLIEAkAcIMJzfthRT6usrui8uGmg=
+github.com/x-cray/logrus-prefixed-formatter v0.5.2/go.mod h1:2duySbKsL6M18s5GU7VPsoEPHyzalCE06qoARUCeBBE=
 go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 h1:4Pp6oUg3+e/6M4C0A/3kJ2VYa++dsWVTtGgLVj5xtHg=
 go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0/go.mod h1:Mjt1i1INqiaoZOMGR1RIUJN+i3ChKoFRqzrRQhlkbs0=
 go.opentelemetry.io/contrib/propagators/b3 v1.31.0 h1:PQPXYscmwbCp76QDvO4hMngF2j8Bx/OTV86laEl8uqo=
@@ -247,6 +242,8 @@ go.uber.org/multierr v1.9.0/go.mod h1:X2jQV1h+kxSjClGpnseKVIxpmcjrj7MNnI0bnlfKTV
 golang.org/x/crypto v0.0.0-20180904163835-0709b304e793/go.mod h1:6SG95UA2DQfeDnfUPMdvaQW0Q7yPrPDi9nlGo2tz2b4=
 golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACkg1iLfiJU5Ep61QUkGW8qpdssI0+w=
 golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392/go.mod h1:/lpIB1dKB+9EgE3H3cr1v9wB50oz8l4C4h62xy7jSTY=
+golang.org/x/crypto v0.31.0 h1:ihbySMvVjLAeSH1IbfcRTkD/iNscyz8rGzjF/E5hV6U=
+golang.org/x/crypto v0.31.0/go.mod h1:kDsLvtWBEx7MV9tJOj9bnXsPbxwJQ6csT/x4KIN4Ssk=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9 h1:GoHiUyI/Tp2nVkLI2mCxVkOjsbSXD66ic0XW0js0R9g=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9/go.mod h1:S2oDrQGGwySpoQPVqRShND87VCbxmc6bL1Yd2oYrm6k=
 golang.org/x/net v0.0.0-20181114220301-adae6a3d119a/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
@@ -256,8 +253,8 @@ golang.org/x/net v0.0.0-20190620200207-3b0461eec859/go.mod h1:z5CRVTTTmAJ677TzLL
 golang.org/x/net v0.0.0-20190923162816-aa69164e4478/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
 golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
 golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1/go.mod h1:9tjilg8BloeKEkVJvy7fQ90B1CfIiPueXVOjqfkSzI8=
-golang.org/x/net v0.30.0 h1:AcW1SDZMkb8IpzCdQUaIq2sP4sZ4zw+55h6ynffypl4=
-golang.org/x/net v0.30.0/go.mod h1:2wGyMJ5iFasEhkwi13ChkO/t1ECNC4X4eBKkVFyYFlU=
+golang.org/x/net v0.33.0 h1:74SYHlV8BIgHIFC/LrYkOGIwL19eTYXQ5wc6TBuO36I=
+golang.org/x/net v0.33.0/go.mod h1:HXLR5J+9DxmrqMwG9qjGCxZ+zKXxBru04zlTvWlWuN4=
 golang.org/x/sync v0.0.0-20181108010431-42b317875d0f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20181221193216-37e7f081c4d4/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20190423024810-112230192c58/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
@@ -284,16 +281,19 @@ golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6/go.mod h1:oPkhp1MJrh7nUepCBc
 golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20220908164124-27713097b956/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
-golang.org/x/sys v0.26.0 h1:KHjCJyddX0LoSTb3J+vWpupP9p0oznkqVk/IfjymZbo=
-golang.org/x/sys v0.26.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
+golang.org/x/sys v0.28.0 h1:Fksou7UEQUWlKvIdsqzJmUmCX3cZuD2+P3XyyzwMhlA=
+golang.org/x/sys v0.28.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
 golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
+golang.org/x/term v0.27.0 h1:WP60Sv1nlK1T6SupCHbXzSaN0b9wUmsPoRS9b61A23Q=
+golang.org/x/term v0.27.0/go.mod h1:iMsnZpn0cago0GOrHO2+Y7u7JPn5AylBrcoWkElMTSM=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
 golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
 golang.org/x/text v0.3.6/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
-golang.org/x/text v0.19.0 h1:kTxAhCbGbxhK0IwgSKiMO5awPoDQ0RpfiVYBfK860YM=
-golang.org/x/text v0.19.0/go.mod h1:BuEKDfySbSR4drPmRPG/7iBdf8hvFMuRexcpahXilzY=
+golang.org/x/text v0.21.0 h1:zyQAAkrwaneQ066sspRyJaG9VNi/YJ1NfzcGB3hZ/qo=
+golang.org/x/text v0.21.0/go.mod h1:4IBbMaMmOPCJ8SecivzSH54+73PCFmPWxNTLm+vZkEQ=
 golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190907020128-2ca718005c18/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
 golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
@@ -302,8 +302,8 @@ google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 h1:
 google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142/go.mod h1:UqMtugtsSgubUsoxbuAoiCXvqvErP7Gf0so0mK9tHxU=
 google.golang.org/grpc v1.67.1 h1:zWnc1Vrcno+lHZCOofnIMvycFcc0QRGIzm9dhnDX68E=
 google.golang.org/grpc v1.67.1/go.mod h1:1gLDyUQU7CTLJI90u3nXZ9ekeghjeM7pTDZlqFNg2AA=
-google.golang.org/protobuf v1.35.1 h1:m3LfL6/Ca+fqnjnlqQXNpFPABW1UD7mjh8KO2mKFytA=
-google.golang.org/protobuf v1.35.1/go.mod h1:9fA7Ob0pmnwhb644+1+CVWFRbNajQ6iRojtC/QF5bRE=
+google.golang.org/protobuf v1.36.1 h1:yBPeRvTftaleIgM3PZ/WBIZ7XM/eEYAaEyCwvyjq/gk=
+google.golang.org/protobuf v1.36.1/go.mod h1:9fA7Ob0pmnwhb644+1+CVWFRbNajQ6iRojtC/QF5bRE=
 gopkg.in/alecthomas/kingpin.v2 v2.2.6/go.mod h1:FMv+mEhP44yOT+4EoQTLFTRgOQ1FBLkstjWtayDeSgw=
 gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
@@ -311,6 +311,8 @@ gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c h1:Hei/4ADfdWqJk1ZMxUNpqntN
 gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c/go.mod h1:JHkPIbrfpd72SG/EVd6muEfDQjcINNoR0C8j2r3qZ4Q=
 gopkg.in/ini.v1 v1.67.0 h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=
 gopkg.in/ini.v1 v1.67.0/go.mod h1:pNLf8WUiyNEtQjuu5G5vTm06TEv9tsIgeAvK8hOrP4k=
+gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=
+gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7/go.mod h1:dt/ZhP58zS4L8KSrWDmTeBkI65Dw0HsyUHuEVlX15mw=
 gopkg.in/yaml.v2 v2.2.1/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.2/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.4/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
index 83a8a90..74a0319 100644
--- a/internal/order/app/command/create_order.go
+++ b/internal/order/app/command/create_order.go
@@ -2,7 +2,6 @@ package command
 
 import (
 	"context"
-	"encoding/json"
 	"fmt"
 
 	"github.com/ghost-yu/go_shop_second/common/broker"
@@ -67,13 +66,8 @@ func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*Creat
 	var err error
 	defer logging.WhenCommandExecute(ctx, "CreateOrderHandler", cmd, err)
 
-	q, err := c.channel.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
-	if err != nil {
-		return nil, err
-	}
-
 	t := otel.Tracer("rabbitmq")
-	ctx, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", q.Name))
+	ctx, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", broker.EventOrderCreated))
 	defer span.End()
 
 	validItems, err := c.validate(ctx, cmd.Items)
@@ -89,19 +83,15 @@ func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*Creat
 		return nil, err
 	}
 
-	marshalledOrder, err := json.Marshal(o)
-	if err != nil {
-		return nil, err
-	}
-	header := broker.InjectRabbitMQHeaders(ctx)
-	err = c.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
-		ContentType:  "application/json",
-		DeliveryMode: amqp.Persistent,
-		Body:         marshalledOrder,
-		Headers:      header,
+	err = broker.PublishEvent(ctx, broker.PublishEventReq{
+		Channel:  c.channel,
+		Routing:  broker.Direct,
+		Queue:    broker.EventOrderCreated,
+		Exchange: "",
+		Body:     o,
 	})
 	if err != nil {
-		return nil, errors.Wrapf(err, "publish event error q.Name=%s", q.Name)
+		return nil, errors.Wrapf(err, "publish event error q.Name=%s", broker.EventOrderCreated)
 	}
 
 	return &CreateOrderResult{OrderID: o.ID}, nil
diff --git a/internal/payment/go.mod b/internal/payment/go.mod
index 5d647ff..452ceaf 100644
--- a/internal/payment/go.mod
+++ b/internal/payment/go.mod
@@ -7,6 +7,7 @@ replace github.com/ghost-yu/go_shop_second/common => ../common
 require (
 	github.com/ghost-yu/go_shop_second/common v0.0.0-00010101000000-000000000000
 	github.com/gin-gonic/gin v1.10.0
+	github.com/pkg/errors v0.9.1
 	github.com/rabbitmq/amqp091-go v1.10.0
 	github.com/sirupsen/logrus v1.8.1
 	github.com/spf13/viper v1.19.0
diff --git a/internal/payment/http.go b/internal/payment/http.go
index 97ade12..0d4eddf 100644
--- a/internal/payment/http.go
+++ b/internal/payment/http.go
@@ -75,32 +75,23 @@ func (h *PaymentHandler) handleWebhook(c *gin.Context) {
 			var items []*orderpb.Item
 			_ = json.Unmarshal([]byte(session.Metadata["items"]), &items)
 
-			marshalledOrder, err := json.Marshal(&domain.Order{
-				ID:          session.Metadata["orderID"],
-				CustomerID:  session.Metadata["customerID"],
-				Status:      string(stripe.CheckoutSessionPaymentStatusPaid),
-				PaymentLink: session.Metadata["paymentLink"],
-				Items:       items,
-			})
-			if err != nil {
-				err = errors.Wrap(err, "error marshal domain.order")
-				c.JSON(http.StatusBadRequest, err.Error())
-				return
-			}
-
-			// TODO: mq logging
 			tr := otel.Tracer("rabbitmq")
 			ctx, span := tr.Start(c.Request.Context(), fmt.Sprintf("rabbitmq.%s.publish", broker.EventOrderPaid))
 			defer span.End()
 
-			headers := broker.InjectRabbitMQHeaders(ctx)
-			_ = h.channel.PublishWithContext(ctx, broker.EventOrderPaid, "", false, false, amqp.Publishing{
-				ContentType:  "application/json",
-				DeliveryMode: amqp.Persistent,
-				Body:         marshalledOrder,
-				Headers:      headers,
+			_ = broker.PublishEvent(ctx, broker.PublishEventReq{
+				Channel:  h.channel,
+				Routing:  broker.FanOut,
+				Queue:    "",
+				Exchange: broker.EventOrderPaid,
+				Body: &domain.Order{
+					ID:          session.Metadata["orderID"],
+					CustomerID:  session.Metadata["customerID"],
+					Status:      string(stripe.CheckoutSessionPaymentStatusPaid),
+					PaymentLink: session.Metadata["paymentLink"],
+					Items:       items,
+				},
 			})
-			logrus.WithContext(c).Infof("message published to %s, body: %s", broker.EventOrderPaid, string(marshalledOrder))
 		}
 	}
 	c.JSON(http.StatusOK, nil)
~~~

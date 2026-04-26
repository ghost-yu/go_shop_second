# Lesson Pair Diff Report

- FromBranch: lesson40
- ToBranch: lesson41

## Short Summary

~~~text
 19 files changed, 809 insertions(+), 46 deletions(-)
~~~

## File Stats

~~~text
 internal/common/decorator/logging.go               |   6 +-
 internal/common/go.mod                             |   1 -
 internal/common/go.sum                             |   2 -
 internal/kitchen/adapters/order_grpc_client.go     |  20 ++
 internal/kitchen/go.mod                            |  61 ++++-
 internal/kitchen/go.sum                            | 302 +++++++++++++++++++++
 .../kitchen/infrastructure/consumer/consumer.go    | 107 ++++++++
 internal/kitchen/main.go                           |  66 +++++
 internal/order/adapters/grpc/stock_grpc.go         |   4 +
 internal/order/tests/create_order_test.go          |   7 +-
 .../stock/app/query/check_if_items_in_stock.go     |  42 ++-
 internal/stock/convertor/convertor.go              |  97 +++++++
 internal/stock/convertor/facade.go                 |  39 +++
 internal/stock/entity/entity.go                    |  21 ++
 internal/stock/go.mod                              |   7 +-
 internal/stock/go.sum                              |  28 +-
 .../stock/infrastructure/integration/stripe.go     |  32 +++
 internal/stock/ports/grpc.go                       |   9 +-
 internal/stock/service/application.go              |   4 +-
 19 files changed, 809 insertions(+), 46 deletions(-)
~~~

## Commit Comparison

~~~text
> bdc08eb kitchen
~~~

## Changed Files

~~~text
internal/common/decorator/logging.go
internal/common/go.mod
internal/common/go.sum
internal/kitchen/adapters/order_grpc_client.go
internal/kitchen/go.mod
internal/kitchen/go.sum
internal/kitchen/infrastructure/consumer/consumer.go
internal/kitchen/main.go
internal/order/adapters/grpc/stock_grpc.go
internal/order/tests/create_order_test.go
internal/stock/app/query/check_if_items_in_stock.go
internal/stock/convertor/convertor.go
internal/stock/convertor/facade.go
internal/stock/entity/entity.go
internal/stock/go.mod
internal/stock/go.sum
internal/stock/infrastructure/integration/stripe.go
internal/stock/ports/grpc.go
internal/stock/service/application.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
internal/common/decorator/logging.go
internal/kitchen/adapters/order_grpc_client.go
internal/kitchen/infrastructure/consumer/consumer.go
internal/kitchen/main.go
internal/order/adapters/grpc/stock_grpc.go
internal/order/tests/create_order_test.go
internal/stock/app/query/check_if_items_in_stock.go
internal/stock/convertor/convertor.go
internal/stock/convertor/facade.go
internal/stock/entity/entity.go
internal/stock/infrastructure/integration/stripe.go
internal/stock/ports/grpc.go
internal/stock/service/application.go
~~~

## Full Diff

~~~diff
diff --git a/internal/common/decorator/logging.go b/internal/common/decorator/logging.go
index dbe402c..32f2f93 100644
--- a/internal/common/decorator/logging.go
+++ b/internal/common/decorator/logging.go
@@ -28,7 +28,8 @@ func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result
 			logger.Error("Failed to execute query", err)
 		}
 	}()
-	return q.base.Handle(ctx, cmd)
+	result, err = q.base.Handle(ctx, cmd)
+	return result, err
 }
 
 type commandLoggingDecorator[C, R any] struct {
@@ -50,7 +51,8 @@ func (q commandLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (resul
 			logger.Error("Failed to execute command", err)
 		}
 	}()
-	return q.base.Handle(ctx, cmd)
+	result, err = q.base.Handle(ctx, cmd)
+	return result, err
 }
 
 func generateActionName(cmd any) string {
diff --git a/internal/common/go.mod b/internal/common/go.mod
index 7bc7434..9ed3597 100644
--- a/internal/common/go.mod
+++ b/internal/common/go.mod
@@ -10,7 +10,6 @@ require (
 	github.com/rabbitmq/amqp091-go v1.10.0
 	github.com/sirupsen/logrus v1.8.1
 	github.com/spf13/viper v1.19.0
-	github.com/zenazn/goji v1.0.1
 	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0
 	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0
 	go.opentelemetry.io/contrib/propagators/b3 v1.31.0
diff --git a/internal/common/go.sum b/internal/common/go.sum
index 9696ba5..9b5cc9c 100644
--- a/internal/common/go.sum
+++ b/internal/common/go.sum
@@ -289,8 +289,6 @@ github.com/ugorji/go/codec v1.2.12 h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65E
 github.com/ugorji/go/codec v1.2.12/go.mod h1:UNopzCgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=
 github.com/yuin/goldmark v1.1.27/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
-github.com/zenazn/goji v1.0.1 h1:4lbD8Mx2h7IvloP7r2C0D6ltZP6Ufip8Hn0wmSK5LR8=
-github.com/zenazn/goji v1.0.1/go.mod h1:7S9M489iMyHBNxwZnk9/EHS098H4/F6TATF2mIxtB1Q=
 go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 h1:0nTRpaCaILLdooXAQnfktlL6Zw1ECKEW9DZGH2byi2c=
 go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0/go.mod h1:A7aFlp4WSLmeOnFRZwf2dMU+40THPc+rsr6KOwZLOcg=
 go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 h1:4Pp6oUg3+e/6M4C0A/3kJ2VYa++dsWVTtGgLVj5xtHg=
diff --git a/internal/kitchen/adapters/order_grpc_client.go b/internal/kitchen/adapters/order_grpc_client.go
new file mode 100644
index 0000000..92402ab
--- /dev/null
+++ b/internal/kitchen/adapters/order_grpc_client.go
@@ -0,0 +1,20 @@
+package adapters
+
+import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+)
+
+type OrderGRPC struct {
+	client orderpb.OrderServiceClient
+}
+
+func NewOrderGRPC(client orderpb.OrderServiceClient) *OrderGRPC {
+	return &OrderGRPC{client: client}
+}
+
+func (g *OrderGRPC) UpdateOrder(ctx context.Context, request *orderpb.Order) error {
+	_, err := g.client.UpdateOrder(ctx, request)
+	return err
+}
diff --git a/internal/kitchen/go.mod b/internal/kitchen/go.mod
index 5449764..d68df76 100644
--- a/internal/kitchen/go.mod
+++ b/internal/kitchen/go.mod
@@ -2,6 +2,61 @@ module github.com/ghost-yu/go_shop_second/kitchen
 
 go 1.22.8
 
-replace (
-	github.com/ghost-yu/go_shop_second/common => ../common
-)
\ No newline at end of file
+replace github.com/ghost-yu/go_shop_second/common => ../common
+
+require (
+	github.com/ghost-yu/go_shop_second/common v0.0.0-00010101000000-000000000000
+	github.com/rabbitmq/amqp091-go v1.10.0
+	github.com/sirupsen/logrus v1.9.3
+	github.com/spf13/viper v1.19.0
+	go.opentelemetry.io/otel v1.31.0
+)
+
+require (
+	github.com/armon/go-metrics v0.4.1 // indirect
+	github.com/fatih/color v1.14.1 // indirect
+	github.com/fsnotify/fsnotify v1.7.0 // indirect
+	github.com/go-logr/logr v1.4.2 // indirect
+	github.com/go-logr/stdr v1.2.2 // indirect
+	github.com/google/uuid v1.6.0 // indirect
+	github.com/hashicorp/consul/api v1.28.2 // indirect
+	github.com/hashicorp/errwrap v1.1.0 // indirect
+	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
+	github.com/hashicorp/go-hclog v1.5.0 // indirect
+	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
+	github.com/hashicorp/go-multierror v1.1.1 // indirect
+	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
+	github.com/hashicorp/golang-lru v0.5.4 // indirect
+	github.com/hashicorp/hcl v1.0.0 // indirect
+	github.com/hashicorp/serf v0.10.1 // indirect
+	github.com/magiconair/properties v1.8.7 // indirect
+	github.com/mattn/go-colorable v0.1.13 // indirect
+	github.com/mattn/go-isatty v0.0.20 // indirect
+	github.com/mitchellh/go-homedir v1.1.0 // indirect
+	github.com/mitchellh/mapstructure v1.5.0 // indirect
+	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
+	github.com/sagikazarmark/locafero v0.4.0 // indirect
+	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
+	github.com/sourcegraph/conc v0.3.0 // indirect
+	github.com/spf13/afero v1.11.0 // indirect
+	github.com/spf13/cast v1.6.0 // indirect
+	github.com/spf13/pflag v1.0.5 // indirect
+	github.com/subosito/gotenv v1.6.0 // indirect
+	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
+	go.opentelemetry.io/contrib/propagators/b3 v1.31.0 // indirect
+	go.opentelemetry.io/otel/exporters/jaeger v1.17.0 // indirect
+	go.opentelemetry.io/otel/metric v1.31.0 // indirect
+	go.opentelemetry.io/otel/sdk v1.31.0 // indirect
+	go.opentelemetry.io/otel/trace v1.31.0 // indirect
+	go.uber.org/atomic v1.9.0 // indirect
+	go.uber.org/multierr v1.9.0 // indirect
+	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
+	golang.org/x/net v0.30.0 // indirect
+	golang.org/x/sys v0.26.0 // indirect
+	golang.org/x/text v0.19.0 // indirect
+	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
+	google.golang.org/grpc v1.67.1 // indirect
+	google.golang.org/protobuf v1.35.1 // indirect
+	gopkg.in/ini.v1 v1.67.0 // indirect
+	gopkg.in/yaml.v3 v3.0.1 // indirect
+)
diff --git a/internal/common/go.sum b/internal/kitchen/go.sum
similarity index 66%
copy from internal/common/go.sum
copy to internal/kitchen/go.sum
index 9696ba5..b98b2e2 100644
--- a/internal/common/go.sum
+++ b/internal/kitchen/go.sum
@@ -1,48 +1,25 @@
-cloud.google.com/go v0.26.0/go.mod h1:aQUYkXzVsufM+DwF1aE+0xfcU+56JwCaLick0ClmMTw=
-github.com/BurntSushi/toml v0.3.1/go.mod h1:xHWCNGjB5oqiDr8zfno3MHue2Ht5sIBksp03qcyfWMU=
 github.com/DataDog/datadog-go v3.2.0+incompatible/go.mod h1:LButxg5PwREeZtORoXG3tL4fMGNddJ+vMq1mwgfaqoQ=
-github.com/RaveNoX/go-jsoncommentstrip v1.0.0/go.mod h1:78ihd09MekBnJnxpICcwzCMzGrKSKYe4AqU6PDYYpjk=
 github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
 github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
 github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
 github.com/alecthomas/units v0.0.0-20190717042225-c3de453c63f4/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
-github.com/apapsch/go-jsonmerge/v2 v2.0.0 h1:axGnT1gRIfimI7gJifB699GoE/oq+F2MU7Dml6nw9rQ=
-github.com/apapsch/go-jsonmerge/v2 v2.0.0/go.mod h1:lvDnEdqiQrp0O42VQGgmlKpxL1AP2+08jFMw88y4klk=
 github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e/go.mod h1:3U/XgcO3hCbHZ8TKRvWD2dDTCfh9M9ya+I9JpbB7O8o=
 github.com/armon/go-metrics v0.0.0-20180917152333-f0300d1749da/go.mod h1:Q73ZrmVTwzkszR9V5SSuryQ31EELlFMUz1kKyl939pY=
 github.com/armon/go-metrics v0.4.1 h1:hR91U9KYmb6bLBYLQjyM+3j+rcd/UhE+G78SFnF8gJA=
 github.com/armon/go-metrics v0.4.1/go.mod h1:E6amYzXo6aW1tqzoZGT755KkbgrJsSdpwZ+3JqfkOG4=
 github.com/armon/go-radix v0.0.0-20180808171621-7fddfc383310/go.mod h1:ufUuZ+zHj4x4TnLV4JWEpy2hxWSpsRywHrMgIH9cCH8=
 github.com/armon/go-radix v1.0.0/go.mod h1:ufUuZ+zHj4x4TnLV4JWEpy2hxWSpsRywHrMgIH9cCH8=
-github.com/benbjohnson/clock v1.1.0/go.mod h1:J11/hYXuz8f4ySSvYwY0FKfm+ezbsZBKZxNJlLklBHA=
 github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973/go.mod h1:Dwedo/Wpr24TaqPxmxbtue+5NUziq4I4S80YR8gNf3Q=
 github.com/beorn7/perks v1.0.0/go.mod h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+CedLV8=
 github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
 github.com/bgentry/speakeasy v0.1.0/go.mod h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=
-github.com/bmatcuk/doublestar v1.1.1/go.mod h1:UD6OnuiIn0yFxxA2le/rnRU1G4RaI4UvFv1sNto9p6w=
-github.com/bytedance/sonic v1.12.3 h1:W2MGa7RCU1QTeYRTPE3+88mVC0yXmsRQRChiyVocVjU=
-github.com/bytedance/sonic v1.12.3/go.mod h1:B8Gt/XvtZ3Fqj+iSKMypzymZxw/FVwgIGKzMzT9r/rk=
-github.com/bytedance/sonic/loader v0.1.1/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
-github.com/bytedance/sonic/loader v0.2.0 h1:zNprn+lsIP06C/IqCHs3gPQIvnvpKbbxyXQP1iU4kWM=
-github.com/bytedance/sonic/loader v0.2.0/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
-github.com/census-instrumentation/opencensus-proto v0.2.1/go.mod h1:f6KPmirojxKA12rnyqOA5BBL4O983OfeGPqjHWSTneU=
 github.com/cespare/xxhash/v2 v2.1.1/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
 github.com/circonus-labs/circonus-gometrics v2.3.1+incompatible/go.mod h1:nmEj6Dob7S7YxXgwXpfOuvO54S+tGdZdw9fuRZt25Ag=
 github.com/circonus-labs/circonusllhist v0.1.3/go.mod h1:kMXHVDlOchFAehlya5ePtbp5jckzBHf4XRpQvBOLI+I=
-github.com/client9/misspell v0.3.4/go.mod h1:qj6jICC3Q7zFZvVWo7KLAzC3yx5G7kyvSDkc90ppPyw=
-github.com/cloudwego/base64x v0.1.4 h1:jwCgWpFanWmN8xoIUHa2rtzmkd5J2plF/dnLS6Xd/0Y=
-github.com/cloudwego/base64x v0.1.4/go.mod h1:0zlkT4Wn5C6NdauXdJRhSKRlJvmclQ1hhJgA0rcu/8w=
-github.com/cloudwego/iasm v0.2.0 h1:1KNIy1I1H9hNNFEEH3DVnI4UujN+1zjpuk6gwHLTssg=
-github.com/cloudwego/iasm v0.2.0/go.mod h1:8rXZaNYT2n95jn+zTI1sDr+IgcD2GVs0nlbbQPiEFhY=
-github.com/cncf/udpa/go v0.0.0-20191209042840-269d4d468f6f/go.mod h1:M8M6+tZqaGXZJjfX53e64911xZQV5JYwmTeXPW+k8Sc=
 github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
 github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
 github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc h1:U9qPSI2PIWSS1VwoXQT9A3Wy9MM3WgvqSxFWenqJduM=
 github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
-github.com/envoyproxy/go-control-plane v0.9.0/go.mod h1:YTl/9mNaCwkRvm6d1a2C3ymFceY/DCBVvsKhRF0iEA4=
-github.com/envoyproxy/go-control-plane v0.9.1-0.20191026205805-5f8ba28d4473/go.mod h1:YTl/9mNaCwkRvm6d1a2C3ymFceY/DCBVvsKhRF0iEA4=
-github.com/envoyproxy/go-control-plane v0.9.4/go.mod h1:6rpuAdCZL397s3pYoYcLgu1mIlRU8Am5FuJP05cCM98=
-github.com/envoyproxy/protoc-gen-validate v0.1.0/go.mod h1:iSmxcyjqTsJpI2R4NaDN7+kN2VEUnK/pcBlmesArF7c=
 github.com/fatih/color v1.7.0/go.mod h1:Zm6kSWBoL9eyXnKyktHP6abPY2pDugNf5KwzbycvMj4=
 github.com/fatih/color v1.9.0/go.mod h1:eQcE1qtQxscV5RaZvpXrrb8Drkc3/DdQ+uUYCNjL+zU=
 github.com/fatih/color v1.13.0/go.mod h1:kLAiJbzzSOZDVNGyDpeOxJ47H46qBXwg5ILebYFFOfk=
@@ -52,60 +29,30 @@ github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHk
 github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
-github.com/gabriel-vasile/mimetype v1.4.5 h1:J7wGKdGu33ocBOhGy0z653k/lFKLFDPJMG8Gql0kxn4=
-github.com/gabriel-vasile/mimetype v1.4.5/go.mod h1:ibHel+/kbxn9x2407k1izTA1S81ku1z/DlgOW2QE0M4=
-github.com/gin-contrib/sse v0.1.0 h1:Y/yl/+YNO8GZSjAhjMsSuLt29uWRFHdHYUb5lYOV9qE=
-github.com/gin-contrib/sse v0.1.0/go.mod h1:RHrZQHXnP2xjPF+u1gW/2HnVO7nvIa9PG3Gm+fLHvGI=
-github.com/gin-gonic/gin v1.10.0 h1:nTuyha1TYqgedzytsKYqna+DfLos46nTv2ygFy86HFU=
-github.com/gin-gonic/gin v1.10.0/go.mod h1:4PMNQiOhvDRa013RKVbsiNwoyezlm2rm0uX/T7kzp5Y=
 github.com/go-kit/kit v0.8.0/go.mod h1:xBxKIO96dXMWWy0MnWVtmwkA9/13aqxPnvrjFYMA2as=
 github.com/go-kit/kit v0.9.0/go.mod h1:xBxKIO96dXMWWy0MnWVtmwkA9/13aqxPnvrjFYMA2as=
-github.com/go-kit/log v0.1.0/go.mod h1:zbhenjAZHb184qTLMA9ZjW7ThYL0H2mk7Q6pNt4vbaY=
 github.com/go-logfmt/logfmt v0.3.0/go.mod h1:Qt1PoO58o5twSAckw1HlFXLmHsOX5/0LbT9GBnD5lWE=
 github.com/go-logfmt/logfmt v0.4.0/go.mod h1:3RMwSq7FuexP4Kalkev3ejPJsZTpXXBr9+V4qmtdjCk=
-github.com/go-logfmt/logfmt v0.5.0/go.mod h1:wCYkCAKZfumFQihp8CzCvQ3paCTfi41vtzG1KdI/P7A=
 github.com/go-logr/logr v1.2.2/go.mod h1:jdQByPbusPIv2/zmleS9BjJVeZ6kBagPoEUsqbVz/1A=
 github.com/go-logr/logr v1.4.2 h1:6pFjapn8bFcIbiKo3XT4j/BhANplGihG6tvd+8rYgrY=
 github.com/go-logr/logr v1.4.2/go.mod h1:9T104GzyrTigFIr8wt5mBrctHMim0Nb2HLGrmQ40KvY=
 github.com/go-logr/stdr v1.2.2 h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=
 github.com/go-logr/stdr v1.2.2/go.mod h1:mMo/vtBO5dYbehREoey6XUKy/eSumjCCveDpRre4VKE=
-github.com/go-playground/assert/v2 v2.2.0 h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=
-github.com/go-playground/assert/v2 v2.2.0/go.mod h1:VDjEfimB/XKnb+ZQfWdccd7VUvScMdVu0Titje2rxJ4=
-github.com/go-playground/locales v0.14.1 h1:EWaQ/wswjilfKLTECiXz7Rh+3BjFhfDFKv/oXslEjJA=
-github.com/go-playground/locales v0.14.1/go.mod h1:hxrqLVvrK65+Rwrd5Fc6F2O76J/NuW9t0sjnWqG1slY=
-github.com/go-playground/universal-translator v0.18.1 h1:Bcnm0ZwsGyWbCzImXv+pAJnYK9S473LQFuzCbDbfSFY=
-github.com/go-playground/universal-translator v0.18.1/go.mod h1:xekY+UJKNuX9WP91TpwSH2VMlDf28Uj24BCp08ZFTUY=
-github.com/go-playground/validator/v10 v10.22.1 h1:40JcKH+bBNGFczGuoBYgX4I6m/i27HYW8P9FDk5PbgA=
-github.com/go-playground/validator/v10 v10.22.1/go.mod h1:dbuPbCMFw/DrkbEynArYaCwl3amGuJotoKCe95atGMM=
 github.com/go-stack/stack v1.8.0/go.mod h1:v0f6uXyyMGvRgIKkXu+yp6POWl0qKG85gN/melR3HDY=
-github.com/goccy/go-json v0.10.3 h1:KZ5WoDbxAIgm2HNbYckL0se1fHD6rz5j4ywS6ebzDqA=
-github.com/goccy/go-json v0.10.3/go.mod h1:oq7eo15ShAhp70Anwd5lgX2pLfOS3QCiwU/PULtXL6M=
 github.com/gogo/protobuf v1.1.1/go.mod h1:r8qH/GZQm5c6nD/R0oafs1akxWv10x8SbQlK7atdtwQ=
-github.com/gogo/protobuf v1.3.2 h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=
-github.com/gogo/protobuf v1.3.2/go.mod h1:P1XiOD3dCwIKUDQYPy72D8LYyHL2YPYrpS2s69NZV8Q=
-github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b/go.mod h1:SBH7ygxi8pfUlaOkMMuAQtPIUF8ecWP5IEl/CR7VP2Q=
-github.com/golang/mock v1.1.1/go.mod h1:oTYuIxOrZwtPieC+H1uAHpcLFnEyAGVDL/k47Jfbm0A=
 github.com/golang/protobuf v1.2.0/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5yJMmIC1U=
 github.com/golang/protobuf v1.3.1/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5yJMmIC1U=
 github.com/golang/protobuf v1.3.2/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5yJMmIC1U=
-github.com/golang/protobuf v1.3.3/go.mod h1:vzj43D7+SQXF/4pzW/hwtAqwc6iTitCiVSaWz5lYuqw=
-github.com/golang/protobuf v1.5.0/go.mod h1:FsONVRAS9T7sI+LIUmWTfcYkHO4aIWwzhcaSAoJOfIk=
-github.com/golang/protobuf v1.5.3 h1:KhyjKVUg7Usr/dYsdSqoFveMYd5ko72D+zANwlG1mmg=
-github.com/golang/protobuf v1.5.3/go.mod h1:XVQd3VNwM+JqD3oG2Ue2ip4fOMUkwXdXDdiuN0vRsmY=
 github.com/google/btree v0.0.0-20180813153112-4030bb1f1f0c/go.mod h1:lNA+9X1NB3Zf8V7Ke586lFgjr2dZNuvo3lPJSGZ5JPQ=
 github.com/google/btree v1.0.1 h1:gK4Kx5IaGY9CD5sPJ36FHiBJ6ZXl0kilRiiCj+jdYp4=
 github.com/google/btree v1.0.1/go.mod h1:xXMiIv4Fb/0kKde4SpL7qlzvu5cMJDRkFDxJfI9uaxA=
-github.com/google/go-cmp v0.2.0/go.mod h1:oXzfMopK8JAjlY9xF4vHSVASa0yLyX7SntLO5aqRK0M=
 github.com/google/go-cmp v0.3.1/go.mod h1:8QqcDgzrUqlUb/G2PQTWiueGozuR1884gddMywk6iLU=
 github.com/google/go-cmp v0.4.0/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
-github.com/google/go-cmp v0.5.5/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
 github.com/google/go-cmp v0.6.0 h1:ofyhxvXcZhMsU5ulbFiLKl/XBFqE1GSq7atu8tAmTRI=
 github.com/google/go-cmp v0.6.0/go.mod h1:17dUlkBOakJ0+DkrSSNjCkIjxS6bF9zb3elmeNGIjoY=
 github.com/google/gofuzz v1.0.0/go.mod h1:dBl0BpW6vV/+mYPU4Po3pmUjxk6FQPldtuIdl/M65Eg=
 github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
 github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
-github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 h1:UH//fgunKIs4JdUbpDl1VZCDaL56wXCB/5+wF6uHfaI=
-github.com/grpc-ecosystem/go-grpc-middleware v1.4.0/go.mod h1:g5qyo/la0ALbONm6Vbp88Yd8NsDy6rZz+RcrMPxvld8=
 github.com/hashicorp/consul/api v1.28.2 h1:mXfkRHrpHN4YY3RqL09nXU1eHKLNiuAN4kHvDQ16k/8=
 github.com/hashicorp/consul/api v1.28.2/go.mod h1:KyzqzgMEya+IZPcD65YFoOVAgPpbfERu4I/tzG6/ueE=
 github.com/hashicorp/consul/sdk v0.16.0 h1:SE9m0W6DEfgIVCJX7xU+iv/hUl4m/nxqMTnCdMxDpJ8=
@@ -154,16 +101,7 @@ github.com/hashicorp/serf v0.10.1 h1:Z1H2J60yRKvfDYAOZLd2MU0ND4AH/WDz7xYHDWQsIPY
 github.com/hashicorp/serf v0.10.1/go.mod h1:yL2t6BqATOLGc5HF7qbFkTfXoPIY0WZdWHfEvMqbG+4=
 github.com/json-iterator/go v1.1.6/go.mod h1:+SdeFBvtyEkXs7REEP0seUULqWtbJapLOCVDaaPEHmU=
 github.com/json-iterator/go v1.1.9/go.mod h1:KdQUCv79m/52Kvf8AW2vK1V8akMuk1QjK/uOdHXbAo4=
-github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
-github.com/json-iterator/go v1.1.12/go.mod h1:e30LSqwooZae/UwlEbR2852Gd8hjQvJoHmT4TnhNGBo=
-github.com/juju/gnuflag v0.0.0-20171113085948-2ce1bb71843d/go.mod h1:2PavIy+JPciBPrBUjwbNvtwB6RQlve+hkpll6QSNmOE=
 github.com/julienschmidt/httprouter v1.2.0/go.mod h1:SYymIcj16QtmaHHD7aYtjjsJG7VTCxuUUipMqKk8s4w=
-github.com/kisielk/errcheck v1.5.0/go.mod h1:pFxgyoBC7bSaBwPgfKdkLd5X25qrDl4LWUI2bnpBCr8=
-github.com/kisielk/gotool v1.0.0/go.mod h1:XhKaO+MFFWcvkIS/tQcRk01m1F5IRFswLeQ+oQHNcck=
-github.com/klauspost/cpuid/v2 v2.0.9/go.mod h1:FInQzS24/EEf25PyTYn52gqo7WaD8xa0213Md/qVLRg=
-github.com/klauspost/cpuid/v2 v2.2.8 h1:+StwCXwm9PdpiEkPyzBXIy+M9KUb4ODm0Zarf1kS5BM=
-github.com/klauspost/cpuid/v2 v2.2.8/go.mod h1:Lcz8mBdAVJIBVzewtcLocK12l3Y+JytZYpaMropDUws=
-github.com/knz/go-libedit v1.10.1/go.mod h1:MZTVkCWyz0oBc7JOWP3wNAzd002ZbM/5hgShxwh4x8M=
 github.com/konsorten/go-windows-terminal-sequences v1.0.1/go.mod h1:T0+1ngSBFLxvqU3pZ+m/2kptfBszLMUkC4ZK/EgS/cQ=
 github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515/go.mod h1:+0opPa2QZZtGFBFZlji/RkVcI2GknAs/DXo4wKdlNEc=
 github.com/kr/pretty v0.1.0/go.mod h1:dAy3ld7l9f0ibDNOQOHHMYYIIbhfbHSm3C4ZsoJORNo=
@@ -173,8 +111,6 @@ github.com/kr/pty v1.1.1/go.mod h1:pFQYn66WHrOpPYNljwOMqo10TkYh1fy3cYio2l3bCsQ=
 github.com/kr/text v0.1.0/go.mod h1:4Jbv+DJW3UT/LiOwJeYQe1efqtUx/iVham/4vfdArNI=
 github.com/kr/text v0.2.0 h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=
 github.com/kr/text v0.2.0/go.mod h1:eLer722TekiGuMkidMxC/pM04lWEeraHUUmBw8l2grE=
-github.com/leodido/go-urn v1.4.0 h1:WT9HwE9SGECu3lg4d/dIA+jxlljEa1/ffXKmRjqdmIQ=
-github.com/leodido/go-urn v1.4.0/go.mod h1:bvxc+MVxLKB4z00jd1z+Dvzr47oO32F/QSNjSBOlFxI=
 github.com/magiconair/properties v1.8.7 h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=
 github.com/magiconair/properties v1.8.7/go.mod h1:Dhd985XPs7jluiymwWYZ0G4Z61jb3vdS329zhj2hYo0=
 github.com/mattn/go-colorable v0.0.9/go.mod h1:9vuHe8Xs5qXnSaW/c/ABM9alt+Vo+STaOChaDxuIBZU=
@@ -203,16 +139,10 @@ github.com/mitchellh/mapstructure v0.0.0-20160808181253-ca63d7c062ee/go.mod h1:F
 github.com/mitchellh/mapstructure v1.5.0 h1:jeMsZIYE/09sWLaz43PL7Gy6RuMjD2eJVyuac5Z2hdY=
 github.com/mitchellh/mapstructure v1.5.0/go.mod h1:bFUtVrKA4DC2yAKiSyO/QUcy7e+RRV2QTWOzhPopBRo=
 github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
-github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=
 github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
 github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3RllmbCylyMrvgv0=
 github.com/modern-go/reflect2 v1.0.1/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3RllmbCylyMrvgv0=
-github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
-github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
 github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
-github.com/oapi-codegen/runtime v1.1.1 h1:EXLHh0DXIJnWhdRPN2w4MXAzFyE4CskzhNLUmtpMYro=
-github.com/oapi-codegen/runtime v1.1.1/go.mod h1:SK9X900oXmPWilYR5/WKPzt3Kqxn/uS/+lbpREv+eCg=
-github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFStlNbqXla1AfSYxPUl2o=
 github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
 github.com/pascaldekloe/goe v0.1.0/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
@@ -232,7 +162,6 @@ github.com/prometheus/client_golang v1.0.0/go.mod h1:db9x61etRT2tGnBNRi70OPL5Fsn
 github.com/prometheus/client_golang v1.4.0/go.mod h1:e9GMxYsXl05ICDXkRhurwBS4Q3OK1iX/F2sw+iXX5zU=
 github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910/go.mod h1:MbSGuTsp3dbXC40dX6PRTWyKYBIrTGTE9sqQNg2J8bo=
 github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
-github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
 github.com/prometheus/client_model v0.2.0/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
 github.com/prometheus/common v0.4.1/go.mod h1:TNfzLD0ON7rHzMJeJkieUDPYmFC7Snx/y86RQel1bk4=
 github.com/prometheus/common v0.9.1/go.mod h1:yhUN8i9wzaXS3w1O07YhxHEBxD+W35wd8bs7vj7HSQ4=
@@ -252,8 +181,8 @@ github.com/sean-/seed v0.0.0-20170313163322-e2103e2c3529 h1:nn5Wsu0esKSJiIVhscUt
 github.com/sean-/seed v0.0.0-20170313163322-e2103e2c3529/go.mod h1:DxrIzT+xaE7yg65j358z/aeFdxmN0P9QXhEzd20vsDc=
 github.com/sirupsen/logrus v1.2.0/go.mod h1:LxeOpSwHxABJmUn/MG1IvRgCAasNZTLOkJPxbbu5VWo=
 github.com/sirupsen/logrus v1.4.2/go.mod h1:tLMulIdttU9McNUspp0xgXVQah82FyeX6MwdIuYE2rE=
-github.com/sirupsen/logrus v1.8.1 h1:dJKuHgqk1NNQlqoA6BTlM1Wf9DOH3NBjQyu0h9+AZZE=
-github.com/sirupsen/logrus v1.8.1/go.mod h1:yWOB1SBYBC5VeMP7gHvWumXLIWorT60ONWic61uBYv0=
+github.com/sirupsen/logrus v1.9.3 h1:dueUQJ1C2q9oE3F7wvmSGAaVtTmUizReu6fjN8uqzbQ=
+github.com/sirupsen/logrus v1.9.3/go.mod h1:naHLuLoDiP4jHNo9R0sCBMtWGeIprob74mVsIT4qYEQ=
 github.com/sourcegraph/conc v0.3.0 h1:OQTbbt6P72L20UqAkXXuLOj79LfEanQ+YQFNpLA9ySo=
 github.com/sourcegraph/conc v0.3.0/go.mod h1:Sdozi7LEKbFPqYX2/J+iBAM6HpqSLTASQIKqDmF7Mt0=
 github.com/spf13/afero v1.11.0 h1:WJQKhtpdm3v2IzqG8VMqrr6Rf3UYpEF239Jy9wNepM8=
@@ -264,35 +193,20 @@ github.com/spf13/pflag v1.0.5 h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=
 github.com/spf13/pflag v1.0.5/go.mod h1:McXfInJRrz4CZXVZOBLb0bTZqETkiAhM9Iw0y3An2Bg=
 github.com/spf13/viper v1.19.0 h1:RWq5SEjt8o25SROyN3z2OrDB9l7RPd3lwTWU8EcEdcI=
 github.com/spf13/viper v1.19.0/go.mod h1:GQUN9bilAbhU/jgc1bKs99f/suXKeUMct8Adx5+Ntkg=
-github.com/spkg/bom v0.0.0-20160624110644-59b7046e48ad/go.mod h1:qLr4V1qq6nMqFKkMo8ZTx3f+BZEkzsRUY10Xsm2mwU0=
 github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.1.1/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
-github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
 github.com/stretchr/objx v0.5.0 h1:1zr/of2m5FGMsad5YfcqgdqdWrIhu+EBEJRhR1U7z/c=
 github.com/stretchr/objx v0.5.0/go.mod h1:Yh+to48EsGEfYuaHDzXPcE3xhTkx73EhmCGUpEOglKo=
 github.com/stretchr/testify v1.2.2/go.mod h1:a8OnRcib4nhh0OaRAV+Yts87kKdq0PP7pXfy6kDkUVs=
 github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
 github.com/stretchr/testify v1.4.0/go.mod h1:j7eGeouHqKxXV5pUuKE4zz7dFj8WfuZ+81PSLYec5m4=
 github.com/stretchr/testify v1.7.0/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
-github.com/stretchr/testify v1.7.1/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
 github.com/stretchr/testify v1.7.2/go.mod h1:R6va5+xMeoiuVRoj+gSkQ7d3FALtqAAGI1FQKckRals=
-github.com/stretchr/testify v1.8.0/go.mod h1:yNjHg4UonilssWZ8iaSj1OCr/vHnekPRkoO+kdMU+MU=
-github.com/stretchr/testify v1.8.1/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o6fzry7u4=
 github.com/stretchr/testify v1.9.0 h1:HtqpIVDClZ4nwg75+f6Lvsy/wHu+3BoSGCbBAcpTsTg=
 github.com/stretchr/testify v1.9.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8C91i36aY=
 github.com/subosito/gotenv v1.6.0 h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=
 github.com/subosito/gotenv v1.6.0/go.mod h1:Dk4QP5c2W3ibzajGcXpNraDfq2IrhjMIvMSWPKKo0FU=
 github.com/tv42/httpunix v0.0.0-20150427012821-b75d8614f926/go.mod h1:9ESjWnEqriFuLhtthL60Sar/7RFoluCcXsuvEwTV5KM=
-github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS4MhqMhdFk5YI=
-github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
-github.com/ugorji/go/codec v1.2.12 h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65EE=
-github.com/ugorji/go/codec v1.2.12/go.mod h1:UNopzCgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=
-github.com/yuin/goldmark v1.1.27/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
-github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
-github.com/zenazn/goji v1.0.1 h1:4lbD8Mx2h7IvloP7r2C0D6ltZP6Ufip8Hn0wmSK5LR8=
-github.com/zenazn/goji v1.0.1/go.mod h1:7S9M489iMyHBNxwZnk9/EHS098H4/F6TATF2mIxtB1Q=
-go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 h1:0nTRpaCaILLdooXAQnfktlL6Zw1ECKEW9DZGH2byi2c=
-go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0/go.mod h1:A7aFlp4WSLmeOnFRZwf2dMU+40THPc+rsr6KOwZLOcg=
 go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 h1:4Pp6oUg3+e/6M4C0A/3kJ2VYa++dsWVTtGgLVj5xtHg=
 go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0/go.mod h1:Mjt1i1INqiaoZOMGR1RIUJN+i3ChKoFRqzrRQhlkbs0=
 go.opentelemetry.io/contrib/propagators/b3 v1.31.0 h1:PQPXYscmwbCp76QDvO4hMngF2j8Bx/OTV86laEl8uqo=
@@ -307,64 +221,36 @@ go.opentelemetry.io/otel/sdk v1.31.0 h1:xLY3abVHYZ5HSfOg3l2E5LUj2Cwva5Y7yGxnSW9H
 go.opentelemetry.io/otel/sdk v1.31.0/go.mod h1:TfRbMdhvxIIr/B2N2LQW2S5v9m3gOQ/08KsbbO5BPT0=
 go.opentelemetry.io/otel/trace v1.31.0 h1:ffjsj1aRouKewfr85U2aGagJ46+MvodynlQ1HYdmJys=
 go.opentelemetry.io/otel/trace v1.31.0/go.mod h1:TXZkRk7SM2ZQLtR6eoAWQFIHPvzQ06FJAsO1tJg480A=
-go.uber.org/atomic v1.7.0/go.mod h1:fEN4uk6kAWBTFdckzkM89CLk9XfWZrxpCo0nPH17wJc=
 go.uber.org/atomic v1.9.0 h1:ECmE8Bn/WFTYwEW/bpKD3M8VtR/zQVbavAoalC1PYyE=
 go.uber.org/atomic v1.9.0/go.mod h1:fEN4uk6kAWBTFdckzkM89CLk9XfWZrxpCo0nPH17wJc=
-go.uber.org/goleak v1.1.10/go.mod h1:8a7PlsEVH3e/a/GLqe5IIrQx6GzcnRmZEufDUTk4A7A=
 go.uber.org/goleak v1.3.0 h1:2K3zAYmnTNqV73imy9J1T3WC+gmCePx2hEGkimedGto=
 go.uber.org/goleak v1.3.0/go.mod h1:CoHD4mav9JJNrW/WLlf7HGZPjdw8EucARQHekz1X6bE=
-go.uber.org/multierr v1.6.0/go.mod h1:cdWPpRnG4AhwMwsgIHip0KRBQjJy5kYEpYjJxpXp9iU=
 go.uber.org/multierr v1.9.0 h1:7fIwc/ZtS0q++VgcfqFDxSBZVv/Xo49/SYnDFupUwlI=
 go.uber.org/multierr v1.9.0/go.mod h1:X2jQV1h+kxSjClGpnseKVIxpmcjrj7MNnI0bnlfKTVQ=
-go.uber.org/zap v1.18.1/go.mod h1:xg/QME4nWcxGxrpdeYfq7UvYrLh66cuVKdrbD1XF/NI=
-golang.org/x/arch v0.11.0 h1:KXV8WWKCXm6tRpLirl2szsO5j/oOODwZf4hATmGVNs4=
-golang.org/x/arch v0.11.0/go.mod h1:FEVrYAQjsQXMVJ1nsMoVVXPZg6p2JE2mx8psSWTDQys=
 golang.org/x/crypto v0.0.0-20180904163835-0709b304e793/go.mod h1:6SG95UA2DQfeDnfUPMdvaQW0Q7yPrPDi9nlGo2tz2b4=
 golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACkg1iLfiJU5Ep61QUkGW8qpdssI0+w=
 golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392/go.mod h1:/lpIB1dKB+9EgE3H3cr1v9wB50oz8l4C4h62xy7jSTY=
-golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
-golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9/go.mod h1:LzIPMQfyMNhhGPhUkYOs5KpL4U8rLKemX1yGLhDgUto=
-golang.org/x/crypto v0.28.0 h1:GBDwsMXVQi34v5CCYUm2jkJvu4cbtru2U4TN2PSyQnw=
-golang.org/x/crypto v0.28.0/go.mod h1:rmgy+3RHxRZMyY0jjAJShp2zgEdOqj2AO7U0pYmeQ7U=
-golang.org/x/exp v0.0.0-20190121172915-509febef88a4/go.mod h1:CJ0aWSM057203Lf6IL+f9T1iT9GByDxfZKAQTCR3kQA=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9 h1:GoHiUyI/Tp2nVkLI2mCxVkOjsbSXD66ic0XW0js0R9g=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9/go.mod h1:S2oDrQGGwySpoQPVqRShND87VCbxmc6bL1Yd2oYrm6k=
-golang.org/x/lint v0.0.0-20181026193005-c67002cb31c3/go.mod h1:UVdnD1Gm6xHRNCYTkRU2/jEulfH38KcIWyp/GAMgvoE=
-golang.org/x/lint v0.0.0-20190227174305-5b3e6a55c961/go.mod h1:wehouNa3lNwaWXcvxsM5YxQ5yQlVC4a0KAMCusXpPoU=
-golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3/go.mod h1:6SW0HCj/g11FgYtHlgUYUwCkIfeOF89ocIRzGO/8vkc=
-golang.org/x/lint v0.0.0-20190930215403-16217165b5de/go.mod h1:6SW0HCj/g11FgYtHlgUYUwCkIfeOF89ocIRzGO/8vkc=
-golang.org/x/mod v0.2.0/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
-golang.org/x/mod v0.3.0/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
-golang.org/x/net v0.0.0-20180724234803-3673e40ba225/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
-golang.org/x/net v0.0.0-20180826012351-8a410e7b638d/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20181114220301-adae6a3d119a/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
-golang.org/x/net v0.0.0-20190213061140-3a22650c66bd/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
-golang.org/x/net v0.0.0-20190311183353-d8887717615a/go.mod h1:t9HGtf8HONx5eT2rtn7q6eTqICYqUVnKs3thJo3Qplg=
 golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3/go.mod h1:t9HGtf8HONx5eT2rtn7q6eTqICYqUVnKs3thJo3Qplg=
 golang.org/x/net v0.0.0-20190613194153-d28f0bde5980/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
 golang.org/x/net v0.0.0-20190620200207-3b0461eec859/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
 golang.org/x/net v0.0.0-20190923162816-aa69164e4478/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
-golang.org/x/net v0.0.0-20200226121028-0de0cce0169b/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
-golang.org/x/net v0.0.0-20201021035429-f5854403a974/go.mod h1:sp8m0HH+o8qH0wwXwYZr8TS3Oi6o0r6Gce1SSxlDquU=
 golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
 golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1/go.mod h1:9tjilg8BloeKEkVJvy7fQ90B1CfIiPueXVOjqfkSzI8=
 golang.org/x/net v0.30.0 h1:AcW1SDZMkb8IpzCdQUaIq2sP4sZ4zw+55h6ynffypl4=
 golang.org/x/net v0.30.0/go.mod h1:2wGyMJ5iFasEhkwi13ChkO/t1ECNC4X4eBKkVFyYFlU=
-golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be/go.mod h1:N/0e6XlmueqKjAGxoOufVs8QHGRruUQn6yWY3a++T0U=
-golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20181108010431-42b317875d0f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20181221193216-37e7f081c4d4/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20190423024810-112230192c58/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
-golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20210220032951-036812b2e83c/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sys v0.0.0-20180823144017-11551d06cbcc/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
-golang.org/x/sys v0.0.0-20180830151530-49385e6e1522/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20180905080454-ebe1bf3edb33/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20181116152217-5ac8a444bdc5/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20190222072716-a9d3bda3a223/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
-golang.org/x/sys v0.0.0-20190412213103-97732733099d/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20190422165155-953cdadca894/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20190922100055-0a153f010e69/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20190924154521-2837fb4f24fe/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
@@ -372,17 +258,15 @@ golang.org/x/sys v0.0.0-20191026070338-33540a1f6037/go.mod h1:h1NjWce9XRLGQEsW7w
 golang.org/x/sys v0.0.0-20200116001909-b77594299b42/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20200122134326-e047566fdf82/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
-golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20201119102817-f84b799fce68/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210303074136-134d130e1a04/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210330210617-4fbd30eecc44/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
-golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
-golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.26.0 h1:KHjCJyddX0LoSTb3J+vWpupP9p0oznkqVk/IfjymZbo=
 golang.org/x/sys v0.26.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
@@ -394,40 +278,17 @@ golang.org/x/text v0.3.6/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
 golang.org/x/text v0.19.0 h1:kTxAhCbGbxhK0IwgSKiMO5awPoDQ0RpfiVYBfK860YM=
 golang.org/x/text v0.19.0/go.mod h1:BuEKDfySbSR4drPmRPG/7iBdf8hvFMuRexcpahXilzY=
 golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
-golang.org/x/tools v0.0.0-20190114222345-bf090417da8b/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
-golang.org/x/tools v0.0.0-20190226205152-f727befe758c/go.mod h1:9Yl7xja0Znq3iFh3HoIrodX9oNMXvdceNzlUR8zjMvY=
-golang.org/x/tools v0.0.0-20190311212946-11955173bddd/go.mod h1:LCzVGOaR6xXOjkQ3onu1FJEFr0SW1gC7cKk1uF8kGRs=
-golang.org/x/tools v0.0.0-20190524140312-2c0ae7006135/go.mod h1:RgjU9mgBXZiqYHBnxXauZ1Gv1EHHAz9KjViQ78xBX0Q=
 golang.org/x/tools v0.0.0-20190907020128-2ca718005c18/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
-golang.org/x/tools v0.0.0-20191108193012-7d206e10da11/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
-golang.org/x/tools v0.0.0-20191119224855-298f0cb1881e/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
-golang.org/x/tools v0.0.0-20200619180055-7c47624df98f/go.mod h1:EkVYQZoAsY45+roYkvgYkIh4xh/qjgUK9TdY2XT94GE=
-golang.org/x/tools v0.0.0-20210106214847-113979e3529a/go.mod h1:emZCQorbCU4vsT4fOWvOPXz4eW1wZW4PmDk9uLelYpA=
 golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
-golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
 golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
-golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
-google.golang.org/appengine v1.1.0/go.mod h1:EbEs0AVv82hx2wNQdGPgUI5lhzA/G0D9YwlJXL52JkM=
-google.golang.org/appengine v1.4.0/go.mod h1:xpcJRLb0r/rnEns0DIKYYv+WjYCduHsrkT7/EB5XEv4=
-google.golang.org/genproto v0.0.0-20180817151627-c66870c02cf8/go.mod h1:JiN7NxoALGmiZfu7CAH4rXhgtRTLTxftemlI0sWmxmc=
-google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55/go.mod h1:DMBHOl98Agz4BDEuKkezgsaosCRResVns1a3J2ZsMNc=
-google.golang.org/genproto v0.0.0-20200423170343-7949de9c1215/go.mod h1:55QSHmfGQM9UVYDPBsyGGes0y52j32PQ3BqQfXhyH3c=
 google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 h1:e7S5W7MGGLaSu8j3YjdezkZ+m1/Nm0uRVRMEMGk26Xs=
 google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142/go.mod h1:UqMtugtsSgubUsoxbuAoiCXvqvErP7Gf0so0mK9tHxU=
-google.golang.org/grpc v1.19.0/go.mod h1:mqu4LbDTu4XGKhr4mRzUsmM4RtVoemTSY81AxZiDr8c=
-google.golang.org/grpc v1.23.0/go.mod h1:Y5yQAOtifL1yxbo5wqy6BxZv8vAUGQwXBOALyacEbxg=
-google.golang.org/grpc v1.25.1/go.mod h1:c3i+UQWmh7LiEpx4sFZnkU36qjEYZ0imhYfXVyQciAY=
-google.golang.org/grpc v1.27.0/go.mod h1:qbnxyOmOxrQa7FizSgH+ReBfzJrCY1pSN7KXBS8abTk=
-google.golang.org/grpc v1.29.1/go.mod h1:itym6AZVZYACWQqET3MqgPpjcuV5QH3BxFS3IjizoKk=
 google.golang.org/grpc v1.67.1 h1:zWnc1Vrcno+lHZCOofnIMvycFcc0QRGIzm9dhnDX68E=
 google.golang.org/grpc v1.67.1/go.mod h1:1gLDyUQU7CTLJI90u3nXZ9ekeghjeM7pTDZlqFNg2AA=
-google.golang.org/protobuf v1.26.0-rc.1/go.mod h1:jlhhOSvTdKEhbULTjvd4ARK9grFBp09yW+WbY/TyQbw=
-google.golang.org/protobuf v1.26.0/go.mod h1:9q0QmTI4eRPtz6boOQmLYwt+qCgq0jsYwAQnmE0givc=
 google.golang.org/protobuf v1.35.1 h1:m3LfL6/Ca+fqnjnlqQXNpFPABW1UD7mjh8KO2mKFytA=
 google.golang.org/protobuf v1.35.1/go.mod h1:9fA7Ob0pmnwhb644+1+CVWFRbNajQ6iRojtC/QF5bRE=
 gopkg.in/alecthomas/kingpin.v2 v2.2.6/go.mod h1:FMv+mEhP44yOT+4EoQTLFTRgOQ1FBLkstjWtayDeSgw=
 gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
-gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 h1:YR8cESwS4TdDjEe65xsg0ogRM/Nc3DYOhEAlW+xobZo=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/ini.v1 v1.67.0 h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=
@@ -436,11 +297,6 @@ gopkg.in/yaml.v2 v2.2.1/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.2/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.4/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.5/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
-gopkg.in/yaml.v2 v2.2.8/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
-gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
 gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
 gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
-honnef.co/go/tools v0.0.0-20190102054323-c2f93a96b099/go.mod h1:rf3lG4BRIbNafJWhAfAdb/ePZxsR/4RtNHQocxwk9r4=
-honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc/go.mod h1:rf3lG4BRIbNafJWhAfAdb/ePZxsR/4RtNHQocxwk9r4=
-nullprogram.com/x/optparse v1.0.0/go.mod h1:KdyPE+Igbe0jQUrVfMqDMeJQIJZEuyV7pjYmp6pbG50=
diff --git a/internal/kitchen/infrastructure/consumer/consumer.go b/internal/kitchen/infrastructure/consumer/consumer.go
new file mode 100644
index 0000000..f17b444
--- /dev/null
+++ b/internal/kitchen/infrastructure/consumer/consumer.go
@@ -0,0 +1,107 @@
+package consumer
+
+import (
+	"context"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"time"
+
+	"github.com/ghost-yu/go_shop_second/common/broker"
+	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	amqp "github.com/rabbitmq/amqp091-go"
+	"github.com/sirupsen/logrus"
+	"go.opentelemetry.io/otel"
+)
+
+type OrderService interface {
+	UpdateOrder(ctx context.Context, request *orderpb.Order) error
+}
+
+type Consumer struct {
+	orderGRPC OrderService
+}
+
+type Order struct {
+	ID          string
+	CustomerID  string
+	Status      string
+	PaymentLink string
+	Items       []*orderpb.Item
+}
+
+func NewConsumer(orderGRPC OrderService) *Consumer {
+	return &Consumer{
+		orderGRPC: orderGRPC,
+	}
+}
+
+func (c *Consumer) Listen(ch *amqp.Channel) {
+	q, err := ch.QueueDeclare("", true, false, true, false, nil)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	if err = ch.QueueBind(q.Name, "", broker.EventOrderPaid, false, nil); err != nil {
+		logrus.Fatal(err)
+	}
+	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
+	if err != nil {
+		logrus.Warnf("fail to consume: queue=%s, err=%v", q.Name, err)
+	}
+
+	var forever chan struct{}
+	go func() {
+		for msg := range msgs {
+			c.handleMessage(ch, msg, q)
+		}
+	}()
+	<-forever
+}
+
+func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
+	var err error
+	logrus.Infof("kitchen receive a message from %s, msg=%v", q.Name, string(msg.Body))
+	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
+	tr := otel.Tracer("rabbitmq")
+	mqCtx, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
+	defer func() {
+		span.End()
+		if err != nil {
+			_ = msg.Nack(false, false)
+		} else {
+			_ = msg.Ack(false)
+		}
+	}()
+
+	o := &Order{}
+	if err := json.Unmarshal(msg.Body, o); err != nil {
+		logrus.Infof("failed to unmarshall msg to order, err=%v", err)
+		return
+	}
+	if o.Status != "paid" {
+		err = errors.New("order not paid, cannot cook")
+		return
+	}
+	cook(o)
+	span.AddEvent(fmt.Sprintf("order_cook: %v", o))
+	if err := c.orderGRPC.UpdateOrder(mqCtx, &orderpb.Order{
+		ID:          o.ID,
+		CustomerID:  o.CustomerID,
+		Status:      "ready",
+		Items:       o.Items,
+		PaymentLink: o.PaymentLink,
+	}); err != nil {
+		if err = broker.HandleRetry(mqCtx, ch, &msg); err != nil {
+			logrus.Warnf("kitchen: error handling retry: err=%v", err)
+		}
+		return
+	}
+	span.AddEvent("kitchen.order.finished.updated")
+	logrus.Info("consume success")
+}
+
+func cook(o *Order) {
+	logrus.Printf("cooking order: %s", o.ID)
+	time.Sleep(5 * time.Second)
+	logrus.Printf("order %s done!", o.ID)
+}
diff --git a/internal/kitchen/main.go b/internal/kitchen/main.go
new file mode 100644
index 0000000..35b4159
--- /dev/null
+++ b/internal/kitchen/main.go
@@ -0,0 +1,66 @@
+package main
+
+import (
+	"context"
+	"os"
+	"os/signal"
+	"syscall"
+
+	"github.com/ghost-yu/go_shop_second/common/broker"
+	grpcClient "github.com/ghost-yu/go_shop_second/common/client"
+	_ "github.com/ghost-yu/go_shop_second/common/config"
+	"github.com/ghost-yu/go_shop_second/common/logging"
+	"github.com/ghost-yu/go_shop_second/common/tracing"
+	"github.com/ghost-yu/go_shop_second/kitchen/adapters"
+	"github.com/ghost-yu/go_shop_second/kitchen/infrastructure/consumer"
+	"github.com/sirupsen/logrus"
+	"github.com/spf13/viper"
+)
+
+func init() {
+	logging.Init()
+}
+
+func main() {
+	serviceName := viper.GetString("kitchen.service-name")
+
+	ctx, cancel := context.WithCancel(context.Background())
+	defer cancel()
+
+	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	defer shutdown(ctx)
+
+	orderClient, closeFunc, err := grpcClient.NewOrderGRPCClient(ctx)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	defer closeFunc()
+
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
+	orderGRPC := adapters.NewOrderGRPC(orderClient)
+	go consumer.NewConsumer(orderGRPC).Listen(ch)
+
+	sigs := make(chan os.Signal, 1)
+	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
+
+	go func() {
+		<-sigs
+		logrus.Infof("receive signal, exiting...")
+		os.Exit(0)
+	}()
+	logrus.Println("to exit, press ctrl+c")
+	select {}
+}
diff --git a/internal/order/adapters/grpc/stock_grpc.go b/internal/order/adapters/grpc/stock_grpc.go
index c2397d6..549bdd1 100644
--- a/internal/order/adapters/grpc/stock_grpc.go
+++ b/internal/order/adapters/grpc/stock_grpc.go
@@ -2,6 +2,7 @@ package grpc
 
 import (
 	"context"
+	"errors"
 
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
@@ -17,6 +18,9 @@ func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
 }
 
 func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error) {
+	if items == nil {
+		return nil, errors.New("grpc items cannot be nil")
+	}
 	resp, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
 	logrus.Info("stock_grpc response", resp)
 	return resp, err
diff --git a/internal/order/tests/create_order_test.go b/internal/order/tests/create_order_test.go
index 2af7d2b..d57d2d2 100644
--- a/internal/order/tests/create_order_test.go
+++ b/internal/order/tests/create_order_test.go
@@ -31,9 +31,13 @@ func TestCreateOrder_success(t *testing.T) {
 		CustomerId: "123",
 		Items: []sw.ItemWithQuantity{
 			{
-				Id:       "test-item-1",
+				Id:       "prod_R3g7MikGYsXKzr",
 				Quantity: 1,
 			},
+			{
+				Id:       "prod_R285C3Wb7FDprc",
+				Quantity: 10,
+			},
 		},
 	})
 	t.Logf("body=%s", string(response.Body))
@@ -57,6 +61,7 @@ func getResponse(t *testing.T, customerID string, body sw.PostCustomerCustomerId
 	if err != nil {
 		t.Fatal(err)
 	}
+	t.Logf("getResponse body=%+v", body)
 	response, err := client.PostCustomerCustomerIdOrdersWithResponse(ctx, customerID, body)
 	if err != nil {
 		t.Fatal(err)
diff --git a/internal/stock/app/query/check_if_items_in_stock.go b/internal/stock/app/query/check_if_items_in_stock.go
index 3d3114e..00b5c4c 100644
--- a/internal/stock/app/query/check_if_items_in_stock.go
+++ b/internal/stock/app/query/check_if_items_in_stock.go
@@ -4,51 +4,59 @@ import (
 	"context"
 
 	"github.com/ghost-yu/go_shop_second/common/decorator"
-	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
+	"github.com/ghost-yu/go_shop_second/stock/entity"
+	"github.com/ghost-yu/go_shop_second/stock/infrastructure/integration"
 	"github.com/sirupsen/logrus"
 )
 
 type CheckIfItemsInStock struct {
-	Items []*orderpb.ItemWithQuantity
+	Items []*entity.ItemWithQuantity
 }
 
-type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*orderpb.Item]
+type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*entity.Item]
 
 type checkIfItemsInStockHandler struct {
 	stockRepo domain.Repository
+	stripeAPI *integration.StripeAPI
 }
 
 func NewCheckIfItemsInStockHandler(
 	stockRepo domain.Repository,
+	stripeAPI *integration.StripeAPI,
 	logger *logrus.Entry,
 	metricClient decorator.MetricsClient,
 ) CheckIfItemsInStockHandler {
 	if stockRepo == nil {
 		panic("nil stockRepo")
 	}
-	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*orderpb.Item](
-		checkIfItemsInStockHandler{stockRepo: stockRepo},
+	if stripeAPI == nil {
+		panic("nil stripeAPI")
+	}
+	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*entity.Item](
+		checkIfItemsInStockHandler{
+			stockRepo: stockRepo,
+			stripeAPI: stripeAPI,
+		},
 		logger,
 		metricClient,
 	)
 }
 
-// TODO: 删掉
+// Deprecated
 var stub = map[string]string{
 	"1": "price_1QBYvXRuyMJmUCSsEyQm2oP7",
 	"2": "price_1QBYl4RuyMJmUCSsWt2tgh6d",
 }
 
-func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*orderpb.Item, error) {
-	var res []*orderpb.Item
+func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*entity.Item, error) {
+	var res []*entity.Item
 	for _, i := range query.Items {
-		// TODO: 改成从数据库 or stripe 获取
-		priceID, ok := stub[i.ID]
-		if !ok {
-			priceID = stub["1"]
+		priceID, err := h.stripeAPI.GetPriceByProductID(ctx, i.ID)
+		if err != nil || priceID == "" {
+			return nil, err
 		}
-		res = append(res, &orderpb.Item{
+		res = append(res, &entity.Item{
 			ID:       i.ID,
 			Quantity: i.Quantity,
 			PriceID:  priceID,
@@ -56,3 +64,11 @@ func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfIte
 	}
 	return res, nil
 }
+
+func getStubPriceID(id string) string {
+	priceID, ok := stub[id]
+	if !ok {
+		priceID = stub["1"]
+	}
+	return priceID
+}
diff --git a/internal/stock/convertor/convertor.go b/internal/stock/convertor/convertor.go
new file mode 100644
index 0000000..be35f65
--- /dev/null
+++ b/internal/stock/convertor/convertor.go
@@ -0,0 +1,97 @@
+package convertor
+
+import (
+	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/stock/entity"
+)
+
+type OrderConvertor struct{}
+type ItemConvertor struct{}
+type ItemWithQuantityConvertor struct{}
+
+func (c *ItemWithQuantityConvertor) EntitiesToProtos(items []*entity.ItemWithQuantity) (res []*orderpb.ItemWithQuantity) {
+	for _, i := range items {
+		res = append(res, c.EntityToProto(i))
+	}
+	return
+}
+
+func (c *ItemWithQuantityConvertor) EntityToProto(i *entity.ItemWithQuantity) *orderpb.ItemWithQuantity {
+	return &orderpb.ItemWithQuantity{
+		ID:       i.ID,
+		Quantity: i.Quantity,
+	}
+}
+
+func (c *ItemWithQuantityConvertor) ProtosToEntities(items []*orderpb.ItemWithQuantity) (res []*entity.ItemWithQuantity) {
+	for _, i := range items {
+		res = append(res, c.ProtoToEntity(i))
+	}
+	return
+}
+
+func (c *ItemWithQuantityConvertor) ProtoToEntity(i *orderpb.ItemWithQuantity) *entity.ItemWithQuantity {
+	return &entity.ItemWithQuantity{
+		ID:       i.ID,
+		Quantity: i.Quantity,
+	}
+}
+
+func (c *OrderConvertor) EntityToProto(o *entity.Order) *orderpb.Order {
+	c.check(o)
+	return &orderpb.Order{
+		ID:          o.ID,
+		CustomerID:  o.CustomerID,
+		Status:      o.Status,
+		Items:       NewItemConvertor().EntitiesToProtos(o.Items),
+		PaymentLink: o.PaymentLink,
+	}
+}
+
+func (c *OrderConvertor) ProtoToEntity(o *orderpb.Order) *entity.Order {
+	c.check(o)
+	return &entity.Order{
+		ID:          o.ID,
+		CustomerID:  o.CustomerID,
+		Status:      o.Status,
+		PaymentLink: o.PaymentLink,
+		Items:       NewItemConvertor().ProtosToEntities(o.Items),
+	}
+}
+func (c *OrderConvertor) check(o interface{}) {
+	if o == nil {
+		panic("connot convert nil order")
+	}
+}
+
+func (c *ItemConvertor) EntitiesToProtos(items []*entity.Item) (res []*orderpb.Item) {
+	for _, i := range items {
+		res = append(res, c.EntityToProto(i))
+	}
+	return
+}
+
+func (c *ItemConvertor) ProtosToEntities(items []*orderpb.Item) (res []*entity.Item) {
+	for _, i := range items {
+		res = append(res, c.ProtoToEntity(i))
+	}
+	return
+}
+
+func (c *ItemConvertor) EntityToProto(i *entity.Item) *orderpb.Item {
+	return &orderpb.Item{
+		ID:       i.ID,
+		Name:     i.Name,
+		Quantity: i.Quantity,
+		PriceID:  i.PriceID,
+	}
+}
+
+func (c *ItemConvertor) ProtoToEntity(i *orderpb.Item) *entity.Item {
+	return &entity.Item{
+		ID:       i.ID,
+		Name:     i.Name,
+		Quantity: i.Quantity,
+		PriceID:  i.PriceID,
+	}
+}
diff --git a/internal/stock/convertor/facade.go b/internal/stock/convertor/facade.go
new file mode 100644
index 0000000..63ad8d0
--- /dev/null
+++ b/internal/stock/convertor/facade.go
@@ -0,0 +1,39 @@
+package convertor
+
+import "sync"
+
+var (
+	orderConvertor *OrderConvertor
+	orderOnce      sync.Once
+)
+
+var (
+	itemConvertor *ItemConvertor
+	itemOnce      sync.Once
+)
+
+var (
+	itemWithQuantityConvertor *ItemWithQuantityConvertor
+	itemWithQuantityOnce      sync.Once
+)
+
+func NewOrderConvertor() *OrderConvertor {
+	orderOnce.Do(func() {
+		orderConvertor = new(OrderConvertor)
+	})
+	return orderConvertor
+}
+
+func NewItemConvertor() *ItemConvertor {
+	itemOnce.Do(func() {
+		itemConvertor = new(ItemConvertor)
+	})
+	return itemConvertor
+}
+
+func NewItemWithQuantityConvertor() *ItemWithQuantityConvertor {
+	itemWithQuantityOnce.Do(func() {
+		itemWithQuantityConvertor = new(ItemWithQuantityConvertor)
+	})
+	return itemWithQuantityConvertor
+}
diff --git a/internal/stock/entity/entity.go b/internal/stock/entity/entity.go
new file mode 100644
index 0000000..671b85b
--- /dev/null
+++ b/internal/stock/entity/entity.go
@@ -0,0 +1,21 @@
+package entity
+
+type Order struct {
+	ID          string
+	CustomerID  string
+	Status      string
+	PaymentLink string
+	Items       []*Item
+}
+
+type Item struct {
+	ID       string
+	Name     string
+	Quantity int32
+	PriceID  string
+}
+
+type ItemWithQuantity struct {
+	ID       string
+	Quantity int32
+}
diff --git a/internal/stock/go.mod b/internal/stock/go.mod
index 4083f63..41d06cf 100644
--- a/internal/stock/go.mod
+++ b/internal/stock/go.mod
@@ -8,10 +8,12 @@ require (
 	github.com/ghost-yu/go_shop_second/common v0.0.0-00010101000000-000000000000
 	github.com/sirupsen/logrus v1.8.1
 	github.com/spf13/viper v1.19.0
+	github.com/stripe/stripe-go/v79 v79.12.0
 	google.golang.org/grpc v1.67.1
 )
 
 require (
+	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
 	github.com/armon/go-metrics v0.4.1 // indirect
 	github.com/bytedance/sonic v1.12.3 // indirect
 	github.com/bytedance/sonic/loader v0.2.0 // indirect
@@ -47,12 +49,11 @@ require (
 	github.com/magiconair/properties v1.8.7 // indirect
 	github.com/mattn/go-colorable v0.1.13 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
-	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
 	github.com/mitchellh/go-homedir v1.1.0 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
 	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
 	github.com/modern-go/reflect2 v1.0.2 // indirect
-	github.com/nxadm/tail v1.4.11 // indirect
+	github.com/oapi-codegen/runtime v1.1.1 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
@@ -63,7 +64,6 @@ require (
 	github.com/subosito/gotenv v1.6.0 // indirect
 	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
 	github.com/ugorji/go/codec v1.2.12 // indirect
-	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
 	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 // indirect
 	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
 	go.opentelemetry.io/contrib/propagators/b3 v1.31.0 // indirect
@@ -79,7 +79,6 @@ require (
 	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
 	golang.org/x/net v0.30.0 // indirect
 	golang.org/x/sys v0.26.0 // indirect
-	golang.org/x/term v0.25.0 // indirect
 	golang.org/x/text v0.19.0 // indirect
 	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
 	google.golang.org/protobuf v1.35.1 // indirect
diff --git a/internal/stock/go.sum b/internal/stock/go.sum
index f649128..17959ae 100644
--- a/internal/stock/go.sum
+++ b/internal/stock/go.sum
@@ -1,10 +1,13 @@
 cloud.google.com/go v0.26.0/go.mod h1:aQUYkXzVsufM+DwF1aE+0xfcU+56JwCaLick0ClmMTw=
 github.com/BurntSushi/toml v0.3.1/go.mod h1:xHWCNGjB5oqiDr8zfno3MHue2Ht5sIBksp03qcyfWMU=
 github.com/DataDog/datadog-go v3.2.0+incompatible/go.mod h1:LButxg5PwREeZtORoXG3tL4fMGNddJ+vMq1mwgfaqoQ=
+github.com/RaveNoX/go-jsoncommentstrip v1.0.0/go.mod h1:78ihd09MekBnJnxpICcwzCMzGrKSKYe4AqU6PDYYpjk=
 github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
 github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
 github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
 github.com/alecthomas/units v0.0.0-20190717042225-c3de453c63f4/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
+github.com/apapsch/go-jsonmerge/v2 v2.0.0 h1:axGnT1gRIfimI7gJifB699GoE/oq+F2MU7Dml6nw9rQ=
+github.com/apapsch/go-jsonmerge/v2 v2.0.0/go.mod h1:lvDnEdqiQrp0O42VQGgmlKpxL1AP2+08jFMw88y4klk=
 github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e/go.mod h1:3U/XgcO3hCbHZ8TKRvWD2dDTCfh9M9ya+I9JpbB7O8o=
 github.com/armon/go-metrics v0.0.0-20180917152333-f0300d1749da/go.mod h1:Q73ZrmVTwzkszR9V5SSuryQ31EELlFMUz1kKyl939pY=
 github.com/armon/go-metrics v0.4.1 h1:hR91U9KYmb6bLBYLQjyM+3j+rcd/UhE+G78SFnF8gJA=
@@ -16,6 +19,7 @@ github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973/go.mod h1:Dwedo/Wpr24
 github.com/beorn7/perks v1.0.0/go.mod h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+CedLV8=
 github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
 github.com/bgentry/speakeasy v0.1.0/go.mod h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=
+github.com/bmatcuk/doublestar v1.1.1/go.mod h1:UD6OnuiIn0yFxxA2le/rnRU1G4RaI4UvFv1sNto9p6w=
 github.com/bytedance/sonic v1.12.3 h1:W2MGa7RCU1QTeYRTPE3+88mVC0yXmsRQRChiyVocVjU=
 github.com/bytedance/sonic v1.12.3/go.mod h1:B8Gt/XvtZ3Fqj+iSKMypzymZxw/FVwgIGKzMzT9r/rk=
 github.com/bytedance/sonic/loader v0.1.1/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
@@ -46,7 +50,6 @@ github.com/fatih/color v1.14.1 h1:qfhVLaG5s+nCROl1zJsZRxFeYrHLqWroPOQ8BWiNb4w=
 github.com/fatih/color v1.14.1/go.mod h1:2oHN61fhTpgcxD3TSWCgKDiH1+x4OiDVVGH8WlgGZGg=
 github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHkI4W8=
 github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
-github.com/fsnotify/fsnotify v1.6.0/go.mod h1:sl3t1tCWJFWoRz9R8WJCbQihKKwmorjAbSClcnxKAGw=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
 github.com/gabriel-vasile/mimetype v1.4.5 h1:J7wGKdGu33ocBOhGy0z653k/lFKLFDPJMG8Gql0kxn4=
@@ -153,6 +156,7 @@ github.com/json-iterator/go v1.1.6/go.mod h1:+SdeFBvtyEkXs7REEP0seUULqWtbJapLOCV
 github.com/json-iterator/go v1.1.9/go.mod h1:KdQUCv79m/52Kvf8AW2vK1V8akMuk1QjK/uOdHXbAo4=
 github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
 github.com/json-iterator/go v1.1.12/go.mod h1:e30LSqwooZae/UwlEbR2852Gd8hjQvJoHmT4TnhNGBo=
+github.com/juju/gnuflag v0.0.0-20171113085948-2ce1bb71843d/go.mod h1:2PavIy+JPciBPrBUjwbNvtwB6RQlve+hkpll6QSNmOE=
 github.com/julienschmidt/httprouter v1.2.0/go.mod h1:SYymIcj16QtmaHHD7aYtjjsJG7VTCxuUUipMqKk8s4w=
 github.com/kisielk/errcheck v1.5.0/go.mod h1:pFxgyoBC7bSaBwPgfKdkLd5X25qrDl4LWUI2bnpBCr8=
 github.com/kisielk/gotool v1.0.0/go.mod h1:XhKaO+MFFWcvkIS/tQcRk01m1F5IRFswLeQ+oQHNcck=
@@ -189,8 +193,6 @@ github.com/mattn/go-isatty v0.0.16/go.mod h1:kYGgaQfpe5nmfYZH+SKPsOc2e4SrIfOl2e/
 github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
 github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
 github.com/matttproud/golang_protobuf_extensions v1.0.1/go.mod h1:D8He9yQNgCq6Z5Ld7szi9bcBfOoFv/3dc6xSMkL2PC0=
-github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d h1:5PJl274Y63IEHC+7izoQE9x6ikvDFZS2mDVS3drnohI=
-github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d/go.mod h1:01TrycV0kFyexm33Z7vhZRXopbI8J3TDReVlkTgMUxE=
 github.com/miekg/dns v1.1.26/go.mod h1:bPDLeHnStXmXAq1m/Ch/hvfNHr14JKNPMBo3VZKjuso=
 github.com/miekg/dns v1.1.41 h1:WMszZWJG0XmzbK9FEmzH2TVcqYzFesusSIB41b8KHxY=
 github.com/miekg/dns v1.1.41/go.mod h1:p6aan82bvRIyn+zDIv9xYNUpwa73JcSh9BKwknJysuI=
@@ -208,12 +210,8 @@ github.com/modern-go/reflect2 v1.0.1/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3Rllmb
 github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
 github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
 github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
-github.com/nxadm/tail v1.4.11 h1:8feyoE3OzPrcshW5/MJ4sGESc5cqmGkGCWlco4l0bqY=
-github.com/nxadm/tail v1.4.11/go.mod h1:OTaG3NK980DZzxbRq6lEuzgU+mug70nY11sMd4JXXHc=
-github.com/onsi/ginkgo v1.16.5 h1:8xi0RTUf59SOSfEtZMvwTvXYMzG4gV23XVHOZiXNtnE=
-github.com/onsi/ginkgo v1.16.5/go.mod h1:+E8gABHa3K6zRBolWtd+ROzc/U5bkGt0FwiG042wbpU=
-github.com/onsi/gomega v1.34.2 h1:pNCwDkzrsv7MS9kpaQvVb1aVLahQXyJ/Tv5oAZMI3i8=
-github.com/onsi/gomega v1.34.2/go.mod h1:v1xfxRgk0KIsG+QOdm7p8UosrOzPYRo60fd3B/1Dukc=
+github.com/oapi-codegen/runtime v1.1.1 h1:EXLHh0DXIJnWhdRPN2w4MXAzFyE4CskzhNLUmtpMYro=
+github.com/oapi-codegen/runtime v1.1.1/go.mod h1:SK9X900oXmPWilYR5/WKPzt3Kqxn/uS/+lbpREv+eCg=
 github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFStlNbqXla1AfSYxPUl2o=
 github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
@@ -264,6 +262,7 @@ github.com/spf13/pflag v1.0.5 h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=
 github.com/spf13/pflag v1.0.5/go.mod h1:McXfInJRrz4CZXVZOBLb0bTZqETkiAhM9Iw0y3An2Bg=
 github.com/spf13/viper v1.19.0 h1:RWq5SEjt8o25SROyN3z2OrDB9l7RPd3lwTWU8EcEdcI=
 github.com/spf13/viper v1.19.0/go.mod h1:GQUN9bilAbhU/jgc1bKs99f/suXKeUMct8Adx5+Ntkg=
+github.com/spkg/bom v0.0.0-20160624110644-59b7046e48ad/go.mod h1:qLr4V1qq6nMqFKkMo8ZTx3f+BZEkzsRUY10Xsm2mwU0=
 github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.1.1/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
@@ -279,6 +278,8 @@ github.com/stretchr/testify v1.8.0/go.mod h1:yNjHg4UonilssWZ8iaSj1OCr/vHnekPRkoO
 github.com/stretchr/testify v1.8.1/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o6fzry7u4=
 github.com/stretchr/testify v1.9.0 h1:HtqpIVDClZ4nwg75+f6Lvsy/wHu+3BoSGCbBAcpTsTg=
 github.com/stretchr/testify v1.9.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8C91i36aY=
+github.com/stripe/stripe-go/v79 v79.12.0 h1:HQs/kxNEB3gYA7FnkSFkp0kSOeez0fsmCWev6SxftYs=
+github.com/stripe/stripe-go/v79 v79.12.0/go.mod h1:cuH6X0zC8peY6f1AubHwgJ/fJSn2dh5pfiCr6CjyKVU=
 github.com/subosito/gotenv v1.6.0 h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=
 github.com/subosito/gotenv v1.6.0/go.mod h1:Dk4QP5c2W3ibzajGcXpNraDfq2IrhjMIvMSWPKKo0FU=
 github.com/tv42/httpunix v0.0.0-20150427012821-b75d8614f926/go.mod h1:9ESjWnEqriFuLhtthL60Sar/7RFoluCcXsuvEwTV5KM=
@@ -286,8 +287,6 @@ github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS
 github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
 github.com/ugorji/go/codec v1.2.12 h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65EE=
 github.com/ugorji/go/codec v1.2.12/go.mod h1:UNopzCgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=
-github.com/x-cray/logrus-prefixed-formatter v0.5.2 h1:00txxvfBM9muc0jiLIEAkAcIMJzfthRT6usrui8uGmg=
-github.com/x-cray/logrus-prefixed-formatter v0.5.2/go.mod h1:2duySbKsL6M18s5GU7VPsoEPHyzalCE06qoARUCeBBE=
 github.com/yuin/goldmark v1.1.27/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
 go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0 h1:0nTRpaCaILLdooXAQnfktlL6Zw1ECKEW9DZGH2byi2c=
@@ -345,6 +344,7 @@ golang.org/x/net v0.0.0-20200226121028-0de0cce0169b/go.mod h1:z5CRVTTTmAJ677TzLL
 golang.org/x/net v0.0.0-20201021035429-f5854403a974/go.mod h1:sp8m0HH+o8qH0wwXwYZr8TS3Oi6o0r6Gce1SSxlDquU=
 golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
 golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1/go.mod h1:9tjilg8BloeKEkVJvy7fQ90B1CfIiPueXVOjqfkSzI8=
+golang.org/x/net v0.0.0-20210520170846-37e1c6afe023/go.mod h1:9nx3DQGgdP8bBQD5qxJ1jj9UTztislL4KSBs9R2vV5Y=
 golang.org/x/net v0.30.0 h1:AcW1SDZMkb8IpzCdQUaIq2sP4sZ4zw+55h6ynffypl4=
 golang.org/x/net v0.30.0/go.mod h1:2wGyMJ5iFasEhkwi13ChkO/t1ECNC4X4eBKkVFyYFlU=
 golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be/go.mod h1:N/0e6XlmueqKjAGxoOufVs8QHGRruUQn6yWY3a++T0U=
@@ -373,20 +373,18 @@ golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f/go.mod h1:h1NjWce9XRLGQEsW7w
 golang.org/x/sys v0.0.0-20201119102817-f84b799fce68/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210303074136-134d130e1a04/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210330210617-4fbd30eecc44/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20210423082822-04245dca01da/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
-golang.org/x/sys v0.0.0-20220908164124-27713097b956/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.26.0 h1:KHjCJyddX0LoSTb3J+vWpupP9p0oznkqVk/IfjymZbo=
 golang.org/x/sys v0.26.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
 golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
-golang.org/x/term v0.25.0 h1:WtHI/ltw4NvSUig5KARz9h521QvRC8RmF/cuYqifU24=
-golang.org/x/term v0.25.0/go.mod h1:RPyXicDX+6vLxogjjRxjgD2TKtmAO6NZBsBRfrOLu7M=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
 golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
@@ -432,8 +430,6 @@ gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 h1:YR8cESwS4TdDjEe65xsg0ogR
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/ini.v1 v1.67.0 h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=
 gopkg.in/ini.v1 v1.67.0/go.mod h1:pNLf8WUiyNEtQjuu5G5vTm06TEv9tsIgeAvK8hOrP4k=
-gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=
-gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7/go.mod h1:dt/ZhP58zS4L8KSrWDmTeBkI65Dw0HsyUHuEVlX15mw=
 gopkg.in/yaml.v2 v2.2.1/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.2/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.4/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
diff --git a/internal/stock/infrastructure/integration/stripe.go b/internal/stock/infrastructure/integration/stripe.go
new file mode 100644
index 0000000..c24cee9
--- /dev/null
+++ b/internal/stock/infrastructure/integration/stripe.go
@@ -0,0 +1,32 @@
+package integration
+
+import (
+	"context"
+
+	_ "github.com/ghost-yu/go_shop_second/common/config"
+	"github.com/sirupsen/logrus"
+	"github.com/spf13/viper"
+	"github.com/stripe/stripe-go/v79"
+	"github.com/stripe/stripe-go/v79/product"
+)
+
+type StripeAPI struct {
+	apiKey string
+}
+
+func NewStripeAPI() *StripeAPI {
+	key := viper.GetString("stripe-key")
+	if key == "" {
+		logrus.Fatal("empty key")
+	}
+	return &StripeAPI{apiKey: key}
+}
+
+func (s *StripeAPI) GetPriceByProductID(ctx context.Context, pid string) (string, error) {
+	stripe.Key = s.apiKey
+	result, err := product.Get(pid, &stripe.ProductParams{})
+	if err != nil {
+		return "", err
+	}
+	return result.DefaultPrice.ID, err
+}
diff --git a/internal/stock/ports/grpc.go b/internal/stock/ports/grpc.go
index 53b5eaf..b2c659d 100644
--- a/internal/stock/ports/grpc.go
+++ b/internal/stock/ports/grpc.go
@@ -1,12 +1,13 @@
 package ports
 
 import (
-	context "context"
+	"context"
 
 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
 	"github.com/ghost-yu/go_shop_second/common/tracing"
 	"github.com/ghost-yu/go_shop_second/stock/app"
 	"github.com/ghost-yu/go_shop_second/stock/app/query"
+	"github.com/ghost-yu/go_shop_second/stock/convertor"
 )
 
 type GRPCServer struct {
@@ -32,12 +33,14 @@ func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.Ch
 	_, span := tracing.Start(ctx, "CheckIfItemsInStock")
 	defer span.End()
 
-	items, err := G.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{Items: request.Items})
+	items, err := G.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{
+		Items: convertor.NewItemWithQuantityConvertor().ProtosToEntities(request.Items),
+	})
 	if err != nil {
 		return nil, err
 	}
 	return &stockpb.CheckIfItemsInStockResponse{
 		InStock: 1,
-		Items:   items,
+		Items:   convertor.NewItemConvertor().EntitiesToProtos(items),
 	}, nil
 }
diff --git a/internal/stock/service/application.go b/internal/stock/service/application.go
index 0ccb6a8..b51f2ad 100644
--- a/internal/stock/service/application.go
+++ b/internal/stock/service/application.go
@@ -7,17 +7,19 @@ import (
 	"github.com/ghost-yu/go_shop_second/stock/adapters"
 	"github.com/ghost-yu/go_shop_second/stock/app"
 	"github.com/ghost-yu/go_shop_second/stock/app/query"
+	"github.com/ghost-yu/go_shop_second/stock/infrastructure/integration"
 	"github.com/sirupsen/logrus"
 )
 
 func NewApplication(_ context.Context) app.Application {
 	stockRepo := adapters.NewMemoryStockRepository()
 	logger := logrus.NewEntry(logrus.StandardLogger())
+	stripeAPI := integration.NewStripeAPI()
 	metricsClient := metrics.TodoMetrics{}
 	return app.Application{
 		Commands: app.Commands{},
 		Queries: app.Queries{
-			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, logger, metricsClient),
+			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, stripeAPI, logger, metricsClient),
 			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricsClient),
 		},
 	}
~~~

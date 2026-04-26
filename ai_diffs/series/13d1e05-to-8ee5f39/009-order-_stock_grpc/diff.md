# Commit Diff Report

- Repo: go_shop_second
- Sequence: 009 / 10
- Commit: 78d6465b2c83f580f21a1edf5ebdf7122458c8a5
- ShortCommit: 78d6465
- Parent: 606345b1bffa9e6e0f221a73bc2671359315a77a
- Subject: order->stock grpc
- Author: ghost-yu <hgfhgfhgfhgfhgfhgf@yeah.net>
- Date: 2024-10-14 23:13:55 +0800
- GeneratedAt: 2026-04-06 17:43:36 +08:00

## Short Summary

~~~text
 7 files changed, 134 insertions(+), 27 deletions(-)
~~~

## File Stats

~~~text
 internal/order/adapters/grpc/stock_grpc.go | 31 ++++++++++++++++++++
 internal/order/app/command/create_order.go |  9 ++++--
 internal/order/app/query/service.go        | 12 ++++++++
 internal/order/go.mod                      | 44 ++++++++++++++---------------
 internal/order/go.sum                      | 45 ++++++++++++++++++++++++++++++
 internal/order/main.go                     |  3 +-
 internal/order/service/application.go      | 17 +++++++++--
 7 files changed, 134 insertions(+), 27 deletions(-)
~~~

## Changed Files

~~~text
internal/order/adapters/grpc/stock_grpc.go
internal/order/app/command/create_order.go
internal/order/app/query/service.go
internal/order/go.mod
internal/order/go.sum
internal/order/main.go
internal/order/service/application.go
~~~

## Focus Files (Excluded: go.mod / go.sum)

~~~text
internal/order/adapters/grpc/stock_grpc.go
internal/order/app/command/create_order.go
internal/order/app/query/service.go
internal/order/main.go
internal/order/service/application.go
~~~

## Patch

~~~diff
diff --git a/internal/order/adapters/grpc/stock_grpc.go b/internal/order/adapters/grpc/stock_grpc.go
new file mode 100644
index 0000000..dfe61e6
--- /dev/null
+++ b/internal/order/adapters/grpc/stock_grpc.go
@@ -0,0 +1,31 @@
+package grpc
+
+import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
+	"github.com/sirupsen/logrus"
+)
+
+type StockGRPC struct {
+	client stockpb.StockServiceClient
+}
+
+func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
+	return &StockGRPC{client: client}
+}
+
+func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) error {
+	resp, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
+	logrus.Info("stock_grpc response", resp)
+	return err
+}
+
+func (s StockGRPC) GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error) {
+	resp, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{ItemIDs: itemIDs})
+	if err != nil {
+		return nil, err
+	}
+	return resp.Items, nil
+}
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
index 17620e2..adb658f 100644
--- a/internal/order/app/command/create_order.go
+++ b/internal/order/app/command/create_order.go
@@ -5,6 +5,7 @@ import (
 
 	"github.com/ghost-yu/go_shop_second/common/decorator"
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/order/app/query"
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
 	"github.com/sirupsen/logrus"
 )
@@ -22,11 +23,12 @@ type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult
 
 type createOrderHandler struct {
 	orderRepo domain.Repository
-	//stockGRPC
+	stockGRPC query.StockService
 }
 
 func NewCreateOrderHandler(
 	orderRepo domain.Repository,
+	stockGRPC query.StockService,
 	logger *logrus.Entry,
 	metricClient decorator.MetricsClient,
 ) CreateOrderHandler {
@@ -34,7 +36,7 @@ func NewCreateOrderHandler(
 		panic("nil orderRepo")
 	}
 	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
-		createOrderHandler{orderRepo: orderRepo},
+		createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC},
 		logger,
 		metricClient,
 	)
@@ -42,6 +44,9 @@ func NewCreateOrderHandler(
 
 func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
 	// TODO: call stock grpc to get items.
+	err := c.stockGRPC.CheckIfItemsInStock(ctx, cmd.Items)
+	resp, err := c.stockGRPC.GetItems(ctx, []string{"123"})
+	logrus.Info("createOrderHandler||resp from stockGRPC.GetItems", resp)
 	var stockResponse []*orderpb.Item
 	for _, item := range cmd.Items {
 		stockResponse = append(stockResponse, &orderpb.Item{
diff --git a/internal/order/app/query/service.go b/internal/order/app/query/service.go
new file mode 100644
index 0000000..3f419a9
--- /dev/null
+++ b/internal/order/app/query/service.go
@@ -0,0 +1,12 @@
+package query
+
+import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+)
+
+type StockService interface {
+	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) error
+	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
+}
diff --git a/internal/order/go.mod b/internal/order/go.mod
index a57a04d..3031f57 100644
--- a/internal/order/go.mod
+++ b/internal/order/go.mod
@@ -8,56 +8,56 @@ require (
 	github.com/ghost-yu/go_shop_second/common v0.0.0-00010101000000-000000000000
 	github.com/gin-gonic/gin v1.10.0
 	github.com/oapi-codegen/runtime v1.1.1
+	github.com/sirupsen/logrus v1.9.3
 	github.com/spf13/viper v1.19.0
-	google.golang.org/grpc v1.62.1
-	google.golang.org/protobuf v1.34.1
+	google.golang.org/grpc v1.67.1
+	google.golang.org/protobuf v1.35.1
 )
 
 require (
 	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
-	github.com/bytedance/sonic v1.11.6 // indirect
-	github.com/bytedance/sonic/loader v0.1.1 // indirect
+	github.com/bytedance/sonic v1.12.3 // indirect
+	github.com/bytedance/sonic/loader v0.2.0 // indirect
 	github.com/cloudwego/base64x v0.1.4 // indirect
 	github.com/cloudwego/iasm v0.2.0 // indirect
 	github.com/fsnotify/fsnotify v1.7.0 // indirect
-	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
+	github.com/gabriel-vasile/mimetype v1.4.6 // indirect
 	github.com/gin-contrib/sse v0.1.0 // indirect
 	github.com/go-playground/locales v0.14.1 // indirect
 	github.com/go-playground/universal-translator v0.18.1 // indirect
-	github.com/go-playground/validator/v10 v10.20.0 // indirect
-	github.com/goccy/go-json v0.10.2 // indirect
-	github.com/golang/protobuf v1.5.3 // indirect
+	github.com/go-playground/validator/v10 v10.22.1 // indirect
+	github.com/goccy/go-json v0.10.3 // indirect
+	github.com/golang/protobuf v1.5.4 // indirect
 	github.com/google/uuid v1.6.0 // indirect
 	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
 	github.com/hashicorp/hcl v1.0.0 // indirect
 	github.com/json-iterator/go v1.1.12 // indirect
-	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
+	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
 	github.com/leodido/go-urn v1.4.0 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
 	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
 	github.com/modern-go/reflect2 v1.0.2 // indirect
-	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
-	github.com/sagikazarmark/locafero v0.4.0 // indirect
+	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
+	github.com/sagikazarmark/locafero v0.6.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
-	github.com/sirupsen/logrus v1.8.1 // indirect
 	github.com/sourcegraph/conc v0.3.0 // indirect
 	github.com/spf13/afero v1.11.0 // indirect
-	github.com/spf13/cast v1.6.0 // indirect
+	github.com/spf13/cast v1.7.0 // indirect
 	github.com/spf13/pflag v1.0.5 // indirect
 	github.com/subosito/gotenv v1.6.0 // indirect
 	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
 	github.com/ugorji/go/codec v1.2.12 // indirect
-	go.uber.org/atomic v1.9.0 // indirect
-	go.uber.org/multierr v1.9.0 // indirect
-	golang.org/x/arch v0.8.0 // indirect
-	golang.org/x/crypto v0.23.0 // indirect
-	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
-	golang.org/x/net v0.25.0 // indirect
-	golang.org/x/sys v0.20.0 // indirect
-	golang.org/x/text v0.15.0 // indirect
-	google.golang.org/genproto/googleapis/rpc v0.0.0-20240314234333-6e1732d8331c // indirect
+	go.uber.org/atomic v1.11.0 // indirect
+	go.uber.org/multierr v1.11.0 // indirect
+	golang.org/x/arch v0.11.0 // indirect
+	golang.org/x/crypto v0.28.0 // indirect
+	golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c // indirect
+	golang.org/x/net v0.30.0 // indirect
+	golang.org/x/sys v0.26.0 // indirect
+	golang.org/x/text v0.19.0 // indirect
+	google.golang.org/genproto/googleapis/rpc v0.0.0-20241007155032-5fefd90f89a9 // indirect
 	gopkg.in/ini.v1 v1.67.0 // indirect
 	gopkg.in/yaml.v3 v3.0.1 // indirect
 )
diff --git a/internal/order/go.sum b/internal/order/go.sum
index d7c5781..7e9ddb3 100644
--- a/internal/order/go.sum
+++ b/internal/order/go.sum
@@ -7,8 +7,12 @@ github.com/benbjohnson/clock v1.1.0/go.mod h1:J11/hYXuz8f4ySSvYwY0FKfm+ezbsZBKZx
 github.com/bmatcuk/doublestar v1.1.1/go.mod h1:UD6OnuiIn0yFxxA2le/rnRU1G4RaI4UvFv1sNto9p6w=
 github.com/bytedance/sonic v1.11.6 h1:oUp34TzMlL+OY1OUWxHqsdkgC/Zfc85zGqw9siXjrc0=
 github.com/bytedance/sonic v1.11.6/go.mod h1:LysEHSvpvDySVdC2f87zGWf6CIKJcAvqab1ZaiQtds4=
+github.com/bytedance/sonic v1.12.3 h1:W2MGa7RCU1QTeYRTPE3+88mVC0yXmsRQRChiyVocVjU=
+github.com/bytedance/sonic v1.12.3/go.mod h1:B8Gt/XvtZ3Fqj+iSKMypzymZxw/FVwgIGKzMzT9r/rk=
 github.com/bytedance/sonic/loader v0.1.1 h1:c+e5Pt1k/cy5wMveRDyk2X4B9hF4g7an8N3zCYjJFNM=
 github.com/bytedance/sonic/loader v0.1.1/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
+github.com/bytedance/sonic/loader v0.2.0 h1:zNprn+lsIP06C/IqCHs3gPQIvnvpKbbxyXQP1iU4kWM=
+github.com/bytedance/sonic/loader v0.2.0/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
 github.com/census-instrumentation/opencensus-proto v0.2.1/go.mod h1:f6KPmirojxKA12rnyqOA5BBL4O983OfeGPqjHWSTneU=
 github.com/client9/misspell v0.3.4/go.mod h1:qj6jICC3Q7zFZvVWo7KLAzC3yx5G7kyvSDkc90ppPyw=
 github.com/cloudwego/base64x v0.1.4 h1:jwCgWpFanWmN8xoIUHa2rtzmkd5J2plF/dnLS6Xd/0Y=
@@ -30,6 +34,8 @@ github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nos
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
 github.com/gabriel-vasile/mimetype v1.4.3 h1:in2uUcidCuFcDKtdcBxlR0rJ1+fsokWf+uqxgUFjbI0=
 github.com/gabriel-vasile/mimetype v1.4.3/go.mod h1:d8uq/6HKRL6CGdk+aubisF/M5GcPfT7nKyLpA0lbSSk=
+github.com/gabriel-vasile/mimetype v1.4.6 h1:3+PzJTKLkvgjeTbts6msPJt4DixhT4YtFNf1gtGe3zc=
+github.com/gabriel-vasile/mimetype v1.4.6/go.mod h1:JX1qVKqZd40hUPpAfiNTe0Sne7hdfKSbOqqmkq8GCXc=
 github.com/gin-contrib/sse v0.1.0 h1:Y/yl/+YNO8GZSjAhjMsSuLt29uWRFHdHYUb5lYOV9qE=
 github.com/gin-contrib/sse v0.1.0/go.mod h1:RHrZQHXnP2xjPF+u1gW/2HnVO7nvIa9PG3Gm+fLHvGI=
 github.com/gin-gonic/gin v1.10.0 h1:nTuyha1TYqgedzytsKYqna+DfLos46nTv2ygFy86HFU=
@@ -44,9 +50,13 @@ github.com/go-playground/universal-translator v0.18.1 h1:Bcnm0ZwsGyWbCzImXv+pAJn
 github.com/go-playground/universal-translator v0.18.1/go.mod h1:xekY+UJKNuX9WP91TpwSH2VMlDf28Uj24BCp08ZFTUY=
 github.com/go-playground/validator/v10 v10.20.0 h1:K9ISHbSaI0lyB2eWMPJo+kOS/FBExVwjEviJTixqxL8=
 github.com/go-playground/validator/v10 v10.20.0/go.mod h1:dbuPbCMFw/DrkbEynArYaCwl3amGuJotoKCe95atGMM=
+github.com/go-playground/validator/v10 v10.22.1 h1:40JcKH+bBNGFczGuoBYgX4I6m/i27HYW8P9FDk5PbgA=
+github.com/go-playground/validator/v10 v10.22.1/go.mod h1:dbuPbCMFw/DrkbEynArYaCwl3amGuJotoKCe95atGMM=
 github.com/go-stack/stack v1.8.0/go.mod h1:v0f6uXyyMGvRgIKkXu+yp6POWl0qKG85gN/melR3HDY=
 github.com/goccy/go-json v0.10.2 h1:CrxCmQqYDkv1z7lO7Wbh2HN93uovUHgrECaO5ZrCXAU=
 github.com/goccy/go-json v0.10.2/go.mod h1:6MelG93GURQebXPDq3khkgXZkazVtN9CRI+MGFi0w8I=
+github.com/goccy/go-json v0.10.3 h1:KZ5WoDbxAIgm2HNbYckL0se1fHD6rz5j4ywS6ebzDqA=
+github.com/goccy/go-json v0.10.3/go.mod h1:oq7eo15ShAhp70Anwd5lgX2pLfOS3QCiwU/PULtXL6M=
 github.com/gogo/protobuf v1.3.2 h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=
 github.com/gogo/protobuf v1.3.2/go.mod h1:P1XiOD3dCwIKUDQYPy72D8LYyHL2YPYrpS2s69NZV8Q=
 github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b/go.mod h1:SBH7ygxi8pfUlaOkMMuAQtPIUF8ecWP5IEl/CR7VP2Q=
@@ -57,6 +67,8 @@ github.com/golang/protobuf v1.3.3/go.mod h1:vzj43D7+SQXF/4pzW/hwtAqwc6iTitCiVSaW
 github.com/golang/protobuf v1.5.0/go.mod h1:FsONVRAS9T7sI+LIUmWTfcYkHO4aIWwzhcaSAoJOfIk=
 github.com/golang/protobuf v1.5.3 h1:KhyjKVUg7Usr/dYsdSqoFveMYd5ko72D+zANwlG1mmg=
 github.com/golang/protobuf v1.5.3/go.mod h1:XVQd3VNwM+JqD3oG2Ue2ip4fOMUkwXdXDdiuN0vRsmY=
+github.com/golang/protobuf v1.5.4 h1:i7eJL8qZTpSEXOPTxNKhASYpMn+8e5Q6AdndVa1dWek=
+github.com/golang/protobuf v1.5.4/go.mod h1:lnTiLA8Wa4RWRcIUkrtSVa5nRhsEGBg48fD6rSs7xps=
 github.com/google/go-cmp v0.2.0/go.mod h1:oXzfMopK8JAjlY9xF4vHSVASa0yLyX7SntLO5aqRK0M=
 github.com/google/go-cmp v0.5.5/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
 github.com/google/go-cmp v0.6.0 h1:ofyhxvXcZhMsU5ulbFiLKl/XBFqE1GSq7atu8tAmTRI=
@@ -76,6 +88,8 @@ github.com/kisielk/gotool v1.0.0/go.mod h1:XhKaO+MFFWcvkIS/tQcRk01m1F5IRFswLeQ+o
 github.com/klauspost/cpuid/v2 v2.0.9/go.mod h1:FInQzS24/EEf25PyTYn52gqo7WaD8xa0213Md/qVLRg=
 github.com/klauspost/cpuid/v2 v2.2.7 h1:ZWSB3igEs+d0qvnxR/ZBzXVmxkgt8DdzP6m9pfuVLDM=
 github.com/klauspost/cpuid/v2 v2.2.7/go.mod h1:Lcz8mBdAVJIBVzewtcLocK12l3Y+JytZYpaMropDUws=
+github.com/klauspost/cpuid/v2 v2.2.8 h1:+StwCXwm9PdpiEkPyzBXIy+M9KUb4ODm0Zarf1kS5BM=
+github.com/klauspost/cpuid/v2 v2.2.8/go.mod h1:Lcz8mBdAVJIBVzewtcLocK12l3Y+JytZYpaMropDUws=
 github.com/knz/go-libedit v1.10.1/go.mod h1:MZTVkCWyz0oBc7JOWP3wNAzd002ZbM/5hgShxwh4x8M=
 github.com/konsorten/go-windows-terminal-sequences v1.0.1/go.mod h1:T0+1ngSBFLxvqU3pZ+m/2kptfBszLMUkC4ZK/EgS/cQ=
 github.com/kr/pretty v0.1.0/go.mod h1:dAy3ld7l9f0ibDNOQOHHMYYIIbhfbHSm3C4ZsoJORNo=
@@ -103,6 +117,8 @@ github.com/oapi-codegen/runtime v1.1.1/go.mod h1:SK9X900oXmPWilYR5/WKPzt3Kqxn/uS
 github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFStlNbqXla1AfSYxPUl2o=
 github.com/pelletier/go-toml/v2 v2.2.2 h1:aYUidT7k73Pcl9nb2gScu7NSrKCSHIDE89b3+6Wq+LM=
 github.com/pelletier/go-toml/v2 v2.2.2/go.mod h1:1t835xjRzz80PqgE6HHgN2JOsmgYu/h4qDAS4n929Rs=
+github.com/pelletier/go-toml/v2 v2.2.3 h1:YmeHyLY8mFWbdkNWwpr+qIL2bEqT0o95WSdkNHvL12M=
+github.com/pelletier/go-toml/v2 v2.2.3/go.mod h1:MfCQTFTvCcUyyvvwm1+G6H/jORL20Xlb6rzQu9GuUkc=
 github.com/pkg/errors v0.8.1/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
 github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
 github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 h1:Jamvg5psRIccs7FGNTlIRMkT8wgtp5eCXdBlqhYGL6U=
@@ -112,17 +128,23 @@ github.com/rogpeppe/go-internal v1.9.0 h1:73kH8U+JUqXU8lRuOHeVHaa/SZPifC7BkcraZV
 github.com/rogpeppe/go-internal v1.9.0/go.mod h1:WtVeX8xhTBvf0smdhujwtBcq4Qrzq/fJaraNFVN+nFs=
 github.com/sagikazarmark/locafero v0.4.0 h1:HApY1R9zGo4DBgr7dqsTH/JJxLTTsOt7u6keLGt6kNQ=
 github.com/sagikazarmark/locafero v0.4.0/go.mod h1:Pe1W6UlPYUk/+wc/6KFhbORCfqzgYEpgQ3O5fPuL3H4=
+github.com/sagikazarmark/locafero v0.6.0 h1:ON7AQg37yzcRPU69mt7gwhFEBwxI6P9T4Qu3N51bwOk=
+github.com/sagikazarmark/locafero v0.6.0/go.mod h1:77OmuIc6VTraTXKXIs/uvUxKGUXjE1GbemJYHqdNjX0=
 github.com/sagikazarmark/slog-shim v0.1.0 h1:diDBnUNK9N/354PgrxMywXnAwEr1QZcOr6gto+ugjYE=
 github.com/sagikazarmark/slog-shim v0.1.0/go.mod h1:SrcSrq8aKtyuqEI1uvTDTK1arOWRIczQRv+GVI1AkeQ=
 github.com/sirupsen/logrus v1.4.2/go.mod h1:tLMulIdttU9McNUspp0xgXVQah82FyeX6MwdIuYE2rE=
 github.com/sirupsen/logrus v1.8.1 h1:dJKuHgqk1NNQlqoA6BTlM1Wf9DOH3NBjQyu0h9+AZZE=
 github.com/sirupsen/logrus v1.8.1/go.mod h1:yWOB1SBYBC5VeMP7gHvWumXLIWorT60ONWic61uBYv0=
+github.com/sirupsen/logrus v1.9.3 h1:dueUQJ1C2q9oE3F7wvmSGAaVtTmUizReu6fjN8uqzbQ=
+github.com/sirupsen/logrus v1.9.3/go.mod h1:naHLuLoDiP4jHNo9R0sCBMtWGeIprob74mVsIT4qYEQ=
 github.com/sourcegraph/conc v0.3.0 h1:OQTbbt6P72L20UqAkXXuLOj79LfEanQ+YQFNpLA9ySo=
 github.com/sourcegraph/conc v0.3.0/go.mod h1:Sdozi7LEKbFPqYX2/J+iBAM6HpqSLTASQIKqDmF7Mt0=
 github.com/spf13/afero v1.11.0 h1:WJQKhtpdm3v2IzqG8VMqrr6Rf3UYpEF239Jy9wNepM8=
 github.com/spf13/afero v1.11.0/go.mod h1:GH9Y3pIexgf1MTIWtNGyogA5MwRIDXGUr+hbWNoBjkY=
 github.com/spf13/cast v1.6.0 h1:GEiTHELF+vaR5dhz3VqZfFSzZjYbgeKDpBxQVS4GYJ0=
 github.com/spf13/cast v1.6.0/go.mod h1:ancEpBxwJDODSW/UG4rDrAqiKolqNNh2DX3mk86cAdo=
+github.com/spf13/cast v1.7.0 h1:ntdiHjuueXFgm5nzDRdOS4yfT43P5Fnud6DH50rz/7w=
+github.com/spf13/cast v1.7.0/go.mod h1:ancEpBxwJDODSW/UG4rDrAqiKolqNNh2DX3mk86cAdo=
 github.com/spf13/pflag v1.0.5 h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=
 github.com/spf13/pflag v1.0.5/go.mod h1:McXfInJRrz4CZXVZOBLb0bTZqETkiAhM9Iw0y3An2Bg=
 github.com/spf13/viper v1.19.0 h1:RWq5SEjt8o25SROyN3z2OrDB9l7RPd3lwTWU8EcEdcI=
@@ -154,22 +176,32 @@ github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9dec
 go.uber.org/atomic v1.7.0/go.mod h1:fEN4uk6kAWBTFdckzkM89CLk9XfWZrxpCo0nPH17wJc=
 go.uber.org/atomic v1.9.0 h1:ECmE8Bn/WFTYwEW/bpKD3M8VtR/zQVbavAoalC1PYyE=
 go.uber.org/atomic v1.9.0/go.mod h1:fEN4uk6kAWBTFdckzkM89CLk9XfWZrxpCo0nPH17wJc=
+go.uber.org/atomic v1.11.0 h1:ZvwS0R+56ePWxUNi+Atn9dWONBPp/AUETXlHW0DxSjE=
+go.uber.org/atomic v1.11.0/go.mod h1:LUxbIzbOniOlMKjJjyPfpl4v+PKK2cNJn91OQbhoJI0=
 go.uber.org/goleak v1.1.10/go.mod h1:8a7PlsEVH3e/a/GLqe5IIrQx6GzcnRmZEufDUTk4A7A=
 go.uber.org/multierr v1.6.0/go.mod h1:cdWPpRnG4AhwMwsgIHip0KRBQjJy5kYEpYjJxpXp9iU=
 go.uber.org/multierr v1.9.0 h1:7fIwc/ZtS0q++VgcfqFDxSBZVv/Xo49/SYnDFupUwlI=
 go.uber.org/multierr v1.9.0/go.mod h1:X2jQV1h+kxSjClGpnseKVIxpmcjrj7MNnI0bnlfKTVQ=
+go.uber.org/multierr v1.11.0 h1:blXXJkSxSSfBVBlC76pxqeO+LN3aDfLQo+309xJstO0=
+go.uber.org/multierr v1.11.0/go.mod h1:20+QtiLqy0Nd6FdQB9TLXag12DsQkrbs3htMFfDN80Y=
 go.uber.org/zap v1.18.1/go.mod h1:xg/QME4nWcxGxrpdeYfq7UvYrLh66cuVKdrbD1XF/NI=
 golang.org/x/arch v0.0.0-20210923205945-b76863e36670/go.mod h1:5om86z9Hs0C8fWVUuoMHwpExlXzs5Tkyp9hOrfG7pp8=
 golang.org/x/arch v0.8.0 h1:3wRIsP3pM4yUptoR96otTUOXI367OS0+c9eeRi9doIc=
 golang.org/x/arch v0.8.0/go.mod h1:FEVrYAQjsQXMVJ1nsMoVVXPZg6p2JE2mx8psSWTDQys=
+golang.org/x/arch v0.11.0 h1:KXV8WWKCXm6tRpLirl2szsO5j/oOODwZf4hATmGVNs4=
+golang.org/x/arch v0.11.0/go.mod h1:FEVrYAQjsQXMVJ1nsMoVVXPZg6p2JE2mx8psSWTDQys=
 golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACkg1iLfiJU5Ep61QUkGW8qpdssI0+w=
 golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
 golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9/go.mod h1:LzIPMQfyMNhhGPhUkYOs5KpL4U8rLKemX1yGLhDgUto=
 golang.org/x/crypto v0.23.0 h1:dIJU/v2J8Mdglj/8rJ6UUOM3Zc9zLZxVZwwxMooUSAI=
 golang.org/x/crypto v0.23.0/go.mod h1:CKFgDieR+mRhux2Lsu27y0fO304Db0wZe70UKqHu0v8=
+golang.org/x/crypto v0.28.0 h1:GBDwsMXVQi34v5CCYUm2jkJvu4cbtru2U4TN2PSyQnw=
+golang.org/x/crypto v0.28.0/go.mod h1:rmgy+3RHxRZMyY0jjAJShp2zgEdOqj2AO7U0pYmeQ7U=
 golang.org/x/exp v0.0.0-20190121172915-509febef88a4/go.mod h1:CJ0aWSM057203Lf6IL+f9T1iT9GByDxfZKAQTCR3kQA=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9 h1:GoHiUyI/Tp2nVkLI2mCxVkOjsbSXD66ic0XW0js0R9g=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9/go.mod h1:S2oDrQGGwySpoQPVqRShND87VCbxmc6bL1Yd2oYrm6k=
+golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c h1:7dEasQXItcW1xKJ2+gg5VOiBnqWrJc+rq0DPKyvvdbY=
+golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c/go.mod h1:NQtJDoLvd6faHhE7m4T/1IY708gDefGGjR/iUW8yQQ8=
 golang.org/x/lint v0.0.0-20181026193005-c67002cb31c3/go.mod h1:UVdnD1Gm6xHRNCYTkRU2/jEulfH38KcIWyp/GAMgvoE=
 golang.org/x/lint v0.0.0-20190227174305-5b3e6a55c961/go.mod h1:wehouNa3lNwaWXcvxsM5YxQ5yQlVC4a0KAMCusXpPoU=
 golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3/go.mod h1:6SW0HCj/g11FgYtHlgUYUwCkIfeOF89ocIRzGO/8vkc=
@@ -186,6 +218,8 @@ golang.org/x/net v0.0.0-20200226121028-0de0cce0169b/go.mod h1:z5CRVTTTmAJ677TzLL
 golang.org/x/net v0.0.0-20201021035429-f5854403a974/go.mod h1:sp8m0HH+o8qH0wwXwYZr8TS3Oi6o0r6Gce1SSxlDquU=
 golang.org/x/net v0.25.0 h1:d/OCCoBEUq33pjydKrGQhw7IlUPI2Oylr+8qLx49kac=
 golang.org/x/net v0.25.0/go.mod h1:JkAGAh7GEvH74S6FOH42FLoXpXbE/aqXSrIQjXgsiwM=
+golang.org/x/net v0.30.0 h1:AcW1SDZMkb8IpzCdQUaIq2sP4sZ4zw+55h6ynffypl4=
+golang.org/x/net v0.30.0/go.mod h1:2wGyMJ5iFasEhkwi13ChkO/t1ECNC4X4eBKkVFyYFlU=
 golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be/go.mod h1:N/0e6XlmueqKjAGxoOufVs8QHGRruUQn6yWY3a++T0U=
 golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20181108010431-42b317875d0f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
@@ -199,14 +233,19 @@ golang.org/x/sys v0.0.0-20190422165155-953cdadca894/go.mod h1:h1NjWce9XRLGQEsW7w
 golang.org/x/sys v0.0.0-20191026070338-33540a1f6037/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.20.0 h1:Od9JTbYCk261bKm4M/mw7AklTlFYIa0bIp9BgSm1S8Y=
 golang.org/x/sys v0.20.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
+golang.org/x/sys v0.26.0 h1:KHjCJyddX0LoSTb3J+vWpupP9p0oznkqVk/IfjymZbo=
+golang.org/x/sys v0.26.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
 golang.org/x/text v0.15.0 h1:h1V/4gjBv8v9cjcR6+AR5+/cIYK5N/WAgiv4xlsEtAk=
 golang.org/x/text v0.15.0/go.mod h1:18ZOQIKpY8NJVqYksKHtTdi31H5itFRjB5/qKTNYzSU=
+golang.org/x/text v0.19.0 h1:kTxAhCbGbxhK0IwgSKiMO5awPoDQ0RpfiVYBfK860YM=
+golang.org/x/text v0.19.0/go.mod h1:BuEKDfySbSR4drPmRPG/7iBdf8hvFMuRexcpahXilzY=
 golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190114222345-bf090417da8b/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190226205152-f727befe758c/go.mod h1:9Yl7xja0Znq3iFh3HoIrodX9oNMXvdceNzlUR8zjMvY=
@@ -227,6 +266,8 @@ google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55/go.mod h1:DMBHOl98
 google.golang.org/genproto v0.0.0-20200423170343-7949de9c1215/go.mod h1:55QSHmfGQM9UVYDPBsyGGes0y52j32PQ3BqQfXhyH3c=
 google.golang.org/genproto/googleapis/rpc v0.0.0-20240314234333-6e1732d8331c h1:lfpJ/2rWPa/kJgxyyXM8PrNnfCzcmxJ265mADgwmvLI=
 google.golang.org/genproto/googleapis/rpc v0.0.0-20240314234333-6e1732d8331c/go.mod h1:WtryC6hu0hhx87FDGxWCDptyssuo68sk10vYjF+T9fY=
+google.golang.org/genproto/googleapis/rpc v0.0.0-20241007155032-5fefd90f89a9 h1:QCqS/PdaHTSWGvupk2F/ehwHtGc0/GYkT+3GAcR1CCc=
+google.golang.org/genproto/googleapis/rpc v0.0.0-20241007155032-5fefd90f89a9/go.mod h1:GX3210XPVPUjJbTUbvwI8f2IpZDMZuPJWDzDuebbviI=
 google.golang.org/grpc v1.19.0/go.mod h1:mqu4LbDTu4XGKhr4mRzUsmM4RtVoemTSY81AxZiDr8c=
 google.golang.org/grpc v1.23.0/go.mod h1:Y5yQAOtifL1yxbo5wqy6BxZv8vAUGQwXBOALyacEbxg=
 google.golang.org/grpc v1.25.1/go.mod h1:c3i+UQWmh7LiEpx4sFZnkU36qjEYZ0imhYfXVyQciAY=
@@ -234,10 +275,14 @@ google.golang.org/grpc v1.27.0/go.mod h1:qbnxyOmOxrQa7FizSgH+ReBfzJrCY1pSN7KXBS8
 google.golang.org/grpc v1.29.1/go.mod h1:itym6AZVZYACWQqET3MqgPpjcuV5QH3BxFS3IjizoKk=
 google.golang.org/grpc v1.62.1 h1:B4n+nfKzOICUXMgyrNd19h/I9oH0L1pizfk1d4zSgTk=
 google.golang.org/grpc v1.62.1/go.mod h1:IWTG0VlJLCh1SkC58F7np9ka9mx/WNkjl4PGJaiq+QE=
+google.golang.org/grpc v1.67.1 h1:zWnc1Vrcno+lHZCOofnIMvycFcc0QRGIzm9dhnDX68E=
+google.golang.org/grpc v1.67.1/go.mod h1:1gLDyUQU7CTLJI90u3nXZ9ekeghjeM7pTDZlqFNg2AA=
 google.golang.org/protobuf v1.26.0-rc.1/go.mod h1:jlhhOSvTdKEhbULTjvd4ARK9grFBp09yW+WbY/TyQbw=
 google.golang.org/protobuf v1.26.0/go.mod h1:9q0QmTI4eRPtz6boOQmLYwt+qCgq0jsYwAQnmE0givc=
 google.golang.org/protobuf v1.34.1 h1:9ddQBjfCyZPOHPUiPxpYESBLc+T8P3E+Vo4IbKZgFWg=
 google.golang.org/protobuf v1.34.1/go.mod h1:c6P6GXX6sHbq/GpV6MGZEdwhWPcYBgnhAHhKbcUYpos=
+google.golang.org/protobuf v1.35.1 h1:m3LfL6/Ca+fqnjnlqQXNpFPABW1UD7mjh8KO2mKFytA=
+google.golang.org/protobuf v1.35.1/go.mod h1:9fA7Ob0pmnwhb644+1+CVWFRbNajQ6iRojtC/QF5bRE=
 gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 h1:YR8cESwS4TdDjEe65xsg0ogRM/Nc3DYOhEAlW+xobZo=
diff --git a/internal/order/main.go b/internal/order/main.go
index 6d2f2c5..55d022b 100644
--- a/internal/order/main.go
+++ b/internal/order/main.go
@@ -26,7 +26,8 @@ func main() {
 	ctx, cancel := context.WithCancel(context.Background())
 	defer cancel()
 
-	application := service.NewApplication(ctx)
+	application, cleanup := service.NewApplication(ctx)
+	defer cleanup()
 
 	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
 		svc := ports.NewGRPCServer(application)
diff --git a/internal/order/service/application.go b/internal/order/service/application.go
index d8fd936..648e438 100644
--- a/internal/order/service/application.go
+++ b/internal/order/service/application.go
@@ -3,21 +3,34 @@ package service
 import (
 	"context"
 
+	grpcClient "github.com/ghost-yu/go_shop_second/common/client"
 	"github.com/ghost-yu/go_shop_second/common/metrics"
 	"github.com/ghost-yu/go_shop_second/order/adapters"
+	"github.com/ghost-yu/go_shop_second/order/adapters/grpc"
 	"github.com/ghost-yu/go_shop_second/order/app"
 	"github.com/ghost-yu/go_shop_second/order/app/command"
 	"github.com/ghost-yu/go_shop_second/order/app/query"
 	"github.com/sirupsen/logrus"
 )
 
-func NewApplication(ctx context.Context) app.Application {
+func NewApplication(ctx context.Context) (app.Application, func()) {
+	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
+	if err != nil {
+		panic(err)
+	}
+	stockGRPC := grpc.NewStockGRPC(stockClient)
+	return newApplication(ctx, stockGRPC), func() {
+		_ = closeStockClient()
+	}
+}
+
+func newApplication(_ context.Context, stockGRPC query.StockService) app.Application {
 	orderRepo := adapters.NewMemoryOrderRepository()
 	logger := logrus.NewEntry(logrus.StandardLogger())
 	metricClient := metrics.TodoMetrics{}
 	return app.Application{
 		Commands: app.Commands{
-			CreateOrder: command.NewCreateOrderHandler(orderRepo, logger, metricClient),
+			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, logger, metricClient),
 			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
 		},
 		Queries: app.Queries{
~~~

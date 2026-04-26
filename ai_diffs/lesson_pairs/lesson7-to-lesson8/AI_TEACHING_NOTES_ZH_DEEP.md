# `lesson7 -> lesson8` 手写教学稿（人工讲解版）

这份笔记是按真实读代码的顺序手写的，不是脚本套模板。你可以把它当成我们坐在一起过 `lesson7 -> lesson8` 这组差异时的讲义。

我先把标准定清楚：
- 我不会把过渡态代码硬说成最佳实践。
- 哪些地方只是“先把链路打通”，我会直说。
- 哪些地方容易踩坑，我会直接指出来。
- `go.mod` 和 `go.sum` 我会保留在 diff 全文里，但正文不逐行讲，因为你前面已经明确说过，这种依赖噪音不是你现在最该花时间的地方。

## 1. 这节课到底在做什么

如果只用一句话概括，这一节不是“订单功能写完了”，而是：

`order` 服务第一次从“有仓储雏形”升级成“有应用层、有用例、有外部库存依赖、有日志和指标包装”的服务骨架。

你要抓住这条主线：

`main.go -> service.NewApplication -> app.Commands / app.Queries -> handler -> repository / grpc adapter`

只要这条链你看懂了，这一节就不会乱。

## 2. 先把完整差异全文贴出来

下面先贴 `git diff lesson7..lesson8` 原文。你先不要急着逐行抠，先建立“哪些文件变了、改动大概往哪走”的整体印象。

```diff
diff --git a/internal/common/decorator/command.go b/internal/common/decorator/command.go
new file mode 100644
index 0000000..53f8c37
--- /dev/null
+++ b/internal/common/decorator/command.go
@@ -0,0 +1,21 @@
+package decorator
+
+import (
+	"context"
+
+	"github.com/sirupsen/logrus"
+)
+
+type CommandHandler[C, R any] interface {
+	Handle(ctx context.Context, cmd C) (R, error)
+}
+
+func ApplyCommandDecorators[C, R any](handler CommandHandler[C, R], logger *logrus.Entry, metricsClient MetricsClient) CommandHandler[C, R] {
+	return queryLoggingDecorator[C, R]{
+		logger: logger,
+		base: queryMetricsDecorator[C, R]{
+			base:   handler,
+			client: metricsClient,
+		},
+	}
+}
diff --git a/internal/common/decorator/logging.go b/internal/common/decorator/logging.go
new file mode 100644
index 0000000..a9a2325
--- /dev/null
+++ b/internal/common/decorator/logging.go
@@ -0,0 +1,34 @@
+package decorator
+
+import (
+	"context"
+	"fmt"
+	"strings"
+
+	"github.com/sirupsen/logrus"
+)
+
+type queryLoggingDecorator[C, R any] struct {
+	logger *logrus.Entry
+	base   QueryHandler[C, R]
+}
+
+func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
+	logger := q.logger.WithFields(logrus.Fields{
+		"query":      generateActionName(cmd),
+		"query_body": fmt.Sprintf("%#v", cmd),
+	})
+	logger.Debug("Executing query")
+	defer func() {
+		if err == nil {
+			logger.Info("Query execute successfully")
+		} else {
+			logger.Error("Failed to execute query", err)
+		}
+	}()
+	return q.base.Handle(ctx, cmd)
+}
+
+func generateActionName(cmd any) string {
+	return strings.Split(fmt.Sprintf("%T", cmd), ".")[1]
+}
diff --git a/internal/common/decorator/metrics.go b/internal/common/decorator/metrics.go
new file mode 100644
index 0000000..fba9a77
--- /dev/null
+++ b/internal/common/decorator/metrics.go
@@ -0,0 +1,32 @@
+package decorator
+
+import (
+	"context"
+	"fmt"
+	"strings"
+	"time"
+)
+
+type MetricsClient interface {
+	Inc(key string, value int)
+}
+
+type queryMetricsDecorator[C, R any] struct {
+	base   QueryHandler[C, R]
+	client MetricsClient
+}
+
+func (q queryMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
+	start := time.Now()
+	actionName := strings.ToLower(generateActionName(cmd))
+	defer func() {
+		end := time.Since(start)
+		q.client.Inc(fmt.Sprintf("querys.%s.duration", actionName), int(end.Seconds()))
+		if err == nil {
+			q.client.Inc(fmt.Sprintf("querys.%s.success", actionName), 1)
+		} else {
+			q.client.Inc(fmt.Sprintf("querys.%s.failure", actionName), 1)
+		}
+	}()
+	return q.base.Handle(ctx, cmd)
+}
diff --git a/internal/common/decorator/query.go b/internal/common/decorator/query.go
new file mode 100644
index 0000000..d3848fe
--- /dev/null
+++ b/internal/common/decorator/query.go
@@ -0,0 +1,23 @@
+package decorator
+
+import (
+	"context"
+
+	"github.com/sirupsen/logrus"
+)
+
+// QueryHandler defines a generic type that receives a Query Q,
+// and returns a result R
+type QueryHandler[Q, R any] interface {
+	Handle(ctx context.Context, query Q) (R, error)
+}
+
+func ApplyQueryDecorators[H, R any](handler QueryHandler[H, R], logger *logrus.Entry, metricsClient MetricsClient) QueryHandler[H, R] {
+	return queryLoggingDecorator[H, R]{
+		logger: logger,
+		base: queryMetricsDecorator[H, R]{
+			base:   handler,
+			client: metricsClient,
+		},
+	}
+}
diff --git a/internal/common/metrics/todo_metrics.go b/internal/common/metrics/todo_metrics.go
new file mode 100644
index 0000000..a05f6b2
--- /dev/null
+++ b/internal/common/metrics/todo_metrics.go
@@ -0,0 +1,6 @@
+package metrics
+
+type TodoMetrics struct{}
+
+func (t TodoMetrics) Inc(_ string, _ int) {
+}
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
diff --git a/internal/order/adapters/order_inmem_repository.go b/internal/order/adapters/order_inmem_repository.go
index f5d8be2..818b51f 100644
--- a/internal/order/adapters/order_inmem_repository.go
+++ b/internal/order/adapters/order_inmem_repository.go
@@ -16,13 +16,21 @@ type MemoryOrderRepository struct {
 }
 
 func NewMemoryOrderRepository() *MemoryOrderRepository {
+	s := make([]*domain.Order, 0)
+	s = append(s, &domain.Order{
+		ID:          "fake-ID",
+		CustomerID:  "fake-customer-id",
+		Status:      "fake-status",
+		PaymentLink: "fake-payment-link",
+		Items:       nil,
+	})
 	return &MemoryOrderRepository{
 		lock:  &sync.RWMutex{},
-		store: make([]*domain.Order, 0),
+		store: s,
 	}
 }
 
-func (m MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
+func (m *MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
 	m.lock.Lock()
 	defer m.lock.Unlock()
 	newOrder := &domain.Order{
@@ -36,23 +44,26 @@ func (m MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*
 	logrus.WithFields(logrus.Fields{
 		"input_order":        order,
 		"store_after_create": m.store,
-	}).Debug("memory_order_repo_create")
+	}).Info("memory_order_repo_create")
 	return newOrder, nil
 }
 
-func (m MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*domain.Order, error) {
+func (m *MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*domain.Order, error) {
+	for i, v := range m.store {
+		logrus.Infof("m.store[%d] = %+v", i, v)
+	}
 	m.lock.RLock()
 	defer m.lock.RUnlock()
 	for _, o := range m.store {
 		if o.ID == id && o.CustomerID == customerID {
-			logrus.Debugf("memory_order_repo_get||found||id=%s||customerID=%s||res=%+v", id, customerID, *o)
+			logrus.Infof("memory_order_repo_get||found||id=%s||customerID=%s||res=%+v", id, customerID, *o)
 			return o, nil
 		}
 	}
 	return nil, domain.NotFoundError{OrderID: id}
 }
 
-func (m MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
+func (m *MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
 	m.lock.Lock()
 	defer m.lock.Unlock()
 	found := false
diff --git a/internal/order/app/app.go b/internal/order/app/app.go
index 42330cd..b2e04ce 100644
--- a/internal/order/app/app.go
+++ b/internal/order/app/app.go
@@ -1,10 +1,20 @@
 package app
 
+import (
+	"github.com/ghost-yu/go_shop_second/order/app/command"
+	"github.com/ghost-yu/go_shop_second/order/app/query"
+)
+
 type Application struct {
 	Commands Commands
 	Queries  Queries
 }
 
-type Commands struct{}
+type Commands struct {
+	CreateOrder command.CreateOrderHandler
+	UpdateOrder command.UpdateOrderHandler
+}
 
-type Queries struct{}
+type Queries struct {
+	GetCustomerOrder query.GetCustomerOrderHandler
+}
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
new file mode 100644
index 0000000..adb658f
--- /dev/null
+++ b/internal/order/app/command/create_order.go
@@ -0,0 +1,65 @@
+package command
+
+import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/common/decorator"
+	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/order/app/query"
+	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
+	"github.com/sirupsen/logrus"
+)
+
+type CreateOrder struct {
+	CustomerID string
+	Items      []*orderpb.ItemWithQuantity
+}
+
+type CreateOrderResult struct {
+	OrderID string
+}
+
+type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]
+
+type createOrderHandler struct {
+	orderRepo domain.Repository
+	stockGRPC query.StockService
+}
+
+func NewCreateOrderHandler(
+	orderRepo domain.Repository,
+	stockGRPC query.StockService,
+	logger *logrus.Entry,
+	metricClient decorator.MetricsClient,
+) CreateOrderHandler {
+	if orderRepo == nil {
+		panic("nil orderRepo")
+	}
+	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
+		createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC},
+		logger,
+		metricClient,
+	)
+}
+
+func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
+	// TODO: call stock grpc to get items.
+	err := c.stockGRPC.CheckIfItemsInStock(ctx, cmd.Items)
+	resp, err := c.stockGRPC.GetItems(ctx, []string{"123"})
+	logrus.Info("createOrderHandler||resp from stockGRPC.GetItems", resp)
+	var stockResponse []*orderpb.Item
+	for _, item := range cmd.Items {
+		stockResponse = append(stockResponse, &orderpb.Item{
+			ID:       item.ID,
+			Quantity: item.Quantity,
+		})
+	}
+	o, err := c.orderRepo.Create(ctx, &domain.Order{
+		CustomerID: cmd.CustomerID,
+		Items:      stockResponse,
+	})
+	if err != nil {
+		return nil, err
+	}
+	return &CreateOrderResult{OrderID: o.ID}, nil
+}
diff --git a/internal/order/app/command/update_order.go b/internal/order/app/command/update_order.go
new file mode 100644
index 0000000..f40716d
--- /dev/null
+++ b/internal/order/app/command/update_order.go
@@ -0,0 +1,47 @@
+package command
+
+import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/common/decorator"
+	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
+	"github.com/sirupsen/logrus"
+)
+
+type UpdateOrder struct {
+	Order    *domain.Order
+	UpdateFn func(context.Context, *domain.Order) (*domain.Order, error)
+}
+
+type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, interface{}]
+
+type updateOrderHandler struct {
+	orderRepo domain.Repository
+	//stockGRPC
+}
+
+func NewUpdateOrderHandler(
+	orderRepo domain.Repository,
+	logger *logrus.Entry,
+	metricClient decorator.MetricsClient,
+) UpdateOrderHandler {
+	if orderRepo == nil {
+		panic("nil orderRepo")
+	}
+	return decorator.ApplyCommandDecorators[UpdateOrder, interface{}](
+		updateOrderHandler{orderRepo: orderRepo},
+		logger,
+		metricClient,
+	)
+}
+
+func (c updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (interface{}, error) {
+	if cmd.UpdateFn == nil {
+		logrus.Warnf("updateOrderHandler got nil UpdateFn, order=%#v", cmd.Order)
+		cmd.UpdateFn = func(_ context.Context, order *domain.Order) (*domain.Order, error) { return order, nil }
+	}
+	if err := c.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFn); err != nil {
+		return nil, err
+	}
+	return nil, nil
+}
diff --git a/internal/order/app/query/get_customer_order.go b/internal/order/app/query/get_customer_order.go
new file mode 100644
index 0000000..b107e01
--- /dev/null
+++ b/internal/order/app/query/get_customer_order.go
@@ -0,0 +1,43 @@
+package query
+
+import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/common/decorator"
+	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
+	"github.com/sirupsen/logrus"
+)
+
+type GetCustomerOrder struct {
+	CustomerID string
+	OrderID    string
+}
+
+type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Order]
+
+type getCustomerOrderHandler struct {
+	orderRepo domain.Repository
+}
+
+func NewGetCustomerOrderHandler(
+	orderRepo domain.Repository,
+	logger *logrus.Entry,
+	metricClient decorator.MetricsClient,
+) GetCustomerOrderHandler {
+	if orderRepo == nil {
+		panic("nil orderRepo")
+	}
+	return decorator.ApplyQueryDecorators[GetCustomerOrder, *domain.Order](
+		getCustomerOrderHandler{orderRepo: orderRepo},
+		logger,
+		metricClient,
+	)
+}
+
+func (g getCustomerOrderHandler) Handle(ctx context.Context, query GetCustomerOrder) (*domain.Order, error) {
+	o, err := g.orderRepo.Get(ctx, query.OrderID, query.CustomerID)
+	if err != nil {
+		return nil, err
+	}
+	return o, nil
+}
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
diff --git a/internal/order/http.go b/internal/order/http.go
index f18fa80..b40adc7 100644
--- a/internal/order/http.go
+++ b/internal/order/http.go
@@ -1,7 +1,12 @@
 package main
 
 import (
+	"net/http"
+
+	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	"github.com/ghost-yu/go_shop_second/order/app"
+	"github.com/ghost-yu/go_shop_second/order/app/command"
+	"github.com/ghost-yu/go_shop_second/order/app/query"
 	"github.com/gin-gonic/gin"
 )
 
@@ -10,11 +15,34 @@ type HTTPServer struct {
 }
 
 func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
-	//TODO implement me
-	panic("implement me")
+	var req orderpb.CreateOrderRequest
+	if err := c.ShouldBindJSON(&req); err != nil {
+		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
+		return
+	}
+	r, err := H.app.Commands.CreateOrder.Handle(c, command.CreateOrder{
+		CustomerID: req.CustomerID,
+		Items:      req.Items,
+	})
+	if err != nil {
+		c.JSON(http.StatusOK, gin.H{"error": err})
+		return
+	}
+	c.JSON(http.StatusOK, gin.H{
+		"message":     "success",
+		"customer_id": req.CustomerID,
+		"order_id":    r.OrderID,
+	})
 }
 
 func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
-	//TODO implement me
-	panic("implement me")
+	o, err := H.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
+		OrderID:    orderID,
+		CustomerID: customerID,
+	})
+	if err != nil {
+		c.JSON(http.StatusOK, gin.H{"error": err})
+		return
+	}
+	c.JSON(http.StatusOK, gin.H{"message": "success", "data": o})
 }
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
index 122b22a..648e438 100644
--- a/internal/order/service/application.go
+++ b/internal/order/service/application.go
@@ -3,9 +3,38 @@ package service
 import (
 	"context"
 
+	grpcClient "github.com/ghost-yu/go_shop_second/common/client"
+	"github.com/ghost-yu/go_shop_second/common/metrics"
+	"github.com/ghost-yu/go_shop_second/order/adapters"
+	"github.com/ghost-yu/go_shop_second/order/adapters/grpc"
 	"github.com/ghost-yu/go_shop_second/order/app"
+	"github.com/ghost-yu/go_shop_second/order/app/command"
+	"github.com/ghost-yu/go_shop_second/order/app/query"
+	"github.com/sirupsen/logrus"
 )
 
-func NewApplication(ctx context.Context) app.Application {
-	return app.Application{}
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
+	orderRepo := adapters.NewMemoryOrderRepository()
+	logger := logrus.NewEntry(logrus.StandardLogger())
+	metricClient := metrics.TodoMetrics{}
+	return app.Application{
+		Commands: app.Commands{
+			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, logger, metricClient),
+			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
+		},
+		Queries: app.Queries{
+			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
+		},
+	}
 }
diff --git a/internal/stock/adapters/stock_inmem_repository.go b/internal/stock/adapters/stock_inmem_repository.go
index cd0b063..4390124 100644
--- a/internal/stock/adapters/stock_inmem_repository.go
+++ b/internal/stock/adapters/stock_inmem_repository.go
@@ -22,7 +22,7 @@ var stub = map[string]*orderpb.Item{
 	},
 }
 
-func NewMemoryOrderRepository() *MemoryStockRepository {
+func NewMemoryStockRepository() *MemoryStockRepository {
 	return &MemoryStockRepository{
 		lock:  &sync.RWMutex{},
 		store: stub,
```

## 3. 正确阅读顺序

不要按文件名字母顺序读，也不要一上来就扎进 `decorator`。这一节最适合 Go 小白的顺序是：

1. `internal/order/main.go`
2. `internal/order/service/application.go`
3. `internal/order/app/app.go`
4. `internal/order/http.go`
5. `internal/order/app/query/get_customer_order.go`
6. `internal/order/app/command/create_order.go`
7. `internal/order/app/command/update_order.go`
8. `internal/order/app/query/service.go`
9. `internal/order/adapters/grpc/stock_grpc.go`
10. `internal/order/domain/order/repository.go`
11. `internal/order/adapters/order_inmem_repository.go`
12. `internal/stock/adapters/stock_inmem_repository.go`
13. `internal/common/decorator/query.go`
14. `internal/common/decorator/command.go`
15. `internal/common/decorator/logging.go`
16. `internal/common/decorator/metrics.go`
17. `internal/common/metrics/todo_metrics.go`

原因很简单：
- 先看入口，知道程序怎么启动。
- 再看装配，知道依赖怎么接起来。
- 再看端口，知道 HTTP 请求落到哪里。
- 再看 command/query，知道真正业务逻辑在哪。
- 最后再看 decorator，理解日志和指标是怎么包上去的。

## 4. 逐文件手写讲解

### 4.1 `internal/order/main.go`

先看代码全文：

```go
package main

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/config"
	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
	"github.com/ghost-yu/go_shop_second/common/server"
	"github.com/ghost-yu/go_shop_second/order/ports"
	"github.com/ghost-yu/go_shop_second/order/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("order.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		svc := ports.NewGRPCServer(application)
		orderpb.RegisterOrderServiceServer(server, svc)
	})

	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, HTTPServer{
			app: application,
		}, ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})

}
```

逐行解释：
- `package main`：说明这个包会被编译成可执行程序。只要你看到 `main` 包，脑子里就把它当成“服务入口”。
- `import (`：开始引入这个入口文件依赖的所有包。
- `"context"`：Go 里做请求生命周期控制、超时、取消传播，几乎都靠它。你现在先把它理解成“给整条调用链带一个可取消的控制器”。
- `"github.com/.../common/config"`：读配置。这里实际是在给 `viper` 做初始化。
- `"github.com/.../common/genproto/orderpb"`：gRPC 生成代码。后面注册 gRPC 服务时要用它。
- `"github.com/.../common/server"`：公共启动器，负责真正跑 HTTP 和 gRPC server。
- `"github.com/.../order/ports"`：端口层。它负责把 HTTP/gRPC 请求翻译成应用层调用。
- `"github.com/.../order/service"`：应用层装配入口。真正把 repo、grpc client、handler 接起来的是它。
- `"github.com/gin-gonic/gin"`：Gin 是 HTTP Web 框架，这里只在注册 HTTP 路由时用到。
- `"github.com/sirupsen/logrus"`：Logrus 是日志库，这里只在启动失败时直接 `Fatal`。
- `"github.com/spf13/viper"`：Viper 是配置库，负责从配置文件里读 `order.service-name` 之类的键。
- `"google.golang.org/grpc"`：gRPC 官方库，注册 gRPC 服务时需要 `*grpc.Server`。
- `func init()`：Go 的 `init` 会在 `main` 之前自动执行。它适合做初始化，但也最容易把副作用藏起来，所以你要知道它有这个风险。
- `if err := config.NewViperConfig(); err != nil {`：先初始化配置。如果这里失败，后面服务名、地址都拿不到，程序没法正常启动。
- `logrus.Fatal(err)`：打印日志并立刻退出进程。不是简单记录一下，而是直接终止。
- `func main()`：程序主入口。
- `serviceName := viper.GetString("order.service-name")`：从配置里读服务名。Viper 如果 key 不存在，会给你空字符串，不会报错，这是一个很容易忽略的点。
- `ctx, cancel := context.WithCancel(context.Background())`：创建根上下文和取消函数。以后你连数据库、gRPC、消息队列时，都会反复看到这种写法。
- `defer cancel()`：`main` 退出前调用取消函数，通知下游资源可以停止了。
- `application, cleanup := service.NewApplication(ctx)`：这一行是本节核心。入口层不直接 new repo、不直接 new grpc client，而是统一交给 `service` 层组装。
- `defer cleanup()`：说明 `NewApplication` 不只是返回了对象，还可能顺带创建了外部连接，所以需要清理。
- `go server.RunGRPCServer(...)`：gRPC 服务放到 goroutine 里单独跑，这样不会阻塞后面的 HTTP server 启动。
- `svc := ports.NewGRPCServer(application)`：把已经装配好的 application 交给 gRPC 端口层。
- `orderpb.RegisterOrderServiceServer(server, svc)`：把订单服务注册到 gRPC server 上。
- `server.RunHTTPServer(...)`：再启动 HTTP server。
- `ports.RegisterHandlersWithOptions(...)`：把 HTTP handler 注册到 Gin 路由器。
- `app: application`：HTTPServer 只持有 application，不直接碰 repo。这就是分层。
- `BaseURL: "/api"`：HTTP 路由统一挂在 `/api` 前缀下。
- `Middlewares: nil, ErrorHandler: nil`：这里先不挂中间件和统一错误处理，说明还在骨架阶段。

为什么这么写：
- `main.go` 应该只负责“服务怎么启动”，不应该直接写业务逻辑。
- 这样 HTTP 和 gRPC 两种入口都能复用同一套应用层。

容易忘和容易错的点：
- `viper.GetString` 读不到配置时不会报错，只会给空值，排查起来很烦。
- `init()` 很方便，但测试时不容易控制副作用。真实项目里如果滥用，会让初始化顺序越来越难追。

### 4.2 `internal/order/service/application.go`

先看代码全文：

```go
package service

import (
	"context"

	grpcClient "github.com/ghost-yu/go_shop_second/common/client"
	"github.com/ghost-yu/go_shop_second/common/metrics"
	"github.com/ghost-yu/go_shop_second/order/adapters"
	"github.com/ghost-yu/go_shop_second/order/adapters/grpc"
	"github.com/ghost-yu/go_shop_second/order/app"
	"github.com/ghost-yu/go_shop_second/order/app/command"
	"github.com/ghost-yu/go_shop_second/order/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	stockGRPC := grpc.NewStockGRPC(stockClient)
	return newApplication(ctx, stockGRPC), func() {
		_ = closeStockClient()
	}
}

func newApplication(_ context.Context, stockGRPC query.StockService) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, logger, metricClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}
```

逐行解释：
- `package service`：这里不是“业务 service 层”那个老式大杂烩概念，而是装配层。
- `import (`：开始引入装配需要的依赖。
- `"context"`：用于把生命周期传给外部依赖创建过程。
- `grpcClient "github.com/.../common/client"`：别名导入，说明这里会用公共的 gRPC client 工厂。
- `"github.com/.../common/metrics"`：指标实现。当前是占位版，但接口先接上了。
- `"github.com/.../order/adapters"`：订单仓储的具体实现。
- `"github.com/.../order/adapters/grpc"`：库存 gRPC 适配器实现。
- `"github.com/.../order/app"`：应用层总容器。
- `"github.com/.../order/app/command"`：命令处理器构造函数在这里。
- `"github.com/.../order/app/query"`：查询处理器构造函数和库存接口定义在这里。
- `"github.com/sirupsen/logrus"`：统一 logger 入口。
- `func NewApplication(ctx context.Context) (app.Application, func())`：返回两个东西，应用对象和清理函数。这个设计很干净，因为创建资源的人也负责告诉外部怎么释放资源。
- `stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)`：先连库存服务。你现在就把它理解成“把外部依赖接进来”。
- `if err != nil { panic(err) }`：这里直接 `panic`，说明作者认为库存依赖没起来时，订单服务继续启动没有意义。
- `stockGRPC := grpc.NewStockGRPC(stockClient)`：把底层生成的 gRPC client 再包成应用层想要的接口对象。这一步就叫 adapter。
- `return newApplication(ctx, stockGRPC), func() { _ = closeStockClient() }`：返回 application，同时把关闭连接的逻辑一起返回。
- `func newApplication(_ context.Context, stockGRPC query.StockService) app.Application`：真正的组装过程放到内部函数里，输入是抽象的 `StockService`，不是具体的 gRPC client。
- `orderRepo := adapters.NewMemoryOrderRepository()`：当前订单仓储还是内存版，便于教学和调试。
- `logger := logrus.NewEntry(logrus.StandardLogger())`：构造一个统一 logger entry，后面所有 handler 都能复用。
- `metricClient := metrics.TodoMetrics{}`：先塞一个空实现，让接口跑通。
- `return app.Application{`：开始组装整个应用容器。
- `Commands: app.Commands{`：命令入口。
- `CreateOrder: command.NewCreateOrderHandler(...)`：创建订单这个命令需要 repo、库存服务、logger、metrics。
- `UpdateOrder: command.NewUpdateOrderHandler(...)`：更新订单当前只依赖 repo、logger、metrics。
- `Queries: app.Queries{`：查询入口。
- `GetCustomerOrder: query.NewGetCustomerOrderHandler(...)`：查询订单当前只依赖 repo、logger、metrics。

为什么这么写：
- 这是典型的依赖注入思路。上层只要知道从这里拿 `Application`，不用知道底层到底接的是内存 repo、MySQL repo 还是 mock。
- `query.StockService` 用接口而不是具体 gRPC 类型，是为了让应用层依赖抽象，不依赖具体通信协议。

容易忘和容易错的点：
- 这里只检查了 `NewStockGRPCClient` 的错误，但没有降级方案，所以库存服务挂了时整个订单服务也会挂。教学项目里可以接受，生产里通常要更明确地评估依赖级别。
- 这版指标实现是空的，所以你看到 decorator 里在记 metrics，不代表真的已经有 Prometheus 或 OpenTelemetry 指标落地了。

### 4.3 `internal/order/app/app.go`

先看代码全文：

```go
package app

import (
	"github.com/ghost-yu/go_shop_second/order/app/command"
	"github.com/ghost-yu/go_shop_second/order/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateOrder command.CreateOrderHandler
	UpdateOrder command.UpdateOrderHandler
}

type Queries struct {
	GetCustomerOrder query.GetCustomerOrderHandler
}
```

逐行解释：
- `package app`：这里放的是应用层对外暴露的能力集合。
- `import (`：只引入 command 和 query 两类能力。
- `type Application struct {`：整个应用容器。
- `Commands Commands`：所有写操作都收拢到这里。
- `Queries Queries`：所有读操作都收拢到这里。
- `type Commands struct {`：命令集合。
- `CreateOrder command.CreateOrderHandler`：创建订单命令入口。
- `UpdateOrder command.UpdateOrderHandler`：更新订单命令入口。
- `type Queries struct {`：查询集合。
- `GetCustomerOrder query.GetCustomerOrderHandler`：查询指定订单入口。

为什么这么写：
- 这就是把系统对外能做什么明确列出来。HTTP、gRPC、定时任务、消息消费者以后都不需要知道 repo 在哪，只要拿着 `Application` 调命令或查询即可。
- 把读写分开，是为了让“查询不产生副作用、命令负责状态变更”这个约束更自然。

### 4.4 `internal/order/http.go`

先看代码全文：

```go
package main

import (
	"net/http"

	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
	"github.com/ghost-yu/go_shop_second/order/app"
	"github.com/ghost-yu/go_shop_second/order/app/command"
	"github.com/ghost-yu/go_shop_second/order/app/query"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	app app.Application
}

func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	var req orderpb.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r, err := H.app.Commands.CreateOrder.Handle(c, command.CreateOrder{
		CustomerID: req.CustomerID,
		Items:      req.Items,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":     "success",
		"customer_id": req.CustomerID,
		"order_id":    r.OrderID,
	})
}

func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	o, err := H.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
		OrderID:    orderID,
		CustomerID: customerID,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": o})
}
```

逐行解释：
- `package main`：因为这个 HTTP handler 跟当前可执行服务放在同一个包里。
- `"net/http"`：用 HTTP 状态码常量。
- `".../orderpb"`：这里用到的是生成出来的请求结构体 `CreateOrderRequest`。
- `".../app"`：HTTPServer 持有的是 application，而不是 repo。
- `".../command"`：POST 请求会转成 `CreateOrder` 命令。
- `".../query"`：GET 请求会转成 `GetCustomerOrder` 查询。
- `"github.com/gin-gonic/gin"`：Gin 的 `Context`、`H` 都来自这里。
- `type HTTPServer struct { app app.Application }`：HTTP 层只持有 application，说明它只负责协议翻译。
- `func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string)`：这是生成的 handler 签名，路径里的 `customerID` 已经被框架解析出来了。
- `var req orderpb.CreateOrderRequest`：准备一个请求对象承接 JSON 解析结果。
- `if err := c.ShouldBindJSON(&req); err != nil {`：Gin 把请求体 JSON 绑定到 Go 结构体。字段名不匹配、类型不对、JSON 不合法，都会在这里报错。
- `c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})`：绑定失败时返回 400，并且这里用了 `err.Error()`，这是对的。
- `r, err := H.app.Commands.CreateOrder.Handle(c, command.CreateOrder{`：把 HTTP 请求翻译成应用层命令，再交给 command handler。
- `CustomerID: req.CustomerID`：这里用了 body 里的 `CustomerID`，但函数参数里已经有路径上的 `customerID` 了。这个差异你一定要注意。
- `Items: req.Items`：把请求里的商品列表传下去。
- `if err != nil { c.JSON(http.StatusOK, gin.H{"error": err}) }`：这里有两个问题。第一，出错时还返回 200，不是很合理。第二，这里直接把 `error` 塞给 JSON，很多 error 类型序列化后并不会给你想要的字符串，最好写成 `err.Error()`。
- `c.JSON(http.StatusOK, gin.H{...})`：成功就返回订单 ID。
- `func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(...)`：读取订单接口。
- `H.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{...})`：把 URL 参数翻译成查询对象。
- `OrderID: orderID, CustomerID: customerID`：这里终于用上了路径参数。
- `if err != nil { c.JSON(http.StatusOK, gin.H{"error": err}) }`：同样有“状态码不对 + error 直接塞 JSON”的问题。
- `c.JSON(http.StatusOK, gin.H{"message": "success", "data": o})`：查到订单就返回。

为什么这么写：
- HTTP 层只做一件事：把外部协议翻译成内部命令/查询。真正业务逻辑不能放在这里，否则以后换成 gRPC 端口时你还得再抄一遍。

外部库和易错点：
- `gin.Context` 既是 HTTP 上下文，也实现了 `context.Context` 的接口，所以它可以直接传给 `Handle`。这个点很多新手第一次会迷糊。
- `ShouldBindJSON` 依赖结构体字段或 tag 能匹配 JSON 字段。绑不上就直接报错。
- Gin 里返回错误时，最好统一用字符串和合适的状态码，不要把原始 `error` 对象直接扔给 JSON。

### 4.5 `internal/order/app/query/get_customer_order.go`

先看代码全文：

```go
package query

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/decorator"
	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
	"github.com/sirupsen/logrus"
)

type GetCustomerOrder struct {
	CustomerID string
	OrderID    string
}

type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Order]

type getCustomerOrderHandler struct {
	orderRepo domain.Repository
}

func NewGetCustomerOrderHandler(
	orderRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) GetCustomerOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}
	return decorator.ApplyQueryDecorators[GetCustomerOrder, *domain.Order](
		getCustomerOrderHandler{orderRepo: orderRepo},
		logger,
		metricClient,
	)
}

func (g getCustomerOrderHandler) Handle(ctx context.Context, query GetCustomerOrder) (*domain.Order, error) {
	o, err := g.orderRepo.Get(ctx, query.OrderID, query.CustomerID)
	if err != nil {
		return nil, err
	}
	return o, nil
}
```

逐行解释：
- `package query`：查询用例放在 query 子包里。
- `"context"`：所有 handler 都接收上下文。
- `".../decorator"`：查询会被日志和指标装饰器包一层。
- `domain ".../order/domain/order"`：查询返回的是领域订单对象。
- `"github.com/sirupsen/logrus"`：构造 handler 时需要注入 logger。
- `type GetCustomerOrder struct {`：这是应用层自己的查询参数对象，不是 HTTP 层请求对象。
- `CustomerID string`：指定客户。
- `OrderID string`：指定订单。
- `type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Order]`：这是个类型别名，意思是“这个 handler 满足泛型查询处理器接口，输入是 `GetCustomerOrder`，输出是 `*domain.Order`”。
- `type getCustomerOrderHandler struct { orderRepo domain.Repository }`：真正业务处理器只依赖仓储。
- `NewGetCustomerOrderHandler(...)`：构造函数。
- `if orderRepo == nil { panic("nil orderRepo") }`：防御式检查，避免装配错了以后运行时才炸。
- `return decorator.ApplyQueryDecorators(...)`：返回的不是裸 handler，而是“业务 handler 外面包了日志和指标”的版本。
- `func (g getCustomerOrderHandler) Handle(...)`：真正执行业务。
- `o, err := g.orderRepo.Get(ctx, query.OrderID, query.CustomerID)`：查仓储。
- `if err != nil { return nil, err }`：查不到就上抛。
- `return o, nil`：查到了就返回。

为什么这么写：
- Query 就该薄。它只负责读数据，不顺手改状态，不顺手调用别的服务，这样测试和推理都简单。

### 4.6 `internal/order/app/command/create_order.go`

先看代码全文：

```go
package command

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/decorator"
	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
	"github.com/ghost-yu/go_shop_second/order/app/query"
	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
	"github.com/sirupsen/logrus"
)

type CreateOrder struct {
	CustomerID string
	Items      []*orderpb.ItemWithQuantity
}

type CreateOrderResult struct {
	OrderID string
}

type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]

type createOrderHandler struct {
	orderRepo domain.Repository
	stockGRPC query.StockService
}

func NewCreateOrderHandler(
	orderRepo domain.Repository,
	stockGRPC query.StockService,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CreateOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}
	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC},
		logger,
		metricClient,
	)
}

func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	// TODO: call stock grpc to get items.
	err := c.stockGRPC.CheckIfItemsInStock(ctx, cmd.Items)
	resp, err := c.stockGRPC.GetItems(ctx, []string{"123"})
	logrus.Info("createOrderHandler||resp from stockGRPC.GetItems", resp)
	var stockResponse []*orderpb.Item
	for _, item := range cmd.Items {
		stockResponse = append(stockResponse, &orderpb.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
		})
	}
	o, err := c.orderRepo.Create(ctx, &domain.Order{
		CustomerID: cmd.CustomerID,
		Items:      stockResponse,
	})
	if err != nil {
		return nil, err
	}
	return &CreateOrderResult{OrderID: o.ID}, nil
}
```

逐行解释：
- `package command`：写操作放到 command 子包。
- `"context"`：命令同样接收上下文。
- `".../decorator"`：命令也会挂日志和指标。
- `".../orderpb"`：这里同时用到了 `ItemWithQuantity` 和 `Item`。
- `".../order/app/query"`：这里依赖的不是具体 gRPC client，而是 query 包里定义的 `StockService` 接口。名字有点绕，但目的是依赖抽象。
- `domain ".../order/domain/order"`：最终要创建的是领域订单。
- `"github.com/sirupsen/logrus"`：调试和日志。
- `type CreateOrder struct { CustomerID string; Items []*orderpb.ItemWithQuantity }`：命令输入模型。你可以把它理解成“应用层版本的下单请求”。
- `type CreateOrderResult struct { OrderID string }`：命令成功后的返回对象。
- `type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]`：命令处理器类型别名。
- `type createOrderHandler struct { orderRepo domain.Repository; stockGRPC query.StockService }`：真正业务处理器依赖两个东西，一个是订单仓储，一个是库存服务。
- `NewCreateOrderHandler(...)`：构造函数。
- `if orderRepo == nil { panic("nil orderRepo") }`：只校验了 repo，没有校验 `stockGRPC`，这点不够完整。
- `return decorator.ApplyCommandDecorators(...)`：和查询一样，外面先包一层日志和指标。
- `func (c createOrderHandler) Handle(...)`：真正的下单逻辑开始。
- `err := c.stockGRPC.CheckIfItemsInStock(ctx, cmd.Items)`：先问库存够不够。这一步方向是对的，下单前先校验库存很合理。
- `resp, err := c.stockGRPC.GetItems(ctx, []string{"123"})`：这里又去拉商品详情，但商品 ID 被硬编码成了 `123`。这说明这版代码还处在“先把 RPC 通路打通”的阶段，不是最终业务版。
- `err` 被这行重新赋值，导致上一行 `CheckIfItemsInStock` 的错误如果发生，会被直接覆盖掉。这是一个明确的 bug/坏味道。
- `logrus.Info("createOrderHandler||resp from stockGRPC.GetItems", resp)`：作者在打印库存 RPC 的返回值，明显是为了确认链路打通。
- `var stockResponse []*orderpb.Item`：准备把请求里的商品列表转换成订单里存的商品对象。
- `for _, item := range cmd.Items { ... }`：循环处理每个商品。
- `ID: item.ID, Quantity: item.Quantity`：这里只保留了 ID 和数量，像商品名、价格 ID 没有从库存服务结果里回填。也就是说，上面拿到的 `resp` 实际上没有真正用起来。
- `o, err := c.orderRepo.Create(ctx, &domain.Order{...})`：把转换好的订单写进仓储。
- `CustomerID: cmd.CustomerID`：订单归属的客户。
- `Items: stockResponse`：订单里的商品快照目前只保留了部分字段。
- `if err != nil { return nil, err }`：仓储写失败就返回错误。
- `return &CreateOrderResult{OrderID: o.ID}, nil`：成功就把订单 ID 返回上去。

为什么这么写：
- 作者的真实意图不是“把创建订单所有细节做到生产级”，而是先把“订单服务会调用库存服务”这条链串起来。
- 对教学来说，这一步很重要，因为从这开始你要学的已经不是单服务 CRUD，而是服务间协作。

外部库和易错点：
- `orderpb` / `stockpb` 是 protobuf 生成类型。你不能把它们当成普通随便改字段名的手写 struct。
- gRPC 调用要始终第一时间检查 `err`，否则很容易像这里一样把前面的错误覆盖掉。
- 这段代码现在有几个明确的待修点：硬编码 `123`、`resp` 未使用、库存校验错误被覆盖、没有校验空商品列表。

### 4.7 `internal/order/app/command/update_order.go`

先看代码全文：

```go
package command

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/decorator"
	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
	"github.com/sirupsen/logrus"
)

type UpdateOrder struct {
	Order    *domain.Order
	UpdateFn func(context.Context, *domain.Order) (*domain.Order, error)
}

type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, interface{}]

type updateOrderHandler struct {
	orderRepo domain.Repository
	//stockGRPC
}

func NewUpdateOrderHandler(
	orderRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) UpdateOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}
	return decorator.ApplyCommandDecorators[UpdateOrder, interface{}](
		updateOrderHandler{orderRepo: orderRepo},
		logger,
		metricClient,
	)
}

func (c updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (interface{}, error) {
	if cmd.UpdateFn == nil {
		logrus.Warnf("updateOrderHandler got nil UpdateFn, order=%#v", cmd.Order)
		cmd.UpdateFn = func(_ context.Context, order *domain.Order) (*domain.Order, error) { return order, nil }
	}
	if err := c.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFn); err != nil {
		return nil, err
	}
	return nil, nil
}
```

逐行解释：
- `package command`：更新订单也是命令。
- `"context"`：上下文照例传递下去。
- `".../decorator"`：命令外层照样可以挂日志和指标。
- `domain ".../order/domain/order"`：更新的对象是领域订单。
- `"github.com/sirupsen/logrus"`：这里只在空 `UpdateFn` 时打 warning。
- `type UpdateOrder struct {`：更新命令。
- `Order *domain.Order`：要更新哪张订单。
- `UpdateFn func(context.Context, *domain.Order) (*domain.Order, error)`：真正怎么改，交给函数参数决定。
- `type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, interface{}]`：返回值现在写成了 `interface{}`，能跑，但语义不够清楚。
- `type updateOrderHandler struct { orderRepo domain.Repository }`：当前只依赖仓储。
- `//stockGRPC`：注释说明作者还在想后面是否把库存依赖接进更新流程。
- `NewUpdateOrderHandler(...)`：构造函数。
- `if orderRepo == nil { panic("nil orderRepo") }`：防御式检查。
- `return decorator.ApplyCommandDecorators(...)`：继续套装饰器。
- `func (c updateOrderHandler) Handle(...)`：真正更新逻辑。
- `if cmd.UpdateFn == nil {`：如果调用方没有传更新函数。
- `logrus.Warnf(...)`：先打一条 warning，告诉你这次更新命令不完整。
- `cmd.UpdateFn = func(_ context.Context, order *domain.Order) (*domain.Order, error) { return order, nil }`：塞一个 no-op 函数，至少让流程不会因为 nil function panic。
- `if err := c.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFn); err != nil {`：真正的更新由 repo 执行。
- `return nil, nil`：当前没有返回额外结果。

为什么这么写：
- “更新订单”不是一个固定动作，而是一类动作。把“怎么更新”作为函数传进来，比预先写死很多个更新方法更灵活。

容易忘和容易错的点：
- `UpdateFn` 这种设计很灵活，但如果调用方把太多业务逻辑塞进去，就会让仓储边界变模糊。
- 返回 `interface{}` 对新手不友好，后面如果能收紧成明确返回类型会更清楚。

### 4.8 `internal/order/app/query/service.go`

先看代码全文：

```go
package query

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
)

type StockService interface {
	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) error
	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
}
```

逐行解释：
- `package query`：这个接口被放在应用层附近，而不是 gRPC 适配器里。
- `"context"`：库存查询同样支持超时和取消传播。
- `".../orderpb"`：这里直接使用 protobuf 生成的商品类型作为参数和返回值。
- `type StockService interface {`：定义库存服务抽象能力。
- `CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) error`：检查库存是否足够。
- `GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)`：批量拿商品详情。

为什么这么写：
- 应用层依赖接口，不依赖具体 gRPC 客户端。以后你可以换成 mock、HTTP、消息总线，应用层不用跟着改。

### 4.9 `internal/order/adapters/grpc/stock_grpc.go`

先看代码全文：

```go
package grpc

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
)

type StockGRPC struct {
	client stockpb.StockServiceClient
}

func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
	return &StockGRPC{client: client}
}

func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) error {
	resp, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
	logrus.Info("stock_grpc response", resp)
	return err
}

func (s StockGRPC) GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error) {
	resp, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{ItemIDs: itemIDs})
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}
```

逐行解释：
- `package grpc`：这个包放的是“通过 gRPC 协议访问外部系统”的适配器。
- `"context"`：RPC 调用天然要带上下文。
- `".../orderpb"`：订单侧和库存侧共用商品相关 protobuf 类型。
- `".../stockpb"`：这里是真正的库存 gRPC 生成 client 所在包。
- `"github.com/sirupsen/logrus"`：当前直接拿来打印响应。
- `type StockGRPC struct { client stockpb.StockServiceClient }`：包了一层底层 gRPC client。
- `NewStockGRPC(...)`：构造适配器。
- `func (s StockGRPC) CheckIfItemsInStock(...) error`：调用库存服务的校验接口。
- `resp, err := s.client.CheckIfItemsInStock(...)`：真正发 gRPC 请求。
- `logrus.Info("stock_grpc response", resp)`：打印响应。能调试，但不够结构化。
- `return err`：这里只把错误返回给上层，不把响应内容返回出去。
- `func (s StockGRPC) GetItems(...)`：调用库存服务的批量查询接口。
- `resp, err := s.client.GetItems(...)`：真正发请求。
- `if err != nil { return nil, err }`：失败直接返回。
- `return resp.Items, nil`：成功只把商品列表交还上层。

为什么这么写：
- 这一层的职责是“翻译协议”。上层不想知道 `stockpb.GetItemsRequest` 这些 gRPC 细节，只想知道“给我商品信息”。

外部库和易错点：
- gRPC 生成的 client 一般都是接口形式，调用方式看起来像本地方法，其实底层是网络请求，所以一定要认真处理 `error` 和上下文超时。
- 这里只记录了响应，没有把调用耗时、请求参数这些结构化字段一起打出来，排查问题时信息会偏少。

### 4.10 `internal/order/domain/order/repository.go`

先看代码全文：

```go
package order

import (
	"context"
	"fmt"
)

type Repository interface {
	Create(context.Context, *Order) (*Order, error)
	Get(ctx context.Context, id, customerID string) (*Order, error)
	Update(
		ctx context.Context,
		o *Order,
		updateFn func(context.Context, *Order) (*Order, error),
	) error
}

type NotFoundError struct {
	OrderID string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("order '%s' not found", e.OrderID)
}
```

逐行解释：
- `package order`：这是订单领域定义。
- `"context"`：仓储方法统一接收上下文。
- `"fmt"`：自定义错误字符串需要用它。
- `type Repository interface {`：定义订单仓储抽象。
- `Create(context.Context, *Order) (*Order, error)`：创建订单。
- `Get(ctx context.Context, id, customerID string) (*Order, error)`：按订单 ID 和客户 ID 查询订单。
- `Update(ctx context.Context, o *Order, updateFn func(context.Context, *Order) (*Order, error)) error`：按更新函数更新订单。
- `type NotFoundError struct { OrderID string }`：定义领域级未找到错误。
- `func (e NotFoundError) Error() string`：实现 `error` 接口。
- `fmt.Sprintf("order '%s' not found", e.OrderID)`：错误文本。

为什么这么写：
- 应用层只依赖仓储接口，具体是内存、MySQL、MongoDB 都可以在外层替换。
- 自定义 `NotFoundError` 能让上层以后根据错误类型决定返回 404 还是别的状态码。

### 4.11 `internal/order/adapters/order_inmem_repository.go`

先看代码全文：

```go
package adapters

import (
	"context"
	"strconv"
	"sync"
	"time"

	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
	"github.com/sirupsen/logrus"
)

type MemoryOrderRepository struct {
	lock  *sync.RWMutex
	store []*domain.Order
}

func NewMemoryOrderRepository() *MemoryOrderRepository {
	s := make([]*domain.Order, 0)
	s = append(s, &domain.Order{
		ID:          "fake-ID",
		CustomerID:  "fake-customer-id",
		Status:      "fake-status",
		PaymentLink: "fake-payment-link",
		Items:       nil,
	})
	return &MemoryOrderRepository{
		lock:  &sync.RWMutex{},
		store: s,
	}
}

func (m *MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	newOrder := &domain.Order{
		ID:          strconv.FormatInt(time.Now().Unix(), 10),
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
		Items:       order.Items,
	}
	m.store = append(m.store, newOrder)
	logrus.WithFields(logrus.Fields{
		"input_order":        order,
		"store_after_create": m.store,
	}).Info("memory_order_repo_create")
	return newOrder, nil
}

func (m *MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*domain.Order, error) {
	for i, v := range m.store {
		logrus.Infof("m.store[%d] = %+v", i, v)
	}
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, o := range m.store {
		if o.ID == id && o.CustomerID == customerID {
			logrus.Infof("memory_order_repo_get||found||id=%s||customerID=%s||res=%+v", id, customerID, *o)
			return o, nil
		}
	}
	return nil, domain.NotFoundError{OrderID: id}
}

func (m *MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	found := false
	for i, o := range m.store {
		if o.ID == order.ID && o.CustomerID == order.CustomerID {
			found = true
			updatedOrder, err := updateFn(ctx, o)
			if err != nil {
				return err
			}
			m.store[i] = updatedOrder
		}
	}
	if !found {
		return domain.NotFoundError{OrderID: order.ID}
	}
	return nil
}
```

逐行解释：
- `package adapters`：具体仓储实现放这里。
- `"context"`：更新函数会用到上下文。
- `"strconv"` 和 `"time"`：当前用 Unix 秒生成订单 ID。
- `"sync"`：内存仓储要靠读写锁保护并发访问。
- `domain ".../order/domain/order"`：真正存的是领域订单对象。
- `"github.com/sirupsen/logrus"`：当前用它做大量调试日志。
- `type MemoryOrderRepository struct { lock *sync.RWMutex; store []*domain.Order }`：仓储内部状态是“一个锁 + 一个切片”。
- `func NewMemoryOrderRepository()`：构造内存仓储。
- `s := make([]*domain.Order, 0)`：初始化切片。
- `s = append(s, &domain.Order{ ... fake data ... })`：预置一条假订单。目的不是业务需要，而是为了让查询链一上来就有数据可查。
- `lock: &sync.RWMutex{}`：读写锁必须是共享的，所以这里用指针持有。
- `store: s`：把初始数据塞进去。
- `func (m *MemoryOrderRepository) Create(...)`：这里改成了指针接收者。这个变化很重要，因为 `Create` 会修改 `store`，如果用值接收者，复制语义会更容易埋坑。
- `m.lock.Lock()` / `defer m.lock.Unlock()`：写操作上写锁。
- `newOrder := &domain.Order{ ... }`：创建新订单对象。
- `ID: strconv.FormatInt(time.Now().Unix(), 10)`：用秒级时间戳当 ID，教学里够用，但并发时同一秒内可能重复，不适合生产。
- `CustomerID / Status / PaymentLink / Items`：把传入订单数据拷贝到新对象。
- `m.store = append(m.store, newOrder)`：追加到切片。
- `logrus.WithFields(...)`：打结构化日志。
- `Info("memory_order_repo_create")`：这次从 Debug 升成了 Info，更容易在默认日志级别下看到。
- `func (m *MemoryOrderRepository) Get(...)`：查询方法也改成了指针接收者。
- `for i, v := range m.store { logrus.Infof(...) }`：这里先打印整个 store。它明显是为了教学/调试，但也有个问题：这段打印发生在加读锁之前，严格来说并不优雅。
- `m.lock.RLock()` / `defer m.lock.RUnlock()`：读操作上读锁。
- `for _, o := range m.store`：遍历内存切片查找。
- `if o.ID == id && o.CustomerID == customerID`：同时按订单 ID 和客户 ID 匹配。
- `return o, nil`：找到就返回。
- `return nil, domain.NotFoundError{OrderID: id}`：找不到就返回领域未找到错误。
- `func (m *MemoryOrderRepository) Update(...)`：更新操作。
- `m.lock.Lock()`：更新要拿写锁。
- `found := false`：先记是否找到目标。
- `for i, o := range m.store`：遍历切片找订单。
- `if o.ID == order.ID && o.CustomerID == order.CustomerID`：匹配同一条订单。
- `updatedOrder, err := updateFn(ctx, o)`：真正怎么更新交给上层传入的函数。
- `m.store[i] = updatedOrder`：把更新后的对象放回去。
- `if !found { return domain.NotFoundError{OrderID: order.ID} }`：最后统一处理没找到的情况。

为什么这么写：
- 这一层的重点是“让应用层不关心数据存在哪”。你以后把它换成 MySQL 仓储，`CreateOrder` 和 `GetCustomerOrder` 基本都不用改。

容易忘和容易错的点：
- 秒级时间戳做 ID 会撞号。
- `Get` 里在加锁前打印 `m.store`，并发场景下不是好习惯。
- 预置假数据是为了教学，不是生产逻辑。

### 4.12 `internal/stock/adapters/stock_inmem_repository.go`

先看代码全文：

```go
package adapters

import (
	"context"
	"sync"

	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
)

type MemoryStockRepository struct {
	lock  *sync.RWMutex
	store map[string]*orderpb.Item
}

var stub = map[string]*orderpb.Item{
	"item_id": {
		ID:       "foo_item",
		Name:     "stub item",
		Quantity: 10000,
		PriceID:  "stub_item_price_id",
	},
}

func NewMemoryStockRepository() *MemoryStockRepository {
	return &MemoryStockRepository{
		lock:  &sync.RWMutex{},
		store: stub,
	}
}

func (m MemoryStockRepository) GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var (
		res     []*orderpb.Item
		missing []string
	)
	for _, id := range ids {
		if item, exist := m.store[id]; exist {
			res = append(res, item)
		} else {
			missing = append(missing, id)
		}
	}
	if len(res) == len(ids) {
		return res, nil
	}
	return res, domain.NotFoundError{Missing: missing}
}
```

逐行解释：
- `package adapters`：库存服务的适配器实现。
- `"context"`：接口统一带上下文。
- `"sync"`：同样用锁保护共享内存数据。
- `".../orderpb"`：库存里存的商品数据也直接用 protobuf 生成结构。
- `domain ".../stock/domain/stock"`：找不到商品时返回库存领域错误。
- `type MemoryStockRepository struct { lock *sync.RWMutex; store map[string]*orderpb.Item }`：库存内存仓储用 map，比切片更适合按商品 ID 查。
- `var stub = map[string]*orderpb.Item{ ... }`：预置一条假库存商品。
- `func NewMemoryStockRepository() *MemoryStockRepository`：这一节只改了这里的名字，之前叫 `NewMemoryOrderRepository` 明显是笔误，现在修正过来了。
- `store: stub`：初始库存数据来自预置 map。
- `func (m MemoryStockRepository) GetItems(...)`：按商品 ID 批量取数据。
- `m.lock.RLock()` / `defer m.lock.RUnlock()`：读锁。
- `var ( res []*orderpb.Item; missing []string )`：分别收集找到的商品和缺失的 ID。
- `for _, id := range ids`：遍历请求的商品 ID。
- `if item, exist := m.store[id]; exist { ... } else { ... }`：找到就收集商品，找不到就记录缺失 ID。
- `if len(res) == len(ids) { return res, nil }`：全部找到才算完全成功。
- `return res, domain.NotFoundError{Missing: missing}`：部分缺失时返回“已找到的商品 + 错误”。这是一种很常见的批量查询设计。

为什么这么写：
- 库存天生是按商品 ID 查，所以用 map 很顺手。
- 批量查询允许“部分成功”，便于调用方自己决定是否继续处理。

### 4.13 `internal/common/decorator/query.go`

先看代码全文：

```go
package decorator

import (
	"context"

	"github.com/sirupsen/logrus"
)

// QueryHandler defines a generic type that receives a Query Q,
// and returns a result R
type QueryHandler[Q, R any] interface {
	Handle(ctx context.Context, query Q) (R, error)
}

func ApplyQueryDecorators[H, R any](handler QueryHandler[H, R], logger *logrus.Entry, metricsClient MetricsClient) QueryHandler[H, R] {
	return queryLoggingDecorator[H, R]{
		logger: logger,
		base: queryMetricsDecorator[H, R]{
			base:   handler,
			client: metricsClient,
		},
	}
}
```

逐行解释：
- `package decorator`：横切逻辑都放在这里。
- `"context"`：handler 统一吃上下文。
- `"github.com/sirupsen/logrus"`：装饰器里需要 logger。
- `type QueryHandler[Q, R any] interface {`：这是泛型接口。`Q` 是查询参数类型，`R` 是返回结果类型，`any` 就是“任意类型都行”。
- `Handle(ctx context.Context, query Q) (R, error)`：只要实现了这个签名，就能被当成 QueryHandler。
- 注释 `// QueryHandler defines ...`：这是作者给泛型接口补的人类说明。
- `func ApplyQueryDecorators[H, R any](...)`：给一个查询 handler 套上装饰器。
- `return queryLoggingDecorator{ ... queryMetricsDecorator{ ... handler ... } }`：执行顺序其实是“外层日志 -> 内层指标 -> 最里面业务 handler”。

为什么这么写：
- 你不想在每个 query handler 里手写一遍日志和指标，所以把它们抽成外层包装。

Go 小白容易卡的点：
- 泛型接口不是魔法，它只是把“同样的 handler 结构”抽象成可复用模板。
- `ApplyQueryDecorators` 返回的仍然是 `QueryHandler`，所以上层根本不需要知道它里面包了几层。

### 4.14 `internal/common/decorator/command.go`

先看代码全文：

```go
package decorator

import (
	"context"

	"github.com/sirupsen/logrus"
)

type CommandHandler[C, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

func ApplyCommandDecorators[C, R any](handler CommandHandler[C, R], logger *logrus.Entry, metricsClient MetricsClient) CommandHandler[C, R] {
	return queryLoggingDecorator[C, R]{
		logger: logger,
		base: queryMetricsDecorator[C, R]{
			base:   handler,
			client: metricsClient,
		},
	}
}
```

逐行解释：
- `type CommandHandler[C, R any] interface { Handle(ctx context.Context, cmd C) (R, error) }`：命令版本的泛型 handler 接口。
- `ApplyCommandDecorators(...)`：给命令套装饰器。
- 这里返回的是 `queryLoggingDecorator` 和 `queryMetricsDecorator`：代码能跑，但语义上不够严谨，因为命令和查询本来应该区分日志字段和指标名字。

为什么这么写：
- 很明显作者是先复用已有装饰器，把链路先跑起来。
- 从教学角度这是可以理解的，但你要知道“能跑”和“语义干净”是两回事。

### 4.15 `internal/common/decorator/logging.go`

先看代码全文：

```go
package decorator

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type queryLoggingDecorator[C, R any] struct {
	logger *logrus.Entry
	base   QueryHandler[C, R]
}

func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	logger := q.logger.WithFields(logrus.Fields{
		"query":      generateActionName(cmd),
		"query_body": fmt.Sprintf("%#v", cmd),
	})
	logger.Debug("Executing query")
	defer func() {
		if err == nil {
			logger.Info("Query execute successfully")
		} else {
			logger.Error("Failed to execute query", err)
		}
	}()
	return q.base.Handle(ctx, cmd)
}

func generateActionName(cmd any) string {
	return strings.Split(fmt.Sprintf("%T", cmd), ".")[1]
}
```

逐行解释：
- `package decorator`：日志装饰器属于横切层。
- `"fmt"`：用来把请求对象和类型名转成字符串。
- `"strings"`：用来切分类型名。
- `type queryLoggingDecorator[C, R any] struct { logger *logrus.Entry; base QueryHandler[C, R] }`：它自己不做业务，只持有 logger 和“下一层 handler”。
- `func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error)`：这里用了命名返回值，主要是为了让 `defer` 里能读取 `err` 最终值。
- `logger := q.logger.WithFields(logrus.Fields{ ... })`：先构造一条带上下文字段的日志对象。
- `"query": generateActionName(cmd)`：记录动作名。
- `"query_body": fmt.Sprintf("%#v", cmd)`：把命令/查询对象按 Go 语法打印出来，方便调试。
- `logger.Debug("Executing query")`：执行前打日志。
- `defer func() { ... }()`：等 handler 执行结束后再统一打成功或失败日志。
- `if err == nil { logger.Info(...) } else { logger.Error(...) }`：根据结果决定日志级别。
- `return q.base.Handle(ctx, cmd)`：真正调用下一层 handler。
- `func generateActionName(cmd any) string`：从类型名里推导动作名。
- `strings.Split(fmt.Sprintf("%T", cmd), ".")[1]`：比如 `query.GetCustomerOrder` 会拿到 `GetCustomerOrder`。但这写法比较脆，如果 `%T` 的结果里没有点号，就会越界。

为什么这么写：
- 统一日志比在每个 handler 手工打日志更整洁。

外部库和易错点：
- Logrus 更推荐 `WithError(err)` 这类结构化写法，而不是简单 `Error("...", err)`。
- `fmt.Sprintf("%#v", cmd)` 虽然方便，但对象很大时日志会非常吵。

### 4.16 `internal/common/decorator/metrics.go`

先看代码全文：

```go
package decorator

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type MetricsClient interface {
	Inc(key string, value int)
}

type queryMetricsDecorator[C, R any] struct {
	base   QueryHandler[C, R]
	client MetricsClient
}

func (q queryMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	start := time.Now()
	actionName := strings.ToLower(generateActionName(cmd))
	defer func() {
		end := time.Since(start)
		q.client.Inc(fmt.Sprintf("querys.%s.duration", actionName), int(end.Seconds()))
		if err == nil {
			q.client.Inc(fmt.Sprintf("querys.%s.success", actionName), 1)
		} else {
			q.client.Inc(fmt.Sprintf("querys.%s.failure", actionName), 1)
		}
	}()
	return q.base.Handle(ctx, cmd)
}
```

逐行解释：
- `package decorator`：指标装饰器也在横切层。
- `"time"`：用来算执行耗时。
- `type MetricsClient interface { Inc(key string, value int) }`：这里把指标系统抽象成一个很简单的接口。当前只支持递增计数。
- `type queryMetricsDecorator[C, R any] struct { base QueryHandler[C, R]; client MetricsClient }`：持有下一层 handler 和指标客户端。
- `start := time.Now()`：记录开始时间。
- `actionName := strings.ToLower(generateActionName(cmd))`：把动作名转成小写，方便拼指标 key。
- `defer func() { ... }()`：执行完成后统一记指标。
- `end := time.Since(start)`：算出耗时。
- `q.client.Inc(fmt.Sprintf("querys.%s.duration", actionName), int(end.Seconds()))`：把耗时作为整数秒记到指标里。这里有两个明显问题：一是 `querys` 拼写不标准，二是秒级整数会丢失小于 1 秒的精度。
- `success` / `failure`：根据 `err` 决定记录成功还是失败计数。
- `return q.base.Handle(ctx, cmd)`：真正调用下一层。

为什么这么写：
- 横切指标和横切日志一样，应该尽量从业务代码里拿掉。

### 4.17 `internal/common/metrics/todo_metrics.go`

先看代码全文：

```go
package metrics

type TodoMetrics struct{}

func (t TodoMetrics) Inc(_ string, _ int) {
}
```

逐行解释：
- `package metrics`：指标实现包。
- `type TodoMetrics struct{}`：空结构体，占位用。
- `func (t TodoMetrics) Inc(_ string, _ int) {}`：空实现，表示“接口先接上，真正的指标系统以后再补”。

为什么这么写：
- 教学项目里很常见，先把依赖形状定下来，后面再接真实监控系统。

## 5. 这节代码里最值得你特别注意的坑

下面这些不是我挑刺，而是你以后自己写 Go 服务时非常容易犯的错：

1. `PostCustomerCustomerIDOrders` 里忽略了路径参数 `customerID`，却信任了 body 里的 `req.CustomerID`。
   这会带来“路径上是 A，body 里是 B，到底以谁为准”的问题。

2. `create_order.go` 里第一次 RPC 的 `err` 被第二次 RPC 覆盖了。
   这是典型的 Go 新手错误。每次拿到 `err` 后都要立刻处理，不要拖。

3. `create_order.go` 里 `GetItems` 用了硬编码 `[]string{"123"}`。
   这说明代码还在调链路，不是最终版逻辑。你读代码时要能识别这种“过渡态实现”。

4. `create_order.go` 里拿到了 `resp` 却没真正用它来构造订单商品。
   也就是说，库存服务已经接进来了，但数据流还没有真正闭环。

5. `http.go` 里业务错误返回了 `200 OK`。
   这在演示项目里常见，但在真实 API 设计里通常不推荐。

6. `http.go` 里 `gin.H{"error": err}` 直接把 `error` 扔给 JSON。
   最稳妥的写法通常是 `err.Error()`，不然前端拿到的格式可能很奇怪。

7. `order_inmem_repository.go` 里查询前打印整个 `store`，而且打印发生在 `RLock` 之前。
   教学阶段方便调试，但并发语义不够严谨。

8. `metrics.go` 里指标 key 用的是 `querys`，还有耗时被强转成整数秒。
   这两个点都说明装饰器现在还是第一版。

9. `command.go` 直接复用了 query 的 logging/metrics decorator。
   功能上可用，但语义上命令和查询最好分开。

10. `generateActionName` 依赖 `strings.Split(..., ".")[1]`。
    只要类型名格式跟预期不一样，就可能越界。这种“靠字符串解析类型名”的技巧要谨慎用。

## 6. 这一节为什么要这么设计

如果你问我：为什么不直接在 `http.go` 里查仓储、调库存服务、写日志、记指标，一把梭？

答案是：因为一开始这么写最快，但后面会非常难维护。

这一节的设计意图其实很清楚：
- 让 `main.go` 只负责启动。
- 让 `application.go` 只负责接线。
- 让 `http.go` 只负责协议翻译。
- 让 `command/query` 只负责业务用例。
- 让 `adapters` 只负责和外部世界交互。
- 让 `decorator` 只负责日志和指标这种横切能力。

这样分开以后：
- 你想加 gRPC 入口，不用重写业务逻辑。
- 你想把内存仓储换成 MySQL，不用重写 HTTP 层。
- 你想把日志/指标增强，不用污染业务 handler。
- 你想给 `CreateOrder` 写单元测试，也可以只 mock `Repository` 和 `StockService`。

这就是为什么我一直强调：这一节真正学的不是“怎么多写几个文件”，而是“怎么把一个服务拆成可维护的层次”。

## 7. 你现在最适合怎么复读

如果你是 Go 小白，我建议你这样回看：

1. 先把 [main.go](/g:/shi/go_shop_second/internal/order/main.go) 和 [application.go](/g:/shi/go_shop_second/internal/order/service/application.go) 连起来看，先搞懂程序怎么启动、依赖怎么装配。
2. 再看 [http.go](/g:/shi/go_shop_second/internal/order/http.go)，理解 HTTP 请求是怎么被翻译成命令/查询的。
3. 然后重点看 [create_order.go](/g:/shi/go_shop_second/internal/order/app/command/create_order.go) 和 [get_customer_order.go](/g:/shi/go_shop_second/internal/order/app/query/get_customer_order.go)，这是这节最核心的业务入口。
4. 最后回头看 [query.go](/g:/shi/go_shop_second/internal/common/decorator/query.go)、[command.go](/g:/shi/go_shop_second/internal/common/decorator/command.go)、[logging.go](/g:/shi/go_shop_second/internal/common/decorator/logging.go)、[metrics.go](/g:/shi/go_shop_second/internal/common/decorator/metrics.go)，你就能看懂横切逻辑为什么单独拿出去。

## 8. 这一节的最终一句话

`lesson7 -> lesson8` 的本质，不是“多了创建和查询订单接口”，而是项目开始形成一个真正像样的服务骨架：入口层、应用层、仓储、外部服务适配器、日志、指标，第一次被串成了一条完整调用链。

这份讲义是手写版，不是模板版。下一轮如果你继续，我会保持同样标准，只讲一组差异，不机械生成。

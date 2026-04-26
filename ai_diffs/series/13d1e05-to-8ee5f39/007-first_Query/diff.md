# Commit Diff Report

- Repo: go_shop_second
- Sequence: 007 / 10
- Commit: 7881a0bb7c726a6ae8ee8f128ba06f421622badb
- ShortCommit: 7881a0b
- Parent: f87d59efe0cc19ac7212056bbcbd3f73c733340d
- Subject: first Query
- Author: ghost-yu <hgfhgfhgfhgfhgfhgf@yeah.net>
- Date: 2024-10-14 21:07:45 +0800
- GeneratedAt: 2026-04-06 17:43:36 +08:00

## Short Summary

~~~text
 10 files changed, 178 insertions(+), 6 deletions(-)
~~~

## File Stats

~~~text
 internal/common/decorator/logging.go              | 34 ++++++++++++++++++
 internal/common/decorator/metrics.go              | 32 +++++++++++++++++
 internal/common/decorator/query.go                | 23 ++++++++++++
 internal/common/metrics/todo_metrics.go           |  6 ++++
 internal/order/adapters/order_inmem_repository.go | 10 +++++-
 internal/order/app/app.go                         |  6 +++-
 internal/order/app/query/get_customer_order.go    | 43 +++++++++++++++++++++++
 internal/order/http.go                            | 14 ++++++--
 internal/order/service/application.go             | 14 +++++++-
 internal/stock/adapters/stock_inmem_repository.go |  2 +-
 10 files changed, 178 insertions(+), 6 deletions(-)
~~~

## Changed Files

~~~text
internal/common/decorator/logging.go
internal/common/decorator/metrics.go
internal/common/decorator/query.go
internal/common/metrics/todo_metrics.go
internal/order/adapters/order_inmem_repository.go
internal/order/app/app.go
internal/order/app/query/get_customer_order.go
internal/order/http.go
internal/order/service/application.go
internal/stock/adapters/stock_inmem_repository.go
~~~

## Focus Files (Excluded: go.mod / go.sum)

~~~text
internal/common/decorator/logging.go
internal/common/decorator/metrics.go
internal/common/decorator/query.go
internal/common/metrics/todo_metrics.go
internal/order/adapters/order_inmem_repository.go
internal/order/app/app.go
internal/order/app/query/get_customer_order.go
internal/order/http.go
internal/order/service/application.go
internal/stock/adapters/stock_inmem_repository.go
~~~

## Patch

~~~diff
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
diff --git a/internal/order/adapters/order_inmem_repository.go b/internal/order/adapters/order_inmem_repository.go
index f5d8be2..667091b 100644
--- a/internal/order/adapters/order_inmem_repository.go
+++ b/internal/order/adapters/order_inmem_repository.go
@@ -16,9 +16,17 @@ type MemoryOrderRepository struct {
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
 
diff --git a/internal/order/app/app.go b/internal/order/app/app.go
index 42330cd..0da5849 100644
--- a/internal/order/app/app.go
+++ b/internal/order/app/app.go
@@ -1,5 +1,7 @@
 package app
 
+import "github.com/ghost-yu/go_shop_second/order/app/query"
+
 type Application struct {
 	Commands Commands
 	Queries  Queries
@@ -7,4 +9,6 @@ type Application struct {
 
 type Commands struct{}
 
-type Queries struct{}
+type Queries struct {
+	GetCustomerOrder query.GetCustomerOrderHandler
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
diff --git a/internal/order/http.go b/internal/order/http.go
index f18fa80..654c9d9 100644
--- a/internal/order/http.go
+++ b/internal/order/http.go
@@ -1,7 +1,10 @@
 package main
 
 import (
+	"net/http"
+
 	"github.com/ghost-yu/go_shop_second/order/app"
+	"github.com/ghost-yu/go_shop_second/order/app/query"
 	"github.com/gin-gonic/gin"
 )
 
@@ -15,6 +18,13 @@ func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID stri
 }
 
 func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
-	//TODO implement me
-	panic("implement me")
+	o, err := H.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
+		OrderID:    "fake-ID",
+		CustomerID: "fake-customer-id",
+	})
+	if err != nil {
+		c.JSON(http.StatusOK, gin.H{"error": err})
+		return
+	}
+	c.JSON(http.StatusOK, gin.H{"message": "success", "data": o})
 }
diff --git a/internal/order/service/application.go b/internal/order/service/application.go
index 122b22a..ece87d4 100644
--- a/internal/order/service/application.go
+++ b/internal/order/service/application.go
@@ -3,9 +3,21 @@ package service
 import (
 	"context"
 
+	"github.com/ghost-yu/go_shop_second/common/metrics"
+	"github.com/ghost-yu/go_shop_second/order/adapters"
 	"github.com/ghost-yu/go_shop_second/order/app"
+	"github.com/ghost-yu/go_shop_second/order/app/query"
+	"github.com/sirupsen/logrus"
 )
 
 func NewApplication(ctx context.Context) app.Application {
-	return app.Application{}
+	orderRepo := adapters.NewMemoryOrderRepository()
+	logger := logrus.NewEntry(logrus.StandardLogger())
+	metricClient := metrics.TodoMetrics{}
+	return app.Application{
+		Commands: app.Commands{},
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
~~~

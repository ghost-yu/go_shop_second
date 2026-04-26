# Commit Diff Report

- Repo: go_shop_second
- Sequence: 008 / 10
- Commit: 606345b1bffa9e6e0f221a73bc2671359315a77a
- ShortCommit: 606345b
- Parent: 7881a0bb7c726a6ae8ee8f128ba06f421622badb
- Subject: order create and update
- Author: ghost-yu <hgfhgfhgfhgfhgfhgf@yeah.net>
- Date: 2024-10-14 22:28:53 +0800
- GeneratedAt: 2026-04-06 17:43:36 +08:00

## Short Summary

~~~text
 7 files changed, 171 insertions(+), 12 deletions(-)
~~~

## File Stats

~~~text
 internal/common/decorator/command.go              | 21 ++++++++
 internal/order/adapters/order_inmem_repository.go | 13 +++--
 internal/order/app/app.go                         | 10 +++-
 internal/order/app/command/create_order.go        | 60 +++++++++++++++++++++++
 internal/order/app/command/update_order.go        | 47 ++++++++++++++++++
 internal/order/http.go                            | 26 ++++++++--
 internal/order/service/application.go             |  6 ++-
 7 files changed, 171 insertions(+), 12 deletions(-)
~~~

## Changed Files

~~~text
internal/common/decorator/command.go
internal/order/adapters/order_inmem_repository.go
internal/order/app/app.go
internal/order/app/command/create_order.go
internal/order/app/command/update_order.go
internal/order/http.go
internal/order/service/application.go
~~~

## Focus Files (Excluded: go.mod / go.sum)

~~~text
internal/common/decorator/command.go
internal/order/adapters/order_inmem_repository.go
internal/order/app/app.go
internal/order/app/command/create_order.go
internal/order/app/command/update_order.go
internal/order/http.go
internal/order/service/application.go
~~~

## Patch

~~~diff
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
diff --git a/internal/order/adapters/order_inmem_repository.go b/internal/order/adapters/order_inmem_repository.go
index 667091b..818b51f 100644
--- a/internal/order/adapters/order_inmem_repository.go
+++ b/internal/order/adapters/order_inmem_repository.go
@@ -30,7 +30,7 @@ func NewMemoryOrderRepository() *MemoryOrderRepository {
 	}
 }
 
-func (m MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
+func (m *MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
 	m.lock.Lock()
 	defer m.lock.Unlock()
 	newOrder := &domain.Order{
@@ -44,23 +44,26 @@ func (m MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*
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
index 0da5849..b2e04ce 100644
--- a/internal/order/app/app.go
+++ b/internal/order/app/app.go
@@ -1,13 +1,19 @@
 package app
 
-import "github.com/ghost-yu/go_shop_second/order/app/query"
+import (
+	"github.com/ghost-yu/go_shop_second/order/app/command"
+	"github.com/ghost-yu/go_shop_second/order/app/query"
+)
 
 type Application struct {
 	Commands Commands
 	Queries  Queries
 }
 
-type Commands struct{}
+type Commands struct {
+	CreateOrder command.CreateOrderHandler
+	UpdateOrder command.UpdateOrderHandler
+}
 
 type Queries struct {
 	GetCustomerOrder query.GetCustomerOrderHandler
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
new file mode 100644
index 0000000..17620e2
--- /dev/null
+++ b/internal/order/app/command/create_order.go
@@ -0,0 +1,60 @@
+package command
+
+import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/common/decorator"
+	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
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
+	//stockGRPC
+}
+
+func NewCreateOrderHandler(
+	orderRepo domain.Repository,
+	logger *logrus.Entry,
+	metricClient decorator.MetricsClient,
+) CreateOrderHandler {
+	if orderRepo == nil {
+		panic("nil orderRepo")
+	}
+	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
+		createOrderHandler{orderRepo: orderRepo},
+		logger,
+		metricClient,
+	)
+}
+
+func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
+	// TODO: call stock grpc to get items.
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
diff --git a/internal/order/http.go b/internal/order/http.go
index 654c9d9..b40adc7 100644
--- a/internal/order/http.go
+++ b/internal/order/http.go
@@ -3,7 +3,9 @@ package main
 import (
 	"net/http"
 
+	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	"github.com/ghost-yu/go_shop_second/order/app"
+	"github.com/ghost-yu/go_shop_second/order/app/command"
 	"github.com/ghost-yu/go_shop_second/order/app/query"
 	"github.com/gin-gonic/gin"
 )
@@ -13,14 +15,30 @@ type HTTPServer struct {
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
 	o, err := H.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
-		OrderID:    "fake-ID",
-		CustomerID: "fake-customer-id",
+		OrderID:    orderID,
+		CustomerID: customerID,
 	})
 	if err != nil {
 		c.JSON(http.StatusOK, gin.H{"error": err})
diff --git a/internal/order/service/application.go b/internal/order/service/application.go
index ece87d4..d8fd936 100644
--- a/internal/order/service/application.go
+++ b/internal/order/service/application.go
@@ -6,6 +6,7 @@ import (
 	"github.com/ghost-yu/go_shop_second/common/metrics"
 	"github.com/ghost-yu/go_shop_second/order/adapters"
 	"github.com/ghost-yu/go_shop_second/order/app"
+	"github.com/ghost-yu/go_shop_second/order/app/command"
 	"github.com/ghost-yu/go_shop_second/order/app/query"
 	"github.com/sirupsen/logrus"
 )
@@ -15,7 +16,10 @@ func NewApplication(ctx context.Context) app.Application {
 	logger := logrus.NewEntry(logrus.StandardLogger())
 	metricClient := metrics.TodoMetrics{}
 	return app.Application{
-		Commands: app.Commands{},
+		Commands: app.Commands{
+			CreateOrder: command.NewCreateOrderHandler(orderRepo, logger, metricClient),
+			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
+		},
 		Queries: app.Queries{
 			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
 		},
~~~

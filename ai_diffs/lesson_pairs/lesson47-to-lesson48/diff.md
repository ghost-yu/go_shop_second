# Lesson Pair Diff Report

- FromBranch: lesson47
- ToBranch: lesson48

## Short Summary

~~~text
 5 files changed, 24 insertions(+), 14 deletions(-)
~~~

## File Stats

~~~text
 internal/order/app/command/create_order.go          |  7 ++++---
 internal/payment/adapters/order_grpc.go             | 14 ++++++++++----
 internal/stock/adapters/stock_mysql_repository.go   |  7 ++++---
 internal/stock/app/query/check_if_items_in_stock.go |  4 ++--
 internal/stock/ports/grpc.go                        |  6 ++++--
 5 files changed, 24 insertions(+), 14 deletions(-)
~~~

## Commit Comparison

~~~text
> ccc9b76 wrap error
~~~

## Changed Files

~~~text
internal/order/app/command/create_order.go
internal/payment/adapters/order_grpc.go
internal/stock/adapters/stock_mysql_repository.go
internal/stock/app/query/check_if_items_in_stock.go
internal/stock/ports/grpc.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
internal/order/app/command/create_order.go
internal/payment/adapters/order_grpc.go
internal/stock/adapters/stock_mysql_repository.go
internal/stock/app/query/check_if_items_in_stock.go
internal/stock/ports/grpc.go
~~~

## Full Diff

~~~diff
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
index 6384a96..998302b 100644
--- a/internal/order/app/command/create_order.go
+++ b/internal/order/app/command/create_order.go
@@ -3,7 +3,6 @@ package command
 import (
 	"context"
 	"encoding/json"
-	"errors"
 	"fmt"
 
 	"github.com/ghost-yu/go_shop_second/common/broker"
@@ -12,9 +11,11 @@ import (
 	"github.com/ghost-yu/go_shop_second/order/convertor"
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
 	"github.com/ghost-yu/go_shop_second/order/entity"
+	"github.com/pkg/errors"
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
 	"go.opentelemetry.io/otel"
+	"google.golang.org/grpc/status"
 )
 
 type CreateOrder struct {
@@ -96,7 +97,7 @@ func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*Creat
 		Headers:      header,
 	})
 	if err != nil {
-		return nil, err
+		return nil, errors.Wrapf(err, "publish event error q.Name=%s", q.Name)
 	}
 
 	return &CreateOrderResult{OrderID: o.ID}, nil
@@ -109,7 +110,7 @@ func (c createOrderHandler) validate(ctx context.Context, items []*entity.ItemWi
 	items = packItems(items)
 	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, convertor.NewItemWithQuantityConvertor().EntitiesToProtos(items))
 	if err != nil {
-		return nil, err
+		return nil, status.Convert(err).Err()
 	}
 	return convertor.NewItemConvertor().ProtosToEntities(resp.Items), nil
 }
diff --git a/internal/payment/adapters/order_grpc.go b/internal/payment/adapters/order_grpc.go
index 3bd171d..6890559 100644
--- a/internal/payment/adapters/order_grpc.go
+++ b/internal/payment/adapters/order_grpc.go
@@ -6,6 +6,7 @@ import (
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	"github.com/ghost-yu/go_shop_second/common/tracing"
 	"github.com/sirupsen/logrus"
+	"google.golang.org/grpc/status"
 )
 
 type OrderGRPC struct {
@@ -16,11 +17,16 @@ func NewOrderGRPC(client orderpb.OrderServiceClient) *OrderGRPC {
 	return &OrderGRPC{client: client}
 }
 
-func (o OrderGRPC) UpdateOrder(ctx context.Context, order *orderpb.Order) error {
+func (o OrderGRPC) UpdateOrder(ctx context.Context, order *orderpb.Order) (err error) {
+	defer func() {
+		if err != nil {
+			logrus.Infof("payment_adapter||update_order,err=%v", err)
+		}
+	}()
+
 	ctx, span := tracing.Start(ctx, "order_grpc.update_order")
 	defer span.End()
 
-	_, err := o.client.UpdateOrder(ctx, order)
-	logrus.Infof("payment_adapter||update_order,err=%v", err)
-	return err
+	_, err = o.client.UpdateOrder(ctx, order)
+	return status.Convert(err).Err()
 }
diff --git a/internal/stock/adapters/stock_mysql_repository.go b/internal/stock/adapters/stock_mysql_repository.go
index 7afb9fa..75175e1 100644
--- a/internal/stock/adapters/stock_mysql_repository.go
+++ b/internal/stock/adapters/stock_mysql_repository.go
@@ -5,6 +5,7 @@ import (
 
 	"github.com/ghost-yu/go_shop_second/stock/entity"
 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent"
+	"github.com/pkg/errors"
 	"github.com/sirupsen/logrus"
 	"gorm.io/gorm"
 )
@@ -25,7 +26,7 @@ func (m MySQLStockRepository) GetItems(ctx context.Context, ids []string) ([]*en
 func (m MySQLStockRepository) GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error) {
 	data, err := m.db.BatchGetStockByID(ctx, ids)
 	if err != nil {
-		return nil, err
+		return nil, errors.Wrap(err, "BatchGetStockByID error")
 	}
 	var result []*entity.ItemWithQuantity
 	for _, d := range data {
@@ -54,7 +55,7 @@ func (m MySQLStockRepository) UpdateStock(
 		}()
 		var dest []*persistent.StockModel
 		if err = tx.Table("o_stock").Where("product_id IN ?", getIDFromEntities(data)).Find(&dest).Error; err != nil {
-			return err
+			return errors.Wrap(err, "failed to find data")
 		}
 		existing := m.unmarshalFromDatabase(dest)
 
@@ -65,7 +66,7 @@ func (m MySQLStockRepository) UpdateStock(
 
 		for _, upd := range updated {
 			if err = tx.Table("o_stock").Where("product_id = ?", upd.ID).Update("quantity", upd.Quantity).Error; err != nil {
-				return err
+				return errors.Wrapf(err, "unable to update %s", upd.ID)
 			}
 		}
 		return nil
diff --git a/internal/stock/app/query/check_if_items_in_stock.go b/internal/stock/app/query/check_if_items_in_stock.go
index 5c48a2b..2ddaa3a 100644
--- a/internal/stock/app/query/check_if_items_in_stock.go
+++ b/internal/stock/app/query/check_if_items_in_stock.go
@@ -10,6 +10,7 @@ import (
 	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
 	"github.com/ghost-yu/go_shop_second/stock/entity"
 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/integration"
+	"github.com/pkg/errors"
 	"github.com/sirupsen/logrus"
 )
 
@@ -58,7 +59,7 @@ var stub = map[string]string{
 
 func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*entity.Item, error) {
 	if err := lock(ctx, getLockKey(query)); err != nil {
-		return nil, err
+		return nil, errors.Wrapf(err, "redis lock error: key=%s", getLockKey(query))
 	}
 	defer func() {
 		if err := unlock(ctx, getLockKey(query)); err != nil {
@@ -78,7 +79,6 @@ func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfIte
 			PriceID:  priceID,
 		})
 	}
-	// TODO: 扣库存
 	if err := h.checkStock(ctx, query.Items); err != nil {
 		return nil, err
 	}
diff --git a/internal/stock/ports/grpc.go b/internal/stock/ports/grpc.go
index 51341e9..13f212e 100644
--- a/internal/stock/ports/grpc.go
+++ b/internal/stock/ports/grpc.go
@@ -8,6 +8,8 @@ import (
 	"github.com/ghost-yu/go_shop_second/stock/app"
 	"github.com/ghost-yu/go_shop_second/stock/app/query"
 	"github.com/ghost-yu/go_shop_second/stock/convertor"
+	"google.golang.org/grpc/codes"
+	"google.golang.org/grpc/status"
 )
 
 type GRPCServer struct {
@@ -24,7 +26,7 @@ func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsReque
 
 	items, err := G.app.Queries.GetItems.Handle(ctx, query.GetItems{ItemIDs: request.ItemIDs})
 	if err != nil {
-		return nil, err
+		return nil, status.Error(codes.Internal, err.Error())
 	}
 	return &stockpb.GetItemsResponse{Items: convertor.NewItemConvertor().EntitiesToProtos(items)}, nil
 }
@@ -37,7 +39,7 @@ func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.Ch
 		Items: convertor.NewItemWithQuantityConvertor().ProtosToEntities(request.Items),
 	})
 	if err != nil {
-		return nil, err
+		return nil, status.Error(codes.Internal, err.Error())
 	}
 	return &stockpb.CheckIfItemsInStockResponse{
 		InStock: 1,
~~~

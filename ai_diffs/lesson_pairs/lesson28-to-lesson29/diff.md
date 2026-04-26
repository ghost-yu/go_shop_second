# Lesson Pair Diff Report

- FromBranch: lesson28
- ToBranch: lesson29

## Short Summary

~~~text
 10 files changed, 279 insertions(+), 47 deletions(-)
~~~

## File Stats

~~~text
 api/openapi/order.yml                             |  14 ++
 internal/common/client/order/openapi_types.gen.go |  22 +--
 internal/order/app/command/create_order.go        |  17 +--
 internal/order/convertor/convertor.go             | 168 ++++++++++++++++++++++
 internal/order/convertor/facade.go                |  39 +++++
 internal/order/domain/order/order.go              |  16 +--
 internal/order/entity/entity.go                   |  13 ++
 internal/order/http.go                            |   3 +-
 internal/order/ports/grpc.go                      |  12 +-
 internal/order/ports/openapi_types.gen.go         |  22 +--
 10 files changed, 279 insertions(+), 47 deletions(-)
~~~

## Commit Comparison

~~~text
> b7bb8ed convertor
~~~

## Changed Files

~~~text
api/openapi/order.yml
internal/common/client/order/openapi_types.gen.go
internal/order/app/command/create_order.go
internal/order/convertor/convertor.go
internal/order/convertor/facade.go
internal/order/domain/order/order.go
internal/order/entity/entity.go
internal/order/http.go
internal/order/ports/grpc.go
internal/order/ports/openapi_types.gen.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
api/openapi/order.yml
internal/order/app/command/create_order.go
internal/order/convertor/convertor.go
internal/order/convertor/facade.go
internal/order/domain/order/order.go
internal/order/entity/entity.go
internal/order/http.go
internal/order/ports/grpc.go
~~~

## Full Diff

~~~diff
diff --git a/api/openapi/order.yml b/api/openapi/order.yml
index f421fce..c65f1cc 100644
--- a/api/openapi/order.yml
+++ b/api/openapi/order.yml
@@ -77,6 +77,12 @@ components:
   schemas:
     Order:
       type: object
+      required:
+        - id
+        - customerID
+        - status
+        - items
+        - paymentLink
       properties:
         id:
           type: string
@@ -93,6 +99,11 @@ components:
 
     Item:
       type: object
+      required:
+        - id
+        - name
+        - quantity
+        - priceID
       properties:
         id:
           type: string
@@ -125,6 +136,9 @@ components:
 
     ItemWithQuantity:
       type: object
+      required:
+        - id
+        - quantity
       properties:
         id:
           type: string
diff --git a/internal/common/client/order/openapi_types.gen.go b/internal/common/client/order/openapi_types.gen.go
index 01e8118..8021b88 100644
--- a/internal/common/client/order/openapi_types.gen.go
+++ b/internal/common/client/order/openapi_types.gen.go
@@ -16,25 +16,25 @@ type Error struct {
 
 // Item defines model for Item.
 type Item struct {
-	Id       *string `json:"id,omitempty"`
-	Name     *string `json:"name,omitempty"`
-	PriceID  *string `json:"priceID,omitempty"`
-	Quantity *int32  `json:"quantity,omitempty"`
+	Id       string `json:"id"`
+	Name     string `json:"name"`
+	PriceID  string `json:"priceID"`
+	Quantity int32  `json:"quantity"`
 }
 
 // ItemWithQuantity defines model for ItemWithQuantity.
 type ItemWithQuantity struct {
-	Id       *string `json:"id,omitempty"`
-	Quantity *int32  `json:"quantity,omitempty"`
+	Id       string `json:"id"`
+	Quantity int32  `json:"quantity"`
 }
 
 // Order defines model for Order.
 type Order struct {
-	CustomerID  *string `json:"customerID,omitempty"`
-	Id          *string `json:"id,omitempty"`
-	Items       *[]Item `json:"items,omitempty"`
-	PaymentLink *string `json:"paymentLink,omitempty"`
-	Status      *string `json:"status,omitempty"`
+	CustomerID  string `json:"customerID"`
+	Id          string `json:"id"`
+	Items       []Item `json:"items"`
+	PaymentLink string `json:"paymentLink"`
+	Status      string `json:"status"`
 }
 
 // PostCustomerCustomerIDOrdersJSONRequestBody defines body for PostCustomerCustomerIDOrders for application/json ContentType.
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
index 282f66c..64f9070 100644
--- a/internal/order/app/command/create_order.go
+++ b/internal/order/app/command/create_order.go
@@ -8,9 +8,10 @@ import (
 
 	"github.com/ghost-yu/go_shop_second/common/broker"
 	"github.com/ghost-yu/go_shop_second/common/decorator"
-	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	"github.com/ghost-yu/go_shop_second/order/app/query"
+	"github.com/ghost-yu/go_shop_second/order/convertor"
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
+	"github.com/ghost-yu/go_shop_second/order/entity"
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
 	"go.opentelemetry.io/otel"
@@ -18,7 +19,7 @@ import (
 
 type CreateOrder struct {
 	CustomerID string
-	Items      []*orderpb.ItemWithQuantity
+	Items      []*entity.ItemWithQuantity
 }
 
 type CreateOrderResult struct {
@@ -100,26 +101,26 @@ func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*Creat
 	return &CreateOrderResult{OrderID: o.ID}, nil
 }
 
-func (c createOrderHandler) validate(ctx context.Context, items []*orderpb.ItemWithQuantity) ([]*orderpb.Item, error) {
+func (c createOrderHandler) validate(ctx context.Context, items []*entity.ItemWithQuantity) ([]*entity.Item, error) {
 	if len(items) == 0 {
 		return nil, errors.New("must have at least one item")
 	}
 	items = packItems(items)
-	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, items)
+	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, convertor.NewItemWithQuantityConvertor().EntitiesToProtos(items))
 	if err != nil {
 		return nil, err
 	}
-	return resp.Items, nil
+	return convertor.NewItemConvertor().ProtosToEntities(resp.Items), nil
 }
 
-func packItems(items []*orderpb.ItemWithQuantity) []*orderpb.ItemWithQuantity {
+func packItems(items []*entity.ItemWithQuantity) []*entity.ItemWithQuantity {
 	merged := make(map[string]int32)
 	for _, item := range items {
 		merged[item.ID] += item.Quantity
 	}
-	var res []*orderpb.ItemWithQuantity
+	var res []*entity.ItemWithQuantity
 	for id, quantity := range merged {
-		res = append(res, &orderpb.ItemWithQuantity{
+		res = append(res, &entity.ItemWithQuantity{
 			ID:       id,
 			Quantity: quantity,
 		})
diff --git a/internal/order/convertor/convertor.go b/internal/order/convertor/convertor.go
new file mode 100644
index 0000000..0217bdb
--- /dev/null
+++ b/internal/order/convertor/convertor.go
@@ -0,0 +1,168 @@
+package convertor
+
+import (
+	client "github.com/ghost-yu/go_shop_second/common/client/order"
+	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
+	"github.com/ghost-yu/go_shop_second/order/entity"
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
+func (c *ItemWithQuantityConvertor) ClientsToEntities(items []client.ItemWithQuantity) (res []*entity.ItemWithQuantity) {
+	for _, i := range items {
+		res = append(res, c.ClientToEntity(i))
+	}
+	return
+}
+
+func (c *ItemWithQuantityConvertor) ClientToEntity(i client.ItemWithQuantity) *entity.ItemWithQuantity {
+	return &entity.ItemWithQuantity{
+		ID:       i.Id,
+		Quantity: i.Quantity,
+	}
+}
+
+func (c *OrderConvertor) EntityToProto(o *domain.Order) *orderpb.Order {
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
+func (c *OrderConvertor) ProtoToEntity(o *orderpb.Order) *domain.Order {
+	c.check(o)
+	return &domain.Order{
+		ID:          o.ID,
+		CustomerID:  o.CustomerID,
+		Status:      o.Status,
+		PaymentLink: o.PaymentLink,
+		Items:       NewItemConvertor().ProtosToEntities(o.Items),
+	}
+}
+
+func (c *OrderConvertor) ClientToEntity(o *client.Order) *domain.Order {
+	c.check(o)
+	return &domain.Order{
+		ID:          o.Id,
+		CustomerID:  o.CustomerID,
+		Status:      o.Status,
+		PaymentLink: o.PaymentLink,
+		Items:       NewItemConvertor().ClientsToEntities(o.Items),
+	}
+}
+
+func (c *OrderConvertor) EntityToClient(o *domain.Order) *client.Order {
+	c.check(o)
+	return &client.Order{
+		Id:          o.ID,
+		CustomerID:  o.CustomerID,
+		Status:      o.Status,
+		PaymentLink: o.PaymentLink,
+		Items:       NewItemConvertor().EntitiesToClients(o.Items),
+	}
+}
+
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
+func (c *ItemConvertor) ClientsToEntities(items []client.Item) (res []*entity.Item) {
+	for _, i := range items {
+		res = append(res, c.ClientToEntity(i))
+	}
+	return
+}
+
+func (c *ItemConvertor) EntitiesToClients(items []*entity.Item) (res []client.Item) {
+	for _, i := range items {
+		res = append(res, c.EntityToClient(i))
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
+
+func (c *ItemConvertor) ClientToEntity(i client.Item) *entity.Item {
+	return &entity.Item{
+		ID:       i.Id,
+		Name:     i.Name,
+		Quantity: i.Quantity,
+		PriceID:  i.PriceID,
+	}
+}
+
+func (c *ItemConvertor) EntityToClient(i *entity.Item) client.Item {
+	return client.Item{
+		Id:       i.ID,
+		Name:     i.Name,
+		Quantity: i.Quantity,
+		PriceID:  i.PriceID,
+	}
+}
diff --git a/internal/order/convertor/facade.go b/internal/order/convertor/facade.go
new file mode 100644
index 0000000..63ad8d0
--- /dev/null
+++ b/internal/order/convertor/facade.go
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
diff --git a/internal/order/domain/order/order.go b/internal/order/domain/order/order.go
index d87e406..08e7e96 100644
--- a/internal/order/domain/order/order.go
+++ b/internal/order/domain/order/order.go
@@ -3,7 +3,7 @@ package order
 import (
 	"fmt"
 
-	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/order/entity"
 	"github.com/pkg/errors"
 	"github.com/stripe/stripe-go/v80"
 )
@@ -13,10 +13,10 @@ type Order struct {
 	CustomerID  string
 	Status      string
 	PaymentLink string
-	Items       []*orderpb.Item
+	Items       []*entity.Item
 }
 
-func NewOrder(id, customerID, status, paymentLink string, items []*orderpb.Item) (*Order, error) {
+func NewOrder(id, customerID, status, paymentLink string, items []*entity.Item) (*Order, error) {
 	if id == "" {
 		return nil, errors.New("empty id")
 	}
@@ -38,16 +38,6 @@ func NewOrder(id, customerID, status, paymentLink string, items []*orderpb.Item)
 	}, nil
 }
 
-func (o *Order) ToProto() *orderpb.Order {
-	return &orderpb.Order{
-		ID:          o.ID,
-		CustomerID:  o.CustomerID,
-		Status:      o.Status,
-		Items:       o.Items,
-		PaymentLink: o.PaymentLink,
-	}
-}
-
 func (o *Order) IsPaid() error {
 	if o.Status == string(stripe.CheckoutSessionPaymentStatusPaid) {
 		return nil
diff --git a/internal/order/entity/entity.go b/internal/order/entity/entity.go
new file mode 100644
index 0000000..51016cc
--- /dev/null
+++ b/internal/order/entity/entity.go
@@ -0,0 +1,13 @@
+package entity
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
diff --git a/internal/order/http.go b/internal/order/http.go
index 7eeee78..eebdb99 100644
--- a/internal/order/http.go
+++ b/internal/order/http.go
@@ -9,6 +9,7 @@ import (
 	"github.com/ghost-yu/go_shop_second/order/app"
 	"github.com/ghost-yu/go_shop_second/order/app/command"
 	"github.com/ghost-yu/go_shop_second/order/app/query"
+	"github.com/ghost-yu/go_shop_second/order/convertor"
 	"github.com/gin-gonic/gin"
 )
 
@@ -27,7 +28,7 @@ func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID stri
 	}
 	r, err := H.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
 		CustomerID: req.CustomerID,
-		Items:      req.Items,
+		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
 	})
 	if err != nil {
 		c.JSON(http.StatusOK, gin.H{"error": err})
diff --git a/internal/order/ports/grpc.go b/internal/order/ports/grpc.go
index 534ce55..b366016 100644
--- a/internal/order/ports/grpc.go
+++ b/internal/order/ports/grpc.go
@@ -7,6 +7,7 @@ import (
 	"github.com/ghost-yu/go_shop_second/order/app"
 	"github.com/ghost-yu/go_shop_second/order/app/command"
 	"github.com/ghost-yu/go_shop_second/order/app/query"
+	"github.com/ghost-yu/go_shop_second/order/convertor"
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
 	"github.com/golang/protobuf/ptypes/empty"
 	"github.com/sirupsen/logrus"
@@ -26,7 +27,7 @@ func NewGRPCServer(app app.Application) *GRPCServer {
 func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
 	_, err := G.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
 		CustomerID: request.CustomerID,
-		Items:      request.Items,
+		Items:      convertor.NewItemWithQuantityConvertor().ProtosToEntities(request.Items),
 	})
 	if err != nil {
 		return nil, status.Error(codes.Internal, err.Error())
@@ -42,12 +43,17 @@ func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderReque
 	if err != nil {
 		return nil, status.Error(codes.NotFound, err.Error())
 	}
-	return o.ToProto(), nil
+	return convertor.NewOrderConvertor().EntityToProto(o), nil
 }
 
 func (G GRPCServer) UpdateOrder(ctx context.Context, request *orderpb.Order) (_ *emptypb.Empty, err error) {
 	logrus.Infof("order_grpc||request_in||request=%+v", request)
-	order, err := domain.NewOrder(request.ID, request.CustomerID, request.Status, request.PaymentLink, request.Items)
+	order, err := domain.NewOrder(
+		request.ID,
+		request.CustomerID,
+		request.Status,
+		request.PaymentLink,
+		convertor.NewItemConvertor().ProtosToEntities(request.Items))
 	if err != nil {
 		err = status.Error(codes.Internal, err.Error())
 		return nil, err
diff --git a/internal/order/ports/openapi_types.gen.go b/internal/order/ports/openapi_types.gen.go
index b6b8b9c..0d69dbf 100644
--- a/internal/order/ports/openapi_types.gen.go
+++ b/internal/order/ports/openapi_types.gen.go
@@ -16,25 +16,25 @@ type Error struct {
 
 // Item defines model for Item.
 type Item struct {
-	Id       *string `json:"id,omitempty"`
-	Name     *string `json:"name,omitempty"`
-	PriceID  *string `json:"priceID,omitempty"`
-	Quantity *int32  `json:"quantity,omitempty"`
+	Id       string `json:"id"`
+	Name     string `json:"name"`
+	PriceID  string `json:"priceID"`
+	Quantity int32  `json:"quantity"`
 }
 
 // ItemWithQuantity defines model for ItemWithQuantity.
 type ItemWithQuantity struct {
-	Id       *string `json:"id,omitempty"`
-	Quantity *int32  `json:"quantity,omitempty"`
+	Id       string `json:"id"`
+	Quantity int32  `json:"quantity"`
 }
 
 // Order defines model for Order.
 type Order struct {
-	CustomerID  *string `json:"customerID,omitempty"`
-	Id          *string `json:"id,omitempty"`
-	Items       *[]Item `json:"items,omitempty"`
-	PaymentLink *string `json:"paymentLink,omitempty"`
-	Status      *string `json:"status,omitempty"`
+	CustomerID  string `json:"customerID"`
+	Id          string `json:"id"`
+	Items       []Item `json:"items"`
+	PaymentLink string `json:"paymentLink"`
+	Status      string `json:"status"`
 }
 
 // PostCustomerCustomerIDOrdersJSONRequestBody defines body for PostCustomerCustomerIDOrders for application/json ContentType.
~~~

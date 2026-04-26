# Lesson Pair Diff Report

- FromBranch: lesson31
- ToBranch: lesson32

## Short Summary

~~~text
 3 files changed, 18 insertions(+), 12 deletions(-)
~~~

## File Stats

~~~text
 internal/order/app/dto/order.go       |  7 +++++++
 internal/order/convertor/convertor.go |  8 ++++----
 internal/order/http.go                | 15 +++++++--------
 3 files changed, 18 insertions(+), 12 deletions(-)
~~~

## Commit Comparison

~~~text
> 6a37c1f dto:
~~~

## Changed Files

~~~text
internal/order/app/dto/order.go
internal/order/convertor/convertor.go
internal/order/http.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
internal/order/app/dto/order.go
internal/order/convertor/convertor.go
internal/order/http.go
~~~

## Full Diff

~~~diff
diff --git a/internal/order/app/dto/order.go b/internal/order/app/dto/order.go
new file mode 100644
index 0000000..704e536
--- /dev/null
+++ b/internal/order/app/dto/order.go
@@ -0,0 +1,7 @@
+package dto
+
+type CreateOrderResponse struct {
+	OrderID     string `json:"order_id"`
+	CustomerID  string `json:"customer_id"`
+	RedirectURL string `json:"redirect_url"`
+}
diff --git a/internal/order/convertor/convertor.go b/internal/order/convertor/convertor.go
index 0217bdb..dc39c0b 100644
--- a/internal/order/convertor/convertor.go
+++ b/internal/order/convertor/convertor.go
@@ -79,7 +79,7 @@ func (c *OrderConvertor) ClientToEntity(o *client.Order) *domain.Order {
 	c.check(o)
 	return &domain.Order{
 		ID:          o.Id,
-		CustomerID:  o.CustomerID,
+		CustomerID:  o.CustomerId,
 		Status:      o.Status,
 		PaymentLink: o.PaymentLink,
 		Items:       NewItemConvertor().ClientsToEntities(o.Items),
@@ -90,7 +90,7 @@ func (c *OrderConvertor) EntityToClient(o *domain.Order) *client.Order {
 	c.check(o)
 	return &client.Order{
 		Id:          o.ID,
-		CustomerID:  o.CustomerID,
+		CustomerId:  o.CustomerID,
 		Status:      o.Status,
 		PaymentLink: o.PaymentLink,
 		Items:       NewItemConvertor().EntitiesToClients(o.Items),
@@ -154,7 +154,7 @@ func (c *ItemConvertor) ClientToEntity(i client.Item) *entity.Item {
 		ID:       i.Id,
 		Name:     i.Name,
 		Quantity: i.Quantity,
-		PriceID:  i.PriceID,
+		PriceID:  i.PriceId,
 	}
 }
 
@@ -163,6 +163,6 @@ func (c *ItemConvertor) EntityToClient(i *entity.Item) client.Item {
 		Id:       i.ID,
 		Name:     i.Name,
 		Quantity: i.Quantity,
-		PriceID:  i.PriceID,
+		PriceId:  i.PriceID,
 	}
 }
diff --git a/internal/order/http.go b/internal/order/http.go
index f06c157..7d43877 100644
--- a/internal/order/http.go
+++ b/internal/order/http.go
@@ -7,6 +7,7 @@ import (
 	client "github.com/ghost-yu/go_shop_second/common/client/order"
 	"github.com/ghost-yu/go_shop_second/order/app"
 	"github.com/ghost-yu/go_shop_second/order/app/command"
+	"github.com/ghost-yu/go_shop_second/order/app/dto"
 	"github.com/ghost-yu/go_shop_second/order/app/query"
 	"github.com/ghost-yu/go_shop_second/order/convertor"
 	"github.com/gin-gonic/gin"
@@ -20,12 +21,8 @@ type HTTPServer struct {
 func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
 	var (
 		req  client.CreateOrderRequest
+		resp dto.CreateOrderResponse
 		err  error
-		resp struct {
-			CustomerID  string `json:"customer_id"`
-			OrderID     string `json:"order_id"`
-			RedirectURL string `json:"redirect_url"`
-		}
 	)
 	defer func() {
 		H.Response(c, err, &resp)
@@ -41,9 +38,11 @@ func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID stri
 	if err != nil {
 		return
 	}
-	resp.CustomerID = req.CustomerId
-	resp.RedirectURL = fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderID)
-	resp.OrderID = r.OrderID
+	resp = dto.CreateOrderResponse{
+		OrderID:     r.OrderID,
+		CustomerID:  req.CustomerId,
+		RedirectURL: fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderID),
+	}
 }
 
 func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
~~~

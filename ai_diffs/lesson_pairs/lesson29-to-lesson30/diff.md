# Lesson Pair Diff Report

- FromBranch: lesson29
- ToBranch: lesson30

## Short Summary

~~~text
 7 files changed, 171 insertions(+), 126 deletions(-)
~~~

## File Stats

~~~text
 api/openapi/order.yml                              | 26 +++---
 internal/common/client/order/openapi_client.gen.go | 98 +++++++++++-----------
 internal/common/client/order/openapi_types.gen.go  | 12 +--
 internal/common/response.go                        | 43 ++++++++++
 internal/order/http.go                             | 58 ++++++-------
 internal/order/ports/openapi_api.gen.go            | 48 +++++------
 internal/order/ports/openapi_types.gen.go          | 12 +--
 7 files changed, 171 insertions(+), 126 deletions(-)
~~~

## Commit Comparison

~~~text
> 1ac2212 response
~~~

## Changed Files

~~~text
api/openapi/order.yml
internal/common/client/order/openapi_client.gen.go
internal/common/client/order/openapi_types.gen.go
internal/common/response.go
internal/order/http.go
internal/order/ports/openapi_api.gen.go
internal/order/ports/openapi_types.gen.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
api/openapi/order.yml
internal/common/response.go
internal/order/http.go
~~~

## Full Diff

~~~diff
diff --git a/api/openapi/order.yml b/api/openapi/order.yml
index c65f1cc..da9fbf8 100644
--- a/api/openapi/order.yml
+++ b/api/openapi/order.yml
@@ -10,18 +10,18 @@ servers:
         default: 127.0.0.1
 
 paths:
-  /customer/{customerID}/orders/{orderID}:
+  /customer/{customer_id}/orders/{order_id}:
     get:
       description: "get order"
       parameters:
         - in: path
-          name: customerID
+          name: customer_id
           schema:
             type: string
           required: true
 
         - in: path
-          name: orderID
+          name: order_id
           schema:
             type: string
           required: true
@@ -41,12 +41,12 @@ paths:
               schema:
                 $ref: '#/components/schemas/Error'
 
-  /customer/{customerID}/orders:
+  /customer/{customer_id}/orders:
     post:
       description: "create order"
       parameters:
         - in: path
-          name: customerID
+          name: customer_id
           schema:
             type: string
           required: true
@@ -79,14 +79,14 @@ components:
       type: object
       required:
         - id
-        - customerID
+        - customer_id
         - status
         - items
-        - paymentLink
+        - payment_link
       properties:
         id:
           type: string
-        customerID:
+        customer_id:
           type: string
         status:
           type: string
@@ -94,7 +94,7 @@ components:
           type: array
           items:
             $ref: '#/components/schemas/Item'
-        paymentLink:
+        payment_link:
           type: string
 
     Item:
@@ -103,7 +103,7 @@ components:
         - id
         - name
         - quantity
-        - priceID
+        - price_id
       properties:
         id:
           type: string
@@ -112,7 +112,7 @@ components:
         quantity:
           type: integer
           format: int32
-        priceID:
+        price_id:
           type: string
 
     Error:
@@ -124,10 +124,10 @@ components:
     CreateOrderRequest:
       type: object
       required:
-        - customerID
+        - customer_id
         - items
       properties:
-        customerID:
+        customer_id:
           type: string
         items:
           type: array
diff --git a/internal/common/client/order/openapi_client.gen.go b/internal/common/client/order/openapi_client.gen.go
index 80dbb66..77f014d 100644
--- a/internal/common/client/order/openapi_client.gen.go
+++ b/internal/common/client/order/openapi_client.gen.go
@@ -89,17 +89,17 @@ func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
 
 // The interface specification for the client above.
 type ClientInterface interface {
-	// PostCustomerCustomerIDOrdersWithBody request with any body
-	PostCustomerCustomerIDOrdersWithBody(ctx context.Context, customerID string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
+	// PostCustomerCustomerIdOrdersWithBody request with any body
+	PostCustomerCustomerIdOrdersWithBody(ctx context.Context, customerId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
 
-	PostCustomerCustomerIDOrders(ctx context.Context, customerID string, body PostCustomerCustomerIDOrdersJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
+	PostCustomerCustomerIdOrders(ctx context.Context, customerId string, body PostCustomerCustomerIdOrdersJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
 
-	// GetCustomerCustomerIDOrdersOrderID request
-	GetCustomerCustomerIDOrdersOrderID(ctx context.Context, customerID string, orderID string, reqEditors ...RequestEditorFn) (*http.Response, error)
+	// GetCustomerCustomerIdOrdersOrderId request
+	GetCustomerCustomerIdOrdersOrderId(ctx context.Context, customerId string, orderId string, reqEditors ...RequestEditorFn) (*http.Response, error)
 }
 
-func (c *Client) PostCustomerCustomerIDOrdersWithBody(ctx context.Context, customerID string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
-	req, err := NewPostCustomerCustomerIDOrdersRequestWithBody(c.Server, customerID, contentType, body)
+func (c *Client) PostCustomerCustomerIdOrdersWithBody(ctx context.Context, customerId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
+	req, err := NewPostCustomerCustomerIdOrdersRequestWithBody(c.Server, customerId, contentType, body)
 	if err != nil {
 		return nil, err
 	}
@@ -110,8 +110,8 @@ func (c *Client) PostCustomerCustomerIDOrdersWithBody(ctx context.Context, custo
 	return c.Client.Do(req)
 }
 
-func (c *Client) PostCustomerCustomerIDOrders(ctx context.Context, customerID string, body PostCustomerCustomerIDOrdersJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
-	req, err := NewPostCustomerCustomerIDOrdersRequest(c.Server, customerID, body)
+func (c *Client) PostCustomerCustomerIdOrders(ctx context.Context, customerId string, body PostCustomerCustomerIdOrdersJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
+	req, err := NewPostCustomerCustomerIdOrdersRequest(c.Server, customerId, body)
 	if err != nil {
 		return nil, err
 	}
@@ -122,8 +122,8 @@ func (c *Client) PostCustomerCustomerIDOrders(ctx context.Context, customerID st
 	return c.Client.Do(req)
 }
 
-func (c *Client) GetCustomerCustomerIDOrdersOrderID(ctx context.Context, customerID string, orderID string, reqEditors ...RequestEditorFn) (*http.Response, error) {
-	req, err := NewGetCustomerCustomerIDOrdersOrderIDRequest(c.Server, customerID, orderID)
+func (c *Client) GetCustomerCustomerIdOrdersOrderId(ctx context.Context, customerId string, orderId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
+	req, err := NewGetCustomerCustomerIdOrdersOrderIdRequest(c.Server, customerId, orderId)
 	if err != nil {
 		return nil, err
 	}
@@ -134,24 +134,24 @@ func (c *Client) GetCustomerCustomerIDOrdersOrderID(ctx context.Context, custome
 	return c.Client.Do(req)
 }
 
-// NewPostCustomerCustomerIDOrdersRequest calls the generic PostCustomerCustomerIDOrders builder with application/json body
-func NewPostCustomerCustomerIDOrdersRequest(server string, customerID string, body PostCustomerCustomerIDOrdersJSONRequestBody) (*http.Request, error) {
+// NewPostCustomerCustomerIdOrdersRequest calls the generic PostCustomerCustomerIdOrders builder with application/json body
+func NewPostCustomerCustomerIdOrdersRequest(server string, customerId string, body PostCustomerCustomerIdOrdersJSONRequestBody) (*http.Request, error) {
 	var bodyReader io.Reader
 	buf, err := json.Marshal(body)
 	if err != nil {
 		return nil, err
 	}
 	bodyReader = bytes.NewReader(buf)
-	return NewPostCustomerCustomerIDOrdersRequestWithBody(server, customerID, "application/json", bodyReader)
+	return NewPostCustomerCustomerIdOrdersRequestWithBody(server, customerId, "application/json", bodyReader)
 }
 
-// NewPostCustomerCustomerIDOrdersRequestWithBody generates requests for PostCustomerCustomerIDOrders with any type of body
-func NewPostCustomerCustomerIDOrdersRequestWithBody(server string, customerID string, contentType string, body io.Reader) (*http.Request, error) {
+// NewPostCustomerCustomerIdOrdersRequestWithBody generates requests for PostCustomerCustomerIdOrders with any type of body
+func NewPostCustomerCustomerIdOrdersRequestWithBody(server string, customerId string, contentType string, body io.Reader) (*http.Request, error) {
 	var err error
 
 	var pathParam0 string
 
-	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "customerID", runtime.ParamLocationPath, customerID)
+	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "customer_id", runtime.ParamLocationPath, customerId)
 	if err != nil {
 		return nil, err
 	}
@@ -181,20 +181,20 @@ func NewPostCustomerCustomerIDOrdersRequestWithBody(server string, customerID st
 	return req, nil
 }
 
-// NewGetCustomerCustomerIDOrdersOrderIDRequest generates requests for GetCustomerCustomerIDOrdersOrderID
-func NewGetCustomerCustomerIDOrdersOrderIDRequest(server string, customerID string, orderID string) (*http.Request, error) {
+// NewGetCustomerCustomerIdOrdersOrderIdRequest generates requests for GetCustomerCustomerIdOrdersOrderId
+func NewGetCustomerCustomerIdOrdersOrderIdRequest(server string, customerId string, orderId string) (*http.Request, error) {
 	var err error
 
 	var pathParam0 string
 
-	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "customerID", runtime.ParamLocationPath, customerID)
+	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "customer_id", runtime.ParamLocationPath, customerId)
 	if err != nil {
 		return nil, err
 	}
 
 	var pathParam1 string
 
-	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "orderID", runtime.ParamLocationPath, orderID)
+	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "order_id", runtime.ParamLocationPath, orderId)
 	if err != nil {
 		return nil, err
 	}
@@ -265,16 +265,16 @@ func WithBaseURL(baseURL string) ClientOption {
 
 // ClientWithResponsesInterface is the interface specification for the client with responses above.
 type ClientWithResponsesInterface interface {
-	// PostCustomerCustomerIDOrdersWithBodyWithResponse request with any body
-	PostCustomerCustomerIDOrdersWithBodyWithResponse(ctx context.Context, customerID string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostCustomerCustomerIDOrdersResponse, error)
+	// PostCustomerCustomerIdOrdersWithBodyWithResponse request with any body
+	PostCustomerCustomerIdOrdersWithBodyWithResponse(ctx context.Context, customerId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostCustomerCustomerIdOrdersResponse, error)
 
-	PostCustomerCustomerIDOrdersWithResponse(ctx context.Context, customerID string, body PostCustomerCustomerIDOrdersJSONRequestBody, reqEditors ...RequestEditorFn) (*PostCustomerCustomerIDOrdersResponse, error)
+	PostCustomerCustomerIdOrdersWithResponse(ctx context.Context, customerId string, body PostCustomerCustomerIdOrdersJSONRequestBody, reqEditors ...RequestEditorFn) (*PostCustomerCustomerIdOrdersResponse, error)
 
-	// GetCustomerCustomerIDOrdersOrderIDWithResponse request
-	GetCustomerCustomerIDOrdersOrderIDWithResponse(ctx context.Context, customerID string, orderID string, reqEditors ...RequestEditorFn) (*GetCustomerCustomerIDOrdersOrderIDResponse, error)
+	// GetCustomerCustomerIdOrdersOrderIdWithResponse request
+	GetCustomerCustomerIdOrdersOrderIdWithResponse(ctx context.Context, customerId string, orderId string, reqEditors ...RequestEditorFn) (*GetCustomerCustomerIdOrdersOrderIdResponse, error)
 }
 
-type PostCustomerCustomerIDOrdersResponse struct {
+type PostCustomerCustomerIdOrdersResponse struct {
 	Body         []byte
 	HTTPResponse *http.Response
 	JSON200      *Order
@@ -282,7 +282,7 @@ type PostCustomerCustomerIDOrdersResponse struct {
 }
 
 // Status returns HTTPResponse.Status
-func (r PostCustomerCustomerIDOrdersResponse) Status() string {
+func (r PostCustomerCustomerIdOrdersResponse) Status() string {
 	if r.HTTPResponse != nil {
 		return r.HTTPResponse.Status
 	}
@@ -290,14 +290,14 @@ func (r PostCustomerCustomerIDOrdersResponse) Status() string {
 }
 
 // StatusCode returns HTTPResponse.StatusCode
-func (r PostCustomerCustomerIDOrdersResponse) StatusCode() int {
+func (r PostCustomerCustomerIdOrdersResponse) StatusCode() int {
 	if r.HTTPResponse != nil {
 		return r.HTTPResponse.StatusCode
 	}
 	return 0
 }
 
-type GetCustomerCustomerIDOrdersOrderIDResponse struct {
+type GetCustomerCustomerIdOrdersOrderIdResponse struct {
 	Body         []byte
 	HTTPResponse *http.Response
 	JSON200      *Order
@@ -305,7 +305,7 @@ type GetCustomerCustomerIDOrdersOrderIDResponse struct {
 }
 
 // Status returns HTTPResponse.Status
-func (r GetCustomerCustomerIDOrdersOrderIDResponse) Status() string {
+func (r GetCustomerCustomerIdOrdersOrderIdResponse) Status() string {
 	if r.HTTPResponse != nil {
 		return r.HTTPResponse.Status
 	}
@@ -313,48 +313,48 @@ func (r GetCustomerCustomerIDOrdersOrderIDResponse) Status() string {
 }
 
 // StatusCode returns HTTPResponse.StatusCode
-func (r GetCustomerCustomerIDOrdersOrderIDResponse) StatusCode() int {
+func (r GetCustomerCustomerIdOrdersOrderIdResponse) StatusCode() int {
 	if r.HTTPResponse != nil {
 		return r.HTTPResponse.StatusCode
 	}
 	return 0
 }
 
-// PostCustomerCustomerIDOrdersWithBodyWithResponse request with arbitrary body returning *PostCustomerCustomerIDOrdersResponse
-func (c *ClientWithResponses) PostCustomerCustomerIDOrdersWithBodyWithResponse(ctx context.Context, customerID string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostCustomerCustomerIDOrdersResponse, error) {
-	rsp, err := c.PostCustomerCustomerIDOrdersWithBody(ctx, customerID, contentType, body, reqEditors...)
+// PostCustomerCustomerIdOrdersWithBodyWithResponse request with arbitrary body returning *PostCustomerCustomerIdOrdersResponse
+func (c *ClientWithResponses) PostCustomerCustomerIdOrdersWithBodyWithResponse(ctx context.Context, customerId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostCustomerCustomerIdOrdersResponse, error) {
+	rsp, err := c.PostCustomerCustomerIdOrdersWithBody(ctx, customerId, contentType, body, reqEditors...)
 	if err != nil {
 		return nil, err
 	}
-	return ParsePostCustomerCustomerIDOrdersResponse(rsp)
+	return ParsePostCustomerCustomerIdOrdersResponse(rsp)
 }
 
-func (c *ClientWithResponses) PostCustomerCustomerIDOrdersWithResponse(ctx context.Context, customerID string, body PostCustomerCustomerIDOrdersJSONRequestBody, reqEditors ...RequestEditorFn) (*PostCustomerCustomerIDOrdersResponse, error) {
-	rsp, err := c.PostCustomerCustomerIDOrders(ctx, customerID, body, reqEditors...)
+func (c *ClientWithResponses) PostCustomerCustomerIdOrdersWithResponse(ctx context.Context, customerId string, body PostCustomerCustomerIdOrdersJSONRequestBody, reqEditors ...RequestEditorFn) (*PostCustomerCustomerIdOrdersResponse, error) {
+	rsp, err := c.PostCustomerCustomerIdOrders(ctx, customerId, body, reqEditors...)
 	if err != nil {
 		return nil, err
 	}
-	return ParsePostCustomerCustomerIDOrdersResponse(rsp)
+	return ParsePostCustomerCustomerIdOrdersResponse(rsp)
 }
 
-// GetCustomerCustomerIDOrdersOrderIDWithResponse request returning *GetCustomerCustomerIDOrdersOrderIDResponse
-func (c *ClientWithResponses) GetCustomerCustomerIDOrdersOrderIDWithResponse(ctx context.Context, customerID string, orderID string, reqEditors ...RequestEditorFn) (*GetCustomerCustomerIDOrdersOrderIDResponse, error) {
-	rsp, err := c.GetCustomerCustomerIDOrdersOrderID(ctx, customerID, orderID, reqEditors...)
+// GetCustomerCustomerIdOrdersOrderIdWithResponse request returning *GetCustomerCustomerIdOrdersOrderIdResponse
+func (c *ClientWithResponses) GetCustomerCustomerIdOrdersOrderIdWithResponse(ctx context.Context, customerId string, orderId string, reqEditors ...RequestEditorFn) (*GetCustomerCustomerIdOrdersOrderIdResponse, error) {
+	rsp, err := c.GetCustomerCustomerIdOrdersOrderId(ctx, customerId, orderId, reqEditors...)
 	if err != nil {
 		return nil, err
 	}
-	return ParseGetCustomerCustomerIDOrdersOrderIDResponse(rsp)
+	return ParseGetCustomerCustomerIdOrdersOrderIdResponse(rsp)
 }
 
-// ParsePostCustomerCustomerIDOrdersResponse parses an HTTP response from a PostCustomerCustomerIDOrdersWithResponse call
-func ParsePostCustomerCustomerIDOrdersResponse(rsp *http.Response) (*PostCustomerCustomerIDOrdersResponse, error) {
+// ParsePostCustomerCustomerIdOrdersResponse parses an HTTP response from a PostCustomerCustomerIdOrdersWithResponse call
+func ParsePostCustomerCustomerIdOrdersResponse(rsp *http.Response) (*PostCustomerCustomerIdOrdersResponse, error) {
 	bodyBytes, err := io.ReadAll(rsp.Body)
 	defer func() { _ = rsp.Body.Close() }()
 	if err != nil {
 		return nil, err
 	}
 
-	response := &PostCustomerCustomerIDOrdersResponse{
+	response := &PostCustomerCustomerIdOrdersResponse{
 		Body:         bodyBytes,
 		HTTPResponse: rsp,
 	}
@@ -379,15 +379,15 @@ func ParsePostCustomerCustomerIDOrdersResponse(rsp *http.Response) (*PostCustome
 	return response, nil
 }
 
-// ParseGetCustomerCustomerIDOrdersOrderIDResponse parses an HTTP response from a GetCustomerCustomerIDOrdersOrderIDWithResponse call
-func ParseGetCustomerCustomerIDOrdersOrderIDResponse(rsp *http.Response) (*GetCustomerCustomerIDOrdersOrderIDResponse, error) {
+// ParseGetCustomerCustomerIdOrdersOrderIdResponse parses an HTTP response from a GetCustomerCustomerIdOrdersOrderIdWithResponse call
+func ParseGetCustomerCustomerIdOrdersOrderIdResponse(rsp *http.Response) (*GetCustomerCustomerIdOrdersOrderIdResponse, error) {
 	bodyBytes, err := io.ReadAll(rsp.Body)
 	defer func() { _ = rsp.Body.Close() }()
 	if err != nil {
 		return nil, err
 	}
 
-	response := &GetCustomerCustomerIDOrdersOrderIDResponse{
+	response := &GetCustomerCustomerIdOrdersOrderIdResponse{
 		Body:         bodyBytes,
 		HTTPResponse: rsp,
 	}
diff --git a/internal/common/client/order/openapi_types.gen.go b/internal/common/client/order/openapi_types.gen.go
index 8021b88..3b659fb 100644
--- a/internal/common/client/order/openapi_types.gen.go
+++ b/internal/common/client/order/openapi_types.gen.go
@@ -5,7 +5,7 @@ package order
 
 // CreateOrderRequest defines model for CreateOrderRequest.
 type CreateOrderRequest struct {
-	CustomerID string             `json:"customerID"`
+	CustomerId string             `json:"customer_id"`
 	Items      []ItemWithQuantity `json:"items"`
 }
 
@@ -18,7 +18,7 @@ type Error struct {
 type Item struct {
 	Id       string `json:"id"`
 	Name     string `json:"name"`
-	PriceID  string `json:"priceID"`
+	PriceId  string `json:"price_id"`
 	Quantity int32  `json:"quantity"`
 }
 
@@ -30,12 +30,12 @@ type ItemWithQuantity struct {
 
 // Order defines model for Order.
 type Order struct {
-	CustomerID  string `json:"customerID"`
+	CustomerId  string `json:"customer_id"`
 	Id          string `json:"id"`
 	Items       []Item `json:"items"`
-	PaymentLink string `json:"paymentLink"`
+	PaymentLink string `json:"payment_link"`
 	Status      string `json:"status"`
 }
 
-// PostCustomerCustomerIDOrdersJSONRequestBody defines body for PostCustomerCustomerIDOrders for application/json ContentType.
-type PostCustomerCustomerIDOrdersJSONRequestBody = CreateOrderRequest
+// PostCustomerCustomerIdOrdersJSONRequestBody defines body for PostCustomerCustomerIdOrders for application/json ContentType.
+type PostCustomerCustomerIdOrdersJSONRequestBody = CreateOrderRequest
diff --git a/internal/common/response.go b/internal/common/response.go
new file mode 100644
index 0000000..7be1fe7
--- /dev/null
+++ b/internal/common/response.go
@@ -0,0 +1,43 @@
+package common
+
+import (
+	"net/http"
+
+	"github.com/ghost-yu/go_shop_second/common/tracing"
+	"github.com/gin-gonic/gin"
+)
+
+type BaseResponse struct{}
+
+type response struct {
+	Errno   int    `json:"errno"`
+	Message string `json:"message"`
+	Data    any    `json:"data"`
+	TraceID string `json:"trace_id"`
+}
+
+func (base *BaseResponse) Response(c *gin.Context, err error, data interface{}) {
+	if err != nil {
+		base.error(c, err)
+	} else {
+		base.success(c, data)
+	}
+}
+
+func (base *BaseResponse) success(c *gin.Context, data interface{}) {
+	c.JSON(http.StatusOK, response{
+		Errno:   0,
+		Message: "success",
+		Data:    data,
+		TraceID: tracing.TraceID(c.Request.Context()),
+	})
+}
+
+func (base *BaseResponse) error(c *gin.Context, err error) {
+	c.JSON(http.StatusOK, response{
+		Errno:   2,
+		Message: err.Error(),
+		Data:    nil,
+		TraceID: tracing.TraceID(c.Request.Context()),
+	})
+}
diff --git a/internal/order/http.go b/internal/order/http.go
index eebdb99..f06c157 100644
--- a/internal/order/http.go
+++ b/internal/order/http.go
@@ -2,10 +2,9 @@ package main
 
 import (
 	"fmt"
-	"net/http"
 
+	"github.com/ghost-yu/go_shop_second/common"
 	client "github.com/ghost-yu/go_shop_second/common/client/order"
-	"github.com/ghost-yu/go_shop_second/common/tracing"
 	"github.com/ghost-yu/go_shop_second/order/app"
 	"github.com/ghost-yu/go_shop_second/order/app/command"
 	"github.com/ghost-yu/go_shop_second/order/app/query"
@@ -14,52 +13,55 @@ import (
 )
 
 type HTTPServer struct {
+	common.BaseResponse
 	app app.Application
 }
 
 func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
-	ctx, span := tracing.Start(c, "PostCustomerCustomerIDOrders")
-	defer span.End()
+	var (
+		req  client.CreateOrderRequest
+		err  error
+		resp struct {
+			CustomerID  string `json:"customer_id"`
+			OrderID     string `json:"order_id"`
+			RedirectURL string `json:"redirect_url"`
+		}
+	)
+	defer func() {
+		H.Response(c, err, &resp)
+	}()
 
-	var req client.CreateOrderRequest
-	if err := c.ShouldBindJSON(&req); err != nil {
-		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
+	if err = c.ShouldBindJSON(&req); err != nil {
 		return
 	}
-	r, err := H.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
-		CustomerID: req.CustomerID,
+	r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
+		CustomerID: req.CustomerId,
 		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
 	})
 	if err != nil {
-		c.JSON(http.StatusOK, gin.H{"error": err})
 		return
 	}
-	c.JSON(http.StatusOK, gin.H{
-		"message":      "success",
-		"trace_id":     tracing.TraceID(ctx),
-		"customer_id":  req.CustomerID,
-		"order_id":     r.OrderID,
-		"redirect_url": fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderID),
-	})
+	resp.CustomerID = req.CustomerId
+	resp.RedirectURL = fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderID)
+	resp.OrderID = r.OrderID
 }
 
 func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
-	ctx, span := tracing.Start(c, "GetCustomerCustomerIDOrdersOrderID")
-	defer span.End()
+	var (
+		err  error
+		resp interface{}
+	)
+	defer func() {
+		H.Response(c, err, resp)
+	}()
 
-	o, err := H.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
+	o, err := H.app.Queries.GetCustomerOrder.Handle(c.Request.Context(), query.GetCustomerOrder{
 		OrderID:    orderID,
 		CustomerID: customerID,
 	})
 	if err != nil {
-		c.JSON(http.StatusOK, gin.H{"error": err})
 		return
 	}
-	c.JSON(http.StatusOK, gin.H{
-		"message":  "success",
-		"trace_id": tracing.TraceID(ctx),
-		"data": gin.H{
-			"Order": o,
-		},
-	})
+
+	resp = convertor.NewOrderConvertor().EntityToClient(o)
 }
diff --git a/internal/order/ports/openapi_api.gen.go b/internal/order/ports/openapi_api.gen.go
index 9cbaed0..33ca9fd 100644
--- a/internal/order/ports/openapi_api.gen.go
+++ b/internal/order/ports/openapi_api.gen.go
@@ -14,11 +14,11 @@ import (
 // ServerInterface represents all server handlers.
 type ServerInterface interface {
 
-	// (POST /customer/{customerID}/orders)
-	PostCustomerCustomerIDOrders(c *gin.Context, customerID string)
+	// (POST /customer/{customer_id}/orders)
+	PostCustomerCustomerIdOrders(c *gin.Context, customerId string)
 
-	// (GET /customer/{customerID}/orders/{orderID})
-	GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string)
+	// (GET /customer/{customer_id}/orders/{order_id})
+	GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerId string, orderId string)
 }
 
 // ServerInterfaceWrapper converts contexts to parameters.
@@ -30,17 +30,17 @@ type ServerInterfaceWrapper struct {
 
 type MiddlewareFunc func(c *gin.Context)
 
-// PostCustomerCustomerIDOrders operation middleware
-func (siw *ServerInterfaceWrapper) PostCustomerCustomerIDOrders(c *gin.Context) {
+// PostCustomerCustomerIdOrders operation middleware
+func (siw *ServerInterfaceWrapper) PostCustomerCustomerIdOrders(c *gin.Context) {
 
 	var err error
 
-	// ------------- Path parameter "customerID" -------------
-	var customerID string
+	// ------------- Path parameter "customer_id" -------------
+	var customerId string
 
-	err = runtime.BindStyledParameterWithOptions("simple", "customerID", c.Param("customerID"), &customerID, runtime.BindStyledParameterOptions{Explode: false, Required: true})
+	err = runtime.BindStyledParameterWithOptions("simple", "customer_id", c.Param("customer_id"), &customerId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
 	if err != nil {
-		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter customerID: %w", err), http.StatusBadRequest)
+		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter customer_id: %w", err), http.StatusBadRequest)
 		return
 	}
 
@@ -51,29 +51,29 @@ func (siw *ServerInterfaceWrapper) PostCustomerCustomerIDOrders(c *gin.Context)
 		}
 	}
 
-	siw.Handler.PostCustomerCustomerIDOrders(c, customerID)
+	siw.Handler.PostCustomerCustomerIdOrders(c, customerId)
 }
 
-// GetCustomerCustomerIDOrdersOrderID operation middleware
-func (siw *ServerInterfaceWrapper) GetCustomerCustomerIDOrdersOrderID(c *gin.Context) {
+// GetCustomerCustomerIdOrdersOrderId operation middleware
+func (siw *ServerInterfaceWrapper) GetCustomerCustomerIdOrdersOrderId(c *gin.Context) {
 
 	var err error
 
-	// ------------- Path parameter "customerID" -------------
-	var customerID string
+	// ------------- Path parameter "customer_id" -------------
+	var customerId string
 
-	err = runtime.BindStyledParameterWithOptions("simple", "customerID", c.Param("customerID"), &customerID, runtime.BindStyledParameterOptions{Explode: false, Required: true})
+	err = runtime.BindStyledParameterWithOptions("simple", "customer_id", c.Param("customer_id"), &customerId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
 	if err != nil {
-		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter customerID: %w", err), http.StatusBadRequest)
+		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter customer_id: %w", err), http.StatusBadRequest)
 		return
 	}
 
-	// ------------- Path parameter "orderID" -------------
-	var orderID string
+	// ------------- Path parameter "order_id" -------------
+	var orderId string
 
-	err = runtime.BindStyledParameterWithOptions("simple", "orderID", c.Param("orderID"), &orderID, runtime.BindStyledParameterOptions{Explode: false, Required: true})
+	err = runtime.BindStyledParameterWithOptions("simple", "order_id", c.Param("order_id"), &orderId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
 	if err != nil {
-		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter orderID: %w", err), http.StatusBadRequest)
+		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter order_id: %w", err), http.StatusBadRequest)
 		return
 	}
 
@@ -84,7 +84,7 @@ func (siw *ServerInterfaceWrapper) GetCustomerCustomerIDOrdersOrderID(c *gin.Con
 		}
 	}
 
-	siw.Handler.GetCustomerCustomerIDOrdersOrderID(c, customerID, orderID)
+	siw.Handler.GetCustomerCustomerIdOrdersOrderId(c, customerId, orderId)
 }
 
 // GinServerOptions provides options for the Gin server.
@@ -114,6 +114,6 @@ func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options
 		ErrorHandler:       errorHandler,
 	}
 
-	router.POST(options.BaseURL+"/customer/:customerID/orders", wrapper.PostCustomerCustomerIDOrders)
-	router.GET(options.BaseURL+"/customer/:customerID/orders/:orderID", wrapper.GetCustomerCustomerIDOrdersOrderID)
+	router.POST(options.BaseURL+"/customer/:customer_id/orders", wrapper.PostCustomerCustomerIdOrders)
+	router.GET(options.BaseURL+"/customer/:customer_id/orders/:order_id", wrapper.GetCustomerCustomerIdOrdersOrderId)
 }
diff --git a/internal/order/ports/openapi_types.gen.go b/internal/order/ports/openapi_types.gen.go
index 0d69dbf..97fe85c 100644
--- a/internal/order/ports/openapi_types.gen.go
+++ b/internal/order/ports/openapi_types.gen.go
@@ -5,7 +5,7 @@ package ports
 
 // CreateOrderRequest defines model for CreateOrderRequest.
 type CreateOrderRequest struct {
-	CustomerID string             `json:"customerID"`
+	CustomerId string             `json:"customer_id"`
 	Items      []ItemWithQuantity `json:"items"`
 }
 
@@ -18,7 +18,7 @@ type Error struct {
 type Item struct {
 	Id       string `json:"id"`
 	Name     string `json:"name"`
-	PriceID  string `json:"priceID"`
+	PriceId  string `json:"price_id"`
 	Quantity int32  `json:"quantity"`
 }
 
@@ -30,12 +30,12 @@ type ItemWithQuantity struct {
 
 // Order defines model for Order.
 type Order struct {
-	CustomerID  string `json:"customerID"`
+	CustomerId  string `json:"customer_id"`
 	Id          string `json:"id"`
 	Items       []Item `json:"items"`
-	PaymentLink string `json:"paymentLink"`
+	PaymentLink string `json:"payment_link"`
 	Status      string `json:"status"`
 }
 
-// PostCustomerCustomerIDOrdersJSONRequestBody defines body for PostCustomerCustomerIDOrders for application/json ContentType.
-type PostCustomerCustomerIDOrdersJSONRequestBody = CreateOrderRequest
+// PostCustomerCustomerIdOrdersJSONRequestBody defines body for PostCustomerCustomerIdOrders for application/json ContentType.
+type PostCustomerCustomerIdOrdersJSONRequestBody = CreateOrderRequest
~~~

# Lesson Pair Diff Report

- FromBranch: lesson3
- ToBranch: lesson4

## Short Summary

~~~text
 23 files changed, 2096 insertions(+), 8 deletions(-)
~~~

## File Stats

~~~text
 .gitignore                                         |   1 +
 Makefile                                           |  10 +
 api/README.md                                      |   2 +
 api/openapi/order.yml                              | 133 ++++++
 api/orderpb/order.proto                            |  41 ++
 internal/common/client/order/openapi_client.gen.go | 413 +++++++++++++++++
 internal/common/client/order/openapi_types.gen.go  |  41 ++
 internal/common/config/global.yaml                 |   4 +
 internal/common/config/viper.go                    |  11 +
 internal/common/genproto/orderpb/order.pb.go       | 508 +++++++++++++++++++++
 internal/common/genproto/orderpb/order_grpc.pb.go  | 176 +++++++
 internal/common/go.mod                             |  37 ++
 internal/common/go.sum                             |  96 ++++
 internal/kitchen/go.mod                            |   7 +
 internal/order/go.mod                              |  36 +-
 internal/order/go.sum                              |  78 +++-
 internal/order/ports/openapi_api.gen.go            | 119 +++++
 internal/order/ports/openapi_types.gen.go          |  41 ++
 internal/payment/go.mod                            |   7 +
 internal/stock/go.mod                              |   7 +
 scripts/genopenapi.sh                              |  54 +++
 scripts/genproto.sh                                |  62 +++
 scripts/lib.sh                                     | 220 +++++++++
 23 files changed, 2096 insertions(+), 8 deletions(-)
~~~

## Commit Comparison

~~~text
> d6b3140 genproto, genopenapi
~~~

## Changed Files

~~~text
.gitignore
Makefile
api/README.md
api/openapi/order.yml
api/orderpb/order.proto
internal/common/client/order/openapi_client.gen.go
internal/common/client/order/openapi_types.gen.go
internal/common/config/global.yaml
internal/common/config/viper.go
internal/common/genproto/orderpb/order.pb.go
internal/common/genproto/orderpb/order_grpc.pb.go
internal/common/go.mod
internal/common/go.sum
internal/kitchen/go.mod
internal/order/go.mod
internal/order/go.sum
internal/order/ports/openapi_api.gen.go
internal/order/ports/openapi_types.gen.go
internal/payment/go.mod
internal/stock/go.mod
scripts/genopenapi.sh
scripts/genproto.sh
scripts/lib.sh
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
.gitignore
Makefile
api/README.md
api/openapi/order.yml
api/orderpb/order.proto
internal/common/config/global.yaml
internal/common/config/viper.go
scripts/genopenapi.sh
scripts/genproto.sh
scripts/lib.sh
~~~

## Full Diff

~~~diff
diff --git a/.gitignore b/.gitignore
index 6f72f89..42e4918 100644
--- a/.gitignore
+++ b/.gitignore
@@ -23,3 +23,4 @@ go.work.sum
 
 # env file
 .env
+.idea
\ No newline at end of file
diff --git a/Makefile b/Makefile
new file mode 100644
index 0000000..6faaf1e
--- /dev/null
+++ b/Makefile
@@ -0,0 +1,10 @@
+.PHONY: gen
+gen: genproto genopenapi
+
+.PHONY: genproto
+genproto:
+	@./scripts/genproto.sh
+
+.PHONY: genopenapi
+genopenapi:
+	@./scripts/genopenapi.sh
\ No newline at end of file
diff --git a/api/README.md b/api/README.md
new file mode 100644
index 0000000..c7b107c
--- /dev/null
+++ b/api/README.md
@@ -0,0 +1,2 @@
+## api
+存放接口相关的通讯协议，例如openapi、proto
\ No newline at end of file
diff --git a/api/openapi/order.yml b/api/openapi/order.yml
new file mode 100644
index 0000000..f421fce
--- /dev/null
+++ b/api/openapi/order.yml
@@ -0,0 +1,133 @@
+openapi: 3.0.3
+info:
+  title: order service
+  description: order service
+  version: 1.0.0
+servers:
+  - url: 'https://{hostname}/api'
+    variables:
+      hostname:
+        default: 127.0.0.1
+
+paths:
+  /customer/{customerID}/orders/{orderID}:
+    get:
+      description: "get order"
+      parameters:
+        - in: path
+          name: customerID
+          schema:
+            type: string
+          required: true
+
+        - in: path
+          name: orderID
+          schema:
+            type: string
+          required: true
+
+      responses:
+        '200':
+          description: todo
+          content:
+            application/json:
+              schema:
+                $ref: '#/components/schemas/Order'
+
+        default:
+          description: todo
+          content:
+            application/json:
+              schema:
+                $ref: '#/components/schemas/Error'
+
+  /customer/{customerID}/orders:
+    post:
+      description: "create order"
+      parameters:
+        - in: path
+          name: customerID
+          schema:
+            type: string
+          required: true
+
+      requestBody:
+        required: true
+        content:
+          application/json:
+            schema:
+              $ref: '#/components/schemas/CreateOrderRequest'
+
+      responses:
+        '200':
+          description: todo
+          content:
+            application/json:
+              schema:
+                $ref: '#/components/schemas/Order'
+
+        default:
+          description: todo
+          content:
+            application/json:
+              schema:
+                $ref: '#/components/schemas/Error'
+
+components:
+  schemas:
+    Order:
+      type: object
+      properties:
+        id:
+          type: string
+        customerID:
+          type: string
+        status:
+          type: string
+        items:
+          type: array
+          items:
+            $ref: '#/components/schemas/Item'
+        paymentLink:
+          type: string
+
+    Item:
+      type: object
+      properties:
+        id:
+          type: string
+        name:
+          type: string
+        quantity:
+          type: integer
+          format: int32
+        priceID:
+          type: string
+
+    Error:
+      type: object
+      properties:
+        message:
+          type: string
+
+    CreateOrderRequest:
+      type: object
+      required:
+        - customerID
+        - items
+      properties:
+        customerID:
+          type: string
+        items:
+          type: array
+          items:
+            $ref: '#/components/schemas/ItemWithQuantity'
+
+    ItemWithQuantity:
+      type: object
+      properties:
+        id:
+          type: string
+        quantity:
+          type: integer
+          format: int32
\ No newline at end of file
diff --git a/api/orderpb/order.proto b/api/orderpb/order.proto
new file mode 100644
index 0000000..24b4e62
--- /dev/null
+++ b/api/orderpb/order.proto
@@ -0,0 +1,41 @@
+syntax = "proto3";
+package orderpb;
+
+option go_package = "github.com/ghost-yu/go_shop_second/internal/common/genproto/orderpb";
+
+import "google/protobuf/empty.proto";
+
+service OrderService {
+  rpc CreateOrder(CreateOrderRequest) returns (google.protobuf.Empty);
+  rpc GetOrder(GetOrderRequest) returns (Order);
+  rpc UpdateOrder(Order) returns (google.protobuf.Empty);
+}
+
+message CreateOrderRequest {
+  string CustomerID = 1;
+  repeated ItemWithQuantity Items = 2;
+}
+
+message GetOrderRequest {
+  string OrderID = 1;
+  string CustomerID = 2;
+}
+
+message ItemWithQuantity {
+  string ID = 1;
+  int32 Quantity = 2;
+}
+
+message Item {
+  string ID = 1;
+  string Name = 2;
+  int32 Quantity = 3;
+  string PriceID = 4;
+}
+
+message Order {
+  string ID = 1;
+  string CustomerID = 2;
+  string Status = 3;
+  repeated Item Items = 4;
+}
\ No newline at end of file
diff --git a/internal/common/client/order/openapi_client.gen.go b/internal/common/client/order/openapi_client.gen.go
new file mode 100644
index 0000000..80dbb66
--- /dev/null
+++ b/internal/common/client/order/openapi_client.gen.go
@@ -0,0 +1,413 @@
+// Package order provides primitives to interact with the openapi HTTP API.
+//
+// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
+package order
+
+import (
+	"bytes"
+	"context"
+	"encoding/json"
+	"fmt"
+	"io"
+	"net/http"
+	"net/url"
+	"strings"
+
+	"github.com/oapi-codegen/runtime"
+)
+
+// RequestEditorFn  is the function signature for the RequestEditor callback function
+type RequestEditorFn func(ctx context.Context, req *http.Request) error
+
+// Doer performs HTTP requests.
+//
+// The standard http.Client implements this interface.
+type HttpRequestDoer interface {
+	Do(req *http.Request) (*http.Response, error)
+}
+
+// Client which conforms to the OpenAPI3 specification for this service.
+type Client struct {
+	// The endpoint of the server conforming to this interface, with scheme,
+	// https://api.deepmap.com for example. This can contain a path relative
+	// to the server, such as https://api.deepmap.com/dev-test, and all the
+	// paths in the swagger spec will be appended to the server.
+	Server string
+
+	// Doer for performing requests, typically a *http.Client with any
+	// customized settings, such as certificate chains.
+	Client HttpRequestDoer
+
+	// A list of callbacks for modifying requests which are generated before sending over
+	// the network.
+	RequestEditors []RequestEditorFn
+}
+
+// ClientOption allows setting custom parameters during construction
+type ClientOption func(*Client) error
+
+// Creates a new Client, with reasonable defaults
+func NewClient(server string, opts ...ClientOption) (*Client, error) {
+	// create a client with sane default values
+	client := Client{
+		Server: server,
+	}
+	// mutate client and add all optional params
+	for _, o := range opts {
+		if err := o(&client); err != nil {
+			return nil, err
+		}
+	}
+	// ensure the server URL always has a trailing slash
+	if !strings.HasSuffix(client.Server, "/") {
+		client.Server += "/"
+	}
+	// create httpClient, if not already present
+	if client.Client == nil {
+		client.Client = &http.Client{}
+	}
+	return &client, nil
+}
+
+// WithHTTPClient allows overriding the default Doer, which is
+// automatically created using http.Client. This is useful for tests.
+func WithHTTPClient(doer HttpRequestDoer) ClientOption {
+	return func(c *Client) error {
+		c.Client = doer
+		return nil
+	}
+}
+
+// WithRequestEditorFn allows setting up a callback function, which will be
+// called right before sending the request. This can be used to mutate the request.
+func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
+	return func(c *Client) error {
+		c.RequestEditors = append(c.RequestEditors, fn)
+		return nil
+	}
+}
+
+// The interface specification for the client above.
+type ClientInterface interface {
+	// PostCustomerCustomerIDOrdersWithBody request with any body
+	PostCustomerCustomerIDOrdersWithBody(ctx context.Context, customerID string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
+
+	PostCustomerCustomerIDOrders(ctx context.Context, customerID string, body PostCustomerCustomerIDOrdersJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
+
+	// GetCustomerCustomerIDOrdersOrderID request
+	GetCustomerCustomerIDOrdersOrderID(ctx context.Context, customerID string, orderID string, reqEditors ...RequestEditorFn) (*http.Response, error)
+}
+
+func (c *Client) PostCustomerCustomerIDOrdersWithBody(ctx context.Context, customerID string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
+	req, err := NewPostCustomerCustomerIDOrdersRequestWithBody(c.Server, customerID, contentType, body)
+	if err != nil {
+		return nil, err
+	}
+	req = req.WithContext(ctx)
+	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
+		return nil, err
+	}
+	return c.Client.Do(req)
+}
+
+func (c *Client) PostCustomerCustomerIDOrders(ctx context.Context, customerID string, body PostCustomerCustomerIDOrdersJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
+	req, err := NewPostCustomerCustomerIDOrdersRequest(c.Server, customerID, body)
+	if err != nil {
+		return nil, err
+	}
+	req = req.WithContext(ctx)
+	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
+		return nil, err
+	}
+	return c.Client.Do(req)
+}
+
+func (c *Client) GetCustomerCustomerIDOrdersOrderID(ctx context.Context, customerID string, orderID string, reqEditors ...RequestEditorFn) (*http.Response, error) {
+	req, err := NewGetCustomerCustomerIDOrdersOrderIDRequest(c.Server, customerID, orderID)
+	if err != nil {
+		return nil, err
+	}
+	req = req.WithContext(ctx)
+	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
+		return nil, err
+	}
+	return c.Client.Do(req)
+}
+
+// NewPostCustomerCustomerIDOrdersRequest calls the generic PostCustomerCustomerIDOrders builder with application/json body
+func NewPostCustomerCustomerIDOrdersRequest(server string, customerID string, body PostCustomerCustomerIDOrdersJSONRequestBody) (*http.Request, error) {
+	var bodyReader io.Reader
+	buf, err := json.Marshal(body)
+	if err != nil {
+		return nil, err
+	}
+	bodyReader = bytes.NewReader(buf)
+	return NewPostCustomerCustomerIDOrdersRequestWithBody(server, customerID, "application/json", bodyReader)
+}
+
+// NewPostCustomerCustomerIDOrdersRequestWithBody generates requests for PostCustomerCustomerIDOrders with any type of body
+func NewPostCustomerCustomerIDOrdersRequestWithBody(server string, customerID string, contentType string, body io.Reader) (*http.Request, error) {
+	var err error
+
+	var pathParam0 string
+
+	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "customerID", runtime.ParamLocationPath, customerID)
+	if err != nil {
+		return nil, err
+	}
+
+	serverURL, err := url.Parse(server)
+	if err != nil {
+		return nil, err
+	}
+
+	operationPath := fmt.Sprintf("/customer/%s/orders", pathParam0)
+	if operationPath[0] == '/' {
+		operationPath = "." + operationPath
+	}
+
+	queryURL, err := serverURL.Parse(operationPath)
+	if err != nil {
+		return nil, err
+	}
+
+	req, err := http.NewRequest("POST", queryURL.String(), body)
+	if err != nil {
+		return nil, err
+	}
+
+	req.Header.Add("Content-Type", contentType)
+
+	return req, nil
+}
+
+// NewGetCustomerCustomerIDOrdersOrderIDRequest generates requests for GetCustomerCustomerIDOrdersOrderID
+func NewGetCustomerCustomerIDOrdersOrderIDRequest(server string, customerID string, orderID string) (*http.Request, error) {
+	var err error
+
+	var pathParam0 string
+
+	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "customerID", runtime.ParamLocationPath, customerID)
+	if err != nil {
+		return nil, err
+	}
+
+	var pathParam1 string
+
+	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "orderID", runtime.ParamLocationPath, orderID)
+	if err != nil {
+		return nil, err
+	}
+
+	serverURL, err := url.Parse(server)
+	if err != nil {
+		return nil, err
+	}
+
+	operationPath := fmt.Sprintf("/customer/%s/orders/%s", pathParam0, pathParam1)
+	if operationPath[0] == '/' {
+		operationPath = "." + operationPath
+	}
+
+	queryURL, err := serverURL.Parse(operationPath)
+	if err != nil {
+		return nil, err
+	}
+
+	req, err := http.NewRequest("GET", queryURL.String(), nil)
+	if err != nil {
+		return nil, err
+	}
+
+	return req, nil
+}
+
+func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
+	for _, r := range c.RequestEditors {
+		if err := r(ctx, req); err != nil {
+			return err
+		}
+	}
+	for _, r := range additionalEditors {
+		if err := r(ctx, req); err != nil {
+			return err
+		}
+	}
+	return nil
+}
+
+// ClientWithResponses builds on ClientInterface to offer response payloads
+type ClientWithResponses struct {
+	ClientInterface
+}
+
+// NewClientWithResponses creates a new ClientWithResponses, which wraps
+// Client with return type handling
+func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
+	client, err := NewClient(server, opts...)
+	if err != nil {
+		return nil, err
+	}
+	return &ClientWithResponses{client}, nil
+}
+
+// WithBaseURL overrides the baseURL.
+func WithBaseURL(baseURL string) ClientOption {
+	return func(c *Client) error {
+		newBaseURL, err := url.Parse(baseURL)
+		if err != nil {
+			return err
+		}
+		c.Server = newBaseURL.String()
+		return nil
+	}
+}
+
+// ClientWithResponsesInterface is the interface specification for the client with responses above.
+type ClientWithResponsesInterface interface {
+	// PostCustomerCustomerIDOrdersWithBodyWithResponse request with any body
+	PostCustomerCustomerIDOrdersWithBodyWithResponse(ctx context.Context, customerID string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostCustomerCustomerIDOrdersResponse, error)
+
+	PostCustomerCustomerIDOrdersWithResponse(ctx context.Context, customerID string, body PostCustomerCustomerIDOrdersJSONRequestBody, reqEditors ...RequestEditorFn) (*PostCustomerCustomerIDOrdersResponse, error)
+
+	// GetCustomerCustomerIDOrdersOrderIDWithResponse request
+	GetCustomerCustomerIDOrdersOrderIDWithResponse(ctx context.Context, customerID string, orderID string, reqEditors ...RequestEditorFn) (*GetCustomerCustomerIDOrdersOrderIDResponse, error)
+}
+
+type PostCustomerCustomerIDOrdersResponse struct {
+	Body         []byte
+	HTTPResponse *http.Response
+	JSON200      *Order
+	JSONDefault  *Error
+}
+
+// Status returns HTTPResponse.Status
+func (r PostCustomerCustomerIDOrdersResponse) Status() string {
+	if r.HTTPResponse != nil {
+		return r.HTTPResponse.Status
+	}
+	return http.StatusText(0)
+}
+
+// StatusCode returns HTTPResponse.StatusCode
+func (r PostCustomerCustomerIDOrdersResponse) StatusCode() int {
+	if r.HTTPResponse != nil {
+		return r.HTTPResponse.StatusCode
+	}
+	return 0
+}
+
+type GetCustomerCustomerIDOrdersOrderIDResponse struct {
+	Body         []byte
+	HTTPResponse *http.Response
+	JSON200      *Order
+	JSONDefault  *Error
+}
+
+// Status returns HTTPResponse.Status
+func (r GetCustomerCustomerIDOrdersOrderIDResponse) Status() string {
+	if r.HTTPResponse != nil {
+		return r.HTTPResponse.Status
+	}
+	return http.StatusText(0)
+}
+
+// StatusCode returns HTTPResponse.StatusCode
+func (r GetCustomerCustomerIDOrdersOrderIDResponse) StatusCode() int {
+	if r.HTTPResponse != nil {
+		return r.HTTPResponse.StatusCode
+	}
+	return 0
+}
+
+// PostCustomerCustomerIDOrdersWithBodyWithResponse request with arbitrary body returning *PostCustomerCustomerIDOrdersResponse
+func (c *ClientWithResponses) PostCustomerCustomerIDOrdersWithBodyWithResponse(ctx context.Context, customerID string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostCustomerCustomerIDOrdersResponse, error) {
+	rsp, err := c.PostCustomerCustomerIDOrdersWithBody(ctx, customerID, contentType, body, reqEditors...)
+	if err != nil {
+		return nil, err
+	}
+	return ParsePostCustomerCustomerIDOrdersResponse(rsp)
+}
+
+func (c *ClientWithResponses) PostCustomerCustomerIDOrdersWithResponse(ctx context.Context, customerID string, body PostCustomerCustomerIDOrdersJSONRequestBody, reqEditors ...RequestEditorFn) (*PostCustomerCustomerIDOrdersResponse, error) {
+	rsp, err := c.PostCustomerCustomerIDOrders(ctx, customerID, body, reqEditors...)
+	if err != nil {
+		return nil, err
+	}
+	return ParsePostCustomerCustomerIDOrdersResponse(rsp)
+}
+
+// GetCustomerCustomerIDOrdersOrderIDWithResponse request returning *GetCustomerCustomerIDOrdersOrderIDResponse
+func (c *ClientWithResponses) GetCustomerCustomerIDOrdersOrderIDWithResponse(ctx context.Context, customerID string, orderID string, reqEditors ...RequestEditorFn) (*GetCustomerCustomerIDOrdersOrderIDResponse, error) {
+	rsp, err := c.GetCustomerCustomerIDOrdersOrderID(ctx, customerID, orderID, reqEditors...)
+	if err != nil {
+		return nil, err
+	}
+	return ParseGetCustomerCustomerIDOrdersOrderIDResponse(rsp)
+}
+
+// ParsePostCustomerCustomerIDOrdersResponse parses an HTTP response from a PostCustomerCustomerIDOrdersWithResponse call
+func ParsePostCustomerCustomerIDOrdersResponse(rsp *http.Response) (*PostCustomerCustomerIDOrdersResponse, error) {
+	bodyBytes, err := io.ReadAll(rsp.Body)
+	defer func() { _ = rsp.Body.Close() }()
+	if err != nil {
+		return nil, err
+	}
+
+	response := &PostCustomerCustomerIDOrdersResponse{
+		Body:         bodyBytes,
+		HTTPResponse: rsp,
+	}
+
+	switch {
+	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
+		var dest Order
+		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
+			return nil, err
+		}
+		response.JSON200 = &dest
+
+	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
+		var dest Error
+		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
+			return nil, err
+		}
+		response.JSONDefault = &dest
+
+	}
+
+	return response, nil
+}
+
+// ParseGetCustomerCustomerIDOrdersOrderIDResponse parses an HTTP response from a GetCustomerCustomerIDOrdersOrderIDWithResponse call
+func ParseGetCustomerCustomerIDOrdersOrderIDResponse(rsp *http.Response) (*GetCustomerCustomerIDOrdersOrderIDResponse, error) {
+	bodyBytes, err := io.ReadAll(rsp.Body)
+	defer func() { _ = rsp.Body.Close() }()
+	if err != nil {
+		return nil, err
+	}
+
+	response := &GetCustomerCustomerIDOrdersOrderIDResponse{
+		Body:         bodyBytes,
+		HTTPResponse: rsp,
+	}
+
+	switch {
+	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
+		var dest Order
+		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
+			return nil, err
+		}
+		response.JSON200 = &dest
+
+	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
+		var dest Error
+		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
+			return nil, err
+		}
+		response.JSONDefault = &dest
+
+	}
+
+	return response, nil
+}
diff --git a/internal/common/client/order/openapi_types.gen.go b/internal/common/client/order/openapi_types.gen.go
new file mode 100644
index 0000000..01e8118
--- /dev/null
+++ b/internal/common/client/order/openapi_types.gen.go
@@ -0,0 +1,41 @@
+// Package order provides primitives to interact with the openapi HTTP API.
+//
+// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
+package order
+
+// CreateOrderRequest defines model for CreateOrderRequest.
+type CreateOrderRequest struct {
+	CustomerID string             `json:"customerID"`
+	Items      []ItemWithQuantity `json:"items"`
+}
+
+// Error defines model for Error.
+type Error struct {
+	Message *string `json:"message,omitempty"`
+}
+
+// Item defines model for Item.
+type Item struct {
+	Id       *string `json:"id,omitempty"`
+	Name     *string `json:"name,omitempty"`
+	PriceID  *string `json:"priceID,omitempty"`
+	Quantity *int32  `json:"quantity,omitempty"`
+}
+
+// ItemWithQuantity defines model for ItemWithQuantity.
+type ItemWithQuantity struct {
+	Id       *string `json:"id,omitempty"`
+	Quantity *int32  `json:"quantity,omitempty"`
+}
+
+// Order defines model for Order.
+type Order struct {
+	CustomerID  *string `json:"customerID,omitempty"`
+	Id          *string `json:"id,omitempty"`
+	Items       *[]Item `json:"items,omitempty"`
+	PaymentLink *string `json:"paymentLink,omitempty"`
+	Status      *string `json:"status,omitempty"`
+}
+
+// PostCustomerCustomerIDOrdersJSONRequestBody defines body for PostCustomerCustomerIDOrders for application/json ContentType.
+type PostCustomerCustomerIDOrdersJSONRequestBody = CreateOrderRequest
diff --git a/internal/common/config/global.yaml b/internal/common/config/global.yaml
new file mode 100644
index 0000000..4e46835
--- /dev/null
+++ b/internal/common/config/global.yaml
@@ -0,0 +1,4 @@
+order:
+  service-name: order
+  server-to-run: http
+  http-addr: 127.0.0.1:8282
diff --git a/internal/common/config/viper.go b/internal/common/config/viper.go
new file mode 100644
index 0000000..103d697
--- /dev/null
+++ b/internal/common/config/viper.go
@@ -0,0 +1,11 @@
+package config
+
+import "github.com/spf13/viper"
+
+func NewViperConfig() error {
+	viper.SetConfigName("global")
+	viper.SetConfigType("yaml")
+	viper.AddConfigPath("../common/config")
+	viper.AutomaticEnv()
+	return viper.ReadInConfig()
+}
diff --git a/internal/common/genproto/orderpb/order.pb.go b/internal/common/genproto/orderpb/order.pb.go
new file mode 100644
index 0000000..f016711
--- /dev/null
+++ b/internal/common/genproto/orderpb/order.pb.go
@@ -0,0 +1,508 @@
+// Code generated by protoc-gen-go. DO NOT EDIT.
+// versions:
+// 	protoc-gen-go v1.28.1
+// 	protoc        v3.21.12
+// source: orderpb/order.proto
+
+package orderpb
+
+import (
+	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
+	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
+	emptypb "google.golang.org/protobuf/types/known/emptypb"
+	reflect "reflect"
+	sync "sync"
+)
+
+const (
+	// Verify that this generated code is sufficiently up-to-date.
+	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
+	// Verify that runtime/protoimpl is sufficiently up-to-date.
+	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
+)
+
+type CreateOrderRequest struct {
+	state         protoimpl.MessageState
+	sizeCache     protoimpl.SizeCache
+	unknownFields protoimpl.UnknownFields
+
+	CustomerID string              `protobuf:"bytes,1,opt,name=CustomerID,proto3" json:"CustomerID,omitempty"`
+	Items      []*ItemWithQuantity `protobuf:"bytes,2,rep,name=Items,proto3" json:"Items,omitempty"`
+}
+
+func (x *CreateOrderRequest) Reset() {
+	*x = CreateOrderRequest{}
+	if protoimpl.UnsafeEnabled {
+		mi := &file_orderpb_order_proto_msgTypes[0]
+		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
+		ms.StoreMessageInfo(mi)
+	}
+}
+
+func (x *CreateOrderRequest) String() string {
+	return protoimpl.X.MessageStringOf(x)
+}
+
+func (*CreateOrderRequest) ProtoMessage() {}
+
+func (x *CreateOrderRequest) ProtoReflect() protoreflect.Message {
+	mi := &file_orderpb_order_proto_msgTypes[0]
+	if protoimpl.UnsafeEnabled && x != nil {
+		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
+		if ms.LoadMessageInfo() == nil {
+			ms.StoreMessageInfo(mi)
+		}
+		return ms
+	}
+	return mi.MessageOf(x)
+}
+
+// Deprecated: Use CreateOrderRequest.ProtoReflect.Descriptor instead.
+func (*CreateOrderRequest) Descriptor() ([]byte, []int) {
+	return file_orderpb_order_proto_rawDescGZIP(), []int{0}
+}
+
+func (x *CreateOrderRequest) GetCustomerID() string {
+	if x != nil {
+		return x.CustomerID
+	}
+	return ""
+}
+
+func (x *CreateOrderRequest) GetItems() []*ItemWithQuantity {
+	if x != nil {
+		return x.Items
+	}
+	return nil
+}
+
+type GetOrderRequest struct {
+	state         protoimpl.MessageState
+	sizeCache     protoimpl.SizeCache
+	unknownFields protoimpl.UnknownFields
+
+	OrderID    string `protobuf:"bytes,1,opt,name=OrderID,proto3" json:"OrderID,omitempty"`
+	CustomerID string `protobuf:"bytes,2,opt,name=CustomerID,proto3" json:"CustomerID,omitempty"`
+}
+
+func (x *GetOrderRequest) Reset() {
+	*x = GetOrderRequest{}
+	if protoimpl.UnsafeEnabled {
+		mi := &file_orderpb_order_proto_msgTypes[1]
+		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
+		ms.StoreMessageInfo(mi)
+	}
+}
+
+func (x *GetOrderRequest) String() string {
+	return protoimpl.X.MessageStringOf(x)
+}
+
+func (*GetOrderRequest) ProtoMessage() {}
+
+func (x *GetOrderRequest) ProtoReflect() protoreflect.Message {
+	mi := &file_orderpb_order_proto_msgTypes[1]
+	if protoimpl.UnsafeEnabled && x != nil {
+		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
+		if ms.LoadMessageInfo() == nil {
+			ms.StoreMessageInfo(mi)
+		}
+		return ms
+	}
+	return mi.MessageOf(x)
+}
+
+// Deprecated: Use GetOrderRequest.ProtoReflect.Descriptor instead.
+func (*GetOrderRequest) Descriptor() ([]byte, []int) {
+	return file_orderpb_order_proto_rawDescGZIP(), []int{1}
+}
+
+func (x *GetOrderRequest) GetOrderID() string {
+	if x != nil {
+		return x.OrderID
+	}
+	return ""
+}
+
+func (x *GetOrderRequest) GetCustomerID() string {
+	if x != nil {
+		return x.CustomerID
+	}
+	return ""
+}
+
+type ItemWithQuantity struct {
+	state         protoimpl.MessageState
+	sizeCache     protoimpl.SizeCache
+	unknownFields protoimpl.UnknownFields
+
+	ID       string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
+	Quantity int32  `protobuf:"varint,2,opt,name=Quantity,proto3" json:"Quantity,omitempty"`
+}
+
+func (x *ItemWithQuantity) Reset() {
+	*x = ItemWithQuantity{}
+	if protoimpl.UnsafeEnabled {
+		mi := &file_orderpb_order_proto_msgTypes[2]
+		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
+		ms.StoreMessageInfo(mi)
+	}
+}
+
+func (x *ItemWithQuantity) String() string {
+	return protoimpl.X.MessageStringOf(x)
+}
+
+func (*ItemWithQuantity) ProtoMessage() {}
+
+func (x *ItemWithQuantity) ProtoReflect() protoreflect.Message {
+	mi := &file_orderpb_order_proto_msgTypes[2]
+	if protoimpl.UnsafeEnabled && x != nil {
+		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
+		if ms.LoadMessageInfo() == nil {
+			ms.StoreMessageInfo(mi)
+		}
+		return ms
+	}
+	return mi.MessageOf(x)
+}
+
+// Deprecated: Use ItemWithQuantity.ProtoReflect.Descriptor instead.
+func (*ItemWithQuantity) Descriptor() ([]byte, []int) {
+	return file_orderpb_order_proto_rawDescGZIP(), []int{2}
+}
+
+func (x *ItemWithQuantity) GetID() string {
+	if x != nil {
+		return x.ID
+	}
+	return ""
+}
+
+func (x *ItemWithQuantity) GetQuantity() int32 {
+	if x != nil {
+		return x.Quantity
+	}
+	return 0
+}
+
+type Item struct {
+	state         protoimpl.MessageState
+	sizeCache     protoimpl.SizeCache
+	unknownFields protoimpl.UnknownFields
+
+	ID       string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
+	Name     string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
+	Quantity int32  `protobuf:"varint,3,opt,name=Quantity,proto3" json:"Quantity,omitempty"`
+	PriceID  string `protobuf:"bytes,4,opt,name=PriceID,proto3" json:"PriceID,omitempty"`
+}
+
+func (x *Item) Reset() {
+	*x = Item{}
+	if protoimpl.UnsafeEnabled {
+		mi := &file_orderpb_order_proto_msgTypes[3]
+		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
+		ms.StoreMessageInfo(mi)
+	}
+}
+
+func (x *Item) String() string {
+	return protoimpl.X.MessageStringOf(x)
+}
+
+func (*Item) ProtoMessage() {}
+
+func (x *Item) ProtoReflect() protoreflect.Message {
+	mi := &file_orderpb_order_proto_msgTypes[3]
+	if protoimpl.UnsafeEnabled && x != nil {
+		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
+		if ms.LoadMessageInfo() == nil {
+			ms.StoreMessageInfo(mi)
+		}
+		return ms
+	}
+	return mi.MessageOf(x)
+}
+
+// Deprecated: Use Item.ProtoReflect.Descriptor instead.
+func (*Item) Descriptor() ([]byte, []int) {
+	return file_orderpb_order_proto_rawDescGZIP(), []int{3}
+}
+
+func (x *Item) GetID() string {
+	if x != nil {
+		return x.ID
+	}
+	return ""
+}
+
+func (x *Item) GetName() string {
+	if x != nil {
+		return x.Name
+	}
+	return ""
+}
+
+func (x *Item) GetQuantity() int32 {
+	if x != nil {
+		return x.Quantity
+	}
+	return 0
+}
+
+func (x *Item) GetPriceID() string {
+	if x != nil {
+		return x.PriceID
+	}
+	return ""
+}
+
+type Order struct {
+	state         protoimpl.MessageState
+	sizeCache     protoimpl.SizeCache
+	unknownFields protoimpl.UnknownFields
+
+	ID         string  `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
+	CustomerID string  `protobuf:"bytes,2,opt,name=CustomerID,proto3" json:"CustomerID,omitempty"`
+	Status     string  `protobuf:"bytes,3,opt,name=Status,proto3" json:"Status,omitempty"`
+	Items      []*Item `protobuf:"bytes,4,rep,name=Items,proto3" json:"Items,omitempty"`
+}
+
+func (x *Order) Reset() {
+	*x = Order{}
+	if protoimpl.UnsafeEnabled {
+		mi := &file_orderpb_order_proto_msgTypes[4]
+		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
+		ms.StoreMessageInfo(mi)
+	}
+}
+
+func (x *Order) String() string {
+	return protoimpl.X.MessageStringOf(x)
+}
+
+func (*Order) ProtoMessage() {}
+
+func (x *Order) ProtoReflect() protoreflect.Message {
+	mi := &file_orderpb_order_proto_msgTypes[4]
+	if protoimpl.UnsafeEnabled && x != nil {
+		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
+		if ms.LoadMessageInfo() == nil {
+			ms.StoreMessageInfo(mi)
+		}
+		return ms
+	}
+	return mi.MessageOf(x)
+}
+
+// Deprecated: Use Order.ProtoReflect.Descriptor instead.
+func (*Order) Descriptor() ([]byte, []int) {
+	return file_orderpb_order_proto_rawDescGZIP(), []int{4}
+}
+
+func (x *Order) GetID() string {
+	if x != nil {
+		return x.ID
+	}
+	return ""
+}
+
+func (x *Order) GetCustomerID() string {
+	if x != nil {
+		return x.CustomerID
+	}
+	return ""
+}
+
+func (x *Order) GetStatus() string {
+	if x != nil {
+		return x.Status
+	}
+	return ""
+}
+
+func (x *Order) GetItems() []*Item {
+	if x != nil {
+		return x.Items
+	}
+	return nil
+}
+
+var File_orderpb_order_proto protoreflect.FileDescriptor
+
+var file_orderpb_order_proto_rawDesc = []byte{
+	0x0a, 0x13, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e,
+	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x70, 0x62, 0x1a, 0x1b,
+	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
+	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x65, 0x0a, 0x12, 0x43,
+	0x72, 0x65, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
+	0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x44, 0x18,
+	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49,
+	0x44, 0x12, 0x2f, 0x0a, 0x05, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
+	0x32, 0x19, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x57,
+	0x69, 0x74, 0x68, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x05, 0x49, 0x74, 0x65,
+	0x6d, 0x73, 0x22, 0x4b, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65,
+	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x44,
+	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x44, 0x12,
+	0x1e, 0x0a, 0x0a, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20,
+	0x01, 0x28, 0x09, 0x52, 0x0a, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x44, 0x22,
+	0x3e, 0x0a, 0x10, 0x49, 0x74, 0x65, 0x6d, 0x57, 0x69, 0x74, 0x68, 0x51, 0x75, 0x61, 0x6e, 0x74,
+	0x69, 0x74, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
+	0x02, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18,
+	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x22,
+	0x60, 0x0a, 0x04, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20,
+	0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18,
+	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x51,
+	0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x51,
+	0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x50, 0x72, 0x69, 0x63, 0x65,
+	0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x50, 0x72, 0x69, 0x63, 0x65, 0x49,
+	0x44, 0x22, 0x74, 0x0a, 0x05, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44,
+	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x75,
+	0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
+	0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x74,
+	0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74,
+	0x75, 0x73, 0x12, 0x23, 0x0a, 0x05, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28,
+	0x0b, 0x32, 0x0d, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x49, 0x74, 0x65, 0x6d,
+	0x52, 0x05, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x32, 0xbf, 0x01, 0x0a, 0x0c, 0x4f, 0x72, 0x64, 0x65,
+	0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x42, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61,
+	0x74, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x1b, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x70,
+	0x62, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71,
+	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
+	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x34, 0x0a, 0x08,
+	0x47, 0x65, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x18, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72,
+	0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
+	0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x4f, 0x72, 0x64,
+	0x65, 0x72, 0x12, 0x35, 0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x64, 0x65,
+	0x72, 0x12, 0x0e, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x4f, 0x72, 0x64, 0x65,
+	0x72, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
+	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x44, 0x5a, 0x42, 0x67, 0x69, 0x74,
+	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65,
+	0x7a, 0x7a, 0x30, 0x30, 0x2f, 0x67, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2d, 0x76, 0x32, 0x2f, 0x69,
+	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x67,
+	0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x70, 0x62, 0x62,
+	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
+}
+
+var (
+	file_orderpb_order_proto_rawDescOnce sync.Once
+	file_orderpb_order_proto_rawDescData = file_orderpb_order_proto_rawDesc
+)
+
+func file_orderpb_order_proto_rawDescGZIP() []byte {
+	file_orderpb_order_proto_rawDescOnce.Do(func() {
+		file_orderpb_order_proto_rawDescData = protoimpl.X.CompressGZIP(file_orderpb_order_proto_rawDescData)
+	})
+	return file_orderpb_order_proto_rawDescData
+}
+
+var file_orderpb_order_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
+var file_orderpb_order_proto_goTypes = []interface{}{
+	(*CreateOrderRequest)(nil), // 0: orderpb.CreateOrderRequest
+	(*GetOrderRequest)(nil),    // 1: orderpb.GetOrderRequest
+	(*ItemWithQuantity)(nil),   // 2: orderpb.ItemWithQuantity
+	(*Item)(nil),               // 3: orderpb.Item
+	(*Order)(nil),              // 4: orderpb.Order
+	(*emptypb.Empty)(nil),      // 5: google.protobuf.Empty
+}
+var file_orderpb_order_proto_depIdxs = []int32{
+	2, // 0: orderpb.CreateOrderRequest.Items:type_name -> orderpb.ItemWithQuantity
+	3, // 1: orderpb.Order.Items:type_name -> orderpb.Item
+	0, // 2: orderpb.OrderService.CreateOrder:input_type -> orderpb.CreateOrderRequest
+	1, // 3: orderpb.OrderService.GetOrder:input_type -> orderpb.GetOrderRequest
+	4, // 4: orderpb.OrderService.UpdateOrder:input_type -> orderpb.Order
+	5, // 5: orderpb.OrderService.CreateOrder:output_type -> google.protobuf.Empty
+	4, // 6: orderpb.OrderService.GetOrder:output_type -> orderpb.Order
+	5, // 7: orderpb.OrderService.UpdateOrder:output_type -> google.protobuf.Empty
+	5, // [5:8] is the sub-list for method output_type
+	2, // [2:5] is the sub-list for method input_type
+	2, // [2:2] is the sub-list for extension type_name
+	2, // [2:2] is the sub-list for extension extendee
+	0, // [0:2] is the sub-list for field type_name
+}
+
+func init() { file_orderpb_order_proto_init() }
+func file_orderpb_order_proto_init() {
+	if File_orderpb_order_proto != nil {
+		return
+	}
+	if !protoimpl.UnsafeEnabled {
+		file_orderpb_order_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
+			switch v := v.(*CreateOrderRequest); i {
+			case 0:
+				return &v.state
+			case 1:
+				return &v.sizeCache
+			case 2:
+				return &v.unknownFields
+			default:
+				return nil
+			}
+		}
+		file_orderpb_order_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
+			switch v := v.(*GetOrderRequest); i {
+			case 0:
+				return &v.state
+			case 1:
+				return &v.sizeCache
+			case 2:
+				return &v.unknownFields
+			default:
+				return nil
+			}
+		}
+		file_orderpb_order_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
+			switch v := v.(*ItemWithQuantity); i {
+			case 0:
+				return &v.state
+			case 1:
+				return &v.sizeCache
+			case 2:
+				return &v.unknownFields
+			default:
+				return nil
+			}
+		}
+		file_orderpb_order_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
+			switch v := v.(*Item); i {
+			case 0:
+				return &v.state
+			case 1:
+				return &v.sizeCache
+			case 2:
+				return &v.unknownFields
+			default:
+				return nil
+			}
+		}
+		file_orderpb_order_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
+			switch v := v.(*Order); i {
+			case 0:
+				return &v.state
+			case 1:
+				return &v.sizeCache
+			case 2:
+				return &v.unknownFields
+			default:
+				return nil
+			}
+		}
+	}
+	type x struct{}
+	out := protoimpl.TypeBuilder{
+		File: protoimpl.DescBuilder{
+			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
+			RawDescriptor: file_orderpb_order_proto_rawDesc,
+			NumEnums:      0,
+			NumMessages:   5,
+			NumExtensions: 0,
+			NumServices:   1,
+		},
+		GoTypes:           file_orderpb_order_proto_goTypes,
+		DependencyIndexes: file_orderpb_order_proto_depIdxs,
+		MessageInfos:      file_orderpb_order_proto_msgTypes,
+	}.Build()
+	File_orderpb_order_proto = out.File
+	file_orderpb_order_proto_rawDesc = nil
+	file_orderpb_order_proto_goTypes = nil
+	file_orderpb_order_proto_depIdxs = nil
+}
diff --git a/internal/common/genproto/orderpb/order_grpc.pb.go b/internal/common/genproto/orderpb/order_grpc.pb.go
new file mode 100644
index 0000000..f18a631
--- /dev/null
+++ b/internal/common/genproto/orderpb/order_grpc.pb.go
@@ -0,0 +1,176 @@
+// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
+// versions:
+// - protoc-gen-go-grpc v1.2.0
+// - protoc             v3.21.12
+// source: orderpb/order.proto
+
+package orderpb
+
+import (
+	context "context"
+	grpc "google.golang.org/grpc"
+	codes "google.golang.org/grpc/codes"
+	status "google.golang.org/grpc/status"
+	emptypb "google.golang.org/protobuf/types/known/emptypb"
+)
+
+// This is a compile-time assertion to ensure that this generated file
+// is compatible with the grpc package it is being compiled against.
+// Requires gRPC-Go v1.32.0 or later.
+const _ = grpc.SupportPackageIsVersion7
+
+// OrderServiceClient is the client API for OrderService service.
+//
+// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
+type OrderServiceClient interface {
+	CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
+	GetOrder(ctx context.Context, in *GetOrderRequest, opts ...grpc.CallOption) (*Order, error)
+	UpdateOrder(ctx context.Context, in *Order, opts ...grpc.CallOption) (*emptypb.Empty, error)
+}
+
+type orderServiceClient struct {
+	cc grpc.ClientConnInterface
+}
+
+func NewOrderServiceClient(cc grpc.ClientConnInterface) OrderServiceClient {
+	return &orderServiceClient{cc}
+}
+
+func (c *orderServiceClient) CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
+	out := new(emptypb.Empty)
+	err := c.cc.Invoke(ctx, "/orderpb.OrderService/CreateOrder", in, out, opts...)
+	if err != nil {
+		return nil, err
+	}
+	return out, nil
+}
+
+func (c *orderServiceClient) GetOrder(ctx context.Context, in *GetOrderRequest, opts ...grpc.CallOption) (*Order, error) {
+	out := new(Order)
+	err := c.cc.Invoke(ctx, "/orderpb.OrderService/GetOrder", in, out, opts...)
+	if err != nil {
+		return nil, err
+	}
+	return out, nil
+}
+
+func (c *orderServiceClient) UpdateOrder(ctx context.Context, in *Order, opts ...grpc.CallOption) (*emptypb.Empty, error) {
+	out := new(emptypb.Empty)
+	err := c.cc.Invoke(ctx, "/orderpb.OrderService/UpdateOrder", in, out, opts...)
+	if err != nil {
+		return nil, err
+	}
+	return out, nil
+}
+
+// OrderServiceServer is the server API for OrderService service.
+// All implementations should embed UnimplementedOrderServiceServer
+// for forward compatibility
+type OrderServiceServer interface {
+	CreateOrder(context.Context, *CreateOrderRequest) (*emptypb.Empty, error)
+	GetOrder(context.Context, *GetOrderRequest) (*Order, error)
+	UpdateOrder(context.Context, *Order) (*emptypb.Empty, error)
+}
+
+// UnimplementedOrderServiceServer should be embedded to have forward compatible implementations.
+type UnimplementedOrderServiceServer struct {
+}
+
+func (UnimplementedOrderServiceServer) CreateOrder(context.Context, *CreateOrderRequest) (*emptypb.Empty, error) {
+	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
+}
+func (UnimplementedOrderServiceServer) GetOrder(context.Context, *GetOrderRequest) (*Order, error) {
+	return nil, status.Errorf(codes.Unimplemented, "method GetOrder not implemented")
+}
+func (UnimplementedOrderServiceServer) UpdateOrder(context.Context, *Order) (*emptypb.Empty, error) {
+	return nil, status.Errorf(codes.Unimplemented, "method UpdateOrder not implemented")
+}
+
+// UnsafeOrderServiceServer may be embedded to opt out of forward compatibility for this service.
+// Use of this interface is not recommended, as added methods to OrderServiceServer will
+// result in compilation errors.
+type UnsafeOrderServiceServer interface {
+	mustEmbedUnimplementedOrderServiceServer()
+}
+
+func RegisterOrderServiceServer(s grpc.ServiceRegistrar, srv OrderServiceServer) {
+	s.RegisterService(&OrderService_ServiceDesc, srv)
+}
+
+func _OrderService_CreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
+	in := new(CreateOrderRequest)
+	if err := dec(in); err != nil {
+		return nil, err
+	}
+	if interceptor == nil {
+		return srv.(OrderServiceServer).CreateOrder(ctx, in)
+	}
+	info := &grpc.UnaryServerInfo{
+		Server:     srv,
+		FullMethod: "/orderpb.OrderService/CreateOrder",
+	}
+	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
+		return srv.(OrderServiceServer).CreateOrder(ctx, req.(*CreateOrderRequest))
+	}
+	return interceptor(ctx, in, info, handler)
+}
+
+func _OrderService_GetOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
+	in := new(GetOrderRequest)
+	if err := dec(in); err != nil {
+		return nil, err
+	}
+	if interceptor == nil {
+		return srv.(OrderServiceServer).GetOrder(ctx, in)
+	}
+	info := &grpc.UnaryServerInfo{
+		Server:     srv,
+		FullMethod: "/orderpb.OrderService/GetOrder",
+	}
+	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
+		return srv.(OrderServiceServer).GetOrder(ctx, req.(*GetOrderRequest))
+	}
+	return interceptor(ctx, in, info, handler)
+}
+
+func _OrderService_UpdateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
+	in := new(Order)
+	if err := dec(in); err != nil {
+		return nil, err
+	}
+	if interceptor == nil {
+		return srv.(OrderServiceServer).UpdateOrder(ctx, in)
+	}
+	info := &grpc.UnaryServerInfo{
+		Server:     srv,
+		FullMethod: "/orderpb.OrderService/UpdateOrder",
+	}
+	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
+		return srv.(OrderServiceServer).UpdateOrder(ctx, req.(*Order))
+	}
+	return interceptor(ctx, in, info, handler)
+}
+
+// OrderService_ServiceDesc is the grpc.ServiceDesc for OrderService service.
+// It's only intended for direct use with grpc.RegisterService,
+// and not to be introspected or modified (even as a copy)
+var OrderService_ServiceDesc = grpc.ServiceDesc{
+	ServiceName: "orderpb.OrderService",
+	HandlerType: (*OrderServiceServer)(nil),
+	Methods: []grpc.MethodDesc{
+		{
+			MethodName: "CreateOrder",
+			Handler:    _OrderService_CreateOrder_Handler,
+		},
+		{
+			MethodName: "GetOrder",
+			Handler:    _OrderService_GetOrder_Handler,
+		},
+		{
+			MethodName: "UpdateOrder",
+			Handler:    _OrderService_UpdateOrder_Handler,
+		},
+	},
+	Streams:  []grpc.StreamDesc{},
+	Metadata: "orderpb/order.proto",
+}
diff --git a/internal/order/go.mod b/internal/common/go.mod
similarity index 65%
copy from internal/order/go.mod
copy to internal/common/go.mod
index 9675054..861e1d8 100644
--- a/internal/order/go.mod
+++ b/internal/common/go.mod
@@ -1,13 +1,19 @@
-module github.com/ghost-yu/go_shop_second/order
+module github.com/ghost-yu/go_shop_second/common
 
 go 1.22.8
 
-replace github.com/ghost-yu/go_shop_second/common => ../common
-
-require github.com/ghost-yu/go_shop_second/common v0.0.0-00010101000000-000000000000
+require (
+	github.com/oapi-codegen/runtime v1.1.1
+	github.com/spf13/viper v1.19.0
+	google.golang.org/grpc v1.62.1
+	google.golang.org/protobuf v1.33.0
+)
 
 require (
+	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
 	github.com/fsnotify/fsnotify v1.7.0 // indirect
+	github.com/golang/protobuf v1.5.3 // indirect
+	github.com/google/uuid v1.6.0 // indirect
 	github.com/hashicorp/hcl v1.0.0 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
@@ -18,13 +24,14 @@ require (
 	github.com/spf13/afero v1.11.0 // indirect
 	github.com/spf13/cast v1.6.0 // indirect
 	github.com/spf13/pflag v1.0.5 // indirect
-	github.com/spf13/viper v1.19.0 // indirect
 	github.com/subosito/gotenv v1.6.0 // indirect
 	go.uber.org/atomic v1.9.0 // indirect
 	go.uber.org/multierr v1.9.0 // indirect
 	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
+	golang.org/x/net v0.23.0 // indirect
 	golang.org/x/sys v0.18.0 // indirect
 	golang.org/x/text v0.14.0 // indirect
+	google.golang.org/genproto/googleapis/rpc v0.0.0-20240314234333-6e1732d8331c // indirect
 	gopkg.in/ini.v1 v1.67.0 // indirect
 	gopkg.in/yaml.v3 v3.0.1 // indirect
 )
diff --git a/internal/order/go.sum b/internal/common/go.sum
similarity index 71%
copy from internal/order/go.sum
copy to internal/common/go.sum
index ec3cc1e..1a450fd 100644
--- a/internal/order/go.sum
+++ b/internal/common/go.sum
@@ -1,3 +1,7 @@
+github.com/RaveNoX/go-jsoncommentstrip v1.0.0/go.mod h1:78ihd09MekBnJnxpICcwzCMzGrKSKYe4AqU6PDYYpjk=
+github.com/apapsch/go-jsonmerge/v2 v2.0.0 h1:axGnT1gRIfimI7gJifB699GoE/oq+F2MU7Dml6nw9rQ=
+github.com/apapsch/go-jsonmerge/v2 v2.0.0/go.mod h1:lvDnEdqiQrp0O42VQGgmlKpxL1AP2+08jFMw88y4klk=
+github.com/bmatcuk/doublestar v1.1.1/go.mod h1:UD6OnuiIn0yFxxA2le/rnRU1G4RaI4UvFv1sNto9p6w=
 github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
 github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
 github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc h1:U9qPSI2PIWSS1VwoXQT9A3Wy9MM3WgvqSxFWenqJduM=
@@ -6,10 +10,17 @@ github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHk
 github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
-github.com/google/go-cmp v0.5.9 h1:O2Tfq5qg4qc4AmwVlvv0oLiVAGB7enBSJ2x2DqQFi38=
-github.com/google/go-cmp v0.5.9/go.mod h1:17dUlkBOakJ0+DkrSSNjCkIjxS6bF9zb3elmeNGIjoY=
+github.com/golang/protobuf v1.5.0/go.mod h1:FsONVRAS9T7sI+LIUmWTfcYkHO4aIWwzhcaSAoJOfIk=
+github.com/golang/protobuf v1.5.3 h1:KhyjKVUg7Usr/dYsdSqoFveMYd5ko72D+zANwlG1mmg=
+github.com/golang/protobuf v1.5.3/go.mod h1:XVQd3VNwM+JqD3oG2Ue2ip4fOMUkwXdXDdiuN0vRsmY=
+github.com/google/go-cmp v0.5.5/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
+github.com/google/go-cmp v0.6.0 h1:ofyhxvXcZhMsU5ulbFiLKl/XBFqE1GSq7atu8tAmTRI=
+github.com/google/go-cmp v0.6.0/go.mod h1:17dUlkBOakJ0+DkrSSNjCkIjxS6bF9zb3elmeNGIjoY=
+github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
+github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
 github.com/hashicorp/hcl v1.0.0 h1:0Anlzjpi4vEasTeNFn2mLJgTSwt0+6sfsiTG8qcWGx4=
 github.com/hashicorp/hcl v1.0.0/go.mod h1:E5yfLk+7swimpb2L/Alb/PJmXilQ/rhwaUYs4T20WEQ=
+github.com/juju/gnuflag v0.0.0-20171113085948-2ce1bb71843d/go.mod h1:2PavIy+JPciBPrBUjwbNvtwB6RQlve+hkpll6QSNmOE=
 github.com/kr/pretty v0.3.1 h1:flRD4NNwYAUpkphVc1HcthR4KEIFJ65n8Mw5qdRn3LE=
 github.com/kr/pretty v0.3.1/go.mod h1:hoEshYVHaxMs3cyo3Yncou5ZscifuDolrwPKZanG3xk=
 github.com/kr/text v0.2.0 h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=
@@ -18,6 +29,8 @@ github.com/magiconair/properties v1.8.7 h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0V
 github.com/magiconair/properties v1.8.7/go.mod h1:Dhd985XPs7jluiymwWYZ0G4Z61jb3vdS329zhj2hYo0=
 github.com/mitchellh/mapstructure v1.5.0 h1:jeMsZIYE/09sWLaz43PL7Gy6RuMjD2eJVyuac5Z2hdY=
 github.com/mitchellh/mapstructure v1.5.0/go.mod h1:bFUtVrKA4DC2yAKiSyO/QUcy7e+RRV2QTWOzhPopBRo=
+github.com/oapi-codegen/runtime v1.1.1 h1:EXLHh0DXIJnWhdRPN2w4MXAzFyE4CskzhNLUmtpMYro=
+github.com/oapi-codegen/runtime v1.1.1/go.mod h1:SK9X900oXmPWilYR5/WKPzt3Kqxn/uS/+lbpREv+eCg=
 github.com/pelletier/go-toml/v2 v2.2.2 h1:aYUidT7k73Pcl9nb2gScu7NSrKCSHIDE89b3+6Wq+LM=
 github.com/pelletier/go-toml/v2 v2.2.2/go.mod h1:1t835xjRzz80PqgE6HHgN2JOsmgYu/h4qDAS4n929Rs=
 github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
@@ -39,6 +52,7 @@ github.com/spf13/pflag v1.0.5 h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=
 github.com/spf13/pflag v1.0.5/go.mod h1:McXfInJRrz4CZXVZOBLb0bTZqETkiAhM9Iw0y3An2Bg=
 github.com/spf13/viper v1.19.0 h1:RWq5SEjt8o25SROyN3z2OrDB9l7RPd3lwTWU8EcEdcI=
 github.com/spf13/viper v1.19.0/go.mod h1:GQUN9bilAbhU/jgc1bKs99f/suXKeUMct8Adx5+Ntkg=
+github.com/spkg/bom v0.0.0-20160624110644-59b7046e48ad/go.mod h1:qLr4V1qq6nMqFKkMo8ZTx3f+BZEkzsRUY10Xsm2mwU0=
 github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
 github.com/stretchr/objx v0.5.0/go.mod h1:Yh+to48EsGEfYuaHDzXPcE3xhTkx73EhmCGUpEOglKo=
@@ -57,10 +71,21 @@ go.uber.org/multierr v1.9.0 h1:7fIwc/ZtS0q++VgcfqFDxSBZVv/Xo49/SYnDFupUwlI=
 go.uber.org/multierr v1.9.0/go.mod h1:X2jQV1h+kxSjClGpnseKVIxpmcjrj7MNnI0bnlfKTVQ=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9 h1:GoHiUyI/Tp2nVkLI2mCxVkOjsbSXD66ic0XW0js0R9g=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9/go.mod h1:S2oDrQGGwySpoQPVqRShND87VCbxmc6bL1Yd2oYrm6k=
+golang.org/x/net v0.23.0 h1:7EYJ93RZ9vYSZAIb2x3lnuvqO5zneoD6IvWjuhfxjTs=
+golang.org/x/net v0.23.0/go.mod h1:JKghWKKOSdJwpW2GEx0Ja7fmaKnMsbu+MWVZTokSYmg=
 golang.org/x/sys v0.18.0 h1:DBdB3niSjOA/O0blCZBqDefyWNYveAYMNF1Wum0DYQ4=
 golang.org/x/sys v0.18.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
 golang.org/x/text v0.14.0 h1:ScX5w1eTa3QqT8oi6+ziP7dTV1S2+ALU0bI+0zXKWiQ=
 golang.org/x/text v0.14.0/go.mod h1:18ZOQIKpY8NJVqYksKHtTdi31H5itFRjB5/qKTNYzSU=
+golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
+google.golang.org/genproto/googleapis/rpc v0.0.0-20240314234333-6e1732d8331c h1:lfpJ/2rWPa/kJgxyyXM8PrNnfCzcmxJ265mADgwmvLI=
+google.golang.org/genproto/googleapis/rpc v0.0.0-20240314234333-6e1732d8331c/go.mod h1:WtryC6hu0hhx87FDGxWCDptyssuo68sk10vYjF+T9fY=
+google.golang.org/grpc v1.62.1 h1:B4n+nfKzOICUXMgyrNd19h/I9oH0L1pizfk1d4zSgTk=
+google.golang.org/grpc v1.62.1/go.mod h1:IWTG0VlJLCh1SkC58F7np9ka9mx/WNkjl4PGJaiq+QE=
+google.golang.org/protobuf v1.26.0-rc.1/go.mod h1:jlhhOSvTdKEhbULTjvd4ARK9grFBp09yW+WbY/TyQbw=
+google.golang.org/protobuf v1.26.0/go.mod h1:9q0QmTI4eRPtz6boOQmLYwt+qCgq0jsYwAQnmE0givc=
+google.golang.org/protobuf v1.33.0 h1:uNO2rsAINq/JlFpSdYEKIZ0uKD/R9cpdv0T+yoGwGmI=
+google.golang.org/protobuf v1.33.0/go.mod h1:c6P6GXX6sHbq/GpV6MGZEdwhWPcYBgnhAHhKbcUYpos=
 gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 h1:YR8cESwS4TdDjEe65xsg0ogRM/Nc3DYOhEAlW+xobZo=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
diff --git a/internal/kitchen/go.mod b/internal/kitchen/go.mod
new file mode 100644
index 0000000..5449764
--- /dev/null
+++ b/internal/kitchen/go.mod
@@ -0,0 +1,7 @@
+module github.com/ghost-yu/go_shop_second/kitchen
+
+go 1.22.8
+
+replace (
+	github.com/ghost-yu/go_shop_second/common => ../common
+)
\ No newline at end of file
diff --git a/internal/order/go.mod b/internal/order/go.mod
index 9675054..898b98b 100644
--- a/internal/order/go.mod
+++ b/internal/order/go.mod
@@ -4,13 +4,36 @@ go 1.22.8
 
 replace github.com/ghost-yu/go_shop_second/common => ../common
 
-require github.com/ghost-yu/go_shop_second/common v0.0.0-00010101000000-000000000000
+require (
+	github.com/ghost-yu/go_shop_second/common v0.0.0-00010101000000-000000000000
+	github.com/gin-gonic/gin v1.10.0
+	github.com/oapi-codegen/runtime v1.1.1
+	github.com/spf13/viper v1.19.0
+)
 
 require (
+	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
+	github.com/bytedance/sonic v1.11.6 // indirect
+	github.com/bytedance/sonic/loader v0.1.1 // indirect
+	github.com/cloudwego/base64x v0.1.4 // indirect
+	github.com/cloudwego/iasm v0.2.0 // indirect
 	github.com/fsnotify/fsnotify v1.7.0 // indirect
+	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
+	github.com/gin-contrib/sse v0.1.0 // indirect
+	github.com/go-playground/locales v0.14.1 // indirect
+	github.com/go-playground/universal-translator v0.18.1 // indirect
+	github.com/go-playground/validator/v10 v10.20.0 // indirect
+	github.com/goccy/go-json v0.10.2 // indirect
+	github.com/google/uuid v1.6.0 // indirect
 	github.com/hashicorp/hcl v1.0.0 // indirect
+	github.com/json-iterator/go v1.1.12 // indirect
+	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
+	github.com/leodido/go-urn v1.4.0 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
+	github.com/mattn/go-isatty v0.0.20 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
+	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
+	github.com/modern-go/reflect2 v1.0.2 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
@@ -18,13 +41,18 @@ require (
 	github.com/spf13/afero v1.11.0 // indirect
 	github.com/spf13/cast v1.6.0 // indirect
 	github.com/spf13/pflag v1.0.5 // indirect
-	github.com/spf13/viper v1.19.0 // indirect
 	github.com/subosito/gotenv v1.6.0 // indirect
+	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
+	github.com/ugorji/go/codec v1.2.12 // indirect
 	go.uber.org/atomic v1.9.0 // indirect
 	go.uber.org/multierr v1.9.0 // indirect
+	golang.org/x/arch v0.8.0 // indirect
+	golang.org/x/crypto v0.23.0 // indirect
 	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
-	golang.org/x/sys v0.18.0 // indirect
-	golang.org/x/text v0.14.0 // indirect
+	golang.org/x/net v0.25.0 // indirect
+	golang.org/x/sys v0.20.0 // indirect
+	golang.org/x/text v0.15.0 // indirect
+	google.golang.org/protobuf v1.34.1 // indirect
 	gopkg.in/ini.v1 v1.67.0 // indirect
 	gopkg.in/yaml.v3 v3.0.1 // indirect
 )
diff --git a/internal/order/go.sum b/internal/order/go.sum
index ec3cc1e..73bf9dd 100644
--- a/internal/order/go.sum
+++ b/internal/order/go.sum
@@ -1,3 +1,15 @@
+github.com/RaveNoX/go-jsoncommentstrip v1.0.0/go.mod h1:78ihd09MekBnJnxpICcwzCMzGrKSKYe4AqU6PDYYpjk=
+github.com/apapsch/go-jsonmerge/v2 v2.0.0 h1:axGnT1gRIfimI7gJifB699GoE/oq+F2MU7Dml6nw9rQ=
+github.com/apapsch/go-jsonmerge/v2 v2.0.0/go.mod h1:lvDnEdqiQrp0O42VQGgmlKpxL1AP2+08jFMw88y4klk=
+github.com/bmatcuk/doublestar v1.1.1/go.mod h1:UD6OnuiIn0yFxxA2le/rnRU1G4RaI4UvFv1sNto9p6w=
+github.com/bytedance/sonic v1.11.6 h1:oUp34TzMlL+OY1OUWxHqsdkgC/Zfc85zGqw9siXjrc0=
+github.com/bytedance/sonic v1.11.6/go.mod h1:LysEHSvpvDySVdC2f87zGWf6CIKJcAvqab1ZaiQtds4=
+github.com/bytedance/sonic/loader v0.1.1 h1:c+e5Pt1k/cy5wMveRDyk2X4B9hF4g7an8N3zCYjJFNM=
+github.com/bytedance/sonic/loader v0.1.1/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
+github.com/cloudwego/base64x v0.1.4 h1:jwCgWpFanWmN8xoIUHa2rtzmkd5J2plF/dnLS6Xd/0Y=
+github.com/cloudwego/base64x v0.1.4/go.mod h1:0zlkT4Wn5C6NdauXdJRhSKRlJvmclQ1hhJgA0rcu/8w=
+github.com/cloudwego/iasm v0.2.0 h1:1KNIy1I1H9hNNFEEH3DVnI4UujN+1zjpuk6gwHLTssg=
+github.com/cloudwego/iasm v0.2.0/go.mod h1:8rXZaNYT2n95jn+zTI1sDr+IgcD2GVs0nlbbQPiEFhY=
 github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
 github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
 github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc h1:U9qPSI2PIWSS1VwoXQT9A3Wy9MM3WgvqSxFWenqJduM=
@@ -6,18 +18,56 @@ github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHk
 github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
 github.com/fsnotify/fsnotify v1.7.0/go.mod h1:40Bi/Hjc2AVfZrqy+aj+yEI+/bRxZnMJyTJwOpGvigM=
+github.com/gabriel-vasile/mimetype v1.4.3 h1:in2uUcidCuFcDKtdcBxlR0rJ1+fsokWf+uqxgUFjbI0=
+github.com/gabriel-vasile/mimetype v1.4.3/go.mod h1:d8uq/6HKRL6CGdk+aubisF/M5GcPfT7nKyLpA0lbSSk=
+github.com/gin-contrib/sse v0.1.0 h1:Y/yl/+YNO8GZSjAhjMsSuLt29uWRFHdHYUb5lYOV9qE=
+github.com/gin-contrib/sse v0.1.0/go.mod h1:RHrZQHXnP2xjPF+u1gW/2HnVO7nvIa9PG3Gm+fLHvGI=
+github.com/gin-gonic/gin v1.10.0 h1:nTuyha1TYqgedzytsKYqna+DfLos46nTv2ygFy86HFU=
+github.com/gin-gonic/gin v1.10.0/go.mod h1:4PMNQiOhvDRa013RKVbsiNwoyezlm2rm0uX/T7kzp5Y=
+github.com/go-playground/assert/v2 v2.2.0 h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=
+github.com/go-playground/assert/v2 v2.2.0/go.mod h1:VDjEfimB/XKnb+ZQfWdccd7VUvScMdVu0Titje2rxJ4=
+github.com/go-playground/locales v0.14.1 h1:EWaQ/wswjilfKLTECiXz7Rh+3BjFhfDFKv/oXslEjJA=
+github.com/go-playground/locales v0.14.1/go.mod h1:hxrqLVvrK65+Rwrd5Fc6F2O76J/NuW9t0sjnWqG1slY=
+github.com/go-playground/universal-translator v0.18.1 h1:Bcnm0ZwsGyWbCzImXv+pAJnYK9S473LQFuzCbDbfSFY=
+github.com/go-playground/universal-translator v0.18.1/go.mod h1:xekY+UJKNuX9WP91TpwSH2VMlDf28Uj24BCp08ZFTUY=
+github.com/go-playground/validator/v10 v10.20.0 h1:K9ISHbSaI0lyB2eWMPJo+kOS/FBExVwjEviJTixqxL8=
+github.com/go-playground/validator/v10 v10.20.0/go.mod h1:dbuPbCMFw/DrkbEynArYaCwl3amGuJotoKCe95atGMM=
+github.com/goccy/go-json v0.10.2 h1:CrxCmQqYDkv1z7lO7Wbh2HN93uovUHgrECaO5ZrCXAU=
+github.com/goccy/go-json v0.10.2/go.mod h1:6MelG93GURQebXPDq3khkgXZkazVtN9CRI+MGFi0w8I=
 github.com/google/go-cmp v0.5.9 h1:O2Tfq5qg4qc4AmwVlvv0oLiVAGB7enBSJ2x2DqQFi38=
 github.com/google/go-cmp v0.5.9/go.mod h1:17dUlkBOakJ0+DkrSSNjCkIjxS6bF9zb3elmeNGIjoY=
+github.com/google/gofuzz v1.0.0/go.mod h1:dBl0BpW6vV/+mYPU4Po3pmUjxk6FQPldtuIdl/M65Eg=
+github.com/google/uuid v1.5.0 h1:1p67kYwdtXjb0gL0BPiP1Av9wiZPo5A8z2cWkTZ+eyU=
+github.com/google/uuid v1.5.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
+github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
 github.com/hashicorp/hcl v1.0.0 h1:0Anlzjpi4vEasTeNFn2mLJgTSwt0+6sfsiTG8qcWGx4=
 github.com/hashicorp/hcl v1.0.0/go.mod h1:E5yfLk+7swimpb2L/Alb/PJmXilQ/rhwaUYs4T20WEQ=
+github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
+github.com/json-iterator/go v1.1.12/go.mod h1:e30LSqwooZae/UwlEbR2852Gd8hjQvJoHmT4TnhNGBo=
+github.com/juju/gnuflag v0.0.0-20171113085948-2ce1bb71843d/go.mod h1:2PavIy+JPciBPrBUjwbNvtwB6RQlve+hkpll6QSNmOE=
+github.com/klauspost/cpuid/v2 v2.0.9/go.mod h1:FInQzS24/EEf25PyTYn52gqo7WaD8xa0213Md/qVLRg=
+github.com/klauspost/cpuid/v2 v2.2.7 h1:ZWSB3igEs+d0qvnxR/ZBzXVmxkgt8DdzP6m9pfuVLDM=
+github.com/klauspost/cpuid/v2 v2.2.7/go.mod h1:Lcz8mBdAVJIBVzewtcLocK12l3Y+JytZYpaMropDUws=
+github.com/knz/go-libedit v1.10.1/go.mod h1:MZTVkCWyz0oBc7JOWP3wNAzd002ZbM/5hgShxwh4x8M=
 github.com/kr/pretty v0.3.1 h1:flRD4NNwYAUpkphVc1HcthR4KEIFJ65n8Mw5qdRn3LE=
 github.com/kr/pretty v0.3.1/go.mod h1:hoEshYVHaxMs3cyo3Yncou5ZscifuDolrwPKZanG3xk=
 github.com/kr/text v0.2.0 h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=
 github.com/kr/text v0.2.0/go.mod h1:eLer722TekiGuMkidMxC/pM04lWEeraHUUmBw8l2grE=
+github.com/leodido/go-urn v1.4.0 h1:WT9HwE9SGECu3lg4d/dIA+jxlljEa1/ffXKmRjqdmIQ=
+github.com/leodido/go-urn v1.4.0/go.mod h1:bvxc+MVxLKB4z00jd1z+Dvzr47oO32F/QSNjSBOlFxI=
 github.com/magiconair/properties v1.8.7 h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=
 github.com/magiconair/properties v1.8.7/go.mod h1:Dhd985XPs7jluiymwWYZ0G4Z61jb3vdS329zhj2hYo0=
+github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
+github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
 github.com/mitchellh/mapstructure v1.5.0 h1:jeMsZIYE/09sWLaz43PL7Gy6RuMjD2eJVyuac5Z2hdY=
 github.com/mitchellh/mapstructure v1.5.0/go.mod h1:bFUtVrKA4DC2yAKiSyO/QUcy7e+RRV2QTWOzhPopBRo=
+github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
+github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=
+github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
+github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
+github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
+github.com/oapi-codegen/runtime v1.1.1 h1:EXLHh0DXIJnWhdRPN2w4MXAzFyE4CskzhNLUmtpMYro=
+github.com/oapi-codegen/runtime v1.1.1/go.mod h1:SK9X900oXmPWilYR5/WKPzt3Kqxn/uS/+lbpREv+eCg=
 github.com/pelletier/go-toml/v2 v2.2.2 h1:aYUidT7k73Pcl9nb2gScu7NSrKCSHIDE89b3+6Wq+LM=
 github.com/pelletier/go-toml/v2 v2.2.2/go.mod h1:1t835xjRzz80PqgE6HHgN2JOsmgYu/h4qDAS4n929Rs=
 github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
@@ -39,28 +89,46 @@ github.com/spf13/pflag v1.0.5 h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=
 github.com/spf13/pflag v1.0.5/go.mod h1:McXfInJRrz4CZXVZOBLb0bTZqETkiAhM9Iw0y3An2Bg=
 github.com/spf13/viper v1.19.0 h1:RWq5SEjt8o25SROyN3z2OrDB9l7RPd3lwTWU8EcEdcI=
 github.com/spf13/viper v1.19.0/go.mod h1:GQUN9bilAbhU/jgc1bKs99f/suXKeUMct8Adx5+Ntkg=
+github.com/spkg/bom v0.0.0-20160624110644-59b7046e48ad/go.mod h1:qLr4V1qq6nMqFKkMo8ZTx3f+BZEkzsRUY10Xsm2mwU0=
 github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
 github.com/stretchr/objx v0.5.0/go.mod h1:Yh+to48EsGEfYuaHDzXPcE3xhTkx73EhmCGUpEOglKo=
 github.com/stretchr/objx v0.5.2/go.mod h1:FRsXN1f5AsAjCGJKqEizvkpNtU+EGNCLh3NxZ/8L+MA=
 github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
+github.com/stretchr/testify v1.7.0/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
 github.com/stretchr/testify v1.7.1/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
 github.com/stretchr/testify v1.8.0/go.mod h1:yNjHg4UonilssWZ8iaSj1OCr/vHnekPRkoO+kdMU+MU=
+github.com/stretchr/testify v1.8.1/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o6fzry7u4=
 github.com/stretchr/testify v1.8.4/go.mod h1:sz/lmYIOXD/1dqDmKjjqLyZ2RngseejIcXlSw2iwfAo=
 github.com/stretchr/testify v1.9.0 h1:HtqpIVDClZ4nwg75+f6Lvsy/wHu+3BoSGCbBAcpTsTg=
 github.com/stretchr/testify v1.9.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8C91i36aY=
 github.com/subosito/gotenv v1.6.0 h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=
 github.com/subosito/gotenv v1.6.0/go.mod h1:Dk4QP5c2W3ibzajGcXpNraDfq2IrhjMIvMSWPKKo0FU=
+github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS4MhqMhdFk5YI=
+github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
+github.com/ugorji/go/codec v1.2.12 h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65EE=
+github.com/ugorji/go/codec v1.2.12/go.mod h1:UNopzCgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=
 go.uber.org/atomic v1.9.0 h1:ECmE8Bn/WFTYwEW/bpKD3M8VtR/zQVbavAoalC1PYyE=
 go.uber.org/atomic v1.9.0/go.mod h1:fEN4uk6kAWBTFdckzkM89CLk9XfWZrxpCo0nPH17wJc=
 go.uber.org/multierr v1.9.0 h1:7fIwc/ZtS0q++VgcfqFDxSBZVv/Xo49/SYnDFupUwlI=
 go.uber.org/multierr v1.9.0/go.mod h1:X2jQV1h+kxSjClGpnseKVIxpmcjrj7MNnI0bnlfKTVQ=
+golang.org/x/arch v0.0.0-20210923205945-b76863e36670/go.mod h1:5om86z9Hs0C8fWVUuoMHwpExlXzs5Tkyp9hOrfG7pp8=
+golang.org/x/arch v0.8.0 h1:3wRIsP3pM4yUptoR96otTUOXI367OS0+c9eeRi9doIc=
+golang.org/x/arch v0.8.0/go.mod h1:FEVrYAQjsQXMVJ1nsMoVVXPZg6p2JE2mx8psSWTDQys=
+golang.org/x/crypto v0.23.0 h1:dIJU/v2J8Mdglj/8rJ6UUOM3Zc9zLZxVZwwxMooUSAI=
+golang.org/x/crypto v0.23.0/go.mod h1:CKFgDieR+mRhux2Lsu27y0fO304Db0wZe70UKqHu0v8=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9 h1:GoHiUyI/Tp2nVkLI2mCxVkOjsbSXD66ic0XW0js0R9g=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9/go.mod h1:S2oDrQGGwySpoQPVqRShND87VCbxmc6bL1Yd2oYrm6k=
-golang.org/x/sys v0.18.0 h1:DBdB3niSjOA/O0blCZBqDefyWNYveAYMNF1Wum0DYQ4=
-golang.org/x/sys v0.18.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
-golang.org/x/text v0.14.0 h1:ScX5w1eTa3QqT8oi6+ziP7dTV1S2+ALU0bI+0zXKWiQ=
-golang.org/x/text v0.14.0/go.mod h1:18ZOQIKpY8NJVqYksKHtTdi31H5itFRjB5/qKTNYzSU=
+golang.org/x/net v0.25.0 h1:d/OCCoBEUq33pjydKrGQhw7IlUPI2Oylr+8qLx49kac=
+golang.org/x/net v0.25.0/go.mod h1:JkAGAh7GEvH74S6FOH42FLoXpXbE/aqXSrIQjXgsiwM=
+golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.20.0 h1:Od9JTbYCk261bKm4M/mw7AklTlFYIa0bIp9BgSm1S8Y=
+golang.org/x/sys v0.20.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
+golang.org/x/text v0.15.0 h1:h1V/4gjBv8v9cjcR6+AR5+/cIYK5N/WAgiv4xlsEtAk=
+golang.org/x/text v0.15.0/go.mod h1:18ZOQIKpY8NJVqYksKHtTdi31H5itFRjB5/qKTNYzSU=
+google.golang.org/protobuf v1.34.1 h1:9ddQBjfCyZPOHPUiPxpYESBLc+T8P3E+Vo4IbKZgFWg=
+google.golang.org/protobuf v1.34.1/go.mod h1:c6P6GXX6sHbq/GpV6MGZEdwhWPcYBgnhAHhKbcUYpos=
 gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 h1:YR8cESwS4TdDjEe65xsg0ogRM/Nc3DYOhEAlW+xobZo=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
@@ -69,3 +137,5 @@ gopkg.in/ini.v1 v1.67.0/go.mod h1:pNLf8WUiyNEtQjuu5G5vTm06TEv9tsIgeAvK8hOrP4k=
 gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
 gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
 gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
+nullprogram.com/x/optparse v1.0.0/go.mod h1:KdyPE+Igbe0jQUrVfMqDMeJQIJZEuyV7pjYmp6pbG50=
+rsc.io/pdf v0.1.1/go.mod h1:n8OzWcQ6Sp37PL01nO98y4iUCRdTGarVfzxY20ICaU4=
diff --git a/internal/order/ports/openapi_api.gen.go b/internal/order/ports/openapi_api.gen.go
new file mode 100644
index 0000000..9cbaed0
--- /dev/null
+++ b/internal/order/ports/openapi_api.gen.go
@@ -0,0 +1,119 @@
+// Package ports provides primitives to interact with the openapi HTTP API.
+//
+// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
+package ports
+
+import (
+	"fmt"
+	"net/http"
+
+	"github.com/gin-gonic/gin"
+	"github.com/oapi-codegen/runtime"
+)
+
+// ServerInterface represents all server handlers.
+type ServerInterface interface {
+
+	// (POST /customer/{customerID}/orders)
+	PostCustomerCustomerIDOrders(c *gin.Context, customerID string)
+
+	// (GET /customer/{customerID}/orders/{orderID})
+	GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string)
+}
+
+// ServerInterfaceWrapper converts contexts to parameters.
+type ServerInterfaceWrapper struct {
+	Handler            ServerInterface
+	HandlerMiddlewares []MiddlewareFunc
+	ErrorHandler       func(*gin.Context, error, int)
+}
+
+type MiddlewareFunc func(c *gin.Context)
+
+// PostCustomerCustomerIDOrders operation middleware
+func (siw *ServerInterfaceWrapper) PostCustomerCustomerIDOrders(c *gin.Context) {
+
+	var err error
+
+	// ------------- Path parameter "customerID" -------------
+	var customerID string
+
+	err = runtime.BindStyledParameterWithOptions("simple", "customerID", c.Param("customerID"), &customerID, runtime.BindStyledParameterOptions{Explode: false, Required: true})
+	if err != nil {
+		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter customerID: %w", err), http.StatusBadRequest)
+		return
+	}
+
+	for _, middleware := range siw.HandlerMiddlewares {
+		middleware(c)
+		if c.IsAborted() {
+			return
+		}
+	}
+
+	siw.Handler.PostCustomerCustomerIDOrders(c, customerID)
+}
+
+// GetCustomerCustomerIDOrdersOrderID operation middleware
+func (siw *ServerInterfaceWrapper) GetCustomerCustomerIDOrdersOrderID(c *gin.Context) {
+
+	var err error
+
+	// ------------- Path parameter "customerID" -------------
+	var customerID string
+
+	err = runtime.BindStyledParameterWithOptions("simple", "customerID", c.Param("customerID"), &customerID, runtime.BindStyledParameterOptions{Explode: false, Required: true})
+	if err != nil {
+		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter customerID: %w", err), http.StatusBadRequest)
+		return
+	}
+
+	// ------------- Path parameter "orderID" -------------
+	var orderID string
+
+	err = runtime.BindStyledParameterWithOptions("simple", "orderID", c.Param("orderID"), &orderID, runtime.BindStyledParameterOptions{Explode: false, Required: true})
+	if err != nil {
+		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter orderID: %w", err), http.StatusBadRequest)
+		return
+	}
+
+	for _, middleware := range siw.HandlerMiddlewares {
+		middleware(c)
+		if c.IsAborted() {
+			return
+		}
+	}
+
+	siw.Handler.GetCustomerCustomerIDOrdersOrderID(c, customerID, orderID)
+}
+
+// GinServerOptions provides options for the Gin server.
+type GinServerOptions struct {
+	BaseURL      string
+	Middlewares  []MiddlewareFunc
+	ErrorHandler func(*gin.Context, error, int)
+}
+
+// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
+func RegisterHandlers(router gin.IRouter, si ServerInterface) {
+	RegisterHandlersWithOptions(router, si, GinServerOptions{})
+}
+
+// RegisterHandlersWithOptions creates http.Handler with additional options
+func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
+	errorHandler := options.ErrorHandler
+	if errorHandler == nil {
+		errorHandler = func(c *gin.Context, err error, statusCode int) {
+			c.JSON(statusCode, gin.H{"msg": err.Error()})
+		}
+	}
+
+	wrapper := ServerInterfaceWrapper{
+		Handler:            si,
+		HandlerMiddlewares: options.Middlewares,
+		ErrorHandler:       errorHandler,
+	}
+
+	router.POST(options.BaseURL+"/customer/:customerID/orders", wrapper.PostCustomerCustomerIDOrders)
+	router.GET(options.BaseURL+"/customer/:customerID/orders/:orderID", wrapper.GetCustomerCustomerIDOrdersOrderID)
+}
diff --git a/internal/order/ports/openapi_types.gen.go b/internal/order/ports/openapi_types.gen.go
new file mode 100644
index 0000000..b6b8b9c
--- /dev/null
+++ b/internal/order/ports/openapi_types.gen.go
@@ -0,0 +1,41 @@
+// Package ports provides primitives to interact with the openapi HTTP API.
+//
+// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
+package ports
+
+// CreateOrderRequest defines model for CreateOrderRequest.
+type CreateOrderRequest struct {
+	CustomerID string             `json:"customerID"`
+	Items      []ItemWithQuantity `json:"items"`
+}
+
+// Error defines model for Error.
+type Error struct {
+	Message *string `json:"message,omitempty"`
+}
+
+// Item defines model for Item.
+type Item struct {
+	Id       *string `json:"id,omitempty"`
+	Name     *string `json:"name,omitempty"`
+	PriceID  *string `json:"priceID,omitempty"`
+	Quantity *int32  `json:"quantity,omitempty"`
+}
+
+// ItemWithQuantity defines model for ItemWithQuantity.
+type ItemWithQuantity struct {
+	Id       *string `json:"id,omitempty"`
+	Quantity *int32  `json:"quantity,omitempty"`
+}
+
+// Order defines model for Order.
+type Order struct {
+	CustomerID  *string `json:"customerID,omitempty"`
+	Id          *string `json:"id,omitempty"`
+	Items       *[]Item `json:"items,omitempty"`
+	PaymentLink *string `json:"paymentLink,omitempty"`
+	Status      *string `json:"status,omitempty"`
+}
+
+// PostCustomerCustomerIDOrdersJSONRequestBody defines body for PostCustomerCustomerIDOrders for application/json ContentType.
+type PostCustomerCustomerIDOrdersJSONRequestBody = CreateOrderRequest
diff --git a/internal/payment/go.mod b/internal/payment/go.mod
new file mode 100644
index 0000000..bb571c2
--- /dev/null
+++ b/internal/payment/go.mod
@@ -0,0 +1,7 @@
+module github.com/ghost-yu/go_shop_second/payment
+
+go 1.22.8
+
+replace (
+	github.com/ghost-yu/go_shop_second/common => ../common
+)
diff --git a/internal/stock/go.mod b/internal/stock/go.mod
new file mode 100644
index 0000000..52221aa
--- /dev/null
+++ b/internal/stock/go.mod
@@ -0,0 +1,7 @@
+module github.com/ghost-yu/go_shop_second/stock
+
+go 1.22.8
+
+replace (
+	github.com/ghost-yu/go_shop_second/common => ../common
+)
\ No newline at end of file
diff --git a/scripts/genopenapi.sh b/scripts/genopenapi.sh
new file mode 100755
index 0000000..421e3f4
--- /dev/null
+++ b/scripts/genopenapi.sh
@@ -0,0 +1,54 @@
+#!/usr/bin/env bash
+
+set -euo pipefail
+
+shopt -s globstar
+
+if ! [[ "$0" =~ scripts/genopenapi.sh ]]; then
+  echo "must be run from repository root"
+  exit 255
+fi
+
+source ./scripts/lib.sh
+
+OPENAPI_ROOT="./api/openapi"
+
+GEN_SERVER=(
+#  "chi-server"
+#  "echo-server"
+  "gin-server"
+)
+
+if [ "${#GEN_SERVER[@]}" -ne 1 ]; then
+  log_error "GEN_SERVER enables more than 1 server, please check."
+  exit 255
+fi
+
+log_callout "Using ${GEN_SERVER[0]}"
+
+function openapi_files {
+  openapi_files=$(ls ${OPENAPI_ROOT})
+  echo "${openapi_files[@]}"
+}
+
+# output_dir, package_name, service_name
+function gen() {
+  local output_dir=$1
+  local package=$2
+  local service=$3
+
+  run mkdir -p "$output_dir"
+  run find "$output_dir" -type f -name "*.gen.go" -delete
+
+  prepare_dir "internal/common/client/$service"
+
+  run oapi-codegen -generate types -o "$output_dir/openapi_types.gen.go" -package "$package" "api/openapi/$service.yml"
+  run oapi-codegen -generate "$GEN_SERVER" -o "$output_dir/openapi_api.gen.go" -package "$package" "api/openapi/$service.yml"
+
+  run oapi-codegen -generate client -o "internal/common/client/$service/openapi_client.gen.go" -package "$service" "api/openapi/$service.yml"
+  run oapi-codegen -generate types -o "internal/common/client/$service/openapi_types.gen.go" -package "$service" "api/openapi/$service.yml"
+}
+
+gen internal/order/ports ports order
+
+log_success "openapi generate success!"
\ No newline at end of file
diff --git a/scripts/genproto.sh b/scripts/genproto.sh
new file mode 100755
index 0000000..64a8e03
--- /dev/null
+++ b/scripts/genproto.sh
@@ -0,0 +1,62 @@
+#!/usr/bin/env bash
+
+set -euo pipefail
+
+shopt -s globstar
+
+if ! [[ "$0" =~ scripts/genproto.sh ]]; then
+  echo "must be run from repository root"
+  exit 255
+fi
+
+source ./scripts/lib.sh
+
+API_ROOT="./api"
+
+function dirs {
+  dirs=()
+  while IFS= read -r dir; do
+    dirs+=("$dir")
+  done < <(find . -type f -name "*.proto" -exec dirname {} \; | xargs -n1 basename | sort -u)
+  echo "${dirs[@]}"
+}
+
+function pb_files {
+  pb_files=$(find . -type f -name '*.proto')
+  echo "${pb_files[@]}"
+}
+
+function gen_for_modules() {
+  local go_out="internal/common/genproto"
+  if [ -d "$go_out" ]; then
+    log_warning "found existing $go_out, cleaning all files under it"
+    run rm -rf $go_out
+  fi
+
+  for dir in $(dirs); do
+    echo "dir=$dir"
+    local service="${dir:0:${#dir}-2}"
+    local pb_file="${service}.proto"
+
+    if [ -d "$go_out/$dir" ]; then
+        log_warning "found existing $go_out/$dir, cleaning all files under it"
+        run rm -rf "$go_out"/$dir/*
+    else
+      run mkdir -p "$go_out/$dir"
+    fi
+    log_info "generating code for $service to $go_out/$dir"
+
+    run protoc \
+      -I="/usr/local/include/" \
+      -I="${API_ROOT}" \
+      "--go_out=${go_out}" --go_opt=paths=source_relative \
+      --go-grpc_opt=require_unimplemented_servers=false \
+      "--go-grpc_out=${go_out}" --go-grpc_opt=paths=source_relative \
+      "${API_ROOT}/${dir}/$pb_file"
+  done
+  log_success "protoc gen done!"
+}
+
+echo "directories containing protos to be built: $(dirs)"
+echo "found pb_files: $(pb_files)"
+gen_for_modules
\ No newline at end of file
diff --git a/scripts/lib.sh b/scripts/lib.sh
new file mode 100755
index 0000000..07321a1
--- /dev/null
+++ b/scripts/lib.sh
@@ -0,0 +1,220 @@
+#!/usr/bin/env bash
+
+set -euo pipefail
+
+ROOT_MODULE="$PWD"
+
+function set_root_dir {
+  ROOT_DIR=$ROOT_MODULE
+#  echo set to $ROOT_DIR
+}
+set_root_dir
+
+COLOR_RED='\033[0;31m'
+COLOR_ORANGE='\033[0;33m'
+COLOR_GREEN='\033[0;32m'
+COLOR_LIGHTCYAN='\033[0;36m'
+COLOR_BLUE='\033[0;94m'
+COLOR_MAGENTA='\033[95m'
+COLOR_BOLD='\033[1m'
+COLOR_NONE='\033[0m' # No Color
+
+
+function log_error {
+  >&2 echo -n -e "${COLOR_BOLD}${COLOR_RED}"
+  >&2 echo "$@"
+  >&2 echo -n -e "${COLOR_NONE}"
+}
+
+function log_warning {
+  >&2 echo -n -e "${COLOR_ORANGE}"
+  >&2 echo "$@"
+  >&2 echo -n -e "${COLOR_NONE}"
+}
+
+function log_callout {
+  >&2 echo -n -e "${COLOR_LIGHTCYAN}"
+  >&2 echo "$@"
+  >&2 echo -n -e "${COLOR_NONE}"
+}
+
+function log_cmd {
+  >&2 echo -n -e "${COLOR_BLUE}"
+  >&2 echo "$@"
+  >&2 echo -n -e "${COLOR_NONE}"
+}
+
+function log_success {
+  >&2 echo -n -e "${COLOR_GREEN}"
+  >&2 echo "$@"
+  >&2 echo -n -e "${COLOR_NONE}"
+}
+
+function log_info {
+  >&2 echo -n -e "${COLOR_NONE}"
+  >&2 echo "$@"
+  >&2 echo -n -e "${COLOR_NONE}"
+}
+
+function modules() {
+  modules=$(ls internal)
+  echo "${modules[@]}"
+}
+
+function modules_exp() {
+  for m in $(modules); do
+    echo -n "${m}/... "
+  done
+}
+
+# From http://stackoverflow.com/a/12498485
+function relativePath {
+  # both $1 and $2 are absolute paths beginning with /
+  # returns relative path to $2 from $1
+  local source=$1
+  local target=$2
+
+  local commonPart=$source
+  local result=""
+
+  while [[ "${target#"$commonPart"}" == "${target}" ]]; do
+    # no match, means that candidate common part is not correct
+    # go up one level (reduce common part)
+    commonPart="$(dirname "$commonPart")"
+    # and record that we went back, with correct / handling
+    if [[ -z $result ]]; then
+      result=".."
+    else
+      result="../$result"
+    fi
+  done
+
+  if [[ $commonPart == "/" ]]; then
+    # special case for root (no common path)
+    result="$result/"
+  fi
+
+  # since we now have identified the common part,
+  # compute the non-common part
+  local forwardPart="${target#"$commonPart"}"
+
+  # and now stick all parts together
+  if [[ -n $result ]] && [[ -n $forwardPart ]]; then
+    result="$result$forwardPart"
+  elif [[ -n $forwardPart ]]; then
+    # extra slash removal
+    result="${forwardPart:1}"
+  fi
+
+  echo "$result"
+}
+
+function module_dirs() {
+  echo "internal/common internal/order internal/stock internal/payment"
+}
+
+function module_subdir {
+  relativePath "${ROOT_DIR}" "${PWD}"
+}
+
+####    Running actions against multiple modules ####
+
+# run [command...] - runs given command, printing it first and
+# again if it failed (in RED). Use to wrap important test commands
+# that user might want to re-execute to shorten the feedback loop when fixing
+# the test.
+function run {
+  local rpath
+  local command
+  rpath=$(module_subdir)
+  # Quoting all components as the commands are fully copy-parsable:
+  command=("${@}")
+  command=("${command[@]@Q}")
+  if [[ "${rpath}" != "." && "${rpath}" != "" ]]; then
+    repro="(cd ${rpath} && ${command[*]})"
+  else
+    repro="${command[*]}"
+  fi
+
+  log_cmd "% ${repro}"
+  "${@}" 2> >(while read -r line; do echo -e "${COLOR_NONE}stderr: ${COLOR_MAGENTA}${line}${COLOR_NONE}">&2; done)
+  local error_code=$?
+  if [ ${error_code} -ne 0 ]; then
+    log_error -e "FAIL: (code:${error_code}):\\n  % ${repro}"
+    return ${error_code}
+  fi
+}
+
+# receives a directory, relative to ROOT_DIR. If it exists, delete all files under it,
+# otherwise create that directory
+function prepare_dir() {
+  local dir="$1"
+  if [ -d "$dir" ]; then
+    log_warning "Directory $dir exists. Deleting all files under $dir."
+    run find "$dir" -mindepth 1 -delete
+  else
+    log_callout "Directory $dir does not exist. Creating directory $dir."
+    # Create the directory
+    run mkdir -p "$dir"
+  fi
+}
+
+# run_for_module [module] [cmd]
+# executes given command in the given module for given pkgs.
+#   module_name - "." (in future: tests, client, server)
+#   cmd         - cmd to be executed - that takes package as last argument
+function run_for_module {
+  local module=${1:-"."}
+  shift 1
+  (
+    cd "${ROOT_DIR}/${module}" && "$@"
+  )
+}
+
+#  run_for_modules [cmd]
+#  run given command across all modules and packages
+#  (unless the set is limited using ${PKG} or / ${USERMOD})
+function run_for_modules {
+  KEEP_GOING_MODULE=${KEEP_GOING_MODULE:-false}
+  local pkg="${PKG:-./...}"
+  log_info "pkg = $pkg"
+  local fail_mod=false
+  if [ -z "${USERMOD:-}" ]; then
+    for m in $(module_dirs); do
+      if run_for_module "${m}" "$@" "${pkg}"; then
+        continue
+      else
+        if [ "$KEEP_GOING_MODULE" = false ]; then
+          log_error "There was a Failure in module ${m}, aborting..."
+          return 1
+        fi
+        log_error "There was a Failure in module ${m}, keep going..."
+        fail_mod=true
+      fi
+    done
+    if [ "$fail_mod" = true ]; then
+      return 1
+    fi
+  else
+    run_for_module "${USERMOD}" "$@" "${pkg}" || return "$?"
+  fi
+}
+
+# generic_checker [cmd...]
+# executes given command in the current module, and clearly fails if it
+# failed or returned output.
+function generic_checker {
+  local cmd=("$@")
+  if ! output=$("${cmd[@]}"); then
+    echo "${output}"
+    log_error -e "FAIL: '${cmd[*]}' checking failed (!=0 return code)"
+    return 255
+  fi
+  if [ -n "${output}" ]; then
+    echo "${output}"
+    log_error -e "FAIL: '${cmd[*]}' checking failed (printed output)"
+    return 255
+  fi
+}
+
+echo "go modules: $(modules)"
\ No newline at end of file
~~~

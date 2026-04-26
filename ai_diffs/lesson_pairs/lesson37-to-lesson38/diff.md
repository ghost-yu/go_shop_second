# Lesson Pair Diff Report

- FromBranch: lesson37
- ToBranch: lesson38

## Short Summary

~~~text
 12 files changed, 128 insertions(+), 23 deletions(-)
~~~

## File Stats

~~~text
 api/openapi/cfg.yaml                               |  2 +
 api/openapi/order.yml                              | 23 +++++++-
 internal/common/client/order/openapi_client.gen.go |  8 +--
 internal/common/client/order/openapi_types.gen.go  |  8 +++
 internal/order/adapters/order_mongo_repository.go  |  2 +-
 internal/order/go.mod                              |  3 +
 internal/order/go.sum                              |  3 +-
 internal/order/ports/openapi_types.gen.go          |  8 +++
 internal/order/tests/create_order_test.go          | 65 ++++++++++++++++++++++
 internal/payment/http.go                           |  1 +
 public/success.html                                | 20 +++----
 scripts/genopenapi.sh                              |  8 +--
 12 files changed, 128 insertions(+), 23 deletions(-)
~~~

## Commit Comparison

~~~text
> 0ee059f unit test
> e7275b6 test all
~~~

## Changed Files

~~~text
api/openapi/cfg.yaml
api/openapi/order.yml
internal/common/client/order/openapi_client.gen.go
internal/common/client/order/openapi_types.gen.go
internal/order/adapters/order_mongo_repository.go
internal/order/go.mod
internal/order/go.sum
internal/order/ports/openapi_types.gen.go
internal/order/tests/create_order_test.go
internal/payment/http.go
public/success.html
scripts/genopenapi.sh
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
api/openapi/cfg.yaml
api/openapi/order.yml
internal/order/adapters/order_mongo_repository.go
internal/order/tests/create_order_test.go
internal/payment/http.go
public/success.html
scripts/genopenapi.sh
~~~

## Full Diff

~~~diff
diff --git a/api/openapi/cfg.yaml b/api/openapi/cfg.yaml
new file mode 100644
index 0000000..3ac1679
--- /dev/null
+++ b/api/openapi/cfg.yaml
@@ -0,0 +1,2 @@
+output-options:
+   skip-prune: true
\ No newline at end of file
diff --git a/api/openapi/order.yml b/api/openapi/order.yml
index da9fbf8..6155cc0 100644
--- a/api/openapi/order.yml
+++ b/api/openapi/order.yml
@@ -32,7 +32,7 @@ paths:
           content:
             application/json:
               schema:
-                $ref: '#/components/schemas/Order'
+                $ref: '#/components/schemas/Response'
 
         default:
           description: todo
@@ -64,7 +64,7 @@ paths:
           content:
             application/json:
               schema:
-                $ref: '#/components/schemas/Order'
+                $ref: '#/components/schemas/Response'
 
         default:
           description: todo
@@ -144,4 +144,21 @@ components:
           type: string
         quantity:
           type: integer
-          format: int32
\ No newline at end of file
+          format: int32
+
+    Response:
+      type: object
+      properties:
+        errno:
+          type: integer
+        message:
+          type: string
+        data:
+          type: object
+        trace_id:
+          type: string
+      required:
+        - errno
+        - message
+        - data
+        - trace_id
\ No newline at end of file
diff --git a/internal/common/client/order/openapi_client.gen.go b/internal/common/client/order/openapi_client.gen.go
index 77f014d..fb45027 100644
--- a/internal/common/client/order/openapi_client.gen.go
+++ b/internal/common/client/order/openapi_client.gen.go
@@ -277,7 +277,7 @@ type ClientWithResponsesInterface interface {
 type PostCustomerCustomerIdOrdersResponse struct {
 	Body         []byte
 	HTTPResponse *http.Response
-	JSON200      *Order
+	JSON200      *Response
 	JSONDefault  *Error
 }
 
@@ -300,7 +300,7 @@ func (r PostCustomerCustomerIdOrdersResponse) StatusCode() int {
 type GetCustomerCustomerIdOrdersOrderIdResponse struct {
 	Body         []byte
 	HTTPResponse *http.Response
-	JSON200      *Order
+	JSON200      *Response
 	JSONDefault  *Error
 }
 
@@ -361,7 +361,7 @@ func ParsePostCustomerCustomerIdOrdersResponse(rsp *http.Response) (*PostCustome
 
 	switch {
 	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
-		var dest Order
+		var dest Response
 		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
 			return nil, err
 		}
@@ -394,7 +394,7 @@ func ParseGetCustomerCustomerIdOrdersOrderIdResponse(rsp *http.Response) (*GetCu
 
 	switch {
 	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
-		var dest Order
+		var dest Response
 		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
 			return nil, err
 		}
diff --git a/internal/common/client/order/openapi_types.gen.go b/internal/common/client/order/openapi_types.gen.go
index 3b659fb..6aece2b 100644
--- a/internal/common/client/order/openapi_types.gen.go
+++ b/internal/common/client/order/openapi_types.gen.go
@@ -37,5 +37,13 @@ type Order struct {
 	Status      string `json:"status"`
 }
 
+// Response defines model for Response.
+type Response struct {
+	Data    map[string]interface{} `json:"data"`
+	Errno   int                    `json:"errno"`
+	Message string                 `json:"message"`
+	TraceId string                 `json:"trace_id"`
+}
+
 // PostCustomerCustomerIdOrdersJSONRequestBody defines body for PostCustomerCustomerIdOrders for application/json ContentType.
 type PostCustomerCustomerIdOrdersJSONRequestBody = CreateOrderRequest
diff --git a/internal/order/adapters/order_mongo_repository.go b/internal/order/adapters/order_mongo_repository.go
index fc4527d..ba202bd 100644
--- a/internal/order/adapters/order_mongo_repository.go
+++ b/internal/order/adapters/order_mongo_repository.go
@@ -115,7 +115,7 @@ func (r *OrderRepositoryMongo) Update(
 	if err != nil {
 		return
 	}
-	updated, err := updateFn(ctx, oldOrder)
+	updated, err := updateFn(ctx, order)
 	if err != nil {
 		return
 	}
diff --git a/internal/order/go.mod b/internal/order/go.mod
index 9cd0692..d780635 100644
--- a/internal/order/go.mod
+++ b/internal/order/go.mod
@@ -13,6 +13,7 @@ require (
 	github.com/rabbitmq/amqp091-go v1.10.0
 	github.com/sirupsen/logrus v1.9.3
 	github.com/spf13/viper v1.19.0
+	github.com/stretchr/testify v1.9.0
 	github.com/stripe/stripe-go/v80 v80.2.0
 	go.mongodb.org/mongo-driver v1.17.1
 	go.opentelemetry.io/otel v1.31.0
@@ -27,6 +28,7 @@ require (
 	github.com/bytedance/sonic/loader v0.2.0 // indirect
 	github.com/cloudwego/base64x v0.1.4 // indirect
 	github.com/cloudwego/iasm v0.2.0 // indirect
+	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
 	github.com/fatih/color v1.14.1 // indirect
 	github.com/fsnotify/fsnotify v1.7.0 // indirect
 	github.com/gabriel-vasile/mimetype v1.4.6 // indirect
@@ -63,6 +65,7 @@ require (
 	github.com/modern-go/reflect2 v1.0.2 // indirect
 	github.com/montanaflynn/stats v0.7.1 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
+	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
 	github.com/sagikazarmark/locafero v0.6.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
 	github.com/sourcegraph/conc v0.3.0 // indirect
diff --git a/internal/order/go.sum b/internal/order/go.sum
index aacb2f0..89f4896 100644
--- a/internal/order/go.sum
+++ b/internal/order/go.sum
@@ -272,8 +272,9 @@ github.com/spkg/bom v0.0.0-20160624110644-59b7046e48ad/go.mod h1:qLr4V1qq6nMqFKk
 github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.1.1/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
-github.com/stretchr/objx v0.5.0 h1:1zr/of2m5FGMsad5YfcqgdqdWrIhu+EBEJRhR1U7z/c=
 github.com/stretchr/objx v0.5.0/go.mod h1:Yh+to48EsGEfYuaHDzXPcE3xhTkx73EhmCGUpEOglKo=
+github.com/stretchr/objx v0.5.2 h1:xuMeJ0Sdp5ZMRXx/aWO6RZxdr3beISkG5/G/aIRr3pY=
+github.com/stretchr/objx v0.5.2/go.mod h1:FRsXN1f5AsAjCGJKqEizvkpNtU+EGNCLh3NxZ/8L+MA=
 github.com/stretchr/testify v1.2.2/go.mod h1:a8OnRcib4nhh0OaRAV+Yts87kKdq0PP7pXfy6kDkUVs=
 github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
 github.com/stretchr/testify v1.4.0/go.mod h1:j7eGeouHqKxXV5pUuKE4zz7dFj8WfuZ+81PSLYec5m4=
diff --git a/internal/order/ports/openapi_types.gen.go b/internal/order/ports/openapi_types.gen.go
index 97fe85c..351e855 100644
--- a/internal/order/ports/openapi_types.gen.go
+++ b/internal/order/ports/openapi_types.gen.go
@@ -37,5 +37,13 @@ type Order struct {
 	Status      string `json:"status"`
 }
 
+// Response defines model for Response.
+type Response struct {
+	Data    map[string]interface{} `json:"data"`
+	Errno   int                    `json:"errno"`
+	Message string                 `json:"message"`
+	TraceId string                 `json:"trace_id"`
+}
+
 // PostCustomerCustomerIdOrdersJSONRequestBody defines body for PostCustomerCustomerIdOrders for application/json ContentType.
 type PostCustomerCustomerIdOrdersJSONRequestBody = CreateOrderRequest
diff --git a/internal/order/tests/create_order_test.go b/internal/order/tests/create_order_test.go
new file mode 100644
index 0000000..2af7d2b
--- /dev/null
+++ b/internal/order/tests/create_order_test.go
@@ -0,0 +1,65 @@
+package tests
+
+import (
+	"context"
+	"fmt"
+	"log"
+	"testing"
+
+	sw "github.com/ghost-yu/go_shop_second/common/client/order"
+	_ "github.com/ghost-yu/go_shop_second/common/config"
+	"github.com/spf13/viper"
+	"github.com/stretchr/testify/assert"
+)
+
+var (
+	ctx    = context.Background()
+	server = fmt.Sprintf("http://%s/api", viper.GetString("order.http-addr"))
+)
+
+func TestMain(m *testing.M) {
+	before()
+	m.Run()
+}
+
+func before() {
+	log.Printf("server=%s", server)
+}
+
+func TestCreateOrder_success(t *testing.T) {
+	response := getResponse(t, "123", sw.PostCustomerCustomerIdOrdersJSONRequestBody{
+		CustomerId: "123",
+		Items: []sw.ItemWithQuantity{
+			{
+				Id:       "test-item-1",
+				Quantity: 1,
+			},
+		},
+	})
+	t.Logf("body=%s", string(response.Body))
+	assert.Equal(t, 200, response.StatusCode())
+
+	assert.Equal(t, 0, response.JSON200.Errno)
+}
+
+func TestCreateOrder_invalidParams(t *testing.T) {
+	response := getResponse(t, "123", sw.PostCustomerCustomerIdOrdersJSONRequestBody{
+		CustomerId: "123",
+		Items:      nil,
+	})
+	assert.Equal(t, 200, response.StatusCode())
+	assert.Equal(t, 2, response.JSON200.Errno)
+}
+
+func getResponse(t *testing.T, customerID string, body sw.PostCustomerCustomerIdOrdersJSONRequestBody) *sw.PostCustomerCustomerIdOrdersResponse {
+	t.Helper()
+	client, err := sw.NewClientWithResponses(server)
+	if err != nil {
+		t.Fatal(err)
+	}
+	response, err := client.PostCustomerCustomerIdOrdersWithResponse(ctx, customerID, body)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return response
+}
diff --git a/internal/payment/http.go b/internal/payment/http.go
index 884e2c9..9671c55 100644
--- a/internal/payment/http.go
+++ b/internal/payment/http.go
@@ -27,6 +27,7 @@ func NewPaymentHandler(ch *amqp.Channel) *PaymentHandler {
 	return &PaymentHandler{channel: ch}
 }
 
+// stripe listen --forward-to localhost:8284/api/webhook
 func (h *PaymentHandler) RegisterRoutes(c *gin.Engine) {
 	c.POST("/api/webhook", h.handleWebhook)
 }
diff --git a/public/success.html b/public/success.html
index 08abb53..1dd4267 100644
--- a/public/success.html
+++ b/public/success.html
@@ -128,22 +128,22 @@
         }
       }
        */
-      if (data.data.Order.Status === 'waiting_for_payment') {
-          order.Status = '等待支付...';
-          document.getElementById('orderStatus').innerText = order.Status;
+      if (data.data.status === 'waiting_for_payment') {
+          order.status = '等待支付...';
+          document.getElementById('orderStatus').innerText = order.status;
           document.querySelector('.after-payment-popup').style.display = 'block';
-          document.getElementById('payment-link').href = data.data.Order.PaymentLink;
+          document.getElementById('payment-link').href = data.data.payment_link;
       }
-      if (data.data.Order.Status === 'paid') {
-          order.Status = '已支付成功，请等待...';
-          document.getElementById('orderStatus').innerText = order.Status;
+      if (data.data.status === 'paid') {
+          order.status = '已支付成功，请等待...';
+          document.getElementById('orderStatus').innerText = order.status;
           setTimeout(getOrder, 5000);
-      } else if (data.data.Order.Status === 'ready') {
-          order.Status = '已完成...';
+      } else if (data.data.status === 'ready') {
+          order.status = '已完成...';
           document.querySelector('.after-payment-popup').style.display = 'none';
           document.querySelector('.ready-popup').style.display = 'block';
           document.getElementById('orderID').innerText = orderID;
-          document.getElementById('orderStatus').innerText = order.Status;
+          document.getElementById('orderStatus').innerText = order.status;
       } else {
           setTimeout(getOrder, 5000);
       }
diff --git a/scripts/genopenapi.sh b/scripts/genopenapi.sh
index 421e3f4..cc3566d 100755
--- a/scripts/genopenapi.sh
+++ b/scripts/genopenapi.sh
@@ -42,11 +42,11 @@ function gen() {
 
   prepare_dir "internal/common/client/$service"
 
-  run oapi-codegen -generate types -o "$output_dir/openapi_types.gen.go" -package "$package" "api/openapi/$service.yml"
-  run oapi-codegen -generate "$GEN_SERVER" -o "$output_dir/openapi_api.gen.go" -package "$package" "api/openapi/$service.yml"
+  run oapi-codegen -config api/openapi/cfg.yaml -generate types -o "$output_dir/openapi_types.gen.go" -package "$package" "api/openapi/$service.yml"
+  run oapi-codegen -config api/openapi/cfg.yaml -generate "$GEN_SERVER" -o "$output_dir/openapi_api.gen.go" -package "$package" "api/openapi/$service.yml"
 
-  run oapi-codegen -generate client -o "internal/common/client/$service/openapi_client.gen.go" -package "$service" "api/openapi/$service.yml"
-  run oapi-codegen -generate types -o "internal/common/client/$service/openapi_types.gen.go" -package "$service" "api/openapi/$service.yml"
+  run oapi-codegen -config api/openapi/cfg.yaml -generate client -o "internal/common/client/$service/openapi_client.gen.go" -package "$service" "api/openapi/$service.yml"
+  run oapi-codegen -config api/openapi/cfg.yaml -generate types -o "internal/common/client/$service/openapi_types.gen.go" -package "$service" "api/openapi/$service.yml"
 }
 
 gen internal/order/ports ports order
~~~

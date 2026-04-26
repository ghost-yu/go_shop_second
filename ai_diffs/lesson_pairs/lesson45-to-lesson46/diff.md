# Lesson Pair Diff Report

- FromBranch: lesson45
- ToBranch: lesson46

## Short Summary

~~~text
 14 files changed, 267 insertions(+), 89 deletions(-)
~~~

## File Stats

~~~text
 api/openapi/order.yml                              |   1 +
 internal/common/config/global.yaml                 |   7 ++
 internal/kitchen/prom.go                           | 133 +++++++++++----------
 internal/order/http.go                             |  13 ++
 internal/stock/adapters/stock_inmem_repository.go  |  10 +-
 internal/stock/adapters/stock_mysql_repository.go  |  36 ++++++
 .../stock/app/query/check_if_items_in_stock.go     |  41 +++++++
 internal/stock/app/query/get_items.go              |   8 +-
 internal/stock/domain/stock/repository.go          |  21 +++-
 internal/stock/go.mod                              |   7 +-
 internal/stock/go.sum                              |  19 +--
 internal/stock/infrastructure/persistent/mysql.go  |  53 ++++++++
 internal/stock/ports/grpc.go                       |   2 +-
 internal/stock/service/application.go              |   5 +-
 14 files changed, 267 insertions(+), 89 deletions(-)
~~~

## Commit Comparison

~~~text
> 98f7f85 stock validate
~~~

## Changed Files

~~~text
api/openapi/order.yml
internal/common/config/global.yaml
internal/kitchen/prom.go
internal/order/http.go
internal/stock/adapters/stock_inmem_repository.go
internal/stock/adapters/stock_mysql_repository.go
internal/stock/app/query/check_if_items_in_stock.go
internal/stock/app/query/get_items.go
internal/stock/domain/stock/repository.go
internal/stock/go.mod
internal/stock/go.sum
internal/stock/infrastructure/persistent/mysql.go
internal/stock/ports/grpc.go
internal/stock/service/application.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
api/openapi/order.yml
internal/common/config/global.yaml
internal/kitchen/prom.go
internal/order/http.go
internal/stock/adapters/stock_inmem_repository.go
internal/stock/adapters/stock_mysql_repository.go
internal/stock/app/query/check_if_items_in_stock.go
internal/stock/app/query/get_items.go
internal/stock/domain/stock/repository.go
internal/stock/infrastructure/persistent/mysql.go
internal/stock/ports/grpc.go
internal/stock/service/application.go
~~~

## Full Diff

~~~diff
diff --git a/api/openapi/order.yml b/api/openapi/order.yml
index 6155cc0..9b198f2 100644
--- a/api/openapi/order.yml
+++ b/api/openapi/order.yml
@@ -145,6 +145,7 @@ components:
         quantity:
           type: integer
           format: int32
+          minimum: 1
 
     Response:
       type: object
diff --git a/internal/common/config/global.yaml b/internal/common/config/global.yaml
index 099b0f5..ddcd94d 100644
--- a/internal/common/config/global.yaml
+++ b/internal/common/config/global.yaml
@@ -40,6 +40,13 @@ mongo:
   db-name: "order"
   coll-name: "order"
 
+mysql:
+  user: root
+  password: root
+  host: localhost
+  port: 3307
+  dbname: "gorder_v2"
+
 
 stripe-key: "${STRIPE_KEY}"
 endpoint-stripe-secret: "${ENDPOINT_STRIPE_SECRET}"
\ No newline at end of file
diff --git a/internal/kitchen/prom.go b/internal/kitchen/prom.go
index bf07668..0a39693 100644
--- a/internal/kitchen/prom.go
+++ b/internal/kitchen/prom.go
@@ -1,68 +1,69 @@
 package main
 
-import (
-	"bytes"
-	"encoding/json"
-	"log"
-	"math/rand"
-	"net/http"
-	"time"
-
-	"github.com/prometheus/client_golang/prometheus"
-	"github.com/prometheus/client_golang/prometheus/collectors"
-	"github.com/prometheus/client_golang/prometheus/promhttp"
-)
-
-const (
-	testAddr = "localhost:9123"
-)
-
-var httpStatusCodeCounter = prometheus.NewCounterVec(
-	prometheus.CounterOpts{
-		Name: "http_status_code_counter",
-		Help: "Count http status code",
-	},
-	[]string{"status_code"},
-)
-
-func main() {
-	go produceData()
-	reg := prometheus.NewRegistry()
-	prometheus.WrapRegistererWith(prometheus.Labels{"serviceName": "demo-service"}, reg).MustRegister(
-		collectors.NewGoCollector(),
-		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
-		httpStatusCodeCounter,
-	)
-	// localhost:9123/metrics
-	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
-	http.HandleFunc("/", sendMetricsHandler)
-	log.Fatal(http.ListenAndServe(testAddr, nil))
-}
-
-func sendMetricsHandler(w http.ResponseWriter, r *http.Request) {
-	var req request
-	defer func() {
-		httpStatusCodeCounter.WithLabelValues(req.StatusCode).Inc()
-		log.Printf("add 1 to %s", req.StatusCode)
-	}()
-	_ = json.NewDecoder(r.Body).Decode(&req)
-	log.Printf("receive req:%+v", req)
-	_, _ = w.Write([]byte(req.StatusCode))
-}
-
-type request struct {
-	StatusCode string
-}
-
-func produceData() {
-	codes := []string{"503", "404", "400", "200", "304", "500"}
-	for {
-		body, _ := json.Marshal(request{
-			StatusCode: codes[rand.Intn(len(codes))],
-		})
-		requestBody := bytes.NewBuffer(body)
-		http.Post("http://"+testAddr, "application/json", requestBody)
-		log.Printf("send request=%s to %s", requestBody.String(), testAddr)
-		time.Sleep(2 * time.Second)
-	}
-}
+//
+//import (
+//	"bytes"
+//	"encoding/json"
+//	"log"
+//	"math/rand"
+//	"net/http"
+//	"time"
+//
+//	"github.com/prometheus/client_golang/prometheus"
+//	"github.com/prometheus/client_golang/prometheus/collectors"
+//	"github.com/prometheus/client_golang/prometheus/promhttp"
+//)
+//
+//const (
+//	testAddr = "localhost:9123"
+//)
+//
+//var httpStatusCodeCounter = prometheus.NewCounterVec(
+//	prometheus.CounterOpts{
+//		Name: "http_status_code_counter",
+//		Help: "Count http status code",
+//	},
+//	[]string{"status_code"},
+//)
+//
+//func main() {
+//	go produceData()
+//	reg := prometheus.NewRegistry()
+//	prometheus.WrapRegistererWith(prometheus.Labels{"serviceName": "demo-service"}, reg).MustRegister(
+//		collectors.NewGoCollector(),
+//		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
+//		httpStatusCodeCounter,
+//	)
+//	// localhost:9123/metrics
+//	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
+//	http.HandleFunc("/", sendMetricsHandler)
+//	log.Fatal(http.ListenAndServe(testAddr, nil))
+//}
+//
+//func sendMetricsHandler(w http.ResponseWriter, r *http.Request) {
+//	var req request
+//	defer func() {
+//		httpStatusCodeCounter.WithLabelValues(req.StatusCode).Inc()
+//		log.Printf("add 1 to %s", req.StatusCode)
+//	}()
+//	_ = json.NewDecoder(r.Body).Decode(&req)
+//	log.Printf("receive req:%+v", req)
+//	_, _ = w.Write([]byte(req.StatusCode))
+//}
+//
+//type request struct {
+//	StatusCode string
+//}
+//
+//func produceData() {
+//	codes := []string{"503", "404", "400", "200", "304", "500"}
+//	for {
+//		body, _ := json.Marshal(request{
+//			StatusCode: codes[rand.Intn(len(codes))],
+//		})
+//		requestBody := bytes.NewBuffer(body)
+//		http.Post("http://"+testAddr, "application/json", requestBody)
+//		log.Printf("send request=%s to %s", requestBody.String(), testAddr)
+//		time.Sleep(2 * time.Second)
+//	}
+//}
diff --git a/internal/order/http.go b/internal/order/http.go
index a7cec4b..a8352ff 100644
--- a/internal/order/http.go
+++ b/internal/order/http.go
@@ -1,6 +1,7 @@
 package main
 
 import (
+	"errors"
 	"fmt"
 
 	"github.com/ghost-yu/go_shop_second/common"
@@ -31,6 +32,9 @@ func (H HTTPServer) PostCustomerCustomerIdOrders(c *gin.Context, customerID stri
 	if err = c.ShouldBindJSON(&req); err != nil {
 		return
 	}
+	if err = H.validate(req); err != nil {
+		return
+	}
 	r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
 		CustomerID: req.CustomerId,
 		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
@@ -64,3 +68,12 @@ func (H HTTPServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerI
 
 	resp = convertor.NewOrderConvertor().EntityToClient(o)
 }
+
+func (H HTTPServer) validate(req client.CreateOrderRequest) error {
+	for _, v := range req.Items {
+		if v.Quantity <= 0 {
+			return errors.New("quantity must be positive")
+		}
+	}
+	return nil
+}
diff --git a/internal/stock/adapters/stock_inmem_repository.go b/internal/stock/adapters/stock_inmem_repository.go
index f4aed23..2426756 100644
--- a/internal/stock/adapters/stock_inmem_repository.go
+++ b/internal/stock/adapters/stock_inmem_repository.go
@@ -4,16 +4,16 @@ import (
 	"context"
 	"sync"
 
-	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
+	"github.com/ghost-yu/go_shop_second/stock/entity"
 )
 
 type MemoryStockRepository struct {
 	lock  *sync.RWMutex
-	store map[string]*orderpb.Item
+	store map[string]*entity.Item
 }
 
-var stub = map[string]*orderpb.Item{
+var stub = map[string]*entity.Item{
 	"item_id": {
 		ID:       "foo_item",
 		Name:     "stub item",
@@ -47,11 +47,11 @@ func NewMemoryStockRepository() *MemoryStockRepository {
 	}
 }
 
-func (m MemoryStockRepository) GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error) {
+func (m MemoryStockRepository) GetItems(ctx context.Context, ids []string) ([]*entity.Item, error) {
 	m.lock.RLock()
 	defer m.lock.RUnlock()
 	var (
-		res     []*orderpb.Item
+		res     []*entity.Item
 		missing []string
 	)
 	for _, id := range ids {
diff --git a/internal/stock/adapters/stock_mysql_repository.go b/internal/stock/adapters/stock_mysql_repository.go
new file mode 100644
index 0000000..91d1e50
--- /dev/null
+++ b/internal/stock/adapters/stock_mysql_repository.go
@@ -0,0 +1,36 @@
+package adapters
+
+import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/stock/entity"
+	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent"
+)
+
+type MySQLStockRepository struct {
+	db *persistent.MySQL
+}
+
+func NewMySQLStockRepository(db *persistent.MySQL) *MySQLStockRepository {
+	return &MySQLStockRepository{db: db}
+}
+
+func (m MySQLStockRepository) GetItems(ctx context.Context, ids []string) ([]*entity.Item, error) {
+	//TODO implement me
+	panic("implement me")
+}
+
+func (m MySQLStockRepository) GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error) {
+	data, err := m.db.BatchGetStockByID(ctx, ids)
+	if err != nil {
+		return nil, err
+	}
+	var result []*entity.ItemWithQuantity
+	for _, d := range data {
+		result = append(result, &entity.ItemWithQuantity{
+			ID:       d.ProductID,
+			Quantity: d.Quantity,
+		})
+	}
+	return result, nil
+}
diff --git a/internal/stock/app/query/check_if_items_in_stock.go b/internal/stock/app/query/check_if_items_in_stock.go
index 00b5c4c..0e86e86 100644
--- a/internal/stock/app/query/check_if_items_in_stock.go
+++ b/internal/stock/app/query/check_if_items_in_stock.go
@@ -50,6 +50,9 @@ var stub = map[string]string{
 }
 
 func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*entity.Item, error) {
+	if err := h.checkStock(ctx, query.Items); err != nil {
+		return nil, err
+	}
 	var res []*entity.Item
 	for _, i := range query.Items {
 		priceID, err := h.stripeAPI.GetPriceByProductID(ctx, i.ID)
@@ -62,9 +65,47 @@ func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfIte
 			PriceID:  priceID,
 		})
 	}
+	// TODO: 扣库存
 	return res, nil
 }
 
+func (h checkIfItemsInStockHandler) checkStock(ctx context.Context, query []*entity.ItemWithQuantity) error {
+	var ids []string
+	for _, i := range query {
+		ids = append(ids, i.ID)
+	}
+	records, err := h.stockRepo.GetStock(ctx, ids)
+	if err != nil {
+		return err
+	}
+	idQuantityMap := make(map[string]int32)
+	for _, r := range records {
+		idQuantityMap[r.ID] += r.Quantity
+	}
+	var (
+		ok       = true
+		failedOn []struct {
+			ID   string
+			Want int32
+			Have int32
+		}
+	)
+	for _, item := range query {
+		if item.Quantity > idQuantityMap[item.ID] {
+			ok = false
+			failedOn = append(failedOn, struct {
+				ID   string
+				Want int32
+				Have int32
+			}{ID: item.ID, Want: item.Quantity, Have: idQuantityMap[item.ID]})
+		}
+	}
+	if ok {
+		return nil
+	}
+	return domain.ExceedStockError{FailedOn: failedOn}
+}
+
 func getStubPriceID(id string) string {
 	priceID, ok := stub[id]
 	if !ok {
diff --git a/internal/stock/app/query/get_items.go b/internal/stock/app/query/get_items.go
index 063a6b0..4de7bd5 100644
--- a/internal/stock/app/query/get_items.go
+++ b/internal/stock/app/query/get_items.go
@@ -4,8 +4,8 @@ import (
 	"context"
 
 	"github.com/ghost-yu/go_shop_second/common/decorator"
-	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
+	"github.com/ghost-yu/go_shop_second/stock/entity"
 	"github.com/sirupsen/logrus"
 )
 
@@ -13,7 +13,7 @@ type GetItems struct {
 	ItemIDs []string
 }
 
-type GetItemsHandler decorator.QueryHandler[GetItems, []*orderpb.Item]
+type GetItemsHandler decorator.QueryHandler[GetItems, []*entity.Item]
 
 type getItemsHandler struct {
 	stockRepo domain.Repository
@@ -27,14 +27,14 @@ func NewGetItemsHandler(
 	if stockRepo == nil {
 		panic("nil stockRepo")
 	}
-	return decorator.ApplyQueryDecorators[GetItems, []*orderpb.Item](
+	return decorator.ApplyQueryDecorators[GetItems, []*entity.Item](
 		getItemsHandler{stockRepo: stockRepo},
 		logger,
 		metricClient,
 	)
 }
 
-func (g getItemsHandler) Handle(ctx context.Context, query GetItems) ([]*orderpb.Item, error) {
+func (g getItemsHandler) Handle(ctx context.Context, query GetItems) ([]*entity.Item, error) {
 	items, err := g.stockRepo.GetItems(ctx, query.ItemIDs)
 	if err != nil {
 		return nil, err
diff --git a/internal/stock/domain/stock/repository.go b/internal/stock/domain/stock/repository.go
index 7a58e6a..c618633 100644
--- a/internal/stock/domain/stock/repository.go
+++ b/internal/stock/domain/stock/repository.go
@@ -5,11 +5,12 @@ import (
 	"fmt"
 	"strings"
 
-	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/stock/entity"
 )
 
 type Repository interface {
-	GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error)
+	GetItems(ctx context.Context, ids []string) ([]*entity.Item, error)
+	GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error)
 }
 
 type NotFoundError struct {
@@ -19,3 +20,19 @@ type NotFoundError struct {
 func (e NotFoundError) Error() string {
 	return fmt.Sprintf("these items not found in stock: %s", strings.Join(e.Missing, ","))
 }
+
+type ExceedStockError struct {
+	FailedOn []struct {
+		ID   string
+		Want int32
+		Have int32
+	}
+}
+
+func (e ExceedStockError) Error() string {
+	var info []string
+	for _, v := range e.FailedOn {
+		info = append(info, fmt.Sprintf("product_id=%s, want %d, have %d", v.ID, v.Want, v.Have))
+	}
+	return fmt.Sprintf("not enough stock for [%s]", strings.Join(info, ","))
+}
diff --git a/internal/stock/go.mod b/internal/stock/go.mod
index 41d06cf..e8dbf38 100644
--- a/internal/stock/go.mod
+++ b/internal/stock/go.mod
@@ -10,10 +10,11 @@ require (
 	github.com/spf13/viper v1.19.0
 	github.com/stripe/stripe-go/v79 v79.12.0
 	google.golang.org/grpc v1.67.1
+	gorm.io/driver/mysql v1.5.7
+	gorm.io/gorm v1.25.12
 )
 
 require (
-	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
 	github.com/armon/go-metrics v0.4.1 // indirect
 	github.com/bytedance/sonic v1.12.3 // indirect
 	github.com/bytedance/sonic/loader v0.2.0 // indirect
@@ -29,6 +30,7 @@ require (
 	github.com/go-playground/locales v0.14.1 // indirect
 	github.com/go-playground/universal-translator v0.18.1 // indirect
 	github.com/go-playground/validator/v10 v10.22.1 // indirect
+	github.com/go-sql-driver/mysql v1.7.0 // indirect
 	github.com/goccy/go-json v0.10.3 // indirect
 	github.com/golang/protobuf v1.5.3 // indirect
 	github.com/google/uuid v1.6.0 // indirect
@@ -43,6 +45,8 @@ require (
 	github.com/hashicorp/golang-lru v0.5.4 // indirect
 	github.com/hashicorp/hcl v1.0.0 // indirect
 	github.com/hashicorp/serf v0.10.1 // indirect
+	github.com/jinzhu/inflection v1.0.0 // indirect
+	github.com/jinzhu/now v1.1.5 // indirect
 	github.com/json-iterator/go v1.1.12 // indirect
 	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
 	github.com/leodido/go-urn v1.4.0 // indirect
@@ -53,7 +57,6 @@ require (
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
 	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
 	github.com/modern-go/reflect2 v1.0.2 // indirect
-	github.com/oapi-codegen/runtime v1.1.1 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
diff --git a/internal/stock/go.sum b/internal/stock/go.sum
index 17959ae..8af5f94 100644
--- a/internal/stock/go.sum
+++ b/internal/stock/go.sum
@@ -1,13 +1,10 @@
 cloud.google.com/go v0.26.0/go.mod h1:aQUYkXzVsufM+DwF1aE+0xfcU+56JwCaLick0ClmMTw=
 github.com/BurntSushi/toml v0.3.1/go.mod h1:xHWCNGjB5oqiDr8zfno3MHue2Ht5sIBksp03qcyfWMU=
 github.com/DataDog/datadog-go v3.2.0+incompatible/go.mod h1:LButxg5PwREeZtORoXG3tL4fMGNddJ+vMq1mwgfaqoQ=
-github.com/RaveNoX/go-jsoncommentstrip v1.0.0/go.mod h1:78ihd09MekBnJnxpICcwzCMzGrKSKYe4AqU6PDYYpjk=
 github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
 github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
 github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
 github.com/alecthomas/units v0.0.0-20190717042225-c3de453c63f4/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
-github.com/apapsch/go-jsonmerge/v2 v2.0.0 h1:axGnT1gRIfimI7gJifB699GoE/oq+F2MU7Dml6nw9rQ=
-github.com/apapsch/go-jsonmerge/v2 v2.0.0/go.mod h1:lvDnEdqiQrp0O42VQGgmlKpxL1AP2+08jFMw88y4klk=
 github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e/go.mod h1:3U/XgcO3hCbHZ8TKRvWD2dDTCfh9M9ya+I9JpbB7O8o=
 github.com/armon/go-metrics v0.0.0-20180917152333-f0300d1749da/go.mod h1:Q73ZrmVTwzkszR9V5SSuryQ31EELlFMUz1kKyl939pY=
 github.com/armon/go-metrics v0.4.1 h1:hR91U9KYmb6bLBYLQjyM+3j+rcd/UhE+G78SFnF8gJA=
@@ -19,7 +16,6 @@ github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973/go.mod h1:Dwedo/Wpr24
 github.com/beorn7/perks v1.0.0/go.mod h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+CedLV8=
 github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
 github.com/bgentry/speakeasy v0.1.0/go.mod h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=
-github.com/bmatcuk/doublestar v1.1.1/go.mod h1:UD6OnuiIn0yFxxA2le/rnRU1G4RaI4UvFv1sNto9p6w=
 github.com/bytedance/sonic v1.12.3 h1:W2MGa7RCU1QTeYRTPE3+88mVC0yXmsRQRChiyVocVjU=
 github.com/bytedance/sonic v1.12.3/go.mod h1:B8Gt/XvtZ3Fqj+iSKMypzymZxw/FVwgIGKzMzT9r/rk=
 github.com/bytedance/sonic/loader v0.1.1/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
@@ -77,6 +73,8 @@ github.com/go-playground/universal-translator v0.18.1 h1:Bcnm0ZwsGyWbCzImXv+pAJn
 github.com/go-playground/universal-translator v0.18.1/go.mod h1:xekY+UJKNuX9WP91TpwSH2VMlDf28Uj24BCp08ZFTUY=
 github.com/go-playground/validator/v10 v10.22.1 h1:40JcKH+bBNGFczGuoBYgX4I6m/i27HYW8P9FDk5PbgA=
 github.com/go-playground/validator/v10 v10.22.1/go.mod h1:dbuPbCMFw/DrkbEynArYaCwl3amGuJotoKCe95atGMM=
+github.com/go-sql-driver/mysql v1.7.0 h1:ueSltNNllEqE3qcWBTD0iQd3IpL/6U+mJxLkazJ7YPc=
+github.com/go-sql-driver/mysql v1.7.0/go.mod h1:OXbVy3sEdcQ2Doequ6Z5BW6fXNQTmx+9S1MCJN5yJMI=
 github.com/go-stack/stack v1.8.0/go.mod h1:v0f6uXyyMGvRgIKkXu+yp6POWl0qKG85gN/melR3HDY=
 github.com/goccy/go-json v0.10.3 h1:KZ5WoDbxAIgm2HNbYckL0se1fHD6rz5j4ywS6ebzDqA=
 github.com/goccy/go-json v0.10.3/go.mod h1:oq7eo15ShAhp70Anwd5lgX2pLfOS3QCiwU/PULtXL6M=
@@ -152,11 +150,14 @@ github.com/hashicorp/memberlist v0.5.0 h1:EtYPN8DpAURiapus508I4n9CzHs2W+8NZGbmmR
 github.com/hashicorp/memberlist v0.5.0/go.mod h1:yvyXLpo0QaGE59Y7hDTsTzDD25JYBZ4mHgHUZ8lrOI0=
 github.com/hashicorp/serf v0.10.1 h1:Z1H2J60yRKvfDYAOZLd2MU0ND4AH/WDz7xYHDWQsIPY=
 github.com/hashicorp/serf v0.10.1/go.mod h1:yL2t6BqATOLGc5HF7qbFkTfXoPIY0WZdWHfEvMqbG+4=
+github.com/jinzhu/inflection v1.0.0 h1:K317FqzuhWc8YvSVlFMCCUb36O/S9MCKRDI7QkRKD/E=
+github.com/jinzhu/inflection v1.0.0/go.mod h1:h+uFLlag+Qp1Va5pdKtLDYj+kHp5pxUVkryuEj+Srlc=
+github.com/jinzhu/now v1.1.5 h1:/o9tlHleP7gOFmsnYNz3RGnqzefHA47wQpKrrdTIwXQ=
+github.com/jinzhu/now v1.1.5/go.mod h1:d3SSVoowX0Lcu0IBviAWJpolVfI5UJVZZ7cO71lE/z8=
 github.com/json-iterator/go v1.1.6/go.mod h1:+SdeFBvtyEkXs7REEP0seUULqWtbJapLOCVDaaPEHmU=
 github.com/json-iterator/go v1.1.9/go.mod h1:KdQUCv79m/52Kvf8AW2vK1V8akMuk1QjK/uOdHXbAo4=
 github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
 github.com/json-iterator/go v1.1.12/go.mod h1:e30LSqwooZae/UwlEbR2852Gd8hjQvJoHmT4TnhNGBo=
-github.com/juju/gnuflag v0.0.0-20171113085948-2ce1bb71843d/go.mod h1:2PavIy+JPciBPrBUjwbNvtwB6RQlve+hkpll6QSNmOE=
 github.com/julienschmidt/httprouter v1.2.0/go.mod h1:SYymIcj16QtmaHHD7aYtjjsJG7VTCxuUUipMqKk8s4w=
 github.com/kisielk/errcheck v1.5.0/go.mod h1:pFxgyoBC7bSaBwPgfKdkLd5X25qrDl4LWUI2bnpBCr8=
 github.com/kisielk/gotool v1.0.0/go.mod h1:XhKaO+MFFWcvkIS/tQcRk01m1F5IRFswLeQ+oQHNcck=
@@ -210,8 +211,6 @@ github.com/modern-go/reflect2 v1.0.1/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3Rllmb
 github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
 github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
 github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
-github.com/oapi-codegen/runtime v1.1.1 h1:EXLHh0DXIJnWhdRPN2w4MXAzFyE4CskzhNLUmtpMYro=
-github.com/oapi-codegen/runtime v1.1.1/go.mod h1:SK9X900oXmPWilYR5/WKPzt3Kqxn/uS/+lbpREv+eCg=
 github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFStlNbqXla1AfSYxPUl2o=
 github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
@@ -262,7 +261,6 @@ github.com/spf13/pflag v1.0.5 h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=
 github.com/spf13/pflag v1.0.5/go.mod h1:McXfInJRrz4CZXVZOBLb0bTZqETkiAhM9Iw0y3An2Bg=
 github.com/spf13/viper v1.19.0 h1:RWq5SEjt8o25SROyN3z2OrDB9l7RPd3lwTWU8EcEdcI=
 github.com/spf13/viper v1.19.0/go.mod h1:GQUN9bilAbhU/jgc1bKs99f/suXKeUMct8Adx5+Ntkg=
-github.com/spkg/bom v0.0.0-20160624110644-59b7046e48ad/go.mod h1:qLr4V1qq6nMqFKkMo8ZTx3f+BZEkzsRUY10Xsm2mwU0=
 github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.1.1/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
@@ -439,6 +437,11 @@ gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C
 gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
 gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
 gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
+gorm.io/driver/mysql v1.5.7 h1:MndhOPYOfEp2rHKgkZIhJ16eVUIRf2HmzgoPmh7FCWo=
+gorm.io/driver/mysql v1.5.7/go.mod h1:sEtPWMiqiN1N1cMXoXmBbd8C6/l+TESwriotuRRpkDM=
+gorm.io/gorm v1.25.7/go.mod h1:hbnx/Oo0ChWMn1BIhpy1oYozzpM15i4YPuHDmfYtwg8=
+gorm.io/gorm v1.25.12 h1:I0u8i2hWQItBq1WfE0o2+WuL9+8L21K9e2HHSTE/0f8=
+gorm.io/gorm v1.25.12/go.mod h1:xh7N7RHfYlNc5EmcI/El95gXusucDrQnHXe0+CgWcLQ=
 honnef.co/go/tools v0.0.0-20190102054323-c2f93a96b099/go.mod h1:rf3lG4BRIbNafJWhAfAdb/ePZxsR/4RtNHQocxwk9r4=
 honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc/go.mod h1:rf3lG4BRIbNafJWhAfAdb/ePZxsR/4RtNHQocxwk9r4=
 nullprogram.com/x/optparse v1.0.0/go.mod h1:KdyPE+Igbe0jQUrVfMqDMeJQIJZEuyV7pjYmp6pbG50=
diff --git a/internal/stock/infrastructure/persistent/mysql.go b/internal/stock/infrastructure/persistent/mysql.go
new file mode 100644
index 0000000..777b770
--- /dev/null
+++ b/internal/stock/infrastructure/persistent/mysql.go
@@ -0,0 +1,53 @@
+package persistent
+
+import (
+	"context"
+	"fmt"
+	"time"
+
+	"github.com/sirupsen/logrus"
+	"github.com/spf13/viper"
+	"gorm.io/driver/mysql"
+	"gorm.io/gorm"
+)
+
+type MySQL struct {
+	db *gorm.DB
+}
+
+func NewMySQL() *MySQL {
+	dsn := fmt.Sprintf(
+		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
+		viper.GetString("mysql.user"),
+		viper.GetString("mysql.password"),
+		viper.GetString("mysql.host"),
+		viper.GetString("mysql.port"),
+		viper.GetString("mysql.dbname"),
+	)
+	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
+	if err != nil {
+		logrus.Panicf("connect to mysql failed, err=%v", err)
+	}
+	return &MySQL{db: db}
+}
+
+type StockModel struct {
+	ID        int64     `gorm:"column:id"`
+	ProductID string    `gorm:"column:product_id"`
+	Quantity  int32     `gorm:"column:quantity"`
+	CreatedAt time.Time `gorm:"column:created_at"`
+	UpdateAt  time.Time `gorm:"column:updated_at"`
+}
+
+func (StockModel) TableName() string {
+	return "o_stock"
+}
+
+func (d MySQL) BatchGetStockByID(ctx context.Context, productIDs []string) ([]StockModel, error) {
+	var result []StockModel
+	tx := d.db.WithContext(ctx).Where("product_id IN ?", productIDs).Find(&result)
+	if tx.Error != nil {
+		return nil, tx.Error
+	}
+	return result, nil
+}
diff --git a/internal/stock/ports/grpc.go b/internal/stock/ports/grpc.go
index b2c659d..51341e9 100644
--- a/internal/stock/ports/grpc.go
+++ b/internal/stock/ports/grpc.go
@@ -26,7 +26,7 @@ func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsReque
 	if err != nil {
 		return nil, err
 	}
-	return &stockpb.GetItemsResponse{Items: items}, nil
+	return &stockpb.GetItemsResponse{Items: convertor.NewItemConvertor().EntitiesToProtos(items)}, nil
 }
 
 func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
diff --git a/internal/stock/service/application.go b/internal/stock/service/application.go
index b51f2ad..5feb83a 100644
--- a/internal/stock/service/application.go
+++ b/internal/stock/service/application.go
@@ -8,11 +8,14 @@ import (
 	"github.com/ghost-yu/go_shop_second/stock/app"
 	"github.com/ghost-yu/go_shop_second/stock/app/query"
 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/integration"
+	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent"
 	"github.com/sirupsen/logrus"
 )
 
 func NewApplication(_ context.Context) app.Application {
-	stockRepo := adapters.NewMemoryStockRepository()
+	//stockRepo := adapters.NewMemoryStockRepository()
+	db := persistent.NewMySQL()
+	stockRepo := adapters.NewMySQLStockRepository(db)
 	logger := logrus.NewEntry(logrus.StandardLogger())
 	stripeAPI := integration.NewStripeAPI()
 	metricsClient := metrics.TodoMetrics{}
~~~

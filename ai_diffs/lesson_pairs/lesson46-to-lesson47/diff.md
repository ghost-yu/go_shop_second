# Lesson Pair Diff Report

- FromBranch: lesson46
- ToBranch: lesson47

## Short Summary

~~~text
 16 files changed, 328 insertions(+), 3 deletions(-)
~~~

## File Stats

~~~text
 .gitignore                                         |  1 +
 docker-compose.yml                                 | 11 +++-
 init.sql                                           |  3 ++
 internal/common/config/global.yaml                 | 11 ++++
 internal/common/go.mod                             |  3 ++
 internal/common/go.sum                             | 10 ++++
 internal/common/handler/factory/singleton.go       | 32 ++++++++++++
 internal/common/handler/redis/client.go            | 55 ++++++++++++++++++++
 internal/common/handler/redis/redis.go             | 60 ++++++++++++++++++++++
 internal/stock/adapters/stock_inmem_repository.go  | 10 ++++
 internal/stock/adapters/stock_mysql_repository.go  | 56 ++++++++++++++++++++
 .../stock/app/query/check_if_items_in_stock.go     | 53 ++++++++++++++++++-
 internal/stock/domain/stock/repository.go          |  9 ++++
 internal/stock/go.mod                              |  3 ++
 internal/stock/go.sum                              | 10 ++++
 internal/stock/infrastructure/persistent/mysql.go  |  4 ++
 16 files changed, 328 insertions(+), 3 deletions(-)
~~~

## Commit Comparison

~~~text
> a11f9d5 update stock
~~~

## Changed Files

~~~text
.gitignore
docker-compose.yml
init.sql
internal/common/config/global.yaml
internal/common/go.mod
internal/common/go.sum
internal/common/handler/factory/singleton.go
internal/common/handler/redis/client.go
internal/common/handler/redis/redis.go
internal/stock/adapters/stock_inmem_repository.go
internal/stock/adapters/stock_mysql_repository.go
internal/stock/app/query/check_if_items_in_stock.go
internal/stock/domain/stock/repository.go
internal/stock/go.mod
internal/stock/go.sum
internal/stock/infrastructure/persistent/mysql.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
.gitignore
docker-compose.yml
init.sql
internal/common/config/global.yaml
internal/common/handler/factory/singleton.go
internal/common/handler/redis/client.go
internal/common/handler/redis/redis.go
internal/stock/adapters/stock_inmem_repository.go
internal/stock/adapters/stock_mysql_repository.go
internal/stock/app/query/check_if_items_in_stock.go
internal/stock/domain/stock/repository.go
internal/stock/infrastructure/persistent/mysql.go
~~~

## Full Diff

~~~diff
diff --git a/.gitignore b/.gitignore
index 41649f8..3672a40 100644
--- a/.gitignore
+++ b/.gitignore
@@ -31,4 +31,5 @@ go.work.sum
 **/tmp/
 **/bin/
 
+data
 mysql_data
\ No newline at end of file
diff --git a/docker-compose.yml b/docker-compose.yml
index 6a3533b..49ea7b0 100644
--- a/docker-compose.yml
+++ b/docker-compose.yml
@@ -61,6 +61,15 @@ services:
     ports:
       - "3000:3000"
 
+  redis:
+    image: redis:latest
+    restart: on-failure
+    ports:
+      - "6379:6379"
+    volumes:
+      - ./data/redis_data:/data
+      - ./redis.conf:/usr/local/etc/redis/redis.conf
+
   mysql:
     image: mysql
     restart: on-failure
@@ -71,6 +80,6 @@ services:
       - MYSQL_PASSWORD=password
     volumes:
       - ./init.sql:/docker-entrypoint-initdb.d/init.sql
-      - ./mysql_data:/var/lib/mysql
+      - ./data/mysql_data:/var/lib/mysql
     ports:
       - "3307:3306"
\ No newline at end of file
diff --git a/init.sql b/init.sql
index ae95cc3..5585bc0 100644
--- a/init.sql
+++ b/init.sql
@@ -10,3 +10,6 @@ CREATE TABLE `o_stock` (
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
 ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
+
+INSERT INTO o_stock (product_id, quantity)
+VALUES ('prod_R3g7MikGYsXKzr', 1000), ('prod_R285C3Wb7FDprc', 500);
\ No newline at end of file
diff --git a/internal/common/config/global.yaml b/internal/common/config/global.yaml
index ddcd94d..4b0b19c 100644
--- a/internal/common/config/global.yaml
+++ b/internal/common/config/global.yaml
@@ -40,6 +40,17 @@ mongo:
   db-name: "order"
   coll-name: "order"
 
+redis:
+  local:
+    ip: 127.0.0.1
+    port: 6379
+    pool_size: 100
+    max_conn: 100
+    conn_timeout: 1000
+    read_timeout: 1000
+    write_timeout: 100
+
+
 mysql:
   user: root
   password: root
diff --git a/internal/common/go.mod b/internal/common/go.mod
index 9ed3597..f7ef6a8 100644
--- a/internal/common/go.mod
+++ b/internal/common/go.mod
@@ -8,6 +8,7 @@ require (
 	github.com/hashicorp/consul/api v1.28.2
 	github.com/oapi-codegen/runtime v1.1.1
 	github.com/rabbitmq/amqp091-go v1.10.0
+	github.com/redis/go-redis/v9 v9.7.0
 	github.com/sirupsen/logrus v1.8.1
 	github.com/spf13/viper v1.19.0
 	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.56.0
@@ -26,8 +27,10 @@ require (
 	github.com/armon/go-metrics v0.4.1 // indirect
 	github.com/bytedance/sonic v1.12.3 // indirect
 	github.com/bytedance/sonic/loader v0.2.0 // indirect
+	github.com/cespare/xxhash/v2 v2.3.0 // indirect
 	github.com/cloudwego/base64x v0.1.4 // indirect
 	github.com/cloudwego/iasm v0.2.0 // indirect
+	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
 	github.com/fatih/color v1.14.1 // indirect
 	github.com/fsnotify/fsnotify v1.7.0 // indirect
 	github.com/gabriel-vasile/mimetype v1.4.5 // indirect
diff --git a/internal/common/go.sum b/internal/common/go.sum
index 9b5cc9c..553c9f1 100644
--- a/internal/common/go.sum
+++ b/internal/common/go.sum
@@ -20,6 +20,10 @@ github.com/beorn7/perks v1.0.0/go.mod h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+Ce
 github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
 github.com/bgentry/speakeasy v0.1.0/go.mod h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=
 github.com/bmatcuk/doublestar v1.1.1/go.mod h1:UD6OnuiIn0yFxxA2le/rnRU1G4RaI4UvFv1sNto9p6w=
+github.com/bsm/ginkgo/v2 v2.12.0 h1:Ny8MWAHyOepLGlLKYmXG4IEkioBysk6GpaRTLC8zwWs=
+github.com/bsm/ginkgo/v2 v2.12.0/go.mod h1:SwYbGRRDovPVboqFv0tPTcG1sN61LM1Z4ARdbAV9g4c=
+github.com/bsm/gomega v1.27.10 h1:yeMWxP2pV2fG3FgAODIY8EiRE3dy0aeFYt4l7wh6yKA=
+github.com/bsm/gomega v1.27.10/go.mod h1:JyEr/xRbxbtgWNi8tIEVPUYZ5Dzef52k01W3YH0H+O0=
 github.com/bytedance/sonic v1.12.3 h1:W2MGa7RCU1QTeYRTPE3+88mVC0yXmsRQRChiyVocVjU=
 github.com/bytedance/sonic v1.12.3/go.mod h1:B8Gt/XvtZ3Fqj+iSKMypzymZxw/FVwgIGKzMzT9r/rk=
 github.com/bytedance/sonic/loader v0.1.1/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
@@ -27,6 +31,8 @@ github.com/bytedance/sonic/loader v0.2.0 h1:zNprn+lsIP06C/IqCHs3gPQIvnvpKbbxyXQP
 github.com/bytedance/sonic/loader v0.2.0/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
 github.com/census-instrumentation/opencensus-proto v0.2.1/go.mod h1:f6KPmirojxKA12rnyqOA5BBL4O983OfeGPqjHWSTneU=
 github.com/cespare/xxhash/v2 v2.1.1/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
+github.com/cespare/xxhash/v2 v2.3.0 h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=
+github.com/cespare/xxhash/v2 v2.3.0/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
 github.com/circonus-labs/circonus-gometrics v2.3.1+incompatible/go.mod h1:nmEj6Dob7S7YxXgwXpfOuvO54S+tGdZdw9fuRZt25Ag=
 github.com/circonus-labs/circonusllhist v0.1.3/go.mod h1:kMXHVDlOchFAehlya5ePtbp5jckzBHf4XRpQvBOLI+I=
 github.com/client9/misspell v0.3.4/go.mod h1:qj6jICC3Q7zFZvVWo7KLAzC3yx5G7kyvSDkc90ppPyw=
@@ -39,6 +45,8 @@ github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSs
 github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
 github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc h1:U9qPSI2PIWSS1VwoXQT9A3Wy9MM3WgvqSxFWenqJduM=
 github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
+github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f h1:lO4WD4F/rVNCu3HqELle0jiPLLBs70cWOduZpkS1E78=
+github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f/go.mod h1:cuUVRXasLTGF7a8hSLbxyZXjz+1KgoB3wDUb6vlszIc=
 github.com/envoyproxy/go-control-plane v0.9.0/go.mod h1:YTl/9mNaCwkRvm6d1a2C3ymFceY/DCBVvsKhRF0iEA4=
 github.com/envoyproxy/go-control-plane v0.9.1-0.20191026205805-5f8ba28d4473/go.mod h1:YTl/9mNaCwkRvm6d1a2C3ymFceY/DCBVvsKhRF0iEA4=
 github.com/envoyproxy/go-control-plane v0.9.4/go.mod h1:6rpuAdCZL397s3pYoYcLgu1mIlRU8Am5FuJP05cCM98=
@@ -241,6 +249,8 @@ github.com/prometheus/procfs v0.0.2/go.mod h1:TjEm7ze935MbeOT/UhFTIMYKhuLP4wbCsT
 github.com/prometheus/procfs v0.0.8/go.mod h1:7Qr8sr6344vo1JqZ6HhLceV9o3AJ1Ff+GxbHq6oeK9A=
 github.com/rabbitmq/amqp091-go v1.10.0 h1:STpn5XsHlHGcecLmMFCtg7mqq0RnD+zFr4uzukfVhBw=
 github.com/rabbitmq/amqp091-go v1.10.0/go.mod h1:Hy4jKW5kQART1u+JkDTF9YYOQUHXqMuhrgxOEeS7G4o=
+github.com/redis/go-redis/v9 v9.7.0 h1:HhLSs+B6O021gwzl+locl0zEDnyNkxMtf/Z3NNBMa9E=
+github.com/redis/go-redis/v9 v9.7.0/go.mod h1:f6zhXITC7JUJIlPEiBOTXxJgPLdZcA93GewI7inzyWw=
 github.com/rogpeppe/go-internal v1.9.0 h1:73kH8U+JUqXU8lRuOHeVHaa/SZPifC7BkcraZVejAe8=
 github.com/rogpeppe/go-internal v1.9.0/go.mod h1:WtVeX8xhTBvf0smdhujwtBcq4Qrzq/fJaraNFVN+nFs=
 github.com/ryanuber/columnize v0.0.0-20160712163229-9b3edd62028f/go.mod h1:sm1tb6uqfes/u+d4ooFouqFdy9/2g9QGwK3SQygK0Ts=
diff --git a/internal/common/handler/factory/singleton.go b/internal/common/handler/factory/singleton.go
new file mode 100644
index 0000000..96f2546
--- /dev/null
+++ b/internal/common/handler/factory/singleton.go
@@ -0,0 +1,32 @@
+package factory
+
+import "sync"
+
+type Supplier func(string) any
+
+type Singleton struct {
+	cache    map[string]any
+	locker   *sync.Mutex
+	supplier Supplier
+}
+
+func NewSingleton(supplier Supplier) *Singleton {
+	return &Singleton{
+		cache:    make(map[string]any),
+		locker:   &sync.Mutex{},
+		supplier: supplier,
+	}
+}
+
+func (s *Singleton) Get(key string) any {
+	if value, hit := s.cache[key]; hit {
+		return value
+	}
+	s.locker.Lock()
+	defer s.locker.Unlock()
+	if value, hit := s.cache[key]; hit {
+		return value
+	}
+	s.cache[key] = s.supplier(key)
+	return s.cache[key]
+}
diff --git a/internal/common/handler/redis/client.go b/internal/common/handler/redis/client.go
new file mode 100644
index 0000000..99c04d3
--- /dev/null
+++ b/internal/common/handler/redis/client.go
@@ -0,0 +1,55 @@
+package redis
+
+import (
+	"context"
+	"errors"
+	"time"
+
+	"github.com/redis/go-redis/v9"
+	"github.com/sirupsen/logrus"
+)
+
+func SetNX(ctx context.Context, client *redis.Client, key, value string, ttl time.Duration) (err error) {
+	now := time.Now()
+	defer func() {
+		l := logrus.WithContext(ctx).WithFields(logrus.Fields{
+			"start": now,
+			"key":   key,
+			"value": value,
+			"err":   err,
+			"cost":  time.Since(now).Milliseconds(),
+		})
+		if err == nil {
+			l.Info("redis_setnx_success")
+		} else {
+			l.Warn("redis_setnx_error")
+		}
+	}()
+	if client == nil {
+		return errors.New("redis client is nil")
+	}
+	_, err = client.SetNX(ctx, key, value, ttl).Result()
+	return err
+}
+
+func Del(ctx context.Context, client *redis.Client, key string) (err error) {
+	now := time.Now()
+	defer func() {
+		l := logrus.WithContext(ctx).WithFields(logrus.Fields{
+			"start": now,
+			"key":   key,
+			"err":   err,
+			"cost":  time.Since(now).Milliseconds(),
+		})
+		if err == nil {
+			l.Info("redis_del_success")
+		} else {
+			l.Warn("redis_del_error")
+		}
+	}()
+	if client == nil {
+		return errors.New("redis client is nil")
+	}
+	_, err = client.Del(ctx, key).Result()
+	return err
+}
diff --git a/internal/common/handler/redis/redis.go b/internal/common/handler/redis/redis.go
new file mode 100644
index 0000000..f76162e
--- /dev/null
+++ b/internal/common/handler/redis/redis.go
@@ -0,0 +1,60 @@
+package redis
+
+import (
+	"fmt"
+	"time"
+
+	"github.com/ghost-yu/go_shop_second/common/handler/factory"
+	"github.com/redis/go-redis/v9"
+	"github.com/spf13/viper"
+)
+
+const (
+	confName      = "redis"
+	localSupplier = "local"
+)
+
+var (
+	singleton = factory.NewSingleton(supplier)
+)
+
+func Init() {
+	conf := viper.GetStringMap(confName)
+	for supplyName := range conf {
+		Client(supplyName)
+	}
+}
+
+func LocalClient() *redis.Client {
+	return Client(localSupplier)
+}
+
+func Client(name string) *redis.Client {
+	return singleton.Get(name).(*redis.Client)
+}
+
+func supplier(key string) any {
+	confKey := confName + "." + key
+	type Section struct {
+		IP           string        `mapstructure:"ip"`
+		Port         string        `mapstructure:"port"`
+		PoolSize     int           `mapstructure:"pool_size"`
+		MaxConn      int           `mapstructure:"max_conn"`
+		ConnTimeout  time.Duration `mapstructure:"conn_timeout"`
+		ReadTimeout  time.Duration `mapstructure:"read_timeout"`
+		WriteTimeout time.Duration `mapstructure:"write_timeout"`
+	}
+	var c Section
+	if err := viper.UnmarshalKey(confKey, &c); err != nil {
+		panic(err)
+	}
+	return redis.NewClient(&redis.Options{
+		Network:         "tcp",
+		Addr:            fmt.Sprintf("%s:%s", c.IP, c.Port),
+		PoolSize:        c.PoolSize,
+		MaxActiveConns:  c.MaxConn,
+		ConnMaxLifetime: c.ConnTimeout * time.Millisecond,
+		ReadTimeout:     c.ReadTimeout * time.Millisecond,
+		WriteTimeout:    c.WriteTimeout * time.Millisecond,
+	})
+}
diff --git a/internal/stock/adapters/stock_inmem_repository.go b/internal/stock/adapters/stock_inmem_repository.go
index 2426756..16f48bd 100644
--- a/internal/stock/adapters/stock_inmem_repository.go
+++ b/internal/stock/adapters/stock_inmem_repository.go
@@ -13,6 +13,16 @@ type MemoryStockRepository struct {
 	store map[string]*entity.Item
 }
 
+func (m MemoryStockRepository) GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error) {
+	//TODO implement me
+	panic("implement me")
+}
+
+func (m MemoryStockRepository) UpdateStock(ctx context.Context, data []*entity.ItemWithQuantity, updateFn func(ctx context.Context, existing []*entity.ItemWithQuantity, query []*entity.ItemWithQuantity) ([]*entity.ItemWithQuantity, error)) error {
+	//TODO implement me
+	panic("implement me")
+}
+
 var stub = map[string]*entity.Item{
 	"item_id": {
 		ID:       "foo_item",
diff --git a/internal/stock/adapters/stock_mysql_repository.go b/internal/stock/adapters/stock_mysql_repository.go
index 91d1e50..7afb9fa 100644
--- a/internal/stock/adapters/stock_mysql_repository.go
+++ b/internal/stock/adapters/stock_mysql_repository.go
@@ -5,6 +5,8 @@ import (
 
 	"github.com/ghost-yu/go_shop_second/stock/entity"
 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent"
+	"github.com/sirupsen/logrus"
+	"gorm.io/gorm"
 )
 
 type MySQLStockRepository struct {
@@ -34,3 +36,57 @@ func (m MySQLStockRepository) GetStock(ctx context.Context, ids []string) ([]*en
 	}
 	return result, nil
 }
+
+func (m MySQLStockRepository) UpdateStock(
+	ctx context.Context,
+	data []*entity.ItemWithQuantity,
+	updateFn func(
+		ctx context.Context,
+		existing []*entity.ItemWithQuantity,
+		query []*entity.ItemWithQuantity,
+	) ([]*entity.ItemWithQuantity, error),
+) error {
+	return m.db.StartTransaction(func(tx *gorm.DB) (err error) {
+		defer func() {
+			if err != nil {
+				logrus.Warnf("update stock transaction err=%v", err)
+			}
+		}()
+		var dest []*persistent.StockModel
+		if err = tx.Table("o_stock").Where("product_id IN ?", getIDFromEntities(data)).Find(&dest).Error; err != nil {
+			return err
+		}
+		existing := m.unmarshalFromDatabase(dest)
+
+		updated, err := updateFn(ctx, existing, data)
+		if err != nil {
+			return err
+		}
+
+		for _, upd := range updated {
+			if err = tx.Table("o_stock").Where("product_id = ?", upd.ID).Update("quantity", upd.Quantity).Error; err != nil {
+				return err
+			}
+		}
+		return nil
+	})
+}
+
+func (m MySQLStockRepository) unmarshalFromDatabase(dest []*persistent.StockModel) []*entity.ItemWithQuantity {
+	var result []*entity.ItemWithQuantity
+	for _, i := range dest {
+		result = append(result, &entity.ItemWithQuantity{
+			ID:       i.ProductID,
+			Quantity: i.Quantity,
+		})
+	}
+	return result
+}
+
+func getIDFromEntities(items []*entity.ItemWithQuantity) []string {
+	var ids []string
+	for _, i := range items {
+		ids = append(ids, i.ID)
+	}
+	return ids
+}
diff --git a/internal/stock/app/query/check_if_items_in_stock.go b/internal/stock/app/query/check_if_items_in_stock.go
index 0e86e86..5c48a2b 100644
--- a/internal/stock/app/query/check_if_items_in_stock.go
+++ b/internal/stock/app/query/check_if_items_in_stock.go
@@ -2,14 +2,21 @@ package query
 
 import (
 	"context"
+	"strings"
+	"time"
 
 	"github.com/ghost-yu/go_shop_second/common/decorator"
+	"github.com/ghost-yu/go_shop_second/common/handler/redis"
 	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
 	"github.com/ghost-yu/go_shop_second/stock/entity"
 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/integration"
 	"github.com/sirupsen/logrus"
 )
 
+const (
+	redisLockPrefix = "check_stock_"
+)
+
 type CheckIfItemsInStock struct {
 	Items []*entity.ItemWithQuantity
 }
@@ -50,9 +57,15 @@ var stub = map[string]string{
 }
 
 func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*entity.Item, error) {
-	if err := h.checkStock(ctx, query.Items); err != nil {
+	if err := lock(ctx, getLockKey(query)); err != nil {
 		return nil, err
 	}
+	defer func() {
+		if err := unlock(ctx, getLockKey(query)); err != nil {
+			logrus.Warnf("redis unlock fail, err=%v", err)
+		}
+	}()
+
 	var res []*entity.Item
 	for _, i := range query.Items {
 		priceID, err := h.stripeAPI.GetPriceByProductID(ctx, i.ID)
@@ -66,9 +79,28 @@ func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfIte
 		})
 	}
 	// TODO: 扣库存
+	if err := h.checkStock(ctx, query.Items); err != nil {
+		return nil, err
+	}
 	return res, nil
 }
 
+func getLockKey(query CheckIfItemsInStock) string {
+	var ids []string
+	for _, i := range query.Items {
+		ids = append(ids, i.ID)
+	}
+	return redisLockPrefix + strings.Join(ids, "_")
+}
+
+func unlock(ctx context.Context, key string) error {
+	return redis.Del(ctx, redis.LocalClient(), key)
+}
+
+func lock(ctx context.Context, key string) error {
+	return redis.SetNX(ctx, redis.LocalClient(), key, "1", 5*time.Minute)
+}
+
 func (h checkIfItemsInStockHandler) checkStock(ctx context.Context, query []*entity.ItemWithQuantity) error {
 	var ids []string
 	for _, i := range query {
@@ -101,7 +133,24 @@ func (h checkIfItemsInStockHandler) checkStock(ctx context.Context, query []*ent
 		}
 	}
 	if ok {
-		return nil
+		return h.stockRepo.UpdateStock(ctx, query, func(
+			ctx context.Context,
+			existing []*entity.ItemWithQuantity,
+			query []*entity.ItemWithQuantity,
+		) ([]*entity.ItemWithQuantity, error) {
+			var newItems []*entity.ItemWithQuantity
+			for _, e := range existing {
+				for _, q := range query {
+					if e.ID == q.ID {
+						newItems = append(newItems, &entity.ItemWithQuantity{
+							ID:       e.ID,
+							Quantity: e.Quantity - q.Quantity,
+						})
+					}
+				}
+			}
+			return newItems, nil
+		})
 	}
 	return domain.ExceedStockError{FailedOn: failedOn}
 }
diff --git a/internal/stock/domain/stock/repository.go b/internal/stock/domain/stock/repository.go
index c618633..ea9d7c3 100644
--- a/internal/stock/domain/stock/repository.go
+++ b/internal/stock/domain/stock/repository.go
@@ -11,6 +11,15 @@ import (
 type Repository interface {
 	GetItems(ctx context.Context, ids []string) ([]*entity.Item, error)
 	GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error)
+	UpdateStock(
+		ctx context.Context,
+		data []*entity.ItemWithQuantity,
+		updateFn func(
+			ctx context.Context,
+			existing []*entity.ItemWithQuantity,
+			query []*entity.ItemWithQuantity,
+		) ([]*entity.ItemWithQuantity, error),
+	) error
 }
 
 type NotFoundError struct {
diff --git a/internal/stock/go.mod b/internal/stock/go.mod
index e8dbf38..3211907 100644
--- a/internal/stock/go.mod
+++ b/internal/stock/go.mod
@@ -18,8 +18,10 @@ require (
 	github.com/armon/go-metrics v0.4.1 // indirect
 	github.com/bytedance/sonic v1.12.3 // indirect
 	github.com/bytedance/sonic/loader v0.2.0 // indirect
+	github.com/cespare/xxhash/v2 v2.3.0 // indirect
 	github.com/cloudwego/base64x v0.1.4 // indirect
 	github.com/cloudwego/iasm v0.2.0 // indirect
+	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
 	github.com/fatih/color v1.14.1 // indirect
 	github.com/fsnotify/fsnotify v1.7.0 // indirect
 	github.com/gabriel-vasile/mimetype v1.4.5 // indirect
@@ -58,6 +60,7 @@ require (
 	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
 	github.com/modern-go/reflect2 v1.0.2 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
+	github.com/redis/go-redis/v9 v9.7.0 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
 	github.com/sourcegraph/conc v0.3.0 // indirect
diff --git a/internal/stock/go.sum b/internal/stock/go.sum
index 8af5f94..cde2c60 100644
--- a/internal/stock/go.sum
+++ b/internal/stock/go.sum
@@ -16,6 +16,10 @@ github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973/go.mod h1:Dwedo/Wpr24
 github.com/beorn7/perks v1.0.0/go.mod h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+CedLV8=
 github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
 github.com/bgentry/speakeasy v0.1.0/go.mod h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=
+github.com/bsm/ginkgo/v2 v2.12.0 h1:Ny8MWAHyOepLGlLKYmXG4IEkioBysk6GpaRTLC8zwWs=
+github.com/bsm/ginkgo/v2 v2.12.0/go.mod h1:SwYbGRRDovPVboqFv0tPTcG1sN61LM1Z4ARdbAV9g4c=
+github.com/bsm/gomega v1.27.10 h1:yeMWxP2pV2fG3FgAODIY8EiRE3dy0aeFYt4l7wh6yKA=
+github.com/bsm/gomega v1.27.10/go.mod h1:JyEr/xRbxbtgWNi8tIEVPUYZ5Dzef52k01W3YH0H+O0=
 github.com/bytedance/sonic v1.12.3 h1:W2MGa7RCU1QTeYRTPE3+88mVC0yXmsRQRChiyVocVjU=
 github.com/bytedance/sonic v1.12.3/go.mod h1:B8Gt/XvtZ3Fqj+iSKMypzymZxw/FVwgIGKzMzT9r/rk=
 github.com/bytedance/sonic/loader v0.1.1/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
@@ -23,6 +27,8 @@ github.com/bytedance/sonic/loader v0.2.0 h1:zNprn+lsIP06C/IqCHs3gPQIvnvpKbbxyXQP
 github.com/bytedance/sonic/loader v0.2.0/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
 github.com/census-instrumentation/opencensus-proto v0.2.1/go.mod h1:f6KPmirojxKA12rnyqOA5BBL4O983OfeGPqjHWSTneU=
 github.com/cespare/xxhash/v2 v2.1.1/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
+github.com/cespare/xxhash/v2 v2.3.0 h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=
+github.com/cespare/xxhash/v2 v2.3.0/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
 github.com/circonus-labs/circonus-gometrics v2.3.1+incompatible/go.mod h1:nmEj6Dob7S7YxXgwXpfOuvO54S+tGdZdw9fuRZt25Ag=
 github.com/circonus-labs/circonusllhist v0.1.3/go.mod h1:kMXHVDlOchFAehlya5ePtbp5jckzBHf4XRpQvBOLI+I=
 github.com/client9/misspell v0.3.4/go.mod h1:qj6jICC3Q7zFZvVWo7KLAzC3yx5G7kyvSDkc90ppPyw=
@@ -35,6 +41,8 @@ github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSs
 github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
 github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc h1:U9qPSI2PIWSS1VwoXQT9A3Wy9MM3WgvqSxFWenqJduM=
 github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
+github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f h1:lO4WD4F/rVNCu3HqELle0jiPLLBs70cWOduZpkS1E78=
+github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f/go.mod h1:cuUVRXasLTGF7a8hSLbxyZXjz+1KgoB3wDUb6vlszIc=
 github.com/envoyproxy/go-control-plane v0.9.0/go.mod h1:YTl/9mNaCwkRvm6d1a2C3ymFceY/DCBVvsKhRF0iEA4=
 github.com/envoyproxy/go-control-plane v0.9.1-0.20191026205805-5f8ba28d4473/go.mod h1:YTl/9mNaCwkRvm6d1a2C3ymFceY/DCBVvsKhRF0iEA4=
 github.com/envoyproxy/go-control-plane v0.9.4/go.mod h1:6rpuAdCZL397s3pYoYcLgu1mIlRU8Am5FuJP05cCM98=
@@ -238,6 +246,8 @@ github.com/prometheus/common v0.9.1/go.mod h1:yhUN8i9wzaXS3w1O07YhxHEBxD+W35wd8b
 github.com/prometheus/procfs v0.0.0-20181005140218-185b4288413d/go.mod h1:c3At6R/oaqEKCNdg8wHV1ftS6bRYblBhIjjI8uT2IGk=
 github.com/prometheus/procfs v0.0.2/go.mod h1:TjEm7ze935MbeOT/UhFTIMYKhuLP4wbCsTZCD3I8kEA=
 github.com/prometheus/procfs v0.0.8/go.mod h1:7Qr8sr6344vo1JqZ6HhLceV9o3AJ1Ff+GxbHq6oeK9A=
+github.com/redis/go-redis/v9 v9.7.0 h1:HhLSs+B6O021gwzl+locl0zEDnyNkxMtf/Z3NNBMa9E=
+github.com/redis/go-redis/v9 v9.7.0/go.mod h1:f6zhXITC7JUJIlPEiBOTXxJgPLdZcA93GewI7inzyWw=
 github.com/rogpeppe/go-internal v1.9.0 h1:73kH8U+JUqXU8lRuOHeVHaa/SZPifC7BkcraZVejAe8=
 github.com/rogpeppe/go-internal v1.9.0/go.mod h1:WtVeX8xhTBvf0smdhujwtBcq4Qrzq/fJaraNFVN+nFs=
 github.com/ryanuber/columnize v0.0.0-20160712163229-9b3edd62028f/go.mod h1:sm1tb6uqfes/u+d4ooFouqFdy9/2g9QGwK3SQygK0Ts=
diff --git a/internal/stock/infrastructure/persistent/mysql.go b/internal/stock/infrastructure/persistent/mysql.go
index 777b770..33637f3 100644
--- a/internal/stock/infrastructure/persistent/mysql.go
+++ b/internal/stock/infrastructure/persistent/mysql.go
@@ -43,6 +43,10 @@ func (StockModel) TableName() string {
 	return "o_stock"
 }
 
+func (d MySQL) StartTransaction(fc func(tx *gorm.DB) error) error {
+	return d.db.Transaction(fc)
+}
+
 func (d MySQL) BatchGetStockByID(ctx context.Context, productIDs []string) ([]StockModel, error) {
 	var result []StockModel
 	tx := d.db.WithContext(ctx).Where("product_id IN ?", productIDs).Find(&result)
~~~

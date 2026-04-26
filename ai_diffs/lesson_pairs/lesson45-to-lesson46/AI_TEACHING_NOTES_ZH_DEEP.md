# `lesson45 -> lesson46` 独立讲义（详细注释版）

这一组开始出现一个非常重要的变化：

`stock 不再只是“有个 MySQL 容器在那儿”，而是第一次真正让库存校验依赖 MySQL 里的真实库存数据。`

如果上一节 `lesson44 -> lesson45` 是“把数据库环境搭起来”，那这一节就是：

- 把数据库配置接进服务
- 把 stock 仓储切到 MySQL 版本
- 在创建订单之前真正检查库存够不够
- 在 HTTP / OpenAPI 两层都提前拦住明显非法的数量输入

所以这组是一个很典型的“从基础设施铺路进入真实业务校验”的转折点。

## 1. 正确阅读顺序

建议按这个顺序看：

1. [api/openapi/order.yml](/g:/shi/go_shop_second/api/openapi/order.yml)
2. [internal/order/http.go](/g:/shi/go_shop_second/internal/order/http.go)
3. [internal/common/config/global.yaml](/g:/shi/go_shop_second/internal/common/config/global.yaml)
4. [internal/stock/domain/stock/repository.go](/g:/shi/go_shop_second/internal/stock/domain/stock/repository.go)
5. [internal/stock/infrastructure/persistent/mysql.go](/g:/shi/go_shop_second/internal/stock/infrastructure/persistent/mysql.go)
6. [internal/stock/adapters/stock_mysql_repository.go](/g:/shi/go_shop_second/internal/stock/adapters/stock_mysql_repository.go)
7. [internal/stock/service/application.go](/g:/shi/go_shop_second/internal/stock/service/application.go)
8. [internal/stock/app/query/check_if_items_in_stock.go](/g:/shi/go_shop_second/internal/stock/app/query/check_if_items_in_stock.go)
9. [internal/stock/ports/grpc.go](/g:/shi/go_shop_second/internal/stock/ports/grpc.go)
10. [internal/stock/app/query/get_items.go](/g:/shi/go_shop_second/internal/stock/app/query/get_items.go)
11. [internal/stock/adapters/stock_inmem_repository.go](/g:/shi/go_shop_second/internal/stock/adapters/stock_inmem_repository.go)

这样看最清楚，因为这组其实有两条线：

1. 入口层参数校验
2. stock 服务内部真正读 MySQL 做库存校验

## 2. 这组到底在干什么

一句话概括：

`订单创建链第一次从“只检查格式”进化到“真正检查库存是否够”。`

前面课程里，`stock` 虽然已经能回商品信息、能查 Stripe price ID，但“库存数量够不够”这件事还没真正落地。

这节开始：

- MySQL 配置被接入服务
- `Repository` 接口新增 `GetStock(...)`
- MySQL 仓储实现开始查询 `o_stock`
- `check_if_items_in_stock` 先查库存，不够就直接报错

同时还多了一层很实用的防御：

- OpenAPI schema 加 `minimum: 1`
- `order/http.go` 再做运行时校验

这其实体现了一个很重要的工程思想：

`输入校验不要只做一层。`

## 3. 总调用链

这一组之后，创建订单的大致链路变成：

`客户端发创建订单请求 -> OpenAPI 契约要求 quantity >= 1 -> order/http.go 再次检查 quantity > 0 -> order 命令调用 stock.CheckIfItemsInStock -> stock gRPC 层把 proto 转成 entity -> stock query 先查 MySQL 库存 -> 库存够才继续查 Stripe priceID -> 返回给 order`

你要特别注意顺序：

不是先查 Stripe 再查库存，而是：

1. 先看库存够不够
2. 再补商品对应的 Stripe 价格 ID

这样做更合理，因为如果库存都不够，就没必要继续做后面的外部依赖调用。

## 4. 关键文件一：`api/openapi/order.yml`

### 4.1 这个文件自己的原始 diff

```diff
diff --git a/api/openapi/order.yml b/api/openapi/order.yml
index 6155cc0..9b198f2 100644
--- a/api/openapi/order.yml
+++ b/api/openapi/order.yml
@@ -145,6 +145,7 @@ components:
         quantity:
           type: integer
           format: int32
+          minimum: 1
```

### 4.2 这段改动在干什么

它给 `CreateOrderRequest.items[].quantity` 加了一条契约限制：

`数量最小值必须是 1。`

这看起来只有一行，但意义很大。因为它是在 API 契约层直接说明：

- 0 不合法
- 负数不合法

也就是说，这个限制不再只是你服务内部“心里这么想”，而是正式写进对外协议了。

### 4.3 带中文注释的关键配置

```yaml
ItemWithQuantity:
  type: object
  required:
    - id
    - quantity
  properties:
    id:
      type: string
    quantity:
      type: integer
      format: int32
      minimum: 1
      # minimum: 1 的意思是：这个字段的数值不能小于 1。
      # 也就是 0、-1、-100 这类值都不符合契约。
```

### 4.4 `minimum` 到底是什么

它是 OpenAPI schema 的一个数值约束关键词。

作用是：

- 给数字字段加下限

为什么这里值得专门讲：

因为很多新手会把 OpenAPI 理解成“只是文档”。

其实不是。它通常还会影响：

- 代码生成
- 客户端 SDK
- 测试数据校验
- 文档展示
- 某些框架里的自动校验

所以这一行不是“注释增强”，而是正式契约增强。

## 5. 关键文件二：`internal/order/http.go`

### 5.1 这个文件自己的原始 diff

```diff
diff --git a/internal/order/http.go b/internal/order/http.go
index a7cec4b..a8352ff 100644
--- a/internal/order/http.go
+++ b/internal/order/http.go
@@ -1,6 +1,7 @@
 package main
 
 import (
+	"errors"
 	"fmt"
@@
 	if err = c.ShouldBindJSON(&req); err != nil {
 		return
 	}
+	if err = H.validate(req); err != nil {
+		return
+	}
@@
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
```

### 5.2 旧代码和新代码的区别

旧代码：

- 只做 `ShouldBindJSON(...)`
- 只要 JSON 结构能绑定成功，就继续往下跑

新代码：

- 绑定 JSON 后，立刻额外做一次业务前置校验
- 只要发现有商品数量 `<= 0`，直接返回错误

也就是说，这节新增的是：

`运行时业务校验，而不只是结构绑定。`

### 5.3 带中文注释的关键代码

```go
func (H HTTPServer) PostCustomerCustomerIdOrders(c *gin.Context, customerID string) {
    var (
        req  client.CreateOrderRequest
        resp dto.CreateOrderResponse
        err  error
    )
    defer func() {
        H.Response(c, err, &resp)
    }()

    // ShouldBindJSON 会尝试把请求体 JSON 绑定到 req 结构体。
    // 它解决的是“格式能不能解析”，不是“业务值合不合理”。
    if err = c.ShouldBindJSON(&req); err != nil {
        return
    }

    // 新增的 validate 解决的是业务约束问题：
    // 比如 quantity 能不能为 0、能不能为负数。
    if err = H.validate(req); err != nil {
        return
    }

    r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
        CustomerID: req.CustomerId,
        Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
    })
    if err != nil {
        return
    }

    resp = dto.CreateOrderResponse{ ... }
}

func (H HTTPServer) validate(req client.CreateOrderRequest) error {
    for _, v := range req.Items {
        if v.Quantity <= 0 {
            // errors.New 用于快速创建一个普通 error。
            return errors.New("quantity must be positive")
        }
    }
    return nil
}
```

### 5.4 为什么 OpenAPI 已经加了 `minimum: 1`，这里还要再校验一次

这是一个非常值得你形成直觉的工程问题。

答案是：

`因为外部契约校验不能替代服务内部防御。`

原因有几个：

1. 客户端可能绕过你预期的 SDK
2. 某些生成代码或网关层不一定帮你挡住所有情况
3. 服务内部自己做一次校验，才最稳

所以这一组体现的不是“重复代码”，而是：

`契约层 + 服务层 双保险。`

## 6. 关键文件三：`internal/common/config/global.yaml`

### 6.1 这个文件自己的原始 diff

```diff
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
```

### 6.2 这段改动在干什么

它把上一组 docker-compose 里加进来的 MySQL 实例，正式变成了服务配置的一部分。

没有这一步，MySQL 容器就只是“跑着”，服务代码根本不知道去哪连它。

### 6.3 带中文注释的配置

```yaml
mysql:
  user: root
  password: root
  host: localhost
  port: 3307
  dbname: "gorder_v2"
```

这里每个字段的作用都很直接：

- `user`：数据库用户名
- `password`：数据库密码
- `host`：数据库地址
- `port`：数据库端口
- `dbname`：要连接的数据库名

你要把这和上一组的 Docker Compose 对起来看：

- Compose 暴露的是宿主机 `3307`
- 这里服务配置就连 `localhost:3307`

这就是配置和基础设施真正对上的地方。

## 7. 关键文件四：`internal/stock/domain/stock/repository.go`

### 7.1 这个文件自己的原始 diff

```diff
diff --git a/internal/stock/domain/stock/repository.go b/internal/stock/domain/stock/repository.go
index 7a58e6a..c618633 100644
--- a/internal/stock/domain/stock/repository.go
+++ b/internal/stock/domain/stock/repository.go
@@ -5,11 +5,12 @@ import (
-	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	"github.com/ghost-yu/go_shop_second/stock/entity"
 )
 
 type Repository interface {
-	GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error)
+	GetItems(ctx context.Context, ids []string) ([]*entity.Item, error)
+	GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error)
 }
@@
+type ExceedStockError struct {
+	FailedOn []struct {
+		ID   string
+		Want int32
+		Have int32
+	}
+}
```

### 7.2 这段改动的意义

这里有两个关键变化：

1. `Repository` 接口新增了 `GetStock(...)`
2. 新增了 `ExceedStockError`

这说明 domain 层正式承认一件事：

`库存检查不是“顺手查一下商品”，而是一个独立能力。`

### 7.3 带中文注释的关键代码

```go
type Repository interface {
    // GetItems 更偏向“拿商品详情”。
    GetItems(ctx context.Context, ids []string) ([]*entity.Item, error)

    // GetStock 明确表示“拿库存数量”。
    // 这说明商品信息和库存信息开始在语义上分开。
    GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error)
}

type ExceedStockError struct {
    FailedOn []struct {
        ID   string
        Want int32
        Have int32
    }
}

func (e ExceedStockError) Error() string {
    var info []string
    for _, v := range e.FailedOn {
        info = append(info, fmt.Sprintf("product_id=%s, want %d, have %d", v.ID, v.Want, v.Have))
    }
    return fmt.Sprintf("not enough stock for [%s]", strings.Join(info, ","))
}
```

### 7.4 为什么这里要专门定义 `ExceedStockError`

因为普通错误虽然也能报“库存不足”，但它不够结构化。

这里把失败详情收集成：

- 哪个商品 ID
- 想要多少
- 实际只有多少

这样后面不管是日志、调试、还是前端提示，信息都会更有用。

这也是领域错误比裸字符串更有价值的地方。

## 8. 关键文件五：`internal/stock/infrastructure/persistent/mysql.go`

### 8.1 这个文件自己的原始 diff

```diff
diff --git a/internal/stock/infrastructure/persistent/mysql.go b/internal/stock/infrastructure/persistent/mysql.go
new file mode 100644
index 0000000..777b770
--- /dev/null
+++ b/internal/stock/infrastructure/persistent/mysql.go
@@ -0,0 +1,53 @@
+package persistent
+...
+func NewMySQL() *MySQL {
+	dsn := fmt.Sprintf(
+		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
+		viper.GetString("mysql.user"),
+		...
+	)
+	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
+	if err != nil {
+		logrus.Panicf("connect to mysql failed, err=%v", err)
+	}
+	return &MySQL{db: db}
+}
+...
+func (d MySQL) BatchGetStockByID(ctx context.Context, productIDs []string) ([]StockModel, error) {
+	var result []StockModel
+	tx := d.db.WithContext(ctx).Where("product_id IN ?", productIDs).Find(&result)
+	...
+}
```

### 8.2 这段代码在项目里做什么

这是 stock 服务第一次真正接入 MySQL 持久层的入口。

职责是：

1. 根据配置拼 MySQL DSN
2. 用 GORM 建立数据库连接
3. 定义数据库表对应的 Go model
4. 提供按商品 ID 批量查库存的方法

### 8.3 带中文注释的关键代码

```go
package persistent

import (
    "context"
    "fmt"
    "time"

    "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type MySQL struct {
    db *gorm.DB
}

func NewMySQL() *MySQL {
    // DSN 是数据库连接串。
    // 它把用户名、密码、地址、端口、数据库名这些信息拼成一个标准格式字符串。
    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        viper.GetString("mysql.user"),
        viper.GetString("mysql.password"),
        viper.GetString("mysql.host"),
        viper.GetString("mysql.port"),
        viper.GetString("mysql.dbname"),
    )

    // gorm.Open(driver, config) 会创建 GORM 数据库实例。
    // mysql.Open(dsn) 是 MySQL 驱动适配器。
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        // Panicf 会先格式化输出，再 panic。
        // 这里表示：数据库连不上，stock 服务就没法正常工作。
        logrus.Panicf("connect to mysql failed, err=%v", err)
    }
    return &MySQL{db: db}
}

type StockModel struct {
    ID        int64     `gorm:"column:id"`
    ProductID string    `gorm:"column:product_id"`
    Quantity  int32     `gorm:"column:quantity"`
    CreatedAt time.Time `gorm:"column:created_at"`
    UpdateAt  time.Time `gorm:"column:updated_at"`
}

func (StockModel) TableName() string {
    return "o_stock"
}

func (d MySQL) BatchGetStockByID(ctx context.Context, productIDs []string) ([]StockModel, error) {
    var result []StockModel

    // WithContext(ctx) 把 context 传给 GORM，便于超时、取消、trace 等上下文继续向下传递。
    // Where("product_id IN ?", productIDs) 会生成 SQL 的 IN 查询。
    // Find(&result) 把查询结果扫进 result 切片。
    tx := d.db.WithContext(ctx).Where("product_id IN ?", productIDs).Find(&result)
    if tx.Error != nil {
        return nil, tx.Error
    }
    return result, nil
}
```

### 8.4 这里必须拆开的库函数 / ORM 概念

#### `gorm.Open(...)` 是什么

来源：GORM

作用：

- 创建数据库访问入口

你可以先把 `*gorm.DB` 理解成：

`一个带 ORM 能力的数据库操作对象。`

#### `mysql.Open(dsn)` 是什么

来源：GORM 的 MySQL 驱动包 `gorm.io/driver/mysql`

作用：

- 告诉 GORM：底层数据库是 MySQL，连接串是这个 dsn

#### DSN 是什么

DSN 全称 Data Source Name。

你可以把它理解成：

`数据库连接地址 + 连接参数的大字符串。`

这里这条 DSN 里最值得你记住的是：

- `charset=utf8mb4`
- `parseTime=True`
- `loc=Local`

##### `parseTime=True` 为什么值得讲

如果不加它，MySQL 里某些时间字段可能不会自动按你预期解析成 Go 的 `time.Time`。

所以这是 Go 连 MySQL 时一个很常见的小坑。

#### `WithContext(ctx)` 是什么

作用：

- 把当前 context 传给这次数据库操作

这样做的好处是：

- 如果上游请求取消了，这里也能感知
- trace / deadline 也更容易往下传

#### `Where("product_id IN ?", productIDs)` 是什么

这是 GORM 的条件查询写法。

它会大致生成类似：

```sql
WHERE product_id IN (...)
```

初学者这里最容易紧张的是：

- 这样写会不会 SQL 注入？

在这种参数绑定风格下，ORM 会帮你做参数化处理，一般比自己字符串拼 SQL 稳得多。

## 9. 关键文件六：`internal/stock/adapters/stock_mysql_repository.go`

### 9.1 这个文件自己的原始 diff

```diff
diff --git a/internal/stock/adapters/stock_mysql_repository.go b/internal/stock/adapters/stock_mysql_repository.go
new file mode 100644
index 0000000..91d1e50
--- /dev/null
+++ b/internal/stock/adapters/stock_mysql_repository.go
@@ -0,0 +1,36 @@
+package adapters
+...
+type MySQLStockRepository struct {
+	db *persistent.MySQL
+}
+...
+func (m MySQLStockRepository) GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error) {
+	data, err := m.db.BatchGetStockByID(ctx, ids)
+	...
+}
```

### 9.2 这段代码在做什么

它把持久层查询结果，转换成 domain / app 层可用的 `entity.ItemWithQuantity`。

这一步非常典型：

- `persistent.MySQL` 更靠近数据库
- `adapter` 负责把数据库模型适配成上层仓储接口要的对象

### 9.3 带中文注释的关键代码

```go
type MySQLStockRepository struct {
    db *persistent.MySQL
}

func NewMySQLStockRepository(db *persistent.MySQL) *MySQLStockRepository {
    return &MySQLStockRepository{db: db}
}

func (m MySQLStockRepository) GetItems(ctx context.Context, ids []string) ([]*entity.Item, error) {
    // 这里还没实现，说明这节重点不是“完整切完 MySQL 仓储”，
    // 而是先把库存校验那条最急的链打通。
    panic("implement me")
}

func (m MySQLStockRepository) GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error) {
    data, err := m.db.BatchGetStockByID(ctx, ids)
    if err != nil {
        return nil, err
    }

    var result []*entity.ItemWithQuantity
    for _, d := range data {
        result = append(result, &entity.ItemWithQuantity{
            ID:       d.ProductID,
            Quantity: d.Quantity,
        })
    }
    return result, nil
}
```

### 9.4 这里最值得你注意的点

`GetItems(...)` 还没实现，直接 `panic("implement me")`。

这说明作者这节的推进方式是：

- 先把“库存够不够”这条线做起来
- 商品详情那条线以后再补

这很像真实开发中的“按最短关键路径推进”。

## 10. 关键文件七：`internal/stock/service/application.go`

### 10.1 这个文件自己的原始 diff

```diff
diff --git a/internal/stock/service/application.go b/internal/stock/service/application.go
index b51f2ad..5feb83a 100644
--- a/internal/stock/service/application.go
+++ b/internal/stock/service/application.go
@@
-	stockRepo := adapters.NewMemoryStockRepository()
+	//stockRepo := adapters.NewMemoryStockRepository()
+	db := persistent.NewMySQL()
+	stockRepo := adapters.NewMySQLStockRepository(db)
```

### 10.2 这段改动的意义

这是最关键的“接线”动作之一。

因为只有这里把 `stockRepo` 从内存实现切成 MySQL 实现，前面写的 `persistent.MySQL` 和 `MySQLStockRepository` 才会真正生效。

也就是说，这一行不是小修，而是：

`stock 服务运行时依赖，正式从内存版切到 MySQL 版。`

## 11. 关键文件八：`internal/stock/app/query/check_if_items_in_stock.go`

### 11.1 这个文件自己的原始 diff

```diff
diff --git a/internal/stock/app/query/check_if_items_in_stock.go b/internal/stock/app/query/check_if_items_in_stock.go
index 00b5c4c..0e86e86 100644
--- a/internal/stock/app/query/check_if_items_in_stock.go
+++ b/internal/stock/app/query/check_if_items_in_stock.go
@@
 func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*entity.Item, error) {
+	if err := h.checkStock(ctx, query.Items); err != nil {
+		return nil, err
+	}
 	var res []*entity.Item
@@
+	// TODO: 扣库存
 	return res, nil
 }
+
+func (h checkIfItemsInStockHandler) checkStock(ctx context.Context, query []*entity.ItemWithQuantity) error {
+	...
+	records, err := h.stockRepo.GetStock(ctx, ids)
+	...
+	if item.Quantity > idQuantityMap[item.ID] {
+		...
+	}
+	...
+	return domain.ExceedStockError{FailedOn: failedOn}
+}
```

### 11.2 这段代码为什么是这组真正的业务核心

因为它让“库存检查”第一次真正有了业务含义。

以前更多只是：

- 看看商品存不存在
- 给你补一个 price ID

现在变成：

- 先看库存够不够
- 不够就直接返回结构化错误
- 够了才继续往下查 Stripe 价格 ID

这就是从“演示型链路”走向“业务约束真正生效”的关键一步。

### 11.3 带中文注释的关键代码

```go
func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*entity.Item, error) {
    // 第一步先查库存是否足够。
    // 只要不够，整条链立刻失败。
    if err := h.checkStock(ctx, query.Items); err != nil {
        return nil, err
    }

    var res []*entity.Item
    for _, i := range query.Items {
        // 库存够了之后，再去 Stripe 拿价格 ID。
        priceID, err := h.stripeAPI.GetPriceByProductID(ctx, i.ID)
        if err != nil || priceID == "" {
            return nil, err
        }
        res = append(res, &entity.Item{
            ID:       i.ID,
            Quantity: i.Quantity,
            PriceID:  priceID,
        })
    }

    // TODO: 扣库存
    // 这里非常关键：当前这节只是“检查库存”，还没有真的扣减库存。
    return res, nil
}

func (h checkIfItemsInStockHandler) checkStock(ctx context.Context, query []*entity.ItemWithQuantity) error {
    var ids []string
    for _, i := range query {
        ids = append(ids, i.ID)
    }

    // 从仓储层批量查库存。
    records, err := h.stockRepo.GetStock(ctx, ids)
    if err != nil {
        return err
    }

    // 先把查询结果整理成 map，后面按商品 ID 直接查更方便。
    idQuantityMap := make(map[string]int32)
    for _, r := range records {
        idQuantityMap[r.ID] += r.Quantity
    }

    var (
        ok       = true
        failedOn []struct {
            ID   string
            Want int32
            Have int32
        }
    )

    for _, item := range query {
        if item.Quantity > idQuantityMap[item.ID] {
            ok = false
            failedOn = append(failedOn, struct {
                ID   string
                Want int32
                Have int32
            }{ID: item.ID, Want: item.Quantity, Have: idQuantityMap[item.ID]})
        }
    }

    if ok {
        return nil
    }

    return domain.ExceedStockError{FailedOn: failedOn}
}
```

### 11.4 这里必须拆开的标准库 / Go 基础

#### `make(map[string]int32)` 是什么

来源：Go 内建函数 `make`

作用：

- 创建一个 map

这里这个 map 的含义是：

- key：商品 ID
- value：当前已知库存数量

为什么这里先转成 map：

- 后面循环校验时，可以 O(1) 按商品 ID 查数量
- 比每次去切片里遍历找匹配项更直接

#### `idQuantityMap[r.ID] += r.Quantity` 为什么能这样写

因为：

- map 里如果 key 不存在，读取零值
- `int32` 的零值是 0

所以第一次遇到某个商品时：

- `idQuantityMap[r.ID]` 默认就是 0
- 加上 `r.Quantity` 后就得到库存数量

这是 Go map 很常见、很顺手的写法。

#### 为什么这里用 `[]struct{...}` 收集失败项

因为作者只想在当前函数内部临时组织一批错误详情，不想再专门为它单独声明一个命名类型。

你可以把它理解成：

- 一个“临时匿名结构体切片”

这种写法在 Go 里不算少见，但对新手可读性一般。更成熟一点的实现，可能会把失败项提成单独类型。

### 11.5 这里还不完美的地方

1. 这里只是“查库存”，还没有真正扣库存。
2. 查库存和扣库存之间还没有事务保护。
3. 并发下可能出现经典问题：两个请求同时都看到库存够，然后都下单成功。
4. 也就是说，这节解决的是“先别让明显超卖发生”，不是完整解决超卖问题。

## 12. `internal/stock/ports/grpc.go` 和 `get_items.go`

### 12.1 这两处在做什么

`GetItems` 这次也统一改成了：

- 应用层返回 `entity.Item`
- gRPC 出口再转成 proto

这其实是在把上一组开始做的边界整理继续补完整，而不是新主题。

### 12.2 为什么值得提一下

因为如果 `GetItems` 这里不一起改，整个 `stock` 服务就会出现：

- 一部分 query 走 entity
- 一部分 query 还直接回 proto

那边界就又会变脏。

## 13. `internal/kitchen/prom.go` 为什么整段被注释掉

这次 diff 里它几乎整段被注释掉了。

这不是这组主线。更像是作者把上一组的 demo 暂时搁置，避免它干扰当前 stock/MySQL 主线。

所以这里我明确告诉你：

- 它不是这组核心内容
- 不要把注意力放错地方

## 14. 这组为什么这么设计

### 14.1 为什么要在 `order/http.go` 和 OpenAPI 两层都校验 `quantity`

因为：

- OpenAPI 契约层负责对外说明规则
- 服务运行时层负责真正守住规则

这是双层防御。

### 14.2 为什么先“查库存”，但还不“扣库存”

因为真正扣库存一旦进来，就会立刻牵涉：

- 并发
- 事务
- 回滚
- 一致性

作者这节先把“校验是否足够”做起来，是一种合理的分阶段推进。

### 14.3 为什么先切 `GetStock(...)`，而 `GetItems(...)` 的 MySQL 版还没做完

因为当前创建订单最关键的业务路径，是“库存够不够”。

作者明显是在沿关键路径推进，而不是一口气把整个仓储所有能力都切完。

## 15. 这组还不完美的地方

1. `MySQLStockRepository.GetItems(...)` 还没实现。
2. `checkStock(...)` 之后还没有真正扣库存。
3. 没有事务和锁，超卖问题并未彻底解决。
4. `order/http.go` 的校验只拦住了 `quantity <= 0`，还没覆盖更多业务规则。
5. MySQL 配置直接写在 `global.yaml`，本地开发方便，但生产一般不会这么直接放明文。

## 16. 这组最该带走的知识点

1. 库存检查真正落地时，仓储接口往往要扩展，不可能一直用最早那一版抽象。
2. OpenAPI 约束和服务内部校验是两层不同职责，最好同时存在。
3. ORM 持久层负责跟数据库打交道，adapter 负责把数据库模型转换成上层实体。
4. `make(map[K]V)` + 累加是一种非常常见的 Go 聚合写法。
5. “检查库存”和“扣减库存”不是一回事，后者要复杂得多。
6. 这组是从“有数据库环境”走向“真正使用数据库做业务判断”的第一步。

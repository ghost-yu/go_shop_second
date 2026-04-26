# `lesson28 -> lesson29` 独立讲义（详细注释版）

这一组的标题虽然是 `convertor`，但它真正想解决的问题不是“怎么写几个转换函数”，而是更底层的一件事：

`不要让外部协议生成出来的类型，直接污染你的业务代码。`

如果你是 Go 小白，这个问题一开始很容易觉得“没必要这么麻烦”。
因为你会直觉地想：
- gRPC 已经帮我生成了 `orderpb.Order`
- OpenAPI 已经帮我生成了 `client.Order`
- 那我直接拿来用不就行了吗？

短期看，当然能用。

但这组课要告诉你的就是：

`能用，不等于应该到处用。`

因为一旦你把这些“协议层类型”直接带进：
- domain
- app command/query
- repository
- consumer
- service

后面就会遇到几个典型问题：

1. 外部协议字段一变，你内部业务全跟着抖
2. 生成代码的字段形态未必适合你内部逻辑
3. 某些类型会带很多你不需要的 tag / optional pointer / protobuf 语义
4. 你的领域层会开始依赖 gRPC / OpenAPI 这种边界技术细节
5. 你以后想同时支持 HTTP 和 gRPC，就会越来越混乱

所以这组真正做的事情是：

`在 order 服务内部建立自己的 entity / domain 类型，并用 convertor 把边界层类型和内部类型隔开。`

这和字幕 `29.txt`、`30.txt` 讲的主线完全一致：
- pb 类型、openapi client 类型不该到处乱用
- 这些东西本质上是 DTO
- 它们应该只活在边界层
- 内部逻辑应该流转你自己的类型

## 1. 这组你应该怎么读

建议按这个顺序看：

1. [api/openapi/order.yml](/g:/shi/go_shop_second/api/openapi/order.yml)
2. [internal/common/client/order/openapi_types.gen.go](/g:/shi/go_shop_second/internal/common/client/order/openapi_types.gen.go)
3. [internal/order/entity/entity.go](/g:/shi/go_shop_second/internal/order/entity/entity.go)
4. [internal/order/domain/order/order.go](/g:/shi/go_shop_second/internal/order/domain/order/order.go)
5. [internal/order/convertor/convertor.go](/g:/shi/go_shop_second/internal/order/convertor/convertor.go)
6. [internal/order/convertor/facade.go](/g:/shi/go_shop_second/internal/order/convertor/facade.go)
7. [internal/order/http.go](/g:/shi/go_shop_second/internal/order/http.go)
8. [internal/order/ports/grpc.go](/g:/shi/go_shop_second/internal/order/ports/grpc.go)
9. [internal/order/app/command/create_order.go](/g:/shi/go_shop_second/internal/order/app/command/create_order.go)

为什么这样看：
- `order.yml` 和生成类型先告诉你“外部世界长什么样”。
- `entity.go` 和 `domain/order.go` 再告诉你“内部世界想长什么样”。
- `convertor.go` 解释这两个世界怎么做转换。
- `http.go` / `grpc.go` 再告诉你“转换应该放在边界层哪里”。
- `create_order.go` 最后用来验证：应用层已经不直接吃 pb 类型了。

## 2. 这组到底在解决什么问题

先说 lesson28 的状态。

那时候代码里有个典型现象：
- HTTP handler 拿到 OpenAPI 生成的 request 类型后，直接往业务层传
- gRPC server 拿到 `orderpb.*` 类型后，直接往业务层传
- domain 里的 `Order.Items` 甚至直接是 `[]*orderpb.Item`
- `Order` 还有一个 `ToProto()`，说明领域对象自己知道 protobuf 长什么样

这意味着什么？

意味着你这个 `domain.Order` 已经不再是一个“纯业务概念”，它身上已经绑了 protobuf 这个技术细节。

这在架构上是很危险的。

因为领域层本来应该表达的是：
- 订单有哪些字段
- 订单状态怎么变
- 订单有哪些业务规则

而不是：
- protobuf 生成的字段名是什么
- OpenAPI 代码生成器怎么生成 struct tag

所以 lesson29 开始干三件事：

1. 补 `entity` 层自己的基础类型
2. 新增 `convertor` 把 client/proto/domain/entity 互转
3. 把 HTTP / gRPC 这些边界层的类型转换责任收回来

这组并不是“业务功能增强”，而是一次典型的“边界清理”。

## 3. 外部协议先收紧：`api/openapi/order.yml`

### 3.1 这个文件自己的原始 diff

```diff
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
```

### 3.2 旧代码做什么，新代码做什么

旧代码：
- OpenAPI schema 里没把这些字段声明成 `required`
- 结果代码生成器会把它们当成“可选字段”
- Go 代码里就会出现一堆 `*string`、`*int32`、`omitempty`

新代码：
- 明确告诉 OpenAPI：这些字段是必填
- 代码生成器就会把它们生成为普通字段，而不是指针字段

### 3.3 为什么这件事很重要

很多新手会觉得：

`*string 也能用啊，无非多解引用一下。`

但实际问题在于：
- 你的业务里大多数时候并不想处理“字段可能不存在”这种三态语义
- 指针字段会让每个地方都开始写 nil 判断
- 生成出来的 DTO 语义会变得很“松”
- 很容易把“空字符串”和“字段没传”混在一起

如果某个字段在协议上本来就是必填，最好的做法就是在 schema 层明确写出来。

### 3.4 当前代码和详细注释

```yaml
components:
  schemas:
    Order:
      type: object
      required:
        - id
        - customerID
        - status
        - items
        - paymentLink
      properties:
        id:
          type: string
        customerID:
          type: string
        status:
          type: string
        items:
          type: array
          items:
            $ref: '#/components/schemas/Item'
        paymentLink:
          type: string
```

这里你要注意的不是 YAML 语法本身，而是 `required` 的语义：

- 它不是给人看的注释
- 它会直接影响代码生成结果
- 它定义了“协议契约”

也就是说：

`先改 schema，后面的生成代码才会变。`

## 4. 生成代码为什么突然从指针变成普通字段：`internal/common/client/order/openapi_types.gen.go`

### 4.1 这个文件自己的原始 diff

```diff
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
```

### 4.2 这段为什么值得讲

虽然这是生成代码，但它恰好能说明一个很重要的工程事实：

`边界契约的改变，会直接改变你所有边界层类型的形态。`

也就是说，上游 schema 一旦不严谨，你下游 Go 类型就会跟着脏。

### 4.3 你作为 Go 小白要记住什么

1. 生成代码通常不要手改
- 因为下次生成就没了

2. 真正应该改的是生成源
- 这里就是 `api/openapi/order.yml`

3. 指针字段不是“更高级”
- 它只是代表“这个字段可以缺席”
- 如果业务不需要这种语义，普通字段更干净

## 5. 先建立内部自己的类型：`internal/order/entity/entity.go`

### 5.1 这个文件自己的原始 diff

```diff
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
```

### 5.2 这段在做什么

这段是在 order 服务内部正式定义自己的“通用业务数据结构”。

注意它不是：
- protobuf 生成类型
- openapi client 类型
- http request 专用类型

而是这个服务内部自己认可的基础结构。

### 5.3 关键代码和详细注释

```go
package entity

// Item 表示订单里真正完整的商品项。
// 这里面已经包含名字、数量、PriceID 等完整信息。
type Item struct {
    ID       string
    Name     string
    Quantity int32
    PriceID  string
}

// ItemWithQuantity 表示“只有商品 ID + 数量”的轻量输入结构。
// 它更适合创建订单这种场景，因为用户下单时通常不会自己把商品全量信息传进来。
type ItemWithQuantity struct {
    ID       string
    Quantity int32
}
```

### 5.4 为什么要拆成 `Item` 和 `ItemWithQuantity`

这是很好的建模意识。

因为这两个东西虽然长得像，但语义不同：
- `ItemWithQuantity`：用户输入的“我买了什么、买多少”
- `Item`：经过库存服务校验和补全后的完整商品信息

如果你只用一个结构体硬扛这两种语义，后面很容易出现：
- 某些字段在某些流程里永远为空
- 一堆“这个阶段先不填”的代码
- 业务含义变混

## 6. domain 不再直接依赖 pb：`internal/order/domain/order/order.go`

### 6.1 这个文件自己的原始 diff

```diff
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
```

### 6.2 旧代码做什么，新代码做什么

旧代码的问题很典型：
- `domain.Order.Items` 直接是 `[]*orderpb.Item`
- `domain.Order` 自己实现了 `ToProto()`

这就等于在说：

`领域对象天然知道 protobuf 长什么样。`

这是不合理的。

领域对象应该知道：
- 自己是什么
- 自己怎么校验
- 自己状态怎么变化

它不应该知道：
- gRPC 要返回的 pb 类型长什么样

新代码把这件事拆开了：
- `Order.Items` 改成内部自己的 `[]*entity.Item`
- `ToProto()` 从领域对象上拿掉
- protobuf 转换责任交给 convertor

### 6.3 关键代码和详细注释

```go
package order

import (
    "fmt"

    "github.com/ghost-yu/go_shop_second/order/entity"
    "github.com/pkg/errors"
    "github.com/stripe/stripe-go/v80"
)

type Order struct {
    ID          string
    CustomerID  string
    Status      string
    PaymentLink string

    // 这里不再是 []*orderpb.Item，而是内部自己的 entity.Item。
    // 这意味着 domain 不再直接耦合 protobuf。
    Items []*entity.Item
}

func NewOrder(id, customerID, status, paymentLink string, items []*entity.Item) (*Order, error) {
    if id == "" {
        return nil, errors.New("empty id")
    }
    if customerID == "" {
        return nil, errors.New("empty customerID")
    }
    if status == "" {
        return nil, errors.New("empty status")
    }
    if items == nil {
        return nil, errors.New("empty items")
    }
    return &Order{
        ID:          id,
        CustomerID:  customerID,
        Status:      status,
        PaymentLink: paymentLink,
        Items:       items,
    }, nil
}
```

### 6.4 为什么删掉 `ToProto()` 很重要

这比看上去重要得多。

因为 `ToProto()` 看起来只是一个方便方法，但它会悄悄把一个错误方向固定下来：

`领域层负责知道外部协议长什么样。`

一旦你后面再加：
- `ToClient()`
- `ToJSON()`
- `ToOpenAPI()`
- `ToKafkaEvent()`

你的领域对象就会越来越胖，越来越像“万能转换器”，最后完全失去边界感。

所以 lesson29 提前刹车，把转换责任移到单独的 convertor 包里，这是对的。

## 7. 这组的主角：`internal/order/convertor/convertor.go`

### 7.1 这个文件自己的原始 diff

```diff
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
+...
```

### 7.2 这段代码为什么是这组课的中心

因为前面说的“边界隔离”，最终都要落到具体代码上。

convertor 的职责就是：
- 在协议类型和内部类型之间做显式转换
- 把耦合集中到一个地方
- 让其他层别再偷偷直接拿 pb/client 类型做业务

### 7.3 关键代码和详细注释

```go
package convertor

import (
    client "github.com/ghost-yu/go_shop_second/common/client/order"
    "github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
    domain "github.com/ghost-yu/go_shop_second/order/domain/order"
    "github.com/ghost-yu/go_shop_second/order/entity"
)

// 这里按类型分了三个 convertor：
// - OrderConvertor
// - ItemConvertor
// - ItemWithQuantityConvertor
// 目的很简单：避免所有转换逻辑糊成一个大文件里的大函数。
type OrderConvertor struct{}
type ItemConvertor struct{}
type ItemWithQuantityConvertor struct{}

func (c *ItemWithQuantityConvertor) EntitiesToProtos(items []*entity.ItemWithQuantity) (res []*orderpb.ItemWithQuantity) {
    for _, i := range items {
        res = append(res, c.EntityToProto(i))
    }
    return
}

func (c *ItemWithQuantityConvertor) EntityToProto(i *entity.ItemWithQuantity) *orderpb.ItemWithQuantity {
    return &orderpb.ItemWithQuantity{
        ID:       i.ID,
        Quantity: i.Quantity,
    }
}

func (c *ItemWithQuantityConvertor) ClientsToEntities(items []client.ItemWithQuantity) (res []*entity.ItemWithQuantity) {
    for _, i := range items {
        res = append(res, c.ClientToEntity(i))
    }
    return
}

func (c *ItemWithQuantityConvertor) ClientToEntity(i client.ItemWithQuantity) *entity.ItemWithQuantity {
    return &entity.ItemWithQuantity{
        ID:       i.Id,
        Quantity: i.Quantity,
    }
}
```

这部分你要学会看出几个模式：

#### 模式 1：单个对象转换 + 列表转换分开

- `EntityToProto` 负责一个对象
- `EntitiesToProtos` 负责切片批量转换

这样做的好处是：
- 单个对象转换逻辑复用
- 列表函数不需要重复写字段映射

#### 模式 2：每种边界都显式命名

比如：
- `ClientToEntity`
- `ProtoToEntity`
- `EntityToClient`
- `EntityToProto`

这比写一个含糊的 `Convert()` 要清晰得多。

因为转换方向在边界代码里非常重要，写不清楚很容易调错。

### 7.4 `OrderConvertor` 为什么不是多余的

再看这块：

```go
func (c *OrderConvertor) EntityToProto(o *domain.Order) *orderpb.Order {
    c.check(o)
    return &orderpb.Order{
        ID:          o.ID,
        CustomerID:  o.CustomerID,
        Status:      o.Status,
        Items:       NewItemConvertor().EntitiesToProtos(o.Items),
        PaymentLink: o.PaymentLink,
    }
}

func (c *OrderConvertor) ProtoToEntity(o *orderpb.Order) *domain.Order {
    c.check(o)
    return &domain.Order{
        ID:          o.ID,
        CustomerID:  o.CustomerID,
        Status:      o.Status,
        PaymentLink: o.PaymentLink,
        Items:       NewItemConvertor().ProtosToEntities(o.Items),
    }
}
```

它的价值在于：
- 当 `Order` 里以后新增字段时，修改点集中
- HTTP / gRPC / consumer 不需要各自重复写一遍组装代码
- 你可以很清楚地看到“边界类型”和“内部类型”之间如何映射

这就是“集中边界映射”的意义。

### 7.5 这里容易踩的坑

#### 坑 1：漏字段

转换代码最常见的问题不是编译不过，而是：

`能编译，但少拷了一个字段。`

比如忘了 `PaymentLink`、忘了 `PriceID`，业务就会变得很诡异。

所以 convertor 代码一定要警惕“静默丢字段”。

#### 坑 2：nil 检查不全

它现在有个：

```go
func (c *OrderConvertor) check(o interface{}) {
    if o == nil {
        panic("connot convert nil order")
    }
}
```

这有两个特点：
- 好处：快速暴露明显错误
- 风险：用 panic 处理这类情况，在业务服务里比较粗暴

更成熟的项目可能会：
- 返回 error
- 或者在上游保证不传 nil

但在课程阶段，这样写是为了简单直接。

#### 坑 3：拼写错误

`connot convert nil order` 明显是拼写错了，应该是 `cannot`。

这种小问题不影响主线，但你要养成注意细节的习惯。

## 8. 为什么这里还搞了个 facade：`internal/order/convertor/facade.go`

### 8.1 这个文件自己的原始 diff

```diff
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
+...
```

### 8.2 这段在做什么

它在做一个很轻量的“单例获取器”：
- `NewOrderConvertor()`
- `NewItemConvertor()`
- `NewItemWithQuantityConvertor()`

内部用 `sync.Once` 确保只初始化一次。

### 8.3 关键代码和详细注释

```go
var (
    orderConvertor *OrderConvertor
    orderOnce      sync.Once
)

func NewOrderConvertor() *OrderConvertor {
    orderOnce.Do(func() {
        orderConvertor = new(OrderConvertor)
    })
    return orderConvertor
}
```

### 8.4 为什么这么做

从“绝对必要性”来看，这里不是必须。

因为这些 convertor 都是无状态对象，你其实完全可以每次直接：

```go
&OrderConvertor{}
```

也能跑。

课程这里这么写，更像是一种统一入口习惯：
- 外部不用关心怎么创建
- 以后如果 convertor 里真的加配置或缓存，入口不需要改 everywhere

但你也要知道，它不是这组的核心。

核心仍然是：

`把转换逻辑集中管理。`

## 9. HTTP 边界开始做自己的转换：`internal/order/http.go`

### 9.1 这个文件自己的原始 diff

```diff
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
```

### 9.2 旧代码的问题

旧代码里 `req.Items` 来自 OpenAPI 生成的 client 类型。
它直接传给 `CreateOrder`，意味着：

`应用层 command 已经开始依赖 HTTP client DTO。`

这就把 HTTP 边界的细节渗透进了应用层。

### 9.3 关键代码和详细注释

```go
func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
    ctx, span := tracing.Start(c, "PostCustomerCustomerIDOrders")
    defer span.End()

    var req client.CreateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    r, err := H.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
        CustomerID: req.CustomerID,

        // 重点在这里：
        // HTTP handler 负责把 OpenAPI client 类型转换成内部 entity 类型。
        // 到了应用层之后，就不应该再看到 client.ItemWithQuantity 了。
        Items: convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
    })
    if err != nil {
        c.JSON(http.StatusOK, gin.H{"error": err})
        return
    }

    c.JSON(http.StatusOK, gin.H{...})
}
```

### 9.4 为什么转换要放在 HTTP handler，而不是放到 command 里

因为 HTTP handler 是边界层。

边界层最该做的事情之一就是：
- 接收外部协议数据
- 验证/绑定
- 转成内部格式
- 再交给应用层

如果把转换塞进 command：
- command 会知道 HTTP client 类型
- 应用层就被 HTTP 耦合了

这正是这组想避免的事情。

## 10. gRPC 边界也要做同样的事：`internal/order/ports/grpc.go`

### 10.1 这个文件自己的原始 diff

```diff
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
@@ -26,7 +27,7 @@ func NewGRPCServer(app app.Application) *GRPCServer {
 func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
 	_, err := G.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
 		CustomerID: request.CustomerID,
-		Items:      request.Items,
+		Items:      convertor.NewItemWithQuantityConvertor().ProtosToEntities(request.Items),
 	})
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
```

### 10.2 这里的主线是什么

和 HTTP 一样：

`gRPC 边界也应该自己负责协议类型和内部类型之间的转换。`

### 10.3 关键代码和详细注释

```go
func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
    _, err := G.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
        CustomerID: request.CustomerID,

        // gRPC request 是 protobuf 类型。
        // 这里在边界层先转换成内部 entity，再交给应用层。
        Items: convertor.NewItemWithQuantityConvertor().ProtosToEntities(request.Items),
    })
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    return &empty.Empty{}, nil
}

func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
    o, err := G.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{...})
    if err != nil {
        return nil, status.Error(codes.NotFound, err.Error())
    }

    // 这里不再让 domain.Order 自己 ToProto()。
    // 转换责任集中交给 convertor。
    return convertor.NewOrderConvertor().EntityToProto(o), nil
}

func (G GRPCServer) UpdateOrder(ctx context.Context, request *orderpb.Order) (_ *emptypb.Empty, err error) {
    order, err := domain.NewOrder(
        request.ID,
        request.CustomerID,
        request.Status,
        request.PaymentLink,

        // request.Items 是 protobuf 类型，先转成内部 entity.Item。
        convertor.NewItemConvertor().ProtosToEntities(request.Items),
    )
    ...
}
```

### 10.4 你要真正理解的一点

很多人学到这里会只记住一句话：

`加个 convertor 就好了。`

这不够。

你要记住真正的规则是：

- HTTP 层处理 OpenAPI/client 类型
- gRPC 层处理 protobuf 类型
- 应用层和领域层处理内部类型

也就是：

`哪一层接收谁，就在哪一层完成转换。`

## 11. 应用层开始真正干净：`internal/order/app/command/create_order.go`

### 11.1 这个文件自己的原始 diff

```diff
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
@@ -18,7 +19,7 @@ import (
 
 type CreateOrder struct {
 	CustomerID string
-	Items      []*orderpb.ItemWithQuantity
+	Items      []*entity.ItemWithQuantity
 }
@@ -100,26 +101,26 @@ func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*Creat
 }
 
-func (c createOrderHandler) validate(ctx context.Context, items []*orderpb.ItemWithQuantity) ([]*orderpb.Item, error) {
+func (c createOrderHandler) validate(ctx context.Context, items []*entity.ItemWithQuantity) ([]*entity.Item, error) {
@@ -100,26 +101,26 @@ func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*Creat
-	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, items)
+	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, convertor.NewItemWithQuantityConvertor().EntitiesToProtos(items))
@@ -100,26 +101,26 @@ func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*Creat
-	return resp.Items, nil
+	return convertor.NewItemConvertor().ProtosToEntities(resp.Items), nil
```

### 11.2 为什么这段是“清理完成”的标志

因为应用层 command 现在终于不直接依赖 `orderpb` 了。

这说明边界已经被收住了：
- command 接收内部 `entity.ItemWithQuantity`
- 真要调用 stock gRPC 时，再临时转成 proto
- 调用结果回来，再转回内部 `entity.Item`

### 11.3 关键代码和详细注释

```go
type CreateOrder struct {
    CustomerID string

    // 这里原来是 []*orderpb.ItemWithQuantity。
    // 现在改成内部自己的 entity 类型，说明 command 不再绑定 protobuf。
    Items []*entity.ItemWithQuantity
}

func (c createOrderHandler) validate(ctx context.Context, items []*entity.ItemWithQuantity) ([]*entity.Item, error) {
    if len(items) == 0 {
        return nil, errors.New("must have at least one item")
    }

    items = packItems(items)

    // 注意这里：
    // 应用层内部用 entity，但一旦要调用 stock gRPC，
    // 就在这一个边界点把 entity 转成 proto。
    resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, convertor.NewItemWithQuantityConvertor().EntitiesToProtos(items))
    if err != nil {
        return nil, err
    }

    // gRPC 回来的是 proto，立刻转回内部 entity。
    return convertor.NewItemConvertor().ProtosToEntities(resp.Items), nil
}
```

### 11.4 这段为什么比以前更好

以前的写法相当于：
- command 内部全是 pb
- stock gRPC 的存在感在应用层里特别强

现在的写法是：
- 应用层主要处理自己的业务类型
- 只有在和外部服务交互时，才显式做一次边界转换

这就叫“把外部协议影响收缩到边界点”。

## 12. 这组最重要的架构结论

### 12.1 `pb` 类型和 `client` 类型本质上都是什么

它们本质上都属于 DTO 一类的东西。

DTO 可以简单理解成：

`为了传输而存在的数据结构。`

它们的主要职责是：
- 对接网络协议
- 对接代码生成器
- 对接边界输入输出

它们的主要职责不是：
- 承载业务核心逻辑
- 作为领域模型在系统里四处流动

### 12.2 为什么“DTO 到处飞”会变成问题

因为 DTO 受外部契约约束太强。

比如 protobuf：
- 字段编号
- 生成风格
- optional 语义
- 代码生成规则

比如 OpenAPI client：
- required / optional
- omitempty
- 指针字段
- JSON tag

这些东西都是边界问题，不是业务问题。

一旦业务层直接使用它们，边界问题就会污染整个系统。

### 12.3 `convertor` 不是多余层，而是“边界隔离层”

很多初学者会觉得：

`这不就是多写了一堆赋值代码吗？`

是的，形式上确实是。

但它换来的是：
- 领域层更纯
- 应用层更稳
- 边界变化影响范围更小
- 多协议接入更容易
- 后面重构更安全

这种代码短期不性感，但长期非常值钱。

## 13. 这组也有一些你要知道的局限

### 13.1 `convertor` 代码很多是机械映射

这类代码的坏处是：
- 容易多
- 容易漏字段
- 改字段时要同步更新多个方向

所以真实项目里通常会：
- 写更严格的测试
- 或者使用一些 mapper 方案
- 或者减少不必要的中间层

但课程这里先手写，是为了把思路讲清楚。

### 13.2 `sync.Once` 单例在这里不是核心收益

它能用，但不要把重点学偏成“这节课主要在讲单例模式”。

不是。

这节课重点还是边界分层。

### 13.3 `domain.NewOrder(...)` 里依然允许 `items == nil` 以外的很多情况

比如：
- `items` 为空切片怎么办
- `PriceID` 空怎么办
- `Quantity <= 0` 怎么办

这些校验后面可能还会继续增强。

所以这组的重点不是“领域校验已完美”，而是“类型边界先分清”。

## 14. 这组你最该记住的话

1. `orderpb.*` 和 OpenAPI client 类型都只是边界 DTO，不应该在业务层里到处流转。
2. HTTP 层收到 client 类型，就在 HTTP 层转成内部类型。
3. gRPC 层收到 pb 类型，就在 gRPC 层转成内部类型。
4. domain 不应该有 `ToProto()` 这种把协议细节带进来的方法。
5. convertor 的价值不在“写赋值语句”，而在“把边界耦合集中起来”。
6. 如果外部协议字段变化，最理想的影响范围应该被限制在边界层和 convertor，而不是整个业务层。

## 15. 你现在应该怎么复习这组

建议你做这个顺序：

1. 先看 [internal/order/domain/order/order.go](/g:/shi/go_shop_second/internal/order/domain/order/order.go)
   - 观察 `Order` 不再依赖 `orderpb.Item`
2. 再看 [internal/order/entity/entity.go](/g:/shi/go_shop_second/internal/order/entity/entity.go)
   - 理解内部到底想流转什么类型
3. 再看 [internal/order/convertor/convertor.go](/g:/shi/go_shop_second/internal/order/convertor/convertor.go)
   - 理解 client/proto/entity/domain 之间的映射方向
4. 再看 [internal/order/http.go](/g:/shi/go_shop_second/internal/order/http.go)
   - 看 HTTP 边界怎么做转换
5. 最后看 [internal/order/ports/grpc.go](/g:/shi/go_shop_second/internal/order/ports/grpc.go)
   - 看 gRPC 边界怎么做转换

如果你下一步继续，我就按同样标准写 `lesson29 -> lesson30`。这两节是连续主题，下一组很可能会继续沿着 convertor / 边界收口 / 类型清理往下推进。
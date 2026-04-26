# `lesson31 -> lesson32` 独立讲义（详细注释版）

这一组代码量不大，但它非常典型，属于那种“看上去只是整理一下文件位置和字段名，实际上是在继续把边界层和应用层分开”的课程。

如果你把上一组 `lesson29 -> lesson30` 的主线记住了，这一组就很好理解：

- 上一组先统一了 HTTP response 外壳
- 同时把 OpenAPI 字段改成 `snake_case`
- 结果生成出来的 client 类型字段名也跟着变了
- 于是这一组要继续做两件事：
  1. 把 response DTO 从 handler 里抽出来
  2. 修正 convertor，适配新的生成字段名

所以这组不是“新业务功能”，而是：

`把上一组接口整理继续做完整。`

字幕 `33.txt` 其实讲得很直白：
- 现在 response 还写在 HTTP handler 里
- 真实项目里这样会膨胀
- 这些东西本质上是 DTO
- 应该放到 `app/dto` 这种位置

这组你要真正学会的是：

`临时写在 handler 里的结构体，一旦成为正式的应用层返回契约，就应该被提炼成独立 DTO。`

## 1. 先说明一个小情况

`lesson30 -> lesson31` 在你这份仓库里没有任何代码变化，所以我直接跳到了下一组真实有 diff 的 `lesson31 -> lesson32`。

## 2. 这组你应该先看什么

建议顺序：

1. [internal/order/http.go](/g:/shi/go_shop_second/internal/order/http.go)
2. [internal/order/app/dto/order.go](/g:/shi/go_shop_second/internal/order/app/dto/order.go)
3. [internal/order/convertor/convertor.go](/g:/shi/go_shop_second/internal/order/convertor/convertor.go)

为什么这样看：
- `http.go` 里你能最直观看到“原来局部匿名结构体”是怎么被提出来的。
- `app/dto/order.go` 能告诉你课程想把 DTO 放到哪里、为什么放这里。
- `convertor.go` 则是在修正上一组 OpenAPI 生成字段改名带来的连锁影响。

## 3. 这组到底在解决什么问题

先回忆 lesson31 的状态。

当时 `PostCustomerCustomerIDOrders` 里已经有统一 response 了，但返回体还是这样写的：

```go
resp struct {
    CustomerID  string `json:"customer_id"`
    OrderID     string `json:"order_id"`
    RedirectURL string `json:"redirect_url"`
}
```

这在课程阶段当然能工作。

但它有几个现实问题：

1. 这个结构已经不再是临时变量，而是接口返回契约的一部分
2. 如果以后别的地方也要返回类似结构，你就要复制一份
3. handler 文件会越来越胖
4. 当 DTO 变复杂时，局部匿名结构会很难维护

所以这组做的第一件事，就是把这个 response 结构正式抽到 `app/dto`。

第二件事则是收尾上一组的 `snake_case` 改动。
因为 OpenAPI 生成代码的字段名从：
- `CustomerID`
- `PriceID`

变成了：
- `CustomerId`
- `PriceId`

convertor 之前还按旧字段名在写，所以要补修。

## 4. 把 response DTO 正式提取出来：`internal/order/app/dto/order.go`

### 4.1 这个文件自己的原始 diff

```diff
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
```

### 4.2 旧代码做什么，新代码做什么

旧代码：
- `CreateOrder` 的返回结构直接写在 `http.go` 的函数局部变量里
- 这个结构没有正式的名字，没有稳定的归属位置

新代码：
- 把它提取成一个具名 DTO：`dto.CreateOrderResponse`
- 放到 `internal/order/app/dto/order.go`

这说明课程开始把它当成“正式应用层数据结构”，而不再是一个临时局部拼装物。

### 4.3 关键代码和详细注释

```go
package dto

// CreateOrderResponse 是“创建订单成功后，对外返回的数据结构”。
// 它不是 domain.Order，也不是 protobuf 生成类型，也不是 OpenAPI client 类型。
// 它是当前应用场景下，专门给“创建订单接口成功响应”准备的 DTO。
type CreateOrderResponse struct {
    OrderID     string `json:"order_id"`
    CustomerID  string `json:"customer_id"`
    RedirectURL string `json:"redirect_url"`
}
```

### 4.4 为什么把 DTO 放在 `app/dto`

这点很值得你建立感觉。

字幕里已经点明了原因：
- HTTP handler 依赖 application
- handler 处理完请求后，要把应用层结果包装成响应
- 这类响应 DTO 更接近“应用层与接口层之间的交换结构”

所以把它放在 `app/dto`，逻辑上是说得通的。

你可以把这层理解成：

`domain 负责业务核心语义，dto 负责某个用例/接口需要交换的数据长相。`

### 4.5 DTO 和 entity / domain 的区别

这是很多初学者容易混的点。

#### `domain.Order`
- 是领域对象
- 关注订单本身的业务属性和规则

#### `entity.Item`
- 是内部通用数据结构
- 更偏业务内部流转

#### `dto.CreateOrderResponse`
- 是某个具体用例对外返回的数据样子
- 只服务于“创建订单成功响应”这个场景

所以 DTO 不是“更高级的实体”，它只是“更面向传输和展示的结构”。

## 5. handler 终于不再塞匿名结构体：`internal/order/http.go`

### 5.1 这个文件自己的原始 diff

```diff
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
```

### 5.2 旧代码做什么，新代码做什么

旧代码：
- handler 里定义匿名 `resp struct`
- 成功后逐字段赋值

新代码：
- handler 里直接使用 `dto.CreateOrderResponse`
- 成功时用结构体字面量一次性赋值

这两个变化都不大，但都很有意义：
- 结构体有了正式名字
- 构造方式更集中、更清晰

### 5.3 关键代码和详细注释

```go
func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
    var (
        req  client.CreateOrderRequest
        resp dto.CreateOrderResponse
        err  error
    )
    defer func() {
        H.Response(c, err, &resp)
    }()

    if err = c.ShouldBindJSON(&req); err != nil {
        return
    }

    r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
        CustomerID: req.CustomerId,
        Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
    })
    if err != nil {
        return
    }

    // 这里改成一次性构造 DTO，而不是先零散赋值。
    // 当字段少的时候两种写法都行；
    // 但结构体字面量更不容易漏字段，也更像是在“构造一个明确结果对象”。
    resp = dto.CreateOrderResponse{
        OrderID:     r.OrderID,
        CustomerID:  req.CustomerId,
        RedirectURL: fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderID),
    }
}
```

### 5.4 为什么这一步值得做

这一步的收益，不在“省了几行代码”，而在于：

1. handler 更轻了
- 结构定义不再塞在函数里

2. DTO 变成正式概念
- 以后别的地方提到“创建订单响应”，就有一个统一类型名

3. 更适合后续增长
- 以后新增字段时，DTO 文件里改
- handler 里只看构造逻辑

### 5.5 为什么用结构体字面量比逐字段赋值更稳一些

比如这两种写法：

```go
resp.CustomerID = req.CustomerId
resp.RedirectURL = ...
resp.OrderID = ...
```

和

```go
resp = dto.CreateOrderResponse{
    OrderID: ...,
    CustomerID: ...,
    RedirectURL: ...,
}
```

第二种的优点是：
- 构造意图更强
- 一眼能看到这个对象最后完整长什么样
- 字段变多时更容易审查
- 漏字段时更容易在 code review 里看出来

## 6. convertor 为什么又要改：`internal/order/convertor/convertor.go`

### 6.1 这个文件自己的原始 diff

```diff
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
```

### 6.2 为什么这段修改不是“小修小补”

表面看只是把：
- `CustomerID` 改成 `CustomerId`
- `PriceID` 改成 `PriceId`

但它本质上是在补上一组 `snake_case` 改完之后的“生成代码字段名适配”。

也就是说：
- OpenAPI 改成了 `customer_id`
- 生成器重新生成 Go 代码
- 生成器字段名变成 `CustomerId`
- convertor 还在按旧字段 `CustomerID` 访问
- 所以要修正

### 6.3 关键代码和详细注释

```go
func (c *OrderConvertor) ClientToEntity(o *client.Order) *domain.Order {
    c.check(o)
    return &domain.Order{
        ID:          o.Id,

        // 这里必须跟着生成代码字段名走。
        // 现在 client.Order 的字段名已经变成 CustomerId，不再是 CustomerID。
        CustomerID:  o.CustomerId,
        Status:      o.Status,
        PaymentLink: o.PaymentLink,
        Items:       NewItemConvertor().ClientsToEntities(o.Items),
    }
}

func (c *OrderConvertor) EntityToClient(o *domain.Order) *client.Order {
    c.check(o)
    return &client.Order{
        Id:         o.ID,

        // 反向转换时也是一样，要写回新的字段名。
        CustomerId: o.CustomerID,
        Status:     o.Status,
        PaymentLink: o.PaymentLink,
        Items:      NewItemConvertor().EntitiesToClients(o.Items),
    }
}
```

以及：

```go
func (c *ItemConvertor) ClientToEntity(i client.Item) *entity.Item {
    return &entity.Item{
        ID:       i.Id,
        Name:     i.Name,
        Quantity: i.Quantity,
        PriceID:  i.PriceId,
    }
}

func (c *ItemConvertor) EntityToClient(i *entity.Item) client.Item {
    return client.Item{
        Id:       i.ID,
        Name:     i.Name,
        Quantity: i.Quantity,
        PriceId:  i.PriceID,
    }
}
```

### 6.4 这里你应该真正学到什么

你要学到的不是“大小写改一下”。

而是：

`边界层生成代码一旦变动，convertor 就是你集中收口和消化变化的地方。`

这正好证明了前面几组课的设计是有收益的。

如果没有 convertor，上一组字段改名之后，你就会在很多地方改：
- handler
- service
- command
- query
- domain
- repository
- consumer

现在因为有 convertor，大部分变化只需要集中修补边界映射。

这就是“隔离层”的价值。

### 6.5 为什么生成器会变成 `CustomerId` 而不是 `CustomerID`

这是代码生成器的命名规则问题。

对于人手写 Go 代码，我们通常更偏好：
- `ID`
- `URL`
- `HTTP`

但很多生成器是按自己的大小写规则来的，它不一定完全遵守 Go 社区的缩写习惯。

所以你要接受一个现实：

`自动生成代码不一定优雅，但你要学会在边界层消化它。`

这也是为什么我们不希望业务层到处直接用这些生成类型。

## 7. 这组为什么继续证明 DTO / convertor 思路是对的

这组虽然改动很小，但它刚好起到一个“验证”作用。

### 7.1 DTO 抽出来后，handler 更清爽

`http.go` 里少了局部结构定义，职责更聚焦：
- 收请求
- 调 command
- 组 DTO
- 统一 Response

### 7.2 convertor 把生成代码变化挡在边界层

OpenAPI 改一轮字段名之后，真正需要修的主要是 convertor，而不是整个服务满地找字段名。

### 7.3 app 层开始长出更明确的“交换结构”

这意味着代码正在逐步形成更稳定的层次感：
- domain: 核心业务对象
- entity: 内部通用数据结构
- dto: 应用场景的对外/对上返回结构
- convertor: 边界映射

## 8. 这组的小坑和边界

### 8.1 现在只抽了 `CreateOrderResponse`

说明 DTO 体系还只是开始，不是已经特别完整。

比如 `GetOrder` 这边现在还是：

```go
resp interface{}
resp = convertor.NewOrderConvertor().EntityToClient(o)
```

这表示课程当前只先把最明显的那个 response 抽出来。
以后可能还会继续整理更多返回 DTO。

### 8.2 DTO 放在 `app/dto` 并不是宇宙唯一答案

这是当前课程选择的一种结构。

真实项目里也有人会放：
- `internal/order/dto`
- `internal/order/app/dto`
- `internal/order/transport/http/dto`

关键不是目录名，而是你的分层逻辑要自洽：
- DTO 到底服务谁
- 转换责任放哪
- 谁依赖谁

### 8.3 `convertor` 依然容易漏字段

这组虽然只是修 `CustomerId / PriceId`，但它再次提醒你：

`convertor 是最容易出现“编译通过但字段映射错了”的地方。`

这类代码最好后面补测试。

## 9. 这组最该记住的话

1. 局部匿名 response 结构一旦稳定下来，就应该提炼成正式 DTO。
2. DTO 不等于 domain，也不等于 entity，它是某个应用场景交换数据的形状。
3. OpenAPI 生成字段名变化后，convertor 是最该优先修的地方。
4. 生成代码不一定符合你理想的 Go 命名风格，这很正常。
5. 真正值钱的是：这些不优雅被挡在边界层，而不是渗透进业务层。

## 10. 你现在应该怎么复习这组

建议顺序：

1. [internal/order/http.go](/g:/shi/go_shop_second/internal/order/http.go)
   - 看匿名 `resp struct` 是怎么被替换成 `dto.CreateOrderResponse` 的
2. [internal/order/app/dto/order.go](/g:/shi/go_shop_second/internal/order/app/dto/order.go)
   - 理解为什么这个类型开始值得有一个正式名字
3. [internal/order/convertor/convertor.go](/g:/shi/go_shop_second/internal/order/convertor/convertor.go)
   - 看上一组 snake_case 改动如何传导到 convertor 修正

如果你继续，我下一组就写 `lesson32 -> lesson33`。字幕映射显示那一组会开始进入 `DLQ / DLX`，也就是 RabbitMQ 死信队列相关内容。
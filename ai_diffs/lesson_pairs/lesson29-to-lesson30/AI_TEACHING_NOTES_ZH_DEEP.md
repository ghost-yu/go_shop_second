# `lesson29 -> lesson30` 独立讲义（详细注释版）

这一组的提交名叫 `response`，但如果你只把它理解成“多了个统一返回结构”，那就看浅了。

这组其实同时在做两件很重要的接口层整理：

1. 给 HTTP 返回值加统一外壳
2. 把 OpenAPI 字段命名统一收口成 `snake_case`

这两件事看起来像“前端接口格式整理”，不像前面那种有明显业务推进的课程，
但它们在真实工程里非常重要。

因为当一个服务开始被：
- 前端调用
- 自动生成客户端调用
- 文档系统消费
- 测试工具验证

接口的“长相”就不再是随便写写了。

你必须开始考虑这些问题：
- 成功返回和失败返回是否统一
- 每个接口的外层格式是不是一致
- trace_id 能不能固定拿到
- JSON 字段命名是不是统一
- OpenAPI 文档、服务端实现、客户端生成代码能不能对齐

这组课就是在做这种“接口契约收口”。

字幕 `31.txt`、`32.txt` 的重点也正好对应这两条：
- `31.txt`：统一 response 返回值
- `32.txt`：把命名统一成 `snake_case`

## 1. 这组你应该怎么读

建议顺序：

1. [internal/common/response.go](/g:/shi/go_shop_second/internal/common/response.go)
2. [internal/order/http.go](/g:/shi/go_shop_second/internal/order/http.go)
3. [api/openapi/order.yml](/g:/shi/go_shop_second/api/openapi/order.yml)
4. [internal/common/client/order/openapi_types.gen.go](/g:/shi/go_shop_second/internal/common/client/order/openapi_types.gen.go)
5. [internal/order/ports/openapi_api.gen.go](/g:/shi/go_shop_second/internal/order/ports/openapi_api.gen.go)
6. [internal/order/ports/openapi_types.gen.go](/g:/shi/go_shop_second/internal/order/ports/openapi_types.gen.go)

为什么这样看：
- `response.go` 是这组最核心的“公共抽象”。
- `http.go` 能让你看到这个抽象怎么落到具体 handler 上。
- `order.yml` 解释为什么路径参数、字段名都开始改成下划线风格。
- 生成代码最后再看，因为它们本质上是“契约改动的结果”，不是这组的源头。

## 2. 这组到底解决什么问题

先看 lesson29 之前的样子。

那时 `order/http.go` 的返回方式比较分散：
- 有时直接 `c.JSON(http.StatusBadRequest, gin.H{"error": ...})`
- 有时成功返回 `message + trace_id + customer_id + order_id + redirect_url`
- 有时查订单时返回的是 `message + trace_id + data: { Order: o }`
- 失败和成功返回的结构不统一
- 有些字段还是 `camelCase` / `mixedCase`

这会带来几个实际问题：

1. 前端或调用方不好统一解析
- 有的接口成功时 `data` 是一层对象
- 有的接口成功时字段直接散在外层
- 有的失败时只有 `error`

2. trace 信息拿起来不一致
- 有的有 `trace_id`
- 有的错误时没有

3. OpenAPI 和真实返回结构容易不对齐
- 文档一套
- 实际 JSON 一套

4. 命名风格不统一
- `customerID`
- `orderID`
- `paymentLink`
- `priceID`

这些字段如果暴露给前端或第三方调用方，会显得很乱。

所以 lesson30 做的事情其实可以总结成一句话：

`把 order 服务对外 HTTP 契约整理得更像一个真正可以长期维护的接口。`

## 3. 统一响应外壳：`internal/common/response.go`

### 3.1 这个文件自己的原始 diff

```diff
diff --git a/internal/common/response.go b/internal/common/response.go
new file mode 100644
index 0000000..bd57fef
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
```

### 3.2 旧代码做什么，新代码做什么

旧代码：
- 每个 handler 自己决定 JSON 怎么长
- 每个 handler 自己决定错误怎么吐
- `trace_id` 也是各自拼接

新代码：
- 提供一个统一的 `BaseResponse`
- handler 只需要把 `err` 和 `data` 传进去
- 成功时自动吐：`errno=0, message=success, data=..., trace_id=...`
- 失败时自动吐：`errno=2, message=err.Error(), data=nil, trace_id=...`

一句话：

`把“怎么返回 JSON”从具体 handler 里抽出来，变成一个公共响应约定。`

### 3.3 关键代码和详细注释

```go
package common

import (
    "net/http"

    "github.com/ghost-yu/go_shop_second/common/tracing"
    "github.com/gin-gonic/gin"
)

// BaseResponse 本身不存状态，它更像一个“能力容器”。
// 把它嵌入到 HTTPServer 里后，HTTPServer 就能直接复用统一响应方法。
type BaseResponse struct{}

// response 这里故意不导出（小写）。
// 含义是：这是服务内部封装用的统一返回外壳，不希望别人从包外直接依赖它的具体实现。
type response struct {
    Errno   int    `json:"errno"`
    Message string `json:"message"`
    Data    any    `json:"data"`
    TraceID string `json:"trace_id"`
}

func (base *BaseResponse) Response(c *gin.Context, err error, data interface{}) {
    // 这是统一出口。
    // handler 不再自己拼 success/error JSON，而是把决策交给这里。
    if err != nil {
        base.error(c, err)
    } else {
        base.success(c, data)
    }
}

func (base *BaseResponse) success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, response{
        Errno:   0,
        Message: "success",
        Data:    data,

        // trace_id 统一从 request context 里拿。
        // 这样无论哪个 handler，只要走这套响应，就都能带上 trace_id。
        TraceID: tracing.TraceID(c.Request.Context()),
    })
}

func (base *BaseResponse) error(c *gin.Context, err error) {
    c.JSON(http.StatusOK, response{
        Errno:   2,
        Message: err.Error(),
        Data:    nil,
        TraceID: tracing.TraceID(c.Request.Context()),
    })
}
```

### 3.4 为什么这里要用“组合”，不是到处复制粘贴

字幕里提到了一个点：Go 没有传统面向对象那种继承，但可以通过嵌入结构体实现类似复用。

这里就是：
- 定义 `BaseResponse`
- 在 `HTTPServer` 里嵌入它
- `HTTPServer` 自动拥有 `Response(...)` 方法

这在 Go 里是非常常见的复用方式。

你不用把它神化成“高级继承技巧”，它本质上就是：

`把公共行为组合进来，而不是每个 handler 手写一遍。`

### 3.5 为什么错误和成功都返回 `200 OK`

这里你要特别注意：

这不是“HTTP 语义最标准”的做法，但它是很多内部系统会采用的一种约定。

也就是说，这套系统选择了：
- HTTP 层统一返回 200
- 业务成功/失败由 `errno` 判定

这有好处：
- 前端/调用方可以统一按一个 JSON 结构解
- 某些网关或旧系统更容易兼容

但也有明显代价：
- HTTP 状态码失去语义信息
- 监控和网关只看状态码时会误判
- REST 风格上不够标准

所以你要理解它是“项目约定”，不是绝对最佳实践。

## 4. handler 怎么开始收口：`internal/order/http.go`

### 4.1 这个文件自己的原始 diff

```diff
diff --git a/internal/order/http.go b/internal/order/http.go
index eebdb99..4df5213 100644
--- a/internal/order/http.go
+++ b/internal/order/http.go
@@ -1,18 +1,17 @@
 package main
 
 import (
 	"fmt"
-	"net/http"
 
+	"github.com/ghost-yu/go_shop_second/common"
 	client "github.com/ghost-yu/go_shop_second/common/client/order"
-	"github.com/ghost-yu/go_shop_second/common/tracing"
 	"github.com/ghost-yu/go_shop_second/order/app"
 	"github.com/ghost-yu/go_shop_second/order/app/command"
 	"github.com/ghost-yu/go_shop_second/order/app/query"
 	"github.com/ghost-yu/go_shop_second/order/convertor"
 	"github.com/gin-gonic/gin"
 )
 
 type HTTPServer struct {
+	common.BaseResponse
 	app app.Application
 }
@@ -20,42 +19,39 @@ type HTTPServer struct {
-func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
-	ctx, span := tracing.Start(c, "PostCustomerCustomerIDOrders")
-	defer span.End()
-
-	var req client.CreateOrderRequest
-	if err := c.ShouldBindJSON(&req); err != nil {
-		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
-		return
-	}
-	r, err := H.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
-		CustomerID: req.CustomerID,
-		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
-	})
-	if err != nil {
-		c.JSON(http.StatusOK, gin.H{"error": err})
-		return
-	}
-	c.JSON(http.StatusOK, gin.H{
-		"message":      "success",
-		"trace_id":     tracing.TraceID(ctx),
-		"customer_id":  req.CustomerID,
-		"order_id":     r.OrderID,
-		"redirect_url": fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderID),
-	})
+func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
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
+
+	if err = c.ShouldBindJSON(&req); err != nil {
+		return
+	}
+	r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
+		CustomerID: req.CustomerId,
+		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
+	})
+	if err != nil {
+		return
+	}
+	resp.CustomerID = req.CustomerId
+	resp.RedirectURL = fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderID)
+	resp.OrderID = r.OrderID
 }
```

### 4.2 旧代码做什么，新代码做什么

旧代码：
- handler 自己开 tracing span
- handler 自己处理 bind error
- handler 自己写 success JSON
- handler 自己写 error JSON

新代码：
- `HTTPServer` 嵌入 `BaseResponse`
- handler 用 `req/err/resp` 三段式组织逻辑
- 用 `defer H.Response(c, err, &resp)` 统一出口
- 请求绑定失败直接给 `err`
- 成功时只负责把 `resp` 填完整

### 4.3 关键代码和详细注释

```go
type HTTPServer struct {
    common.BaseResponse // 嵌入后，HTTPServer 直接拥有统一响应能力
    app app.Application
}

func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
    var (
        req client.CreateOrderRequest
        err error

        // 这里把成功返回的数据定义成一个局部匿名结构体。
        // 它只服务于这个接口，不必为了一个很小的返回结构单独建全局类型。
        resp struct {
            CustomerID  string `json:"customer_id"`
            OrderID     string `json:"order_id"`
            RedirectURL string `json:"redirect_url"`
        }
    )

    // 非常关键：统一 defer 出口。
    // 这样函数内部只需要更新 err 和 resp，最后统一交给 BaseResponse 决定怎么吐 JSON。
    defer func() {
        H.Response(c, err, &resp)
    }()

    if err = c.ShouldBindJSON(&req); err != nil {
        // 这里只赋 err，不自己 c.JSON。
        // defer 会在函数退出时统一处理。
        return
    }

    r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
        CustomerID: req.CustomerId,
        Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
    })
    if err != nil {
        return
    }

    // 成功时就只是填 resp。
    resp.CustomerID = req.CustomerId
    resp.RedirectURL = fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderID)
    resp.OrderID = r.OrderID
}
```

### 4.4 这里为什么用 `defer` 统一返回

这个技巧非常值得你学。

它带来的好处是：
- 减少重复 `c.JSON(...)`
- 减少成功/失败分支里到处手写返回格式
- 让函数主体更聚焦业务流程

你可以把它理解成：

`先把“怎么返回”统一放到 defer，函数里面只管“做什么”和“是否报错”。`

但它也有一个前提：
- 你得能接受所有返回都走同一个出口
- 你要小心 `err` 和 `resp` 的赋值时机

### 4.5 这里有个你必须知道的 Go 细节

`defer` 捕获变量时，捕获的是变量本身，不是当前值的拷贝（对这里的闭包写法来说）。

所以这段代码之所以成立，是因为：
- `err` 后续可以不断被更新
- `resp` 后续可以不断被填值
- defer 执行时看到的是“最终状态”

如果你对 Go 闭包和 defer 不熟，这里很容易误解。

### 4.6 为什么这里把 `tracing.Start(...)` 去掉了

旧版本 handler 手动开 span，新版本直接用 `c.Request.Context()`。

这说明课程此时更偏向：
- 把 trace 统一从 request context 取
- 不在 handler 里额外维护一套独立返回逻辑

这不代表 tracing 不重要，而是这组重点转移到了“统一 response”。

## 5. 查询接口也被同样收口：`internal/order/http.go` 的 `Get...`

### 5.1 关键变化

旧代码：
- `Get` 成功时直接返回 `gin.H{"message":..., "trace_id":..., "data": gin.H{"Order": o}}`
- 结构比较随意，甚至 `Order` 这个 key 还是大写开头

新代码：
- 统一走 `resp interface{}` + `H.Response(c, err, resp)`
- 把实体转成 client 类型后直接放进 `data`

### 5.2 关键代码和详细注释

```go
func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
    var (
        err  error
        resp interface{}
    )
    defer func() {
        H.Response(c, err, resp)
    }()

    o, err := H.app.Queries.GetCustomerOrder.Handle(c.Request.Context(), query.GetCustomerOrder{
        OrderID:    orderID,
        CustomerID: customerID,
    })
    if err != nil {
        return
    }

    // 注意这里直接转成 client.Order，
    // 也就是 HTTP 对外返回的是“面向 API 的类型”，不是 domain 原始对象。
    resp = convertor.NewOrderConvertor().EntityToClient(o)
}
```

### 5.3 为什么不再返回 `data: { "Order": o }`

因为那种写法很随意：
- key 叫 `Order`，还是大写
- 外层又套一层 `data`
- 调用方得写 `data.Order.xxx`

新代码把 `data` 直接设成 `client.Order`，结构更规整：
- 调用方解析更统一
- 文档对齐更容易
- 少一层没必要的包装

## 6. OpenAPI 为什么开始全量改成 `snake_case`：`api/openapi/order.yml`

### 6.1 这个文件自己的原始 diff

```diff
diff --git a/api/openapi/order.yml b/api/openapi/order.yml
index c65f1cc..da9fbf8 100644
--- a/api/openapi/order.yml
+++ b/api/openapi/order.yml
@@ -10,18 +10,18 @@ servers:
 paths:
-	/customer/{customerID}/orders/{orderID}:
+	/customer/{customer_id}/orders/{order_id}:
 ...
-		name: customerID
+		name: customer_id
 ...
-		name: orderID
+		name: order_id
 ...
-	/customer/{customerID}/orders:
+	/customer/{customer_id}/orders:
 ...
-		name: customerID
+		name: customer_id
 ...
-		- customerID
+		- customer_id
 ...
-		- paymentLink
+		- payment_link
 ...
-		customerID:
+		customer_id:
 ...
-		paymentLink:
+		payment_link:
 ...
-		- priceID
+		- price_id
 ...
-		priceID:
+		price_id:
```

### 6.2 这段到底在改什么

它不是简单改几个字段名，而是在统一 API 契约风格：
- path 参数改成 `customer_id` / `order_id`
- JSON 字段改成 `customer_id` / `payment_link` / `price_id`
- request body 里的字段名也同步改

### 6.3 为什么要改成 `snake_case`

Go 代码里的字段名一般是驼峰，因为 Go 导出字段必须大写开头。

但 JSON / HTTP API 世界里，很多团队更偏向 `snake_case`，原因包括：
- 前后端语言混用时可读性更稳定
- 数据库字段、日志字段常常也是下划线风格
- OpenAPI 文档里看起来更统一

更重要的是：

`接口风格一旦选定，就应该尽量一致，不要同一个 API 里一半 camelCase 一半 snake_case。`

### 6.4 你要明白的一个事实

这里改的是协议，不是 Go 内部字段。

也就是说：
- Go 内部继续可以是 `CustomerID`
- 但对外 JSON tag、OpenAPI schema 可以是 `customer_id`

这正是边界层和内部层分开的价值。

## 7. 生成代码为什么跟着大面积变化

### 7.1 `internal/common/client/order/openapi_types.gen.go`

当前 lesson30 的生成结果是：

```go
type CreateOrderRequest struct {
    CustomerId string             `json:"customer_id"`
    Items      []ItemWithQuantity `json:"items"`
}

type Order struct {
    CustomerId  string `json:"customer_id"`
    Id          string `json:"id"`
    Items       []Item `json:"items"`
    PaymentLink string `json:"payment_link"`
    Status      string `json:"status"`
}
```

你要注意两个层面：

1. JSON tag 已经变成 snake_case
2. Go 字段名被生成器处理成了 `CustomerId`、`PriceId` 这种形式

### 7.2 为什么会出现 `CustomerId` 而不是 `CustomerID`

这是代码生成器的命名规则问题。

很多生成器对 `id` 这种缩写的处理未必完全符合 Go 社区手写代码的习惯。

Go 社区里通常更推荐：
- `ID`
- `URL`
- `HTTP`

但生成器未必都这么智能。

所以你会看到：
- `CustomerId`
- `PriceId`

这看上去不够优雅，但这不是主线问题。

真正主线仍然是：

`边界类型变了，内部逻辑不需要跟着大改，因为前一组已经有 convertor 做隔离。`

这就是 lesson29 铺垫出来的收益。

## 8. 路由和接口名为什么也变了：`internal/order/ports/openapi_api.gen.go`

### 8.1 当前生成结果里的关键点

```go
// (POST /customer/{customer_id}/orders)
PostCustomerCustomerIdOrders(c *gin.Context, customerId string)

// (GET /customer/{customer_id}/orders/{order_id})
GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerId string, orderId string)
```

以及：

```go
router.POST(options.BaseURL+"/customer/:customer_id/orders", wrapper.PostCustomerCustomerIdOrders)
router.GET(options.BaseURL+"/customer/:customer_id/orders/:order_id", wrapper.GetCustomerCustomerIdOrdersOrderId)
```

### 8.2 这说明什么

说明 OpenAPI 里的命名变化会一路传导到：
- 生成的 client 方法名
- 生成的 server interface 方法名
- router 参数名绑定
- request body 结构字段

所以你不能把 schema 改动看成“只是文档层面的事”。

它在这种代码生成体系里，本质上就是“接口单一事实来源”。

## 9. 这一组为什么值得你认真学

因为它其实在训练你一个重要工程习惯：

`把“接口怎么长”当成正式设计对象，而不是临时输出。`

很多初学者写接口时容易这样：
- 先能跑
- 成功了就随手返回一个 `gin.H`
- 出错了再临时返回另一个 `gin.H`
- 字段名想到什么写什么

短期当然省事。

但一旦系统开始有：
- 多个接口
- 多个调用方
- 自动文档
- 自动客户端
- 接口测试

这种随意性会迅速变成维护负担。

lesson30 做的就是把这种随意性收一收。

## 10. 这组也有你必须知道的局限

### 10.1 `errno` 现在还是硬编码

```go
Errno: 0
Errno: 2
```

这说明错误体系还没完全体系化。

更成熟的做法可能会有：
- 业务错误码枚举
- 不同错误映射不同 errno
- 错误分类和国际化消息

当前版本只是先把“统一外壳”搭起来。

### 10.2 所有错误都直接 `err.Error()` 返回给外部

这在课程阶段方便调试，但真实生产里要慎重：
- 可能泄露内部细节
- 错误文案不稳定
- 不利于前端按错误码处理

### 10.3 所有 HTTP 都返回 200 并不总是理想

前面说过，这是一种项目约定，不是 universally best practice。

### 10.4 `response` 是匿名结构 + 局部 DTO，适合当前体量

当前 `Post` 里直接定义局部匿名 `resp struct` 是合理的。
但如果字段越来越多，或者多个接口共享返回体，后面可能会继续抽成独立 DTO。

## 11. 这组最该记住的结论

1. `BaseResponse` 的核心价值是统一成功/失败响应格式和 trace_id 注入。
2. `defer H.Response(c, err, data)` 是这组最值得学的 Go 写法之一。
3. OpenAPI schema 一旦改名，生成的 client/server/types 会连锁变化。
4. `snake_case` 改动不是审美问题，而是 API 契约一致性问题。
5. 前一组做了 convertor，这一组才能放心改边界字段名而不把内部逻辑一起拖乱。
6. 统一 response 外壳是接口工程化的重要一步，不是“只是前端喜欢这样”。

## 12. 你现在应该怎么复习这组

建议按这个顺序回看：

1. [internal/common/response.go](/g:/shi/go_shop_second/internal/common/response.go)
   - 先理解统一外壳长什么样
2. [internal/order/http.go](/g:/shi/go_shop_second/internal/order/http.go)
   - 理解 defer 统一出口怎么落地
3. [api/openapi/order.yml](/g:/shi/go_shop_second/api/openapi/order.yml)
   - 看命名如何统一成 snake_case
4. [internal/common/client/order/openapi_types.gen.go](/g:/shi/go_shop_second/internal/common/client/order/openapi_types.gen.go)
   - 看 schema 改动如何传导到生成代码
5. [internal/order/ports/openapi_api.gen.go](/g:/shi/go_shop_second/internal/order/ports/openapi_api.gen.go)
   - 看 path 参数名如何影响服务端接口和路由绑定

如果你继续，我下一组就写 `lesson30 -> lesson31`。
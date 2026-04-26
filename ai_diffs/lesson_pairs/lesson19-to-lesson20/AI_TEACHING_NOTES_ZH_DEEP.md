# `lesson19 -> lesson20` 独立讲义（重写版）

这一组的重点不是“前端写得漂亮不漂亮”，而是：

`上一组已经能在后端生成 Stripe 支付链接了，这一组开始把这条链做成一个人能实际操作、能肉眼观察状态变化的页面。`

所以你要把它理解成一组“流程可视化”差异，而不是“前端课程插播”。

它解决的核心问题是：

1. 创建订单后，客户端下一步该去哪一个页面继续操作
2. 这个页面怎么知道自己对应哪张订单
3. 页面怎么不断查询后端，知道订单现在是“等待支付”、“已支付”还是“已完成”

## 1. 原始 diff 去哪里看

对照文件还在这里：
[diff.md](/g:/shi/go_shop_second/ai_diffs/lesson_pairs/lesson19-to-lesson20/diff.md)

这一版讲义按你的要求改了：
- 不在开头贴整份 diff
- 每讲一个文件，先贴这个文件自己的 diff
- 然后再贴“带中文注释的代码”
- 再解释为什么这么做

## 2. 正确阅读顺序

这组要按这个顺序读：

1. [internal/order/http.go](/g:/shi/go_shop_second/internal/order/http.go)
2. [internal/order/main.go](/g:/shi/go_shop_second/internal/order/main.go)
3. [public/success.html](/g:/shi/go_shop_second/public/success.html)

原因很简单：
- `http.go` 决定后端到底返回什么给前端
- `main.go` 决定浏览器能不能访问这个页面
- `success.html` 决定页面如何把后端状态展示出来

## 3. 总调用链

这组最重要的调用链是：

```text
POST /api/customer/{customerID}/orders
-> order/http.go 创建订单
-> 返回 order_id + redirect_url
-> 浏览器打开 /success?customerID=...&orderID=...
-> success.html 从 URL 里取出 customerID/orderID
-> success.html 轮询 GET /api/customer/{customerID}/orders/{orderID}
-> 如果状态是 waiting_for_payment，就显示 payment link
-> 用户点链接跳去 Stripe
-> 页面继续轮询，等待后续状态变化
```

这条链里你一定要记住一句话：

`success.html 只是一个观察者页面，它自己不会修改订单状态。`

也就是说，页面只是去查状态，然后显示。真正改状态的动作，还是在后端别的链路里做的。

## 4. 关键文件详细讲解

### 4.1 [internal/order/http.go](/g:/shi/go_shop_second/internal/order/http.go)

这份文件是这组最重要的后端入口，因为前端页面后面能不能跑，全靠它现在返回的字段。

先看这个文件自己的原始 diff。

```diff
diff --git a/internal/order/http.go b/internal/order/http.go
index b40adc7..2073a4f 100644
--- a/internal/order/http.go
+++ b/internal/order/http.go
@@ -1,6 +1,7 @@
 package main
 
 import (
+	"fmt"
 	"net/http"
@@ -29,9 +30,10 @@ func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID stri
 		return
 	}
 	c.JSON(http.StatusOK, gin.H{
-		"message":     "success",
-		"customer_id": req.CustomerID,
-		"order_id":    r.OrderID,
+		"message":      "success",
+		"customer_id":  req.CustomerID,
+		"order_id":     r.OrderID,
+		"redirect_url": fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderID),
 	})
 }
@@ -44,5 +46,10 @@ func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerI
 		c.JSON(http.StatusOK, gin.H{"error": err})
 		return
 	}
-	c.JSON(http.StatusOK, gin.H{"message": "success", "data": o})
+	c.JSON(http.StatusOK, gin.H{
+		"message": "success",
+		"data": gin.H{
+			"Order": o,
+		},
+	})
 }
```

这个 diff 里其实只有两件真正重要的事。

#### 第一件事：创建订单成功后，多返回了一个 `redirect_url`

下面我把当前文件里和这次差异最相关的代码贴出来，并直接写中文注释。

```go
package main

import (
    "fmt" // 这次新加 fmt，就是为了拼接 redirect_url

    "github.com/ghost-yu/go_shop_second/common"
    client "github.com/ghost-yu/go_shop_second/common/client/order"
    "github.com/ghost-yu/go_shop_second/common/consts"
    "github.com/ghost-yu/go_shop_second/common/convertor"
    "github.com/ghost-yu/go_shop_second/common/handler/errors"
    "github.com/ghost-yu/go_shop_second/order/app"
    "github.com/ghost-yu/go_shop_second/order/app/command"
    "github.com/ghost-yu/go_shop_second/order/app/dto"
    "github.com/ghost-yu/go_shop_second/order/app/query"
    "github.com/gin-gonic/gin"
)

type HTTPServer struct {
    common.BaseResponse // 这一层不是直接 c.JSON 了，而是统一走 BaseResponse 封装
    app app.Application
}

func (H HTTPServer) PostCustomerCustomerIdOrders(c *gin.Context, customerID string) {
    var (
        req  client.CreateOrderRequest // 这里已经不是以前直接绑 protobuf 结构了，而是绑 client 侧 DTO
        resp dto.CreateOrderResponse   // 返回值也不再是临时 gin.H，而是 DTO
        err  error
    )

    // 这里 defer 的意思是：
    // 函数最后统一把 resp / err 交给 Response 去输出成 HTTP 响应。
    defer func() {
        H.Response(c, err, &resp)
    }()

    // 先从请求体里解析 JSON。
    if err = c.ShouldBindJSON(&req); err != nil {
        err = errors.NewWithError(consts.ErrnoBindRequestError, err)
        return
    }

    // 再做一层基础校验，比如 quantity 不能 <= 0。
    if err = H.validate(req); err != nil {
        err = errors.NewWithError(consts.ErrnoRequestValidateError, err)
        return
    }

    // 真正创建订单还是走应用层 command。
    r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
        CustomerID: req.CustomerId,
        Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
    })
    if err != nil {
        return
    }

    // 这组最关键的新字段就在这里：RedirectURL。
    // 后端不只返回 OrderID，还告诉前端：
    // 你下一步应该去哪个页面继续这个支付流程。
    resp = dto.CreateOrderResponse{
        OrderID:     r.OrderID,
        CustomerID:  req.CustomerId,
        RedirectURL: fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderID),
    }
}
```

这里你一定要看懂 `RedirectURL` 这一行。

旧代码做的事情是：
- 只返回 `message/customer_id/order_id`
- 前端如果要继续流程，得自己知道该跳去哪里

新代码做的事情是：
- 明确告诉前端“下一步去 `/success` 页面”
- 并且把 `customerID`、`orderID` 这两个页面后续查询订单必须用到的参数，一起拼进 URL

为什么这么做？

因为这个页面后面要轮询订单状态，它至少得知道：
- 我在查谁的订单
- 我在查哪一张订单

如果后端不把这两个值带过去，页面就不知道自己该查询谁。

#### 这里补一个 Go 小白容易忘的点：`fmt.Sprintf` 是干什么的

`fmt.Sprintf` 不是打印日志，它是“格式化生成字符串”。

例如：

```go
fmt.Sprintf("id=%s", "123")
```

结果是：

```text
id=123
```

这里作者就是用它来拼这个 URL：

```text
http://localhost:8282/success?customerID=xxx&orderID=yyy
```

所以 `fmt` 这次新增，不是无关 import，而是为了构造跳转链接。

#### 这里的不足也要主动指出

这一行能跑，但不成熟：

1. `localhost:8282` 被写死了
2. URL 参数是手工拼字符串，不够统一
3. 换环境时会很麻烦

生产里通常会：
- 从配置里读前端地址
- 统一拼 URL
- 或交给前端路由自己处理

---

#### 第二件事：查询订单接口的响应结构变了

再看这个函数：

```go
func (H HTTPServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerID string, orderID string) {
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

    // 这里把应用层返回的 order，再转成 client.Order 返回给前端。
    // 也就是说，前端看到的不是 domain 对象原样，而是面向接口的 DTO。
    resp = client.Order{
        CustomerId:  o.CustomerID,
        Id:          o.ID,
        Items:       convertor.NewItemConvertor().EntitiesToClients(o.Items),
        PaymentLink: o.PaymentLink,
        Status:      o.Status,
    }
}
```

注意：你当前工作区里的代码，比 `lesson19 -> lesson20` 当时更往后了。

也就是说：
- 这次 diff 里显示的是“把 `data` 包成 `data.Order`”
- 你当前代码里已经进一步演进成统一 `BaseResponse + DTO` 方案了

这也是为什么你以后读讲义时，要区分：
- `diff` 讲的是当时这组具体改了什么
- 当前工作区代码可能已经比这组更先进一两步

所以这次差异在“当时”的真实含义是：

旧代码：
```json
{
  "message": "success",
  "data": o
}
```

新代码：
```json
{
  "message": "success",
  "data": {
    "Order": o
  }
}
```

为什么这么改？

因为前端页面想更稳定地从 `data` 里取一个明确命名的资源对象，而不是直接拿一个裸对象。

你可以把它理解成“开始做一点点响应结构规范化”。

#### 这里要补的基础知识：`gin.H` 是什么

`gin.H` 本质上就是：

```go
map[string]any
```

只是 Gin 给了它一个简写。

所以：

```go
gin.H{"message": "success"}
```

其实就是在手写一个 JSON 对象。

这类写法很方便，但坏处也明显：
- 字段名很容易到处乱飞
- 响应结构容易每个接口都不一样

所以你当前工作区后面才会演进到 `BaseResponse`、DTO 这种更统一的方式。

### 这一文件你最该带走的结论

`http.go` 这组最重要的变化，不是语法层面的，而是它开始把“创建订单之后前端要怎么继续走流程”这件事明确编码进接口返回值里。`

---

### 4.2 [internal/order/main.go](/g:/shi/go_shop_second/internal/order/main.go)

这个文件 diff 只有一行，但没有这一行，前端页面根本打不开。

先看这个文件自己的 diff。

```diff
diff --git a/internal/order/main.go b/internal/order/main.go
index 96f2ad3..b35cdc8 100644
--- a/internal/order/main.go
+++ b/internal/order/main.go
@@ -46,6 +46,7 @@ func main() {
 	})
 
 	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
+		router.StaticFile("/success", "../../public/success.html")
 		ports.RegisterHandlersWithOptions(router, HTTPServer{
 			app: application,
 		}, ports.GinServerOptions{
```

再看当前代码里的相关部分，并直接注释：

```go
func main() {
    serviceName := viper.GetString("order.service-name")

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // 下面这些初始化：tracing、service discovery、MQ consumer、gRPC server
    // 虽然这一组没重点改它们，但说明 order 服务本身已经是微服务链中的一个节点。
    // 这页 success.html 不是一个独立前端项目，而是临时由 order 服务自己托管。

    server.RunHTTPServer(serviceName, func(router *gin.Engine) {
        // 这一行就是本组核心。
        // 访问 /success 时，Gin 直接返回本地 public/success.html 文件。
        router.StaticFile("/success", "../../public/success.html")

        ports.RegisterHandlersWithOptions(router, HTTPServer{
            app: application,
        }, ports.GinServerOptions{
            BaseURL:      "/api",
            Middlewares:  nil,
            ErrorHandler: nil,
        })
    })
}
```

这行 `StaticFile` 到底是什么意思？

你可以把它理解成：

```text
浏览器 GET /success
=> 不走业务 handler
=> 直接把 success.html 这个静态文件返回回去
```

也就是说，这组不是单独起了个前端服务，而是直接让 `order` 服务顺手托管这张 HTML 页面。

#### 这里补一个 Gin 基础知识：`StaticFile`

Gin 不只是能写接口，还能映射静态文件。

这行：

```go
router.StaticFile("/success", "../../public/success.html")
```

左边是访问路径：
- `/success`

右边是服务器磁盘上的文件：
- `../../public/success.html`

这就是为什么前面 `redirect_url` 拼的是 `/success?...`。

因为服务端已经把 `/success` 这条路由让给这个静态页面了。

#### 新手最容易搞混的点

`/success?customerID=...&orderID=...` 里面：
- `/success` 是路由路径
- `?customerID=...&orderID=...` 是查询参数

Gin 匹配静态文件只看路径 `/success`。

后面的查询参数不是用来找文件的，而是页面加载后由 JavaScript 自己读取的。

这点如果你没分清，就会不明白“一个静态 HTML 为什么还能知道订单号”。

答案是：
- 服务端只负责把 HTML 发出去
- 浏览器里的 JS 再从当前 URL 读参数

### 这一文件你最该带走的结论

`main.go` 这一行改动虽小，但它把“后端 API”真正连到了“浏览器页面入口”上。没有它，前面返回 redirect_url 也只是空谈。`

---

### 4.3 [public/success.html](/g:/shi/go_shop_second/public/success.html)

这是这组最大的新文件，也是你最容易因为 HTML/CSS/JS 混在一起看晕的地方。

先看这个文件自己的 diff。

```diff
diff --git a/public/success.html b/public/success.html
new file mode 100644
index 0000000..08abb53
--- /dev/null
+++ b/public/success.html
@@ -0,0 +1,153 @@
+<!DOCTYPE html>
+<html lang="en">
+<head>
+    <meta charset="UTF-8">
+    <title>Gorder</title>
+</head>
+<body>
+<section>
+  <p>
+    您已成功下单！
+  </p>
+  <p>
+    订单状态：<span id="orderStatus">等待中...</span>
+  </p>
+  <div class="ready-popup">
+    <p>您的订单正在处理中...</p>
+    <p style="color:burlywood; margin:12px">
+      订单号：<b><span id="orderID"></span></b>
+    </p>
+
+    <button class="close-btn" onclick="document.querySelector('.ready-popup').style.display = 'none'">
+      关闭
+    </button>
+  </div>
+
+  <div class="after-payment-popup">
+    <p>等待支付中...</p>
+    <a id="payment-link" href="#">去支付</a>
+  </div>
+</section>
+</body>
+
+<style>
+  ...
+</style>
+
+<script>
+  const urlParam = new URLSearchParams(window.location.search);
+  const customerID = urlParam.get('customerID');
+  const orderID = urlParam.get('orderID');
+  const order = {
+      customerID,
+      orderID,
+      status: 'pending'
+  };
+  const getOrder = async() => {
+      const res = await fetch(`/api/customer/${customerID}/orders/${orderID}`);
+      const data = await res.json();
+
+      if (data.data.Order.Status === 'waiting_for_payment') {
+          order.Status = '等待支付...';
+          document.getElementById('orderStatus').innerText = order.Status;
+          document.querySelector('.after-payment-popup').style.display = 'block';
+          document.getElementById('payment-link').href = data.data.Order.PaymentLink;
+      }
+      if (data.data.Order.Status === 'paid') {
+          order.Status = '已支付成功，请等待...';
+          document.getElementById('orderStatus').innerText = order.Status;
+          setTimeout(getOrder, 5000);
+      } else if (data.data.Order.Status === 'ready') {
+          order.Status = '已完成...';
+          document.querySelector('.after-payment-popup').style.display = 'none';
+          document.querySelector('.ready-popup').style.display = 'block';
+          document.getElementById('orderID').innerText = orderID;
+          document.getElementById('orderStatus').innerText = order.Status;
+      } else {
+          setTimeout(getOrder, 5000);
+      }
+  }
+  getOrder();
+</script>
+</html>
```

这个文件你别一口气全读。正确方式是拆成三块：

1. HTML 骨架
2. CSS 样式
3. JavaScript 状态轮询逻辑

#### 第一块：HTML 骨架到底在干什么

下面我把骨架部分摘出来并写注释：

```html
<section>
  <!-- 页面最上面先提示：下单已经成功了 -->
  <p>
    您已成功下单！
  </p>

  <!-- 这里专门留一个 span，用来动态显示订单状态 -->
  <p>
    订单状态：<span id="orderStatus">等待中...</span>
  </p>

  <!-- 订单最终 ready 时弹出的提示框 -->
  <div class="ready-popup">
    <p>您的订单正在处理中...</p>
    <p style="color:burlywood; margin:12px">
      订单号：<b><span id="orderID"></span></b>
    </p>

    <!-- 点关闭时，直接把 ready-popup 隐藏掉 -->
    <button class="close-btn" onclick="document.querySelector('.ready-popup').style.display = 'none'">
      关闭
    </button>
  </div>

  <!-- 当订单进入 waiting_for_payment 时，显示这个弹框 -->
  <div class="after-payment-popup">
    <p>等待支付中...</p>

    <!-- 这个链接一开始是空的，后面 JS 会把它改成真正的 Stripe 支付链接 -->
    <a id="payment-link" href="#">去支付</a>
  </div>
</section>
```

这里你要看懂的是：

这页不是复杂前端，它本质上只是预先摆好了几个“占位区域”：
- 一个状态文字
- 一个 ready 弹窗
- 一个支付弹窗

后面 JS 根据订单状态，决定显示谁、隐藏谁、往里面填什么数据。

#### 第二块：CSS 在这里的作用是什么

这组你不用死啃 CSS 语法，但至少要懂它在这页里干嘛。

比如：

```css
.ready-popup {
    display: none;
    ...
}

.after-payment-popup {
    display: none;
    ...
}
```

这两段最重要的不是颜色，也不是圆角，而是：

`display: none`

意思是页面刚加载时，这两个弹窗都先隐藏。

为什么？

因为页面刚打开时，前端还没查到订单状态，不知道现在该显示“去支付”还是“已完成”。

等 JS 查完接口之后，才决定：
- 把哪个弹窗改成 `block`
- 让它真正显示出来

所以在这页里，CSS 不是重点知识点，但它承担了“默认隐藏、按状态显示”的作用。

#### 第三块：JavaScript 才是这页的核心

下面我把核心逻辑按注释方式重写给你看：

```js
// 读取当前页面 URL 上的查询参数。
// 例如：/success?customerID=123&orderID=456
const urlParam = new URLSearchParams(window.location.search);

// 从查询参数里拿出 customerID。
const customerID = urlParam.get('customerID');

// 从查询参数里拿出 orderID。
const orderID = urlParam.get('orderID');

// 本地维护一个简单的页面状态对象。
// 它不是后端权威数据，只是页面临时保存状态用。
const order = {
    customerID,
    orderID,
    status: 'pending'
};

// 异步函数：去后端查一次订单状态。
const getOrder = async() => {
    // 调用 order 服务查询订单接口。
    const res = await fetch(`/api/customer/${customerID}/orders/${orderID}`);

    // 把响应体解析成 JSON。
    const data = await res.json();

    // 如果状态是 waiting_for_payment：
    // 说明 payment 服务已经把 Stripe 的 PaymentLink 回写回来了。
    if (data.data.Order.Status === 'waiting_for_payment') {
        order.Status = '等待支付...';

        // 更新页面顶部状态文字。
        document.getElementById('orderStatus').innerText = order.Status;

        // 显示支付弹窗。
        document.querySelector('.after-payment-popup').style.display = 'block';

        // 把后端返回的 PaymentLink 填进“去支付”这个超链接。
        document.getElementById('payment-link').href = data.data.Order.PaymentLink;
    }

    // 如果状态已经 paid：
    // 说明支付已经成功，但订单可能还没进入最终 ready。
    if (data.data.Order.Status === 'paid') {
        order.Status = '已支付成功，请等待...';
        document.getElementById('orderStatus').innerText = order.Status;

        // 5 秒后再查一次。
        setTimeout(getOrder, 5000);

    // 如果状态已经 ready：
    // 说明整个流程基本完成，可以给用户弹完成提示。
    } else if (data.data.Order.Status === 'ready') {
        order.Status = '已完成...';

        // 把支付弹窗隐藏掉。
        document.querySelector('.after-payment-popup').style.display = 'none';

        // 把 ready 弹窗显示出来。
        document.querySelector('.ready-popup').style.display = 'block';

        // 把订单号填进去。
        document.getElementById('orderID').innerText = orderID;
        document.getElementById('orderStatus').innerText = order.Status;

    } else {
        // 如果还不是上述状态，就继续每 5 秒轮询一次。
        setTimeout(getOrder, 5000);
    }
}

// 页面一加载就先查一次订单状态。
getOrder();
```

这段代码你要看懂 3 个最重要的点。

##### 点 1：它怎么知道自己在查哪张订单

靠的是：

```js
const customerID = urlParam.get('customerID');
const orderID = urlParam.get('orderID');
```

也就是说，前面 `http.go` 里返回的 `redirect_url` 不是随便拼着玩的。

它的真正作用是：
- 把订单上下文带到前端页面
- 让这个页面知道“我该轮询谁”

##### 点 2：它怎么拿到支付链接

不是自己算，也不是自己生成。

它只是查：

```js
data.data.Order.PaymentLink
```

这个值是上一组 `payment` 服务调用 Stripe 之后，回写到 `order` 服务里的。

所以这组和上一组的关系非常紧：
- 上一组解决“生成支付链接”
- 这一组解决“把支付链接展示给人点击”

##### 点 3：它为什么一直 `setTimeout(getOrder, 5000)`

因为这是最简单的轮询方案。

什么叫轮询？

不是后端主动告诉页面状态变了，而是页面每隔一段时间自己问一次：

`现在订单怎么样了？`

这里用 5 秒一轮，非常教学版，但好理解。

#### 这里要补的前端基础知识

##### `URLSearchParams`

它是浏览器原生提供的“查询参数解析器”。

你可以把它理解成：
- URL 里 `?a=1&b=2`
- 它帮你把这些值读出来

##### `fetch`

它是浏览器里发 HTTP 请求的 API。

你可以把它理解成“前端世界里调用后端接口”的标准写法。

##### `await`

表示：
- 先等这个异步操作完成
- 再继续执行下面的代码

你先不用把它想得太复杂，当前这页里它就是“等接口返回”的意思。

### 这个文件当前还不完美的地方

这页能跑，但明显只是教学版。

#### 1. 错误处理几乎没有

比如：

```js
const res = await fetch(...)
const data = await res.json()
```

这里没有：
- `try/catch`
- 网络失败处理
- `res.ok` 判断
- 后端返回错误码处理

一旦接口挂了，页面很容易直接异常。

#### 2. 轮询会一直打接口

`setTimeout(getOrder, 5000)` 很简单，但代价是：
- 页面开着就一直请求
- 没有停止条件
- 没有限流和退避策略

教学演示可以，生产不够好。

#### 3. 状态分支不完整

这里只处理了：
- `waiting_for_payment`
- `paid`
- `ready`

没有处理：
- 支付失败
- 订单不存在
- PaymentLink 为空
- 查询异常

所以它只覆盖了 happy path。

#### 4. 页面和接口结构强耦合

它直接依赖：

```js
data.data.Order.Status
```

也就是说，后端 JSON 结构一改，前端也得跟着改。

## 5. 第三方库和易错点

### Gin

这一组 Gin 主要做两件事：
- 在 `http.go` 里处理接口
- 在 `main.go` 里托管静态页面

新手容易忽略：

`Gin 不只是 REST API 框架，也能顺手返回静态文件。`

### 浏览器原生 API

这组虽然没有前端框架，但用到了浏览器原生能力：
- `URLSearchParams`
- `fetch`
- `document.getElementById`
- `document.querySelector`
- `setTimeout`

这说明：

`作者是故意避开复杂前端工程，用最小成本把后端链路展示出来。`

## 6. 为什么这么设计

这组设计非常务实：

`先不要引入 React/Vue，不要先搭完整前端项目，先用一张静态 HTML 页把“下单 -> 查询状态 -> 去支付”这条链跑出来。`

这么做的好处是：

1. 快
2. 简单
3. 后端学习者更容易理解
4. 能立刻验证上一组 Stripe 接入到底有没有真成功

所以这组不是“前端正规化”，而是“后端链路可视化”。

## 7. 当前还不完美的地方

这一组你必须主动看到这些问题：

1. `redirect_url` 写死了地址
2. 页面是演示级 UI，不是产品级 UI
3. 轮询方案简单但低效
4. 页面不负责最终支付确认，真正确认还得靠后面的 webhook
5. 前后端字段风格还不统一，后面还会继续收敛

## 8. 这组最该带走的知识点

1. 一条后端链路跑通后，最好尽快给它做一个最小可观察界面
2. `redirect_url` 的本质是把后续流程入口和上下文一起交给前端
3. 静态 HTML + Gin `StaticFile` 是教学期非常省力的方案
4. 轮询是最简单的状态同步方式，但不是最终最佳实践
5. 页面只负责展示状态，不负责改变状态
6. 本组能显示“去支付”，完全依赖上一组已经把 `PaymentLink` 回写进订单了

## 9. 一句话收住这组

`lesson19 -> lesson20` 的本质，不是“加了个前端页面”，而是“把原本只能靠日志和接口验证的支付链，第一次做成了一个人能直接点击和观察的演示流程”。`
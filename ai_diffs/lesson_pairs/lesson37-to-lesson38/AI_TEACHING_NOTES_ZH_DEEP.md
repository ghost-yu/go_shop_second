# `lesson37 -> lesson38` 独立讲义（详细注释版）

这一组的主题，表面上看是“加了测试”。

但如果你只把它理解成“多写了一个 `go test` 文件”，那还是看浅了。

它真正做的是把前面几组已经逐步成形的 HTTP 契约，第一次系统地闭环起来：

1. 统一 response 外壳，不再只是运行时约定
2. OpenAPI 文档正式声明返回值是 `Response`
3. 代码生成器同步生成新的 client/server 类型
4. 前端成功页脚本按新返回结构修正
5. 测试开始真正调用 HTTP 接口，验证这套契约没有写崩

所以这组的核心不是“单元测试技术细节”，而是：

`让“接口长什么样”从口头约定，变成文档、生成代码、前端脚本、测试代码四方一致。`

字幕 `39.TXT` 其实就在讲这件事：
- 用生成好的 OpenAPI client 来做测试
- 把 response 结构补进 OpenAPI 文档
- 让测试自动验证接口

这对 Go 小白非常重要，因为很多人刚开始学接口开发时会犯一个典型错误：
- 服务端返回一套 JSON
- 文档写另一套
- 前端按第三套结构取字段
- 大家都觉得“差不多”

最后问题就会集中爆发。

这组课就是在防这个。

## 1. 先说明一个小情况

`lesson35 -> lesson36` 和 `lesson36 -> lesson37` 在你这份仓库里都没有代码差异，所以我直接跳到了下一组真实有变化的 `lesson37 -> lesson38`。

## 2. 这组你应该怎么读

建议顺序：

1. [api/openapi/order.yml](/g:/shi/go_shop_second/api/openapi/order.yml)
2. [api/openapi/cfg.yaml](/g:/shi/go_shop_second/api/openapi/cfg.yaml)
3. [scripts/genopenapi.sh](/g:/shi/go_shop_second/scripts/genopenapi.sh)
4. [internal/order/tests/create_order_test.go](/g:/shi/go_shop_second/internal/order/tests/create_order_test.go)
5. [public/success.html](/g:/shi/go_shop_second/public/success.html)
6. [internal/order/adapters/order_mongo_repository.go](/g:/shi/go_shop_second/internal/order/adapters/order_mongo_repository.go)
7. [internal/payment/http.go](/g:/shi/go_shop_second/internal/payment/http.go)

为什么这样看：
- `order.yml` 是这组“契约源头”的核心。
- `cfg.yaml + genopenapi.sh` 解释为什么生成代码终于能保留 `Response` 这种当前 schema 没直接引用的类型。
- `create_order_test.go` 展示“怎么用生成 client 做接口测试”。
- `success.html` 则说明前端脚本也必须跟着 response 结构变化调整。
- `order_mongo_repository.go` 那一行虽然小，但会影响测试和更新行为是否正确。

## 3. 这组到底解决什么问题

先回忆一下 lesson35 之前的状态。

前面几组已经做了这些事：
- HTTP 返回被统一包进 `BaseResponse`
- JSON 字段名被统一成 `snake_case`
- OpenAPI 路径/字段也改成了 `snake_case`
- order 服务已经能用 Mongo 持久化

但这里仍有一个明显缺口：

`OpenAPI 文档里“200 返回什么”，还没有真正跟运行时 response 外壳对齐。`

也就是说，当时会出现这种情况：
- 代码实际返回：`{ errno, message, data, trace_id }`
- OpenAPI 文档可能还写的是：直接返回 `Order`
- 生成 client 也会以为 `200` 返回的是 `Order`

这会导致：
- 测试写不顺
- 前端和 client 很容易按错结构解析
- 生成代码的 `JSON200` 类型和真实返回体对不上

所以 lesson38 的主要动作是：

1. 把 `Response` 结构正式写进 OpenAPI 文档
2. 调整代码生成配置，让这个 schema 不会被自动裁掉
3. 写一个真正调用接口的测试，用生成 client 去验证 response
4. 顺手修前端成功页读取字段的方式

## 4. 契约源头先修正：`api/openapi/order.yml`

### 4.1 这个文件自己的原始 diff

```diff
diff --git a/api/openapi/order.yml b/api/openapi/order.yml
index da9fbf8..6155cc0 100644
--- a/api/openapi/order.yml
+++ b/api/openapi/order.yml
@@ -32,7 +32,7 @@ paths:
           content:
             application/json:
               schema:
-				$ref: '#/components/schemas/Order'
+				$ref: '#/components/schemas/Response'
@@ -64,7 +64,7 @@ paths:
           content:
             application/json:
               schema:
-				$ref: '#/components/schemas/Order'
+				$ref: '#/components/schemas/Response'
@@ -144,4 +144,21 @@ components:
           type: string
         quantity:
           type: integer
-		  format: int32
+		  format: int32
+
+	Response:
+	  type: object
+	  properties:
+		errno:
+		  type: integer
+		message:
+		  type: string
+		data:
+		  type: object
+		trace_id:
+		  type: string
+	  required:
+		- errno
+		- message
+		- data
+		- trace_id
```

### 4.2 旧代码做什么，新代码做什么

旧文档：
- `GET /orders/{id}` 的 `200` 响应写成 `Order`
- `POST /orders` 的 `200` 响应也写成 `Order`

但这和真实运行时已经不一致了。

因为真实接口其实返回的是：

```json
{
  "errno": 0,
  "message": "success",
  "data": {...},
  "trace_id": "..."
}
```

新文档：
- 明确把 `200` 返回写成 `Response`
- 并正式定义 `Response` schema

这一步非常关键，因为它是在把“运行时代码里的约定”正式升格成“文档契约”。

### 4.3 关键代码和详细注释

```yaml
paths:
  /customer/{customer_id}/orders/{order_id}:
    get:
      responses:
        '200':
          description: todo
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Response'

  /customer/{customer_id}/orders:
    post:
      responses:
        '200':
          description: todo
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Response'
```

这段最重要的不是 YAML 语法，而是语义：

`成功响应不再假装自己是裸 Order，而是明确声明：外层是统一 response 包装。`

再看 schema：

```yaml
Response:
  type: object
  properties:
    errno:
      type: integer
    message:
      type: string
    data:
      type: object
    trace_id:
      type: string
  required:
    - errno
    - message
    - data
    - trace_id
```

这里你要注意几个点：

#### 1. `data` 被定义成 `object`

这意味着当前 OpenAPI 文档对 `data` 的类型描述还比较宽。

它没有精确到：
- 创建订单时 `data` 长什么样
- 获取订单时 `data` 长什么样

这说明课程当前先追求的是“先把统一 response 外壳声明出来”，
还没有进一步把每个接口的 `data` 做成特别精确的泛型/嵌套 schema。

#### 2. `required` 很重要

它表示：
- 这四个字段在契约层都是必须存在的
- 所以后面生成代码时，`Response` 会长得更稳定

### 4.4 为什么这一步值得认真学

因为很多初学者写接口文档时会偷懒：
- 代码怎么返回是一回事
- OpenAPI 里先随便写个 `Order`
- 觉得“调用方自己看着办”

这会直接导致生成 client 错误。

而这一组正好在纠这个问题：

`文档不是装饰品，它会影响生成代码，所以必须和真实返回保持一致。`

## 5. 为什么要加 `cfg.yaml`：`api/openapi/cfg.yaml`

### 5.1 这个文件自己的原始 diff

```diff
diff --git a/api/openapi/cfg.yaml b/api/openapi/cfg.yaml
new file mode 100644
index 0000000..3ac1679
--- /dev/null
+++ b/api/openapi/cfg.yaml
@@ -0,0 +1,2 @@
+output-options:
+   skip-prune: true
```

### 5.2 这段在做什么

这是给 `oapi-codegen` 的配置文件。

里面最关键的一行是：

```yaml
skip-prune: true
```

### 5.3 `prune` 是什么，为什么要 `skip-prune`

你可以把 `prune` 理解成：

`代码生成器在生成时，会尝试把“它觉得没被引用的 schema”裁掉。`

这样做通常是为了让生成结果更干净，不要生成一堆没用类型。

但这次课程碰到的问题是：
- 你虽然在 schema 里定义了 `Response`
- 但生成器不一定总能按你预期保留它
- 或者某些生成路径里，它可能把没直接引用到的类型裁掉

所以这里显式加上：

```yaml
skip-prune: true
```

含义就是：

`别替我自动裁类型，我宁愿多保留一些，也不要把我需要的 Response 干掉。`

### 5.4 为什么这对初学者很重要

因为你会第一次碰到一个非常真实的工程问题：

`不是你写了 OpenAPI schema，生成器就一定完全按你直觉保留所有东西。`

代码生成器有自己的策略。

所以遇到“为什么我 schema 写了，但生成代码里没有这个类型”的情况时，
你要知道可能不是你代码错了，而是生成配置的问题。

## 6. 生成脚本为什么也要跟着改：`scripts/genopenapi.sh`

### 6.1 这个文件自己的原始 diff

```diff
diff --git a/scripts/genopenapi.sh b/scripts/genopenapi.sh
index 421e3f4..cc3566d 100755
--- a/scripts/genopenapi.sh
+++ b/scripts/genopenapi.sh
@@ -42,11 +42,11 @@ function gen() {
-  run oapi-codegen -generate types -o "$output_dir/openapi_types.gen.go" -package "$package" "api/openapi/$service.yml"
-  run oapi-codegen -generate "$GEN_SERVER" -o "$output_dir/openapi_api.gen.go" -package "$package" "api/openapi/$service.yml"
+  run oapi-codegen -config api/openapi/cfg.yaml -generate types -o "$output_dir/openapi_types.gen.go" -package "$package" "api/openapi/$service.yml"
+  run oapi-codegen -config api/openapi/cfg.yaml -generate "$GEN_SERVER" -o "$output_dir/openapi_api.gen.go" -package "$package" "api/openapi/$service.yml"
 
-  run oapi-codegen -generate client -o "internal/common/client/$service/openapi_client.gen.go" -package "$service" "api/openapi/$service.yml"
-  run oapi-codegen -generate types -o "internal/common/client/$service/openapi_types.gen.go" -package "$service" "api/openapi/$service.yml"
+  run oapi-codegen -config api/openapi/cfg.yaml -generate client -o "internal/common/client/$service/openapi_client.gen.go" -package "$service" "api/openapi/$service.yml"
+  run oapi-codegen -config api/openapi/cfg.yaml -generate types -o "internal/common/client/$service/openapi_types.gen.go" -package "$service" "api/openapi/$service.yml"
 }
```

### 6.2 这段为什么重要

因为你新增了 `cfg.yaml`，如果生成脚本不使用它，那这个配置文件就是摆设。

所以这里真正做的是：

`把“生成配置”纳入正式生成流程。`

### 6.3 关键代码和详细注释

```bash
run oapi-codegen -config api/openapi/cfg.yaml -generate types -o "$output_dir/openapi_types.gen.go" -package "$package" "api/openapi/$service.yml"
run oapi-codegen -config api/openapi/cfg.yaml -generate "$GEN_SERVER" -o "$output_dir/openapi_api.gen.go" -package "$package" "api/openapi/$service.yml"

run oapi-codegen -config api/openapi/cfg.yaml -generate client -o "internal/common/client/$service/openapi_client.gen.go" -package "$service" "api/openapi/$service.yml"
run oapi-codegen -config api/openapi/cfg.yaml -generate types -o "internal/common/client/$service/openapi_types.gen.go" -package "$service" "api/openapi/$service.yml"
```

这几行看似只是多了个 `-config`，但它意味着：
- 以后所有 server/client/types 生成都统一遵循同一套规则
- 不是“手工本地偶尔用某个参数”，而是正式脚本级约束

这在工程里非常重要，因为你不希望：
- A 同学本地生成一套
- B 同学本地生成另一套
- CI 又生成第三套

## 7. 生成出来的 client 为什么终于能正确认识 `Response`

### 7.1 生成结果的变化

这组生成后，像这样的类型就出现了：

```go
type Response struct {
    Data    map[string]interface{} `json:"data"`
    Errno   int                    `json:"errno"`
    Message string                 `json:"message"`
    TraceId string                 `json:"trace_id"`
}
```

以及：

```go
type PostCustomerCustomerIdOrdersResponse struct {
    Body         []byte
    HTTPResponse *http.Response
    JSON200      *Response
    JSONDefault  *Error
}
```

### 7.2 这意味着什么

意味着从这组开始，生成 client 不再以为：
- `200` 返回的是 `Order`

而是知道：
- `200` 返回的是统一 `Response`

这对测试非常关键，因为后面测试代码才可以写：

```go
assert.Equal(t, 0, response.JSON200.Errno)
```

如果生成 client 还以为 `JSON200` 是 `Order`，这类断言就根本写不通。

## 8. 这组的主角之一：测试终于真正开始校验 HTTP 契约

### 8.1 `internal/order/tests/create_order_test.go`

#### 8.1.1 这个文件自己的原始 diff

```diff
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
+...
```

#### 8.1.2 这段的真正意义

这不是“纯内部函数单元测试”。

它更接近一种：
- 用生成好的 OpenAPI client 调实际 HTTP 接口
- 再检查返回结构是否符合预期

严格说更偏向轻量集成测试 / API 测试。

但课程里叫单元测试也能理解，因为它是在服务模块内部第一次认真写自动化测试。

#### 8.1.3 关键代码和详细注释

```go
var (
    ctx    = context.Background()
    server = fmt.Sprintf("http://%s/api", viper.GetString("order.http-addr"))
)
```

这里做的事很简单：
- 用配置拼出 order 服务地址
- 测试不手写死完整 URL，只依赖配置

这说明测试已经开始和服务真实启动方式接轨。

继续看入口：

```go
func TestMain(m *testing.M) {
    before()
    m.Run()
}

func before() {
    log.Printf("server=%s", server)
}
```

这里你要知道 `TestMain` 是什么：
- 它是 Go 测试的自定义入口
- 可以在整个测试包跑之前做一些初始化
- 跑完后也可以做清理

课程这里先只打印日志，但它给你演示了一个真实入口：
- 以后你完全可以在这里做数据库初始化、测试环境检查、清理逻辑

再看第一个测试：

```go
func TestCreateOrder_success(t *testing.T) {
    response := getResponse(t, "123", sw.PostCustomerCustomerIdOrdersJSONRequestBody{
        CustomerId: "123",
        Items: []sw.ItemWithQuantity{
            {
                Id:       "test-item-1",
                Quantity: 1,
            },
        },
    })
    t.Logf("body=%s", string(response.Body))
    assert.Equal(t, 200, response.StatusCode())

    assert.Equal(t, 0, response.JSON200.Errno)
}
```

这段测试你至少要看懂三层：

##### 第 1 层：它真的是在用生成 client 调 HTTP

不是伪造 handler 输入，不是直接调函数。

而是：
- `NewClientWithResponses(server)`
- `PostCustomerCustomerIdOrdersWithResponse(...)`

这意味着它在校验：
- 路由
- JSON 序列化
- HTTP 状态码
- 返回结构解析

##### 第 2 层：它在校验统一 response 外壳

```go
assert.Equal(t, 0, response.JSON200.Errno)
```

注意它没有直接断言 `order_id` 之类字段。
它先断言的是：
- `JSON200` 能被解析成 `Response`
- `Errno == 0`

也就是说，这组测试最优先关心的是：

`统一 response 契约是不是对的。`

##### 第 3 层：它还没有把断言做得特别深

比如它没有继续断言：
- `response.JSON200.Data` 里面具体有哪些字段
- `redirect_url` 是否存在

这说明测试是第一版，先把最关键的成功/失败外壳校验起来。

再看失败参数测试：

```go
func TestCreateOrder_invalidParams(t *testing.T) {
    response := getResponse(t, "123", sw.PostCustomerCustomerIdOrdersJSONRequestBody{
        CustomerId: "123",
        Items:      nil,
    })
    assert.Equal(t, 200, response.StatusCode())
    assert.Equal(t, 2, response.JSON200.Errno)
}
```

这段很重要，因为它再次验证了前面那条接口约定：
- 即使参数错误，HTTP 还是 200
- 真正的业务失败看 `errno`

这就把“项目的响应语义”明确固定下来了。

### 8.2 为什么测试里要用生成 client，而不是手写 `http.NewRequest`

因为这有几个好处：

1. 它直接验证 OpenAPI 生成代码是否能正常使用
2. 测试代码更贴近真实调用方体验
3. 能更早发现“文档契约和服务实现不一致”的问题

这其实是这组最有价值的工程思路之一：

`用文档生成的 client 去测服务本身。`

这样文档和实现如果偏了，测试会第一时间把你揪出来。

## 9. 前端成功页为什么也要改：`public/success.html`

### 9.1 这个文件自己的原始 diff

```diff
diff --git a/public/success.html b/public/success.html
index 08abb53..1dd4267 100644
--- a/public/success.html
+++ b/public/success.html
@@ -128,22 +128,22 @@
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
       }
```

### 9.2 为什么这段很重要

这段是这组最直接的“前端适配证据”。

它说明前面统一 response 外壳、snake_case 命名、客户端类型调整这些事，不只是后端自己觉得整齐。

它真的会影响前端取字段方式。

### 9.3 关键代码和详细注释

```js
if (data.data.status === 'waiting_for_payment') {
    order.status = '等待支付...';
    document.getElementById('orderStatus').innerText = order.status;
    document.querySelector('.after-payment-popup').style.display = 'block';
    document.getElementById('payment-link').href = data.data.payment_link;
}
```

这段意味着：
- 现在 `data` 里直接就是订单字段
- 不再是 `data.Order.Status`
- 字段名也跟 OpenAPI/JSON 一致变成了 `status` / `payment_link`

这其实正好反过来证明上一组 lesson30/32 的整理不是白做的。

因为一旦契约真正统一，前端代码就会更简单：
- 少一层 `Order`
- 字段名统一
- 不需要到处猜大小写

### 9.4 这里你要学会一个工程事实

前后端接口改动不是“后端自己改就完了”。

只要返回 JSON 结构变了，通常要同步影响：
- OpenAPI 文档
- 生成 client
- 测试
- 前端脚本

lesson38 正是在完整演示这条连锁反应。

## 10. 那一行 Mongo 仓储的小改动为什么值得讲

### 10.1 `internal/order/adapters/order_mongo_repository.go`

#### 10.1.1 这个文件自己的原始 diff

```diff
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
```

### 10.2 这行改动在做什么

旧代码：
- `updateFn` 拿到的是 `oldOrder`

新代码：
- `updateFn` 拿到的是传进来的 `order`

### 10.3 为什么这行会影响行为

你要先记住：
- `oldOrder` 是仓储刚从数据库里查出来的旧值
- `order` 是调用方传进来的目标对象

这两者不是一回事。

如果 `updateFn` 是基于新值去决定更新内容，那么把 `oldOrder` 塞进去，语义就会有偏差。

这行看起来小，但本质上是在修：

`Update 回调到底应该基于谁来生成最终更新对象。`

这会直接影响 Mongo 更新逻辑是不是按预期工作，尤其在支付状态更新、支付链接回写这些流程里非常关键。

### 10.4 为什么这类 bug 容易漏掉

因为它通常：
- 能编译
- 代码看着也“差不多”
- 但实际业务结果会偏

这类 bug 最适合被什么发现？
- 测试

所以它和这组新增测试其实是呼应的：

`当系统开始变复杂，小而真实的行为测试就越来越有价值。`

## 11. 这一组还顺手动了什么

### 11.1 `internal/payment/http.go`

这个文件只多了一点小改动，不是这组主线核心。

### 11.2 `openapi_types.gen.go / openapi_client.gen.go`

这些都是 `order.yml + genopenapi.sh + cfg.yaml` 连锁生成出来的结果。

你要记住：
- 这些文件是结果
- 不是源头

真正要学的是：
- 文档怎么改
- 生成配置怎么改
- 为什么需要这么改

## 12. 这组最大的价值是什么

如果你让我用一句话总结 lesson38 的价值，我会说：

`它第一次把“接口契约一致性”从口头整理，推进成了文档、生成代码、前端和测试都同时对齐的状态。`

这对工程质量非常关键。

因为系统一旦开始有：
- OpenAPI
- 自动生成 client
- 前端脚本
- 集成测试

你就不能再靠“差不多能跑”活着了。

## 13. 这组的局限和你要注意的点

### 13.1 `Response.data` 现在还是 `object`

这意味着 OpenAPI 还没有把不同接口的 `data` 描述得特别精确。

比如：
- 创建订单时 `data` 的字段集合
- 获取订单时 `data` 的字段集合

还没有拆成更强类型的 schema。

所以这组解决的是“先统一外壳”，不是“response 泛型建模已完美”。

### 13.2 测试更像轻量集成测试，不是纯粹 unit test

因为它真在调用 HTTP 服务。

这不是坏事，但你要知道它和“直接测某个函数”的单元测试不是一回事。

### 13.3 测试断言还比较浅

现在主要断言：
- HTTP 200
- `errno` 正确

但还没有深挖：
- `data` 具体字段
- `trace_id` 是否存在
- 错误 message 是否符合预期

这说明测试体系是刚起步。

## 14. 这组最该记住的话

1. 统一 response 外壳，必须同步写进 OpenAPI 文档，否则生成 client 会错。
2. `skip-prune: true` 的价值是避免生成器把你需要的 schema 裁掉。
3. 用生成好的 OpenAPI client 去测服务，是非常实用的工程手段。
4. 前端读取字段方式会直接受接口契约变化影响。
5. 那一行 Mongo `updateFn(ctx, order)` 的改动虽然小，但属于典型行为修复。

## 15. 你现在应该怎么复习这组

建议顺序：

1. [api/openapi/order.yml](/g:/shi/go_shop_second/api/openapi/order.yml)
   - 先看 `200` 返回为什么改成 `Response`
2. [api/openapi/cfg.yaml](/g:/shi/go_shop_second/api/openapi/cfg.yaml)
   - 理解 `skip-prune` 在解决什么问题
3. [scripts/genopenapi.sh](/g:/shi/go_shop_second/scripts/genopenapi.sh)
   - 看生成流程怎么正式接入这个配置
4. [internal/order/tests/create_order_test.go](/g:/shi/go_shop_second/internal/order/tests/create_order_test.go)
   - 看生成 client 怎么被真正拿来测接口
5. [public/success.html](/g:/shi/go_shop_second/public/success.html)
   - 看前端为什么必须跟着契约改
6. [internal/order/adapters/order_mongo_repository.go](/g:/shi/go_shop_second/internal/order/adapters/order_mongo_repository.go)
   - 理解那个一行更新为什么会影响实际行为

如果你继续，我下一组就写 `lesson38 -> lesson39`。
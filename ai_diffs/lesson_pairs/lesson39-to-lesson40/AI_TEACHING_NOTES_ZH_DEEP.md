# `lesson39 -> lesson40` 独立讲义（详细注释版）

这一组不是在加新业务，而是在补“排障能力”。

如果你把前面课程理解成“先把订单、支付、库存、MQ、追踪都串起来”，那这一组做的事就是：

- 请求刚进服务时，记录一份更完整的入站日志
- 请求处理完时，再记录一份出站日志
- 把最终响应内容也挂到日志里
- 把 `command` 和 `query` 的日志、指标彻底分开

这类改动对初学者最容易被低估，因为它不直接新增 API，也不新增数据库表。但系统一旦复杂起来，没有这类可观测性代码，后面排查问题会非常痛苦。

## 1. 正确阅读顺序

建议按这个顺序读：

1. [internal/common/middleware/request.go](/g:/shi/go_shop_second/internal/common/middleware/request.go)
2. [internal/common/response.go](/g:/shi/go_shop_second/internal/common/response.go)
3. [internal/common/server/http.go](/g:/shi/go_shop_second/internal/common/server/http.go)
4. [internal/common/decorator/logging.go](/g:/shi/go_shop_second/internal/common/decorator/logging.go)
5. [internal/common/decorator/metrics.go](/g:/shi/go_shop_second/internal/common/decorator/metrics.go)
6. [internal/common/decorator/command.go](/g:/shi/go_shop_second/internal/common/decorator/command.go)
7. [internal/common/middleware/logger.go](/g:/shi/go_shop_second/internal/common/middleware/logger.go)

原因是：

- `request.go` 是这组最核心的新能力。
- `response.go` 是它的配套。没有它，`request_out` 拿不到响应体。
- `http.go` 负责把中间件真正挂上 HTTP 服务。
- 后面的 `logging.go`、`metrics.go`、`command.go` 是把应用层装饰器的语义修正干净。

## 2. 先把整条调用链理顺

这组之后，一次 HTTP 请求的大致路径变成：

`gin 收到请求 -> RequestLog 中间件记录 request_in -> 业务 handler 执行 -> BaseResponse 统一写响应并把 response 放进 gin.Context -> defer 的 requestOut 记录 request_out -> 你在日志系统里同时看到入站和出站`

应用层另一条线是：

`handler 调 command/query -> 对应 decorator 记录日志和指标 -> command 不再伪装成 query`

这两条线一条在 HTTP 入口，一条在应用层执行点。它们组合起来，日志才真正有上下文。

## 3. 这组到底解决什么问题

### 3.1 之前的问题

之前虽然已经有：

- trace
- 基础日志
- 统一响应结构
- command/query decorator

但还缺几块关键内容：

1. HTTP 请求体看不清
2. HTTP 响应体看不清
3. 只能看到“某个 handler 跑了”，不容易快速对上输入输出
4. `command` 的日志和指标还沿用了 `query` 的命名，语义不准确

### 3.2 这一组的核心目标

这一组不是让业务更强，而是让你更容易回答这些问题：

- 用户刚才到底传了什么参数？
- 服务最后到底返回了什么？
- 这个请求耗时多少毫秒？
- 当前执行的是 command 还是 query？
- 指标里到底统计的是查询还是命令？

这就是“可观测性”的实际价值。它不是一个抽象口号，而是让你出问题时有证据可查。

## 4. 关键文件一：`internal/common/middleware/request.go`

### 4.1 这个文件自己的原始 diff

```diff
diff --git a/internal/common/middleware/request.go b/internal/common/middleware/request.go
new file mode 100644
index 0000000..3b90df0
--- /dev/null
+++ b/internal/common/middleware/request.go
@@ -0,0 +1,44 @@
+package middleware
+
+import (
+	"bytes"
+	"encoding/json"
+	"io"
+	"time"
+
+	"github.com/gin-gonic/gin"
+	"github.com/sirupsen/logrus"
+)
+
+func RequestLog(l *logrus.Entry) gin.HandlerFunc {
+	return func(c *gin.Context) {
+		requestIn(c, l)
+		defer requestOut(c, l)
+		c.Next()
+	}
+}
+
+func requestOut(c *gin.Context, l *logrus.Entry) {
+	response, _ := c.Get("response")
+	start, _ := c.Get("request_start")
+	startTime := start.(time.Time)
+		l.WithContext(c.Request.Context()).WithFields(logrus.Fields{
+		"proc_time_ms": time.Since(startTime).Milliseconds(),
+		"response":     response,
+	}).Info("__request_out")
+}
+
+func requestIn(c *gin.Context, l *logrus.Entry) {
+	c.Set("request_start", time.Now())
+	body := c.Request.Body
+	bodyBytes, _ := io.ReadAll(body)
+	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
+	var compactJson bytes.Buffer
+	_ = json.Compact(&compactJson, bodyBytes)
+	l.WithContext(c.Request.Context()).WithFields(logrus.Fields{
+		"start": time.Now().Unix(),
+		"args":  compactJson.String(),
+		"from":  c.RemoteIP(),
+		"uri":   c.Request.RequestURI,
+	}).Info("__request_in")
+}
```

### 4.2 这段新代码在项目里干什么

它定义了一个新的 Gin 中间件 `RequestLog`。

中间件你可以先简单理解成：

`在真正业务 handler 前后插一层公共逻辑`

这里插进去的公共逻辑就是：

- 业务开始前打 `__request_in`
- 业务结束后打 `__request_out`

这样每个 HTTP 请求都会有一进一出两条日志。

### 4.3 带中文注释的关键代码

```go
package middleware

import (
    // bytes 包用于处理字节切片和内存缓冲区。
    // 这里主要用它把读出来的 body 再包装回一个“可再次读取”的 Reader。
    "bytes"

    // encoding/json 提供 JSON 编解码能力。
    // 这里不是做业务反序列化，而是做日志里的 JSON 压缩整理。
    "encoding/json"

    // io 提供统一 I/O 抽象。
    // 这里用到了 io.ReadAll 和 io.NopCloser。
    "io"

    // time 用来记开始时间、计算处理耗时。
    "time"

    // gin 是当前项目的 HTTP 框架。
    // gin.Context 是一次请求的上下文对象，很多数据都通过它传递。
    "github.com/gin-gonic/gin"

    // logrus 是结构化日志库。
    // Entry 可以理解成“带上下文字段的 logger”。
    "github.com/sirupsen/logrus"
)

func RequestLog(l *logrus.Entry) gin.HandlerFunc {
    // gin.HandlerFunc 本质上是：func(*gin.Context)
    // 也就是“接收当前请求上下文的处理函数”。
    return func(c *gin.Context) {
        // 先记录请求进入时的信息。
        requestIn(c, l)

        // defer 的意思是：当前函数 return 之前，一定执行这句后面的函数。
        // 这里就形成了一个经典模式：
        // 进入时做一件事，退出时再做一件事。
        defer requestOut(c, l)

        // c.Next() 是 Gin 中间件里非常关键的方法。
        // 它的意思不是“下一行”，而是“把控制权交给后续中间件和最终 handler 继续处理”。
        // 如果不调用它，这个请求链就会停在这里，真正业务不会执行。
        c.Next()
    }
}

func requestIn(c *gin.Context, l *logrus.Entry) {
    // c.Set(key, value) 是 Gin 提供的上下文存值方法。
    // 它把数据挂到当前请求的 Context 上，供后面的代码复用。
    // 这里存的是请求开始时间，给 requestOut 算耗时用。
    c.Set("request_start", time.Now())

    // c.Request 是标准库 http.Request。
    // Body 是请求体，是一个 io.ReadCloser。
    // ReadCloser 你可以先理解成：既能读，又能关闭的流。
    body := c.Request.Body

    // io.ReadAll 会把 body 里的全部内容一次性读成 []byte。
    // 输入：一个 Reader
    // 输出：完整字节切片 + error
    // 这里这样做的原因是：我们想把请求 JSON 原文打进日志。
    bodyBytes, _ := io.ReadAll(body)

    // 这是整个文件里最关键、也最容易漏掉的一行。
    // 原因：Request.Body 默认是“读一遍就没了”的流。
    // 如果这里把 body 读完，却不重新放回去，后面的 ShouldBindJSON / BindJSON 就读不到内容了。
    // bytes.NewReader(bodyBytes) 会基于这份字节切片创建一个新的 Reader。
    // io.NopCloser(...) 会给这个 Reader 包一层假的 Close 方法，让它满足 io.ReadCloser 接口。
    // 最终效果就是：我们把“读过的 body”重新塞回请求对象，后续 handler 还能正常继续读。
    c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

    // bytes.Buffer 是一个可增长的内存缓冲区，经常用来接收字符串或字节输出。
    var compactJson bytes.Buffer

    // json.Compact(dst, src) 会把 JSON 压成紧凑格式：
    // 去掉多余空格、换行、缩进。
    // 这里的目的不是业务转换，而是让日志里的一整段 JSON 更好搜、更省空间。
    _ = json.Compact(&compactJson, bodyBytes)

    // WithContext(ctx) 把标准 context 附着到日志条目上。
    // 如果后面 logger hook 或 tracing 集成要读 context，这里会有用。
    // WithFields(...) 是 logrus 的结构化日志核心能力：
    // 它不是简单拼字符串，而是挂一组 key-value 字段。
    l.WithContext(c.Request.Context()).WithFields(logrus.Fields{
        // Unix() 返回秒级时间戳。
        // 这里作者另外记了一个开始时间字段，便于在日志系统里直接看原始时间。
        "start": time.Now().Unix(),

        // compactJson.String() 把压缩后的 JSON 缓冲区转成字符串。
        "args": compactJson.String(),

        // RemoteIP() 是 Gin 提供的方法，用于拿远端地址。
        // 初学者要注意：真实线上环境里它是否可信，跟反向代理配置有关。
        "from": c.RemoteIP(),

        // RequestURI 是标准库 http.Request 的字段，包含原始请求 URI。
        "uri": c.Request.RequestURI,
    }).Info("__request_in")
}

func requestOut(c *gin.Context, l *logrus.Entry) {
    // c.Get(key) 从 Gin 上下文取值。
    // 返回两个值：value 和 是否存在。
    // 这里代码直接忽略了第二个布尔值，这是课程代码里常见的简化写法。
    response, _ := c.Get("response")
    start, _ := c.Get("request_start")

    // start.(time.Time) 是类型断言。
    // 因为 c.Get 返回的是 any，真正使用前要还原成具体类型。
    // 如果这里拿到的不是 time.Time，会 panic。
    // 之所以敢这么写，是因为作者默认 requestIn 一定已经提前塞入了 time.Time。
    startTime := start.(time.Time)

    l.WithContext(c.Request.Context()).WithFields(logrus.Fields{
        // time.Since(t) 等价于 time.Now().Sub(t)。
        // Milliseconds() 再把 duration 转成毫秒整数。
        "proc_time_ms": time.Since(startTime).Milliseconds(),

        // response 是前面 response.go 里塞进 gin.Context 的响应 JSON 字符串。
        "response": response,
    }).Info("__request_out")
}
```

### 4.4 这里涉及到的库函数，单独拆开讲

#### `c.Next()` 是什么

来源：`github.com/gin-gonic/gin`

作用：

- 让 Gin 继续执行后续中间件和最终业务 handler。

为什么这里必须调用：

- 中间件不是自动“穿透”的。
- 你不调用 `c.Next()`，请求就被卡在这里，后面的接口逻辑不会跑。

初学者最容易误解：

- 以为它只是“执行下一行逻辑”。不是。
- 它实际上是在推动整个 Gin handler 链继续前进。

#### `c.Set()` / `c.Get()` 是什么

来源：`gin.Context`

作用：

- 在同一次请求范围内，给多个中间件 / handler 共享数据。

这里为什么用它：

- `requestIn` 里存开始时间
- `response.go` 里存响应 JSON
- `requestOut` 再把它们取出来做日志

常见坑：

- `c.Get()` 返回的是 `any`，后面要自己做类型断言。
- key 写错不会编译报错，只会运行时取不到值。

#### `io.ReadAll()` 是什么

来源：标准库 `io`

作用：

- 把一个 Reader 里的内容全部读出来，变成 `[]byte`。

这里为什么用它：

- 因为日志想打印整个请求体。

常见坑：

- 对大 body 直接 `ReadAll` 会把所有内容都放内存，不适合无限大请求。
- 这里是教学项目，所以写法直接，但生产里通常还会考虑大小限制和敏感字段脱敏。

#### `io.NopCloser()` 是什么

来源：标准库 `io`

作用：

- 给一个只有 `Read` 能力的 Reader 包装出一个假的 `Close` 方法。
- 包装后它就满足 `io.ReadCloser` 接口。

为什么这里一定要它：

- `c.Request.Body` 的类型要求是 `io.ReadCloser`。
- `bytes.NewReader(bodyBytes)` 只是 Reader，不带 Close。
- 所以要用 `io.NopCloser` 补齐接口。

#### `json.Compact()` 是什么

来源：标准库 `encoding/json`

作用：

- 把 JSON 压缩成紧凑格式。

这里为什么用它：

- 为了日志更干净，避免多行缩进 JSON 搞乱日志输出。

常见坑：

- 如果请求体不是合法 JSON，`Compact` 会报错。
- 这份课程代码直接忽略 error，所以非 JSON body 时日志字段可能为空或不准确。

#### `WithContext().WithFields().Info()` 是什么

来源：`logrus`

链路拆开理解：

- `WithContext(ctx)`：把 context 挂进去
- `WithFields(fields)`：追加结构化字段
- `Info(msg)`：按 Info 级别真正输出日志

为什么不用 `fmt.Println`：

- 因为结构化日志可以被日志平台按字段检索，比如按 `uri`、`proc_time_ms`、`command` 搜索。

### 4.5 这份实现还不完美的地方

1. `io.ReadAll(body)` 的错误被忽略了。
2. `json.Compact(...)` 的错误也被忽略了。
3. `start.(time.Time)` 如果没取到值会 panic。
4. 没有做敏感字段脱敏，密码、token 一类信息如果进 body，会直接被打日志。
5. 不是所有请求体都一定适合完整打印。

也就是说，这一版是“把观察能力先建立起来”，不是生产级日志治理终稿。

## 5. 配套文件：`internal/common/response.go`

### 5.1 这个文件自己的原始 diff

```diff
diff --git a/internal/common/response.go b/internal/common/response.go
index 7be1fe7..4abf56e 100644
--- a/internal/common/response.go
+++ b/internal/common/response.go
@@ -1,6 +1,7 @@
 package common
 
 import (
+	"encoding/json"
 	"net/http"
@@
 func (base *BaseResponse) success(c *gin.Context, data interface{}) {
-	c.JSON(http.StatusOK, response{
+	r := response{
 		Errno:   0,
 		Message: "success",
 		Data:    data,
 		TraceID: tracing.TraceID(c.Request.Context()),
-	})
+	}
+	resp, _ := json.Marshal(r)
+	c.Set("response", string(resp))
+	c.JSON(http.StatusOK, r)
 }
@@
 func (base *BaseResponse) error(c *gin.Context, err error) {
-	c.JSON(http.StatusOK, response{
+	r := response{
 		Errno:   2,
 		Message: err.Error(),
 		Data:    nil,
 		TraceID: tracing.TraceID(c.Request.Context()),
-	})
+	}
+	resp, _ := json.Marshal(r)
+	c.Set("response", string(resp))
+	c.JSON(http.StatusOK, r)
 }
```

### 5.2 旧代码和新代码的区别

旧代码：

- 直接 `c.JSON(...)` 输出响应。
- 目标只有一个：把 HTTP 响应发给客户端。

新代码：

- 先构造 `r := response{...}`。
- 再 `json.Marshal(r)` 得到 JSON 字符串。
- 把这个字符串存进 `gin.Context`。
- 最后再 `c.JSON(...)` 返回给客户端。

你可以把这个变化理解成：

`响应现在不只是“发出去”，还要“留痕”给日志中间件用。`

### 5.3 带中文注释的关键代码

```go
package common

import (
    // 这里新增 encoding/json，不是为了客户端响应本身，
    // 而是为了先把响应对象手动序列化成字符串，方便塞进 gin.Context。
    "encoding/json"
    "net/http"

    "github.com/ghost-yu/go_shop_second/common/tracing"
    "github.com/gin-gonic/gin"
)

func (base *BaseResponse) success(c *gin.Context, data interface{}) {
    r := response{
        Errno:   0,
        Message: "success",
        Data:    data,
        TraceID: tracing.TraceID(c.Request.Context()),
    }

    // json.Marshal 把 Go 值编码成 JSON 字节切片。
    // 输入：结构体 r
    // 输出：[]byte + error
    // 这里不是为了替代 c.JSON，而是为了提前拿到“最终响应长什么样”的字符串版本。
    resp, _ := json.Marshal(r)

    // 把序列化后的响应挂到 gin.Context，后面的 requestOut 日志中间件再取出来。
    c.Set("response", string(resp))

    // c.JSON 才是真正给客户端写 HTTP 响应的地方。
    // StatusOK 这里仍然保持 200，业务状态放在 errno 里。
    c.JSON(http.StatusOK, r)
}

func (base *BaseResponse) error(c *gin.Context, err error) {
    r := response{
        Errno:   2,
        Message: err.Error(),
        Data:    nil,
        TraceID: tracing.TraceID(c.Request.Context()),
    }

    resp, _ := json.Marshal(r)
    c.Set("response", string(resp))
    c.JSON(http.StatusOK, r)
}
```

### 5.4 库函数单独解释

#### `json.Marshal()` 是什么

来源：标准库 `encoding/json`

作用：

- 把 Go 的结构体、map、slice 等值编码成 JSON。

这里为什么要它：

- 因为 `requestOut` 想直接记录最终响应文本。
- `c.JSON` 会直接写到 HTTP 输出流，但不方便再拿回来复用。
- 所以这里先自己 `Marshal` 一次，得到字符串副本。

常见坑：

- 某些字段如果不可序列化会报错。
- 这里课程代码把 error 忽略了，生产里通常应该处理。

#### `c.JSON()` 是什么

来源：`gin.Context`

作用：

- 设置响应头、状态码，并把对象编码成 JSON 写回客户端。

这里为什么还要保留它：

- `json.Marshal` 只是为了“留日志副本”。
- 真正给浏览器/前端回包还是要靠 `c.JSON`。

### 5.5 为什么不直接在中间件里拦截 ResponseWriter

更成熟的做法通常是：

- 自定义一个包装版 `ResponseWriter`
- 拦截写出的 body
- 在中间件层统一采集响应

但这节课没这么做，而是选择：

- 既然项目已经有统一响应出口 `BaseResponse`
- 那就在这里顺手把响应保存一份

这是典型的课程风格：

`先用最少改动打通链路，再考虑更通用的抽象。`

## 6. `internal/common/server/http.go`：新中间件是怎么挂进去的

### 6.1 这个文件自己的原始 diff

```diff
diff --git a/internal/common/server/http.go b/internal/common/server/http.go
index 77af698..840c0a5 100644
--- a/internal/common/server/http.go
+++ b/internal/common/server/http.go
@@ -24,5 +24,6 @@ func setMiddlewares(r *gin.Engine) {
 	r.Use(middleware.StructuredLog(logrus.NewEntry(logrus.StandardLogger())))
 	r.Use(gin.Recovery())
+	r.Use(middleware.RequestLog(logrus.NewEntry(logrus.StandardLogger())))
 	r.Use(otelgin.Middleware("default_server"))
 }
```

### 6.2 带中文注释的代码

```go
func setMiddlewares(r *gin.Engine) {
    // 旧的结构化日志中间件还保留着，但这节里它已经逐渐不再承担主角角色。
    r.Use(middleware.StructuredLog(logrus.NewEntry(logrus.StandardLogger())))

    // gin.Recovery() 是 Gin 官方提供的恢复中间件。
    // 如果 handler panic，它会拦住 panic，避免整个进程直接崩掉，
    // 并返回 500。
    r.Use(gin.Recovery())

    // 这里是本节真正新增的中间件。
    // 从此以后，每个 HTTP 请求都会走 request_in / request_out 两条日志。
    r.Use(middleware.RequestLog(logrus.NewEntry(logrus.StandardLogger())))

    // OpenTelemetry 的 Gin 中间件，负责 trace 采集。
    r.Use(otelgin.Middleware("default_server"))
}
```

### 6.3 库函数解释

#### `r.Use()` 是什么

来源：`gin.Engine`

作用：

- 往 Gin 路由引擎里注册中间件。

你可以理解成：

`以后所有请求都要先经过这些函数。`

#### `gin.Recovery()` 是什么

来源：Gin 官方中间件。

作用：

- 捕获 panic，防止整个 HTTP 服务崩掉。

为什么重要：

- Go 里如果 panic 没被 recover，当前 goroutine 会炸掉。
- 对 HTTP 服务来说，没有 recovery 会非常危险。

## 7. `internal/common/decorator/logging.go`：command 不再伪装成 query

### 7.1 这个文件自己的原始 diff

```diff
diff --git a/internal/common/decorator/logging.go b/internal/common/decorator/logging.go
index ac4eed0..c25240d 100644
--- a/internal/common/decorator/logging.go
+++ b/internal/common/decorator/logging.go
@@ -27,3 +27,26 @@ func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result
 	return q.base.Handle(ctx, cmd)
 }
+
+type commandLoggingDecorator[C, R any] struct {
+	logger *logrus.Entry
+	base   CommandHandler[C, R]
+}
+
+func (q commandLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
+	body, _ := json.Marshal(cmd)
+	logger := q.logger.WithFields(logrus.Fields{
+		"command":      generateActionName(cmd),
+		"command_body": string(body),
+	})
+	logger.Debug("Executing command")
+	defer func() {
+		if err == nil {
+			logger.Info("Command execute successfully")
+		} else {
+			logger.Error("Failed to execute command", err)
+		}
+	}()
+	return q.base.Handle(ctx, cmd)
+}
```

### 7.2 这段改动的真实意义

以前项目里虽然有 `ApplyCommandDecorators(...)`，但内部实际复用了 query 的日志装饰器。这样会导致：

- 你明明在执行 command
- 日志字段却写成 `query`
- body 字段名也叫 `query_body`

从功能上它也许还能跑，但从语义上已经错了。

日志语义错，后面排查时就会误导人。

### 7.3 带中文注释的代码

```go
type commandLoggingDecorator[C, R any] struct {
    // logger 是 logrus 的日志入口。
    logger *logrus.Entry

    // base 是被装饰的真正 command handler。
    // 装饰器模式的意思是：我不改真正业务逻辑，
    // 我在外面再包一层，额外加日志、指标、权限等横切逻辑。
    base CommandHandler[C, R]
}

func (q commandLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
    // 这里把 command 结构体转成 JSON 字符串，方便日志输出。
    // 如果你直接 %v 打印，复杂结构体往往不直观。
    body, _ := json.Marshal(cmd)

    logger := q.logger.WithFields(logrus.Fields{
        // generateActionName(cmd) 会从类型名里提取动作名。
        // 比如 CreateOrderCommand -> 记录成 createordercommand 对应动作名的一部分。
        "command": generateActionName(cmd),

        // command_body 明确表示“这是命令入参”，不再混成 query_body。
        "command_body": string(body),
    })

    logger.Debug("Executing command")

    // defer 的目的：不管 Handle 成功还是失败，最后都补一条收尾日志。
    defer func() {
        if err == nil {
            logger.Info("Command execute successfully")
        } else {
            // 这里用 Error(msg, err) 在 logrus 里不是最优雅写法，
            // 更常见是 WithError(err).Error(msg)。
            // 但课程里先把语义修正过来，已经比之前强很多。
            logger.Error("Failed to execute command", err)
        }
    }()

    return q.base.Handle(ctx, cmd)
}
```

### 7.4 这里涉及的库和技巧

#### `json.Marshal(cmd)` 为什么比 `%v` 更适合日志

- command 往往是结构体。
- `%v` 打印出来不一定稳定、可读、可搜索。
- JSON 形式更容易进日志平台检索。

#### `defer` 在这里的价值

这里不是为了“代码好看”，而是为了保证：

- 无论 handler 最后成功还是失败
- 收尾日志都能打出来

这是一种非常常见的 Go 写法。

## 8. `internal/common/decorator/metrics.go`：指标语义也跟着修正

### 8.1 这个文件自己的原始 diff

```diff
diff --git a/internal/common/decorator/metrics.go b/internal/common/decorator/metrics.go
index 6453218..97239c4 100644
--- a/internal/common/decorator/metrics.go
+++ b/internal/common/decorator/metrics.go
@@ -24,3 +24,24 @@ func (q queryMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result
 	return q.base.Handle(ctx, cmd)
 }
+
+type commandMetricsDecorator[C, R any] struct {
+	base   CommandHandler[C, R]
+	client MetricsClient
+}
+
+func (q commandMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
+	start := time.Now()
+	actionName := strings.ToLower(generateActionName(cmd))
+	defer func() {
+		end := time.Since(start)
+		q.client.Inc(fmt.Sprintf("command.%s.duration", actionName), int(end.Seconds()))
+		if err == nil {
+			q.client.Inc(fmt.Sprintf("command.%s.success", actionName), 1)
+		} else {
+			q.client.Inc(fmt.Sprintf("command.%s.failure", actionName), 1)
+		}
+	}()
+	return q.base.Handle(ctx, cmd)
+}
```

### 8.2 带中文注释的代码

```go
type commandMetricsDecorator[C, R any] struct {
    base   CommandHandler[C, R]
    client MetricsClient
}

func (q commandMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
    start := time.Now()

    // generateActionName 提取动作名，再统一转小写，便于拼指标 key。
    actionName := strings.ToLower(generateActionName(cmd))

    defer func() {
        end := time.Since(start)

        // 这里的 key 终于变成 command.xxx，而不是沿用 query.xxx。
        q.client.Inc(fmt.Sprintf("command.%s.duration", actionName), int(end.Seconds()))

        if err == nil {
            q.client.Inc(fmt.Sprintf("command.%s.success", actionName), 1)
        } else {
            q.client.Inc(fmt.Sprintf("command.%s.failure", actionName), 1)
        }
    }()

    return q.base.Handle(ctx, cmd)
}
```

### 8.3 要补给初学者的基础设施解释

#### 指标和日志的区别

- 日志：更适合看“某次具体请求发生了什么”
- 指标：更适合看“整体趋势怎么样”

例如：

- 日志适合查某个订单为什么失败
- 指标适合看过去 10 分钟 `CreateOrder` 失败率是否上升

#### `MetricsClient.Inc()` 在这里是什么思路

这不是标准库，而是项目自己抽出来的指标接口。含义很简单：

- 某个 key 增加一个数值

这样做的好处是：

- 当前课程里可以先用简单实现
- 以后接 Prometheus、StatsD、Datadog 时，上层代码不用大改

### 8.4 这里还有什么不完美

1. `duration` 这里用的是 `int(end.Seconds())`，精度比较粗，秒以下会被截断。
2. query 那边还沿用了 `querys` 这个拼写，不够规范。
3. 只是把语义分开了，还没谈更完整的指标标签体系。

## 9. `internal/common/decorator/command.go`：真正接入新的 command 装饰器

### 9.1 这个文件自己的原始 diff

```diff
diff --git a/internal/common/decorator/command.go b/internal/common/decorator/command.go
index 9e2adbd..4ea2abc 100644
--- a/internal/common/decorator/command.go
+++ b/internal/common/decorator/command.go
@@ -11,10 +11,10 @@ type CommandHandler[C, R any] interface {
 }
 
 func ApplyCommandDecorators[C, R any](handler CommandHandler[C, R], logger *logrus.Entry, metricsClient MetricsClient) CommandHandler[C, R] {
-	return queryLoggingDecorator[C, R]{
+	return commandLoggingDecorator[C, R]{
 		logger: logger,
-		base: queryMetricsDecorator[C, R]{
+		base: commandMetricsDecorator[C, R]{
 			base:   handler,
 			client: metricsClient,
 		},
```

### 9.2 这段改动为什么重要

前面的 `logging.go` 和 `metrics.go` 只是新增了两个“新零件”。

`command.go` 才是把这些零件真的接上线的地方。

也就是说：

- 没这一步，新 decorator 只是定义了，但没人用
- 加了这一步，command 链路才真正换成新的语义

### 9.3 带中文注释的代码

```go
func ApplyCommandDecorators[C, R any](handler CommandHandler[C, R], logger *logrus.Entry, metricsClient MetricsClient) CommandHandler[C, R] {
    return commandLoggingDecorator[C, R]{
        logger: logger,
        base: commandMetricsDecorator[C, R]{
            base:   handler,
            client: metricsClient,
        },
    }
}
```

这其实就是一个“套娃”结构：

- 最里层：真正业务 handler
- 外面一层：metrics decorator
- 再外面一层：logging decorator

最终执行时，相当于：

`logging -> metrics -> real handler`

这也是装饰器模式在 Go 项目里很常见的写法。

## 10. `internal/common/middleware/logger.go`：旧的 StructuredLog 退居次要位置

### 10.1 这个文件自己的原始 diff

```diff
diff --git a/internal/common/middleware/logger.go b/internal/common/middleware/logger.go
index 44df58c..8aab01c 100644
--- a/internal/common/middleware/logger.go
+++ b/internal/common/middleware/logger.go
@@ -1,8 +1,10 @@
 package middleware
 
 import (
+	// "time"
 	"github.com/gin-gonic/gin"
 	"github.com/sirupsen/logrus"
 )
@@
 	return func(c *gin.Context) {
+		//t := time.Now()
 		c.Next()
+		//elapsed := time.Since(t)
 		//l.WithFields(logrus.Fields{
 		//	"time_elapsed_ms": elapsed.Milliseconds(),
```

### 10.2 这说明了什么

这不是“乱删代码”，而是一个明显的演进信号：

- 原来的 `StructuredLog` 只打一条比较粗的请求日志
- 新的 `RequestLog` 已经开始承担更完整的入站/出站日志职责
- 所以旧逻辑开始淡出

它还没被正式删掉，说明课程还处在迁移期。

这类代码你以后会经常看到：

`新实现已经上线，旧实现暂时保留，等确认稳定后再彻底清理。`

## 11. 第三方库和基础设施，这组必须理解的点

### 11.1 Gin 在这里到底扮演什么角色

Gin 是 HTTP Web 框架。

这组里你至少要记住它的几个角色：

- `gin.Engine`：整个 HTTP 路由引擎
- `gin.Context`：单次请求的上下文
- `gin.HandlerFunc`：handler / middleware 的函数类型
- `r.Use(...)`：注册中间件
- `c.Next()`：继续执行后续链路
- `c.Set/Get()`：在同一请求内传数据
- `c.JSON()`：返回 JSON 响应

也就是说，这组课本质上是在“更深入地使用 Gin 的中间件模型”。

### 11.2 Logrus 在这里到底扮演什么角色

Logrus 是结构化日志库。

这里不是简单 `println`，而是要把日志打成：

- 消息文本
- 一组可检索字段
- 可附带 context

这样后面如果接日志平台，才能真正按字段搜索和聚合。

### 11.3 为什么这组和 tracing 不是一回事

初学者容易把日志和 trace 混在一起。

你可以这样分：

- trace：看一条请求跨服务怎么流动
- 日志：看某个节点内部到底发生了什么

这一组主要补的是日志，不是 trace。本质上它是在让单服务内部“看得更清楚”。

## 12. 这一组还不完美的地方

1. body 和 response 都可能包含敏感信息，后面应该考虑脱敏。
2. `io.ReadAll` 直接全读，对超大请求不够稳。
3. 错误处理大量被忽略，是课程简化写法。
4. response 采集依赖统一 `BaseResponse` 出口，通用性还不够强。
5. 日志字段命名、错误记录方式还可以继续规范化。

## 13. 这组最该带走的知识点

1. Gin 中间件不是魔法，本质上就是在 handler 前后包逻辑。
2. `Request.Body` 是一次性流，读完后如果还想给下游用，必须重新塞回去。
3. `gin.Context` 不只是拿参数，它还能在请求链里传递临时数据。
4. `json.Marshal` 和 `c.JSON` 不是重复，它们这次分别服务于“日志留痕”和“HTTP 回包”。
5. command 和 query 的日志/指标必须语义一致，否则会直接误导排查。
6. 可观测性代码不是装饰品，而是系统复杂后必须补的基础能力。

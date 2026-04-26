# Lesson Pair Diff Report

- FromBranch: lesson10
- ToBranch: lesson11

## Short Summary

~~~text
 28 files changed, 853 insertions(+), 4 deletions(-)
~~~

## File Stats

~~~text
 ...273\243\347\240\201\350\247\243\350\257\273.md" | 716 +++++++++++++++++++++
 internal/common/config/viper.go                    |  19 +
 internal/common/decorator/command.go               |   3 +
 internal/common/decorator/logging.go               |   4 +
 internal/common/decorator/metrics.go               |   4 +
 internal/common/decorator/query.go                 |   3 +
 internal/common/metrics/todo_metrics.go            |   2 +
 internal/common/server/gprc.go                     |  10 +
 internal/common/server/http.go                     |   5 +
 internal/order/adapters/grpc/stock_grpc.go         |   4 +
 internal/order/adapters/order_inmem_repository.go  |   7 +
 internal/order/app/app.go                          |   3 +
 internal/order/app/command/create_order.go         |   7 +
 internal/order/app/command/update_order.go         |   3 +
 internal/order/app/query/get_customer_order.go     |   6 +
 internal/order/app/query/service.go                |   2 +
 internal/order/domain/order/order.go               |  14 +-
 internal/order/domain/order/repository.go          |   2 +
 internal/order/http.go                             |   5 +
 internal/order/main.go                             |   6 +
 internal/order/ports/grpc.go                       |   5 +
 internal/order/service/application.go              |   4 +
 internal/stock/adapters/stock_inmem_repository.go  |   5 +
 internal/stock/app/app.go                          |   4 +
 internal/stock/domain/stock/repository.go          |   3 +
 internal/stock/main.go                             |   5 +
 internal/stock/ports/grpc.go                       |   4 +
 internal/stock/service/application.go              |   2 +
 28 files changed, 853 insertions(+), 4 deletions(-)
~~~

## Commit Comparison

~~~text
> 3163c56 docs: add comments to viper config loader
> 853d6f0 docs: add learning comments for backend internship study
> 2d6db2e docs: gprc.go & http.go 补全完整代码+逐段详细注释
> e2e4b04 docs: 扩展 Lesson5-11 为文件级复现手册
> 6d9bcf7 docs: 添加 Lesson5-11 关键代码解读（Go 小白版）
~~~

## Changed Files

~~~text
"docs/lesson-notes/lesson5-lesson11-\345\205\263\351\224\256\344\273\243\347\240\201\350\247\243\350\257\273.md"
internal/common/config/viper.go
internal/common/decorator/command.go
internal/common/decorator/logging.go
internal/common/decorator/metrics.go
internal/common/decorator/query.go
internal/common/metrics/todo_metrics.go
internal/common/server/gprc.go
internal/common/server/http.go
internal/order/adapters/grpc/stock_grpc.go
internal/order/adapters/order_inmem_repository.go
internal/order/app/app.go
internal/order/app/command/create_order.go
internal/order/app/command/update_order.go
internal/order/app/query/get_customer_order.go
internal/order/app/query/service.go
internal/order/domain/order/order.go
internal/order/domain/order/repository.go
internal/order/http.go
internal/order/main.go
internal/order/ports/grpc.go
internal/order/service/application.go
internal/stock/adapters/stock_inmem_repository.go
internal/stock/app/app.go
internal/stock/domain/stock/repository.go
internal/stock/main.go
internal/stock/ports/grpc.go
internal/stock/service/application.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
"docs/lesson-notes/lesson5-lesson11-\345\205\263\351\224\256\344\273\243\347\240\201\350\247\243\350\257\273.md"
internal/common/config/viper.go
internal/common/decorator/command.go
internal/common/decorator/logging.go
internal/common/decorator/metrics.go
internal/common/decorator/query.go
internal/common/metrics/todo_metrics.go
internal/common/server/gprc.go
internal/common/server/http.go
internal/order/adapters/grpc/stock_grpc.go
internal/order/adapters/order_inmem_repository.go
internal/order/app/app.go
internal/order/app/command/create_order.go
internal/order/app/command/update_order.go
internal/order/app/query/get_customer_order.go
internal/order/app/query/service.go
internal/order/domain/order/order.go
internal/order/domain/order/repository.go
internal/order/http.go
internal/order/main.go
internal/order/ports/grpc.go
internal/order/service/application.go
internal/stock/adapters/stock_inmem_repository.go
internal/stock/app/app.go
internal/stock/domain/stock/repository.go
internal/stock/main.go
internal/stock/ports/grpc.go
internal/stock/service/application.go
~~~

## Full Diff

~~~diff
diff --git "a/docs/lesson-notes/lesson5-lesson11-\345\205\263\351\224\256\344\273\243\347\240\201\350\247\243\350\257\273.md" "b/docs/lesson-notes/lesson5-lesson11-\345\205\263\351\224\256\344\273\243\347\240\201\350\247\243\350\257\273.md"
new file mode 100644
index 0000000..fc31db2
--- /dev/null
+++ "b/docs/lesson-notes/lesson5-lesson11-\345\205\263\351\224\256\344\273\243\347\240\201\350\247\243\350\257\273.md"
@@ -0,0 +1,716 @@
+# gorder-v2 Lesson5-Lesson11 复现级代码解读（Go 小白版）
+
+这版文档目标是“可复现”，不是只做概念导读。
+你可以按本文件逐个文件阅读并运行，每个文件都给出：
+1. 文件职责
+2. 关键代码必须理解点
+3. 上下游依赖（谁调用它，它调用谁）
+4. 复现检查动作（你可以自己执行验证）
+
+## 0. 先说明边界
+
+Lesson5-Lesson11 里涉及的文件很多，其中有三类：
+1. 业务代码文件：需要详细读懂。
+2. 生成文件（*.pb.go, *_grpc.pb.go）：需要理解结构和用途，不建议逐行硬背。
+3. 依赖锁文件（go.sum）：不需要逐条解释，每行本质是模块校验和。
+
+因此本文对所有文件都做说明，但会按“学习价值”给不同颗粒度，确保你能复现整个阶段。
+
+## 1. Lesson5 到 Lesson11 的总调用图
+
+1. 配置加载：internal/common/config/global.yaml -> viper
+2. 进程入口：internal/order/main.go、internal/stock/main.go
+3. 基础服务启动：internal/common/server/http.go、internal/common/server/gprc.go
+4. 端口层：internal/order/http.go、internal/order/ports/grpc.go、internal/stock/ports/grpc.go
+5. 应用层：internal/order/app/...、internal/stock/app/...
+6. 领域接口：internal/order/domain/order/repository.go、internal/stock/domain/stock/repository.go
+7. 适配器实现：in-memory repo + grpc 适配器
+8. 协议层：api/*.proto + internal/common/genproto/*.pb.go
+
+---
+
+## 2. 按文件详细解释（Lesson5-Lesson11）
+
+### A. 协议与配置文件
+
+#### api/stockpb/stock.proto
+文件职责：定义库存服务 RPC 契约。
+关键理解：
+1. service StockService 里定义了两个 RPC：GetItems、CheckIfItemsInStock。
+2. 请求和响应消息是跨服务通信唯一标准，服务端和客户端必须一致。
+3. 引用 orderpb/order.proto 是为了复用订单侧 Item 结构，避免重复建模。
+依赖关系：
+1. 被 internal/common/genproto/stockpb/stock.pb.go 和 stock_grpc.pb.go 生成工具使用。
+2. 被 order 侧 grpc 客户端调用路径间接使用。
+复现检查：
+1. 运行 proto 生成脚本后，确认生成代码存在对应消息与服务接口。
+
+#### api/orderpb/order.proto
+文件职责：定义订单服务对外协议和公共消息（如 Item）。
+关键理解：
+1. Item 和 ItemWithQuantity 是库存校验与下单共用的数据结构。
+2. OrderService 方法签名决定了 order 端口层函数形态。
+依赖关系：
+1. 被 OpenAPI/GRPC 生成代码使用。
+2. 被 stock.proto import。
+复现检查：
+1. 改字段名后重新生成，会导致调用方编译错误，证明契约约束生效。
+
+#### internal/common/config/global.yaml
+文件职责：集中配置各服务地址与运行模式。
+关键理解：
+1. order 与 stock 的 grpc/http 地址是独立的。
+2. fallback-grpc-addr 是防御性配置，避免空地址直接崩溃。
+依赖关系：
+1. 被 config/viper 初始化后读取。
+2. 被 server/http.go 和 server/gprc.go 间接消费。
+复现检查：
+1. 改 grpc-addr 后重启服务，端口变化应立即生效。
+
+#### .gitignore
+文件职责：控制 Git 不追踪哪些本地文件。
+关键理解：
+1. 用于忽略临时构建产物和本地工具缓存。
+2. 防止开发机噪音提交进仓库。
+依赖关系：
+1. Git 客户端直接读取。
+复现检查：
+1. 新建被忽略文件后 git status 不应出现。
+
+### B. 生成代码文件
+
+#### internal/common/genproto/orderpb/order.pb.go
+文件职责：order.proto 的消息结构与序列化代码。
+关键理解：
+1. 这类文件由 protoc 自动生成，不手改。
+2. 包含 message struct、反射元数据、编解码实现。
+依赖关系：
+1. 被业务代码 import，作为请求/响应 DTO。
+复现检查：
+1. 删除后项目会出现大量类型缺失错误。
+
+#### internal/common/genproto/stockpb/stock.pb.go
+文件职责：stock.proto 消息生成代码。
+关键理解：
+1. 与 order.pb.go 同类，承载跨服务消息类型。
+依赖关系：
+1. 被 stock 服务端与 order 客户端共同依赖。
+复现检查：
+1. 运行 stock grpc 相关代码需依赖此文件类型。
+
+#### internal/common/genproto/stockpb/stock_grpc.pb.go
+文件职责：stock.proto 的 gRPC 客户端/服务端接口生成代码。
+关键理解：
+1. 生成 StockServiceClient 和 StockServiceServer 接口。
+2. order 侧 grpc 适配器通过 client 发请求。
+3. stock 侧端口实现通过 server 接口对外暴露。
+依赖关系：
+1. 被 internal/order/adapters/grpc/stock_grpc.go 使用。
+2. 被 internal/stock/ports/grpc.go 实现。
+复现检查：
+1. 在 main.go 注册 stockpb.RegisterStockServiceServer 即依赖此文件。
+
+### C. 基础设施启动文件
+
+#### internal/common/server/gprc.go
+文件职责：统一 gRPC server 启动流程，所有 gRPC 服务都通过这里启动，业务服务只需提供「注册业务」这一个回调。
+
+完整代码+逐段注释：
+```go
+package server
+
+import (
+    "net"
+    grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
+    grpc_tags   "github.com/grpc-ecosystem/go-grpc-middleware/tags"
+    "github.com/sirupsen/logrus"
+    "github.com/spf13/viper"
+    "google.golang.org/grpc"
+)
+
+// init 是 Go 的特殊函数，包被 import 时自动执行，不需要手动调用。
+// 这里把 gRPC 框架内部日志替换成项目统一使用的 logrus，同时把级别设为 Warn，
+// 减少 gRPC 框架内部的噪音日志，保持日志格式一致，方便集中收集和排查问题。
+func init() {
+    logger := logrus.New()
+    logger.SetLevel(logrus.WarnLevel)                      // 只打 warn 及以上
+    grpc_logrus.ReplaceGrpcLogger(logrus.NewEntry(logger)) // 替换 gRPC 内置 logger
+}
+
+// RunGRPCServer 是对外的唯一入口。
+// 参数 serviceName：对应 global.yaml 里的 key，比如 "order" 或 "stock"。
+// 参数 registerServer：业务方传入的回调，用于把自己的 handler 注册进 grpc.Server。
+//
+// 调用方写法（order/main.go）示例：
+//   server.RunGRPCServer("order", func(s *grpc.Server) {
+//       orderpb.RegisterOrderServiceServer(s, ports.NewGRPCServer(app))
+//   })
+func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server)) {
+    // viper.Sub("order").GetString("grpc-addr") 相当于读 global.yaml 里 order.grpc-addr 的值。
+    addr := viper.Sub(serviceName).GetString("grpc-addr")
+    if addr == "" {
+        // 地址缺失时用全局兜底地址，防止配置漏填导致进程直接崩溃。
+        addr = viper.GetString("fallback-grpc-addr")
+    }
+    RunGRPCServerOnAddr(addr, registerServer)
+}
+
+// RunGRPCServerOnAddr 真正创建并运行 gRPC server。
+// 独立函数的原因：测试时可以直接传 "127.0.0.1:0"（随机端口），不依赖配置文件。
+func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
+    logrusEntry := logrus.NewEntry(logrus.StandardLogger())
+
+    // grpc.NewServer(...) 创建 gRPC 服务器，括号里配置拦截器。
+    // 拦截器 = HTTP 框架里的中间件，每次 RPC 调用进来和出去时都会经过它们，
+    // 适合放日志、追踪、鉴权、panic 恢复等和业务无关的横切逻辑。
+    grpcServer := grpc.NewServer(
+
+        // ChainUnaryInterceptor：配置「一元 RPC」拦截器链（普通请求-响应模式，最常用）。
+        // 执行顺序：请求进来时从左到右，响应返回时从右到左。
+        grpc.ChainUnaryInterceptor(
+            // ① 从 proto 生成的请求结构里自动提取字段（如 customer_id）打到日志 Tag 里。
+            //   后续该请求的所有日志都自动携带这些字段，排查问题时不用手动传。
+            grpc_tags.UnaryServerInterceptor(
+                grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor),
+            ),
+            // ② 自动记录每次 RPC 的方法名、耗时、gRPC 状态码到 logrus。
+            grpc_logrus.UnaryServerInterceptor(logrusEntry),
+
+            // 下面是后续 lesson 会逐步打开的能力，当前注释保留位置：
+            // otelgrpc.UnaryServerInterceptor(),       // 分布式链路追踪（Lesson25+）
+            // srvMetrics.UnaryServerInterceptor(...),  // Prometheus 指标采集（Lesson43+）
+            // selector.UnaryServerInterceptor(...),    // 选择性鉴权（某些路由才需要鉴权）
+            // recovery.UnaryServerInterceptor(...),    // panic 自动恢复（防止一个请求崩掉整个进程）
+        ),
+
+        // ChainStreamInterceptor：配置「流式 RPC」拦截器链（双向流/服务端推送场景）。
+        // 当前项目主要用一元 RPC，这里配置和一元保持同步即可。
+        grpc.ChainStreamInterceptor(
+            grpc_tags.StreamServerInterceptor(
+                grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor),
+            ),
+            grpc_logrus.StreamServerInterceptor(logrusEntry),
+            // 同上，预留位置暂时注释
+        ),
+    )
+
+    // 执行业务注册回调，把 OrderService/StockService 的具体实现挂到 grpcServer 上。
+    // 这行之前 grpcServer 是空的（不知道有哪些方法），执行后才知道怎么处理请求。
+    registerServer(grpcServer)
+
+    // net.Listen 在 TCP 层监听指定地址，此时还没开始接受 gRPC 请求。
+    listen, err := net.Listen("tcp", addr)
+    if err != nil {
+        // 端口被占用、权限不足等情况，直接 panic 让进程退出，而不是静默失败。
+        logrus.Panic(err)
+    }
+    logrus.Infof("Starting gRPC server, Listening: %s", addr)
+
+    // grpcServer.Serve 是阻塞调用，进入后一直处理 gRPC 请求，直到 server Stop 或出错。
+    if err := grpcServer.Serve(listen); err != nil {
+        logrus.Panic(err)
+    }
+}
+```
+
+核心设计要点：
+1. **两函数分离**：RunGRPCServer 负责读配置（依赖 viper），RunGRPCServerOnAddr 只依赖
+   地址字符串。单元测试时可以绕过配置直接调用 RunGRPCServerOnAddr，不需要 global.yaml。
+2. **registerServer 回调模式**：框架不知道业务注册了什么，业务不知道框架怎么启动，
+   通过这一个函数参数完成解耦。这是 Go 里高频出现的依赖注入写法，全项目会反复看到。
+3. **拦截器链是功能扩展点**：目前只有日志+标签两层，后续每个注释行都是可独立打开的能力，
+   打开时完全不需要改任何业务代码，只需取消注释并引入对应包。
+
+依赖关系：
+- 上游调用方：order/main.go 和 stock/main.go 传入 serviceName 和 registerServer 调用本函数。
+- 下游依赖：registerServer 回调由业务服务提供（例如 orderpb.RegisterOrderServiceServer）。
+
+复现检查：
+1. 启动后终端出现 `Starting gRPC server, Listening: 127.0.0.1:5002`。
+2. 改 global.yaml 的 grpc-addr 后重启，监听端口随之变化。
+3. 用重复端口启动第二个实例，net.Listen 会报 bind: address already in use 并 panic。
+
+---
+
+#### internal/common/server/http.go
+文件职责：统一 HTTP server 启动流程，所有 HTTP 服务都通过这里启动，业务方只需传入路由注册回调。
+
+完整代码+逐段注释：
+```go
+package server
+
+import (
+    "github.com/gin-gonic/gin"
+    "github.com/spf13/viper"
+)
+
+// RunHTTPServer 按服务名从 viper 读 http-addr，然后启动 HTTP 服务。
+// 参数 wrapper：业务方提供的路由注册函数，接收 *gin.Engine 并往里注册路由和 handler。
+//
+// 调用方写法（order/main.go）示例：
+//   server.RunHTTPServer("order", func(router *gin.Engine) {
+//       ports.RegisterHandlersWithOptions(router, HTTPServer{app: application}, opts)
+//   })
+func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
+    // viper.Sub("order").GetString("http-addr") 读取 global.yaml 里 order.http-addr 的值。
+    addr := viper.Sub(serviceName).GetString("http-addr")
+    if addr == "" {
+        // TODO：后续应加 warning 日志提醒配置缺失。
+        // 注意：addr 为空时 gin.Run("") 会监听 :80，非 root 权限下会权限错误。
+    }
+    RunHTTPServerOnAddr(addr, wrapper)
+}
+
+// RunHTTPServerOnAddr 真正创建 gin 引擎并运行。
+// 独立出来的原因：测试时可以直接传 ":0"（随机端口），不依赖 viper 配置。
+func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
+    // gin.New() 创建一个没有任何默认中间件的干净 gin 引擎。
+    // 区别于 gin.Default()：Default() 会自动加 Logger 和 Recovery 中间件，
+    // 这里用 New() 是为了完全掌控中间件，后续 lesson 会按需手动添加。
+    apiRouter := gin.New()
+
+    // 执行业务方传入的路由注册函数，把路由和 handler 挂到 apiRouter 上。
+    // 执行完后 apiRouter 才知道"POST /api/customer/.../orders"对应哪个 handler。
+    wrapper(apiRouter)
+
+    // ⚠️ 注意：下面这行实际上是无效代码（遗留 bug）。
+    // apiRouter.Group("/api") 返回一个新的 RouterGroup 对象，但没有赋值给变量，
+    // 返回值被直接丢弃，对已有路由没有任何影响，相当于执行了一个空操作。
+    // 真正的 /api 前缀是在 wrapper 内部通过 RegisterHandlersWithOptions 的 BaseURL:"/api" 设置的。
+    // 理解这一点可以防止你自己写出同类 bug。
+    apiRouter.Group("/api")
+
+    // apiRouter.Run 是阻塞调用，开始监听并处理 HTTP 请求。
+    // 出错（端口冲突、权限不足）时直接 panic，不要静默失败。
+    if err := apiRouter.Run(addr); err != nil {
+        panic(err)
+    }
+}
+```
+
+核心设计要点：
+1. **gin.New() 而非 gin.Default()**：主动掌控中间件，不被框架默认行为绑架。
+   需要加 cors、recovery、自定义日志时直接往 apiRouter 上加，与 grpc.go 的拦截器链思路一致。
+2. **wrapper 回调模式**：与 gprc.go 的 registerServer 完全相同的解耦思路，
+   http.go 不知道有哪些路由，order/main.go 负责告诉它。
+3. **apiRouter.Group("/api") 是无效代码**：这是代码里的遗留问题，不影响功能运行，
+   因为路由前缀已经由 RegisterHandlersWithOptions 的 BaseURL 参数正确设置了。
+
+依赖关系：
+- 上游调用方：order/main.go 传入 "order" 和路由注册 wrapper 调用本函数。
+- 下游依赖：wrapper 函数由 order/main.go 提供，其中调用 ports.RegisterHandlersWithOptions。
+
+复现检查：
+1. 启动 order 服务后，`curl http://127.0.0.1:8282/api/customer/test/orders` 能到达 handler。
+2. 修改 http-addr 端口后重启，访问旧端口应收到连接拒绝。
+3. 将端口改为已占用端口，gin.Run 返回 error 并触发 panic。
+
+### D. 入口与端口层文件
+
+#### internal/order/main.go
+文件职责：Order 服务入口。
+关键理解：
+1. init 先加载配置，main 再启动服务。
+2. 同时启动 grpc 和 http，说明订单服务双协议暴露。
+3. HTTP 注册走 OpenAPI 生成的 RegisterHandlersWithOptions。
+依赖关系：
+1. 调用 common/server。
+2. 调用 order/ports + order/http。
+3. 调用 order/service.NewApplication（Lesson7 后）。
+复现检查：
+1. 运行后应同时看到 HTTP 与 gRPC 监听日志。
+
+#### internal/order/http.go
+文件职责：Order HTTP 端口适配层。
+关键理解：
+1. 方法名来自 OpenAPI 生成接口，不是手写命名。
+2. 正式业务版本里只做三件事：绑定参数、调用应用层、返回响应。
+依赖关系：
+1. 被 main.go 注册为 HTTP handler。
+2. 调用 app.Commands / app.Queries。
+复现检查：
+1. POST 下单和 GET 查询调用时会进入对应 handler。
+
+#### internal/order/ports/grpc.go
+文件职责：Order gRPC 端口实现。
+关键理解：
+1. 实现 orderpb.OrderServiceServer 接口。
+2. 该阶段多为 TODO，后续 lesson 逐步填充。
+依赖关系：
+1. 被 main.go 的 RegisterOrderServiceServer 注册。
+复现检查：
+1. 若方法未实现，调用对应 RPC 会 panic。
+
+#### internal/stock/main.go
+文件职责：Stock 服务入口。
+关键理解：
+1. 读取 stock.server-to-run 决定跑 grpc 还是 http。
+2. 当前阶段主要跑 grpc。
+依赖关系：
+1. 调用 service.NewApplication。
+2. 调用 common/server.RunGRPCServer。
+复现检查：
+1. 改 server-to-run 值可触发不同分支行为。
+
+#### internal/stock/ports/grpc.go
+文件职责：Stock gRPC 端口实现。
+关键理解：
+1. 需要实现 GetItems 和 CheckIfItemsInStock。
+2. 这层把外部请求转成应用层调用。
+依赖关系：
+1. 被 stock main 注册。
+2. 供 order 的 grpc client 调用。
+复现检查：
+1. 用 grpc 客户端请求库存接口能否返回。
+
+### E. 应用层组装与门面文件
+
+#### internal/stock/app/app.go
+文件职责：Stock 应用层门面（Commands、Queries 聚合）。
+关键理解：
+1. 通过一个 Application struct 对外暴露能力。
+2. 减少上层对内部 handler 细节感知。
+依赖关系：
+1. 被 stock service/application.go 填充。
+2. 被 stock ports 层消费。
+复现检查：
+1. main 中传递 Application 后端口层能拿到功能入口。
+
+#### internal/stock/service/application.go
+文件职责：Stock 的依赖注入与应用组装。
+关键理解：
+1. 这是组合根，负责把 repo、handler 组装成 app.Application。
+2. main 不直接 new 各种依赖，职责更清晰。
+依赖关系：
+1. 被 stock main 调用。
+2. 依赖 stock adapters/domain/app。
+复现检查：
+1. 替换 repo 实现时只改此处注入。
+
+#### internal/order/app/app.go
+文件职责：Order 应用层门面。
+关键理解：
+1. Commands 放写流程（Create/Update）。
+2. Queries 放读流程（GetCustomerOrder）。
+依赖关系：
+1. 被 order/service/application.go 赋值。
+2. 被 HTTP/gRPC 端口调用。
+复现检查：
+1. 端口层只依赖 app.Application 就能调用所有用例。
+
+#### internal/order/service/application.go
+文件职责：Order 依赖注入核心。
+关键理解：
+1. 初始化 orderRepo、logger、metricsClient。
+2. 初始化 stock grpc client 并封装成 StockService。
+3. 把 command/query handler 统一装配到 app.Application。
+依赖关系：
+1. 被 order main 调用。
+2. 依赖 adapters、decorator、query/service 接口。
+复现检查：
+1. 替换 stock 调用实现时，这里改注入即可。
+
+### F. 领域层与仓储实现
+
+#### internal/order/domain/order/order.go
+文件职责：订单领域模型。
+关键理解：
+1. Order 代表核心业务实体。
+2. Items 使用 proto 结构，便于和接口层衔接。
+依赖关系：
+1. 被 command/query/repository 使用。
+复现检查：
+1. 新增字段后需要同步影响创建与查询返回。
+
+#### internal/order/domain/order/repository.go
+文件职责：订单仓储接口定义。
+关键理解：
+1. 定义 Create/Get/Update 三个能力。
+2. NotFoundError 提供业务语义化错误。
+依赖关系：
+1. 被应用层依赖。
+2. 被 adapters 实现。
+复现检查：
+1. 任何 repo 实现都必须满足该接口。
+
+#### internal/order/adapters/order_inmem_repository.go
+文件职责：订单仓储 in-memory 实现。
+关键理解：
+1. RWMutex 保障并发安全。
+2. Create 用时间戳生成 ID 并写入内存切片。
+3. Get 按订单号和客户号匹配。
+4. Update 使用高阶函数 updateFn 控制更新逻辑。
+依赖关系：
+1. 实现 domain.Repository。
+2. 被 order/service/application.go 注入到 handler。
+复现检查：
+1. 连续创建订单，内存 store 长度增加。
+2. 查不存在订单返回 NotFoundError。
+
+#### internal/stock/domain/stock/repository.go
+文件职责：库存仓储接口定义。
+关键理解：
+1. 定义 GetItems 与 CheckIfItemsInStock 等库存查询能力。
+2. 约束 stock 应用层只依赖接口。
+依赖关系：
+1. 被 stock 适配器实现。
+复现检查：
+1. repo 实现变更不影响应用层签名。
+
+#### internal/stock/adapters/stock_inmem_repository.go
+文件职责：库存仓储 in-memory 实现。
+关键理解：
+1. 用内存数据模拟库存系统。
+2. 返回商品详情和库存校验结果。
+依赖关系：
+1. 实现 stock domain repo。
+2. 被 stock service 注入。
+复现检查：
+1. 调整初始库存数据，库存接口结果应变化。
+
+### G. Command / Query 与装饰器
+
+#### internal/common/decorator/query.go
+文件职责：QueryHandler 泛型接口与装饰器装配入口。
+关键理解：
+1. 统一查询处理器形态。
+2. ApplyQueryDecorators 按顺序套日志和指标。
+依赖关系：
+1. 被 query handler 构造函数调用。
+复现检查：
+1. 执行查询时可看到日志与指标上报被触发。
+
+#### internal/common/decorator/command.go
+文件职责：CommandHandler 泛型接口与装饰器装配入口。
+关键理解：
+1. 与 query 对称，作用于写操作。
+2. 保证 command 也拥有统一日志/指标能力。
+依赖关系：
+1. 被 command handler 构造函数调用。
+复现检查：
+1. 下单成功/失败时都有对应统计路径。
+
+#### internal/common/decorator/logging.go
+文件职责：通用日志装饰器。
+关键理解：
+1. 记录 action 名、请求体、执行结果。
+2. 错误场景会打印失败日志，方便排查。
+依赖关系：
+1. 被 query.go / command.go 装配调用。
+复现检查：
+1. 故意触发参数错误，日志中应有失败记录。
+
+#### internal/common/decorator/metrics.go
+文件职责：通用指标装饰器。
+关键理解：
+1. 统计耗时、成功次数、失败次数。
+2. action 名来自请求类型。
+依赖关系：
+1. 依赖 MetricsClient 接口。
+2. 被 query/command 装配。
+复现检查：
+1. 打断点可看到 Inc 被调用三类 key。
+
+#### internal/common/metrics/todo_metrics.go
+文件职责：指标客户端占位实现。
+关键理解：
+1. 当前 Inc 空实现，仅保证接口连通。
+2. 后续可替换为 Prometheus/StatsD 真正上报。
+依赖关系：
+1. 被 application.go 注入到装饰器。
+复现检查：
+1. 替换 Inc 内容后可观察指标输出变化。
+
+#### internal/order/app/query/get_customer_order.go
+文件职责：查询订单用例。
+关键理解：
+1. Handle 只做一件事：调用 repo.Get。
+2. 构造函数里应用了 query decorators。
+依赖关系：
+1. 被 HTTP Get handler 调用。
+2. 依赖 order repository 接口。
+复现检查：
+1. 用存在/不存在订单号分别调用，观察成功与错误路径。
+
+#### internal/order/app/query/service.go
+文件职责：定义 Order 侧需要的库存服务端口接口。
+关键理解：
+1. 这是“防腐层端口”，隔离外部 grpc 客户端细节。
+2. create_order 依赖该接口而非具体实现。
+依赖关系：
+1. 被 create_order.go 依赖。
+2. 被 adapters/grpc/stock_grpc.go 实现。
+复现检查：
+1. 用 mock 实现该接口可做单元测试。
+
+#### internal/order/app/command/create_order.go
+文件职责：创建订单用例。
+关键理解：
+1. validate 先检查 items 非空。
+2. packItems 合并重复商品数量，减少重复校验。
+3. 调用 StockService.CheckIfItemsInStock 进行库存校验。
+4. 校验通过后调用 orderRepo.Create。
+依赖关系：
+1. 被 HTTP POST handler 调用。
+2. 依赖 order repo + stock service 接口。
+复现检查：
+1. 传空 items 应报错。
+2. 传重复 item id 时会被合并后再请求库存服务。
+
+#### internal/order/app/command/update_order.go
+文件职责：更新订单用例。
+关键理解：
+1. UpdateFn 由上层传入，决定如何改订单。
+2. handler 只负责流程控制和仓储调用。
+依赖关系：
+1. 被上层更新场景调用。
+2. 依赖 order repo.Update。
+复现检查：
+1. UpdateFn 为空时走默认 no-op 逻辑。
+
+#### internal/order/adapters/grpc/stock_grpc.go
+文件职责：StockService 的 gRPC 适配器实现。
+关键理解：
+1. 封装 stockpb.StockServiceClient。
+2. 把领域需要的方法映射为具体 RPC 请求。
+依赖关系：
+1. 实现 order/app/query/service.go 的 StockService。
+2. 被 order/service/application.go 注入。
+复现检查：
+1. 断开 stock 服务后，下单库存校验应报远程调用错误。
+
+### H. 模块与工具配置文件
+
+#### internal/common/go.mod
+文件职责：common 模块依赖声明。
+关键理解：
+1. require 指明编译期直接依赖。
+2. indirect 是传递依赖。
+依赖关系：
+1. 与 common/go.sum 配套。
+复现检查：
+1. go mod tidy 后依赖变化会更新。
+
+#### internal/order/go.mod
+文件职责：order 模块依赖声明。
+关键理解：
+1. 包含 gin、grpc、viper 等运行依赖。
+2. 通过 replace（若存在）可指向本地模块。
+依赖关系：
+1. 与 order/go.sum 配套。
+复现检查：
+1. 新增 import 后 go mod tidy 会增补依赖。
+
+#### internal/stock/go.mod
+文件职责：stock 模块依赖声明。
+关键理解：
+1. 与 order 模块类似，但按 stock 需求维护。
+依赖关系：
+1. 与 stock/go.sum 配套。
+复现检查：
+1. 删除必要依赖后编译会失败。
+
+#### internal/common/go.sum
+文件职责：common 模块依赖校验和锁文件。
+关键理解：
+1. 每行是某个模块版本的 hash，不是业务代码。
+2. 用于防篡改和可重复构建。
+依赖关系：
+1. 被 Go module 下载与校验流程读取。
+复现检查：
+1. 删除后重新 tidy 会再生成。
+
+#### internal/order/go.sum
+文件职责：order 模块依赖校验和。
+关键理解：
+1. 本质与 common/go.sum 相同。
+依赖关系：
+1. 由 go 命令自动维护。
+复现检查：
+1. 拉新依赖后会新增对应 checksum 行。
+
+#### internal/stock/go.sum
+文件职责：stock 模块依赖校验和。
+关键理解：
+1. 本质与其他 go.sum 相同。
+依赖关系：
+1. 由 go 命令自动维护。
+复现检查：
+1. tidy 后内容变化属于正常现象。
+
+#### internal/stock/.air.toml
+文件职责：stock 服务本地热重载配置。
+关键理解：
+1. 指定监听目录、忽略目录、重建命令。
+2. 开发时改代码自动重启。
+依赖关系：
+1. 被 air 工具读取。
+复现检查：
+1. 运行 air 后修改 go 文件观察自动重编译。
+
+#### internal/order/.air.toml
+文件职责：order 服务热重载配置。
+关键理解：
+1. 与 stock/.air.toml 同类。
+依赖关系：
+1. 被 air 工具读取。
+复现检查：
+1. 修改 order 代码后自动重启。
+
+#### internal/kitchen/.air.toml
+文件职责：kitchen 服务热重载配置。
+关键理解：
+1. 给 kitchen 模块开发阶段提速。
+依赖关系：
+1. 被 air 工具读取。
+复现检查：
+1. 进入 kitchen 模块运行 air 验证。
+
+#### internal/payment/.air.toml
+文件职责：payment 服务热重载配置。
+关键理解：
+1. 给 payment 模块开发阶段提速。
+依赖关系：
+1. 被 air 工具读取。
+复现检查：
+1. 进入 payment 模块运行 air 验证。
+
+---
+
+## 3. 复现步骤（按 Lesson5-11）
+
+1. 启动 stock 服务（grpc）。
+2. 启动 order 服务（http + grpc）。
+3. 调用 order 下单 HTTP 接口。
+4. order 的 create_order handler 内会通过 StockService 走 grpc 校验库存。
+5. 校验通过后写入 order in-memory repo。
+6. 再调用查询接口，验证能读到订单。
+
+如果你复现失败，优先按这个顺序查：
+1. 配置地址是否一致（global.yaml）。
+2. stock 服务是否先启动。
+3. order/service/application.go 是否成功注入 stock grpc client。
+4. proto 生成代码是否和 proto 契约一致。
+
+---
+
+## 4. Lesson5-11 的“必须会改”文件清单
+
+你练习时，优先动这些文件：
+1. internal/order/http.go
+2. internal/order/app/command/create_order.go
+3. internal/order/app/query/get_customer_order.go
+4. internal/order/adapters/order_inmem_repository.go
+5. internal/order/service/application.go
+6. internal/order/adapters/grpc/stock_grpc.go
+7. internal/stock/ports/grpc.go
+
+这些文件改熟了，后面的数据库、消息队列、可观测性才容易接上。
+
+---
+
+## 5. 下一批输出方式（你确认后继续）
+
+下一批我按同样模板写 Lesson12-Lesson20，并继续提交到 Git：
+1. 每个文件都有“职责 + 关键代码 + 调用关系 + 复现检查”。
+2. 对生成文件/go.sum 继续采用“用途级详细说明”，不做低价值逐行复述。
diff --git a/internal/common/config/viper.go b/internal/common/config/viper.go
index 103d697..5c2fe30 100644
--- a/internal/common/config/viper.go
+++ b/internal/common/config/viper.go
@@ -2,10 +2,29 @@ package config
 
 import "github.com/spf13/viper"
 
+// NewViperConfig 统一加载全局配置文件到内存。
+// 这个函数在 init() 里被调用，确保 main 启动前配置已经就位，
+// 后续各个函数可以直接用 viper.Get... 获取配置，而不用每次都读文件。
 func NewViperConfig() error {
+	// SetConfigName 指定配置文件名（不包括扩展名）。
+	// 这里是 "global"，对应 global.yaml。
 	viper.SetConfigName("global")
+
+	// SetConfigType 告诉 viper 配置格式是什么。
+	// viper 支持 json/yaml/toml/hcl 等多种格式。
 	viper.SetConfigType("yaml")
+
+	// AddConfigPath 告诉 viper 去哪个目录找配置文件。
+	// "../common/config" 是相对当前执行目录的路径，
+	// 假设你从项目根目录启动 order 服务，viper 会去 internal/common/config 目录找 global.yaml。
 	viper.AddConfigPath("../common/config")
+
+	// AutomaticEnv 让环境变量自动覆盖配置文件中的值。
+	// 举例：如果环境里设置了 ORDER_GRPC_ADDR=":6666"，会覆盖 yaml 里的 order.grpc-addr。
+	// 这样本地开发、测试、生产环境可以动态调整配置，而不改代码。
 	viper.AutomaticEnv()
+
+	// ReadInConfig 真正去磁盘读配置文件，把内容解析到 viper 的内存结构里。
+	// 如果文件不存在或格式有问题，这里会返回 error。
 	return viper.ReadInConfig()
 }
diff --git a/internal/common/decorator/command.go b/internal/common/decorator/command.go
index 53f8c37..e820637 100644
--- a/internal/common/decorator/command.go
+++ b/internal/common/decorator/command.go
@@ -6,10 +6,13 @@ import (
 	"github.com/sirupsen/logrus"
 )
 
+// CommandHandler 和 QueryHandler 对称，只是语义上代表“会改状态”的操作。
 type CommandHandler[C, R any] interface {
 	Handle(ctx context.Context, cmd C) (R, error)
 }
 
+// ApplyCommandDecorators 让所有 command 统一拥有日志和指标能力。
+// 这样新增一个写操作时，只要实现业务 handler，不用重复写埋点代码。
 func ApplyCommandDecorators[C, R any](handler CommandHandler[C, R], logger *logrus.Entry, metricsClient MetricsClient) CommandHandler[C, R] {
 	return queryLoggingDecorator[C, R]{
 		logger: logger,
diff --git a/internal/common/decorator/logging.go b/internal/common/decorator/logging.go
index a9a2325..2b158ae 100644
--- a/internal/common/decorator/logging.go
+++ b/internal/common/decorator/logging.go
@@ -8,12 +8,15 @@ import (
 	"github.com/sirupsen/logrus"
 )
 
+// queryLoggingDecorator 在真实 handler 外层补一圈日志。
+// “Decorator” 不是框架魔法，本质就是包一层结构体再转调 base。
 type queryLoggingDecorator[C, R any] struct {
 	logger *logrus.Entry
 	base   QueryHandler[C, R]
 }
 
 func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
+	// generateActionName 会把类型名拿出来，作为统一的日志 action 名。
 	logger := q.logger.WithFields(logrus.Fields{
 		"query":      generateActionName(cmd),
 		"query_body": fmt.Sprintf("%#v", cmd),
@@ -29,6 +32,7 @@ func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result
 	return q.base.Handle(ctx, cmd)
 }
 
+// generateActionName 直接复用请求类型名，避免每个 handler 手动写一次 action 常量。
 func generateActionName(cmd any) string {
 	return strings.Split(fmt.Sprintf("%T", cmd), ".")[1]
 }
diff --git a/internal/common/decorator/metrics.go b/internal/common/decorator/metrics.go
index fba9a77..ba04887 100644
--- a/internal/common/decorator/metrics.go
+++ b/internal/common/decorator/metrics.go
@@ -7,16 +7,20 @@ import (
 	"time"
 )
 
+// MetricsClient 是一个很小的抽象，方便先接假实现，后续再切到真实监控系统。
 type MetricsClient interface {
 	Inc(key string, value int)
 }
 
+// queryMetricsDecorator 统计耗时、成功、失败次数。
+// 它和日志装饰器可以自由组合，因为两者都只依赖同一个 handler 接口。
 type queryMetricsDecorator[C, R any] struct {
 	base   QueryHandler[C, R]
 	client MetricsClient
 }
 
 func (q queryMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
+	// 用 defer 统一收尾，确保无论成功还是失败都能上报指标。
 	start := time.Now()
 	actionName := strings.ToLower(generateActionName(cmd))
 	defer func() {
diff --git a/internal/common/decorator/query.go b/internal/common/decorator/query.go
index d3848fe..b5e9cf0 100644
--- a/internal/common/decorator/query.go
+++ b/internal/common/decorator/query.go
@@ -8,10 +8,13 @@ import (
 
 // QueryHandler defines a generic type that receives a Query Q,
 // and returns a result R
+// 泛型接口的好处是：同一套日志/指标装饰逻辑可以复用在不同查询上。
 type QueryHandler[Q, R any] interface {
 	Handle(ctx context.Context, query Q) (R, error)
 }
 
+// ApplyQueryDecorators 按固定顺序把日志和指标能力包在真实 handler 外面。
+// 这就是装饰器模式：不改业务代码，也能统一加横切能力。
 func ApplyQueryDecorators[H, R any](handler QueryHandler[H, R], logger *logrus.Entry, metricsClient MetricsClient) QueryHandler[H, R] {
 	return queryLoggingDecorator[H, R]{
 		logger: logger,
diff --git a/internal/common/metrics/todo_metrics.go b/internal/common/metrics/todo_metrics.go
index a05f6b2..b379f64 100644
--- a/internal/common/metrics/todo_metrics.go
+++ b/internal/common/metrics/todo_metrics.go
@@ -1,5 +1,7 @@
 package metrics
 
+// TodoMetrics 是占位实现，用来先打通接口而不真正上报指标。
+// 初学时先理解“依赖抽象”比马上接 Prometheus 更重要。
 type TodoMetrics struct{}
 
 func (t TodoMetrics) Inc(_ string, _ int) {
diff --git a/internal/common/server/gprc.go b/internal/common/server/gprc.go
index 817cfcb..ba9f6e1 100644
--- a/internal/common/server/gprc.go
+++ b/internal/common/server/gprc.go
@@ -10,12 +10,16 @@ import (
 	"google.golang.org/grpc"
 )
 
+// init 会在包加载时执行，这里统一替换 gRPC 默认日志实现，
+// 让框架日志和业务日志都走 logrus，便于初学者排查问题时只看一种格式。
 func init() {
 	logger := logrus.New()
 	logger.SetLevel(logrus.WarnLevel)
 	grpc_logrus.ReplaceGrpcLogger(logrus.NewEntry(logger))
 }
 
+// RunGRPCServer 负责按服务名读取配置，再把真正的启动动作委托给 RunGRPCServerOnAddr。
+// 这样测试时可以直接传地址，避免每次都依赖 viper 配置。
 func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server)) {
 	addr := viper.Sub(serviceName).GetString("grpc-addr")
 	if addr == "" {
@@ -25,10 +29,13 @@ func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server))
 	RunGRPCServerOnAddr(addr, registerServer)
 }
 
+// RunGRPCServerOnAddr 创建 gRPC server，并通过 registerServer 回调注册具体业务服务。
+// 这类“框架负责启动，业务负责注册”的回调模式，在 Go 后端里非常常见。
 func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
 	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
 	grpcServer := grpc.NewServer(
 		grpc.ChainUnaryInterceptor(
+			// 一元拦截器链类似 HTTP 中间件，用来放日志、指标、鉴权这类横切逻辑。
 			grpc_tags.UnaryServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
 			grpc_logrus.UnaryServerInterceptor(logrusEntry),
 			//otelgrpc.UnaryServerInterceptor(),
@@ -38,6 +45,7 @@ func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server))
 			//recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
 		),
 		grpc.ChainStreamInterceptor(
+			// 流式 RPC 和一元 RPC 分开配置，是因为两类调用的拦截接口不同。
 			grpc_tags.StreamServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
 			grpc_logrus.StreamServerInterceptor(logrusEntry),
 			//otelgrpc.StreamServerInterceptor(),
@@ -47,8 +55,10 @@ func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server))
 			//recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
 		),
 	)
+	// 业务方在这里把 OrderService/StockService 的实现挂到 grpcServer 上。
 	registerServer(grpcServer)
 
+	// net.Listen 只负责占住端口；真正开始接收 gRPC 请求要等 Serve 调用后才发生。
 	listen, err := net.Listen("tcp", addr)
 	if err != nil {
 		logrus.Panic(err)
diff --git a/internal/common/server/http.go b/internal/common/server/http.go
index 7f39359..1c561ce 100644
--- a/internal/common/server/http.go
+++ b/internal/common/server/http.go
@@ -5,6 +5,7 @@ import (
 	"github.com/spf13/viper"
 )
 
+// RunHTTPServer 按服务名从配置中读取监听地址，再交给 RunHTTPServerOnAddr 真正启动。
 func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
 	addr := viper.Sub(serviceName).GetString("http-addr")
 	if addr == "" {
@@ -13,9 +14,13 @@ func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
 	RunHTTPServerOnAddr(addr, wrapper)
 }
 
+// RunHTTPServerOnAddr 创建 gin.Engine，并把路由注册动作交给 wrapper 回调。
 func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
+	// gin.New 不带默认中间件，适合教学阶段观察“哪些能力是手动加进去的”。
 	apiRouter := gin.New()
+	// wrapper 负责把 OpenAPI 生成的 handler 绑定到具体路由。
 	wrapper(apiRouter)
+	// 这一行不会修改已有路由，只是创建了一个没有被接住的 RouterGroup。
 	apiRouter.Group("/api")
 	if err := apiRouter.Run(addr); err != nil {
 		panic(err)
diff --git a/internal/order/adapters/grpc/stock_grpc.go b/internal/order/adapters/grpc/stock_grpc.go
index c2397d6..f3acb2b 100644
--- a/internal/order/adapters/grpc/stock_grpc.go
+++ b/internal/order/adapters/grpc/stock_grpc.go
@@ -8,6 +8,8 @@ import (
 	"github.com/sirupsen/logrus"
 )
 
+// StockGRPC 是 order 侧访问 stock 服务的远程适配器。
+// 它实现的是应用层定义的 StockService 接口，而不是把 proto client 直接暴露出去。
 type StockGRPC struct {
 	client stockpb.StockServiceClient
 }
@@ -17,12 +19,14 @@ func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
 }
 
 func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error) {
+	// 这里负责把应用层参数翻译成 gRPC 请求对象。
 	resp, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
 	logrus.Info("stock_grpc response", resp)
 	return resp, err
 }
 
 func (s StockGRPC) GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error) {
+	// 对调用方来说这里只是“查商品”，至于底层走 gRPC 还是别的协议都被屏蔽了。
 	resp, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{ItemIDs: itemIDs})
 	if err != nil {
 		return nil, err
diff --git a/internal/order/adapters/order_inmem_repository.go b/internal/order/adapters/order_inmem_repository.go
index 818b51f..94364ec 100644
--- a/internal/order/adapters/order_inmem_repository.go
+++ b/internal/order/adapters/order_inmem_repository.go
@@ -10,12 +10,15 @@ import (
 	"github.com/sirupsen/logrus"
 )
 
+// MemoryOrderRepository 是 domain.Repository 的内存版实现。
+// 教学项目先用它跑通流程，后面换数据库时只需要替换这一层。
 type MemoryOrderRepository struct {
 	lock  *sync.RWMutex
 	store []*domain.Order
 }
 
 func NewMemoryOrderRepository() *MemoryOrderRepository {
+	// 这里放一条假数据，方便一开始就能演示“查询已有订单”的路径。
 	s := make([]*domain.Order, 0)
 	s = append(s, &domain.Order{
 		ID:          "fake-ID",
@@ -31,9 +34,11 @@ func NewMemoryOrderRepository() *MemoryOrderRepository {
 }
 
 func (m *MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
+	// 写操作要加互斥锁，避免多个请求并发修改切片时产生数据竞争。
 	m.lock.Lock()
 	defer m.lock.Unlock()
 	newOrder := &domain.Order{
+		// 当前用 Unix 时间戳凑一个简单 ID，真实项目里通常会换成雪花算法或 UUID。
 		ID:          strconv.FormatInt(time.Now().Unix(), 10),
 		CustomerID:  order.CustomerID,
 		Status:      order.Status,
@@ -52,6 +57,7 @@ func (m *MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*
 	for i, v := range m.store {
 		logrus.Infof("m.store[%d] = %+v", i, v)
 	}
+	// 读操作使用读锁，允许多个查询并发进行。
 	m.lock.RLock()
 	defer m.lock.RUnlock()
 	for _, o := range m.store {
@@ -64,6 +70,7 @@ func (m *MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*
 }
 
 func (m *MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
+	// UpdateFn 把“怎么改”交给上层，把“在哪里存”留给仓储层，是一种职责分离。
 	m.lock.Lock()
 	defer m.lock.Unlock()
 	found := false
diff --git a/internal/order/app/app.go b/internal/order/app/app.go
index b2e04ce..5963e7c 100644
--- a/internal/order/app/app.go
+++ b/internal/order/app/app.go
@@ -5,16 +5,19 @@ import (
 	"github.com/ghost-yu/go_shop_second/order/app/query"
 )
 
+// Application 是应用层门面，让上层只依赖一个总入口，而不用感知每个 handler 的构造细节。
 type Application struct {
 	Commands Commands
 	Queries  Queries
 }
 
+// Commands 聚合“会改状态”的用例，典型如创建、更新订单。
 type Commands struct {
 	CreateOrder command.CreateOrderHandler
 	UpdateOrder command.UpdateOrderHandler
 }
 
+// Queries 聚合“只读”的用例，便于和 Commands 做职责分离。
 type Queries struct {
 	GetCustomerOrder query.GetCustomerOrderHandler
 }
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
index 1c72b04..859e170 100644
--- a/internal/order/app/command/create_order.go
+++ b/internal/order/app/command/create_order.go
@@ -11,17 +11,20 @@ import (
 	"github.com/sirupsen/logrus"
 )
 
+// CreateOrder 是创建订单用例的输入模型。
 type CreateOrder struct {
 	CustomerID string
 	Items      []*orderpb.ItemWithQuantity
 }
 
+// CreateOrderResult 是返回给端口层的结果，避免直接暴露底层仓储对象。
 type CreateOrderResult struct {
 	OrderID string
 }
 
 type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]
 
+// createOrderHandler 只持有完成下单流程所需的两个依赖：订单仓储和库存服务。
 type createOrderHandler struct {
 	orderRepo domain.Repository
 	stockGRPC query.StockService
@@ -36,6 +39,7 @@ func NewCreateOrderHandler(
 	if orderRepo == nil {
 		panic("nil orderRepo")
 	}
+	// 下单也是 command，所以同样套上统一的日志和指标装饰器。
 	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
 		createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC},
 		logger,
@@ -44,6 +48,7 @@ func NewCreateOrderHandler(
 }
 
 func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
+	// 先校验库存，再真正创建订单，避免生成一张无法履约的脏订单。
 	validItems, err := c.validate(ctx, cmd.Items)
 	if err != nil {
 		return nil, err
@@ -62,6 +67,7 @@ func (c createOrderHandler) validate(ctx context.Context, items []*orderpb.ItemW
 	if len(items) == 0 {
 		return nil, errors.New("must have at least one item")
 	}
+	// packItems 先合并重复商品，避免把同一个商品重复发给库存服务校验。
 	items = packItems(items)
 	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, items)
 	if err != nil {
@@ -71,6 +77,7 @@ func (c createOrderHandler) validate(ctx context.Context, items []*orderpb.ItemW
 }
 
 func packItems(items []*orderpb.ItemWithQuantity) []*orderpb.ItemWithQuantity {
+	// merged 以商品 ID 为 key，把数量累加起来。
 	merged := make(map[string]int32)
 	for _, item := range items {
 		merged[item.ID] += item.Quantity
diff --git a/internal/order/app/command/update_order.go b/internal/order/app/command/update_order.go
index f40716d..3cce1e7 100644
--- a/internal/order/app/command/update_order.go
+++ b/internal/order/app/command/update_order.go
@@ -8,6 +8,7 @@ import (
 	"github.com/sirupsen/logrus"
 )
 
+// UpdateOrder 把“更新哪个订单”和“如何更新”一起传给 handler。
 type UpdateOrder struct {
 	Order    *domain.Order
 	UpdateFn func(context.Context, *domain.Order) (*domain.Order, error)
@@ -15,6 +16,7 @@ type UpdateOrder struct {
 
 type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, interface{}]
 
+// updateOrderHandler 负责流程控制，具体修改细节由 UpdateFn 注入。
 type updateOrderHandler struct {
 	orderRepo domain.Repository
 	//stockGRPC
@@ -36,6 +38,7 @@ func NewUpdateOrderHandler(
 }
 
 func (c updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (interface{}, error) {
+	// 给 nil UpdateFn 一个 no-op 默认值，避免直接调用时发生空指针问题。
 	if cmd.UpdateFn == nil {
 		logrus.Warnf("updateOrderHandler got nil UpdateFn, order=%#v", cmd.Order)
 		cmd.UpdateFn = func(_ context.Context, order *domain.Order) (*domain.Order, error) { return order, nil }
diff --git a/internal/order/app/query/get_customer_order.go b/internal/order/app/query/get_customer_order.go
index b107e01..19b7ffc 100644
--- a/internal/order/app/query/get_customer_order.go
+++ b/internal/order/app/query/get_customer_order.go
@@ -8,13 +8,17 @@ import (
 	"github.com/sirupsen/logrus"
 )
 
+// GetCustomerOrder 是查询订单的输入对象。
+// 把参数收成一个结构体后，后续加字段不会破坏函数签名。
 type GetCustomerOrder struct {
 	CustomerID string
 	OrderID    string
 }
 
+// GetCustomerOrderHandler 是一个已经套好泛型的查询处理器类型别名。
 type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Order]
 
+// getCustomerOrderHandler 才是真正的业务实现，字段里只放它需要的依赖。
 type getCustomerOrderHandler struct {
 	orderRepo domain.Repository
 }
@@ -27,6 +31,7 @@ func NewGetCustomerOrderHandler(
 	if orderRepo == nil {
 		panic("nil orderRepo")
 	}
+	// 构造函数里统一加装饰器，这样调用方拿到的就是“带日志和指标能力”的 handler。
 	return decorator.ApplyQueryDecorators[GetCustomerOrder, *domain.Order](
 		getCustomerOrderHandler{orderRepo: orderRepo},
 		logger,
@@ -35,6 +40,7 @@ func NewGetCustomerOrderHandler(
 }
 
 func (g getCustomerOrderHandler) Handle(ctx context.Context, query GetCustomerOrder) (*domain.Order, error) {
+	// 查询用例本身很薄，只负责描述流程，把数据访问交给仓储接口。
 	o, err := g.orderRepo.Get(ctx, query.OrderID, query.CustomerID)
 	if err != nil {
 		return nil, err
diff --git a/internal/order/app/query/service.go b/internal/order/app/query/service.go
index 2e3e4f4..0d33749 100644
--- a/internal/order/app/query/service.go
+++ b/internal/order/app/query/service.go
@@ -7,6 +7,8 @@ import (
 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
 )
 
+// StockService 是 order 应用层眼里的“库存能力端口”。
+// 它隔离了底层 gRPC 细节，让 create_order 用例只表达业务意图。
 type StockService interface {
 	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error)
 	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
diff --git a/internal/order/domain/order/order.go b/internal/order/domain/order/order.go
index 6a9e94b..073b8f0 100644
--- a/internal/order/domain/order/order.go
+++ b/internal/order/domain/order/order.go
@@ -2,10 +2,16 @@ package order
 
 import "github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 
+// Order 是订单领域对象，代表业务里真正被创建、查询、更新的核心实体。
 type Order struct {
-	ID          string
-	CustomerID  string
-	Status      string
+	// ID 是订单主键，由仓储层创建时生成。
+	ID string
+	// CustomerID 表示订单属于哪个客户。
+	CustomerID string
+	// Status 预留给支付中、已支付等订单状态流转。
+	Status string
+	// PaymentLink 预留给支付服务返回的支付链接。
 	PaymentLink string
-	Items       []*orderpb.Item
+	// Items 直接复用 proto 里的商品结构，减少服务边界两侧的对象转换成本。
+	Items []*orderpb.Item
 }
diff --git a/internal/order/domain/order/repository.go b/internal/order/domain/order/repository.go
index 04b783f..685eddc 100644
--- a/internal/order/domain/order/repository.go
+++ b/internal/order/domain/order/repository.go
@@ -5,6 +5,7 @@ import (
 	"fmt"
 )
 
+// Repository 定义订单持久化能力，应用层只依赖这个接口，而不关心底层是内存还是数据库。
 type Repository interface {
 	Create(context.Context, *Order) (*Order, error)
 	Get(ctx context.Context, id, customerID string) (*Order, error)
@@ -15,6 +16,7 @@ type Repository interface {
 	) error
 }
 
+// NotFoundError 是带业务语义的错误类型，调用方可以明确知道“订单不存在”而不是普通异常。
 type NotFoundError struct {
 	OrderID string
 }
diff --git a/internal/order/http.go b/internal/order/http.go
index b40adc7..8b5d70b 100644
--- a/internal/order/http.go
+++ b/internal/order/http.go
@@ -10,16 +10,20 @@ import (
 	"github.com/gin-gonic/gin"
 )
 
+// HTTPServer 是 OpenAPI 生成接口在 order 服务里的具体实现。
+// 它本身不处理复杂业务，只负责把 HTTP 请求翻译成应用层调用。
 type HTTPServer struct {
 	app app.Application
 }
 
 func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
 	var req orderpb.CreateOrderRequest
+	// ShouldBindJSON 负责把请求体反序列化为 proto 请求对象。
 	if err := c.ShouldBindJSON(&req); err != nil {
 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
 		return
 	}
+	// 进入应用层前，把 HTTP 层对象转换成 command 对象，避免业务逻辑依赖 gin。
 	r, err := H.app.Commands.CreateOrder.Handle(c, command.CreateOrder{
 		CustomerID: req.CustomerID,
 		Items:      req.Items,
@@ -36,6 +40,7 @@ func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID stri
 }
 
 func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
+	// 查询场景走 Queries 分组，体现读写分离的组织方式。
 	o, err := H.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
 		OrderID:    orderID,
 		CustomerID: customerID,
diff --git a/internal/order/main.go b/internal/order/main.go
index 55d022b..9807588 100644
--- a/internal/order/main.go
+++ b/internal/order/main.go
@@ -14,6 +14,7 @@ import (
 	"google.golang.org/grpc"
 )
 
+// init 先加载全局配置，这样 main 里启动服务时就能直接从 viper 读取地址和服务名。
 func init() {
 	if err := config.NewViperConfig(); err != nil {
 		logrus.Fatal(err)
@@ -21,19 +22,24 @@ func init() {
 }
 
 func main() {
+	// serviceName 对应 global.yaml 里的 order.service-name，后续 server 包会继续用它找端口配置。
 	serviceName := viper.GetString("order.service-name")
 
+	// ctx 用来把进程级生命周期传给下游依赖，例如 gRPC 客户端连接。
 	ctx, cancel := context.WithCancel(context.Background())
 	defer cancel()
 
+	// NewApplication 是组合根：在这里把 repo、远程 client、handler 全部装配好。
 	application, cleanup := service.NewApplication(ctx)
 	defer cleanup()
 
+	// gRPC 和 HTTP 同时启动，说明 order 服务对内和对外用了两套协议暴露能力。
 	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
 		svc := ports.NewGRPCServer(application)
 		orderpb.RegisterOrderServiceServer(server, svc)
 	})
 
+	// HTTP 侧通过 OpenAPI 生成的 RegisterHandlersWithOptions 完成路由到实现的映射。
 	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
 		ports.RegisterHandlersWithOptions(router, HTTPServer{
 			app: application,
diff --git a/internal/order/ports/grpc.go b/internal/order/ports/grpc.go
index e1e8621..4841586 100644
--- a/internal/order/ports/grpc.go
+++ b/internal/order/ports/grpc.go
@@ -8,6 +8,8 @@ import (
 	"google.golang.org/protobuf/types/known/emptypb"
 )
 
+// GRPCServer 是 orderpb.OrderServiceServer 的端口适配器实现。
+// 这一层的职责和 HTTPServer 一样，都是把外部协议请求转给应用层。
 type GRPCServer struct {
 	app app.Application
 }
@@ -17,16 +19,19 @@ func NewGRPCServer(app app.Application) *GRPCServer {
 }
 
 func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
+	// lesson 当前故意保留 TODO，方便后续逐步实现 gRPC 版本的下单流程。
 	//TODO implement me
 	panic("implement me")
 }
 
 func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
+	// 这里最终会把 request 转为 query，再返回 proto 层定义的 orderpb.Order。
 	//TODO implement me
 	panic("implement me")
 }
 
 func (G GRPCServer) UpdateOrder(ctx context.Context, order *orderpb.Order) (*emptypb.Empty, error) {
+	// 更新接口通常会调用 Commands.UpdateOrder；当前阶段先保留占位。
 	//TODO implement me
 	panic("implement me")
 }
diff --git a/internal/order/service/application.go b/internal/order/service/application.go
index 648e438..faf1f24 100644
--- a/internal/order/service/application.go
+++ b/internal/order/service/application.go
@@ -13,11 +13,14 @@ import (
 	"github.com/sirupsen/logrus"
 )
 
+// NewApplication 是 order 服务的组合根。
+// 它负责创建外部依赖、构造适配器，并返回一个已经装配好的应用层门面。
 func NewApplication(ctx context.Context) (app.Application, func()) {
 	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
 	if err != nil {
 		panic(err)
 	}
+	// 这里把生成的 gRPC client 再包一层适配器，避免应用层直接依赖 proto 细节。
 	stockGRPC := grpc.NewStockGRPC(stockClient)
 	return newApplication(ctx, stockGRPC), func() {
 		_ = closeStockClient()
@@ -25,6 +28,7 @@ func NewApplication(ctx context.Context) (app.Application, func()) {
 }
 
 func newApplication(_ context.Context, stockGRPC query.StockService) app.Application {
+	// 组合根里统一决定“接口对应哪种实现”，后面要切数据库或 mock 时只改这里。
 	orderRepo := adapters.NewMemoryOrderRepository()
 	logger := logrus.NewEntry(logrus.StandardLogger())
 	metricClient := metrics.TodoMetrics{}
diff --git a/internal/stock/adapters/stock_inmem_repository.go b/internal/stock/adapters/stock_inmem_repository.go
index 4390124..c38b115 100644
--- a/internal/stock/adapters/stock_inmem_repository.go
+++ b/internal/stock/adapters/stock_inmem_repository.go
@@ -8,11 +8,14 @@ import (
 	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
 )
 
+// MemoryStockRepository 用 map 模拟库存数据源。
+// 相比切片，map 更适合按商品 ID 做快速查找。
 type MemoryStockRepository struct {
 	lock  *sync.RWMutex
 	store map[string]*orderpb.Item
 }
 
+// stub 是教学阶段的假数据，用来模拟一个已经存在的库存商品。
 var stub = map[string]*orderpb.Item{
 	"item_id": {
 		ID:       "foo_item",
@@ -30,6 +33,7 @@ func NewMemoryStockRepository() *MemoryStockRepository {
 }
 
 func (m MemoryStockRepository) GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error) {
+	// 这里只有读场景，所以拿读锁即可。
 	m.lock.RLock()
 	defer m.lock.RUnlock()
 	var (
@@ -46,5 +50,6 @@ func (m MemoryStockRepository) GetItems(ctx context.Context, ids []string) ([]*o
 	if len(res) == len(ids) {
 		return res, nil
 	}
+	// 返回“部分结果 + 缺失错误”可以帮助上层更明确地决定如何提示用户。
 	return res, domain.NotFoundError{Missing: missing}
 }
diff --git a/internal/stock/app/app.go b/internal/stock/app/app.go
index 42330cd..c159f54 100644
--- a/internal/stock/app/app.go
+++ b/internal/stock/app/app.go
@@ -1,10 +1,14 @@
 package app
 
+// Application 预留给 stock 服务聚合 Commands 和 Queries。
+// 当前 lesson 里 stock 逻辑还比较薄，所以两个分组都是空壳结构。
 type Application struct {
 	Commands Commands
 	Queries  Queries
 }
 
+// Commands 未来承载库存写操作，例如扣减库存。
 type Commands struct{}
 
+// Queries 未来承载库存读操作，例如按商品 ID 查询库存。
 type Queries struct{}
diff --git a/internal/stock/domain/stock/repository.go b/internal/stock/domain/stock/repository.go
index 7a58e6a..b0f1796 100644
--- a/internal/stock/domain/stock/repository.go
+++ b/internal/stock/domain/stock/repository.go
@@ -8,10 +8,13 @@ import (
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 )
 
+// Repository 定义库存服务需要提供的最小数据访问能力。
+// order 服务只要通过应用层依赖这个接口，就不需要知道库存数据放在哪里。
 type Repository interface {
 	GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error)
 }
 
+// NotFoundError 用来明确告诉调用方哪些商品在库存里不存在。
 type NotFoundError struct {
 	Missing []string
 }
diff --git a/internal/stock/main.go b/internal/stock/main.go
index 78ca718..e1fe459 100644
--- a/internal/stock/main.go
+++ b/internal/stock/main.go
@@ -13,6 +13,7 @@ import (
 	"google.golang.org/grpc"
 )
 
+// stock 服务和 order 服务一样，在进程入口先加载配置，避免后续各层自己读文件。
 func init() {
 	if err := config.NewViperConfig(); err != nil {
 		logrus.Fatal(err)
@@ -20,17 +21,21 @@ func init() {
 }
 
 func main() {
+	// serviceName 决定去哪个配置分组下找 grpc-addr 等参数。
 	serviceName := viper.GetString("stock.service-name")
+	// serverType 让同一个二进制可以选择跑哪种协议，当前 lesson 主要走 grpc 分支。
 	serverType := viper.GetString("stock.server-to-run")
 
 	logrus.Info(serverType)
 
+	// 这里先组装应用层，再根据配置决定挂到哪个端口适配器上。
 	ctx, cancel := context.WithCancel(context.Background())
 	defer cancel()
 
 	application := service.NewApplication(ctx)
 	switch serverType {
 	case "grpc":
+		// gRPC 分支把 ports.GRPCServer 注册为 StockService 的服务端实现。
 		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
 			svc := ports.NewGRPCServer(application)
 			stockpb.RegisterStockServiceServer(server, svc)
diff --git a/internal/stock/ports/grpc.go b/internal/stock/ports/grpc.go
index fb41c1f..72bc356 100644
--- a/internal/stock/ports/grpc.go
+++ b/internal/stock/ports/grpc.go
@@ -7,6 +7,8 @@ import (
 	"github.com/ghost-yu/go_shop_second/stock/app"
 )
 
+// GRPCServer 负责实现 stockpb.StockServiceServer 接口。
+// 外部请求进入后，应该在这里完成参数转换，再转给应用层处理。
 type GRPCServer struct {
 	app app.Application
 }
@@ -16,11 +18,13 @@ func NewGRPCServer(app app.Application) *GRPCServer {
 }
 
 func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
+	// 这里后续会把商品 ID 列表转给库存查询用例。
 	//TODO implement me
 	panic("implement me")
 }
 
 func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
+	// 这一步对应 order 下单前的库存校验，是跨服务调用的关键入口。
 	//TODO implement me
 	panic("implement me")
 }
diff --git a/internal/stock/service/application.go b/internal/stock/service/application.go
index 423b368..23daf9a 100644
--- a/internal/stock/service/application.go
+++ b/internal/stock/service/application.go
@@ -6,6 +6,8 @@ import (
 	"github.com/ghost-yu/go_shop_second/stock/app"
 )
 
+// NewApplication 当前先返回一个空的应用层门面。
+// lesson 后续扩展库存能力时，repo 和 handler 的装配点也会放在这里。
 func NewApplication(ctx context.Context) app.Application {
 	return app.Application{}
 }
~~~

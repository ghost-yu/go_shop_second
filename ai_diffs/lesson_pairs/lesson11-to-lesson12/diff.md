# Lesson Pair Diff Report

- FromBranch: lesson11
- ToBranch: lesson12

## Short Summary

~~~text
 41 files changed, 928 insertions(+), 903 deletions(-)
~~~

## File Stats

~~~text
 docker-compose.yml                                 |   8 +
 ...273\243\347\240\201\350\247\243\350\257\273.md" | 716 ---------------------
 internal/common/client/grpc.go                     |  29 +
 internal/common/config/viper.go                    |  19 -
 internal/common/decorator/command.go               |   3 -
 internal/common/decorator/logging.go               |   4 -
 internal/common/decorator/metrics.go               |   4 -
 internal/common/decorator/query.go                 |   3 -
 internal/common/discovery/consul/consul.go         |  88 +++
 internal/common/discovery/discovery.go             |  20 +
 internal/common/discovery/grpc.go                  |  37 ++
 internal/common/go.mod                             |  27 +-
 internal/common/go.sum                             | 194 +++++-
 internal/common/metrics/todo_metrics.go            |   2 -
 internal/common/server/gprc.go                     |  10 -
 internal/common/server/http.go                     |   5 -
 internal/order/adapters/grpc/stock_grpc.go         |   4 -
 internal/order/adapters/order_inmem_repository.go  |   7 -
 internal/order/app/app.go                          |   3 -
 internal/order/app/command/create_order.go         |   7 -
 internal/order/app/command/update_order.go         |   3 -
 internal/order/app/query/get_customer_order.go     |   6 -
 internal/order/app/query/service.go                |   2 -
 internal/order/domain/order/order.go               |  14 +-
 internal/order/domain/order/repository.go          |   2 -
 internal/order/go.mod                              |  13 +
 internal/order/go.sum                              | 167 +++++
 internal/order/http.go                             |   5 -
 internal/order/main.go                             |  15 +-
 internal/order/ports/grpc.go                       |   5 -
 internal/order/service/application.go              |   4 -
 internal/stock/adapters/stock_inmem_repository.go  |  23 +-
 internal/stock/app/app.go                          |  11 +-
 .../stock/app/query/check_if_items_in_stock.go     |  46 ++
 internal/stock/app/query/get_items.go              |  43 ++
 internal/stock/domain/stock/repository.go          |   3 -
 internal/stock/go.mod                              |  29 +-
 internal/stock/go.sum                              | 194 +++++-
 internal/stock/main.go                             |  15 +-
 internal/stock/ports/grpc.go                       |  22 +-
 internal/stock/service/application.go              |  19 +-
 41 files changed, 928 insertions(+), 903 deletions(-)
~~~

## Commit Comparison

~~~text
< 3163c56 docs: add comments to viper config loader
< 853d6f0 docs: add learning comments for backend internship study
< 2d6db2e docs: gprc.go & http.go 补全完整代码+逐段详细注释
< e2e4b04 docs: 扩展 Lesson5-11 为文件级复现手册
< 6d9bcf7 docs: 添加 Lesson5-11 关键代码解读（Go 小白版）
> 9a37c1c consul
~~~

## Changed Files

~~~text
docker-compose.yml
"docs/lesson-notes/lesson5-lesson11-\345\205\263\351\224\256\344\273\243\347\240\201\350\247\243\350\257\273.md"
internal/common/client/grpc.go
internal/common/config/viper.go
internal/common/decorator/command.go
internal/common/decorator/logging.go
internal/common/decorator/metrics.go
internal/common/decorator/query.go
internal/common/discovery/consul/consul.go
internal/common/discovery/discovery.go
internal/common/discovery/grpc.go
internal/common/go.mod
internal/common/go.sum
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
internal/order/go.mod
internal/order/go.sum
internal/order/http.go
internal/order/main.go
internal/order/ports/grpc.go
internal/order/service/application.go
internal/stock/adapters/stock_inmem_repository.go
internal/stock/app/app.go
internal/stock/app/query/check_if_items_in_stock.go
internal/stock/app/query/get_items.go
internal/stock/domain/stock/repository.go
internal/stock/go.mod
internal/stock/go.sum
internal/stock/main.go
internal/stock/ports/grpc.go
internal/stock/service/application.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
docker-compose.yml
"docs/lesson-notes/lesson5-lesson11-\345\205\263\351\224\256\344\273\243\347\240\201\350\247\243\350\257\273.md"
internal/common/client/grpc.go
internal/common/config/viper.go
internal/common/decorator/command.go
internal/common/decorator/logging.go
internal/common/decorator/metrics.go
internal/common/decorator/query.go
internal/common/discovery/consul/consul.go
internal/common/discovery/discovery.go
internal/common/discovery/grpc.go
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
internal/stock/app/query/check_if_items_in_stock.go
internal/stock/app/query/get_items.go
internal/stock/domain/stock/repository.go
internal/stock/main.go
internal/stock/ports/grpc.go
internal/stock/service/application.go
~~~

## Full Diff

~~~diff
diff --git a/docker-compose.yml b/docker-compose.yml
new file mode 100644
index 0000000..079748a
--- /dev/null
+++ b/docker-compose.yml
@@ -0,0 +1,8 @@
+version: "3.9"
+services:
+  consul:
+    image: hashicorp/consul
+    command: agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0
+    ports:
+      - 8500:8500
+      - 8600:8600/udp
\ No newline at end of file
diff --git "a/docs/lesson-notes/lesson5-lesson11-\345\205\263\351\224\256\344\273\243\347\240\201\350\247\243\350\257\273.md" "b/docs/lesson-notes/lesson5-lesson11-\345\205\263\351\224\256\344\273\243\347\240\201\350\247\243\350\257\273.md"
deleted file mode 100644
index fc31db2..0000000
--- "a/docs/lesson-notes/lesson5-lesson11-\345\205\263\351\224\256\344\273\243\347\240\201\350\247\243\350\257\273.md"
+++ /dev/null
@@ -1,716 +0,0 @@
-# gorder-v2 Lesson5-Lesson11 复现级代码解读（Go 小白版）
-
-这版文档目标是“可复现”，不是只做概念导读。
-你可以按本文件逐个文件阅读并运行，每个文件都给出：
-1. 文件职责
-2. 关键代码必须理解点
-3. 上下游依赖（谁调用它，它调用谁）
-4. 复现检查动作（你可以自己执行验证）
-
-## 0. 先说明边界
-
-Lesson5-Lesson11 里涉及的文件很多，其中有三类：
-1. 业务代码文件：需要详细读懂。
-2. 生成文件（*.pb.go, *_grpc.pb.go）：需要理解结构和用途，不建议逐行硬背。
-3. 依赖锁文件（go.sum）：不需要逐条解释，每行本质是模块校验和。
-
-因此本文对所有文件都做说明，但会按“学习价值”给不同颗粒度，确保你能复现整个阶段。
-
-## 1. Lesson5 到 Lesson11 的总调用图
-
-1. 配置加载：internal/common/config/global.yaml -> viper
-2. 进程入口：internal/order/main.go、internal/stock/main.go
-3. 基础服务启动：internal/common/server/http.go、internal/common/server/gprc.go
-4. 端口层：internal/order/http.go、internal/order/ports/grpc.go、internal/stock/ports/grpc.go
-5. 应用层：internal/order/app/...、internal/stock/app/...
-6. 领域接口：internal/order/domain/order/repository.go、internal/stock/domain/stock/repository.go
-7. 适配器实现：in-memory repo + grpc 适配器
-8. 协议层：api/*.proto + internal/common/genproto/*.pb.go
-
----
-
-## 2. 按文件详细解释（Lesson5-Lesson11）
-
-### A. 协议与配置文件
-
-#### api/stockpb/stock.proto
-文件职责：定义库存服务 RPC 契约。
-关键理解：
-1. service StockService 里定义了两个 RPC：GetItems、CheckIfItemsInStock。
-2. 请求和响应消息是跨服务通信唯一标准，服务端和客户端必须一致。
-3. 引用 orderpb/order.proto 是为了复用订单侧 Item 结构，避免重复建模。
-依赖关系：
-1. 被 internal/common/genproto/stockpb/stock.pb.go 和 stock_grpc.pb.go 生成工具使用。
-2. 被 order 侧 grpc 客户端调用路径间接使用。
-复现检查：
-1. 运行 proto 生成脚本后，确认生成代码存在对应消息与服务接口。
-
-#### api/orderpb/order.proto
-文件职责：定义订单服务对外协议和公共消息（如 Item）。
-关键理解：
-1. Item 和 ItemWithQuantity 是库存校验与下单共用的数据结构。
-2. OrderService 方法签名决定了 order 端口层函数形态。
-依赖关系：
-1. 被 OpenAPI/GRPC 生成代码使用。
-2. 被 stock.proto import。
-复现检查：
-1. 改字段名后重新生成，会导致调用方编译错误，证明契约约束生效。
-
-#### internal/common/config/global.yaml
-文件职责：集中配置各服务地址与运行模式。
-关键理解：
-1. order 与 stock 的 grpc/http 地址是独立的。
-2. fallback-grpc-addr 是防御性配置，避免空地址直接崩溃。
-依赖关系：
-1. 被 config/viper 初始化后读取。
-2. 被 server/http.go 和 server/gprc.go 间接消费。
-复现检查：
-1. 改 grpc-addr 后重启服务，端口变化应立即生效。
-
-#### .gitignore
-文件职责：控制 Git 不追踪哪些本地文件。
-关键理解：
-1. 用于忽略临时构建产物和本地工具缓存。
-2. 防止开发机噪音提交进仓库。
-依赖关系：
-1. Git 客户端直接读取。
-复现检查：
-1. 新建被忽略文件后 git status 不应出现。
-
-### B. 生成代码文件
-
-#### internal/common/genproto/orderpb/order.pb.go
-文件职责：order.proto 的消息结构与序列化代码。
-关键理解：
-1. 这类文件由 protoc 自动生成，不手改。
-2. 包含 message struct、反射元数据、编解码实现。
-依赖关系：
-1. 被业务代码 import，作为请求/响应 DTO。
-复现检查：
-1. 删除后项目会出现大量类型缺失错误。
-
-#### internal/common/genproto/stockpb/stock.pb.go
-文件职责：stock.proto 消息生成代码。
-关键理解：
-1. 与 order.pb.go 同类，承载跨服务消息类型。
-依赖关系：
-1. 被 stock 服务端与 order 客户端共同依赖。
-复现检查：
-1. 运行 stock grpc 相关代码需依赖此文件类型。
-
-#### internal/common/genproto/stockpb/stock_grpc.pb.go
-文件职责：stock.proto 的 gRPC 客户端/服务端接口生成代码。
-关键理解：
-1. 生成 StockServiceClient 和 StockServiceServer 接口。
-2. order 侧 grpc 适配器通过 client 发请求。
-3. stock 侧端口实现通过 server 接口对外暴露。
-依赖关系：
-1. 被 internal/order/adapters/grpc/stock_grpc.go 使用。
-2. 被 internal/stock/ports/grpc.go 实现。
-复现检查：
-1. 在 main.go 注册 stockpb.RegisterStockServiceServer 即依赖此文件。
-
-### C. 基础设施启动文件
-
-#### internal/common/server/gprc.go
-文件职责：统一 gRPC server 启动流程，所有 gRPC 服务都通过这里启动，业务服务只需提供「注册业务」这一个回调。
-
-完整代码+逐段注释：
-```go
-package server
-
-import (
-    "net"
-    grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
-    grpc_tags   "github.com/grpc-ecosystem/go-grpc-middleware/tags"
-    "github.com/sirupsen/logrus"
-    "github.com/spf13/viper"
-    "google.golang.org/grpc"
-)
-
-// init 是 Go 的特殊函数，包被 import 时自动执行，不需要手动调用。
-// 这里把 gRPC 框架内部日志替换成项目统一使用的 logrus，同时把级别设为 Warn，
-// 减少 gRPC 框架内部的噪音日志，保持日志格式一致，方便集中收集和排查问题。
-func init() {
-    logger := logrus.New()
-    logger.SetLevel(logrus.WarnLevel)                      // 只打 warn 及以上
-    grpc_logrus.ReplaceGrpcLogger(logrus.NewEntry(logger)) // 替换 gRPC 内置 logger
-}
-
-// RunGRPCServer 是对外的唯一入口。
-// 参数 serviceName：对应 global.yaml 里的 key，比如 "order" 或 "stock"。
-// 参数 registerServer：业务方传入的回调，用于把自己的 handler 注册进 grpc.Server。
-//
-// 调用方写法（order/main.go）示例：
-//   server.RunGRPCServer("order", func(s *grpc.Server) {
-//       orderpb.RegisterOrderServiceServer(s, ports.NewGRPCServer(app))
-//   })
-func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server)) {
-    // viper.Sub("order").GetString("grpc-addr") 相当于读 global.yaml 里 order.grpc-addr 的值。
-    addr := viper.Sub(serviceName).GetString("grpc-addr")
-    if addr == "" {
-        // 地址缺失时用全局兜底地址，防止配置漏填导致进程直接崩溃。
-        addr = viper.GetString("fallback-grpc-addr")
-    }
-    RunGRPCServerOnAddr(addr, registerServer)
-}
-
-// RunGRPCServerOnAddr 真正创建并运行 gRPC server。
-// 独立函数的原因：测试时可以直接传 "127.0.0.1:0"（随机端口），不依赖配置文件。
-func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
-    logrusEntry := logrus.NewEntry(logrus.StandardLogger())
-
-    // grpc.NewServer(...) 创建 gRPC 服务器，括号里配置拦截器。
-    // 拦截器 = HTTP 框架里的中间件，每次 RPC 调用进来和出去时都会经过它们，
-    // 适合放日志、追踪、鉴权、panic 恢复等和业务无关的横切逻辑。
-    grpcServer := grpc.NewServer(
-
-        // ChainUnaryInterceptor：配置「一元 RPC」拦截器链（普通请求-响应模式，最常用）。
-        // 执行顺序：请求进来时从左到右，响应返回时从右到左。
-        grpc.ChainUnaryInterceptor(
-            // ① 从 proto 生成的请求结构里自动提取字段（如 customer_id）打到日志 Tag 里。
-            //   后续该请求的所有日志都自动携带这些字段，排查问题时不用手动传。
-            grpc_tags.UnaryServerInterceptor(
-                grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor),
-            ),
-            // ② 自动记录每次 RPC 的方法名、耗时、gRPC 状态码到 logrus。
-            grpc_logrus.UnaryServerInterceptor(logrusEntry),
-
-            // 下面是后续 lesson 会逐步打开的能力，当前注释保留位置：
-            // otelgrpc.UnaryServerInterceptor(),       // 分布式链路追踪（Lesson25+）
-            // srvMetrics.UnaryServerInterceptor(...),  // Prometheus 指标采集（Lesson43+）
-            // selector.UnaryServerInterceptor(...),    // 选择性鉴权（某些路由才需要鉴权）
-            // recovery.UnaryServerInterceptor(...),    // panic 自动恢复（防止一个请求崩掉整个进程）
-        ),
-
-        // ChainStreamInterceptor：配置「流式 RPC」拦截器链（双向流/服务端推送场景）。
-        // 当前项目主要用一元 RPC，这里配置和一元保持同步即可。
-        grpc.ChainStreamInterceptor(
-            grpc_tags.StreamServerInterceptor(
-                grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor),
-            ),
-            grpc_logrus.StreamServerInterceptor(logrusEntry),
-            // 同上，预留位置暂时注释
-        ),
-    )
-
-    // 执行业务注册回调，把 OrderService/StockService 的具体实现挂到 grpcServer 上。
-    // 这行之前 grpcServer 是空的（不知道有哪些方法），执行后才知道怎么处理请求。
-    registerServer(grpcServer)
-
-    // net.Listen 在 TCP 层监听指定地址，此时还没开始接受 gRPC 请求。
-    listen, err := net.Listen("tcp", addr)
-    if err != nil {
-        // 端口被占用、权限不足等情况，直接 panic 让进程退出，而不是静默失败。
-        logrus.Panic(err)
-    }
-    logrus.Infof("Starting gRPC server, Listening: %s", addr)
-
-    // grpcServer.Serve 是阻塞调用，进入后一直处理 gRPC 请求，直到 server Stop 或出错。
-    if err := grpcServer.Serve(listen); err != nil {
-        logrus.Panic(err)
-    }
-}
-```
-
-核心设计要点：
-1. **两函数分离**：RunGRPCServer 负责读配置（依赖 viper），RunGRPCServerOnAddr 只依赖
-   地址字符串。单元测试时可以绕过配置直接调用 RunGRPCServerOnAddr，不需要 global.yaml。
-2. **registerServer 回调模式**：框架不知道业务注册了什么，业务不知道框架怎么启动，
-   通过这一个函数参数完成解耦。这是 Go 里高频出现的依赖注入写法，全项目会反复看到。
-3. **拦截器链是功能扩展点**：目前只有日志+标签两层，后续每个注释行都是可独立打开的能力，
-   打开时完全不需要改任何业务代码，只需取消注释并引入对应包。
-
-依赖关系：
-- 上游调用方：order/main.go 和 stock/main.go 传入 serviceName 和 registerServer 调用本函数。
-- 下游依赖：registerServer 回调由业务服务提供（例如 orderpb.RegisterOrderServiceServer）。
-
-复现检查：
-1. 启动后终端出现 `Starting gRPC server, Listening: 127.0.0.1:5002`。
-2. 改 global.yaml 的 grpc-addr 后重启，监听端口随之变化。
-3. 用重复端口启动第二个实例，net.Listen 会报 bind: address already in use 并 panic。
-
----
-
-#### internal/common/server/http.go
-文件职责：统一 HTTP server 启动流程，所有 HTTP 服务都通过这里启动，业务方只需传入路由注册回调。
-
-完整代码+逐段注释：
-```go
-package server
-
-import (
-    "github.com/gin-gonic/gin"
-    "github.com/spf13/viper"
-)
-
-// RunHTTPServer 按服务名从 viper 读 http-addr，然后启动 HTTP 服务。
-// 参数 wrapper：业务方提供的路由注册函数，接收 *gin.Engine 并往里注册路由和 handler。
-//
-// 调用方写法（order/main.go）示例：
-//   server.RunHTTPServer("order", func(router *gin.Engine) {
-//       ports.RegisterHandlersWithOptions(router, HTTPServer{app: application}, opts)
-//   })
-func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
-    // viper.Sub("order").GetString("http-addr") 读取 global.yaml 里 order.http-addr 的值。
-    addr := viper.Sub(serviceName).GetString("http-addr")
-    if addr == "" {
-        // TODO：后续应加 warning 日志提醒配置缺失。
-        // 注意：addr 为空时 gin.Run("") 会监听 :80，非 root 权限下会权限错误。
-    }
-    RunHTTPServerOnAddr(addr, wrapper)
-}
-
-// RunHTTPServerOnAddr 真正创建 gin 引擎并运行。
-// 独立出来的原因：测试时可以直接传 ":0"（随机端口），不依赖 viper 配置。
-func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
-    // gin.New() 创建一个没有任何默认中间件的干净 gin 引擎。
-    // 区别于 gin.Default()：Default() 会自动加 Logger 和 Recovery 中间件，
-    // 这里用 New() 是为了完全掌控中间件，后续 lesson 会按需手动添加。
-    apiRouter := gin.New()
-
-    // 执行业务方传入的路由注册函数，把路由和 handler 挂到 apiRouter 上。
-    // 执行完后 apiRouter 才知道"POST /api/customer/.../orders"对应哪个 handler。
-    wrapper(apiRouter)
-
-    // ⚠️ 注意：下面这行实际上是无效代码（遗留 bug）。
-    // apiRouter.Group("/api") 返回一个新的 RouterGroup 对象，但没有赋值给变量，
-    // 返回值被直接丢弃，对已有路由没有任何影响，相当于执行了一个空操作。
-    // 真正的 /api 前缀是在 wrapper 内部通过 RegisterHandlersWithOptions 的 BaseURL:"/api" 设置的。
-    // 理解这一点可以防止你自己写出同类 bug。
-    apiRouter.Group("/api")
-
-    // apiRouter.Run 是阻塞调用，开始监听并处理 HTTP 请求。
-    // 出错（端口冲突、权限不足）时直接 panic，不要静默失败。
-    if err := apiRouter.Run(addr); err != nil {
-        panic(err)
-    }
-}
-```
-
-核心设计要点：
-1. **gin.New() 而非 gin.Default()**：主动掌控中间件，不被框架默认行为绑架。
-   需要加 cors、recovery、自定义日志时直接往 apiRouter 上加，与 grpc.go 的拦截器链思路一致。
-2. **wrapper 回调模式**：与 gprc.go 的 registerServer 完全相同的解耦思路，
-   http.go 不知道有哪些路由，order/main.go 负责告诉它。
-3. **apiRouter.Group("/api") 是无效代码**：这是代码里的遗留问题，不影响功能运行，
-   因为路由前缀已经由 RegisterHandlersWithOptions 的 BaseURL 参数正确设置了。
-
-依赖关系：
-- 上游调用方：order/main.go 传入 "order" 和路由注册 wrapper 调用本函数。
-- 下游依赖：wrapper 函数由 order/main.go 提供，其中调用 ports.RegisterHandlersWithOptions。
-
-复现检查：
-1. 启动 order 服务后，`curl http://127.0.0.1:8282/api/customer/test/orders` 能到达 handler。
-2. 修改 http-addr 端口后重启，访问旧端口应收到连接拒绝。
-3. 将端口改为已占用端口，gin.Run 返回 error 并触发 panic。
-
-### D. 入口与端口层文件
-
-#### internal/order/main.go
-文件职责：Order 服务入口。
-关键理解：
-1. init 先加载配置，main 再启动服务。
-2. 同时启动 grpc 和 http，说明订单服务双协议暴露。
-3. HTTP 注册走 OpenAPI 生成的 RegisterHandlersWithOptions。
-依赖关系：
-1. 调用 common/server。
-2. 调用 order/ports + order/http。
-3. 调用 order/service.NewApplication（Lesson7 后）。
-复现检查：
-1. 运行后应同时看到 HTTP 与 gRPC 监听日志。
-
-#### internal/order/http.go
-文件职责：Order HTTP 端口适配层。
-关键理解：
-1. 方法名来自 OpenAPI 生成接口，不是手写命名。
-2. 正式业务版本里只做三件事：绑定参数、调用应用层、返回响应。
-依赖关系：
-1. 被 main.go 注册为 HTTP handler。
-2. 调用 app.Commands / app.Queries。
-复现检查：
-1. POST 下单和 GET 查询调用时会进入对应 handler。
-
-#### internal/order/ports/grpc.go
-文件职责：Order gRPC 端口实现。
-关键理解：
-1. 实现 orderpb.OrderServiceServer 接口。
-2. 该阶段多为 TODO，后续 lesson 逐步填充。
-依赖关系：
-1. 被 main.go 的 RegisterOrderServiceServer 注册。
-复现检查：
-1. 若方法未实现，调用对应 RPC 会 panic。
-
-#### internal/stock/main.go
-文件职责：Stock 服务入口。
-关键理解：
-1. 读取 stock.server-to-run 决定跑 grpc 还是 http。
-2. 当前阶段主要跑 grpc。
-依赖关系：
-1. 调用 service.NewApplication。
-2. 调用 common/server.RunGRPCServer。
-复现检查：
-1. 改 server-to-run 值可触发不同分支行为。
-
-#### internal/stock/ports/grpc.go
-文件职责：Stock gRPC 端口实现。
-关键理解：
-1. 需要实现 GetItems 和 CheckIfItemsInStock。
-2. 这层把外部请求转成应用层调用。
-依赖关系：
-1. 被 stock main 注册。
-2. 供 order 的 grpc client 调用。
-复现检查：
-1. 用 grpc 客户端请求库存接口能否返回。
-
-### E. 应用层组装与门面文件
-
-#### internal/stock/app/app.go
-文件职责：Stock 应用层门面（Commands、Queries 聚合）。
-关键理解：
-1. 通过一个 Application struct 对外暴露能力。
-2. 减少上层对内部 handler 细节感知。
-依赖关系：
-1. 被 stock service/application.go 填充。
-2. 被 stock ports 层消费。
-复现检查：
-1. main 中传递 Application 后端口层能拿到功能入口。
-
-#### internal/stock/service/application.go
-文件职责：Stock 的依赖注入与应用组装。
-关键理解：
-1. 这是组合根，负责把 repo、handler 组装成 app.Application。
-2. main 不直接 new 各种依赖，职责更清晰。
-依赖关系：
-1. 被 stock main 调用。
-2. 依赖 stock adapters/domain/app。
-复现检查：
-1. 替换 repo 实现时只改此处注入。
-
-#### internal/order/app/app.go
-文件职责：Order 应用层门面。
-关键理解：
-1. Commands 放写流程（Create/Update）。
-2. Queries 放读流程（GetCustomerOrder）。
-依赖关系：
-1. 被 order/service/application.go 赋值。
-2. 被 HTTP/gRPC 端口调用。
-复现检查：
-1. 端口层只依赖 app.Application 就能调用所有用例。
-
-#### internal/order/service/application.go
-文件职责：Order 依赖注入核心。
-关键理解：
-1. 初始化 orderRepo、logger、metricsClient。
-2. 初始化 stock grpc client 并封装成 StockService。
-3. 把 command/query handler 统一装配到 app.Application。
-依赖关系：
-1. 被 order main 调用。
-2. 依赖 adapters、decorator、query/service 接口。
-复现检查：
-1. 替换 stock 调用实现时，这里改注入即可。
-
-### F. 领域层与仓储实现
-
-#### internal/order/domain/order/order.go
-文件职责：订单领域模型。
-关键理解：
-1. Order 代表核心业务实体。
-2. Items 使用 proto 结构，便于和接口层衔接。
-依赖关系：
-1. 被 command/query/repository 使用。
-复现检查：
-1. 新增字段后需要同步影响创建与查询返回。
-
-#### internal/order/domain/order/repository.go
-文件职责：订单仓储接口定义。
-关键理解：
-1. 定义 Create/Get/Update 三个能力。
-2. NotFoundError 提供业务语义化错误。
-依赖关系：
-1. 被应用层依赖。
-2. 被 adapters 实现。
-复现检查：
-1. 任何 repo 实现都必须满足该接口。
-
-#### internal/order/adapters/order_inmem_repository.go
-文件职责：订单仓储 in-memory 实现。
-关键理解：
-1. RWMutex 保障并发安全。
-2. Create 用时间戳生成 ID 并写入内存切片。
-3. Get 按订单号和客户号匹配。
-4. Update 使用高阶函数 updateFn 控制更新逻辑。
-依赖关系：
-1. 实现 domain.Repository。
-2. 被 order/service/application.go 注入到 handler。
-复现检查：
-1. 连续创建订单，内存 store 长度增加。
-2. 查不存在订单返回 NotFoundError。
-
-#### internal/stock/domain/stock/repository.go
-文件职责：库存仓储接口定义。
-关键理解：
-1. 定义 GetItems 与 CheckIfItemsInStock 等库存查询能力。
-2. 约束 stock 应用层只依赖接口。
-依赖关系：
-1. 被 stock 适配器实现。
-复现检查：
-1. repo 实现变更不影响应用层签名。
-
-#### internal/stock/adapters/stock_inmem_repository.go
-文件职责：库存仓储 in-memory 实现。
-关键理解：
-1. 用内存数据模拟库存系统。
-2. 返回商品详情和库存校验结果。
-依赖关系：
-1. 实现 stock domain repo。
-2. 被 stock service 注入。
-复现检查：
-1. 调整初始库存数据，库存接口结果应变化。
-
-### G. Command / Query 与装饰器
-
-#### internal/common/decorator/query.go
-文件职责：QueryHandler 泛型接口与装饰器装配入口。
-关键理解：
-1. 统一查询处理器形态。
-2. ApplyQueryDecorators 按顺序套日志和指标。
-依赖关系：
-1. 被 query handler 构造函数调用。
-复现检查：
-1. 执行查询时可看到日志与指标上报被触发。
-
-#### internal/common/decorator/command.go
-文件职责：CommandHandler 泛型接口与装饰器装配入口。
-关键理解：
-1. 与 query 对称，作用于写操作。
-2. 保证 command 也拥有统一日志/指标能力。
-依赖关系：
-1. 被 command handler 构造函数调用。
-复现检查：
-1. 下单成功/失败时都有对应统计路径。
-
-#### internal/common/decorator/logging.go
-文件职责：通用日志装饰器。
-关键理解：
-1. 记录 action 名、请求体、执行结果。
-2. 错误场景会打印失败日志，方便排查。
-依赖关系：
-1. 被 query.go / command.go 装配调用。
-复现检查：
-1. 故意触发参数错误，日志中应有失败记录。
-
-#### internal/common/decorator/metrics.go
-文件职责：通用指标装饰器。
-关键理解：
-1. 统计耗时、成功次数、失败次数。
-2. action 名来自请求类型。
-依赖关系：
-1. 依赖 MetricsClient 接口。
-2. 被 query/command 装配。
-复现检查：
-1. 打断点可看到 Inc 被调用三类 key。
-
-#### internal/common/metrics/todo_metrics.go
-文件职责：指标客户端占位实现。
-关键理解：
-1. 当前 Inc 空实现，仅保证接口连通。
-2. 后续可替换为 Prometheus/StatsD 真正上报。
-依赖关系：
-1. 被 application.go 注入到装饰器。
-复现检查：
-1. 替换 Inc 内容后可观察指标输出变化。
-
-#### internal/order/app/query/get_customer_order.go
-文件职责：查询订单用例。
-关键理解：
-1. Handle 只做一件事：调用 repo.Get。
-2. 构造函数里应用了 query decorators。
-依赖关系：
-1. 被 HTTP Get handler 调用。
-2. 依赖 order repository 接口。
-复现检查：
-1. 用存在/不存在订单号分别调用，观察成功与错误路径。
-
-#### internal/order/app/query/service.go
-文件职责：定义 Order 侧需要的库存服务端口接口。
-关键理解：
-1. 这是“防腐层端口”，隔离外部 grpc 客户端细节。
-2. create_order 依赖该接口而非具体实现。
-依赖关系：
-1. 被 create_order.go 依赖。
-2. 被 adapters/grpc/stock_grpc.go 实现。
-复现检查：
-1. 用 mock 实现该接口可做单元测试。
-
-#### internal/order/app/command/create_order.go
-文件职责：创建订单用例。
-关键理解：
-1. validate 先检查 items 非空。
-2. packItems 合并重复商品数量，减少重复校验。
-3. 调用 StockService.CheckIfItemsInStock 进行库存校验。
-4. 校验通过后调用 orderRepo.Create。
-依赖关系：
-1. 被 HTTP POST handler 调用。
-2. 依赖 order repo + stock service 接口。
-复现检查：
-1. 传空 items 应报错。
-2. 传重复 item id 时会被合并后再请求库存服务。
-
-#### internal/order/app/command/update_order.go
-文件职责：更新订单用例。
-关键理解：
-1. UpdateFn 由上层传入，决定如何改订单。
-2. handler 只负责流程控制和仓储调用。
-依赖关系：
-1. 被上层更新场景调用。
-2. 依赖 order repo.Update。
-复现检查：
-1. UpdateFn 为空时走默认 no-op 逻辑。
-
-#### internal/order/adapters/grpc/stock_grpc.go
-文件职责：StockService 的 gRPC 适配器实现。
-关键理解：
-1. 封装 stockpb.StockServiceClient。
-2. 把领域需要的方法映射为具体 RPC 请求。
-依赖关系：
-1. 实现 order/app/query/service.go 的 StockService。
-2. 被 order/service/application.go 注入。
-复现检查：
-1. 断开 stock 服务后，下单库存校验应报远程调用错误。
-
-### H. 模块与工具配置文件
-
-#### internal/common/go.mod
-文件职责：common 模块依赖声明。
-关键理解：
-1. require 指明编译期直接依赖。
-2. indirect 是传递依赖。
-依赖关系：
-1. 与 common/go.sum 配套。
-复现检查：
-1. go mod tidy 后依赖变化会更新。
-
-#### internal/order/go.mod
-文件职责：order 模块依赖声明。
-关键理解：
-1. 包含 gin、grpc、viper 等运行依赖。
-2. 通过 replace（若存在）可指向本地模块。
-依赖关系：
-1. 与 order/go.sum 配套。
-复现检查：
-1. 新增 import 后 go mod tidy 会增补依赖。
-
-#### internal/stock/go.mod
-文件职责：stock 模块依赖声明。
-关键理解：
-1. 与 order 模块类似，但按 stock 需求维护。
-依赖关系：
-1. 与 stock/go.sum 配套。
-复现检查：
-1. 删除必要依赖后编译会失败。
-
-#### internal/common/go.sum
-文件职责：common 模块依赖校验和锁文件。
-关键理解：
-1. 每行是某个模块版本的 hash，不是业务代码。
-2. 用于防篡改和可重复构建。
-依赖关系：
-1. 被 Go module 下载与校验流程读取。
-复现检查：
-1. 删除后重新 tidy 会再生成。
-
-#### internal/order/go.sum
-文件职责：order 模块依赖校验和。
-关键理解：
-1. 本质与 common/go.sum 相同。
-依赖关系：
-1. 由 go 命令自动维护。
-复现检查：
-1. 拉新依赖后会新增对应 checksum 行。
-
-#### internal/stock/go.sum
-文件职责：stock 模块依赖校验和。
-关键理解：
-1. 本质与其他 go.sum 相同。
-依赖关系：
-1. 由 go 命令自动维护。
-复现检查：
-1. tidy 后内容变化属于正常现象。
-
-#### internal/stock/.air.toml
-文件职责：stock 服务本地热重载配置。
-关键理解：
-1. 指定监听目录、忽略目录、重建命令。
-2. 开发时改代码自动重启。
-依赖关系：
-1. 被 air 工具读取。
-复现检查：
-1. 运行 air 后修改 go 文件观察自动重编译。
-
-#### internal/order/.air.toml
-文件职责：order 服务热重载配置。
-关键理解：
-1. 与 stock/.air.toml 同类。
-依赖关系：
-1. 被 air 工具读取。
-复现检查：
-1. 修改 order 代码后自动重启。
-
-#### internal/kitchen/.air.toml
-文件职责：kitchen 服务热重载配置。
-关键理解：
-1. 给 kitchen 模块开发阶段提速。
-依赖关系：
-1. 被 air 工具读取。
-复现检查：
-1. 进入 kitchen 模块运行 air 验证。
-
-#### internal/payment/.air.toml
-文件职责：payment 服务热重载配置。
-关键理解：
-1. 给 payment 模块开发阶段提速。
-依赖关系：
-1. 被 air 工具读取。
-复现检查：
-1. 进入 payment 模块运行 air 验证。
-
----
-
-## 3. 复现步骤（按 Lesson5-11）
-
-1. 启动 stock 服务（grpc）。
-2. 启动 order 服务（http + grpc）。
-3. 调用 order 下单 HTTP 接口。
-4. order 的 create_order handler 内会通过 StockService 走 grpc 校验库存。
-5. 校验通过后写入 order in-memory repo。
-6. 再调用查询接口，验证能读到订单。
-
-如果你复现失败，优先按这个顺序查：
-1. 配置地址是否一致（global.yaml）。
-2. stock 服务是否先启动。
-3. order/service/application.go 是否成功注入 stock grpc client。
-4. proto 生成代码是否和 proto 契约一致。
-
----
-
-## 4. Lesson5-11 的“必须会改”文件清单
-
-你练习时，优先动这些文件：
-1. internal/order/http.go
-2. internal/order/app/command/create_order.go
-3. internal/order/app/query/get_customer_order.go
-4. internal/order/adapters/order_inmem_repository.go
-5. internal/order/service/application.go
-6. internal/order/adapters/grpc/stock_grpc.go
-7. internal/stock/ports/grpc.go
-
-这些文件改熟了，后面的数据库、消息队列、可观测性才容易接上。
-
----
-
-## 5. 下一批输出方式（你确认后继续）
-
-下一批我按同样模板写 Lesson12-Lesson20，并继续提交到 Git：
-1. 每个文件都有“职责 + 关键代码 + 调用关系 + 复现检查”。
-2. 对生成文件/go.sum 继续采用“用途级详细说明”，不做低价值逐行复述。
diff --git a/internal/common/client/grpc.go b/internal/common/client/grpc.go
new file mode 100644
index 0000000..f9f514a
--- /dev/null
+++ b/internal/common/client/grpc.go
@@ -0,0 +1,29 @@
+package client
+
+import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
+	"github.com/spf13/viper"
+	"google.golang.org/grpc"
+	"google.golang.org/grpc/credentials/insecure"
+)
+
+func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
+	grpcAddr := viper.GetString("stock.grpc-addr")
+	opts, err := grpcDialOpts(grpcAddr)
+	if err != nil {
+		return nil, func() error { return nil }, err
+	}
+	conn, err := grpc.NewClient(grpcAddr, opts...)
+	if err != nil {
+		return nil, func() error { return nil }, err
+	}
+	return stockpb.NewStockServiceClient(conn), conn.Close, nil
+}
+
+func grpcDialOpts(addr string) ([]grpc.DialOption, error) {
+	return []grpc.DialOption{
+		grpc.WithTransportCredentials(insecure.NewCredentials()),
+	}, nil
+}
diff --git a/internal/common/config/viper.go b/internal/common/config/viper.go
index 5c2fe30..103d697 100644
--- a/internal/common/config/viper.go
+++ b/internal/common/config/viper.go
@@ -2,29 +2,10 @@ package config
 
 import "github.com/spf13/viper"
 
-// NewViperConfig 统一加载全局配置文件到内存。
-// 这个函数在 init() 里被调用，确保 main 启动前配置已经就位，
-// 后续各个函数可以直接用 viper.Get... 获取配置，而不用每次都读文件。
 func NewViperConfig() error {
-	// SetConfigName 指定配置文件名（不包括扩展名）。
-	// 这里是 "global"，对应 global.yaml。
 	viper.SetConfigName("global")
-
-	// SetConfigType 告诉 viper 配置格式是什么。
-	// viper 支持 json/yaml/toml/hcl 等多种格式。
 	viper.SetConfigType("yaml")
-
-	// AddConfigPath 告诉 viper 去哪个目录找配置文件。
-	// "../common/config" 是相对当前执行目录的路径，
-	// 假设你从项目根目录启动 order 服务，viper 会去 internal/common/config 目录找 global.yaml。
 	viper.AddConfigPath("../common/config")
-
-	// AutomaticEnv 让环境变量自动覆盖配置文件中的值。
-	// 举例：如果环境里设置了 ORDER_GRPC_ADDR=":6666"，会覆盖 yaml 里的 order.grpc-addr。
-	// 这样本地开发、测试、生产环境可以动态调整配置，而不改代码。
 	viper.AutomaticEnv()
-
-	// ReadInConfig 真正去磁盘读配置文件，把内容解析到 viper 的内存结构里。
-	// 如果文件不存在或格式有问题，这里会返回 error。
 	return viper.ReadInConfig()
 }
diff --git a/internal/common/decorator/command.go b/internal/common/decorator/command.go
index e820637..53f8c37 100644
--- a/internal/common/decorator/command.go
+++ b/internal/common/decorator/command.go
@@ -6,13 +6,10 @@ import (
 	"github.com/sirupsen/logrus"
 )
 
-// CommandHandler 和 QueryHandler 对称，只是语义上代表“会改状态”的操作。
 type CommandHandler[C, R any] interface {
 	Handle(ctx context.Context, cmd C) (R, error)
 }
 
-// ApplyCommandDecorators 让所有 command 统一拥有日志和指标能力。
-// 这样新增一个写操作时，只要实现业务 handler，不用重复写埋点代码。
 func ApplyCommandDecorators[C, R any](handler CommandHandler[C, R], logger *logrus.Entry, metricsClient MetricsClient) CommandHandler[C, R] {
 	return queryLoggingDecorator[C, R]{
 		logger: logger,
diff --git a/internal/common/decorator/logging.go b/internal/common/decorator/logging.go
index 2b158ae..a9a2325 100644
--- a/internal/common/decorator/logging.go
+++ b/internal/common/decorator/logging.go
@@ -8,15 +8,12 @@ import (
 	"github.com/sirupsen/logrus"
 )
 
-// queryLoggingDecorator 在真实 handler 外层补一圈日志。
-// “Decorator” 不是框架魔法，本质就是包一层结构体再转调 base。
 type queryLoggingDecorator[C, R any] struct {
 	logger *logrus.Entry
 	base   QueryHandler[C, R]
 }
 
 func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
-	// generateActionName 会把类型名拿出来，作为统一的日志 action 名。
 	logger := q.logger.WithFields(logrus.Fields{
 		"query":      generateActionName(cmd),
 		"query_body": fmt.Sprintf("%#v", cmd),
@@ -32,7 +29,6 @@ func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result
 	return q.base.Handle(ctx, cmd)
 }
 
-// generateActionName 直接复用请求类型名，避免每个 handler 手动写一次 action 常量。
 func generateActionName(cmd any) string {
 	return strings.Split(fmt.Sprintf("%T", cmd), ".")[1]
 }
diff --git a/internal/common/decorator/metrics.go b/internal/common/decorator/metrics.go
index ba04887..fba9a77 100644
--- a/internal/common/decorator/metrics.go
+++ b/internal/common/decorator/metrics.go
@@ -7,20 +7,16 @@ import (
 	"time"
 )
 
-// MetricsClient 是一个很小的抽象，方便先接假实现，后续再切到真实监控系统。
 type MetricsClient interface {
 	Inc(key string, value int)
 }
 
-// queryMetricsDecorator 统计耗时、成功、失败次数。
-// 它和日志装饰器可以自由组合，因为两者都只依赖同一个 handler 接口。
 type queryMetricsDecorator[C, R any] struct {
 	base   QueryHandler[C, R]
 	client MetricsClient
 }
 
 func (q queryMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
-	// 用 defer 统一收尾，确保无论成功还是失败都能上报指标。
 	start := time.Now()
 	actionName := strings.ToLower(generateActionName(cmd))
 	defer func() {
diff --git a/internal/common/decorator/query.go b/internal/common/decorator/query.go
index b5e9cf0..d3848fe 100644
--- a/internal/common/decorator/query.go
+++ b/internal/common/decorator/query.go
@@ -8,13 +8,10 @@ import (
 
 // QueryHandler defines a generic type that receives a Query Q,
 // and returns a result R
-// 泛型接口的好处是：同一套日志/指标装饰逻辑可以复用在不同查询上。
 type QueryHandler[Q, R any] interface {
 	Handle(ctx context.Context, query Q) (R, error)
 }
 
-// ApplyQueryDecorators 按固定顺序把日志和指标能力包在真实 handler 外面。
-// 这就是装饰器模式：不改业务代码，也能统一加横切能力。
 func ApplyQueryDecorators[H, R any](handler QueryHandler[H, R], logger *logrus.Entry, metricsClient MetricsClient) QueryHandler[H, R] {
 	return queryLoggingDecorator[H, R]{
 		logger: logger,
diff --git a/internal/common/discovery/consul/consul.go b/internal/common/discovery/consul/consul.go
new file mode 100644
index 0000000..b0ef77f
--- /dev/null
+++ b/internal/common/discovery/consul/consul.go
@@ -0,0 +1,88 @@
+package consul
+
+import (
+	"context"
+	"errors"
+	"fmt"
+	"strconv"
+	"strings"
+	"sync"
+
+	"github.com/hashicorp/consul/api"
+	"github.com/sirupsen/logrus"
+)
+
+type Registry struct {
+	client *api.Client
+}
+
+var (
+	consulClient *Registry
+	once         sync.Once
+	initErr      error
+)
+
+func New(consulAddr string) (*Registry, error) {
+	once.Do(func() {
+		config := api.DefaultConfig()
+		config.Address = consulAddr
+		client, err := api.NewClient(config)
+		if err != nil {
+			initErr = err
+			return
+		}
+		consulClient = &Registry{
+			client: client,
+		}
+	})
+	if initErr != nil {
+		return nil, initErr
+	}
+	return consulClient, nil
+}
+
+func (r *Registry) Register(_ context.Context, instanceID, serviceName, hostPort string) error {
+	parts := strings.Split(hostPort, ":")
+	if len(parts) != 2 {
+		return errors.New("invalid host:port format")
+	}
+	host := parts[0]
+	port, _ := strconv.Atoi(parts[1])
+	return r.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
+		ID:      instanceID,
+		Address: host,
+		Port:    port,
+		Name:    serviceName,
+		Check: &api.AgentServiceCheck{
+			CheckID:                        instanceID,
+			TLSSkipVerify:                  false,
+			TTL:                            "5s",
+			Timeout:                        "5s",
+			DeregisterCriticalServiceAfter: "10s",
+		},
+	})
+}
+
+func (r *Registry) Deregister(_ context.Context, instanceID, serviceName string) error {
+	logrus.WithFields(logrus.Fields{
+		"instanceID":  instanceID,
+		"serviceName": serviceName,
+	}).Info("deregister from consul")
+	return r.client.Agent().CheckDeregister(instanceID)
+}
+
+func (r *Registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
+	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
+	if err != nil {
+		return nil, err
+	}
+	var ips []string
+	for _, e := range entries {
+		ips = append(ips, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
+	}
+	return ips, nil
+}
+
+func (r *Registry) HealthCheck(instanceID, serviceName string) error {
+	return r.client.Agent().UpdateTTL(instanceID, "online", api.HealthPassing)
+}
diff --git a/internal/common/discovery/discovery.go b/internal/common/discovery/discovery.go
new file mode 100644
index 0000000..7cf6f58
--- /dev/null
+++ b/internal/common/discovery/discovery.go
@@ -0,0 +1,20 @@
+package discovery
+
+import (
+	"context"
+	"fmt"
+	"math/rand"
+	"time"
+)
+
+type Registry interface {
+	Register(ctx context.Context, instanceID, serviceName, hostPort string) error
+	Deregister(ctx context.Context, instanceID, serviceName string) error
+	Discover(ctx context.Context, serviceName string) ([]string, error)
+	HealthCheck(instanceID, serviceName string) error
+}
+
+func GenerateInstanceID(serviceName string) string {
+	x := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
+	return fmt.Sprintf("%s-%d", serviceName, x)
+}
diff --git a/internal/common/discovery/grpc.go b/internal/common/discovery/grpc.go
new file mode 100644
index 0000000..cda522a
--- /dev/null
+++ b/internal/common/discovery/grpc.go
@@ -0,0 +1,37 @@
+package discovery
+
+import (
+	"context"
+	"time"
+
+	"github.com/ghost-yu/go_shop_second/common/discovery/consul"
+	"github.com/sirupsen/logrus"
+	"github.com/spf13/viper"
+)
+
+func RegisterToConsul(ctx context.Context, serviceName string) (func() error, error) {
+	registry, err := consul.New(viper.GetString("consul.addr"))
+	if err != nil {
+		return func() error { return nil }, err
+	}
+	instanceID := GenerateInstanceID(serviceName)
+	grpcAddr := viper.Sub(serviceName).GetString("grpc-addr")
+	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
+		return func() error { return nil }, err
+	}
+	go func() {
+		for {
+			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
+				logrus.Panicf("no heartbeat from %s to registry, err=%v", serviceName, err)
+			}
+			time.Sleep(1 * time.Second)
+		}
+	}()
+	logrus.WithFields(logrus.Fields{
+		"serviceName": serviceName,
+		"addr":        grpcAddr,
+	}).Info("registered to consul")
+	return func() error {
+		return registry.Deregister(ctx, instanceID, serviceName)
+	}, nil
+}
diff --git a/internal/common/go.mod b/internal/common/go.mod
index 8355ff8..440cdc8 100644
--- a/internal/common/go.mod
+++ b/internal/common/go.mod
@@ -5,18 +5,21 @@ go 1.22.8
 require (
 	github.com/gin-gonic/gin v1.9.1
 	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
+	github.com/hashicorp/consul/api v1.28.2
 	github.com/oapi-codegen/runtime v1.1.1
 	github.com/sirupsen/logrus v1.8.1
 	github.com/spf13/viper v1.19.0
-	google.golang.org/grpc v1.62.1
-	google.golang.org/protobuf v1.33.0
+	google.golang.org/grpc v1.67.1
+	google.golang.org/protobuf v1.35.1
 )
 
 require (
 	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
+	github.com/armon/go-metrics v0.4.1 // indirect
 	github.com/bytedance/sonic v1.10.0-rc3 // indirect
 	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
 	github.com/chenzhuoyu/iasm v0.9.0 // indirect
+	github.com/fatih/color v1.14.1 // indirect
 	github.com/fsnotify/fsnotify v1.7.0 // indirect
 	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
 	github.com/gin-contrib/sse v0.1.0 // indirect
@@ -26,12 +29,22 @@ require (
 	github.com/goccy/go-json v0.10.2 // indirect
 	github.com/golang/protobuf v1.5.3 // indirect
 	github.com/google/uuid v1.6.0 // indirect
+	github.com/hashicorp/errwrap v1.1.0 // indirect
+	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
+	github.com/hashicorp/go-hclog v1.5.0 // indirect
+	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
+	github.com/hashicorp/go-multierror v1.1.1 // indirect
+	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
+	github.com/hashicorp/golang-lru v0.5.4 // indirect
 	github.com/hashicorp/hcl v1.0.0 // indirect
+	github.com/hashicorp/serf v0.10.1 // indirect
 	github.com/json-iterator/go v1.1.12 // indirect
 	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
 	github.com/leodido/go-urn v1.2.4 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
+	github.com/mattn/go-colorable v0.1.13 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
+	github.com/mitchellh/go-homedir v1.1.0 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
 	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
 	github.com/modern-go/reflect2 v1.0.2 // indirect
@@ -48,12 +61,12 @@ require (
 	go.uber.org/atomic v1.9.0 // indirect
 	go.uber.org/multierr v1.9.0 // indirect
 	golang.org/x/arch v0.4.0 // indirect
-	golang.org/x/crypto v0.21.0 // indirect
+	golang.org/x/crypto v0.26.0 // indirect
 	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
-	golang.org/x/net v0.23.0 // indirect
-	golang.org/x/sys v0.18.0 // indirect
-	golang.org/x/text v0.14.0 // indirect
-	google.golang.org/genproto/googleapis/rpc v0.0.0-20240314234333-6e1732d8331c // indirect
+	golang.org/x/net v0.28.0 // indirect
+	golang.org/x/sys v0.24.0 // indirect
+	golang.org/x/text v0.17.0 // indirect
+	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
 	gopkg.in/ini.v1 v1.67.0 // indirect
 	gopkg.in/yaml.v3 v3.0.1 // indirect
 )
diff --git a/internal/common/go.sum b/internal/common/go.sum
index a404a00..bfe7af2 100644
--- a/internal/common/go.sum
+++ b/internal/common/go.sum
@@ -1,21 +1,39 @@
 cloud.google.com/go v0.26.0/go.mod h1:aQUYkXzVsufM+DwF1aE+0xfcU+56JwCaLick0ClmMTw=
 github.com/BurntSushi/toml v0.3.1/go.mod h1:xHWCNGjB5oqiDr8zfno3MHue2Ht5sIBksp03qcyfWMU=
+github.com/DataDog/datadog-go v3.2.0+incompatible/go.mod h1:LButxg5PwREeZtORoXG3tL4fMGNddJ+vMq1mwgfaqoQ=
 github.com/RaveNoX/go-jsoncommentstrip v1.0.0/go.mod h1:78ihd09MekBnJnxpICcwzCMzGrKSKYe4AqU6PDYYpjk=
+github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
+github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
+github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
+github.com/alecthomas/units v0.0.0-20190717042225-c3de453c63f4/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
 github.com/apapsch/go-jsonmerge/v2 v2.0.0 h1:axGnT1gRIfimI7gJifB699GoE/oq+F2MU7Dml6nw9rQ=
 github.com/apapsch/go-jsonmerge/v2 v2.0.0/go.mod h1:lvDnEdqiQrp0O42VQGgmlKpxL1AP2+08jFMw88y4klk=
+github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e/go.mod h1:3U/XgcO3hCbHZ8TKRvWD2dDTCfh9M9ya+I9JpbB7O8o=
+github.com/armon/go-metrics v0.0.0-20180917152333-f0300d1749da/go.mod h1:Q73ZrmVTwzkszR9V5SSuryQ31EELlFMUz1kKyl939pY=
+github.com/armon/go-metrics v0.4.1 h1:hR91U9KYmb6bLBYLQjyM+3j+rcd/UhE+G78SFnF8gJA=
+github.com/armon/go-metrics v0.4.1/go.mod h1:E6amYzXo6aW1tqzoZGT755KkbgrJsSdpwZ+3JqfkOG4=
+github.com/armon/go-radix v0.0.0-20180808171621-7fddfc383310/go.mod h1:ufUuZ+zHj4x4TnLV4JWEpy2hxWSpsRywHrMgIH9cCH8=
+github.com/armon/go-radix v1.0.0/go.mod h1:ufUuZ+zHj4x4TnLV4JWEpy2hxWSpsRywHrMgIH9cCH8=
 github.com/benbjohnson/clock v1.1.0/go.mod h1:J11/hYXuz8f4ySSvYwY0FKfm+ezbsZBKZxNJlLklBHA=
+github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973/go.mod h1:Dwedo/Wpr24TaqPxmxbtue+5NUziq4I4S80YR8gNf3Q=
+github.com/beorn7/perks v1.0.0/go.mod h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+CedLV8=
+github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
+github.com/bgentry/speakeasy v0.1.0/go.mod h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=
 github.com/bmatcuk/doublestar v1.1.1/go.mod h1:UD6OnuiIn0yFxxA2le/rnRU1G4RaI4UvFv1sNto9p6w=
 github.com/bytedance/sonic v1.5.0/go.mod h1:ED5hyg4y6t3/9Ku1R6dU/4KyJ48DZ4jPhfY1O2AihPM=
 github.com/bytedance/sonic v1.10.0-rc/go.mod h1:ElCzW+ufi8qKqNW0FY314xriJhyJhuoJ3gFZdAHF7NM=
 github.com/bytedance/sonic v1.10.0-rc3 h1:uNSnscRapXTwUgTyOF0GVljYD08p9X/Lbr9MweSV3V0=
 github.com/bytedance/sonic v1.10.0-rc3/go.mod h1:iZcSUejdk5aukTND/Eu/ivjQuEL0Cu9/rf50Hi0u/g4=
 github.com/census-instrumentation/opencensus-proto v0.2.1/go.mod h1:f6KPmirojxKA12rnyqOA5BBL4O983OfeGPqjHWSTneU=
+github.com/cespare/xxhash/v2 v2.1.1/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
 github.com/chenzhuoyu/base64x v0.0.0-20211019084208-fb5309c8db06/go.mod h1:DH46F32mSOjUmXrMHnKwZdA8wcEefY7UVqBKYGjpdQY=
 github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311/go.mod h1:b583jCggY9gE99b6G5LEC39OIiVsWj+R97kbl5odCEk=
 github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d h1:77cEq6EriyTZ0g/qfRdp61a3Uu/AWrgIq2s0ClJV1g0=
 github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d/go.mod h1:8EPpVsBuRksnlj1mLy4AWzRNQYxauNi62uWcE3to6eA=
 github.com/chenzhuoyu/iasm v0.9.0 h1:9fhXjVzq5hUy2gkhhgHl95zG2cEAhw9OSGs8toWWAwo=
 github.com/chenzhuoyu/iasm v0.9.0/go.mod h1:Xjy2NpN3h7aUqeqM+woSuuvxmIe6+DDsiNLIrkAmYog=
+github.com/circonus-labs/circonus-gometrics v2.3.1+incompatible/go.mod h1:nmEj6Dob7S7YxXgwXpfOuvO54S+tGdZdw9fuRZt25Ag=
+github.com/circonus-labs/circonusllhist v0.1.3/go.mod h1:kMXHVDlOchFAehlya5ePtbp5jckzBHf4XRpQvBOLI+I=
 github.com/client9/misspell v0.3.4/go.mod h1:qj6jICC3Q7zFZvVWo7KLAzC3yx5G7kyvSDkc90ppPyw=
 github.com/cncf/udpa/go v0.0.0-20191209042840-269d4d468f6f/go.mod h1:M8M6+tZqaGXZJjfX53e64911xZQV5JYwmTeXPW+k8Sc=
 github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
@@ -26,6 +44,11 @@ github.com/envoyproxy/go-control-plane v0.9.0/go.mod h1:YTl/9mNaCwkRvm6d1a2C3ymF
 github.com/envoyproxy/go-control-plane v0.9.1-0.20191026205805-5f8ba28d4473/go.mod h1:YTl/9mNaCwkRvm6d1a2C3ymFceY/DCBVvsKhRF0iEA4=
 github.com/envoyproxy/go-control-plane v0.9.4/go.mod h1:6rpuAdCZL397s3pYoYcLgu1mIlRU8Am5FuJP05cCM98=
 github.com/envoyproxy/protoc-gen-validate v0.1.0/go.mod h1:iSmxcyjqTsJpI2R4NaDN7+kN2VEUnK/pcBlmesArF7c=
+github.com/fatih/color v1.7.0/go.mod h1:Zm6kSWBoL9eyXnKyktHP6abPY2pDugNf5KwzbycvMj4=
+github.com/fatih/color v1.9.0/go.mod h1:eQcE1qtQxscV5RaZvpXrrb8Drkc3/DdQ+uUYCNjL+zU=
+github.com/fatih/color v1.13.0/go.mod h1:kLAiJbzzSOZDVNGyDpeOxJ47H46qBXwg5ILebYFFOfk=
+github.com/fatih/color v1.14.1 h1:qfhVLaG5s+nCROl1zJsZRxFeYrHLqWroPOQ8BWiNb4w=
+github.com/fatih/color v1.14.1/go.mod h1:2oHN61fhTpgcxD3TSWCgKDiH1+x4OiDVVGH8WlgGZGg=
 github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHkI4W8=
 github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
@@ -36,7 +59,11 @@ github.com/gin-contrib/sse v0.1.0 h1:Y/yl/+YNO8GZSjAhjMsSuLt29uWRFHdHYUb5lYOV9qE
 github.com/gin-contrib/sse v0.1.0/go.mod h1:RHrZQHXnP2xjPF+u1gW/2HnVO7nvIa9PG3Gm+fLHvGI=
 github.com/gin-gonic/gin v1.9.1 h1:4idEAncQnU5cB7BeOkPtxjfCSye0AAm1R0RVIqJ+Jmg=
 github.com/gin-gonic/gin v1.9.1/go.mod h1:hPrL7YrpYKXt5YId3A/Tnip5kqbEAP+KLuI3SUcPTeU=
+github.com/go-kit/kit v0.8.0/go.mod h1:xBxKIO96dXMWWy0MnWVtmwkA9/13aqxPnvrjFYMA2as=
+github.com/go-kit/kit v0.9.0/go.mod h1:xBxKIO96dXMWWy0MnWVtmwkA9/13aqxPnvrjFYMA2as=
 github.com/go-kit/log v0.1.0/go.mod h1:zbhenjAZHb184qTLMA9ZjW7ThYL0H2mk7Q6pNt4vbaY=
+github.com/go-logfmt/logfmt v0.3.0/go.mod h1:Qt1PoO58o5twSAckw1HlFXLmHsOX5/0LbT9GBnD5lWE=
+github.com/go-logfmt/logfmt v0.4.0/go.mod h1:3RMwSq7FuexP4Kalkev3ejPJsZTpXXBr9+V4qmtdjCk=
 github.com/go-logfmt/logfmt v0.5.0/go.mod h1:wCYkCAKZfumFQihp8CzCvQ3paCTfi41vtzG1KdI/P7A=
 github.com/go-playground/assert/v2 v2.2.0 h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=
 github.com/go-playground/assert/v2 v2.2.0/go.mod h1:VDjEfimB/XKnb+ZQfWdccd7VUvScMdVu0Titje2rxJ4=
@@ -49,17 +76,24 @@ github.com/go-playground/validator/v10 v10.14.1/go.mod h1:9iXMNT7sEkjXb0I+enO7QX
 github.com/go-stack/stack v1.8.0/go.mod h1:v0f6uXyyMGvRgIKkXu+yp6POWl0qKG85gN/melR3HDY=
 github.com/goccy/go-json v0.10.2 h1:CrxCmQqYDkv1z7lO7Wbh2HN93uovUHgrECaO5ZrCXAU=
 github.com/goccy/go-json v0.10.2/go.mod h1:6MelG93GURQebXPDq3khkgXZkazVtN9CRI+MGFi0w8I=
+github.com/gogo/protobuf v1.1.1/go.mod h1:r8qH/GZQm5c6nD/R0oafs1akxWv10x8SbQlK7atdtwQ=
 github.com/gogo/protobuf v1.3.2 h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=
 github.com/gogo/protobuf v1.3.2/go.mod h1:P1XiOD3dCwIKUDQYPy72D8LYyHL2YPYrpS2s69NZV8Q=
 github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b/go.mod h1:SBH7ygxi8pfUlaOkMMuAQtPIUF8ecWP5IEl/CR7VP2Q=
 github.com/golang/mock v1.1.1/go.mod h1:oTYuIxOrZwtPieC+H1uAHpcLFnEyAGVDL/k47Jfbm0A=
 github.com/golang/protobuf v1.2.0/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5yJMmIC1U=
+github.com/golang/protobuf v1.3.1/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5yJMmIC1U=
 github.com/golang/protobuf v1.3.2/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5yJMmIC1U=
 github.com/golang/protobuf v1.3.3/go.mod h1:vzj43D7+SQXF/4pzW/hwtAqwc6iTitCiVSaWz5lYuqw=
 github.com/golang/protobuf v1.5.0/go.mod h1:FsONVRAS9T7sI+LIUmWTfcYkHO4aIWwzhcaSAoJOfIk=
 github.com/golang/protobuf v1.5.3 h1:KhyjKVUg7Usr/dYsdSqoFveMYd5ko72D+zANwlG1mmg=
 github.com/golang/protobuf v1.5.3/go.mod h1:XVQd3VNwM+JqD3oG2Ue2ip4fOMUkwXdXDdiuN0vRsmY=
+github.com/google/btree v0.0.0-20180813153112-4030bb1f1f0c/go.mod h1:lNA+9X1NB3Zf8V7Ke586lFgjr2dZNuvo3lPJSGZ5JPQ=
+github.com/google/btree v1.0.1 h1:gK4Kx5IaGY9CD5sPJ36FHiBJ6ZXl0kilRiiCj+jdYp4=
+github.com/google/btree v1.0.1/go.mod h1:xXMiIv4Fb/0kKde4SpL7qlzvu5cMJDRkFDxJfI9uaxA=
 github.com/google/go-cmp v0.2.0/go.mod h1:oXzfMopK8JAjlY9xF4vHSVASa0yLyX7SntLO5aqRK0M=
+github.com/google/go-cmp v0.3.1/go.mod h1:8QqcDgzrUqlUb/G2PQTWiueGozuR1884gddMywk6iLU=
+github.com/google/go-cmp v0.4.0/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
 github.com/google/go-cmp v0.5.5/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
 github.com/google/go-cmp v0.6.0 h1:ofyhxvXcZhMsU5ulbFiLKl/XBFqE1GSq7atu8tAmTRI=
 github.com/google/go-cmp v0.6.0/go.mod h1:17dUlkBOakJ0+DkrSSNjCkIjxS6bF9zb3elmeNGIjoY=
@@ -68,11 +102,58 @@ github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
 github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
 github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 h1:UH//fgunKIs4JdUbpDl1VZCDaL56wXCB/5+wF6uHfaI=
 github.com/grpc-ecosystem/go-grpc-middleware v1.4.0/go.mod h1:g5qyo/la0ALbONm6Vbp88Yd8NsDy6rZz+RcrMPxvld8=
+github.com/hashicorp/consul/api v1.28.2 h1:mXfkRHrpHN4YY3RqL09nXU1eHKLNiuAN4kHvDQ16k/8=
+github.com/hashicorp/consul/api v1.28.2/go.mod h1:KyzqzgMEya+IZPcD65YFoOVAgPpbfERu4I/tzG6/ueE=
+github.com/hashicorp/consul/sdk v0.16.0 h1:SE9m0W6DEfgIVCJX7xU+iv/hUl4m/nxqMTnCdMxDpJ8=
+github.com/hashicorp/consul/sdk v0.16.0/go.mod h1:7pxqqhqoaPqnBnzXD1StKed62LqJeClzVsUEy85Zr0A=
+github.com/hashicorp/errwrap v1.0.0/go.mod h1:YH+1FKiLXxHSkmPseP+kNlulaMuP3n2brvKWEqk/Jc4=
+github.com/hashicorp/errwrap v1.1.0 h1:OxrOeh75EUXMY8TBjag2fzXGZ40LB6IKw45YeGUDY2I=
+github.com/hashicorp/errwrap v1.1.0/go.mod h1:YH+1FKiLXxHSkmPseP+kNlulaMuP3n2brvKWEqk/Jc4=
+github.com/hashicorp/go-cleanhttp v0.5.0/go.mod h1:JpRdi6/HCYpAwUzNwuwqhbovhLtngrth3wmdIIUrZ80=
+github.com/hashicorp/go-cleanhttp v0.5.2 h1:035FKYIWjmULyFRBKPs8TBQoi0x6d9G4xc9neXJWAZQ=
+github.com/hashicorp/go-cleanhttp v0.5.2/go.mod h1:kO/YDlP8L1346E6Sodw+PrpBSV4/SoxCXGY6BqNFT48=
+github.com/hashicorp/go-hclog v1.5.0 h1:bI2ocEMgcVlz55Oj1xZNBsVi900c7II+fWDyV9o+13c=
+github.com/hashicorp/go-hclog v1.5.0/go.mod h1:W4Qnvbt70Wk/zYJryRzDRU/4r0kIg0PVHBcfoyhpF5M=
+github.com/hashicorp/go-immutable-radix v1.0.0/go.mod h1:0y9vanUI8NX6FsYoO3zeMjhV/C5i9g4Q3DwcSNZ4P60=
+github.com/hashicorp/go-immutable-radix v1.3.1 h1:DKHmCUm2hRBK510BaiZlwvpD40f8bJFeZnpfm2KLowc=
+github.com/hashicorp/go-immutable-radix v1.3.1/go.mod h1:0y9vanUI8NX6FsYoO3zeMjhV/C5i9g4Q3DwcSNZ4P60=
+github.com/hashicorp/go-msgpack v0.5.3/go.mod h1:ahLV/dePpqEmjfWmKiqvPkv/twdG7iPBM1vqhUKIvfM=
+github.com/hashicorp/go-msgpack v0.5.5 h1:i9R9JSrqIz0QVLz3sz+i3YJdT7TTSLcfLLzJi9aZTuI=
+github.com/hashicorp/go-msgpack v0.5.5/go.mod h1:ahLV/dePpqEmjfWmKiqvPkv/twdG7iPBM1vqhUKIvfM=
+github.com/hashicorp/go-multierror v1.0.0/go.mod h1:dHtQlpGsu+cZNNAkkCN/P3hoUDHhCYQXV3UM06sGGrk=
+github.com/hashicorp/go-multierror v1.1.0/go.mod h1:spPvp8C1qA32ftKqdAHm4hHTbPw+vmowP0z+KUhOZdA=
+github.com/hashicorp/go-multierror v1.1.1 h1:H5DkEtf6CXdFp0N0Em5UCwQpXMWke8IA0+lD48awMYo=
+github.com/hashicorp/go-multierror v1.1.1/go.mod h1:iw975J/qwKPdAO1clOe2L8331t/9/fmwbPZ6JB6eMoM=
+github.com/hashicorp/go-retryablehttp v0.5.3/go.mod h1:9B5zBasrRhHXnJnui7y6sL7es7NDiJgTc6Er0maI1Xs=
+github.com/hashicorp/go-rootcerts v1.0.2 h1:jzhAVGtqPKbwpyCPELlgNWhE1znq+qwJtW5Oi2viEzc=
+github.com/hashicorp/go-rootcerts v1.0.2/go.mod h1:pqUvnprVnM5bf7AOirdbb01K4ccR319Vf4pU3K5EGc8=
+github.com/hashicorp/go-sockaddr v1.0.0/go.mod h1:7Xibr9yA9JjQq1JpNB2Vw7kxv8xerXegt+ozgdvDeDU=
+github.com/hashicorp/go-sockaddr v1.0.2 h1:ztczhD1jLxIRjVejw8gFomI1BQZOe2WoVOu0SyteCQc=
+github.com/hashicorp/go-sockaddr v1.0.2/go.mod h1:rB4wwRAUzs07qva3c5SdrY/NEtAUjGlgmH/UkBUC97A=
+github.com/hashicorp/go-syslog v1.0.0/go.mod h1:qPfqrKkXGihmCqbJM2mZgkZGvKG1dFdvsLplgctolz4=
+github.com/hashicorp/go-uuid v1.0.0/go.mod h1:6SBZvOh/SIDV7/2o3Jml5SYk/TvGqwFJ/bN7x4byOro=
+github.com/hashicorp/go-uuid v1.0.1/go.mod h1:6SBZvOh/SIDV7/2o3Jml5SYk/TvGqwFJ/bN7x4byOro=
+github.com/hashicorp/go-uuid v1.0.3 h1:2gKiV6YVmrJ1i2CKKa9obLvRieoRGviZFL26PcT/Co8=
+github.com/hashicorp/go-uuid v1.0.3/go.mod h1:6SBZvOh/SIDV7/2o3Jml5SYk/TvGqwFJ/bN7x4byOro=
+github.com/hashicorp/go-version v1.2.1 h1:zEfKbn2+PDgroKdiOzqiE8rsmLqU2uwi5PB5pBJ3TkI=
+github.com/hashicorp/go-version v1.2.1/go.mod h1:fltr4n8CU8Ke44wwGCBoEymUuxUHl09ZGVZPK5anwXA=
+github.com/hashicorp/golang-lru v0.5.0/go.mod h1:/m3WP610KZHVQ1SGc6re/UDhFvYD7pJ4Ao+sR/qLZy8=
+github.com/hashicorp/golang-lru v0.5.4 h1:YDjusn29QI/Das2iO9M0BHnIbxPeyuCHsjMW+lJfyTc=
+github.com/hashicorp/golang-lru v0.5.4/go.mod h1:iADmTwqILo4mZ8BN3D2Q6+9jd8WM5uGBxy+E8yxSoD4=
 github.com/hashicorp/hcl v1.0.0 h1:0Anlzjpi4vEasTeNFn2mLJgTSwt0+6sfsiTG8qcWGx4=
 github.com/hashicorp/hcl v1.0.0/go.mod h1:E5yfLk+7swimpb2L/Alb/PJmXilQ/rhwaUYs4T20WEQ=
+github.com/hashicorp/logutils v1.0.0/go.mod h1:QIAnNjmIWmVIIkWDTG1z5v++HQmx9WQRO+LraFDTW64=
+github.com/hashicorp/mdns v1.0.4/go.mod h1:mtBihi+LeNXGtG8L9dX59gAEa12BDtBQSp4v/YAJqrc=
+github.com/hashicorp/memberlist v0.5.0 h1:EtYPN8DpAURiapus508I4n9CzHs2W+8NZGbmmR/prTM=
+github.com/hashicorp/memberlist v0.5.0/go.mod h1:yvyXLpo0QaGE59Y7hDTsTzDD25JYBZ4mHgHUZ8lrOI0=
+github.com/hashicorp/serf v0.10.1 h1:Z1H2J60yRKvfDYAOZLd2MU0ND4AH/WDz7xYHDWQsIPY=
+github.com/hashicorp/serf v0.10.1/go.mod h1:yL2t6BqATOLGc5HF7qbFkTfXoPIY0WZdWHfEvMqbG+4=
+github.com/json-iterator/go v1.1.6/go.mod h1:+SdeFBvtyEkXs7REEP0seUULqWtbJapLOCVDaaPEHmU=
+github.com/json-iterator/go v1.1.9/go.mod h1:KdQUCv79m/52Kvf8AW2vK1V8akMuk1QjK/uOdHXbAo4=
 github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
 github.com/json-iterator/go v1.1.12/go.mod h1:e30LSqwooZae/UwlEbR2852Gd8hjQvJoHmT4TnhNGBo=
 github.com/juju/gnuflag v0.0.0-20171113085948-2ce1bb71843d/go.mod h1:2PavIy+JPciBPrBUjwbNvtwB6RQlve+hkpll6QSNmOE=
+github.com/julienschmidt/httprouter v1.2.0/go.mod h1:SYymIcj16QtmaHHD7aYtjjsJG7VTCxuUUipMqKk8s4w=
 github.com/kisielk/errcheck v1.5.0/go.mod h1:pFxgyoBC7bSaBwPgfKdkLd5X25qrDl4LWUI2bnpBCr8=
 github.com/kisielk/gotool v1.0.0/go.mod h1:XhKaO+MFFWcvkIS/tQcRk01m1F5IRFswLeQ+oQHNcck=
 github.com/klauspost/cpuid/v2 v2.0.9/go.mod h1:FInQzS24/EEf25PyTYn52gqo7WaD8xa0213Md/qVLRg=
@@ -80,6 +161,7 @@ github.com/klauspost/cpuid/v2 v2.2.5 h1:0E5MSMDEoAulmXNFquVs//DdoomxaoTY1kUhbc/q
 github.com/klauspost/cpuid/v2 v2.2.5/go.mod h1:Lcz8mBdAVJIBVzewtcLocK12l3Y+JytZYpaMropDUws=
 github.com/knz/go-libedit v1.10.1/go.mod h1:MZTVkCWyz0oBc7JOWP3wNAzd002ZbM/5hgShxwh4x8M=
 github.com/konsorten/go-windows-terminal-sequences v1.0.1/go.mod h1:T0+1ngSBFLxvqU3pZ+m/2kptfBszLMUkC4ZK/EgS/cQ=
+github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515/go.mod h1:+0opPa2QZZtGFBFZlji/RkVcI2GknAs/DXo4wKdlNEc=
 github.com/kr/pretty v0.1.0/go.mod h1:dAy3ld7l9f0ibDNOQOHHMYYIIbhfbHSm3C4ZsoJORNo=
 github.com/kr/pretty v0.3.1 h1:flRD4NNwYAUpkphVc1HcthR4KEIFJ65n8Mw5qdRn3LE=
 github.com/kr/pretty v0.3.1/go.mod h1:hoEshYVHaxMs3cyo3Yncou5ZscifuDolrwPKZanG3xk=
@@ -91,31 +173,78 @@ github.com/leodido/go-urn v1.2.4 h1:XlAE/cm/ms7TE/VMVoduSpNBoyc2dOxHs5MZSwAN63Q=
 github.com/leodido/go-urn v1.2.4/go.mod h1:7ZrI8mTSeBSHl/UaRyKQW1qZeMgak41ANeCNaVckg+4=
 github.com/magiconair/properties v1.8.7 h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=
 github.com/magiconair/properties v1.8.7/go.mod h1:Dhd985XPs7jluiymwWYZ0G4Z61jb3vdS329zhj2hYo0=
+github.com/mattn/go-colorable v0.0.9/go.mod h1:9vuHe8Xs5qXnSaW/c/ABM9alt+Vo+STaOChaDxuIBZU=
+github.com/mattn/go-colorable v0.1.4/go.mod h1:U0ppj6V5qS13XJ6of8GYAs25YV2eR4EVcfRqFIhoBtE=
+github.com/mattn/go-colorable v0.1.6/go.mod h1:u6P/XSegPjTcexA+o6vUJrdnUu04hMope9wVRipJSqc=
+github.com/mattn/go-colorable v0.1.9/go.mod h1:u6P/XSegPjTcexA+o6vUJrdnUu04hMope9wVRipJSqc=
+github.com/mattn/go-colorable v0.1.12/go.mod h1:u5H1YNBxpqRaxsYJYSkiCWKzEfiAb1Gb520KVy5xxl4=
+github.com/mattn/go-colorable v0.1.13 h1:fFA4WZxdEF4tXPZVKMLwD8oUnCTTo08duU7wxecdEvA=
+github.com/mattn/go-colorable v0.1.13/go.mod h1:7S9/ev0klgBDR4GtXTXX8a3vIGJpMovkB8vQcUbaXHg=
+github.com/mattn/go-isatty v0.0.3/go.mod h1:M+lRXTBqGeGNdLjl/ufCoiOlB5xdOkqRJdNxMWT7Zi4=
+github.com/mattn/go-isatty v0.0.8/go.mod h1:Iq45c/XA43vh69/j3iqttzPXn0bhXyGjM0Hdxcsrc5s=
+github.com/mattn/go-isatty v0.0.11/go.mod h1:PhnuNfih5lzO57/f3n+odYbM4JtupLOxQOAqxQCu2WE=
+github.com/mattn/go-isatty v0.0.12/go.mod h1:cbi8OIDigv2wuxKPP5vlRcQ1OAZbq2CE4Kysco4FUpU=
+github.com/mattn/go-isatty v0.0.14/go.mod h1:7GGIvUiUoEMVVmxf/4nioHXj79iQHKdU27kJ6hsGG94=
+github.com/mattn/go-isatty v0.0.16/go.mod h1:kYGgaQfpe5nmfYZH+SKPsOc2e4SrIfOl2e/yFXSvRLM=
 github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
 github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
+github.com/matttproud/golang_protobuf_extensions v1.0.1/go.mod h1:D8He9yQNgCq6Z5Ld7szi9bcBfOoFv/3dc6xSMkL2PC0=
+github.com/miekg/dns v1.1.26/go.mod h1:bPDLeHnStXmXAq1m/Ch/hvfNHr14JKNPMBo3VZKjuso=
+github.com/miekg/dns v1.1.41 h1:WMszZWJG0XmzbK9FEmzH2TVcqYzFesusSIB41b8KHxY=
+github.com/miekg/dns v1.1.41/go.mod h1:p6aan82bvRIyn+zDIv9xYNUpwa73JcSh9BKwknJysuI=
+github.com/mitchellh/cli v1.1.0/go.mod h1:xcISNoH86gajksDmfB23e/pu+B+GeFRMYmoHXxx3xhI=
+github.com/mitchellh/go-homedir v1.1.0 h1:lukF9ziXFxDFPkA1vsr5zpc1XuPDn/wFntq5mG+4E0Y=
+github.com/mitchellh/go-homedir v1.1.0/go.mod h1:SfyaCUpYCn1Vlf4IUYiD9fPX4A5wJrkLzIz1N1q0pr0=
+github.com/mitchellh/mapstructure v0.0.0-20160808181253-ca63d7c062ee/go.mod h1:FVVH3fgwuzCH5S8UJGiWEs2h04kUh9fWfEaFds41c1Y=
 github.com/mitchellh/mapstructure v1.5.0 h1:jeMsZIYE/09sWLaz43PL7Gy6RuMjD2eJVyuac5Z2hdY=
 github.com/mitchellh/mapstructure v1.5.0/go.mod h1:bFUtVrKA4DC2yAKiSyO/QUcy7e+RRV2QTWOzhPopBRo=
 github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
 github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=
 github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
+github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3RllmbCylyMrvgv0=
+github.com/modern-go/reflect2 v1.0.1/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3RllmbCylyMrvgv0=
 github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
 github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
+github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
 github.com/oapi-codegen/runtime v1.1.1 h1:EXLHh0DXIJnWhdRPN2w4MXAzFyE4CskzhNLUmtpMYro=
 github.com/oapi-codegen/runtime v1.1.1/go.mod h1:SK9X900oXmPWilYR5/WKPzt3Kqxn/uS/+lbpREv+eCg=
 github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFStlNbqXla1AfSYxPUl2o=
+github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
+github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
+github.com/pascaldekloe/goe v0.1.0/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pelletier/go-toml/v2 v2.2.2 h1:aYUidT7k73Pcl9nb2gScu7NSrKCSHIDE89b3+6Wq+LM=
 github.com/pelletier/go-toml/v2 v2.2.2/go.mod h1:1t835xjRzz80PqgE6HHgN2JOsmgYu/h4qDAS4n929Rs=
+github.com/pkg/errors v0.8.0/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
 github.com/pkg/errors v0.8.1/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
+github.com/pkg/errors v0.9.1 h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=
+github.com/pkg/errors v0.9.1/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
 github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
 github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 h1:Jamvg5psRIccs7FGNTlIRMkT8wgtp5eCXdBlqhYGL6U=
 github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
+github.com/posener/complete v1.1.1/go.mod h1:em0nMJCgc9GFtwrmVmEMR/ZL6WyhyjMBndrE9hABlRI=
+github.com/posener/complete v1.2.3/go.mod h1:WZIdtGGp+qx0sLrYKtIRAruyNpv6hFCicSgv7Sy7s/s=
+github.com/prometheus/client_golang v0.9.1/go.mod h1:7SWBe2y4D6OKWSNQJUaRYU/AaXPKyh/dDVn+NZz0KFw=
+github.com/prometheus/client_golang v1.0.0/go.mod h1:db9x61etRT2tGnBNRi70OPL5FsnadC4Ky3P0J6CfImo=
+github.com/prometheus/client_golang v1.4.0/go.mod h1:e9GMxYsXl05ICDXkRhurwBS4Q3OK1iX/F2sw+iXX5zU=
+github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910/go.mod h1:MbSGuTsp3dbXC40dX6PRTWyKYBIrTGTE9sqQNg2J8bo=
+github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
 github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
+github.com/prometheus/client_model v0.2.0/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
+github.com/prometheus/common v0.4.1/go.mod h1:TNfzLD0ON7rHzMJeJkieUDPYmFC7Snx/y86RQel1bk4=
+github.com/prometheus/common v0.9.1/go.mod h1:yhUN8i9wzaXS3w1O07YhxHEBxD+W35wd8bs7vj7HSQ4=
+github.com/prometheus/procfs v0.0.0-20181005140218-185b4288413d/go.mod h1:c3At6R/oaqEKCNdg8wHV1ftS6bRYblBhIjjI8uT2IGk=
+github.com/prometheus/procfs v0.0.2/go.mod h1:TjEm7ze935MbeOT/UhFTIMYKhuLP4wbCsTZCD3I8kEA=
+github.com/prometheus/procfs v0.0.8/go.mod h1:7Qr8sr6344vo1JqZ6HhLceV9o3AJ1Ff+GxbHq6oeK9A=
 github.com/rogpeppe/go-internal v1.9.0 h1:73kH8U+JUqXU8lRuOHeVHaa/SZPifC7BkcraZVejAe8=
 github.com/rogpeppe/go-internal v1.9.0/go.mod h1:WtVeX8xhTBvf0smdhujwtBcq4Qrzq/fJaraNFVN+nFs=
+github.com/ryanuber/columnize v0.0.0-20160712163229-9b3edd62028f/go.mod h1:sm1tb6uqfes/u+d4ooFouqFdy9/2g9QGwK3SQygK0Ts=
 github.com/sagikazarmark/locafero v0.4.0 h1:HApY1R9zGo4DBgr7dqsTH/JJxLTTsOt7u6keLGt6kNQ=
 github.com/sagikazarmark/locafero v0.4.0/go.mod h1:Pe1W6UlPYUk/+wc/6KFhbORCfqzgYEpgQ3O5fPuL3H4=
 github.com/sagikazarmark/slog-shim v0.1.0 h1:diDBnUNK9N/354PgrxMywXnAwEr1QZcOr6gto+ugjYE=
 github.com/sagikazarmark/slog-shim v0.1.0/go.mod h1:SrcSrq8aKtyuqEI1uvTDTK1arOWRIczQRv+GVI1AkeQ=
+github.com/sean-/seed v0.0.0-20170313163322-e2103e2c3529 h1:nn5Wsu0esKSJiIVhscUtVbo7ada43DJhG55ua/hjS5I=
+github.com/sean-/seed v0.0.0-20170313163322-e2103e2c3529/go.mod h1:DxrIzT+xaE7yg65j358z/aeFdxmN0P9QXhEzd20vsDc=
+github.com/sirupsen/logrus v1.2.0/go.mod h1:LxeOpSwHxABJmUn/MG1IvRgCAasNZTLOkJPxbbu5VWo=
 github.com/sirupsen/logrus v1.4.2/go.mod h1:tLMulIdttU9McNUspp0xgXVQah82FyeX6MwdIuYE2rE=
 github.com/sirupsen/logrus v1.8.1 h1:dJKuHgqk1NNQlqoA6BTlM1Wf9DOH3NBjQyu0h9+AZZE=
 github.com/sirupsen/logrus v1.8.1/go.mod h1:yWOB1SBYBC5VeMP7gHvWumXLIWorT60ONWic61uBYv0=
@@ -134,12 +263,14 @@ github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+
 github.com/stretchr/objx v0.1.1/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
 github.com/stretchr/objx v0.5.0/go.mod h1:Yh+to48EsGEfYuaHDzXPcE3xhTkx73EhmCGUpEOglKo=
+github.com/stretchr/objx v0.5.2 h1:xuMeJ0Sdp5ZMRXx/aWO6RZxdr3beISkG5/G/aIRr3pY=
 github.com/stretchr/objx v0.5.2/go.mod h1:FRsXN1f5AsAjCGJKqEizvkpNtU+EGNCLh3NxZ/8L+MA=
 github.com/stretchr/testify v1.2.2/go.mod h1:a8OnRcib4nhh0OaRAV+Yts87kKdq0PP7pXfy6kDkUVs=
 github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
 github.com/stretchr/testify v1.4.0/go.mod h1:j7eGeouHqKxXV5pUuKE4zz7dFj8WfuZ+81PSLYec5m4=
 github.com/stretchr/testify v1.7.0/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
 github.com/stretchr/testify v1.7.1/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
+github.com/stretchr/testify v1.7.2/go.mod h1:R6va5+xMeoiuVRoj+gSkQ7d3FALtqAAGI1FQKckRals=
 github.com/stretchr/testify v1.8.0/go.mod h1:yNjHg4UonilssWZ8iaSj1OCr/vHnekPRkoO+kdMU+MU=
 github.com/stretchr/testify v1.8.1/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o6fzry7u4=
 github.com/stretchr/testify v1.8.2/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o6fzry7u4=
@@ -148,6 +279,7 @@ github.com/stretchr/testify v1.9.0 h1:HtqpIVDClZ4nwg75+f6Lvsy/wHu+3BoSGCbBAcpTsT
 github.com/stretchr/testify v1.9.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8C91i36aY=
 github.com/subosito/gotenv v1.6.0 h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=
 github.com/subosito/gotenv v1.6.0/go.mod h1:Dk4QP5c2W3ibzajGcXpNraDfq2IrhjMIvMSWPKKo0FU=
+github.com/tv42/httpunix v0.0.0-20150427012821-b75d8614f926/go.mod h1:9ESjWnEqriFuLhtthL60Sar/7RFoluCcXsuvEwTV5KM=
 github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS4MhqMhdFk5YI=
 github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
 github.com/ugorji/go/codec v1.2.11 h1:BMaWp1Bb6fHwEtbplGBGJ498wD+LKlNSl25MjdZY4dU=
@@ -165,11 +297,13 @@ go.uber.org/zap v1.18.1/go.mod h1:xg/QME4nWcxGxrpdeYfq7UvYrLh66cuVKdrbD1XF/NI=
 golang.org/x/arch v0.0.0-20210923205945-b76863e36670/go.mod h1:5om86z9Hs0C8fWVUuoMHwpExlXzs5Tkyp9hOrfG7pp8=
 golang.org/x/arch v0.4.0 h1:A8WCeEWhLwPBKNbFi5Wv5UTCBx5zzubnXDlMOFAzFMc=
 golang.org/x/arch v0.4.0/go.mod h1:5om86z9Hs0C8fWVUuoMHwpExlXzs5Tkyp9hOrfG7pp8=
+golang.org/x/crypto v0.0.0-20180904163835-0709b304e793/go.mod h1:6SG95UA2DQfeDnfUPMdvaQW0Q7yPrPDi9nlGo2tz2b4=
 golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACkg1iLfiJU5Ep61QUkGW8qpdssI0+w=
+golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392/go.mod h1:/lpIB1dKB+9EgE3H3cr1v9wB50oz8l4C4h62xy7jSTY=
 golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
 golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9/go.mod h1:LzIPMQfyMNhhGPhUkYOs5KpL4U8rLKemX1yGLhDgUto=
-golang.org/x/crypto v0.21.0 h1:X31++rzVUdKhX5sWmSOFZxx8UW/ldWx55cbf08iNAMA=
-golang.org/x/crypto v0.21.0/go.mod h1:0BP7YvVV9gBbVKyeTG0Gyn+gZm94bibOW5BjDEYAOMs=
+golang.org/x/crypto v0.26.0 h1:RrRspgV4mU+YwB4FYnuBoKsUapNIL5cohGAmSH3azsw=
+golang.org/x/crypto v0.26.0/go.mod h1:GY7jblb9wI+FOo5y8/S2oY4zWP07AkOJ4+jxCqdqn54=
 golang.org/x/exp v0.0.0-20190121172915-509febef88a4/go.mod h1:CJ0aWSM057203Lf6IL+f9T1iT9GByDxfZKAQTCR3kQA=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9 h1:GoHiUyI/Tp2nVkLI2mCxVkOjsbSXD66ic0XW0js0R9g=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9/go.mod h1:S2oDrQGGwySpoQPVqRShND87VCbxmc6bL1Yd2oYrm6k=
@@ -181,40 +315,68 @@ golang.org/x/mod v0.2.0/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
 golang.org/x/mod v0.3.0/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
 golang.org/x/net v0.0.0-20180724234803-3673e40ba225/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20180826012351-8a410e7b638d/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
+golang.org/x/net v0.0.0-20181114220301-adae6a3d119a/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20190213061140-3a22650c66bd/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20190311183353-d8887717615a/go.mod h1:t9HGtf8HONx5eT2rtn7q6eTqICYqUVnKs3thJo3Qplg=
 golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3/go.mod h1:t9HGtf8HONx5eT2rtn7q6eTqICYqUVnKs3thJo3Qplg=
+golang.org/x/net v0.0.0-20190613194153-d28f0bde5980/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
 golang.org/x/net v0.0.0-20190620200207-3b0461eec859/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
+golang.org/x/net v0.0.0-20190923162816-aa69164e4478/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
 golang.org/x/net v0.0.0-20200226121028-0de0cce0169b/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
 golang.org/x/net v0.0.0-20201021035429-f5854403a974/go.mod h1:sp8m0HH+o8qH0wwXwYZr8TS3Oi6o0r6Gce1SSxlDquU=
-golang.org/x/net v0.23.0 h1:7EYJ93RZ9vYSZAIb2x3lnuvqO5zneoD6IvWjuhfxjTs=
-golang.org/x/net v0.23.0/go.mod h1:JKghWKKOSdJwpW2GEx0Ja7fmaKnMsbu+MWVZTokSYmg=
+golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
+golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1/go.mod h1:9tjilg8BloeKEkVJvy7fQ90B1CfIiPueXVOjqfkSzI8=
+golang.org/x/net v0.28.0 h1:a9JDOJc5GMUJ0+UDqmLT86WiEy7iWyIhz8gz8E4e5hE=
+golang.org/x/net v0.28.0/go.mod h1:yqtgsTWOOnlGLG9GFRrK3++bGOUEkNBoHZc8MEDWPNg=
 golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be/go.mod h1:N/0e6XlmueqKjAGxoOufVs8QHGRruUQn6yWY3a++T0U=
 golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20181108010431-42b317875d0f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
+golang.org/x/sync v0.0.0-20181221193216-37e7f081c4d4/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20190423024810-112230192c58/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
+golang.org/x/sync v0.0.0-20210220032951-036812b2e83c/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
+golang.org/x/sys v0.0.0-20180823144017-11551d06cbcc/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20180830151530-49385e6e1522/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
+golang.org/x/sys v0.0.0-20180905080454-ebe1bf3edb33/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
+golang.org/x/sys v0.0.0-20181116152217-5ac8a444bdc5/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
+golang.org/x/sys v0.0.0-20190222072716-a9d3bda3a223/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20190412213103-97732733099d/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20190422165155-953cdadca894/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20190922100055-0a153f010e69/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20190924154521-2837fb4f24fe/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20191026070338-33540a1f6037/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20200116001909-b77594299b42/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20200122134326-e047566fdf82/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20201119102817-f84b799fce68/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20210303074136-134d130e1a04/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20210330210617-4fbd30eecc44/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
-golang.org/x/sys v0.18.0 h1:DBdB3niSjOA/O0blCZBqDefyWNYveAYMNF1Wum0DYQ4=
-golang.org/x/sys v0.18.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
+golang.org/x/sys v0.24.0 h1:Twjiwq9dn6R1fQcyiK+wQyHWfaz/BJB+YIpzU/Cv3Xg=
+golang.org/x/sys v0.24.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
+golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
+golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
-golang.org/x/text v0.14.0 h1:ScX5w1eTa3QqT8oi6+ziP7dTV1S2+ALU0bI+0zXKWiQ=
-golang.org/x/text v0.14.0/go.mod h1:18ZOQIKpY8NJVqYksKHtTdi31H5itFRjB5/qKTNYzSU=
+golang.org/x/text v0.3.6/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
+golang.org/x/text v0.17.0 h1:XtiM5bkSOt+ewxlOE/aE/AKEHibwj/6gvWMl9Rsh0Qc=
+golang.org/x/text v0.17.0/go.mod h1:BuEKDfySbSR4drPmRPG/7iBdf8hvFMuRexcpahXilzY=
 golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190114222345-bf090417da8b/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190226205152-f727befe758c/go.mod h1:9Yl7xja0Znq3iFh3HoIrodX9oNMXvdceNzlUR8zjMvY=
 golang.org/x/tools v0.0.0-20190311212946-11955173bddd/go.mod h1:LCzVGOaR6xXOjkQ3onu1FJEFr0SW1gC7cKk1uF8kGRs=
 golang.org/x/tools v0.0.0-20190524140312-2c0ae7006135/go.mod h1:RgjU9mgBXZiqYHBnxXauZ1Gv1EHHAz9KjViQ78xBX0Q=
+golang.org/x/tools v0.0.0-20190907020128-2ca718005c18/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
 golang.org/x/tools v0.0.0-20191108193012-7d206e10da11/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
 golang.org/x/tools v0.0.0-20191119224855-298f0cb1881e/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
 golang.org/x/tools v0.0.0-20200619180055-7c47624df98f/go.mod h1:EkVYQZoAsY45+roYkvgYkIh4xh/qjgUK9TdY2XT94GE=
@@ -228,26 +390,30 @@ google.golang.org/appengine v1.4.0/go.mod h1:xpcJRLb0r/rnEns0DIKYYv+WjYCduHsrkT7
 google.golang.org/genproto v0.0.0-20180817151627-c66870c02cf8/go.mod h1:JiN7NxoALGmiZfu7CAH4rXhgtRTLTxftemlI0sWmxmc=
 google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55/go.mod h1:DMBHOl98Agz4BDEuKkezgsaosCRResVns1a3J2ZsMNc=
 google.golang.org/genproto v0.0.0-20200423170343-7949de9c1215/go.mod h1:55QSHmfGQM9UVYDPBsyGGes0y52j32PQ3BqQfXhyH3c=
-google.golang.org/genproto/googleapis/rpc v0.0.0-20240314234333-6e1732d8331c h1:lfpJ/2rWPa/kJgxyyXM8PrNnfCzcmxJ265mADgwmvLI=
-google.golang.org/genproto/googleapis/rpc v0.0.0-20240314234333-6e1732d8331c/go.mod h1:WtryC6hu0hhx87FDGxWCDptyssuo68sk10vYjF+T9fY=
+google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 h1:e7S5W7MGGLaSu8j3YjdezkZ+m1/Nm0uRVRMEMGk26Xs=
+google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142/go.mod h1:UqMtugtsSgubUsoxbuAoiCXvqvErP7Gf0so0mK9tHxU=
 google.golang.org/grpc v1.19.0/go.mod h1:mqu4LbDTu4XGKhr4mRzUsmM4RtVoemTSY81AxZiDr8c=
 google.golang.org/grpc v1.23.0/go.mod h1:Y5yQAOtifL1yxbo5wqy6BxZv8vAUGQwXBOALyacEbxg=
 google.golang.org/grpc v1.25.1/go.mod h1:c3i+UQWmh7LiEpx4sFZnkU36qjEYZ0imhYfXVyQciAY=
 google.golang.org/grpc v1.27.0/go.mod h1:qbnxyOmOxrQa7FizSgH+ReBfzJrCY1pSN7KXBS8abTk=
 google.golang.org/grpc v1.29.1/go.mod h1:itym6AZVZYACWQqET3MqgPpjcuV5QH3BxFS3IjizoKk=
-google.golang.org/grpc v1.62.1 h1:B4n+nfKzOICUXMgyrNd19h/I9oH0L1pizfk1d4zSgTk=
-google.golang.org/grpc v1.62.1/go.mod h1:IWTG0VlJLCh1SkC58F7np9ka9mx/WNkjl4PGJaiq+QE=
+google.golang.org/grpc v1.67.1 h1:zWnc1Vrcno+lHZCOofnIMvycFcc0QRGIzm9dhnDX68E=
+google.golang.org/grpc v1.67.1/go.mod h1:1gLDyUQU7CTLJI90u3nXZ9ekeghjeM7pTDZlqFNg2AA=
 google.golang.org/protobuf v1.26.0-rc.1/go.mod h1:jlhhOSvTdKEhbULTjvd4ARK9grFBp09yW+WbY/TyQbw=
 google.golang.org/protobuf v1.26.0/go.mod h1:9q0QmTI4eRPtz6boOQmLYwt+qCgq0jsYwAQnmE0givc=
-google.golang.org/protobuf v1.33.0 h1:uNO2rsAINq/JlFpSdYEKIZ0uKD/R9cpdv0T+yoGwGmI=
-google.golang.org/protobuf v1.33.0/go.mod h1:c6P6GXX6sHbq/GpV6MGZEdwhWPcYBgnhAHhKbcUYpos=
+google.golang.org/protobuf v1.35.1 h1:m3LfL6/Ca+fqnjnlqQXNpFPABW1UD7mjh8KO2mKFytA=
+google.golang.org/protobuf v1.35.1/go.mod h1:9fA7Ob0pmnwhb644+1+CVWFRbNajQ6iRojtC/QF5bRE=
+gopkg.in/alecthomas/kingpin.v2 v2.2.6/go.mod h1:FMv+mEhP44yOT+4EoQTLFTRgOQ1FBLkstjWtayDeSgw=
 gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 h1:YR8cESwS4TdDjEe65xsg0ogRM/Nc3DYOhEAlW+xobZo=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/ini.v1 v1.67.0 h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=
 gopkg.in/ini.v1 v1.67.0/go.mod h1:pNLf8WUiyNEtQjuu5G5vTm06TEv9tsIgeAvK8hOrP4k=
+gopkg.in/yaml.v2 v2.2.1/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.2/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
+gopkg.in/yaml.v2 v2.2.4/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
+gopkg.in/yaml.v2 v2.2.5/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.8/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
 gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
diff --git a/internal/common/metrics/todo_metrics.go b/internal/common/metrics/todo_metrics.go
index b379f64..a05f6b2 100644
--- a/internal/common/metrics/todo_metrics.go
+++ b/internal/common/metrics/todo_metrics.go
@@ -1,7 +1,5 @@
 package metrics
 
-// TodoMetrics 是占位实现，用来先打通接口而不真正上报指标。
-// 初学时先理解“依赖抽象”比马上接 Prometheus 更重要。
 type TodoMetrics struct{}
 
 func (t TodoMetrics) Inc(_ string, _ int) {
diff --git a/internal/common/server/gprc.go b/internal/common/server/gprc.go
index ba9f6e1..817cfcb 100644
--- a/internal/common/server/gprc.go
+++ b/internal/common/server/gprc.go
@@ -10,16 +10,12 @@ import (
 	"google.golang.org/grpc"
 )
 
-// init 会在包加载时执行，这里统一替换 gRPC 默认日志实现，
-// 让框架日志和业务日志都走 logrus，便于初学者排查问题时只看一种格式。
 func init() {
 	logger := logrus.New()
 	logger.SetLevel(logrus.WarnLevel)
 	grpc_logrus.ReplaceGrpcLogger(logrus.NewEntry(logger))
 }
 
-// RunGRPCServer 负责按服务名读取配置，再把真正的启动动作委托给 RunGRPCServerOnAddr。
-// 这样测试时可以直接传地址，避免每次都依赖 viper 配置。
 func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server)) {
 	addr := viper.Sub(serviceName).GetString("grpc-addr")
 	if addr == "" {
@@ -29,13 +25,10 @@ func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server))
 	RunGRPCServerOnAddr(addr, registerServer)
 }
 
-// RunGRPCServerOnAddr 创建 gRPC server，并通过 registerServer 回调注册具体业务服务。
-// 这类“框架负责启动，业务负责注册”的回调模式，在 Go 后端里非常常见。
 func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
 	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
 	grpcServer := grpc.NewServer(
 		grpc.ChainUnaryInterceptor(
-			// 一元拦截器链类似 HTTP 中间件，用来放日志、指标、鉴权这类横切逻辑。
 			grpc_tags.UnaryServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
 			grpc_logrus.UnaryServerInterceptor(logrusEntry),
 			//otelgrpc.UnaryServerInterceptor(),
@@ -45,7 +38,6 @@ func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server))
 			//recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
 		),
 		grpc.ChainStreamInterceptor(
-			// 流式 RPC 和一元 RPC 分开配置，是因为两类调用的拦截接口不同。
 			grpc_tags.StreamServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
 			grpc_logrus.StreamServerInterceptor(logrusEntry),
 			//otelgrpc.StreamServerInterceptor(),
@@ -55,10 +47,8 @@ func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server))
 			//recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
 		),
 	)
-	// 业务方在这里把 OrderService/StockService 的实现挂到 grpcServer 上。
 	registerServer(grpcServer)
 
-	// net.Listen 只负责占住端口；真正开始接收 gRPC 请求要等 Serve 调用后才发生。
 	listen, err := net.Listen("tcp", addr)
 	if err != nil {
 		logrus.Panic(err)
diff --git a/internal/common/server/http.go b/internal/common/server/http.go
index 1c561ce..7f39359 100644
--- a/internal/common/server/http.go
+++ b/internal/common/server/http.go
@@ -5,7 +5,6 @@ import (
 	"github.com/spf13/viper"
 )
 
-// RunHTTPServer 按服务名从配置中读取监听地址，再交给 RunHTTPServerOnAddr 真正启动。
 func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
 	addr := viper.Sub(serviceName).GetString("http-addr")
 	if addr == "" {
@@ -14,13 +13,9 @@ func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
 	RunHTTPServerOnAddr(addr, wrapper)
 }
 
-// RunHTTPServerOnAddr 创建 gin.Engine，并把路由注册动作交给 wrapper 回调。
 func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
-	// gin.New 不带默认中间件，适合教学阶段观察“哪些能力是手动加进去的”。
 	apiRouter := gin.New()
-	// wrapper 负责把 OpenAPI 生成的 handler 绑定到具体路由。
 	wrapper(apiRouter)
-	// 这一行不会修改已有路由，只是创建了一个没有被接住的 RouterGroup。
 	apiRouter.Group("/api")
 	if err := apiRouter.Run(addr); err != nil {
 		panic(err)
diff --git a/internal/order/adapters/grpc/stock_grpc.go b/internal/order/adapters/grpc/stock_grpc.go
index f3acb2b..c2397d6 100644
--- a/internal/order/adapters/grpc/stock_grpc.go
+++ b/internal/order/adapters/grpc/stock_grpc.go
@@ -8,8 +8,6 @@ import (
 	"github.com/sirupsen/logrus"
 )
 
-// StockGRPC 是 order 侧访问 stock 服务的远程适配器。
-// 它实现的是应用层定义的 StockService 接口，而不是把 proto client 直接暴露出去。
 type StockGRPC struct {
 	client stockpb.StockServiceClient
 }
@@ -19,14 +17,12 @@ func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
 }
 
 func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error) {
-	// 这里负责把应用层参数翻译成 gRPC 请求对象。
 	resp, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
 	logrus.Info("stock_grpc response", resp)
 	return resp, err
 }
 
 func (s StockGRPC) GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error) {
-	// 对调用方来说这里只是“查商品”，至于底层走 gRPC 还是别的协议都被屏蔽了。
 	resp, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{ItemIDs: itemIDs})
 	if err != nil {
 		return nil, err
diff --git a/internal/order/adapters/order_inmem_repository.go b/internal/order/adapters/order_inmem_repository.go
index 94364ec..818b51f 100644
--- a/internal/order/adapters/order_inmem_repository.go
+++ b/internal/order/adapters/order_inmem_repository.go
@@ -10,15 +10,12 @@ import (
 	"github.com/sirupsen/logrus"
 )
 
-// MemoryOrderRepository 是 domain.Repository 的内存版实现。
-// 教学项目先用它跑通流程，后面换数据库时只需要替换这一层。
 type MemoryOrderRepository struct {
 	lock  *sync.RWMutex
 	store []*domain.Order
 }
 
 func NewMemoryOrderRepository() *MemoryOrderRepository {
-	// 这里放一条假数据，方便一开始就能演示“查询已有订单”的路径。
 	s := make([]*domain.Order, 0)
 	s = append(s, &domain.Order{
 		ID:          "fake-ID",
@@ -34,11 +31,9 @@ func NewMemoryOrderRepository() *MemoryOrderRepository {
 }
 
 func (m *MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
-	// 写操作要加互斥锁，避免多个请求并发修改切片时产生数据竞争。
 	m.lock.Lock()
 	defer m.lock.Unlock()
 	newOrder := &domain.Order{
-		// 当前用 Unix 时间戳凑一个简单 ID，真实项目里通常会换成雪花算法或 UUID。
 		ID:          strconv.FormatInt(time.Now().Unix(), 10),
 		CustomerID:  order.CustomerID,
 		Status:      order.Status,
@@ -57,7 +52,6 @@ func (m *MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*
 	for i, v := range m.store {
 		logrus.Infof("m.store[%d] = %+v", i, v)
 	}
-	// 读操作使用读锁，允许多个查询并发进行。
 	m.lock.RLock()
 	defer m.lock.RUnlock()
 	for _, o := range m.store {
@@ -70,7 +64,6 @@ func (m *MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*
 }
 
 func (m *MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
-	// UpdateFn 把“怎么改”交给上层，把“在哪里存”留给仓储层，是一种职责分离。
 	m.lock.Lock()
 	defer m.lock.Unlock()
 	found := false
diff --git a/internal/order/app/app.go b/internal/order/app/app.go
index 5963e7c..b2e04ce 100644
--- a/internal/order/app/app.go
+++ b/internal/order/app/app.go
@@ -5,19 +5,16 @@ import (
 	"github.com/ghost-yu/go_shop_second/order/app/query"
 )
 
-// Application 是应用层门面，让上层只依赖一个总入口，而不用感知每个 handler 的构造细节。
 type Application struct {
 	Commands Commands
 	Queries  Queries
 }
 
-// Commands 聚合“会改状态”的用例，典型如创建、更新订单。
 type Commands struct {
 	CreateOrder command.CreateOrderHandler
 	UpdateOrder command.UpdateOrderHandler
 }
 
-// Queries 聚合“只读”的用例，便于和 Commands 做职责分离。
 type Queries struct {
 	GetCustomerOrder query.GetCustomerOrderHandler
 }
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
index 859e170..1c72b04 100644
--- a/internal/order/app/command/create_order.go
+++ b/internal/order/app/command/create_order.go
@@ -11,20 +11,17 @@ import (
 	"github.com/sirupsen/logrus"
 )
 
-// CreateOrder 是创建订单用例的输入模型。
 type CreateOrder struct {
 	CustomerID string
 	Items      []*orderpb.ItemWithQuantity
 }
 
-// CreateOrderResult 是返回给端口层的结果，避免直接暴露底层仓储对象。
 type CreateOrderResult struct {
 	OrderID string
 }
 
 type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]
 
-// createOrderHandler 只持有完成下单流程所需的两个依赖：订单仓储和库存服务。
 type createOrderHandler struct {
 	orderRepo domain.Repository
 	stockGRPC query.StockService
@@ -39,7 +36,6 @@ func NewCreateOrderHandler(
 	if orderRepo == nil {
 		panic("nil orderRepo")
 	}
-	// 下单也是 command，所以同样套上统一的日志和指标装饰器。
 	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
 		createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC},
 		logger,
@@ -48,7 +44,6 @@ func NewCreateOrderHandler(
 }
 
 func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
-	// 先校验库存，再真正创建订单，避免生成一张无法履约的脏订单。
 	validItems, err := c.validate(ctx, cmd.Items)
 	if err != nil {
 		return nil, err
@@ -67,7 +62,6 @@ func (c createOrderHandler) validate(ctx context.Context, items []*orderpb.ItemW
 	if len(items) == 0 {
 		return nil, errors.New("must have at least one item")
 	}
-	// packItems 先合并重复商品，避免把同一个商品重复发给库存服务校验。
 	items = packItems(items)
 	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, items)
 	if err != nil {
@@ -77,7 +71,6 @@ func (c createOrderHandler) validate(ctx context.Context, items []*orderpb.ItemW
 }
 
 func packItems(items []*orderpb.ItemWithQuantity) []*orderpb.ItemWithQuantity {
-	// merged 以商品 ID 为 key，把数量累加起来。
 	merged := make(map[string]int32)
 	for _, item := range items {
 		merged[item.ID] += item.Quantity
diff --git a/internal/order/app/command/update_order.go b/internal/order/app/command/update_order.go
index 3cce1e7..f40716d 100644
--- a/internal/order/app/command/update_order.go
+++ b/internal/order/app/command/update_order.go
@@ -8,7 +8,6 @@ import (
 	"github.com/sirupsen/logrus"
 )
 
-// UpdateOrder 把“更新哪个订单”和“如何更新”一起传给 handler。
 type UpdateOrder struct {
 	Order    *domain.Order
 	UpdateFn func(context.Context, *domain.Order) (*domain.Order, error)
@@ -16,7 +15,6 @@ type UpdateOrder struct {
 
 type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, interface{}]
 
-// updateOrderHandler 负责流程控制，具体修改细节由 UpdateFn 注入。
 type updateOrderHandler struct {
 	orderRepo domain.Repository
 	//stockGRPC
@@ -38,7 +36,6 @@ func NewUpdateOrderHandler(
 }
 
 func (c updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (interface{}, error) {
-	// 给 nil UpdateFn 一个 no-op 默认值，避免直接调用时发生空指针问题。
 	if cmd.UpdateFn == nil {
 		logrus.Warnf("updateOrderHandler got nil UpdateFn, order=%#v", cmd.Order)
 		cmd.UpdateFn = func(_ context.Context, order *domain.Order) (*domain.Order, error) { return order, nil }
diff --git a/internal/order/app/query/get_customer_order.go b/internal/order/app/query/get_customer_order.go
index 19b7ffc..b107e01 100644
--- a/internal/order/app/query/get_customer_order.go
+++ b/internal/order/app/query/get_customer_order.go
@@ -8,17 +8,13 @@ import (
 	"github.com/sirupsen/logrus"
 )
 
-// GetCustomerOrder 是查询订单的输入对象。
-// 把参数收成一个结构体后，后续加字段不会破坏函数签名。
 type GetCustomerOrder struct {
 	CustomerID string
 	OrderID    string
 }
 
-// GetCustomerOrderHandler 是一个已经套好泛型的查询处理器类型别名。
 type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Order]
 
-// getCustomerOrderHandler 才是真正的业务实现，字段里只放它需要的依赖。
 type getCustomerOrderHandler struct {
 	orderRepo domain.Repository
 }
@@ -31,7 +27,6 @@ func NewGetCustomerOrderHandler(
 	if orderRepo == nil {
 		panic("nil orderRepo")
 	}
-	// 构造函数里统一加装饰器，这样调用方拿到的就是“带日志和指标能力”的 handler。
 	return decorator.ApplyQueryDecorators[GetCustomerOrder, *domain.Order](
 		getCustomerOrderHandler{orderRepo: orderRepo},
 		logger,
@@ -40,7 +35,6 @@ func NewGetCustomerOrderHandler(
 }
 
 func (g getCustomerOrderHandler) Handle(ctx context.Context, query GetCustomerOrder) (*domain.Order, error) {
-	// 查询用例本身很薄，只负责描述流程，把数据访问交给仓储接口。
 	o, err := g.orderRepo.Get(ctx, query.OrderID, query.CustomerID)
 	if err != nil {
 		return nil, err
diff --git a/internal/order/app/query/service.go b/internal/order/app/query/service.go
index 0d33749..2e3e4f4 100644
--- a/internal/order/app/query/service.go
+++ b/internal/order/app/query/service.go
@@ -7,8 +7,6 @@ import (
 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
 )
 
-// StockService 是 order 应用层眼里的“库存能力端口”。
-// 它隔离了底层 gRPC 细节，让 create_order 用例只表达业务意图。
 type StockService interface {
 	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error)
 	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
diff --git a/internal/order/domain/order/order.go b/internal/order/domain/order/order.go
index 073b8f0..6a9e94b 100644
--- a/internal/order/domain/order/order.go
+++ b/internal/order/domain/order/order.go
@@ -2,16 +2,10 @@ package order
 
 import "github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 
-// Order 是订单领域对象，代表业务里真正被创建、查询、更新的核心实体。
 type Order struct {
-	// ID 是订单主键，由仓储层创建时生成。
-	ID string
-	// CustomerID 表示订单属于哪个客户。
-	CustomerID string
-	// Status 预留给支付中、已支付等订单状态流转。
-	Status string
-	// PaymentLink 预留给支付服务返回的支付链接。
+	ID          string
+	CustomerID  string
+	Status      string
 	PaymentLink string
-	// Items 直接复用 proto 里的商品结构，减少服务边界两侧的对象转换成本。
-	Items []*orderpb.Item
+	Items       []*orderpb.Item
 }
diff --git a/internal/order/domain/order/repository.go b/internal/order/domain/order/repository.go
index 685eddc..04b783f 100644
--- a/internal/order/domain/order/repository.go
+++ b/internal/order/domain/order/repository.go
@@ -5,7 +5,6 @@ import (
 	"fmt"
 )
 
-// Repository 定义订单持久化能力，应用层只依赖这个接口，而不关心底层是内存还是数据库。
 type Repository interface {
 	Create(context.Context, *Order) (*Order, error)
 	Get(ctx context.Context, id, customerID string) (*Order, error)
@@ -16,7 +15,6 @@ type Repository interface {
 	) error
 }
 
-// NotFoundError 是带业务语义的错误类型，调用方可以明确知道“订单不存在”而不是普通异常。
 type NotFoundError struct {
 	OrderID string
 }
diff --git a/internal/order/go.mod b/internal/order/go.mod
index bd794ca..3156da5 100644
--- a/internal/order/go.mod
+++ b/internal/order/go.mod
@@ -16,10 +16,12 @@ require (
 
 require (
 	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
+	github.com/armon/go-metrics v0.4.1 // indirect
 	github.com/bytedance/sonic v1.12.3 // indirect
 	github.com/bytedance/sonic/loader v0.2.0 // indirect
 	github.com/cloudwego/base64x v0.1.4 // indirect
 	github.com/cloudwego/iasm v0.2.0 // indirect
+	github.com/fatih/color v1.14.1 // indirect
 	github.com/fsnotify/fsnotify v1.7.0 // indirect
 	github.com/gabriel-vasile/mimetype v1.4.6 // indirect
 	github.com/gin-contrib/sse v0.1.0 // indirect
@@ -30,12 +32,23 @@ require (
 	github.com/golang/protobuf v1.5.4 // indirect
 	github.com/google/uuid v1.6.0 // indirect
 	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
+	github.com/hashicorp/consul/api v1.28.2 // indirect
+	github.com/hashicorp/errwrap v1.1.0 // indirect
+	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
+	github.com/hashicorp/go-hclog v1.5.0 // indirect
+	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
+	github.com/hashicorp/go-multierror v1.1.1 // indirect
+	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
+	github.com/hashicorp/golang-lru v0.5.4 // indirect
 	github.com/hashicorp/hcl v1.0.0 // indirect
+	github.com/hashicorp/serf v0.10.1 // indirect
 	github.com/json-iterator/go v1.1.12 // indirect
 	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
 	github.com/leodido/go-urn v1.4.0 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
+	github.com/mattn/go-colorable v0.1.13 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
+	github.com/mitchellh/go-homedir v1.1.0 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
 	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
 	github.com/modern-go/reflect2 v1.0.2 // indirect
diff --git a/internal/order/go.sum b/internal/order/go.sum
index 5e77cf7..fc15276 100644
--- a/internal/order/go.sum
+++ b/internal/order/go.sum
@@ -1,9 +1,24 @@
 cloud.google.com/go v0.26.0/go.mod h1:aQUYkXzVsufM+DwF1aE+0xfcU+56JwCaLick0ClmMTw=
 github.com/BurntSushi/toml v0.3.1/go.mod h1:xHWCNGjB5oqiDr8zfno3MHue2Ht5sIBksp03qcyfWMU=
+github.com/DataDog/datadog-go v3.2.0+incompatible/go.mod h1:LButxg5PwREeZtORoXG3tL4fMGNddJ+vMq1mwgfaqoQ=
 github.com/RaveNoX/go-jsoncommentstrip v1.0.0/go.mod h1:78ihd09MekBnJnxpICcwzCMzGrKSKYe4AqU6PDYYpjk=
+github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
+github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
+github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
+github.com/alecthomas/units v0.0.0-20190717042225-c3de453c63f4/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
 github.com/apapsch/go-jsonmerge/v2 v2.0.0 h1:axGnT1gRIfimI7gJifB699GoE/oq+F2MU7Dml6nw9rQ=
 github.com/apapsch/go-jsonmerge/v2 v2.0.0/go.mod h1:lvDnEdqiQrp0O42VQGgmlKpxL1AP2+08jFMw88y4klk=
+github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e/go.mod h1:3U/XgcO3hCbHZ8TKRvWD2dDTCfh9M9ya+I9JpbB7O8o=
+github.com/armon/go-metrics v0.0.0-20180917152333-f0300d1749da/go.mod h1:Q73ZrmVTwzkszR9V5SSuryQ31EELlFMUz1kKyl939pY=
+github.com/armon/go-metrics v0.4.1 h1:hR91U9KYmb6bLBYLQjyM+3j+rcd/UhE+G78SFnF8gJA=
+github.com/armon/go-metrics v0.4.1/go.mod h1:E6amYzXo6aW1tqzoZGT755KkbgrJsSdpwZ+3JqfkOG4=
+github.com/armon/go-radix v0.0.0-20180808171621-7fddfc383310/go.mod h1:ufUuZ+zHj4x4TnLV4JWEpy2hxWSpsRywHrMgIH9cCH8=
+github.com/armon/go-radix v1.0.0/go.mod h1:ufUuZ+zHj4x4TnLV4JWEpy2hxWSpsRywHrMgIH9cCH8=
 github.com/benbjohnson/clock v1.1.0/go.mod h1:J11/hYXuz8f4ySSvYwY0FKfm+ezbsZBKZxNJlLklBHA=
+github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973/go.mod h1:Dwedo/Wpr24TaqPxmxbtue+5NUziq4I4S80YR8gNf3Q=
+github.com/beorn7/perks v1.0.0/go.mod h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+CedLV8=
+github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
+github.com/bgentry/speakeasy v0.1.0/go.mod h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=
 github.com/bmatcuk/doublestar v1.1.1/go.mod h1:UD6OnuiIn0yFxxA2le/rnRU1G4RaI4UvFv1sNto9p6w=
 github.com/bytedance/sonic v1.12.3 h1:W2MGa7RCU1QTeYRTPE3+88mVC0yXmsRQRChiyVocVjU=
 github.com/bytedance/sonic v1.12.3/go.mod h1:B8Gt/XvtZ3Fqj+iSKMypzymZxw/FVwgIGKzMzT9r/rk=
@@ -11,6 +26,9 @@ github.com/bytedance/sonic/loader v0.1.1/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4
 github.com/bytedance/sonic/loader v0.2.0 h1:zNprn+lsIP06C/IqCHs3gPQIvnvpKbbxyXQP1iU4kWM=
 github.com/bytedance/sonic/loader v0.2.0/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
 github.com/census-instrumentation/opencensus-proto v0.2.1/go.mod h1:f6KPmirojxKA12rnyqOA5BBL4O983OfeGPqjHWSTneU=
+github.com/cespare/xxhash/v2 v2.1.1/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
+github.com/circonus-labs/circonus-gometrics v2.3.1+incompatible/go.mod h1:nmEj6Dob7S7YxXgwXpfOuvO54S+tGdZdw9fuRZt25Ag=
+github.com/circonus-labs/circonusllhist v0.1.3/go.mod h1:kMXHVDlOchFAehlya5ePtbp5jckzBHf4XRpQvBOLI+I=
 github.com/client9/misspell v0.3.4/go.mod h1:qj6jICC3Q7zFZvVWo7KLAzC3yx5G7kyvSDkc90ppPyw=
 github.com/cloudwego/base64x v0.1.4 h1:jwCgWpFanWmN8xoIUHa2rtzmkd5J2plF/dnLS6Xd/0Y=
 github.com/cloudwego/base64x v0.1.4/go.mod h1:0zlkT4Wn5C6NdauXdJRhSKRlJvmclQ1hhJgA0rcu/8w=
@@ -25,6 +43,11 @@ github.com/envoyproxy/go-control-plane v0.9.0/go.mod h1:YTl/9mNaCwkRvm6d1a2C3ymF
 github.com/envoyproxy/go-control-plane v0.9.1-0.20191026205805-5f8ba28d4473/go.mod h1:YTl/9mNaCwkRvm6d1a2C3ymFceY/DCBVvsKhRF0iEA4=
 github.com/envoyproxy/go-control-plane v0.9.4/go.mod h1:6rpuAdCZL397s3pYoYcLgu1mIlRU8Am5FuJP05cCM98=
 github.com/envoyproxy/protoc-gen-validate v0.1.0/go.mod h1:iSmxcyjqTsJpI2R4NaDN7+kN2VEUnK/pcBlmesArF7c=
+github.com/fatih/color v1.7.0/go.mod h1:Zm6kSWBoL9eyXnKyktHP6abPY2pDugNf5KwzbycvMj4=
+github.com/fatih/color v1.9.0/go.mod h1:eQcE1qtQxscV5RaZvpXrrb8Drkc3/DdQ+uUYCNjL+zU=
+github.com/fatih/color v1.13.0/go.mod h1:kLAiJbzzSOZDVNGyDpeOxJ47H46qBXwg5ILebYFFOfk=
+github.com/fatih/color v1.14.1 h1:qfhVLaG5s+nCROl1zJsZRxFeYrHLqWroPOQ8BWiNb4w=
+github.com/fatih/color v1.14.1/go.mod h1:2oHN61fhTpgcxD3TSWCgKDiH1+x4OiDVVGH8WlgGZGg=
 github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHkI4W8=
 github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
@@ -35,7 +58,11 @@ github.com/gin-contrib/sse v0.1.0 h1:Y/yl/+YNO8GZSjAhjMsSuLt29uWRFHdHYUb5lYOV9qE
 github.com/gin-contrib/sse v0.1.0/go.mod h1:RHrZQHXnP2xjPF+u1gW/2HnVO7nvIa9PG3Gm+fLHvGI=
 github.com/gin-gonic/gin v1.10.0 h1:nTuyha1TYqgedzytsKYqna+DfLos46nTv2ygFy86HFU=
 github.com/gin-gonic/gin v1.10.0/go.mod h1:4PMNQiOhvDRa013RKVbsiNwoyezlm2rm0uX/T7kzp5Y=
+github.com/go-kit/kit v0.8.0/go.mod h1:xBxKIO96dXMWWy0MnWVtmwkA9/13aqxPnvrjFYMA2as=
+github.com/go-kit/kit v0.9.0/go.mod h1:xBxKIO96dXMWWy0MnWVtmwkA9/13aqxPnvrjFYMA2as=
 github.com/go-kit/log v0.1.0/go.mod h1:zbhenjAZHb184qTLMA9ZjW7ThYL0H2mk7Q6pNt4vbaY=
+github.com/go-logfmt/logfmt v0.3.0/go.mod h1:Qt1PoO58o5twSAckw1HlFXLmHsOX5/0LbT9GBnD5lWE=
+github.com/go-logfmt/logfmt v0.4.0/go.mod h1:3RMwSq7FuexP4Kalkev3ejPJsZTpXXBr9+V4qmtdjCk=
 github.com/go-logfmt/logfmt v0.5.0/go.mod h1:wCYkCAKZfumFQihp8CzCvQ3paCTfi41vtzG1KdI/P7A=
 github.com/go-playground/assert/v2 v2.2.0 h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=
 github.com/go-playground/assert/v2 v2.2.0/go.mod h1:VDjEfimB/XKnb+ZQfWdccd7VUvScMdVu0Titje2rxJ4=
@@ -48,16 +75,23 @@ github.com/go-playground/validator/v10 v10.22.1/go.mod h1:dbuPbCMFw/DrkbEynArYaC
 github.com/go-stack/stack v1.8.0/go.mod h1:v0f6uXyyMGvRgIKkXu+yp6POWl0qKG85gN/melR3HDY=
 github.com/goccy/go-json v0.10.3 h1:KZ5WoDbxAIgm2HNbYckL0se1fHD6rz5j4ywS6ebzDqA=
 github.com/goccy/go-json v0.10.3/go.mod h1:oq7eo15ShAhp70Anwd5lgX2pLfOS3QCiwU/PULtXL6M=
+github.com/gogo/protobuf v1.1.1/go.mod h1:r8qH/GZQm5c6nD/R0oafs1akxWv10x8SbQlK7atdtwQ=
 github.com/gogo/protobuf v1.3.2 h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=
 github.com/gogo/protobuf v1.3.2/go.mod h1:P1XiOD3dCwIKUDQYPy72D8LYyHL2YPYrpS2s69NZV8Q=
 github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b/go.mod h1:SBH7ygxi8pfUlaOkMMuAQtPIUF8ecWP5IEl/CR7VP2Q=
 github.com/golang/mock v1.1.1/go.mod h1:oTYuIxOrZwtPieC+H1uAHpcLFnEyAGVDL/k47Jfbm0A=
 github.com/golang/protobuf v1.2.0/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5yJMmIC1U=
+github.com/golang/protobuf v1.3.1/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5yJMmIC1U=
 github.com/golang/protobuf v1.3.2/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5yJMmIC1U=
 github.com/golang/protobuf v1.3.3/go.mod h1:vzj43D7+SQXF/4pzW/hwtAqwc6iTitCiVSaWz5lYuqw=
 github.com/golang/protobuf v1.5.4 h1:i7eJL8qZTpSEXOPTxNKhASYpMn+8e5Q6AdndVa1dWek=
 github.com/golang/protobuf v1.5.4/go.mod h1:lnTiLA8Wa4RWRcIUkrtSVa5nRhsEGBg48fD6rSs7xps=
+github.com/google/btree v0.0.0-20180813153112-4030bb1f1f0c/go.mod h1:lNA+9X1NB3Zf8V7Ke586lFgjr2dZNuvo3lPJSGZ5JPQ=
+github.com/google/btree v1.0.1 h1:gK4Kx5IaGY9CD5sPJ36FHiBJ6ZXl0kilRiiCj+jdYp4=
+github.com/google/btree v1.0.1/go.mod h1:xXMiIv4Fb/0kKde4SpL7qlzvu5cMJDRkFDxJfI9uaxA=
 github.com/google/go-cmp v0.2.0/go.mod h1:oXzfMopK8JAjlY9xF4vHSVASa0yLyX7SntLO5aqRK0M=
+github.com/google/go-cmp v0.3.1/go.mod h1:8QqcDgzrUqlUb/G2PQTWiueGozuR1884gddMywk6iLU=
+github.com/google/go-cmp v0.4.0/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
 github.com/google/go-cmp v0.6.0 h1:ofyhxvXcZhMsU5ulbFiLKl/XBFqE1GSq7atu8tAmTRI=
 github.com/google/go-cmp v0.6.0/go.mod h1:17dUlkBOakJ0+DkrSSNjCkIjxS6bF9zb3elmeNGIjoY=
 github.com/google/gofuzz v1.0.0/go.mod h1:dBl0BpW6vV/+mYPU4Po3pmUjxk6FQPldtuIdl/M65Eg=
@@ -65,11 +99,58 @@ github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
 github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
 github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 h1:UH//fgunKIs4JdUbpDl1VZCDaL56wXCB/5+wF6uHfaI=
 github.com/grpc-ecosystem/go-grpc-middleware v1.4.0/go.mod h1:g5qyo/la0ALbONm6Vbp88Yd8NsDy6rZz+RcrMPxvld8=
+github.com/hashicorp/consul/api v1.28.2 h1:mXfkRHrpHN4YY3RqL09nXU1eHKLNiuAN4kHvDQ16k/8=
+github.com/hashicorp/consul/api v1.28.2/go.mod h1:KyzqzgMEya+IZPcD65YFoOVAgPpbfERu4I/tzG6/ueE=
+github.com/hashicorp/consul/sdk v0.16.0 h1:SE9m0W6DEfgIVCJX7xU+iv/hUl4m/nxqMTnCdMxDpJ8=
+github.com/hashicorp/consul/sdk v0.16.0/go.mod h1:7pxqqhqoaPqnBnzXD1StKed62LqJeClzVsUEy85Zr0A=
+github.com/hashicorp/errwrap v1.0.0/go.mod h1:YH+1FKiLXxHSkmPseP+kNlulaMuP3n2brvKWEqk/Jc4=
+github.com/hashicorp/errwrap v1.1.0 h1:OxrOeh75EUXMY8TBjag2fzXGZ40LB6IKw45YeGUDY2I=
+github.com/hashicorp/errwrap v1.1.0/go.mod h1:YH+1FKiLXxHSkmPseP+kNlulaMuP3n2brvKWEqk/Jc4=
+github.com/hashicorp/go-cleanhttp v0.5.0/go.mod h1:JpRdi6/HCYpAwUzNwuwqhbovhLtngrth3wmdIIUrZ80=
+github.com/hashicorp/go-cleanhttp v0.5.2 h1:035FKYIWjmULyFRBKPs8TBQoi0x6d9G4xc9neXJWAZQ=
+github.com/hashicorp/go-cleanhttp v0.5.2/go.mod h1:kO/YDlP8L1346E6Sodw+PrpBSV4/SoxCXGY6BqNFT48=
+github.com/hashicorp/go-hclog v1.5.0 h1:bI2ocEMgcVlz55Oj1xZNBsVi900c7II+fWDyV9o+13c=
+github.com/hashicorp/go-hclog v1.5.0/go.mod h1:W4Qnvbt70Wk/zYJryRzDRU/4r0kIg0PVHBcfoyhpF5M=
+github.com/hashicorp/go-immutable-radix v1.0.0/go.mod h1:0y9vanUI8NX6FsYoO3zeMjhV/C5i9g4Q3DwcSNZ4P60=
+github.com/hashicorp/go-immutable-radix v1.3.1 h1:DKHmCUm2hRBK510BaiZlwvpD40f8bJFeZnpfm2KLowc=
+github.com/hashicorp/go-immutable-radix v1.3.1/go.mod h1:0y9vanUI8NX6FsYoO3zeMjhV/C5i9g4Q3DwcSNZ4P60=
+github.com/hashicorp/go-msgpack v0.5.3/go.mod h1:ahLV/dePpqEmjfWmKiqvPkv/twdG7iPBM1vqhUKIvfM=
+github.com/hashicorp/go-msgpack v0.5.5 h1:i9R9JSrqIz0QVLz3sz+i3YJdT7TTSLcfLLzJi9aZTuI=
+github.com/hashicorp/go-msgpack v0.5.5/go.mod h1:ahLV/dePpqEmjfWmKiqvPkv/twdG7iPBM1vqhUKIvfM=
+github.com/hashicorp/go-multierror v1.0.0/go.mod h1:dHtQlpGsu+cZNNAkkCN/P3hoUDHhCYQXV3UM06sGGrk=
+github.com/hashicorp/go-multierror v1.1.0/go.mod h1:spPvp8C1qA32ftKqdAHm4hHTbPw+vmowP0z+KUhOZdA=
+github.com/hashicorp/go-multierror v1.1.1 h1:H5DkEtf6CXdFp0N0Em5UCwQpXMWke8IA0+lD48awMYo=
+github.com/hashicorp/go-multierror v1.1.1/go.mod h1:iw975J/qwKPdAO1clOe2L8331t/9/fmwbPZ6JB6eMoM=
+github.com/hashicorp/go-retryablehttp v0.5.3/go.mod h1:9B5zBasrRhHXnJnui7y6sL7es7NDiJgTc6Er0maI1Xs=
+github.com/hashicorp/go-rootcerts v1.0.2 h1:jzhAVGtqPKbwpyCPELlgNWhE1znq+qwJtW5Oi2viEzc=
+github.com/hashicorp/go-rootcerts v1.0.2/go.mod h1:pqUvnprVnM5bf7AOirdbb01K4ccR319Vf4pU3K5EGc8=
+github.com/hashicorp/go-sockaddr v1.0.0/go.mod h1:7Xibr9yA9JjQq1JpNB2Vw7kxv8xerXegt+ozgdvDeDU=
+github.com/hashicorp/go-sockaddr v1.0.2 h1:ztczhD1jLxIRjVejw8gFomI1BQZOe2WoVOu0SyteCQc=
+github.com/hashicorp/go-sockaddr v1.0.2/go.mod h1:rB4wwRAUzs07qva3c5SdrY/NEtAUjGlgmH/UkBUC97A=
+github.com/hashicorp/go-syslog v1.0.0/go.mod h1:qPfqrKkXGihmCqbJM2mZgkZGvKG1dFdvsLplgctolz4=
+github.com/hashicorp/go-uuid v1.0.0/go.mod h1:6SBZvOh/SIDV7/2o3Jml5SYk/TvGqwFJ/bN7x4byOro=
+github.com/hashicorp/go-uuid v1.0.1/go.mod h1:6SBZvOh/SIDV7/2o3Jml5SYk/TvGqwFJ/bN7x4byOro=
+github.com/hashicorp/go-uuid v1.0.3 h1:2gKiV6YVmrJ1i2CKKa9obLvRieoRGviZFL26PcT/Co8=
+github.com/hashicorp/go-uuid v1.0.3/go.mod h1:6SBZvOh/SIDV7/2o3Jml5SYk/TvGqwFJ/bN7x4byOro=
+github.com/hashicorp/go-version v1.2.1 h1:zEfKbn2+PDgroKdiOzqiE8rsmLqU2uwi5PB5pBJ3TkI=
+github.com/hashicorp/go-version v1.2.1/go.mod h1:fltr4n8CU8Ke44wwGCBoEymUuxUHl09ZGVZPK5anwXA=
+github.com/hashicorp/golang-lru v0.5.0/go.mod h1:/m3WP610KZHVQ1SGc6re/UDhFvYD7pJ4Ao+sR/qLZy8=
+github.com/hashicorp/golang-lru v0.5.4 h1:YDjusn29QI/Das2iO9M0BHnIbxPeyuCHsjMW+lJfyTc=
+github.com/hashicorp/golang-lru v0.5.4/go.mod h1:iADmTwqILo4mZ8BN3D2Q6+9jd8WM5uGBxy+E8yxSoD4=
 github.com/hashicorp/hcl v1.0.0 h1:0Anlzjpi4vEasTeNFn2mLJgTSwt0+6sfsiTG8qcWGx4=
 github.com/hashicorp/hcl v1.0.0/go.mod h1:E5yfLk+7swimpb2L/Alb/PJmXilQ/rhwaUYs4T20WEQ=
+github.com/hashicorp/logutils v1.0.0/go.mod h1:QIAnNjmIWmVIIkWDTG1z5v++HQmx9WQRO+LraFDTW64=
+github.com/hashicorp/mdns v1.0.4/go.mod h1:mtBihi+LeNXGtG8L9dX59gAEa12BDtBQSp4v/YAJqrc=
+github.com/hashicorp/memberlist v0.5.0 h1:EtYPN8DpAURiapus508I4n9CzHs2W+8NZGbmmR/prTM=
+github.com/hashicorp/memberlist v0.5.0/go.mod h1:yvyXLpo0QaGE59Y7hDTsTzDD25JYBZ4mHgHUZ8lrOI0=
+github.com/hashicorp/serf v0.10.1 h1:Z1H2J60yRKvfDYAOZLd2MU0ND4AH/WDz7xYHDWQsIPY=
+github.com/hashicorp/serf v0.10.1/go.mod h1:yL2t6BqATOLGc5HF7qbFkTfXoPIY0WZdWHfEvMqbG+4=
+github.com/json-iterator/go v1.1.6/go.mod h1:+SdeFBvtyEkXs7REEP0seUULqWtbJapLOCVDaaPEHmU=
+github.com/json-iterator/go v1.1.9/go.mod h1:KdQUCv79m/52Kvf8AW2vK1V8akMuk1QjK/uOdHXbAo4=
 github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
 github.com/json-iterator/go v1.1.12/go.mod h1:e30LSqwooZae/UwlEbR2852Gd8hjQvJoHmT4TnhNGBo=
 github.com/juju/gnuflag v0.0.0-20171113085948-2ce1bb71843d/go.mod h1:2PavIy+JPciBPrBUjwbNvtwB6RQlve+hkpll6QSNmOE=
+github.com/julienschmidt/httprouter v1.2.0/go.mod h1:SYymIcj16QtmaHHD7aYtjjsJG7VTCxuUUipMqKk8s4w=
 github.com/kisielk/errcheck v1.5.0/go.mod h1:pFxgyoBC7bSaBwPgfKdkLd5X25qrDl4LWUI2bnpBCr8=
 github.com/kisielk/gotool v1.0.0/go.mod h1:XhKaO+MFFWcvkIS/tQcRk01m1F5IRFswLeQ+oQHNcck=
 github.com/klauspost/cpuid/v2 v2.0.9/go.mod h1:FInQzS24/EEf25PyTYn52gqo7WaD8xa0213Md/qVLRg=
@@ -77,6 +158,7 @@ github.com/klauspost/cpuid/v2 v2.2.8 h1:+StwCXwm9PdpiEkPyzBXIy+M9KUb4ODm0Zarf1kS
 github.com/klauspost/cpuid/v2 v2.2.8/go.mod h1:Lcz8mBdAVJIBVzewtcLocK12l3Y+JytZYpaMropDUws=
 github.com/knz/go-libedit v1.10.1/go.mod h1:MZTVkCWyz0oBc7JOWP3wNAzd002ZbM/5hgShxwh4x8M=
 github.com/konsorten/go-windows-terminal-sequences v1.0.1/go.mod h1:T0+1ngSBFLxvqU3pZ+m/2kptfBszLMUkC4ZK/EgS/cQ=
+github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515/go.mod h1:+0opPa2QZZtGFBFZlji/RkVcI2GknAs/DXo4wKdlNEc=
 github.com/kr/pretty v0.1.0/go.mod h1:dAy3ld7l9f0ibDNOQOHHMYYIIbhfbHSm3C4ZsoJORNo=
 github.com/kr/pretty v0.3.1 h1:flRD4NNwYAUpkphVc1HcthR4KEIFJ65n8Mw5qdRn3LE=
 github.com/kr/pretty v0.3.1/go.mod h1:hoEshYVHaxMs3cyo3Yncou5ZscifuDolrwPKZanG3xk=
@@ -88,31 +170,78 @@ github.com/leodido/go-urn v1.4.0 h1:WT9HwE9SGECu3lg4d/dIA+jxlljEa1/ffXKmRjqdmIQ=
 github.com/leodido/go-urn v1.4.0/go.mod h1:bvxc+MVxLKB4z00jd1z+Dvzr47oO32F/QSNjSBOlFxI=
 github.com/magiconair/properties v1.8.7 h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=
 github.com/magiconair/properties v1.8.7/go.mod h1:Dhd985XPs7jluiymwWYZ0G4Z61jb3vdS329zhj2hYo0=
+github.com/mattn/go-colorable v0.0.9/go.mod h1:9vuHe8Xs5qXnSaW/c/ABM9alt+Vo+STaOChaDxuIBZU=
+github.com/mattn/go-colorable v0.1.4/go.mod h1:U0ppj6V5qS13XJ6of8GYAs25YV2eR4EVcfRqFIhoBtE=
+github.com/mattn/go-colorable v0.1.6/go.mod h1:u6P/XSegPjTcexA+o6vUJrdnUu04hMope9wVRipJSqc=
+github.com/mattn/go-colorable v0.1.9/go.mod h1:u6P/XSegPjTcexA+o6vUJrdnUu04hMope9wVRipJSqc=
+github.com/mattn/go-colorable v0.1.12/go.mod h1:u5H1YNBxpqRaxsYJYSkiCWKzEfiAb1Gb520KVy5xxl4=
+github.com/mattn/go-colorable v0.1.13 h1:fFA4WZxdEF4tXPZVKMLwD8oUnCTTo08duU7wxecdEvA=
+github.com/mattn/go-colorable v0.1.13/go.mod h1:7S9/ev0klgBDR4GtXTXX8a3vIGJpMovkB8vQcUbaXHg=
+github.com/mattn/go-isatty v0.0.3/go.mod h1:M+lRXTBqGeGNdLjl/ufCoiOlB5xdOkqRJdNxMWT7Zi4=
+github.com/mattn/go-isatty v0.0.8/go.mod h1:Iq45c/XA43vh69/j3iqttzPXn0bhXyGjM0Hdxcsrc5s=
+github.com/mattn/go-isatty v0.0.11/go.mod h1:PhnuNfih5lzO57/f3n+odYbM4JtupLOxQOAqxQCu2WE=
+github.com/mattn/go-isatty v0.0.12/go.mod h1:cbi8OIDigv2wuxKPP5vlRcQ1OAZbq2CE4Kysco4FUpU=
+github.com/mattn/go-isatty v0.0.14/go.mod h1:7GGIvUiUoEMVVmxf/4nioHXj79iQHKdU27kJ6hsGG94=
+github.com/mattn/go-isatty v0.0.16/go.mod h1:kYGgaQfpe5nmfYZH+SKPsOc2e4SrIfOl2e/yFXSvRLM=
 github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
 github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
+github.com/matttproud/golang_protobuf_extensions v1.0.1/go.mod h1:D8He9yQNgCq6Z5Ld7szi9bcBfOoFv/3dc6xSMkL2PC0=
+github.com/miekg/dns v1.1.26/go.mod h1:bPDLeHnStXmXAq1m/Ch/hvfNHr14JKNPMBo3VZKjuso=
+github.com/miekg/dns v1.1.41 h1:WMszZWJG0XmzbK9FEmzH2TVcqYzFesusSIB41b8KHxY=
+github.com/miekg/dns v1.1.41/go.mod h1:p6aan82bvRIyn+zDIv9xYNUpwa73JcSh9BKwknJysuI=
+github.com/mitchellh/cli v1.1.0/go.mod h1:xcISNoH86gajksDmfB23e/pu+B+GeFRMYmoHXxx3xhI=
+github.com/mitchellh/go-homedir v1.1.0 h1:lukF9ziXFxDFPkA1vsr5zpc1XuPDn/wFntq5mG+4E0Y=
+github.com/mitchellh/go-homedir v1.1.0/go.mod h1:SfyaCUpYCn1Vlf4IUYiD9fPX4A5wJrkLzIz1N1q0pr0=
+github.com/mitchellh/mapstructure v0.0.0-20160808181253-ca63d7c062ee/go.mod h1:FVVH3fgwuzCH5S8UJGiWEs2h04kUh9fWfEaFds41c1Y=
 github.com/mitchellh/mapstructure v1.5.0 h1:jeMsZIYE/09sWLaz43PL7Gy6RuMjD2eJVyuac5Z2hdY=
 github.com/mitchellh/mapstructure v1.5.0/go.mod h1:bFUtVrKA4DC2yAKiSyO/QUcy7e+RRV2QTWOzhPopBRo=
 github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
 github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=
 github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
+github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3RllmbCylyMrvgv0=
+github.com/modern-go/reflect2 v1.0.1/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3RllmbCylyMrvgv0=
 github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
 github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
+github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
 github.com/oapi-codegen/runtime v1.1.1 h1:EXLHh0DXIJnWhdRPN2w4MXAzFyE4CskzhNLUmtpMYro=
 github.com/oapi-codegen/runtime v1.1.1/go.mod h1:SK9X900oXmPWilYR5/WKPzt3Kqxn/uS/+lbpREv+eCg=
 github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFStlNbqXla1AfSYxPUl2o=
+github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
+github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
+github.com/pascaldekloe/goe v0.1.0/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pelletier/go-toml/v2 v2.2.3 h1:YmeHyLY8mFWbdkNWwpr+qIL2bEqT0o95WSdkNHvL12M=
 github.com/pelletier/go-toml/v2 v2.2.3/go.mod h1:MfCQTFTvCcUyyvvwm1+G6H/jORL20Xlb6rzQu9GuUkc=
+github.com/pkg/errors v0.8.0/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
 github.com/pkg/errors v0.8.1/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
+github.com/pkg/errors v0.9.1 h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=
+github.com/pkg/errors v0.9.1/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
 github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
 github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 h1:Jamvg5psRIccs7FGNTlIRMkT8wgtp5eCXdBlqhYGL6U=
 github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
+github.com/posener/complete v1.1.1/go.mod h1:em0nMJCgc9GFtwrmVmEMR/ZL6WyhyjMBndrE9hABlRI=
+github.com/posener/complete v1.2.3/go.mod h1:WZIdtGGp+qx0sLrYKtIRAruyNpv6hFCicSgv7Sy7s/s=
+github.com/prometheus/client_golang v0.9.1/go.mod h1:7SWBe2y4D6OKWSNQJUaRYU/AaXPKyh/dDVn+NZz0KFw=
+github.com/prometheus/client_golang v1.0.0/go.mod h1:db9x61etRT2tGnBNRi70OPL5FsnadC4Ky3P0J6CfImo=
+github.com/prometheus/client_golang v1.4.0/go.mod h1:e9GMxYsXl05ICDXkRhurwBS4Q3OK1iX/F2sw+iXX5zU=
+github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910/go.mod h1:MbSGuTsp3dbXC40dX6PRTWyKYBIrTGTE9sqQNg2J8bo=
+github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
 github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
+github.com/prometheus/client_model v0.2.0/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
+github.com/prometheus/common v0.4.1/go.mod h1:TNfzLD0ON7rHzMJeJkieUDPYmFC7Snx/y86RQel1bk4=
+github.com/prometheus/common v0.9.1/go.mod h1:yhUN8i9wzaXS3w1O07YhxHEBxD+W35wd8bs7vj7HSQ4=
+github.com/prometheus/procfs v0.0.0-20181005140218-185b4288413d/go.mod h1:c3At6R/oaqEKCNdg8wHV1ftS6bRYblBhIjjI8uT2IGk=
+github.com/prometheus/procfs v0.0.2/go.mod h1:TjEm7ze935MbeOT/UhFTIMYKhuLP4wbCsTZCD3I8kEA=
+github.com/prometheus/procfs v0.0.8/go.mod h1:7Qr8sr6344vo1JqZ6HhLceV9o3AJ1Ff+GxbHq6oeK9A=
 github.com/rogpeppe/go-internal v1.9.0 h1:73kH8U+JUqXU8lRuOHeVHaa/SZPifC7BkcraZVejAe8=
 github.com/rogpeppe/go-internal v1.9.0/go.mod h1:WtVeX8xhTBvf0smdhujwtBcq4Qrzq/fJaraNFVN+nFs=
+github.com/ryanuber/columnize v0.0.0-20160712163229-9b3edd62028f/go.mod h1:sm1tb6uqfes/u+d4ooFouqFdy9/2g9QGwK3SQygK0Ts=
 github.com/sagikazarmark/locafero v0.6.0 h1:ON7AQg37yzcRPU69mt7gwhFEBwxI6P9T4Qu3N51bwOk=
 github.com/sagikazarmark/locafero v0.6.0/go.mod h1:77OmuIc6VTraTXKXIs/uvUxKGUXjE1GbemJYHqdNjX0=
 github.com/sagikazarmark/slog-shim v0.1.0 h1:diDBnUNK9N/354PgrxMywXnAwEr1QZcOr6gto+ugjYE=
 github.com/sagikazarmark/slog-shim v0.1.0/go.mod h1:SrcSrq8aKtyuqEI1uvTDTK1arOWRIczQRv+GVI1AkeQ=
+github.com/sean-/seed v0.0.0-20170313163322-e2103e2c3529 h1:nn5Wsu0esKSJiIVhscUtVbo7ada43DJhG55ua/hjS5I=
+github.com/sean-/seed v0.0.0-20170313163322-e2103e2c3529/go.mod h1:DxrIzT+xaE7yg65j358z/aeFdxmN0P9QXhEzd20vsDc=
+github.com/sirupsen/logrus v1.2.0/go.mod h1:LxeOpSwHxABJmUn/MG1IvRgCAasNZTLOkJPxbbu5VWo=
 github.com/sirupsen/logrus v1.4.2/go.mod h1:tLMulIdttU9McNUspp0xgXVQah82FyeX6MwdIuYE2rE=
 github.com/sirupsen/logrus v1.9.3 h1:dueUQJ1C2q9oE3F7wvmSGAaVtTmUizReu6fjN8uqzbQ=
 github.com/sirupsen/logrus v1.9.3/go.mod h1:naHLuLoDiP4jHNo9R0sCBMtWGeIprob74mVsIT4qYEQ=
@@ -130,18 +259,21 @@ github.com/spkg/bom v0.0.0-20160624110644-59b7046e48ad/go.mod h1:qLr4V1qq6nMqFKk
 github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.1.1/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
+github.com/stretchr/objx v0.5.0 h1:1zr/of2m5FGMsad5YfcqgdqdWrIhu+EBEJRhR1U7z/c=
 github.com/stretchr/objx v0.5.0/go.mod h1:Yh+to48EsGEfYuaHDzXPcE3xhTkx73EhmCGUpEOglKo=
 github.com/stretchr/testify v1.2.2/go.mod h1:a8OnRcib4nhh0OaRAV+Yts87kKdq0PP7pXfy6kDkUVs=
 github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
 github.com/stretchr/testify v1.4.0/go.mod h1:j7eGeouHqKxXV5pUuKE4zz7dFj8WfuZ+81PSLYec5m4=
 github.com/stretchr/testify v1.7.0/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
 github.com/stretchr/testify v1.7.1/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
+github.com/stretchr/testify v1.7.2/go.mod h1:R6va5+xMeoiuVRoj+gSkQ7d3FALtqAAGI1FQKckRals=
 github.com/stretchr/testify v1.8.0/go.mod h1:yNjHg4UonilssWZ8iaSj1OCr/vHnekPRkoO+kdMU+MU=
 github.com/stretchr/testify v1.8.1/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o6fzry7u4=
 github.com/stretchr/testify v1.9.0 h1:HtqpIVDClZ4nwg75+f6Lvsy/wHu+3BoSGCbBAcpTsTg=
 github.com/stretchr/testify v1.9.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8C91i36aY=
 github.com/subosito/gotenv v1.6.0 h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=
 github.com/subosito/gotenv v1.6.0/go.mod h1:Dk4QP5c2W3ibzajGcXpNraDfq2IrhjMIvMSWPKKo0FU=
+github.com/tv42/httpunix v0.0.0-20150427012821-b75d8614f926/go.mod h1:9ESjWnEqriFuLhtthL60Sar/7RFoluCcXsuvEwTV5KM=
 github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS4MhqMhdFk5YI=
 github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
 github.com/ugorji/go/codec v1.2.12 h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65EE=
@@ -156,7 +288,9 @@ go.uber.org/multierr v1.11.0/go.mod h1:20+QtiLqy0Nd6FdQB9TLXag12DsQkrbs3htMFfDN8
 go.uber.org/zap v1.18.1/go.mod h1:xg/QME4nWcxGxrpdeYfq7UvYrLh66cuVKdrbD1XF/NI=
 golang.org/x/arch v0.11.0 h1:KXV8WWKCXm6tRpLirl2szsO5j/oOODwZf4hATmGVNs4=
 golang.org/x/arch v0.11.0/go.mod h1:FEVrYAQjsQXMVJ1nsMoVVXPZg6p2JE2mx8psSWTDQys=
+golang.org/x/crypto v0.0.0-20180904163835-0709b304e793/go.mod h1:6SG95UA2DQfeDnfUPMdvaQW0Q7yPrPDi9nlGo2tz2b4=
 golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACkg1iLfiJU5Ep61QUkGW8qpdssI0+w=
+golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392/go.mod h1:/lpIB1dKB+9EgE3H3cr1v9wB50oz8l4C4h62xy7jSTY=
 golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
 golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9/go.mod h1:LzIPMQfyMNhhGPhUkYOs5KpL4U8rLKemX1yGLhDgUto=
 golang.org/x/crypto v0.28.0 h1:GBDwsMXVQi34v5CCYUm2jkJvu4cbtru2U4TN2PSyQnw=
@@ -172,33 +306,61 @@ golang.org/x/mod v0.2.0/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
 golang.org/x/mod v0.3.0/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
 golang.org/x/net v0.0.0-20180724234803-3673e40ba225/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20180826012351-8a410e7b638d/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
+golang.org/x/net v0.0.0-20181114220301-adae6a3d119a/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20190213061140-3a22650c66bd/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20190311183353-d8887717615a/go.mod h1:t9HGtf8HONx5eT2rtn7q6eTqICYqUVnKs3thJo3Qplg=
 golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3/go.mod h1:t9HGtf8HONx5eT2rtn7q6eTqICYqUVnKs3thJo3Qplg=
+golang.org/x/net v0.0.0-20190613194153-d28f0bde5980/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
 golang.org/x/net v0.0.0-20190620200207-3b0461eec859/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
+golang.org/x/net v0.0.0-20190923162816-aa69164e4478/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
 golang.org/x/net v0.0.0-20200226121028-0de0cce0169b/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
 golang.org/x/net v0.0.0-20201021035429-f5854403a974/go.mod h1:sp8m0HH+o8qH0wwXwYZr8TS3Oi6o0r6Gce1SSxlDquU=
+golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
+golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1/go.mod h1:9tjilg8BloeKEkVJvy7fQ90B1CfIiPueXVOjqfkSzI8=
 golang.org/x/net v0.30.0 h1:AcW1SDZMkb8IpzCdQUaIq2sP4sZ4zw+55h6ynffypl4=
 golang.org/x/net v0.30.0/go.mod h1:2wGyMJ5iFasEhkwi13ChkO/t1ECNC4X4eBKkVFyYFlU=
 golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be/go.mod h1:N/0e6XlmueqKjAGxoOufVs8QHGRruUQn6yWY3a++T0U=
 golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20181108010431-42b317875d0f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
+golang.org/x/sync v0.0.0-20181221193216-37e7f081c4d4/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20190423024810-112230192c58/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
+golang.org/x/sync v0.0.0-20210220032951-036812b2e83c/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
+golang.org/x/sys v0.0.0-20180823144017-11551d06cbcc/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20180830151530-49385e6e1522/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
+golang.org/x/sys v0.0.0-20180905080454-ebe1bf3edb33/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
+golang.org/x/sys v0.0.0-20181116152217-5ac8a444bdc5/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
+golang.org/x/sys v0.0.0-20190222072716-a9d3bda3a223/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20190412213103-97732733099d/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20190422165155-953cdadca894/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20190922100055-0a153f010e69/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20190924154521-2837fb4f24fe/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20191026070338-33540a1f6037/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20200116001909-b77594299b42/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20200122134326-e047566fdf82/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20201119102817-f84b799fce68/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20210303074136-134d130e1a04/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20210330210617-4fbd30eecc44/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.26.0 h1:KHjCJyddX0LoSTb3J+vWpupP9p0oznkqVk/IfjymZbo=
 golang.org/x/sys v0.26.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
+golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
+golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
+golang.org/x/text v0.3.6/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
 golang.org/x/text v0.19.0 h1:kTxAhCbGbxhK0IwgSKiMO5awPoDQ0RpfiVYBfK860YM=
 golang.org/x/text v0.19.0/go.mod h1:BuEKDfySbSR4drPmRPG/7iBdf8hvFMuRexcpahXilzY=
 golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
@@ -206,6 +368,7 @@ golang.org/x/tools v0.0.0-20190114222345-bf090417da8b/go.mod h1:n7NCudcB/nEzxVGm
 golang.org/x/tools v0.0.0-20190226205152-f727befe758c/go.mod h1:9Yl7xja0Znq3iFh3HoIrodX9oNMXvdceNzlUR8zjMvY=
 golang.org/x/tools v0.0.0-20190311212946-11955173bddd/go.mod h1:LCzVGOaR6xXOjkQ3onu1FJEFr0SW1gC7cKk1uF8kGRs=
 golang.org/x/tools v0.0.0-20190524140312-2c0ae7006135/go.mod h1:RgjU9mgBXZiqYHBnxXauZ1Gv1EHHAz9KjViQ78xBX0Q=
+golang.org/x/tools v0.0.0-20190907020128-2ca718005c18/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
 golang.org/x/tools v0.0.0-20191108193012-7d206e10da11/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
 golang.org/x/tools v0.0.0-20191119224855-298f0cb1881e/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
 golang.org/x/tools v0.0.0-20200619180055-7c47624df98f/go.mod h1:EkVYQZoAsY45+roYkvgYkIh4xh/qjgUK9TdY2XT94GE=
@@ -230,13 +393,17 @@ google.golang.org/grpc v1.67.1 h1:zWnc1Vrcno+lHZCOofnIMvycFcc0QRGIzm9dhnDX68E=
 google.golang.org/grpc v1.67.1/go.mod h1:1gLDyUQU7CTLJI90u3nXZ9ekeghjeM7pTDZlqFNg2AA=
 google.golang.org/protobuf v1.35.1 h1:m3LfL6/Ca+fqnjnlqQXNpFPABW1UD7mjh8KO2mKFytA=
 google.golang.org/protobuf v1.35.1/go.mod h1:9fA7Ob0pmnwhb644+1+CVWFRbNajQ6iRojtC/QF5bRE=
+gopkg.in/alecthomas/kingpin.v2 v2.2.6/go.mod h1:FMv+mEhP44yOT+4EoQTLFTRgOQ1FBLkstjWtayDeSgw=
 gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 h1:YR8cESwS4TdDjEe65xsg0ogRM/Nc3DYOhEAlW+xobZo=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/ini.v1 v1.67.0 h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=
 gopkg.in/ini.v1 v1.67.0/go.mod h1:pNLf8WUiyNEtQjuu5G5vTm06TEv9tsIgeAvK8hOrP4k=
+gopkg.in/yaml.v2 v2.2.1/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.2/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
+gopkg.in/yaml.v2 v2.2.4/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
+gopkg.in/yaml.v2 v2.2.5/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.8/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
 gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
diff --git a/internal/order/http.go b/internal/order/http.go
index 8b5d70b..b40adc7 100644
--- a/internal/order/http.go
+++ b/internal/order/http.go
@@ -10,20 +10,16 @@ import (
 	"github.com/gin-gonic/gin"
 )
 
-// HTTPServer 是 OpenAPI 生成接口在 order 服务里的具体实现。
-// 它本身不处理复杂业务，只负责把 HTTP 请求翻译成应用层调用。
 type HTTPServer struct {
 	app app.Application
 }
 
 func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
 	var req orderpb.CreateOrderRequest
-	// ShouldBindJSON 负责把请求体反序列化为 proto 请求对象。
 	if err := c.ShouldBindJSON(&req); err != nil {
 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
 		return
 	}
-	// 进入应用层前，把 HTTP 层对象转换成 command 对象，避免业务逻辑依赖 gin。
 	r, err := H.app.Commands.CreateOrder.Handle(c, command.CreateOrder{
 		CustomerID: req.CustomerID,
 		Items:      req.Items,
@@ -40,7 +36,6 @@ func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID stri
 }
 
 func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
-	// 查询场景走 Queries 分组，体现读写分离的组织方式。
 	o, err := H.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
 		OrderID:    orderID,
 		CustomerID: customerID,
diff --git a/internal/order/main.go b/internal/order/main.go
index 9807588..cb08cc4 100644
--- a/internal/order/main.go
+++ b/internal/order/main.go
@@ -4,6 +4,7 @@ import (
 	"context"
 
 	"github.com/ghost-yu/go_shop_second/common/config"
+	"github.com/ghost-yu/go_shop_second/common/discovery"
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 	"github.com/ghost-yu/go_shop_second/common/server"
 	"github.com/ghost-yu/go_shop_second/order/ports"
@@ -14,7 +15,6 @@ import (
 	"google.golang.org/grpc"
 )
 
-// init 先加载全局配置，这样 main 里启动服务时就能直接从 viper 读取地址和服务名。
 func init() {
 	if err := config.NewViperConfig(); err != nil {
 		logrus.Fatal(err)
@@ -22,24 +22,27 @@ func init() {
 }
 
 func main() {
-	// serviceName 对应 global.yaml 里的 order.service-name，后续 server 包会继续用它找端口配置。
 	serviceName := viper.GetString("order.service-name")
 
-	// ctx 用来把进程级生命周期传给下游依赖，例如 gRPC 客户端连接。
 	ctx, cancel := context.WithCancel(context.Background())
 	defer cancel()
 
-	// NewApplication 是组合根：在这里把 repo、远程 client、handler 全部装配好。
 	application, cleanup := service.NewApplication(ctx)
 	defer cleanup()
 
-	// gRPC 和 HTTP 同时启动，说明 order 服务对内和对外用了两套协议暴露能力。
+	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	defer func() {
+		_ = deregisterFunc()
+	}()
+
 	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
 		svc := ports.NewGRPCServer(application)
 		orderpb.RegisterOrderServiceServer(server, svc)
 	})
 
-	// HTTP 侧通过 OpenAPI 生成的 RegisterHandlersWithOptions 完成路由到实现的映射。
 	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
 		ports.RegisterHandlersWithOptions(router, HTTPServer{
 			app: application,
diff --git a/internal/order/ports/grpc.go b/internal/order/ports/grpc.go
index 4841586..e1e8621 100644
--- a/internal/order/ports/grpc.go
+++ b/internal/order/ports/grpc.go
@@ -8,8 +8,6 @@ import (
 	"google.golang.org/protobuf/types/known/emptypb"
 )
 
-// GRPCServer 是 orderpb.OrderServiceServer 的端口适配器实现。
-// 这一层的职责和 HTTPServer 一样，都是把外部协议请求转给应用层。
 type GRPCServer struct {
 	app app.Application
 }
@@ -19,19 +17,16 @@ func NewGRPCServer(app app.Application) *GRPCServer {
 }
 
 func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
-	// lesson 当前故意保留 TODO，方便后续逐步实现 gRPC 版本的下单流程。
 	//TODO implement me
 	panic("implement me")
 }
 
 func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
-	// 这里最终会把 request 转为 query，再返回 proto 层定义的 orderpb.Order。
 	//TODO implement me
 	panic("implement me")
 }
 
 func (G GRPCServer) UpdateOrder(ctx context.Context, order *orderpb.Order) (*emptypb.Empty, error) {
-	// 更新接口通常会调用 Commands.UpdateOrder；当前阶段先保留占位。
 	//TODO implement me
 	panic("implement me")
 }
diff --git a/internal/order/service/application.go b/internal/order/service/application.go
index faf1f24..648e438 100644
--- a/internal/order/service/application.go
+++ b/internal/order/service/application.go
@@ -13,14 +13,11 @@ import (
 	"github.com/sirupsen/logrus"
 )
 
-// NewApplication 是 order 服务的组合根。
-// 它负责创建外部依赖、构造适配器，并返回一个已经装配好的应用层门面。
 func NewApplication(ctx context.Context) (app.Application, func()) {
 	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
 	if err != nil {
 		panic(err)
 	}
-	// 这里把生成的 gRPC client 再包一层适配器，避免应用层直接依赖 proto 细节。
 	stockGRPC := grpc.NewStockGRPC(stockClient)
 	return newApplication(ctx, stockGRPC), func() {
 		_ = closeStockClient()
@@ -28,7 +25,6 @@ func NewApplication(ctx context.Context) (app.Application, func()) {
 }
 
 func newApplication(_ context.Context, stockGRPC query.StockService) app.Application {
-	// 组合根里统一决定“接口对应哪种实现”，后面要切数据库或 mock 时只改这里。
 	orderRepo := adapters.NewMemoryOrderRepository()
 	logger := logrus.NewEntry(logrus.StandardLogger())
 	metricClient := metrics.TodoMetrics{}
diff --git a/internal/stock/adapters/stock_inmem_repository.go b/internal/stock/adapters/stock_inmem_repository.go
index c38b115..f4aed23 100644
--- a/internal/stock/adapters/stock_inmem_repository.go
+++ b/internal/stock/adapters/stock_inmem_repository.go
@@ -8,14 +8,11 @@ import (
 	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
 )
 
-// MemoryStockRepository 用 map 模拟库存数据源。
-// 相比切片，map 更适合按商品 ID 做快速查找。
 type MemoryStockRepository struct {
 	lock  *sync.RWMutex
 	store map[string]*orderpb.Item
 }
 
-// stub 是教学阶段的假数据，用来模拟一个已经存在的库存商品。
 var stub = map[string]*orderpb.Item{
 	"item_id": {
 		ID:       "foo_item",
@@ -23,6 +20,24 @@ var stub = map[string]*orderpb.Item{
 		Quantity: 10000,
 		PriceID:  "stub_item_price_id",
 	},
+	"item1": {
+		ID:       "item1",
+		Name:     "stub item 1",
+		Quantity: 10000,
+		PriceID:  "stub_item1_price_id",
+	},
+	"item2": {
+		ID:       "item2",
+		Name:     "stub item 2",
+		Quantity: 10000,
+		PriceID:  "stub_item2_price_id",
+	},
+	"item3": {
+		ID:       "item3",
+		Name:     "stub item 3",
+		Quantity: 10000,
+		PriceID:  "stub_item3_price_id",
+	},
 }
 
 func NewMemoryStockRepository() *MemoryStockRepository {
@@ -33,7 +48,6 @@ func NewMemoryStockRepository() *MemoryStockRepository {
 }
 
 func (m MemoryStockRepository) GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error) {
-	// 这里只有读场景，所以拿读锁即可。
 	m.lock.RLock()
 	defer m.lock.RUnlock()
 	var (
@@ -50,6 +64,5 @@ func (m MemoryStockRepository) GetItems(ctx context.Context, ids []string) ([]*o
 	if len(res) == len(ids) {
 		return res, nil
 	}
-	// 返回“部分结果 + 缺失错误”可以帮助上层更明确地决定如何提示用户。
 	return res, domain.NotFoundError{Missing: missing}
 }
diff --git a/internal/stock/app/app.go b/internal/stock/app/app.go
index c159f54..94a18fb 100644
--- a/internal/stock/app/app.go
+++ b/internal/stock/app/app.go
@@ -1,14 +1,15 @@
 package app
 
-// Application 预留给 stock 服务聚合 Commands 和 Queries。
-// 当前 lesson 里 stock 逻辑还比较薄，所以两个分组都是空壳结构。
+import "github.com/ghost-yu/go_shop_second/stock/app/query"
+
 type Application struct {
 	Commands Commands
 	Queries  Queries
 }
 
-// Commands 未来承载库存写操作，例如扣减库存。
 type Commands struct{}
 
-// Queries 未来承载库存读操作，例如按商品 ID 查询库存。
-type Queries struct{}
+type Queries struct {
+	CheckIfItemsInStock query.CheckIfItemsInStockHandler
+	GetItems            query.GetItemsHandler
+}
diff --git a/internal/stock/app/query/check_if_items_in_stock.go b/internal/stock/app/query/check_if_items_in_stock.go
new file mode 100644
index 0000000..d1078f0
--- /dev/null
+++ b/internal/stock/app/query/check_if_items_in_stock.go
@@ -0,0 +1,46 @@
+package query
+
+import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/common/decorator"
+	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
+	"github.com/sirupsen/logrus"
+)
+
+type CheckIfItemsInStock struct {
+	Items []*orderpb.ItemWithQuantity
+}
+
+type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*orderpb.Item]
+
+type checkIfItemsInStockHandler struct {
+	stockRepo domain.Repository
+}
+
+func NewCheckIfItemsInStockHandler(
+	stockRepo domain.Repository,
+	logger *logrus.Entry,
+	metricClient decorator.MetricsClient,
+) CheckIfItemsInStockHandler {
+	if stockRepo == nil {
+		panic("nil stockRepo")
+	}
+	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*orderpb.Item](
+		checkIfItemsInStockHandler{stockRepo: stockRepo},
+		logger,
+		metricClient,
+	)
+}
+
+func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*orderpb.Item, error) {
+	var res []*orderpb.Item
+	for _, i := range query.Items {
+		res = append(res, &orderpb.Item{
+			ID:       i.ID,
+			Quantity: i.Quantity,
+		})
+	}
+	return res, nil
+}
diff --git a/internal/stock/app/query/get_items.go b/internal/stock/app/query/get_items.go
new file mode 100644
index 0000000..063a6b0
--- /dev/null
+++ b/internal/stock/app/query/get_items.go
@@ -0,0 +1,43 @@
+package query
+
+import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/common/decorator"
+	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
+	"github.com/sirupsen/logrus"
+)
+
+type GetItems struct {
+	ItemIDs []string
+}
+
+type GetItemsHandler decorator.QueryHandler[GetItems, []*orderpb.Item]
+
+type getItemsHandler struct {
+	stockRepo domain.Repository
+}
+
+func NewGetItemsHandler(
+	stockRepo domain.Repository,
+	logger *logrus.Entry,
+	metricClient decorator.MetricsClient,
+) GetItemsHandler {
+	if stockRepo == nil {
+		panic("nil stockRepo")
+	}
+	return decorator.ApplyQueryDecorators[GetItems, []*orderpb.Item](
+		getItemsHandler{stockRepo: stockRepo},
+		logger,
+		metricClient,
+	)
+}
+
+func (g getItemsHandler) Handle(ctx context.Context, query GetItems) ([]*orderpb.Item, error) {
+	items, err := g.stockRepo.GetItems(ctx, query.ItemIDs)
+	if err != nil {
+		return nil, err
+	}
+	return items, nil
+}
diff --git a/internal/stock/domain/stock/repository.go b/internal/stock/domain/stock/repository.go
index b0f1796..7a58e6a 100644
--- a/internal/stock/domain/stock/repository.go
+++ b/internal/stock/domain/stock/repository.go
@@ -8,13 +8,10 @@ import (
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
 )
 
-// Repository 定义库存服务需要提供的最小数据访问能力。
-// order 服务只要通过应用层依赖这个接口，就不需要知道库存数据放在哪里。
 type Repository interface {
 	GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error)
 }
 
-// NotFoundError 用来明确告诉调用方哪些商品在库存里不存在。
 type NotFoundError struct {
 	Missing []string
 }
diff --git a/internal/stock/go.mod b/internal/stock/go.mod
index 83d4ef8..4a63762 100644
--- a/internal/stock/go.mod
+++ b/internal/stock/go.mod
@@ -6,15 +6,18 @@ replace github.com/ghost-yu/go_shop_second/common => ../common
 
 require (
 	github.com/ghost-yu/go_shop_second/common v0.0.0-00010101000000-000000000000
+	github.com/sirupsen/logrus v1.8.1
 	github.com/spf13/viper v1.19.0
-	google.golang.org/grpc v1.62.1
+	google.golang.org/grpc v1.67.1
 )
 
 require (
+	github.com/armon/go-metrics v0.4.1 // indirect
 	github.com/bytedance/sonic v1.11.6 // indirect
 	github.com/bytedance/sonic/loader v0.1.1 // indirect
 	github.com/cloudwego/base64x v0.1.4 // indirect
 	github.com/cloudwego/iasm v0.2.0 // indirect
+	github.com/fatih/color v1.14.1 // indirect
 	github.com/fsnotify/fsnotify v1.7.0 // indirect
 	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
 	github.com/gin-contrib/sse v0.1.0 // indirect
@@ -25,19 +28,29 @@ require (
 	github.com/goccy/go-json v0.10.2 // indirect
 	github.com/golang/protobuf v1.5.3 // indirect
 	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
+	github.com/hashicorp/consul/api v1.28.2 // indirect
+	github.com/hashicorp/errwrap v1.1.0 // indirect
+	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
+	github.com/hashicorp/go-hclog v1.5.0 // indirect
+	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
+	github.com/hashicorp/go-multierror v1.1.1 // indirect
+	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
+	github.com/hashicorp/golang-lru v0.5.4 // indirect
 	github.com/hashicorp/hcl v1.0.0 // indirect
+	github.com/hashicorp/serf v0.10.1 // indirect
 	github.com/json-iterator/go v1.1.12 // indirect
 	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
 	github.com/leodido/go-urn v1.4.0 // indirect
 	github.com/magiconair/properties v1.8.7 // indirect
+	github.com/mattn/go-colorable v0.1.13 // indirect
 	github.com/mattn/go-isatty v0.0.20 // indirect
+	github.com/mitchellh/go-homedir v1.1.0 // indirect
 	github.com/mitchellh/mapstructure v1.5.0 // indirect
 	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
 	github.com/modern-go/reflect2 v1.0.2 // indirect
 	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
 	github.com/sagikazarmark/locafero v0.4.0 // indirect
 	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
-	github.com/sirupsen/logrus v1.8.1 // indirect
 	github.com/sourcegraph/conc v0.3.0 // indirect
 	github.com/spf13/afero v1.11.0 // indirect
 	github.com/spf13/cast v1.6.0 // indirect
@@ -48,13 +61,13 @@ require (
 	go.uber.org/atomic v1.9.0 // indirect
 	go.uber.org/multierr v1.9.0 // indirect
 	golang.org/x/arch v0.8.0 // indirect
-	golang.org/x/crypto v0.23.0 // indirect
+	golang.org/x/crypto v0.26.0 // indirect
 	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
-	golang.org/x/net v0.25.0 // indirect
-	golang.org/x/sys v0.20.0 // indirect
-	golang.org/x/text v0.15.0 // indirect
-	google.golang.org/genproto/googleapis/rpc v0.0.0-20240314234333-6e1732d8331c // indirect
-	google.golang.org/protobuf v1.34.1 // indirect
+	golang.org/x/net v0.28.0 // indirect
+	golang.org/x/sys v0.24.0 // indirect
+	golang.org/x/text v0.17.0 // indirect
+	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
+	google.golang.org/protobuf v1.35.1 // indirect
 	gopkg.in/ini.v1 v1.67.0 // indirect
 	gopkg.in/yaml.v3 v3.0.1 // indirect
 )
diff --git a/internal/stock/go.sum b/internal/stock/go.sum
index b380353..59d9f63 100644
--- a/internal/stock/go.sum
+++ b/internal/stock/go.sum
@@ -1,11 +1,29 @@
 cloud.google.com/go v0.26.0/go.mod h1:aQUYkXzVsufM+DwF1aE+0xfcU+56JwCaLick0ClmMTw=
 github.com/BurntSushi/toml v0.3.1/go.mod h1:xHWCNGjB5oqiDr8zfno3MHue2Ht5sIBksp03qcyfWMU=
+github.com/DataDog/datadog-go v3.2.0+incompatible/go.mod h1:LButxg5PwREeZtORoXG3tL4fMGNddJ+vMq1mwgfaqoQ=
+github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
+github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
+github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
+github.com/alecthomas/units v0.0.0-20190717042225-c3de453c63f4/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
+github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e/go.mod h1:3U/XgcO3hCbHZ8TKRvWD2dDTCfh9M9ya+I9JpbB7O8o=
+github.com/armon/go-metrics v0.0.0-20180917152333-f0300d1749da/go.mod h1:Q73ZrmVTwzkszR9V5SSuryQ31EELlFMUz1kKyl939pY=
+github.com/armon/go-metrics v0.4.1 h1:hR91U9KYmb6bLBYLQjyM+3j+rcd/UhE+G78SFnF8gJA=
+github.com/armon/go-metrics v0.4.1/go.mod h1:E6amYzXo6aW1tqzoZGT755KkbgrJsSdpwZ+3JqfkOG4=
+github.com/armon/go-radix v0.0.0-20180808171621-7fddfc383310/go.mod h1:ufUuZ+zHj4x4TnLV4JWEpy2hxWSpsRywHrMgIH9cCH8=
+github.com/armon/go-radix v1.0.0/go.mod h1:ufUuZ+zHj4x4TnLV4JWEpy2hxWSpsRywHrMgIH9cCH8=
 github.com/benbjohnson/clock v1.1.0/go.mod h1:J11/hYXuz8f4ySSvYwY0FKfm+ezbsZBKZxNJlLklBHA=
+github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973/go.mod h1:Dwedo/Wpr24TaqPxmxbtue+5NUziq4I4S80YR8gNf3Q=
+github.com/beorn7/perks v1.0.0/go.mod h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+CedLV8=
+github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
+github.com/bgentry/speakeasy v0.1.0/go.mod h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=
 github.com/bytedance/sonic v1.11.6 h1:oUp34TzMlL+OY1OUWxHqsdkgC/Zfc85zGqw9siXjrc0=
 github.com/bytedance/sonic v1.11.6/go.mod h1:LysEHSvpvDySVdC2f87zGWf6CIKJcAvqab1ZaiQtds4=
 github.com/bytedance/sonic/loader v0.1.1 h1:c+e5Pt1k/cy5wMveRDyk2X4B9hF4g7an8N3zCYjJFNM=
 github.com/bytedance/sonic/loader v0.1.1/go.mod h1:ncP89zfokxS5LZrJxl5z0UJcsk4M4yY2JpfqGeCtNLU=
 github.com/census-instrumentation/opencensus-proto v0.2.1/go.mod h1:f6KPmirojxKA12rnyqOA5BBL4O983OfeGPqjHWSTneU=
+github.com/cespare/xxhash/v2 v2.1.1/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
+github.com/circonus-labs/circonus-gometrics v2.3.1+incompatible/go.mod h1:nmEj6Dob7S7YxXgwXpfOuvO54S+tGdZdw9fuRZt25Ag=
+github.com/circonus-labs/circonusllhist v0.1.3/go.mod h1:kMXHVDlOchFAehlya5ePtbp5jckzBHf4XRpQvBOLI+I=
 github.com/client9/misspell v0.3.4/go.mod h1:qj6jICC3Q7zFZvVWo7KLAzC3yx5G7kyvSDkc90ppPyw=
 github.com/cloudwego/base64x v0.1.4 h1:jwCgWpFanWmN8xoIUHa2rtzmkd5J2plF/dnLS6Xd/0Y=
 github.com/cloudwego/base64x v0.1.4/go.mod h1:0zlkT4Wn5C6NdauXdJRhSKRlJvmclQ1hhJgA0rcu/8w=
@@ -20,6 +38,11 @@ github.com/envoyproxy/go-control-plane v0.9.0/go.mod h1:YTl/9mNaCwkRvm6d1a2C3ymF
 github.com/envoyproxy/go-control-plane v0.9.1-0.20191026205805-5f8ba28d4473/go.mod h1:YTl/9mNaCwkRvm6d1a2C3ymFceY/DCBVvsKhRF0iEA4=
 github.com/envoyproxy/go-control-plane v0.9.4/go.mod h1:6rpuAdCZL397s3pYoYcLgu1mIlRU8Am5FuJP05cCM98=
 github.com/envoyproxy/protoc-gen-validate v0.1.0/go.mod h1:iSmxcyjqTsJpI2R4NaDN7+kN2VEUnK/pcBlmesArF7c=
+github.com/fatih/color v1.7.0/go.mod h1:Zm6kSWBoL9eyXnKyktHP6abPY2pDugNf5KwzbycvMj4=
+github.com/fatih/color v1.9.0/go.mod h1:eQcE1qtQxscV5RaZvpXrrb8Drkc3/DdQ+uUYCNjL+zU=
+github.com/fatih/color v1.13.0/go.mod h1:kLAiJbzzSOZDVNGyDpeOxJ47H46qBXwg5ILebYFFOfk=
+github.com/fatih/color v1.14.1 h1:qfhVLaG5s+nCROl1zJsZRxFeYrHLqWroPOQ8BWiNb4w=
+github.com/fatih/color v1.14.1/go.mod h1:2oHN61fhTpgcxD3TSWCgKDiH1+x4OiDVVGH8WlgGZGg=
 github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHkI4W8=
 github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
 github.com/fsnotify/fsnotify v1.7.0 h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=
@@ -30,7 +53,11 @@ github.com/gin-contrib/sse v0.1.0 h1:Y/yl/+YNO8GZSjAhjMsSuLt29uWRFHdHYUb5lYOV9qE
 github.com/gin-contrib/sse v0.1.0/go.mod h1:RHrZQHXnP2xjPF+u1gW/2HnVO7nvIa9PG3Gm+fLHvGI=
 github.com/gin-gonic/gin v1.10.0 h1:nTuyha1TYqgedzytsKYqna+DfLos46nTv2ygFy86HFU=
 github.com/gin-gonic/gin v1.10.0/go.mod h1:4PMNQiOhvDRa013RKVbsiNwoyezlm2rm0uX/T7kzp5Y=
+github.com/go-kit/kit v0.8.0/go.mod h1:xBxKIO96dXMWWy0MnWVtmwkA9/13aqxPnvrjFYMA2as=
+github.com/go-kit/kit v0.9.0/go.mod h1:xBxKIO96dXMWWy0MnWVtmwkA9/13aqxPnvrjFYMA2as=
 github.com/go-kit/log v0.1.0/go.mod h1:zbhenjAZHb184qTLMA9ZjW7ThYL0H2mk7Q6pNt4vbaY=
+github.com/go-logfmt/logfmt v0.3.0/go.mod h1:Qt1PoO58o5twSAckw1HlFXLmHsOX5/0LbT9GBnD5lWE=
+github.com/go-logfmt/logfmt v0.4.0/go.mod h1:3RMwSq7FuexP4Kalkev3ejPJsZTpXXBr9+V4qmtdjCk=
 github.com/go-logfmt/logfmt v0.5.0/go.mod h1:wCYkCAKZfumFQihp8CzCvQ3paCTfi41vtzG1KdI/P7A=
 github.com/go-playground/assert/v2 v2.2.0 h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=
 github.com/go-playground/assert/v2 v2.2.0/go.mod h1:VDjEfimB/XKnb+ZQfWdccd7VUvScMdVu0Titje2rxJ4=
@@ -43,27 +70,81 @@ github.com/go-playground/validator/v10 v10.20.0/go.mod h1:dbuPbCMFw/DrkbEynArYaC
 github.com/go-stack/stack v1.8.0/go.mod h1:v0f6uXyyMGvRgIKkXu+yp6POWl0qKG85gN/melR3HDY=
 github.com/goccy/go-json v0.10.2 h1:CrxCmQqYDkv1z7lO7Wbh2HN93uovUHgrECaO5ZrCXAU=
 github.com/goccy/go-json v0.10.2/go.mod h1:6MelG93GURQebXPDq3khkgXZkazVtN9CRI+MGFi0w8I=
+github.com/gogo/protobuf v1.1.1/go.mod h1:r8qH/GZQm5c6nD/R0oafs1akxWv10x8SbQlK7atdtwQ=
 github.com/gogo/protobuf v1.3.2 h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=
 github.com/gogo/protobuf v1.3.2/go.mod h1:P1XiOD3dCwIKUDQYPy72D8LYyHL2YPYrpS2s69NZV8Q=
 github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b/go.mod h1:SBH7ygxi8pfUlaOkMMuAQtPIUF8ecWP5IEl/CR7VP2Q=
 github.com/golang/mock v1.1.1/go.mod h1:oTYuIxOrZwtPieC+H1uAHpcLFnEyAGVDL/k47Jfbm0A=
 github.com/golang/protobuf v1.2.0/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5yJMmIC1U=
+github.com/golang/protobuf v1.3.1/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5yJMmIC1U=
 github.com/golang/protobuf v1.3.2/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5yJMmIC1U=
 github.com/golang/protobuf v1.3.3/go.mod h1:vzj43D7+SQXF/4pzW/hwtAqwc6iTitCiVSaWz5lYuqw=
 github.com/golang/protobuf v1.5.0/go.mod h1:FsONVRAS9T7sI+LIUmWTfcYkHO4aIWwzhcaSAoJOfIk=
 github.com/golang/protobuf v1.5.3 h1:KhyjKVUg7Usr/dYsdSqoFveMYd5ko72D+zANwlG1mmg=
 github.com/golang/protobuf v1.5.3/go.mod h1:XVQd3VNwM+JqD3oG2Ue2ip4fOMUkwXdXDdiuN0vRsmY=
+github.com/google/btree v0.0.0-20180813153112-4030bb1f1f0c/go.mod h1:lNA+9X1NB3Zf8V7Ke586lFgjr2dZNuvo3lPJSGZ5JPQ=
+github.com/google/btree v1.0.1 h1:gK4Kx5IaGY9CD5sPJ36FHiBJ6ZXl0kilRiiCj+jdYp4=
+github.com/google/btree v1.0.1/go.mod h1:xXMiIv4Fb/0kKde4SpL7qlzvu5cMJDRkFDxJfI9uaxA=
 github.com/google/go-cmp v0.2.0/go.mod h1:oXzfMopK8JAjlY9xF4vHSVASa0yLyX7SntLO5aqRK0M=
+github.com/google/go-cmp v0.3.1/go.mod h1:8QqcDgzrUqlUb/G2PQTWiueGozuR1884gddMywk6iLU=
+github.com/google/go-cmp v0.4.0/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
 github.com/google/go-cmp v0.5.5/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
 github.com/google/go-cmp v0.6.0 h1:ofyhxvXcZhMsU5ulbFiLKl/XBFqE1GSq7atu8tAmTRI=
 github.com/google/go-cmp v0.6.0/go.mod h1:17dUlkBOakJ0+DkrSSNjCkIjxS6bF9zb3elmeNGIjoY=
 github.com/google/gofuzz v1.0.0/go.mod h1:dBl0BpW6vV/+mYPU4Po3pmUjxk6FQPldtuIdl/M65Eg=
 github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 h1:UH//fgunKIs4JdUbpDl1VZCDaL56wXCB/5+wF6uHfaI=
 github.com/grpc-ecosystem/go-grpc-middleware v1.4.0/go.mod h1:g5qyo/la0ALbONm6Vbp88Yd8NsDy6rZz+RcrMPxvld8=
+github.com/hashicorp/consul/api v1.28.2 h1:mXfkRHrpHN4YY3RqL09nXU1eHKLNiuAN4kHvDQ16k/8=
+github.com/hashicorp/consul/api v1.28.2/go.mod h1:KyzqzgMEya+IZPcD65YFoOVAgPpbfERu4I/tzG6/ueE=
+github.com/hashicorp/consul/sdk v0.16.0 h1:SE9m0W6DEfgIVCJX7xU+iv/hUl4m/nxqMTnCdMxDpJ8=
+github.com/hashicorp/consul/sdk v0.16.0/go.mod h1:7pxqqhqoaPqnBnzXD1StKed62LqJeClzVsUEy85Zr0A=
+github.com/hashicorp/errwrap v1.0.0/go.mod h1:YH+1FKiLXxHSkmPseP+kNlulaMuP3n2brvKWEqk/Jc4=
+github.com/hashicorp/errwrap v1.1.0 h1:OxrOeh75EUXMY8TBjag2fzXGZ40LB6IKw45YeGUDY2I=
+github.com/hashicorp/errwrap v1.1.0/go.mod h1:YH+1FKiLXxHSkmPseP+kNlulaMuP3n2brvKWEqk/Jc4=
+github.com/hashicorp/go-cleanhttp v0.5.0/go.mod h1:JpRdi6/HCYpAwUzNwuwqhbovhLtngrth3wmdIIUrZ80=
+github.com/hashicorp/go-cleanhttp v0.5.2 h1:035FKYIWjmULyFRBKPs8TBQoi0x6d9G4xc9neXJWAZQ=
+github.com/hashicorp/go-cleanhttp v0.5.2/go.mod h1:kO/YDlP8L1346E6Sodw+PrpBSV4/SoxCXGY6BqNFT48=
+github.com/hashicorp/go-hclog v1.5.0 h1:bI2ocEMgcVlz55Oj1xZNBsVi900c7II+fWDyV9o+13c=
+github.com/hashicorp/go-hclog v1.5.0/go.mod h1:W4Qnvbt70Wk/zYJryRzDRU/4r0kIg0PVHBcfoyhpF5M=
+github.com/hashicorp/go-immutable-radix v1.0.0/go.mod h1:0y9vanUI8NX6FsYoO3zeMjhV/C5i9g4Q3DwcSNZ4P60=
+github.com/hashicorp/go-immutable-radix v1.3.1 h1:DKHmCUm2hRBK510BaiZlwvpD40f8bJFeZnpfm2KLowc=
+github.com/hashicorp/go-immutable-radix v1.3.1/go.mod h1:0y9vanUI8NX6FsYoO3zeMjhV/C5i9g4Q3DwcSNZ4P60=
+github.com/hashicorp/go-msgpack v0.5.3/go.mod h1:ahLV/dePpqEmjfWmKiqvPkv/twdG7iPBM1vqhUKIvfM=
+github.com/hashicorp/go-msgpack v0.5.5 h1:i9R9JSrqIz0QVLz3sz+i3YJdT7TTSLcfLLzJi9aZTuI=
+github.com/hashicorp/go-msgpack v0.5.5/go.mod h1:ahLV/dePpqEmjfWmKiqvPkv/twdG7iPBM1vqhUKIvfM=
+github.com/hashicorp/go-multierror v1.0.0/go.mod h1:dHtQlpGsu+cZNNAkkCN/P3hoUDHhCYQXV3UM06sGGrk=
+github.com/hashicorp/go-multierror v1.1.0/go.mod h1:spPvp8C1qA32ftKqdAHm4hHTbPw+vmowP0z+KUhOZdA=
+github.com/hashicorp/go-multierror v1.1.1 h1:H5DkEtf6CXdFp0N0Em5UCwQpXMWke8IA0+lD48awMYo=
+github.com/hashicorp/go-multierror v1.1.1/go.mod h1:iw975J/qwKPdAO1clOe2L8331t/9/fmwbPZ6JB6eMoM=
+github.com/hashicorp/go-retryablehttp v0.5.3/go.mod h1:9B5zBasrRhHXnJnui7y6sL7es7NDiJgTc6Er0maI1Xs=
+github.com/hashicorp/go-rootcerts v1.0.2 h1:jzhAVGtqPKbwpyCPELlgNWhE1znq+qwJtW5Oi2viEzc=
+github.com/hashicorp/go-rootcerts v1.0.2/go.mod h1:pqUvnprVnM5bf7AOirdbb01K4ccR319Vf4pU3K5EGc8=
+github.com/hashicorp/go-sockaddr v1.0.0/go.mod h1:7Xibr9yA9JjQq1JpNB2Vw7kxv8xerXegt+ozgdvDeDU=
+github.com/hashicorp/go-sockaddr v1.0.2 h1:ztczhD1jLxIRjVejw8gFomI1BQZOe2WoVOu0SyteCQc=
+github.com/hashicorp/go-sockaddr v1.0.2/go.mod h1:rB4wwRAUzs07qva3c5SdrY/NEtAUjGlgmH/UkBUC97A=
+github.com/hashicorp/go-syslog v1.0.0/go.mod h1:qPfqrKkXGihmCqbJM2mZgkZGvKG1dFdvsLplgctolz4=
+github.com/hashicorp/go-uuid v1.0.0/go.mod h1:6SBZvOh/SIDV7/2o3Jml5SYk/TvGqwFJ/bN7x4byOro=
+github.com/hashicorp/go-uuid v1.0.1/go.mod h1:6SBZvOh/SIDV7/2o3Jml5SYk/TvGqwFJ/bN7x4byOro=
+github.com/hashicorp/go-uuid v1.0.3 h1:2gKiV6YVmrJ1i2CKKa9obLvRieoRGviZFL26PcT/Co8=
+github.com/hashicorp/go-uuid v1.0.3/go.mod h1:6SBZvOh/SIDV7/2o3Jml5SYk/TvGqwFJ/bN7x4byOro=
+github.com/hashicorp/go-version v1.2.1 h1:zEfKbn2+PDgroKdiOzqiE8rsmLqU2uwi5PB5pBJ3TkI=
+github.com/hashicorp/go-version v1.2.1/go.mod h1:fltr4n8CU8Ke44wwGCBoEymUuxUHl09ZGVZPK5anwXA=
+github.com/hashicorp/golang-lru v0.5.0/go.mod h1:/m3WP610KZHVQ1SGc6re/UDhFvYD7pJ4Ao+sR/qLZy8=
+github.com/hashicorp/golang-lru v0.5.4 h1:YDjusn29QI/Das2iO9M0BHnIbxPeyuCHsjMW+lJfyTc=
+github.com/hashicorp/golang-lru v0.5.4/go.mod h1:iADmTwqILo4mZ8BN3D2Q6+9jd8WM5uGBxy+E8yxSoD4=
 github.com/hashicorp/hcl v1.0.0 h1:0Anlzjpi4vEasTeNFn2mLJgTSwt0+6sfsiTG8qcWGx4=
 github.com/hashicorp/hcl v1.0.0/go.mod h1:E5yfLk+7swimpb2L/Alb/PJmXilQ/rhwaUYs4T20WEQ=
+github.com/hashicorp/logutils v1.0.0/go.mod h1:QIAnNjmIWmVIIkWDTG1z5v++HQmx9WQRO+LraFDTW64=
+github.com/hashicorp/mdns v1.0.4/go.mod h1:mtBihi+LeNXGtG8L9dX59gAEa12BDtBQSp4v/YAJqrc=
+github.com/hashicorp/memberlist v0.5.0 h1:EtYPN8DpAURiapus508I4n9CzHs2W+8NZGbmmR/prTM=
+github.com/hashicorp/memberlist v0.5.0/go.mod h1:yvyXLpo0QaGE59Y7hDTsTzDD25JYBZ4mHgHUZ8lrOI0=
+github.com/hashicorp/serf v0.10.1 h1:Z1H2J60yRKvfDYAOZLd2MU0ND4AH/WDz7xYHDWQsIPY=
+github.com/hashicorp/serf v0.10.1/go.mod h1:yL2t6BqATOLGc5HF7qbFkTfXoPIY0WZdWHfEvMqbG+4=
+github.com/json-iterator/go v1.1.6/go.mod h1:+SdeFBvtyEkXs7REEP0seUULqWtbJapLOCVDaaPEHmU=
+github.com/json-iterator/go v1.1.9/go.mod h1:KdQUCv79m/52Kvf8AW2vK1V8akMuk1QjK/uOdHXbAo4=
 github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
 github.com/json-iterator/go v1.1.12/go.mod h1:e30LSqwooZae/UwlEbR2852Gd8hjQvJoHmT4TnhNGBo=
+github.com/julienschmidt/httprouter v1.2.0/go.mod h1:SYymIcj16QtmaHHD7aYtjjsJG7VTCxuUUipMqKk8s4w=
 github.com/kisielk/errcheck v1.5.0/go.mod h1:pFxgyoBC7bSaBwPgfKdkLd5X25qrDl4LWUI2bnpBCr8=
 github.com/kisielk/gotool v1.0.0/go.mod h1:XhKaO+MFFWcvkIS/tQcRk01m1F5IRFswLeQ+oQHNcck=
 github.com/klauspost/cpuid/v2 v2.0.9/go.mod h1:FInQzS24/EEf25PyTYn52gqo7WaD8xa0213Md/qVLRg=
@@ -71,6 +152,7 @@ github.com/klauspost/cpuid/v2 v2.2.7 h1:ZWSB3igEs+d0qvnxR/ZBzXVmxkgt8DdzP6m9pfuV
 github.com/klauspost/cpuid/v2 v2.2.7/go.mod h1:Lcz8mBdAVJIBVzewtcLocK12l3Y+JytZYpaMropDUws=
 github.com/knz/go-libedit v1.10.1/go.mod h1:MZTVkCWyz0oBc7JOWP3wNAzd002ZbM/5hgShxwh4x8M=
 github.com/konsorten/go-windows-terminal-sequences v1.0.1/go.mod h1:T0+1ngSBFLxvqU3pZ+m/2kptfBszLMUkC4ZK/EgS/cQ=
+github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515/go.mod h1:+0opPa2QZZtGFBFZlji/RkVcI2GknAs/DXo4wKdlNEc=
 github.com/kr/pretty v0.1.0/go.mod h1:dAy3ld7l9f0ibDNOQOHHMYYIIbhfbHSm3C4ZsoJORNo=
 github.com/kr/pretty v0.3.1 h1:flRD4NNwYAUpkphVc1HcthR4KEIFJ65n8Mw5qdRn3LE=
 github.com/kr/pretty v0.3.1/go.mod h1:hoEshYVHaxMs3cyo3Yncou5ZscifuDolrwPKZanG3xk=
@@ -82,29 +164,76 @@ github.com/leodido/go-urn v1.4.0 h1:WT9HwE9SGECu3lg4d/dIA+jxlljEa1/ffXKmRjqdmIQ=
 github.com/leodido/go-urn v1.4.0/go.mod h1:bvxc+MVxLKB4z00jd1z+Dvzr47oO32F/QSNjSBOlFxI=
 github.com/magiconair/properties v1.8.7 h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=
 github.com/magiconair/properties v1.8.7/go.mod h1:Dhd985XPs7jluiymwWYZ0G4Z61jb3vdS329zhj2hYo0=
+github.com/mattn/go-colorable v0.0.9/go.mod h1:9vuHe8Xs5qXnSaW/c/ABM9alt+Vo+STaOChaDxuIBZU=
+github.com/mattn/go-colorable v0.1.4/go.mod h1:U0ppj6V5qS13XJ6of8GYAs25YV2eR4EVcfRqFIhoBtE=
+github.com/mattn/go-colorable v0.1.6/go.mod h1:u6P/XSegPjTcexA+o6vUJrdnUu04hMope9wVRipJSqc=
+github.com/mattn/go-colorable v0.1.9/go.mod h1:u6P/XSegPjTcexA+o6vUJrdnUu04hMope9wVRipJSqc=
+github.com/mattn/go-colorable v0.1.12/go.mod h1:u5H1YNBxpqRaxsYJYSkiCWKzEfiAb1Gb520KVy5xxl4=
+github.com/mattn/go-colorable v0.1.13 h1:fFA4WZxdEF4tXPZVKMLwD8oUnCTTo08duU7wxecdEvA=
+github.com/mattn/go-colorable v0.1.13/go.mod h1:7S9/ev0klgBDR4GtXTXX8a3vIGJpMovkB8vQcUbaXHg=
+github.com/mattn/go-isatty v0.0.3/go.mod h1:M+lRXTBqGeGNdLjl/ufCoiOlB5xdOkqRJdNxMWT7Zi4=
+github.com/mattn/go-isatty v0.0.8/go.mod h1:Iq45c/XA43vh69/j3iqttzPXn0bhXyGjM0Hdxcsrc5s=
+github.com/mattn/go-isatty v0.0.11/go.mod h1:PhnuNfih5lzO57/f3n+odYbM4JtupLOxQOAqxQCu2WE=
+github.com/mattn/go-isatty v0.0.12/go.mod h1:cbi8OIDigv2wuxKPP5vlRcQ1OAZbq2CE4Kysco4FUpU=
+github.com/mattn/go-isatty v0.0.14/go.mod h1:7GGIvUiUoEMVVmxf/4nioHXj79iQHKdU27kJ6hsGG94=
+github.com/mattn/go-isatty v0.0.16/go.mod h1:kYGgaQfpe5nmfYZH+SKPsOc2e4SrIfOl2e/yFXSvRLM=
 github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
 github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
+github.com/matttproud/golang_protobuf_extensions v1.0.1/go.mod h1:D8He9yQNgCq6Z5Ld7szi9bcBfOoFv/3dc6xSMkL2PC0=
+github.com/miekg/dns v1.1.26/go.mod h1:bPDLeHnStXmXAq1m/Ch/hvfNHr14JKNPMBo3VZKjuso=
+github.com/miekg/dns v1.1.41 h1:WMszZWJG0XmzbK9FEmzH2TVcqYzFesusSIB41b8KHxY=
+github.com/miekg/dns v1.1.41/go.mod h1:p6aan82bvRIyn+zDIv9xYNUpwa73JcSh9BKwknJysuI=
+github.com/mitchellh/cli v1.1.0/go.mod h1:xcISNoH86gajksDmfB23e/pu+B+GeFRMYmoHXxx3xhI=
+github.com/mitchellh/go-homedir v1.1.0 h1:lukF9ziXFxDFPkA1vsr5zpc1XuPDn/wFntq5mG+4E0Y=
+github.com/mitchellh/go-homedir v1.1.0/go.mod h1:SfyaCUpYCn1Vlf4IUYiD9fPX4A5wJrkLzIz1N1q0pr0=
+github.com/mitchellh/mapstructure v0.0.0-20160808181253-ca63d7c062ee/go.mod h1:FVVH3fgwuzCH5S8UJGiWEs2h04kUh9fWfEaFds41c1Y=
 github.com/mitchellh/mapstructure v1.5.0 h1:jeMsZIYE/09sWLaz43PL7Gy6RuMjD2eJVyuac5Z2hdY=
 github.com/mitchellh/mapstructure v1.5.0/go.mod h1:bFUtVrKA4DC2yAKiSyO/QUcy7e+RRV2QTWOzhPopBRo=
 github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
 github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=
 github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
+github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3RllmbCylyMrvgv0=
+github.com/modern-go/reflect2 v1.0.1/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3RllmbCylyMrvgv0=
 github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
 github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
+github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
 github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFStlNbqXla1AfSYxPUl2o=
+github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
+github.com/pascaldekloe/goe v0.1.0 h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=
+github.com/pascaldekloe/goe v0.1.0/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
 github.com/pelletier/go-toml/v2 v2.2.2 h1:aYUidT7k73Pcl9nb2gScu7NSrKCSHIDE89b3+6Wq+LM=
 github.com/pelletier/go-toml/v2 v2.2.2/go.mod h1:1t835xjRzz80PqgE6HHgN2JOsmgYu/h4qDAS4n929Rs=
+github.com/pkg/errors v0.8.0/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
 github.com/pkg/errors v0.8.1/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
+github.com/pkg/errors v0.9.1 h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=
+github.com/pkg/errors v0.9.1/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
 github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
 github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 h1:Jamvg5psRIccs7FGNTlIRMkT8wgtp5eCXdBlqhYGL6U=
 github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
+github.com/posener/complete v1.1.1/go.mod h1:em0nMJCgc9GFtwrmVmEMR/ZL6WyhyjMBndrE9hABlRI=
+github.com/posener/complete v1.2.3/go.mod h1:WZIdtGGp+qx0sLrYKtIRAruyNpv6hFCicSgv7Sy7s/s=
+github.com/prometheus/client_golang v0.9.1/go.mod h1:7SWBe2y4D6OKWSNQJUaRYU/AaXPKyh/dDVn+NZz0KFw=
+github.com/prometheus/client_golang v1.0.0/go.mod h1:db9x61etRT2tGnBNRi70OPL5FsnadC4Ky3P0J6CfImo=
+github.com/prometheus/client_golang v1.4.0/go.mod h1:e9GMxYsXl05ICDXkRhurwBS4Q3OK1iX/F2sw+iXX5zU=
+github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910/go.mod h1:MbSGuTsp3dbXC40dX6PRTWyKYBIrTGTE9sqQNg2J8bo=
+github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
 github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
+github.com/prometheus/client_model v0.2.0/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
+github.com/prometheus/common v0.4.1/go.mod h1:TNfzLD0ON7rHzMJeJkieUDPYmFC7Snx/y86RQel1bk4=
+github.com/prometheus/common v0.9.1/go.mod h1:yhUN8i9wzaXS3w1O07YhxHEBxD+W35wd8bs7vj7HSQ4=
+github.com/prometheus/procfs v0.0.0-20181005140218-185b4288413d/go.mod h1:c3At6R/oaqEKCNdg8wHV1ftS6bRYblBhIjjI8uT2IGk=
+github.com/prometheus/procfs v0.0.2/go.mod h1:TjEm7ze935MbeOT/UhFTIMYKhuLP4wbCsTZCD3I8kEA=
+github.com/prometheus/procfs v0.0.8/go.mod h1:7Qr8sr6344vo1JqZ6HhLceV9o3AJ1Ff+GxbHq6oeK9A=
 github.com/rogpeppe/go-internal v1.9.0 h1:73kH8U+JUqXU8lRuOHeVHaa/SZPifC7BkcraZVejAe8=
 github.com/rogpeppe/go-internal v1.9.0/go.mod h1:WtVeX8xhTBvf0smdhujwtBcq4Qrzq/fJaraNFVN+nFs=
+github.com/ryanuber/columnize v0.0.0-20160712163229-9b3edd62028f/go.mod h1:sm1tb6uqfes/u+d4ooFouqFdy9/2g9QGwK3SQygK0Ts=
 github.com/sagikazarmark/locafero v0.4.0 h1:HApY1R9zGo4DBgr7dqsTH/JJxLTTsOt7u6keLGt6kNQ=
 github.com/sagikazarmark/locafero v0.4.0/go.mod h1:Pe1W6UlPYUk/+wc/6KFhbORCfqzgYEpgQ3O5fPuL3H4=
 github.com/sagikazarmark/slog-shim v0.1.0 h1:diDBnUNK9N/354PgrxMywXnAwEr1QZcOr6gto+ugjYE=
 github.com/sagikazarmark/slog-shim v0.1.0/go.mod h1:SrcSrq8aKtyuqEI1uvTDTK1arOWRIczQRv+GVI1AkeQ=
+github.com/sean-/seed v0.0.0-20170313163322-e2103e2c3529 h1:nn5Wsu0esKSJiIVhscUtVbo7ada43DJhG55ua/hjS5I=
+github.com/sean-/seed v0.0.0-20170313163322-e2103e2c3529/go.mod h1:DxrIzT+xaE7yg65j358z/aeFdxmN0P9QXhEzd20vsDc=
+github.com/sirupsen/logrus v1.2.0/go.mod h1:LxeOpSwHxABJmUn/MG1IvRgCAasNZTLOkJPxbbu5VWo=
 github.com/sirupsen/logrus v1.4.2/go.mod h1:tLMulIdttU9McNUspp0xgXVQah82FyeX6MwdIuYE2rE=
 github.com/sirupsen/logrus v1.8.1 h1:dJKuHgqk1NNQlqoA6BTlM1Wf9DOH3NBjQyu0h9+AZZE=
 github.com/sirupsen/logrus v1.8.1/go.mod h1:yWOB1SBYBC5VeMP7gHvWumXLIWorT60ONWic61uBYv0=
@@ -122,12 +251,14 @@ github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+
 github.com/stretchr/objx v0.1.1/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
 github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
 github.com/stretchr/objx v0.5.0/go.mod h1:Yh+to48EsGEfYuaHDzXPcE3xhTkx73EhmCGUpEOglKo=
+github.com/stretchr/objx v0.5.2 h1:xuMeJ0Sdp5ZMRXx/aWO6RZxdr3beISkG5/G/aIRr3pY=
 github.com/stretchr/objx v0.5.2/go.mod h1:FRsXN1f5AsAjCGJKqEizvkpNtU+EGNCLh3NxZ/8L+MA=
 github.com/stretchr/testify v1.2.2/go.mod h1:a8OnRcib4nhh0OaRAV+Yts87kKdq0PP7pXfy6kDkUVs=
 github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
 github.com/stretchr/testify v1.4.0/go.mod h1:j7eGeouHqKxXV5pUuKE4zz7dFj8WfuZ+81PSLYec5m4=
 github.com/stretchr/testify v1.7.0/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
 github.com/stretchr/testify v1.7.1/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
+github.com/stretchr/testify v1.7.2/go.mod h1:R6va5+xMeoiuVRoj+gSkQ7d3FALtqAAGI1FQKckRals=
 github.com/stretchr/testify v1.8.0/go.mod h1:yNjHg4UonilssWZ8iaSj1OCr/vHnekPRkoO+kdMU+MU=
 github.com/stretchr/testify v1.8.1/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o6fzry7u4=
 github.com/stretchr/testify v1.8.4/go.mod h1:sz/lmYIOXD/1dqDmKjjqLyZ2RngseejIcXlSw2iwfAo=
@@ -135,6 +266,7 @@ github.com/stretchr/testify v1.9.0 h1:HtqpIVDClZ4nwg75+f6Lvsy/wHu+3BoSGCbBAcpTsT
 github.com/stretchr/testify v1.9.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8C91i36aY=
 github.com/subosito/gotenv v1.6.0 h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=
 github.com/subosito/gotenv v1.6.0/go.mod h1:Dk4QP5c2W3ibzajGcXpNraDfq2IrhjMIvMSWPKKo0FU=
+github.com/tv42/httpunix v0.0.0-20150427012821-b75d8614f926/go.mod h1:9ESjWnEqriFuLhtthL60Sar/7RFoluCcXsuvEwTV5KM=
 github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS4MhqMhdFk5YI=
 github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
 github.com/ugorji/go/codec v1.2.12 h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65EE=
@@ -152,11 +284,13 @@ go.uber.org/zap v1.18.1/go.mod h1:xg/QME4nWcxGxrpdeYfq7UvYrLh66cuVKdrbD1XF/NI=
 golang.org/x/arch v0.0.0-20210923205945-b76863e36670/go.mod h1:5om86z9Hs0C8fWVUuoMHwpExlXzs5Tkyp9hOrfG7pp8=
 golang.org/x/arch v0.8.0 h1:3wRIsP3pM4yUptoR96otTUOXI367OS0+c9eeRi9doIc=
 golang.org/x/arch v0.8.0/go.mod h1:FEVrYAQjsQXMVJ1nsMoVVXPZg6p2JE2mx8psSWTDQys=
+golang.org/x/crypto v0.0.0-20180904163835-0709b304e793/go.mod h1:6SG95UA2DQfeDnfUPMdvaQW0Q7yPrPDi9nlGo2tz2b4=
 golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACkg1iLfiJU5Ep61QUkGW8qpdssI0+w=
+golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392/go.mod h1:/lpIB1dKB+9EgE3H3cr1v9wB50oz8l4C4h62xy7jSTY=
 golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
 golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9/go.mod h1:LzIPMQfyMNhhGPhUkYOs5KpL4U8rLKemX1yGLhDgUto=
-golang.org/x/crypto v0.23.0 h1:dIJU/v2J8Mdglj/8rJ6UUOM3Zc9zLZxVZwwxMooUSAI=
-golang.org/x/crypto v0.23.0/go.mod h1:CKFgDieR+mRhux2Lsu27y0fO304Db0wZe70UKqHu0v8=
+golang.org/x/crypto v0.26.0 h1:RrRspgV4mU+YwB4FYnuBoKsUapNIL5cohGAmSH3azsw=
+golang.org/x/crypto v0.26.0/go.mod h1:GY7jblb9wI+FOo5y8/S2oY4zWP07AkOJ4+jxCqdqn54=
 golang.org/x/exp v0.0.0-20190121172915-509febef88a4/go.mod h1:CJ0aWSM057203Lf6IL+f9T1iT9GByDxfZKAQTCR3kQA=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9 h1:GoHiUyI/Tp2nVkLI2mCxVkOjsbSXD66ic0XW0js0R9g=
 golang.org/x/exp v0.0.0-20230905200255-921286631fa9/go.mod h1:S2oDrQGGwySpoQPVqRShND87VCbxmc6bL1Yd2oYrm6k=
@@ -168,40 +302,68 @@ golang.org/x/mod v0.2.0/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
 golang.org/x/mod v0.3.0/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
 golang.org/x/net v0.0.0-20180724234803-3673e40ba225/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20180826012351-8a410e7b638d/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
+golang.org/x/net v0.0.0-20181114220301-adae6a3d119a/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20190213061140-3a22650c66bd/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
 golang.org/x/net v0.0.0-20190311183353-d8887717615a/go.mod h1:t9HGtf8HONx5eT2rtn7q6eTqICYqUVnKs3thJo3Qplg=
 golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3/go.mod h1:t9HGtf8HONx5eT2rtn7q6eTqICYqUVnKs3thJo3Qplg=
+golang.org/x/net v0.0.0-20190613194153-d28f0bde5980/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
 golang.org/x/net v0.0.0-20190620200207-3b0461eec859/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
+golang.org/x/net v0.0.0-20190923162816-aa69164e4478/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
 golang.org/x/net v0.0.0-20200226121028-0de0cce0169b/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
 golang.org/x/net v0.0.0-20201021035429-f5854403a974/go.mod h1:sp8m0HH+o8qH0wwXwYZr8TS3Oi6o0r6Gce1SSxlDquU=
-golang.org/x/net v0.25.0 h1:d/OCCoBEUq33pjydKrGQhw7IlUPI2Oylr+8qLx49kac=
-golang.org/x/net v0.25.0/go.mod h1:JkAGAh7GEvH74S6FOH42FLoXpXbE/aqXSrIQjXgsiwM=
+golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
+golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1/go.mod h1:9tjilg8BloeKEkVJvy7fQ90B1CfIiPueXVOjqfkSzI8=
+golang.org/x/net v0.28.0 h1:a9JDOJc5GMUJ0+UDqmLT86WiEy7iWyIhz8gz8E4e5hE=
+golang.org/x/net v0.28.0/go.mod h1:yqtgsTWOOnlGLG9GFRrK3++bGOUEkNBoHZc8MEDWPNg=
 golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be/go.mod h1:N/0e6XlmueqKjAGxoOufVs8QHGRruUQn6yWY3a++T0U=
 golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20181108010431-42b317875d0f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
+golang.org/x/sync v0.0.0-20181221193216-37e7f081c4d4/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20190423024810-112230192c58/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
 golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
+golang.org/x/sync v0.0.0-20210220032951-036812b2e83c/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
+golang.org/x/sys v0.0.0-20180823144017-11551d06cbcc/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20180830151530-49385e6e1522/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
+golang.org/x/sys v0.0.0-20180905080454-ebe1bf3edb33/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
+golang.org/x/sys v0.0.0-20181116152217-5ac8a444bdc5/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
+golang.org/x/sys v0.0.0-20190222072716-a9d3bda3a223/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
 golang.org/x/sys v0.0.0-20190412213103-97732733099d/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20190422165155-953cdadca894/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20190922100055-0a153f010e69/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20190924154521-2837fb4f24fe/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20191026070338-33540a1f6037/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20200116001909-b77594299b42/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20200122134326-e047566fdf82/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20201119102817-f84b799fce68/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20210303074136-134d130e1a04/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20210330210617-4fbd30eecc44/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
+golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
-golang.org/x/sys v0.20.0 h1:Od9JTbYCk261bKm4M/mw7AklTlFYIa0bIp9BgSm1S8Y=
-golang.org/x/sys v0.20.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
+golang.org/x/sys v0.24.0 h1:Twjiwq9dn6R1fQcyiK+wQyHWfaz/BJB+YIpzU/Cv3Xg=
+golang.org/x/sys v0.24.0/go.mod h1:/VUhepiaJMQUp4+oa/7Zr1D23ma6VTLIYjOOTFZPUcA=
+golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
 golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
+golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
 golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
-golang.org/x/text v0.15.0 h1:h1V/4gjBv8v9cjcR6+AR5+/cIYK5N/WAgiv4xlsEtAk=
-golang.org/x/text v0.15.0/go.mod h1:18ZOQIKpY8NJVqYksKHtTdi31H5itFRjB5/qKTNYzSU=
+golang.org/x/text v0.3.6/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
+golang.org/x/text v0.17.0 h1:XtiM5bkSOt+ewxlOE/aE/AKEHibwj/6gvWMl9Rsh0Qc=
+golang.org/x/text v0.17.0/go.mod h1:BuEKDfySbSR4drPmRPG/7iBdf8hvFMuRexcpahXilzY=
 golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190114222345-bf090417da8b/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
 golang.org/x/tools v0.0.0-20190226205152-f727befe758c/go.mod h1:9Yl7xja0Znq3iFh3HoIrodX9oNMXvdceNzlUR8zjMvY=
 golang.org/x/tools v0.0.0-20190311212946-11955173bddd/go.mod h1:LCzVGOaR6xXOjkQ3onu1FJEFr0SW1gC7cKk1uF8kGRs=
 golang.org/x/tools v0.0.0-20190524140312-2c0ae7006135/go.mod h1:RgjU9mgBXZiqYHBnxXauZ1Gv1EHHAz9KjViQ78xBX0Q=
+golang.org/x/tools v0.0.0-20190907020128-2ca718005c18/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
 golang.org/x/tools v0.0.0-20191108193012-7d206e10da11/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
 golang.org/x/tools v0.0.0-20191119224855-298f0cb1881e/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
 golang.org/x/tools v0.0.0-20200619180055-7c47624df98f/go.mod h1:EkVYQZoAsY45+roYkvgYkIh4xh/qjgUK9TdY2XT94GE=
@@ -215,26 +377,30 @@ google.golang.org/appengine v1.4.0/go.mod h1:xpcJRLb0r/rnEns0DIKYYv+WjYCduHsrkT7
 google.golang.org/genproto v0.0.0-20180817151627-c66870c02cf8/go.mod h1:JiN7NxoALGmiZfu7CAH4rXhgtRTLTxftemlI0sWmxmc=
 google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55/go.mod h1:DMBHOl98Agz4BDEuKkezgsaosCRResVns1a3J2ZsMNc=
 google.golang.org/genproto v0.0.0-20200423170343-7949de9c1215/go.mod h1:55QSHmfGQM9UVYDPBsyGGes0y52j32PQ3BqQfXhyH3c=
-google.golang.org/genproto/googleapis/rpc v0.0.0-20240314234333-6e1732d8331c h1:lfpJ/2rWPa/kJgxyyXM8PrNnfCzcmxJ265mADgwmvLI=
-google.golang.org/genproto/googleapis/rpc v0.0.0-20240314234333-6e1732d8331c/go.mod h1:WtryC6hu0hhx87FDGxWCDptyssuo68sk10vYjF+T9fY=
+google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 h1:e7S5W7MGGLaSu8j3YjdezkZ+m1/Nm0uRVRMEMGk26Xs=
+google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142/go.mod h1:UqMtugtsSgubUsoxbuAoiCXvqvErP7Gf0so0mK9tHxU=
 google.golang.org/grpc v1.19.0/go.mod h1:mqu4LbDTu4XGKhr4mRzUsmM4RtVoemTSY81AxZiDr8c=
 google.golang.org/grpc v1.23.0/go.mod h1:Y5yQAOtifL1yxbo5wqy6BxZv8vAUGQwXBOALyacEbxg=
 google.golang.org/grpc v1.25.1/go.mod h1:c3i+UQWmh7LiEpx4sFZnkU36qjEYZ0imhYfXVyQciAY=
 google.golang.org/grpc v1.27.0/go.mod h1:qbnxyOmOxrQa7FizSgH+ReBfzJrCY1pSN7KXBS8abTk=
 google.golang.org/grpc v1.29.1/go.mod h1:itym6AZVZYACWQqET3MqgPpjcuV5QH3BxFS3IjizoKk=
-google.golang.org/grpc v1.62.1 h1:B4n+nfKzOICUXMgyrNd19h/I9oH0L1pizfk1d4zSgTk=
-google.golang.org/grpc v1.62.1/go.mod h1:IWTG0VlJLCh1SkC58F7np9ka9mx/WNkjl4PGJaiq+QE=
+google.golang.org/grpc v1.67.1 h1:zWnc1Vrcno+lHZCOofnIMvycFcc0QRGIzm9dhnDX68E=
+google.golang.org/grpc v1.67.1/go.mod h1:1gLDyUQU7CTLJI90u3nXZ9ekeghjeM7pTDZlqFNg2AA=
 google.golang.org/protobuf v1.26.0-rc.1/go.mod h1:jlhhOSvTdKEhbULTjvd4ARK9grFBp09yW+WbY/TyQbw=
 google.golang.org/protobuf v1.26.0/go.mod h1:9q0QmTI4eRPtz6boOQmLYwt+qCgq0jsYwAQnmE0givc=
-google.golang.org/protobuf v1.34.1 h1:9ddQBjfCyZPOHPUiPxpYESBLc+T8P3E+Vo4IbKZgFWg=
-google.golang.org/protobuf v1.34.1/go.mod h1:c6P6GXX6sHbq/GpV6MGZEdwhWPcYBgnhAHhKbcUYpos=
+google.golang.org/protobuf v1.35.1 h1:m3LfL6/Ca+fqnjnlqQXNpFPABW1UD7mjh8KO2mKFytA=
+google.golang.org/protobuf v1.35.1/go.mod h1:9fA7Ob0pmnwhb644+1+CVWFRbNajQ6iRojtC/QF5bRE=
+gopkg.in/alecthomas/kingpin.v2 v2.2.6/go.mod h1:FMv+mEhP44yOT+4EoQTLFTRgOQ1FBLkstjWtayDeSgw=
 gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 h1:YR8cESwS4TdDjEe65xsg0ogRM/Nc3DYOhEAlW+xobZo=
 gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
 gopkg.in/ini.v1 v1.67.0 h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=
 gopkg.in/ini.v1 v1.67.0/go.mod h1:pNLf8WUiyNEtQjuu5G5vTm06TEv9tsIgeAvK8hOrP4k=
+gopkg.in/yaml.v2 v2.2.1/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.2/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
+gopkg.in/yaml.v2 v2.2.4/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
+gopkg.in/yaml.v2 v2.2.5/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v2 v2.2.8/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
 gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
 gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
diff --git a/internal/stock/main.go b/internal/stock/main.go
index e1fe459..8990bf4 100644
--- a/internal/stock/main.go
+++ b/internal/stock/main.go
@@ -4,6 +4,7 @@ import (
 	"context"
 
 	"github.com/ghost-yu/go_shop_second/common/config"
+	"github.com/ghost-yu/go_shop_second/common/discovery"
 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
 	"github.com/ghost-yu/go_shop_second/common/server"
 	"github.com/ghost-yu/go_shop_second/stock/ports"
@@ -13,7 +14,6 @@ import (
 	"google.golang.org/grpc"
 )
 
-// stock 服务和 order 服务一样，在进程入口先加载配置，避免后续各层自己读文件。
 func init() {
 	if err := config.NewViperConfig(); err != nil {
 		logrus.Fatal(err)
@@ -21,21 +21,26 @@ func init() {
 }
 
 func main() {
-	// serviceName 决定去哪个配置分组下找 grpc-addr 等参数。
 	serviceName := viper.GetString("stock.service-name")
-	// serverType 让同一个二进制可以选择跑哪种协议，当前 lesson 主要走 grpc 分支。
 	serverType := viper.GetString("stock.server-to-run")
 
 	logrus.Info(serverType)
 
-	// 这里先组装应用层，再根据配置决定挂到哪个端口适配器上。
 	ctx, cancel := context.WithCancel(context.Background())
 	defer cancel()
 
 	application := service.NewApplication(ctx)
+
+	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+	defer func() {
+		_ = deregisterFunc()
+	}()
+
 	switch serverType {
 	case "grpc":
-		// gRPC 分支把 ports.GRPCServer 注册为 StockService 的服务端实现。
 		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
 			svc := ports.NewGRPCServer(application)
 			stockpb.RegisterStockServiceServer(server, svc)
diff --git a/internal/stock/ports/grpc.go b/internal/stock/ports/grpc.go
index 72bc356..91e34e8 100644
--- a/internal/stock/ports/grpc.go
+++ b/internal/stock/ports/grpc.go
@@ -5,10 +5,9 @@ import (
 
 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
 	"github.com/ghost-yu/go_shop_second/stock/app"
+	"github.com/ghost-yu/go_shop_second/stock/app/query"
 )
 
-// GRPCServer 负责实现 stockpb.StockServiceServer 接口。
-// 外部请求进入后，应该在这里完成参数转换，再转给应用层处理。
 type GRPCServer struct {
 	app app.Application
 }
@@ -18,13 +17,20 @@ func NewGRPCServer(app app.Application) *GRPCServer {
 }
 
 func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
-	// 这里后续会把商品 ID 列表转给库存查询用例。
-	//TODO implement me
-	panic("implement me")
+	items, err := G.app.Queries.GetItems.Handle(ctx, query.GetItems{ItemIDs: request.ItemIDs})
+	if err != nil {
+		return nil, err
+	}
+	return &stockpb.GetItemsResponse{Items: items}, nil
 }
 
 func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
-	// 这一步对应 order 下单前的库存校验，是跨服务调用的关键入口。
-	//TODO implement me
-	panic("implement me")
+	items, err := G.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{Items: request.Items})
+	if err != nil {
+		return nil, err
+	}
+	return &stockpb.CheckIfItemsInStockResponse{
+		InStock: 1,
+		Items:   items,
+	}, nil
 }
diff --git a/internal/stock/service/application.go b/internal/stock/service/application.go
index 23daf9a..0ccb6a8 100644
--- a/internal/stock/service/application.go
+++ b/internal/stock/service/application.go
@@ -3,11 +3,22 @@ package service
 import (
 	"context"
 
+	"github.com/ghost-yu/go_shop_second/common/metrics"
+	"github.com/ghost-yu/go_shop_second/stock/adapters"
 	"github.com/ghost-yu/go_shop_second/stock/app"
+	"github.com/ghost-yu/go_shop_second/stock/app/query"
+	"github.com/sirupsen/logrus"
 )
 
-// NewApplication 当前先返回一个空的应用层门面。
-// lesson 后续扩展库存能力时，repo 和 handler 的装配点也会放在这里。
-func NewApplication(ctx context.Context) app.Application {
-	return app.Application{}
+func NewApplication(_ context.Context) app.Application {
+	stockRepo := adapters.NewMemoryStockRepository()
+	logger := logrus.NewEntry(logrus.StandardLogger())
+	metricsClient := metrics.TodoMetrics{}
+	return app.Application{
+		Commands: app.Commands{},
+		Queries: app.Queries{
+			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, logger, metricsClient),
+			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricsClient),
+		},
+	}
 }
~~~

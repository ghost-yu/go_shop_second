# gorder-v2 Lesson5-Lesson11 复现级代码解读（Go 小白版）

这版文档目标是“可复现”，不是只做概念导读。
你可以按本文件逐个文件阅读并运行，每个文件都给出：
1. 文件职责
2. 关键代码必须理解点
3. 上下游依赖（谁调用它，它调用谁）
4. 复现检查动作（你可以自己执行验证）

## 0. 先说明边界

Lesson5-Lesson11 里涉及的文件很多，其中有三类：
1. 业务代码文件：需要详细读懂。
2. 生成文件（*.pb.go, *_grpc.pb.go）：需要理解结构和用途，不建议逐行硬背。
3. 依赖锁文件（go.sum）：不需要逐条解释，每行本质是模块校验和。

因此本文对所有文件都做说明，但会按“学习价值”给不同颗粒度，确保你能复现整个阶段。

## 1. Lesson5 到 Lesson11 的总调用图

1. 配置加载：internal/common/config/global.yaml -> viper
2. 进程入口：internal/order/main.go、internal/stock/main.go
3. 基础服务启动：internal/common/server/http.go、internal/common/server/gprc.go
4. 端口层：internal/order/http.go、internal/order/ports/grpc.go、internal/stock/ports/grpc.go
5. 应用层：internal/order/app/...、internal/stock/app/...
6. 领域接口：internal/order/domain/order/repository.go、internal/stock/domain/stock/repository.go
7. 适配器实现：in-memory repo + grpc 适配器
8. 协议层：api/*.proto + internal/common/genproto/*.pb.go

---

## 2. 按文件详细解释（Lesson5-Lesson11）

### A. 协议与配置文件

#### api/stockpb/stock.proto
文件职责：定义库存服务 RPC 契约。
关键理解：
1. service StockService 里定义了两个 RPC：GetItems、CheckIfItemsInStock。
2. 请求和响应消息是跨服务通信唯一标准，服务端和客户端必须一致。
3. 引用 orderpb/order.proto 是为了复用订单侧 Item 结构，避免重复建模。
依赖关系：
1. 被 internal/common/genproto/stockpb/stock.pb.go 和 stock_grpc.pb.go 生成工具使用。
2. 被 order 侧 grpc 客户端调用路径间接使用。
复现检查：
1. 运行 proto 生成脚本后，确认生成代码存在对应消息与服务接口。

#### api/orderpb/order.proto
文件职责：定义订单服务对外协议和公共消息（如 Item）。
关键理解：
1. Item 和 ItemWithQuantity 是库存校验与下单共用的数据结构。
2. OrderService 方法签名决定了 order 端口层函数形态。
依赖关系：
1. 被 OpenAPI/GRPC 生成代码使用。
2. 被 stock.proto import。
复现检查：
1. 改字段名后重新生成，会导致调用方编译错误，证明契约约束生效。

#### internal/common/config/global.yaml
文件职责：集中配置各服务地址与运行模式。
关键理解：
1. order 与 stock 的 grpc/http 地址是独立的。
2. fallback-grpc-addr 是防御性配置，避免空地址直接崩溃。
依赖关系：
1. 被 config/viper 初始化后读取。
2. 被 server/http.go 和 server/gprc.go 间接消费。
复现检查：
1. 改 grpc-addr 后重启服务，端口变化应立即生效。

#### .gitignore
文件职责：控制 Git 不追踪哪些本地文件。
关键理解：
1. 用于忽略临时构建产物和本地工具缓存。
2. 防止开发机噪音提交进仓库。
依赖关系：
1. Git 客户端直接读取。
复现检查：
1. 新建被忽略文件后 git status 不应出现。

### B. 生成代码文件

#### internal/common/genproto/orderpb/order.pb.go
文件职责：order.proto 的消息结构与序列化代码。
关键理解：
1. 这类文件由 protoc 自动生成，不手改。
2. 包含 message struct、反射元数据、编解码实现。
依赖关系：
1. 被业务代码 import，作为请求/响应 DTO。
复现检查：
1. 删除后项目会出现大量类型缺失错误。

#### internal/common/genproto/stockpb/stock.pb.go
文件职责：stock.proto 消息生成代码。
关键理解：
1. 与 order.pb.go 同类，承载跨服务消息类型。
依赖关系：
1. 被 stock 服务端与 order 客户端共同依赖。
复现检查：
1. 运行 stock grpc 相关代码需依赖此文件类型。

#### internal/common/genproto/stockpb/stock_grpc.pb.go
文件职责：stock.proto 的 gRPC 客户端/服务端接口生成代码。
关键理解：
1. 生成 StockServiceClient 和 StockServiceServer 接口。
2. order 侧 grpc 适配器通过 client 发请求。
3. stock 侧端口实现通过 server 接口对外暴露。
依赖关系：
1. 被 internal/order/adapters/grpc/stock_grpc.go 使用。
2. 被 internal/stock/ports/grpc.go 实现。
复现检查：
1. 在 main.go 注册 stockpb.RegisterStockServiceServer 即依赖此文件。

### C. 基础设施启动文件

#### internal/common/server/gprc.go
文件职责：统一 gRPC server 启动流程，所有 gRPC 服务都通过这里启动，业务服务只需提供「注册业务」这一个回调。

完整代码+逐段注释：
```go
package server

import (
    "net"
    grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
    grpc_tags   "github.com/grpc-ecosystem/go-grpc-middleware/tags"
    "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
    "google.golang.org/grpc"
)

// init 是 Go 的特殊函数，包被 import 时自动执行，不需要手动调用。
// 这里把 gRPC 框架内部日志替换成项目统一使用的 logrus，同时把级别设为 Warn，
// 减少 gRPC 框架内部的噪音日志，保持日志格式一致，方便集中收集和排查问题。
func init() {
    logger := logrus.New()
    logger.SetLevel(logrus.WarnLevel)                      // 只打 warn 及以上
    grpc_logrus.ReplaceGrpcLogger(logrus.NewEntry(logger)) // 替换 gRPC 内置 logger
}

// RunGRPCServer 是对外的唯一入口。
// 参数 serviceName：对应 global.yaml 里的 key，比如 "order" 或 "stock"。
// 参数 registerServer：业务方传入的回调，用于把自己的 handler 注册进 grpc.Server。
//
// 调用方写法（order/main.go）示例：
//   server.RunGRPCServer("order", func(s *grpc.Server) {
//       orderpb.RegisterOrderServiceServer(s, ports.NewGRPCServer(app))
//   })
func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server)) {
    // viper.Sub("order").GetString("grpc-addr") 相当于读 global.yaml 里 order.grpc-addr 的值。
    addr := viper.Sub(serviceName).GetString("grpc-addr")
    if addr == "" {
        // 地址缺失时用全局兜底地址，防止配置漏填导致进程直接崩溃。
        addr = viper.GetString("fallback-grpc-addr")
    }
    RunGRPCServerOnAddr(addr, registerServer)
}

// RunGRPCServerOnAddr 真正创建并运行 gRPC server。
// 独立函数的原因：测试时可以直接传 "127.0.0.1:0"（随机端口），不依赖配置文件。
func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
    logrusEntry := logrus.NewEntry(logrus.StandardLogger())

    // grpc.NewServer(...) 创建 gRPC 服务器，括号里配置拦截器。
    // 拦截器 = HTTP 框架里的中间件，每次 RPC 调用进来和出去时都会经过它们，
    // 适合放日志、追踪、鉴权、panic 恢复等和业务无关的横切逻辑。
    grpcServer := grpc.NewServer(

        // ChainUnaryInterceptor：配置「一元 RPC」拦截器链（普通请求-响应模式，最常用）。
        // 执行顺序：请求进来时从左到右，响应返回时从右到左。
        grpc.ChainUnaryInterceptor(
            // ① 从 proto 生成的请求结构里自动提取字段（如 customer_id）打到日志 Tag 里。
            //   后续该请求的所有日志都自动携带这些字段，排查问题时不用手动传。
            grpc_tags.UnaryServerInterceptor(
                grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor),
            ),
            // ② 自动记录每次 RPC 的方法名、耗时、gRPC 状态码到 logrus。
            grpc_logrus.UnaryServerInterceptor(logrusEntry),

            // 下面是后续 lesson 会逐步打开的能力，当前注释保留位置：
            // otelgrpc.UnaryServerInterceptor(),       // 分布式链路追踪（Lesson25+）
            // srvMetrics.UnaryServerInterceptor(...),  // Prometheus 指标采集（Lesson43+）
            // selector.UnaryServerInterceptor(...),    // 选择性鉴权（某些路由才需要鉴权）
            // recovery.UnaryServerInterceptor(...),    // panic 自动恢复（防止一个请求崩掉整个进程）
        ),

        // ChainStreamInterceptor：配置「流式 RPC」拦截器链（双向流/服务端推送场景）。
        // 当前项目主要用一元 RPC，这里配置和一元保持同步即可。
        grpc.ChainStreamInterceptor(
            grpc_tags.StreamServerInterceptor(
                grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor),
            ),
            grpc_logrus.StreamServerInterceptor(logrusEntry),
            // 同上，预留位置暂时注释
        ),
    )

    // 执行业务注册回调，把 OrderService/StockService 的具体实现挂到 grpcServer 上。
    // 这行之前 grpcServer 是空的（不知道有哪些方法），执行后才知道怎么处理请求。
    registerServer(grpcServer)

    // net.Listen 在 TCP 层监听指定地址，此时还没开始接受 gRPC 请求。
    listen, err := net.Listen("tcp", addr)
    if err != nil {
        // 端口被占用、权限不足等情况，直接 panic 让进程退出，而不是静默失败。
        logrus.Panic(err)
    }
    logrus.Infof("Starting gRPC server, Listening: %s", addr)

    // grpcServer.Serve 是阻塞调用，进入后一直处理 gRPC 请求，直到 server Stop 或出错。
    if err := grpcServer.Serve(listen); err != nil {
        logrus.Panic(err)
    }
}
```

核心设计要点：
1. **两函数分离**：RunGRPCServer 负责读配置（依赖 viper），RunGRPCServerOnAddr 只依赖
   地址字符串。单元测试时可以绕过配置直接调用 RunGRPCServerOnAddr，不需要 global.yaml。
2. **registerServer 回调模式**：框架不知道业务注册了什么，业务不知道框架怎么启动，
   通过这一个函数参数完成解耦。这是 Go 里高频出现的依赖注入写法，全项目会反复看到。
3. **拦截器链是功能扩展点**：目前只有日志+标签两层，后续每个注释行都是可独立打开的能力，
   打开时完全不需要改任何业务代码，只需取消注释并引入对应包。

依赖关系：
- 上游调用方：order/main.go 和 stock/main.go 传入 serviceName 和 registerServer 调用本函数。
- 下游依赖：registerServer 回调由业务服务提供（例如 orderpb.RegisterOrderServiceServer）。

复现检查：
1. 启动后终端出现 `Starting gRPC server, Listening: 127.0.0.1:5002`。
2. 改 global.yaml 的 grpc-addr 后重启，监听端口随之变化。
3. 用重复端口启动第二个实例，net.Listen 会报 bind: address already in use 并 panic。

---

#### internal/common/server/http.go
文件职责：统一 HTTP server 启动流程，所有 HTTP 服务都通过这里启动，业务方只需传入路由注册回调。

完整代码+逐段注释：
```go
package server

import (
    "github.com/gin-gonic/gin"
    "github.com/spf13/viper"
)

// RunHTTPServer 按服务名从 viper 读 http-addr，然后启动 HTTP 服务。
// 参数 wrapper：业务方提供的路由注册函数，接收 *gin.Engine 并往里注册路由和 handler。
//
// 调用方写法（order/main.go）示例：
//   server.RunHTTPServer("order", func(router *gin.Engine) {
//       ports.RegisterHandlersWithOptions(router, HTTPServer{app: application}, opts)
//   })
func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
    // viper.Sub("order").GetString("http-addr") 读取 global.yaml 里 order.http-addr 的值。
    addr := viper.Sub(serviceName).GetString("http-addr")
    if addr == "" {
        // TODO：后续应加 warning 日志提醒配置缺失。
        // 注意：addr 为空时 gin.Run("") 会监听 :80，非 root 权限下会权限错误。
    }
    RunHTTPServerOnAddr(addr, wrapper)
}

// RunHTTPServerOnAddr 真正创建 gin 引擎并运行。
// 独立出来的原因：测试时可以直接传 ":0"（随机端口），不依赖 viper 配置。
func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
    // gin.New() 创建一个没有任何默认中间件的干净 gin 引擎。
    // 区别于 gin.Default()：Default() 会自动加 Logger 和 Recovery 中间件，
    // 这里用 New() 是为了完全掌控中间件，后续 lesson 会按需手动添加。
    apiRouter := gin.New()

    // 执行业务方传入的路由注册函数，把路由和 handler 挂到 apiRouter 上。
    // 执行完后 apiRouter 才知道"POST /api/customer/.../orders"对应哪个 handler。
    wrapper(apiRouter)

    // ⚠️ 注意：下面这行实际上是无效代码（遗留 bug）。
    // apiRouter.Group("/api") 返回一个新的 RouterGroup 对象，但没有赋值给变量，
    // 返回值被直接丢弃，对已有路由没有任何影响，相当于执行了一个空操作。
    // 真正的 /api 前缀是在 wrapper 内部通过 RegisterHandlersWithOptions 的 BaseURL:"/api" 设置的。
    // 理解这一点可以防止你自己写出同类 bug。
    apiRouter.Group("/api")

    // apiRouter.Run 是阻塞调用，开始监听并处理 HTTP 请求。
    // 出错（端口冲突、权限不足）时直接 panic，不要静默失败。
    if err := apiRouter.Run(addr); err != nil {
        panic(err)
    }
}
```

核心设计要点：
1. **gin.New() 而非 gin.Default()**：主动掌控中间件，不被框架默认行为绑架。
   需要加 cors、recovery、自定义日志时直接往 apiRouter 上加，与 grpc.go 的拦截器链思路一致。
2. **wrapper 回调模式**：与 gprc.go 的 registerServer 完全相同的解耦思路，
   http.go 不知道有哪些路由，order/main.go 负责告诉它。
3. **apiRouter.Group("/api") 是无效代码**：这是代码里的遗留问题，不影响功能运行，
   因为路由前缀已经由 RegisterHandlersWithOptions 的 BaseURL 参数正确设置了。

依赖关系：
- 上游调用方：order/main.go 传入 "order" 和路由注册 wrapper 调用本函数。
- 下游依赖：wrapper 函数由 order/main.go 提供，其中调用 ports.RegisterHandlersWithOptions。

复现检查：
1. 启动 order 服务后，`curl http://127.0.0.1:8282/api/customer/test/orders` 能到达 handler。
2. 修改 http-addr 端口后重启，访问旧端口应收到连接拒绝。
3. 将端口改为已占用端口，gin.Run 返回 error 并触发 panic。

### D. 入口与端口层文件

#### internal/order/main.go
文件职责：Order 服务入口。
关键理解：
1. init 先加载配置，main 再启动服务。
2. 同时启动 grpc 和 http，说明订单服务双协议暴露。
3. HTTP 注册走 OpenAPI 生成的 RegisterHandlersWithOptions。
依赖关系：
1. 调用 common/server。
2. 调用 order/ports + order/http。
3. 调用 order/service.NewApplication（Lesson7 后）。
复现检查：
1. 运行后应同时看到 HTTP 与 gRPC 监听日志。

#### internal/order/http.go
文件职责：Order HTTP 端口适配层。
关键理解：
1. 方法名来自 OpenAPI 生成接口，不是手写命名。
2. 正式业务版本里只做三件事：绑定参数、调用应用层、返回响应。
依赖关系：
1. 被 main.go 注册为 HTTP handler。
2. 调用 app.Commands / app.Queries。
复现检查：
1. POST 下单和 GET 查询调用时会进入对应 handler。

#### internal/order/ports/grpc.go
文件职责：Order gRPC 端口实现。
关键理解：
1. 实现 orderpb.OrderServiceServer 接口。
2. 该阶段多为 TODO，后续 lesson 逐步填充。
依赖关系：
1. 被 main.go 的 RegisterOrderServiceServer 注册。
复现检查：
1. 若方法未实现，调用对应 RPC 会 panic。

#### internal/stock/main.go
文件职责：Stock 服务入口。
关键理解：
1. 读取 stock.server-to-run 决定跑 grpc 还是 http。
2. 当前阶段主要跑 grpc。
依赖关系：
1. 调用 service.NewApplication。
2. 调用 common/server.RunGRPCServer。
复现检查：
1. 改 server-to-run 值可触发不同分支行为。

#### internal/stock/ports/grpc.go
文件职责：Stock gRPC 端口实现。
关键理解：
1. 需要实现 GetItems 和 CheckIfItemsInStock。
2. 这层把外部请求转成应用层调用。
依赖关系：
1. 被 stock main 注册。
2. 供 order 的 grpc client 调用。
复现检查：
1. 用 grpc 客户端请求库存接口能否返回。

### E. 应用层组装与门面文件

#### internal/stock/app/app.go
文件职责：Stock 应用层门面（Commands、Queries 聚合）。
关键理解：
1. 通过一个 Application struct 对外暴露能力。
2. 减少上层对内部 handler 细节感知。
依赖关系：
1. 被 stock service/application.go 填充。
2. 被 stock ports 层消费。
复现检查：
1. main 中传递 Application 后端口层能拿到功能入口。

#### internal/stock/service/application.go
文件职责：Stock 的依赖注入与应用组装。
关键理解：
1. 这是组合根，负责把 repo、handler 组装成 app.Application。
2. main 不直接 new 各种依赖，职责更清晰。
依赖关系：
1. 被 stock main 调用。
2. 依赖 stock adapters/domain/app。
复现检查：
1. 替换 repo 实现时只改此处注入。

#### internal/order/app/app.go
文件职责：Order 应用层门面。
关键理解：
1. Commands 放写流程（Create/Update）。
2. Queries 放读流程（GetCustomerOrder）。
依赖关系：
1. 被 order/service/application.go 赋值。
2. 被 HTTP/gRPC 端口调用。
复现检查：
1. 端口层只依赖 app.Application 就能调用所有用例。

#### internal/order/service/application.go
文件职责：Order 依赖注入核心。
关键理解：
1. 初始化 orderRepo、logger、metricsClient。
2. 初始化 stock grpc client 并封装成 StockService。
3. 把 command/query handler 统一装配到 app.Application。
依赖关系：
1. 被 order main 调用。
2. 依赖 adapters、decorator、query/service 接口。
复现检查：
1. 替换 stock 调用实现时，这里改注入即可。

### F. 领域层与仓储实现

#### internal/order/domain/order/order.go
文件职责：订单领域模型。
关键理解：
1. Order 代表核心业务实体。
2. Items 使用 proto 结构，便于和接口层衔接。
依赖关系：
1. 被 command/query/repository 使用。
复现检查：
1. 新增字段后需要同步影响创建与查询返回。

#### internal/order/domain/order/repository.go
文件职责：订单仓储接口定义。
关键理解：
1. 定义 Create/Get/Update 三个能力。
2. NotFoundError 提供业务语义化错误。
依赖关系：
1. 被应用层依赖。
2. 被 adapters 实现。
复现检查：
1. 任何 repo 实现都必须满足该接口。

#### internal/order/adapters/order_inmem_repository.go
文件职责：订单仓储 in-memory 实现。
关键理解：
1. RWMutex 保障并发安全。
2. Create 用时间戳生成 ID 并写入内存切片。
3. Get 按订单号和客户号匹配。
4. Update 使用高阶函数 updateFn 控制更新逻辑。
依赖关系：
1. 实现 domain.Repository。
2. 被 order/service/application.go 注入到 handler。
复现检查：
1. 连续创建订单，内存 store 长度增加。
2. 查不存在订单返回 NotFoundError。

#### internal/stock/domain/stock/repository.go
文件职责：库存仓储接口定义。
关键理解：
1. 定义 GetItems 与 CheckIfItemsInStock 等库存查询能力。
2. 约束 stock 应用层只依赖接口。
依赖关系：
1. 被 stock 适配器实现。
复现检查：
1. repo 实现变更不影响应用层签名。

#### internal/stock/adapters/stock_inmem_repository.go
文件职责：库存仓储 in-memory 实现。
关键理解：
1. 用内存数据模拟库存系统。
2. 返回商品详情和库存校验结果。
依赖关系：
1. 实现 stock domain repo。
2. 被 stock service 注入。
复现检查：
1. 调整初始库存数据，库存接口结果应变化。

### G. Command / Query 与装饰器

#### internal/common/decorator/query.go
文件职责：QueryHandler 泛型接口与装饰器装配入口。
关键理解：
1. 统一查询处理器形态。
2. ApplyQueryDecorators 按顺序套日志和指标。
依赖关系：
1. 被 query handler 构造函数调用。
复现检查：
1. 执行查询时可看到日志与指标上报被触发。

#### internal/common/decorator/command.go
文件职责：CommandHandler 泛型接口与装饰器装配入口。
关键理解：
1. 与 query 对称，作用于写操作。
2. 保证 command 也拥有统一日志/指标能力。
依赖关系：
1. 被 command handler 构造函数调用。
复现检查：
1. 下单成功/失败时都有对应统计路径。

#### internal/common/decorator/logging.go
文件职责：通用日志装饰器。
关键理解：
1. 记录 action 名、请求体、执行结果。
2. 错误场景会打印失败日志，方便排查。
依赖关系：
1. 被 query.go / command.go 装配调用。
复现检查：
1. 故意触发参数错误，日志中应有失败记录。

#### internal/common/decorator/metrics.go
文件职责：通用指标装饰器。
关键理解：
1. 统计耗时、成功次数、失败次数。
2. action 名来自请求类型。
依赖关系：
1. 依赖 MetricsClient 接口。
2. 被 query/command 装配。
复现检查：
1. 打断点可看到 Inc 被调用三类 key。

#### internal/common/metrics/todo_metrics.go
文件职责：指标客户端占位实现。
关键理解：
1. 当前 Inc 空实现，仅保证接口连通。
2. 后续可替换为 Prometheus/StatsD 真正上报。
依赖关系：
1. 被 application.go 注入到装饰器。
复现检查：
1. 替换 Inc 内容后可观察指标输出变化。

#### internal/order/app/query/get_customer_order.go
文件职责：查询订单用例。
关键理解：
1. Handle 只做一件事：调用 repo.Get。
2. 构造函数里应用了 query decorators。
依赖关系：
1. 被 HTTP Get handler 调用。
2. 依赖 order repository 接口。
复现检查：
1. 用存在/不存在订单号分别调用，观察成功与错误路径。

#### internal/order/app/query/service.go
文件职责：定义 Order 侧需要的库存服务端口接口。
关键理解：
1. 这是“防腐层端口”，隔离外部 grpc 客户端细节。
2. create_order 依赖该接口而非具体实现。
依赖关系：
1. 被 create_order.go 依赖。
2. 被 adapters/grpc/stock_grpc.go 实现。
复现检查：
1. 用 mock 实现该接口可做单元测试。

#### internal/order/app/command/create_order.go
文件职责：创建订单用例。
关键理解：
1. validate 先检查 items 非空。
2. packItems 合并重复商品数量，减少重复校验。
3. 调用 StockService.CheckIfItemsInStock 进行库存校验。
4. 校验通过后调用 orderRepo.Create。
依赖关系：
1. 被 HTTP POST handler 调用。
2. 依赖 order repo + stock service 接口。
复现检查：
1. 传空 items 应报错。
2. 传重复 item id 时会被合并后再请求库存服务。

#### internal/order/app/command/update_order.go
文件职责：更新订单用例。
关键理解：
1. UpdateFn 由上层传入，决定如何改订单。
2. handler 只负责流程控制和仓储调用。
依赖关系：
1. 被上层更新场景调用。
2. 依赖 order repo.Update。
复现检查：
1. UpdateFn 为空时走默认 no-op 逻辑。

#### internal/order/adapters/grpc/stock_grpc.go
文件职责：StockService 的 gRPC 适配器实现。
关键理解：
1. 封装 stockpb.StockServiceClient。
2. 把领域需要的方法映射为具体 RPC 请求。
依赖关系：
1. 实现 order/app/query/service.go 的 StockService。
2. 被 order/service/application.go 注入。
复现检查：
1. 断开 stock 服务后，下单库存校验应报远程调用错误。

### H. 模块与工具配置文件

#### internal/common/go.mod
文件职责：common 模块依赖声明。
关键理解：
1. require 指明编译期直接依赖。
2. indirect 是传递依赖。
依赖关系：
1. 与 common/go.sum 配套。
复现检查：
1. go mod tidy 后依赖变化会更新。

#### internal/order/go.mod
文件职责：order 模块依赖声明。
关键理解：
1. 包含 gin、grpc、viper 等运行依赖。
2. 通过 replace（若存在）可指向本地模块。
依赖关系：
1. 与 order/go.sum 配套。
复现检查：
1. 新增 import 后 go mod tidy 会增补依赖。

#### internal/stock/go.mod
文件职责：stock 模块依赖声明。
关键理解：
1. 与 order 模块类似，但按 stock 需求维护。
依赖关系：
1. 与 stock/go.sum 配套。
复现检查：
1. 删除必要依赖后编译会失败。

#### internal/common/go.sum
文件职责：common 模块依赖校验和锁文件。
关键理解：
1. 每行是某个模块版本的 hash，不是业务代码。
2. 用于防篡改和可重复构建。
依赖关系：
1. 被 Go module 下载与校验流程读取。
复现检查：
1. 删除后重新 tidy 会再生成。

#### internal/order/go.sum
文件职责：order 模块依赖校验和。
关键理解：
1. 本质与 common/go.sum 相同。
依赖关系：
1. 由 go 命令自动维护。
复现检查：
1. 拉新依赖后会新增对应 checksum 行。

#### internal/stock/go.sum
文件职责：stock 模块依赖校验和。
关键理解：
1. 本质与其他 go.sum 相同。
依赖关系：
1. 由 go 命令自动维护。
复现检查：
1. tidy 后内容变化属于正常现象。

#### internal/stock/.air.toml
文件职责：stock 服务本地热重载配置。
关键理解：
1. 指定监听目录、忽略目录、重建命令。
2. 开发时改代码自动重启。
依赖关系：
1. 被 air 工具读取。
复现检查：
1. 运行 air 后修改 go 文件观察自动重编译。

#### internal/order/.air.toml
文件职责：order 服务热重载配置。
关键理解：
1. 与 stock/.air.toml 同类。
依赖关系：
1. 被 air 工具读取。
复现检查：
1. 修改 order 代码后自动重启。

#### internal/kitchen/.air.toml
文件职责：kitchen 服务热重载配置。
关键理解：
1. 给 kitchen 模块开发阶段提速。
依赖关系：
1. 被 air 工具读取。
复现检查：
1. 进入 kitchen 模块运行 air 验证。

#### internal/payment/.air.toml
文件职责：payment 服务热重载配置。
关键理解：
1. 给 payment 模块开发阶段提速。
依赖关系：
1. 被 air 工具读取。
复现检查：
1. 进入 payment 模块运行 air 验证。

---

## 3. 复现步骤（按 Lesson5-11）

1. 启动 stock 服务（grpc）。
2. 启动 order 服务（http + grpc）。
3. 调用 order 下单 HTTP 接口。
4. order 的 create_order handler 内会通过 StockService 走 grpc 校验库存。
5. 校验通过后写入 order in-memory repo。
6. 再调用查询接口，验证能读到订单。

如果你复现失败，优先按这个顺序查：
1. 配置地址是否一致（global.yaml）。
2. stock 服务是否先启动。
3. order/service/application.go 是否成功注入 stock grpc client。
4. proto 生成代码是否和 proto 契约一致。

---

## 4. Lesson5-11 的“必须会改”文件清单

你练习时，优先动这些文件：
1. internal/order/http.go
2. internal/order/app/command/create_order.go
3. internal/order/app/query/get_customer_order.go
4. internal/order/adapters/order_inmem_repository.go
5. internal/order/service/application.go
6. internal/order/adapters/grpc/stock_grpc.go
7. internal/stock/ports/grpc.go

这些文件改熟了，后面的数据库、消息队列、可观测性才容易接上。

---

## 5. 下一批输出方式（你确认后继续）

下一批我按同样模板写 Lesson12-Lesson20，并继续提交到 Git：
1. 每个文件都有“职责 + 关键代码 + 调用关系 + 复现检查”。
2. 对生成文件/go.sum 继续采用“用途级详细说明”，不做低价值逐行复述。

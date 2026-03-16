# gorder-v2 Lesson5-Lesson11 关键代码解读（Go 小白版）

这份文档不追求每行都解释，而是只解释开发必须理解的核心代码。
阅读顺序建议：先看启动流程，再看接口定义，再看业务处理，再看依赖注入。

## 一、从外到内的整体调用链

### 1) Order 服务启动链
1. 进程入口：internal/order/main.go
2. 读取配置：internal/common/config/global.yaml
3. 启动 gRPC 公共框架：internal/common/server/gprc.go
4. 启动 HTTP 公共框架：internal/common/server/http.go
5. 注册 HTTP 路由：internal/order/http.go（后续 lesson 才逐步填充业务）
6. 注册 gRPC 端口：internal/order/ports/grpc.go（这阶段还是占位）

### 2) 创建订单业务链（Lesson10-11 开始成型）
1. HTTP Handler 接收请求
2. 调用 app.Commands.CreateOrder
3. 先校验商品和库存（经 StockService）
4. 再写入 OrderRepository
5. 返回 order_id

这条链路体现了分层思想：
- 端口层（HTTP/gRPC）只做入参出参处理
- 应用层（command/query）只编排业务流程
- 领域层（domain）定义模型和接口
- 适配器层（adapters）实现接口（如 in-memory、gRPC）

---

## 二、Lesson5：服务骨架搭建（http, grpc 服务搭建）

### 关键文件与作用

#### api/stockpb/stock.proto
作用：定义库存服务 gRPC 协议。
必须理解点：
1. service StockService：定义了可远程调用的接口。
2. GetItems：通过商品 ID 列表取商品详情。
3. CheckIfItemsInStock：校验商品数量是否有库存。
4. import orderpb/order.proto：复用订单协议里的 Item 结构，避免重复定义。

为什么重要：
- 这是跨服务通信的“合同”。
- 只要合同不破坏兼容，服务可以独立演进。

#### internal/common/server/gprc.go
作用：公共 gRPC 启动器。
必须理解点：
1. RunGRPCServer(serviceName, registerServer)
   - 按服务名读取配置中的 grpc-addr。
   - 若缺失地址，使用 fallback-grpc-addr。
2. RunGRPCServerOnAddr(addr, registerServer)
   - 创建 grpc.Server。
   - 挂载 Unary/Stream 拦截器（日志、标签等）。
   - 调用 registerServer 把业务服务注册进去。
   - 监听并开始服务。
3. 拦截器链相当于中间件：以后可以加 tracing、metrics、鉴权、panic 恢复。

为什么重要：
- 业务服务无需关心底层启动细节。
- 统一中间件入口，便于全链路治理。

#### internal/common/server/http.go
作用：公共 HTTP 启动器。
必须理解点：
1. RunHTTPServer(serviceName, wrapper)
   - 按服务名读取 http-addr。
2. RunHTTPServerOnAddr(addr, wrapper)
   - 创建 gin router。
   - 执行 wrapper 注册路由。
   - 启动监听。

为什么重要：
- 路由注册和启动分离。
- 后续每个服务都可以复用这套启动逻辑。

#### internal/common/config/global.yaml
作用：中心化管理服务配置。
必须理解点：
1. order 与 stock 分别有独立地址。
2. fallback-grpc-addr 提供兜底值，防止空配置直接崩。

为什么重要：
- 服务地址、运行模式不硬编码在代码中。
- 后续可迁移到不同环境（本地/测试/生产）。

#### internal/order/main.go
作用：Order 服务主入口。
必须理解点：
1. init() 中加载配置。
2. main() 中同时启动 gRPC 和 HTTP。
3. gRPC 通过 RegisterOrderServiceServer 注册服务实现。
4. HTTP 通过 RegisterHandlersWithOptions 注册 OpenAPI 生成的路由。

为什么重要：
- 它连接了配置层、公共 server 层、端口层。
- 是理解服务生命周期的第一站。

#### internal/order/http.go
作用：HTTP 端口层（当时为 TODO 骨架）。
必须理解点：
1. 方法签名来自 OpenAPI 生成的接口约定。
2. 这层不应写复杂业务，应该只做参数转换和调用应用层。

为什么重要：
- 你后续会在这里看到“接口适配”而不是“业务实现”。

#### internal/order/ports/grpc.go、internal/stock/ports/grpc.go
作用：gRPC 端口层骨架。
必须理解点：
1. 先把服务入口位置搭好。
2. 真正业务在后续 lesson 再逐步实现。

---

## 三、Lesson6-7：引入 Application 组装层

### 关键新增
- internal/stock/app/app.go
- internal/stock/service/application.go
- internal/order/app/app.go
- internal/order/service/application.go

### 必须理解的设计目的
1. app/app.go
   - 定义 Application 聚合对象（Commands、Queries）。
   - 让上层调用时只依赖 Application，不感知内部实现细节。
2. service/application.go
   - 集中做依赖注入：repo、logger、外部 client、handler。
   - main.go 只负责启动，不负责业务对象创建。

为什么重要：
- 这是典型“组合根（Composition Root）”。
- 依赖关系集中管理后，测试和替换实现会简单很多。

---

## 四、Lesson8：引入领域接口 + 内存仓储

### 关键新增
- internal/order/domain/order/order.go
- internal/order/domain/order/repository.go
- internal/order/adapters/order_inmem_repository.go
- internal/stock/domain/stock/repository.go
- internal/stock/adapters/stock_inmem_repository.go

### 必须理解点

#### domain/repository.go（接口定义）
1. Repository 是“能力契约”，定义 Create/Get/Update。
2. 应用层只依赖接口，不依赖 in-memory 这种具体实现。

#### adapters/order_inmem_repository.go（接口实现）
1. 用 RWMutex 保护内存数据，避免并发读写问题。
2. Create 会构造新订单并放入 store。
3. Get 按 orderID + customerID 查找。
4. Update 通过 updateFn 回调更新目标订单。

为什么重要：
- 你学到的是“面向接口编程”，不是把逻辑写死在一个 struct 里。
- 先用 in-memory 快速跑通流程，后续再替换 MySQL/Mongo 都不影响应用层。

---

## 五、Lesson9：Query 模式 + 装饰器

### 关键新增
- internal/order/app/query/get_customer_order.go
- internal/common/decorator/query.go
- internal/common/decorator/logging.go
- internal/common/decorator/metrics.go
- internal/common/metrics/todo_metrics.go

### 必须理解点
1. QueryHandler[Q, R]
   - 泛型接口，统一查询处理器形态。
2. ApplyQueryDecorators
   - 给 handler 套上日志和指标能力。
3. get_customer_order.go
   - 只做一件事：调用 orderRepo.Get。

为什么重要：
- 业务处理器保持单一职责。
- 日志/指标这种横切逻辑不污染业务代码。

---

## 六、Lesson10：Command 模式（创建/更新订单）

### 关键新增
- internal/order/app/command/create_order.go
- internal/order/app/command/update_order.go
- internal/common/decorator/command.go

### 必须理解点

#### create_order.go
1. Handle() 主流程：
   - validate(items)
   - orderRepo.Create(...)
   - 返回 CreateOrderResult
2. validate()：
   - 先检查 items 非空
   - packItems 合并重复商品数量
   - 调用 stockGRPC.CheckIfItemsInStock

这是创建订单最核心流程，建议反复读。

#### update_order.go
1. 通过 UpdateFn 把“更新细节”以回调方式传入。
2. 仓储层只负责持久化更新，不负责决定如何更新字段。

为什么重要：
- command 和 query 分离后，读写逻辑更清晰。
- 以后做权限、审计、重试都更好扩展。

---

## 七、Lesson11：Order 调 Stock 的 gRPC 适配

### 关键新增
- internal/order/app/query/service.go
- internal/order/adapters/grpc/stock_grpc.go

### 必须理解点
1. service.go 定义 StockService 接口（端口）。
2. stock_grpc.go 是该接口的 gRPC 版本实现（适配器）。
3. create_order.go 通过 StockService 工作，而不是直接依赖 grpc client。
4. application.go 负责把 grpc client 注入给 handler。

为什么重要：
- 这就是“六边形架构/端口适配器”的核心实践。
- 你可以轻松换成 mock、HTTP、消息队列实现，不改业务流程。

---

## 八、给 Go 小白的检查清单（完成 Lesson5-11 后）

你需要确认自己已经能回答下面问题：
1. 为什么 main.go 不应该直接 new 所有业务对象？
2. 为什么先定义 Repository 接口再写 in-memory 实现？
3. 为什么 command/query 要加 decorator？
4. 为什么 order 调 stock 要先定义 StockService 接口？
5. 如果把 in-memory 换成 MySQL，哪些层会变，哪些层不变？

如果这 5 个问题都能讲清楚，说明你已经掌握这个阶段的核心设计思想。

---

## 九、下一批计划

下一批会继续写入 Lesson12-Lesson20，重点是：
1. payment 服务接入
2. 事件发布与消费
3. webhook 与 processor
4. 早期可观测性实践

我会保持同样格式：
- 每个 lesson 的关键文件
- 这些文件干了什么
- 必须理解的代码点
- 调用链

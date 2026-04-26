# 逐行详细解读（Part 1）

范围：`13d1e05 -> 8ee5f39`
说明：本文件按“代码行号”解释，先覆盖前两次关键提交。

---

## Commit: d6b3140 (`genproto, genopenapi`)

### 文件：`internal/common/config/viper.go`

代码（提交当时版本）：

```go
1: package config
2:
3: import "github.com/spf13/viper"
4:
5: func NewViperConfig() error {
6: 	viper.SetConfigName("global")
7: 	viper.SetConfigType("yaml")
8: 	viper.AddConfigPath("../common/config")
9: 	viper.AutomaticEnv()
10:	return viper.ReadInConfig()
11: }
```

逐行解释：

- `1`：声明包名 `config`，表示这个文件负责“配置加载”能力。
- `2`：空行，只是分隔结构，提升可读性。
- `3`：引入 `viper`，这是 Go 里常见的配置库。
- `4`：空行，分隔 import 与函数定义。
- `5`：定义 `NewViperConfig()`，返回 `error`，意味着调用方必须处理配置加载失败。
- `6`：配置文件名是 `global`（不含后缀），后续会与类型拼成 `global.yaml`。
- `7`：配置文件类型设为 `yaml`，告诉 viper 用 YAML 解析器。
- `8`：配置文件搜索目录是 `../common/config`，这是相对运行目录的路径。
- `9`：开启环境变量覆盖能力，允许通过环境变量覆盖配置值。
- `10`：真正读取配置文件，成功返回 `nil`，失败返回具体错误。
- `11`：函数结束。

这段代码的本质：把“配置初始化流程”集中在一个函数里，后续所有服务 `main` 都可复用。

---

## Commit: 0059266 (`http, grpc 服务搭建`)

### 文件：`internal/common/server/http.go`

代码（提交当时版本）：

```go
1: package server
2:
3: import (
4: 	"github.com/gin-gonic/gin"
5: 	"github.com/spf13/viper"
6: )
7:
8: func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
9: 	addr := viper.Sub(serviceName).GetString("http-addr")
10:	if addr == "" {
11:		// TODO: Warning log
12:	}
13:	RunHTTPServerOnAddr(addr, wrapper)
14: }
15:
16: func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
17:	apiRouter := gin.New()
18:	wrapper(apiRouter)
19:	apiRouter.Group("/api")
20:	if err := apiRouter.Run(addr); err != nil {
21:		panic(err)
22:	}
23: }
```

逐行解释：

- `1`：公共服务器能力放在 `server` 包。
- `2`：空行。
- `3`：开始 import 块。
- `4`：引入 Gin，用于 HTTP 路由与服务。
- `5`：引入 Viper，用于读配置中的监听地址。
- `6`：import 块结束。
- `7`：空行。
- `8`：定义统一 HTTP 启动函数；`wrapper` 参数用于把业务路由注册进来。
- `9`：从 `serviceName` 子配置里读取 `http-addr`。
- `10`：判断地址是否为空。
- `11`：TODO 注释，表示后续应记录告警日志。
- `12`：if 结束。
- `13`：把读取到的地址交给底层函数启动。
- `14`：函数结束。
- `15`：空行。
- `16`：定义底层按地址启动函数。
- `17`：创建全新 Gin 引擎（无默认中间件）。
- `18`：执行 `wrapper`，由外部注入路由。
- `19`：创建 `/api` 路由组，但这里没有保存返回值，实际没有继续挂载子路由（后续可优化）。
- `20`：启动 HTTP 服务监听 `addr`。
- `21`：启动失败直接 panic，进程退出。
- `22`：if 结束。
- `23`：函数结束。

---

### 文件：`internal/common/server/gprc.go`

代码（提交当时版本）：

```go
1: package server
2:
3: import (
4: 	"net"
5:
6: 	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
7: 	grpc_tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
8: 	"github.com/sirupsen/logrus"
9: 	"github.com/spf13/viper"
10:	"google.golang.org/grpc"
11: )
12:
13: func init() {
14:	logger := logrus.New()
15:	logger.SetLevel(logrus.WarnLevel)
16:	grpc_logrus.ReplaceGrpcLogger(logrus.NewEntry(logger))
17: }
18:
19: func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server)) {
20:	addr := viper.Sub(serviceName).GetString("grpc-addr")
21:	if addr == "" {
22:		// TODO: Warning log
23:		addr = viper.GetString("fallback-grpc-addr")
24:	}
25:	RunGRPCServerOnAddr(addr, registerServer)
26: }
27:
28: func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
29:	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
30:	grpcServer := grpc.NewServer(
31:		grpc.ChainUnaryInterceptor(
32:			grpc_tags.UnaryServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
33:			grpc_logrus.UnaryServerInterceptor(logrusEntry),
34:			//otelgrpc.UnaryServerInterceptor(),
35:			//srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
36:			//logging.UnaryServerInterceptor(interceptorLogger(rpcLogger), logging.WithFieldsFromContext(logTraceID)),
37:			//selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authFn), selector.MatchFunc(allButHealthZ)),
38:			//recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
39:		),
40:		grpc.ChainStreamInterceptor(
41:			grpc_tags.StreamServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
42:			grpc_logrus.StreamServerInterceptor(logrusEntry),
43:			//otelgrpc.StreamServerInterceptor(),
44:			//srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
45:			//logging.StreamServerInterceptor(interceptorLogger(rpcLogger), logging.WithFieldsFromContext(logTraceID)),
46:			//selector.StreamServerInterceptor(auth.StreamServerInterceptor(authFn), selector.MatchFunc(allButHealthZ)),
47:			//recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
48:		),
49:	)
50:	registerServer(grpcServer)
51:
52:	listen, err := net.Listen("tcp", addr)
53:	if err != nil {
54:		logrus.Panic(err)
55:	}
56:	logrus.Infof("Starting gRPC server, Listening: %s", addr)
57:	if err := grpcServer.Serve(listen); err != nil {
58:		logrus.Panic(err)
59:	}
60: }
```

逐行解释：

- `1`：同属 `server` 包，代表 HTTP/gRPC 启动能力在同一层。
- `2`：空行。
- `3`：开始 import 块。
- `4`：`net` 用于 `net.Listen` 创建 TCP 监听器。
- `5`：空行，分隔标准库和第三方库。
- `6`：引入 gRPC + logrus 的日志中间件。
- `7`：引入 gRPC tags 中间件，用于把请求字段打到日志上下文。
- `8`：引入 logrus 日志库。
- `9`：引入 viper 读配置。
- `10`：引入 gRPC 主包。
- `11`：import 结束。
- `12`：空行。
- `13`：`init()` 在包加载时执行，用于初始化 gRPC 日志行为。
- `14`：创建一个新的 logrus logger。
- `15`：把日志级别设为 `Warn`，默认过滤较低级别日志。
- `16`：把 gRPC 内部 logger 替换成 logrus。
- `17`：`init` 结束。
- `18`：空行。
- `19`：高层启动函数，按服务名读取配置并注册业务服务。
- `20`：读取当前服务的 `grpc-addr`。
- `21`：如果地址为空。
- `22`：TODO：应记录配置缺失告警。
- `23`：回退到全局 `fallback-grpc-addr`。
- `24`：if 结束。
- `25`：调用底层地址启动函数。
- `26`：函数结束。
- `27`：空行。
- `28`：底层 gRPC 启动实现。
- `29`：创建日志 entry 给中间件复用。
- `30`：创建 gRPC server。
- `31`：开始装配 unary 拦截器链。
- `32`：第一个 unary 拦截器：提取请求字段用于标签。
- `33`：第二个 unary 拦截器：记录日志。
- `34`：注释位：未来可接入 OpenTelemetry unary 拦截器。
- `35`：注释位：未来可接入指标拦截器。
- `36`：注释位：未来可接入结构化 logging 拦截器。
- `37`：注释位：未来可接入鉴权拦截器。
- `38`：注释位：未来可接入 panic recovery 拦截器。
- `39`：unary 链结束。
- `40`：开始装配 stream 拦截器链。
- `41`：stream 的 tags 拦截器。
- `42`：stream 的 logrus 拦截器。
- `43`：注释位：OpenTelemetry stream 拦截器。
- `44`：注释位：指标 stream 拦截器。
- `45`：注释位：结构化日志 stream 拦截器。
- `46`：注释位：鉴权 stream 拦截器。
- `47`：注释位：recovery stream 拦截器。
- `48`：stream 链结束。
- `49`：`grpc.NewServer` 参数结束。
- `50`：调用外部注入函数注册具体业务服务实现。
- `51`：空行。
- `52`：监听 TCP 地址。
- `53`：监听失败则进入错误分支。
- `54`：panic 退出，避免服务处于不确定状态。
- `55`：if 结束。
- `56`：打印启动成功日志。
- `57`：开始阻塞服务处理请求，若返回错误进入分支。
- `58`：服务异常退出时 panic。
- `59`：if 结束。
- `60`：函数结束。

---

### 文件：`internal/order/main.go`

代码（提交当时版本）：

```go
1: package main
2:
3: import (
4: 	"log"
5:
6: 	"github.com/ghost-yu/go_shop_second/common/config"
7: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
8: 	"github.com/ghost-yu/go_shop_second/common/server"
9: 	"github.com/ghost-yu/go_shop_second/order/ports"
10:	"github.com/gin-gonic/gin"
11:	"github.com/spf13/viper"
12:	"google.golang.org/grpc"
13: )
14:
15: func init() {
16:	if err := config.NewViperConfig(); err != nil {
17:		log.Fatal(err)
18:	}
19: }
20:
21: func main() {
22:	serviceName := viper.GetString("order.service-name")
23:
24:	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
25:		svc := ports.NewGRPCServer()
26:		orderpb.RegisterOrderServiceServer(server, svc)
27:	})
28:
29:	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
30:		ports.RegisterHandlersWithOptions(router, HTTPServer{}, ports.GinServerOptions{
31:			BaseURL:      "/api",
32:			Middlewares:  nil,
33:			ErrorHandler: nil,
34:		})
35:	})
36:
37: }
```

逐行解释：

- `1`：Order 服务入口包。
- `2`：空行。
- `3`：开始 import。
- `4`：标准库 `log`，用于启动失败直接输出并退出。
- `5`：空行。
- `6`：引入配置初始化函数。
- `7`：引入 gRPC 生成的 order 服务接口定义。
- `8`：引入通用 server 启动封装。
- `9`：引入 order 传输层实现（ports）。
- `10`：Gin 引擎类型用于 HTTP wrapper。
- `11`：Viper 用于读取 `service-name`。
- `12`：gRPC server 类型。
- `13`：import 结束。
- `14`：空行。
- `15`：`init` 阶段先加载配置。
- `16`：尝试读取配置文件。
- `17`：失败即 `Fatal`，进程退出（fail-fast）。
- `18`：if 结束。
- `19`：`init` 结束。
- `20`：空行。
- `21`：主函数入口。
- `22`：读取 order 服务名，用于查配置分组。
- `23`：空行。
- `24`：开 goroutine 启动 gRPC 服务，不阻塞主线程。
- `25`：构造 gRPC 端口服务实现对象。
- `26`：把实现注册到 gRPC server。
- `27`：gRPC wrapper 回调结束。
- `28`：空行。
- `29`：主线程启动 HTTP 服务（阻塞）。
- `30`：通过 OpenAPI 生成的路由注册函数绑定处理器。
- `31`：HTTP 路由前缀设置为 `/api`。
- `32`：中间件暂未注入。
- `33`：自定义错误处理器暂未注入。
- `34`：HTTP 选项结构结束。
- `35`：HTTP wrapper 回调结束。
- `36`：空行。
- `37`：main 结束。

---

### 文件：`internal/order/http.go`

代码（提交当时版本）：

```go
1: package main
2:
3: import (
4: 	"github.com/gin-gonic/gin"
5: )
6:
7: type HTTPServer struct{}
8:
9: func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
10:	//TODO implement me
11:	panic("implement me")
12: }
13:
14: func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
15:	//TODO implement me
16:	panic("implement me")
17: }
```

逐行解释：

- `1`：和入口同包，便于直接在 `main` 中实例化。
- `2`：空行。
- `3`：import 开始。
- `4`：引入 Gin 请求上下文类型。
- `5`：import 结束。
- `6`：空行。
- `7`：HTTP 处理器结构体，目前无字段。
- `8`：空行。
- `9`：POST 下单接口方法签名（由 OpenAPI 生成接口约束）。
- `10`：TODO 说明业务逻辑尚未实现。
- `11`：显式 panic，避免“未实现却默默成功”。
- `12`：方法结束。
- `13`：空行。
- `14`：GET 查询订单接口方法签名。
- `15`：TODO 未实现。
- `16`：panic，提醒当前不可用。
- `17`：方法结束。

---

### 文件：`internal/order/ports/grpc.go`

代码（提交当时版本）：

```go
1: package ports
2:
3: import (
4: 	context "context"
5:
6: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
7: 	"google.golang.org/protobuf/types/known/emptypb"
8: )
9:
10: type GRPCServer struct {
11: }
12:
13: func NewGRPCServer() *GRPCServer {
14:	return &GRPCServer{}
15: }
16:
17: func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
18:	//TODO implement me
19:	panic("implement me")
20: }
21:
22: func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
23:	//TODO implement me
24:	panic("implement me")
25: }
26:
27: func (G GRPCServer) UpdateOrder(ctx context.Context, order *orderpb.Order) (*emptypb.Empty, error) {
28:	//TODO implement me
29:	panic("implement me")
30: }
```

逐行解释：

- `1`：放在 `ports` 包，表示“传输适配层”。
- `2`：空行。
- `3`：import 开始。
- `4`：引入 `context`，处理请求级超时/取消。
- `5`：空行。
- `6`：引入 order gRPC 协议类型。
- `7`：引入 protobuf 空响应类型。
- `8`：import 结束。
- `9`：空行。
- `10`：定义 gRPC server 实现结构体。
- `11`：当前无字段，后续通常会挂 `app.Application`。
- `12`：空行。
- `13`：构造函数。
- `14`：返回空实现对象。
- `15`：构造函数结束。
- `16`：空行。
- `17`：实现 `CreateOrder` RPC 方法签名。
- `18`：TODO。
- `19`：panic 代表未实现。
- `20`：方法结束。
- `21`：空行。
- `22`：实现 `GetOrder` RPC 方法签名。
- `23`：TODO。
- `24`：panic。
- `25`：方法结束。
- `26`：空行。
- `27`：实现 `UpdateOrder` RPC 方法签名。
- `28`：TODO。
- `29`：panic。
- `30`：方法结束。

---

### 文件：`internal/stock/main.go`

代码（提交当时版本）：

```go
1: package main
2:
3: import (
4: 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
5: 	"github.com/ghost-yu/go_shop_second/common/server"
6: 	"github.com/ghost-yu/go_shop_second/stock/ports"
7: 	"github.com/spf13/viper"
8: 	"google.golang.org/grpc"
9: )
10:
11: func main() {
12:	serviceName := viper.GetString("stock.service-name")
13:	serverType := viper.GetString("stock.server-to-run")
14:
15:	switch serverType {
16:	case "grpc":
17:		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
18:			svc := ports.NewGRPCServer()
19:			stockpb.RegisterStockServiceServer(server, svc)
20:		})
21:	case "http":
22:		// 暂时不用
23:	default:
24:		panic("unexpected server type")
25:	}
26: }
```

逐行解释：

- `1`：Stock 服务入口。
- `2`：空行。
- `3`：import 开始。
- `4`：引入 stock gRPC 生成协议。
- `5`：引入通用 server 启动能力。
- `6`：引入 stock 端口实现。
- `7`：读配置。
- `8`：gRPC server 类型。
- `9`：import 结束。
- `10`：空行。
- `11`：主函数入口。
- `12`：读取 stock 服务名。
- `13`：读取要运行哪种协议服务（grpc/http）。
- `14`：空行。
- `15`：按配置分支启动。
- `16`：grpc 模式。
- `17`：调用统一 gRPC 启动器。
- `18`：构造 stock gRPC 实现。
- `19`：注册 stock 服务。
- `20`：回调结束。
- `21`：http 模式分支。
- `22`：注释说明暂不实现。
- `23`：兜底分支。
- `24`：未知模式直接 panic。
- `25`：switch 结束。
- `26`：main 结束。

---

### 文件：`internal/stock/ports/grpc.go`

代码（提交当时版本）：

```go
1: package ports
2:
3: import (
4: 	context "context"
5:
6: 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
7: )
8:
9: type GRPCServer struct {
10: }
11:
12: func NewGRPCServer() *GRPCServer {
13:	return &GRPCServer{}
14: }
15:
16: func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
17:	//TODO implement me
18:	panic("implement me")
19: }
20:
21: func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
22:	//TODO implement me
23:	panic("implement me")
24: }
```

逐行解释：

- `1`：传输层包。
- `2`：空行。
- `3`：import 开始。
- `4`：请求上下文。
- `5`：空行。
- `6`：stock 协议类型。
- `7`：import 结束。
- `8`：空行。
- `9`：gRPC 服务实现结构体。
- `10`：当前无依赖字段。
- `11`：空行。
- `12`：构造函数。
- `13`：返回实例。
- `14`：结束。
- `15`：空行。
- `16`：实现获取商品接口签名。
- `17`：TODO。
- `18`：panic。
- `19`：方法结束。
- `20`：空行。
- `21`：实现库存校验接口签名。
- `22`：TODO。
- `23`：panic。
- `24`：方法结束。

---

下一步我会继续 Part 2：`49bfa8e` 到 `606345b`，仍然按这种逐行方式写。

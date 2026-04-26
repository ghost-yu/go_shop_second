# Line-by-Line Notes (Auto Generated)

- Repo: go_shop_second
- FromCommit: 13d1e05
- ToCommit: 8ee5f39
- GeneratedAt: 2026-04-06 18:03:33 +08:00
- Scope: *.go only, excluding *.pb.go and *.gen.go

## Commit 1: [13d1e05] config, proto, openapi

### File: internal/order/main.go

```go
   1: package main
   2: 
   3: import (
   4: 	"log"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/config"
   7: 	"github.com/spf13/viper"
   8: )
   9: 
  10: func init() {
  11: 	if err := config.NewViperConfig(); err != nil {
  12: 		log.Fatal(err)
  13: 	}
  14: }
  15: 
  16: func main() {
  17: 	log.Printf("%v", viper.Get("order"))
  18: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Block end: closes import or parameter list block. |
| 9 | Blank line: separates code blocks for readability. |
| 10 | Function/method declaration: defines executable logic. |
| 11 | Conditional branch: executes following block when condition is true. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Block end: closes previous function/branch/type block. |
| 14 | Block end: closes previous function/branch/type block. |
| 15 | Blank line: separates code blocks for readability. |
| 16 | Function/method declaration: defines executable logic. |
| 17 | Executable statement: participates in current implementation logic. |
| 18 | Block end: closes previous function/branch/type block. |

## Commit 2: [d6b3140] genproto, genopenapi

### File: internal/common/config/viper.go

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
  10: 	return viper.ReadInConfig()
  11: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Single import: brings one dependency into scope. |
| 4 | Blank line: separates code blocks for readability. |
| 5 | Function/method declaration: defines executable logic. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Return statement: returns values (and often error) to caller. |
| 11 | Block end: closes previous function/branch/type block. |

## Commit 3: [0059266] http, grpc 服务搭建

### File: internal/common/server/gprc.go

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
  10: 	"google.golang.org/grpc"
  11: )
  12: 
  13: func init() {
  14: 	logger := logrus.New()
  15: 	logger.SetLevel(logrus.WarnLevel)
  16: 	grpc_logrus.ReplaceGrpcLogger(logrus.NewEntry(logger))
  17: }
  18: 
  19: func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server)) {
  20: 	addr := viper.Sub(serviceName).GetString("grpc-addr")
  21: 	if addr == "" {
  22: 		// TODO: Warning log
  23: 		addr = viper.GetString("fallback-grpc-addr")
  24: 	}
  25: 	RunGRPCServerOnAddr(addr, registerServer)
  26: }
  27: 
  28: func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
  29: 	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
  30: 	grpcServer := grpc.NewServer(
  31: 		grpc.ChainUnaryInterceptor(
  32: 			grpc_tags.UnaryServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
  33: 			grpc_logrus.UnaryServerInterceptor(logrusEntry),
  34: 			//otelgrpc.UnaryServerInterceptor(),
  35: 			//srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
  36: 			//logging.UnaryServerInterceptor(interceptorLogger(rpcLogger), logging.WithFieldsFromContext(logTraceID)),
  37: 			//selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authFn), selector.MatchFunc(allButHealthZ)),
  38: 			//recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
  39: 		),
  40: 		grpc.ChainStreamInterceptor(
  41: 			grpc_tags.StreamServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
  42: 			grpc_logrus.StreamServerInterceptor(logrusEntry),
  43: 			//otelgrpc.StreamServerInterceptor(),
  44: 			//srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
  45: 			//logging.StreamServerInterceptor(interceptorLogger(rpcLogger), logging.WithFieldsFromContext(logTraceID)),
  46: 			//selector.StreamServerInterceptor(auth.StreamServerInterceptor(authFn), selector.MatchFunc(allButHealthZ)),
  47: 			//recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
  48: 		),
  49: 	)
  50: 	registerServer(grpcServer)
  51: 
  52: 	listen, err := net.Listen("tcp", addr)
  53: 	if err != nil {
  54: 		logrus.Panic(err)
  55: 	}
  56: 	logrus.Infof("Starting gRPC server, Listening: %s", addr)
  57: 	if err := grpcServer.Serve(listen); err != nil {
  58: 		logrus.Panic(err)
  59: 	}
  60: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Block end: closes import or parameter list block. |
| 12 | Blank line: separates code blocks for readability. |
| 13 | Function/method declaration: defines executable logic. |
| 14 | Short declaration: defines and initializes local variable. |
| 15 | Executable statement: participates in current implementation logic. |
| 16 | Executable statement: participates in current implementation logic. |
| 17 | Block end: closes previous function/branch/type block. |
| 18 | Blank line: separates code blocks for readability. |
| 19 | Function/method declaration: defines executable logic. |
| 20 | Short declaration: defines and initializes local variable. |
| 21 | Conditional branch: executes following block when condition is true. |
| 22 | Comment: explains intent, todo, or context. |
| 23 | Assignment: sets or updates variable/field value. |
| 24 | Block end: closes previous function/branch/type block. |
| 25 | Executable statement: participates in current implementation logic. |
| 26 | Block end: closes previous function/branch/type block. |
| 27 | Blank line: separates code blocks for readability. |
| 28 | Function/method declaration: defines executable logic. |
| 29 | Short declaration: defines and initializes local variable. |
| 30 | Short declaration: defines and initializes local variable. |
| 31 | Executable statement: participates in current implementation logic. |
| 32 | Executable statement: participates in current implementation logic. |
| 33 | Executable statement: participates in current implementation logic. |
| 34 | Comment: explains intent, todo, or context. |
| 35 | Comment: explains intent, todo, or context. |
| 36 | Comment: explains intent, todo, or context. |
| 37 | Comment: explains intent, todo, or context. |
| 38 | Comment: explains intent, todo, or context. |
| 39 | Executable statement: participates in current implementation logic. |
| 40 | Executable statement: participates in current implementation logic. |
| 41 | Executable statement: participates in current implementation logic. |
| 42 | Executable statement: participates in current implementation logic. |
| 43 | Comment: explains intent, todo, or context. |
| 44 | Comment: explains intent, todo, or context. |
| 45 | Comment: explains intent, todo, or context. |
| 46 | Comment: explains intent, todo, or context. |
| 47 | Comment: explains intent, todo, or context. |
| 48 | Executable statement: participates in current implementation logic. |
| 49 | Block end: closes import or parameter list block. |
| 50 | Executable statement: participates in current implementation logic. |
| 51 | Blank line: separates code blocks for readability. |
| 52 | Short declaration: defines and initializes local variable. |
| 53 | Conditional branch: executes following block when condition is true. |
| 54 | Executable statement: participates in current implementation logic. |
| 55 | Block end: closes previous function/branch/type block. |
| 56 | Executable statement: participates in current implementation logic. |
| 57 | Conditional branch: executes following block when condition is true. |
| 58 | Executable statement: participates in current implementation logic. |
| 59 | Block end: closes previous function/branch/type block. |
| 60 | Block end: closes previous function/branch/type block. |

### File: internal/common/server/http.go

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
  10: 	if addr == "" {
  11: 		// TODO: Warning log
  12: 	}
  13: 	RunHTTPServerOnAddr(addr, wrapper)
  14: }
  15: 
  16: func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
  17: 	apiRouter := gin.New()
  18: 	wrapper(apiRouter)
  19: 	apiRouter.Group("/api")
  20: 	if err := apiRouter.Run(addr); err != nil {
  21: 		panic(err)
  22: 	}
  23: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Executable statement: participates in current implementation logic. |
| 6 | Block end: closes import or parameter list block. |
| 7 | Blank line: separates code blocks for readability. |
| 8 | Function/method declaration: defines executable logic. |
| 9 | Short declaration: defines and initializes local variable. |
| 10 | Conditional branch: executes following block when condition is true. |
| 11 | Comment: explains intent, todo, or context. |
| 12 | Block end: closes previous function/branch/type block. |
| 13 | Executable statement: participates in current implementation logic. |
| 14 | Block end: closes previous function/branch/type block. |
| 15 | Blank line: separates code blocks for readability. |
| 16 | Function/method declaration: defines executable logic. |
| 17 | Short declaration: defines and initializes local variable. |
| 18 | Executable statement: participates in current implementation logic. |
| 19 | Executable statement: participates in current implementation logic. |
| 20 | Conditional branch: executes following block when condition is true. |
| 21 | Panic: aborts normal control flow with unrecoverable error. |
| 22 | Block end: closes previous function/branch/type block. |
| 23 | Block end: closes previous function/branch/type block. |

### File: internal/order/http.go

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
  10: 	//TODO implement me
  11: 	panic("implement me")
  12: }
  13: 
  14: func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
  15: 	//TODO implement me
  16: 	panic("implement me")
  17: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Block end: closes import or parameter list block. |
| 6 | Blank line: separates code blocks for readability. |
| 7 | Type declaration: defines struct/interface/type alias. |
| 8 | Blank line: separates code blocks for readability. |
| 9 | Function/method declaration: defines executable logic. |
| 10 | Comment: explains intent, todo, or context. |
| 11 | Panic: aborts normal control flow with unrecoverable error. |
| 12 | Block end: closes previous function/branch/type block. |
| 13 | Blank line: separates code blocks for readability. |
| 14 | Function/method declaration: defines executable logic. |
| 15 | Comment: explains intent, todo, or context. |
| 16 | Panic: aborts normal control flow with unrecoverable error. |
| 17 | Block end: closes previous function/branch/type block. |

### File: internal/order/main.go

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
  10: 	"github.com/gin-gonic/gin"
  11: 	"github.com/spf13/viper"
  12: 	"google.golang.org/grpc"
  13: )
  14: 
  15: func init() {
  16: 	if err := config.NewViperConfig(); err != nil {
  17: 		log.Fatal(err)
  18: 	}
  19: }
  20: 
  21: func main() {
  22: 	serviceName := viper.GetString("order.service-name")
  23: 
  24: 	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
  25: 		svc := ports.NewGRPCServer()
  26: 		orderpb.RegisterOrderServiceServer(server, svc)
  27: 	})
  28: 
  29: 	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
  30: 		ports.RegisterHandlersWithOptions(router, HTTPServer{}, ports.GinServerOptions{
  31: 			BaseURL:      "/api",
  32: 			Middlewares:  nil,
  33: 			ErrorHandler: nil,
  34: 		})
  35: 	})
  36: 
  37: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Executable statement: participates in current implementation logic. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Block end: closes import or parameter list block. |
| 14 | Blank line: separates code blocks for readability. |
| 15 | Function/method declaration: defines executable logic. |
| 16 | Conditional branch: executes following block when condition is true. |
| 17 | Executable statement: participates in current implementation logic. |
| 18 | Block end: closes previous function/branch/type block. |
| 19 | Block end: closes previous function/branch/type block. |
| 20 | Blank line: separates code blocks for readability. |
| 21 | Function/method declaration: defines executable logic. |
| 22 | Short declaration: defines and initializes local variable. |
| 23 | Blank line: separates code blocks for readability. |
| 24 | Goroutine launch: runs function asynchronously. |
| 25 | Short declaration: defines and initializes local variable. |
| 26 | Executable statement: participates in current implementation logic. |
| 27 | Block end: closes previous function/branch/type block. |
| 28 | Blank line: separates code blocks for readability. |
| 29 | Executable statement: participates in current implementation logic. |
| 30 | Executable statement: participates in current implementation logic. |
| 31 | Executable statement: participates in current implementation logic. |
| 32 | Executable statement: participates in current implementation logic. |
| 33 | Executable statement: participates in current implementation logic. |
| 34 | Block end: closes previous function/branch/type block. |
| 35 | Block end: closes previous function/branch/type block. |
| 36 | Blank line: separates code blocks for readability. |
| 37 | Block end: closes previous function/branch/type block. |

### File: internal/order/ports/grpc.go

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
  14: 	return &GRPCServer{}
  15: }
  16: 
  17: func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
  18: 	//TODO implement me
  19: 	panic("implement me")
  20: }
  21: 
  22: func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
  23: 	//TODO implement me
  24: 	panic("implement me")
  25: }
  26: 
  27: func (G GRPCServer) UpdateOrder(ctx context.Context, order *orderpb.Order) (*emptypb.Empty, error) {
  28: 	//TODO implement me
  29: 	panic("implement me")
  30: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Block end: closes import or parameter list block. |
| 9 | Blank line: separates code blocks for readability. |
| 10 | Type declaration: defines struct/interface/type alias. |
| 11 | Block end: closes previous function/branch/type block. |
| 12 | Blank line: separates code blocks for readability. |
| 13 | Function/method declaration: defines executable logic. |
| 14 | Return statement: returns values (and often error) to caller. |
| 15 | Block end: closes previous function/branch/type block. |
| 16 | Blank line: separates code blocks for readability. |
| 17 | Function/method declaration: defines executable logic. |
| 18 | Comment: explains intent, todo, or context. |
| 19 | Panic: aborts normal control flow with unrecoverable error. |
| 20 | Block end: closes previous function/branch/type block. |
| 21 | Blank line: separates code blocks for readability. |
| 22 | Function/method declaration: defines executable logic. |
| 23 | Comment: explains intent, todo, or context. |
| 24 | Panic: aborts normal control flow with unrecoverable error. |
| 25 | Block end: closes previous function/branch/type block. |
| 26 | Blank line: separates code blocks for readability. |
| 27 | Function/method declaration: defines executable logic. |
| 28 | Comment: explains intent, todo, or context. |
| 29 | Panic: aborts normal control flow with unrecoverable error. |
| 30 | Block end: closes previous function/branch/type block. |

### File: internal/stock/main.go

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
  12: 	serviceName := viper.GetString("stock.service-name")
  13: 	serverType := viper.GetString("stock.server-to-run")
  14: 
  15: 	switch serverType {
  16: 	case "grpc":
  17: 		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
  18: 			svc := ports.NewGRPCServer()
  19: 			stockpb.RegisterStockServiceServer(server, svc)
  20: 		})
  21: 	case "http":
  22: 		// 暂时不用
  23: 	default:
  24: 		panic("unexpected server type")
  25: 	}
  26: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Executable statement: participates in current implementation logic. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Block end: closes import or parameter list block. |
| 10 | Blank line: separates code blocks for readability. |
| 11 | Function/method declaration: defines executable logic. |
| 12 | Short declaration: defines and initializes local variable. |
| 13 | Short declaration: defines and initializes local variable. |
| 14 | Blank line: separates code blocks for readability. |
| 15 | Switch branch: multi-way branching by expression value. |
| 16 | Case branch: logic for one matching switch value. |
| 17 | Executable statement: participates in current implementation logic. |
| 18 | Short declaration: defines and initializes local variable. |
| 19 | Executable statement: participates in current implementation logic. |
| 20 | Block end: closes previous function/branch/type block. |
| 21 | Case branch: logic for one matching switch value. |
| 22 | Comment: explains intent, todo, or context. |
| 23 | Default branch: fallback when no case matches. |
| 24 | Panic: aborts normal control flow with unrecoverable error. |
| 25 | Block end: closes previous function/branch/type block. |
| 26 | Block end: closes previous function/branch/type block. |

### File: internal/stock/ports/grpc.go

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
  13: 	return &GRPCServer{}
  14: }
  15: 
  16: func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
  17: 	//TODO implement me
  18: 	panic("implement me")
  19: }
  20: 
  21: func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
  22: 	//TODO implement me
  23: 	panic("implement me")
  24: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Block end: closes import or parameter list block. |
| 8 | Blank line: separates code blocks for readability. |
| 9 | Type declaration: defines struct/interface/type alias. |
| 10 | Block end: closes previous function/branch/type block. |
| 11 | Blank line: separates code blocks for readability. |
| 12 | Function/method declaration: defines executable logic. |
| 13 | Return statement: returns values (and often error) to caller. |
| 14 | Block end: closes previous function/branch/type block. |
| 15 | Blank line: separates code blocks for readability. |
| 16 | Function/method declaration: defines executable logic. |
| 17 | Comment: explains intent, todo, or context. |
| 18 | Panic: aborts normal control flow with unrecoverable error. |
| 19 | Block end: closes previous function/branch/type block. |
| 20 | Blank line: separates code blocks for readability. |
| 21 | Function/method declaration: defines executable logic. |
| 22 | Comment: explains intent, todo, or context. |
| 23 | Panic: aborts normal control flow with unrecoverable error. |
| 24 | Block end: closes previous function/branch/type block. |

## Commit 4: [49bfa8e] add app to servers

### File: internal/stock/app/app.go

```go
   1: package app
   2: 
   3: type Application struct {
   4: 	Commands Commands
   5: 	Queries  Queries
   6: }
   7: 
   8: type Commands struct{}
   9: 
  10: type Queries struct{}
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Type declaration: defines struct/interface/type alias. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Executable statement: participates in current implementation logic. |
| 6 | Block end: closes previous function/branch/type block. |
| 7 | Blank line: separates code blocks for readability. |
| 8 | Type declaration: defines struct/interface/type alias. |
| 9 | Blank line: separates code blocks for readability. |
| 10 | Type declaration: defines struct/interface/type alias. |

### File: internal/stock/main.go

```go
   1: package main
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/config"
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
   8: 	"github.com/ghost-yu/go_shop_second/common/server"
   9: 	"github.com/ghost-yu/go_shop_second/stock/ports"
  10: 	"github.com/ghost-yu/go_shop_second/stock/service"
  11: 	"github.com/sirupsen/logrus"
  12: 	"github.com/spf13/viper"
  13: 	"google.golang.org/grpc"
  14: )
  15: 
  16: func init() {
  17: 	if err := config.NewViperConfig(); err != nil {
  18: 		logrus.Fatal(err)
  19: 	}
  20: }
  21: 
  22: func main() {
  23: 	serviceName := viper.GetString("stock.service-name")
  24: 	serverType := viper.GetString("stock.server-to-run")
  25: 
  26: 	logrus.Info(serverType)
  27: 
  28: 	ctx, cancel := context.WithCancel(context.Background())
  29: 	defer cancel()
  30: 
  31: 	application := service.NewApplication(ctx)
  32: 	switch serverType {
  33: 	case "grpc":
  34: 		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
  35: 			svc := ports.NewGRPCServer(application)
  36: 			stockpb.RegisterStockServiceServer(server, svc)
  37: 		})
  38: 	case "http":
  39: 		// 暂时不用
  40: 	default:
  41: 		panic("unexpected server type")
  42: 	}
  43: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Executable statement: participates in current implementation logic. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Executable statement: participates in current implementation logic. |
| 14 | Block end: closes import or parameter list block. |
| 15 | Blank line: separates code blocks for readability. |
| 16 | Function/method declaration: defines executable logic. |
| 17 | Conditional branch: executes following block when condition is true. |
| 18 | Executable statement: participates in current implementation logic. |
| 19 | Block end: closes previous function/branch/type block. |
| 20 | Block end: closes previous function/branch/type block. |
| 21 | Blank line: separates code blocks for readability. |
| 22 | Function/method declaration: defines executable logic. |
| 23 | Short declaration: defines and initializes local variable. |
| 24 | Short declaration: defines and initializes local variable. |
| 25 | Blank line: separates code blocks for readability. |
| 26 | Executable statement: participates in current implementation logic. |
| 27 | Blank line: separates code blocks for readability. |
| 28 | Short declaration: defines and initializes local variable. |
| 29 | Deferred call: executes when function returns. |
| 30 | Blank line: separates code blocks for readability. |
| 31 | Short declaration: defines and initializes local variable. |
| 32 | Switch branch: multi-way branching by expression value. |
| 33 | Case branch: logic for one matching switch value. |
| 34 | Executable statement: participates in current implementation logic. |
| 35 | Short declaration: defines and initializes local variable. |
| 36 | Executable statement: participates in current implementation logic. |
| 37 | Block end: closes previous function/branch/type block. |
| 38 | Case branch: logic for one matching switch value. |
| 39 | Comment: explains intent, todo, or context. |
| 40 | Default branch: fallback when no case matches. |
| 41 | Panic: aborts normal control flow with unrecoverable error. |
| 42 | Block end: closes previous function/branch/type block. |
| 43 | Block end: closes previous function/branch/type block. |

### File: internal/stock/ports/grpc.go

```go
   1: package ports
   2: 
   3: import (
   4: 	context "context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
   7: 	"github.com/ghost-yu/go_shop_second/stock/app"
   8: )
   9: 
  10: type GRPCServer struct {
  11: 	app app.Application
  12: }
  13: 
  14: func NewGRPCServer(app app.Application) *GRPCServer {
  15: 	return &GRPCServer{app: app}
  16: }
  17: 
  18: func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
  19: 	//TODO implement me
  20: 	panic("implement me")
  21: }
  22: 
  23: func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
  24: 	//TODO implement me
  25: 	panic("implement me")
  26: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Block end: closes import or parameter list block. |
| 9 | Blank line: separates code blocks for readability. |
| 10 | Type declaration: defines struct/interface/type alias. |
| 11 | Executable statement: participates in current implementation logic. |
| 12 | Block end: closes previous function/branch/type block. |
| 13 | Blank line: separates code blocks for readability. |
| 14 | Function/method declaration: defines executable logic. |
| 15 | Return statement: returns values (and often error) to caller. |
| 16 | Block end: closes previous function/branch/type block. |
| 17 | Blank line: separates code blocks for readability. |
| 18 | Function/method declaration: defines executable logic. |
| 19 | Comment: explains intent, todo, or context. |
| 20 | Panic: aborts normal control flow with unrecoverable error. |
| 21 | Block end: closes previous function/branch/type block. |
| 22 | Blank line: separates code blocks for readability. |
| 23 | Function/method declaration: defines executable logic. |
| 24 | Comment: explains intent, todo, or context. |
| 25 | Panic: aborts normal control flow with unrecoverable error. |
| 26 | Block end: closes previous function/branch/type block. |

### File: internal/stock/service/application.go

```go
   1: package service
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/stock/app"
   7: )
   8: 
   9: func NewApplication(ctx context.Context) app.Application {
  10: 	return app.Application{}
  11: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Block end: closes import or parameter list block. |
| 8 | Blank line: separates code blocks for readability. |
| 9 | Function/method declaration: defines executable logic. |
| 10 | Return statement: returns values (and often error) to caller. |
| 11 | Block end: closes previous function/branch/type block. |

## Commit 5: [49f4e56] add app to servers && air

### File: internal/order/app/app.go

```go
   1: package app
   2: 
   3: type Application struct {
   4: 	Commands Commands
   5: 	Queries  Queries
   6: }
   7: 
   8: type Commands struct{}
   9: 
  10: type Queries struct{}
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Type declaration: defines struct/interface/type alias. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Executable statement: participates in current implementation logic. |
| 6 | Block end: closes previous function/branch/type block. |
| 7 | Blank line: separates code blocks for readability. |
| 8 | Type declaration: defines struct/interface/type alias. |
| 9 | Blank line: separates code blocks for readability. |
| 10 | Type declaration: defines struct/interface/type alias. |

### File: internal/order/http.go

```go
   1: package main
   2: 
   3: import (
   4: 	"github.com/ghost-yu/go_shop_second/order/app"
   5: 	"github.com/gin-gonic/gin"
   6: )
   7: 
   8: type HTTPServer struct {
   9: 	app app.Application
  10: }
  11: 
  12: func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
  13: 	//TODO implement me
  14: 	panic("implement me")
  15: }
  16: 
  17: func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
  18: 	//TODO implement me
  19: 	panic("implement me")
  20: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Executable statement: participates in current implementation logic. |
| 6 | Block end: closes import or parameter list block. |
| 7 | Blank line: separates code blocks for readability. |
| 8 | Type declaration: defines struct/interface/type alias. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Block end: closes previous function/branch/type block. |
| 11 | Blank line: separates code blocks for readability. |
| 12 | Function/method declaration: defines executable logic. |
| 13 | Comment: explains intent, todo, or context. |
| 14 | Panic: aborts normal control flow with unrecoverable error. |
| 15 | Block end: closes previous function/branch/type block. |
| 16 | Blank line: separates code blocks for readability. |
| 17 | Function/method declaration: defines executable logic. |
| 18 | Comment: explains intent, todo, or context. |
| 19 | Panic: aborts normal control flow with unrecoverable error. |
| 20 | Block end: closes previous function/branch/type block. |

### File: internal/order/main.go

```go
   1: package main
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/config"
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   8: 	"github.com/ghost-yu/go_shop_second/common/server"
   9: 	"github.com/ghost-yu/go_shop_second/order/ports"
  10: 	"github.com/ghost-yu/go_shop_second/order/service"
  11: 	"github.com/gin-gonic/gin"
  12: 	"github.com/sirupsen/logrus"
  13: 	"github.com/spf13/viper"
  14: 	"google.golang.org/grpc"
  15: )
  16: 
  17: func init() {
  18: 	if err := config.NewViperConfig(); err != nil {
  19: 		logrus.Fatal(err)
  20: 	}
  21: }
  22: 
  23: func main() {
  24: 	serviceName := viper.GetString("order.service-name")
  25: 
  26: 	ctx, cancel := context.WithCancel(context.Background())
  27: 	defer cancel()
  28: 
  29: 	application := service.NewApplication(ctx)
  30: 
  31: 	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
  32: 		svc := ports.NewGRPCServer(application)
  33: 		orderpb.RegisterOrderServiceServer(server, svc)
  34: 	})
  35: 
  36: 	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
  37: 		ports.RegisterHandlersWithOptions(router, HTTPServer{
  38: 			app: application,
  39: 		}, ports.GinServerOptions{
  40: 			BaseURL:      "/api",
  41: 			Middlewares:  nil,
  42: 			ErrorHandler: nil,
  43: 		})
  44: 	})
  45: 
  46: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Executable statement: participates in current implementation logic. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Executable statement: participates in current implementation logic. |
| 14 | Executable statement: participates in current implementation logic. |
| 15 | Block end: closes import or parameter list block. |
| 16 | Blank line: separates code blocks for readability. |
| 17 | Function/method declaration: defines executable logic. |
| 18 | Conditional branch: executes following block when condition is true. |
| 19 | Executable statement: participates in current implementation logic. |
| 20 | Block end: closes previous function/branch/type block. |
| 21 | Block end: closes previous function/branch/type block. |
| 22 | Blank line: separates code blocks for readability. |
| 23 | Function/method declaration: defines executable logic. |
| 24 | Short declaration: defines and initializes local variable. |
| 25 | Blank line: separates code blocks for readability. |
| 26 | Short declaration: defines and initializes local variable. |
| 27 | Deferred call: executes when function returns. |
| 28 | Blank line: separates code blocks for readability. |
| 29 | Short declaration: defines and initializes local variable. |
| 30 | Blank line: separates code blocks for readability. |
| 31 | Goroutine launch: runs function asynchronously. |
| 32 | Short declaration: defines and initializes local variable. |
| 33 | Executable statement: participates in current implementation logic. |
| 34 | Block end: closes previous function/branch/type block. |
| 35 | Blank line: separates code blocks for readability. |
| 36 | Executable statement: participates in current implementation logic. |
| 37 | Executable statement: participates in current implementation logic. |
| 38 | Executable statement: participates in current implementation logic. |
| 39 | Block end: closes previous function/branch/type block. |
| 40 | Executable statement: participates in current implementation logic. |
| 41 | Executable statement: participates in current implementation logic. |
| 42 | Executable statement: participates in current implementation logic. |
| 43 | Block end: closes previous function/branch/type block. |
| 44 | Block end: closes previous function/branch/type block. |
| 45 | Blank line: separates code blocks for readability. |
| 46 | Block end: closes previous function/branch/type block. |

### File: internal/order/ports/grpc.go

```go
   1: package ports
   2: 
   3: import (
   4: 	context "context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   7: 	"github.com/ghost-yu/go_shop_second/order/app"
   8: 	"google.golang.org/protobuf/types/known/emptypb"
   9: )
  10: 
  11: type GRPCServer struct {
  12: 	app app.Application
  13: }
  14: 
  15: func NewGRPCServer(app app.Application) *GRPCServer {
  16: 	return &GRPCServer{app: app}
  17: }
  18: 
  19: func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
  20: 	//TODO implement me
  21: 	panic("implement me")
  22: }
  23: 
  24: func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
  25: 	//TODO implement me
  26: 	panic("implement me")
  27: }
  28: 
  29: func (G GRPCServer) UpdateOrder(ctx context.Context, order *orderpb.Order) (*emptypb.Empty, error) {
  30: 	//TODO implement me
  31: 	panic("implement me")
  32: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Block end: closes import or parameter list block. |
| 10 | Blank line: separates code blocks for readability. |
| 11 | Type declaration: defines struct/interface/type alias. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Block end: closes previous function/branch/type block. |
| 14 | Blank line: separates code blocks for readability. |
| 15 | Function/method declaration: defines executable logic. |
| 16 | Return statement: returns values (and often error) to caller. |
| 17 | Block end: closes previous function/branch/type block. |
| 18 | Blank line: separates code blocks for readability. |
| 19 | Function/method declaration: defines executable logic. |
| 20 | Comment: explains intent, todo, or context. |
| 21 | Panic: aborts normal control flow with unrecoverable error. |
| 22 | Block end: closes previous function/branch/type block. |
| 23 | Blank line: separates code blocks for readability. |
| 24 | Function/method declaration: defines executable logic. |
| 25 | Comment: explains intent, todo, or context. |
| 26 | Panic: aborts normal control flow with unrecoverable error. |
| 27 | Block end: closes previous function/branch/type block. |
| 28 | Blank line: separates code blocks for readability. |
| 29 | Function/method declaration: defines executable logic. |
| 30 | Comment: explains intent, todo, or context. |
| 31 | Panic: aborts normal control flow with unrecoverable error. |
| 32 | Block end: closes previous function/branch/type block. |

### File: internal/order/service/application.go

```go
   1: package service
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/order/app"
   7: )
   8: 
   9: func NewApplication(ctx context.Context) app.Application {
  10: 	return app.Application{}
  11: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Block end: closes import or parameter list block. |
| 8 | Blank line: separates code blocks for readability. |
| 9 | Function/method declaration: defines executable logic. |
| 10 | Return statement: returns values (and often error) to caller. |
| 11 | Block end: closes previous function/branch/type block. |

## Commit 6: [f87d59e] order stock inmem repo

### File: internal/order/adapters/order_inmem_repository.go

```go
   1: package adapters
   2: 
   3: import (
   4: 	"context"
   5: 	"strconv"
   6: 	"sync"
   7: 	"time"
   8: 
   9: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
  10: 	"github.com/sirupsen/logrus"
  11: )
  12: 
  13: type MemoryOrderRepository struct {
  14: 	lock  *sync.RWMutex
  15: 	store []*domain.Order
  16: }
  17: 
  18: func NewMemoryOrderRepository() *MemoryOrderRepository {
  19: 	return &MemoryOrderRepository{
  20: 		lock:  &sync.RWMutex{},
  21: 		store: make([]*domain.Order, 0),
  22: 	}
  23: }
  24: 
  25: func (m MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
  26: 	m.lock.Lock()
  27: 	defer m.lock.Unlock()
  28: 	newOrder := &domain.Order{
  29: 		ID:          strconv.FormatInt(time.Now().Unix(), 10),
  30: 		CustomerID:  order.CustomerID,
  31: 		Status:      order.Status,
  32: 		PaymentLink: order.PaymentLink,
  33: 		Items:       order.Items,
  34: 	}
  35: 	m.store = append(m.store, newOrder)
  36: 	logrus.WithFields(logrus.Fields{
  37: 		"input_order":        order,
  38: 		"store_after_create": m.store,
  39: 	}).Debug("memory_order_repo_create")
  40: 	return newOrder, nil
  41: }
  42: 
  43: func (m MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*domain.Order, error) {
  44: 	m.lock.RLock()
  45: 	defer m.lock.RUnlock()
  46: 	for _, o := range m.store {
  47: 		if o.ID == id && o.CustomerID == customerID {
  48: 			logrus.Debugf("memory_order_repo_get||found||id=%s||customerID=%s||res=%+v", id, customerID, *o)
  49: 			return o, nil
  50: 		}
  51: 	}
  52: 	return nil, domain.NotFoundError{OrderID: id}
  53: }
  54: 
  55: func (m MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
  56: 	m.lock.Lock()
  57: 	defer m.lock.Unlock()
  58: 	found := false
  59: 	for i, o := range m.store {
  60: 		if o.ID == order.ID && o.CustomerID == order.CustomerID {
  61: 			found = true
  62: 			updatedOrder, err := updateFn(ctx, o)
  63: 			if err != nil {
  64: 				return err
  65: 			}
  66: 			m.store[i] = updatedOrder
  67: 		}
  68: 	}
  69: 	if !found {
  70: 		return domain.NotFoundError{OrderID: order.ID}
  71: 	}
  72: 	return nil
  73: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Executable statement: participates in current implementation logic. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Blank line: separates code blocks for readability. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Block end: closes import or parameter list block. |
| 12 | Blank line: separates code blocks for readability. |
| 13 | Type declaration: defines struct/interface/type alias. |
| 14 | Executable statement: participates in current implementation logic. |
| 15 | Executable statement: participates in current implementation logic. |
| 16 | Block end: closes previous function/branch/type block. |
| 17 | Blank line: separates code blocks for readability. |
| 18 | Function/method declaration: defines executable logic. |
| 19 | Return statement: returns values (and often error) to caller. |
| 20 | Executable statement: participates in current implementation logic. |
| 21 | Executable statement: participates in current implementation logic. |
| 22 | Block end: closes previous function/branch/type block. |
| 23 | Block end: closes previous function/branch/type block. |
| 24 | Blank line: separates code blocks for readability. |
| 25 | Function/method declaration: defines executable logic. |
| 26 | Executable statement: participates in current implementation logic. |
| 27 | Deferred call: executes when function returns. |
| 28 | Short declaration: defines and initializes local variable. |
| 29 | Executable statement: participates in current implementation logic. |
| 30 | Executable statement: participates in current implementation logic. |
| 31 | Executable statement: participates in current implementation logic. |
| 32 | Executable statement: participates in current implementation logic. |
| 33 | Executable statement: participates in current implementation logic. |
| 34 | Block end: closes previous function/branch/type block. |
| 35 | Assignment: sets or updates variable/field value. |
| 36 | Executable statement: participates in current implementation logic. |
| 37 | Executable statement: participates in current implementation logic. |
| 38 | Executable statement: participates in current implementation logic. |
| 39 | Block end: closes previous function/branch/type block. |
| 40 | Return statement: returns values (and often error) to caller. |
| 41 | Block end: closes previous function/branch/type block. |
| 42 | Blank line: separates code blocks for readability. |
| 43 | Function/method declaration: defines executable logic. |
| 44 | Executable statement: participates in current implementation logic. |
| 45 | Deferred call: executes when function returns. |
| 46 | Loop: iterates collection or repeats logic. |
| 47 | Conditional branch: executes following block when condition is true. |
| 48 | Assignment: sets or updates variable/field value. |
| 49 | Return statement: returns values (and often error) to caller. |
| 50 | Block end: closes previous function/branch/type block. |
| 51 | Block end: closes previous function/branch/type block. |
| 52 | Return statement: returns values (and often error) to caller. |
| 53 | Block end: closes previous function/branch/type block. |
| 54 | Blank line: separates code blocks for readability. |
| 55 | Function/method declaration: defines executable logic. |
| 56 | Executable statement: participates in current implementation logic. |
| 57 | Deferred call: executes when function returns. |
| 58 | Short declaration: defines and initializes local variable. |
| 59 | Loop: iterates collection or repeats logic. |
| 60 | Conditional branch: executes following block when condition is true. |
| 61 | Assignment: sets or updates variable/field value. |
| 62 | Short declaration: defines and initializes local variable. |
| 63 | Conditional branch: executes following block when condition is true. |
| 64 | Return statement: returns values (and often error) to caller. |
| 65 | Block end: closes previous function/branch/type block. |
| 66 | Assignment: sets or updates variable/field value. |
| 67 | Block end: closes previous function/branch/type block. |
| 68 | Block end: closes previous function/branch/type block. |
| 69 | Conditional branch: executes following block when condition is true. |
| 70 | Return statement: returns values (and often error) to caller. |
| 71 | Block end: closes previous function/branch/type block. |
| 72 | Return statement: returns values (and often error) to caller. |
| 73 | Block end: closes previous function/branch/type block. |

### File: internal/order/domain/order/order.go

```go
   1: package order
   2: 
   3: import "github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   4: 
   5: type Order struct {
   6: 	ID          string
   7: 	CustomerID  string
   8: 	Status      string
   9: 	PaymentLink string
  10: 	Items       []*orderpb.Item
  11: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Single import: brings one dependency into scope. |
| 4 | Blank line: separates code blocks for readability. |
| 5 | Type declaration: defines struct/interface/type alias. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Block end: closes previous function/branch/type block. |

### File: internal/order/domain/order/repository.go

```go
   1: package order
   2: 
   3: import (
   4: 	"context"
   5: 	"fmt"
   6: )
   7: 
   8: type Repository interface {
   9: 	Create(context.Context, *Order) (*Order, error)
  10: 	Get(ctx context.Context, id, customerID string) (*Order, error)
  11: 	Update(
  12: 		ctx context.Context,
  13: 		o *Order,
  14: 		updateFn func(context.Context, *Order) (*Order, error),
  15: 	) error
  16: }
  17: 
  18: type NotFoundError struct {
  19: 	OrderID string
  20: }
  21: 
  22: func (e NotFoundError) Error() string {
  23: 	return fmt.Sprintf("order '%s' not found", e.OrderID)
  24: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Executable statement: participates in current implementation logic. |
| 6 | Block end: closes import or parameter list block. |
| 7 | Blank line: separates code blocks for readability. |
| 8 | Type declaration: defines struct/interface/type alias. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Executable statement: participates in current implementation logic. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Executable statement: participates in current implementation logic. |
| 14 | Executable statement: participates in current implementation logic. |
| 15 | Executable statement: participates in current implementation logic. |
| 16 | Block end: closes previous function/branch/type block. |
| 17 | Blank line: separates code blocks for readability. |
| 18 | Type declaration: defines struct/interface/type alias. |
| 19 | Executable statement: participates in current implementation logic. |
| 20 | Block end: closes previous function/branch/type block. |
| 21 | Blank line: separates code blocks for readability. |
| 22 | Function/method declaration: defines executable logic. |
| 23 | Return statement: returns values (and often error) to caller. |
| 24 | Block end: closes previous function/branch/type block. |

### File: internal/stock/adapters/stock_inmem_repository.go

```go
   1: package adapters
   2: 
   3: import (
   4: 	"context"
   5: 	"sync"
   6: 
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   8: 	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
   9: )
  10: 
  11: type MemoryStockRepository struct {
  12: 	lock  *sync.RWMutex
  13: 	store map[string]*orderpb.Item
  14: }
  15: 
  16: var stub = map[string]*orderpb.Item{
  17: 	"item_id": {
  18: 		ID:       "foo_item",
  19: 		Name:     "stub item",
  20: 		Quantity: 10000,
  21: 		PriceID:  "stub_item_price_id",
  22: 	},
  23: }
  24: 
  25: func NewMemoryOrderRepository() *MemoryStockRepository {
  26: 	return &MemoryStockRepository{
  27: 		lock:  &sync.RWMutex{},
  28: 		store: stub,
  29: 	}
  30: }
  31: 
  32: func (m MemoryStockRepository) GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error) {
  33: 	m.lock.RLock()
  34: 	defer m.lock.RUnlock()
  35: 	var (
  36: 		res     []*orderpb.Item
  37: 		missing []string
  38: 	)
  39: 	for _, id := range ids {
  40: 		if item, exist := m.store[id]; exist {
  41: 			res = append(res, item)
  42: 		} else {
  43: 			missing = append(missing, id)
  44: 		}
  45: 	}
  46: 	if len(res) == len(ids) {
  47: 		return res, nil
  48: 	}
  49: 	return res, domain.NotFoundError{Missing: missing}
  50: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Executable statement: participates in current implementation logic. |
| 6 | Blank line: separates code blocks for readability. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Block end: closes import or parameter list block. |
| 10 | Blank line: separates code blocks for readability. |
| 11 | Type declaration: defines struct/interface/type alias. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Executable statement: participates in current implementation logic. |
| 14 | Block end: closes previous function/branch/type block. |
| 15 | Blank line: separates code blocks for readability. |
| 16 | Assignment: sets or updates variable/field value. |
| 17 | Executable statement: participates in current implementation logic. |
| 18 | Executable statement: participates in current implementation logic. |
| 19 | Executable statement: participates in current implementation logic. |
| 20 | Executable statement: participates in current implementation logic. |
| 21 | Executable statement: participates in current implementation logic. |
| 22 | Block end: closes previous function/branch/type block. |
| 23 | Block end: closes previous function/branch/type block. |
| 24 | Blank line: separates code blocks for readability. |
| 25 | Function/method declaration: defines executable logic. |
| 26 | Return statement: returns values (and often error) to caller. |
| 27 | Executable statement: participates in current implementation logic. |
| 28 | Executable statement: participates in current implementation logic. |
| 29 | Block end: closes previous function/branch/type block. |
| 30 | Block end: closes previous function/branch/type block. |
| 31 | Blank line: separates code blocks for readability. |
| 32 | Function/method declaration: defines executable logic. |
| 33 | Executable statement: participates in current implementation logic. |
| 34 | Deferred call: executes when function returns. |
| 35 | Executable statement: participates in current implementation logic. |
| 36 | Executable statement: participates in current implementation logic. |
| 37 | Executable statement: participates in current implementation logic. |
| 38 | Block end: closes import or parameter list block. |
| 39 | Loop: iterates collection or repeats logic. |
| 40 | Conditional branch: executes following block when condition is true. |
| 41 | Assignment: sets or updates variable/field value. |
| 42 | Block end: closes previous function/branch/type block. |
| 43 | Assignment: sets or updates variable/field value. |
| 44 | Block end: closes previous function/branch/type block. |
| 45 | Block end: closes previous function/branch/type block. |
| 46 | Conditional branch: executes following block when condition is true. |
| 47 | Return statement: returns values (and often error) to caller. |
| 48 | Block end: closes previous function/branch/type block. |
| 49 | Return statement: returns values (and often error) to caller. |
| 50 | Block end: closes previous function/branch/type block. |

### File: internal/stock/domain/stock/repository.go

```go
   1: package stock
   2: 
   3: import (
   4: 	"context"
   5: 	"fmt"
   6: 	"strings"
   7: 
   8: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   9: )
  10: 
  11: type Repository interface {
  12: 	GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error)
  13: }
  14: 
  15: type NotFoundError struct {
  16: 	Missing []string
  17: }
  18: 
  19: func (e NotFoundError) Error() string {
  20: 	return fmt.Sprintf("these items not found in stock: %s", strings.Join(e.Missing, ","))
  21: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Executable statement: participates in current implementation logic. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Blank line: separates code blocks for readability. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Block end: closes import or parameter list block. |
| 10 | Blank line: separates code blocks for readability. |
| 11 | Type declaration: defines struct/interface/type alias. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Block end: closes previous function/branch/type block. |
| 14 | Blank line: separates code blocks for readability. |
| 15 | Type declaration: defines struct/interface/type alias. |
| 16 | Executable statement: participates in current implementation logic. |
| 17 | Block end: closes previous function/branch/type block. |
| 18 | Blank line: separates code blocks for readability. |
| 19 | Function/method declaration: defines executable logic. |
| 20 | Return statement: returns values (and often error) to caller. |
| 21 | Block end: closes previous function/branch/type block. |

## Commit 7: [7881a0b] first Query

### File: internal/common/decorator/logging.go

```go
   1: package decorator
   2: 
   3: import (
   4: 	"context"
   5: 	"fmt"
   6: 	"strings"
   7: 
   8: 	"github.com/sirupsen/logrus"
   9: )
  10: 
  11: type queryLoggingDecorator[C, R any] struct {
  12: 	logger *logrus.Entry
  13: 	base   QueryHandler[C, R]
  14: }
  15: 
  16: func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
  17: 	logger := q.logger.WithFields(logrus.Fields{
  18: 		"query":      generateActionName(cmd),
  19: 		"query_body": fmt.Sprintf("%#v", cmd),
  20: 	})
  21: 	logger.Debug("Executing query")
  22: 	defer func() {
  23: 		if err == nil {
  24: 			logger.Info("Query execute successfully")
  25: 		} else {
  26: 			logger.Error("Failed to execute query", err)
  27: 		}
  28: 	}()
  29: 	return q.base.Handle(ctx, cmd)
  30: }
  31: 
  32: func generateActionName(cmd any) string {
  33: 	return strings.Split(fmt.Sprintf("%T", cmd), ".")[1]
  34: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Executable statement: participates in current implementation logic. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Blank line: separates code blocks for readability. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Block end: closes import or parameter list block. |
| 10 | Blank line: separates code blocks for readability. |
| 11 | Type declaration: defines struct/interface/type alias. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Executable statement: participates in current implementation logic. |
| 14 | Block end: closes previous function/branch/type block. |
| 15 | Blank line: separates code blocks for readability. |
| 16 | Function/method declaration: defines executable logic. |
| 17 | Short declaration: defines and initializes local variable. |
| 18 | Executable statement: participates in current implementation logic. |
| 19 | Executable statement: participates in current implementation logic. |
| 20 | Block end: closes previous function/branch/type block. |
| 21 | Executable statement: participates in current implementation logic. |
| 22 | Deferred call: executes when function returns. |
| 23 | Conditional branch: executes following block when condition is true. |
| 24 | Executable statement: participates in current implementation logic. |
| 25 | Block end: closes previous function/branch/type block. |
| 26 | Executable statement: participates in current implementation logic. |
| 27 | Block end: closes previous function/branch/type block. |
| 28 | Block end: closes previous function/branch/type block. |
| 29 | Return statement: returns values (and often error) to caller. |
| 30 | Block end: closes previous function/branch/type block. |
| 31 | Blank line: separates code blocks for readability. |
| 32 | Function/method declaration: defines executable logic. |
| 33 | Return statement: returns values (and often error) to caller. |
| 34 | Block end: closes previous function/branch/type block. |

### File: internal/common/decorator/metrics.go

```go
   1: package decorator
   2: 
   3: import (
   4: 	"context"
   5: 	"fmt"
   6: 	"strings"
   7: 	"time"
   8: )
   9: 
  10: type MetricsClient interface {
  11: 	Inc(key string, value int)
  12: }
  13: 
  14: type queryMetricsDecorator[C, R any] struct {
  15: 	base   QueryHandler[C, R]
  16: 	client MetricsClient
  17: }
  18: 
  19: func (q queryMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
  20: 	start := time.Now()
  21: 	actionName := strings.ToLower(generateActionName(cmd))
  22: 	defer func() {
  23: 		end := time.Since(start)
  24: 		q.client.Inc(fmt.Sprintf("querys.%s.duration", actionName), int(end.Seconds()))
  25: 		if err == nil {
  26: 			q.client.Inc(fmt.Sprintf("querys.%s.success", actionName), 1)
  27: 		} else {
  28: 			q.client.Inc(fmt.Sprintf("querys.%s.failure", actionName), 1)
  29: 		}
  30: 	}()
  31: 	return q.base.Handle(ctx, cmd)
  32: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Executable statement: participates in current implementation logic. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Block end: closes import or parameter list block. |
| 9 | Blank line: separates code blocks for readability. |
| 10 | Type declaration: defines struct/interface/type alias. |
| 11 | Executable statement: participates in current implementation logic. |
| 12 | Block end: closes previous function/branch/type block. |
| 13 | Blank line: separates code blocks for readability. |
| 14 | Type declaration: defines struct/interface/type alias. |
| 15 | Executable statement: participates in current implementation logic. |
| 16 | Executable statement: participates in current implementation logic. |
| 17 | Block end: closes previous function/branch/type block. |
| 18 | Blank line: separates code blocks for readability. |
| 19 | Function/method declaration: defines executable logic. |
| 20 | Short declaration: defines and initializes local variable. |
| 21 | Short declaration: defines and initializes local variable. |
| 22 | Deferred call: executes when function returns. |
| 23 | Short declaration: defines and initializes local variable. |
| 24 | Executable statement: participates in current implementation logic. |
| 25 | Conditional branch: executes following block when condition is true. |
| 26 | Executable statement: participates in current implementation logic. |
| 27 | Block end: closes previous function/branch/type block. |
| 28 | Executable statement: participates in current implementation logic. |
| 29 | Block end: closes previous function/branch/type block. |
| 30 | Block end: closes previous function/branch/type block. |
| 31 | Return statement: returns values (and often error) to caller. |
| 32 | Block end: closes previous function/branch/type block. |

### File: internal/common/decorator/query.go

```go
   1: package decorator
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/sirupsen/logrus"
   7: )
   8: 
   9: // QueryHandler defines a generic type that receives a Query Q,
  10: // and returns a result R
  11: type QueryHandler[Q, R any] interface {
  12: 	Handle(ctx context.Context, query Q) (R, error)
  13: }
  14: 
  15: func ApplyQueryDecorators[H, R any](handler QueryHandler[H, R], logger *logrus.Entry, metricsClient MetricsClient) QueryHandler[H, R] {
  16: 	return queryLoggingDecorator[H, R]{
  17: 		logger: logger,
  18: 		base: queryMetricsDecorator[H, R]{
  19: 			base:   handler,
  20: 			client: metricsClient,
  21: 		},
  22: 	}
  23: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Block end: closes import or parameter list block. |
| 8 | Blank line: separates code blocks for readability. |
| 9 | Comment: explains intent, todo, or context. |
| 10 | Comment: explains intent, todo, or context. |
| 11 | Type declaration: defines struct/interface/type alias. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Block end: closes previous function/branch/type block. |
| 14 | Blank line: separates code blocks for readability. |
| 15 | Function/method declaration: defines executable logic. |
| 16 | Return statement: returns values (and often error) to caller. |
| 17 | Executable statement: participates in current implementation logic. |
| 18 | Executable statement: participates in current implementation logic. |
| 19 | Executable statement: participates in current implementation logic. |
| 20 | Executable statement: participates in current implementation logic. |
| 21 | Block end: closes previous function/branch/type block. |
| 22 | Block end: closes previous function/branch/type block. |
| 23 | Block end: closes previous function/branch/type block. |

### File: internal/common/metrics/todo_metrics.go

```go
   1: package metrics
   2: 
   3: type TodoMetrics struct{}
   4: 
   5: func (t TodoMetrics) Inc(_ string, _ int) {
   6: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Type declaration: defines struct/interface/type alias. |
| 4 | Blank line: separates code blocks for readability. |
| 5 | Function/method declaration: defines executable logic. |
| 6 | Block end: closes previous function/branch/type block. |

### File: internal/order/adapters/order_inmem_repository.go

```go
   1: package adapters
   2: 
   3: import (
   4: 	"context"
   5: 	"strconv"
   6: 	"sync"
   7: 	"time"
   8: 
   9: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
  10: 	"github.com/sirupsen/logrus"
  11: )
  12: 
  13: type MemoryOrderRepository struct {
  14: 	lock  *sync.RWMutex
  15: 	store []*domain.Order
  16: }
  17: 
  18: func NewMemoryOrderRepository() *MemoryOrderRepository {
  19: 	s := make([]*domain.Order, 0)
  20: 	s = append(s, &domain.Order{
  21: 		ID:          "fake-ID",
  22: 		CustomerID:  "fake-customer-id",
  23: 		Status:      "fake-status",
  24: 		PaymentLink: "fake-payment-link",
  25: 		Items:       nil,
  26: 	})
  27: 	return &MemoryOrderRepository{
  28: 		lock:  &sync.RWMutex{},
  29: 		store: s,
  30: 	}
  31: }
  32: 
  33: func (m MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
  34: 	m.lock.Lock()
  35: 	defer m.lock.Unlock()
  36: 	newOrder := &domain.Order{
  37: 		ID:          strconv.FormatInt(time.Now().Unix(), 10),
  38: 		CustomerID:  order.CustomerID,
  39: 		Status:      order.Status,
  40: 		PaymentLink: order.PaymentLink,
  41: 		Items:       order.Items,
  42: 	}
  43: 	m.store = append(m.store, newOrder)
  44: 	logrus.WithFields(logrus.Fields{
  45: 		"input_order":        order,
  46: 		"store_after_create": m.store,
  47: 	}).Debug("memory_order_repo_create")
  48: 	return newOrder, nil
  49: }
  50: 
  51: func (m MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*domain.Order, error) {
  52: 	m.lock.RLock()
  53: 	defer m.lock.RUnlock()
  54: 	for _, o := range m.store {
  55: 		if o.ID == id && o.CustomerID == customerID {
  56: 			logrus.Debugf("memory_order_repo_get||found||id=%s||customerID=%s||res=%+v", id, customerID, *o)
  57: 			return o, nil
  58: 		}
  59: 	}
  60: 	return nil, domain.NotFoundError{OrderID: id}
  61: }
  62: 
  63: func (m MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
  64: 	m.lock.Lock()
  65: 	defer m.lock.Unlock()
  66: 	found := false
  67: 	for i, o := range m.store {
  68: 		if o.ID == order.ID && o.CustomerID == order.CustomerID {
  69: 			found = true
  70: 			updatedOrder, err := updateFn(ctx, o)
  71: 			if err != nil {
  72: 				return err
  73: 			}
  74: 			m.store[i] = updatedOrder
  75: 		}
  76: 	}
  77: 	if !found {
  78: 		return domain.NotFoundError{OrderID: order.ID}
  79: 	}
  80: 	return nil
  81: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Executable statement: participates in current implementation logic. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Blank line: separates code blocks for readability. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Block end: closes import or parameter list block. |
| 12 | Blank line: separates code blocks for readability. |
| 13 | Type declaration: defines struct/interface/type alias. |
| 14 | Executable statement: participates in current implementation logic. |
| 15 | Executable statement: participates in current implementation logic. |
| 16 | Block end: closes previous function/branch/type block. |
| 17 | Blank line: separates code blocks for readability. |
| 18 | Function/method declaration: defines executable logic. |
| 19 | Short declaration: defines and initializes local variable. |
| 20 | Assignment: sets or updates variable/field value. |
| 21 | Executable statement: participates in current implementation logic. |
| 22 | Executable statement: participates in current implementation logic. |
| 23 | Executable statement: participates in current implementation logic. |
| 24 | Executable statement: participates in current implementation logic. |
| 25 | Executable statement: participates in current implementation logic. |
| 26 | Block end: closes previous function/branch/type block. |
| 27 | Return statement: returns values (and often error) to caller. |
| 28 | Executable statement: participates in current implementation logic. |
| 29 | Executable statement: participates in current implementation logic. |
| 30 | Block end: closes previous function/branch/type block. |
| 31 | Block end: closes previous function/branch/type block. |
| 32 | Blank line: separates code blocks for readability. |
| 33 | Function/method declaration: defines executable logic. |
| 34 | Executable statement: participates in current implementation logic. |
| 35 | Deferred call: executes when function returns. |
| 36 | Short declaration: defines and initializes local variable. |
| 37 | Executable statement: participates in current implementation logic. |
| 38 | Executable statement: participates in current implementation logic. |
| 39 | Executable statement: participates in current implementation logic. |
| 40 | Executable statement: participates in current implementation logic. |
| 41 | Executable statement: participates in current implementation logic. |
| 42 | Block end: closes previous function/branch/type block. |
| 43 | Assignment: sets or updates variable/field value. |
| 44 | Executable statement: participates in current implementation logic. |
| 45 | Executable statement: participates in current implementation logic. |
| 46 | Executable statement: participates in current implementation logic. |
| 47 | Block end: closes previous function/branch/type block. |
| 48 | Return statement: returns values (and often error) to caller. |
| 49 | Block end: closes previous function/branch/type block. |
| 50 | Blank line: separates code blocks for readability. |
| 51 | Function/method declaration: defines executable logic. |
| 52 | Executable statement: participates in current implementation logic. |
| 53 | Deferred call: executes when function returns. |
| 54 | Loop: iterates collection or repeats logic. |
| 55 | Conditional branch: executes following block when condition is true. |
| 56 | Assignment: sets or updates variable/field value. |
| 57 | Return statement: returns values (and often error) to caller. |
| 58 | Block end: closes previous function/branch/type block. |
| 59 | Block end: closes previous function/branch/type block. |
| 60 | Return statement: returns values (and often error) to caller. |
| 61 | Block end: closes previous function/branch/type block. |
| 62 | Blank line: separates code blocks for readability. |
| 63 | Function/method declaration: defines executable logic. |
| 64 | Executable statement: participates in current implementation logic. |
| 65 | Deferred call: executes when function returns. |
| 66 | Short declaration: defines and initializes local variable. |
| 67 | Loop: iterates collection or repeats logic. |
| 68 | Conditional branch: executes following block when condition is true. |
| 69 | Assignment: sets or updates variable/field value. |
| 70 | Short declaration: defines and initializes local variable. |
| 71 | Conditional branch: executes following block when condition is true. |
| 72 | Return statement: returns values (and often error) to caller. |
| 73 | Block end: closes previous function/branch/type block. |
| 74 | Assignment: sets or updates variable/field value. |
| 75 | Block end: closes previous function/branch/type block. |
| 76 | Block end: closes previous function/branch/type block. |
| 77 | Conditional branch: executes following block when condition is true. |
| 78 | Return statement: returns values (and often error) to caller. |
| 79 | Block end: closes previous function/branch/type block. |
| 80 | Return statement: returns values (and often error) to caller. |
| 81 | Block end: closes previous function/branch/type block. |

### File: internal/order/app/app.go

```go
   1: package app
   2: 
   3: import "github.com/ghost-yu/go_shop_second/order/app/query"
   4: 
   5: type Application struct {
   6: 	Commands Commands
   7: 	Queries  Queries
   8: }
   9: 
  10: type Commands struct{}
  11: 
  12: type Queries struct {
  13: 	GetCustomerOrder query.GetCustomerOrderHandler
  14: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Single import: brings one dependency into scope. |
| 4 | Blank line: separates code blocks for readability. |
| 5 | Type declaration: defines struct/interface/type alias. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Block end: closes previous function/branch/type block. |
| 9 | Blank line: separates code blocks for readability. |
| 10 | Type declaration: defines struct/interface/type alias. |
| 11 | Blank line: separates code blocks for readability. |
| 12 | Type declaration: defines struct/interface/type alias. |
| 13 | Executable statement: participates in current implementation logic. |
| 14 | Block end: closes previous function/branch/type block. |

### File: internal/order/app/query/get_customer_order.go

```go
   1: package query
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/decorator"
   7: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
   8: 	"github.com/sirupsen/logrus"
   9: )
  10: 
  11: type GetCustomerOrder struct {
  12: 	CustomerID string
  13: 	OrderID    string
  14: }
  15: 
  16: type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Order]
  17: 
  18: type getCustomerOrderHandler struct {
  19: 	orderRepo domain.Repository
  20: }
  21: 
  22: func NewGetCustomerOrderHandler(
  23: 	orderRepo domain.Repository,
  24: 	logger *logrus.Entry,
  25: 	metricClient decorator.MetricsClient,
  26: ) GetCustomerOrderHandler {
  27: 	if orderRepo == nil {
  28: 		panic("nil orderRepo")
  29: 	}
  30: 	return decorator.ApplyQueryDecorators[GetCustomerOrder, *domain.Order](
  31: 		getCustomerOrderHandler{orderRepo: orderRepo},
  32: 		logger,
  33: 		metricClient,
  34: 	)
  35: }
  36: 
  37: func (g getCustomerOrderHandler) Handle(ctx context.Context, query GetCustomerOrder) (*domain.Order, error) {
  38: 	o, err := g.orderRepo.Get(ctx, query.OrderID, query.CustomerID)
  39: 	if err != nil {
  40: 		return nil, err
  41: 	}
  42: 	return o, nil
  43: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Block end: closes import or parameter list block. |
| 10 | Blank line: separates code blocks for readability. |
| 11 | Type declaration: defines struct/interface/type alias. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Executable statement: participates in current implementation logic. |
| 14 | Block end: closes previous function/branch/type block. |
| 15 | Blank line: separates code blocks for readability. |
| 16 | Type declaration: defines struct/interface/type alias. |
| 17 | Blank line: separates code blocks for readability. |
| 18 | Type declaration: defines struct/interface/type alias. |
| 19 | Executable statement: participates in current implementation logic. |
| 20 | Block end: closes previous function/branch/type block. |
| 21 | Blank line: separates code blocks for readability. |
| 22 | Function/method declaration: defines executable logic. |
| 23 | Executable statement: participates in current implementation logic. |
| 24 | Executable statement: participates in current implementation logic. |
| 25 | Executable statement: participates in current implementation logic. |
| 26 | Executable statement: participates in current implementation logic. |
| 27 | Conditional branch: executes following block when condition is true. |
| 28 | Panic: aborts normal control flow with unrecoverable error. |
| 29 | Block end: closes previous function/branch/type block. |
| 30 | Return statement: returns values (and often error) to caller. |
| 31 | Executable statement: participates in current implementation logic. |
| 32 | Executable statement: participates in current implementation logic. |
| 33 | Executable statement: participates in current implementation logic. |
| 34 | Block end: closes import or parameter list block. |
| 35 | Block end: closes previous function/branch/type block. |
| 36 | Blank line: separates code blocks for readability. |
| 37 | Function/method declaration: defines executable logic. |
| 38 | Short declaration: defines and initializes local variable. |
| 39 | Conditional branch: executes following block when condition is true. |
| 40 | Return statement: returns values (and often error) to caller. |
| 41 | Block end: closes previous function/branch/type block. |
| 42 | Return statement: returns values (and often error) to caller. |
| 43 | Block end: closes previous function/branch/type block. |

### File: internal/order/http.go

```go
   1: package main
   2: 
   3: import (
   4: 	"net/http"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/order/app"
   7: 	"github.com/ghost-yu/go_shop_second/order/app/query"
   8: 	"github.com/gin-gonic/gin"
   9: )
  10: 
  11: type HTTPServer struct {
  12: 	app app.Application
  13: }
  14: 
  15: func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
  16: 	//TODO implement me
  17: 	panic("implement me")
  18: }
  19: 
  20: func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
  21: 	o, err := H.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
  22: 		OrderID:    "fake-ID",
  23: 		CustomerID: "fake-customer-id",
  24: 	})
  25: 	if err != nil {
  26: 		c.JSON(http.StatusOK, gin.H{"error": err})
  27: 		return
  28: 	}
  29: 	c.JSON(http.StatusOK, gin.H{"message": "success", "data": o})
  30: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Block end: closes import or parameter list block. |
| 10 | Blank line: separates code blocks for readability. |
| 11 | Type declaration: defines struct/interface/type alias. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Block end: closes previous function/branch/type block. |
| 14 | Blank line: separates code blocks for readability. |
| 15 | Function/method declaration: defines executable logic. |
| 16 | Comment: explains intent, todo, or context. |
| 17 | Panic: aborts normal control flow with unrecoverable error. |
| 18 | Block end: closes previous function/branch/type block. |
| 19 | Blank line: separates code blocks for readability. |
| 20 | Function/method declaration: defines executable logic. |
| 21 | Short declaration: defines and initializes local variable. |
| 22 | Executable statement: participates in current implementation logic. |
| 23 | Executable statement: participates in current implementation logic. |
| 24 | Block end: closes previous function/branch/type block. |
| 25 | Conditional branch: executes following block when condition is true. |
| 26 | Executable statement: participates in current implementation logic. |
| 27 | Return statement: returns values (and often error) to caller. |
| 28 | Block end: closes previous function/branch/type block. |
| 29 | Executable statement: participates in current implementation logic. |
| 30 | Block end: closes previous function/branch/type block. |

### File: internal/order/service/application.go

```go
   1: package service
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/metrics"
   7: 	"github.com/ghost-yu/go_shop_second/order/adapters"
   8: 	"github.com/ghost-yu/go_shop_second/order/app"
   9: 	"github.com/ghost-yu/go_shop_second/order/app/query"
  10: 	"github.com/sirupsen/logrus"
  11: )
  12: 
  13: func NewApplication(ctx context.Context) app.Application {
  14: 	orderRepo := adapters.NewMemoryOrderRepository()
  15: 	logger := logrus.NewEntry(logrus.StandardLogger())
  16: 	metricClient := metrics.TodoMetrics{}
  17: 	return app.Application{
  18: 		Commands: app.Commands{},
  19: 		Queries: app.Queries{
  20: 			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
  21: 		},
  22: 	}
  23: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Block end: closes import or parameter list block. |
| 12 | Blank line: separates code blocks for readability. |
| 13 | Function/method declaration: defines executable logic. |
| 14 | Short declaration: defines and initializes local variable. |
| 15 | Short declaration: defines and initializes local variable. |
| 16 | Short declaration: defines and initializes local variable. |
| 17 | Return statement: returns values (and often error) to caller. |
| 18 | Executable statement: participates in current implementation logic. |
| 19 | Executable statement: participates in current implementation logic. |
| 20 | Executable statement: participates in current implementation logic. |
| 21 | Block end: closes previous function/branch/type block. |
| 22 | Block end: closes previous function/branch/type block. |
| 23 | Block end: closes previous function/branch/type block. |

### File: internal/stock/adapters/stock_inmem_repository.go

```go
   1: package adapters
   2: 
   3: import (
   4: 	"context"
   5: 	"sync"
   6: 
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   8: 	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
   9: )
  10: 
  11: type MemoryStockRepository struct {
  12: 	lock  *sync.RWMutex
  13: 	store map[string]*orderpb.Item
  14: }
  15: 
  16: var stub = map[string]*orderpb.Item{
  17: 	"item_id": {
  18: 		ID:       "foo_item",
  19: 		Name:     "stub item",
  20: 		Quantity: 10000,
  21: 		PriceID:  "stub_item_price_id",
  22: 	},
  23: }
  24: 
  25: func NewMemoryStockRepository() *MemoryStockRepository {
  26: 	return &MemoryStockRepository{
  27: 		lock:  &sync.RWMutex{},
  28: 		store: stub,
  29: 	}
  30: }
  31: 
  32: func (m MemoryStockRepository) GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error) {
  33: 	m.lock.RLock()
  34: 	defer m.lock.RUnlock()
  35: 	var (
  36: 		res     []*orderpb.Item
  37: 		missing []string
  38: 	)
  39: 	for _, id := range ids {
  40: 		if item, exist := m.store[id]; exist {
  41: 			res = append(res, item)
  42: 		} else {
  43: 			missing = append(missing, id)
  44: 		}
  45: 	}
  46: 	if len(res) == len(ids) {
  47: 		return res, nil
  48: 	}
  49: 	return res, domain.NotFoundError{Missing: missing}
  50: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Executable statement: participates in current implementation logic. |
| 6 | Blank line: separates code blocks for readability. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Block end: closes import or parameter list block. |
| 10 | Blank line: separates code blocks for readability. |
| 11 | Type declaration: defines struct/interface/type alias. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Executable statement: participates in current implementation logic. |
| 14 | Block end: closes previous function/branch/type block. |
| 15 | Blank line: separates code blocks for readability. |
| 16 | Assignment: sets or updates variable/field value. |
| 17 | Executable statement: participates in current implementation logic. |
| 18 | Executable statement: participates in current implementation logic. |
| 19 | Executable statement: participates in current implementation logic. |
| 20 | Executable statement: participates in current implementation logic. |
| 21 | Executable statement: participates in current implementation logic. |
| 22 | Block end: closes previous function/branch/type block. |
| 23 | Block end: closes previous function/branch/type block. |
| 24 | Blank line: separates code blocks for readability. |
| 25 | Function/method declaration: defines executable logic. |
| 26 | Return statement: returns values (and often error) to caller. |
| 27 | Executable statement: participates in current implementation logic. |
| 28 | Executable statement: participates in current implementation logic. |
| 29 | Block end: closes previous function/branch/type block. |
| 30 | Block end: closes previous function/branch/type block. |
| 31 | Blank line: separates code blocks for readability. |
| 32 | Function/method declaration: defines executable logic. |
| 33 | Executable statement: participates in current implementation logic. |
| 34 | Deferred call: executes when function returns. |
| 35 | Executable statement: participates in current implementation logic. |
| 36 | Executable statement: participates in current implementation logic. |
| 37 | Executable statement: participates in current implementation logic. |
| 38 | Block end: closes import or parameter list block. |
| 39 | Loop: iterates collection or repeats logic. |
| 40 | Conditional branch: executes following block when condition is true. |
| 41 | Assignment: sets or updates variable/field value. |
| 42 | Block end: closes previous function/branch/type block. |
| 43 | Assignment: sets or updates variable/field value. |
| 44 | Block end: closes previous function/branch/type block. |
| 45 | Block end: closes previous function/branch/type block. |
| 46 | Conditional branch: executes following block when condition is true. |
| 47 | Return statement: returns values (and often error) to caller. |
| 48 | Block end: closes previous function/branch/type block. |
| 49 | Return statement: returns values (and often error) to caller. |
| 50 | Block end: closes previous function/branch/type block. |

## Commit 8: [606345b] order create and update

### File: internal/common/decorator/command.go

```go
   1: package decorator
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/sirupsen/logrus"
   7: )
   8: 
   9: type CommandHandler[C, R any] interface {
  10: 	Handle(ctx context.Context, cmd C) (R, error)
  11: }
  12: 
  13: func ApplyCommandDecorators[C, R any](handler CommandHandler[C, R], logger *logrus.Entry, metricsClient MetricsClient) CommandHandler[C, R] {
  14: 	return queryLoggingDecorator[C, R]{
  15: 		logger: logger,
  16: 		base: queryMetricsDecorator[C, R]{
  17: 			base:   handler,
  18: 			client: metricsClient,
  19: 		},
  20: 	}
  21: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Block end: closes import or parameter list block. |
| 8 | Blank line: separates code blocks for readability. |
| 9 | Type declaration: defines struct/interface/type alias. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Block end: closes previous function/branch/type block. |
| 12 | Blank line: separates code blocks for readability. |
| 13 | Function/method declaration: defines executable logic. |
| 14 | Return statement: returns values (and often error) to caller. |
| 15 | Executable statement: participates in current implementation logic. |
| 16 | Executable statement: participates in current implementation logic. |
| 17 | Executable statement: participates in current implementation logic. |
| 18 | Executable statement: participates in current implementation logic. |
| 19 | Block end: closes previous function/branch/type block. |
| 20 | Block end: closes previous function/branch/type block. |
| 21 | Block end: closes previous function/branch/type block. |

### File: internal/order/adapters/order_inmem_repository.go

```go
   1: package adapters
   2: 
   3: import (
   4: 	"context"
   5: 	"strconv"
   6: 	"sync"
   7: 	"time"
   8: 
   9: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
  10: 	"github.com/sirupsen/logrus"
  11: )
  12: 
  13: type MemoryOrderRepository struct {
  14: 	lock  *sync.RWMutex
  15: 	store []*domain.Order
  16: }
  17: 
  18: func NewMemoryOrderRepository() *MemoryOrderRepository {
  19: 	s := make([]*domain.Order, 0)
  20: 	s = append(s, &domain.Order{
  21: 		ID:          "fake-ID",
  22: 		CustomerID:  "fake-customer-id",
  23: 		Status:      "fake-status",
  24: 		PaymentLink: "fake-payment-link",
  25: 		Items:       nil,
  26: 	})
  27: 	return &MemoryOrderRepository{
  28: 		lock:  &sync.RWMutex{},
  29: 		store: s,
  30: 	}
  31: }
  32: 
  33: func (m *MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
  34: 	m.lock.Lock()
  35: 	defer m.lock.Unlock()
  36: 	newOrder := &domain.Order{
  37: 		ID:          strconv.FormatInt(time.Now().Unix(), 10),
  38: 		CustomerID:  order.CustomerID,
  39: 		Status:      order.Status,
  40: 		PaymentLink: order.PaymentLink,
  41: 		Items:       order.Items,
  42: 	}
  43: 	m.store = append(m.store, newOrder)
  44: 	logrus.WithFields(logrus.Fields{
  45: 		"input_order":        order,
  46: 		"store_after_create": m.store,
  47: 	}).Info("memory_order_repo_create")
  48: 	return newOrder, nil
  49: }
  50: 
  51: func (m *MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*domain.Order, error) {
  52: 	for i, v := range m.store {
  53: 		logrus.Infof("m.store[%d] = %+v", i, v)
  54: 	}
  55: 	m.lock.RLock()
  56: 	defer m.lock.RUnlock()
  57: 	for _, o := range m.store {
  58: 		if o.ID == id && o.CustomerID == customerID {
  59: 			logrus.Infof("memory_order_repo_get||found||id=%s||customerID=%s||res=%+v", id, customerID, *o)
  60: 			return o, nil
  61: 		}
  62: 	}
  63: 	return nil, domain.NotFoundError{OrderID: id}
  64: }
  65: 
  66: func (m *MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
  67: 	m.lock.Lock()
  68: 	defer m.lock.Unlock()
  69: 	found := false
  70: 	for i, o := range m.store {
  71: 		if o.ID == order.ID && o.CustomerID == order.CustomerID {
  72: 			found = true
  73: 			updatedOrder, err := updateFn(ctx, o)
  74: 			if err != nil {
  75: 				return err
  76: 			}
  77: 			m.store[i] = updatedOrder
  78: 		}
  79: 	}
  80: 	if !found {
  81: 		return domain.NotFoundError{OrderID: order.ID}
  82: 	}
  83: 	return nil
  84: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Executable statement: participates in current implementation logic. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Blank line: separates code blocks for readability. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Block end: closes import or parameter list block. |
| 12 | Blank line: separates code blocks for readability. |
| 13 | Type declaration: defines struct/interface/type alias. |
| 14 | Executable statement: participates in current implementation logic. |
| 15 | Executable statement: participates in current implementation logic. |
| 16 | Block end: closes previous function/branch/type block. |
| 17 | Blank line: separates code blocks for readability. |
| 18 | Function/method declaration: defines executable logic. |
| 19 | Short declaration: defines and initializes local variable. |
| 20 | Assignment: sets or updates variable/field value. |
| 21 | Executable statement: participates in current implementation logic. |
| 22 | Executable statement: participates in current implementation logic. |
| 23 | Executable statement: participates in current implementation logic. |
| 24 | Executable statement: participates in current implementation logic. |
| 25 | Executable statement: participates in current implementation logic. |
| 26 | Block end: closes previous function/branch/type block. |
| 27 | Return statement: returns values (and often error) to caller. |
| 28 | Executable statement: participates in current implementation logic. |
| 29 | Executable statement: participates in current implementation logic. |
| 30 | Block end: closes previous function/branch/type block. |
| 31 | Block end: closes previous function/branch/type block. |
| 32 | Blank line: separates code blocks for readability. |
| 33 | Function/method declaration: defines executable logic. |
| 34 | Executable statement: participates in current implementation logic. |
| 35 | Deferred call: executes when function returns. |
| 36 | Short declaration: defines and initializes local variable. |
| 37 | Executable statement: participates in current implementation logic. |
| 38 | Executable statement: participates in current implementation logic. |
| 39 | Executable statement: participates in current implementation logic. |
| 40 | Executable statement: participates in current implementation logic. |
| 41 | Executable statement: participates in current implementation logic. |
| 42 | Block end: closes previous function/branch/type block. |
| 43 | Assignment: sets or updates variable/field value. |
| 44 | Executable statement: participates in current implementation logic. |
| 45 | Executable statement: participates in current implementation logic. |
| 46 | Executable statement: participates in current implementation logic. |
| 47 | Block end: closes previous function/branch/type block. |
| 48 | Return statement: returns values (and often error) to caller. |
| 49 | Block end: closes previous function/branch/type block. |
| 50 | Blank line: separates code blocks for readability. |
| 51 | Function/method declaration: defines executable logic. |
| 52 | Loop: iterates collection or repeats logic. |
| 53 | Assignment: sets or updates variable/field value. |
| 54 | Block end: closes previous function/branch/type block. |
| 55 | Executable statement: participates in current implementation logic. |
| 56 | Deferred call: executes when function returns. |
| 57 | Loop: iterates collection or repeats logic. |
| 58 | Conditional branch: executes following block when condition is true. |
| 59 | Assignment: sets or updates variable/field value. |
| 60 | Return statement: returns values (and often error) to caller. |
| 61 | Block end: closes previous function/branch/type block. |
| 62 | Block end: closes previous function/branch/type block. |
| 63 | Return statement: returns values (and often error) to caller. |
| 64 | Block end: closes previous function/branch/type block. |
| 65 | Blank line: separates code blocks for readability. |
| 66 | Function/method declaration: defines executable logic. |
| 67 | Executable statement: participates in current implementation logic. |
| 68 | Deferred call: executes when function returns. |
| 69 | Short declaration: defines and initializes local variable. |
| 70 | Loop: iterates collection or repeats logic. |
| 71 | Conditional branch: executes following block when condition is true. |
| 72 | Assignment: sets or updates variable/field value. |
| 73 | Short declaration: defines and initializes local variable. |
| 74 | Conditional branch: executes following block when condition is true. |
| 75 | Return statement: returns values (and often error) to caller. |
| 76 | Block end: closes previous function/branch/type block. |
| 77 | Assignment: sets or updates variable/field value. |
| 78 | Block end: closes previous function/branch/type block. |
| 79 | Block end: closes previous function/branch/type block. |
| 80 | Conditional branch: executes following block when condition is true. |
| 81 | Return statement: returns values (and often error) to caller. |
| 82 | Block end: closes previous function/branch/type block. |
| 83 | Return statement: returns values (and often error) to caller. |
| 84 | Block end: closes previous function/branch/type block. |

### File: internal/order/app/app.go

```go
   1: package app
   2: 
   3: import (
   4: 	"github.com/ghost-yu/go_shop_second/order/app/command"
   5: 	"github.com/ghost-yu/go_shop_second/order/app/query"
   6: )
   7: 
   8: type Application struct {
   9: 	Commands Commands
  10: 	Queries  Queries
  11: }
  12: 
  13: type Commands struct {
  14: 	CreateOrder command.CreateOrderHandler
  15: 	UpdateOrder command.UpdateOrderHandler
  16: }
  17: 
  18: type Queries struct {
  19: 	GetCustomerOrder query.GetCustomerOrderHandler
  20: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Executable statement: participates in current implementation logic. |
| 6 | Block end: closes import or parameter list block. |
| 7 | Blank line: separates code blocks for readability. |
| 8 | Type declaration: defines struct/interface/type alias. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Block end: closes previous function/branch/type block. |
| 12 | Blank line: separates code blocks for readability. |
| 13 | Type declaration: defines struct/interface/type alias. |
| 14 | Executable statement: participates in current implementation logic. |
| 15 | Executable statement: participates in current implementation logic. |
| 16 | Block end: closes previous function/branch/type block. |
| 17 | Blank line: separates code blocks for readability. |
| 18 | Type declaration: defines struct/interface/type alias. |
| 19 | Executable statement: participates in current implementation logic. |
| 20 | Block end: closes previous function/branch/type block. |

### File: internal/order/app/command/create_order.go

```go
   1: package command
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/decorator"
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   8: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
   9: 	"github.com/sirupsen/logrus"
  10: )
  11: 
  12: type CreateOrder struct {
  13: 	CustomerID string
  14: 	Items      []*orderpb.ItemWithQuantity
  15: }
  16: 
  17: type CreateOrderResult struct {
  18: 	OrderID string
  19: }
  20: 
  21: type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]
  22: 
  23: type createOrderHandler struct {
  24: 	orderRepo domain.Repository
  25: 	//stockGRPC
  26: }
  27: 
  28: func NewCreateOrderHandler(
  29: 	orderRepo domain.Repository,
  30: 	logger *logrus.Entry,
  31: 	metricClient decorator.MetricsClient,
  32: ) CreateOrderHandler {
  33: 	if orderRepo == nil {
  34: 		panic("nil orderRepo")
  35: 	}
  36: 	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
  37: 		createOrderHandler{orderRepo: orderRepo},
  38: 		logger,
  39: 		metricClient,
  40: 	)
  41: }
  42: 
  43: func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
  44: 	// TODO: call stock grpc to get items.
  45: 	var stockResponse []*orderpb.Item
  46: 	for _, item := range cmd.Items {
  47: 		stockResponse = append(stockResponse, &orderpb.Item{
  48: 			ID:       item.ID,
  49: 			Quantity: item.Quantity,
  50: 		})
  51: 	}
  52: 	o, err := c.orderRepo.Create(ctx, &domain.Order{
  53: 		CustomerID: cmd.CustomerID,
  54: 		Items:      stockResponse,
  55: 	})
  56: 	if err != nil {
  57: 		return nil, err
  58: 	}
  59: 	return &CreateOrderResult{OrderID: o.ID}, nil
  60: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Block end: closes import or parameter list block. |
| 11 | Blank line: separates code blocks for readability. |
| 12 | Type declaration: defines struct/interface/type alias. |
| 13 | Executable statement: participates in current implementation logic. |
| 14 | Executable statement: participates in current implementation logic. |
| 15 | Block end: closes previous function/branch/type block. |
| 16 | Blank line: separates code blocks for readability. |
| 17 | Type declaration: defines struct/interface/type alias. |
| 18 | Executable statement: participates in current implementation logic. |
| 19 | Block end: closes previous function/branch/type block. |
| 20 | Blank line: separates code blocks for readability. |
| 21 | Type declaration: defines struct/interface/type alias. |
| 22 | Blank line: separates code blocks for readability. |
| 23 | Type declaration: defines struct/interface/type alias. |
| 24 | Executable statement: participates in current implementation logic. |
| 25 | Comment: explains intent, todo, or context. |
| 26 | Block end: closes previous function/branch/type block. |
| 27 | Blank line: separates code blocks for readability. |
| 28 | Function/method declaration: defines executable logic. |
| 29 | Executable statement: participates in current implementation logic. |
| 30 | Executable statement: participates in current implementation logic. |
| 31 | Executable statement: participates in current implementation logic. |
| 32 | Executable statement: participates in current implementation logic. |
| 33 | Conditional branch: executes following block when condition is true. |
| 34 | Panic: aborts normal control flow with unrecoverable error. |
| 35 | Block end: closes previous function/branch/type block. |
| 36 | Return statement: returns values (and often error) to caller. |
| 37 | Executable statement: participates in current implementation logic. |
| 38 | Executable statement: participates in current implementation logic. |
| 39 | Executable statement: participates in current implementation logic. |
| 40 | Block end: closes import or parameter list block. |
| 41 | Block end: closes previous function/branch/type block. |
| 42 | Blank line: separates code blocks for readability. |
| 43 | Function/method declaration: defines executable logic. |
| 44 | Comment: explains intent, todo, or context. |
| 45 | Executable statement: participates in current implementation logic. |
| 46 | Loop: iterates collection or repeats logic. |
| 47 | Assignment: sets or updates variable/field value. |
| 48 | Executable statement: participates in current implementation logic. |
| 49 | Executable statement: participates in current implementation logic. |
| 50 | Block end: closes previous function/branch/type block. |
| 51 | Block end: closes previous function/branch/type block. |
| 52 | Short declaration: defines and initializes local variable. |
| 53 | Executable statement: participates in current implementation logic. |
| 54 | Executable statement: participates in current implementation logic. |
| 55 | Block end: closes previous function/branch/type block. |
| 56 | Conditional branch: executes following block when condition is true. |
| 57 | Return statement: returns values (and often error) to caller. |
| 58 | Block end: closes previous function/branch/type block. |
| 59 | Return statement: returns values (and often error) to caller. |
| 60 | Block end: closes previous function/branch/type block. |

### File: internal/order/app/command/update_order.go

```go
   1: package command
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/decorator"
   7: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
   8: 	"github.com/sirupsen/logrus"
   9: )
  10: 
  11: type UpdateOrder struct {
  12: 	Order    *domain.Order
  13: 	UpdateFn func(context.Context, *domain.Order) (*domain.Order, error)
  14: }
  15: 
  16: type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, interface{}]
  17: 
  18: type updateOrderHandler struct {
  19: 	orderRepo domain.Repository
  20: 	//stockGRPC
  21: }
  22: 
  23: func NewUpdateOrderHandler(
  24: 	orderRepo domain.Repository,
  25: 	logger *logrus.Entry,
  26: 	metricClient decorator.MetricsClient,
  27: ) UpdateOrderHandler {
  28: 	if orderRepo == nil {
  29: 		panic("nil orderRepo")
  30: 	}
  31: 	return decorator.ApplyCommandDecorators[UpdateOrder, interface{}](
  32: 		updateOrderHandler{orderRepo: orderRepo},
  33: 		logger,
  34: 		metricClient,
  35: 	)
  36: }
  37: 
  38: func (c updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (interface{}, error) {
  39: 	if cmd.UpdateFn == nil {
  40: 		logrus.Warnf("updateOrderHandler got nil UpdateFn, order=%#v", cmd.Order)
  41: 		cmd.UpdateFn = func(_ context.Context, order *domain.Order) (*domain.Order, error) { return order, nil }
  42: 	}
  43: 	if err := c.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFn); err != nil {
  44: 		return nil, err
  45: 	}
  46: 	return nil, nil
  47: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Block end: closes import or parameter list block. |
| 10 | Blank line: separates code blocks for readability. |
| 11 | Type declaration: defines struct/interface/type alias. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Executable statement: participates in current implementation logic. |
| 14 | Block end: closes previous function/branch/type block. |
| 15 | Blank line: separates code blocks for readability. |
| 16 | Type declaration: defines struct/interface/type alias. |
| 17 | Blank line: separates code blocks for readability. |
| 18 | Type declaration: defines struct/interface/type alias. |
| 19 | Executable statement: participates in current implementation logic. |
| 20 | Comment: explains intent, todo, or context. |
| 21 | Block end: closes previous function/branch/type block. |
| 22 | Blank line: separates code blocks for readability. |
| 23 | Function/method declaration: defines executable logic. |
| 24 | Executable statement: participates in current implementation logic. |
| 25 | Executable statement: participates in current implementation logic. |
| 26 | Executable statement: participates in current implementation logic. |
| 27 | Executable statement: participates in current implementation logic. |
| 28 | Conditional branch: executes following block when condition is true. |
| 29 | Panic: aborts normal control flow with unrecoverable error. |
| 30 | Block end: closes previous function/branch/type block. |
| 31 | Return statement: returns values (and often error) to caller. |
| 32 | Executable statement: participates in current implementation logic. |
| 33 | Executable statement: participates in current implementation logic. |
| 34 | Executable statement: participates in current implementation logic. |
| 35 | Block end: closes import or parameter list block. |
| 36 | Block end: closes previous function/branch/type block. |
| 37 | Blank line: separates code blocks for readability. |
| 38 | Function/method declaration: defines executable logic. |
| 39 | Conditional branch: executes following block when condition is true. |
| 40 | Assignment: sets or updates variable/field value. |
| 41 | Assignment: sets or updates variable/field value. |
| 42 | Block end: closes previous function/branch/type block. |
| 43 | Conditional branch: executes following block when condition is true. |
| 44 | Return statement: returns values (and often error) to caller. |
| 45 | Block end: closes previous function/branch/type block. |
| 46 | Return statement: returns values (and often error) to caller. |
| 47 | Block end: closes previous function/branch/type block. |

### File: internal/order/http.go

```go
   1: package main
   2: 
   3: import (
   4: 	"net/http"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   7: 	"github.com/ghost-yu/go_shop_second/order/app"
   8: 	"github.com/ghost-yu/go_shop_second/order/app/command"
   9: 	"github.com/ghost-yu/go_shop_second/order/app/query"
  10: 	"github.com/gin-gonic/gin"
  11: )
  12: 
  13: type HTTPServer struct {
  14: 	app app.Application
  15: }
  16: 
  17: func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
  18: 	var req orderpb.CreateOrderRequest
  19: 	if err := c.ShouldBindJSON(&req); err != nil {
  20: 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
  21: 		return
  22: 	}
  23: 	r, err := H.app.Commands.CreateOrder.Handle(c, command.CreateOrder{
  24: 		CustomerID: req.CustomerID,
  25: 		Items:      req.Items,
  26: 	})
  27: 	if err != nil {
  28: 		c.JSON(http.StatusOK, gin.H{"error": err})
  29: 		return
  30: 	}
  31: 	c.JSON(http.StatusOK, gin.H{
  32: 		"message":     "success",
  33: 		"customer_id": req.CustomerID,
  34: 		"order_id":    r.OrderID,
  35: 	})
  36: }
  37: 
  38: func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
  39: 	o, err := H.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
  40: 		OrderID:    orderID,
  41: 		CustomerID: customerID,
  42: 	})
  43: 	if err != nil {
  44: 		c.JSON(http.StatusOK, gin.H{"error": err})
  45: 		return
  46: 	}
  47: 	c.JSON(http.StatusOK, gin.H{"message": "success", "data": o})
  48: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Block end: closes import or parameter list block. |
| 12 | Blank line: separates code blocks for readability. |
| 13 | Type declaration: defines struct/interface/type alias. |
| 14 | Executable statement: participates in current implementation logic. |
| 15 | Block end: closes previous function/branch/type block. |
| 16 | Blank line: separates code blocks for readability. |
| 17 | Function/method declaration: defines executable logic. |
| 18 | Executable statement: participates in current implementation logic. |
| 19 | Conditional branch: executes following block when condition is true. |
| 20 | Executable statement: participates in current implementation logic. |
| 21 | Return statement: returns values (and often error) to caller. |
| 22 | Block end: closes previous function/branch/type block. |
| 23 | Short declaration: defines and initializes local variable. |
| 24 | Executable statement: participates in current implementation logic. |
| 25 | Executable statement: participates in current implementation logic. |
| 26 | Block end: closes previous function/branch/type block. |
| 27 | Conditional branch: executes following block when condition is true. |
| 28 | Executable statement: participates in current implementation logic. |
| 29 | Return statement: returns values (and often error) to caller. |
| 30 | Block end: closes previous function/branch/type block. |
| 31 | Executable statement: participates in current implementation logic. |
| 32 | Executable statement: participates in current implementation logic. |
| 33 | Executable statement: participates in current implementation logic. |
| 34 | Executable statement: participates in current implementation logic. |
| 35 | Block end: closes previous function/branch/type block. |
| 36 | Block end: closes previous function/branch/type block. |
| 37 | Blank line: separates code blocks for readability. |
| 38 | Function/method declaration: defines executable logic. |
| 39 | Short declaration: defines and initializes local variable. |
| 40 | Executable statement: participates in current implementation logic. |
| 41 | Executable statement: participates in current implementation logic. |
| 42 | Block end: closes previous function/branch/type block. |
| 43 | Conditional branch: executes following block when condition is true. |
| 44 | Executable statement: participates in current implementation logic. |
| 45 | Return statement: returns values (and often error) to caller. |
| 46 | Block end: closes previous function/branch/type block. |
| 47 | Executable statement: participates in current implementation logic. |
| 48 | Block end: closes previous function/branch/type block. |

### File: internal/order/service/application.go

```go
   1: package service
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/metrics"
   7: 	"github.com/ghost-yu/go_shop_second/order/adapters"
   8: 	"github.com/ghost-yu/go_shop_second/order/app"
   9: 	"github.com/ghost-yu/go_shop_second/order/app/command"
  10: 	"github.com/ghost-yu/go_shop_second/order/app/query"
  11: 	"github.com/sirupsen/logrus"
  12: )
  13: 
  14: func NewApplication(ctx context.Context) app.Application {
  15: 	orderRepo := adapters.NewMemoryOrderRepository()
  16: 	logger := logrus.NewEntry(logrus.StandardLogger())
  17: 	metricClient := metrics.TodoMetrics{}
  18: 	return app.Application{
  19: 		Commands: app.Commands{
  20: 			CreateOrder: command.NewCreateOrderHandler(orderRepo, logger, metricClient),
  21: 			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
  22: 		},
  23: 		Queries: app.Queries{
  24: 			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
  25: 		},
  26: 	}
  27: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Executable statement: participates in current implementation logic. |
| 12 | Block end: closes import or parameter list block. |
| 13 | Blank line: separates code blocks for readability. |
| 14 | Function/method declaration: defines executable logic. |
| 15 | Short declaration: defines and initializes local variable. |
| 16 | Short declaration: defines and initializes local variable. |
| 17 | Short declaration: defines and initializes local variable. |
| 18 | Return statement: returns values (and often error) to caller. |
| 19 | Executable statement: participates in current implementation logic. |
| 20 | Executable statement: participates in current implementation logic. |
| 21 | Executable statement: participates in current implementation logic. |
| 22 | Block end: closes previous function/branch/type block. |
| 23 | Executable statement: participates in current implementation logic. |
| 24 | Executable statement: participates in current implementation logic. |
| 25 | Block end: closes previous function/branch/type block. |
| 26 | Block end: closes previous function/branch/type block. |
| 27 | Block end: closes previous function/branch/type block. |

## Commit 9: [78d6465] order->stock grpc

### File: internal/order/adapters/grpc/stock_grpc.go

```go
   1: package grpc
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
   8: 	"github.com/sirupsen/logrus"
   9: )
  10: 
  11: type StockGRPC struct {
  12: 	client stockpb.StockServiceClient
  13: }
  14: 
  15: func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
  16: 	return &StockGRPC{client: client}
  17: }
  18: 
  19: func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) error {
  20: 	resp, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
  21: 	logrus.Info("stock_grpc response", resp)
  22: 	return err
  23: }
  24: 
  25: func (s StockGRPC) GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error) {
  26: 	resp, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{ItemIDs: itemIDs})
  27: 	if err != nil {
  28: 		return nil, err
  29: 	}
  30: 	return resp.Items, nil
  31: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Block end: closes import or parameter list block. |
| 10 | Blank line: separates code blocks for readability. |
| 11 | Type declaration: defines struct/interface/type alias. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Block end: closes previous function/branch/type block. |
| 14 | Blank line: separates code blocks for readability. |
| 15 | Function/method declaration: defines executable logic. |
| 16 | Return statement: returns values (and often error) to caller. |
| 17 | Block end: closes previous function/branch/type block. |
| 18 | Blank line: separates code blocks for readability. |
| 19 | Function/method declaration: defines executable logic. |
| 20 | Short declaration: defines and initializes local variable. |
| 21 | Executable statement: participates in current implementation logic. |
| 22 | Return statement: returns values (and often error) to caller. |
| 23 | Block end: closes previous function/branch/type block. |
| 24 | Blank line: separates code blocks for readability. |
| 25 | Function/method declaration: defines executable logic. |
| 26 | Short declaration: defines and initializes local variable. |
| 27 | Conditional branch: executes following block when condition is true. |
| 28 | Return statement: returns values (and often error) to caller. |
| 29 | Block end: closes previous function/branch/type block. |
| 30 | Return statement: returns values (and often error) to caller. |
| 31 | Block end: closes previous function/branch/type block. |

### File: internal/order/app/command/create_order.go

```go
   1: package command
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/decorator"
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   8: 	"github.com/ghost-yu/go_shop_second/order/app/query"
   9: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
  10: 	"github.com/sirupsen/logrus"
  11: )
  12: 
  13: type CreateOrder struct {
  14: 	CustomerID string
  15: 	Items      []*orderpb.ItemWithQuantity
  16: }
  17: 
  18: type CreateOrderResult struct {
  19: 	OrderID string
  20: }
  21: 
  22: type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]
  23: 
  24: type createOrderHandler struct {
  25: 	orderRepo domain.Repository
  26: 	stockGRPC query.StockService
  27: }
  28: 
  29: func NewCreateOrderHandler(
  30: 	orderRepo domain.Repository,
  31: 	stockGRPC query.StockService,
  32: 	logger *logrus.Entry,
  33: 	metricClient decorator.MetricsClient,
  34: ) CreateOrderHandler {
  35: 	if orderRepo == nil {
  36: 		panic("nil orderRepo")
  37: 	}
  38: 	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
  39: 		createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC},
  40: 		logger,
  41: 		metricClient,
  42: 	)
  43: }
  44: 
  45: func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
  46: 	// TODO: call stock grpc to get items.
  47: 	err := c.stockGRPC.CheckIfItemsInStock(ctx, cmd.Items)
  48: 	resp, err := c.stockGRPC.GetItems(ctx, []string{"123"})
  49: 	logrus.Info("createOrderHandler||resp from stockGRPC.GetItems", resp)
  50: 	var stockResponse []*orderpb.Item
  51: 	for _, item := range cmd.Items {
  52: 		stockResponse = append(stockResponse, &orderpb.Item{
  53: 			ID:       item.ID,
  54: 			Quantity: item.Quantity,
  55: 		})
  56: 	}
  57: 	o, err := c.orderRepo.Create(ctx, &domain.Order{
  58: 		CustomerID: cmd.CustomerID,
  59: 		Items:      stockResponse,
  60: 	})
  61: 	if err != nil {
  62: 		return nil, err
  63: 	}
  64: 	return &CreateOrderResult{OrderID: o.ID}, nil
  65: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Block end: closes import or parameter list block. |
| 12 | Blank line: separates code blocks for readability. |
| 13 | Type declaration: defines struct/interface/type alias. |
| 14 | Executable statement: participates in current implementation logic. |
| 15 | Executable statement: participates in current implementation logic. |
| 16 | Block end: closes previous function/branch/type block. |
| 17 | Blank line: separates code blocks for readability. |
| 18 | Type declaration: defines struct/interface/type alias. |
| 19 | Executable statement: participates in current implementation logic. |
| 20 | Block end: closes previous function/branch/type block. |
| 21 | Blank line: separates code blocks for readability. |
| 22 | Type declaration: defines struct/interface/type alias. |
| 23 | Blank line: separates code blocks for readability. |
| 24 | Type declaration: defines struct/interface/type alias. |
| 25 | Executable statement: participates in current implementation logic. |
| 26 | Executable statement: participates in current implementation logic. |
| 27 | Block end: closes previous function/branch/type block. |
| 28 | Blank line: separates code blocks for readability. |
| 29 | Function/method declaration: defines executable logic. |
| 30 | Executable statement: participates in current implementation logic. |
| 31 | Executable statement: participates in current implementation logic. |
| 32 | Executable statement: participates in current implementation logic. |
| 33 | Executable statement: participates in current implementation logic. |
| 34 | Executable statement: participates in current implementation logic. |
| 35 | Conditional branch: executes following block when condition is true. |
| 36 | Panic: aborts normal control flow with unrecoverable error. |
| 37 | Block end: closes previous function/branch/type block. |
| 38 | Return statement: returns values (and often error) to caller. |
| 39 | Executable statement: participates in current implementation logic. |
| 40 | Executable statement: participates in current implementation logic. |
| 41 | Executable statement: participates in current implementation logic. |
| 42 | Block end: closes import or parameter list block. |
| 43 | Block end: closes previous function/branch/type block. |
| 44 | Blank line: separates code blocks for readability. |
| 45 | Function/method declaration: defines executable logic. |
| 46 | Comment: explains intent, todo, or context. |
| 47 | Short declaration: defines and initializes local variable. |
| 48 | Short declaration: defines and initializes local variable. |
| 49 | Executable statement: participates in current implementation logic. |
| 50 | Executable statement: participates in current implementation logic. |
| 51 | Loop: iterates collection or repeats logic. |
| 52 | Assignment: sets or updates variable/field value. |
| 53 | Executable statement: participates in current implementation logic. |
| 54 | Executable statement: participates in current implementation logic. |
| 55 | Block end: closes previous function/branch/type block. |
| 56 | Block end: closes previous function/branch/type block. |
| 57 | Short declaration: defines and initializes local variable. |
| 58 | Executable statement: participates in current implementation logic. |
| 59 | Executable statement: participates in current implementation logic. |
| 60 | Block end: closes previous function/branch/type block. |
| 61 | Conditional branch: executes following block when condition is true. |
| 62 | Return statement: returns values (and often error) to caller. |
| 63 | Block end: closes previous function/branch/type block. |
| 64 | Return statement: returns values (and often error) to caller. |
| 65 | Block end: closes previous function/branch/type block. |

### File: internal/order/app/query/service.go

```go
   1: package query
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   7: )
   8: 
   9: type StockService interface {
  10: 	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) error
  11: 	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
  12: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Block end: closes import or parameter list block. |
| 8 | Blank line: separates code blocks for readability. |
| 9 | Type declaration: defines struct/interface/type alias. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Executable statement: participates in current implementation logic. |
| 12 | Block end: closes previous function/branch/type block. |

### File: internal/order/main.go

```go
   1: package main
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/config"
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   8: 	"github.com/ghost-yu/go_shop_second/common/server"
   9: 	"github.com/ghost-yu/go_shop_second/order/ports"
  10: 	"github.com/ghost-yu/go_shop_second/order/service"
  11: 	"github.com/gin-gonic/gin"
  12: 	"github.com/sirupsen/logrus"
  13: 	"github.com/spf13/viper"
  14: 	"google.golang.org/grpc"
  15: )
  16: 
  17: func init() {
  18: 	if err := config.NewViperConfig(); err != nil {
  19: 		logrus.Fatal(err)
  20: 	}
  21: }
  22: 
  23: func main() {
  24: 	serviceName := viper.GetString("order.service-name")
  25: 
  26: 	ctx, cancel := context.WithCancel(context.Background())
  27: 	defer cancel()
  28: 
  29: 	application, cleanup := service.NewApplication(ctx)
  30: 	defer cleanup()
  31: 
  32: 	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
  33: 		svc := ports.NewGRPCServer(application)
  34: 		orderpb.RegisterOrderServiceServer(server, svc)
  35: 	})
  36: 
  37: 	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
  38: 		ports.RegisterHandlersWithOptions(router, HTTPServer{
  39: 			app: application,
  40: 		}, ports.GinServerOptions{
  41: 			BaseURL:      "/api",
  42: 			Middlewares:  nil,
  43: 			ErrorHandler: nil,
  44: 		})
  45: 	})
  46: 
  47: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Executable statement: participates in current implementation logic. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Executable statement: participates in current implementation logic. |
| 14 | Executable statement: participates in current implementation logic. |
| 15 | Block end: closes import or parameter list block. |
| 16 | Blank line: separates code blocks for readability. |
| 17 | Function/method declaration: defines executable logic. |
| 18 | Conditional branch: executes following block when condition is true. |
| 19 | Executable statement: participates in current implementation logic. |
| 20 | Block end: closes previous function/branch/type block. |
| 21 | Block end: closes previous function/branch/type block. |
| 22 | Blank line: separates code blocks for readability. |
| 23 | Function/method declaration: defines executable logic. |
| 24 | Short declaration: defines and initializes local variable. |
| 25 | Blank line: separates code blocks for readability. |
| 26 | Short declaration: defines and initializes local variable. |
| 27 | Deferred call: executes when function returns. |
| 28 | Blank line: separates code blocks for readability. |
| 29 | Short declaration: defines and initializes local variable. |
| 30 | Deferred call: executes when function returns. |
| 31 | Blank line: separates code blocks for readability. |
| 32 | Goroutine launch: runs function asynchronously. |
| 33 | Short declaration: defines and initializes local variable. |
| 34 | Executable statement: participates in current implementation logic. |
| 35 | Block end: closes previous function/branch/type block. |
| 36 | Blank line: separates code blocks for readability. |
| 37 | Executable statement: participates in current implementation logic. |
| 38 | Executable statement: participates in current implementation logic. |
| 39 | Executable statement: participates in current implementation logic. |
| 40 | Block end: closes previous function/branch/type block. |
| 41 | Executable statement: participates in current implementation logic. |
| 42 | Executable statement: participates in current implementation logic. |
| 43 | Executable statement: participates in current implementation logic. |
| 44 | Block end: closes previous function/branch/type block. |
| 45 | Block end: closes previous function/branch/type block. |
| 46 | Blank line: separates code blocks for readability. |
| 47 | Block end: closes previous function/branch/type block. |

### File: internal/order/service/application.go

```go
   1: package service
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	grpcClient "github.com/ghost-yu/go_shop_second/common/client"
   7: 	"github.com/ghost-yu/go_shop_second/common/metrics"
   8: 	"github.com/ghost-yu/go_shop_second/order/adapters"
   9: 	"github.com/ghost-yu/go_shop_second/order/adapters/grpc"
  10: 	"github.com/ghost-yu/go_shop_second/order/app"
  11: 	"github.com/ghost-yu/go_shop_second/order/app/command"
  12: 	"github.com/ghost-yu/go_shop_second/order/app/query"
  13: 	"github.com/sirupsen/logrus"
  14: )
  15: 
  16: func NewApplication(ctx context.Context) (app.Application, func()) {
  17: 	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
  18: 	if err != nil {
  19: 		panic(err)
  20: 	}
  21: 	stockGRPC := grpc.NewStockGRPC(stockClient)
  22: 	return newApplication(ctx, stockGRPC), func() {
  23: 		_ = closeStockClient()
  24: 	}
  25: }
  26: 
  27: func newApplication(_ context.Context, stockGRPC query.StockService) app.Application {
  28: 	orderRepo := adapters.NewMemoryOrderRepository()
  29: 	logger := logrus.NewEntry(logrus.StandardLogger())
  30: 	metricClient := metrics.TodoMetrics{}
  31: 	return app.Application{
  32: 		Commands: app.Commands{
  33: 			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, logger, metricClient),
  34: 			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
  35: 		},
  36: 		Queries: app.Queries{
  37: 			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
  38: 		},
  39: 	}
  40: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Executable statement: participates in current implementation logic. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Executable statement: participates in current implementation logic. |
| 14 | Block end: closes import or parameter list block. |
| 15 | Blank line: separates code blocks for readability. |
| 16 | Function/method declaration: defines executable logic. |
| 17 | Short declaration: defines and initializes local variable. |
| 18 | Conditional branch: executes following block when condition is true. |
| 19 | Panic: aborts normal control flow with unrecoverable error. |
| 20 | Block end: closes previous function/branch/type block. |
| 21 | Short declaration: defines and initializes local variable. |
| 22 | Return statement: returns values (and often error) to caller. |
| 23 | Assignment: sets or updates variable/field value. |
| 24 | Block end: closes previous function/branch/type block. |
| 25 | Block end: closes previous function/branch/type block. |
| 26 | Blank line: separates code blocks for readability. |
| 27 | Function/method declaration: defines executable logic. |
| 28 | Short declaration: defines and initializes local variable. |
| 29 | Short declaration: defines and initializes local variable. |
| 30 | Short declaration: defines and initializes local variable. |
| 31 | Return statement: returns values (and often error) to caller. |
| 32 | Executable statement: participates in current implementation logic. |
| 33 | Executable statement: participates in current implementation logic. |
| 34 | Executable statement: participates in current implementation logic. |
| 35 | Block end: closes previous function/branch/type block. |
| 36 | Executable statement: participates in current implementation logic. |
| 37 | Executable statement: participates in current implementation logic. |
| 38 | Block end: closes previous function/branch/type block. |
| 39 | Block end: closes previous function/branch/type block. |
| 40 | Block end: closes previous function/branch/type block. |

## Commit 10: [8ee5f39] stock query and grpc

### File: internal/order/adapters/grpc/stock_grpc.go

```go
   1: package grpc
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
   8: 	"github.com/sirupsen/logrus"
   9: )
  10: 
  11: type StockGRPC struct {
  12: 	client stockpb.StockServiceClient
  13: }
  14: 
  15: func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
  16: 	return &StockGRPC{client: client}
  17: }
  18: 
  19: func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error) {
  20: 	resp, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
  21: 	logrus.Info("stock_grpc response", resp)
  22: 	return resp, err
  23: }
  24: 
  25: func (s StockGRPC) GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error) {
  26: 	resp, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{ItemIDs: itemIDs})
  27: 	if err != nil {
  28: 		return nil, err
  29: 	}
  30: 	return resp.Items, nil
  31: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Block end: closes import or parameter list block. |
| 10 | Blank line: separates code blocks for readability. |
| 11 | Type declaration: defines struct/interface/type alias. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Block end: closes previous function/branch/type block. |
| 14 | Blank line: separates code blocks for readability. |
| 15 | Function/method declaration: defines executable logic. |
| 16 | Return statement: returns values (and often error) to caller. |
| 17 | Block end: closes previous function/branch/type block. |
| 18 | Blank line: separates code blocks for readability. |
| 19 | Function/method declaration: defines executable logic. |
| 20 | Short declaration: defines and initializes local variable. |
| 21 | Executable statement: participates in current implementation logic. |
| 22 | Return statement: returns values (and often error) to caller. |
| 23 | Block end: closes previous function/branch/type block. |
| 24 | Blank line: separates code blocks for readability. |
| 25 | Function/method declaration: defines executable logic. |
| 26 | Short declaration: defines and initializes local variable. |
| 27 | Conditional branch: executes following block when condition is true. |
| 28 | Return statement: returns values (and often error) to caller. |
| 29 | Block end: closes previous function/branch/type block. |
| 30 | Return statement: returns values (and often error) to caller. |
| 31 | Block end: closes previous function/branch/type block. |

### File: internal/order/app/command/create_order.go

```go
   1: package command
   2: 
   3: import (
   4: 	"context"
   5: 	"errors"
   6: 
   7: 	"github.com/ghost-yu/go_shop_second/common/decorator"
   8: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   9: 	"github.com/ghost-yu/go_shop_second/order/app/query"
  10: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
  11: 	"github.com/sirupsen/logrus"
  12: )
  13: 
  14: type CreateOrder struct {
  15: 	CustomerID string
  16: 	Items      []*orderpb.ItemWithQuantity
  17: }
  18: 
  19: type CreateOrderResult struct {
  20: 	OrderID string
  21: }
  22: 
  23: type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]
  24: 
  25: type createOrderHandler struct {
  26: 	orderRepo domain.Repository
  27: 	stockGRPC query.StockService
  28: }
  29: 
  30: func NewCreateOrderHandler(
  31: 	orderRepo domain.Repository,
  32: 	stockGRPC query.StockService,
  33: 	logger *logrus.Entry,
  34: 	metricClient decorator.MetricsClient,
  35: ) CreateOrderHandler {
  36: 	if orderRepo == nil {
  37: 		panic("nil orderRepo")
  38: 	}
  39: 	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
  40: 		createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC},
  41: 		logger,
  42: 		metricClient,
  43: 	)
  44: }
  45: 
  46: func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
  47: 	validItems, err := c.validate(ctx, cmd.Items)
  48: 	if err != nil {
  49: 		return nil, err
  50: 	}
  51: 	o, err := c.orderRepo.Create(ctx, &domain.Order{
  52: 		CustomerID: cmd.CustomerID,
  53: 		Items:      validItems,
  54: 	})
  55: 	if err != nil {
  56: 		return nil, err
  57: 	}
  58: 	return &CreateOrderResult{OrderID: o.ID}, nil
  59: }
  60: 
  61: func (c createOrderHandler) validate(ctx context.Context, items []*orderpb.ItemWithQuantity) ([]*orderpb.Item, error) {
  62: 	if len(items) == 0 {
  63: 		return nil, errors.New("must have at least one item")
  64: 	}
  65: 	items = packItems(items)
  66: 	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, items)
  67: 	if err != nil {
  68: 		return nil, err
  69: 	}
  70: 	return resp.Items, nil
  71: }
  72: 
  73: func packItems(items []*orderpb.ItemWithQuantity) []*orderpb.ItemWithQuantity {
  74: 	merged := make(map[string]int32)
  75: 	for _, item := range items {
  76: 		merged[item.ID] += item.Quantity
  77: 	}
  78: 	var res []*orderpb.ItemWithQuantity
  79: 	for id, quantity := range merged {
  80: 		res = append(res, &orderpb.ItemWithQuantity{
  81: 			ID:       id,
  82: 			Quantity: quantity,
  83: 		})
  84: 	}
  85: 	return res
  86: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Executable statement: participates in current implementation logic. |
| 6 | Blank line: separates code blocks for readability. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Executable statement: participates in current implementation logic. |
| 9 | Executable statement: participates in current implementation logic. |
| 10 | Executable statement: participates in current implementation logic. |
| 11 | Executable statement: participates in current implementation logic. |
| 12 | Block end: closes import or parameter list block. |
| 13 | Blank line: separates code blocks for readability. |
| 14 | Type declaration: defines struct/interface/type alias. |
| 15 | Executable statement: participates in current implementation logic. |
| 16 | Executable statement: participates in current implementation logic. |
| 17 | Block end: closes previous function/branch/type block. |
| 18 | Blank line: separates code blocks for readability. |
| 19 | Type declaration: defines struct/interface/type alias. |
| 20 | Executable statement: participates in current implementation logic. |
| 21 | Block end: closes previous function/branch/type block. |
| 22 | Blank line: separates code blocks for readability. |
| 23 | Type declaration: defines struct/interface/type alias. |
| 24 | Blank line: separates code blocks for readability. |
| 25 | Type declaration: defines struct/interface/type alias. |
| 26 | Executable statement: participates in current implementation logic. |
| 27 | Executable statement: participates in current implementation logic. |
| 28 | Block end: closes previous function/branch/type block. |
| 29 | Blank line: separates code blocks for readability. |
| 30 | Function/method declaration: defines executable logic. |
| 31 | Executable statement: participates in current implementation logic. |
| 32 | Executable statement: participates in current implementation logic. |
| 33 | Executable statement: participates in current implementation logic. |
| 34 | Executable statement: participates in current implementation logic. |
| 35 | Executable statement: participates in current implementation logic. |
| 36 | Conditional branch: executes following block when condition is true. |
| 37 | Panic: aborts normal control flow with unrecoverable error. |
| 38 | Block end: closes previous function/branch/type block. |
| 39 | Return statement: returns values (and often error) to caller. |
| 40 | Executable statement: participates in current implementation logic. |
| 41 | Executable statement: participates in current implementation logic. |
| 42 | Executable statement: participates in current implementation logic. |
| 43 | Block end: closes import or parameter list block. |
| 44 | Block end: closes previous function/branch/type block. |
| 45 | Blank line: separates code blocks for readability. |
| 46 | Function/method declaration: defines executable logic. |
| 47 | Short declaration: defines and initializes local variable. |
| 48 | Conditional branch: executes following block when condition is true. |
| 49 | Return statement: returns values (and often error) to caller. |
| 50 | Block end: closes previous function/branch/type block. |
| 51 | Short declaration: defines and initializes local variable. |
| 52 | Executable statement: participates in current implementation logic. |
| 53 | Executable statement: participates in current implementation logic. |
| 54 | Block end: closes previous function/branch/type block. |
| 55 | Conditional branch: executes following block when condition is true. |
| 56 | Return statement: returns values (and often error) to caller. |
| 57 | Block end: closes previous function/branch/type block. |
| 58 | Return statement: returns values (and often error) to caller. |
| 59 | Block end: closes previous function/branch/type block. |
| 60 | Blank line: separates code blocks for readability. |
| 61 | Function/method declaration: defines executable logic. |
| 62 | Conditional branch: executes following block when condition is true. |
| 63 | Return statement: returns values (and often error) to caller. |
| 64 | Block end: closes previous function/branch/type block. |
| 65 | Assignment: sets or updates variable/field value. |
| 66 | Short declaration: defines and initializes local variable. |
| 67 | Conditional branch: executes following block when condition is true. |
| 68 | Return statement: returns values (and often error) to caller. |
| 69 | Block end: closes previous function/branch/type block. |
| 70 | Return statement: returns values (and often error) to caller. |
| 71 | Block end: closes previous function/branch/type block. |
| 72 | Blank line: separates code blocks for readability. |
| 73 | Function/method declaration: defines executable logic. |
| 74 | Short declaration: defines and initializes local variable. |
| 75 | Loop: iterates collection or repeats logic. |
| 76 | Assignment: sets or updates variable/field value. |
| 77 | Block end: closes previous function/branch/type block. |
| 78 | Executable statement: participates in current implementation logic. |
| 79 | Loop: iterates collection or repeats logic. |
| 80 | Assignment: sets or updates variable/field value. |
| 81 | Executable statement: participates in current implementation logic. |
| 82 | Executable statement: participates in current implementation logic. |
| 83 | Block end: closes previous function/branch/type block. |
| 84 | Block end: closes previous function/branch/type block. |
| 85 | Return statement: returns values (and often error) to caller. |
| 86 | Block end: closes previous function/branch/type block. |

### File: internal/order/app/query/service.go

```go
   1: package query
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
   8: )
   9: 
  10: type StockService interface {
  11: 	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error)
  12: 	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
  13: }
```

| Line | Explanation |
| --- | --- |
| 1 | Package declaration: defines which package this file belongs to. |
| 2 | Blank line: separates code blocks for readability. |
| 3 | Import block start: declares dependencies used by this file. |
| 4 | Executable statement: participates in current implementation logic. |
| 5 | Blank line: separates code blocks for readability. |
| 6 | Executable statement: participates in current implementation logic. |
| 7 | Executable statement: participates in current implementation logic. |
| 8 | Block end: closes import or parameter list block. |
| 9 | Blank line: separates code blocks for readability. |
| 10 | Type declaration: defines struct/interface/type alias. |
| 11 | Executable statement: participates in current implementation logic. |
| 12 | Executable statement: participates in current implementation logic. |
| 13 | Block end: closes previous function/branch/type block. |



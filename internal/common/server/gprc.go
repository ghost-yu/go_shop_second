package server

import (
	"net"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// init 会在包加载时执行，这里统一替换 gRPC 默认日志实现，
// 让框架日志和业务日志都走 logrus，便于初学者排查问题时只看一种格式。
func init() {
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)
	grpc_logrus.ReplaceGrpcLogger(logrus.NewEntry(logger))
}

// RunGRPCServer 负责按服务名读取配置，再把真正的启动动作委托给 RunGRPCServerOnAddr。
// 这样测试时可以直接传地址，避免每次都依赖 viper 配置。
func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server)) {
	addr := viper.Sub(serviceName).GetString("grpc-addr")
	if addr == "" {
		// TODO: Warning log
		addr = viper.GetString("fallback-grpc-addr")
	}
	RunGRPCServerOnAddr(addr, registerServer)
}

// RunGRPCServerOnAddr 创建 gRPC server，并通过 registerServer 回调注册具体业务服务。
// 这类“框架负责启动，业务负责注册”的回调模式，在 Go 后端里非常常见。
func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			// 一元拦截器链类似 HTTP 中间件，用来放日志、指标、鉴权这类横切逻辑。
			grpc_tags.UnaryServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
			//otelgrpc.UnaryServerInterceptor(),
			//srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			//logging.UnaryServerInterceptor(interceptorLogger(rpcLogger), logging.WithFieldsFromContext(logTraceID)),
			//selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authFn), selector.MatchFunc(allButHealthZ)),
			//recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
		grpc.ChainStreamInterceptor(
			// 流式 RPC 和一元 RPC 分开配置，是因为两类调用的拦截接口不同。
			grpc_tags.StreamServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(logrusEntry),
			//otelgrpc.StreamServerInterceptor(),
			//srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			//logging.StreamServerInterceptor(interceptorLogger(rpcLogger), logging.WithFieldsFromContext(logTraceID)),
			//selector.StreamServerInterceptor(auth.StreamServerInterceptor(authFn), selector.MatchFunc(allButHealthZ)),
			//recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
	)
	// 业务方在这里把 OrderService/StockService 的实现挂到 grpcServer 上。
	registerServer(grpcServer)

	// net.Listen 只负责占住端口；真正开始接收 gRPC 请求要等 Serve 调用后才发生。
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Infof("Starting gRPC server, Listening: %s", addr)
	if err := grpcServer.Serve(listen); err != nil {
		logrus.Panic(err)
	}
}

package main

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/config"
	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
	"github.com/ghost-yu/go_shop_second/common/server"
	"github.com/ghost-yu/go_shop_second/order/ports"
	"github.com/ghost-yu/go_shop_second/order/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// init 先加载全局配置，这样 main 里启动服务时就能直接从 viper 读取地址和服务名。
func init() {
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	// serviceName 对应 global.yaml 里的 order.service-name，后续 server 包会继续用它找端口配置。
	serviceName := viper.GetString("order.service-name")

	// ctx 用来把进程级生命周期传给下游依赖，例如 gRPC 客户端连接。
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// NewApplication 是组合根：在这里把 repo、远程 client、handler 全部装配好。
	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	// gRPC 和 HTTP 同时启动，说明 order 服务对内和对外用了两套协议暴露能力。
	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		svc := ports.NewGRPCServer(application)
		orderpb.RegisterOrderServiceServer(server, svc)
	})

	// HTTP 侧通过 OpenAPI 生成的 RegisterHandlersWithOptions 完成路由到实现的映射。
	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, HTTPServer{
			app: application,
		}, ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})

}

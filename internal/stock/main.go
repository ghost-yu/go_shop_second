package main

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/config"
	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
	"github.com/ghost-yu/go_shop_second/common/server"
	"github.com/ghost-yu/go_shop_second/stock/ports"
	"github.com/ghost-yu/go_shop_second/stock/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// stock 服务和 order 服务一样，在进程入口先加载配置，避免后续各层自己读文件。
func init() {
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	// serviceName 决定去哪个配置分组下找 grpc-addr 等参数。
	serviceName := viper.GetString("stock.service-name")
	// serverType 让同一个二进制可以选择跑哪种协议，当前 lesson 主要走 grpc 分支。
	serverType := viper.GetString("stock.server-to-run")

	logrus.Info(serverType)

	// 这里先组装应用层，再根据配置决定挂到哪个端口适配器上。
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := service.NewApplication(ctx)
	switch serverType {
	case "grpc":
		// gRPC 分支把 ports.GRPCServer 注册为 StockService 的服务端实现。
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			svc := ports.NewGRPCServer(application)
			stockpb.RegisterStockServiceServer(server, svc)
		})
	case "http":
		// 暂时不用
	default:
		panic("unexpected server type")
	}
}

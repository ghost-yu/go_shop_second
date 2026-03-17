package service

import (
	"context"

	grpcClient "github.com/ghost-yu/go_shop_second/common/client"
	"github.com/ghost-yu/go_shop_second/common/metrics"
	"github.com/ghost-yu/go_shop_second/order/adapters"
	"github.com/ghost-yu/go_shop_second/order/adapters/grpc"
	"github.com/ghost-yu/go_shop_second/order/app"
	"github.com/ghost-yu/go_shop_second/order/app/command"
	"github.com/ghost-yu/go_shop_second/order/app/query"
	"github.com/sirupsen/logrus"
)

// NewApplication 是 order 服务的组合根。
// 它负责创建外部依赖、构造适配器，并返回一个已经装配好的应用层门面。
func NewApplication(ctx context.Context) (app.Application, func()) {
	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	// 这里把生成的 gRPC client 再包一层适配器，避免应用层直接依赖 proto 细节。
	stockGRPC := grpc.NewStockGRPC(stockClient)
	return newApplication(ctx, stockGRPC), func() {
		_ = closeStockClient()
	}
}

func newApplication(_ context.Context, stockGRPC query.StockService) app.Application {
	// 组合根里统一决定“接口对应哪种实现”，后面要切数据库或 mock 时只改这里。
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, logger, metricClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}

package service

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/broker"
	grpcClient "github.com/ghost-yu/go_shop_second/common/client"
	"github.com/ghost-yu/go_shop_second/common/metrics"
	"github.com/ghost-yu/go_shop_second/order/adapters"
	"github.com/ghost-yu/go_shop_second/order/adapters/grpc"
	"github.com/ghost-yu/go_shop_second/order/app"
	"github.com/ghost-yu/go_shop_second/order/app/command"
	"github.com/ghost-yu/go_shop_second/order/app/query"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	stockGRPC := grpc.NewStockGRPC(stockClient)
	return newApplication(ctx, stockGRPC, ch), func() {
		_ = closeStockClient()
		_ = closeCh()
		_ = ch.Close()
	}
}

func newApplication(_ context.Context, stockGRPC query.StockService, ch *amqp.Channel) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, ch, logger, metricClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}

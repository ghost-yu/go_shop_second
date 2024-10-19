package service

import (
	"context"

	grpcClient "github.com/ghost-yu/go_shop_second/common/client"
	"github.com/ghost-yu/go_shop_second/common/metrics"
	"github.com/ghost-yu/go_shop_second/payment/adapters"
	"github.com/ghost-yu/go_shop_second/payment/app"
	"github.com/ghost-yu/go_shop_second/payment/app/command"
	"github.com/ghost-yu/go_shop_second/payment/domain"
	"github.com/ghost-yu/go_shop_second/payment/infrastructure/processor"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	orderGRPC := adapters.NewOrderGRPC(orderClient)
	memoryProcessor := processor.NewInmemProcessor()
	return newApplication(ctx, orderGRPC, memoryProcessor), func() {
		_ = closeOrderClient()
	}
}

func newApplication(ctx context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logger, metricClient),
		},
	}
}

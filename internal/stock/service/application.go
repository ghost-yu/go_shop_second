package service

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/metrics"
	"github.com/ghost-yu/go_shop_second/stock/adapters"
	"github.com/ghost-yu/go_shop_second/stock/app"
	"github.com/ghost-yu/go_shop_second/stock/app/query"
	"github.com/ghost-yu/go_shop_second/stock/infrastructure/integration"
	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApplication(_ context.Context) app.Application {
	//stockRepo := adapters.NewMemoryStockRepository()
	db := persistent.NewMySQL()
	stockRepo := adapters.NewMySQLStockRepository(db)
	stripeAPI := integration.NewStripeAPI()
	metricsClient := metrics.NewPrometheusMetricsClient(&metrics.PrometheusMetricsClientConfig{
		Host:        viper.GetString("stock.metrics_export_addr"),
		ServiceName: viper.GetString("stock.service-name"),
	})
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, stripeAPI, logrus.StandardLogger(), metricsClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logrus.StandardLogger(), metricsClient),
		},
	}
}

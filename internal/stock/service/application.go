package service

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/metrics"
	"github.com/ghost-yu/go_shop_second/stock/adapters"
	"github.com/ghost-yu/go_shop_second/stock/app"
	"github.com/ghost-yu/go_shop_second/stock/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(_ context.Context) app.Application {
	stockRepo := adapters.NewMemoryStockRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, logger, metricsClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricsClient),
		},
	}
}

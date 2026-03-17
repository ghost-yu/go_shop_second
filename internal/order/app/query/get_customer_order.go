package query

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/decorator"
	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
	"github.com/sirupsen/logrus"
)

// GetCustomerOrder 是查询订单的输入对象。
// 把参数收成一个结构体后，后续加字段不会破坏函数签名。
type GetCustomerOrder struct {
	CustomerID string
	OrderID    string
}

// GetCustomerOrderHandler 是一个已经套好泛型的查询处理器类型别名。
type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Order]

// getCustomerOrderHandler 才是真正的业务实现，字段里只放它需要的依赖。
type getCustomerOrderHandler struct {
	orderRepo domain.Repository
}

func NewGetCustomerOrderHandler(
	orderRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) GetCustomerOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}
	// 构造函数里统一加装饰器，这样调用方拿到的就是“带日志和指标能力”的 handler。
	return decorator.ApplyQueryDecorators[GetCustomerOrder, *domain.Order](
		getCustomerOrderHandler{orderRepo: orderRepo},
		logger,
		metricClient,
	)
}

func (g getCustomerOrderHandler) Handle(ctx context.Context, query GetCustomerOrder) (*domain.Order, error) {
	// 查询用例本身很薄，只负责描述流程，把数据访问交给仓储接口。
	o, err := g.orderRepo.Get(ctx, query.OrderID, query.CustomerID)
	if err != nil {
		return nil, err
	}
	return o, nil
}

package command

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/decorator"
	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
	"github.com/sirupsen/logrus"
)

// UpdateOrder 把“更新哪个订单”和“如何更新”一起传给 handler。
type UpdateOrder struct {
	Order    *domain.Order
	UpdateFn func(context.Context, *domain.Order) (*domain.Order, error)
}

type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, interface{}]

// updateOrderHandler 负责流程控制，具体修改细节由 UpdateFn 注入。
type updateOrderHandler struct {
	orderRepo domain.Repository
	//stockGRPC
}

func NewUpdateOrderHandler(
	orderRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) UpdateOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}
	return decorator.ApplyCommandDecorators[UpdateOrder, interface{}](
		updateOrderHandler{orderRepo: orderRepo},
		logger,
		metricClient,
	)
}

func (c updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (interface{}, error) {
	// 给 nil UpdateFn 一个 no-op 默认值，避免直接调用时发生空指针问题。
	if cmd.UpdateFn == nil {
		logrus.Warnf("updateOrderHandler got nil UpdateFn, order=%#v", cmd.Order)
		cmd.UpdateFn = func(_ context.Context, order *domain.Order) (*domain.Order, error) { return order, nil }
	}
	if err := c.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFn); err != nil {
		return nil, err
	}
	return nil, nil
}

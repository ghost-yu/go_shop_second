package command

import (
	"context"
	"errors"

	"github.com/ghost-yu/go_shop_second/common/decorator"
	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
	"github.com/ghost-yu/go_shop_second/order/app/query"
	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
	"github.com/sirupsen/logrus"
)

// CreateOrder 是创建订单用例的输入模型。
type CreateOrder struct {
	CustomerID string
	Items      []*orderpb.ItemWithQuantity
}

// CreateOrderResult 是返回给端口层的结果，避免直接暴露底层仓储对象。
type CreateOrderResult struct {
	OrderID string
}

type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]

// createOrderHandler 只持有完成下单流程所需的两个依赖：订单仓储和库存服务。
type createOrderHandler struct {
	orderRepo domain.Repository
	stockGRPC query.StockService
}

func NewCreateOrderHandler(
	orderRepo domain.Repository,
	stockGRPC query.StockService,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CreateOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}
	// 下单也是 command，所以同样套上统一的日志和指标装饰器。
	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC},
		logger,
		metricClient,
	)
}

func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	// 先校验库存，再真正创建订单，避免生成一张无法履约的脏订单。
	validItems, err := c.validate(ctx, cmd.Items)
	if err != nil {
		return nil, err
	}
	o, err := c.orderRepo.Create(ctx, &domain.Order{
		CustomerID: cmd.CustomerID,
		Items:      validItems,
	})
	if err != nil {
		return nil, err
	}
	return &CreateOrderResult{OrderID: o.ID}, nil
}

func (c createOrderHandler) validate(ctx context.Context, items []*orderpb.ItemWithQuantity) ([]*orderpb.Item, error) {
	if len(items) == 0 {
		return nil, errors.New("must have at least one item")
	}
	// packItems 先合并重复商品，避免把同一个商品重复发给库存服务校验。
	items = packItems(items)
	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, items)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func packItems(items []*orderpb.ItemWithQuantity) []*orderpb.ItemWithQuantity {
	// merged 以商品 ID 为 key，把数量累加起来。
	merged := make(map[string]int32)
	for _, item := range items {
		merged[item.ID] += item.Quantity
	}
	var res []*orderpb.ItemWithQuantity
	for id, quantity := range merged {
		res = append(res, &orderpb.ItemWithQuantity{
			ID:       id,
			Quantity: quantity,
		})
	}
	return res
}

package command

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, order *orderpb.Order) error
}

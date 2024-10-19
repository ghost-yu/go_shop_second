package domain

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
)

type Processor interface {
	CreatePaymentLink(context.Context, *orderpb.Order) (string, error)
}

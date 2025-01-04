package processor

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
)

type InmemProcessor struct {
}

func NewInmemProcessor() *InmemProcessor {
	return &InmemProcessor{}
}

func (i InmemProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
	return "inmem-payment-link", nil
}

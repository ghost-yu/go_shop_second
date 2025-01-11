package domain

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/entity"
)

type Processor interface {
	CreatePaymentLink(context.Context, *entity.Order) (string, error)
}

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*entity.Item
}

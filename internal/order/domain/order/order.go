package order

import "github.com/ghost-yu/go_shop_second/common/genproto/orderpb"

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*orderpb.Item
}

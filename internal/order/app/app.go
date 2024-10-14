package app

import "github.com/ghost-yu/go_shop_second/order/app/query"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct{}

type Queries struct {
	GetCustomerOrder query.GetCustomerOrderHandler
}

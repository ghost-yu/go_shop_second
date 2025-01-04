package app

import "github.com/ghost-yu/go_shop_second/stock/app/query"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct{}

type Queries struct {
	CheckIfItemsInStock query.CheckIfItemsInStockHandler
	GetItems            query.GetItemsHandler
}

package app

import (
	"github.com/ghost-yu/go_shop_second/order/app/command"
	"github.com/ghost-yu/go_shop_second/order/app/query"
)

// Application 是应用层门面，让上层只依赖一个总入口，而不用感知每个 handler 的构造细节。
type Application struct {
	Commands Commands
	Queries  Queries
}

// Commands 聚合“会改状态”的用例，典型如创建、更新订单。
type Commands struct {
	CreateOrder command.CreateOrderHandler
	UpdateOrder command.UpdateOrderHandler
}

// Queries 聚合“只读”的用例，便于和 Commands 做职责分离。
type Queries struct {
	GetCustomerOrder query.GetCustomerOrderHandler
}

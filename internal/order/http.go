package main

import (
	"net/http"

	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
	"github.com/ghost-yu/go_shop_second/order/app"
	"github.com/ghost-yu/go_shop_second/order/app/command"
	"github.com/ghost-yu/go_shop_second/order/app/query"
	"github.com/gin-gonic/gin"
)

// HTTPServer 是 OpenAPI 生成接口在 order 服务里的具体实现。
// 它本身不处理复杂业务，只负责把 HTTP 请求翻译成应用层调用。
type HTTPServer struct {
	app app.Application
}

func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	var req orderpb.CreateOrderRequest
	// ShouldBindJSON 负责把请求体反序列化为 proto 请求对象。
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 进入应用层前，把 HTTP 层对象转换成 command 对象，避免业务逻辑依赖 gin。
	r, err := H.app.Commands.CreateOrder.Handle(c, command.CreateOrder{
		CustomerID: req.CustomerID,
		Items:      req.Items,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":     "success",
		"customer_id": req.CustomerID,
		"order_id":    r.OrderID,
	})
}

func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	// 查询场景走 Queries 分组，体现读写分离的组织方式。
	o, err := H.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
		OrderID:    orderID,
		CustomerID: customerID,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": o})
}

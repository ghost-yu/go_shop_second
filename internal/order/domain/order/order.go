package order

import "github.com/ghost-yu/go_shop_second/common/genproto/orderpb"

// Order 是订单领域对象，代表业务里真正被创建、查询、更新的核心实体。
type Order struct {
	// ID 是订单主键，由仓储层创建时生成。
	ID string
	// CustomerID 表示订单属于哪个客户。
	CustomerID string
	// Status 预留给支付中、已支付等订单状态流转。
	Status string
	// PaymentLink 预留给支付服务返回的支付链接。
	PaymentLink string
	// Items 直接复用 proto 里的商品结构，减少服务边界两侧的对象转换成本。
	Items []*orderpb.Item
}

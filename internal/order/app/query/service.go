package query

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
)

// StockService 是 order 应用层眼里的“库存能力端口”。
// 它隔离了底层 gRPC 细节，让 create_order 用例只表达业务意图。
type StockService interface {
	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error)
	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
}

package ports

import (
	context "context"

	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
	"github.com/ghost-yu/go_shop_second/stock/app"
)

// GRPCServer 负责实现 stockpb.StockServiceServer 接口。
// 外部请求进入后，应该在这里完成参数转换，再转给应用层处理。
type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	// 这里后续会把商品 ID 列表转给库存查询用例。
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	// 这一步对应 order 下单前的库存校验，是跨服务调用的关键入口。
	//TODO implement me
	panic("implement me")
}

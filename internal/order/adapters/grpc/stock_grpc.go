package grpc

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
)

// StockGRPC 是 order 侧访问 stock 服务的远程适配器。
// 它实现的是应用层定义的 StockService 接口，而不是把 proto client 直接暴露出去。
type StockGRPC struct {
	client stockpb.StockServiceClient
}

func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
	return &StockGRPC{client: client}
}

func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error) {
	// 这里负责把应用层参数翻译成 gRPC 请求对象。
	resp, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
	logrus.Info("stock_grpc response", resp)
	return resp, err
}

func (s StockGRPC) GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error) {
	// 对调用方来说这里只是“查商品”，至于底层走 gRPC 还是别的协议都被屏蔽了。
	resp, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{ItemIDs: itemIDs})
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

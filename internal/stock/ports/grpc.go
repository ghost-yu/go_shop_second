package ports

import (
	"context"

	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
	"github.com/ghost-yu/go_shop_second/common/tracing"
	"github.com/ghost-yu/go_shop_second/stock/app"
	"github.com/ghost-yu/go_shop_second/stock/app/query"
	"github.com/ghost-yu/go_shop_second/stock/convertor"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	_, span := tracing.Start(ctx, "GetItems")
	defer span.End()

	items, err := G.app.Queries.GetItems.Handle(ctx, query.GetItems{ItemIDs: request.ItemIDs})
	if err != nil {
		return nil, err
	}
	return &stockpb.GetItemsResponse{Items: items}, nil
}

func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	_, span := tracing.Start(ctx, "CheckIfItemsInStock")
	defer span.End()

	items, err := G.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{
		Items: convertor.NewItemWithQuantityConvertor().ProtosToEntities(request.Items),
	})
	if err != nil {
		return nil, err
	}
	return &stockpb.CheckIfItemsInStockResponse{
		InStock: 1,
		Items:   convertor.NewItemConvertor().EntitiesToProtos(items),
	}, nil
}

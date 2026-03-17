package ports

import (
	context "context"

	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
	"github.com/ghost-yu/go_shop_second/order/app"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GRPCServer 是 orderpb.OrderServiceServer 的端口适配器实现。
// 这一层的职责和 HTTPServer 一样，都是把外部协议请求转给应用层。
type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
	// lesson 当前故意保留 TODO，方便后续逐步实现 gRPC 版本的下单流程。
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	// 这里最终会把 request 转为 query，再返回 proto 层定义的 orderpb.Order。
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) UpdateOrder(ctx context.Context, order *orderpb.Order) (*emptypb.Empty, error) {
	// 更新接口通常会调用 Commands.UpdateOrder；当前阶段先保留占位。
	//TODO implement me
	panic("implement me")
}

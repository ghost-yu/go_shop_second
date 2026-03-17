package order

import (
	"context"
	"fmt"
)

// Repository 定义订单持久化能力，应用层只依赖这个接口，而不关心底层是内存还是数据库。
type Repository interface {
	Create(context.Context, *Order) (*Order, error)
	Get(ctx context.Context, id, customerID string) (*Order, error)
	Update(
		ctx context.Context,
		o *Order,
		updateFn func(context.Context, *Order) (*Order, error),
	) error
}

// NotFoundError 是带业务语义的错误类型，调用方可以明确知道“订单不存在”而不是普通异常。
type NotFoundError struct {
	OrderID string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("order '%s' not found", e.OrderID)
}

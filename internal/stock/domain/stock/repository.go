package stock

import (
	"context"
	"fmt"
	"strings"

	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
)

// Repository 定义库存服务需要提供的最小数据访问能力。
// order 服务只要通过应用层依赖这个接口，就不需要知道库存数据放在哪里。
type Repository interface {
	GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error)
}

// NotFoundError 用来明确告诉调用方哪些商品在库存里不存在。
type NotFoundError struct {
	Missing []string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("these items not found in stock: %s", strings.Join(e.Missing, ","))
}

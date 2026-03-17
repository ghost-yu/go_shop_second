package adapters

import (
	"context"
	"sync"

	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
)

// MemoryStockRepository 用 map 模拟库存数据源。
// 相比切片，map 更适合按商品 ID 做快速查找。
type MemoryStockRepository struct {
	lock  *sync.RWMutex
	store map[string]*orderpb.Item
}

// stub 是教学阶段的假数据，用来模拟一个已经存在的库存商品。
var stub = map[string]*orderpb.Item{
	"item_id": {
		ID:       "foo_item",
		Name:     "stub item",
		Quantity: 10000,
		PriceID:  "stub_item_price_id",
	},
}

func NewMemoryStockRepository() *MemoryStockRepository {
	return &MemoryStockRepository{
		lock:  &sync.RWMutex{},
		store: stub,
	}
}

func (m MemoryStockRepository) GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error) {
	// 这里只有读场景，所以拿读锁即可。
	m.lock.RLock()
	defer m.lock.RUnlock()
	var (
		res     []*orderpb.Item
		missing []string
	)
	for _, id := range ids {
		if item, exist := m.store[id]; exist {
			res = append(res, item)
		} else {
			missing = append(missing, id)
		}
	}
	if len(res) == len(ids) {
		return res, nil
	}
	// 返回“部分结果 + 缺失错误”可以帮助上层更明确地决定如何提示用户。
	return res, domain.NotFoundError{Missing: missing}
}

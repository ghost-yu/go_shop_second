package adapters

import (
	"context"
	"strconv"
	"sync"
	"time"

	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
	"github.com/sirupsen/logrus"
)

// MemoryOrderRepository 是 domain.Repository 的内存版实现。
// 教学项目先用它跑通流程，后面换数据库时只需要替换这一层。
type MemoryOrderRepository struct {
	lock  *sync.RWMutex
	store []*domain.Order
}

func NewMemoryOrderRepository() *MemoryOrderRepository {
	// 这里放一条假数据，方便一开始就能演示“查询已有订单”的路径。
	s := make([]*domain.Order, 0)
	s = append(s, &domain.Order{
		ID:          "fake-ID",
		CustomerID:  "fake-customer-id",
		Status:      "fake-status",
		PaymentLink: "fake-payment-link",
		Items:       nil,
	})
	return &MemoryOrderRepository{
		lock:  &sync.RWMutex{},
		store: s,
	}
}

func (m *MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
	// 写操作要加互斥锁，避免多个请求并发修改切片时产生数据竞争。
	m.lock.Lock()
	defer m.lock.Unlock()
	newOrder := &domain.Order{
		// 当前用 Unix 时间戳凑一个简单 ID，真实项目里通常会换成雪花算法或 UUID。
		ID:          strconv.FormatInt(time.Now().Unix(), 10),
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
		Items:       order.Items,
	}
	m.store = append(m.store, newOrder)
	logrus.WithFields(logrus.Fields{
		"input_order":        order,
		"store_after_create": m.store,
	}).Info("memory_order_repo_create")
	return newOrder, nil
}

func (m *MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*domain.Order, error) {
	for i, v := range m.store {
		logrus.Infof("m.store[%d] = %+v", i, v)
	}
	// 读操作使用读锁，允许多个查询并发进行。
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, o := range m.store {
		if o.ID == id && o.CustomerID == customerID {
			logrus.Infof("memory_order_repo_get||found||id=%s||customerID=%s||res=%+v", id, customerID, *o)
			return o, nil
		}
	}
	return nil, domain.NotFoundError{OrderID: id}
}

func (m *MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
	// UpdateFn 把“怎么改”交给上层，把“在哪里存”留给仓储层，是一种职责分离。
	m.lock.Lock()
	defer m.lock.Unlock()
	found := false
	for i, o := range m.store {
		if o.ID == order.ID && o.CustomerID == order.CustomerID {
			found = true
			updatedOrder, err := updateFn(ctx, o)
			if err != nil {
				return err
			}
			m.store[i] = updatedOrder
		}
	}
	if !found {
		return domain.NotFoundError{OrderID: order.ID}
	}
	return nil
}

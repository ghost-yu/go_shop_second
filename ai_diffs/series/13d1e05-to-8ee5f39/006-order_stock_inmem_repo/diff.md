# Commit Diff Report

- Repo: go_shop_second
- Sequence: 006 / 10
- Commit: f87d59efe0cc19ac7212056bbcbd3f73c733340d
- ShortCommit: f87d59e
- Parent: 49f4e56b544b754e388515fdd7f8b15bf1341e20
- Subject: order stock inmem repo
- Author: ghost-yu <hgfhgfhgfhgfhgfhgf@yeah.net>
- Date: 2024-10-14 01:21:55 +0800
- GeneratedAt: 2026-04-06 17:43:36 +08:00

## Short Summary

~~~text
 5 files changed, 179 insertions(+)
~~~

## File Stats

~~~text
 internal/order/adapters/order_inmem_repository.go | 73 +++++++++++++++++++++++
 internal/order/domain/order/order.go              | 11 ++++
 internal/order/domain/order/repository.go         | 24 ++++++++
 internal/stock/adapters/stock_inmem_repository.go | 50 ++++++++++++++++
 internal/stock/domain/stock/repository.go         | 21 +++++++
 5 files changed, 179 insertions(+)
~~~

## Changed Files

~~~text
internal/order/adapters/order_inmem_repository.go
internal/order/domain/order/order.go
internal/order/domain/order/repository.go
internal/stock/adapters/stock_inmem_repository.go
internal/stock/domain/stock/repository.go
~~~

## Focus Files (Excluded: go.mod / go.sum)

~~~text
internal/order/adapters/order_inmem_repository.go
internal/order/domain/order/order.go
internal/order/domain/order/repository.go
internal/stock/adapters/stock_inmem_repository.go
internal/stock/domain/stock/repository.go
~~~

## Patch

~~~diff
diff --git a/internal/order/adapters/order_inmem_repository.go b/internal/order/adapters/order_inmem_repository.go
new file mode 100644
index 0000000..f5d8be2
--- /dev/null
+++ b/internal/order/adapters/order_inmem_repository.go
@@ -0,0 +1,73 @@
+package adapters
+
+import (
+	"context"
+	"strconv"
+	"sync"
+	"time"
+
+	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
+	"github.com/sirupsen/logrus"
+)
+
+type MemoryOrderRepository struct {
+	lock  *sync.RWMutex
+	store []*domain.Order
+}
+
+func NewMemoryOrderRepository() *MemoryOrderRepository {
+	return &MemoryOrderRepository{
+		lock:  &sync.RWMutex{},
+		store: make([]*domain.Order, 0),
+	}
+}
+
+func (m MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
+	m.lock.Lock()
+	defer m.lock.Unlock()
+	newOrder := &domain.Order{
+		ID:          strconv.FormatInt(time.Now().Unix(), 10),
+		CustomerID:  order.CustomerID,
+		Status:      order.Status,
+		PaymentLink: order.PaymentLink,
+		Items:       order.Items,
+	}
+	m.store = append(m.store, newOrder)
+	logrus.WithFields(logrus.Fields{
+		"input_order":        order,
+		"store_after_create": m.store,
+	}).Debug("memory_order_repo_create")
+	return newOrder, nil
+}
+
+func (m MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*domain.Order, error) {
+	m.lock.RLock()
+	defer m.lock.RUnlock()
+	for _, o := range m.store {
+		if o.ID == id && o.CustomerID == customerID {
+			logrus.Debugf("memory_order_repo_get||found||id=%s||customerID=%s||res=%+v", id, customerID, *o)
+			return o, nil
+		}
+	}
+	return nil, domain.NotFoundError{OrderID: id}
+}
+
+func (m MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
+	m.lock.Lock()
+	defer m.lock.Unlock()
+	found := false
+	for i, o := range m.store {
+		if o.ID == order.ID && o.CustomerID == order.CustomerID {
+			found = true
+			updatedOrder, err := updateFn(ctx, o)
+			if err != nil {
+				return err
+			}
+			m.store[i] = updatedOrder
+		}
+	}
+	if !found {
+		return domain.NotFoundError{OrderID: order.ID}
+	}
+	return nil
+}
diff --git a/internal/order/domain/order/order.go b/internal/order/domain/order/order.go
new file mode 100644
index 0000000..6a9e94b
--- /dev/null
+++ b/internal/order/domain/order/order.go
@@ -0,0 +1,11 @@
+package order
+
+import "github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+
+type Order struct {
+	ID          string
+	CustomerID  string
+	Status      string
+	PaymentLink string
+	Items       []*orderpb.Item
+}
diff --git a/internal/order/domain/order/repository.go b/internal/order/domain/order/repository.go
new file mode 100644
index 0000000..04b783f
--- /dev/null
+++ b/internal/order/domain/order/repository.go
@@ -0,0 +1,24 @@
+package order
+
+import (
+	"context"
+	"fmt"
+)
+
+type Repository interface {
+	Create(context.Context, *Order) (*Order, error)
+	Get(ctx context.Context, id, customerID string) (*Order, error)
+	Update(
+		ctx context.Context,
+		o *Order,
+		updateFn func(context.Context, *Order) (*Order, error),
+	) error
+}
+
+type NotFoundError struct {
+	OrderID string
+}
+
+func (e NotFoundError) Error() string {
+	return fmt.Sprintf("order '%s' not found", e.OrderID)
+}
diff --git a/internal/stock/adapters/stock_inmem_repository.go b/internal/stock/adapters/stock_inmem_repository.go
new file mode 100644
index 0000000..cd0b063
--- /dev/null
+++ b/internal/stock/adapters/stock_inmem_repository.go
@@ -0,0 +1,50 @@
+package adapters
+
+import (
+	"context"
+	"sync"
+
+	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
+)
+
+type MemoryStockRepository struct {
+	lock  *sync.RWMutex
+	store map[string]*orderpb.Item
+}
+
+var stub = map[string]*orderpb.Item{
+	"item_id": {
+		ID:       "foo_item",
+		Name:     "stub item",
+		Quantity: 10000,
+		PriceID:  "stub_item_price_id",
+	},
+}
+
+func NewMemoryOrderRepository() *MemoryStockRepository {
+	return &MemoryStockRepository{
+		lock:  &sync.RWMutex{},
+		store: stub,
+	}
+}
+
+func (m MemoryStockRepository) GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error) {
+	m.lock.RLock()
+	defer m.lock.RUnlock()
+	var (
+		res     []*orderpb.Item
+		missing []string
+	)
+	for _, id := range ids {
+		if item, exist := m.store[id]; exist {
+			res = append(res, item)
+		} else {
+			missing = append(missing, id)
+		}
+	}
+	if len(res) == len(ids) {
+		return res, nil
+	}
+	return res, domain.NotFoundError{Missing: missing}
+}
diff --git a/internal/stock/domain/stock/repository.go b/internal/stock/domain/stock/repository.go
new file mode 100644
index 0000000..7a58e6a
--- /dev/null
+++ b/internal/stock/domain/stock/repository.go
@@ -0,0 +1,21 @@
+package stock
+
+import (
+	"context"
+	"fmt"
+	"strings"
+
+	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
+)
+
+type Repository interface {
+	GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error)
+}
+
+type NotFoundError struct {
+	Missing []string
+}
+
+func (e NotFoundError) Error() string {
+	return fmt.Sprintf("these items not found in stock: %s", strings.Join(e.Missing, ","))
+}
~~~

# Lesson Pair Diff Report

- FromBranch: lesson67
- ToBranch: lesson68

## Short Summary

~~~text
 6 files changed, 143 insertions(+), 26 deletions(-)
~~~

## File Stats

~~~text
 internal/order/app/command/create_order.go   | 36 +++++++---------------
 internal/order/domain/order/aggregate.go     | 34 +++++++++++++++++++++
 internal/order/domain/order/event.go         | 13 ++++++++
 internal/order/domain/service/order.go       | 45 ++++++++++++++++++++++++++++
 internal/order/infrastructure/mq/rabbitmq.go | 37 +++++++++++++++++++++++
 internal/order/service/application.go        |  4 ++-
 6 files changed, 143 insertions(+), 26 deletions(-)
~~~

## Commit Comparison

~~~text
> 24561fe aggregate root
~~~

## Changed Files

~~~text
internal/order/app/command/create_order.go
internal/order/domain/order/aggregate.go
internal/order/domain/order/event.go
internal/order/domain/service/order.go
internal/order/infrastructure/mq/rabbitmq.go
internal/order/service/application.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
internal/order/app/command/create_order.go
internal/order/domain/order/aggregate.go
internal/order/domain/order/event.go
internal/order/domain/service/order.go
internal/order/infrastructure/mq/rabbitmq.go
internal/order/service/application.go
~~~

## Full Diff

~~~diff
diff --git a/internal/order/app/command/create_order.go b/internal/order/app/command/create_order.go
index df86d4a..e4d27aa 100644
--- a/internal/order/app/command/create_order.go
+++ b/internal/order/app/command/create_order.go
@@ -11,8 +11,8 @@ import (
 	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/ghost-yu/go_shop_second/order/app/query"
 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
+	"github.com/ghost-yu/go_shop_second/order/domain/service"
 	"github.com/pkg/errors"
-	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
 	"go.opentelemetry.io/otel"
 	"google.golang.org/grpc/status"
@@ -30,26 +30,26 @@ type CreateOrderResult struct {
 type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]
 
 type createOrderHandler struct {
-	orderRepo domain.Repository
-	stockGRPC query.StockService
-	channel   *amqp.Channel
+	orderRepo      domain.Repository
+	stockGRPC      query.StockService
+	eventPublisher domain.EventPublisher
 }
 
-func NewCreateOrderHandler(orderRepo domain.Repository, stockGRPC query.StockService, channel *amqp.Channel, logger *logrus.Logger, metricClient decorator.MetricsClient) CreateOrderHandler {
+func NewCreateOrderHandler(orderRepo domain.Repository, stockGRPC query.StockService, eventPublisher domain.EventPublisher, logger *logrus.Logger, metricClient decorator.MetricsClient) CreateOrderHandler {
 	if orderRepo == nil {
 		panic("nil orderRepo")
 	}
 	if stockGRPC == nil {
 		panic("nil stockGRPC")
 	}
-	if channel == nil {
-		panic("nil channel ")
+	if eventPublisher == nil {
+		panic("nil eventPublisher")
 	}
 	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
 		createOrderHandler{
-			orderRepo: orderRepo,
-			stockGRPC: stockGRPC,
-			channel:   channel,
+			orderRepo:      orderRepo,
+			stockGRPC:      stockGRPC,
+			eventPublisher: eventPublisher,
 		},
 		logger,
 		metricClient,
@@ -72,22 +72,8 @@ func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*Creat
 	if err != nil {
 		return nil, err
 	}
-	o, err := c.orderRepo.Create(ctx, pendingOrder)
-	if err != nil {
-		return nil, err
-	}
-
-	err = broker.PublishEvent(ctx, broker.PublishEventReq{
-		Channel:  c.channel,
-		Routing:  broker.Direct,
-		Queue:    broker.EventOrderCreated,
-		Exchange: "",
-		Body:     o,
-	})
-	if err != nil {
-		return nil, errors.Wrapf(err, "publish event error q.Name=%s", broker.EventOrderCreated)
-	}
 
+	o, err := service.NewOrderDomainService(c.orderRepo, c.eventPublisher).CreateOrder(ctx, *pendingOrder)
 	return &CreateOrderResult{OrderID: o.ID}, nil
 }
 
diff --git a/internal/order/domain/order/aggregate.go b/internal/order/domain/order/aggregate.go
new file mode 100644
index 0000000..b8dcebd
--- /dev/null
+++ b/internal/order/domain/order/aggregate.go
@@ -0,0 +1,34 @@
+package order
+
+import "github.com/pkg/errors"
+
+type Identity struct {
+	CustomerID string
+	OrderID    string
+}
+
+type AggregateRoot struct {
+	Identity    Identity
+	OrderEntity *Order
+}
+
+func NewAggregateRoot(identity Identity, orderEntity *Order) *AggregateRoot {
+	return &AggregateRoot{Identity: identity, OrderEntity: orderEntity}
+}
+
+func (a *AggregateRoot) BusinessIdentity() Identity {
+	return Identity{
+		CustomerID: a.OrderEntity.CustomerID,
+		OrderID:    a.OrderEntity.ID,
+	}
+}
+
+func (a *AggregateRoot) Validate() error {
+	if a.Identity.OrderID == "" || a.Identity.CustomerID == "" {
+		return errors.New("invalid identity")
+	}
+	if a.OrderEntity != nil {
+		return errors.New("empty order")
+	}
+	return nil
+}
diff --git a/internal/order/domain/order/event.go b/internal/order/domain/order/event.go
new file mode 100644
index 0000000..e4bbf42
--- /dev/null
+++ b/internal/order/domain/order/event.go
@@ -0,0 +1,13 @@
+package order
+
+import "context"
+
+type DomainEvent struct {
+	Dest string
+	Data any
+}
+
+type EventPublisher interface {
+	Publish(ctx context.Context, event DomainEvent) error
+	Broadcast(ctx context.Context, event DomainEvent) error
+}
diff --git a/internal/order/domain/service/order.go b/internal/order/domain/service/order.go
new file mode 100644
index 0000000..d6b729b
--- /dev/null
+++ b/internal/order/domain/service/order.go
@@ -0,0 +1,45 @@
+package service
+
+import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/common/broker"
+	"github.com/ghost-yu/go_shop_second/common/entity"
+	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
+	"github.com/pkg/errors"
+)
+
+type OrderDomainService struct {
+	Repo           domain.Repository
+	EventPublisher domain.EventPublisher
+}
+
+func NewOrderDomainService(repo domain.Repository, eventPublisher domain.EventPublisher) *OrderDomainService {
+	return &OrderDomainService{Repo: repo, EventPublisher: eventPublisher}
+}
+
+func (s *OrderDomainService) CreateOrder(ctx context.Context, order domain.Order) (res *entity.Order, err error) {
+	root := domain.NewAggregateRoot(
+		domain.Identity{CustomerID: order.CustomerID, OrderID: order.ID},
+		&order,
+	)
+	o, err := s.Repo.Create(ctx, root.OrderEntity)
+	if err != nil {
+		return nil, err
+	}
+
+	if err = s.EventPublisher.Publish(ctx, domain.DomainEvent{
+		Dest: broker.EventOrderCreated,
+		Data: o,
+	}); err != nil {
+		return nil, errors.Wrapf(err, "publish event error q.Name=%s", broker.EventOrderCreated)
+	}
+
+	return &entity.Order{
+		ID:          o.ID,
+		CustomerID:  o.CustomerID,
+		Status:      o.Status,
+		PaymentLink: o.PaymentLink,
+		Items:       o.Items,
+	}, nil
+}
diff --git a/internal/order/infrastructure/mq/rabbitmq.go b/internal/order/infrastructure/mq/rabbitmq.go
new file mode 100644
index 0000000..6c83aa2
--- /dev/null
+++ b/internal/order/infrastructure/mq/rabbitmq.go
@@ -0,0 +1,37 @@
+package mq
+
+import (
+	"context"
+
+	"github.com/ghost-yu/go_shop_second/common/broker"
+	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
+	amqp "github.com/rabbitmq/amqp091-go"
+)
+
+type RabbitMQEventPublisher struct {
+	Channel *amqp.Channel
+}
+
+func NewRabbitMQEventPublisher(channel *amqp.Channel) *RabbitMQEventPublisher {
+	return &RabbitMQEventPublisher{Channel: channel}
+}
+
+func (r *RabbitMQEventPublisher) Publish(ctx context.Context, event domain.DomainEvent) error {
+	return broker.PublishEvent(ctx, broker.PublishEventReq{
+		Channel:  r.Channel,
+		Routing:  broker.Direct,
+		Queue:    event.Dest,
+		Exchange: "",
+		Body:     event.Data,
+	})
+}
+
+func (r *RabbitMQEventPublisher) Broadcast(ctx context.Context, event domain.DomainEvent) error {
+	return broker.PublishEvent(ctx, broker.PublishEventReq{
+		Channel:  r.Channel,
+		Routing:  broker.FanOut,
+		Queue:    event.Dest,
+		Exchange: "",
+		Body:     event.Data,
+	})
+}
diff --git a/internal/order/service/application.go b/internal/order/service/application.go
index 2879b67..d4fdf41 100644
--- a/internal/order/service/application.go
+++ b/internal/order/service/application.go
@@ -13,6 +13,7 @@ import (
 	"github.com/ghost-yu/go_shop_second/order/app"
 	"github.com/ghost-yu/go_shop_second/order/app/command"
 	"github.com/ghost-yu/go_shop_second/order/app/query"
+	"github.com/ghost-yu/go_shop_second/order/infrastructure/mq"
 	amqp "github.com/rabbitmq/amqp091-go"
 	"github.com/sirupsen/logrus"
 	"github.com/spf13/viper"
@@ -45,9 +46,10 @@ func newApplication(_ context.Context, stockGRPC query.StockService, ch *amqp.Ch
 	mongoClient := newMongoClient()
 	orderRepo := adapters.NewOrderRepositoryMongo(mongoClient)
 	metricClient := metrics.TodoMetrics{}
+	eventPublisher := mq.NewRabbitMQEventPublisher(ch)
 	return app.Application{
 		Commands: app.Commands{
-			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, ch, logrus.StandardLogger(), metricClient),
+			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, eventPublisher, logrus.StandardLogger(), metricClient),
 			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logrus.StandardLogger(), metricClient),
 		},
 		Queries: app.Queries{
~~~

package order

import "github.com/pkg/errors"

type Identity struct {
	CustomerID string
	OrderID    string
}

type AggregateRoot struct {
	Identity    Identity
	OrderEntity *Order
}

func NewAggregateRoot(identity Identity, orderEntity *Order) *AggregateRoot {
	return &AggregateRoot{Identity: identity, OrderEntity: orderEntity}
}

func (a *AggregateRoot) BusinessIdentity() Identity {
	return Identity{
		CustomerID: a.OrderEntity.CustomerID,
		OrderID:    a.OrderEntity.ID,
	}
}

func (a *AggregateRoot) Validate() error {
	if a.Identity.OrderID == "" || a.Identity.CustomerID == "" {
		return errors.New("invalid identity")
	}
	if a.OrderEntity != nil {
		return errors.New("empty order")
	}
	return nil
}

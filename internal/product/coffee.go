package product

import (
	"coffy/internal/event"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Coffee struct {
	AggregateID string
	Type        string
	price       float64
	events      []event.Event
}

func (b *Coffee) Events() []event.Event {
	return b.events
}

func (b *Coffee) Price() float64 {
	return b.price
}

// ChangePrice updates the price of the current Coffee.
//
// Only values greater or equal zero are allowed.
func (b *Coffee) ChangePrice(p float64, reason string) error {
	if p <= 0 {
		return errors.New("invalid price")
	}
	e := NewPriceUpdated(b.AggregateID, p, reason)
	if err := b.apply(*e); err != nil {
		return errors.Join(
			fmt.Errorf("could not change price for %s [id: %s]",
				b.Type, b.AggregateID), err)
	}
	return nil
}

// Clear empties the current event cache of the coffee and removes all previously appended events.
func (b *Coffee) Clear() {
	b.events = []event.Event{}
}

// Load sets the state of the current coffee by applying all events iteratively.
//
// After all events have been applied to the account, the event cache is emptied.
func (b *Coffee) Load(events []event.Event) error {
	for _, e := range events {
		if err := b.apply(e); err != nil {
			return fmt.Errorf("could not apply event: %v", err)
		}
	}
	b.Clear()
	return nil
}

func (b *Coffee) apply(e event.Event) error {
	switch theEvent := e.(type) {
	case CoffeeCreated:
		b.applyCreated(theEvent)
	case PriceUpdated:
		if err := b.applyNewPrice(theEvent); err != nil {
			return err
		}
	default:
		return fmt.Errorf("cannot apply event: unknown event '%T'", e)
	}
	return nil
}

func (b *Coffee) applyNewPrice(e PriceUpdated) error {
	if e.AggregateID() != b.AggregateID {
		return fmt.Errorf("beverage ids do not match: expected %s, actual %s", b.AggregateID, e.AggregateID())
	}
	b.price = e.Price
	b.events = append(b.events, e)
	return nil
}

func (b *Coffee) applyCreated(e CoffeeCreated) {
	b.AggregateID = e.ID
	b.Type = e.BeverageType
	b.price = e.Price
	b.events = append(b.events, e)
}

type CoffeeCreated struct {
	ID           string
	BeverageType string
	Price        float64
	OccurredOn   time.Time
}

func NewCoffee(coffeeType string, price float64) (*Coffee, error) {
	if coffeeType == "" {
		return nil, errors.New("beverage type cannot be empty")
	}
	if price <= 0 {
		return nil, errors.New("price must be greater than zero")
	}
	beverage := &Coffee{}
	created := NewCoffeeCreated(uuid.NewString(), coffeeType, price)
	if err := beverage.apply(*created); err != nil {
		return nil, err
	}
	return beverage, nil
}

func NewCoffeeCreated(id string, coffeeType string, price float64) *CoffeeCreated {
	return &CoffeeCreated{id, coffeeType, price, time.Now()}
}

func (b CoffeeCreated) AggregateID() string {
	return b.ID
}

func (b CoffeeCreated) Type() string {
	return "CoffeeCreated"
}

func (b CoffeeCreated) Occurred() time.Time {
	return b.OccurredOn
}

type PriceUpdated struct {
	ID         string
	Price      float64
	Reason     string
	OccurredOn time.Time
}

func NewPriceUpdated(aggregateID string, price float64, reason string) *PriceUpdated {
	return &PriceUpdated{aggregateID, price, reason, time.Now()}
}

func (e PriceUpdated) AggregateID() string {
	return e.ID
}

func (e PriceUpdated) Type() string {
	return "PriceUpdated"
}

func (e PriceUpdated) Occurred() time.Time {
	return e.OccurredOn
}

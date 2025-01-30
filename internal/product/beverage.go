package product

import (
	"coffy/internal/event"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Beverage struct {
	AggregateID  string
	BeverageType string
	price        float64
	events       []event.Event
}

func (b *Beverage) Events() []event.Event {
	return b.events
}

func (b *Beverage) Price() float64 {
	return b.price
}

// ChangePrice updates the price of the current Beverage.
//
// Only values greater or equal zero are allowed.
func (b *Beverage) ChangePrice(p float64, reason string) error {
	if p <= 0 {
		return errors.New("invalid price")
	}
	e := NewPriceUpdated(b.AggregateID, p, reason)
	if err := b.Apply(*e); err != nil {
		return errors.Join(
			fmt.Errorf("could not change price for %s [id: %s]",
				b.BeverageType, b.AggregateID), err)
	}
	return nil
}

func (b *Beverage) Apply(e event.Event) error {
	switch theEvent := e.(type) {
	case BeverageCreated:
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

func (b *Beverage) applyNewPrice(e PriceUpdated) error {
	if e.AggregateID() != b.AggregateID {
		return fmt.Errorf("beverage ids do not match: expected %s, actual %s", b.AggregateID, e.AggregateID())
	}
	b.price = e.Price
	b.events = append(b.events, e)
	return nil
}

func (b *Beverage) applyCreated(e BeverageCreated) {
	b.AggregateID = e.ID
	b.BeverageType = e.BeverageType
	b.price = e.Price
	b.events = append(b.events, e)
}

type BeverageCreated struct {
	ID           string
	BeverageType string
	Price        float64
	OccurredOn   time.Time
}

func NewBeverage(beverageType string, price float64) (*Beverage, error) {
	if beverageType == "" {
		return nil, errors.New("beverage type cannot be empty")
	}
	if price <= 0 {
		return nil, errors.New("price must be greater than zero")
	}
	beverage := &Beverage{}
	created := NewBeverageCreated(uuid.NewString(), beverageType, price)
	if err := beverage.Apply(*created); err != nil {
		return nil, err
	}
	return beverage, nil
}

func NewBeverageCreated(id string, beverageType string, price float64) *BeverageCreated {
	return &BeverageCreated{id, beverageType, price, time.Now()}
}

func (b BeverageCreated) AggregateID() string {
	return b.ID
}

func (b BeverageCreated) Type() string {
	return "BeverageCreated"
}

func (b BeverageCreated) Occurred() time.Time {
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

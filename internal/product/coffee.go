package product

import (
	"coffy/internal/event"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// Coffee is the actual representation of the black gold.
type Coffee struct {
	AggregateID string        // The ID to identify the coffee in the system
	Type        string        // What type of coffee it is, e.g. black coffee, espresso, ...
	price       float64       // The current price of the coffee in €, e.g. 0.50 for 50 cents
	cva         CuppingScore  // The coffee value assessment result, currently the CuppingScore
	details     Details       // A more detailed description about the coffee
	events      []event.Event // Uncommitted events of the aggregate
}

// CoffeeValue provides the assessed Value of the current coffee.
// It is represented by a CuppingScore, which is the SCA standard metric
// after a coffee Value assessment (CVA).
func (c *Coffee) CoffeeValue() *CuppingScore {
	return &c.cva
}

// Price returns the latest price of the coffee.
func (c *Coffee) Price() float64 {
	return c.price
}

// Events returns all uncommitted events of the current coffee aggregate
func (c *Coffee) Events() []event.Event {
	return c.events
}

// CuppingScore holds an assessment Value resulting from a coffee Value assessment, which is between 58 and 100 (both inclusive).
// A score of 58 is the lowest Value possible, 100 the highest and represents the best coffee Value possible.
type CuppingScore struct {
	Value int
}

func newCuppingScore(value int) (*CuppingScore, error) {
	if value < 58 || value > 100 {
		return nil, errors.New("invalid cupping score")
	}
	return &CuppingScore{Value: value}, nil
}

// SetCuppingScore sets a standardized sensor-based quality metric that resulted from a coffee Value assessment (CVA).
// The SCA normalised the score to be between an inclusive range of 58 (worst) and 100 (best). Values outside this range
// will result in an error.
func (c *Coffee) SetCuppingScore(value int) error {
	score, err := newCuppingScore(value)
	if err != nil {
		return err
	}
	e := newCvaProvided(c.AggregateID, score.Value)
	if err := c.apply(e); err != nil {
		return errors.Join(
			fmt.Errorf("could not set cupping score for %s [id: %s]",
				c.Type, c.AggregateID), err)
	}
	return nil
}

func newCvaProvided(id string, value int) CvaProvided {
	return CvaProvided{ID: id, Value: value, OccurredOn: time.Now()}
}

// CvaProvided an event record about change in a coffee's Value assessment.
//
// In this case it is the CuppingScore, represented as a simple integer Value for the event.
type CvaProvided struct {
	ID         string
	Value      int
	OccurredOn time.Time
}

func (e CvaProvided) AggregateID() string {
	return e.ID
}

func (e CvaProvided) Type() string {
	return "CvaProvided"
}

func (e CvaProvided) Occurred() time.Time {
	return e.OccurredOn
}

// Details enables a containerised description with more detail of a Coffee.
type Details struct {
	Origin      string            `json:"origin"`      // The country the coffee has been produced
	Description string            `json:"description"` // Some detailed description about the coffee
	RoastHouse  string            `json:"roast_house"` // The location the coffee has been roasted
	Misc        map[string]string `json:"misc"`        // An unstructured collection of key:values to provide more details
}

func (c *Coffee) Details() Details {
	return c.details
}

// UpdateDetails sets some more detailed information for the current Coffee.
func (c *Coffee) UpdateDetails(details Details) error {
	e := NewDetailsUpdated(c.AggregateID, details)
	if err := c.apply(e); err != nil {
		return errors.Join(fmt.Errorf("could not update details for %s [id: %s]", c.Type, c.AggregateID), err)
	}
	return nil
}

type DetailsUpdated struct {
	ID         string
	Details    Details
	OccurredOn time.Time
}

func (e DetailsUpdated) AggregateID() string {
	return e.ID
}

func (e DetailsUpdated) Type() string {
	return "DetailsUpdated"
}

func (e DetailsUpdated) Occurred() time.Time {
	return e.OccurredOn
}

func NewDetailsUpdated(id string, details Details) DetailsUpdated {
	return DetailsUpdated{ID: id, Details: details, OccurredOn: time.Now()}
}

// ChangePrice updates the price of the current Coffee.
//
// Only values greater or equal zero are allowed.
func (c *Coffee) ChangePrice(p float64, reason string) error {
	if p <= 0 {
		return errors.New("invalid price")
	}
	e := NewPriceUpdated(c.AggregateID, p, reason)
	if err := c.apply(*e); err != nil {
		return errors.Join(
			fmt.Errorf("could not change price for %s [id: %s]",
				c.Type, c.AggregateID), err)
	}
	return nil
}

// Clear empties the current event cache of the coffee and removes all previously appended events.
func (c *Coffee) Clear() {
	c.events = []event.Event{}
}

// Load sets the state of the current coffee by applying all events iteratively.
//
// After all events have been applied to the account, the event cache is emptied.
func (c *Coffee) Load(events []event.Event) error {
	for _, e := range events {
		if err := c.apply(e); err != nil {
			return fmt.Errorf("could not apply event: %v", err)
		}
	}
	c.Clear()
	return nil
}

func (c *Coffee) apply(e event.Event) error {
	switch theEvent := e.(type) {
	case CoffeeCreated:
		c.applyCreated(theEvent)
	case PriceUpdated:
		if err := c.applyNewPrice(theEvent); err != nil {
			return err
		}
	case CvaProvided:
		if err := c.applyCva(theEvent); err != nil {
			return err
		}
	case DetailsUpdated:
		if err := c.applyDetails(theEvent); err != nil {
			return err
		}
	default:
		return fmt.Errorf("cannot apply event: unknown event '%T'", e)
	}
	return nil
}

func (c *Coffee) applyNewPrice(e PriceUpdated) error {
	if e.AggregateID() != c.AggregateID {
		return fmt.Errorf("coffee ids do not match: expected %s, actual %s", c.AggregateID, e.AggregateID())
	}
	c.price = e.Price
	c.events = append(c.events, e)
	return nil
}

func (c *Coffee) applyCreated(e CoffeeCreated) {
	c.AggregateID = e.ID
	c.Type = e.BeverageType
	c.price = e.Price
	c.events = append(c.events, e)
}

func (c *Coffee) applyCva(e CvaProvided) error {
	if e.AggregateID() != c.AggregateID {
		return fmt.Errorf("coffee ids do not match: expected %s, actual %s", c.AggregateID, e.AggregateID())
	}
	score, err := newCuppingScore(e.Value)
	if err != nil {
		return err
	}
	c.cva = *score
	c.events = append(c.events, e)
	return nil
}

func (c *Coffee) applyDetails(e DetailsUpdated) error {
	c.details = e.Details
	c.events = append(c.events, e)
	return nil
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

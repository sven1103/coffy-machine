package account

import (
	"coffy/internal/event"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

type Account struct {
	id      string
	owner   string
	balance float64
	events  []event.Event
}

func NewAccount(owner string) (*Account, error) {
	created := newAccountCreated(uuid.New().String(), time.Now(), owner)
	a := Account{}
	if err := a.Apply(*created); err != nil {
		return nil, err
	}
	return &a, nil
}

func (a *Account) Apply(e event.Event) error {
	switch theEvent := e.(type) {
	case coffyConsumed:
		return a.applyConsumed(theEvent)
	case accountCreated:
		return a.createAccount(theEvent)
	case incomingPayment:
		return a.applyPayment(theEvent)
	default:
		return fmt.Errorf("unknown event: %v", e)
	}
}

func (a *Account) applyConsumed(e coffyConsumed) error {
	if e.AggregateID() != a.id {
		return fmt.Errorf("event aggregate id does not match current aggregate")
	}
	a.balance -= e.Costs
	a.events = append(a.events, e)
	return nil
}

func (a *Account) Consume(price float64, coffeeType string) error {
	if price < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	e := newCoffyConsumed(a.id, coffeeType, price)
	if err := a.Apply(*e); err != nil {
		return err
	}
	return nil
}

func (a *Account) Pay(amount float64, reason string) error {
	if amount < 0 {
		return fmt.Errorf("payment amount cannot be negative")
	}
	e := newIncomingPayment(a.id, amount, reason)
	if err := a.Apply(*e); err != nil {
		log.Printf("Error: %v", err)
		return fmt.Errorf("error paying %.2f to Account ID '%s'", amount, a.id)
	}
	return nil
}

func (a *Account) ID() string {
	return a.id
}

func (a *Account) Owner() string {
	return a.owner
}

func (a *Account) Balance() float64 {
	return a.balance
}

func (a *Account) Events() []event.Event {
	return a.events
}

func (a *Account) createAccount(e accountCreated) error {
	if a.owner != "" {
		return fmt.Errorf("Account already exists")
	}
	a.owner = e.Owner
	a.id = e.AggregateID()
	a.events = append(a.events, e)
	return nil
}

func (a *Account) applyPayment(e incomingPayment) error {
	if a.id != e.AggregateID() {
		return fmt.Errorf("event aggregate id does not match current aggregate")
	}
	a.balance += e.Amount
	a.events = append(a.events, e)
	return nil
}

type accountCreated struct {
	AccountID  string
	OccurredOn time.Time
	EventType  string
	Owner      string
}

func newAccountCreated(accountID string, occurredOn time.Time, owner string) *accountCreated {
	return &accountCreated{AccountID: accountID, OccurredOn: occurredOn, EventType: "AccountCreated", Owner: owner}
}

func (e accountCreated) Type() string {
	return e.EventType
}

func (e accountCreated) Occurred() time.Time {
	return e.OccurredOn
}

func (e accountCreated) AggregateID() string {
	return e.AccountID
}

// The coffyConsumed event records a coffee consumption event. Next to the common properties of Event, it
// also records the coffee type (CoffyType) that has been consumed to increase the transparency and the
// associated costs (Costs).
type coffyConsumed struct {
	AccountID  string
	OccurredOn time.Time
	eventType  string
	CoffyType  string
	Costs      float64
}

func newCoffyConsumed(accountID string, coffyType string, costs float64) *coffyConsumed {
	return &coffyConsumed{accountID, time.Now(), "CoffyConsumed", coffyType, costs}
}

// The incomingPayment event records an effort to pay someone's outstanding coffy debts. Next to the common properties
// of Event, it also records a reason (Reason), e.g. one just paid to the bank of coffy, or purchased
// some maintenance materials, like descaling agent etc. The amount (Amount) represents the amount of money that
// has been used to pay to the Account.
type incomingPayment struct {
	AccountID  string
	OccurredOn time.Time
	EventType  string
	Amount     float64
	Reason     string
}

func newIncomingPayment(accountID string, amount float64, reason string) *incomingPayment {
	return &incomingPayment{accountID, time.Now(), "IncomingPayment", amount, reason}
}

func (e coffyConsumed) AggregateID() string {
	return e.AccountID
}

func (e coffyConsumed) Occurred() time.Time {
	return e.OccurredOn
}

func (e coffyConsumed) Type() string {
	return e.eventType
}

func (e incomingPayment) AggregateID() string {
	return e.AccountID
}

func (e incomingPayment) Occurred() time.Time {
	return e.OccurredOn
}

func (e incomingPayment) Type() string {
	return e.EventType
}

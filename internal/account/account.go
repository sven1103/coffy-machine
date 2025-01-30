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
	created := NewAccountCreated(uuid.New().String(), time.Now(), owner)
	a := Account{}
	if err := a.Apply(*created); err != nil {
		return nil, err
	}
	return &a, nil
}

// Apply takes an event.Event and applies it to the current account.
//
// The function is intended to be used only for object creation. For interactions
// with the account, use any other public function (e.g. Consume or Pay).
func (a *Account) Apply(e event.Event) error {
	switch theEvent := e.(type) {
	case CoffyConsumed:
		return a.applyConsumed(theEvent)
	case AccountCreated:
		return a.createAccount(theEvent)
	case IncomingPayment:
		return a.applyPayment(theEvent)
	default:
		return fmt.Errorf("unknown event: %v", e)
	}
}

func (a *Account) applyConsumed(e CoffyConsumed) error {
	if e.AggregateID() != a.id {
		return fmt.Errorf("event aggregate id does not match current aggregate")
	}
	a.balance -= e.Costs
	a.events = append(a.events, e)
	return nil
}

// Consume charges the account with the price of the coffee consumed
// and stores the type of coffee.
//
// The value for the price must be greater or equal zero.
func (a *Account) Consume(price float64, coffeeType string) error {
	if price < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	e := NewCoffyConsumed(a.id, coffeeType, price)
	if err := a.Apply(*e); err != nil {
		return err
	}
	return nil
}

// ConsumeN charges the account with the price of multiple coffee (n) consumed and records
// the type of coffee.
//
// Note: this function calls Consume n times, so every coffee will be a single recorded event.
//
// The value for the price must be greater or equal zero
// and the number of coffee must be greater or equal zero.
func (a *Account) ConsumeN(price float64, coffeeType string, n int) error {
	if n < 0 {
		return fmt.Errorf("n must be positive")
	}
	for range n {
		if err := a.Consume(price, coffeeType); err != nil {
			return err
		}
	}
	return nil
}

// Pay balances the account with a given amount and reason to the account.
// The payment is a deposit to the account and the reason serves as semantic context of the payment.
//
// Only values greater or equal 0 are allowed. The value for reason can be left empty if not required.
func (a *Account) Pay(amount float64, reason string) error {
	if amount < 0 {
		return fmt.Errorf("payment amount cannot be negative")
	}
	e := NewIncomingPayment(a.id, amount, reason)
	if err := a.Apply(*e); err != nil {
		log.Printf("Error: %v", err)
		return fmt.Errorf("error paying %.2f to Account ID '%s'", amount, a.id)
	}
	return nil
}

// ConsumedTotal returns the total amount of coffee consumed.
// Only events of type CoffyConsumed are considered.
func (a *Account) ConsumedTotal() int {
	consumed := 0
	for _, e := range a.events {
		if e.Type() == "CoffyConsumed" {
			consumed += 1
		}
	}
	return consumed
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

func (a *Account) createAccount(e AccountCreated) error {
	if a.owner != "" {
		return fmt.Errorf("Account already exists")
	}
	a.owner = e.Owner
	a.id = e.AggregateID()
	a.events = append(a.events, e)
	return nil
}

func (a *Account) applyPayment(e IncomingPayment) error {
	if a.id != e.AggregateID() {
		return fmt.Errorf("event aggregate id does not match current aggregate")
	}
	a.balance += e.Amount
	a.events = append(a.events, e)
	return nil
}

type AccountCreated struct {
	AccountID  string    `json:"accountID"`
	OccurredOn time.Time `json:"occurredOn"`
	EventType  string    `json:"eventType"`
	Owner      string    `json:"owner"`
}

func NewAccountCreated(accountID string, occurredOn time.Time, owner string) *AccountCreated {
	return &AccountCreated{AccountID: accountID, OccurredOn: occurredOn, EventType: "AccountCreated", Owner: owner}
}

func (e AccountCreated) Type() string {
	return e.EventType
}

func (e AccountCreated) Occurred() time.Time {
	return e.OccurredOn
}

func (e AccountCreated) AggregateID() string {
	return e.AccountID
}

// The CoffyConsumed event records a coffee consumption event. Next to the common properties of Event, it
// also records the coffee type (CoffyType) that has been consumed to increase the transparency and the
// associated costs (Costs).
type CoffyConsumed struct {
	AccountID  string    `json:"accountID"`
	OccurredOn time.Time `json:"occurredOn"`
	EventType  string    `json:"eventType"`
	CoffyType  string    `json:"coffyType"`
	Costs      float64   `json:"costs"`
}

func NewCoffyConsumed(accountID string, coffyType string, costs float64) *CoffyConsumed {
	return &CoffyConsumed{accountID, time.Now(), "CoffyConsumed", coffyType, costs}
}

// The IncomingPayment event records an effort to pay someone's outstanding coffy debts. Next to the common properties
// of Event, it also records a reason (Reason), e.g. one just paid to the bank of coffy, or purchased
// some maintenance materials, like descaling agent etc. The amount (Amount) represents the amount of money that
// has been used to pay to the Account.
type IncomingPayment struct {
	AccountID  string    `json:"accountID"`
	OccurredOn time.Time `json:"occurredOn"`
	EventType  string    `json:"eventType"`
	Amount     float64   `json:"amount"`
	Reason     string    `json:"reason"`
}

func NewIncomingPayment(accountID string, amount float64, reason string) *IncomingPayment {
	return &IncomingPayment{accountID, time.Now(), "IncomingPayment", amount, reason}
}

func (e CoffyConsumed) AggregateID() string {
	return e.AccountID
}

func (e CoffyConsumed) Occurred() time.Time {
	return e.OccurredOn
}

func (e CoffyConsumed) Type() string {
	return e.EventType
}

func (e IncomingPayment) AggregateID() string {
	return e.AccountID
}

func (e IncomingPayment) Occurred() time.Time {
	return e.OccurredOn
}

func (e IncomingPayment) Type() string {
	return e.EventType
}

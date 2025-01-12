package account

import (
	"coffy/internal/event"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type account struct {
	owner   string
	balance float64
	events  []event.Event
}

func newAccount(owner string) (*account, error) {
	created := newAccountCreated(uuid.New().String(), time.Now(), owner)
	a := account{}
	if err := a.Apply(*created); err != nil {
		return nil, err
	}
	return &a, nil
}

func (a *account) Apply(e event.Event) error {
	switch theEvent := e.(type) {
	case accountCreated:
		if err := a.createAccount(theEvent); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unknown event: %v", e)
	}
}

func (a *account) Owner() string {
	return a.owner
}

func (a *account) Events() []event.Event {
	return a.events
}

func (a *account) createAccount(e accountCreated) error {
	if a.owner != "" {
		return fmt.Errorf("account already exists")
	}
	a.owner = e.owner
	a.events = append(a.events, e)
	return nil
}

type accountCreated struct {
	accountID  string
	occurredOn time.Time
	eventType  string
	owner      string
}

func newAccountCreated(accountID string, occurredOn time.Time, owner string) *accountCreated {
	return &accountCreated{accountID: accountID, occurredOn: occurredOn, eventType: "AccountCreated", owner: owner}
}

func (e accountCreated) Type() string {
	return e.eventType
}

func (e accountCreated) OccurredOn() time.Time {
	return e.occurredOn
}

func (e accountCreated) AggregateID() string {
	return e.accountID
}

// The coffyConsumed event records a coffee consumption event. Next to the common properties of Event, it
// also records the coffee type (CoffyType) that has been consumed to increase the transparency and the
// associated costs (Costs).
type coffyConsumed struct {
	accountID  string
	occurredOn time.Time
	CoffyType  string
	Costs      float64
}

func newCoffyConsumed(accountID string, coffyType string, costs float64) *coffyConsumed {
	return &coffyConsumed{accountID, time.Now(), coffyType, costs}
}

// The incomingPayment event records an effort to pay someone's outstanding coffy debts. Next to the common properties
// of Event, it also records a reason (Reason), e.g. one just paid to the bank of coffy, or purchased
// some maintenance materials, like descaling agent etc. The amount (Amount) represents the amount of money that
// has been used to pay to the account.
type incomingPayment struct {
	accountID  string
	occurredOn time.Time
	Amount     float64
	Reason     string
}

func newIncomingPayment(accountID string, amount float64, reason string) *incomingPayment {
	return &incomingPayment{accountID, time.Now(), amount, reason}
}

func (e coffyConsumed) AggregateID() string {
	return e.accountID
}

func (e coffyConsumed) OccurredOn() time.Time {
	return e.occurredOn
}

func (e incomingPayment) AggregateID() string {
	return e.accountID
}

func (e incomingPayment) OccurredOn() time.Time {
	return e.occurredOn
}

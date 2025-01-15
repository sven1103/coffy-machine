package account

import (
	"coffy/internal/event"
	"coffy/internal/storage"
	"encoding/json"
	"fmt"
)

type Accounting struct {
	repo storage.EventRepository
}

func (a *Accounting) Find(accountID string) (*Account, error) {
	query, err := a.repo.LoadAll(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to load account: %w", err)
	}
	events := make([]event.Event, 0)
	for _, entry := range query {
		evnt, err := convert(entry)
		if err != nil {
			return nil, fmt.Errorf("failed to convert event: %w", err)
		}
		events = append(events, evnt)
	}
	acc := &Account{}
	for _, e := range events {
		if err := acc.Apply(e); err != nil {
			return nil, fmt.Errorf("failed to apply event: %w", err)
		}
	}
	return acc, nil
}

func convert(entry storage.EventEntry) (event.Event, error) {
	switch entry.EventType {
	case "AccountCreated":
		evnt, err := toCreate(entry)
		if err != nil {
			return nil, err
		}
		return evnt, nil
	case "IncomingPayment":
		evnt, err := toPayment(entry)
		if err != nil {
			return nil, err
		}
		return evnt, nil
	case "CoffyConsumed":
		evnt, err := toConsume(entry)
		if err != nil {
			return nil, err
		}
		return evnt, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s", entry.EventType)
	}
}

func toCreate(entry storage.EventEntry) (event.Event, error) {
	e := accountCreated{}
	if err := json.Unmarshal(entry.EventData, &e); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data as AccountCreated: %w", err)
	}
	return e, nil
}

func toPayment(entry storage.EventEntry) (event.Event, error) {
	e := incomingPayment{}
	if err := json.Unmarshal(entry.EventData, &e); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event data as IncomingPayment: %w", err)
	}
	return e, nil
}

func toConsume(entry storage.EventEntry) (event.Event, error) {
	e := coffyConsumed{}
	if err := json.Unmarshal(entry.EventData, &e); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event data as CoffyConsumed: %w", err)
	}
	return e, nil
}

func NewAccounting(store *storage.EventRepository) *Accounting {
	service := &Accounting{}
	service.repo = *store
	return service
}

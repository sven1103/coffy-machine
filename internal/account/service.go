package account

import (
	"coffy/internal/event"
	"coffy/internal/storage"
	"encoding/json"
	"errors"
	"fmt"
)

var ErrorNotFound = errors.New("account not found")

type Accounting struct {
	repo storage.EventRepository
}

func (a *Accounting) Create(owner string) (*Account, error) {
	account, err := NewAccount(owner)
	if err != nil {
		return nil, fmt.Errorf("error creating account: %w", err)
	}
	entries, err := a.convertAll(account.events)
	if err != nil {
		return nil, fmt.Errorf("error converting events: %w", err)
	}
	if err := a.repo.SaveAll(entries); err != nil {
		return nil, fmt.Errorf("error saving events: %w", err)
	}
	return account, nil
}

func (a *Accounting) convertAll(events []event.Event) ([]storage.EventEntry, error) {
	entries := make([]storage.EventEntry, 0)
	for _, e := range events {
		entry, err := toEventEntry(e)
		if err != nil {
			return entries, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (a *Accounting) Find(accountID string) (*Account, error) {
	query, err := a.repo.LoadAll(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to load account: %w", err)
	}
	if query == nil || len(query) == 0 {
		return nil, ErrorNotFound
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

func (a *Accounting) ListAll() ([]string, error) {
	query, err := a.repo.FetchByEventType("AccountCreated")
	if err != nil {
		return nil, fmt.Errorf("failed to load accounts: %w", err)
	}
	ids := make([]string, 0)
	for _, entry := range query {
		ids = append(ids, entry.AggregateID)
	}
	return ids, nil
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

func toEventEntry(event event.Event) (storage.EventEntry, error) {
	switch t := event.(type) {
	case AccountCreated:
		data, err := json.Marshal(t)
		if err != nil {
			return storage.EventEntry{}, err
		}
		return storage.EventEntry{AggregateID: t.AccountID, EventType: t.EventType, EventData: data}, nil
	case IncomingPayment:
		data, err := json.Marshal(t)
		if err != nil {
			return storage.EventEntry{}, err
		}
		return storage.EventEntry{AggregateID: t.AccountID, EventType: t.EventType, EventData: data}, nil
	case CoffyConsumed:
		data, err := json.Marshal(t)
		if err != nil {
			return storage.EventEntry{}, err
		}
		return storage.EventEntry{AggregateID: t.AccountID, EventType: t.EventType, EventData: data}, nil
	default:
		return storage.EventEntry{}, errors.New("unknown event type")
	}
}

func toCreate(entry storage.EventEntry) (event.Event, error) {
	e := AccountCreated{}
	if err := json.Unmarshal(entry.EventData, &e); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data as AccountCreated: %w", err)
	}
	return e, nil
}

func toPayment(entry storage.EventEntry) (event.Event, error) {
	e := IncomingPayment{}
	if err := json.Unmarshal(entry.EventData, &e); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event data as IncomingPayment: %w", err)
	}
	return e, nil
}

func toConsume(entry storage.EventEntry) (event.Event, error) {
	e := CoffyConsumed{}
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

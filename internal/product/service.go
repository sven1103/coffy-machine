package product

import (
	"coffy/internal/event"
	"coffy/internal/storage"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type Service struct {
	repo storage.EventRepository
}

func NewService(repo *storage.EventRepository) *Service {
	return &Service{*repo}
}

func (s *Service) ListAll() ([]string, error) {
	query, err := s.repo.FetchByEventType("BeverageCreated")
	if err != nil {
		return nil, fmt.Errorf("failed to load beverages: %w", err)
	}
	ids := make([]string, 0)
	for _, entry := range query {
		ids = append(ids, entry.AggregateID)
	}
	return ids, nil
}

func (s *Service) Find(beverageID string) (*Beverage, error) {
	entries, err := s.repo.LoadAll(beverageID)
	const errorMsg = "failed to load beverage '%s'"
	if err != nil {
		return nil, fmt.Errorf(errorMsg, beverageID)
	}

	// No entries mean the beverage was not found, thus we can return an error here already
	if len(entries) == 0 {
		return nil, fmt.Errorf(errorMsg, beverageID)
	}

	events := make([]event.Event, len(entries))
	for _, entry := range entries {
		e, err := convert(entry)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf(errorMsg, beverageID)
		}
		events = append(events, e)
	}

	b := &Beverage{}
	for _, e := range events {
		if err := b.Apply(e); err != nil {
			log.Println(err)
			return nil, fmt.Errorf(errorMsg, beverageID)
		}
	}
	return b, nil
}

func convert(entry storage.EventEntry) (event.Event, error) {
	switch entry.EventType {
	case "BeverageCreated":
		evnt, err := toBeverageCreated(entry)
		if err != nil {
			return nil, err
		}
		return evnt, nil
	case "PriceUpdated":
		evnt, err := toPriceUpdated(entry)
		if err != nil {
			return nil, err
		}
		return evnt, nil
	default:
		return nil, fmt.Errorf("unknown event type '%s'", entry.EventType)
	}
}

func (s *Service) Create(name string, price float64) (*Beverage, error) {
	b, err := NewBeverage(name, price)
	if err != nil {
		return nil, err
	}

	entries := make([]storage.EventEntry, 0)
	for _, e := range b.Events() {
		entry, err := toEventEntry(e)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	if err := s.repo.SaveAll(entries); err != nil {
		return b, err
	}
	return b, nil
}

func toEventEntry(event event.Event) (storage.EventEntry, error) {
	switch t := event.(type) {
	case BeverageCreated:
		data, err := json.Marshal(t)
		if err != nil {
			return storage.EventEntry{}, err
		}
		return storage.EventEntry{AggregateID: t.AggregateID(), EventType: "BeverageCreated", EventData: data}, nil
	case PriceUpdated:
		data, err := json.Marshal(t)
		if err != nil {
			return storage.EventEntry{}, err
		}
		return storage.EventEntry{AggregateID: t.AggregateID(), EventType: "PriceUpdated", EventData: data}, nil
	default:
		return storage.EventEntry{}, errors.New("unknown event")
	}
}

func toBeverageCreated(entry storage.EventEntry) (event.Event, error) {
	e := BeverageCreated{}
	if err := json.Unmarshal(entry.EventData, &e); err != nil {
		return nil, fmt.Errorf("could not unmarshal event data as BeverageCreated: %w", err)
	}
	return e, nil
}

func toPriceUpdated(entry storage.EventEntry) (event.Event, error) {
	e := PriceUpdated{}
	if err := json.Unmarshal(entry.EventData, &e); err != nil {
		return nil, fmt.Errorf("could not unmarshal event data as PriceUpdated: %w", err)
	}
	return e, nil
}

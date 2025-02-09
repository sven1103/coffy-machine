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

func (s *Service) ListAll() ([]Coffee, error) {
	query, err := s.repo.FetchByEventType("CoffeeCreated")
	if err != nil {
		return nil, fmt.Errorf("failed to load beverages: %w", err)
	}
	bev := make([]Coffee, 0)
	for _, entry := range query {
		r, err := s.Find(entry.AggregateID)
		if err != nil {
			return nil, errors.Join(errors.New("failed to load beverage"), err)
		}
		bev = append(bev, *r)
	}

	return bev, nil
}

func (s *Service) Find(coffeeId string) (*Coffee, error) {
	entries, err := s.repo.LoadAll(coffeeId)
	const errorMsg = "failed to load coffee '%s'"
	if err != nil {
		return nil, fmt.Errorf(errorMsg, coffeeId)
	}

	// No entries mean the beverage was not found, thus we can return an error here already
	if len(entries) == 0 {
		return nil, fmt.Errorf(errorMsg, coffeeId)
	}

	events := make([]event.Event, 0)
	for _, entry := range entries {
		e, err := convert(entry)
		if err != nil {
			return nil, fmt.Errorf(errorMsg, coffeeId)
		}
		events = append(events, e)
	}

	b := &Coffee{}
	for _, e := range events {
		if err := b.apply(e); err != nil {
			log.Println(err)
			return nil, fmt.Errorf(errorMsg, coffeeId)
		}
	}
	return b, nil
}

func convert(entry storage.EventEntry) (event.Event, error) {
	switch entry.EventType {
	case "CoffeeCreated":
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
	case "CvaProvided":
		evnt, err := toCvaProvided(entry)
		if err != nil {
			return nil, err
		}
		return evnt, nil
	case "DetailsUpdated":
		evnt, err := toDetailsUpdated(entry)
		if err != nil {
			return nil, err
		}
		return evnt, nil
	default:
		return nil, fmt.Errorf("unknown event type '%s'", entry.EventType)
	}
}

func toDetailsUpdated(entry storage.EventEntry) (DetailsUpdated, error) {
	e := DetailsUpdated{}
	err := json.Unmarshal(entry.EventData, &e)
	if err != nil {
		return e, fmt.Errorf("failed to unmarshal expected DetailsUpdated event data: %w", err)
	}
	return e, nil
}

var InvalidPropertyError = errors.New("invalid property")

func (s *Service) Create(name string, price float64, score *int, details *CoffeeDetails) (*Coffee, error) {
	b, err := NewCoffee(name, price)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", InvalidPropertyError, err.Error())
	}
	if score != nil {
		err = b.SetCuppingScore(*score)
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %s", InvalidPropertyError, err.Error())
	}

	// Details are optional
	if details != nil {
		d := Details{
			Origin:      details.Origin,
			Description: details.Description,
			RoastHouse:  details.RoastHouse,
			Misc:        details.Misc}
		if err = b.UpdateDetails(d); err != nil {
			return nil, fmt.Errorf("%w: %s", InvalidPropertyError, err.Error())
		}
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
	case CoffeeCreated:
		data, err := json.Marshal(t)
		if err != nil {
			return storage.EventEntry{}, err
		}
		return storage.EventEntry{AggregateID: t.AggregateID(), Date: t.OccurredOn, EventType: "CoffeeCreated", EventData: data}, nil
	case PriceUpdated:
		data, err := json.Marshal(t)
		if err != nil {
			return storage.EventEntry{}, err
		}
		return storage.EventEntry{AggregateID: t.AggregateID(), Date: t.OccurredOn, EventType: "PriceUpdated", EventData: data}, nil
	case CvaProvided:
		data, err := json.Marshal(t)
		if err != nil {
			return storage.EventEntry{}, err
		}
		return storage.EventEntry{AggregateID: t.AggregateID(), Date: t.OccurredOn, EventType: "CvaProvided", EventData: data}, nil
	case DetailsUpdated:
		data, err := json.Marshal(t)
		if err != nil {
			return storage.EventEntry{}, err
		}
		return storage.EventEntry{AggregateID: t.AggregateID(), Date: t.OccurredOn, EventType: "DetailsUpdated", EventData: data}, nil
	default:
		return storage.EventEntry{}, fmt.Errorf("unknown event type: %T", t)
	}
}

func toBeverageCreated(entry storage.EventEntry) (event.Event, error) {
	e := CoffeeCreated{}
	if err := json.Unmarshal(entry.EventData, &e); err != nil {
		return nil, fmt.Errorf("could not unmarshal event data as CoffeeCreated: %w", err)
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

func toCvaProvided(entry storage.EventEntry) (event.Event, error) {
	e := CvaProvided{}
	if err := json.Unmarshal(entry.EventData, &e); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
	}
	return e, nil
}

type CoffeeDetails struct {
	Origin      string            `json:"origin"`
	Description string            `json:"description"`
	RoastHouse  string            `json:"roast_house"`
	Misc        map[string]string `json:"misc"`
}

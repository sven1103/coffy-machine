package equipment

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

func NewService(repo *storage.EventRepository) *Service { return &Service{*repo} }

func (s *Service) ListAll() ([]Machine, error) {
	query, err := s.repo.FetchByEventType("MachineCreated")
	if err != nil {
		return nil, fmt.Errorf("failed to load machines: %w", err)
	}
	m := make([]Machine, 0)
	for _, entry := range query {
		r, err := s.FindById(entry.AggregateID)
		if err != nil {
			return nil, errors.Join(errors.New("failed to load machines"), err)
		}
		m = append(m, *r)
	}
	return m, nil
}

func (s *Service) FindById(machineId string) (*Machine, error) {
	entries, err := s.repo.LoadAll(machineId)
	const errorMsg = "failed to load machine '%s'"
	if err != nil {
		return nil, fmt.Errorf(errorMsg, machineId)
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf(errorMsg, machineId)
	}

	events := make([]event.Event, 0)
	for _, entry := range entries {
		e, err := convert(entry)
		if err != nil {
			return nil, fmt.Errorf(errorMsg, machineId)
		}
		events = append(events, e)
	}

	m := &Machine{}
	for _, e := range events {
		if err := m.apply(e); err != nil {
			log.Println(err)
			return nil, fmt.Errorf(errorMsg, machineId)
		}
	}
	m.Clear()
	return m, nil
}

func (s *Service) Create(brand string, model string) (*Machine, error) {
	m, err := NewMachine(brand, model)
	if err != nil {
		return &Machine{}, err
	}
	entries := make([]storage.EventEntry, 0)
	for _, e := range m.Events() {
		entry, err := toEventEntry(e)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if err = s.repo.SaveAll(entries); err != nil {
		return &Machine{}, err
	}
	m.Clear()
	return m, nil

}

func (s *Service) LoadCoffee(machineId string, coffeeId string) (*Machine, error) {
	m, err := s.FindById(machineId)
	if err != nil {
		return nil, fmt.Errorf("failed to load machine '%s': %w", coffeeId, err)
	}
	err = m.Load(coffeeId)
	if err != nil {
		return nil, fmt.Errorf("failed to load coffee '%s': %w", coffeeId, err)
	}
	events := m.Events()
	entries := make([]storage.EventEntry, 0)
	for _, e := range events {
		entry, err := toEventEntry(e)
		if err != nil {
			return nil, errors.Join(errors.New("failed to load coffee"), err)
		}
		entries = append(entries, entry)
	}
	err = s.repo.SaveAll(entries)
	if err != nil {
		return nil, fmt.Errorf("failed to save update'%s': %w", coffeeId, err)
	}
	m.Clear()
	return m, nil
}

func convert(entry storage.EventEntry) (event.Event, error) {
	switch entry.EventType {
	case "MachineCreated":
		evnt, err := toMachineCreated(entry)
		if err != nil {
			return nil, err
		}
		return evnt, nil
	case "CoffeeLoaded":
		evnt, err := toCoffeeLoaded(entry)
		if err != nil {
			return nil, err
		}
		return evnt, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s", entry.EventType)

	}
}

func toCoffeeLoaded(entry storage.EventEntry) (CoffeeLoaded, error) {
	evnt := CoffeeLoaded{}
	err := json.Unmarshal(entry.EventData, &evnt)
	if err != nil {
		return CoffeeLoaded{}, fmt.Errorf("failed to unmarshal event data: %w", err)
	}
	return evnt, nil
}

func toEventEntry(e event.Event) (storage.EventEntry, error) {
	switch t := e.(type) {
	case MachineCreated:
		data, err := json.Marshal(t)
		if err != nil {
			return storage.EventEntry{}, err
		}
		return storage.EventEntry{AggregateID: e.AggregateID(), EventType: "MachineCreated", Date: e.Occurred(), EventData: data}, nil
	case CoffeeLoaded:
		data, err := json.Marshal(t)
		if err != nil {
			return storage.EventEntry{}, err
		}
		return storage.EventEntry{AggregateID: e.AggregateID(), EventType: "CoffeeLoaded", Date: e.Occurred(), EventData: data}, nil
	default:
		return storage.EventEntry{}, fmt.Errorf("failed to convert event to entry: unkown event type '%T'", t)
	}
}

func toMachineCreated(entry storage.EventEntry) (MachineCreated, error) {
	evnt := MachineCreated{}
	err := json.Unmarshal(entry.EventData, &evnt)
	if err != nil {
		return MachineCreated{}, fmt.Errorf("failed to unmarshal event data: %w", err)
	}
	return evnt, nil
}

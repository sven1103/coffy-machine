package equipment

import (
	"coffy/internal/event"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Machine struct {
	AggregateID string        // machines unique ID in Coffy
	Brand       string        // brand of the machine, e.g. Phillips
	Model       string        // model of the machine, e.g. EP2334/10
	coffee      string        // current loaded product.Coffee, its ID
	events      []event.Event // object's event cache with uncommited event.Event
}

func NewMachine(brand string, model string) (*Machine, error) {
	created := MachineCreated{ID: uuid.New().String(), Brand: brand, Model: model, OccurredOn: time.Now()}
	m := Machine{}
	if err := m.apply(created); err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *Machine) Events() []event.Event {
	return m.events
}

func (m *Machine) Clear() {
	m.events = []event.Event{}
}

func (m *Machine) Load(coffeeID string) error {
	e := CoffeeLoaded{CoffeeID: coffeeID, MachineID: m.AggregateID, OccurredOn: time.Now()}
	if err := m.apply(e); err != nil {
		return err
	}
	return nil
}

func (m *Machine) Coffee() (string, error) {
	if len(m.coffee) == 0 {
		return "", fmt.Errorf("no coffee loaded")
	}
	return m.coffee, nil
}

func (m *Machine) apply(e event.Event) error {
	switch theEvent := e.(type) {
	case MachineCreated:
		return m.applyMachineCreated(theEvent)
	case CoffeeLoaded:
		return m.applyCoffeeLoaded(theEvent)
	default:
		return fmt.Errorf("unknown event type '%T'", theEvent)
	}
}

func (m *Machine) applyMachineCreated(e MachineCreated) error {
	m.AggregateID = e.ID
	m.Brand = e.Brand
	m.Model = e.Model
	m.events = append(m.events, e)
	return nil
}

func (m *Machine) applyCoffeeLoaded(e CoffeeLoaded) error {
	if e.MachineID != m.AggregateID {
		return fmt.Errorf("event does not belong to this aggregate")
	}
	m.coffee = e.CoffeeID
	m.events = append(m.events, e)
	return nil
}

type MachineCreated struct {
	ID         string    `json:"id"`
	Brand      string    `json:"brand"`
	Model      string    `json:"model"`
	OccurredOn time.Time `json:"occurredOn"`
}

func (e MachineCreated) AggregateID() string {
	return e.ID
}

func (e MachineCreated) Occurred() time.Time {
	return e.OccurredOn
}

func (e MachineCreated) Type() string {
	return "MachineCreated"
}

type CoffeeLoaded struct {
	MachineID  string    `json:"id"`         // the affected machine
	CoffeeID   string    `json:"coffee_id"`  // the loaded coffee (ID)
	OccurredOn time.Time `json:"occurredOn"` // the event time point
}

func (e CoffeeLoaded) AggregateID() string {
	return e.MachineID
}

func (e CoffeeLoaded) Type() string {
	return "CoffeeLoaded"
}

func (e CoffeeLoaded) Occurred() time.Time {
	return e.OccurredOn
}

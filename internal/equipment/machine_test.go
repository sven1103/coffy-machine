package equipment

import (
	"github.com/google/uuid"
	"testing"
)

func TestMachineCreated(t *testing.T) {
	m, err := NewMachine("Philips", "EP2334")
	if err != nil {
		t.Errorf("Could not create machine: %s", err)
	}
	e := m.Events()[0]
	switch exp := e.(type) {
	case MachineCreated:
		break
	default:
		t.Errorf("Wrong event type, expected MachineCreated got: %T", exp)
		return
	}

	if m.Model != "EP2334" {
		t.Errorf("Expected model '%s' but got '%s'", "EP2334", m.Model)
	}

	if m.Brand != "Philips" {
		t.Errorf("Expected brand '%s' but got '%s'", "Philips", m.Brand)
	}
}

func TestCoffeeLoaded(t *testing.T) {
	m, err := NewMachine("Philips", "EP2334")
	if err != nil {
		t.Errorf("Could not create machine: %s", err)
	}
	m.Clear()
	coffeeRef := uuid.New().String()
	if err := m.Load(coffeeRef); err != nil {
		t.Errorf("failed to load coffee: %s", err)
		return
	}

	e := m.Events()[0]
	switch exp := e.(type) {
	case CoffeeLoaded:
		break
	default:
		t.Errorf("wrong event type, expected CoffeeLoaded  got: %T", exp)
		return
	}

	c, err := m.Coffee()
	if err != nil {
		t.Errorf("expected machine to be loaded")
		return
	}
	if c != coffeeRef {
		t.Errorf("Mismatching coffee ref. Expected '%s' but was '%s'", coffeeRef, c)
	}
}

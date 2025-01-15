package event

import (
	"testing"
	"time"
)

type exampleEvent struct {
	ID         string
	occurredOn time.Time
	eventType  string
}

func (e *exampleEvent) AggregateID() string {
	return e.ID
}

func (e *exampleEvent) Occurred() time.Time {
	return e.occurredOn
}

func (e *exampleEvent) Type() string {
	return e.eventType
}

func TestEvent(t *testing.T) {
	currentTime := time.Now()
	exampleEvent := exampleEvent{ID: "123", occurredOn: currentTime}

	returnedID := exampleEvent.AggregateID()
	returnedTime := exampleEvent.Occurred()

	if returnedID != "123" {
		t.Errorf("Event returned ID should be '123', got '%s'", returnedID)
	}
	if returnedTime != currentTime {
		t.Errorf("Event returned Time should be '%s', got '%s'", currentTime, returnedTime)
	}
}

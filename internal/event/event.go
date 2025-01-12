package event

import (
	"time"
)

// An Event provides a common interface for coffy domain events. An Event always contains the aggregate ID of the aggregate
// that emitted the Event.
//
// The second property is the recorded time when the Event has happened.
type Event interface {
	AggregateID() string
	OccurredOn() time.Time
	Type() string
}

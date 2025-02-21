package bus

import "github.com/oklog/ulid/v2"

type EventData interface {
	EventType() string
}

type Event struct {
	EventData
	ID ulid.ULID
}

type StateChangeEvent struct {
	Key      uint32
	NewState any
}

// EventType implements EventData.
func (s *StateChangeEvent) EventType() string {
	return "state_change"
}

var _ EventData = (*StateChangeEvent)(nil)

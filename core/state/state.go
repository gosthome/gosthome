package state

import (
	"context"
	"log/slog"

	"github.com/gosthome/gosthome/core/bus"
	"github.com/gosthome/gosthome/core/entity"
)

type State_[T any] struct {
	emitter bus.Emitter[bus.StateChangeEvent, *bus.StateChangeEvent]

	e     entity.Entity
	value T
}

func NewState[T any](ctx context.Context, e entity.Entity, initial T) (*State_[T], error) {
	b := bus.Get(ctx)

	return &State_[T]{
		emitter: bus.MakeEventEmitter[bus.StateChangeEvent](b),
		e:       e,
		value:   initial,
	}, nil
}

// State implements entity.WithState.
func (s *State_[T]) State() T {
	return s.value
}

func (s *State_[T]) SetState(nv T) {
	s.value = nv
	ns := s.value
	slog.Info("Sending state", "id", s.e.ID(), "state", ns)
	s.emitter.Emit(&bus.StateChangeEvent{
		Key:      s.e.HashID(),
		NewState: &ns,
	})
}

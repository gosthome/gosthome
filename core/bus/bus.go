package bus

import (
	"context"
	"log/slog"
	"sync"
	"weak"

	"github.com/majfault/signal"
	"github.com/majfault/signal/dispatcher"
	"github.com/oklog/ulid/v2"
)

type eventSignal = signal.Signal1[*Event]
type EventSubsciption struct {
	sig  weak.Pointer[eventSignal]
	slot weak.Pointer[signal.Slot1[*Event]]
}

func (b *EventSubsciption) Close() {
	sig := b.sig.Value()
	if sig == nil {
		return
	}
	slot := b.slot.Value()
	if slot == nil {
		return
	}
	sig.Disconnect(slot)
}

type serviceCallSignal = signal.Signal1[*serviceRequest]

type Bus struct {
	mux      sync.RWMutex
	events   map[string]*eventSignal
	services map[string]*serviceCallSignal
}

func New() *Bus {
	return &Bus{
		events:   make(map[string]*eventSignal),
		services: make(map[string]*serviceCallSignal),
	}
}

type busCtxKey struct{}

func Context(ctx context.Context, b *Bus) context.Context {
	return context.WithValue(ctx, busCtxKey{}, b)
}

func Get(ctx context.Context) *Bus {
	v := ctx.Value(busCtxKey{})
	if v == nil {
		return nil
	}
	b, ok := v.(*Bus)
	if !ok {
		return nil
	}
	return b
}

func rlocked[R any](b *Bus, f func() R) R {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return f()
}

func locked[R any](b *Bus, f func() R) R {
	b.mux.Lock()
	defer b.mux.Unlock()
	return f()
}

func eventKey(e EventData) string {
	return "event/" + e.EventType()
}

func serviceKey(e ServiceRequestData) string {
	return "event/" + e.ServiceType()
}

type Emitter[T any, PT interface {
	*T
	EventData
}] interface {
	Emit(PT)
}

type emitter[T any, PT interface {
	*T
	EventData
}] struct {
	sig weak.Pointer[eventSignal]
}

func (e *emitter[T, PT]) Emit(data PT) {
	es := e.sig.Value()
	if es == nil {
		return
	}
	event := &Event{
		EventData: data,
		ID:        ulid.Make(),
	}
	slog.Debug("Bus emitting", "event", event)
	es.Emit(event)
	slog.Debug("Bus emitted", "event", event)
}

func MakeEventEmitter[T any, PT interface {
	*T
	EventData
}](b *Bus) Emitter[T, PT] {
	ek := eventKey(PT(nil))
	es := locked(b, func() *eventSignal {
		es, ok := b.events[ek]
		if !ok {
			es = &eventSignal{}
			b.events[ek] = es
		}
		return es
	})
	return &emitter[T, PT]{
		sig: weak.Make(es),
	}
}

func (b *Bus) emitEvent(data EventData, source string) {

}

type eventHandler struct {
	etype  string
	handle func(*Event)
}

func EventHandler[T any, PT interface {
	EventData
	*T
}](f func(t PT)) eventHandler {
	return eventHandler{
		etype: eventKey(PT(nil)),
		handle: func(e *Event) {
			slog.Debug("Bus handling", "event", e)
			f(e.EventData.(PT))
		},
	}
}

func (b *Bus) HandleEvents(h eventHandler) EventSubsciption {
	es := locked(b, func() *eventSignal {
		es, ok := b.events[h.etype]
		if !ok {
			es = &eventSignal{}
			b.events[h.etype] = es
		}
		return es
	})
	slot := es.Connect(&dispatcher.Queued{
		BlockOnFull: true,
	}, h.handle)
	return EventSubsciption{
		sig:  weak.Make(es),
		slot: weak.Make(slot),
	}
}

func (b *Bus) CallService(data ServiceRequestData) *ulid.ULID {
	sk := serviceKey(data)
	scID := ulid.Make()
	service := &serviceRequest{
		ID:   scID,
		Data: data,
	}
	var es *serviceCallSignal
	if !locked(b, func() (ok bool) {
		es, ok = b.services[sk]
		return
	}) {
		slog.Error("calling unknown service", "key", sk, "event", service)
		return nil
	}
	es.Emit(service)
	return &scID
}

type serviceHandler struct {
	stype  string
	handle func(*serviceRequest)
}

func ServiceHandler[T any, PT interface {
	ServiceRequestData
	*T
}](f func(t PT)) serviceHandler {
	return serviceHandler{
		stype: serviceKey(PT(nil)),
		handle: func(e *serviceRequest) {
			f(e.Data.(PT))
		},
	}
}

func ServiceHandlerWithRespose[R any, T any, PT interface {
	ServiceRequestData
	*T
}](b *Bus, f func(t PT) R) serviceHandler {
	em := MakeEventEmitter[ServiceResponseEvent](b)
	return serviceHandler{
		stype: serviceKey(PT(nil)),
		handle: func(e *serviceRequest) {
			r := f(e.Data.(PT))
			em.Emit(&ServiceResponseEvent{
				RequestID: e.ID,
				Response:  r,
			})
		},
	}
}

func (b *Bus) HandleServiceCalls(h serviceHandler) {
	es := locked(b, func() *serviceCallSignal {
		es, ok := b.services[h.stype]
		if !ok {
			es = &serviceCallSignal{}
			b.services[h.stype] = es
		}
		return es
	})
	es.Connect(&dispatcher.Queued{
		BlockOnFull: true,
	}, h.handle)
}

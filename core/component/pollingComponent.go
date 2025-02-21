package component

import (
	"context"
	"log/slog"
	"reflect"
	"sync"
	"time"
	"weak"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type PollingComponentConfig[T any, PT interface {
	*T
	Poller
}] struct {
	UpdateInterval time.Duration `yaml:"update_interval"`
}

// Validate implements Config.
func (p *PollingComponentConfig[T, PT]) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(
			&p.UpdateInterval,
			validation.Required,
			validation.Min(0*time.Second).Exclusive()),
	)
}

type Poller interface {
	Component
	Poll()
}

type PollingComponent[T any, PT interface {
	*T
	Poller
}] struct {
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	updateInterval time.Duration
	poll           weak.Pointer[T]
}

func NewPollingComponent[T any, PT interface {
	*T
	Poller
}](ctx context.Context, poll PT, cfg *PollingComponentConfig[T, PT]) (*PollingComponent[T, PT], error) {
	ctx, cancel := context.WithCancel(ctx)
	return &PollingComponent[T, PT]{
		ctx:            ctx,
		cancel:         cancel,
		wg:             sync.WaitGroup{},
		poll:           weak.Make(poll),
		updateInterval: cfg.UpdateInterval,
	}, nil
}

// Setup implements Component.
func (p *PollingComponent[T, PT]) Setup() {
	p.wg.Add(1)
	done := p.ctx.Done()
	go func() {
		defer p.wg.Done()
		type_ := reflect.TypeOf(PT(nil)).String()
		timer := time.NewTicker(p.updateInterval)
		slog.Debug("Starting polling for", "type", type_)
		for {
			select {
			case <-timer.C:
				t := p.poll.Value()
				if t != nil {
					slog.Debug("Poll", "type", type_)
					PT(t).Poll()
				} else {
					slog.Warn("Unable to poll for expired", "type", type_)
				}
			case <-done:
				timer.Stop()
				slog.Debug("Stopped polling for", "type", type_)
				return
			}
		}
	}()
}

// Close implements Component.
func (p *PollingComponent[T, PT]) Close() error {
	p.cancel()
	p.wg.Wait()
	return nil
}

// InitializationPriority implements Component.
func (p *PollingComponent[T, PT]) InitializationPriority() InitializationPriority {
	panic("abstract method")
}

package demo

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/gosthome/gosthome/components/binarysensor"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/entity"
)

type DemoBinarySensorConfig struct {
	binarysensor.BaseBinarySensorConfig[DemoBinarySensor, *DemoBinarySensor] `yaml:",inline"`

	poll component.PollingComponentConfig[DemoBinarySensor, *DemoBinarySensor] `yaml:",inline"`
}

func NewDemoBinarySensorConfig() DemoBinarySensorConfig {
	return DemoBinarySensorConfig{
		poll: component.PollingComponentConfig[DemoBinarySensor, *DemoBinarySensor]{
			UpdateInterval: 10 * time.Second,
		},
	}
}

type DemoBinarySensor struct {
	binarysensor.BaseBinarySensor[DemoBinarySensor, *DemoBinarySensor]
	poll *component.PollingComponent[DemoBinarySensor, *DemoBinarySensor]
	rand *rand.Rand
}

func NewDemoBinarySensor(ctx context.Context, src rand.Source, cfg *DemoBinarySensorConfig) (*DemoBinarySensor, error) {
	d := &DemoBinarySensor{
		rand: rand.New(src),
	}
	bs, err := binarysensor.NewBaseBinarySensor(ctx, d, &cfg.BaseBinarySensorConfig)
	if err != nil {
		return nil, err
	}
	d.BaseBinarySensor = *bs
	d.poll, err = component.NewPollingComponent(ctx, d, &cfg.poll)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// InitializationPriority implements component.Poller.
func (d *DemoBinarySensor) InitializationPriority() component.InitializationPriority {
	return component.InitializationPriorityData
}

// Setup implements component.Poller.
func (d *DemoBinarySensor) Setup() {
	d.Poll()
	d.poll.Setup()
}

// Poll implements component.Poller.
func (d *DemoBinarySensor) Poll() {
	d.SetState(entity.BinarySensorState{
		State:   d.rand.Float32() < 0.5,
		Missing: false,
	})
}

// Close implements component.Poller.
func (d *DemoBinarySensor) Close() error {
	return d.poll.Close()
}

var _ (entity.BinarySensor) = (*DemoBinarySensor)(nil)
var _ (component.Poller) = (*DemoBinarySensor)(nil)

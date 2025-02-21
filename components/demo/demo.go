package demo

import (
	"context"
	"math/rand/v2"
	"sync"

	"github.com/gosthome/gosthome/core"
	"github.com/gosthome/gosthome/core/component"
)

type safePCG struct {
	mx  sync.Mutex
	src rand.PCG
}

func (s *safePCG) Seed(seed1, seed2 uint64) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.src.Seed(seed1, seed2)
}

// Uint64 implements rand.Source.
func (s *safePCG) Uint64() uint64 {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.src.Uint64()
}

func newSafePCG(s1, s2 uint64) *safePCG {
	return &safePCG{
		mx:  sync.Mutex{},
		src: *rand.NewPCG(s1, s2),
	}
}

var _ rand.Source = (*safePCG)(nil)

type Demo struct {
	r *safePCG
}

func New(ctx context.Context, cfg *Config) ([]component.Component, error) {
	node := core.GetNode(ctx)
	if node == nil {
		panic("No node in config during binary_sensors initialization")
	}
	d := &Demo{
		r: newSafePCG(cfg.Seeds[0], cfg.Seeds[1]),
	}
	ret := []component.Component{}
	for _, bsc := range cfg.BinarySensors {
		bs, err := NewDemoBinarySensor(ctx, d.r, &bsc)
		if err != nil {
			return nil, err
		}
		ret = append(ret, bs)
		err = node.RegisterBinarySensor(bs)
		if err != nil {
			return nil, err
		}
	}
	for _, bc := range cfg.Buttons {
		b, err := NewDemoButton(ctx, d, &bc)
		if err != nil {
			return nil, err
		}
		ret = append(ret, b)
		err = node.RegisterButton(b)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

// Setup implements component.Component.
func (c *Demo) Setup() {

}

// Close implements component.Component.
func (c *Demo) Close() error {
	return nil
}

// InitializationPriority implements component.Component.
func (c *Demo) InitializationPriority() component.InitializationPriority {
	return component.InitializationPriorityHardware
}

var _ component.Component = (*Demo)(nil)

package demo

import (
	"context"
	"log/slog"
	"math/rand/v2"

	"github.com/gosthome/gosthome/components/button"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/entity"
)

type DemoButtonConfig struct {
	button.BaseButtonConfig[DemoButton, *DemoButton] `yaml:",inline"`
}

func NewDemoButtonConfig() DemoButtonConfig {
	return DemoButtonConfig{}
}

func (t *DemoButtonConfig) ValidateWithContext(ctx context.Context) error {
	return t.BaseButtonConfig.ValidateWithContext(ctx)
}

type DemoButton struct {
	button.BaseButton[DemoButton, *DemoButton]
	demo *Demo
}

func (t *DemoButton) Press(ctx context.Context) error {
	slog.Info("Demo button pressed, reseeding hash")
	r := rand.New(t.demo.r)
	t.demo.r.Seed(r.Uint64(), r.Uint64())
	return nil
}

// Close implements component.Component.
func (t *DemoButton) Close() error {
	return nil
}

// InitializationPriority implements component.Component.
func (t *DemoButton) InitializationPriority() component.InitializationPriority {
	return component.InitializationPriorityProcessor
}

// Setup implements component.Component.
func (t *DemoButton) Setup() {
}

func NewDemoButton(ctx context.Context, d *Demo, cfg *DemoButtonConfig) (*DemoButton, error) {
	ret := &DemoButton{
		demo: d,
	}
	bs, err := button.NewBaseButton(ctx, ret, &cfg.BaseButtonConfig)
	if err != nil {
		return nil, err
	}
	ret.BaseButton = *bs
	return ret, nil
}

var _ component.Component = (*DemoButton)(nil)
var _ entity.Button = (*DemoButton)(nil)

package uart

import (
	"context"
	"log/slog"

	"github.com/gosthome/gosthome/components/button"
	"github.com/gosthome/gosthome/core/bus"
	"github.com/gosthome/gosthome/core/component"
	cv "github.com/gosthome/gosthome/core/configvalidation"
	"github.com/gosthome/gosthome/core/entity"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ButtonConfig struct {
	button.BaseButtonConfig[Button, *Button] `yaml:",inline"`

	UARTID string    `yaml:"uart_id"`
	Data   *cv.Bytes `yaml:"data"`
}

func NewButtonConfig() *ButtonConfig {
	return &ButtonConfig{}
}

// Validate implements validation.Validatable.
func (c *ButtonConfig) ValidateWithContext(ctx context.Context) error {
	return cv.ValidateEmbedded(
		c.BaseButtonConfig.ValidateWithContext(ctx),
		validation.ValidateStructWithContext(
			ctx, c, validation.Field(&c.UARTID, cv.String(cv.Optional(cv.Name()))),
		),
	)
}

var _ component.Config = (*ButtonConfig)(nil)

type Button struct {
	button.BaseButton[Button, *Button]
	b      *bus.Bus
	uartID uint32
	data   []byte
}

func NewButton(ctx context.Context, cfg *ButtonConfig) ([]component.Component, error) {
	slog.Info("Init btn", "cfg", cfg)
	ret := &Button{}
	bb, err := button.NewBaseButton(ctx, ret, &cfg.BaseButtonConfig)
	if err != nil {
		return nil, err
	}
	ret.BaseButton = *bb
	ret.b = bus.Get(ctx)
	if cfg.UARTID != "" {
		ret.uartID = entity.HashID(cfg.UARTID)
	}
	ret.data = cfg.Data.Data
	return []component.Component{ret}, nil
}

// Setup implements component.Component.
func (b *Button) Setup() {

}

func (b *Button) Press(ctx context.Context) error {
	b.b.CallService(&UARTWrite{
		Key:  b.uartID,
		Data: b.data,
	})
	return nil
}

// Close implements component.Component.
func (b *Button) Close() error {

	return nil
}

// InitializationPriority implements component.Component.
func (c *Button) InitializationPriority() component.InitializationPriority {
	return component.InitializationPriorityHardware
}

var _ component.Component = (*Button)(nil)

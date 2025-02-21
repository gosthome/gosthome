package button

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gosthome/gosthome/core"
	"github.com/gosthome/gosthome/core/bus"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/config"
	"github.com/gosthome/gosthome/core/entity"
)

type Config struct {
	component.ConfigOf[entity.ButtonDomain, *entity.ButtonDomain]
	config.PlatformConfig
}

// Validate implements component.Config.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return c.PlatformConfig.ValidateWithContext(ctx)
}

func NewConfig() *Config {
	return &Config{
		PlatformConfig: config.PlatformConfig{
			DomainType: entity.DomainTypeButton,
		},
	}
}

type ButtonPress struct {
	Key uint32
}

// ServiceType implements bus.ServiceRequestData.
func (b *ButtonPress) ServiceType() string {
	return "button.press"
}

var _ bus.ServiceRequestData = (*ButtonPress)(nil)

func New(ctx context.Context, c *Config) ([]component.Component, error) {
	node := core.GetNode(ctx)
	if node == nil {
		panic("No node in config during button initialization")
	}
	domain := &entity.ButtonDomain{}
	ret := []component.Component{domain}
	slog.Info("button domain", "c", fmt.Sprintf("%#v", c))
	for _, platformConfig := range c.Configs {
		cd, ok := node.Config.Registry.GetEntityComponent(entity.DomainTypeButton, platformConfig.Platform)
		if !ok {
			panic("unregistered button platform in config " + platformConfig.Platform)
		}
		comp, err := cd.Component(ctx, platformConfig.Config.Config)
		if err != nil {
			return nil, err
		}
		for _, c := range comp {
			domain.Register(c.(entity.Button))
		}
		ret = append(ret, comp...)
	}
	slog.Info("Initialized button domain")
	node.ButtonDomain = domain
	b := bus.Get(ctx)
	if b == nil {
		panic("No bus in config during button initialization")
	}
	b.HandleServiceCalls(bus.ServiceHandlerWithRespose(b, func(t *ButtonPress) error {
		button, ok := domain.FindByKey(t.Key)
		if !ok {
			slog.Error("Tried to press nonexisting button", "key", t.Key)
			return fmt.Errorf("tried to press nonexisting button %d", t.Key)
		}
		return button.Press(ctx)
	}))
	return ret, nil
}

var _ component.Config = (*Config)(nil)

package demo

import (
	"context"

	"github.com/gosthome/gosthome/components/binarysensor"
	"github.com/gosthome/gosthome/components/button"
	"github.com/gosthome/gosthome/core/component"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	component.ConfigOf[Demo, *Demo]

	Seeds [2]uint64

	BinarySensors []DemoBinarySensorConfig `yaml:"binary_sensors"`
	Buttons       []DemoButtonConfig       `yaml:"buttons"`
}

func modify[T any](c T, f func(c *T)) T {
	f(&c)
	return c
}

func NewConfig() *Config {
	return &Config{
		Seeds: [2]uint64{1, 2},
		BinarySensors: []DemoBinarySensorConfig{
			modify(NewDemoBinarySensorConfig(), func(c *DemoBinarySensorConfig) {
				c.Name = "Demo Basement Floor Wet"
				c.DeviceClass = string(binarysensor.DeviceClassMoisture)
			}),
			modify(NewDemoBinarySensorConfig(), func(c *DemoBinarySensorConfig) {
				c.Name = "Demo Movement Backyard"
				c.DeviceClass = string(binarysensor.DeviceClassMotion)
			}),
		},
		Buttons: []DemoButtonConfig{
			modify(NewDemoButtonConfig(), func(c *DemoButtonConfig) {
				c.Name = "Demo Regenerate Seed"
				c.DeviceClass = string(button.DeviceClassRestart)
			}),
		},
	}
}

// Validate implements validation.Validatable.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c, validation.Field(&c.BinarySensors))
}

// AutoLoad implements component.AutoLoader.
func (c *Config) AutoLoad() []string {
	return []string{
		binarysensor.COMPONENT_KEY,
		button.COMPONENT_KEY,
	}
}

var _ component.Config = (*Config)(nil)
var _ component.AutoLoader = (*Config)(nil)

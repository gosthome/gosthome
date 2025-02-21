package file

import (
	"context"

	"github.com/gosthome/gosthome/components/binarysensor"
	"github.com/gosthome/gosthome/core/component"
)

type Config struct {
	component.ConfigOf[File, *File]
}

func NewConfig() *Config {
	return &Config{}
}

// Validate implements validation.Validatable.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	// return validation.ValidateStruct(c, validation.Field(&c.BinarySensors))
	return nil
}

// AutoLoad implements component.AutoLoader.
func (c *Config) AutoLoad() []string {
	return []string{
		binarysensor.COMPONENT_KEY,
	}
}

var _ component.Config = (*Config)(nil)
var _ component.AutoLoader = (*Config)(nil)

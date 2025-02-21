package component

import (
	"context"
	"fmt"

	cv "github.com/gosthome/gosthome/core/configvalidation"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

type ConfigDecoder struct {
	Config    Config
	Marshal   func(context.Context, Config) (interface{}, error)
	Unmarshal func(ctx context.Context, cfg Config, unmarshal func(interface{}) error) error
}

// Validate implements validation.Validatable.
func (c *ConfigDecoder) ValidateWithContext(ctx context.Context) error {
	return c.Config.ValidateWithContext(ctx)
}

// MarshalYAML implements yaml.InterfaceMarshalerContext.
func (c *ConfigDecoder) MarshalYAML(context.Context) (interface{}, error) {
	panic("unimplemented")
}

// UnmarshalYAML implements yaml.InterfaceUnmarshalerContext.
func (c *ConfigDecoder) UnmarshalYAML(ctx context.Context, node ast.Node) error {
	if c.Unmarshal == nil {
		return fmt.Errorf("component decoder wih nil Unmarshal")
	}
	dec := ctx.Value(cv.ConfigYAMLDecoderKey{}).(*yaml.Decoder)
	return c.Unmarshal(ctx, c.Config, func(i interface{}) error {
		return dec.DecodeFromNodeContext(ctx, node, i)
	})
}

var _ yaml.InterfaceMarshalerContext = (*ConfigDecoder)(nil)
var _ yaml.NodeUnmarshalerContext = (*ConfigDecoder)(nil)
var _ cv.Validatable = (*ConfigDecoder)(nil)

func Marshal[T any, PT interface {
	*T
	Config
}](context.Context, Config) (interface{}, error) {
	panic("unimplemented")
}

func Unmarshal[T any, PT interface {
	*T
	Config
}](ctx context.Context, cfg Config, unmarshal func(interface{}) error) error {
	t := cfg.(PT)
	return unmarshal(t)
}

func NewConfigDecoder[T any, PT interface {
	*T
	Config
}](config PT) *ConfigDecoder {
	return &ConfigDecoder{
		Config:    config,
		Marshal:   Marshal[T, PT],
		Unmarshal: Unmarshal[T, PT],
	}
}

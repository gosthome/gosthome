package config

import (
	"context"

	"github.com/gosthome/gosthome/core/component"
	cv "github.com/gosthome/gosthome/core/configvalidation"
	"github.com/gosthome/gosthome/core/entity"
	"github.com/gosthome/gosthome/core/registry"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

type Platform struct {
	Platform string                   `yaml:"platform"`
	Config   *component.ConfigDecoder `yaml:"-"`
}

func (p *Platform) ValidateWithContext(ctx context.Context) error {
	return nil
}

type PlatformConfig struct {
	DomainType entity.DomainType
	Configs    []*Platform
}

// Validate implements component.Config.
func (p *PlatformConfig) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateWithContext(ctx, p.Configs)
}

// MarshalYAML implements yaml.InterfaceMarshalerContext.
func (p *PlatformConfig) MarshalYAML(context.Context) (interface{}, error) {
	panic("unimplemented")
}

// UnmarshalYAML implements yaml.NodeUnmarshaler.
func (p *PlatformConfig) UnmarshalYAML(ctx context.Context, src ast.Node) error {
	list := []ast.Node{}
	dec := ctx.Value(cv.ConfigYAMLDecoderKey{}).(*yaml.Decoder)
	err := dec.DecodeFromNodeContext(ctx, src, &list)
	if err != nil {
		return err
	}
	cr := ctx.Value(cv.ComponentRegistryKey{}).(*registry.Registry)
	data := make([]*Platform, 0, len(list))
	for i, node := range list {
		if node == nil {
			return cv.ErrUnknownField(src, "%s list element %d is empty, should contain at least `platform` key", p.DomainType, i)
		}
		platform := &Platform{}
		err = dec.DecodeFromNodeContext(ctx, node, platform)
		if err != nil {
			return err
		}
		cd, ok := cr.GetEntityComponent(p.DomainType, platform.Platform)
		if !ok {
			return cv.ErrUnknownField(node, "%s cannot parse unknown platform: %s", p.DomainType, platform.Platform)
		}
		platform.Config = cd.Config()
		err = dec.DecodeFromNodeContext(ctx, node, platform.Config)
		if err != nil {
			return err
		}
		data = append(data, platform)
	}
	p.Configs = data
	return nil
}

var _ yaml.InterfaceMarshalerContext = (*PlatformConfig)(nil)
var _ yaml.NodeUnmarshalerContext = (*PlatformConfig)(nil)
var _ cv.Validatable = (*PlatformConfig)(nil)

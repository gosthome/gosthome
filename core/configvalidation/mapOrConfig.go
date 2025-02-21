package cv

import (
	"context"
	"log/slog"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

type MapOrValue[Config any, PConfig interface {
	*Config
}] struct {
	factory func() PConfig
	Configs []PConfig
}

// MarshalYAML implements yaml.InterfaceMarshalerContext.
func (m *MapOrValue[Config, PConfig]) MarshalYAML(context.Context) (interface{}, error) {
	panic("unimplemented")
}

// UnmarshalYAML implements yaml.InterfaceUnmarshalerContext.
func (m *MapOrValue[Config, PConfig]) UnmarshalYAML(ctx context.Context, node ast.Node) error {
	if node == nil {
		m.Configs = []PConfig{m.factory()}
		return nil
	}
	dec := ctx.Value(ConfigYAMLDecoderKey{}).(*yaml.Decoder)
	if node.Type() == ast.MappingType {
		m.Configs = make([]PConfig, 0, 1)
		nc := m.factory()
		err := dec.DecodeFromNodeContext(ctx, node, nc)
		if err != nil {
			slog.Error("Here", "err", err)
			return err
		}
		m.Configs = append(m.Configs, nc)
		return nil
	}
	var sequence []ast.Node
	serr := dec.DecodeFromNodeContext(ctx, node, &sequence)
	if serr != nil {
		return serr
	}
	m.Configs = make([]PConfig, 0, len(sequence))
	for _, node := range sequence {
		nc := m.factory()
		serr = dec.DecodeFromNodeContext(ctx, node, nc)
		if serr != nil {
			return serr
		}
		m.Configs = append(m.Configs, nc)
	}
	return nil
}

var _ yaml.InterfaceMarshalerContext = (*MapOrValue[string, *string])(nil)
var _ yaml.NodeUnmarshalerContext = (*MapOrValue[string, *string])(nil)

func NewMapOrValue[Config any, PConfig interface {
	*Config
}](factory func() PConfig) MapOrValue[Config, PConfig] {
	return MapOrValue[Config, PConfig]{
		factory: factory,
	}
}

package config

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/gosthome/gosthome/core/component"
	cv "github.com/gosthome/gosthome/core/configvalidation"
	"github.com/gosthome/gosthome/core/registry"
)

type Configs map[string]*component.ConfigDecoder

// Validate implements validation.Validatable.
func (c *Configs) ValidateWithContext(ctx context.Context) error {
	errs := make(validation.Errors)
	for k, v := range *c {
		err := v.ValidateWithContext(ctx)
		if err != nil {
			errs[k] = err
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// MarshalYAML implements yaml.InterfaceMarshalerContext.
func (c *Configs) MarshalYAML(context.Context) (interface{}, error) {
	panic("unimplemented")
}

// UnmarshalYAML implements yaml.InterfaceUnmarshaler.
func (c *Configs) UnmarshalYAML(ctx context.Context, src ast.Node) error {
	// var keys map[ast.Node]ast.Node
	dec := ctx.Value(cv.ConfigYAMLDecoderKey{}).(*yaml.Decoder)
	// err := dec.DecodeFromNodeContext(ctx, src, &keys)
	// if err != nil {
	// 	return err
	// }
	cr := ctx.Value(cv.ComponentRegistryKey{}).(*registry.Registry)
	data := make(map[string]*component.ConfigDecoder)
	mn, ok := src.(*ast.MappingNode)
	if !ok {
		return &yaml.UnexpectedNodeTypeError{
			Actual:   src.Type(),
			Expected: ast.MappingType,
			Token:    src.GetToken(),
		}
	}
	for _, mapValue := range mn.Values {
		knode := mapValue.Key
		node := mapValue.Value
		var k string
		err := dec.DecodeFromNodeContext(ctx, knode, &k)
		if err != nil {
			return err
		}
		if k == "gosthome" {
			continue
		}
		enode := node
		if enode == nil {
			enode = knode
		}
		if enode == nil {
			enode = src
		}
		declaration, ok := cr.Get(k)
		if !ok {
			return cv.ErrUnknownField(enode, "cannot parse unknown component: %s", k)
		}
		_, ok = (data)[k]
		if ok {
			// enode := node
			// if enode == nil {
			// 	enode = src
			// }
			return cv.ErrDuplicateKey(node, "component specified twice!: %s", k)
		}
		cfg := declaration.Config()
		data[k] = cfg
		if node != nil && node.Type() != ast.NullType {
			err = dec.DecodeFromNodeContext(ctx, node, cfg.Config)
			if err != nil {
				return err
			}
		}
		err = cfg.ValidateWithContext(ctx)
		if err != nil {
			return &yaml.SyntaxError{
				Token:   enode.GetToken(),
				Message: err.Error(),
			}
		}
	}
	*c = data
	return nil
}

var _ yaml.InterfaceMarshalerContext = (*Configs)(nil)
var _ yaml.NodeUnmarshalerContext = (*Configs)(nil)
var _ cv.Validatable = (*Configs)(nil)

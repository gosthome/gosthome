package binarysensor

import (
	"context"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/goccy/go-yaml"
	"github.com/gosthome/gosthome/core/component"
	cv "github.com/gosthome/gosthome/core/configvalidation"
	"github.com/gosthome/gosthome/core/entity"
	"github.com/gosthome/gosthome/core/registry"
)

type testBinarySensorConfig struct {
	BaseBinarySensorConfig[testBinarySensor, *testBinarySensor] `yaml:",inline"`

	Abc string `yaml:"abc"`
}

func (t *testBinarySensorConfig) ValidateWithContext(ctx context.Context) error {
	return cv.ValidateEmbedded(
		t.BaseBinarySensorConfig.ValidateWithContext(ctx),
		validation.ValidateStructWithContext(ctx, t,
			validation.Field(t.Abc, cv.String(cv.OneOf("yay", "nay"))),
		),
	)
}

type testBinarySensor struct {
	BaseBinarySensor[testBinarySensor, *testBinarySensor]
}

// Close implements component.Component.
func (t *testBinarySensor) Close() error {
	return nil
}

// InitializationPriority implements component.Component.
func (t *testBinarySensor) InitializationPriority() component.InitializationPriority {
	return component.InitializationPriorityProcessor
}

// Setup implements component.Component.
func (t *testBinarySensor) Setup() {
}

func newTestBinarySensor(ctx context.Context, cfg *testBinarySensorConfig) (*testBinarySensor, error) {
	ret := &testBinarySensor{}
	bs, err := NewBaseBinarySensor(ctx, ret, &cfg.BaseBinarySensorConfig)
	if err != nil {
		return nil, err
	}
	ret.BaseBinarySensor = *bs
	return ret, nil
}

var _ component.Component = (*testBinarySensor)(nil)
var _ entity.BinarySensor = (*testBinarySensor)(nil)

var _ yaml.NodeUnmarshalerContext = (*BaseBinarySensor[testBinarySensor, *testBinarySensor])(nil)

type testBinarySensorDeclaration struct {
}

// Component implements component.Declaration.
func (t *testBinarySensorDeclaration) Component(ctx context.Context, cfg component.Config) ([]component.Component, error) {
	panic("unimplemented")
}

// Config implements component.Declaration.
func (t *testBinarySensorDeclaration) Config() *component.ConfigDecoder {
	return &component.ConfigDecoder{
		Config:    &testBinarySensorConfig{},
		Marshal:   component.Marshal[testBinarySensorConfig],
		Unmarshal: component.Unmarshal[testBinarySensorConfig],
	}
}

var _ component.Declaration = (*testBinarySensorDeclaration)(nil)

func TestBaseBinaryComponent(t *testing.T) {
	reg := registry.NewRegistry()
	reg.RegisterEntityComponent(entity.DomainTypeBinarySensor, "test", &testBinarySensorDeclaration{})
}

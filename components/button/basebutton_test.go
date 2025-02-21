package button

import (
	"context"
	"log/slog"
	"testing"

	"github.com/gosthome/gosthome/core/component"
	cv "github.com/gosthome/gosthome/core/configvalidation"
	"github.com/gosthome/gosthome/core/entity"
	"github.com/gosthome/gosthome/core/registry"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type testButtonConfig struct {
	BaseButtonConfig[testButton, *testButton] `yaml:",inline"`

	Abc string `yaml:"abc"`
}

func (t *testButtonConfig) ValidateWithContext(ctx context.Context) error {
	return cv.ValidateEmbedded(
		t.BaseButtonConfig.ValidateWithContext(ctx),
		validation.ValidateStruct(t,
			validation.Field(t.Abc, cv.String(cv.OneOf("yay", "nay"))),
		),
	)
}

type testButton struct {
	BaseButton[testButton, *testButton]
}

func (t *testButton) Press(ctx context.Context) error {
	slog.Info("Button was pressed")
	return nil
}

// Close implements component.Component.
func (t *testButton) Close() error {
	return nil
}

// InitializationPriority implements component.Component.
func (t *testButton) InitializationPriority() component.InitializationPriority {
	return component.InitializationPriorityProcessor
}

// Setup implements component.Component.
func (t *testButton) Setup() {
}

func newTestButton(ctx context.Context, cfg *testButtonConfig) (*testButton, error) {
	ret := &testButton{}
	bs, err := NewBaseButton(ctx, ret, &cfg.BaseButtonConfig)
	if err != nil {
		return nil, err
	}
	ret.BaseButton = *bs
	return ret, nil
}

var _ component.Component = (*testButton)(nil)
var _ entity.Button = (*testButton)(nil)

type testButtonDeclaration struct {
}

// Component implements component.Declaration.
func (t *testButtonDeclaration) Component(ctx context.Context, cfg component.Config) ([]component.Component, error) {
	panic("unimplemented")
}

// Config implements component.Declaration.
func (t *testButtonDeclaration) Config() *component.ConfigDecoder {
	return &component.ConfigDecoder{
		Config:    &testButtonConfig{},
		Marshal:   component.Marshal[testButtonConfig],
		Unmarshal: component.Unmarshal[testButtonConfig],
	}
}

var _ component.Declaration = (*testButtonDeclaration)(nil)

func TestBaseBinaryComponent(t *testing.T) {
	reg := registry.NewRegistry()
	reg.RegisterEntityComponent(entity.DomainTypeButton, "test", &testButtonDeclaration{})
}

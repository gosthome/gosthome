package button

import (
	"context"

	"github.com/gosthome/gosthome/core/component"
	cv "github.com/gosthome/gosthome/core/configvalidation"
	"github.com/gosthome/gosthome/core/entity"
)

//go:generate go-enum --names --values

// ENUM(
// identify,//The button is used to identify a device.
// restart,//The button restarts the device.
// update,//The button updates the software of the device.
// )
type DeviceClass string

// DeviceClassValues implements entity.DeviceClassValues.
func (x *DeviceClass) DeviceClassValues() []string {
	return DeviceClassNames()
}

var _ (entity.DeviceClassValues) = (*DeviceClass)(nil)

type BaseButtonConfig[T any, PT interface {
	*T
	component.Component
	entity.Button
}] struct {
	component.ConfigOf[T, PT]
	entity.EntityConfig                                      `yaml:",inline"`
	entity.DeviceClassMixinConfig[DeviceClass, *DeviceClass] `yaml:",inline"`
	entity.IconMixinConfig                                   `yaml:",inline"`
}

func (bsc *BaseButtonConfig[T, PT]) ValidateWithContext(ctx context.Context) error {
	return cv.ValidateEmbedded(
		bsc.EntityConfig.ValidateWithContext(ctx),
		bsc.DeviceClassMixinConfig.ValidateWithContext(ctx),
		bsc.IconMixinConfig.ValidateWithContext(ctx))
}

type BaseButton[T any, PT interface {
	*T
	component.Component
	entity.Button
}] struct {
	entity.BaseEntity
	entity.DeviceClassMixin
	entity.IconMixin
}

func NewBaseButton[T any, PT interface {
	*T
	component.Component
	entity.Button
}](ctx context.Context, t PT, cfg *BaseButtonConfig[T, PT]) (*BaseButton[T, PT], error) {
	ret := &BaseButton[T, PT]{}
	ret.BaseEntity = *entity.NewBaseEntity(entity.DomainTypeButton, &cfg.EntityConfig)
	ret.DeviceClassMixin = entity.NewDeviceClassMixin(&cfg.DeviceClassMixinConfig)
	ret.IconMixin = entity.NewIconMixin(&cfg.IconMixinConfig)
	return ret, nil
}

package binarysensor

import (
	"context"

	"github.com/gosthome/gosthome/core/component"
	cv "github.com/gosthome/gosthome/core/configvalidation"
	"github.com/gosthome/gosthome/core/entity"
	"github.com/gosthome/gosthome/core/state"
)

//go:generate go-enum --names --values

// ENUM(
// battery,
// battery_charging,
// carbon_monoxide,
// cold,
// connectivity,
// door,
// empty,
// garage_door,
// gas,
// heat,
// light,
// lock,
// moisture,
// motion,
// moving,
// occupancy,
// opening,
// plug,
// power,
// presence,
// problem,
// running,
// safety,
// smoke,
// sound,
// tamper,
// update,
// vibration,
// window,
// )
type DeviceClass string

// DeviceClassValues implements entity.DeviceClassValues.
func (x *DeviceClass) DeviceClassValues() []string {
	return DeviceClassNames()
}

var _ (entity.DeviceClassValues) = (*DeviceClass)(nil)

type BaseBinarySensorConfig[T any, PT interface {
	*T
	component.Component
	entity.BinarySensor
}] struct {
	component.ConfigOf[T, PT]
	entity.EntityConfig                                      `yaml:",inline"`
	entity.DeviceClassMixinConfig[DeviceClass, *DeviceClass] `yaml:",inline"`
	entity.IconMixinConfig                                   `yaml:",inline"`
}

func (bsc *BaseBinarySensorConfig[T, PT]) ValidateWithContext(ctx context.Context) error {
	return cv.ValidateEmbedded(
		bsc.EntityConfig.ValidateWithContext(ctx),
		bsc.DeviceClassMixinConfig.ValidateWithContext(ctx),
		bsc.IconMixinConfig.ValidateWithContext(ctx),
	)
}

type BaseBinarySensor[T any, PT interface {
	*T
	component.Component
	entity.BinarySensor
}] struct {
	entity.BaseEntity
	entity.DeviceClassMixin
	entity.IconMixin
	state.State_[entity.BinarySensorState]
}

func NewBaseBinarySensor[T any, PT interface {
	*T
	component.Component
	entity.BinarySensor
}](ctx context.Context, t PT, cfg *BaseBinarySensorConfig[T, PT]) (*BaseBinarySensor[T, PT], error) {
	ret := &BaseBinarySensor[T, PT]{}
	ret.BaseEntity = *entity.NewBaseEntity(entity.DomainTypeBinarySensor, &cfg.EntityConfig)
	ret.DeviceClassMixin = entity.NewDeviceClassMixin(&cfg.DeviceClassMixinConfig)
	ret.IconMixin = entity.NewIconMixin(&cfg.IconMixinConfig)
	s, err := state.NewState(ctx, t, entity.BinarySensorState{
		State:   false,
		Missing: true,
	})
	if err != nil {
		return nil, err
	}
	ret.State_ = *s

	return ret, nil
}

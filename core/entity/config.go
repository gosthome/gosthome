package entity

import (
	"context"
	"log/slog"
	"regexp"

	cv "github.com/gosthome/gosthome/core/configvalidation"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type EntityConfig struct {
	ID                string         `yaml:"id"`
	Name              string         `yaml:"name"`
	Internal          bool           `yaml:"internal"`
	Category          EntityCategory `yaml:"category"`
	DisabledByDefault bool           `yaml:"disabled_by_default"`
}

func (ec *EntityConfig) ValidateWithContext(ctx context.Context) error {
	slog.Info("ec", "ec", ec)
	return validation.ValidateStructWithContext(ctx, ec,
		validation.Field(&ec.ID, validation.When(ec.Name != "").Else(validation.Required)),
		validation.Field(&ec.Name, validation.When(ec.ID != "").Else(validation.Required)),
	)
}

type IconMixinConfig struct {
	Icon string `yaml:"icon"`
}

var iconRE = regexp.MustCompile(`^$|^[\w\-]+:[\w\-]+$`)

// Validate implements validation.Validatable.
func (i *IconMixinConfig) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, i,
		validation.Field(&i.Icon, validation.Match(iconRE).Error(`Icons must match the format "[icon pack]:[icon]", e.g. "mdi:home-assistant"`)),
	)
}

type DeviceClassValues interface {
	DeviceClassValues() []string
}

type DeviceClassMixinConfig[Enum any, PE interface {
	DeviceClassValues
	*Enum
}] struct {
	DeviceClass string `yaml:"device_class"`
}

// Validate implements validation.Validatable.
func (d *DeviceClassMixinConfig[Enum, PE]) ValidateWithContext(ctx context.Context) error {
	values := PE(nil).DeviceClassValues()
	return validation.ValidateStructWithContext(ctx, d,
		validation.Field(&d.DeviceClass, cv.String(cv.Optional(cv.OneOf(values...)))))
}

type UnitOfMeasurementMixinConfig struct {
	UnitOfMeasurement string
}

// Validate implements validation.Validatable.
func (u *UnitOfMeasurementMixinConfig) ValidateWithContext(ctx context.Context) error {
	panic("unimplemented")
}

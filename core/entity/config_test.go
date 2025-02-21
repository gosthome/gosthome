package entity

import (
	cv "github.com/gosthome/gosthome/core/configvalidation"
)

var _ cv.Validatable = (*EntityConfig)(nil)
var _ cv.Validatable = (*IconMixinConfig)(nil)

type testDeviceClassEnum int

// DeviceClassValues implements DeviceClassValues.
func (t *testDeviceClassEnum) DeviceClassValues() []string {
	panic("unimplemented")
}

var _ DeviceClassValues = (*testDeviceClassEnum)(nil)

var _ cv.Validatable = (*DeviceClassMixinConfig[testDeviceClassEnum, *testDeviceClassEnum])(nil)
var _ cv.Validatable = (*UnitOfMeasurementMixinConfig)(nil)

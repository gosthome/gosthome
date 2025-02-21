package component

import (
	"context"
	"reflect"

	cv "github.com/gosthome/gosthome/core/configvalidation"
)

type Component interface {
	Setup()
	InitializationPriority() InitializationPriority
	Close() error
}

type Config interface {
	ComponentType() reflect.Type
	cv.Validatable
}

// AutoLoader is a component config showing need for another component during startup
// E.g. component is registering EntityComponents
type AutoLoader interface {
	AutoLoad() []string
}

type ConfigOf[T any, PT interface {
	*T
	Component
}] struct{}

func (*ConfigOf[T, PT]) ComponentType() reflect.Type {
	return reflect.TypeOf(PT(nil))
}

type Declaration interface {
	Config() *ConfigDecoder
	Component(ctx context.Context, cfg Config) ([]Component, error)
}

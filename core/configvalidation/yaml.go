package cv

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Validatable = validation.ValidatableWithContext

type Validator struct {
	Context context.Context
}

func (valid *Validator) Struct(v interface{}) error {
	validatable, ok := v.(Validatable)
	if !ok {
		return validation.NewInternalError(fmt.Errorf("the struct %T should be a gosthome/core/cv.Validatable with pointer reciever", v))
	}
	return ConvertErrors(validation.ValidateWithContext(valid.Context, validatable))
}

type ConfigYAMLDecoderKey struct{}
type ComponentRegistryKey struct{}

// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package entity

import (
	"errors"
	"fmt"
)

const (
	// CategoryNone is a Category of type None.
	CategoryNone Category = iota
	// CategoryConfig is a Category of type Config.
	CategoryConfig
	// CategoryDisgnostic is a Category of type Disgnostic.
	CategoryDisgnostic
)

var ErrInvalidCategory = errors.New("not a valid Category")

const _CategoryName = "noneconfigdisgnostic"

var _CategoryMap = map[Category]string{
	CategoryNone:       _CategoryName[0:4],
	CategoryConfig:     _CategoryName[4:10],
	CategoryDisgnostic: _CategoryName[10:20],
}

// String implements the Stringer interface.
func (x Category) String() string {
	if str, ok := _CategoryMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Category(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Category) IsValid() bool {
	_, ok := _CategoryMap[x]
	return ok
}

var _CategoryValue = map[string]Category{
	_CategoryName[0:4]:   CategoryNone,
	_CategoryName[4:10]:  CategoryConfig,
	_CategoryName[10:20]: CategoryDisgnostic,
}

// ParseCategory attempts to convert a string to a Category.
func ParseCategory(name string) (Category, error) {
	if x, ok := _CategoryValue[name]; ok {
		return x, nil
	}
	return Category(0), fmt.Errorf("%s is %w", name, ErrInvalidCategory)
}

// MarshalText implements the text marshaller method.
func (x Category) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Category) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseCategory(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

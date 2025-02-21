package cv

import (
	"fmt"
	"slices"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// type ConfigValue[T any] interface {
// 	ValidateWithContext(ctx context.Context) error
// 	Equal(other *ConfigPassword) bool
// 	Valid() bool
// 	MarshalText(text []byte) (err error)
// 	UnmarshalText(text []byte) (err error)
// }

type StringRule interface {
	Validate(value string) error
}

func String(rules ...StringRule) validation.Rule {
	return &stringRule{
		rules: rules,
	}
}

type stringRule struct {
	rules []StringRule
}

// Validate implements validation.Rule.
func (s *stringRule) Validate(ivalue interface{}) error {
	value, ok := ivalue.(string)
	if !ok {
		return validation.NewError("cv_not_a_string", "this value should be a string")
	}
	for _, rule := range s.rules {
		err := rule.Validate(value)
		if err != nil {
			return err
		}
	}
	return nil
}

func Optional(rules ...StringRule) StringRule {
	return &optionalRule{
		rules: rules,
	}
}

type optionalRule struct {
	rules []StringRule
}

// Validate implements validation.Rule.
func (o *optionalRule) Validate(value string) error {
	if value == "" {
		return nil
	}
	for _, rule := range o.rules {
		err := rule.Validate(value)
		if err != nil {
			return err
		}
	}
	return nil
}

type stringWithChars struct {
	chars string
}

// Validate implements StringRule.
func (n *stringWithChars) Validate(value string) error {
	for _, c := range value {
		if strings.ContainsRune(n.chars, c) {
			return validation.NewError("cv_string_has_illegal_chars", fmt.Sprintf(
				"'%c' is an invalid character for names. Valid characters are: %s (lowercase, no spaces)",
				c, ALLOWED_NAME_CHARS,
			))
		}
	}
	return nil
}

func Name() StringRule {
	return &stringWithChars{
		chars: ALLOWED_NAME_CHARS,
	}
}

func OneOf(variants ...string) StringRule {
	return &oneOfRule{
		variants: variants,
	}
}

type oneOfRule struct {
	variants []string
}

// Validate implements StringRule.
func (o *oneOfRule) Validate(value string) error {
	if slices.Contains(o.variants, value) {
		return nil
	}
	return validation.NewError("cv_not_one_of", fmt.Sprintf(
		"%s value should be one of %v", value, o.variants))
}

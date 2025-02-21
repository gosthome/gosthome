package cv

import (
	"errors"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

func ValidateEmbedded(errs ...error) error {
	nonverr := []error{}
	verrs := make(validation.Errors)
	for _, err := range errs {
		if err == nil {
			continue
		}
		verr, ok := err.(validation.Errors)
		if !ok {
			nonverr = append(nonverr, err)
			continue
		}
		for k, v := range verr {
			verrs[k] = v
		}
	}
	if len(nonverr) == 0 {
		if len(verrs) == 0 {
			return nil
		}
		return verrs
	}
	if len(verrs) != 0 {
		nonverr = append(nonverr, verrs)
	}
	return errors.Join(nonverr...)
}

func ErrDuplicateKey(node ast.Node, msg string, format ...any) error {
	return &yaml.DuplicateKeyError{
		Message: fmt.Sprintf(msg, format...),
		Token:   node.GetToken(),
	}
}

func ErrUnknownField(node ast.Node, msg string, format ...any) error {
	if node == nil || node.GetToken() == nil {
		return fmt.Errorf(msg, format...)
		// return fmt.Errorf("Node@%q \n"+msg, append([]any{node.GetPath()}, format)...)
	}
	return &yaml.UnknownFieldError{
		Message: fmt.Sprintf(msg, format...),
		Token:   node.GetToken(),
	}
}

type Error struct {
	err       error
	fieldName string
}

// Error implements error.
func (e *Error) Error() string {
	// return e.err.Error()
	return fmt.Sprintf("field %s was wrong: %s", e.fieldName, e.err.Error())
}

// StructField implements yaml.FieldError.
func (e Error) StructField() string {
	return e.fieldName
}

type Errors []Error

// Error implements error.
func (e Errors) Error() string {
	b := strings.Builder{}
	for _, err := range e {
		b.WriteString(err.Error())
		b.WriteByte('\n')
	}
	return b.String()
}

func ConvertErrors(err error) error {
	verrs, ok := err.(validation.Errors)
	if !ok {
		return err
	}
	ret := Errors{}

	for k, verr := range verrs {
		ret = append(ret, Error{
			err:       verr,
			fieldName: k,
		})
	}
	return ret
}

var _ error = (Errors)(nil)
var _ error = (*Error)(nil)
var _ yaml.FieldError = (*Error)(nil)

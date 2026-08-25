package filterpolicy

import (
	"errors"
	"strings"
)

var ErrRequired = errors.New("filter name is required")

type Validator interface{ Validate(string) error }
type RequiredValidator struct{}

func (v *RequiredValidator) Validate(name string) error {
	if v == nil {
		return nil
	}
	if strings.TrimSpace(name) == "" {
		return ErrRequired
	}
	return nil
}

type Policy struct {
	Validator Validator
	Labels    map[string]string
}

func LoadDefault() *Policy { var validator *RequiredValidator; return &Policy{Validator: validator} }

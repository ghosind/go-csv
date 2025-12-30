package csv

import (
	"errors"
	"reflect"
)

var (
	ErrInvalidType      = errors.New("csv: invalid type")
	ErrUnsupportedType  = errors.New("csv: unsupported type")
	ErrCannotSet        = errors.New("csv: cannot set value to nil pointer")
	ErrInvalidUnmarshal = errors.New("csv: Unmarshal(nil)")
)

func newInvalidUnmarshalError(rv reflect.Value) error {
	if !rv.IsValid() {
		return ErrInvalidUnmarshal
	}

	return errors.New("csv: Unmarshal(nil " + rv.Type().Name() + ")")
}

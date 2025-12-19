package csv

import (
	"errors"
	"reflect"
)

var (
	ErrInvalidType     = errors.New("invalid type for CSV operation")
	ErrUnsupportedType = errors.New("unsupported type for CSV operation")
)

func newInvalidUnmarshalError(rv reflect.Value) error {
	if !rv.IsValid() {
		return errors.New("csv: Unmarshal(nil)")
	}

	return errors.New("csv: Unmarshal(nil " + rv.Type().Name() + ")")
}

package csv

import (
	"errors"
	"fmt"
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

type DecodeError struct {
	line   int
	column int
	field  string
	value  string
	err    error
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("csv: line %d, column %d (field: %s, value: %s): %v",
		e.line, e.column, e.field, e.value, e.err)
}

func (e *DecodeError) Unwrap() error {
	return e.err
}

func (e *DecodeError) Row() int {
	return e.line
}

func (e *DecodeError) Col() int {
	return e.column
}

func (e *DecodeError) Field() string {
	return e.field
}

func (e *DecodeError) Value() string {
	return e.value
}

func newDecodeError(line, column int, field, value string, err error) *DecodeError {
	return &DecodeError{
		line:   line,
		column: column,
		field:  field,
		value:  value,
		err:    err,
	}
}

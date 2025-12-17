package csv

import "errors"

var (
	ErrInvalidType     = errors.New("invalid type for CSV operation")
	ErrUnsupportedType = errors.New("unsupported type for CSV operation")
)

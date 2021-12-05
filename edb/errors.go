package edb

import "errors"

// package errors
var (
	ErrNotFound           = errors.New("edb: not found")
	ErrNotSupported       = errors.New("edb: not supported")
	ErrTableNotSpecified  = errors.New("edb: table not specified")
	ErrColumnNotSpecified = errors.New("edb: column not specified")
	ErrInvalidPointer     = errors.New("edb: attempt to load into an invalid pointer")
	ErrPlaceholderCount   = errors.New("edb: wrong placeholder count")
	ErrInvalidSliceLength = errors.New("edb: length of slice is 0. length must be >= 1")
	ErrCantConvertToTime  = errors.New("edb: can't convert to time.Time")
	ErrInvalidTimestring  = errors.New("edb: invalid time string")
)

package errors

import "errors"

var (
	ErrByteOverFlow = errors.New("no of bytes exceed capacity for an individual value ")
	ErrIOConnNil    = errors.New("IO connection is nil")
)

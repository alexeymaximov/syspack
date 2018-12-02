package mmap

import (
	"fmt"

	"github.com/alexeymaximov/syspack"
)

// Error occurred when mapping closed.
type ErrorClosed struct{}

// Get error message.
func (err *ErrorClosed) Error() string {
	return "mmap: mapping closed"
}

// Error occurred when mapping mode is invalid.
type ErrorInvalidMode struct{ Mode Mode }

// Get error message.
func (err *ErrorInvalidMode) Error() string {
	return fmt.Sprintf("mmap: invalid mode 0x%x", err.Mode)
}

// Error occurred when offset is invalid.
type ErrorInvalidOffset struct{ Offset syspack.Offset }

// Get error message.
func (err *ErrorInvalidOffset) Error() string {
	return fmt.Sprintf("mmap: invalid offset 0x%x", err.Offset)
}

// Error occurred when offset range is invalid.
type ErrorInvalidOffsetRange struct{ Low, High syspack.Offset }

// Get error message.
func (err *ErrorInvalidOffsetRange) Error() string {
	return fmt.Sprintf("mmap: invalid offset range 0x%x..0x%x", err.Low, err.High)
}

// Error occurred when size is invalid.
type ErrorInvalidSize struct{ Size syspack.Size }

// Get error message.
func (err *ErrorInvalidSize) Error() string {
	return fmt.Sprintf("mmap: invalid size %d", err.Size)
}

// Error occurred when operation is not allowed.
type ErrorNotAllowed struct{ Operation string }

// Get error message.
func (err *ErrorNotAllowed) Error() string {
	return fmt.Sprintf("mmap: %s is not allowed", err.Operation)
}

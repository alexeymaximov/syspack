package mmap

import (
	"io"

	"github.com/alexeymaximov/syspack"
)

// Memory access mode.
type Mode int

// Available modes.
const (
	ModeReadOnly Mode = iota
	ModeReadWrite
	ModeReadWritePrivate
)

type Options struct {
	// Mapping options.

	// Memory access mode.
	Mode Mode

	// Memory is executable.
	Executable bool
}

// Get mapping length.
func (mapping *Mapping) Len() int {
	return len(mapping.data)
}

// Whether is reading from mapping allowed.
func (mapping *Mapping) CanRead() bool {
	return true
}

// Whether is writing to mapping allowed.
func (mapping *Mapping) CanWrite() bool {
	return mapping.canWrite
}

// Whether is mapping execution allowed.
func (mapping *Mapping) CanExecute() bool {
	return mapping.canExecute
}

// Get direct byte slice in offset range [low, high).
func (mapping *Mapping) Direct(low, high syspack.Offset) ([]byte, error) {
	if mapping.data == nil {
		return nil, &ErrorClosed{}
	}
	length := syspack.Off(mapping.data)
	if low < 0 || low >= length {
		return nil, &ErrorInvalidOffset{Offset: low}
	}
	if high < 1 || high > length {
		return nil, &ErrorInvalidOffset{Offset: high}
	}
	if low >= high {
		return nil, &ErrorInvalidOffsetRange{Low: low, High: high - 1}
	}
	return mapping.data[low:high], nil
}

// Read single byte from mapping at given offset.
func (mapping *Mapping) ReadByteAt(offset syspack.Offset) (byte, error) {
	if mapping.data == nil {
		return 0, &ErrorClosed{}
	}
	if offset < 0 || offset >= syspack.Off(mapping.data) {
		return 0, &ErrorInvalidOffset{Offset: offset}
	}
	return mapping.data[offset], nil
}

// Write single byte to mapping at given offset.
func (mapping *Mapping) WriteByteAt(byte byte, offset syspack.Offset) error {
	if mapping.data == nil {
		return &ErrorClosed{}
	}
	if !mapping.canWrite {
		return &ErrorNotAllowed{Operation: "write"}
	}
	if offset < 0 || offset >= syspack.Off(mapping.data) {
		return &ErrorInvalidOffset{Offset: offset}
	}
	mapping.data[offset] = byte
	return nil
}

// Read len(buffer) bytes from mapping at given offset.
func (mapping *Mapping) ReadAt(buffer []byte, offset syspack.Offset) (int, error) {
	if mapping.data == nil {
		return 0, &ErrorClosed{}
	}
	if offset < 0 || offset >= syspack.Off(mapping.data) {
		return 0, &ErrorInvalidOffset{Offset: offset}
	}
	n := copy(buffer, mapping.data[offset:])
	if n < len(buffer) {
		return n, io.EOF
	}
	return n, nil
}

// Write len(buffer) bytes to mapping at given offset.
func (mapping *Mapping) WriteAt(buffer []byte, offset syspack.Offset) (int, error) {
	if mapping.data == nil {
		return 0, &ErrorClosed{}
	}
	if !mapping.canWrite {
		return 0, &ErrorNotAllowed{Operation: "write"}
	}
	if offset < 0 || offset >= syspack.Off(mapping.data) {
		return 0, &ErrorInvalidOffset{Offset: offset}
	}
	n := copy(mapping.data[offset:], buffer)
	if n < len(buffer) {
		return n, io.EOF
	}
	return n, nil
}

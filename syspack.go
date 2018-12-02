package syspack

import (
	"math"
	"syscall"
)

type (
	Dword  = uint32
	Offset = int64
	Qword  = uint64
	Size   = uintptr
)

const (
	MaxDword   = Dword(math.MaxUint32)
	MaxInt     = int(MaxUint >> 1)
	MaxOffset  = Offset(math.MaxInt64)
	MaxQword   = Qword(math.MaxUint64)
	MaxSize    = ^Size(0)
	MaxUint    = ^uint(0)
	MaxUintptr = ^uintptr(0)
	MinDword   = Dword(0)
	MinInt     = -MaxInt - 1
	MinOffset  = Offset(math.MinInt64)
	MinQword   = Qword(0)
	MinSize    = Size(0)
	MinUint    = uint(0)
	MinUintptr = uintptr(0)
)

// Get system call error number.
func Errno(err error) error {
	if err != nil {
		if errno, ok := err.(syscall.Errno); ok && errno == 0 {
			return syscall.EINVAL
		}
		return err
	}
	return syscall.EINVAL
}

// Get high and low bits.
func HighLow(value Qword) (high, low Dword) {
	high = Dword(value >> 32)
	low = Dword(value & Qword(MaxDword))
	return
}

// Get length of byte array.
func Len(value []byte) Size {
	return Size(len(value))
}
func Off(value []byte) Offset {
	return Offset(len(value))
}

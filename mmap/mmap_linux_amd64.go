package mmap

import (
	"os"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/alexeymaximov/syspack"
)

type Mapping struct {
	// Memory mapping.

	// Aligned address.
	alignedAddress uintptr

	// Aligned size.
	alignedSize syspack.Size

	// Data.
	data []byte

	// Writing is allowed.
	canWrite bool

	// Execution is allowed.
	canExecute bool
}

// Make new mapping.
func NewMapping(fd uintptr, offset syspack.Offset, size syspack.Size, options *Options) (*Mapping, error) {
	if offset < 0 {
		return nil, &ErrorInvalidOffset{Offset: offset}
	}
	if size > syspack.Size(syspack.MaxInt) {
		return nil, &ErrorInvalidSize{Size: size}
	}
	mapping := &Mapping{}
	protection := syscall.PROT_READ
	flags := syscall.MAP_SHARED
	if options != nil {
		if options.Mode < ModeReadOnly || options.Mode > ModeReadWritePrivate {
			return nil, &ErrorInvalidMode{Mode: options.Mode}
		}
		if options.Mode > ModeReadOnly {
			protection |= syscall.PROT_WRITE
			mapping.canWrite = true
		}
		if options.Mode == ModeReadWritePrivate {
			flags = syscall.MAP_PRIVATE
		}
		if options.Executable {
			protection |= syscall.PROT_EXEC
			mapping.canExecute = true
		}
	}
	pageSize := syspack.Offset(os.Getpagesize())
	if pageSize < 0 {
		return nil, os.NewSyscallError("getpagesize", syscall.EINVAL)
	}
	outerOffset := offset / pageSize
	innerOffset := offset % pageSize
	mapping.alignedSize = syspack.Size(innerOffset) + size
	var err error
	mapping.alignedAddress, err = syspack.MmapE(0, mapping.alignedSize, protection, flags, fd, outerOffset)
	if err != nil {
		return nil, err
	}
	var sliceHeader struct {
		data uintptr
		len  int
		cap  int
	}
	sliceHeader.data = mapping.alignedAddress + uintptr(innerOffset)
	sliceHeader.len = int(size)
	sliceHeader.cap = sliceHeader.len
	mapping.data = *(*[]byte)(unsafe.Pointer(&sliceHeader))
	runtime.SetFinalizer(mapping, (*Mapping).Close)
	return mapping, nil
}

// Lock mapping.
func (mapping *Mapping) Lock() error {
	if mapping.data == nil {
		return &ErrorClosed{}
	}
	return syspack.MlockE(mapping.alignedAddress, mapping.alignedSize)
}

// Unlock mapping.
func (mapping *Mapping) Unlock() error {
	if mapping.data == nil {
		return &ErrorClosed{}
	}
	return syspack.MunlockE(mapping.alignedAddress, mapping.alignedSize)
}

// Sync mapping.
func (mapping *Mapping) Sync() error {
	if mapping.data == nil {
		return &ErrorClosed{}
	}
	if !mapping.canWrite {
		return &ErrorNotAllowed{Operation: "sync"}
	}
	return syspack.MsyncE(mapping.alignedAddress, mapping.alignedSize)
}

// Close mapping.
func (mapping *Mapping) Close() error {
	if mapping.data == nil {
		return &ErrorClosed{}
	}

	// Maybe unnecessary.
	if mapping.canWrite {
		if err := mapping.Sync(); err != nil {
			return err
		}
	}

	if err := syspack.MunmapE(mapping.alignedAddress, mapping.alignedSize); err != nil {
		return err
	}
	mapping.data = nil
	runtime.SetFinalizer(mapping, nil)
	return nil
}

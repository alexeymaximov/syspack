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

	// Process handle.
	hProcess syspack.Handle

	// File handle.
	hFile syspack.Handle

	// Mapping handle.
	hMapping syspack.Handle

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
	protection := syspack.Dword(syscall.PAGE_READONLY)
	access := syspack.Dword(syscall.FILE_MAP_READ)
	if options != nil {
		switch options.Mode {
		case ModeReadOnly:
			// NOOP
		case ModeReadWrite:
			protection = syscall.PAGE_READWRITE
			access = syscall.FILE_MAP_WRITE
			mapping.canWrite = true
		case ModeReadWritePrivate:
			protection = syscall.PAGE_WRITECOPY
			access = syscall.FILE_MAP_COPY
			mapping.canWrite = true
		default:
			return nil, &ErrorInvalidMode{Mode: options.Mode}
		}
		if options.Executable {
			protection <<= 4
			access |= syscall.FILE_MAP_EXECUTE
			mapping.canExecute = true
		}
	}
	var err error
	mapping.hProcess, err = syspack.GetCurrentProcessE()
	if err != nil {
		return nil, err
	}
	err = syspack.DuplicateHandleE(
		mapping.hProcess, syspack.Handle(fd),
		mapping.hProcess, &mapping.hFile,
		0, true, syscall.DUPLICATE_SAME_ACCESS,
	)
	if err != nil {
		return nil, err
	}
	pageSize := syspack.Offset(os.Getpagesize())
	if pageSize < 0 {
		return nil, os.NewSyscallError("getpagesize", syscall.EINVAL)
	}
	outerOffset := offset / pageSize
	innerOffset := offset % pageSize
	mapping.alignedSize = syspack.Size(innerOffset) + size
	maxSizeHigh, maxSizeLow := syspack.HighLow(syspack.Qword(outerOffset) + syspack.Qword(mapping.alignedSize))
	mapping.hMapping, err = syspack.CreateFileMappingE(mapping.hFile, nil, protection, maxSizeHigh, maxSizeLow, nil)
	if err != nil {
		return nil, err
	}
	offsetHigh, offsetLow := syspack.HighLow(syspack.Qword(outerOffset))
	mapping.alignedAddress, err = syspack.MapViewOfFileE(
		mapping.hMapping, access,
		offsetHigh, offsetLow, mapping.alignedSize,
	)
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

// Ensure process quota for mapping.
func (mapping *Mapping) EnsureQuota() error {
	if mapping.data == nil {
		return &ErrorClosed{}
	}
	var minSize, maxSize syspack.Size
	if err := syspack.GetProcessWorkingSetSizeE(mapping.hProcess, &minSize, &maxSize); err != nil {
		return err
	}
	minSize += mapping.alignedSize
	maxSize += mapping.alignedSize
	if err := syspack.SetProcessWorkingSetSizeE(mapping.hProcess, minSize, maxSize); err != nil {
		return err
	}
	return nil
}

// Lock mapping.
func (mapping *Mapping) Lock() error {
	if mapping.data == nil {
		return &ErrorClosed{}
	}
	return syspack.VirtualLockE(mapping.alignedAddress, mapping.alignedSize)
}

// Unlock mapping.
func (mapping *Mapping) Unlock() error {
	if mapping.data == nil {
		return &ErrorClosed{}
	}
	return syspack.VirtualUnlockE(mapping.alignedAddress, mapping.alignedSize)
}

// Sync mapping.
func (mapping *Mapping) Sync() error {
	if mapping.data == nil {
		return &ErrorClosed{}
	}
	if !mapping.canWrite {
		return &ErrorNotAllowed{Operation: "sync"}
	}
	if err := syspack.FlushViewOfFileE(mapping.alignedAddress, mapping.alignedSize); err != nil {
		return err
	}
	if err := syspack.FlushFileBuffersE(mapping.hFile); err != nil {
		return err
	}
	return nil
}

// Close mapping.
func (mapping *Mapping) Close() error {
	if mapping.data == nil {
		return &ErrorClosed{}
	}
	if mapping.canWrite {
		if err := mapping.Sync(); err != nil {
			return err
		}
	}
	if err := syspack.UnmapViewOfFileE(mapping.alignedAddress); err != nil {
		return err
	}
	if err := syspack.CloseHandle(mapping.hMapping); err != nil {
		return err
	}
	if err := syspack.CloseHandle(mapping.hFile); err != nil {
		return err
	}
	mapping.data = nil
	runtime.SetFinalizer(mapping, nil)
	return nil
}

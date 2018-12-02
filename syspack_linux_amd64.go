package syspack

import (
	"os"
	"syscall"
)

const (
	SymbolMlock   = "mlock"
	SymbolMmap    = "mmap"
	SymbolMsync   = "msync"
	SymbolMunlock = "munlock"
	SymbolMunmap  = "munmap"
)

func Mlock(addr uintptr, length Size) error {
	_, _, err := syscall.Syscall(syscall.SYS_MLOCK, addr, length, 0)
	if err != 0 {
		return Errno(err)
	}
	return err
}
func MlockE(addr uintptr, length Size) error {
	return os.NewSyscallError(SymbolMlock, Mlock(addr, length))
}

func Mmap(addr uintptr, length Size, prot, flags int, fd uintptr, offset Offset) (uintptr, error) {
	if prot < 0 || flags < 0 || offset < 0 {
		return 0, syscall.EINVAL
	}
	result, _, err := syscall.Syscall6(syscall.SYS_MMAP, addr, length, uintptr(prot), uintptr(flags), fd, uintptr(offset))
	if err != 0 {
		return 0, Errno(err)
	}
	return result, nil
}
func MmapE(addr uintptr, length Size, prot, flags int, fd uintptr, offset Offset) (uintptr, error) {
	memory, err := Mmap(addr, length, prot, flags, fd, offset)
	if err != nil {
		return memory, os.NewSyscallError(SymbolMmap, err)
	}
	return memory, nil
}

func Msync(addr uintptr, length Size) error {
	_, _, err := syscall.Syscall(syscall.SYS_MSYNC, addr, length, syscall.MS_SYNC)
	if err != 0 {
		return Errno(err)
	}
	return nil
}
func MsyncE(addr uintptr, length Size) error {
	return os.NewSyscallError(SymbolMsync, Msync(addr, length))
}

func Munlock(addr uintptr, length Size) error {
	_, _, err := syscall.Syscall(syscall.SYS_MUNLOCK, addr, length, 0)
	if err != 0 {
		return Errno(err)
	}
	return nil
}
func MunlockE(addr uintptr, length Size) error {
	return os.NewSyscallError(SymbolMunlock, Munlock(addr, length))
}

func Munmap(addr uintptr, length Size) error {
	_, _, err := syscall.Syscall(syscall.SYS_MUNMAP, addr, length, 0)
	if err != 0 {
		return Errno(err)
	}
	return nil
}
func MunmapE(addr uintptr, length Size) error {
	return os.NewSyscallError(SymbolMunmap, Munmap(addr, length))
}

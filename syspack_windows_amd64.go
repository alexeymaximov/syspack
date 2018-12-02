package syspack

import (
	"os"
	"syscall"
	"unsafe"
)

const (
	SymbolKernel32 = "kernel32.dll"

	SymbolCloseHandle              = "CloseHandle"
	SymbolCreateFileMapping        = "CreateFileMapping"
	SymbolDuplicateHandle          = "DuplicateHandle"
	SymbolFlushFileBuffers         = "FlushFileBuffers"
	SymbolFlushViewOfFile          = "FlushViewOfFile"
	SymbolGenerateConsoleCtrlEvent = "GenerateConsoleCtrlEvent"
	SymbolGetCurrentProcess        = "GetCurrentProcess"
	SymbolGetProcessWorkingSetSize = "GetProcessWorkingSetSize"
	SymbolMapViewOfFile            = "MapViewOfFile"
	SymbolSetProcessWorkingSetSize = "SetProcessWorkingSetSize"
	SymbolUnmapViewOfFile          = "UnmapViewOfFile"
	SymbolVirtualLock              = "VirtualLock"
	SymbolVirtualUnlock            = "VirtualUnlock"
)

var (
	modkernel32 = syscall.NewLazyDLL(SymbolKernel32)

	procGenerateConsoleCtrlEvent = modkernel32.NewProc(SymbolGenerateConsoleCtrlEvent)
	procGetProcessWorkingSetSize = modkernel32.NewProc(SymbolGetProcessWorkingSetSize)
	procSetProcessWorkingSetSize = modkernel32.NewProc(SymbolSetProcessWorkingSetSize)
)

type Handle = syscall.Handle

func CloseHandle(hObject Handle) error {
	return syscall.CloseHandle(hObject)
}
func CloseHandleE(hObject Handle) error {
	return os.NewSyscallError(SymbolCloseHandle, syscall.CloseHandle(hObject))
}

func CreateFileMapping(
	hFile Handle, lpFileMappingAttributes *syscall.SecurityAttributes,
	flProtect, dwMaximumSizeHigh, dwMaximumSizeLow Dword, lpName *uint16,
) (Handle, error) {
	return syscall.CreateFileMapping(
		hFile, lpFileMappingAttributes,
		flProtect, dwMaximumSizeHigh, dwMaximumSizeLow, lpName,
	)
}
func CreateFileMappingE(
	hFile Handle, lpFileMappingAttributes *syscall.SecurityAttributes,
	flProtect, dwMaximumSizeHigh, dwMaximumSizeLow Dword, lpName *uint16,
) (Handle, error) {
	handle, err := syscall.CreateFileMapping(
		hFile, lpFileMappingAttributes,
		flProtect, dwMaximumSizeHigh, dwMaximumSizeLow, lpName,
	)
	if err != nil {
		return handle, os.NewSyscallError(SymbolCreateFileMapping, err)
	}
	return handle, nil
}

func DuplicateHandle(
	hSourceProcessHandle, hSourceHandle, hTargetProcessHandle Handle,
	lpTargetHandle *Handle, dwDesiredAccess Dword, bInheritHandle bool, dwOptions Dword,
) error {
	return syscall.DuplicateHandle(
		hSourceProcessHandle, hSourceHandle, hTargetProcessHandle,
		lpTargetHandle, dwDesiredAccess, bInheritHandle, dwOptions,
	)
}
func DuplicateHandleE(
	hSourceProcessHandle, hSourceHandle, hTargetProcessHandle Handle,
	lpTargetHandle *Handle, dwDesiredAccess Dword, bInheritHandle bool, dwOptions Dword,
) error {
	return os.NewSyscallError(
		SymbolDuplicateHandle,
		syscall.DuplicateHandle(
			hSourceProcessHandle, hSourceHandle, hTargetProcessHandle,
			lpTargetHandle, dwDesiredAccess, bInheritHandle, dwOptions,
		),
	)
}

func FlushFileBuffers(hFile Handle) error {
	return syscall.FlushFileBuffers(hFile)
}
func FlushFileBuffersE(hFile Handle) error {
	return os.NewSyscallError(SymbolFlushFileBuffers, syscall.FlushFileBuffers(hFile))
}

func FlushViewOfFile(lpBaseAddress uintptr, dwNumberOfBytesToFlush Size) error {
	return syscall.FlushViewOfFile(lpBaseAddress, dwNumberOfBytesToFlush)
}
func FlushViewOfFileE(lpBaseAddress uintptr, dwNumberOfBytesToFlush Size) error {
	return os.NewSyscallError(SymbolFlushViewOfFile, syscall.FlushViewOfFile(lpBaseAddress, dwNumberOfBytesToFlush))
}

func GenerateConsoleCtrlEvent(dwCtrlEvent, dwProcessGroupId Dword) error {
	result, _, err := procGenerateConsoleCtrlEvent.Call(uintptr(dwCtrlEvent), uintptr(dwProcessGroupId))
	if result == 0 {
		return Errno(err)
	}
	return nil
}
func GenerateConsoleCtrlEventE(dwCtrlEvent, dwProcessGroupId Dword) error {
	return os.NewSyscallError(
		SymbolGenerateConsoleCtrlEvent,
		GenerateConsoleCtrlEvent(dwCtrlEvent, dwProcessGroupId),
	)
}

func GetCurrentProcess() (Handle, error) {
	return syscall.GetCurrentProcess()
}
func GetCurrentProcessE() (Handle, error) {
	pseudoHandle, err := syscall.GetCurrentProcess()
	if err != nil {
		return pseudoHandle, os.NewSyscallError(SymbolGetCurrentProcess, err)
	}
	return pseudoHandle, nil
}

func GetProcessWorkingSetSize(hProcess Handle, lpMinimumWorkingSetSize, lpMaximumWorkingSetSize *Size) error {
	result, _, err := procGetProcessWorkingSetSize.Call(
		uintptr(hProcess),
		uintptr(unsafe.Pointer(lpMinimumWorkingSetSize)),
		uintptr(unsafe.Pointer(lpMaximumWorkingSetSize)),
	)
	if result == 0 {
		return Errno(err)
	}
	return nil
}
func GetProcessWorkingSetSizeE(hProcess Handle, lpMinimumWorkingSetSize, lpMaximumWorkingSetSize *Size) error {
	return os.NewSyscallError(
		SymbolGetProcessWorkingSetSize,
		GetProcessWorkingSetSize(hProcess, lpMinimumWorkingSetSize, lpMaximumWorkingSetSize),
	)
}

func MapViewOfFile(
	hFileMappingObject Handle, dwDesiredAccess, dwFileOffsetHigh, dwFileOffsetLow Dword,
	dwNumberOfBytesToMap Size,
) (uintptr, error) {
	return syscall.MapViewOfFile(
		hFileMappingObject, dwDesiredAccess, dwFileOffsetHigh, dwFileOffsetLow,
		dwNumberOfBytesToMap,
	)
}
func MapViewOfFileE(
	hFileMappingObject Handle, dwDesiredAccess, dwFileOffsetHigh, dwFileOffsetLow Dword,
	dwNumberOfBytesToMap Size,
) (uintptr, error) {
	addr, err := syscall.MapViewOfFile(
		hFileMappingObject, dwDesiredAccess, dwFileOffsetHigh, dwFileOffsetLow,
		dwNumberOfBytesToMap,
	)
	if err != nil {
		return addr, os.NewSyscallError(SymbolMapViewOfFile, err)
	}
	return addr, nil
}

func SetProcessWorkingSetSize(hProcess Handle, dwMinimumWorkingSetSize, dwMaximumWorkingSetSize Size) error {
	result, _, err := procSetProcessWorkingSetSize.Call(
		uintptr(hProcess),
		dwMinimumWorkingSetSize,
		dwMaximumWorkingSetSize,
	)
	if result == 0 {
		return Errno(err)
	}
	return nil
}
func SetProcessWorkingSetSizeE(hProcess Handle, dwMinimumWorkingSetSize, dwMaximumWorkingSetSize Size) error {
	return os.NewSyscallError(
		SymbolSetProcessWorkingSetSize,
		SetProcessWorkingSetSize(hProcess, dwMinimumWorkingSetSize, dwMaximumWorkingSetSize),
	)
}

func UnmapViewOfFile(lpBaseAddress uintptr) error {
	return syscall.UnmapViewOfFile(lpBaseAddress)
}
func UnmapViewOfFileE(lpBaseAddress uintptr) error {
	return os.NewSyscallError(SymbolUnmapViewOfFile, syscall.UnmapViewOfFile(lpBaseAddress))
}

func VirtualLock(lpAddress uintptr, dwSize Size) error {
	return syscall.VirtualLock(lpAddress, dwSize)
}
func VirtualLockE(lpAddress uintptr, dwSize Size) error {
	return os.NewSyscallError(SymbolVirtualLock, syscall.VirtualLock(lpAddress, dwSize))
}

func VirtualUnlock(lpAddress uintptr, dwSize Size) error {
	return syscall.VirtualUnlock(lpAddress, dwSize)
}
func VirtualUnlockE(lpAddress uintptr, dwSize Size) error {
	return os.NewSyscallError(SymbolVirtualUnlock, syscall.VirtualUnlock(lpAddress, dwSize))
}

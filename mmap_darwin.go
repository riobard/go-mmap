package mmap

import (
	"syscall"
	"unsafe"
)

const (
	MAP_NOCACHE      MapFlag = syscall.MAP_NOCACHE
	MAP_HASSEMAPHORE         = syscall.MAP_HASSEMAPHORE
)

const (
	MADV_FREE             Advice = syscall.MADV_FREE
	MADV_ZERO_WIRED_PAGES        = syscall.MADV_ZERO_WIRED_PAGES
)

const (
	MINCORE_REFERENCED       MincoreState = 0x2
	MINCORE_MODIFIED                      = 0x4
	MINCORE_REFERENCED_OTHER              = 0x8
	MINCORE_MODIFIED_OTHER                = 0x10
)

func (m Mmap) Advise(advice Advice) error {
	_, _, errno := syscall.Syscall(syscall.SYS_MADVISE, uintptr(unsafe.Pointer(&m[0])), uintptr(len(m)), uintptr(advice))
	if errno != 0 {
		return errno
	}
	return nil
}

func (m Mmap) Lock() error {
	_, _, errno := syscall.Syscall(syscall.SYS_MLOCK, uintptr(unsafe.Pointer(&m[0])), uintptr(len(m)), 0)
	if errno != 0 {
		return errno
	}
	return nil
}

func (m Mmap) Unlock() error {
	_, _, errno := syscall.Syscall(syscall.SYS_MUNLOCK, uintptr(unsafe.Pointer(&m[0])), uintptr(len(m)), 0)
	if errno != 0 {
		return errno
	}
	return nil
}

func (m Mmap) Protect(prot Prot) error {
	_, _, errno := syscall.Syscall(syscall.SYS_MPROTECT, uintptr(unsafe.Pointer(&m[0])), uintptr(len(m)), uintptr(prot))
	if errno != 0 {
		return errno
	}
	return nil
}

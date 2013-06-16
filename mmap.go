package mmap

import (
	"os"
	"syscall"
	"unsafe"
)

type Prot int

const (
	PROT_NONE  Prot = syscall.PROT_NONE
	PROT_READ       = syscall.PROT_READ
	PROT_WRITE      = syscall.PROT_WRITE
	PROT_EXEC       = syscall.PROT_EXEC
)

type MapFlag int

const (
	MAP_ANON    MapFlag = syscall.MAP_ANON
	MAP_FILE            = syscall.MAP_FILE
	MAP_FIXED           = syscall.MAP_FIXED
	MAP_PRIVATE         = syscall.MAP_PRIVATE
	MAP_SHARED          = syscall.MAP_SHARED
)

type SyncFlag int

const (
	MS_ASYNC      SyncFlag = syscall.MS_ASYNC
	MS_SYNC                = syscall.MS_SYNC
	MS_INVALIDATE          = syscall.MS_INVALIDATE
)

type Advice int

const (
	MADV_DONTNEED   Advice = syscall.MADV_DONTNEED
	MADV_NORMAL            = syscall.MADV_NORMAL
	MADV_RANDOM            = syscall.MADV_RANDOM
	MADV_SEQUENTIAL        = syscall.MADV_SEQUENTIAL
	MADV_WILLNEED          = syscall.MADV_WILLNEED
)

type MincoreState byte

const (
	MINCORE_INCORE MincoreState = 0x1
)

type Mmap []byte

func Map(f *os.File, offset int64, len int, prot Prot, flags MapFlag) (Mmap, error) {
	return syscall.Mmap(int(f.Fd()), offset, len, int(prot), int(flags))
}

func (m Mmap) Unmap() error {
	err := syscall.Munmap(m)
	m = nil
	return err
}

func (m Mmap) Sync(flags SyncFlag) error {
	_, _, errno := syscall.Syscall(syscall.SYS_MSYNC, uintptr(unsafe.Pointer(&m[0])), uintptr(len(m)), uintptr(flags))
	if errno != 0 {
		return errno
	}
	return nil
}

func (m Mmap) Incore() ([]MincoreState, error) {
	pageSize := os.Getpagesize()
	vec := make([]MincoreState, (len(m)+pageSize-1)/pageSize)
	_, _, errno := syscall.Syscall(syscall.SYS_MINCORE, uintptr(unsafe.Pointer(&m[0])), uintptr(len(m)), uintptr(unsafe.Pointer(&vec[0])))
	if errno != 0 {
		return nil, errno
	}
	return vec, nil
}

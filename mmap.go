/*
Complete support for memory-mapped files.
*/
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
	MADV_NORMAL     Advice = syscall.MADV_NORMAL
	MADV_RANDOM            = syscall.MADV_RANDOM
	MADV_SEQUENTIAL        = syscall.MADV_SEQUENTIAL
	MADV_WILLNEED          = syscall.MADV_WILLNEED
	MADV_DONTNEED          = syscall.MADV_DONTNEED
)

type MincoreState byte

const (
	MINCORE_INCORE MincoreState = 0x1 // the memory page is core resident at the time of the mincore() call
)

// An mmap region. Subslicing of this type should be done at memory page
// boundary (usually 4096 bytes) for most methods to work properly. The size
// of the mmap is limited to 2GB on Go verions prior to 1.1 or 32-bit
// platforms due to `int` being 32-bit.
type Mmap []byte

// Create an mmap backed by a file. Offset must be multiples of memory page size.
func Map(f *os.File, offset int64, len int, prot Prot, flags MapFlag) (Mmap, error) {
	return syscall.Mmap(int(f.Fd()), offset, len, int(prot), int(flags))
}

// Create an anonymous mmap without backing file.
func AnonMap(len int, prot Prot, flags MapFlag) (Mmap, error) {
	flags |= MAP_ANON // force anonymous
	return syscall.Mmap(-1, 0, len, int(prot), int(flags))
}

// Unmap the mmap. After unmap, program will crash if any slices based on the mmap are used.
func (m Mmap) Unmap() error {
	err := syscall.Munmap(m)
	m = nil
	return err
}

// Flush the changes in memory to the backing file. If MS_ASYNC flag is used,
// Sync() will return immediately; actual flushing will happen later.
func (m Mmap) Sync(flags SyncFlag) error {
	_, _, errno := syscall.Syscall(syscall.SYS_MSYNC, uintptr(unsafe.Pointer(&m[0])), uintptr(len(m)), uintptr(flags))
	if errno != 0 {
		return errno
	}
	return nil
}

// Return a slice of MincoreState describing the in-core status of memory pages
// in the mmap. You should bit OR with the various `MINCORE_*` flags to find
// out the in-core state.
func (m Mmap) Incore() ([]MincoreState, error) {
	pageSize := os.Getpagesize()
	vec := make([]MincoreState, (len(m)+pageSize-1)/pageSize)
	_, _, errno := syscall.Syscall(syscall.SYS_MINCORE, uintptr(unsafe.Pointer(&m[0])), uintptr(len(m)), uintptr(unsafe.Pointer(&vec[0])))
	if errno != 0 {
		return nil, errno
	}
	return vec, nil
}

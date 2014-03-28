package mmap

import (
	"syscall"
)

const (
	MAP_32BIT         = MapFlag(syscall.MAP_32BIT) // since Linux 2.4.20, 2.6
	MAP_ANONYMOUS     = MapFlag(syscall.MAP_ANONYMOUS)
	MAP_DENYWRITE     = MapFlag(syscall.MAP_DENYWRITE)
	MAP_NORESERVE     = MapFlag(syscall.MAP_NORESERVE)
	MAP_GROWSDOWN     = MapFlag(syscall.MAP_GROWSDOWN)
	MAP_EXECUTABLE    = MapFlag(syscall.MAP_EXECUTABLE)
	MAP_LOCKED        = MapFlag(syscall.MAP_LOCKED)   // since Linux 2.5.37
	MAP_NONBLOCK      = MapFlag(syscall.MAP_NONBLOCK) // since Linux 2.5.46
	MAP_POPULATE      = MapFlag(syscall.MAP_POPULATE) // since Linux 2.5.46
	MAP_STACK         = MapFlag(syscall.MAP_STACK)    // since Linux 2.6.27
	MAP_HUGETLB       = MapFlag(syscall.MAP_HUGETLB)  // since Linux 2.6.32
	MAP_UNINITIALIZED = MapFlag(0x4000000)            // since Linux 2.6.33
)

const (
	MADV_REMOVE       = Advice(syscall.MADV_REMOVE)      // since Linux 2.6.16
	MADV_DONTFORK     = Advice(syscall.MADV_DONTFORK)    // since Linux 2.6.16
	MADV_DOFORK       = Advice(syscall.MADV_DOFORK)      // since Linux 2.6.16
	MADV_HWPOISON     = Advice(syscall.MADV_HWPOISON)    // since Linux 2.6.32
	MADV_MERGEABLE    = Advice(syscall.MADV_MERGEABLE)   // since Linux 2.6.32
	MADV_UNMERGABLE   = Advice(syscall.MADV_UNMERGEABLE) // since Linux 2.6.32
	MADV_SOFT_OFFLINE = Advice(101)                      // since Linux 2.6.33
	MADV_HUGEPAGE     = Advice(syscall.MADV_HUGEPAGE)    // since Linux 2.6.38
	MADV_NOHUGEPAGE   = Advice(syscall.MADV_NOHUGEPAGE)  // since Linux 2.6.38
	MADV_DONTDUMP     = Advice(16)                       // since Linux 3.4
	MADV_DODUMP       = Advice(17)                       // since Linux 3.4
)

type LockallFlag int

const (
	MCL_CURRENT = LockallFlag(syscall.MCL_CURRENT)
	MCL_FUTURE  = LockallFlag(syscall.MCL_FUTURE)
)

// Give kernel hints about usage pattern.
func (m Mmap) Advise(advice Advice) error {
	return syscall.Madvise(m, int(advice))
}

// Prevent the mmap from being paged to the swap area.
func (m Mmap) Lock() error {
	return syscall.Mlock(m)
}

// Undo lock.
func (m Mmap) Unlock() error {
	return syscall.Munlock(m)
}

// Change the protect mode of the mmap.
func (m Mmap) Protect(prot int) error {
	return syscall.Mprotect(m, prot)
}

func Lockall(flags LockallFlag) error {
	return syscall.Mlockall(int(flags))
}

func Unlockall() error {
	return syscall.Munlockall()
}

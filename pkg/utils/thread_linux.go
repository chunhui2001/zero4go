package utils

import "syscall"

// ThreadID 返回当前 OS 线程 ID（Linux）
func ThreadID() uint64 {
	return uint64(syscall.Gettid())
}

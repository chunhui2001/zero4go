package utils

import (
	"runtime"
	"strconv"
)

// GoroutineID 返回当前 goroutine 的 ID
func GoroutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	stack := buf[:n]
	// stack: "goroutine 18 [running]:\n"
	stack = stack[len("goroutine "):]
	i := 0
	for ; i < len(stack); i++ {
		if stack[i] == ' ' {
			break
		}
	}
	id, _ := strconv.ParseInt(string(stack[:i]), 10, 64)
	return id
}

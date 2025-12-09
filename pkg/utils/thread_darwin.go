package utils

/*
#include <pthread.h>
#include <stdint.h>

uint64_t gettid_wrapper() {
    uint64_t tid;
    pthread_threadid_np(NULL, &tid);
    return tid;
}
*/
import "C"

// ThreadID 返回当前 OS 线程 ID（macOS）
func ThreadID() uint64 {
	return uint64(C.gettid_wrapper())
}

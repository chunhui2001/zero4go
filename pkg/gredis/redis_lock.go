package gredis

import (
	"context"
	"time"

	"github.com/bsm/redislock"

	. "github.com/chunhui2001/zero4go/pkg/logs"
)

//// 分布式任务示例
//func DistributedJob() {
//	lock, cancel, err := ObtainLockWithAutoRefresh("distributed-cron-job", 10*time.Second, 3*time.Second)
//	if err != nil {
//		fmt.Println("Another node is running the job or error:", err)
//		return
//	}
//
//	defer func() {
//		cancel() // 停止自动刷新
//		if err := lock.Release(ctx); err != nil {
//			fmt.Println("Failed to release lock:", err)
//		}
//	}()
//
//	// 执行任务
//	fmt.Println("Executing distributed job...")
//	time.Sleep(8 * time.Second)
//	fmt.Println("Job finished")
//}

// 获取锁并自动刷新 TTL
// ttl ≠ 锁最多活多久
// ttl = 单次租约时长
// 自动刷新 = 不断续租
// 超过 ttl 时间没有 Refresh
// Redis 会删除锁
// 其他节点可以重新抢锁
// 你的任务可能 被并发执行
// ttl = 你能接受的“任务失联最大时间”
// | 场景        | ttl | refreshInterval |
// | --------- | --- | --------------- |
// | 普通定时任务    | 10s | 3s              |
// | IO / 网络任务 | 30s | 10s             |
// | 重计算任务     | 60s | 15–20s          |
// ttl: 10*time.Second, // 👉 如果我 10 秒内没心跳，锁就释放
// 3*time.Second,  // 👉 每 3 秒续一次 10 秒的租

// ObtainLock 如果 “我拿到了这个锁，如果我连续 10 秒没刷新，就当我死了，
// 其他节点可以接管；
// 只要我活着，每 3 秒告诉 Redis 一次。”
// ttl 不是“任务预计执行时间”
// TTL = 心跳容忍时间
// 执行时间 = 无上限（靠 Refresh）
// ttl 是 锁的租约时间
// Refresh 会 不断重置 ttl
// ttl 决定：
// 崩溃后多久能被别人接管
// 系统是否会出现并发执行
// “如果持有锁的节点在 ttl 时间内消失，这把锁会被自动回收”
// 如何进一步“硬防” double execute（进阶）
// 任务内部再加一次 fencing token（推荐）
// 如果你是「每秒最多执行一次」，那么：
// ttl := 1 * time.Second
// 不要在任务结束时立即 Release
// 让 TTL 自然过期
// defer cancel()     // 停止刷新
// // 不调用 lock.Release()
// 一秒内，不可能再次获得锁
// ttl := 1 * time.Second
//
// lock, cancel, err := ObtainLockWithAutoRefresh(key, ttl)
// if err != nil {
// return
// }
//
// defer cancel()
// // 不调用 Release，让 TTL 控制执行频率

// LeaseLock 租约锁
// key: job key, 用一个 key 保证「这一秒只能执行一次」, 例如: job:2025-12-14T10:00:01
// ttl: 不是“任务预计执行时间”, 加入 ttl=10, 即: 如果我 10 秒内没心跳，锁就释放, 业务上保证 ttl >= 1s
// hook: 拿到锁后执行的函数
func LeaseLock(key string, ttl time.Duration, hook func()) {
	locker := redislock.New(RedisClient)

	// SET key value NX PX ttl
	lock, err := locker.Obtain(ctx, key, ttl, nil)

	if err != nil {
		Log.Errorf("Failed to refresh lock: Error=%s", err.Error())

		return // 获取不到锁，直接返回
	}

	// 创建取消函数，用于停止刷新
	refreshCtxTTl := 600 * time.Millisecond
	refreshCtx, cancel := context.WithCancel(ctx)

	// 启动 goroutine 自动刷新锁 TTL
	go func() {
		ticker := time.NewTicker(refreshCtxTTl)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := lock.Refresh(refreshCtx, ttl, nil); err != nil {
					// Refresh 失败时，必须停止业务（或者标记失效）
					Log.Errorf("Failed to refresh lock: Error=%s", err.Error())
				}
			case <-refreshCtx.Done():
				return
			}
		}
	}()

	defer func() {
		defer cancel() // 解除租约
		// lock.Release(ctx) // 不要在任务结束时立即 Release, 让 TTL 自然过期
	}()

	// 拿到锁执行业务逻辑
	hook()
}

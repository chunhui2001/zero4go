package gredis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	. "github.com/chunhui2001/zero4go/pkg/logs"
	"github.com/redis/go-redis/v9"
)

func Pipeline() redis.Pipeliner {
	return RedisClient.Pipeline()
}

func Exec(pipe redis.Pipeliner) {
	pipe.Exec(ctx)
}

func Expire(key string, expiration int) error {
	if err := RedisClient.Expire(ctx, key, time.Duration(expiration)*time.Second).Err(); err != nil {
		return err
	}

	return nil
}

func Del(key ...string) error {
	if err := RedisClient.Del(ctx, key...).Err(); err != nil {
		return err
	}

	return nil
}

func Ttl(key string) (int64, error) {
	val, err := RedisClient.TTL(ctx, key).Result()

	switch {
	case errors.Is(err, redis.Nil):

		return 0, nil
	case err != nil:
		Log.Errorf(`Redis-Get-Key-Error: Key=%s, Error=%s`, key, err.Error())

		return 0, err
	}

	return val.Nanoseconds(), nil
}

func Exists(key string) (bool, error) {

	val, err := RedisClient.Exists(ctx, key).Result()

	switch {
	case errors.Is(err, redis.Nil):

		return false, nil
	case err != nil:
		Log.Errorf(`Redis-Get-Key-Error: Key=%s, Error=%s`, key, err.Error())

		return false, err
	}

	return val != 0, nil
}

func Get(key string) (string, error) {
	val, err := RedisClient.Get(ctx, key).Result()

	switch {
	case errors.Is(err, redis.Nil):

		return "", nil
	case err != nil:
		Log.Errorf(`Redis-Get-Key-Error: Key=%s, Error=%s`, key, err.Error())

		return "", err
	case val == "":
		return "", nil
	}

	return val, nil
}

// expir 0 代表无过期时间, 过期时间单位是秒
func Set(key string, value string, expir int) bool {
	if err := RedisClient.Set(ctx, key, value, time.Duration(expir)*time.Second).Err(); err != nil {
		Log.Errorf(`Redis-Set-Error: Key=%s, Error=%s`, key, err.Error())

		return false
	}

	return true
}

func SetNX(key string, value string, expir int) (bool, error) {
	result, err := RedisClient.SetNX(ctx, key, value, time.Duration(expir)*time.Second).Result()

	if err != nil {
		Log.Errorf(`Redis-SetNX-Error: Key=%s, Error=%s`, key, err.Error())

		return false, err
	}

	return result, nil
}

// 将给定 key 的值设为 value ，并返回 key 的旧值(old value)。
// 当 key 存在但不是字符串类型时，返回一个错误。
// 当 key 没有旧值时，也即是，key 不存在时，返回 null 的同时将当前key设置为新值
func GetSet(key string, value string) (string, error) {
	val, err := RedisClient.GetSet(ctx, key, value).Result()

	switch {
	case errors.Is(err, redis.Nil):

		return "", nil
	case err != nil:
		Log.Errorf(`Redis-GetSet-Key-Error: Key=%s, Error=%s`, key, err.Error())

		return "", err
	case val == "":

		return "", nil
	}

	return val, nil
}

// 查询列表元素索引,没找到返回-1
// The command returns the index of matching elements inside a Redis list.
// maxLen: 最多找几个
func LindexOf(key string, value string, maxLen int64) (int64, error) {
	val, err := RedisClient.LPos(ctx, key, value, redis.LPosArgs{Rank: 0, MaxLen: maxLen}).Result()

	//redis.LPosArgs{}

	switch {
	case errors.Is(err, redis.Nil):

		return -1, nil
	case err != nil:
		Log.Errorf(`Redis-LindexOf-Error: Key=%s, Error=%s`, key, err.Error())

		return 0, err
	}

	return val, nil
}

// 列表操作
func Lpush(key string, values ...interface{}) bool {
	if err := RedisClient.LPush(ctx, key, values...).Err(); err != nil {
		Log.Errorf(`Redis-Lpush-Error: Key=%s, Error=%s`, key, err.Error())

		return false
	}

	return true
}

// 列表操作
func Rpush(key string, values ...interface{}) bool {
	if err := RedisClient.RPush(ctx, key, values...).Err(); err != nil {
		Log.Errorf(`Redis-Rpush-Error: Key=%s, Error=%s`, key, err.Error())

		return false
	}

	return true
}

// 读取列表元素: end=-1, 读取所有
func Lrange(key string, start int64, end int64) []string {
	val, err := RedisClient.LRange(ctx, key, start, end).Result()

	switch {
	case errors.Is(err, redis.Nil):

		return []string{}
	case err != nil:
		Log.Errorf(`Redis-Lrange-Error: Key=%s, Error=%s`, key, err.Error())

		return nil
	}

	return val
}

// 删除指定范围的列表元素
// start=100, end=-1, 将第100个之前的全部删除, 即保留100个之后的元素
func Ltrim(key string, start int64, end int64) bool {
	if err := RedisClient.LTrim(ctx, key, start, end).Err(); err != nil {
		Log.Errorf(`Redis-Lpop-Error: Key=%s, Error=%s`, key, err.Error())

		return false
	}

	return true
}

func Lpop(key string) (string, error) {
	val, err := RedisClient.LPop(ctx, key).Result()

	switch {
	case errors.Is(err, redis.Nil):

		return "", nil
	case err != nil:
		Log.Errorf(`Redis-Lpop-Error: Key=%s, Error=%s`, key, err.Error())

		return "", err
	}

	return val, nil
}

func Rpop(key string) (string, error) {
	val, err := RedisClient.RPop(ctx, key).Result()

	switch {
	case errors.Is(err, redis.Nil):

		return "", nil
	case err != nil:
		Log.Errorf(`Redis-Rpop-Error: Key=%s, Error=%s`, key, err.Error())

		return "", err
	}

	return val, nil
}

func Llen(key string) (int64, error) {
	val, err := RedisClient.LLen(ctx, key).Result()

	switch {
	case errors.Is(err, redis.Nil):

		return 0, nil
	case err != nil:
		Log.Errorf(`Redis-Llen-Error: Key=%s, Error=%s`, key, err.Error())

		return 0, err
	}

	return val, nil

}

func Hset(key string, values ...interface{}) bool {
	if err := RedisClient.HSet(ctx, key, values...).Err(); err != nil {
		Log.Errorf(`Redis-Hget-Error: Key=%s, Error=%s`, key, err.Error())

		return false
	}

	return true
}

func Hsetnx(key string, field string, value interface{}) bool {
	if err := RedisClient.HSetNX(ctx, key, field, value).Err(); err != nil {
		Log.Errorf(`Redis-Hget-Error: Key=%s, Error=%s`, key, err.Error())

		return false
	}

	return true
}

func Hget(key string, field string) (string, error) {
	val, err := RedisClient.HGet(ctx, key, field).Result()

	switch {
	case errors.Is(err, redis.Nil):
		return "", nil
	case err != nil:
		Log.Errorf(`Redis-Hget-Error: Key=%s, Error=%s`, key, err.Error())
		return "", err
	}

	return val, nil
}

func Hgetall(key string) (map[string]string, error) {
	val, err := RedisClient.HGetAll(ctx, key).Result()

	switch {
	case errors.Is(err, redis.Nil):
		return nil, nil
	case err != nil:
		Log.Errorf(`Redis-Hgetall-Error: Key=%s, Error=%s`, key, err.Error())

		return nil, err
	}

	return val, nil
}

func Hvals(key string) ([]string, error) {
	val, err := RedisClient.HVals(ctx, key).Result()

	switch {
	case errors.Is(err, redis.Nil):
		return nil, nil
	case err != nil:
		Log.Errorf(`Redis-Hvals-Error: Key=%s, Error=%s`, key, err.Error())
		return nil, err
	}

	return val, nil
}

func Zincr(key string) (int64, error) {
	result, err := RedisClient.Incr(ctx, key).Result()

	if err != nil {
		Log.Errorf(`Redis-Zincr-Error: Key=%s, Error=%s`, key, err.Error())
		return 0, err
	}

	return result, nil
}

//func Pub(channel string, payload string) {
//
//	err := RedisClient.Publish(ctx, channel, payload).Err()
//
//	if err != nil {
//		Log.Error(fmt.Sprintf("Redis-Publish-Error: channel=%s, Error=%v", channel, err))
//	}
//}

//func Sub(channel string, handler MessageHandler) {
//
//	if channel == "" {
//		return
//	}
//
//	if conf != nil && conf.Mode == Disabled {
//		panic(errors.New("Redis-Not-Enabled"))
//	}
//
//	if !connected {
//		logger.Info("Redis-Not-Connected: connected=" + utils.ToString(connected))
//		return
//	}
//
//	var pubSub *redis.PubSub
//
//	if redisClient != nil {
//		pubSub = redisClient.Subscribe(ctx, channel)
//	} else if redisCluster != nil {
//		pubSub = redisCluster.Subscribe(ctx, channel)
//	} else {
//		panic(errors.New("Redis-Client-Not-Initializable"))
//	}
//
//	// defer pubSub.Close()
//
//	logger.Info("Redis-Subscribe-A-Channel: channel=" + channel)
//
//	go LoopMessage(pubSub, channel, handler)
//
//}

//func LoopMessage(pubSub *redis.PubSub, channel string, handler MessageHandler) {
//
//	for {
//
//		msg, err := pubSub.ReceiveMessage(ctx)
//
//		if err != nil {
//			logger.Error(fmt.Sprintf("Redis-ReceiveMessage-Error: channel=%s, errorMessage=%s", channel, utils.ErrorToString(err)))
//		} else {
//			if handler == nil {
//				logger.Info("Redis-ReceivedMessage: channel=" + msg.Channel + ", payload=" + msg.Payload)
//			} else {
//				go handler(channel, msg.Payload)
//			}
//		}
//	}
//}

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

// ObtainLockWithAutoRefresh 获取锁并自动刷新 TTL
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

// 如果 “我拿到了这个锁，如果我连续 10 秒没刷新，就当我死了，
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
func ObtainLockWithAutoRefresh(key string, ttl time.Duration) (*redislock.Lock, context.CancelFunc, error) {
	locker := redislock.New(RedisClient)

	// SET key value NX PX ttl
	lock, err := locker.Obtain(ctx, key, ttl, nil)

	if err != nil {
		return nil, nil, err // 获取不到锁，直接返回
	}

	// 创建取消函数，用于停止刷新
	refreshCtxTTl := 300 * time.Millisecond
	refreshCtx, cancel := context.WithCancel(ctx)

	// 启动 goroutine 自动刷新锁 TTL
	go func() {
		ticker := time.NewTicker(refreshCtxTTl)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := lock.Refresh(refreshCtx, ttl, nil); err != nil {
					fmt.Println("Failed to refresh lock:", err)
				}
			case <-refreshCtx.Done():
				return
			}
		}
	}()

	return lock, cancel, nil
}

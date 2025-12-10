package gredis

import (
	"errors"
	"time"

	. "github.com/chunhui2001/zero4go/pkg/logs"
	"github.com/go-redis/redis/v8"
)

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

package gredis

import (
	"context"
	"strconv"
	"time"

	. "github.com/chunhui2001/zero4go/pkg/logs"
)

type String interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type RedisBlockingQueue[T String] struct {
	Key string
	Len int64
}

func (q *RedisBlockingQueue[T]) Push(values ...T) {
	if len(values) == 0 {

		return
	}

	// Lua 原子安全 push，固定容量
	var _lua = `
local maxLen = tonumber(ARGV[1])
local curLen = redis.call('LLEN', KEYS[1])
local pushCount = #ARGV - 1

if curLen + pushCount > maxLen then

    return 0
end

local values = {}

for i = 2, #ARGV do
    values[#values + 1] = ARGV[i]
end

redis.call('RPUSH', KEYS[1], unpack(values))

return pushCount
`

	args := make([]interface{}, 0, len(values)+1)
	args = append(args, q.Len)

	for _, v := range values {
		args = append(args, toString(v))
	}

	for {
		count, err := RedisClient.Eval(context.Background(), _lua, []string{q.Key}, args...).Int64()

		if err != nil {
			Log.Errorf("Push Failed: Key=%s, Error=%+v", q.Key, err)
			time.Sleep(time.Millisecond * 100)

			continue
		}

		var realCount = RedisClient.LLen(context.Background(), q.Key)

		Log.Infof("Redis阻塞队列: Key=%s, RealSize=%d, PushSize=%d, 队列支持的最大Size=%d", q.Key, realCount, len(values), q.Len)

		if count > 0 {
			return
		}

		Log.Infof("Redis阻塞队列满，等待: Key=%s, RealSize=%d, PushSize=%d, 队列支持的最大Size=%d", q.Key, realCount, len(values), q.Len)

		// 队列满，等待
		time.Sleep(time.Millisecond * 1000)
	}
}

func (q *RedisBlockingQueue[T]) Pop(batchSize int) []T {
	// Lua 原子安全, 批量取
	var _lua = `
local vals = redis.call('LRANGE', KEYS[1], 0, ARGV[1]-1)

if #vals > 0 then
    redis.call('LTRIM', KEYS[1], #vals, -1)
end

return vals
`

	for {
		res, err := RedisClient.Eval(context.Background(), _lua, []string{q.Key}, batchSize).Result()

		if err != nil {
			Log.Errorf("Pop Failed: Key=%s, Error=%+v", q.Key, err)

			return make([]T, 0)
		}

		vals := res.([]interface{})

		if len(vals) == 0 {

			return make([]T, 0)
		}

		out := make([]T, len(vals))

		for i, v := range vals {
			out[i] = fromString[T](v.(string))
		}

		return out
	}
}

func toString[T String](v T) string {
	switch x := any(v).(type) {
	case string:
		return x
	case int:
		return strconv.Itoa(x)
	case int8:
		return strconv.FormatInt(int64(x), 10)
	case int16:
		return strconv.FormatInt(int64(x), 10)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case uint:
		return strconv.FormatUint(uint64(x), 10)
	case uint8:
		return strconv.FormatUint(uint64(x), 10)
	case uint16:
		return strconv.FormatUint(uint64(x), 10)
	case uint32:
		return strconv.FormatUint(uint64(x), 10)
	case uint64:
		return strconv.FormatUint(x, 10)
	default:
		panic("unsupported type")
	}
}

func fromString[T String](s string) T {
	var zero T

	switch any(zero).(type) {
	case string:
		return any(s).(T)

	case int:
		v, _ := strconv.Atoi(s)
		return any(v).(T)

	case int8:
		v, _ := strconv.ParseInt(s, 10, 8)
		return any(int8(v)).(T)

	case int16:
		v, _ := strconv.ParseInt(s, 10, 16)
		return any(int16(v)).(T)

	case int32:
		v, _ := strconv.ParseInt(s, 10, 32)
		return any(int32(v)).(T)

	case int64:
		v, _ := strconv.ParseInt(s, 10, 64)
		return any(v).(T)

	case uint:
		v, _ := strconv.ParseUint(s, 10, 0)
		return any(uint(v)).(T)

	case uint8:
		v, _ := strconv.ParseUint(s, 10, 8)
		return any(uint8(v)).(T)

	case uint16:
		v, _ := strconv.ParseUint(s, 10, 16)
		return any(uint16(v)).(T)

	case uint32:
		v, _ := strconv.ParseUint(s, 10, 32)
		return any(uint32(v)).(T)

	case uint64:
		v, _ := strconv.ParseUint(s, 10, 64)
		return any(v).(T)

	default:
		panic("unsupported type")
	}
}

package gredis

import (
	"context"
	"strconv"

	. "github.com/chunhui2001/zero4go/pkg/logs" //nolint:staticcheck
)

type Number interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type RedisBlockingQueue[T Number] struct {
	Key      string
	MaxCount int64
	MaxVal   int64 // 当前队列里的最大值
	LastVal  int64 // 当前队列里的最后一个值
}

func (q *RedisBlockingQueue[T]) Push(values ...T) (int64, error) {
	if len(values) == 0 {

		return 0, nil
	}

	// Lua 原子安全 push，固定容量
	var _lua = `
local maxLen 	= tonumber(ARGV[1])
local curLen 	= redis.call('LLEN', KEYS[1])
local maxKey 	= KEYS[1] .. "_max"
local lastKey 	= KEYS[1] .. "_last"

local pushCount = #ARGV - 1

local curMax 	= redis.call('GET', maxKey)
curMax 			= curMax and tonumber(curMax) or 0

if curLen + pushCount > maxLen then
	local last 		= redis.call('LINDEX', KEYS[1], -1)
	last 			= last and tonumber(last) or 0

    return {0, curMax, last}
end

local values = {}

for i = 2, #ARGV do
    local v = tonumber(ARGV[i])

    if v then
        values[#values + 1] = ARGV[i]

        if not curMax or v > curMax then
            curMax = v
        end
    end
end

if #values > 0 then
    redis.call('RPUSH', KEYS[1], unpack(values))
    redis.call('SET', maxKey, curMax)

    local ttl = redis.call('TTL', KEYS[1])

    if ttl > 0 then
        redis.call('EXPIRE', maxKey, ttl)
    end
end

local last 		= redis.call('LINDEX', KEYS[1], -1)
last 			= last and tonumber(last) or 0

return {#values, curMax, last}
`

	args := make([]interface{}, 0, len(values)+1)
	args = append(args, q.MaxCount)

	for _, v := range values {
		args = append(args, toString(v))
	}

	res, err := RedisClient.Eval(context.Background(), _lua, []string{q.Key}, args...).Result()

	if err != nil {

		return 0, err
	}

	var realCount, _ = RedisClient.LLen(context.Background(), q.Key).Result()

	Log.Infof("%+v", res)

	vals := res.([]interface{})
	pushCount := vals[0].(int64)
	q.MaxVal = vals[1].(int64)
	q.LastVal = vals[2].(int64)

	if pushCount > 0 {
		return pushCount, nil
	}

	Log.Warnf("Redis阻塞队列满: Key=%s, 队列支持的最大Size=%d, RealSize=%d, PushSize=%d, MaxVal=%d", q.Key, q.MaxCount, realCount, len(values), q.MaxVal)

	return 0, nil
}

func (q *RedisBlockingQueue[T]) Pop(count int) []T {
	// Lua 原子安全, 批量取
	var _lua = `
local vals = redis.call('LRANGE', KEYS[1], 0, ARGV[1]-1)

if #vals > 0 then
    redis.call('LTRIM', KEYS[1], #vals, -1)
end

return vals
`

	res, err := RedisClient.Eval(context.Background(), _lua, []string{q.Key}, count).Result()

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

func toString[T Number](v T) string {
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

func fromString[T Number](s string) T {
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

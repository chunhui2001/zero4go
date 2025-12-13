package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const TimeStampFormat = "2006-01-02T15:04:05.000Z07:00"

func NowTimestamp() int64 {
	return time.Now().UTC().UnixNano() / int64(time.Millisecond)
}

func RootDir() string {
	var appRoot = os.Getenv("WORK_DIR")

	if appRoot != "" {
		return appRoot
	}

	dir, _ := os.Getwd()

	return dir
}

func TempDir() string {
	return os.TempDir()
}

func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

func ToString(s any) string {

	switch s.(type) {
	case float64, float32:
		return fmt.Sprint(s)
	case string:
		return fmt.Sprintf("%s", s)
	case bool:
		return fmt.Sprintf("%t", s)
	case byte:
		return fmt.Sprintf("%x", s)
	case []uint8:
		return string(s.([]byte))
	default:
		return fmt.Sprintf("%d", s)
	}
}

func ToJsonString(v any) string {
	if v == nil {
		return ""
	}

	if reflect.TypeOf(v).String() == "string" {
		return v.(string)
	}

	b, err := json.Marshal(v)

	if err != nil {
		panic(err)
	}

	return string(b)
}

func NormalizeSpace(s string) string {
	// 1. 去掉所有换行符
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ") // Windows 换行

	// 2. 将一个或多个空格替换为单个空格
	re := regexp.MustCompile(`\s+`)
	s = re.ReplaceAllString(s, " ")

	// 3. 去掉首尾空格
	return strings.TrimSpace(s)
}

func ToDateTimeUTCString(tm time.Time) string {
	return tm.Format(TimeStampFormat)
}

func OfMap(kv ...any) map[string]any {
	if kv == nil {
		return make(map[string]any)
	}

	if len(kv)%2 != 0 {
		panic(errors.New("Invalid map size: currentSize=" + ToString(len(kv))))
	}

	m := make(map[string]any)

	for i := 0; i < len(kv); i++ {
		k := ToString(kv[i])
		m[k] = kv[i+1]

		i++
	}

	return m
}

func ToBase64String(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func FromBase64String(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

package utils

import (
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

func RootDir() string {

	var appRoot = os.Getenv("APP_ROOT")

	if appRoot != "" {
		return appRoot
	}

	dir, _ := os.Getwd()

	return dir
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

func ToJsonString(v interface{}) string {
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

func NowTimestamp() int64 {
	return time.Now().UTC().UnixNano() / int64(time.Millisecond)
}

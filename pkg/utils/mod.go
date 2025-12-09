package utils

import (
	"errors"
	"fmt"
	"os"
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

func ToDateTimeUTCString(tm time.Time) string {
	return tm.Format(TimeStampFormat)
}

func Max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

package logs

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/chunhui2001/zero4go/pkg/stdout"
)

var SkipFieldsForConsole = []string{"@timestamp", "@version", "app", "build_git_hash", "build_git_version", "build_name", "build_timestamp", "captain_gen", "captain_seq", "env", "hostname", "logger_name", "thread_name"}

func NewLogger(conf *LogConf) Logger {

	fileEnable := strings.Contains(conf.LogOutput, "file")
	consoleEnable := strings.Contains(conf.LogOutput, "console")
	kafkaEnable := strings.Contains(conf.LogOutput, "kafka")

	var writers []io.Writer

	if consoleEnable {
		consoleWriter := zerolog.NewConsoleWriter()
		consoleWriter.NoColor = !conf.LogColor
		consoleWriter.FieldsExclude = SkipFieldsForConsole

		consoleWriter.FormatTimestamp = func(i interface{}) string {
			return time.Now().Format(stdout.TimeStampFormat) // 自动包含时区
		}
		consoleWriter.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("%s", i))
		}

		writers = append(writers, consoleWriter)
	}

	if fileEnable {
		rotate := &lumberjack.Logger{
			Compress:   true, // gzip 压缩旧文件
			Filename:   conf.LogFilePath,
			MaxSize:    Max(conf.LogFileMaxSize, 10),    // MB，超过后切分
			MaxBackups: Max(conf.LogFileMaxBackups, 10), // 最多保留 10 个旧文件
			MaxAge:     Max(conf.LogFileMaxAge, 30),     // 保留 30 天
		}

		writers = append(writers, rotate)
	}

	if kafkaEnable {
		kWriter := NewKafkaWriter(strings.Split(conf.LogKafkaServer, ","), conf.LogKafkaTopic)

		writers = append(writers, kWriter)
	}

	mw := zerolog.MultiLevelWriter(
		writers...,
	)

	l := zerolog.New(mw).
		Level(zerolog.DebugLevel).
		With().
		// 修复封装log后 “导致 caller 跳帧” 问题
		//Caller().
		// CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount+1).
		Str("app", os.Getenv("APP_NAME")).
		Str("build_git_hash", os.Getenv("BUILD_GIT_VERSION")).
		Str("captain_gen", os.Getenv("CAPTAIN_GEN")).
		Str("env", os.Getenv("ENV")).
		Logger()

	logger := Logger{
		Logger: l,
	}

	logger.Infof("logger initialized: %v", conf)

	return logger
}

func Max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

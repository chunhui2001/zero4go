package logs

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/chunhui2001/zero4go/pkg/stdout"
	"github.com/chunhui2001/zero4go/pkg/utils"
)

type LogConf struct {
	LogLevel          string `mapstructure:"LOG_LEVEL" json:"log_level"`                       // debug, into, error, warn
	LogOutput         string `mapstructure:"LOG_OUTPUT" json:"log_output"`                     // console,file,kafka
	LogColor          bool   `mapstructure:"LOG_COLOR" json:"log_color"`                       // true
	LogFilePath       string `mapstructure:"LOG_FILE_PATH" json:"log_file_path"`               // /tmp/logs
	LogFileMaxSize    int    `mapstructure:"LOG_FILE_MAX_SIZE" json:"log_file_max_size"`       // 5 (MB)
	LogFileMaxBackups int    `mapstructure:"LOG_FILE_MAX_BACKUPS" json:"log_file_max_backups"` // 10
	LogFileMaxAge     int    `mapstructure:"LOG_FILE_MAX_AGE" json:"log_file_max_age"`         // 30
	LogKafkaServer    string `mapstructure:"LOG_KAFKA_SERVER" json:"log_kafka_server"`         // 127.0.0.1:9092
	LogKafkaTopic     string `mapstructure:"LOG_KAFKA_TOPIC" json:"log_kafka_topic"`           // app_log
}

var LogSetting = &LogConf{
	LogLevel:          "debug",
	LogOutput:         "console,file",
	LogColor:          false,
	LogFilePath:       "logs/app.log",
	LogFileMaxSize:    20, // 单位: MB，超过后切分
	LogFileMaxBackups: 10, // 最多保留 10 个旧文件
	LogFileMaxAge:     30, // 保留 30 天
}

var Log Logger

func InitLog() {
	Log = NewLogger(LogSetting)
}

func OnChange(conf *LogConf) {
	Log = NewLogger(conf)
}

type Logger struct {
	zerolog.Logger
}

func (l Logger) WithThreadName0() Logger {
	return Logger{
		Logger: l.Logger.With().
			Str("thread_name", fmt.Sprintf("thread-g%d-t%d", utils.GoroutineID(), utils.ThreadID())).
			Str("@timestamp", time.Now().Format(stdout.TimeStampFormat)).
			Caller().
			Logger(),
	}
}

func (l Logger) WithThreadName1() Logger {
	return Logger{
		Logger: l.Logger.With().
			CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount+1).
			Str("thread_name", fmt.Sprintf("thread-g%d-t%d", utils.GoroutineID(), utils.ThreadID())).
			Str("@timestamp", time.Now().Format(stdout.TimeStampFormat)).
			Logger(),
	}
}

func (l Logger) WithThreadName3() Logger {
	return Logger{
		Logger: l.Logger.With().
			CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount+3).
			Str("thread_name", fmt.Sprintf("thread-g%d-t%d", utils.GoroutineID(), utils.ThreadID())).
			Str("@timestamp", time.Now().Format(stdout.TimeStampFormat)).
			Logger(),
	}
}

func (l Logger) infol() *zerolog.Event {
	return l.Logger.Info()
}

func (l Logger) Infoe() *zerolog.Event {
	return l.WithThreadName0().infol()
}

func (l Logger) Infoe1() *zerolog.Event {
	return l.WithThreadName1().infol()
}

func (l Logger) Infoe3() *zerolog.Event {
	return l.WithThreadName3().infol()
}

func (l Logger) Infof(format string, args ...interface{}) {
	l.Infoe1().Msgf(format, args...)
}

func (l Logger) Info(msg string) {
	l.Infoe1().Msg(msg)
}

func (l Logger) warnl() *zerolog.Event {
	return l.Logger.Warn()
}

func (l Logger) Warne() *zerolog.Event {
	return l.WithThreadName0().warnl()
}

func (l Logger) Warne1() *zerolog.Event {
	return l.WithThreadName1().warnl()
}

func (l Logger) Warne3() *zerolog.Event {
	return l.WithThreadName3().warnl()
}

func (l Logger) Warnf(format string, args ...interface{}) {
	l.Warne1().Msgf(format, args...)
}

func (l Logger) Warn(msg string) {
	l.Warne1().Msg(msg)
}

func (l Logger) errorl() *zerolog.Event {
	return l.Logger.Error()
}

func (l Logger) Errore() *zerolog.Event {
	return l.WithThreadName0().errorl()
}

func (l Logger) Errore1() *zerolog.Event {
	return l.WithThreadName1().errorl()
}

func (l Logger) Errore3() *zerolog.Event {
	return l.WithThreadName3().errorl()
}

func (l Logger) Errorf(format string, args ...interface{}) {
	l.Errore1().Msgf(format, args...)
}

func (l Logger) Error(msg string) {
	l.Errore1().Msg(msg)
}

func (l Logger) debugl() *zerolog.Event {
	return l.Logger.Debug()
}

func (l Logger) Debuge() *zerolog.Event {
	return l.WithThreadName0().debugl()
}

func (l Logger) Debuge1() *zerolog.Event {
	return l.WithThreadName1().debugl()
}

func (l Logger) Debuge3() *zerolog.Event {
	return l.WithThreadName3().debugl()
}

func (l Logger) Debugf(format string, args ...interface{}) {
	l.Debuge1().Msgf(format, args...)
}

func (l Logger) Debug(msg string) {
	l.Debuge1().Msg(msg)
}

func (l Logger) fatall() *zerolog.Event {
	return l.Logger.Fatal()
}

func (l Logger) Fatale() *zerolog.Event {
	return l.WithThreadName0().fatall()
}

func (l Logger) Fatale1() *zerolog.Event {
	return l.WithThreadName1().fatall()
}

func (l Logger) Fatale3() *zerolog.Event {
	return l.WithThreadName3().fatall()
}

func (l Logger) Fatalf(format string, args ...interface{}) {
	l.Fatale1().Msgf(format, args...)
}

func (l Logger) Fatal(format string) {
	l.Fatale1().Msg(format)
}

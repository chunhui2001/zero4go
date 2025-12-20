package gtask

import (
	"os"
	"time"

	"github.com/robfig/cron"

	"github.com/chunhui2001/zero4go/pkg/gredis"
	. "github.com/chunhui2001/zero4go/pkg/logs" //nolint:staticcheck
)

var c = cron.New()

func init() {
	c.Start()
}

// AddTask 添加一个定时任务
// "* * * * * *" -- 每秒1次
// "0/5 * * * * *" -- 每5秒
// "0 30 * * * *" -- 每半小时1次
// "15 * * * * *" -- 每15秒1次
// "@hourly" -- Every hour
// "@every 1h30m" -- Every hour thirty
// "@daily" -- Every day
// ###################################################################################
// -----                  | -----------                                | -------------
// Entry                  | Description                                | Equivalent To
// -----                  | -----------                                | -------------
// @yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 0 1 1 *
// @monthly               | Run once a month, midnight, first of month | 0 0 0 1 * *
// @weekly                | Run once a week, midnight between Sat/Sun  | 0 0 0 * * 0
// @daily (or @midnight)  | Run once a day, midnight                   | 0 0 0 * * *
// @hourly                | Run once an hour, beginning of hour        | 0 0 * * * *
// ###################################################################################
func AddTask(name string, JobID string, spec string, tasks func(key string)) {
	_ = c.AddFunc(spec, func() {
		gredis.Lock(JobID, 1*time.Second, 330*time.Millisecond, func() {
			var _key = JobID + "#" + time.Now().UTC().Format("2006-01-02T15:04:05")

			Log.Debugf("执行定时任务: Name=%s, Expr=%s, Key=%s, PID=%d", name, spec, _key, os.Getppid())

			tasks(_key)
		})
	})

	Log.Infof(`注册了一个定时任务: Name=%s, JobID=%s, Expr='%s'`, name, JobID, spec)
}

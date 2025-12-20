package build_info

import (
	"fmt"
	"log"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

// go version devel +475d92ba4d Thu Oct 5 10:50:18 2017 +0000 darwin/amd64

type BuildInfo struct {
	Name   string
	Author string
	Commit string
	OS     string
	Time   string
}

var INFO = &BuildInfo{
	OS: readBuildInfo("GOOS") + "/" + readBuildInfo("GOARCH"),
}

func (b *BuildInfo) Info() string {
	var s = fmt.Sprintf("%s version devel +%s %s [%s]", b.Name, b.Commit, b.OS, b.Time)
	log.Println(s)
	
	return s
}

var (
	layout = "Mon Jan 2 15:04:05 -0700 MST 2006"
	Commit string
)

func init() {
	ldflags := readBuildInfo("-ldflags")
	flags := strings.Split(ldflags, "-X ")

	for _, f := range flags {

		pair := strings.Split(f, "=")

		if len(pair) != 2 {
			continue
		}

		if strings.Trim(pair[0], " ") == "main.Author" {
			INFO.Author = strings.Trim(pair[1], " ")
		} else if strings.Trim(pair[0], " ") == "main.Commit" {
			if pair[1] == "" {
				INFO.Commit = "Unknown"
			} else {
				if len(pair[1]) >= 10 {
					INFO.Commit = strings.Trim(pair[1], " ")[0:10]
				} else {
					INFO.Commit = pair[1]
				}
			}

			Commit = INFO.Commit
		} else if strings.Trim(pair[0], " ") == "main.Name" {
			INFO.Name = strings.Trim(pair[1], " ")
		} else if strings.Trim(pair[0], " ") == "main.Time" {
			i, err := strconv.ParseInt(pair[1], 10, 64)
			if err == nil {
				tm := time.Unix(i, 0)
				loc, _ := time.LoadLocation("UTC")
				INFO.Time = tm.In(loc).Format(layout)
			}
		}
	}
}

func readBuildInfo(key string) string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == key {
				return setting.Value
			}
		}
	}

	return ""
}

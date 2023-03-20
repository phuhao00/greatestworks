package consul

import (
	"encoding/json"
	"net"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"
)

var messageSendTime int64

var memoryUsed uint64

var lastReadMemStatsTime time.Time

const (
	SignKey                  = "$m%s@y61*Qd@le"
	ReadReadMemStatsInterval = 2 * time.Minute
	DayBeginHour             = 5
	OneDaySec                = 24 * 60 * 60
	BeginHourSec             = DayBeginHour * 60 * 60
)

type Performance struct {
	Zid          int    `json:"zid"`
	Pid          int    `json:"pid"`
	PIdx         int    `json:"PIdx"`
	SvrPort      int    `json:"svrPort"`
	PlayerNum    int32  `json:"playerNum"`
	MaxPlayerNum int32  `json:"maxPlayerNum"`
	Mem          uint64 `json:"mem"`
	StartTM      int64  `json:"startTM"`
	SvrAddr      string `json:"svrAddr"`
	SvrId        string `json:"svrId"`
	SvrPath      string `json:"svrPath"`
	InnerAddr    string `json:"innerAddr"`
}

func GetPerformance(performance *Performance) string {
	host, _, _ := net.SplitHostPort(performance.SvrAddr)
	pid := os.Getpid()
	if time.Since(lastReadMemStatsTime) > ReadReadMemStatsInterval {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		memoryUsed = mem.HeapAlloc
		lastReadMemStatsTime = time.Now()
	}
	performance.Pid = pid
	performance.Mem = memoryUsed
	performance.SvrAddr = host
	data, err := json.Marshal(performance)
	if err != nil {
		return ""
	}
	return string(data)
}

func GetOsUserName() string {
	var name string
	u, err := user.Current()
	if err != nil {
		name = "unknown"
	} else {
		name = u.Username
	}
	sl := strings.Split(name, "\\")
	name = sl[len(sl)-1]
	return name
}

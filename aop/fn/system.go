package fn

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"greatestworks/aop/logger"
	"net"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	SignKey      = "hh$m%s@y61*Qd@le"
	ReadMemTime  = 2 * time.Minute
	DayStartHour = 5
	OneDaySec    = 24 * 60 * 60
	StartHourSec = DayStartHour * 60 * 60
)

var (
	lastReadMemTime time.Time
	memoryUsed      uint64
)

func GetUser() string {
	var userName string
	u, err := user.Current()
	if err != nil {
		userName = "unknow"
	} else {
		userName = u.Username
	}
	sl := strings.Split(userName, "\\")
	userName = sl[len(sl)-1]
	return userName
}

func GetMD5Sign(userID uint64, token string) string {
	appSecret := "f1382332e76bc73a34e4d635a20cb952"
	md5Str := fmt.Sprintf("userID=%vtoken=%v%v", userID, token, appSecret)
	md5Sign := md5.Sum([]byte(md5Str))
	return hex.EncodeToString(md5Sign[:])
}

func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func CheckSign(data, sign string) bool {
	token := md5.Sum([]byte(data + SignKey))
	return hex.EncodeToString(token[:]) == sign
}

func GetPrivateIPv4() (string, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range as {
		ipNet, ok := a.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue
		}

		ip := ipNet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip.String(), nil
		}
	}
	return "", errors.New("no private ip address")
}

func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Error(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func IpAddressStringToUint64(ipaddr string) uint64 {
	var digit uint64
	data := strings.Split(ipaddr, ":")
	port := SplitStringToUint32Slice(data[1], ".")
	val := SplitStringToUint32Slice(data[0], ".")
	digit = uint64(val[0]<<24) | uint64(val[1]<<16) | uint64(val[2]<<8) | uint64(val[3])
	digit = digit<<32 | uint64(port[0])
	return digit
}

type Performance struct {
	Zid          int // gateway, online用的分区 - 对全服聊天分区，登录展示分区
	Pid          int
	PIdx         int
	Mem          uint64
	StartTM      int64
	PlayerNum    int32
	MaxPlayerNum int32
	SvrID        string
	SvrAddr      string
	SvrPort      int
	SvrPath      string
	InnerAddr    string
}

func GetPerformanceExt(svrID string, svrAddr string, svrPort int, svrPath string, startTM int64, playerNum int32,
	processIdx int, zoneId int, innerIp string) string {
	host, _, _ := net.SplitHostPort(svrAddr)
	pid := os.Getpid()
	if time.Since(lastReadMemTime) > ReadMemTime {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		memoryUsed = mem.HeapAlloc
		lastReadMemTime = time.Now()
	}
	performance := &Performance{
		Zid:       zoneId,
		Pid:       pid,
		PIdx:      processIdx,
		Mem:       memoryUsed,
		StartTM:   startTM,
		PlayerNum: playerNum,
		SvrID:     svrID,
		SvrAddr:   host,
		SvrPort:   svrPort,
		SvrPath:   svrPath,
		InnerAddr: innerIp,
	}

	data, err := json.Marshal(performance)
	if err != nil {
		return ""
	}
	return string(data)
}

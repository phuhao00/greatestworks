package fn

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"os/user"
	"strings"
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

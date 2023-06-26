package redis

import (
	"fmt"
	"strconv"
)

func MakeGatewayKey(addr string) string {
	return fmt.Sprintf("GatewayConn:%v", addr)
}

func MakeTokenKey(userid uint64) string {
	return "token:" + strconv.FormatInt(int64(userid), 10)
}

package rediskey

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net"
	"time"
)

func CheckLive(redisCluster *redis.ClusterClient, userid uint64) bool {
	key := MakePlayerCacheKey(userid)
	val, _ := redisCluster.HGet(context.TODO(), key, "tick").Int64()
	return val+int64(2*2*KeepAliveFreq) > time.Now().Unix()
}

func FetchServiceID(svcName string, svcAddr string) string {
	host, port, err := net.SplitHostPort(svcAddr)
	if err != nil {
		return svcName + "-" + svcAddr
	}
	ipInt := ip2Int(host)
	return fmt.Sprintf("%v-%v", ipInt, port)
}

func ip2Int(ip string) uint32 {
	var ipInt uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &ipInt)
	return ipInt
}

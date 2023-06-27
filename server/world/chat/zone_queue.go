package chat

import (
	"context"
	"fmt"
	goRedis "github.com/go-redis/redis/v8"
	"greatestworks/aop/logger"
	"greatestworks/aop/redis"
	"greatestworks/server/world/server"
	"math"

	"time"
)

var (
	cfgZoneWindowSize int64 = 1000
	cfgZoneQueueSize  int64 = 2200
	cfgTimeSlice      int64 = 5
)

func getZoneQueueKeyName() string {
	return fmt.Sprintf("ZoneChatQueue-%v", server.Oasis.Config.Global.ZoneId)
}

func getZoneWindowKeyName(rankTime int64) string {
	return fmt.Sprintf("ZoneChatWindow-%v-%v", server.Oasis.Config.Global.ZoneId, rankTime)
}

func getNTimeSlice() int64 {
	return time.Now().Unix() / cfgTimeSlice
}

func getZoneQueueLength() int64 {
	ret := redis.RateLimitRedis().ZCard(context.TODO(), getZoneQueueKeyName())
	if ret.Err() != nil && ret.Err() != goRedis.Nil {
		return math.MaxInt64
	}
	return ret.Val()
}

func getZoneQueueRank(uidStr string) (int64, bool) {
	ZoneQueueKeyName := getZoneQueueKeyName()
	ret1 := redis.RateLimitRedis().ZRank(context.TODO(), ZoneQueueKeyName, uidStr)
	if ret1.Err() != nil {
		return 0, false
	}
	return ret1.Val() + 1, true
}

func genZoneQueueScore() int64 {
	return time.Now().Unix()
}

func ZoneQueueEnqueue(uidStr string) int64 {
	var curRank int64
	rankScore := genZoneQueueScore()
	member := &goRedis.Z{
		Score:  float64(rankScore),
		Member: uidStr,
	}
	ZoneQueueKeyName := getZoneQueueKeyName()
	redis.RateLimitRedis().ZAdd(context.TODO(), ZoneQueueKeyName, member)

	ret := redis.RateLimitRedis().ZRank(context.TODO(), ZoneQueueKeyName, uidStr)
	if ret.Err() != nil {
		logger.Error("err:%v", ret.Err().Error)
		return math.MaxInt64
	}

	curRank = ret.Val() + 1

	return curRank
}

func ZoneQueueDequeue(uidStr string) {
	ZoneQueueKeyName := getZoneQueueKeyName()
	redis.RateLimitRedis().ZRem(context.TODO(), ZoneQueueKeyName, uidStr)
}

func getZoneWindowSize(rankTime int64) int64 {
	windowKey := getZoneWindowKeyName(rankTime)
	windowSize, err := redis.CacheRedis().Get(context.TODO(), windowKey).Int64()

	if err != nil && err != goRedis.Nil {
		return 0
	}

	windowSize = cfgZoneWindowSize - windowSize
	if windowSize < 0 {
		windowSize = 0
	}

	return windowSize
}

func occupyingAZoneWindowSeat(rankTime int64) bool {
	windowKey := getZoneWindowKeyName(rankTime)
	curValue := redis.CacheRedis().Incr(context.TODO(), windowKey)
	if curValue.Err() != nil {
		return false
	}

	redis.RateLimitRedis().Expire(context.TODO(), windowKey, time.Duration(2*cfgTimeSlice)*time.Second)

	if curValue.Val() > cfgZoneWindowSize {
		return false
	} else {
		return true
	}
}

func checkSendZoneChat(uidStr string) int {

	curRank := ZoneQueueEnqueue(uidStr)

	if curRank > cfgZoneQueueSize {
		ZoneQueueDequeue(uidStr)
		return 2
	}

	rankTime := getNTimeSlice()
	ZoneWindowSize := getZoneWindowSize(rankTime)

	if curRank <= ZoneWindowSize {
		if occupyingAZoneWindowSeat(rankTime) {
			ZoneQueueDequeue(uidStr)
			return 0
		}
	}

	return 1
}

func retryCheckSendZoneChat(uidStr string) int {

	curRank, ok := getZoneQueueRank(uidStr)
	if !ok {
		return 2 // 当前已经不在队列中了不需要再重试了
	}

	rankTime := getNTimeSlice()

	ZoneWindowSize := getZoneWindowSize(rankTime)
	if curRank <= ZoneWindowSize {
		if occupyingAZoneWindowSeat(rankTime) {
			ZoneQueueDequeue(uidStr)
			return 0
		}
	}

	return 1
}

func clearZoneQueueTimeoutMemberTimer() {
	clearZoneQueueTimeoutMember(2 * 60)
}

func clearZoneQueueTimeoutMember(timeout int64) {

	var (
		minValue string
		maxValue string
	)

	if timeout == 0 {
		minValue = "0"
		curZoneQueueScore := genZoneQueueScore()
		maxValue = fmt.Sprintf("%v", curZoneQueueScore)
	} else {
		curZoneQueueScore := genZoneQueueScore()
		timeoutScore := curZoneQueueScore - timeout
		if timeoutScore < 0 {
			return
		}
		minValue = "0"
		maxValue = fmt.Sprintf("%v", timeoutScore)
	}

	ZoneQueueKeyName := getZoneQueueKeyName()
	redis.RateLimitRedis().ZRemRangeByScore(context.TODO(), ZoneQueueKeyName, minValue, maxValue)
}

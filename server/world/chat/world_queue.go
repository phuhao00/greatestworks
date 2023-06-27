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
	cfgWorldWindowSize int64 = 600
	cfgWorldQueueSize  int64 = 1200
)

func getWorldQueueKeyName() string {
	return "WorldChatQueue"
}

func getWorldWindowKeyName(rankTime int64) string {
	return fmt.Sprintf("WorldChatWindow-%v", rankTime)
}

func getWorldQueueLength() int64 {
	ret := redis.RateLimitRedis().ZCard(context.TODO(), getWorldQueueKeyName())
	if ret.Err() != nil && ret.Err() != goRedis.Nil {
		return math.MaxInt64
	}
	return ret.Val()
}

func getWorldQueueRank(uidStr string) (int64, bool) {
	WorldQueueKeyName := getWorldQueueKeyName()
	ret1 := redis.RateLimitRedis().ZRank(context.TODO(), WorldQueueKeyName, uidStr)
	if ret1.Err() != nil {
		return 0, false
	}
	return ret1.Val() + 1, true
}

func genWorldQueueScore() int64 {
	return time.Now().Unix()
}
func WorldQueueEnqueue(uidStr string) int64 {
	var curRank int64
	rankScore := genWorldQueueScore()
	member := &goRedis.Z{
		Score:  float64(rankScore),
		Member: uidStr,
	}
	WorldQueueKeyName := getWorldQueueKeyName()
	redis.RateLimitRedis().ZAdd(context.TODO(), WorldQueueKeyName, member)

	ret := redis.RateLimitRedis().ZRank(context.TODO(), WorldQueueKeyName, uidStr)
	if ret.Err() != nil {
		logger.Error("err:%v", ret.Err().Error)
		return math.MaxInt64
	}

	curRank = ret.Val() + 1

	return curRank
}

func WorldQueueDequeue(uidStr string) {
	WorldQueueKeyName := getWorldQueueKeyName()
	redis.RateLimitRedis().ZRem(context.TODO(), WorldQueueKeyName, uidStr)
}

func getWorldWindowSize(rankTime int64) int64 {
	windowKey := getWorldWindowKeyName(rankTime)
	windowSize, err := redis.CacheRedis().Get(context.TODO(), windowKey).Int64()
	if err != nil && err != goRedis.Nil {
		return 0
	}

	windowSize = cfgWorldWindowSize - windowSize
	if windowSize < 0 {
		windowSize = 0
	}

	return windowSize
}

func occupyingAWorldWindowSeat(rankTime int64) bool {
	windowKey := getWorldWindowKeyName(rankTime)
	curValue := redis.CacheRedis().Incr(context.TODO(), windowKey)
	if curValue.Err() != nil {
		return false
	}
	redis.RateLimitRedis().Expire(context.TODO(), windowKey, time.Duration(2*cfgTimeSlice)*time.Second)

	if curValue.Val() > cfgWorldWindowSize {
		return false
	} else {
		return true
	}
}

func checkSendWorldChat(uidStr string) int {
	curRank := WorldQueueEnqueue(uidStr)
	if curRank > cfgWorldQueueSize {
		WorldQueueDequeue(uidStr)
		return 2
	}

	rankTime := getNTimeSlice()
	WorldWindowSize := getWorldWindowSize(rankTime)

	if curRank <= WorldWindowSize {
		if occupyingAWorldWindowSeat(rankTime) {
			WorldQueueDequeue(uidStr)
			return 0
		}
	}

	return 1
}

func retryCheckSendWorldChat(uidStr string) int {
	curRank, ok := getWorldQueueRank(uidStr)
	if !ok {
		return 2
	}

	rankTime := getNTimeSlice()

	WorldWindowSize := getWorldWindowSize(rankTime)

	if curRank <= WorldWindowSize {
		if occupyingAWorldWindowSeat(rankTime) {
			WorldQueueDequeue(uidStr)
			return 0 // 重试成功了,发送喇叭消息
		}
	}

	return 1
}

func clearWorldQueueTimeoutMemberTimer() {
	if server.Oasis.Pid == 1 {
		clearWorldQueueTimeoutMember(2 * 60)
	}
}

func clearWorldQueueTimeoutMember(timeout int64) {

	if server.Oasis.Pid == 1 {
		var (
			minValue string
			maxValue string
		)

		if timeout == 0 {
			minValue = "0"
			curWorldQueueScore := genWorldQueueScore()
			maxValue = fmt.Sprintf("%v", curWorldQueueScore)
		} else {
			curWorldQueueScore := genWorldQueueScore()
			timeoutScore := curWorldQueueScore - timeout
			if timeoutScore < 0 {
				return
			}
			minValue = "0"
			maxValue = fmt.Sprintf("%v", timeoutScore)
		}

		WorldQueueKeyName := getWorldQueueKeyName()
		redis.RateLimitRedis().ZRemRangeByScore(context.TODO(), WorldQueueKeyName, minValue, maxValue)
	}
}

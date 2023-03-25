package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "greatestworks/aop/redis"
	"greatestworks/server/login/config"
	"math"
)

import "time"

var redisCluster = redis.NewClusterClient(nil)

func GetQueueName() string {
	return "login-queue"
}

func GetWindowKey(timestamp int64) string {
	return fmt.Sprintf("LoginWindow-%v", timestamp)
}

func GetQueueLen() int64 {
	return redisCluster.ZCard(context.TODO(), GetQueueName()).Val()
}

func Dequeue(account string) {
	redisCluster.ZRem(context.TODO(), GetQueueName(), account)
}
func Enqueue(account string) (int64, int64, int64) {
	var (
		curRank     int64
		curScore    int64
		curWaitTime int64
	)
	score := genScore(0, 0, int64(config.MaxRetryWaitTime))
	mem := &redis.Z{
		Score:  float64(score),
		Member: account,
	}
	queueName := GetQueueName()
	//todo
	redis.NewClusterClient(nil).ZAdd(context.TODO(), queueName, mem)
	return curRank, curScore, curWaitTime
}

func EnqueueWithWaitTime(dev string, curScore, lastWaitTime, waitTime int64) {
	rankScore := genScore(curScore, lastWaitTime, waitTime)
	mem := &redis.Z{
		Score:  float64(rankScore),
		Member: dev,
	}
	queueName := GetQueueName()
	redis.NewClusterClient(nil).ZAdd(context.TODO(), queueName, mem)
}

func GetQueueRank(account string) (rank, waitTime, score int64, err error) {
	queueName := GetQueueName()
	idxRet := redisCluster.ZRank(context.TODO(), queueName, account)
	if idxRet.Err() == redis.Nil {
		return rank, score, waitTime, idxRet.Err()
	} else if idxRet.Err() != nil {
		return math.MaxInt64, 0, 0, idxRet.Err()
	} else {
		rank = idxRet.Val() + 1
	}

	scoreRet := redisCluster.ZScore(context.TODO(), queueName, account)
	if scoreRet.Err() == redis.Nil {
		return rank, score, waitTime, scoreRet.Err()
	} else if scoreRet.Err() != nil {
		return math.MaxInt64, 0, 0, nil
	} else {
		score, waitTime = getQueueScore(int64(scoreRet.Val()))
	}
	return rank, score, waitTime, nil

}

func genScore(score, lastWaitTime, waitTime int64) int64 {
	nowTime := time.Now().Unix()
	if score == 0 {
		score = nowTime
	}
	if lastWaitTime != 0 {
		lastWaitTime = nowTime - score
	}
	return (score << config.WaitTimeBitLen) | (lastWaitTime + waitTime)
}

func getQueueScore(score int64) (int64, int64) {
	var (
		realScore    int64
		waitTime     int64
		waitTimeMask = int64(1<<config.WaitTimeBitLen - 1)
	)
	realScore = score >> config.WaitTimeBitLen
	waitTime = score & waitTimeMask
	return realScore, waitTime
}

func LoginTimeStamp() int64 {
	cfgTimeSlice := int64(GetServer().Conf.Me.TimeStamp)
	if cfgTimeSlice > config.LoginTimeStampMcu {
		cfgTimeSlice = config.LoginTimeStampMcu
	}
	return time.Now().Unix() / cfgTimeSlice
}

func getWindowSize(rankTime int64) int64 {
	windowKey := GetWindowKey(rankTime)
	windowSize, err := redisCluster.Get(context.TODO(), windowKey).Int64()
	if err != nil && err != redis.Nil {
		return 0
	}
	windowSize = config.LoginWindowSize - windowSize
	if windowSize < 0 {
		windowSize = 0
	}
	return windowSize
}

func occupyWindowSeat(rankTime int64) bool {
	windowKey := GetWindowKey(rankTime)
	curValue := redisCluster.Incr(context.TODO(), windowKey)

	if curValue.Err() != nil {
		return false
	}

	cfgTimeSlice := int64(GetServer().Conf.Me.TimeStamp)
	if cfgTimeSlice > config.LoginTimeStampMcu {
		cfgTimeSlice = config.LoginTimeStampMcu
	}

	redisCluster.Expire(context.TODO(), windowKey, time.Duration(2*cfgTimeSlice)*time.Second)

	if curValue.Val() > config.LoginWindowSize {
		return false
	} else {
		return true
	}
}

func checkWaitLevel(account string) (int, int) {

	if !GetServer().Conf.Me.EnableLoginQueue {
		return 0, 0
	}

	var waitLevel int

	rankTime := LoginTimeStamp()

	cfgTimeStamp := int64(GetServer().Conf.Me.TimeStamp)
	if cfgTimeStamp > config.LoginTimeStampMcu {
		cfgTimeStamp = config.LoginTimeStampMcu
	}

	waitLevelRatio := config.LoginTimeStampMcu / cfgTimeStamp

	curRank, curScore, curWaitTime, err := GetQueueRank(account)
	if err != nil {
		if config.QueueLength > config.MaxLoginQueueLength {
			waitRankTime := config.QueueLength / config.LoginWindowSize
			waitLevel = int(waitRankTime/waitLevelRatio + 1)
			return waitLevel, config.MaxRetryWaitTime
		}
		curRank, curScore, curWaitTime = Enqueue(account)
	}

	windowSize := getWindowSize(rankTime)

	if windowSize <= 0 {
		waitRankTime := curRank / config.LoginWindowSize
		waitLevel = int(waitRankTime/waitLevelRatio + 1)
	} else if curRank <= windowSize {
		if occupyWindowSeat(rankTime) {
			waitLevel = 0
		} else {
			waitRankTime := curRank / config.LoginWindowSize
			waitLevel = int(waitRankTime/waitLevelRatio + 1)
		}
	} else {
		waitRankTime := curRank / config.LoginWindowSize
		waitLevel = int(waitRankTime/waitLevelRatio + 1)
	}
	waitRatio := math.Ceil(float64(curRank) / float64(config.LoginWindowSize*waitLevelRatio))
	cfgLoginRetryRatio := GetServer().Conf.Me.RetryRatio

	waitTime := int(math.Ceil(cfgLoginRetryRatio * waitRatio))

	if waitTime > config.MaxRetryWaitTime {
		waitTime = config.MaxRetryWaitTime
	}

	if waitLevel == 0 {
		Dequeue(account)
	}

	if waitLevel > 0 {
		EnqueueWithWaitTime(account, curScore, curWaitTime, int64(waitTime))
	}

	return waitLevel, waitTime
}

func QueueSizeTimer() {
	config.QueueLength = GetQueueLen()
	if GetServer().ProcessId == 1 && config.QueueLength > 0 {
		fmt.Println(fmt.Sprintf("[login] queue长度:%v", config.QueueLength))
	}
}

func clearTimeoutMember() {

	if GetServer().ProcessId != 1 {
		return
	}

	queueName := GetQueueName()
	topData := redisCluster.ZRangeWithScores(context.TODO(), queueName, 0, config.LoginWindowSize).Val()
	if len(topData) <= 0 {
		return
	}

	nowTime := time.Now().Unix()
	clears := make([]interface{}, 0, len(topData))
	for _, data := range topData {
		dev, ok := data.Member.(string)
		if !ok {
			continue
		}
		rankTime, waitTime := getQueueScore(int64(data.Score))
		retryTime := rankTime + waitTime
		if retryTime+config.WaitTimeLimit < nowTime {
			clears = append(clears, dev)
		}
	}

	l := len(clears)
	if l > 0 {
		redisCluster.ZRem(context.TODO(), queueName, clears...)
		fmt.Println(fmt.Sprintf("[login] 清除长度:%v", l))
	}
}

package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	nsqpb "github.com/phuhao00/greatestworks-proto/nsq"
	"google.golang.org/protobuf/proto"
	"greatestworks/aop/nsq"
	"greatestworks/server/login/config"
	"runtime"
	"time"
)

type DailyRegisterController struct {
	LastFetchIncrTime int64
	DailyLimitCnt     int64
	DailyIncrCnt      int64
}

var (
	dailyRegCtrl = DailyRegisterController{
		LastFetchIncrTime: 0,
		DailyLimitCnt:     0,
		DailyIncrCnt:      0,
	}
	loginCounter   int64 // 当前正在登录的人数
	loginCounterTM int64 // 计数的时间
)

func todayDailyRegisterKey(channel string) string {
	today := time.Now().Format("2006-01-02")
	return fmt.Sprintf("RegdCnt:%s:%s", channel, today)
}

func dailyLimitKey(channel string) string {
	return fmt.Sprintf("LimitCnt:%s", channel)
}

func (dc *DailyRegisterController) checkDailyRegCntLimit(channel string) bool {
	limitCntKey := dailyLimitKey(channel)
	regKey := todayDailyRegisterKey(channel)
	limitCnt, e := redisCluster.HGet(context.TODO(), config.DailyLimitKey, limitCntKey).Int64()
	if e != nil && e != redis.Nil {
		return false
	}
	if limitCnt <= 0 {
		return true
	}
	regCnt, re := redisCluster.HGet(context.TODO(), config.DailyLimitKey, regKey).Int64()
	if re != nil && re != redis.Nil {
		return false
	}
	if regCnt > limitCnt {
		return false
	} else {
		return true
	}
}

func (dc *DailyRegisterController) increaseDailyRegCnt(channel string, cnt int64) {
	limitCntKey := dailyLimitKey(channel)
	limitCnt, e := redisCluster.HGet(context.TODO(), config.DailyLimitKey, limitCntKey).Int64()
	if e != nil && e != redis.Nil {
	}

	if limitCnt > 0 {
		regKey := todayDailyRegisterKey(channel)
		result := redisCluster.HIncrBy(context.TODO(), config.DailyLimitKey, regKey, cnt)
		if result == nil {

		}
	}
}

func sendMQ(command nsqpb.NsqCommand, data []byte) {
	msg := nsqpb.ComplexMessage{Cmd: command, Data: data, Time: time.Now().Unix()}
	sdata, _ := proto.Marshal(&msg)
	nsq.PublishAsync(nsq.Logic, nsq.Complex, sdata, nil)
}

func checkSvrStat() int {
	ret := 0
	numGo := runtime.NumGoroutine()
	if numGo > GetServer().Conf.Me.LimitGoroutinesNum*6/5 {
		ret = 2
	} else if numGo >= GetServer().Conf.Me.LimitGoroutinesNum {
		ret = 1
	} else if numGo > GetServer().Conf.Me.LimitGoroutinesNum*2/3 {
		ret = 0
	} else {
	}
	return ret
}

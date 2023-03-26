package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	loginpb "github.com/phuhao00/greatestworks-proto/login"
	nsqpb "github.com/phuhao00/greatestworks-proto/nsq"
	"google.golang.org/protobuf/proto"
	"greatestworks/aop/nsq"
	"greatestworks/internal/rediskey"
	"greatestworks/server/login/config"
	"runtime"
	"strconv"
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
	loginCounter   int64
	loginCounterTM int64
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

func checkInCd(account string) bool {
	return len(redisCluster.Get(context.TODO(), "CdKeyList:"+account).Val()) != 0
}

func checkBanUser(userid uint64) bool {
	banKey := rediskey.MakeBanUserKey(userid)
	ret := redisCluster.Get(context.TODO(), banKey)
	banUserID, err := ret.Uint64()
	if err != nil {
		return false
	}
	if userid == banUserID {
		return true
	}
	return false
}

func checkInWhiteList(account string) bool {
	return redisCluster.SIsMember(context.TODO(), "WhiteList", account).Val()
}

func recommendGatewayForPlayer(zoneId int, userid uint64, inner bool, clientIP, userChannel string) (ok bool, recZone int, ip string, port int) {
	if zoneId > 0 && GetZoneManager().getGateWay(zoneId) == nil {
		return false, 0, "", 0
	}
	key := rediskey.MakeAccountKey(int64(userid))
	prevZoneId, zErr := redisCluster.HGet(context.TODO(), key, "ZoneId").Int()
	if zErr != nil {
		prevZoneId = 0
	}
	if prevZoneId > 0 && zoneId > 0 && prevZoneId != zoneId && rediskey.CheckLive(redisCluster, userid) {
		kickOutInfo := &loginpb.KickOutPlayer{UserID: userid, KickOutReason: loginpb.KickOutReason_RemoteLogin, Reason: "remote login"}
		data, mErr := proto.Marshal(kickOutInfo)
		if mErr == nil {
			sendMQ(nsqpb.NsqCommand_KickOutPlayer, data)
		} else {
		}

		for sleepTimes := 1; prevZoneId > 0 && prevZoneId != zoneId && sleepTimes < 6 && rediskey.CheckLive(redisCluster, userid); sleepTimes++ {
			time.Sleep(2 * time.Second)
			prevZoneId, zErr = redisCluster.HGet(context.TODO(), key, "ZoneId").Int()
			if zErr != nil {
				prevZoneId = 0
			}
		}
	}

	if zoneId == 0 && prevZoneId > 0 {
		zoneId = prevZoneId
	}

	ret := redisCluster.HGet(context.TODO(), key, "Gateway")
	var gatewayEndpoint *config.EndPoint
	var err error
	var exist bool
	if ret.Val() != "" {
		srvID := rediskey.FetchServiceID("gateway-tcp", ret.Val())
		gatewayEndpoint, exist = GetZoneManager().IsExistGateway(zoneId, srvID)
		if !exist {
			gatewayEndpoint, err = GetZoneManager().RecommendGateway(zoneId)
			if err != nil {
				return false, 0, "", 0
			}
		}
	} else {
		gatewayEndpoint, err = GetZoneManager().RecommendGateway(zoneId)
		if err != nil {
			return false, 0, "", 0
		}
	}

	recZone = gatewayEndpoint.ZoneId
	if inner && len(gatewayEndpoint.InnerIP) >= len("0.0.0.0") {
		ip = gatewayEndpoint.InnerIP
	} else {
		ip = gatewayEndpoint.IP
	}
	port = gatewayEndpoint.Port

	tcpAddr := gatewayEndpoint.IP + ":" + strconv.Itoa(gatewayEndpoint.Port)
	hSetVal := make(map[string]interface{}, 2)
	hSetVal["Gateway"] = tcpAddr
	hSetVal["ZoneId"] = gatewayEndpoint.ZoneId
	err = redisCluster.HMSet(context.TODO(), key, hSetVal).Err()
	if err != nil {
		return false, 0, "", 0
	}

	ok = true
	GetZoneManager().UpdateGatewayLocalWeight(gatewayEndpoint)
	return
}

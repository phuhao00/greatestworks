package main

import (
	"context"
	"encoding/json"
	loginpb "github.com/phuhao00/greatestworks-proto/login"
	"greatestworks/aop/fn"
	"greatestworks/internal/rediskey"
	"greatestworks/server/login/config"
	"net/http"
	"strings"
	"time"
)

func RandomName() {
	//todo 随机名字
}

func Register() {
	preRegister()
	registerReward()
}

func preRegister() {
	//todo 预注册

}

func registerReward() {

}

func returnHandler(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	if b, err := json.Marshal(data); err == nil {
		w.Write(b)
	} else {
		http.Error(w, err.Error(), 400)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	loginInfo := &loginpb.LoginData{Result: config.Succ}
	loginInfo.ServerTime = time.Now().Unix()
	loginInfo.RegRegTimeStart = config.StartPreRegisterTime
	loginInfo.RegRegTimeEnd = config.EndPreRegisterTime
	loginInfo.ServerOpenTime = config.ServerOpenTime
	loginInfo.IsIpInWhiteList = gm.IsIpInWhiteList(fn.ClientIP(r))
	loginInfo.ShuShuGameID = ""
	defer returnHandler(w, r, loginInfo)
	if r.Body == nil {
		loginInfo.Result = config.UnknownErr
		return
	}

	var accData config.AccountData
	err := json.NewDecoder(r.Body).Decode(&accData)
	if err != nil {
		loginInfo.Result = config.DecodeErr
		return
	}

	if len(accData.Account) == 0 ||
		len(accData.Password) == 0 ||
		len(accData.Sign) == 0 {
		loginInfo.Result = config.UnknownErr
		return
	}

	if fn.CheckSign(accData.Account, accData.Sign) == false {
		loginInfo.Result = config.VerifyTokenErr
		return
	}

	if GetServer().Conf.WhiteList.Check && !checkInWhiteList(accData.Account) {
		loginInfo.Result = config.WhiteListErr
		return
	}

	waitLevel, waitTime := checkWaitLevel(accData.Account)
	if waitLevel != 0 {
		loginInfo.Result = config.LoginBusy
		loginInfo.BusyLevel = int32(waitLevel)
		loginInfo.BusyWaitTime = int32(waitTime)
		return
	}

	if GetServer().Conf.WhiteList.TokenCheck {
		var ltToken, ltSid string
		if len(accData.Token) > 0 && len(accData.Sid) > 0 {
			ltToken = accData.Token
			ltSid = accData.Sid
		} else {
			ltSid, ltToken = GetToken()
		}
		if len(ltToken) > 0 && len(ltSid) > 0 {
			//todo third party token check
		}
	}

	if !dailyRegCtrl.checkDailyRegCntLimit("channel") {
		loginInfo.Result = config.DailyIncrOver
		return
	}
	//todo db check exist
	dailyRegCtrl.increaseDailyRegCnt("channel", 1)

	config.GetIdxByLimitInfo(nil, 0)

	//todo ..
	loginInfo.IsPreReged = true
	//todo 预注册奖励
	if accData.ZoneId > 0 && !GetZoneManager().ZoneExistForOnline(accData.ZoneId) {
		loginInfo.Result = config.ZoneError
		return
	}

	loginInfo.UserID = 0
	loginInfo.Token = ""
	loginInfo.SessionID = loginInfo.Token
	if accData.ZoneId == 0 {
		loginInfo.ZoneId = int32(GetZoneManager().recommendZone())
	} else {
		loginInfo.ZoneId = int32(accData.ZoneId)
	}
	inner := false
	if !GetServer().Conf.WhiteList.UsePublicNetwork &&
		strings.HasPrefix(accData.Account, "robot") && accData.Password == "robot" {
		inner = true
	}

	ok, zid, ip, port := recommendGatewayForPlayer(int(loginInfo.ZoneId), loginInfo.UserID, inner, fn.ClientIP(r), "channel")
	if ok == true {
		loginInfo.IP, loginInfo.Port = ip, int32(port)
		loginInfo.ZoneId = int32(zid)
		loginInfo.ZoneList = GetZoneManager().GetZonesInfo()

		rok, rList := GetZoneManager().RecommendZoneWorlds(zid)
		if rok != true {
		}
		loginInfo.RecommendWorld = rList

		loginInfo.WorldList = GetZoneManager().GetZoneOnlineList(int(loginInfo.ZoneId))
		loginInfo.Result = config.Succ
		tokenKey := rediskey.MakeTokenKey(loginInfo.UserID)
		redisCluster.Set(context.TODO(), tokenKey, loginInfo.Token, config.TokenExpireDuration)
		if inner {
			key := rediskey.MakeAccountKey(0) //todo
			redisCluster.HSet(context.TODO(), key, "AccountRobot", accData.Account)
		}
	} else {
		loginInfo.Result = config.ZoneError
	}
}

func ThirdPartyLogin() {

}

func GetGateWay() {

}

func GetWorldServers() {

}

func HealthyCheck() {

}

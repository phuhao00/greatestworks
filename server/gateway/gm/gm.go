package gm

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/phuhao00/greatestworks-proto/chat"
	"github.com/phuhao00/greatestworks-proto/messageId"
	"github.com/phuhao00/greatestworks-proto/server_common"
	"greatestworks/aop/fn"
	"greatestworks/aop/logger"
	"greatestworks/aop/redis"
	"greatestworks/server/gateway/client"
	"greatestworks/server/gateway/server"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GM struct {
	IsOpenNow       bool
	localStartTM    int64
	fetchInterval   int64
	localServerId   string
	localServerType string

	opt         string
	afterSec    int64
	optSec      int64
	optsDesc    string
	kickoffUser int64

	lastFetchTime  int64
	optTime        int64
	lastNotifyTime int64
	showTips       string

	ipWhiteList  *sync.Map
	uidWhiteList *sync.Map

	localZone int32
}

var (
	GMInstance *GM
)

func Init(startTM int64, serverId string, serverType string, fetchInterval int64, IsOpenNow bool) {
	GMInstance = &GM{}
	GMInstance.IsOpenNow = IsOpenNow
	GMInstance.localStartTM = startTM
	GMInstance.fetchInterval = fetchInterval
	GMInstance.localServerId = strings.ToLower(serverId)
	GMInstance.localServerType = strings.ToLower(serverType)

	GMInstance.opt = ""
	GMInstance.afterSec = 0
	GMInstance.optSec = 0
	GMInstance.kickoffUser = 0

	GMInstance.lastFetchTime = 0
	GMInstance.lastNotifyTime = 0
	GMInstance.optTime = 0
	GMInstance.showTips = ""

	GMInstance.ipWhiteList = &sync.Map{}
	defaultIP := "127.0.0.1"
	GMInstance.ipWhiteList.Store(defaultIP, true)

	GMInstance.uidWhiteList = &sync.Map{}

	GMInstance.localZone = int32(server.GetServer().Config.Global.ZoneId)
}

func (gm *GM) OnTimer() {
	now := time.Now().Unix()
	if now-gm.lastFetchTime >= gm.fetchInterval {
		gm.lastFetchTime = now

		gm.fetchInfoFromRedis(now)
	}

	if len(gm.opt) > 0 {
		gm.countDownUpdate(now)
	}
}

func (gm *GM) fetchInfoFromRedis(now int64) {
	gm.fetchOpenCloseDoorInfo(now)

	gm.fetchIpWhiteList(now)
	gm.fetchUseridWhiteList(now)
}

func (gm *GM) fetchIpWhiteList(now int64) {
	ipList := redis.NonCacheRedis().HGetAll(context.TODO(), "GatewayIpWhiteList").Val()
	for k, v := range ipList {
		if v == "1" {
			gm.ipWhiteList.Store(k, true)
		} else {
			gm.ipWhiteList.Store(k, false)
		}
	}
}

func (gm *GM) fetchUseridWhiteList(now int64) {
	uidList := redis.NonCacheRedis().HGetAll(context.TODO(), "GatewayUidWhiteList").Val()
	for k, v := range uidList {
		if v == "1" {
			gm.uidWhiteList.Store(k, true)
		} else {
			gm.uidWhiteList.Store(k, false)
		}
	}
}

func (gm *GM) fetchOpenCloseDoorInfo(now int64) {
	notify := redis.NonCacheRedis().HGet(context.TODO(), "DoorOpenClose", gm.localServerType).Val()

	if len(notify) == 0 {
		return
	}

	req := &DealGmCommandRequest{}
	err := json.Unmarshal([]byte(notify), req)
	if err != nil {
		logger.Error("[fetchOpenCloseDoorInfo] got data from redis cannot json.Unmarshal : %v", notify)
		return
	}

	target := strings.ToLower(req.Target)
	if target != "gateway" {
		logger.Debug("[fetchOpenCloseDoorInfo] the target server is: %v but not this gateway, so escape.",
			req.Target)
		return
	}

	if req.Zones == "all" {
	} else {
		wanted := false
		zids := fn.SplitStringToInt32Slice(req.Zones, ",")
		for _, zid := range zids {
			if zid >= 0 && zid == gm.localZone {
				wanted = true
				break
			}
		}
		if !wanted {
			//logger.Debugf("fetchOpenCloseDoorInfo wanted zones: %v but not this, so escape.", req.Zones)
			return
		}
	}
	if req.OptTime <= gm.optTime {
		//logger.Errorf("fetchOpenCloseDoorInfo got notify data time %v has been process in pretime : %v, so escape.",
		//	req.OptTime, gm.optTime)
		return
	}

	logger.Debug("[fetchOpenCloseDoorInfo]: %v", notify)
	opt := strings.ToLower(req.Opt)
	if opt == "close_door" {
		gm.updateCloseDoorInfo(now, req)
	} else if opt == "open_door" {
		gm.updateOpenDoorInfo(now, req)
	} else {
		logger.Error("[fetchOpenCloseDoorInfo] got unSupported gameCenter gm opt: %v", opt)
	}
}

func (gm *GM) updateCloseDoorInfo(now int64, info *DealGmCommandRequest) {
	logger.Debug("[updateCloseDoorInfo]...")
	afterMin, err := strconv.ParseInt(info.AfterMin, 10, 64)
	if err != nil {
		logger.Error("[updateCloseDoorInfo] but info.AfterMin %v error: %v",
			info.AfterMin, err)
		return
	}
	gm.afterSec = afterMin*60 + int64(rand.Intn(60)+1)

	opsMin, err := strconv.ParseInt(info.OpsMin, 10, 64)
	if err != nil {
		logger.Error("[updateCloseDoorInfo] but info.OpsMin %v error: %v",
			info.OpsMin, err)
		return
	}
	gm.optSec = opsMin * 60

	kickoff, err := strconv.ParseInt(info.KickoffUser, 10, 64)
	if err != nil {
		logger.Error("[updateCloseDoorInfo] but info.kickoffUser %v error: %v",
			info.KickoffUser, err)
		return
	}
	gm.optTime = info.OptTime
	gm.opt = strings.ToLower(info.Opt)
	gm.kickoffUser = kickoff
	gm.optsDesc = info.OpsDesc
	logger.Debug("[updateCloseDoorInfo] finished.")
}

func (gm *GM) updateOpenDoorInfo(now int64, info *DealGmCommandRequest) {
	logger.Debug("[updateOpenDoorInfo] now!")

	gm.IsOpenNow = true
	gm.optTime = info.OptTime
	gm.opt = ""
	gm.afterSec = 0
	gm.optSec = 0
	gm.kickoffUser = 0
	gm.showTips = ""
}

func (gm *GM) countDownUpdate(now int64) {
	if gm.opt == "close_door" {
		gm.closeDoorCountDown(now)
	} else if gm.opt == "open_door" {
		gm.openDoorCountDown(now)
	} else {
		logger.Error("[countDownUpdate] but now opt is:%v, unsupportted.", gm.opt)
	}
}

func (gm *GM) closeDoorCountDown(now int64) {
	if gm.opt != "close_door" {
		logger.Error("[closeDoorCountDown] but now opt is:%v, return.", gm.opt)
		return
	}
	closeTime := gm.optTime + gm.afterSec
	closeRestSec := closeTime - now

	restNotifySec := now - gm.lastNotifyTime

	var tip string
	if closeRestSec >= 5*60 {
		if restNotifySec >= 2*60 {
			if gm.optSec > 0 {
				tip = fmt.Sprintf("%d 分钟后游戏进入维护, 预计维护时长 %d 分钟。%s", closeRestSec/60, gm.optSec/60, gm.optsDesc)
			} else {
				tip = fmt.Sprintf("%d 分钟后游戏进入维护。%s", closeRestSec/60, gm.optsDesc)
			}
			gm.updateNotifyTip(now, tip)

			gm.notifyAllOnlinePlayer(now, tip)
		}
	} else if closeRestSec >= 3*60 {
		if restNotifySec >= 1*60 {
			if gm.optSec > 0 {
				tip = fmt.Sprintf("%d 分钟后游戏进入维护, 预计维护时长 %d 分钟。%s", closeRestSec/60, gm.optSec/60, gm.optsDesc)
			} else {
				tip = fmt.Sprintf("%d 分钟后游戏进入维护。%s", closeRestSec/60, gm.optsDesc)
			}
			gm.updateNotifyTip(now, tip)
			gm.notifyAllOnlinePlayer(now, tip)
		}
	} else if closeRestSec >= 1*60 {
		if restNotifySec >= 30 {
			if gm.optSec > 0 {
				tip = fmt.Sprintf("%d 分钟后游戏进入维护, 预计维护时长 %d 分钟。%s", closeRestSec/60, gm.optSec/60, gm.optsDesc)
			} else {
				tip = fmt.Sprintf("%d 分钟后游戏进入维护。%s", closeRestSec/60, gm.optsDesc)
			}
			gm.updateNotifyTip(now, tip)

			gm.notifyAllOnlinePlayer(now, tip)
		}
	} else if closeRestSec > 0 {
		if restNotifySec >= 10 {
			if gm.optSec > 0 {
				tip = fmt.Sprintf("%d 秒后游戏进入维护, 预计维护时长 %d 分钟。%s", closeRestSec, gm.optSec/60, gm.optsDesc)
			} else {
				tip = fmt.Sprintf("%d 秒后游戏进入维护。%s", closeRestSec, gm.optsDesc)
			}
			gm.updateNotifyTip(now, tip)

			gm.notifyAllOnlinePlayer(now, tip)
		}
	} else {
		gm.IsOpenNow = false
		if gm.kickoffUser != 0 {
			gm.kickOffOnlinePlayer()
		}
		if gm.optSec > 0 {
			gm.opt = "open_door"
		} else {
			gm.opt = ""
		}
	}

	if len(tip) > 0 {
		logger.Info(tip)
	}
}

func (gm *GM) openDoorCountDown(now int64) {
	if gm.opt != "open_door" {
		logger.Error("[openDoorCountDown] but now opt is:%v, return.", gm.opt)
		return
	}
	openTime := gm.optTime + gm.afterSec + gm.optSec
	openRestSec := openTime - now
	if openRestSec <= 0 {
		gm.IsOpenNow = true
		gm.opt = ""
		gm.afterSec = 0
		gm.optSec = 0
		gm.kickoffUser = 0
		gm.showTips = ""
		logger.Info("[openDoorCountDown] open door now.")
	} else {
		var tip string
		if openRestSec > 60 {
			tip = fmt.Sprintf("系统维护中, 预计维护时长 %d 分钟。%s", openRestSec/60, gm.optsDesc)
		} else {
			tip = fmt.Sprintf("系统维护中, 敬请期待。")
		}
		gm.updateNotifyTip(now, tip)
	}
}

func (gm *GM) updateNotifyTip(now int64, tip string) string {
	logger.Debug("[updateNotifyTip]: %v", tip)
	if gm.IsOpenNow {
		return ""
	}
	gm.showTips = tip

	return gm.showTips
}

func (gm *GM) notifyAllOnlinePlayer(now int64, tips string) {
	gm.lastNotifyTime = now
	cmd := messageId.MessageId_SCSystemMessage
	msg := &chat.SCSystemMessage{MsgType: 2, Content: tips}
	client.GetMe().SendMsgToAllPlayer(cmd, msg)
}

func (gm *GM) kickOffOnlinePlayer() {
	logger.Debug("[kickOffOnlinePlayer]")
	cmd := messageId.MessageId_SCKick
	msg := &server_common.SCKick{Kick: server_common.KickReason_ServerClose, Contend: "服务器关闭进入维护..."}
	client.GetMe().SendMsgToAllPlayer(cmd, msg)
}

func (gm *GM) IsIpInWhiteList(ip string) bool {
	ret, ok := gm.ipWhiteList.Load(ip)
	return ok && ret.(bool)
}

func (gm *GM) IsUidInWhiteList(uid string) bool {
	ret, ok := gm.uidWhiteList.Load(uid)
	return ok && ret.(bool)
}

package mail

import (
	"context"
	"github.com/phuhao00/greatestworks-proto/mail"
	"go.mongodb.org/mongo-driver/bson"
	"greatestworks/aop/container"
	"greatestworks/aop/logger"
	"greatestworks/aop/mongo"
	"greatestworks/internal/communicate/player"
	"time"
)

type System struct {
	owner   *player.Player
	sendCnt int16
	data    container.IDelegate
}

func NewMailSystem(player *player.Player) *System {
	email := &System{owner: player, data: container.NewMailContainer(player.PlayerID)}
	email.loadFromMgo()
	return email
}

func (email *System) loadFromMgo() {
	data := mongo.MailSystem{}
	query := bson.M{mongo.PrimaryKey: email.owner.PlayerID}
	result := mongo.Client.FindOne(context.TODO(), data.DB(), data.C(), query)
	if result.Err() != nil {
		_, err := mongo.Client.InsertOne(context.TODO(), data.DB(), data.C(), &mongo.MailSystem{OwnerID: email.owner.PlayerID})
		if err != nil {
			logger.Error("")

		}
	}
	email.sendCnt = data.SendCnt
	email.data.Set("normal", data.Normal)
	email.data.Set("collect", data.Collect)
	email.data.Set("recycle", data.Recycle)
	email.data.Set("history", data.History)
	data.Normal = make([]mongo.MailInfo, 0)
	data.Collect = make([]mongo.MailInfo, 0)
	data.Recycle = make([]mongo.MailInfo, 0)
	data.History = make([]uint64, 0)
	email.delExpireRecycle()
	email.loadGlobalEmail()
	email.balanceNormal(false)
}

func (email *System) syncSendCnt() {

}

func (email *System) clearSendCnt() {
	email.updateSendCnt(0)
	email.syncSendCnt()
}

func (email *System) updateSendCnt(val int16) {
	email.sendCnt = val
	set := bson.M{"scnt": email.sendCnt}
	email.syncMgo(bson.M{mongo.PrimaryKey: email.owner.PlayerID}, bson.M{"$set": set})
}

func (email *System) delExpireRecycle() {
	now := time.Now().Unix()
	expire := int64(7 * 24 * 60 * 60) // 7天
	dels := []mongo.MailInfo{}
	data := email.data.Get("recycle").([]mongo.MailInfo)
	for _, val := range data {
		if val.MTime+expire < now {
			dels = append(dels, mongo.MailInfo{MUuid: val.MUuid, MType: int32(mail.EmailType_RECYCLE)})
		}
	}
}

func (email *System) delTheMail(muuid uint64, mtype mail.EmailType) {

}

func (email *System) loadGlobalEmail() {

}

func (email *System) balanceNormal(sendmsg bool) {

}

func (email *System) chgEmailStatus(uuids []uint64, mode mail.EmailType, status mail.EmailStatus) {

}

func (email *System) move(smode, dmode mail.EmailType, suuids []uint64) bool {

	return false
}

func (email *System) checkInHistory(euuid uint64) bool {
	data := email.data.Get("history").([]uint64)
	for _, id := range data {
		if id == euuid {
			return true
		}
	}
	return false
}

func (email *System) addHistory(euuid uint64, mailid uint32) {

}

func (email *System) packMailInfo(info *mongo.MailInfo) *mail.MailInfo {
	data := &mail.MailInfo{}
	for _, val := range info.MItems {
		data.Goods = append(data.Goods, &mail.GoodsInfo{ItemId: val.ItemId, Num: val.Num})
	}
	return data
}

func (email *System) syncMgo(query, option interface{}) {
}

func (email *System) addMail(mail *mail.MailInfo) {

}

func (email *System) checkTime(reg []int64, vali []int64) bool {

	return true
}

func (email *System) canAdd(euuid uint64) bool {
	if email.checkInHistory(euuid) {
		return false
	}
	info := email.data.Get("normal").([]mongo.MailInfo)
	info = append(info, email.data.Get("collect").([]mongo.MailInfo)...)
	info = append(info, email.data.Get("recycle").([]mongo.MailInfo)...)
	for _, val := range info {
		if val.MUuid == euuid {
			return false
		}
	}
	return true
}

func (email *System) checkChan(chans []string) bool {
	return false
}

func (email *System) checkVersion(sv string) bool {
	if len(sv) == 0 {
		return true
	}
	return false
}

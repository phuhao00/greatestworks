package container

import (
	"context"
	"errors"
	"github.com/phuhao00/greatestworks-proto/mail"
	"go.mongodb.org/mongo-driver/bson"
	"greatestworks/aop/fn"
	"greatestworks/aop/logger"
	"greatestworks/aop/mongo"
	"reflect"
)

var empty = errors.New("empty")

type MailContainer struct {
	data *mongo.MailSystem
}

func NewMailContainer(userid uint64) *MailContainer {
	return &MailContainer{data: &mongo.MailSystem{OwnerID: userid}}
}

func (u *MailContainer) Save(query, update interface{}) {
	_, err := mongo.Client.UpdateOne(context.TODO(), u.data.DB(), u.data.C(), query, update)
	if err != nil {
		logger.Error("err:%v", err.Error())

	}
}

func (u *MailContainer) Set(tag string, val interface{}) {
	for pos := 0; pos < reflect.TypeOf(*u.data).NumField(); pos++ {
		if reflect.TypeOf(*u.data).Field(pos).Tag.Get("bson") == tag {
			if reflect.ValueOf(u.data).Elem().Field(pos).Type() == reflect.TypeOf(val) {
				reflect.ValueOf(u.data).Elem().Field(pos).Set(reflect.ValueOf(val))
			}
		}
	}
}

func (u *MailContainer) GetItem(val interface{}) interface{} {
	ret := mongo.MailInfo{}
	if msg, ok := val.(mongo.MailInfo); ok {
		if mail.EmailType(msg.MType) == mail.EmailType_COLLECT {
			data, _, _ := u.getCollect(msg.MUuid)
			err := fn.Copy(data, &ret)
			if err != nil {
				logger.Error("err:%v", err.Error())
			}
		} else if mail.EmailType(msg.MType) == mail.EmailType_NORMAL {
			data, _, _ := u.getNormal(msg.MUuid)
			err := fn.Copy(data, &ret)
			if err != nil {
				logger.Error("err:%v", err.Error())
			}
		} else if mail.EmailType(msg.MType) == mail.EmailType_RECYCLE {
			data, _, _ := u.getRecycle(msg.MUuid)
			err := fn.Copy(data, &ret)
			if err != nil {
				logger.Error("err:%v", err.Error())
			}
		}
	}
	return ret
}

func (u *MailContainer) Get(tag string) interface{} {
	for pos := 0; pos < reflect.TypeOf(*u.data).NumField(); pos++ {
		if reflect.TypeOf(*u.data).Field(pos).Tag.Get("bson") == tag {
			return reflect.ValueOf(*u.data).Field(pos).Interface()
		}
	}
	return nil
}

func (u *MailContainer) Add(val interface{}) {
	addMail := make(map[string][]mongo.MailInfo)
	if val, ok := val.([]mongo.MailInfo); ok {
		for _, info := range val {
			logger.Debug("MailContainer add 222 %v", info)
			if mail.EmailType(info.MType) == mail.EmailType_COLLECT {
				u.addCollect(&info)
				key := u.getKey(mail.EmailType_COLLECT)
				addMail[key] = append(addMail[key], info)
			} else if mail.EmailType(info.MType) == mail.EmailType_NORMAL {
				u.addNormal(&info)
				key := u.getKey(mail.EmailType_NORMAL)
				addMail[key] = append(addMail[key], info)
			} else if mail.EmailType(info.MType) == mail.EmailType_RECYCLE {
				u.addRecycle(&info)
				key := u.getKey(mail.EmailType_RECYCLE)
				addMail[key] = append(addMail[key], info)
			}
		}
	}
	for key, val := range addMail {
		option := bson.M{"$push": bson.M{key: bson.M{"$each": val}}}
		u.Save(bson.M{mongo.PrimaryKey: u.data.OwnerID}, option)
	}
}

func (u *MailContainer) Del(val interface{}) {
	delMail := make(map[string][]uint64)
	logger.Debug("MailContainer %v", reflect.TypeOf(val))
	if val, ok := val.([]mongo.MailInfo); ok {
		for _, info := range val {
			if mail.EmailType(info.MType) == mail.EmailType_COLLECT {
				u.delCollect(info.MUuid)
				key := u.getKey(mail.EmailType_COLLECT)
				delMail[key] = append(delMail[key], info.MUuid)
			} else if mail.EmailType(info.MType) == mail.EmailType_NORMAL {
				u.delNormal(info.MUuid)
				key := u.getKey(mail.EmailType_NORMAL)
				delMail[key] = append(delMail[key], info.MUuid)
			} else if mail.EmailType(info.MType) == mail.EmailType_RECYCLE {
				u.delRecycle(info.MUuid)
				key := u.getKey(mail.EmailType_RECYCLE)
				delMail[key] = append(delMail[key], info.MUuid)
			}
		}
	}
	for key, uuids := range delMail {
		logger.Info("MailContainer del %v %v", key, uuids)
		option := bson.M{"$pull": bson.M{key: bson.M{"muid": bson.M{"$in": uuids}}}}
		u.Save(bson.M{mongo.PrimaryKey: u.data.OwnerID}, option)
	}
}

func (u *MailContainer) getKey(mode mail.EmailType) string {
	keys := []string{"normal", "collect", "recycle"}
	return keys[mode]
}

func (u *MailContainer) getNormal(uuid uint64) (*mongo.MailInfo, int, error) {
	for index, val := range u.data.Normal {
		if val.MUuid == uuid {
			return &u.data.Normal[index], index, nil
		}
	}
	return nil, -1, empty
}

func (u *MailContainer) getCollect(uuid uint64) (*mongo.MailInfo, int, error) {
	for index, val := range u.data.Collect {
		if val.MUuid == uuid {
			return &u.data.Collect[index], index, nil
		}
	}
	return nil, -1, empty
}

func (u *MailContainer) getRecycle(uuid uint64) (*mongo.MailInfo, int, error) {
	for index, val := range u.data.Recycle {
		if val.MUuid == uuid {
			return &u.data.Recycle[index], index, nil
		}
	}
	return nil, -1, empty
}

func (u *MailContainer) addNormal(info *mongo.MailInfo) {
	u.data.Normal = append(u.data.Normal, *info)
}

func (u *MailContainer) addCollect(info *mongo.MailInfo) {
	u.data.Collect = append(u.data.Collect, *info)
}

func (u *MailContainer) addRecycle(info *mongo.MailInfo) {
	u.data.Recycle = append(u.data.Recycle, *info)
}

func (u *MailContainer) delNormal(uuid uint64) bool {
	for index, val := range u.data.Normal {
		if val.MUuid == uuid {
			u.data.Normal = append(u.data.Normal[:index], u.data.Normal[index+1:]...)
			return true
		}
	}
	return false
}

func (u *MailContainer) delCollect(uuid uint64) bool {
	for index, val := range u.data.Collect {
		if val.MUuid == uuid {
			u.data.Collect = append(u.data.Collect[:index], u.data.Collect[index+1:]...)
			return true
		}
	}
	return false
}

func (u *MailContainer) delRecycle(uuid uint64) bool {
	for index, val := range u.data.Recycle {
		if val.MUuid == uuid {
			u.data.Recycle = append(u.data.Recycle[:index], u.data.Recycle[index+1:]...)
			return true
		}
	}
	return false
}

func (u *MailContainer) SetItem(val interface{}, items interface{}) {
	status, ok := val.(mail.EmailStatus)
	if !ok {
		logger.Error("SetItem err 1 %v", reflect.TypeOf(val))
		return
	}
	if item, ok := items.([]mongo.MailInfo); ok {
		updates := make([]interface{}, 0, len(item))
		for _, info := range item {
			if mail.EmailType(info.MType) == mail.EmailType_COLLECT {
				if _, index, err := u.getCollect(info.MUuid); err == nil {
					u.data.Collect[index].MStatus = uint32(status)
					update := bson.M{"$set": bson.M{"collect.$.stat": uint32(status)}}
					query := bson.M{mongo.PrimaryKey: u.data.OwnerID,
						"collect.muid": info.MUuid}

					updates = append(updates, query, update)
				}
			} else if mail.EmailType(info.MType) == mail.EmailType_NORMAL {
				if _, index, err := u.getNormal(info.MUuid); err == nil {
					u.data.Normal[index].MStatus = uint32(status)

					update := bson.M{"$set": bson.M{"normal.$.stat": uint32(status)}}
					query := bson.M{mongo.PrimaryKey: u.data.OwnerID,
						"normal.muid": info.MUuid}
					updates = append(updates, query, update)
				}
			} else if mail.EmailType(info.MType) == mail.EmailType_RECYCLE {
				if _, index, err := u.getRecycle(info.MUuid); err == nil {
					u.data.Recycle[index].MStatus = uint32(status)
					update := bson.M{"$set": bson.M{"recycle.$.stat": uint32(status)}}
					query := bson.M{mongo.PrimaryKey: u.data.OwnerID,
						"recycle.muid": info.MUuid}
					updates = append(updates, query, update)
				}
			}
		}
		if len(updates) > 0 {
			//mongo.Client.BulkWrite(context.TODO(),u.data.DB(),u.data.C(), updates,nil)
		}
	}
}

package server

import (
	"context"
	"greatestworks/aop/redis"
	"strconv"
	"strings"
	"sync"
)

var (
	systemCloseList *SystemCloseList
	closeListOnce   sync.Once
)

type SystemCloseList struct {
	closeList *sync.Map
}

func systemCloseListGetMe() *SystemCloseList {
	closeListOnce.Do(func() {
		systemCloseList = &SystemCloseList{
			closeList: &sync.Map{},
		}
	})
	return systemCloseList
}

func (scl *SystemCloseList) loadFromRedis() {
	closeList := redis.NonCacheRedis().HGetAll(context.TODO(), redis.SystemCloseList).Val()
	if closeList == nil {
		return
	}

	if len(closeList) == 0 {
		return
	}

	for k, v := range closeList {

		isClose, err := strconv.ParseBool(v)
		if err != nil {
			continue
		}

		if isClose {
			keys := strings.Split(k, ":")
			if len(keys) != 2 {
				continue
			}

			sysID, err := strconv.ParseUint(keys[1], 10, 64)
			if err != nil {
				continue
			}

			scl.closeList.Store(uint32(sysID), isClose)
		}
	}
}

func (scl *SystemCloseList) setCloseSystem(sys uint32, close bool) {
	if close {
		scl.closeList.Store(sys, close)
	} else {
		scl.delCloseSystem(sys)
	}
}

func (scl *SystemCloseList) delCloseSystem(sys uint32) {
	_, ok := scl.closeList.Load(sys)
	if !ok {
		return
	}

	scl.closeList.Delete(sys)
}

func (scl *SystemCloseList) getSystemInCloseList(sys uint32) bool {
	sysClose, ok := scl.closeList.Load(sys)
	if !ok {
		return false
	}

	return sysClose.(bool)
}

func (scl *SystemCloseList) getSystemCloseList() []uint32 {
	closeList := make([]uint32, 0, 8)
	scl.closeList.Range(func(k, v interface{}) bool {

		if v.(bool) {
			closeList = append(closeList, k.(uint32))
		}

		return true
	})

	return closeList
}

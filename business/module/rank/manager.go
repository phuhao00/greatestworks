package rank

import (
	"fmt"
	"greatestworks/aop/redis"
	"sync"
	"time"
)

type Manager struct {
	initFlag          bool
	bMutex            sync.Mutex
	rlsMutex          sync.Mutex
	cache             sync.Map
	rankLastScoreList map[uint32]int64
	blackList         map[uint32]*BlackList
}

func (m *Manager) Init() error {

	stringStartTime := "2023-01-01 00:00:00"
	loc, _ := time.LoadLocation("Local")
	start, err := time.ParseInLocation("2006-01-02 15:04:05", stringStartTime, loc)

	if err != nil {
		return err
	}

	startTime = start.Unix()
	stringFinalTime := "2100-01-01 00:00:00"
	final, err := time.ParseInLocation("2006-01-02 15:04:05", stringFinalTime, loc)

	if err != nil {
		return err
	}
	finalTime = final.Unix()
	m.blackList = make(map[uint32]*BlackList, 16)
	m.rankLastScoreList = make(map[uint32]int64)

	m.initFlag = true

	return nil
}

func (m *Manager) GetRank(rankId uint32, playerId uint64, sortType SortType) {
	conf := &Config{}
	rankName := conf.getRankName(rankId)
	if sortType == Aes {
		redis.GetMockInstance()
		//todo ZRank(rankName, playerId)
	}
	if sortType == Des {
		//todo ZRevRank(rankName, playerId)
	}
	//todo ZScore(rankName, playerId)
	fmt.Println(rankName)

}

func (m *Manager) GetZCard(rankId uint32) int64 {
	conf := &Config{}
	rankName := conf.getRankName(rankId)
	//todo  ZCard(rankName)
	fmt.Println(rankName)
	return 0
}

func (m *Manager) Clear(rankId uint32) {

}

func (m *Manager) Save() {

}

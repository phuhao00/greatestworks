package rank

import "fmt"

type Manager struct {
}

func (m *Manager) Init() {

}

func (m *Manager) GetRank(rankId uint32, playerId uint64, sortType SortType) {
	conf := &Config{}
	rankName := conf.getRankName(rankId)
	if sortType == Des {
		//todo ZRank(rankName, playerId)
	}
	if sortType == Aes {
		//todo ZRevRank(rankName, playerId)
	}
	//todo ZScore(rankName, playerId)
	fmt.Println(rankName)
}

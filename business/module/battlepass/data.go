package battlepass

// BattlePass 战斗通行证
type BattlePass struct {
	HistorySeasons []*Season //保留多少根据策划要求
	Cur            *Season   //当前赛季
}

// Season 赛季
type Season struct {
	ConfigId uint32
	Score    uint64
	Records  map[CardCategory][]uint8
	Card     map[CardCategory]bool
}

func (p *BattlePass) Tag() {

}

func (p *BattlePass) Restore() {

}

func (p *BattlePass) Reload() {

}

func (p *BattlePass) Refresh() {

}

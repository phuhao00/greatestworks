package player

type BaseInfo struct {
	UId    uint64 `json:"uid"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender int    `json:"gender"`
}

func (p *Player) Load() {

}

func (p *Player) Save() {
	
}

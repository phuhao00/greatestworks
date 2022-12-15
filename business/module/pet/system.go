package pet

// System pet system
type System struct {
	pets       map[uint32]Abstract   //宠物实例
	Pictorials map[uint32]*Pictorial //图鉴
	Fragments  map[uint32]*Fragment  //碎片
	Skins      map[uint32]*Skin      //皮肤
}

func NewSystem() *System {
	return &System{
		pets:       nil,
		Pictorials: nil,
		Fragments:  nil,
		Skins:      nil,
	}
}

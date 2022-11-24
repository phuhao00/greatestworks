package pet

// System pet system
type System struct {
	pets       map[uint32]Abstract
	Pictorials map[uint32]*Pictorial
	Fragments  map[uint32]*Fragment
	Skins      map[uint32]*Skin
}

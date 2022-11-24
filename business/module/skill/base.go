package skill

type Base struct {
	Id       uint32 `json:"id"`
	IsUnlock uint32 `json:"isUnlock"`
	IsFixed  bool   `json:"isFixed"`
}

func (b *Base) GetAttack() int64 {
	//TODO implement me
	panic("implement me")
}

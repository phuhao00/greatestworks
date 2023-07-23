package card

type IPlayer interface {
	GetCard() *Data
	Reward(conf interface{})
}

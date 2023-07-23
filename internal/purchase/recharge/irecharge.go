package recharge

type IRecharge interface {
	CheckProductLimit(id uint32) bool
	CreateBuyOrder(id uint32, channel string, useCoupon bool)
	ProcOpenOrder(order *Order) bool
	CloseOrder(order *Order, closeStatus OrderStatus, IsFirstBuy bool)
	GetPayProduct(order *Order)
	SyncPaymentData()
	AddRecord(id uint32, price uint32)
	ClearDailyRecord()
	GetDailyPayCount(id uint32) uint64
	GetTotalPayCnt(id uint32) uint32
	Refund(orderId, id string, amount int)
	Recover(orderId, id string, amount int)
}

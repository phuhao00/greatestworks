package recharge

type Order struct {
	OrderID       string        `json:"order_id"`
	PlayerID      uint64        `json:"player_id"`
	PayStatus     PayStatus     `json:"pay_status"`
	OrderSt       OrderStatus   `json:"order_st"`
	ProductID     uint32        `json:"product_id"`
	MoneyCategory MoneyCategory `json:"money_category"`
	MoneyNum      uint32        `json:"money_num"`
	Channel       string        `json:"channel"`
	ThirdOrderID  string        `json:"third_order_id"`
	CreateTM      string        `json:"create_tm"`
	CloseTM       string        `json:"close_tm"`
	PayTime       string        `json:"pay_time"`
	Price         uint32        `json:"price"`
	ReviseTM      string        `json:"revise_tm"`
	Revised       bool          `json:"revised"`
	IsCoupon      bool          `json:"is_coupon"`
}

package pay

type Recharge struct {
	Id        uint32 `json:"id"`
	Desc      string `json:"desc"`
	MoneyType int    `json:"moneyType"`
	Count     int    `json:"count"`
}

func init() {
	var RechargeRecord = Recharge{}
	_ = RechargeRecord
}

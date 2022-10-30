package recharge

const (
	PayOrderTimeOut = 7 * 86400 // 7天订单未支付则超时
)

type PayStatus uint32

const (
	PayStatusNone    PayStatus = 0 // 空
	PayStatusWait    PayStatus = 1 // 等待支付
	PayStatusSuccess PayStatus = 2 // 支付验证成功
)

type OrderStatus uint32

const (
	OrderStatusNone      OrderStatus = 0 // 空
	OrderStatusOpen      OrderStatus = 1 // 开启
	OrderStatusSuccess   OrderStatus = 2 // 领取奖励成功
	OrderStatusTimeOut   OrderStatus = 3 // 超时订单关闭
	OrderStatusError     OrderStatus = 4 // 订单出错
	OrderStatusRefunding OrderStatus = 5 // 恢复订单 没发货
	OrderStatusRefunded  OrderStatus = 6 // 恢复订单 已发货
)

type MoneyCategory uint32

const (
	MoneyCategoryCNY = 1
	MoneyCategoryUSD = 2
)

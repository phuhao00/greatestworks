package card

import (
	"errors"
	"github.com/phuhao00/greatestworks-proto/purchase"
	"greatestworks/aop/fn"
	"time"
)

type Data struct {
	Cards map[Category]Card //
}

type Category uint8

const (
	NotDefine Category = 0
	Weekly    Category = 1
	Monthly   Category = 2
)

func (c Category) GetAddExpireTime() int64 {
	if Weekly == c {
		return 7 * 24 * 60 * 60
	}
	if Monthly == c {
		return 30 * 24 * 60 * 60
	}
	return 0
}

type Card struct {
	CanReceivedTimes int
	BuyTime          int64 //
	ExpireTime       int64 //做续费判断
	LastReceivedTime int64 //领取时间
	Category
}

// - 每天可以领取奖励
//  - 几天没有领取
//  - 续费 （快到期了2天内可以续费）
//  - 购买

func (d *Data) Execute(category Category, action purchase.CardAction) error {
	card := d.Cards[category]
	switch action {
	case purchase.CardAction_DailyReceived:
		card.DailyReceive()
	case purchase.CardAction_Renew:
		card.Renew()
	case purchase.CardAction_Buy:
		return card.Buy()
	default:
		return errors.New("action not exist")
	}
	return nil
}

func (c *Card) checkIsSameDay() bool {
	return fn.IsSameDay(c.LastReceivedTime, time.Now().Unix())
}

func (c *Card) checkIsExpired() bool {
	return false
}

func (c *Card) DailyReceive() error {
	//todo check
	if c.checkIsSameDay() && c.checkIsExpired() {

	}
	//
	c.LastReceivedTime = time.Now().Unix()

	//todo
	return nil
}

func (c *Card) Renew() error {
	///todo 是否可以续费的检查
	c.BuyTime = time.Now().Unix()
	c.ExpireTime = c.ExpireTime + c.Category.GetAddExpireTime()
	return nil

}

func (c *Card) checkCanRenew() bool {
	if c.CanReceivedTimes != 0 {
		return false
	}
	return true
}

func (c *Card) Buy() error {
	///todo 是否可以购买的检查
	c.BuyTime = time.Now().Unix()
	c.ExpireTime = time.Now().Unix() + c.Category.GetAddExpireTime()
	return nil
}

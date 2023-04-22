package card

import (
	"errors"
	"fmt"
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
		return card.DailyReceive()
	case purchase.CardAction_Renew:
		return card.Renew()
	case purchase.CardAction_Buy:
		return card.buy()
	default:
		return errors.New("action not exist")
	}
}

func (c *Card) checkIsSameDay() bool {
	return fn.IsSameDay(c.LastReceivedTime, time.Now().Unix())
}

// checkIsExpired check card is expired
func (c *Card) checkIsExpired() bool {
	return time.Now().Unix() > c.ExpireTime
}

// DailyReceive  daily received reward
func (c *Card) DailyReceive() error {
	hadReceived := fn.IsSameDay(c.LastReceivedTime, time.Now().Unix())
	if hadReceived || c.checkIsExpired() {
		return errors.New("today had received")
	}
	c.LastReceivedTime = time.Now().Unix()

	//give daily  reward
	cardConf := getCardConf(c.Category)
	fmt.Println(cardConf.DailyReward)
	return nil
}

func (c *Card) Renew() error {
	if !c.checkCanRenew() {
		return errors.New("can not renew ")
	}
	c.BuyTime = time.Now().Unix()
	c.ExpireTime = c.ExpireTime + c.Category.GetAddExpireTime()
	return nil

}

// checkCanRenew check can renew
func (c *Card) checkCanRenew() bool {
	if c.CanReceivedTimes != 0 {
		return false
	}
	cardConf := getCardConf(c.Category)
	RenewInterval := int64(cardConf.RenewInterval * 24 * 60 * 60)
	if c.ExpireTime-fn.GetTimeStampDay0Time(time.Now().Unix()) > RenewInterval {
		return false
	}
	return true
}

// buy ...
func (c *Card) buy() error {
	///todo 是否可以购买的检查
	c.BuyTime = time.Now().Unix()
	c.ExpireTime = time.Now().Unix() + c.Category.GetAddExpireTime()
	return nil
}

// checkCanBuy  check can buy
func (c *Card) checkCanBuy() bool {
	if c.CanReceivedTimes != 0 || time.Now().Unix() < c.ExpireTime {
		return false
	}
	return true
}

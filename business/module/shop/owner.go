package shop

type Owner interface {
	GetShopInfo()
	BuyItem()
	RefreshShop()
	GetItemInfo()
	GetBuyRecords()
}

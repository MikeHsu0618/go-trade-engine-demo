package constants

type PriceType int
type OrderSide int

const (
	PriceTypeLimit          PriceType = 0
	PriceTypeMarket         PriceType = 1
	PriceTypeMarketQuantity PriceType = 2
	PriceTypeMarketAmount   PriceType = 3

	OrderSideBuy  OrderSide = 0
	OrderSideSell OrderSide = 1
)

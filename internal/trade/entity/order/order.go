package order

import (
	"github.com/shopspring/decimal"
	"go_trade_engine_demo/internal/trade/constants"
)

type Order struct {
	OrderId    string
	Price      decimal.Decimal
	Quantity   decimal.Decimal
	CreateTime int64
	index      int

	PriceType constants.PriceType
	Amount    decimal.Decimal
}

type CreateOrderRequest struct {
	OrderId   string `json:"order_id"`
	OrderType string `json:"order_type" binding:"required"`
	PriceType string `json:"price_type" binding:"required"`
	Price     string `json:"price" binding:"required"`
	Quantity  string `json:"quantity" binding:"required"`
	Amount    string `json:"amount"`
}

type DeleteOrderRequest struct {
	OrderId string `json:"order_id" binding:"required"`
}

func (o *Order) GetIndex() int {
	return o.index
}

func (o *Order) SetIndex(index int) {
	o.index = index
}

func (o *Order) SetQuantity(qnt decimal.Decimal) {
	o.Quantity = qnt
}

func (o *Order) SetAmount(amount decimal.Decimal) {
	o.Amount = amount
}

func (o *Order) GetUniqueId() string {
	return o.OrderId
}

func (o *Order) GetPrice() decimal.Decimal {
	return o.Price
}

func (o *Order) GetQuantity() decimal.Decimal {
	return o.Quantity
}

func (o *Order) GetCreateTime() int64 {
	return o.CreateTime
}

func (o *Order) GetPriceType() constants.PriceType {
	return o.PriceType
}
func (o *Order) GetAmount() decimal.Decimal {
	return o.Amount
}

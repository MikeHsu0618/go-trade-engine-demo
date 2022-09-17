package repository

import (
	"github.com/shopspring/decimal"
	"go_trade_engine_demo/internal/trade/constants"
	"go_trade_engine_demo/internal/trade/entity/order"
)

type AskRepository interface {
	CreateAskItem(pt constants.PriceType, uniqId string, price, quantity, amount decimal.Decimal, createTime int64) *order.AskItem
	CreateAskLimitItem(uniq string, price, quantity decimal.Decimal, createTime int64) *order.AskItem
	CreateAskMarketQtyItem(uniq string, quantity decimal.Decimal, createTime int64) *order.AskItem
	CreateAskMarketAmountItem(uniq string, amount, maxHoldQty decimal.Decimal, createTime int64) *order.AskItem
}

type askRepository struct{}

func NewAskRepository() AskRepository {
	return &askRepository{}
}

func (a *askRepository) CreateAskItem(pt constants.PriceType, uniqId string, price, quantity, amount decimal.Decimal, createTime int64) *order.AskItem {
	return &order.AskItem{
		Order: order.Order{
			OrderId:    uniqId,
			Price:      price,
			Quantity:   quantity,
			CreateTime: createTime,
			PriceType:  pt,
			Amount:     amount,
		},
	}
}

func (a *askRepository) CreateAskLimitItem(uniq string, price, quantity decimal.Decimal, createTime int64) *order.AskItem {
	return a.CreateAskItem(constants.PriceTypeLimit, uniq, price, quantity, decimal.Zero, createTime)
}

func (a *askRepository) CreateAskMarketQtyItem(uniq string, quantity decimal.Decimal, createTime int64) *order.AskItem {
	return a.CreateAskItem(constants.PriceTypeMarketQuantity, uniq, decimal.Zero, quantity, decimal.Zero, createTime)
}

//市价 按金额卖出订单时，需要用户持有交易物的数量，在撮合时候防止超卖
func (a *askRepository) CreateAskMarketAmountItem(uniq string, amount, maxHoldQty decimal.Decimal, createTime int64) *order.AskItem {
	return a.CreateAskItem(constants.PriceTypeMarketAmount, uniq, decimal.Zero, maxHoldQty, amount, createTime)
}

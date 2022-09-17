package repository

import (
	"github.com/shopspring/decimal"
	"go_trade_engine_demo/internal/trade/constants"
	"go_trade_engine_demo/internal/trade/entity/order"
)

type BidRepository interface {
	CreateBidItem(pt constants.PriceType, uniqId string, price, quantity, amount decimal.Decimal, createTime int64) *order.BidItem
	CreateBidMarketQtyItem(uniq string, quantity, maxAmount decimal.Decimal, createTime int64) *order.BidItem
	CreateBidMarketAmountItem(uniq string, amount decimal.Decimal, createTime int64) *order.BidItem
}

type bidRepository struct{}

func NewBidRepository() BidRepository {
	return &bidRepository{}
}

func (b *bidRepository) CreateBidItem(pt constants.PriceType, uniqId string, price, quantity, amount decimal.Decimal, createTime int64) *order.BidItem {
	return &order.BidItem{
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

func (b *bidRepository) CreateBidLimitItem(uniq string, price, quantity decimal.Decimal, createTime int64) *order.BidItem {
	return b.CreateBidItem(constants.PriceTypeLimit, uniq, price, quantity, decimal.Zero, createTime)
}

func (b *bidRepository) CreateBidMarketQtyItem(uniq string, quantity, maxAmount decimal.Decimal, createTime int64) *order.BidItem {
	return b.CreateBidItem(constants.PriceTypeMarketQuantity, uniq, decimal.Zero, quantity, maxAmount, createTime)
}

//市价 按金额卖出订单时，需要用户持有交易物的数量，在撮合时候防止超卖
func (b *bidRepository) CreateBidMarketAmountItem(uniq string, amount decimal.Decimal, createTime int64) *order.BidItem {
	return b.CreateBidItem(constants.PriceTypeMarketAmount, uniq, decimal.Zero, decimal.Zero, amount, createTime)
}

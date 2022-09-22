package order

import (
	"go_trade_engine_demo/internal/trade/constants"
	"go_trade_engine_demo/internal/trade/pkg/queue"
)

type BidItem struct {
	Order
}

func (b *BidItem) Less(other queue.QueueItem) bool {
	// 價格優先，時間優先原則
	// 價格高的在最上面
	return (b.Price.Cmp(other.(*BidItem).Price) == 1) || (b.Price.Cmp(other.(*BidItem).Price) == 0 && b.CreateTime < other.(*BidItem).CreateTime)
}

func (b *BidItem) GetOrderSide() constants.OrderSide {
	return constants.OrderSideBuy
}

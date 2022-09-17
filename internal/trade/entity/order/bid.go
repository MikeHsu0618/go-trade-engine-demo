package order

import (
	"go_trade_engine_demo/internal/trade/constants"
	"go_trade_engine_demo/internal/trade/entity/queue"
)

type BidItem struct {
	Order
}

func (a *BidItem) Less(other queue.QueueItem) bool {
	//价格优先，时间优先原则
	//价格高的在最上面
	return (a.Price.Cmp(other.(*BidItem).Price) == 1) || (a.Price.Cmp(other.(*BidItem).Price) == 0 && a.CreateTime < other.(*BidItem).CreateTime)
}

func (a *BidItem) GetOrderSide() constants.OrderSide {
	return constants.OrderSideBuy
}

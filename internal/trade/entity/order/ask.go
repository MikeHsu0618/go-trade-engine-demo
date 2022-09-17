package order

import (
	"go_trade_engine_demo/internal/trade/constants"
	"go_trade_engine_demo/internal/trade/entity/queue"
)

type AskItem struct {
	Order
}

func (a *AskItem) Less(other queue.QueueItem) bool {
	//价格优先，时间优先原则
	//价格低的在最上面
	return (a.Price.Cmp(other.(*AskItem).Price) == -1) || (a.Price.Cmp(other.(*AskItem).Price) == 0 && a.CreateTime < other.(*AskItem).CreateTime)
}

func (a *AskItem) GetOrderSide() constants.OrderSide {
	return constants.OrderSideSell
}

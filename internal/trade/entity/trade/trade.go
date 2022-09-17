package trade

import (
	"github.com/shopspring/decimal"
	"go_trade_engine_demo/internal/trade/entity/queue"
	"go_trade_engine_demo/internal/trade/util"
	"sync"
)

type Result struct {
	Symbol        string          `json:"symbol"`
	AskOrderId    string          `json:"ask_order_id"`
	BidOrderId    string          `json:"bid_order_id"`
	TradeQuantity decimal.Decimal `json:"trade_quantity"`
	TradePrice    decimal.Decimal `json:"trade_price"`
	TradeAmount   decimal.Decimal `json:"trade_amount"`
	TradeTime     int64           `json:"trade_time"`
}

type Pair struct {
	Symbol           string
	TradeResultChan  chan Result
	NewOrderChan     chan queue.QueueItem
	CancelResultChan chan string
	RecentTrade      []interface{}

	PriceDigit    int
	QuantityDigit int
	LatestPrice   decimal.Decimal

	AskQueue *queue.OrderQueue
	BidQueue *queue.OrderQueue

	sync.Mutex
}

func (t *Pair) Price2String(price decimal.Decimal) string {
	return util.FormatDecimal2String(price, t.PriceDigit)
}

func (t *Pair) Qty2String(qty decimal.Decimal) string {
	return util.FormatDecimal2String(qty, t.QuantityDigit)
}

func (t *Pair) Depth(queue *queue.OrderQueue, size int) [][2]string {
	queue.Lock()
	defer queue.Unlock()

	max := len(queue.Depth)
	if size <= 0 || size > max {
		size = max
	}

	return queue.Depth[0:size]
}

func (t *Pair) AskLen() int {
	t.Lock()
	defer t.Unlock()

	return t.AskQueue.Len()
}

func (t *Pair) BidLen() int {
	t.Lock()
	defer t.Unlock()

	return t.BidQueue.Len()
}

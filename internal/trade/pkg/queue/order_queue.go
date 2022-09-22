package queue

import (
	"container/heap"
	"github.com/shopspring/decimal"
	"go_trade_engine_demo/internal/trade/constants"
	"sync"
)

type QueueItem interface {
	SetIndex(index int)
	SetQuantity(quantity decimal.Decimal)
	SetAmount(amount decimal.Decimal)
	Less(item QueueItem) bool
	GetIndex() int
	GetUniqueId() string
	GetPrice() decimal.Decimal
	GetQuantity() decimal.Decimal
	GetCreateTime() int64
	GetOrderSide() constants.OrderSide
	GetPriceType() constants.PriceType
	GetAmount() decimal.Decimal // 訂單金額，在市價單下單才用的到，限價單不需要
}

func CreateQueue() *OrderQueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	queue := OrderQueue{
		Pq: &pq,
		m:  make(map[string]*QueueItem),
	}
	return &queue
}

type OrderQueue struct {
	Pq *PriorityQueue
	m  map[string]*QueueItem
	sync.Mutex

	Depth [][2]string
}

func (o *OrderQueue) Len() int {
	return o.Pq.Len()
}

func (o *OrderQueue) Push(item QueueItem) (exist bool) {
	o.Lock()
	defer o.Unlock()

	id := item.GetUniqueId()
	if _, ok := o.m[id]; ok {
		return true
	}

	heap.Push(o.Pq, item)
	o.m[id] = &item
	return false
}

func (o *OrderQueue) Get(index int) QueueItem {
	n := o.Pq.Len()
	if n <= index {
		return nil
	}

	return (*o.Pq)[index]
}

func (o *OrderQueue) Top() QueueItem {
	return o.Get(0)
}

func (o *OrderQueue) Remove(uniqId string) QueueItem {
	o.Lock()
	defer o.Unlock()

	old, ok := o.m[uniqId]
	if !ok {
		return nil
	}

	item := heap.Remove(o.Pq, (*old).GetIndex())
	delete(o.m, uniqId)
	return item.(QueueItem)
}

func (o *OrderQueue) clean() {
	o.Lock()
	defer o.Unlock()

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	o.Pq = &pq
	o.m = make(map[string]*QueueItem)
}

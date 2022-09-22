package repository

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go_trade_engine_demo/internal/trade/constants"
	"go_trade_engine_demo/internal/trade/entity/trade"
	"go_trade_engine_demo/internal/trade/pkg/log"
	"go_trade_engine_demo/internal/trade/pkg/queue"
	"go_trade_engine_demo/internal/trade/pkg/wss"
	"go_trade_engine_demo/internal/trade/util"
	"time"
)

type TradeRepository interface {
	GetPair() *trade.Pair
	GetAskDepth(size int) [][2]string
	GetBidDepth(size int) [][2]string
	SendMessage(tag string, data interface{})
	DeleteOrder(side constants.OrderSide, uniq string)
}

type tradeRepository struct {
	pair   *trade.Pair
	wssHub *wss.Hub
	logger *log.Logger
}

var Debug = true

func NewTradeRepository(pair *trade.Pair, wssHub *wss.Hub, logger *log.Logger) TradeRepository {
	repo := &tradeRepository{
		pair:   pair,
		wssHub: wssHub,
		logger: logger,
	}
	go repo.pushDepth()
	go repo.depthTicker(repo.pair.AskQueue)
	go repo.depthTicker(repo.pair.BidQueue)
	go repo.matching()
	go repo.watchTradeLog()
	return repo
}

func (r *tradeRepository) GetPair() *trade.Pair {
	return r.pair
}

func (r *tradeRepository) GetAskDepth(size int) [][2]string {
	return r.pair.Depth(r.pair.AskQueue, size)
}

func (r *tradeRepository) GetBidDepth(size int) [][2]string {
	return r.pair.Depth(r.pair.BidQueue, size)
}

func (r *tradeRepository) matching() {
	for {
		select {
		case newOrder := <-r.pair.NewOrderChan:
			go r.handlerNewOrder(newOrder)
		default:
			r.handlerLimitOrder()
		}
	}
}

func (r *tradeRepository) handlerNewOrder(newOrder queue.QueueItem) {
	r.pair.Lock()
	defer r.pair.Unlock()

	if newOrder.GetPriceType() == constants.PriceTypeLimit {
		if newOrder.GetOrderSide() == constants.OrderSideSell {
			r.pair.AskQueue.Push(newOrder)
		} else {
			r.pair.BidQueue.Push(newOrder)
		}
	} else {
		// 市價單處理
		if newOrder.GetOrderSide() == constants.OrderSideSell {
			r.doMarketSell(newOrder)
		} else {
			r.doMarketBuy(newOrder)
		}
	}
}

func (r *tradeRepository) handlerLimitOrder() {
	ok := func() bool {
		r.pair.Lock()
		defer r.pair.Unlock()

		if r.pair.AskQueue == nil || r.pair.BidQueue == nil {
			return false
		}

		if r.pair.AskQueue.Len() == 0 || r.pair.BidQueue.Len() == 0 {
			return false
		}

		askTop := r.pair.AskQueue.Top()
		bidTop := r.pair.BidQueue.Top()

		defer func() {
			if askTop.GetQuantity().Equal(decimal.Zero) {
				r.pair.AskQueue.Remove(askTop.GetUniqueId())
			}
			if bidTop.GetQuantity().Equal(decimal.Zero) {
				r.pair.BidQueue.Remove(bidTop.GetUniqueId())
			}
		}()

		if bidTop.GetPrice().Cmp(askTop.GetPrice()) >= 0 {
			curTradeQty := decimal.Zero
			curTradePrice := decimal.Zero
			if bidTop.GetQuantity().Cmp(askTop.GetQuantity()) >= 0 {
				curTradeQty = askTop.GetQuantity()
			} else if bidTop.GetQuantity().Cmp(askTop.GetQuantity()) == -1 {
				curTradeQty = bidTop.GetQuantity()
			}
			askTop.SetQuantity(askTop.GetQuantity().Sub(curTradeQty))
			bidTop.SetQuantity(bidTop.GetQuantity().Sub(curTradeQty))

			if askTop.GetCreateTime() >= bidTop.GetCreateTime() {
				curTradePrice = bidTop.GetPrice()
			} else {
				curTradePrice = askTop.GetPrice()
			}

			r.sendTradeResultNotify(askTop, bidTop, curTradePrice, curTradeQty)
			return true
		} else {
			return false
		}

	}()

	if !ok {
		time.Sleep(time.Duration(200) * time.Millisecond)
	}
}

func (r *tradeRepository) doMarketBuy(item queue.QueueItem) {

	for {
		ok := func() bool {

			if r.pair.AskQueue.Len() == 0 {
				return false
			}

			ask := r.pair.AskQueue.Top()
			if item.GetPriceType() == constants.PriceTypeMarketQuantity {
				// TODO 根據用戶資產計算出當前價格能買的最大數量
				// maxTradeQty := item.GetAmount().Div(ask.GetPrice())        // 當前資產可買進的最大數目
				// maxTradeQty = decimal.Min(maxTradeQty, item.GetQuantity()) // 比較是否超過下單數目
				maxTradeQty := item.GetQuantity()
				curTradeQty := decimal.Zero

				// 市價按買入數量
				if maxTradeQty.Cmp(decimal.New(1, int32(-r.pair.QuantityDigit))) < 0 {
					return false
				}

				if ask.GetQuantity().Cmp(maxTradeQty) <= 0 {
					curTradeQty = ask.GetQuantity()
					defer r.pair.AskQueue.Remove(ask.GetUniqueId())
				} else {
					curTradeQty = maxTradeQty
					ask.SetQuantity(ask.GetQuantity().Sub(curTradeQty))
				}

				r.sendTradeResultNotify(ask, item, ask.GetPrice(), curTradeQty)
				item.SetQuantity(item.GetQuantity().Sub(curTradeQty))
				item.SetAmount(item.GetAmount().Sub(curTradeQty.Mul(ask.GetPrice())))
				return true
			} else if item.GetPriceType() == constants.PriceTypeMarketAmount {
				//  TODO 市價-按成交金額
				// 成交金額不包含手續費，手續費應該由上層系統計算提前預留
				// 撮合會針對輸入金額最大限度買入
				if ask.GetPrice().Cmp(decimal.Zero) <= 0 {
					return false
				}

				maxTradeQty := item.GetAmount().Div(ask.GetPrice())
				curTradeQty := decimal.Zero

				if maxTradeQty.Cmp(decimal.New(1, int32(-r.pair.QuantityDigit))) < 0 {
					return false
				}
				if ask.GetQuantity().Cmp(maxTradeQty) <= 0 {
					curTradeQty = ask.GetQuantity()
					defer r.pair.AskQueue.Remove(ask.GetUniqueId())
				} else {
					curTradeQty = maxTradeQty
					ask.SetQuantity(ask.GetQuantity().Sub(curTradeQty))
				}

				r.sendTradeResultNotify(ask, item, ask.GetPrice(), curTradeQty)
				// 部分成交了，需要更新此單剩餘可成交金額，用於下一輪重新计算最大成交量
				item.SetAmount(item.GetAmount().Sub(curTradeQty.Mul(ask.GetPrice())))
				item.SetQuantity(item.GetQuantity().Add(curTradeQty))
				return true
			}

			return false
		}()

		if !ok {
			// 市價單不管是否完全成交，都觸發一次撤單操作
			r.pair.CancelResultChan <- item.GetUniqueId()
			break
		}

	}
}
func (r *tradeRepository) doMarketSell(item queue.QueueItem) {

	for {
		ok := func() bool {

			if r.pair.BidQueue.Len() == 0 {
				return false
			}

			bid := r.pair.BidQueue.Top()
			if item.GetPriceType() == constants.PriceTypeMarketQuantity {

				curTradeQuantity := decimal.Zero
				// 市價單按賣出數量
				// TODO 按持有數量限制買入
				if item.GetQuantity().Equal(decimal.Zero) {
					return false
				}

				if bid.GetQuantity().Cmp(item.GetQuantity()) <= 0 {
					curTradeQuantity = bid.GetQuantity()
					defer r.pair.BidQueue.Remove(bid.GetUniqueId())
				} else {
					curTradeQuantity = item.GetQuantity()
					bid.SetQuantity(bid.GetQuantity().Sub(curTradeQuantity))
				}

				r.sendTradeResultNotify(item, bid, bid.GetPrice(), curTradeQuantity)
				item.SetQuantity(item.GetQuantity().Sub(curTradeQuantity))

				return true
			} else if item.GetPriceType() == constants.PriceTypeMarketAmount {
				// 市價-按成交金額成交
				if bid.GetPrice().Cmp(decimal.Zero) <= 0 {
					return false
				}

				// TODO 還需要用戶是否持有那麼多資產來賣出的條件限制
				//maxTradeQty := item.GetAmount().Div(bid.GetPrice()).Truncate(int32(r.pair.QuantityDigit))
				//maxTradeQty = decimal.Min(maxTradeQty, item.GetQuantity())

				maxTradeQty := item.GetQuantity()

				curTradeQty := decimal.Zero
				if maxTradeQty.Cmp(decimal.New(1, int32(-r.pair.QuantityDigit))) < 0 {
					return false
				}

				if bid.GetQuantity().Cmp(maxTradeQty) <= 0 {
					curTradeQty = bid.GetQuantity()
					defer r.pair.BidQueue.Remove(bid.GetUniqueId())
				} else {
					curTradeQty = maxTradeQty
					bid.SetQuantity(bid.GetQuantity().Sub(curTradeQty))
				}

				r.sendTradeResultNotify(item, bid, bid.GetPrice(), curTradeQty)
				item.SetAmount(item.GetAmount().Sub(curTradeQty.Mul(bid.GetPrice())))

				// 市價單按成交額賣出時，需要用户持有的資產数量来進行限制
				item.SetQuantity(item.GetQuantity().Sub(curTradeQty))

				return true
			}

			return false
		}()

		if !ok {
			r.pair.CancelResultChan <- item.GetUniqueId()
			break
		}

	}
}

func (r *tradeRepository) sendTradeResultNotify(ask, bid queue.QueueItem, price, tradeQty decimal.Decimal) {
	tradelog := trade.Result{}
	tradelog.Symbol = r.pair.Symbol
	tradelog.AskOrderId = ask.GetUniqueId()
	tradelog.BidOrderId = bid.GetUniqueId()
	tradelog.TradeQuantity = tradeQty
	tradelog.TradePrice = price
	tradelog.TradeTime = time.Now().UnixNano()
	tradelog.TradeAmount = tradeQty.Mul(price)

	r.pair.LatestPrice = price

	if Debug {
		r.logger.Info(fmt.Sprintf("%s tradelog: %+v", r.pair.Symbol, tradelog))
	}

	r.pair.TradeResultChan <- tradelog
}

func (r *tradeRepository) depthTicker(que *queue.OrderQueue) {

	ticker := time.NewTicker(time.Duration(50) * time.Millisecond)

	for {
		<-ticker.C
		func() {
			r.pair.Lock()
			defer r.pair.Unlock()

			que.Lock()
			defer que.Unlock()
			que.Depth = [][2]string{}
			depthMap := make(map[string]string)

			if que.Pq.Len() > 0 {

				for i := 0; i < que.Pq.Len(); i++ {
					item := (*que.Pq)[i]

					price := util.FormatDecimal2String(item.GetPrice(), r.pair.PriceDigit)

					if _, ok := depthMap[price]; !ok {
						depthMap[price] = util.FormatDecimal2String(item.GetQuantity(), r.pair.QuantityDigit)
					} else {
						oldQuantity, _ := decimal.NewFromString(depthMap[price])
						depthMap[price] = util.FormatDecimal2String(oldQuantity.Add(item.GetQuantity()), r.pair.QuantityDigit)
					}
				}

				// 按價格排序的 MAP
				que.Depth = util.SortMap2Slice(depthMap, que.Top().GetOrderSide())
			}
		}()
	}
}

func (r *tradeRepository) pushDepth() {
	for {
		ask := r.GetAskDepth(10)
		bid := r.GetBidDepth(10)

		r.SendMessage("depth", gin.H{
			"ask": ask,
			"bid": bid,
		})

		time.Sleep(time.Duration(150) * time.Millisecond)
	}
}

func (r *tradeRepository) SendMessage(tag string, data interface{}) {
	msg := gin.H{
		"tag":  tag,
		"data": data,
	}
	msgByte, _ := json.Marshal(msg)
	r.wssHub.Send(msgByte)
}

func (r *tradeRepository) watchTradeLog() {
	for {
		select {
		case log, ok := <-r.pair.TradeResultChan:
			if ok {
				relog := gin.H{
					"trade_price":    r.pair.Price2String(log.TradePrice),
					"trade_amount":   r.pair.Price2String(log.TradeAmount),
					"trade_quantity": r.pair.Qty2String(log.TradeQuantity),
					"trade_time":     log.TradeTime,
					"ask_order_id":   log.AskOrderId,
					"bid_order_id":   log.BidOrderId,
				}
				r.SendMessage("trade", relog)

				if len(r.pair.RecentTrade) >= 10 {
					r.pair.RecentTrade = r.pair.RecentTrade[1:]
				}
				r.pair.RecentTrade = append(r.pair.RecentTrade, relog)

				// latest price
				r.SendMessage("latest_price", gin.H{
					"latest_price": r.pair.Price2String(log.TradePrice),
				})

			}
		case cancelOrderId := <-r.pair.CancelResultChan:
			r.SendMessage("cancel_order", gin.H{
				"OrderId": cancelOrderId,
			})
		default:
			time.Sleep(time.Duration(100) * time.Millisecond)
		}

	}
}

func (r *tradeRepository) DeleteOrder(side constants.OrderSide, uniq string) {
	// TODO 最好根據訂單編號知道是買單還是賣單，方便直接查找相應的對列並從中刪除
	if side == constants.OrderSideSell {
		r.pair.AskQueue.Remove(uniq)
	} else {
		r.pair.BidQueue.Remove(uniq)
	}
	// 删除成功后需要發送通知
	r.pair.CancelResultChan <- uniq
}

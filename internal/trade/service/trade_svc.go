package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go_trade_engine_demo/internal/trade/constants"
	"go_trade_engine_demo/internal/trade/entity/order"
	"go_trade_engine_demo/internal/trade/pkg/httputil"
	"go_trade_engine_demo/internal/trade/repository"
	"go_trade_engine_demo/internal/trade/util"
	"strconv"
	"strings"
	"time"
)

type TradeService interface {
	GetDepth(c *gin.Context)
	GetTradeLog(c *gin.Context)
	CreateOrder(c *gin.Context)
	DeleteOrder(c *gin.Context)
	CreateRandomOrder(c *gin.Context)
}

type tradeService struct {
	tradeRepo repository.TradeRepository
	askRepo   repository.AskRepository
	bidRepo   repository.BidRepository
}

func NewTradeService(tradeRepo repository.TradeRepository, askRepo repository.AskRepository, bidRepo repository.BidRepository) TradeService {
	return &tradeService{
		tradeRepo: tradeRepo,
		askRepo:   askRepo,
		bidRepo:   bidRepo,
	}
}

func (s *tradeService) GetDepth(c *gin.Context) {
	limit := c.Query("limit")
	limitInt, _ := strconv.Atoi(limit)
	if limitInt <= 0 || limitInt > 100 {
		limitInt = 10
	}
	a := s.tradeRepo.GetAskDepth(limitInt)
	b := s.tradeRepo.GetBidDepth(limitInt)

	httputil.NewSuccess(c, gin.H{
		"ask": a,
		"bid": b,
	})
}

func (s *tradeService) CreateOrder(c *gin.Context) {
	var param order.CreateOrderRequest
	err := c.ShouldBind(&param)
	if err != nil {
		httputil.NewError(c, 404, fmt.Sprintf("Validation Error: %v", err))
		return
	}
	tradePair := s.tradeRepo.GetPair()
	orderId := uuid.NewString()
	param.OrderId = orderId

	amount := util.String2decimal(param.Amount)
	price := util.String2decimal(param.Price)
	quantity := util.String2decimal(param.Quantity)

	var pt constants.PriceType
	if param.PriceType == "market" {
		param.Price = "0"
		pt = constants.PriceTypeMarket
		if param.Amount != "" {
			pt = constants.PriceTypeMarketAmount
			//市价按成交金额卖出时，默认持有该资产1000个
			param.Quantity = "100"
			if amount.Cmp(decimal.NewFromFloat(100000000)) > 0 || amount.Cmp(decimal.Zero) <= 0 {
				httputil.NewError(c, 429, "金額必須大於零，且不能超過 100000000")
				return
			}

		} else if param.Quantity != "" {
			pt = constants.PriceTypeMarketQuantity
			// TODO 市價按數量買入資產時，需要用戶帳戶所有可用資產數量，默認 100
			param.Amount = "100"
			if quantity.Cmp(decimal.NewFromFloat(100000000)) > 0 || quantity.Cmp(decimal.Zero) <= 0 {
				httputil.NewError(c, 429, "數量必須大於零，且不能超過 100000000")
				return
			}
		}
	} else {
		pt = constants.PriceTypeLimit
		param.Amount = "0"
		if price.Cmp(decimal.NewFromFloat(100000000)) > 0 || price.Cmp(decimal.Zero) < 0 {
			httputil.NewError(c, 429, "價格必須大於等於零，且不能超過 100000000")
			return
		}
		if quantity.Cmp(decimal.NewFromFloat(100000000)) > 0 || quantity.Cmp(decimal.Zero) <= 0 {
			httputil.NewError(c, 429, "數量必須大於零，且不能超過 100000000")
			return
		}
	}

	if strings.ToLower(param.OrderType) == "ask" {
		param.OrderId = fmt.Sprintf("a-%s", orderId)
		item := s.askRepo.CreateAskItem(pt, param.OrderId, util.String2decimal(param.Price), util.String2decimal(param.Quantity), util.String2decimal(param.Amount), time.Now().UnixNano())
		tradePair.NewOrderChan <- item

	} else {
		param.OrderId = fmt.Sprintf("b-%s", orderId)
		item := s.bidRepo.CreateBidItem(pt, param.OrderId, util.String2decimal(param.Price), util.String2decimal(param.Quantity), util.String2decimal(param.Amount), time.Now().UnixNano())
		tradePair.NewOrderChan <- item
	}

	go s.tradeRepo.SendMessage("new_order", param)
	httputil.NewSuccess(c, gin.H{
		"ask_len": tradePair.AskLen(),
		"bid_len": tradePair.BidLen(),
	})
}

func (s *tradeService) DeleteOrder(c *gin.Context) {
	var param order.DeleteOrderRequest
	err := c.ShouldBindJSON(&param)
	if err != nil {
		httputil.NewError(c, 404, "Validation Error")
		return
	}

	if strings.HasPrefix(param.OrderId, "a-") {
		s.tradeRepo.DeleteOrder(constants.OrderSideSell, param.OrderId)
	} else {
		s.tradeRepo.DeleteOrder(constants.OrderSideBuy, param.OrderId)
	}

	go s.tradeRepo.SendMessage("cancel_order", param)

	httputil.NewSuccess(c, "success")
}

func (s *tradeService) GetTradeLog(c *gin.Context) {
	pair := s.tradeRepo.GetPair()
	httputil.NewSuccess(c, gin.H{
		"trade_log":    pair.RecentTrade,
		"latest_price": pair.LatestPrice,
	})
}

func (s *tradeService) CreateRandomOrder(c *gin.Context) {
	type params struct {
		OrderType string `json:"order_type"`
	}
	var input params
	err := c.ShouldBindJSON(&input)
	if err != nil {
		httputil.NewError(c, 429, fmt.Sprintf("Validation Error: %v", err))
		return
	}
	op := strings.ToLower(input.OrderType)
	if op != "ask" {
		op = "bid"
	}
	pair := s.tradeRepo.GetPair()
	func() {
		cnt := 10
		for i := 0; i < cnt; i++ {
			orderId := uuid.NewString()
			if op == "ask" {
				orderId = fmt.Sprintf("a-%s", orderId)
				item := s.askRepo.CreateAskLimitItem(orderId, util.RandDecimal(20, 50), util.RandDecimal(20, 100), time.Now().UnixNano())
				pair.NewOrderChan <- item
			} else {
				orderId = fmt.Sprintf("b-%s", orderId)
				item := s.bidRepo.CreateBidLimitItem(orderId, util.RandDecimal(1, 20), util.RandDecimal(20, 100), time.Now().UnixNano())
				pair.NewOrderChan <- item
			}
		}
	}()
	httputil.NewSuccess(c, gin.H{
		"ask_len": pair.AskLen(),
		"bid_len": pair.BidLen(),
	})
}

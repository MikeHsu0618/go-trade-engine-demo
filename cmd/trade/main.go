package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"go_trade_engine_demo/internal/trade/entity/queue"
	"go_trade_engine_demo/internal/trade/entity/trade"
	"go_trade_engine_demo/internal/trade/handler/httpHandler"
	"go_trade_engine_demo/internal/trade/pkg/log"
	"go_trade_engine_demo/internal/trade/pkg/wss"
	"go_trade_engine_demo/internal/trade/repository"
	"go_trade_engine_demo/internal/trade/service"
)

var port = flag.String("port", ":8888", "port")
var wssHub = wss.NewHub()

func main() {
	flag.Parse()
	logger := log.New()

	askRepo := repository.NewAskRepository()
	bidRepo := repository.NewBidRepository()
	tradeRepo := repository.NewTradeRepository(&trade.Pair{
		Symbol:           "TSM",
		TradeResultChan:  make(chan trade.Result, 10),
		NewOrderChan:     make(chan queue.QueueItem),
		CancelResultChan: make(chan string, 10),
		RecentTrade:      make([]interface{}, 0),
		PriceDigit:       2,
		QuantityDigit:    4,
		AskQueue:         queue.CreateQueue(),
		BidQueue:         queue.CreateQueue(),
	}, wssHub, logger)
	tradeSvc := service.NewTradeService(tradeRepo, askRepo, bidRepo)

	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(initWss)
	httpHandler.NewHandler(&httpHandler.Config{
		Router:   r,
		TradeSvc: tradeSvc,
	})

	go wssHub.Run()
	r.Run(*port)
}

func initWss(c *gin.Context) {
	c.Set("wssHub", wssHub)
}

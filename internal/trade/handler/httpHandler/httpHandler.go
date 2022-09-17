package httpHandler

import (
	"github.com/gin-gonic/gin"
	"go_trade_engine_demo/internal/trade/pkg/wss"
	"go_trade_engine_demo/internal/trade/service"
)

type Handler struct {
	tradeSvc service.TradeService
}

type Config struct {
	Router   *gin.Engine
	TradeSvc service.TradeService
}

func NewHandler(c *Config) {
	h := &Handler{
		tradeSvc: c.TradeSvc,
	}

	v1Group := c.Router.Group("/api/v1")
	trade := v1Group.Group("trade")
	{
		trade.GET("/depth", h.tradeSvc.GetDepth)
		trade.POST("/orders", h.tradeSvc.CreateOrder)
		trade.DELETE("/orders", h.tradeSvc.DeleteOrder)
		trade.GET("/log", h.tradeSvc.GetTradeLog)
		trade.GET("/test_rand", h.success)
	}
	// wss
	c.Router.GET("/ws", wss.ServeWs)
}

func (h *Handler) success(c *gin.Context) {
	c.JSON(200, gin.H{"data": "success"})
}

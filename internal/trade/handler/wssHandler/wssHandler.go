package httpHandler

import (
	"go_trade_engine_demo/internal/trade/pkg/wss"
)

type Handler struct {
	WssHub *wss.Hub
}

type Config struct {
	WssHub *wss.Hub
}

var sendMsg = make(chan []byte, 100)

func NewHandler(c *Config) {
	h := &Handler{
		WssHub: c.WssHub,
	}
	h.WssHub.Run()
}

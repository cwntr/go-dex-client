package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cwntr/go-dex-client/lncli"
	"github.com/cwntr/go-dex-client/lssdrpc"
	"github.com/cwntr/go-dex-client/pkg/common"
	"github.com/cwntr/go-dex-client/trading"
	"github.com/gin-gonic/gin"
)

// Order struct holds instructions how to place an order.
// Side is either "sell" or "buy"
// TradingPair is either "XSN_LTC" or "XSN_BTC" (for now)
// PriceRangeStart, PriceRangeEnd will be the order placement loop stat and end parameter whereas FixedFunding will be the actual funds that will be put on the current price of iteration.
// Type enum either "single" or "repeat" or "arbitrage"
type Order struct {
	Side               string `json:"side"`
	TradingPair        string `json:"tradingPair"`
	PriceRangeStart    int    `json:"priceRangeStart"`
	PriceRangeEnd      int    `json:"priceRangeEnd"`
	PriceRangeStepSize int    `json:"priceRangeStepSize"`
	FixedFunding       int    `json:"fixedFunding"`
	Type               string `json:"type"`
	Label              string `json:"label,omitempty"`
}

// CancelTradingPair struct holds information how to cancel placed orders.
// Currency is mandatory ("XSN_LTC" or "XSN_BTC")
// (optional) DeleteAll if true, will delete all "own" orders from current orderbook
// (optional) OrderIDs will delete all orders with specified orderIDs
type CancelTradingPair struct {
	TradingPair string   `json:"tradingPair"`
	DeleteAll   bool     `json:"deleteAll,omitempty"`
	OrderIDs    []string `json:"orderIds,omitempty"`
}

// OrderCancel structs holds all the CancelTradingPair objects
type OrderCancel struct {
	TradingPair []CancelTradingPair `json:"cancelTradingPairs"`
}

// addRoutes will provide endpoints to interact with your client
func addRoutes(r *gin.Engine) {
	//orderbook
	r.GET("/orderbook/:tradingPair", APIGetOrderbook)

	//my balances
	r.GET("/balances/:coin", APIGetBalance)

	// orders
	r.GET("/orders/:tradingPair", APIGetOrders)
	r.POST("/orders", APIPostMyOrders)
	r.POST("/orders/cancel", APIPostOrdersCancel)
}

func APIGetOrderbook(c *gin.Context) {
	onlyOwn := false
	ownValue := c.Query("only-own")
	if ownValue == "" {
		b, err := strconv.ParseBool(ownValue)
		if err == nil && b {
			onlyOwn = true
		}
	}

	tp := c.Param("tradingPair")
	if common.IsValidTradingPair(strings.ToUpper(tp)) {
		logger.Infoln("trading pair invalid")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	orders, err := bot.ListOrders(strings.ToUpper(tp), true, true)
	if err != nil {
		logger.Errorf("err while listing the orderbook %v", err)
		return
	}

	if onlyOwn {
		var temp []lssdrpc.Order
		for _, o := range orders {
			if o.IsOwnOrder {
				temp = append(temp, o)
			}
		}
		orders = temp
	}
	c.JSON(http.StatusOK, orders)
	return
}

func APIPostMyOrders(c *gin.Context) {
	var orders []Order
	err := c.BindJSON(&orders)
	if err != nil {
		logger.Errorf("error binding order payload: %v", err)
		c.Error(err)
		return
	}

	for _, order := range orders {
		if order.Type == "" {
			logger.Errorln("order type missing")
			continue
		}

		if order.PriceRangeStart == 0 || order.PriceRangeEnd == 0 || order.PriceRangeStepSize == 0 || order.FixedFunding == 0 {
			logger.Errorln("price range config: cannot have any '0' value")
			continue
		}

		if order.Side != "sell" && order.Side != "buy" {
			logger.Errorln("err: order wrong side - must be either `sell` or `buy`")
			continue
		}

		tp := strings.ToUpper(order.TradingPair)
		if common.IsValidTradingPair(strings.ToUpper(tp)) {
			logger.Errorln("trading pair invalid")
			continue
		}

		//Iterate over order price configs
		for _, price := range makeRange(order.PriceRangeStart, order.PriceRangeEnd, order.PriceRangeStepSize) {
			//resolve side
			var side lssdrpc.OrderSide
			if order.Side == "sell" {
				side = lssdrpc.OrderSide_sell
			} else if order.Side == "buy" {
				side = lssdrpc.OrderSide_buy
			}

			//Place the order
			res, err := bot.PlaceOrder(order.TradingPair, int64(price), int64(order.FixedFunding), side, order.Type, order.Label)
			if err != nil {
				logger.Errorf("err while placing an order %v", err)
			} else {
				logger.Infof("Added order, outcome: %v", res.Outcome)
			}
		}
	}
	c.JSON(http.StatusOK, orders)
	return
}

func APIGetOrders(c *gin.Context) {
	tp := c.Param("tradingPair")
	tp = strings.ToUpper(tp)
	if common.IsValidTradingPair(strings.ToUpper(tp)) {
		logger.Infoln("trading pair invalid")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	cachedOrders := trading.GetMyOrdersCache()
	coll, found := cachedOrders.Get(tp)
	if !found {
		logger.Errorf("unable to get orders cache for trading-pair: %s", tp)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	cached, ok := coll.(trading.OrderCollection)
	if !ok {
		logger.Errorf("unable to convert to OrderCollection")
		c.Error(fmt.Errorf("unable to convert to OrderCollection"))
		return
	}
	c.JSON(http.StatusOK, cached)
	return
}

func APIGetBalance(c *gin.Context) {
	coin := c.Param("coin")
	if coin == "" {
		logger.Infoln("APIGetBalance - missing coin")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	var b lncli.Balance
	var err error
	switch strings.ToUpper(coin) {
	case common.CurrencyBTC:
		b, err = lncli.GetWalletBalance(cfg.Bot.LNCLIPath, cfg.BTC.Directory, cfg.BTC.Host, cfg.BTC.Port)
		if err != nil {
			logger.Errorf("err while GetWalletBalance %v for coin: %s", err, coin)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		break
	case common.CurrencyXSN:
		b, err = lncli.GetWalletBalance(cfg.Bot.LNCLIPath, cfg.XSN.Directory, cfg.XSN.Host, cfg.XSN.Port)
		if err != nil {
			logger.Errorf("err while GetWalletBalance %v for coin: %s", err, coin)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		break
	case common.CurrencyLTC:
		b, err = lncli.GetWalletBalance(cfg.Bot.LNCLIPath, cfg.LTC.Directory, cfg.LTC.Host, cfg.LTC.Port)
		if err != nil {
			logger.Errorf("err while GetWalletBalance %v for coin: %s", err, coin)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		break
	default:
		logger.Errorf("err while GetWalletBalance %v for coin: %s", err, coin)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, b)
	return
}

func APIPostOrdersCancel(c *gin.Context) {
	var orderCancel OrderCancel
	err := c.BindJSON(&orderCancel)
	if err != nil {
		logger.Errorf("error binding order cancel payload: %v", err)
		c.Error(err)
		return
	}
	for _, cancel := range orderCancel.TradingPair {
		for _, id := range cancel.OrderIDs {
			resp, err := bot.CancelOrder(strings.ToUpper(cancel.TradingPair), id)
			logger.Debugf("cancelResponse: %v", resp)
			if err != nil {
				logger.Errorf("err while canceling order: %s, err %v", id, err)
			} else {
				logger.Debugf("canceled orderId: %s", id)
			}
		}
		if cancel.DeleteAll {
			orders, err := bot.ListOrders(strings.ToUpper(cancel.TradingPair), true, false)
			if err != nil {
				logger.Errorf("err while fetching all orders from orderbook, err %v", err)
				return
			}
			for _, o := range orders {
				if o.IsOwnOrder {
					resp, err := bot.CancelOrder(strings.ToUpper(cancel.TradingPair), o.OrderId)
					logger.Debugf("cancelResponse: %v", resp)
					if err != nil {
						logger.Errorf("err while canceling order: %s, err %v", o.OrderId, err)
					} else {
						logger.Debugf("canceled orderId: %s", o.OrderId)
					}
				}
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"ok": "orders cancelled"})
	return
}

func makeRange(min int, max int, step int) []int {
	var rangeList []int
	for i := min; i <= max; i += step {
		rangeList = append(rangeList, i)
	}
	return rangeList
}

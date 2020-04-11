package trading

import (
	"encoding/gob"
	"fmt"
	"github.com/cwntr/go-dex-client/pkg/common"
	"os"
	"sync"
	"time"

	"github.com/cwntr/go-dex-client/lssdrpc"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

var (
	// will be backed up in file that it can be recovered after restart
	myOrders     *cache.Cache //map[trading_pair] = OrderCollection
	marketOrders *cache.Cache //map[trading_pair] = OrderCollection

	cachePathMyOrders = "cache/myOrders.dat"

	//temporary mapping until its stored in persisted cache
	orderStates sync.Map
	orderTypes  sync.Map //map[string]string
	orderLabels sync.Map //map[string]string

	logger *logrus.Logger
	bot    *Bot
)

func init() {
	gob.Register(OrderCollection{})

	orderStates = sync.Map{}
	orderTypes = sync.Map{}
	orderLabels = sync.Map{}
}

func SetLogger(logLvl string) {
	logger = logrus.New()
	logLevel, err := logrus.ParseLevel(logLvl)
	if err != nil {
		fmt.Printf("err: unable to parse logLevel, err: %v", err)
	}
	logger.Level = logLevel
}
func GetMyOrdersCache() *cache.Cache {
	return myOrders
}

func GetMyOrderByID(tradingPair string, orderId string) *Order {
	cacheCollection, isFound := myOrders.Get(tradingPair)
	if !isFound {
		return nil
	} else {
		cached, ok := cacheCollection.(OrderCollection)
		if ok {
			for _, or := range cached {
				if or.OrderId == orderId {
					return &or
				}
			}
		}
	}
	return nil
}

// SetupCache func will read previously used cache files to get back to the order state it was before shutting down
func SetupCache(tradingBot *Bot) {
	bot = tradingBot

	if _, err := os.Stat("cache"); os.IsNotExist(err) {
		err = os.Mkdir("cache", os.ModePerm)
		if err != nil {
			logger.Errorf("err setting up cache dir, err: %v", err)
		}
	}

	myOrders = cache.New(cache.NoExpiration, time.Duration(0))
	if fileExists(cachePathMyOrders) {
		err := myOrders.LoadFile(cachePathMyOrders)
		if err != nil {
			logger.Errorf("err loading myOrders cache to file, err: %v", err)
		} else {
			logger.Debugf("loaded ownOrder items (%d) from cache ", len(myOrders.Items()))
			myOrders = cache.NewFrom(cache.NoExpiration, time.Duration(0), myOrders.Items())
		}
	}
	marketOrders = cache.New(cache.NoExpiration, time.Duration(0))
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func BackupCache() {
	for _, pair := range common.GetPairs() {
		orders := GetMemoryEnrichedMyOrders(pair)
		myOrders.Set(pair, orders, 0)
	}

	err := myOrders.SaveFile(cachePathMyOrders)
	if err != nil {
		logger.Errorf("err saving myOrders cache to file, err: %v", err)
		return
	}
}

func syncMyOrderBook(c *cache.Cache, orderBook OrderCollection, tradingPair string) {
	cacheCollection, isFound := c.Get(tradingPair)
	if !isFound {
		// no cache yet, first iteration
		c.Set(tradingPair, orderBook, 0)
		return
	} else {
		cached, ok := cacheCollection.(OrderCollection)
		if ok {
			logger.Debugf("syncMyOrderBook found cached orderCollection: %d", len(cached))
			// sync all orderIds
			newOrderCount := 0
			cachedUpdated := &cached
			for _, item := range orderBook {
				isNew := cached.GetByOrderID(item.OrderId)
				if isNew == nil {
					//not found -> add because its new in order book
					if val, ok := orderTypes.Load(item.OrderId); ok {
						item.OrderType = val.(string)
					}

					if val, ok := orderLabels.Load(item.OrderId); ok {
						item.Label = val.(string)
					}
					item.Status = OrderStatusPlaced
					cachedUpdated = cached.Add(item)
					newOrderCount++
				}
			}
			logger.Debugf("newOrderCount: %d", newOrderCount)

			// sync order status
			coll := &OrderCollection{}
			for _, oc := range *cachedUpdated {
				if oc.Status == OrderStatusPlaced {

					//Not in current orderBook, but status was placed before? -> completed if it was not cancelled
					if orderBook.GetByOrderID(oc.OrderId) == nil {
						isCanceled := false
						if val, ok := orderStates.Load(oc.OrderId); ok {
							if val == OrderStatusCanceled {
								isCanceled = true
								oc.Status = OrderStatusCanceled
								logger.Debugf("order moved to status canceled: %s", oc.OrderId)
								if val, ok := orderTypes.Load(oc.OrderId); ok {
									oc.OrderType = val.(string)
								}
								if val, ok := orderLabels.Load(oc.OrderId); ok {
									oc.Label = val.(string)
								}
							}
						}
						if !isCanceled {
							logger.Debugf("order moved to status completed: %s", oc.OrderId)
							oc.Status = OrderStatusCompleted
							if val, ok := orderTypes.Load(oc.OrderId); ok {
								oc.OrderType = val.(string)
							}
							if val, ok := orderLabels.Load(oc.OrderId); ok {
								oc.Label = val.(string)
							}
						}
					}
				}
				coll = coll.Add(oc)
			}
			c.Set(tradingPair, *coll, 0)
		} else {
			logger.Errorln("syncMyOrderBook orderCollection convert failed: ")
		}
	}
}

func syncMarketOrderBook(c *cache.Cache, orderBook OrderCollection, tradingPair string) {
	cacheCollection, isFound := c.Get(tradingPair)
	if !isFound {
		// no cache yet, first iteration
		c.Set(tradingPair, orderBook, 0)
		return
	} else {
		cached, ok := cacheCollection.(OrderCollection)
		if ok {
			logger.Debugf("syncMarketOrderBook found cached orderCollection: %d", len(cached))
			if len(cached) != len(orderBook) {
				logger.Debugf("market change detected -> cache %d != new orderBook %d", len(cached), len(orderBook))
			}
			c.Set(tradingPair, orderBook, 0)
		} else {
			logger.Errorln("syncMarketOrderBook orderCollection convert failed..")
		}
	}
}

func updateOrderStatus(c *cache.Cache, order lssdrpc.Order, tradingPair string, status string, orderType string, label string) error {
	defer BackupCache()
	cacheCollection, isAlreadyInCache := c.Get(tradingPair)
	if !isAlreadyInCache {
		return nil
	}
	cached, ok := cacheCollection.(OrderCollection)
	if !ok {
		return nil
	}

	if status == OrderStatusPlaced {
		o := toOrder(order)
		o.OrderType = orderType
		o.Label = label
		o.Status = OrderStatusPlaced
		newColl := cached.Add(o)
		myOrders.Set(tradingPair, *newColl, 0)
		return nil
	}

	if status == OrderStatusCanceled {
		for _, o := range cached {
			if o.OrderId == order.OrderId {
				orderStates.Store(order.OrderId, OrderStatusCanceled)
				orderTypes.Store(order.OrderId, o.OrderType)
				orderLabels.Store(order.OrderId, o.Label)
			}
		}
		return nil
	}

	isAlreadyInCache = false
	for _, oldItem := range cached {
		if oldItem.OrderId == order.OrderId {
			logger.Debugf("updated order (%s) status from %s to %s", order.OrderId, oldItem.Status, status)
			oldItem.Status = status
			oldItem.OrderType = orderType
			oldItem.Label = label

			orderStates.Store(order.OrderId, oldItem.Status)
			orderTypes.Store(order.OrderId, oldItem.OrderType)
			orderLabels.Store(order.OrderId, label)
			isAlreadyInCache = true
		}
	}

	if !isAlreadyInCache {
		o := toOrder(order)

		o.Status = status
		o.OrderType = orderType
		o.Label = label

		orderStates.Store(order.OrderId, o.Status)
		orderTypes.Store(order.OrderId, o.OrderType)
		orderLabels.Store(order.OrderId, label)

		logger.Debugf("new order details: %v", o)
		logger.Debugf("added order (%s) new order with status %s, orderType: %s", o.OrderId, o.Status, o.OrderType)
	}
	return nil
}

package trading

import (
	"sort"
	"time"

	"strconv"

	"github.com/cwntr/go-dex-client/lssdrpc"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

const (
	OrderSideSell = "sell"
	OrderSideBuy  = "buy"

	OrderStatusPlaced    = "placed"
	OrderStatusCanceled  = "canceled"
	OrderStatusCompleted = "completed"
)

type Order struct {
	PairId    string `json:"pairId"`
	OrderId   string `json:"orderId"`
	Price     int64  `json:"price"` //satoshis
	Funds     int64  `json:"funds"` //satoshis
	CreatedAt uint64 `json:"createdAt"`
	Side      string `json:"side"`
	Status    string `json:"status"`
	OrderType string `json:"type"`
	Label     string `json:"label,omitempty"`
	IsOwn     bool   `json:"isOwn"`
}

type OrderCollection []Order

func (o *OrderCollection) GetByOrderID(orderId string) *Order {
	for _, item := range *o {
		if item.OrderId == orderId {
			return &item
		}
	}
	return nil
}

func (o *OrderCollection) Unique() *OrderCollection {
	keys := make(map[string]bool)
	newColl := OrderCollection{}
	for _, order := range *o {
		if _, value := keys[order.OrderId]; !value {
			keys[order.OrderId] = true
			newColl = append(newColl, order)
		}
	}
	return &newColl
}

func (o *OrderCollection) Add(order Order) *OrderCollection {
	oc := *o
	oc = append(oc, order)
	return &oc
}

func (o *OrderCollection) GetPriceTopOrdersSorted(side string, limit int, excludingMyOrders bool) *OrderCollection {
	coll := OrderCollection{}
	for _, order := range *o {
		if order.Side == side {
			if !order.IsOwn {
				coll = append(coll, order)
			}
		}
	}

	if side == OrderSideBuy {
		sort.Slice(coll, func(i, j int) bool {
			if coll[i].Price > coll[j].Price {
				return true
			}
			if coll[i].Price < coll[j].Price {
				return false
			}
			return coll[i].Price > coll[j].Price
		})
	} else {
		sort.Slice(coll, func(i, j int) bool {
			if coll[i].Price < coll[j].Price {
				return true
			}
			if coll[i].Price > coll[j].Price {
				return false
			}
			return coll[i].Price < coll[j].Price
		})
	}
	if limit > 0 && limit <= len(coll)-1 {
		coll = coll[0:limit]
	}
	return &coll
}

func (o *OrderCollection) SumFunds(status string) (sum int64) {
	if o == nil {
		return
	}
	for _, order := range *o {
		if order.Status == status {
			sum += order.Funds
		}
	}
	return
}

func toMyCollection(orders []lssdrpc.Order, tradingPair string) (oc OrderCollection) {
	myOrderColl := GetMemoryEnrichedMyOrders(tradingPair)
	for _, o := range orders {
		if !o.IsOwnOrder {
			continue
		}
		myOrder := myOrderColl.GetByOrderID(o.OrderId)
		if myOrder == nil {
			continue
		}
		order := toOrder(o)
		oc = append(oc, order)
	}
	return oc
}

func toExternalCollection(orders []lssdrpc.Order, tradingPair string) (oc OrderCollection) {
	myOrderColl := GetMemoryEnrichedMyOrders(tradingPair)
	for _, o := range orders {
		myOrder := myOrderColl.GetByOrderID(o.OrderId)
		if myOrder != nil {
			continue
		}
		if o.IsOwnOrder {
			continue
		}
		order := toOrder(o)
		oc = append(oc, order)
	}
	return oc
}

func toOrder(o lssdrpc.Order) Order {
	order := Order{PairId: o.PairId}
	order.IsOwn = o.IsOwnOrder
	order.OrderId = o.OrderId
	p, _ := strconv.Atoi(o.Price.Value)
	order.Price = int64(p)
	f, _ := strconv.Atoi(o.Funds.Value)
	order.Funds = int64(f)
	order.CreatedAt = o.CreatedAt
	order.Side = o.Side.String()
	cachedOrder := GetFromCacheByOrderId(myOrders, o.PairId, o.OrderId)
	if cachedOrder != nil {
		co := *cachedOrder
		order.OrderType = co.OrderType
	}
	return order
}

func PeriodicOrderListener(bot *Bot, tradingPair string, withOwnOrders bool, interval time.Duration, logger *logrus.Entry) {
	ticker := time.NewTicker(interval)

	go func() {
		for {
			select {
			case _ = <-ticker.C:
				defer BackupCache()
				if withOwnOrders {
					orders, err := bot.ListOrders(tradingPair, true, false)
					logger.Debugf("currently (with mine: %v) orders: (%d)", withOwnOrders, len(orders))
					if err != nil {
						logger.Errorf("listener err: %v", err)
						break
					}
					oc := toMyCollection(orders, tradingPair)
					syncMyOrderBook(GetMyOrdersCache(), oc, tradingPair)
				} else {
					orders, err := bot.ListOrders(tradingPair, false, false)
					logger.Debugf("currently (with mine: %v) orders: (%d)", withOwnOrders, len(orders))
					if err != nil {
						logger.Errorf("listener err: %v", err)
						break
					}
					oc := toExternalCollection(orders, tradingPair)
					syncMarketOrderBook(marketOrders, oc, tradingPair)
				}
			}
		}
	}()
}

func GetFromCacheByOrderId(c *cache.Cache, tradingPair string, orderId string) *Order {
	cacheCollection, isFound := c.Get(tradingPair)
	if !isFound {
		return nil
	} else {
		cached, ok := cacheCollection.(OrderCollection)
		if ok {
			return cached.GetByOrderID(orderId)
		}
	}
	return nil
}

func GetMemoryEnrichedMyOrders(tp string) OrderCollection {
	cacheCollection, isFound := myOrders.Get(tp)
	if !isFound {
		return OrderCollection{}
	} else {
		cached, ok := cacheCollection.(OrderCollection)
		if ok {
			var enriched OrderCollection
			uniqueCache := cached.Unique()
			for _, item := range *uniqueCache {
				newItem := item
				if val, ok := orderStates.Load(item.OrderId); ok {
					if val != "" {
						newItem.Status = val.(string)
					}
				}
				if val, ok := orderTypes.Load(item.OrderId); ok {
					if val != "" {
						newItem.OrderType = val.(string)
					}
				}
				if val, ok := orderLabels.Load(item.OrderId); ok {
					if val != "" {
						newItem.Label = val.(string)
					}
				}
				enriched = append(enriched, newItem)
			}
			return enriched
		}
	}
	return OrderCollection{}
}

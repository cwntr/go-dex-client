package trading

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cwntr/go-dex-client/lssdrpc"
)

const (
	LSSDConnectionProblem = "connect: connection refused"
)

// Bot struct holds all the gRPC client connections, the resolved config and the logger
type Bot struct {
	OrderClient       lssdrpc.OrdersClient
	SwapClient        lssdrpc.SwapsClient
	CurrencyClient    lssdrpc.CurrenciesClient
	TradingPairClient lssdrpc.TradingPairsClient

	LNDConfig LNDConfig
	Log       *logrus.Entry
	LogLevel  string

	OrderLimit int
	Timeout    time.Duration

	SentryDSN string
}

// Init func will read the contents of the given certificate paths
func (t *Bot) Init() error {
	t.LNDConfig.Certs = make(map[string]string, 0)
	for currency, path := range t.LNDConfig.TLSPaths {
		b, err := ReadFile(path)
		if err != nil {
			return err
		}
		t.LNDConfig.Certs[currency] = string(b)
	}
	return nil
}

// AddCurrencies func will register all the currencies by performing a lssd grpc request for all certificates are specified in the config
func (t *Bot) AddCurrencies() error {
	ctx, cancel := context.WithTimeout(context.Background(), t.Timeout)
	defer cancel()
	defer track(runningTime("bot.AddCurrencies"))

	for currency, cert := range t.LNDConfig.Certs {
		c := t.LNDConfig.Connections
		x := c[currency]
		cr := &lssdrpc.AddCurrencyRequest{
			Currency:   currency,
			LndChannel: x.Format(),
		}
		cr.TlsCert = &lssdrpc.AddCurrencyRequest_RawCert{RawCert: cert}
		_, err := t.CurrencyClient.AddCurrency(ctx, cr)
		if err != nil {
			return err
		}
	}
	return nil
}

// SubscribeSwaps func will subscribe to swaps
func (t *Bot) SubscribeSwaps() error {
	ctx, cancel := context.WithTimeout(context.Background(), t.Timeout)
	defer cancel()
	defer track(runningTime("bot.SubscribeSwaps"))

	_, err := t.SwapClient.SubscribeSwaps(ctx, &lssdrpc.SubscribeSwapsRequest{})
	if err != nil {
		if IsLSSDConnProblem(err) {
			//@TODO graceful shutdown of active orders
			logger.Errorf("err connecting to LSSD, err: %v", err)
		} else {
			logger.Errorf("err subscribing to swaps : %v", err)
		}
	}
	return err
}

// EnableTradingPair func will enable the given tradingPair
func (t *Bot) EnableTradingPair(pair string) error {
	ctx, cancel := context.WithTimeout(context.Background(), t.Timeout)
	defer cancel()
	defer track(runningTime("bot.EnableTradingPair"))

	_, err := bot.TradingPairClient.EnableTradingPair(ctx, &lssdrpc.EnableTradingPairRequest{PairId: pair})
	if err != nil {
		if IsLSSDConnProblem(err) {
			//@TODO graceful shutdown of active orders
			logger.Errorf("err connecting to LSSD, err: %v", err)
		} else {
			logger.Errorf("err enabling trading pairs : %v", err)
		}
	}
	return err
}

// SubscribeOrders func will subscribe to orders
func (t *Bot) SubscribeOrders() error {
	ctx, cancel := context.WithTimeout(context.Background(), t.Timeout)
	defer cancel()
	defer track(runningTime("bot.SubscribeOrders"))

	_, err := t.OrderClient.SubscribeOrders(ctx, &lssdrpc.SubscribeOrdersRequest{})
	if err != nil {
		if IsLSSDConnProblem(err) {
			//@TODO graceful shutdown of active orders
			logger.Errorf("err connecting to LSSD, err: %v", err)
		} else {
			logger.Errorf("err subscribing to orders : %v", err)
		}
	}
	return err
}

// ListOrders func will perform a lssd gRPC request to the ListOrders endpoint of the Orders service for a specified trading pair
func (t *Bot) ListOrders(tradingPair string, myOrders bool, printCLI bool) ([]lssdrpc.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), t.Timeout)
	defer cancel()
	defer track(runningTime("bot.ListOrders"))

	var orders []lssdrpc.Order
	var processed uint32
	var limit = uint32(t.OrderLimit)
	for {
		res, err := t.OrderClient.ListOrders(ctx, &lssdrpc.ListOrdersRequest{
			PairId:           tradingPair,
			IncludeOwnOrders: myOrders,
			Skip:             processed,
			Limit:            limit,
		})
		if err != nil {
			if IsLSSDConnProblem(err) {
				//@TODO graceful shutdown of active orders
				logger.Errorf("err connecting to LSSD, err: %v", err)
			} else {
				logger.Errorf("err listing orders : %v", err)
			}
			return nil, err
		}
		if len(res.Orders) == 0 {
			break
		}
		for _, o := range res.Orders {
			if o != nil {
				orders = append(orders, *o)
			}
		}
		processed += uint32(len(res.Orders))
	}
	if printCLI {
		for i, o := range orders {
			fmt.Printf("id: %d | pair: %s | side: %s | orderId: %s | price: %v | funds: %v | isMy: %v \n", i, o.PairId, o.Side, o.OrderId, o.Price, o.Funds, o.IsOwnOrder)
		}
	}
	return orders, nil
}

// PlaceOrder func will perform a lssd gRPC request to the PlaceOrder endpoint of the Orders service for a specified trading pair and order settings
func (t *Bot) PlaceOrder(tradingPair string, price int64, funds int64, side lssdrpc.OrderSide, orderType string, label string) (*lssdrpc.PlaceOrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), t.Timeout)
	defer cancel()
	defer track(runningTime("bot.PlaceOrder"))

	p := strconv.FormatInt(price, 10)
	priceBigInt := &lssdrpc.BigInteger{Value: p}

	s := strconv.FormatInt(funds, 10)
	fundsBigInt := &lssdrpc.BigInteger{Value: s}
	order := &lssdrpc.PlaceOrderRequest{
		PairId: tradingPair,
		Side:   side,
		Funds:  fundsBigInt,
		Price:  priceBigInt,
	}
	res, err := t.OrderClient.PlaceOrder(ctx, order)
	if err != nil {
		if IsLSSDConnProblem(err) {
			//@TODO graceful shutdown of active orders
			logger.Errorf("err connecting to LSSD, err: %v", err)
		} else {
			logger.Errorf("err placing order: %v", err)
		}
		return nil, err
	}
	// check if the order was placed
	outcomeOrder, isPlaced := res.Outcome.(*lssdrpc.PlaceOrderResponse_Order)
	if isPlaced {
		if outcomeOrder.Order != nil {
			_ = updateOrderStatus(myOrders, *outcomeOrder.Order, tradingPair, OrderStatusPlaced, orderType, label)
		}
		return res, nil
	}

	// check if the order was swapped already
	outcomeSwapped, isSwapped := res.Outcome.(*lssdrpc.PlaceOrderResponse_SwapSuccess)
	if isSwapped {
		if outcomeSwapped.SwapSuccess != nil {
			o := lssdrpc.Order{}
			o.OrderId = outcomeSwapped.SwapSuccess.OrderId
			o.Funds = outcomeSwapped.SwapSuccess.Funds
			o.Price = outcomeSwapped.SwapSuccess.Price
			o.CreatedAt = uint64(time.Now().Unix())
			o.Side = side
			_ = updateOrderStatus(myOrders, o, tradingPair, OrderStatusCompleted, orderType, label)
		}
		return res, nil
	}

	// check if the order failed (swap / orderbook error)
	outcomeFailure, isFailure := res.Outcome.(*lssdrpc.PlaceOrderResponse_Failure)
	if isFailure {
		if outcomeFailure.Failure.Failure != nil {
			//try OrderbookFailure
			orderbookFailure, isSwapFailure := outcomeFailure.Failure.Failure.(*lssdrpc.PlaceOrderFailure_OrderbookFailure)
			if isSwapFailure && orderbookFailure.OrderbookFailure != nil {
				logger.Errorf(
					"order outcome failure [orderbook-failure] (placingOrderDetails: funds: %d, price: %d, side: %s ), failureDetails: (funds: %s, requiredFee: %s, failureReason: %s)",
					funds,
					price,
					side.String(),
					orderbookFailure.OrderbookFailure.Funds.Value,
					orderbookFailure.OrderbookFailure.RequiredFee.Value,
					orderbookFailure.OrderbookFailure.FailureReason,
				)
				return res, nil
			} else if isSwapFailure {
				logger.Errorf(
					"order outcome failure [orderbook-failure] (placingOrderDetails: funds: %d, price: %d, side: %s ), failureDetails: (empty)",
					funds,
					price,
					side.String(),
				)
				return res, nil
			}

			//try SwapFailure
			swapFailure, isSwapFailure := outcomeFailure.Failure.Failure.(*lssdrpc.PlaceOrderFailure_SwapFailure)
			if isSwapFailure && swapFailure.SwapFailure != nil {
				logger.Errorf(
					"order outcome failure [swap-failure] (placingOrderDetails: funds: %d, price: %d, side: %s ), failureDetails: (funds: %s, orderId: %s, failureReason: %s)",
					funds,
					price,
					side.String(),
					swapFailure.SwapFailure.Funds,
					swapFailure.SwapFailure.OrderId,
					swapFailure.SwapFailure.FailureReason,
				)
				if strings.Contains(swapFailure.SwapFailure.FailureReason, "map::at") {
					InitLSSD()
					time.Sleep(time.Second * 5)
				}
				return res, nil
			} else if isSwapFailure {
				logger.Errorf(
					"order outcome failure [swap-failure] (placingOrderDetails: funds: %d, price: %d, side: %s ), failureDetails: (empty)",
					funds,
					price,
					side.String(),
				)
				return res, nil
			}
		}

		// its a failure but failure-type cannot be determined
		logger.Errorf("order outcome failure [?-failure] (funds: %d, price: %d, side: %s ), cannot determine failure type", funds, price, side.String())
		return res, nil
	}
	return res, nil
}

// CancelOrder func will perform a lssd gRPC request to the PlaceOrder endpoint of the Orders service for a specified trading pair and order settings
func (t *Bot) CancelOrder(tradingPair string, orderID string) (*lssdrpc.CancelOrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), t.Timeout)
	defer cancel()
	defer track(runningTime("bot.CancelOrder"))

	order := &lssdrpc.CancelOrderRequest{
		PairId:  tradingPair,
		OrderId: orderID,
	}
	o := lssdrpc.Order{OrderId: orderID}
	cacheOrder := GetMyOrderByID(tradingPair, orderID)

	res, err := t.OrderClient.CancelOrder(ctx, order)
	if err != nil {
		if IsLSSDConnProblem(err) {
			//@TODO graceful shutdown of active orders
			logger.Errorf("err connecting to LSSD, err: %v", err)
		} else {
			logger.Errorf("err canceling order : %v", err)
		}
		return nil, err
	}
	if cacheOrder != nil {
		updateOrderStatus(myOrders, o, tradingPair, OrderStatusCanceled, cacheOrder.OrderType, cacheOrder.Label)
	}
	return res, nil
}

// NewBot func will create a new Bot object based on all the lssd client connections and the resolved config
func NewBot(o lssdrpc.OrdersClient, s lssdrpc.SwapsClient, c lssdrpc.CurrenciesClient, t lssdrpc.TradingPairsClient, lndConfig LNDConfig, timeout time.Duration, orderLimit int) (*Bot, error) {
	if lndConfig.IsEmpty() {
		return nil, fmt.Errorf("lndConfig is empty")
	}

	b := &Bot{
		OrderClient:       o,
		SwapClient:        s,
		CurrencyClient:    c,
		TradingPairClient: t,
		LNDConfig:         lndConfig,
		Timeout:           timeout,
		OrderLimit:        orderLimit,
	}
	b.Log = logrus.WithFields(logrus.Fields{"context": "bot"})
	err := b.Init()
	return b, err
}

// ReadFile will read the content of a given filepath, e.g. the certificate contents
func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func runningTime(s string) (string, time.Time) {
	logger.Debugf("-- bot time [%s] start", s)
	return s, time.Now()
}

func track(s string, startTime time.Time) {
	endTime := time.Now()
	logger.Debugf("-- bot time [%s] end | took: %s", s, endTime.Sub(startTime))
}

func IsLSSDConnProblem(err error) bool {
	if strings.Contains(err.Error(), LSSDConnectionProblem) {
		return true
	}
	return false
}

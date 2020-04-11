package trading

import (
	"github.com/cwntr/go-dex-client/lssdrpc"
	"time"

	"github.com/cwntr/go-dex-client/pkg/common"
	"github.com/sirupsen/logrus"
)

func InitLSSD() error {
	//Subscribe Swaps
	err := bot.SubscribeSwaps()
	if err != nil {
		logger.Errorf("error subscribing to swaps, err %v", err)
		return err
	}
	logger.Debugln("subscribed swaps")

	//Subscribe Orders
	err = bot.SubscribeOrders()
	if err != nil {
		logger.Errorf("error subscribing to orders, err %v", err)
		return err
	}
	logger.Debugln("subscribed orders")

	//Add Currencies
	err = bot.AddCurrencies()
	if err != nil {
		logger.Errorf("err: %v", err)
		return err
	}

	err = bot.EnableTradingPair(common.PairLTCBTC)
	if err != nil {
		logger.Errorf("err: %v", err)
		return err
	}

	err = bot.EnableTradingPair(common.PairXSNBTC)
	if err != nil {
		logger.Errorf("err: %v", err)
		return err
	}
	err = bot.EnableTradingPair(common.PairXSNLTC)
	if err != nil {
		logger.Errorf("err: %v", err)
		return err
	}

	logger.Debugln("added currencies")
	return nil
}

// Periodically perform request to check if LSSD is running
// Its using "enableTradingPair" instead of an actual ping endpoint but its just for the concept
func PeriodicPing(interval time.Duration, entry *logrus.Entry) {
	ticker := time.NewTicker(interval)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				entry.Debugf("ping %s", "LSSD")

				err := bot.EnableTradingPair(common.PairXSNLTC)
				entry.Debugf("ping EnableTradingPair error: %v", err)
				if err != nil && IsLSSDConnProblem(err) {
					break
				}
				//@Todo LSSD service is back up -> clean up
			}
		}
	}()
}

func toRPCSide(s string) (side lssdrpc.OrderSide) {
	if s == OrderSideSell {
		side = lssdrpc.OrderSide_sell
	} else {
		side = lssdrpc.OrderSide_buy
	}
	return side
}

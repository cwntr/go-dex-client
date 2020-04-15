package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"

	"github.com/cwntr/go-dex-client/pkg/common"

	"github.com/sirupsen/logrus"

	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cwntr/go-dex-client/trading"
)

var (
	cfg    Config
	logger *logrus.Entry
	bot    *trading.Bot
)

func main() {
	//Read global config
	err := readConfig()
	if err != nil {
		logger.Errorf("error reading global config, err %v", err)
		os.Exit(1)
		return
	}

	logLevel, err := logrus.ParseLevel(cfg.Bot.LogLevel)
	if err != nil {
		fmt.Printf("err: unable to parse logLevel, err: %v", err)
		os.Exit(1)
		return
	}
	logrus.SetLevel(logLevel)
	trading.SetLogger(cfg.Bot.LogLevel)
	logger = logrus.WithFields(logrus.Fields{"context": "main"})

	logger.Infoln("global config loaded")

	//Initialize Clients
	tpClient, tpConn := createTradingPairClient()
	defer tpConn.Close()
	oClient, oConn := createOrdersClient()
	defer oConn.Close()
	cClient, cConn := createCurrencyClient()
	defer cConn.Close()
	sClient, sConn := createSwapClient()
	defer sConn.Close()
	logger.Infoln("clients initiated")

	//Initialize LNDConfig
	tradingCfg := trading.NewConfig()
	err = tradingCfg.Add(common.CurrencyXSN, cfg.XSN.CertPath, cfg.XSN.Host, cfg.XSN.Port)
	if err != nil {
		logger.Errorf("error adding XSN to trading config, err %v", err)
		return
	}
	err = tradingCfg.Add(common.CurrencyLTC, cfg.LTC.CertPath, cfg.LTC.Host, cfg.LTC.Port)
	if err != nil {
		logger.Errorf("error adding LTC to trading config, err %v", err)
		return
	}
	err = tradingCfg.Add(common.CurrencyBTC, cfg.BTC.CertPath, cfg.BTC.Host, cfg.BTC.Port)
	if err != nil {
		logger.Errorf("error adding BTC to trading config, err %v", err)
		return
	}
	logger.Infoln("trading config loaded")

	//Initialize Bot
	bot, err = trading.NewBot(oClient, sClient, cClient, tpClient, tradingCfg, cfg.LSSDConfig.Timeout, cfg.Bot.OrderLimit)
	if err != nil {
		logger.Errorf("error initializing trading bot, err %v", err)
		return
	}
	logger.Infoln("trading bot initialized")

	//Perform infrastructure checks (channels to hub peers have enough capacity)
	_, err = checkInfra(common.CurrencyXSN, false)
	if err != nil {
		logger.Errorf("infra check, err %v", err)
		return
	}
	_, err = checkInfra(common.CurrencyLTC, false)
	if err != nil {
		logger.Errorf("infra check, err %v", err)
		return
	}
	_, err = checkInfra(common.CurrencyBTC, false)
	if err != nil {
		logger.Errorf("infra check, err %v", err)
		return
	}

	// Initialize Cache of orders
	trading.SetupCache(bot)

	//Start subscribing all mandatory services
	err = trading.InitLSSD()
	if err != nil {
		logger.Errorf("subscription failed, err %v", err)
		return
	}

	defer trading.BackupCache()

	//Register jobs
	registerJobs(bot)

	// listen to termination signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// centralized done channel
	done := make(chan struct{})
	errSig := make(chan struct{})
	// wait for all services to close gracefully
	wgClose := &sync.WaitGroup{}
	defer wgClose.Wait()

	//Setup routes
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	addRoutes(router)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Bot.Host, cfg.Bot.Port),
		Handler: router,
	}

	// setup server
	go initHTTPServer(srv, done, errSig, wgClose)

	defer close(done)
	select {
	case <-signalChan:
		return
	case <-errSig:
		return
	}
}

// Handle gracefully shutdown
func initHTTPServer(srv *http.Server, done chan struct{}, errSig chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	defer func() {
		srv.Shutdown(context.Background())
		wg.Done()
	}()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
		errSig <- struct{}{}
	}()
	<-done
}

func registerJobs(bot *trading.Bot) {
	//Periodic check if LSSD crashed, and back up again to reactive strategies
	logPing := logrus.WithFields(logrus.Fields{"context": "ping"})
	go trading.PeriodicPing(time.Second*30, logPing)

	//Periodic checkers for XSN_LTC trading pair
	logOrderXSNLTC := logrus.WithFields(logrus.Fields{"context": "my-orders", "trading-pair": common.PairXSNLTC})
	logMarketXSNLTC := logrus.WithFields(logrus.Fields{"context": "market-orders", "trading-pair": common.PairXSNLTC})
	go trading.PeriodicOrderListener(bot, common.PairXSNLTC, true, cfg.Bot.JobInterval, logOrderXSNLTC)
	go trading.PeriodicOrderListener(bot, common.PairXSNLTC, false, cfg.Bot.JobInterval+time.Second*1, logMarketXSNLTC)

	//Periodic checkers for XSN_BTC trading pair
	logOrderXSNBTC := logrus.WithFields(logrus.Fields{"context": "my-orders", "trading-pair": common.PairXSNBTC})
	logMarketXSNBTC := logrus.WithFields(logrus.Fields{"context": "market-orders", "trading-pair": common.PairXSNBTC})
	go trading.PeriodicOrderListener(bot, common.PairXSNBTC, true, cfg.Bot.JobInterval, logOrderXSNBTC)
	go trading.PeriodicOrderListener(bot, common.PairXSNBTC, false, cfg.Bot.JobInterval+time.Second*1, logMarketXSNBTC)

	//Periodic checkers for LTC_BTC trading pair
	logOrderLTCBTC := logrus.WithFields(logrus.Fields{"context": "my-orders", "trading-pair": common.PairLTCBTC})
	logMarketLTCBTC := logrus.WithFields(logrus.Fields{"context": "market-orders", "trading-pair": common.PairLTCBTC})
	go trading.PeriodicOrderListener(bot, common.PairLTCBTC, true, cfg.Bot.JobInterval, logOrderLTCBTC)
	go trading.PeriodicOrderListener(bot, common.PairLTCBTC, false, cfg.Bot.JobInterval+time.Second*1, logMarketLTCBTC)
}

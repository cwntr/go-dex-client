package main

import (
	"encoding/json"
	"os"
	"time"
)

// HubPeer struct holds the remoteKey and the address of the peer
type HubPeer struct {
	RemoteKey string `json:"remoteKey"`
	Address   string `json:"address"`
}

// LNDConfig struct holds mandatory information how to connect to the lncli
type LNDConfig struct {
	Directory string    `json:"lndDir"`
	CertPath  string    `json:"certPath"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	HubPeers  []HubPeer `json:"hubPeers"`
}

// LSSDConfig struct holds mandatory information how connect to the LSSD which is a bridge to the Stakenet swaps and Stakenet orderbook
type LSSDConfig struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	TimeoutStr string `json:"timeout"`
	Timeout    time.Duration
}

// BotConfig struct holds mandatory data of the how the bot is configured
// Also, it contains the path to the lncli to fetch useful information for the respective LND
type BotConfig struct {
	JobIntervalStr string `json:"jobInterval"`
	JobInterval    time.Duration
	Host           string `json:"host"`
	Port           int    `json:"port"`
	LNCLIPath      string `json:"lnCLIPath"`
	LogLevel       string `json:"logLevel"`
	OrderLimit     int    `json:"orderLimit"`
}

// Config struct keeps all the individual configs together
type Config struct {
	Bot        BotConfig  `json:"botCfg"`
	LSSDConfig LSSDConfig `json:"lssdConfig"`
	XSN        LNDConfig  `json:"xsnLNDConfig"`
	LTC        LNDConfig  `json:"ltcLNDConfig"`
	BTC        LNDConfig  `json:"btcLNDConfig"`
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

// Reads the entire config, "cfg.json" is hardcoded and must be placed on same level as the application binary
func readConfig() error {
	file, err := os.Open("cfg.json")
	if err != nil {
		logger.Fatalf("can't open config file: %v", err)
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		logger.Fatalf("can't decode config JSON: %v", err)
		return err
	}

	cfg.LSSDConfig.Timeout, err = time.ParseDuration(cfg.LSSDConfig.TimeoutStr)
	if err != nil {
		return err
	}
	cfg.Bot.JobInterval, err = time.ParseDuration(cfg.Bot.JobIntervalStr)
	if err != nil {
		return err
	}
	return nil
}

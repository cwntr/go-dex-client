package main

import (
	"github.com/cwntr/go-dex-client/pkg/common"
	"strconv"

	"fmt"

	"github.com/cwntr/go-dex-client/lncli"
)

type InfraChecks []Check
type Check struct {
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Result   bool   `json:"result"`
	Details  string `json:"details"`
}

type Channel struct {
	Currency         string `json:"currency"`
	Active           bool   `json:"active"`
	ChanID           string `json:"chan_id"`
	Capacity         string `json:"capacity"`
	LocalBalance     string `json:"local_balance"`
	RemoteBalance    string `json:"remote_balance"`
	UnsettledBalance string `json:"unsettled_balance"`
	ChanStatusFlags  string `json:"chan_status_flags"`
}

func getMinimumSatoshis(currency string) int64 {
	switch currency {
	case common.CurrencyXSN:
		return 60000
	case common.CurrencyLTC:
		return 275000
	case common.CurrencyBTC:
		return 20000
	}
	return 0
}

func getHubPeers(currency string) []HubPeer {
	switch currency {
	case common.CurrencyXSN:
		return cfg.XSN.HubPeers
	case common.CurrencyLTC:
		return cfg.LTC.HubPeers
	case common.CurrencyBTC:
		return cfg.BTC.HubPeers
	}
	return []HubPeer{}
}

func getChannelCapacities(currency string) (min int64, max int64) {
	switch currency {
	case common.CurrencyXSN:
		min = 60000
		max = 100000000000
		return
	case common.CurrencyLTC:
		min = 275000
		max = 1000000000
		return
	case common.CurrencyBTC:
		min = 20000
		max = 1600000
		return
	}
	return
}

const (
	InfraCheckNamePeers    = "HubPeers"
	InfraCheckNameChannels = "HubChannels"
)

func getChannels(currency string) []Channel {
	var lndCfg LNDConfig

	switch currency {
	case common.CurrencyXSN:
		lndCfg = cfg.XSN
	case common.CurrencyLTC:
		lndCfg = cfg.LTC
	case common.CurrencyBTC:
		lndCfg = cfg.BTC
	default:
		return []Channel{}
	}

	// Hub Channels
	lc, err := lncli.GetListChannels(cfg.Bot.LNCLIPath, lndCfg.Directory, lndCfg.Host, lndCfg.Port)
	if err != nil {
		logger.Errorf("unable to get listChannels, err: %v", err)
		return []Channel{}
	}

	var channels []Channel
	for _, c := range lc.Channels {
		channels = append(channels, Channel{
			Currency:         currency,
			Active:           c.Active,
			ChanID:           c.ChanID,
			Capacity:         c.Capacity,
			LocalBalance:     c.LocalBalance,
			RemoteBalance:    c.RemoteBalance,
			UnsettledBalance: c.UnsettledBalance,
			ChanStatusFlags:  c.ChanStatusFlags,
		})
	}
	return channels
}

func checkInfra(currency string, printCLI bool) (InfraChecks, error) {
	var infra InfraChecks
	var lndCfg LNDConfig

	switch currency {
	case common.CurrencyXSN:
		lndCfg = cfg.XSN
	case common.CurrencyLTC:
		lndCfg = cfg.LTC
	case common.CurrencyBTC:
		lndCfg = cfg.BTC
	default:
		return infra, fmt.Errorf("unable to resolve trading pair")
	}
	hubPeersCheck := Check{Name: InfraCheckNamePeers, Currency: currency}
	hubChannelsCheck := Check{Name: InfraCheckNameChannels, Currency: currency}

	// HubPeers
	lp, err := lncli.GetPeers(cfg.Bot.LNCLIPath, lndCfg.Directory, lndCfg.Host, lndCfg.Port)
	if err != nil {
		logger.Errorf("unable to list peers, err: %v", err)
		hubPeersCheck.Details = "unable to list peers"
		return append(infra, hubPeersCheck), err
	}

	peerCheck := false
	hubPeers := getHubPeers(currency)
	for _, p := range lp.Peers {
		for _, hp := range hubPeers {
			if p.PubKey == hp.RemoteKey && p.Address == hp.Address {
				if printCLI {
					logger.Debugf("hub peer found: %s@%s", p.PubKey, p.Address)
				}
				hubPeersCheck.Details = fmt.Sprintf("hub peer found: %s@%s", p.PubKey, p.Address)
				peerCheck = true
			}
		}
	}
	if peerCheck {
		if printCLI {
			logger.Debugf("Infra check: %s - peers (1/2) OK ", currency)
		}
		hubPeersCheck.Result = true

	} else {
		logger.Errorf("Infra check: %s - peers (1/2) FAILED", currency)
		return append(infra, hubPeersCheck), err
	}

	infra = append(infra, hubPeersCheck)

	// Hub Channels
	lc, err := lncli.GetListChannels(cfg.Bot.LNCLIPath, lndCfg.Directory, lndCfg.Host, lndCfg.Port)
	if err != nil {
		logger.Errorf("unable to get listChannels, err: %v", err)
		return infra, err
	}

	capacityMin, capacityMax := getChannelCapacities(currency)
	channelCheck := false
	cntHubChannels := 0
	for _, c := range lc.Channels {
		for _, hp := range hubPeers {
			if hp.RemoteKey == c.RemotePubkey {
				capacity, _ := strconv.ParseInt(c.Capacity, 10, 64)

				if capacity >= capacityMin && capacity <= capacityMax {
					localBalance, _ := strconv.ParseInt(c.LocalBalance, 10, 64)
					lbf := float64(localBalance)

					remoteBalance, _ := strconv.ParseInt(c.RemoteBalance, 10, 64)
					rbf := float64(remoteBalance)
					if printCLI {
						logger.Debugf("[%s] channel capacity to hub is OK (chanID: %s) local_balance: %d (%.8f %s), remote_balance: %d (%.8f %s) ",
							currency,
							c.ChanID,
							localBalance,
							lbf/1e8,
							currency,
							remoteBalance,
							rbf/1e8,
							currency,
						)
					}
					channelCheck = true
					cntHubChannels++
				} else {
					logger.Errorf("[%s] channel capacity NOT OK -> your capacity %d must be between %d and %d ", currency, capacity, capacityMin, capacityMax)
				}
			}
		}
	}

	if cntHubChannels > 0 && cntHubChannels%2 != 0 {
		logger.Errorf("Infra check: %s - channels (2/2) FAILED", currency)
		hubChannelsCheck.Result = false
		hubChannelsCheck.Details = fmt.Sprintf("[%s] channels NOT OK: missing back channel(s)", currency)
		return append(infra, hubChannelsCheck), err
	}

	if channelCheck {
		if printCLI {
			logger.Debugf("Infra check: %s - channels (2/2) OK ", currency)
		}
		hubChannelsCheck.Result = true
		hubChannelsCheck.Details = "hub channel balances are ok"
	} else {
		logger.Errorf("Infra check: %s - channels (2/2) FAILED", currency)
		hubChannelsCheck.Details = fmt.Sprintf("channels capacity NOT OK: capacity must be between %d and %d ", capacityMin, capacityMax)
		return append(infra, hubChannelsCheck), err
	}
	if printCLI {
		logger.Debugf("Infra check: %s complete", currency)
		logger.Debugln("-------------------------------------")
	}
	infra = append(infra, hubChannelsCheck)
	return infra, nil
}

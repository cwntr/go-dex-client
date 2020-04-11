package lncli

import (
	"fmt"
	"strconv"

	"github.com/cwntr/go-stakenet/lnd"
)

// Balance has the current LND wallet in satoshi and the resolved currency with 1e8
type Balance struct {
	UnconfirmedSat int64   `json:"unconfirmed_sat"`
	ConfirmedSat   int64   `json:"confirmed_sat"`
	TotalSat       int64   `json:"total_sat"`
	Unconfirmed    float64 `json:"unconfirmed"`
	Confirmed      float64 `json:"confirmed"`
	Total          float64 `json:"total"`
}

// GetWalletBalance will connect to the lncli by passing the lnd directory and execute the walletbalance command
func GetWalletBalance(lncli string, directory string, host string, port int) (Balance, error) {
	if lncli == "" || directory == "" || port == 0 || host == "" {
		return Balance{}, fmt.Errorf("missing config for lncli / lnd config")
	}
	lnOptions := []string{fmt.Sprintf("--lnddir=%s", directory), "--no-macaroons", fmt.Sprintf("--rpcserver=%s:%d", host, port)}
	wb, err := lnd.GetWalletBalance(lncli, false, lnOptions...)
	if err != nil {
		return Balance{}, err
	}

	cb, err := strconv.ParseInt(wb.ConfirmedBalance, 10, 64)
	if err != nil {
		return Balance{}, err
	}
	tb, err := strconv.ParseInt(wb.TotalBalance, 10, 64)
	if err != nil {
		return Balance{}, err
	}
	ub, err := strconv.ParseInt(wb.UnconfirmedBalance, 10, 64)
	if err != nil {
		return Balance{}, err
	}
	ubf := float64(ub)
	cbf := float64(cb)
	tbf := float64(tb)
	b := Balance{
		UnconfirmedSat: ub,
		ConfirmedSat:   cb,
		TotalSat:       tb,
		Unconfirmed:    ubf / 1e8,
		Confirmed:      cbf / 1e8,
		Total:          tbf / 1e8,
	}
	return b, err
}

// GetListChannels will connect to the lncli by passing the lnd directory and execute the listchannels command
func GetListChannels(lncli string, directory string, host string, port int) (lnd.ListChannels, error) {
	if lncli == "" || directory == "" || port == 0 || host == "" {
		return lnd.ListChannels{}, fmt.Errorf("missing config for lncli / lnd config")
	}
	lnOptions := []string{fmt.Sprintf("--lnddir=%s", directory), "--no-macaroons", fmt.Sprintf("--rpcserver=%s:%d", host, port)}
	lc, err := lnd.GetListChannels(lncli, false, lnOptions...)
	if err != nil {
		return lnd.ListChannels{}, err
	}
	return lc, err
}

// GetPeers will connect to the lncli by passing the lnd directory and execute the listpeers command
func GetPeers(lncli string, directory string, host string, port int) (lnd.ListPeers, error) {
	if lncli == "" || directory == "" || port == 0 || host == "" {
		return lnd.ListPeers{}, fmt.Errorf("missing config for lncli / lnd config")
	}
	lnOptions := []string{fmt.Sprintf("--lnddir=%s", directory), "--no-macaroons", fmt.Sprintf("--rpcserver=%s:%d", host, port)}
	lp, err := lnd.GetListPeers(lncli, false, lnOptions...)
	if err != nil {
		return lnd.ListPeers{}, err
	}
	return lp, err
}

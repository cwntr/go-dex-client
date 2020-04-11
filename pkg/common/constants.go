package common

const (
	// PairXSNLTC trading pair for buy and sell XSN & LTC
	PairXSNLTC = "XSN_LTC"

	// PairXSNBTC trading pair for buy and sell XSN & BTC
	PairXSNBTC = "XSN_BTC"

	// PairXSNBTC trading pair for buy and sell LTC & BTC
	PairLTCBTC = "LTC_BTC"

	// CurrencyLTC implemented crypto-currency Litecoin - LTC
	CurrencyLTC = "LTC"

	// CurrencyXSN implemented crypto-currency Stakenet - XSN
	CurrencyXSN = "XSN"

	// CurrencyBTC implemented crypto-currency Bitcoin - BTC
	CurrencyBTC = "BTC"
)

// GetPairs func will return all active trading pairs
func GetPairs() []string {
	return []string{PairXSNLTC, PairXSNBTC, PairLTCBTC}
}

// GetPairs func will return all active currencies
func GetCurrencies() []string {
	return []string{CurrencyLTC, CurrencyXSN, CurrencyBTC}
}

// IsValidTrading func will evaluate if the passed pairId belongs to the list of active trading pairs or not
func IsValidTradingPair(pair string) bool {
	isFound := false
	for _, c := range GetCurrencies() {
		if pair == c {
			isFound = true
		}
	}
	return isFound
}

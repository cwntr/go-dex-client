package trading

import (
	"fmt"
	"github.com/cwntr/go-dex-client/pkg/common"
)

// Connection struct holds connection data
type Connection struct {
	Host string
	Port int
}

// Format func will format in a common host:port scheme
func (c *Connection) Format() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// LNDConfig holds mandatory information how to connection / authenticate with a LND
type LNDConfig struct {
	TLSPaths    map[string]string     // map[LTC]-> /path/to/lnd_ltc.tls.cert
	Certs       map[string]string     // map[LTC]-> -----BEGIN CERTIFICATE-----...
	Connections map[string]Connection // map[LTC]{Host: localhost, Port: 10001}
}

// NewConfig func will create a config with initialized maps
func NewConfig() LNDConfig {
	cfg := LNDConfig{}
	cfg.TLSPaths = make(map[string]string, 0)
	cfg.Certs = make(map[string]string, 0)
	cfg.Connections = make(map[string]Connection, 0)
	return cfg
}

// IsEmpty func checks whether there are certificate paths already specified
func (c *LNDConfig) IsEmpty() bool {
	return len(c.TLSPaths) == 0
}

// Add will add a new currency config to the LNDConfig
func (c *LNDConfig) Add(currency string, certPath string, host string, port int) error {
	if !ValidateCurrency(currency) {
		return fmt.Errorf("currency '%s' not allowed", currency)
	}
	if !fileExists(certPath) {
		return fmt.Errorf("certPath '%s' file does not exist", certPath)
	}
	if host == "" {
		return fmt.Errorf("host is empty")
	}
	if port == 0 {
		return fmt.Errorf("port is empty")
	}

	c.TLSPaths[currency] = certPath
	c.Connections[currency] = Connection{Host: host, Port: port}
	return nil
}

// ValidateCurrency checks if the given currency is allowed or not
func ValidateCurrency(currency string) bool {
	switch currency {
	case common.CurrencyLTC:
		return true
	case common.CurrencyBTC:
		return true
	case common.CurrencyXSN:
		return true
	}
	return false
}

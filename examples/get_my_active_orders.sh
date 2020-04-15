#!/usr/bin/env bash

# Get active placed orders by trading_pair XSN_LTC
curl http://localhost:9999/orders/XSN_LTC

# Get active placed orders by trading_pair XSN_BTC
curl http://localhost:9999/orders/XSN_BTC

# Get active placed orders by trading_pair LTC_BTC
curl http://localhost:9999/orders/LTC_BTC
#!/usr/bin/env bash

# XSN_LTC place 10 SELL orders with fixed quantity (funding of 0.1 XSN) with a price increase of 1 sat
# single means -> only one order placed, no follow up logic
curl -X POST \
  http://localhost:9999/orders \
  -H 'content-type: application/json' \
  -d '[{
    "side": "sell",
    "tradingPair": "XSN_LTC",
    "priceRangeStart": 101230,
    "priceRangeEnd": 101250,
    "priceRangeStepSize": 1,
    "fixedFunding": 10000000,
    "type": "single"
}]'

# XSN_LTC place 10 BUY orders with fixed quantity (funding of 0.0002 LTC) with a price increase of 1 sat
# IMPORTANT: buy orders need to have the funds converted to their base currency -> in this case its LTC
# single means -> only one order placed, no follow up logic
curl -X POST \
  http://localhost:9999/orders \
  -H 'content-type: application/json' \
  -d '[{
    "side": "buy",
    "tradingPair": "XSN_LTC",
    "priceRangeStart": 1012215,
    "priceRangeEnd": 1012229,
    "priceRangeStepSize": 1,
    "fixedFunding": 20000,
    "type": "single"
}]'


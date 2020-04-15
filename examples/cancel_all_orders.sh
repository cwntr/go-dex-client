#!/usr/bin/env bash

#cancel orders by orderId for XSN_LTC
curl -X POST http://localhost:9999/orders/cancel \
  -d '{
	"cancelTradingPairs": [
		{
			"tradingPair": "XSN_LTC",
			"orderIds": ["7a5f862c-9651-472f-aa4f-ac91c221a0c6"]
		}
	]
}'

#cancel all orders for XSN_LTC
curl -X POST \
  http://localhost:9999/orders/cancel \
  -d '{
	"cancelTradingPairs": [
		{
			"tradingPair": "XSN_LTC",
			"deleteAll": true
		}
	]
}'

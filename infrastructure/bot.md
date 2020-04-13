## Simple DEX Client installation

This "bot" will not place any orders or perform any actions but connecting to the API to check whether the infrastructure is setup correctly.
The actual logic how to place orders / active trading needs to be done by your customized implementation.

**Mandatory:** check out [infrastructure guide](infrastructure.md) to get your VM installed with the required components. 

You could run the example of a DexAPI Client following the next steps described in [client installation](bot.md)
The following "bot" is just a simple application to show how to deal with the API. It has basic implementations which should only be considered a starting point. All the trading / algorithms are left up for development.

`cd ~/bot`

`wget https://github.com/cwntr/go-dex-client/releases/download/v1.0.0/bot`

`chmod +x bot`

#### copy tls to local bot path

`cp ~/.lnd_xsn/tls.cert ~/bot/certs/xsn.cert`

`cp ~/.lnd_ltc/tls.cert ~/bot/certs/ltc.cert`

`cp ~/.lnd_btc/tls.cert ~/bot/certs/btc.cert`


##### configure the bot

`touch cfg.json`

`nano cfg.json` -default:

```
{
  "botCfg": {
    "host":"localhost",
    "port":9999,
    "lnCLIPath": "/home/ubuntu/lnds/lncli",
    "jobInterval": "5s",
    "logLevel": "debug",
    "orderLimit": 10000
  },
  "lssdConfig": {
    "host": "",
    "port": 50051,
    "timeout": "500s"
  },
  "xsnLNDConfig": {
    "lndDir": "/home/ubuntu/.lnd_xsn/",
    "certPath":"certs/xsn.cert",
    "host": "localhost",
    "port": 10003,
    "hubPeers": [
      {"remoteKey": "0396ca2f7cec03d3d179464acd57b4e6eabebb5f201705fa56e83363e3ccc622bb", "address": "134.209.164.91:11384"},
      {"remoteKey": "03bc3a97ffad197796fc2ea99fc63131b2fd6158992f174860c696af9f215b5cf1", "address": "134.209.164.91:21384"}
    ]
  },
  "ltcLNDConfig": {
    "lndDir":"/home/ubuntu/.lnd_ltc/",
    "certPath":"certs/ltc.cert",
    "host": "localhost",
    "port": 10001,
    "hubPeers": [
      {"remoteKey": "0375e7d882b442785aa697d57c3ed3aef523eb2743193389bd205f9ae0c609e6f3", "address": "134.209.164.91:11002"},
      {"remoteKey": "0211eeda84950d7078aa62383c7b91def5cf6c5bb52d209a324cda0482dbfbe4d2", "address": "134.209.164.91:21002"}
    ]
  },
  "btcLNDConfig": {
    "lndDir":"/home/ubuntu/.lnd_btc/",
    "certPath":"certs/btc.cert",
    "host": "localhost",
    "port": 10002,
    "hubPeers" : [
      {"remoteKey": "03757b80302c8dfe38a127c252700ec3052e5168a7ec6ba183cdab2ac7adad3910", "address":"134.209.164.91:11000"},
      {"remoteKey": "02bfe54c7b2ce6f737f0074062a2f2aaf855f81741474c05fd4836a33595960e18", "address":"134.209.164.91:21000"}
    ]
  }
}
```

## 17.) Start the client

##### Start the mandatory services that the client can operate, if you not done yet:

1. `sudo systemctl start lnd_xsn`

2. `sudo systemctl start lnd_ltc`

3. `sudo systemctl start lnd_btc`

4. `sudo systemctl start lssd`

##### Actual bot start 
`sudo systemctl start bot`

##### Stop bot
`sudo systemctl stop bot`

# Guide #5 - Setup Lightning channels

[< back to Guide 4 - LSSD](guide_4_lssd.md)

Now we do have 3x LN Nodes (BTC, LTC, XSN) but none of them is able to do any payment yet. We need to make channels to the hub manually and the hub operator will also then need to make channels back to you.

## Basic flow
1) fund your LN wallet
2) connect to the hub you want to make channels with
3) open a channel from your node with a local amount to the hub
4) tell the hub operator your channel capacity and with the channel id that he can match it and open back a channel
5) check that you have at least 2 channels to hub (one with local amount, one with equal amount of remote amount)

### Stakenet hub nodes

Find all the Stakenet hub nodes here:

https://github.com/X9Developers/DexAPI/blob/master/LNDCONFIGURATION.md#btc-hub-nodes

BTC only uses 1 hub node (the 2nd can be ignored)


## Example flow with XSN

#### 1) fund your LN wallet

Check your current balance (on-chain)

`lnxsn walletbalance` this will show you the balance that arrived with funds you send via the core chain (xsnd)

Probably so far everything will be zero amounts. So to get started, you will need to send some XSN to your XSN LND wallet. 

Firstly, you need to create a new address:

`lnxsn newaddress p2wkh`

This will give you an address you can fund on-chain. Beware: This is address is in a bench32 format, make sure the wallet you are sending the funds with is compatible with this format. (e.g. Coinomi Wallet would be an option for BTC & LTC) for XSN you can send if from the Core wallet.


#### 2) connect to the hub you want to make channels with

`lnxsn connect 0396ca2f7cec03d3d179464acd57b4e6eabebb5f201705fa56e83363e3ccc622bb@178.128.97.48:11384`

you can check with `lnxsn listpeers` if you see that node key in your connected peers


#### 3) open a channel from your node with a local amount to the hub

In this case you'll open a channel to the XSN Hub having a local balance of 5 XSN. You need to wait until the Hub opens back a channel to you to perform actual swaps.

`lnxsn openchannel --local_amt=500000000 --node_key=0396ca2f7cec03d3d179464acd57b4e6eabebb5f201705fa56e83363e3ccc622bb`

you can check the status / progress of your recent open channel with

`lnxsn pendingchannels` (visible here if its waiting for the minimum confirmations)

`lnxsn listchannels` (visible here if the funding transaction has enough confirmations, not visible in pendingchannels anymore)

#### 4) tell the hub operator you channel capacity and with the channel id that he can match it and open back a channel

Join us on [Discord](https://discord.gg/cyF5yCA) and get in contact with the Stakenet team. 
You just need to provide information like channel-id and capacity you used for opening that channel, you will get a channel back with same capacity.
You need to wait until the Hub opens back a channel to you, to perform actual swaps.

#### 5) check that you have at least 2 channels to hub (one with local amount, one with equal amount of remote amount)

If you got the ok that it was opened back to you, just regularly check with `listchannels` if you have channels that have a remote amount and are connected to the same remote node pubkey.


## Repeat everything for BTC and LTC 

**Important:** Make sure to always use the shortcut `lnltc` or `lnbtc` when working on the other setups otherwise you may lose funds.

Just follow the same steps as for XSN also for BTC and LTC. Please consider that for BTC there will be a longer period to get the all the minimum confirmations for your channel funding transaction. Especially with lower amounts this may take several hours.

[ > to next section (Guide 6 - Dex Clients)](guide_6_dex_clients.md)












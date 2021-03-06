# Infrastructure guide to run a DexAPI Client 

In order to get the client running, you need to install some mandatory components. This guide will follow you through the full installation process with step by step commands to be executed.

VM requirements:
| Minimum  | Recommended |
| ------------- | ------------- |
| 4x vCPUs  | 8x vCPUs |
| 8 GB of RAM | 16 GB of RAM  |
| 500 GB HDD | 600+ GB HDD |
 
to get the components (fully synchronized BTC, LTC and XSN chains and 3x Lightning Network Daemon) running.

### Components Overview
![alt text](components.png)


So log on to your Ubuntu VM and get started with the following steps. It describes the installation for a user named "ubuntu".

**Remark**: if decide to install everything as a `root` user (not recommended) you need to do make sure to dont use any /home/* directories and always keep in mind you are running everything on `/root` directory.

## 1.) Create basic folder structure:
`mkdir ~/bot`

`mkdir ~/coins`

`mkdir ~/lnds`

`mkdir ~/lssd`

## 2.) Install coins
`cd ~/coins`


#### XSN 
`wget https://github.com/X9Developers/XSN/releases/download/v1.0.21/xsn-1.0.21-x86_64-linux-gnu.tar.gz`

`tar xvzf xsn-1.0.21-x86_64-linux-gnu.tar.gz`

`rm xsn-1.0.21-x86_64-linux-gnu.tar.gz`

`mv xsn-1.0.21 xsn`

#### LTC
`wget https://download.litecoin.org/litecoin-0.17.1/linux/litecoin-0.17.1-x86_64-linux-gnu.tar.gz`

`tar xvzf litecoin-0.17.1-x86_64-linux-gnu.tar.gz`

`rm litecoin-0.17.1-x86_64-linux-gnu.tar.gz`

`mv litecoin-0.17.1 litecoin`

#### BTC
`wget https://bitcoin.org/bin/bitcoin-core-0.19.1/bitcoin-0.19.1-x86_64-linux-gnu.tar.gz`

`tar xvzf bitcoin-0.19.1-x86_64-linux-gnu.tar.gz`

`rm bitcoin-0.19.1-x86_64-linux-gnu.tar.gz`

`mv bitcoin-0.19.1 bitcoin`

## 3.) Start daemons in background and let them sync

You could also download the respective bootstrap but for simplicity lets just start the daemons and let them sync in the background and proceed with further steps.
 
#### XSN
`cd ~/coins/xsn/bin`

`./xsnd -daemon`

#### LTC
`cd ~/coins/litecoin/bin`

`./litecoind -daemon`

#### BTC
`cd ~/coins/bitcoin/bin`

`./bitcoind -daemon`


## 4.) Check synchronization status
The chains are fulled synchronized when the "blocks" count do match the "headers" count.
You can proceed with the next steps already and let it sync in background.
 
#### XSN
`cd ~/coins/xsn/bin`

`./xsn-cli getblockchaininfo`

#### LTC
`cd ~/coins/litecoin/bin`

`./litecoin-cli getblockchaininfo`

#### BTC
`cd ~/coins/bitcoin/bin`

`./bitcoin-cli getblockchaininfo`


## 5.) Install zmq (https://zeromq.org/download/)
`sudo apt-get install libzmq3-dev`

## 6.) Install unzip
`sudo apt-get install unzip`

## 7.) Download lnd's
`cd ~/lnds`

`wget https://github.com/X9Developers/DexAPI/releases/download/latest/lnds_0.8.2.5.zip`

`unzip lnds_0.8.2.5.zip`

`rm lnds_0.8.2.5.zip`

## 8.) Create basic LND folder structure:

Just start the lnd_* binaries once that they will create you the folder structures. You will see an error when executing it because no configuration is in place yet.
So just execute the binary and terminate (Ctrl+C) it as soon as you see some error / it hangs.

`cd ~/lnds`

`./lnd_ltc --lnddir=/home/ubuntu/lnds/.lnd_ltc`

`./lnd_btc --lnddir=/home/ubuntu/lnds/.lnd_btc`

`./lnd_xsn --lnddir=/home/ubuntu/lnds/.lnd_xsn`

## 9.) Configure XSN Lightning Daemon:

**Important**: noseedbackup=1 means a default password will be used for wallet encryption -> this should only be used for testing purposes. If you want to specify your own password for wallet encryption, use noseedbackup=0 and later on at step 11.)  and you can use `./lncli create` to also see the backup seed.

`cd ~/.lnd_xsn`

`touch lnd.conf`

`nano lnd.conf` paste:

```
noseedbackup=0
rpclisten=localhost:10003
listen=0.0.0.0:8403
restlisten=127.0.0.1:8083
nobootstrap=1
xsncoin.active=1
xsncoin.mainnet=1
xsncoin.defaultchanconfs=6
xsncoin.node=xsnd
xsnd.rpcuser=XSNDUSER
xsnd.rpcpass=XSNDPASSWORD123123
xsnd.zmqpubrawblock=tcp://127.0.0.1:28332
xsnd.zmqpubrawtx=tcp://127.0.0.1:28333
debuglevel=debug
maxpendingchannels=50
chan-enable-timeout=1m
max-cltv-expiry=10080
maxlogfiles=10
```


## 10.) Configure XSND - XSN Core wallet
`cd ~/coins/xsn/bin`

`./xsn-cli stop`

`cd ~/.xsncore/`

`touch xsn.conf`

`nano xsn.conf` -paste:

```
rpcuser=XSNDUSER
rpcpassword=XSNDPASSWORD123123
rpcallowip=127.0.0.1
listen=1
server=1
daemon=1
maxconnections=264
zmqpubrawblock=tcp://127.0.0.1:28332
zmqpubrawtx=tcp://127.0.0.1:28333
txindex=1
```

`cd ~/coins/xsn/bin`

`./xsnd`

[wait ~30s] and verify the auth is working by executing
`./xsn-cli getblockchaininfo`

which should give you once again the "blocks" and "headers" information


## 10.) Download lncli 
`cd ~/lnds`

`wget https://github.com/X9Developers/DexAPI/releases/download/v2020.01.23/lncli`

`chmod +x lncli`


## 11.) Starting XSN Lightning Node
###### Preparation:
- Add the all `*.service` files from [infratructure](systemd) to your `systemd` (/etc/systemd/system) with sudo permission, make sure to properly replace the fields of: `User`, `ExecStart` and `WorkingDirectory` to your user's name and the path is matching to locate the binaries.
- Add the shortcuts that will make your life easier operating the different lnd's.
  - `sudo nano ~/.bashrc` and scroll to end of file and paste the content of [bash profile](bash_profile)
  - once the file is saved, simple execute `bash` in the command line and its updated. Which allows you to use shortcut like 
    - `lnxsn <command>` which will only have an effect to the XSN LND
    - `lnltc <command>` which will only have an effect to the LTC LND
    - `lnbtc <command>` which will only have an effect to the BTC LND



Starting the lnd_xsn with:

`sudo systemctl start lnd_xsn`

Check if it's working (should be active state):

`sudo systemctl status lnd_xsn`

If its not working yet, you can execute the following command to get further details from the logs:

`sudo journalctl -f -u lnd_xsn`

##### If you configured your LND wallet with nobackupseed=0 in lnd.conf:
You need to consider the following things:
1) First time starting the LND: `lnxsn create` -> this will go through the encryption process where you set your encryption password and also the backup seed will be outputted which you can store somewhere safe.
2) Once you enabled the encryption, you need to do `lnxsn unlock` every time the LND has been restarted because it will be locked by default after restart.

##### If you configured your LND wallet with nobackupseed=1 in lnd.conf:
1) No need to create a further encryption because it is already encrypted with a default password 
2) No need to "unlock" the wallet upon restart of the LND

To check if the LND is active and unlocked try a command and check if you get proper outputs, e.g `walletbalance` command:

`lnxsn walletbalance`

Connect to the XSN Lightning Hubs:

`lnxsn connect 0396ca2f7cec03d3d179464acd57b4e6eabebb5f201705fa56e83363e3ccc622bb@134.209.164.91:11384`

`lnxsn connect 03bc3a97ffad197796fc2ea99fc63131b2fd6158992f174860c696af9f215b5cf1@134.209.164.91:21384`

[wait a few seconds]
Check whether it's connected properly to the network by outputting the network graph.

`lnxsn describegraph`

## 12.) Fund your XSN Node

`lnxsn newaddress p2wkh`
 
This will give you an address you can fund on-chain. Beware: This is address is in a bench32 format, make sure the wallet you are sending the funds with is compatible with this format. 
(e.g. Coinomi Wallet would be an option for BTC & LTC) for XSN you can send if from the Core wallet.

## 13.) Open a channel to a XSN Hub

In this case you'll open a channel to the XSN Hub having a local balance of 5 XSN. You need to wait until the Hub opens back a channel to you, to perform actual swaps.

`lnxsn openchannel --local_amt=500000000 --node_key=0396ca2f7cec03d3d179464acd57b4e6eabebb5f201705fa56e83363e3ccc622bb`

[wait ~5m]

`lnxsn listchannels` check that you have 2 channels: 1 channel filled with local balance, 1 channel filled with remote balance

**From this point on, you are done with setting up everything for XSN. What's to follow is to do the same steps 9.) - 13.) for LTC and BTC.**

**_Important_**: Make sure to always use the shortcut `lnltc` or `lnbtc` when working on the other setups otherwise you may lose funds.

## 14.) LTC: same steps as 9.) -13.) with different configs:

##### LTC LND config:
```
noseedbackup=0
rpclisten=localhost:10001
listen=0.0.0.0:8401
restlisten=127.0.0.1:8081
nobootstrap=1
litecoin.active=1
litecoin.mainnet=1
litecoin.defaultchanconfs=6
litecoin.node=litecoind
litecoind.rpcuser=LITECOINDUSER
litecoind.rpcpass=LITECOINDPASSWORD123123
litecoind.zmqpubrawblock=tcp://127.0.0.1:28336
litecoind.zmqpubrawtx=tcp://127.0.0.1:28337
debuglevel=debug
maxpendingchannels=50
chan-enable-timeout=1m
max-cltv-expiry=4032
maxlogfiles=10
```

##### Litecoind config:
```
rpcuser=LITECOINDUSER
rpcpassword=LITECOINDPASSWORD123123
rpcallowip=127.0.0.1
listen=1
server=1
daemon=1
maxconnections=264
zmqpubrawblock=tcp://127.0.0.1:28336
zmqpubrawtx=tcp://127.0.0.1:28337
txindex=1
```

Funding your LTC wallet: by generating a new address and sending LTC to it.

`lnltc newaddress p2wkh`

##### Connect to LTC Hub Peers:

`lnltc connect 0375e7d882b442785aa697d57c3ed3aef523eb2743193389bd205f9ae0c609e6f3@134.209.164.91:11002`

`lnltc connect 0211eeda84950d7078aa62383c7b91def5cf6c5bb52d209a324cda0482dbfbe4d2@134.209.164.91:21002`


##### Open Channel:
This would open a channel with a local amount of 0.015 LTC to the LTC HUB

`lnltc openchannel --local_amt=1500000 --node_key=0375e7d882b442785aa697d57c3ed3aef523eb2743193389bd205f9ae0c609e6f3`

## 15.) BTC: same steps as 9.) -13.) with different configs:
Have some patience with this one, it takes some time for the transactions to be confirmed.

##### BTC LND config:
```
noseedbackup=0
rpclisten=localhost:10002
listen=0.0.0.0:8402
restlisten=127.0.0.1:8082
nobootstrap=1
bitcoin.active=1
bitcoin.mainnet=1
bitcoin.defaultchanconfs=6
bitcoin.node=bitcoind
bitcoind.rpcuser=BITCOINDUSER
bitcoind.rpcpass=BITCOINDPASSWORD123123
bitcoind.zmqpubrawblock=tcp://127.0.0.1:28338
bitcoind.zmqpubrawtx=tcp://127.0.0.1:28339
debuglevel=debug
maxpendingchannels=50
chan-enable-timeout=1m
maxlogfiles=10
```

##### Bitcoind config:
```
rpcuser=BITCOINDUSER
rpcpassword=BITCOINDPASSWORD123123
rpcallowip=127.0.0.1
listen=1
server=1
daemon=1
maxconnections=264
zmqpubrawblock=tcp://127.0.0.1:28338
zmqpubrawtx=tcp://127.0.0.1:28339
txindex=1
```


Fund your BTC LND balance: by generating a new address and sending BTC to it.

`lnbtc newaddress p2wkh`

##### Connect to BTC Hub Peers:

`lnbtc connect 03757b80302c8dfe38a127c252700ec3052e5168a7ec6ba183cdab2ac7adad3910@134.209.164.91:11000`

`lnbtc connect 02bfe54c7b2ce6f737f0074062a2f2aaf855f81741474c05fd4836a33595960e18@134.209.164.91:21000`

##### Open Channel:
This will open a channel to the BTC hub with a local amount of 0.00092 BTC

`lnbtc openchannel --local_amt=92000 --node_key=03757b80302c8dfe38a127c252700ec3052e5168a7ec6ba183cdab2ac7adad3910`

## 15.) LSSD installation

`cd ~/lssd`

`wget https://github.com/X9Developers/DexAPI/releases/download/latest/lssd.zip`

`unzip lssd.zip`

`rm lssd.zip`

`sudo systemctl start lssd`

Check if it's working:

`sudo systemctl status lssd`

## Done

Congratulations, you just setup all mandatory things to get to work with the Stakenet DexAPI. 
What next is to run a client that consumes this API.
You could:
  - run the example of a DexAPI client following the next steps described in [client installation](bot.md) which allows you list, place & cancel orders
  - checkout the [scala implementation from X9](https://github.com/X9Developers/DexAPI/tree/master/trading-bot)
  - take this repo as a template and build some trading logic on top of it 
 

## Extra: LND config

Current XSN Lightning Network Explorer: https://ln.stakenet.info/

You can give your Lightning nodes a name (alias) and an unique color to easier find your node in network explorers.
You can do that by modifying the respective `lnd.conf` from the `.lnd_* directory` and `sudo systemctl stop lnd_xsn` & `sudo systemctl start lnd_xsn` start your LND node.

```
alias=stakenet.info
color=#9d03fc
```

For example with this alias and color the Lightning Node will look like that:

![alt text](node_alias.JPG)

---

If you find this page any useful, feel free to spare some coin:

**XSN:** XfiBSgdNJSVWsfZzKFa1HrcU8WgFCCUXGn

**ETH:** 0x28260E316Fa2491Af0005b7dF0085073B844B97b

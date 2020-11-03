# Guide #2 - Core blockchains install + running

[< back to Guide 1 - prerequisites](guide_1_prerequisites.md)

### Download and install coin backend's

Download all to a common folder e.g. `"~/coins"`

#### XSN 

`cd ~/coins`

`wget https://github.com/X9Developers/XSN/releases/download/v1.0.26/xsn-1.0.26-x86_64-linux-gnu.tar.gz`

`tar xvzf xsn-1.0.26-x86_64-linux-gnu.tar.gz`

`rm xsn-1.0.26-x86_64-linux-gnu.tar.gz`

`mv xsn-1.0.26 xsn`

#### LTC

`cd ~/coins`

`wget https://download.litecoin.org/litecoin-0.17.1/linux/litecoin-0.17.1-x86_64-linux-gnu.tar.gz`

`tar xvzf litecoin-0.17.1-x86_64-linux-gnu.tar.gz`

`rm litecoin-0.17.1-x86_64-linux-gnu.tar.gz`

`mv litecoin-0.17.1 litecoin`

#### BTC

`cd ~/coins`

`wget https://bitcoin.org/bin/bitcoin-core-0.19.1/bitcoin-0.19.1-x86_64-linux-gnu.tar.gz`

`tar xvzf bitcoin-0.19.1-x86_64-linux-gnu.tar.gz`

`rm bitcoin-0.19.1-x86_64-linux-gnu.tar.gz`

`mv bitcoin-0.19.1 bitcoin`


### Run the daemons

You only need to quickly run every single one of them for first time. They will create you basic folder structure for the blockchain
on your home folder. Kill the process when it seems to be running, to make proper config changes in next step.

#### XSN
`cd ~/coins/xsn/bin`

`./xsnd -daemon`
 
#### LTC
`cd ~/coins/litecoin/bin`

`./litecoind -daemon`
 
#### BTC
`cd ~/coins/bitcoin/bin`

`./bitcoind -daemon`

## Config adaptions

XSN -> `/root/.xsncore/`

LTC -> `/root/.litecoin/`

BTC -> `/root/.bitcoin/`

In these directories you can copy the `*.conf` files from  [/root](../installation_guide/root) to your respective
blockchain folder from home directory. Once the .conf file is updated, you must restart the daemon to apply it and check if everything works.

## Use systemd to easily start and stop the core blockchains

Copy 

`core_litecoin.service`

`core_bitcoin.service`

`core_xsn.service`
 
 files from [infrastructure](../installation_guide/etc/systemd/system) to your `/etc/systemd/system` directory on your server
 
## Check if all core blockchain daemons are stopped

`ps aux | grep bitcoin`

`ps aux | grep xsn`

`ps aux | grep litecoin` 


if you dont see any process running, its all good if not, please kill everything that was still open from manual starting earlier.

## Control daemons via systemd

#### Starting & Stopping Bitcoin core daemon:

`systemctl start core_bitcoin`

`systemctl stop core_bitcoin`  

#### Starting & Stopping Litecoin core daemon:

`systemctl start core_litecoin`

`systemctl stop core_litecoin`  

#### Starting & Stopping XSN core daemon:

`systemctl start core_xsn`

`systemctl stop core_xsn`


#### Troubleshooting
You see some error after starting a daemon? You can debug with `systemctl status core_bitcoin`. Probably a path is wrong or permissions are insufficient.

Also, you can directly observe the debug file of the daemon if there is more info, e.g. for bitcoin:
 
`tail -f ~/.bitcoin/debug.log`

Still no info? Check if your pruning works and you have enough disk space

`df -H` will list current disk space consumption. There should be a row with several GB total and check the % usage of it.


# All done?

Before continuing to next steps its important to check all is good for this section. You can check it like this:

### BTC

`/root/coins/bitcoin/bin/bitcoin-cli getblockchaininfo` -> "blocks" == "headers" -> all synced

Also "automatic_pruning": true and "pruneheight:" has a value which is several thousand blocks lower than current synced block.

`systemctl status core_bitcoin` -> status active

#### LTC

`/root/coins/litecoin/bin/litecoin-cli getblockchaininfo` -> "blocks" == "headers" -> all synced

`systemctl status core_litecoin` -> status active


#### XSN

`/root/coins/xsn/bin/xsn-cli getblockchaininfo` -> "blocks" == "headers" -> all synced

`systemctl status core_xsn` -> status active


If all of that is also reflected on your server, go ahead to next step, if not maybe check the Troubleshooting section.

[ > to next section (Guide 3 - Lightning Daemons)](guide_3_lightning_daemons.md)



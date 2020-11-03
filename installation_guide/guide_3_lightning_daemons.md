# Guide #3 - LNDs Installation + Running

[< back to Guide 2 - Core Chains](guide_2_core_chains.md)

## Download

Download the recent LND binaries

`cd ~/lnds`

`wget https://github.com/X9Developers/DexAPI/releases/download/0.4.0.4/lnd.tar.gz`

`sudo tar -xvf lnd.tar.gz`

`sudo rm lnd.tar.gz`

## LND wallet folders

Create folders that will contain your Lightning wallets

`mkdir ~/.lnd_btc`

`mkdir ~/.lnd_xsn`

`mkdir ~/.lnd_ltc`

## Configure your Lightning daemon to connect to the previously installed core Blockchains

Just copy every `lnd.conf` file from every folder from 

[/.lnd_btc](../installation_guide/root/.lnd_btc)
 
[/.lnd_ltc](../installation_guide/root/.lnd_ltc)

[/.lnd_xsn](../installation_guide/root/.lnd_xsn)

All 3x `lnd.conf` file need to be copied to their folders which you created in previous step. They have the config
how to connect to the core blockchain.

## Use lncli

Add shortcut functions to your bashrc. Your bashrc can be found at ~/.bashrc which probably already has some content in it.
Just edit this file and go to very end of the file and add the content from 

[/.bashrc](../installation_guide/root/.bashrc) 

Save the file and in your command line just type `bash` and enter. It will refresh the bash profile and if you type
`lnx` and press `tab` it will auto-complete it to `lnxsn` which is your shortcut for your LND XSN wallet.

Once that works you have now 3 shortcuts to control your LND wallets

`lnbtc` -> LND BTC wallet

`lnltc` -> LND LTC wallet

`lnxsn` -> LND XSN wallet

With this command you can check balances, channels, create invoices, send payments, send coins, etc. 

## Systemd for your LNDs

Copy 

`lnd_btc.service`

`lnd_ltc.service`

`lnd_xsn.service`
 
 files from [infratructure](../installation_guide/etc/systemd/system) to your `/etc/systemd/system` directory on your server
 
## Control daemons via systemd
 
####Starting & Stopping LND BTC

`systemctl start lnd_btc`

`systemctl stop lnd_btc`  

####Starting & Stopping LND LTC

`systemctl start lnd_ltc`

`systemctl stop lnd_ltc`  

####Starting & Stopping LND XSN:

`systemctl start lnd_xsn`

`systemctl stop lnd_xsn`


## All Running?

So start all of your LNDs with `systemctl start lnd_*` and make sure they are running. 

You can check with `ps aux | grep lnd_btc` if you see a process, or with `systemctl status lnd_btc` if its active.

Also, for debugging purposes you can observe the live logs with `tail -f /root/.lnd_btc/logs/bitcoin/mainnet/lnd.log`


## Create actual wallets
Now we have all the preparation to handle multiple LNDs. Lets create every single LND wallet.

**IMPORTANT**: the next commands will create you seeds which you will only see here once. If you plan to use the wallet actively with any funds
make sure to copy the printed seeds. Otherwise you may lose it all.

The following commands will open a small interactive dialog with the wallet where you can enter a passphrase to encrypt your wallet. Also make sure to keep it somewhere.

**IMPORTANT**: please make sure to also save your passphrases somewhere safe and keep them reachable. For every restart of the LND, you need to unlock it with the passphrase.

`lnxsn create`

`lnltc create`

`lnbtc create`

After all 3 dialogs you have the wallet created and Lightning can be setup.

## Restart everything

#### stop
`systemctl stop lnd_btc`

`systemctl stop lnd_xsn`

`systemctl stop lnd_ltc`

#### start
`systemctl start lnd_ltc`

`systemctl start lnd_xsn`

`systemctl start lnd_btc`

#### unlock

(may wait 10s) and have your passwords ready.

`lnxsn unlock` (+ enter your password in dialog)

`lnltc unlock` (+ enter your password in dialog)

`lnbtc unlock` (+ enter your password in dialog)

#### check

(may wait 10s) 

Lets check if all daemons are working fine. 

`lnxsn walletbalance`

`lnbtc walletbalance`

`lnltc walletbalance`

If you see `0` values for your balances, all is good. If not, probably wait a bit and try again. If it continues check again if your core wallets are running properly.

## Sync with network and build graph

Right now, you just setup the wallet but dont have any network data. Lets change that by connecting to fixed public nodes to download the graphs.
This may need up to 15 min but not much longer.

**BTC**

`lnbtc connect 03757b80302c8dfe38a127c252700ec3052e5168a7ec6ba183cdab2ac7adad3910@178.128.97.48:11000`
 

**LTC**

`lnltc connect 0375e7d882b442785aa697d57c3ed3aef523eb2743193389bd205f9ae0c609e6f3@178.128.97.48:11002`
 
`lnltc connect 0211eeda84950d7078aa62383c7b91def5cf6c5bb52d209a324cda0482dbfbe4d2@178.128.97.48:21002`

**XSN**

`lnxsn connect 0396ca2f7cec03d3d179464acd57b4e6eabebb5f201705fa56e83363e3ccc622bb@178.128.97.48:11384`
 
`lnxsn connect 0396ca2f7cec03d3d179464acd57b4e6eabebb5f201705fa56e83363e3ccc622bb@178.128.97.48:11384`


# All done?

Lets check if we can continue to next guide section. 

`lnxsn getinfo`

`lnbtc getinfo`

`lnltc getinfo`

for every of the 3 commands from above you must see

` "synced_to_chain": true,`
 
` "synced_to_graph": true,`


And also just double checking that all is still in `active`

`systemctl status lnd_btc`

`systemctl status lnd_ltc`

`systemctl status lnd_xsn`

Now continue to next guide section. If something is off, try restarting + unlocking the LNDs and also make sure your core wallets are running properly.

[ > to next section (Guide 4 - LSSD)](guide_4_lssd.md)

# Guide #4 - LSSD Install + Run

LSSD will require active LND interaction so make sure you have all the previous guide steps completed, otherwise nothing will work.


## Install
`cd ~/lssd`

`wget https://github.com/X9Developers/DexAPI/releases/download/latest/lssd.zip`

`unzip lssd.zip`

`rm lssd.zip`

## Systemd Configuration

Copy 

`lssd.service`

 files from [infratructure](../installation_guide/etc/systemd/system) to your `/etc/systemd/system` directory on your server

## Run LSSD
`sudo systemctl start lssd`

Optionally, you can ask the Stakenet team for orderbook api-key that will lower your order fee. Simply adjust in your 

`/etc/systemd/system/lssd.service` adapting this line with key.

`ExecStart=/usr/bin/env /root/lssd/AppRun --orderbookAPISecret XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX`


## Check if all is active and running

`sudo systemctl status lssd` -> should be `active`

Also, you can observe the current log files with 

`taillssd` which is an alias from the example `~/.bashrc` file

## Troubleshooting

Some errors while starting the `lssd`? Stop everything and manually start to check if an error is visible with

`sudo systemctl stop lssd`

`cd ~/lssd`

`./AppRun`

and observe what is happening here.

 

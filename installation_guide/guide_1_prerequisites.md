# Guide #1 -  Prerequisites

[< back to overview](README.md)

Important before you start:
- make sure your server os `Ubuntu 18.04+`, otherwise the `lssd` wont be running.
- this guide is written in a perspective as `root` user, if you want to have a different setup please consider to manually adapt the paths from the provided examples.
- the files provided are organized how it would be on located on your server (given default settings as root)

### 0.) update apt-get:
 
`sudo apt-get update`

### 1.) install zmq: 

`sudo apt-get install libzmq3-dev`

### 2.) install unzip: 

`sudo apt-get install unzip`

### 3.) folder structure

`mkdir ~/bot`

`mkdir ~/coins`

`mkdir ~/lnds`

`mkdir ~/lssd`

### 4.) install lncli

`cd ~/lnds`

`wget https://github.com/lightningnetwork/lnd/releases/download/v0.11.1-beta.rc3/lnd-linux-amd64-v0.11.1-beta.rc3.tar.gz`

`tar -xvf lnd-linux-amd64-v0.11.1-beta.rc3.tar.gz`

`rm lnd-linux-amd64-v0.11.1-beta.rc3.tar.gz`

`mv lnd-linux-amd64-v0.11.1-beta.rc3/lncli lncli`

`rm -r lnd-linux-amd64-v0.11.1-beta.rc3/`

`chmod +x lncli`

[ > to next section (Guide 2 - Core Chains)](guide_2_core_chains.md)

# go-dex-trading-bot

[![License](http://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/miguelmota/cwntr/go-crypto-tools/LICENSE.md)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/1ab794872bef48d59e09f8e3160d6326)](https://www.codacy.com/manual/cwntr/go-dex-client?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=cwntr/go-dex-client&amp;utm_campaign=Badge_Grade)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](#contributing)

A trading bot for Stakenet's XSN DexAPI (Decentralized Exchange) written in golang.

The official XSN DexAPI with a trading bot written in Scala can be found here: [**github.com/X9Developers/DexAPI**](https://github.com/X9Developers/DexAPI). It also has an extensive documentation how to set up all the required components on your machine which is mandatory to get this bot running, since this repository only provides an alternative trading bot implementation.

## Components
![alt text](infrastructure/components.png)

## Run the bot
 1. Download the **UNIX** `bot` binary for your from the release
 2. Create a directory `certs` and paste all the lnd's `tls.cert` files
 3. Copy the default `cfg.json` from the repository and modify it based on your setup.
 4. Execute the binary `./bot`
 
By default, it will print the infrastructure setup checks. You need to make sure to have all mandatory components for the XSN, LTC or BTC lnd's running and the channels to the Stakenet HUB are working fine and enough capacity is available.

### Examples to place orders, deal with the orderbook
To interact with the bot you can send HTTP requests to fetch orderbook data, place orders, cancel orders or to retrieve your wallet LND's balance. Have a look in `examples` directory that contains example curl requests that interact with your local web server.  

## Create stub via protoc (if not up-to-date)
Use the following link to install the prerequisites ([**https://grpc.io/docs/quickstart/go/**](https://grpc.io/docs/quickstart/go/)):

 1. install `protoc compiler` (3.6.1+) 
 2. install `protoc-gen-go` compiler plugin 

### Generate a stub by using the lssdrpc.proto file
The latest **lssdrpc.proto** file can be found on: 
[**github.com/X9Developers/DexAPI/releases**](https://github.com/X9Developers/DexAPI/releases)

This **lssdrpc.proto** has to be copied to the `lssdrpc` directory and the following commands have to be executed.

Go to the project root, execute the following command to generate a go client for the lssdrpc API:

 1. `cd lssdrpc/`
 2. `protoc -I . lssdrpc.proto --go_out=plugins=grpc:.`

which will output a **lssdrpc.rb.go** that has client and server connectors automatically generated.
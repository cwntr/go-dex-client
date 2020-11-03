function lnltc() {
    ~/lnds/lncli --lnddir=~/.lnd_ltc --macaroonpath=/root/.lnd_ltc/data/chain/litecoin/mainnet/admin.macaroon -rpcserver=localhost:10001 $1 $2 $3 $4
}

function lnxsn() {
    ~/lnds/lncli --lnddir=~/.lnd_xsn --macaroonpath=/root/.lnd_xsn/data/chain/xsncoin/mainnet/admin.macaroon -rpcserver=localhost:10003 $1 $2 $3 $4
}

function lnbtc() {
    ~/lnds/lncli --lnddir=~/.lnd_btc --macaroonpath=/root/.lnd_btc/data/chain/bitcoin/mainnet/admin.macaroon -rpcserver=localhost:10002 $1 $2 $3 $4
}

function taillssd() {
    sudo journalctl -f -u lssd
}
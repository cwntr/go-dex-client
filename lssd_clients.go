package main

import (
	"fmt"
	"github.com/cwntr/go-dex-client/lssdrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"time"
)

// createSwapClient creates a gRPC connection for the Swaps service
func createSwapClient() (lssdrpc.SwapsClient, *grpc.ClientConn) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "", cfg.LSSDConfig.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return lssdrpc.NewSwapsClient(conn), conn
}

// createSwapClient creates a gRPC connection for the Currency service
func createCurrencyClient() (lssdrpc.CurrenciesClient, *grpc.ClientConn) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "", cfg.LSSDConfig.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return lssdrpc.NewCurrenciesClient(conn), conn
}

// createTradingPairClient creates a gRPC connection for the Currency service
func createTradingPairClient() (lssdrpc.TradingPairsClient, *grpc.ClientConn) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "", cfg.LSSDConfig.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return lssdrpc.NewTradingPairsClient(conn), conn
}

// createOrdersClient creates a gRPC connection for the Order service
func createOrdersClient() (lssdrpc.OrdersClient, *grpc.ClientConn) {
	var conn *grpc.ClientConn
	params := keepalive.ClientParameters{}
	params.Timeout = time.Duration(time.Second * 350)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "", cfg.LSSDConfig.Port), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithKeepaliveParams(params))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return lssdrpc.NewOrdersClient(conn), conn
}

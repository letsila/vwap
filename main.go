package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/letsila/vwap/internal"
	"github.com/letsila/vwap/internal/vwap"
	"github.com/letsila/vwap/internal/websocket/exchanges/coinbase"
)

const (
	defaultTradingPairs = "BTC-USD,ETH-USD,ETH-BTC"
	defaultWindowSize   = 200
)

func main() {
	ctx := context.Background()

	tradingPairs := flag.String("trading-pairs", defaultTradingPairs, "trading pairs to subscribe to")

	wsURL := flag.String("ws-url", coinbase.DefaultURL, "coinbase websocket url")

	fmt.Printf("wsURL %s\n", *wsURL)

	windowSize := flag.Uint("window-size", defaultWindowSize, "window size")

	tradingPairsArr := strings.Split(*tradingPairs, ",")

	flag.Parse()

	// Intercepting shutdown signals.
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

		s := <-quit

		log.Printf("received signal: %s", s)

		os.Exit(0)
	}()

	ws, err := coinbase.NewClient(*wsURL)
	if err != nil {
		log.Fatal(err)
	}

	list, err := vwap.NewList([]vwap.DataPoint{}, *windowSize)
	if err != nil {
		log.Fatal(err)
	}

	service := internal.NewService(ws, tradingPairsArr, &list)

	err = service.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

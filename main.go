package main

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

var infuraKey = ""

func main() {

	c, err := ethclient.DialContext(context.Background(), infuraKey)
	if err != nil {
		log.Fatal("failed to connect to infura client")
	}

	defer c.Close()

	lastBlock, err := c.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	lastBlockNumber := lastBlock.Number()

	println(lastBlockNumber)
}

package main

import (
	"context"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	INFURA := os.Getenv("INFURA_ETH_SEP")

	c, err := ethclient.DialContext(context.Background(), INFURA)
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

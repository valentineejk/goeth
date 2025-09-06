package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/joho/godotenv/autoload"

	"github.com/ethereum/go-ethereum/ethclient"
)

var GANACHE = "HTTP://127.0.0.1:7545"

func main() {

	rpcUrl := os.Getenv("INFURA_ETH_SEP1")
	if rpcUrl == "" {
		rpcUrl = GANACHE
	}

	c, err := ethclient.DialContext(context.Background(), rpcUrl)
	if err != nil {
		log.Fatal("failed to connect to infura client")
	}

	defer c.Close()

	// Verify connection
	chainID, err := c.ChainID(context.Background())
	if err != nil {
		log.Fatalf("failed to get chain ID: %v", err)
	}
	fmt.Printf("Connected to chain ID: %s\n", chainID.String())

	//last block
	lastBlockNumber := getBlockNumber(c)
	fmt.Println("last block: ", lastBlockNumber)

	//addr
	//addr := common.HexToAddress("0xEef248FBe38b25657AeFeBCcafCF2C5Aa942fed4")

	//keys
	priKey, pubKey := createKEYs()
	fmt.Println("private key: ", priKey)
	fmt.Println("public key: ", pubKey)
	addr := common.HexToAddress(pubKey)

	//balance
	bal := getBalance(c, addr)
	fmt.Println("balance: ", bal)

	//secure wallet gen
	//addr2 := keyStoreFile("Mark01")
	//fmt.Println("secure address: ", addr2)

	location := "./icloud/UTC--2025-09-06T14-31-35.882458000Z--6720f1b9a6e326ee118a27e942858de3a51460ab"

	priKey2, pubKey2, addr2, err := decryptKeystore("Mark01", location)
	if err != nil {
		log.Printf("Error decrypting keystore: %v", err)
		return
	}
	fmt.Println("decrypted private key2: ", priKey2)
	fmt.Println("decrypted pub key2: ", pubKey2)
	fmt.Println("decrypted addr key2: ", addr2)

}

func getBlockNumber(c *ethclient.Client) *big.Int {
	lastBlock, err := c.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	lastBlockNumber := lastBlock.Number()

	return lastBlockNumber
}

func getBalance(c *ethclient.Client, addr common.Address) string {
	bal, err := c.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		log.Fatal(err)
	}

	ethBalance := weiToEth(bal)

	return ethBalance.String()
}

func weiToEth(wei *big.Int) *big.Float {
	// 1 ETH = 10^18 Wei
	eth := new(big.Float).SetInt(wei)
	return eth.Quo(eth, big.NewFloat(1e18))
}

func createKEYs() (string, string) {
	key, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(key)
	PriKey := hexutil.Encode(privateKeyBytes)

	address := crypto.PubkeyToAddress(key.PublicKey).Hex()

	return PriKey, address
}

//func keyStoreFile(password string) string {
//	key := keystore.NewKeyStore("./icloud", keystore.StandardScryptN, keystore.StandardScryptP)
//	addr, err := key.NewAccount(password)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return addr.Address.Hex()
//
//}

func decryptKeystore(password, dir string) (string, string, string, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return "", "", "", fmt.Errorf("keystore file not found: %s", dir)
	}

	file, err := os.ReadFile(dir)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read keystore file: %w", err)
	}

	key, err := keystore.DecryptKey(file, password)
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(key.PrivateKey)
	PriKey := hexutil.Encode(privateKeyBytes)

	pubKeyBytes := crypto.FromECDSAPub(&key.PrivateKey.PublicKey)
	pubKey := hexutil.Encode(pubKeyBytes)

	return PriKey, pubKey, key.Address.Hex(), nil
}

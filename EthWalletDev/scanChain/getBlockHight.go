package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ethAdder = "0x98DB50E04D81A09cD4e63e77558a2d7698ab60ee"
)

func getBlockHeight() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/fzgfj4QuLlNEyn2LrLZsseBAClGdMnyP")
	if err != nil {
		log.Fatal(err)
	}

	// Get the balance of an account
	account := common.HexToAddress(ethAdder)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Account balance: %d\n", balance) // 25893180161173005034

	// Get the latest known block
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Latest block: %d\n", block.Number().Uint64())
}
func getTransactionCount() {

	url := "https://eth-sepolia.g.alchemy.com/v2/fzgfj4QuLlNEyn2LrLZsseBAClGdMnyP"

	payload := strings.NewReader(fmt.Sprintf("{\"id\":1,\"jsonrpc\":\"2.0\",\"params\":[\"%s\",\"latest\"],\"method\":\"eth_getTransactionCount\"}", ethAdder))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))

}
func main() {
	getBlockHeight()
	getTransactionCount()
}

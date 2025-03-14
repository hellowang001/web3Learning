package sol_official

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"log"
)

func NewClient() *client.Client {
	//ENDPOINT := "https://rough-frequent-tent.SOLANA_MAINNET.quiknode.pro/1f710f95479f793bf564bba4ababea2c969137f6/"
	DEV := "https://solana-devnet.g.alchemy.com/v2/wqZxT7UnY6AgrzV42CtGgGQ7ZGM-UrTq"

	c := client.NewClient(DEV)
	resp, err := c.GetVersion(context.TODO())
	if err != nil {
		log.Fatalf("failed to version info, err: %v", err)
	}

	fmt.Println("version", resp.SolanaCore)
	return c
}

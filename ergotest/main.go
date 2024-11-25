package main

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	fmt.Printf("hhhh")
	priv, err := crypto.GenerateKey()
	if err != nil {
		fmt.Printf("err:%v", err)
	}
	fmt.Printf("priv:%v", hex.EncodeToString(priv))
}

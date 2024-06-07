package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
)

// 恢复一个账户出来
func privateKeyToAccount(PrivateKeyStr string) string {
	fmt.Println("ttttt")
	// 将私钥字符串转化成字节
	privateKey, _ := crypto.HexToECDSA(PrivateKeyStr)

	// 从私钥拿到公钥
	publicKey := privateKey.PublicKey

	// 直接用crypto里的方法将公钥推出地址
	ethAddre := crypto.PubkeyToAddress(publicKey).Hex() // Hex是16进制转字符串
	fmt.Println("ethAddre", ethAddre)                   // ethAddre 0x38B59D6D4ef6A4991926Cf04c7c2092a0E86140F
	return ethAddre
}

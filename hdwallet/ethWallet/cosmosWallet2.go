package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"log"

	"golang.org/x/crypto/ripemd160"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/go-bip39"
)

func main() {
	fmt.Println("开始")

	mnemonic := "library enter must rural laptop wrestle slot beyond guide nominee scissors nerve"
	// 通过助记词生成种子Seed
	seed := bip39.NewSeed(mnemonic, "") // password盐值不要加
	fmt.Println("seed", hex.EncodeToString(seed))

	// 接下来就是将种子恢复出主私钥 masterKey 这里进入到了BIP-32了
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("masterKey", masterKey.String())

	// 现在要派生出对应Cosmos的子私钥，遵循BIP-44
	// 完成派生,完成对应path参数 m / purpose' / coin_type' / account' / change / address_index
	purpose, err := masterKey.Child(44 + hdkeychain.HardenedKeyStart)
	if err != nil {
		log.Fatal(err)
	}
	coinType, err := purpose.Child(118 + hdkeychain.HardenedKeyStart) // 118 是 Cosmos 的币种标识符
	if err != nil {
		log.Fatal(err)
	}
	account, err := coinType.Child(0 + hdkeychain.HardenedKeyStart)
	if err != nil {
		log.Fatal(err)
	}
	change, err := account.Child(0)
	if err != nil {
		log.Fatal(err)
	}
	addressIndex, err := change.Child(0)
	if err != nil {
		log.Fatal(err)
	}

	// 获取子私钥
	privateKey, err := addressIndex.ECPrivKey()
	if err != nil {
		log.Fatal(err)
	}
	cosmosPrivateKey := hex.EncodeToString(privateKey.Serialize())
	fmt.Println("privateKey", cosmosPrivateKey)

	// 从私钥生成公钥
	publicKey := privateKey.PubKey()
	pubKeyBytes := publicKey.SerializeCompressed()
	fmt.Println("pubKeyBytes str ", hex.EncodeToString(pubKeyBytes))

	fmt.Println("-----------------------------------")
	// 对公钥进行 SHA-256 哈希处理
	sha256Hasher := sha256.New()
	sha256Hasher.Write(pubKeyBytes)
	sha256Hash := sha256Hasher.Sum(nil)

	// 对 SHA-256 哈希结果进行 RIPEMD-160 哈希处理
	ripemd160Hasher := ripemd160.New()
	ripemd160Hasher.Write(sha256Hash)
	pubKeyHash := ripemd160Hasher.Sum(nil)

	// Bech32 编码
	bech32Addr, err := bech32.ConvertAndEncode("cosmos", pubKeyHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("cosmosAddress", bech32Addr)
}

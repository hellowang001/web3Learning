package main

import (
	"fmt"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func main() {
	// 由熵源生成助记词
	// @参数 128 => 12个单词
	// @参数 256 => 24个单词
	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	fmt.Println("助记词：", mnemonic)

	// 由助记词生成种子(Seed)
	seed := bip39.NewSeed(mnemonic, "salt")
	fmt.Println("seed", seed)

	// 由种子生成主账户私钥
	masterKey, _ := bip32.NewMasterKey(seed)
	fmt.Println("masterKey", masterKey)

	pubkey := masterKey.PublicKey()
	fmt.Println("pubkey", pubkey)
	// 解压缩公钥
	pubKey1, _ := crypto.DecompressPubkey(PubKeyToChild.Key)
	// 生成子账户地址
	addr1 := crypto.PubkeyToAddress(*pubKey1)

}

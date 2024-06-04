package main

import (
	"fmt"
)

func main() {
	// 生成随机熵，熵的长度可以是128位、256位等，必须是32位的倍数, 16 bytes = 128 bits / 8
	//entropy, err := bip39.NewEntropy(128)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// 使用熵生成助记词
	//mnemonic, err := bip39.NewMnemonic(entropy)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println("Generated Mnemonic:", mnemonic)
	//
	//// 验证助记词
	//if !bip39.IsMnemonicValid(mnemonic) {
	//	log.Fatal("verification failed")
	//}
	//
	//// 通过助记词生成种子
	//seed := bip39.NewSeed(mnemonic, "")
	//fmt.Println("Seed:", seed)
	fmt.Println("mnemonictest22")
}

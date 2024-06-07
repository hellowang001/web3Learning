package main

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func main() {
	// 由熵源生成助记词
	// @参数 128 => 12个单词
	// @参数 256 => 24个单词
	entropy, _ := bip39.NewEntropy(128)
	//fmt.Println("entropy", entropy)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	//mnemonic := "distance cool level useful print pet yard measure paddle bird solve inner"
	//mnemonic := "arctic prefer palm meat fruit love stick achieve glance jelly clerk skirt"
	fmt.Println("助记词：", mnemonic)
	fmt.Printf("助记词：type: %T\n", mnemonic)

	// 由助记词生成种子(Seed)
	//seed := bip39.NewSeed(mnemonic, "salt")
	seed := bip39.NewSeed(mnemonic, "")
	fmt.Println("seed", seed)
	fmt.Printf("type: %T\n", seed)
	//hex := fmt.Sprintf("%X", seed)
	//fmt.Println(hex)

	// 由种子生成主账户私钥
	masterKey, _ := bip32.NewMasterKey(seed)
	fmt.Println("这里是主私钥masterKey:", masterKey) // 这里是主私钥

	// 由主账户私钥生成子账户私钥
	// @参数 索引号
	//childKey1, _ := masterKey.NewChildKey(1)
	//fmt.Println("childKey1:", childKey1)

	// 派生第二个子账户
	//childKey2, _ := masterKey.NewChildKey(2)
	//fmt.Println("childKey2:", childKey2)

	// 用主账户公钥 派生 子账户公钥（没有私钥）
	//publicKey := masterKey.PublicKey() // 这里是拿到主账户公钥，待会下面会派生
	//fmt.Println("publicKey:", publicKey)

	//PubKeyToChild, _ := publicKey.NewChildKey(1) // 通过主账户公钥派生出子账户公钥
	//fmt.Println("PubKeyToChild:", PubKeyToChild)

	// 用主账户私钥，派生出子账户私钥，再生产子账户公钥，然后判断一下是否一致
	//childKey01, _ := masterKey.NewChildKey(1)
	//fmt.Println("childKey01:", childKey01)
	//publicKey1 := childKey01.PublicKey()
	//fmt.Println("publicKey3:", publicKey1)
	//fmt.Println(bytes.Equal(PubKeyToChild.Key, publicKey1.Key)) // 返回的ture

	// 由公钥推出地址（解压缩的过程，就是转化然后去头去尾的）
	// 先解压缩公钥，再去推地址
	//pubKey1, _ := crypto.DecompressPubkey(PubKeyToChild.Key)
	// 生成子账户地址
	//addre1 := crypto.PubkeyToAddress(*pubKey1)
	//fmt.Println("addre1:", addre1)
	// 以太坊的币种类型是60
	// FirstHardenedChild = uint32(0x80000000) 是一个常量
	// 以路径（path: "m/44'/60'/0'/0/0"）为例
	key, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)  // 这一步是在做强化派生， 对应 purpose' 44是符合BIP-44标准
	key, _ = key.NewChildKey(bip32.FirstHardenedChild + uint32(60)) // 还是强化派生，对应 coin_type类型60'
	key, _ = key.NewChildKey(bip32.FirstHardenedChild + uint32(0))  // 还是强化派生，对应account
	key, _ = key.NewChildKey(uint32(0))                             // 常规派生，对应 change
	key, _ = key.NewChildKey(uint32(0))                             // 常规派生，对应的 address_index

	// path已经完毕，生成地址
	ethPublicKey, _ := crypto.DecompressPubkey(key.PublicKey().Key)
	fmt.Println("ethPublicKey:", ethPublicKey)
	ethAddre := crypto.PubkeyToAddress(*ethPublicKey).Hex()
	fmt.Println("ethAddre:", ethAddre)
	// 尝试打印一下这个key的私钥
	fmt.Println("key:", key.Key)                 // 这个时候你直接打印Key是十进制的字节切片
	privateKeyHex := hex.EncodeToString(key.Key) // 这里在做的就是转化成16进制的字符，
	fmt.Println("key:", privateKeyHex)
}

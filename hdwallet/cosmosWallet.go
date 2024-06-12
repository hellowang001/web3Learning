package hdwallet

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ripemd160"
)

// 生成一个以太坊的钱包
func cosmosWalletPrivate() {
	fmt.Println("开始")
	// 首先你可以生成一个随机熵，熵源助记词是BIP-39，
	//entropy, _ := bip39.NewEntropy(128)
	//fmt.Println("entroy:", entropy)
	// 通过熵源生成助记词  ==> 注意，不一定要有助记词才有种子，只是助记词方便备份，可以转成种子，你要直接由种子也行，但不好记
	//mnemonic, _ := bip39.NewMnemonic(entropy)
	mnemonic := "rural neither robot good glove bracket fee harsh bird iron segment rug"
	fmt.Println("mnemonic:", mnemonic)
	// 通过助记词生成种子Seed
	seed := bip39.NewSeed(mnemonic, "") // password盐值不要加
	fmt.Println("seed", seed)

	// 接下来就是将种子恢复出主私钥 masterKey 这里进入到了BIP-32了 a
	masterKey, _ := bip32.NewMasterKey(seed)
	// 注意，此时还是主私钥，接下来要派生子私钥，派生出来的子私钥才是真正的“私钥”才能对应链的公钥，才能解压缩出地址
	fmt.Println("masterKey", masterKey)

	// 现在要派生出对应以太坊的子私钥，遵循BIP-44
	// 接下来进入BIP-44 完成派生,完成对应path参数 m / purpose' / coin_type' / account' / change / address_index
	// 通过主私钥派生出子私钥,FirstHardenedChild = uint32(0x80000000) 是一个常量，对应强化派生范围
	key, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)   // purpose' ： 44 是固定值，即BIP-44标准，强化派生
	key, _ = key.NewChildKey(bip32.FirstHardenedChild + uint32(118)) // coin_type' ：118是cosmos标识'， 继续强化派生
	key, _ = key.NewChildKey(bip32.FirstHardenedChild + uint32(0))   // account' : 0 标记账户类型，从0开始，强化派生
	key, _ = key.NewChildKey(uint32(0))                              // change :0 外部可见地址， 1 找零地址（外部不可见），通常是 0，普通派生
	key, _ = key.NewChildKey(uint32(0))                              // 地址索引 0 1 2 3 这样索引，普通派生
	// 派生完毕，对应的path 就是 " m/44'/118'/0'/0/0 "

	// 子私钥key 已经出来了，先打印私钥，key.Key就是私钥，注意要转化进制
	cosmosPrivateKey := hex.EncodeToString(key.Key) // 编码成字符串 , 这里打印出来的私钥就可以
	fmt.Println("privateKey", cosmosPrivateKey)

	fmt.Println("--------------COSMOS在公钥转地址的地方这里会有些不同---------------------")
	// 要对公钥进行 SHA-256处理----> cosmos 的公钥要进行sha256 然后 ripemd哈希
	sha256Hasher := sha256.New()            // 相当于在创建sha256这个对象
	sha256Hasher.Write(key.PublicKey().Key) // 这里就是在对公钥的字节进行sha256哈希加密
	sha256Hash := sha256Hasher.Sum(nil)     //

	// 然后对 SHA-256 哈希结果进行 RIPEMD-160 哈希加密处理
	ripemd160Hasher := ripemd160.New()     //  同样 创建ripemd-160这个对象
	ripemd160Hasher.Write(sha256Hash)      // 这里是在对之前sha256加密后的公钥再进行RIPEMD-160
	pubKeyHash := ripemd160Hasher.Sum(nil) //

	// Bech32 编码
	bech32Addr, err := bech32.ConvertAndEncode("cosmos", pubKeyHash)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("cosmosAddress", bech32Addr)
	// 最后把 cosmosPrivateKey 导入到 开普勒钱包 里面去验证下，私钥和地址对上了没

}

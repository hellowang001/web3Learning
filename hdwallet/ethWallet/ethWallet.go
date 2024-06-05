package main

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// 生成一个以太坊的钱包
func main() {
	fmt.Println("开始")
	// 首先你可以生成一个随机熵，熵源助记词是BIP-39，
	entropy, _ := bip39.NewEntropy(128)
	fmt.Println("entroy", entropy)
	// 通过熵源生成助记词  ==> 注意，不一定要有助记词才有种子，只是助记词方便备份，可以转成种子，你要直接由种子也行，但不好记
	mnemonic, _ := bip39.NewMnemonic(entropy)
	fmt.Println("mnemonic", mnemonic)
	// 通过助记词生成种子Seed
	seed := bip39.NewSeed(mnemonic, "") // password盐值不要加
	fmt.Println("seed", seed)

	// 接下来就是将种子恢复出主私钥 masterKey 这里进入到了BIP-32了
	masterKey, _ := bip32.NewMasterKey(seed)
	// 注意，此时还是主私钥，接下来要派生子私钥，派生出来的子私钥才是真正的“私钥”才能对应链的公钥，才能解压缩出地址
	fmt.Println("masterKey", masterKey)

	// 现在要派生出对应以太坊的子私钥，遵循BIP-44
	// 接下来进入BIP-44 完成派生,完成对应path参数 m / purpose' / coin_type' / account' / change / address_index
	// 通过主私钥派生出子私钥,FirstHardenedChild = uint32(0x80000000) 是一个常量，对应强化派生范围
	key, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)  // purpose' ： 44 是固定值，即BIP-44标准，强化派生
	key, _ = key.NewChildKey(bip32.FirstHardenedChild + uint32(60)) // coin_type' ：60是以太坊标识'， 继续强化派生
	key, _ = key.NewChildKey(bip32.FirstHardenedChild + uint32(0))  // account' : 0 标记账户类型，从0开始，强化派生
	key, _ = key.NewChildKey(uint32(0))                             // change :0 外部可见地址， 1 找零地址（外部不可见），通常是 0，普通派生
	key, _ = key.NewChildKey(uint32(0))                             // 地址索引 0 1 2 3 这样索引，普通派生
	// 派生完毕，对应的path 就是 " m/44'/60'/0'/0/0 "

	// 子私钥key 已经出来了，先打印私钥，key.Key就是私钥，注意要转化进制
	ethPrivateKey := hex.EncodeToString(key.Key) // 编码成字符串 , 这里打印出来的私钥就可以
	fmt.Println("privateKey", ethPrivateKey)
	// 子私钥 key 里面就包含他对应的公钥属性，拿出来即可,因为公钥是私钥椭圆曲线加密得来的结果
	ethPublicKey := hex.EncodeToString(key.PublicKey().Key) // 编码成字符串
	fmt.Println("ethPublicKey", ethPublicKey)

	// 现在拿到公钥了，先对公钥进行压缩 keccak256 压缩成 32 byte
	compressPubKey, _ := crypto.DecompressPubkey(key.PublicKey().Key)
	fmt.Println("compressPubKey", compressPubKey)
	// 压缩32字节后取最后 20 byte 就是地址了
	ethAddre := crypto.PubkeyToAddress(*compressPubKey).Hex() // Hex是16进制转字符串
	fmt.Println("ethAddre", ethAddre)                         // ethAddre 0x38B59D6D4ef6A4991926Cf04c7c2092a0E86140F

	// 最后把 ethPrivateKey 导入到 metaMask 里面去验证下，私钥和地址对上了没
}

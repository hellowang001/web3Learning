package stackswallet

import (
	"github.com/sirupsen/logrus"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func GenerateOfflineAddress() {
	// 生成随机墒---bit39
	entropy, _ := bip39.NewEntropy(256)
	// 墒生成助记词
	mnemonic, _ := bip39.NewMnemonic(entropy)

	// 从助记词生成种子

	logrus.Infof("mnemonic: %s", mnemonic)
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, _ := bip32.NewMasterKey(seed)
	key, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)    // 这一步是在做强化派生， 对应 purpose' 44是符合BIP-44标准
	key, _ = key.NewChildKey(bip32.FirstHardenedChild + uint32(5757)) // 还是强化派生，对应 coin_type类型60'
	key, _ = key.NewChildKey(bip32.FirstHardenedChild + uint32(0))    // 还是强化派生，对应account
	key, _ = key.NewChildKey(uint32(0))                               // 常规派生，对应 change
	key, _ = key.NewChildKey(uint32(0))                               // 常规派生，对应的 address_index，

}

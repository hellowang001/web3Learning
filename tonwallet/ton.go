package tonwallet

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// 离线地址生成
func GenerateOfflineAddress() {
	// 生成随机墒---bit39
	//entropy, _ := bip39.NewEntropy(256)
	// 墒生成助记词
	//mnemonic, err := bip39.NewMnemonic(entropy)
	//if err != nil {
	//	logrus.Warnf("mnemonic error: %s", err)
	//}
	mnemonic := "fish fever crash stove pause sign diagram fresh bus gasp velvet energy truck oven remain sell change glance hawk fatal electric addict biology moral"
	logrus.Infof("mnemonic: %s", mnemonic)
	// 助记词到
	seed := bip39.NewSeed(mnemonic, "")
	// 种子生成key--->bip32
	masterKey, _ := bip32.NewMasterKey(seed)

	// MasterKey 到child key --bip44
	key, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)   // 这一步是在做强化派生， 对应 purpose' 44是符合BIP-44标准
	key, _ = key.NewChildKey(bip32.FirstHardenedChild + uint32(607)) // 还是强化派生，对应 coin_type类型60'
	key, _ = key.NewChildKey(bip32.FirstHardenedChild + uint32(0))   // 还是强化派生，对应account
	key, _ = key.NewChildKey(uint32(0))                              // 常规派生，对应 change
	key, _ = key.NewChildKey(uint32(0))                              // 常规派生，对应的 address_index
	// child key 到私钥
	tonPrivateKey := hex.EncodeToString(key.Key) // 编码成字符串 , 这里打印出来的私钥就可以
	//tonPrivateKey := "55ae45784b6179a74c1e45b4c724b32ed7c8f0f675d6f947dc321a28b4dae172" // 编码成字符串 , 这里打印出来的私钥就可以
	logrus.Infof("tonPrivateKey %s", tonPrivateKey)
	piv_key := ed25519.PrivateKey(tonPrivateKey)
	publicKey := ed25519.PrivateKey(tonPrivateKey).Public().(ed25519.PublicKey)
	logrus.Infof("publicKey: %s", publicKey)

	publicKeyHex := hex.EncodeToString(publicKey)
	logrus.Infof("publicKeyHex: %s", publicKeyHex)
	//client, err := liteapi.NewClientWithDefaultTestnet()
	client, err := liteapi.NewClientWithDefaultMainnet()
	w, err := wallet.New(piv_key, wallet.V4R2, client)
	if err != nil {
		logrus.Warnf("生成钱包错误 %v", err)
	}

	addr := w.GetAddress()
	logrus.Infof("address: %s", addr)
	// 生成钱包地址
	address := generateTonAddress(publicKey)
	logrus.Infof("TON Wallet Address:%v", address)
	b, _ := w.GetBalance(context.Background())
	logrus.Infof("Wallet balance: %v\n", b)
	const OneTON tlb.Grams = 1_000_000
	recipientAddr, _ := ton.AccountIDFromRaw("0:507dea7d606f22d9e85678d3eede39bbe133a868d2a0e3e07f5502cb70b8a512")
	simpleTransfer := wallet.SimpleTransfer{
		Amount: ton.OneTON,
		//Address: tongo.MustParseAccountID("EQBszTJahYw3lpP64ryqscKQaDGk4QpsO7RO6LYVvKHSINS0"),
		Address: recipientAddr,
		Comment: "hi! hope it will be enough for buying a yacht",
	}
	logrus.Infof("simpleTransfer: %v", simpleTransfer)
	bg := context.Background()
	if err = w.Send(bg, simpleTransfer); err != nil {
		logrus.Warnf("send trx error: %v", err)
	}

}

func generateTonAddress(publicKey ed25519.PublicKey) string {
	// 计算 SHA-256 哈希
	hash := sha256.Sum256(publicKey)

	// 设置工作链 ID，通常为 0
	workchain := byte(0)

	// 构建原始地址
	rawAddress := append([]byte{workchain}, hash[:]...)

	// Base64 URL 编码地址
	encodedAddress := base64.RawURLEncoding.EncodeToString(rawAddress)

	// 确保没有填充字符
	return encodedAddress
}

func BuildTrx() {

}

package suiwallet

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/utils"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/blake2b"
)

func SuiWallet() {
	fmt.Println("开始")
	//entropy, _ := bip39.NewEntropy(128)
	//fmt.Println("entroy:", entropy)
	// 通过熵源生成助记词  ==> 注意，不一定要有助记词才有种子，只是助记词方便备份，可以转成种子，你要直接由种子也行，但不好记
	//mnemonic, _ := bip39.NewMnemonic(entropy)
	mnemonic := "found ramp year coconut depend undo salon hybrid ocean tissue install senior"
	//mnemonic := "fish fever crash stove pause sign diagram fresh bus gasp velvet energy truck oven remain sell change glance hawk fatal electric addict biology moral"

	fmt.Println("mnemonic:", mnemonic)
	// 通过助记词生成种子Seed
	seed := bip39.NewSeed(mnemonic, "") // password盐值不要加
	//fmt.Println("seed", seed)

	// 接下来就是将种子恢复出主私钥 masterKey
	masterKey, _ := bip32.NewMasterKey(seed)
	fmt.Println("masterKey", masterKey)

	// sui的BIP44路径: m/44'/195'/0'/0/0
	masterKey, _ = masterKey.NewChildKey(bip32.FirstHardenedChild + 44)  // purpose' ： 44 是固定值，即BIP-44标准
	masterKey, _ = masterKey.NewChildKey(bip32.FirstHardenedChild + 784) // coin_type' ：195是sui的标识
	masterKey, _ = masterKey.NewChildKey(bip32.FirstHardenedChild + 0)   // account' : 0
	masterKey, _ = masterKey.NewChildKey(0)                              // change :0
	masterKey, _ = masterKey.NewChildKey(0)                              // address_index : 0

	//
	privateKeyBytes := masterKey.Key
	privateKey := ed25519.NewKeyFromSeed(privateKeyBytes[:])
	privateKeyStr := hex.EncodeToString(privateKey)
	fmt.Printf("privateKeyStr=%s\n", privateKeyStr)

	// 使用公钥生成地址
	//priKey := ed25519.NewKeyFromSeed(seed[:])r
	pubKey := privateKey.Public().(ed25519.PublicKey)

	tmp := []byte{byte(0)}
	tmp = append(tmp, pubKey...)
	addrBytes := blake2b.Sum256(tmp)
	suiAddressStr := "0x" + hex.EncodeToString(addrBytes[:])[:64]
	//// 1.对公钥进行blake2b进行哈希，(sui标准版）
	//hasher, _ := blake2b.New(32, nil)
	//hasher.Write(publicKey)
	//publicKeyHash := hasher.Sum(nil)
	//
	//// 2. 取哈希的钱32字节作为地址 (sui地址规则)
	//addressBytes := publicKeyHash[:32]
	////suiAddressStr := "0x" + hex.EncodeToString(addressBytes)
	//suiAddressStr := "0x" + hex.EncodeToString(addressBytes[:])[:64]
	fmt.Printf("Sui Address: %s\n", suiAddressStr)

}

func SuiWalletSDK() *signer.Signer {
	//var ctx = context.Background()
	//var cli = sui.NewSuiClient(constant.BvTestnetEndpoint)
	mnemonic := "found ramp year coconut depend undo salon hybrid ocean tissue install senior"

	signerAccount, err := signer.NewSignertWithMnemonic(mnemonic)
	if err != nil {
		fmt.Printf("err=%v", err)
	}
	//suiPrivateKey := hex.EncodeToString(signerAccount.PriKey) // 编码成字符串 , 这里打印出来的私钥就可以

	//fmt.Println("privateKey", suiPrivateKey)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	//priKey := signerAccount.PriKey
	fmt.Printf("signerAccount.Address: %s\n", signerAccount.Address)
	return signerAccount
}

func RequestDevNetSuiFromFaucet() {
	// 水龙头领
	faucetHost, err := sui.GetFaucetHost(constant.SuiTestnet)
	if err != nil {
		fmt.Printf("get faucetHost %s", err)
	}
	fmt.Println("faucetHost:", faucetHost)
	header := map[string]string{}
	recipient := "0x1554d2f25f3b7bcf21b1f9fa3a2914f153c050e14ecad72327731388c900e38a"
	err = sui.RequestSuiFromFaucet(faucetHost, recipient, header)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// the successful transaction block url: https://suiscan.xyz/devnet/account/0x1554d2f25f3b7bcf21b1f9fa3a2914f153c050e14ecad72327731388c900e38a
	fmt.Println("Request DevNet Sui From Faucet success")
}

const ZANTestnetEndpoint = "https://api.zan.top/node/v1/sui/testnet/75189005b6174ea99b7214faf1a13cb5"

var ctx = context.Background()

var cli = sui.NewSuiClient(ZANTestnetEndpoint)

//var cli = sui.NewSuiClient(constant.BvTestnetEndpoint)

//var cli = sui.NewSuiClient(constant.BvMainnetEndpoint)

//
//var cli = sui.NewSuiClient(constant.SuiTestnetEndpoint)

func SendSui() {
	//var ctx = context.Background()
	////var cli = sui.NewSuiClient(constant.BvTestnetEndpoint)
	//
	//var cli = sui.NewSuiClient(constant.SuiTestnetEndpoint)
	//signerAccount, err := signer.NewSignertWithMnemonic("input your mnemonic")
	alice := SuiWalletSDK()
	fmt.Printf("alice address =%s\n ", alice.Address)
	rsp, err := cli.TransferSui(ctx, models.TransferSuiRequest{
		Signer:      alice.Address,
		SuiObjectId: "0x7d20dcdb2bca4f508ea9613994683eb4e76e9c4ed371169677c1be02aaf0b58e",
		GasBudget:   "1000000",
		Recipient:   "0xb7f98d327f19f674347e1e40641408253142d6e7e5093a7c96eda8cdfd7d9bb5",
		Amount:      "1",
	})
	if err != nil {
		fmt.Printf("TransferSui = %s", err.Error())
		panic("错啦")
	}
	//fmt.Printf("TxnMetaData =%s", TxnMetaData.TxBytes)
	//utils.PrettyPrint(rsp)
	//see the successful transaction url: https://explorer.sui.io/txblock/C7iYsH4tU5RdY1KBeNax4mCBn3XLZ5UswsuDpKrVkcH6?network=testnet
	rsp2, err := cli.SignAndExecuteTransactionBlock(ctx, models.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: rsp,
		PriKey:      alice.PriKey,
		// only fetch the effects field
		Options: models.SuiTransactionBlockOptions{
			ShowInput:    true,
			ShowRawInput: true,
			ShowEffects:  true,
		},
		RequestType: "WaitForLocalExecution",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	utils.PrettyPrint(rsp2)
}
func GetCoin() {
	//cli.SuiXGetCoins()
	//alice := SuiWalletSDK()

	//coins, err := cli.SuiXGetCoins(ctx, models.SuiXGetCoinsRequest{
	//	Owner:    alice.Address,
	//	CoinType: "0x2::sui::SUI", // 标准SUI Coin类型
	//	Limit:    10,
	//})
	//fmt.Printf("len = %d\n", len(coins.Data))
	//// 打印可用的Coin Object IDs
	//for _, coin := range coins.Data {
	//	fmt.Printf("Coin ID: %s, Balance: %s\n", coin.CoinObjectId, coin.Balance)
	//}
	//
	//if err != nil {
	//	fmt.Println(err.Error())
	//	panic("获取代币信息错误")
	//	return
	//}
	balances, err := cli.SuiXGetAllBalance(ctx, models.SuiXGetAllBalanceRequest{
		//Owner: alice.Address,
		//Owner: "0xd6a00995db029a7cea10487e8a2afb7db5b2372e6a122295c3b7eeb8e71b4a43",
		Owner: "0x875440044b941d198c8bae9e3d0fb31127f5acc3054b682c65f0e740c49b8a61",
		//CoinType: "0x2::sui::SUI", // 标准SUI Coin类型
	})

	if err != nil {
		fmt.Printf("报错啦: %v\n", err.Error())
		//fmt.Println(err.Error())
		//panic("获取代币信息错误")
		return
	}
	fmt.Printf("balances= %v", balances)

}

//func GetBlock() {
//	cli.SuiXGetCoins()
//}

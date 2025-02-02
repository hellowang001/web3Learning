package solwallet

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/pkg/hdwallet"
	"github.com/blocto/solana-go-sdk/program/system"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/sirupsen/logrus"
	"github.com/tyler-smith/go-bip39"
)

func GenerateOfflineAddress() types.Account {
	// 生成随机墒---bit39
	//entropy, _ := bip39.NewEntropy(256)
	// 墒生成助记词
	//mnemonic, _ := bip39.NewMnemonic(entropy)
	mnemonic := "fish fever crash stove pause sign diagram fresh bus gasp velvet energy truck oven remain sell change glance hawk fatal electric addict biology moral"
	// 从助记词生成种子

	logrus.Infof("mnemonic: %s", mnemonic)
	seed := bip39.NewSeed(mnemonic, "")

	path := `m/44'/501'/0'/0'`
	derivedKey, _ := hdwallet.Derived(path, seed)
	//accountFromSeed, err := types.AccountFromBytes(derivedKey)
	accountFromSeed, err := types.AccountFromSeed(derivedKey.PrivateKey)
	if err != nil {
		fmt.Printf("err=%v", err)
	}
	publicKey := accountFromSeed.PublicKey.ToBase58()
	fmt.Println("Solana Wallet Address:", publicKey)
	return accountFromSeed
}

func SendTrx() {
	Account := GenerateOfflineAddress()
	solClient := getClient()
	bg := context.Background()
	slot, err := solClient.GetSlot(bg)
	if err != nil {
		logrus.Warnf("err:=%s", err)
	}
	logrus.Infof("slot: %v", slot)

	block, err := solClient.GetBlock(bg, slot)
	if err != nil {
		logrus.Warnf("err:=%s", err)
	}
	logrus.Infof("block hash: %v", block.Blockhash)

	//instruction := system.Transfer(system.TransferParam{
	//	From:   Account.PublicKey,
	//	To:     common.PublicKeyFromString("6cPnfGr9Y4bZK7ykNpxe2hkKfaPPgsy6Tu5ahyGhzQLt"),
	//	Amount: 123000,
	//})

	//msg := types.NewMessageParam{
	//	FeePayer:        Account.PublicKey,
	//	RecentBlockhash: block.Blockhash,
	//	Instructions:    []types.Instruction{instruction},
	//}
	//MsgTx := types.NewMessage(msg)

	NewTx := types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        Account.PublicKey,
			RecentBlockhash: block.Blockhash,
			Instructions: []types.Instruction{system.Transfer(system.TransferParam{
				From:   Account.PublicKey,
				To:     common.PublicKeyFromString("6cPnfGr9Y4bZK7ykNpxe2hkKfaPPgsy6Tu5ahyGhzQLt"),
				Amount: 12300000,
			})},
		}),
		Signers: []types.Account{Account},
	}

	tx, err := types.NewTransaction(NewTx)
	if err != nil {
		logrus.Warnf("new transaction error: %s", err)
	}

	hash, err := solClient.SendTransaction(bg, tx)
	if err != nil {
		logrus.Warnf("SendTransaction error: %s", err)
	}
	logrus.Infof("hash: %v", hash)

}

func TokenTx() {
	Account := GenerateOfflineAddress()
	solClient := getClient()
	bg := context.Background()
	slot, err := solClient.GetSlot(bg)
	if err != nil {
		logrus.Warnf("err:=%s", err)
	}
	block, err := solClient.GetBlock(bg, slot)
	if err != nil {
		logrus.Warnf("err:=%s", err)
	}
	logrus.Infof("block: %v", block)
	tokenTx := token.TransferCheckedParam{
		From:     Account.PublicKey,
		To:       common.PublicKeyFromString("6cPnfGr9Y4bZK7ykNpxe2hkKfaPPgsy6Tu5ahyGhzQLt"),
		Mint:     common.PublicKeyFromString("Df6yfrKC8kZE3KNkrHERKzAetSxbrWeniQfyJY4Jpump"),
		Auth:     Account.PublicKey,
		Signers:  []common.PublicKey{},
		Amount:   12345,
		Decimals: 9,
	}
	instruction := token.TransferChecked(tokenTx)
	NewTx := types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        Account.PublicKey,
			RecentBlockhash: block.Blockhash,
			Instructions:    []types.Instruction{instruction},
		}),
		Signers: []types.Account{Account},
	}
	tx, err := types.NewTransaction(NewTx)
	if err != nil {
		logrus.Warnf("new transaction error: %s", err)
	}

	hash, err := solClient.SendTransaction(bg, tx)
	if err != nil {
		logrus.Warnf("new transaction error: %s", err)
	}
	logrus.Infof("hash: %v", hash)
}

func getAccountInfo() {
	solClient := getClient()
	bg := context.Background()
	address := "4bQ3w1H8UUsmTKuXZyzVQmpErSVKmPuJhuN4xvkTFRJ5"
	accountInfo, err := solClient.GetAccountInfo(bg, address)
	if err != nil {
		logrus.Warnf("err:=%s", err)
	}
	logrus.Infof("accountInfo: %v", accountInfo)
	GetBalance, err := solClient.GetBalance(bg, address)
	if err != nil {
		logrus.Warnf("err:=%s", err)
	}
	logrus.Infof("GetBalance: %v", GetBalance)
}

func getRecentBlockHash() string {
	solClient := getClient()
	bg := context.Background()
	count, err := solClient.GetSlot(bg)
	if err != nil {
		logrus.Warnf("err:=%s", err)
	}
	block, err := solClient.GetBlock(bg, count)
	if err != nil {
		logrus.Warnf("err:=%s", err)
	}
	logrus.Infof("block: %v", block)
	return block.Blockhash
}
func getClient() *client.Client {
	//ENDPOINT := "https://solana-mainnet.g.alchemy.com/v2/xxxxxxxxxxxxxxxxxxxxxxxx"
	//ENDPOINT := "https://solana-mainnet.g.alchemy.com/v2/fzgfj4QuLlNEyn2LrLZsseBAClGdMnyP"
	ENDPOINT := "https://rough-frequent-tent.SOLANA_MAINNET.quiknode.pro/1f710f95479f793bf564bba4ababea2c969137f6/"
	solClient := client.NewClient(ENDPOINT)
	return solClient
}

package tronwallet

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/fbsobreira/gotron-sdk/pkg/abi"
	"github.com/fbsobreira/gotron-sdk/pkg/account"
	"github.com/fbsobreira/gotron-sdk/pkg/client/transaction"
	"github.com/fbsobreira/gotron-sdk/pkg/store"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"

	"github.com/fbsobreira/gotron-sdk/pkg/client"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"google.golang.org/grpc"
)

func TronWallet() (piv, addr string) {
	fmt.Println("开始")
	// entropy, _ := bip39.NewEntropy(128)
	// fmt.Println("entroy:", entropy)
	// 通过熵源生成助记词  ==> 注意，不一定要有助记词才有种子，只是助记词方便备份，可以转成种子，你要直接由种子也行，但不好记
	// mnemonic, _ := bip39.NewMnemonic(entropy)
	mnemonic := "fish fever crash stove pause sign diagram fresh bus gasp velvet energy truck oven remain sell change glance hawk fatal electric addict biology moral"

	fmt.Println("mnemonic:", mnemonic)
	// 通过助记词生成种子Seed
	seed := bip39.NewSeed(mnemonic, "") // password盐值不要加
	fmt.Println("seed", seed)

	// 接下来就是将种子恢复出主私钥 masterKey
	masterKey, _ := bip32.NewMasterKey(seed)
	fmt.Println("masterKey", masterKey)

	// TRON的BIP44路径: m/44'/195'/0'/0/0
	masterKey, _ = masterKey.NewChildKey(bip32.FirstHardenedChild + 44)          // purpose' ： 44 是固定值，即BIP-44标准
	masterKey, _ = masterKey.NewChildKey(bip32.FirstHardenedChild + uint32(195)) // coin_type' ：195是TRON的标识
	masterKey, _ = masterKey.NewChildKey(bip32.FirstHardenedChild + uint32(0))   // account' : 0
	masterKey, _ = masterKey.NewChildKey(uint32(0))                              // change :0
	masterKey, _ = masterKey.NewChildKey(uint32(0))                              // address_index : 0

	// 获取私钥
	tronPrivateKey := hex.EncodeToString(masterKey.Key)
	fmt.Println("tronPrivateKey:", tronPrivateKey)

	// 获取公钥
	tronPublicKey := hex.EncodeToString(masterKey.PublicKey().Key)
	fmt.Println("tronPublicKey:", tronPublicKey)

	// 从公钥生成TRON地址
	// 1. 先对公钥进行Keccak-256哈希
	pubKeyBytes := masterKey.PublicKey().Key
	fmt.Printf("len: %v\n", len(pubKeyBytes))
	if len(pubKeyBytes) == 33 {
		// 如果是压缩格式的公钥，需要解压缩
		decompressedPubKey, _ := crypto.DecompressPubkey(pubKeyBytes)
		pubKeyBytes = crypto.FromECDSAPub(decompressedPubKey) // 转换为未压缩格式
	}

	hash := crypto.Keccak256(pubKeyBytes[1:]) // 去掉0x04前缀,对公钥进行 Keccak-256 哈希运算。

	// 2. 取最后20字节作为地址
	addressBytes := hash[len(hash)-20:]

	// 3. 添加前缀0x41
	prefix := []byte{0x41}
	rawAddr := append(prefix, addressBytes...)

	// 4. 计算两次SHA256得到校验和
	h := sha256.New()
	h.Write(rawAddr)
	hash1 := h.Sum(nil)

	h.Reset()
	h.Write(hash1)
	hash2 := h.Sum(nil)

	// 5. 取校验和的前4个字节
	checksum := hash2[:4]

	// 6. 将地址和校验和合并
	fullAddr := append(rawAddr, checksum...)

	// 7. Base58编码得到最终地址
	mainnetAddress := base58.Encode(fullAddr)
	fmt.Println("Mainnet Address:", mainnetAddress)

	// 8. 生成测试网地址（将0x41替换为0x27）
	testPrefix := []byte{0x27}
	testRawAddr := append(testPrefix, addressBytes...)
	h.Reset()
	h.Write(testRawAddr)
	testHash1 := h.Sum(nil)

	h.Reset()
	h.Write(testHash1)
	testHash2 := h.Sum(nil)
	testChecksum := testHash2[:4]
	testFullAddr := append(testRawAddr, testChecksum...)
	testnetAddress := base58.Encode(testFullAddr)
	fmt.Println("Testnet Address:", testnetAddress)
	return tronPrivateKey, mainnetAddress
}

func WalletGen() {
	err := account.RemoveAccount("alice")
	if err != nil {
		fmt.Println("RemoveAccount删除账户失败:", err)

	}
	alliceName, err := account.ImportFromPrivateKey("950f9fb7c95364c599c80a70e28e0d0d980f67276ce7e271a0a8eff124262510", "alice", "")
	if err != nil {
		fmt.Println("ImportFromPrivateKey恢复账户失败:", err)

	}
	fmt.Println("alliceName", alliceName)

}
func FindAccount() {
	fmt.Println("开始查找账户")
	finedAccount := store.LocalAccounts()
	fmt.Println("finedAccount", finedAccount)
}
func VerifyAccount() {
	fmt.Println("开始验证账户")
	keyStore := store.FromAccountName("alice")
	fmt.Println("keyStore", keyStore.Accounts()[0].Address.String())
}

func SendTrx() bool {
	fmt.Println("开始发送TRX") // api key b862ca9d-563b-4eec-8108-2d436b614561
	opts := make([]grpc.DialOption, 0)
	opts = append(opts, grpc.WithInsecure())
	c := client.NewGrpcClient("grpc.nile.trongrid.io:50051")

	if err := c.Start(opts...); err != nil {
		fmt.Println("Start err:", err)
		return false
	}
	if err := c.SetAPIKey("b862ca9d-563b-4eec-8108-2d436b614561"); err != nil {
		fmt.Println("SetAPIKey err:", err)
		return false
	}
	alice := "TBSssNBoAX5isgjaR9cLZ4i84rf1Q7A5gr"
	bob := "TJF5RmJW2io8pcbHVXAGf9i5uXncJQBxX3"

	// 1. 创建转账交易
	tx, err := c.Transfer(alice, bob, 1000)
	// 没有找到自定义带宽和能量的方法，只能用默认的
	if err != nil {
		fmt.Println("创建交易失败:", err)
		return false
	}
	// 2.恢复签名者
	keyStore, keyStoreAccount, err := store.UnlockedKeystore(alice, "")
	if err != nil {
		fmt.Printf("解锁账户失败=%v", err)
		return false
	}
	// 3. 签名交易
	controller := transaction.NewController(c, keyStore, keyStoreAccount, tx.Transaction)
	// 4. 广播交易--->这里很特别，新建控制器，通过控制器去执行交易，是一种设计模式，把步骤变成一种行为
	if err = controller.ExecuteTransaction(); err != nil {
		logrus.Warnf("ExecuteTransaction error: %s", err)
		return false
	}
	txHash, err := controller.TransactionHash()
	if err != nil {
		fmt.Println("TransactionHash err:", err)
		return false
	}
	fmt.Printf("https://nile.tronscan.org/#/transaction/%s", txHash[2:])
	return controller.Result.Result

}
func SendTrc20() bool {
	fmt.Println("开始发送TRC20")
	c := getClient()
	alice := "TBSssNBoAX5isgjaR9cLZ4i84rf1Q7A5gr"
	bob := "TJF5RmJW2io8pcbHVXAGf9i5uXncJQBxX3"
	contractAddress := "TF17BgPaZYbz8oxbjhriubPDsA7ArKoLX3"
	account, err := c.GetAccount(alice)
	if err != nil {
		fmt.Println("GetAccount err:", err)
		return false
	}
	fmt.Println("account", account)
	// 1. 构建交易发送TRC20
	tx, err := c.TRC20Send(alice, bob, contractAddress, big.NewInt(100000000), 25000000)
	if err != nil {
		fmt.Println("TRC20Send err:", err)
		return false
	}
	// 2.恢复签名者
	keyStore, keyStoreAccount, err := store.UnlockedKeystore(alice, "")
	if err != nil {
		fmt.Printf("解锁账户失败=%v", err)
		return false
	}
	signedTx, err := keyStore.SignTx(*keyStoreAccount, tx.Transaction)
	if err != nil {
		fmt.Println("SignTx err:", err)
		return false

	}
	rawData := signedTx.GetRawData()
	rawDataPro, _ := proto.Marshal(rawData)
	h256h := sha256.New()
	h256h.Write(rawDataPro)
	hash := h256h.Sum(nil)
	//fmt.Printf("hash=%v", common.BytesToHexString(hash))
	fmt.Printf("https://nile.tronscan.org/#/transaction/%s", hex.EncodeToString(hash))
	res, err := c.Broadcast(signedTx)
	if err != nil {
		fmt.Println("Broadcast err:", err)
		return false

	}
	return res.Result

	//
	//// 3. 签名交易
	//controller := transaction.NewController(c, keyStore, keyStoreAccount, tx.Transaction)
	//// 4. 广播交易--->这里很特别，新建控制器，通过控制器去执行交易，是一种设计模式，把步骤变成一种行为
	//if err = controller.ExecuteTransaction(); err != nil {
	//	logrus.Warnf("ExecuteTransaction error: %s", err)
	//	return false
	//}
	//txHash, err := controller.TransactionHash()
	//if err != nil {
	//	fmt.Println("TransactionHash err:", err)
	//	return false
	//}
	//fmt.Println("txHash", txHash)
	//
	//return controller.Result.Result

}

func TransactionBuilder() bool {
	fmt.Println("开始构建交易")
	c := getClient()
	alice := "TBSssNBoAX5isgjaR9cLZ4i84rf1Q7A5gr"
	bob := "TJF5RmJW2io8pcbHVXAGf9i5uXncJQBxX3"

	from := alice
	contractAddress := "TF17BgPaZYbz8oxbjhriubPDsA7ArKoLX3"
	method := "transfer(address,uint256)"

	// 1 计算能源消耗
	estimateEnergyMessage, err := c.EstimateEnergy(from, contractAddress, method, fmt.Sprintf(`
	[
		{"address":"TBSssNBoAX5isgjaR9cLZ4i84rf1Q7A5gr"},
		{"uint256":"100000000000000000"}
	]
	`), 0, "", 0)
	if err != nil {
		fmt.Println("EstimateEnergy err:", err)
		return false
	}
	fmt.Printf("txExtention: %+v\n", estimateEnergyMessage)
	// 能源消耗=estimateEnergyMessage.EnergyRequired  result:{result:true} energy_required:16915
	// 能源消耗乘以单价210
	fee_limit := estimateEnergyMessage.EnergyRequired * 180
	fmt.Println("energyCost", fee_limit)
	// 2组装交易
	tx, err := c.TRC20Send(alice, bob, contractAddress, big.NewInt(100000000000000000), fee_limit)
	if err != nil {
		fmt.Println("TRC20Send err:", err)
		return false
	}
	// 3.恢复签名者
	keyStore, keyStoreAccount, err := store.UnlockedKeystore(alice, "")
	if err != nil {
		fmt.Printf("解锁账户失败=%v", err)
		return false
	}
	// 4.签名交易
	signedTx, err := keyStore.SignTx(*keyStoreAccount, tx.Transaction)
	if err != nil {
		fmt.Println("SignTx err:", err)
		return false

	}
	rawData := signedTx.GetRawData()
	rawDataPro, _ := proto.Marshal(rawData)
	h256h := sha256.New()
	h256h.Write(rawDataPro)
	hash := h256h.Sum(nil)
	//fmt.Printf("hash=%v", common.BytesToHexString(hash))
	fmt.Printf("https://nile.tronscan.org/#/transaction/%s", hex.EncodeToString(hash))
	// 5.广播交易
	res, err := c.Broadcast(signedTx)
	if err != nil {
		fmt.Println("Broadcast err:", err)
		return false

	}
	return res.Result
}
func getClient() *client.GrpcClient {
	opts := make([]grpc.DialOption, 0)
	opts = append(opts, grpc.WithInsecure())
	c := client.NewGrpcClient("grpc.nile.trongrid.io:50051")
	if err := c.Start(opts...); err != nil {
		fmt.Println("Start err:", err)
	}
	if err := c.SetAPIKey("b862ca9d-563b-4eec-8108-2d436b614561"); err != nil {
		fmt.Println("SetAPIKey err:", err)
	}
	return c
}
func getTranscation(txHash string) {
	// 0x1e66d8290ff533c497af42dc6b0bda3cbaf9fa294c439daaba496db2783a7b13
	c := getClient()
	tx, err := c.GetTransactionByID(txHash)
	if err != nil {
		fmt.Println("GetTransaction err:", err)
	}
	fmt.Println("tx", tx)
}
func GetContractABI(contractAddress string) {
	c := getClient()

	// fmt.Println("abi", abi)
	// abiJson, _ := json.MarshalIndent(&abitest, "", "  ")
	// fmt.Println("abiJson", string(abiJson))
	method := "transfer"
	abitest, err := c.GetContractABI(contractAddress)
	if err != nil {
		fmt.Println("GetContractABI err:", err)
	}
	Arguments, err := abi.GetInputsParser(abitest, method)
	if err != nil {
		fmt.Println("GetParser err:", err)
	}
	// fmt.Println("Arguments", Arguments)
	for _, arg := range Arguments {
		fmt.Println("arg", arg)

	}
}

const (
	ResourceCode_BANDWIDTH  = 0 // 带宽
	ResourceCode_ENERGY     = 1 // 能量
	ResourceCode_TRON_POWER = 2 // 超级节点投票权 不可用，
)

func DelegateTrx() bool {
	fmt.Println("开始委托TRX")
	c := getClient()
	alice := "TBSssNBoAX5isgjaR9cLZ4i84rf1Q7A5gr"
	// bob := "TJF5RmJW2io8pcbHVXAGf9i5uXncJQBxX3"
	// 质押TRX 有三种：BANDWIDTH带宽，ENERGY能量，TRON_POWER超级节点投票权
	// 第一种：带宽
	// tx, err := c.FreezeBalanceV2(alice, ResourceCode_BANDWIDTH, 297000000)
	tx, err := c.FreezeBalanceV2(alice, ResourceCode_ENERGY, 297000000)
	// tx, err := c.FreezeBalanceV2(alice, ResourceCode_TRON_POWER, 297000000) // 不可用
	if err != nil {
		fmt.Println("FreezeBalanceV2 err:", err)
		return false
	}
	fmt.Println("tx", tx)
	// 2.恢复签名者
	keyStore, keyStoreAccount, err := store.UnlockedKeystore(alice, "")
	if err != nil {
		fmt.Printf("解锁账户失败=%v", err)
		return false
	}
	// 3. 签名交易
	controller := transaction.NewController(c, keyStore, keyStoreAccount, tx.Transaction)
	// 4. 广播交易--->这里很特别，新建控制器，通过控制器去执行交易，是一种设计模式，把步骤变成一种行为
	if err = controller.ExecuteTransaction(); err != nil {
		logrus.Warnf("ExecuteTransaction error: %s", err)
		return false
	}
	txHash, err := controller.TransactionHash()
	if err != nil {
		fmt.Println("TransactionHash err:", err)
		return false
	}
	fmt.Println("txHash", txHash)
	fmt.Printf("https://nile.tronscan.org/#/transaction/%s", txHash[2:])

	return controller.Result.Result
}

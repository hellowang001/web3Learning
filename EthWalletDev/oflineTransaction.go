package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

var (
	PrivateKeyStr = "c535facd6873ca2b3718e3ede4f626ae126c99b4a26b4354da704d8dc78b43c1"      // go生成的账户 发送发私钥
	fromAddrStr   = "0xa3856a939A623EdBde8f908037d3F33FceBC5408"                            // 对应go生成的账户 发送方地址
	toAddrStr     = "0x38B59D6D4ef6A4991926Cf04c7c2092a0E86140F"                            // 接收方地址
	ethUrl        = "https://eth-sepolia.g.alchemy.com/v2/fzgfj4QuLlNEyn2LrLZsseBAClGdMnyP" //API URL
)

func main() {

	// 把私钥个地址转码一下
	privateKey, _ := crypto.HexToECDSA(PrivateKeyStr)
	toAddr := common.HexToAddress(toAddrStr)

	// 新建连接，连上以太坊
	client, err := ethclient.Dial(ethUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Get the balance of an account 第0步：查询账户余额，要先确认自己有没有钱
	fromAddr := common.HexToAddress(fromAddrStr)
	balance, err := client.BalanceAt(context.Background(), fromAddr, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Account balance: %d\n", balance) // 1008465873901900000

	// Get the latest known block 第0步：获取最新块高，刷新
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Latest block: %d\n", block.Number().Uint64()) //6055020

	// 生成交易数据,这里注意fromAddr要转成 account common.Address对象

	// 第一步拿nonce，这个拿nonce的的方法已经封装好了的，直接用就可，也可以去网站上copy url下来发请求，不过那样不优雅
	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("nonce: %d\n", nonce)

	// 拿到gasPrice ，也可以不用拿，直接指定就可以
	gasPriceNow, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to retrieve gas price: %v", err)
	}
	fmt.Printf("gasPriceNow: %d\n", gasPriceNow)

	gasPrice := big.NewInt(2000000000)
	amount := big.NewInt(1100000000000000)
	gasLimit := uint64(21208)
	//gasPrice := big.NewInt(1000)
	data := []byte("Hello, World!")

	// 第二步 构建交易对象 ---> 生成一个未签名的交易
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddr,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	// 第三步 对交易进行签名，签名前先拿到链ID，即你要上的链
	chainID, err := client.NetworkID(context.Background()) // 拿到链ID
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey) // 对交易进行签名
	if err != nil {
		log.Fatal(err)
	}

	// 第四步  发送交易,这一步就是广播
	err = client.SendTransaction(context.Background(), signedTx) // 广播出去
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("转账交易已发送，交易哈希\n：%s\n", signedTx.Hash().Hex())
	// 第一笔交易哈希：0xb4f9ff541c3423e5bbffbb63f78fed2b384cde60b74b9a3838f05eb172c65e90 ， 第一次发交易的时候gasPrice 花费的比较低，所以第二次改了gasPrice后再发
	// 第二笔交易哈希：转账交易已发送，交易哈希：0x99962924698a903a7ce7b8b36ad8fb6cc04c1a6201789e3ce769f9e48a2995d6
	// 第三笔交易哈希：转账交易已发送，交易哈希：0x53785349f7f7247431dddc0e60dc3e270db4313b715ab50a9022ca90bd510890
	// 第四笔交易哈希：转账交易已发送，交易哈希：0xb6d4131e7b6deae72775198343173a585f5c4753147b630138dd3455824aed08
	// 第五笔交易哈希：转账交易已发送，交易哈希：0x544691942a682344271139954991e8da4d9aea933791f5064e98c1eae8a1f55c
	// 第六笔交易哈希：转账交易已发送，交易哈希：0x5da55b4791d6e701c2012babae955a6a6508fbc335fbfb043be5321ad0397fd5
	// 第七笔交易哈希：转账交易已发送，交易哈希：0x93afb74b4c4f0bca6c1ccdf4a13fe555f7a8768201804c0e21ec33b50febec46
	// 第八笔交易哈希：转账交易已发送，交易哈希：0x9b438259c6d412adc4ae036db8d006759d5e300ecbb94191417cef380aeb951c
	// 第八笔交易哈希：转账交易已发送，交易哈希：0xf951d0190718ccac2d20421e7b9e805b08be9aa3427e60f6d964165b2a914582
	// 第八笔交易哈希：转账交易已发送，交易哈希：0xf94951ea83e5971738d448d1ea7a06e945325cc9cf741fb2d415f3507c85f1f8
	// 前面的一直在pending中 有点问题，
	// 第八笔交易哈希：转账交易已发送，交易哈希：0xd65407692e1f872d47ce1231a80d9afbf0dbcf8336362f78262d4fd872a55c81
	// 第八笔交易哈希：转账交易已发送，交易哈希：0xcb297210b49a5a3747c27783ab3610b8b287134792e50ccbf6afe86c6c96edfd
	// 第八笔交易哈希：转账交易已发送，交易哈希：0xc922a93fd5a5eabe7819623a65aa6b572e755658f3a9f24e71a64f3d533b1a87
	// 第八笔交易哈希：转账交易已发送，交易哈希：0x4a8e649836a69d22e2b979c7330f03b97e7017b4a89d27e9c646051725f89a2c
	// 第八笔交易哈希：转账交易已发送，交易哈希：0xae0b39c8cb2e02f7c647778f4325f0cca6392fe54c8f144a8e4ae58da79484eb
	// 显示pending 然后转indexing状态 最后Success成功
	// pending 状态： This txn hash was found in our secondary node and should be picked up by our indexer in a short while.
	// indexing 状态：This transaction has been included and will be reflected in a short while.

}

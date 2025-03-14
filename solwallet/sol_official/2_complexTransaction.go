package sol_official

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/compute_budget"
	"github.com/blocto/solana-go-sdk/program/system"
	"github.com/blocto/solana-go-sdk/types"
	"log"
)

func CreateComplexTx() {

	c := NewClient()
	alice, err := types.AccountFromBase58(Piv)
	if err != nil {
		fmt.Printf("err=%v", err)
	}
	fmt.Printf("alice=%v\n", alice.PublicKey.String())
	// 查询下地址的当前余额
	balances, err := c.GetBalance(context.Background(), alice.PublicKey.String())
	if err != nil {
		fmt.Printf("get balances err = %v", err)
	}
	fmt.Printf("balances = %v\n", balances)
	// 2、组装交易
	// 2.1 拿recentBlockhashResponse
	recentBlockhashResponse, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}

	// 最小转账金额
	minimumBalanceForRentExemption, err := c.GetMinimumBalanceForRentExemption(context.Background(), 0)
	fmt.Printf("minimumBalanceForRentExemption = %v\n", minimumBalanceForRentExemption)
	// 2.2 定义Instruction然后填充它
	ins := make([]types.Instruction, 0, 4)
	// 交易分成2个大部份 Transactions 和 Instructions 指令，Instructions可以有多个
	// 追加优先费
	ComputerUnitPrice := uint64(2000000) //一个计算单元的价格
	ComputerUnitLimit := uint32(2000000) //计算单元的限制
	ins = append(ins, compute_budget.SetComputeUnitPrice(compute_budget.SetComputeUnitPriceParam{
		MicroLamports: ComputerUnitPrice,
	}))

	ins = append(ins, compute_budget.SetComputeUnitLimit(compute_budget.SetComputeUnitLimitParam{
		Units: ComputerUnitLimit,
	}))

	// DOINGS 我在做的事情是创建一个地址并为其在链上创建一个账号,参考createAccount方法
	// 创建测试帐户指令，生成一个钱包，先不管他的地址
	bob := types.NewAccount()
	fmt.Printf("bob address=%s\n", bob.PublicKey.String())
	createAccountIns2 := system.CreateAccount(system.CreateAccountParam{
		From:     alice.PublicKey,
		New:      bob.PublicKey,
		Owner:    common.SystemProgramID,
		Lamports: minimumBalanceForRentExemption + 1,
		Space:    0,
	})
	ins = append(ins, createAccountIns2)
	// 第二个指令，给静态地址转账
	Cindy := common.PublicKeyFromString("2e7MJy7rh3mr7myjEtyD6Bjyyc9fYcYVwwhGhgVDbx5U") // 随便一个静态地址
	ins = append(ins, system.Transfer(system.TransferParam{
		From:   alice.PublicKey,
		To:     Cindy,
		Amount: 100000000,
	}))
	// 第三个指令，给测试帐户转账
	ins = append(ins, system.Transfer(system.TransferParam{
		From:   alice.PublicKey,
		To:     bob.PublicKey,
		Amount: 200000000,
	}))
	// 第四个指令，给静态地址转账
	ins = append(ins, system.Transfer(system.TransferParam{
		From:   alice.PublicKey,
		To:     Cindy,
		Amount: 300000000,
	}))
	// 消息包含指令、区块哈希、fee支付者
	message := types.NewMessage(
		types.NewMessageParam{
			FeePayer: alice.PublicKey,
			//FeePayer:        common.PublicKeyFromString(alice.PublicKey.String()), //先试试这个，再试试上面的
			Instructions:    ins,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
		})
	// 交易包含签名和消息
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: message,
		Signers: []types.Account{alice, bob},
	})
	if err != nil {
		log.Fatalf("failed to new transaction, err: %v", err)
	}
	// 3、广播上链
	txhash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to SendTransaction, err: %v", err)
	}
	fmt.Println("tx hash", txhash)
}

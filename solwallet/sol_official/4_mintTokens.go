package sol_official

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/compute_budget"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/types"
	"log"
)

func MintToken() {

	c := NewClient()
	alice, err := types.AccountFromBase58(Piv)

	if err != nil {
		fmt.Printf("err=%v", err)
	}
	fmt.Printf("alice=%v\n", alice.PublicKey.String())
	//mintKeypair =6GVCVzbiGNfjATr91bWPTmiToZXUkACbyLtNStTmYAK9// 这个就是代币地址
	mintKeypair2 := common.PublicKeyFromString("6GVCVzbiGNfjATr91bWPTmiToZXUkACbyLtNStTmYAK9")
	//ATAmetadata := common.PublicKeyFromString("Esby6ub71djgjGzhkfujaCQrL7nvLTrsMnFFEG4Jtvm3")
	// 2、组装交易
	// 2.1 拿recentBlockhashResponse
	recentBlockhashResponse, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}

	// 2.2 定义Instruction然后填充它
	ins := make([]types.Instruction, 0, 2)

	// 追加优先费
	ComputerUnitPrice := uint64(2000000) //一个计算单元的价格
	ComputerUnitLimit := uint32(2000000) //计算单元的限制
	ins = append(ins, compute_budget.SetComputeUnitPrice(compute_budget.SetComputeUnitPriceParam{
		MicroLamports: ComputerUnitPrice,
	}))

	ins = append(ins, compute_budget.SetComputeUnitLimit(compute_budget.SetComputeUnitLimitParam{
		Units: ComputerUnitLimit,
	}))
	mintSize := uint64(82)
	minimumBalanceForRentExemption, err := c.GetMinimumBalanceForRentExemption(context.Background(), mintSize)
	fmt.Printf("minimumBalanceForRentExemption = %v\n", minimumBalanceForRentExemption)
	// 为某个地址铸造token出来
	// DOINGS 我们先完成前创建token和元数据的3个指令 createMintAccountInstruction,initializeMintInstruction,createMetadataInstruction
	// 拿到ATA地址
	toTokenATAAlice, _, err := common.FindAssociatedTokenAddress(alice.PublicKey, mintKeypair2)
	// 第一次没有创建ATA地址，所以需要创建ATA地址，后面就不用了，这个创建是要上链的
	//ins = append(ins, associated_token_account.Create(
	//	associated_token_account.CreateParam{
	//		Funder:                 alice.PublicKey,
	//		Owner:                  alice.PublicKey,
	//		Mint:                   mintKeypair2,
	//		AssociatedTokenAccount: toTokenATAAlice,
	//		//ProgramID:              programID,
	//	}))

	mintIns := token.MintTo(token.MintToParam{
		Mint:    mintKeypair2,
		To:      toTokenATAAlice,
		Auth:    alice.PublicKey,
		Signers: []common.PublicKey{alice.PublicKey},
		Amount:  10000,
	})
	ins = append(ins, mintIns)
	// 消息包含指令、区块哈希、fee支付者
	message := types.NewMessage(
		types.NewMessageParam{
			FeePayer:        alice.PublicKey,
			Instructions:    ins,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
		})
	// 交易包含签名和消息
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: message,
		Signers: []types.Account{alice},
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
	// mintKeypair address=6GVCVzbiGNfjATr91bWPTmiToZXUkACbyLtNStTmYAK9
	// ATAmetadata address  =Esby6ub71djgjGzhkfujaCQrL7nvLTrsMnFFEG4Jtvm3
}

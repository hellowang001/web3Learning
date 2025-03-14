package sol_official

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/pkg/pointer"
	"github.com/blocto/solana-go-sdk/program/compute_budget"
	"github.com/blocto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/blocto/solana-go-sdk/types"
	"log"
)

func UpdateMetadata() {
	c := NewClient()
	alice, err := types.AccountFromBase58(Piv)
	mintKeypair := common.PublicKeyFromString("6GVCVzbiGNfjATr91bWPTmiToZXUkACbyLtNStTmYAK9")
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
	// DOINGS 创建更新代币元数据的指令
	tokenData := token_metadata.DataV2{
		Name:                 "New wang test token 2",
		Symbol:               "NEWWANGTTT",
		Uri:                  "https://www.solana.com",
		SellerFeeBasisPoints: 0,
		Creators: &[]token_metadata.Creator{
			{
				Address:  alice.PublicKey,
				Verified: true,
				Share:    100,
			},
		},
		//Creators:   nil,
		Collection: nil, // 这里是再铸造nft的时候用的，nft就是总供应量为1的非同质化代币，这个代币归属哪个集合
		Uses:       nil,
	}
	seeds := [][]byte{}
	seeds = append(seeds, []byte("metadata"))
	seeds = append(seeds, common.MetaplexTokenMetaProgramID.Bytes())
	seeds = append(seeds, mintKeypair.Bytes())
	metadata, _, err := common.FindProgramAddress(seeds, common.MetaplexTokenMetaProgramID) // 派生PDA地址原数据
	if err != nil {
		fmt.Printf("err=%v", err)
	}
	fmt.Printf("ATA metadata address  =%s\n", metadata.String())
	updateMetadataInstruction := token_metadata.UpdateMetadataAccountV2(token_metadata.UpdateMetadataAccountV2Param{
		MetadataAccount:     metadata,
		UpdateAuthority:     alice.PublicKey,
		Data:                &tokenData,
		NewUpdateAuthority:  &alice.PublicKey,
		PrimarySaleHappened: nil,
		IsMutable:           pointer.Get[bool](true), // bool 的指针
	})
	ins = append(ins, updateMetadataInstruction)

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
}

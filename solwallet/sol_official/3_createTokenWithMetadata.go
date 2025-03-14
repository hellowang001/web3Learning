package sol_official

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/compute_budget"
	"github.com/blocto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/blocto/solana-go-sdk/program/system"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/types"
	"log"
)

func CreateTokenWithMetadata() {

	c := NewClient()
	alice, err := types.AccountFromBase58(Piv)
	//testWallet := common.PublicKeyFromString("2e7MJy7rh3mr7myjEtyD6Bjyyc9fYcYVwwhGhgVDbx5U")
	//spltoken:=types.NewAccount()
	if err != nil {
		fmt.Printf("err=%v", err)
	}
	fmt.Printf("alice=%v\n", alice.PublicKey.String())

	// 2、组装交易
	// 2.1 拿recentBlockhashResponse
	recentBlockhashResponse, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}

	// 2.2 定义Instruction然后填充它
	ins := make([]types.Instruction, 0, 6)

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
	// DOINGS 我们先完成前创建token和元数据的3个指令 createMintAccountInstruction,initializeMintInstruction,createMetadataInstruction
	// createMintAccountInstruction // create instruction for the token mint account
	// 生成一个钱包，先不管他的地址
	mintKeypair := types.NewAccount()
	fmt.Printf("mintKeypair address=%s\n", mintKeypair.PublicKey.String())
	createMintAccountInstruction := system.CreateAccount(system.CreateAccountParam{
		From:     alice.PublicKey,
		New:      mintKeypair.PublicKey,
		Owner:    common.TokenProgramID,
		Space:    mintSize,
		Lamports: minimumBalanceForRentExemption,
	})
	ins = append(ins, createMintAccountInstruction)
	// initializeMintInstruction
	initializeMintInstruction := token.InitializeMint(token.InitializeMintParam{
		Decimals:   2,                // token精度
		FreezeAuth: &alice.PublicKey, //冻结权限\冻结权威，-->只有被指定为冻结权威的地址才能冻结token
		//FreezeAuth: nil,                   //冻结权限，可以是空,空应该就是丢权威，空就不会有人再铸造token
		Mint:     mintKeypair.PublicKey, //帐户地址
		MintAuth: alice.PublicKey,       // 帐户权限\铸造权威，-->只有被指定为铸造权威的地址才能铸造更多的token
	})
	ins = append(ins, initializeMintInstruction)
	// createMetadataInstruction
	//metadata, _, err := common.FindAssociatedTokenAddress(common.MetaplexTokenMetaProgramID, mintKeypair.PublicKey)
	seeds := [][]byte{}
	seeds = append(seeds, []byte("metadata"))
	seeds = append(seeds, common.MetaplexTokenMetaProgramID.Bytes())
	seeds = append(seeds, mintKeypair.PublicKey.Bytes())
	metadata, _, err := common.FindProgramAddress(seeds, common.MetaplexTokenMetaProgramID) // 派生PDA地址原数据
	if err != nil {
		fmt.Printf("err=%v", err)
	}
	fmt.Printf("ATA metadata address  =%s\n", metadata.String())
	tokenData := token_metadata.DataV2{
		Name:                 "wang test token 2",
		Symbol:               "WANGTTT",
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
	createMetadataInstruction := token_metadata.CreateMetadataAccountV3(token_metadata.CreateMetadataAccountV3Param{
		Metadata:                metadata,
		Mint:                    mintKeypair.PublicKey,
		MintAuthority:           alice.PublicKey,
		Payer:                   alice.PublicKey,
		UpdateAuthority:         alice.PublicKey,
		UpdateAuthorityIsSigner: true,
		IsMutable:               true,
		Data:                    tokenData,
		CollectionDetails:       nil,
	})
	ins = append(ins, createMetadataInstruction)
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
		Signers: []types.Account{alice, mintKeypair},
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

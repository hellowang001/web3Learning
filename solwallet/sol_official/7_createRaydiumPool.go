package sol_official

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/pkg/bincode"
	"github.com/blocto/solana-go-sdk/program/compute_budget"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/near/borsh-go"
	"log"
	"time"
)
import (
	"crypto/sha256"
)

//discriminator := GetInstructionDiscriminator("create_pool")

// ---------- 辅助函数和类型定义 ----------
// 计算指令 Discriminator
func getInstructionDiscriminator(name string) [8]byte {
	seed := []byte("global:" + name)
	hash := sha256.Sum256(seed)
	var discriminator [8]byte
	copy(discriminator[:], hash[:8])
	return discriminator
}

// 定义 uint128 类型
type uint128 struct {
	Lower uint64 // 低 64 位
	Upper uint64 // 高 64 位
}

// 定义指令数据结构
type CreatePoolInstructionData struct {
	Discriminator [8]byte
	SqrtPriceX64  uint128
	OpenTime      uint64
}

func CreateRaydiumCLPMPoolV2() {
	// ---------- 1. 设置必要的常量和变量 ----------
	// 替换为您的实际程序 ID
	programID := common.PublicKeyFromString("devi51mZmdwUJGU9hjN27vEz64Gps7uUefqxg27EAtH")

	// 替换为实际的公钥
	alice, err := types.AccountFromBase58(Piv)
	if err != nil {
		fmt.Printf("err=%v", err)
	}
	poolCreatorPubkey := alice.PublicKey
	ammConfigPubkey := common.PublicKeyFromString("CQYbhr6amxUER4p5SC44C63R4qw4NFc9Z4Db9vF4tZwG")
	tokenMint0Pubkey := common.PublicKeyFromString("So11111111111111111111111111111111111111112")
	//tokenMint1Pubkey := common.PublicKeyFromString("So11111111111111111111111111111111111111112")
	tokenMint1Pubkey := common.PublicKeyFromString("6GVCVzbiGNfjATr91bWPTmiToZXUkACbyLtNStTmYAK9")
	//tokenMint1Pubkey := common.PublicKeyFromString("8w52rwTTDxE8XQgXnBFLfNd9FJjQBUXUHTsq1rvNpump")
	//tokenMint0Pubkey := common.PublicKeyFromString("6GVCVzbiGNfjATr91bWPTmiToZXUkACbyLtNStTmYAK9")

	// 确定 tokenMint0 和 tokenMint1 的顺序
	var lowerMintPubkey, higherMintPubkey common.PublicKey
	if tokenMint0Pubkey.String() < tokenMint1Pubkey.String() {
		lowerMintPubkey = tokenMint0Pubkey
		higherMintPubkey = tokenMint1Pubkey
	} else {
		//lowerMintPubkey = tokenMint1Pubkey
		//higherMintPubkey = tokenMint0Pubkey
		lowerMintPubkey = tokenMint0Pubkey
		higherMintPubkey = tokenMint1Pubkey
	}

	// sqrtPriceX64 和 openTime 的值，您需要根据实际情况设置
	sqrtPriceX64 := uint128{
		Lower: 100,
		Upper: 1000,
	}
	openTime := uint64(time.Now().Unix() - 60)

	// ---------- 2. 计算指令 Discriminator ----------
	discriminator := getInstructionDiscriminator("create_pool")

	// ---------- 3. 准备指令数据 ----------
	instructionData := CreatePoolInstructionData{
		Discriminator: discriminator,
		SqrtPriceX64:  sqrtPriceX64,
		OpenTime:      openTime,
	}

	// 序列化指令数据
	serializedData, err := borsh.Serialize(instructionData)
	if err != nil {
		log.Fatalf("Failed to serialize instruction data: %v", err)
	}

	// ---------- 4. 计算 PDA ----------
	// pool_state PDA
	poolStateSeeds := [][]byte{
		[]byte("pool"),
		ammConfigPubkey.Bytes(),
		lowerMintPubkey.Bytes(),
		higherMintPubkey.Bytes(),
	}

	poolStatePubkey, _, err := common.FindProgramAddress(poolStateSeeds, programID)
	if err != nil {
		log.Fatalf("Failed to find poolState PDA: %v", err)
	}

	// token_vault_0 PDA
	tokenVault0Seeds := [][]byte{
		[]byte("pool_vault"),
		poolStatePubkey.Bytes(),
		lowerMintPubkey.Bytes(),
	}

	tokenVault0Pubkey, _, err := common.FindProgramAddress(tokenVault0Seeds, programID)
	if err != nil {
		log.Fatalf("Failed to find tokenVault0 PDA: %v", err)
	}

	// token_vault_1 PDA
	tokenVault1Seeds := [][]byte{
		[]byte("pool_vault"),
		poolStatePubkey.Bytes(),
		higherMintPubkey.Bytes(),
	}

	tokenVault1Pubkey, _, err := common.FindProgramAddress(tokenVault1Seeds, programID)
	if err != nil {
		log.Fatalf("Failed to find tokenVault1 PDA: %v", err)
	}

	// observation_state PDA
	observationStateSeeds := [][]byte{
		[]byte("observation"),
		poolStatePubkey.Bytes(),
	}

	observationStatePubkey, _, err := common.FindProgramAddress(observationStateSeeds, programID)
	if err != nil {
		log.Fatalf("Failed to find observationState PDA: %v", err)
	}

	// tick_array_bitmap PDA
	tickArrayBitmapSeeds := [][]byte{
		[]byte("pool_tick_array_bitmap_extension"),
		poolStatePubkey.Bytes(),
	}

	tickArrayBitmapPubkey, _, err := common.FindProgramAddress(tickArrayBitmapSeeds, programID)
	if err != nil {
		log.Fatalf("Failed to find tickArrayBitmap PDA: %v", err)
	}

	// ---------- 5. 准备账户元数据 ----------
	accounts := []types.AccountMeta{
		// pool_creator
		{
			PubKey:     poolCreatorPubkey,
			IsSigner:   true,
			IsWritable: true,
		},
		// amm_config
		{
			PubKey:     ammConfigPubkey,
			IsSigner:   false,
			IsWritable: false,
		},
		// pool_state
		{
			PubKey:     poolStatePubkey,
			IsSigner:   false,
			IsWritable: true,
		},
		// token_mint_0
		{
			PubKey:     lowerMintPubkey,
			IsSigner:   false,
			IsWritable: false,
		},
		// token_mint_1
		{
			PubKey:     higherMintPubkey,
			IsSigner:   false,
			IsWritable: false,
		},
		// token_vault_0
		{
			PubKey:     tokenVault0Pubkey,
			IsSigner:   false,
			IsWritable: true,
		},
		// token_vault_1
		{
			PubKey:     tokenVault1Pubkey,
			IsSigner:   false,
			IsWritable: true,
		},
		// observation_state
		{
			PubKey:     observationStatePubkey,
			IsSigner:   false,
			IsWritable: true,
		},
		// tick_array_bitmap
		{
			PubKey:     tickArrayBitmapPubkey,
			IsSigner:   false,
			IsWritable: true,
		},
		// token_program_0
		{
			PubKey:     common.TokenProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
		// token_program_1
		{
			PubKey:     common.TokenProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
		// system_program
		{
			PubKey:     common.SystemProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
		// rent
		{
			PubKey:     common.SysVarRentPubkey,
			IsSigner:   false,
			IsWritable: false,
		},
	}

	// ---------- 6. 创建指令 ----------
	instruction := types.Instruction{
		ProgramID: programID,
		Accounts:  accounts,
		Data:      serializedData,
	}

	// ---------- 7. 创建并发送交易 ----------
	// 加载签名者账户
	//poolCreatorPrivateKey := []byte{
	//	// 您的私钥字节数组，长度为64字节（32字节私钥 + 32字节公钥）
	//	// 请确保安全地加载您的私钥
	//}

	//poolCreatorAccount, err := types.AccountFromBytes(poolCreatorPrivateKey)
	//if err != nil {
	//	log.Fatalf("Failed to load pool creator account: %v", err)
	//}

	// 创建客户端
	c := NewClient()

	// 获取最近的区块哈希
	recentBlockhashResponse, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("Failed to get recent blockhash: %v", err)
	}
	recentBlockhash := recentBlockhashResponse.Blockhash

	// 创建交易
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        poolCreatorPubkey,
			Instructions:    []types.Instruction{instruction},
			RecentBlockhash: recentBlockhash,
		}),
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

func CreateRaydiumCLPMPool() {

	c := NewClient()
	alice, err := types.AccountFromBase58(Piv)

	if err != nil {
		fmt.Printf("err=%v", err)
	}
	fmt.Printf("alice=%v\n", alice.PublicKey.String())
	//mintKeypair =6GVCVzbiGNfjATr91bWPTmiToZXUkACbyLtNStTmYAK9// 这个就是代币地址
	//mintKeypair2 := common.PublicKeyFromString("6GVCVzbiGNfjATr91bWPTmiToZXUkACbyLtNStTmYAK9")
	//ATAmetadata := common.PublicKeyFromString("Esby6ub71djgjGzhkfujaCQrL7nvLTrsMnFFEG4Jtvm3")
	// 2、组装交易
	// 2.1 拿recentBlockhashResponse
	recentBlockhashResponse, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}

	// 2.2 定义Instruction然后填充它
	ins := make([]types.Instruction, 0, 3)

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
	// DOINGS 调用Raydium合约createPool指令--->Account 列表需要13个地址
	// 0、programID
	programID := common.PublicKeyFromString("devi51mZmdwUJGU9hjN27vEz64Gps7uUefqxg27EAtH")
	// 1,Accounts 第一个，是 pool_creator
	poolCreator := alice.PublicKey
	// 2,Accounts 第2个，是 amm_config
	ammConfig := common.PublicKeyFromString("CQYbhr6amxUER4p5SC44C63R4qw4NFc9Z4Db9vF4tZwG")
	// 3,Accounts 第3个，是 pool_state
	//poolState
	//#[account(
	//init,
	//seeds = [
	//POOL_SEED.as_bytes(),
	//amm_config.key().as_ref(),
	//token_mint_0.key().as_ref(),
	//token_mint_1.key().as_ref(),
	//],
	//bump,
	//payer = pool_creator,
	//space = PoolState::LEN
	//)]
	// 4,Accounts 第4个，token_mint_0 WSOL
	tokenMint0 := common.PublicKeyFromString("So11111111111111111111111111111111111111112")
	// 5,Accounts 第5个，token_mint_1
	tokenMint1 := common.PublicKeyFromString("6GVCVzbiGNfjATr91bWPTmiToZXUkACbyLtNStTmYAK9")
	seeds := [][]byte{}
	seeds = append(seeds, []byte("pool"))
	seeds = append(seeds, ammConfig.Bytes())
	seeds = append(seeds, tokenMint0.Bytes())
	seeds = append(seeds, tokenMint1.Bytes())
	poolState, _, err := common.FindProgramAddress(seeds, programID) // 派生PDA地址原数据
	if err != nil {
		fmt.Printf("err=%v", err)
	}
	// 6,Accounts 第6个，token_vault_0
	seeds2 := [][]byte{}
	seeds2 = append(seeds2, []byte("pool_vault"))
	seeds2 = append(seeds2, poolState.Bytes())
	seeds2 = append(seeds2, tokenMint0.Bytes())
	tokenVault0, _, err := common.FindProgramAddress(seeds2, programID) // 派生PDA地址原数据
	// 7,Accounts 第6个，token_vault_1
	seeds3 := [][]byte{}
	seeds3 = append(seeds3, []byte("pool_vault"))
	seeds3 = append(seeds3, poolState.Bytes())
	seeds3 = append(seeds3, tokenMint1.Bytes())
	tokenVault1, _, err := common.FindProgramAddress(seeds3, programID) // 派生PDA地址原数据
	// 8,Accounts 第6个，observation_state
	seeds4 := [][]byte{}
	seeds4 = append(seeds4, []byte("observation"))
	seeds4 = append(seeds4, poolState.Bytes())
	observationState, _, err := common.FindProgramAddress(seeds3, programID) // 派生PDA地址原数据
	// 9,Accounts 第6个，tick_array_bitmap
	seeds5 := [][]byte{}
	seeds5 = append(seeds5, []byte("pool_tick_array_bitmap_extension"))
	seeds5 = append(seeds5, poolState.Bytes())
	tickArrayBitmap, _, err := common.FindProgramAddress(seeds3, programID) // 派生PDA地址原数据
	// 10,Accounts 第6个，token_program_0
	tokenProgram0 := common.TokenProgramID
	// 11,Accounts 第6个，token_program_1
	tokenProgram1 := common.TokenProgramID
	// 12,Accounts 第6个，system_program
	systemProgram := common.SystemProgramID
	// 13,Accounts 第6个，rent
	rent := common.SysVarRentPubkey

	// 第一次没有创建ATA地址，所以需要创建ATA地址，后面就不用了，这个创建是要上链的
	Accounts := []types.AccountMeta{
		{PubKey: poolCreator, IsSigner: true, IsWritable: true},       // 1 pool_creator
		{PubKey: ammConfig, IsSigner: false, IsWritable: true},        // 2 amm_config
		{PubKey: poolState, IsSigner: false, IsWritable: true},        // 3 pool_state
		{PubKey: tokenMint0, IsSigner: false, IsWritable: false},      // 4 token_mint_0
		{PubKey: tokenMint1, IsSigner: false, IsWritable: false},      // 5 token_mint_1
		{PubKey: tokenVault0, IsSigner: false, IsWritable: true},      // 6 token_vault_0
		{PubKey: tokenVault1, IsSigner: false, IsWritable: true},      // 7 token_vault_1
		{PubKey: observationState, IsSigner: false, IsWritable: true}, // 8 observation_state
		{PubKey: tickArrayBitmap, IsSigner: false, IsWritable: true},  //9 tick_array_bitmap
		{PubKey: tokenProgram0, IsSigner: false, IsWritable: false},   //10 token_program_0
		{PubKey: tokenProgram1, IsSigner: false, IsWritable: false},   //11 token_program_1
		{PubKey: systemProgram, IsSigner: false, IsWritable: false},   //12  system_program
		{PubKey: rent, IsSigner: false, IsWritable: false},            //13  rent
	}
	//anchor.IDL{}
	data, err := bincode.SerializeData(struct {
		Instruction  uint8
		sqrtPriceX64 uint64
		openTime     uint64
	}{
		Instruction:  2,
		sqrtPriceX64: 1,
		openTime:     1234556,
	})
	createPoolIns := types.Instruction{
		ProgramID: programID,
		Accounts:  Accounts,
		Data:      data,
	}
	ins = append(ins, createPoolIns)

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
	// mintKeypair address= 6GVCVzbiGNfjATr91bWPTmiToZXUkACbyLtNStTmYAK9
	// ATAmetadata address  =Esby6ub71djgjGzhkfujaCQrL7nvLTrsMnFFEG4Jtvm3
}
func IncreaseLiquidityV2() {
	// ---------- 1. 设置必要的常量和变量 ----------
	// 替换为您的实际程序 ID
	programID := common.PublicKeyFromString("devi51mZmdwUJGU9hjN27vEz64Gps7uUefqxg27EAtH")

	// 替换为实际的公钥
	alice, err := types.AccountFromBase58(Piv)
	if err != nil {
		fmt.Printf("err=%v", err)
	}
	nftOwner := alice.PublicKey
	ammConfigPubkey := common.PublicKeyFromString("CQYbhr6amxUER4p5SC44C63R4qw4NFc9Z4Db9vF4tZwG")
	tokenMint0Pubkey := common.PublicKeyFromString("So11111111111111111111111111111111111111112")
	//tokenMint1Pubkey := common.PublicKeyFromString("So11111111111111111111111111111111111111112")
	tokenMint1Pubkey := common.PublicKeyFromString("6GVCVzbiGNfjATr91bWPTmiToZXUkACbyLtNStTmYAK9")
	//tokenMint1Pubkey := common.PublicKeyFromString("8w52rwTTDxE8XQgXnBFLfNd9FJjQBUXUHTsq1rvNpump")
	//tokenMint0Pubkey := common.PublicKeyFromString("6GVCVzbiGNfjATr91bWPTmiToZXUkACbyLtNStTmYAK9")

	// 确定 tokenMint0 和 tokenMint1 的顺序
	var lowerMintPubkey, higherMintPubkey common.PublicKey
	if tokenMint0Pubkey.String() < tokenMint1Pubkey.String() {
		lowerMintPubkey = tokenMint0Pubkey
		higherMintPubkey = tokenMint1Pubkey
	} else {
		//lowerMintPubkey = tokenMint1Pubkey
		//higherMintPubkey = tokenMint0Pubkey
		lowerMintPubkey = tokenMint0Pubkey
		higherMintPubkey = tokenMint1Pubkey
	}

	// sqrtPriceX64 和 openTime 的值，您需要根据实际情况设置
	sqrtPriceX64 := uint128{
		Lower: 100,
		Upper: 1000,
	}
	openTime := uint64(time.Now().Unix() - 60)

	// ---------- 2. 计算指令 Discriminator ----------
	discriminator := getInstructionDiscriminator("create_pool")

	// ---------- 3. 准备指令数据 ----------
	instructionData := CreatePoolInstructionData{
		Discriminator: discriminator,
		SqrtPriceX64:  sqrtPriceX64,
		OpenTime:      openTime,
	}

	// 序列化指令数据
	serializedData, err := borsh.Serialize(instructionData)
	if err != nil {
		log.Fatalf("Failed to serialize instruction data: %v", err)
	}

	// ---------- 4. 计算 PDA ----------
	// pool_state PDA
	poolStateSeeds := [][]byte{
		[]byte("pool"),
		ammConfigPubkey.Bytes(),
		lowerMintPubkey.Bytes(),
		higherMintPubkey.Bytes(),
	}

	poolStatePubkey, _, err := common.FindProgramAddress(poolStateSeeds, programID)
	if err != nil {
		log.Fatalf("Failed to find poolState PDA: %v", err)
	}

	// token_vault_0 PDA
	tokenVault0Seeds := [][]byte{
		[]byte("pool_vault"),
		poolStatePubkey.Bytes(),
		lowerMintPubkey.Bytes(),
	}

	tokenVault0Pubkey, _, err := common.FindProgramAddress(tokenVault0Seeds, programID)
	if err != nil {
		log.Fatalf("Failed to find tokenVault0 PDA: %v", err)
	}

	// token_vault_1 PDA
	tokenVault1Seeds := [][]byte{
		[]byte("pool_vault"),
		poolStatePubkey.Bytes(),
		higherMintPubkey.Bytes(),
	}

	tokenVault1Pubkey, _, err := common.FindProgramAddress(tokenVault1Seeds, programID)
	if err != nil {
		log.Fatalf("Failed to find tokenVault1 PDA: %v", err)
	}

	// observation_state PDA
	observationStateSeeds := [][]byte{
		[]byte("observation"),
		poolStatePubkey.Bytes(),
	}

	observationStatePubkey, _, err := common.FindProgramAddress(observationStateSeeds, programID)
	if err != nil {
		log.Fatalf("Failed to find observationState PDA: %v", err)
	}

	// tick_array_bitmap PDA
	tickArrayBitmapSeeds := [][]byte{
		[]byte("pool_tick_array_bitmap_extension"),
		poolStatePubkey.Bytes(),
	}

	tickArrayBitmapPubkey, _, err := common.FindProgramAddress(tickArrayBitmapSeeds, programID)
	if err != nil {
		log.Fatalf("Failed to find tickArrayBitmap PDA: %v", err)
	}

	// ---------- 5. 准备账户元数据 ----------
	accounts := []types.AccountMeta{
		// nftOwner
		{
			PubKey:     nftOwner,
			IsSigner:   true,
			IsWritable: true,
		},
		// nftAccount
		{
			PubKey:     ammConfigPubkey,
			IsSigner:   false,
			IsWritable: false,
		},
		// poolState
		{
			PubKey:     poolStatePubkey,
			IsSigner:   false,
			IsWritable: true,
		},
		// protocolPosition
		{
			PubKey:     lowerMintPubkey,
			IsSigner:   false,
			IsWritable: false,
		},
		// personalPosition
		{
			PubKey:     higherMintPubkey,
			IsSigner:   false,
			IsWritable: false,
		},
		// tickArrayLower
		{
			PubKey:     tokenVault0Pubkey,
			IsSigner:   false,
			IsWritable: true,
		},
		// tickArrayUpper
		{
			PubKey:     tokenVault1Pubkey,
			IsSigner:   false,
			IsWritable: true,
		},
		// tokenAccount0
		{
			PubKey:     observationStatePubkey,
			IsSigner:   false,
			IsWritable: true,
		},
		// tokenAccount1
		{
			PubKey:     tickArrayBitmapPubkey,
			IsSigner:   false,
			IsWritable: true,
		},
		// tokenVault0
		{
			PubKey:     common.TokenProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
		// tokenVault1
		{
			PubKey:     common.TokenProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
		// tokenProgram
		{
			PubKey:     common.SystemProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
		// tokenProgram2022
		{
			PubKey:     common.SysVarRentPubkey,
			IsSigner:   false,
			IsWritable: false,
		},
		// vault0Mint
		{
			PubKey:     common.SysVarRentPubkey,
			IsSigner:   false,
			IsWritable: false,
		},
		// vault1Mint
		{
			PubKey:     common.SysVarRentPubkey,
			IsSigner:   false,
			IsWritable: false,
		},
	}

	// ---------- 6. 创建指令 ----------
	instruction := types.Instruction{
		ProgramID: programID,
		Accounts:  accounts,
		Data:      serializedData,
	}

	// ---------- 7. 创建并发送交易 ----------
	// 加载签名者账户
	//poolCreatorPrivateKey := []byte{
	//	// 您的私钥字节数组，长度为64字节（32字节私钥 + 32字节公钥）
	//	// 请确保安全地加载您的私钥
	//}

	//poolCreatorAccount, err := types.AccountFromBytes(poolCreatorPrivateKey)
	//if err != nil {
	//	log.Fatalf("Failed to load pool creator account: %v", err)
	//}

	// 创建客户端
	c := NewClient()

	// 获取最近的区块哈希
	recentBlockhashResponse, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("Failed to get recent blockhash: %v", err)
	}
	recentBlockhash := recentBlockhashResponse.Blockhash

	// 创建交易
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        alice.PublicKey,
			Instructions:    []types.Instruction{instruction},
			RecentBlockhash: recentBlockhash,
		}),
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

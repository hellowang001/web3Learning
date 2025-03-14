package solwallet

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/pkg/bincode"
	"github.com/blocto/solana-go-sdk/program/compute_budget"
	"github.com/blocto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/blocto/solana-go-sdk/program/system"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/types"
	"log"
)

const (
	//DEV = "https://solana-devnet.g.alchemy.com/v2/wqZxT7UnY6AgrzV42CtGgGQ7ZGM-UrTq"
	Piv = "46M2pAp4z3mNPuTh7jS8XHSn69TC4FAnX33Avjx7wqy3W1zzZKiSoBmNTH5PEBDKu7xR2rPa9ocSyzGWYFK7VRF2"
)

// DOINGS 这个文件主要是官网学习的 sol

func CreateAccountSimpleTx() {
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
	fmt.Printf("recentBlockhashResponse = %v\n", recentBlockhashResponse)

	// 最小转账金额
	minimumBalanceForRentExemption, err := c.GetMinimumBalanceForRentExemption(context.Background(), 0)
	fmt.Printf("minimumBalanceForRentExemption = %v\n", minimumBalanceForRentExemption)
	// 2.2 定义Instruction然后填充它
	ins := make([]types.Instruction, 0, 2)
	// 交易分成2个大部份 Transactions 和 Instructions 指令，Instructions可以有多个
	// 追加优先费
	ComputerUnitPrice := uint64(200000) //一个计算单元的价格
	ComputerUnitLimit := uint32(200000) //计算单元的限制
	ins = append(ins, compute_budget.SetComputeUnitPrice(compute_budget.SetComputeUnitPriceParam{
		MicroLamports: ComputerUnitPrice,
	}))

	ins = append(ins, compute_budget.SetComputeUnitLimit(compute_budget.SetComputeUnitLimitParam{
		Units: ComputerUnitLimit,
	}))

	// DOINGS 我在做的事情是创建一个地址并为其在链上创建一个账号,参考createAccount方法
	// 生成一个钱包，先不管他的地址
	bob := types.NewAccount()
	fmt.Printf("bob address=%s\n", bob.PublicKey.String())
	// 指令三纬度，ProgramId 、 Accounts、Data
	Accounts := []types.AccountMeta{
		{PubKey: common.PublicKeyFromString(alice.PublicKey.String()), IsSigner: true, IsWritable: true},
		{PubKey: common.PublicKeyFromString(bob.PublicKey.String()), IsSigner: true, IsWritable: true},
	}
	data, err := bincode.SerializeData(struct {
		Instruction system.Instruction
		Lamports    uint64
		Space       uint64
		Owner       common.PublicKey
	}{
		Instruction: system.InstructionCreateAccount,
		Lamports:    minimumBalanceForRentExemption + 1,
		//Space:       system.NonceAccountSize,
		Space: 0,
		Owner: common.SystemProgramID,
	})
	createAccountIns := types.Instruction{
		ProgramID: common.SystemProgramID,
		Accounts:  Accounts,
		//Data:      []byte{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 200, 0, 0, 0, 0, 0, 0, 0, 6, 161, 216, 23, 145, 55, 84, 42, 152, 52, 55, 189, 254, 42, 122, 178, 85, 127, 83, 92, 138, 120, 114, 43, 104, 164, 157, 192, 0, 0, 0, 0},
		Data: data,
	}
	ins = append(ins, createAccountIns)
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
func CreateAccountSystemTx() {
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
	fmt.Printf("recentBlockhashResponse = %v\n", recentBlockhashResponse)

	// 最小转账金额
	minimumBalanceForRentExemption, err := c.GetMinimumBalanceForRentExemption(context.Background(), 0)
	fmt.Printf("minimumBalanceForRentExemption = %v\n", minimumBalanceForRentExemption)
	// 2.2 定义Instruction然后填充它
	ins := make([]types.Instruction, 0, 2)
	// 交易分成2个大部份 Transactions 和 Instructions 指令，Instructions可以有多个
	// 追加优先费
	ComputerUnitPrice := uint64(200000) //一个计算单元的价格
	ComputerUnitLimit := uint32(200000) //计算单元的限制
	ins = append(ins, compute_budget.SetComputeUnitPrice(compute_budget.SetComputeUnitPriceParam{
		MicroLamports: ComputerUnitPrice,
	}))

	ins = append(ins, compute_budget.SetComputeUnitLimit(compute_budget.SetComputeUnitLimitParam{
		Units: ComputerUnitLimit,
	}))

	// DOINGS 我在做的事情是创建一个地址并为其在链上创建一个账号,参考createAccount方法
	// 生成一个钱包，先不管他的地址
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
	// mintKeypair address= 6GVCVzbiGNfjATr91bWPTmiToZXUkACbyLtNStTmYAK9
	// ATAmetadata address  =Esby6ub71djgjGzhkfujaCQrL7nvLTrsMnFFEG4Jtvm3
}
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

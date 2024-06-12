package cosmos

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typetx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/proto"
)

// 重写Cosmos 发交易

func transaction2() {
	// 导入你的私钥
	privateKeySecp, _ := Import2("32205066c02dc335b48e6f4d86fb11cae80d61d713472580d6ff7e8020928663")

	// 拿到私钥对应的地址
	var privateKeyAccountAddr sdk.AccAddress = privateKeySecp.PubKey().Address().Bytes()
	fmt.Println(privateKeyAccountAddr)
	fmt.Printf("type=%T\n", privateKeyAccountAddr)

	// 定义交易的接收方地址和代币数量
	toAddr := sdk.MustAccAddressFromBech32("cosmos1962tp3xz6d8mt8pvc405dfrz9hh5z059wc8fq4")
	coin := sdk.NewInt64Coin("uatom", 200)
	amount := sdk.NewCoins(coin)
	fee := sdk.NewInt64Coin("uatom", 1753)
	gasLimt := uint64(100000)

	// 创建一个转移代币的消息，去定义msg
	msg := banktypes.NewMsgSend(privateKeyAccountAddr, toAddr, amount)

	// 建立gRPC客户端
	grpcConn, _ := NewClient()
	authClient := authtypes.NewQueryClient(grpcConn)

	// 查询发送方的信息，Account returns account details based on address
	fromAddrInfo, err := authClient.Account(context.Background(), &authtypes.QueryAccountRequest{Address: privateKeyAccountAddr.String()})
	if err != nil {
		panic(err)
	}
	var AccountInfo authtypes.BaseAccount

	err = proto.Unmarshal(fromAddrInfo.Account.Value, &AccountInfo)
	if err != nil {
		panic(err)
	}

	// 1、要拼txBody出来Missing len argument in the make function
	txBodyMessage := make([]*types.Any, 0)
	msgLen := []sdk.Msg{msg}
	for i := 0; i < len(msgLen); i++ {
		msgAnyValue, err := types.NewAnyWithValue(msgLen[i])
		if err != nil {
			panic(err)
		}
		txBodyMessage = append(txBodyMessage, msgAnyValue)
	}
	txBody := &typetx.TxBody{
		Messages:                    txBodyMessage,
		Memo:                        "hello Cosmos",
		TimeoutHeight:               0,
		ExtensionOptions:            nil,
		NonCriticalExtensionOptions: nil,
	}
	txBodyBytes, err := proto.Marshal(txBody)
	if err != nil {
		panic(err)
	}
	pubAny, err := types.NewAnyWithValue(privateKeySecp.PubKey())
	authInfo := &typetx.AuthInfo{
		SignerInfos: []*typetx.SignerInfo{
			{
				PublicKey: pubAny,
				ModeInfo: &typetx.ModeInfo{
					Sum: &typetx.ModeInfo_Single_{
						Single: &typetx.ModeInfo_Single{Mode: signing.SignMode_SIGN_MODE_DIRECT},
					},
				},
				Sequence: AccountInfo.Sequence,
			},
		},
		Fee: &typetx.Fee{
			Amount:   sdk.NewCoins(fee),
			GasLimit: uint64(gasLimt),
			Payer:    "",
			Granter:  "",
		},
	}
	txAuthInfoBytes, err := proto.Marshal(authInfo)
	if err != nil {
		panic(err)
	}
	chainId := "cosmoshub-4"
	accountNumber := AccountInfo.AccountNumber
	// 2、signDoc
	signDoc := &typetx.SignDoc{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: txAuthInfoBytes,
		ChainId:       chainId,
		AccountNumber: accountNumber,
	}

	// 3、 sign
	signatures, err := proto.Marshal(signDoc)
	if err != nil {
		panic(err)
	}
	sign, err := privateKeySecp.Sign(signatures)
	if err != nil {
		panic(err)
	}
	// 要拼这三个东西出来
	txRaw := &typetx.TxRaw{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: signDoc.AuthInfoBytes,
		Signatures:    [][]byte{sign},
	}
	// 序列化交易，准备广播出去
	txBytes, err := proto.Marshal(txRaw)
	if err != nil {
		panic(err)
	}
	// 广播到链上去
	txClient := typetx.NewServiceClient(grpcConn)
	// 初始化一个请求模型
	req := &typetx.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    typetx.BroadcastMode_BROADCAST_MODE_SYNC, // 同步广播模式
	}
	// 把初始化的广播请求，用初始化的广播响应去接收
	txResp, err := txClient.BroadcastTx(context.Background(), req)
	if err != nil {
		panic(err)
	}
	// 把你的广播后得到的内容打印一下
	fmt.Println("广播后接收到的txResp:", txResp)
}

func Import2(privateKeyStr string) (*secp256k1.PrivKey, error) {
	// 将私钥字符串解码成字节
	privateKeyBytes, err := hex.DecodeString(privateKeyStr)
	if err != nil {
		panic(err)
	}
	// 初始化 secp256k1私钥
	privateKey := &secp256k1.PrivKey{Key: privateKeyBytes}
	return privateKey, err
}

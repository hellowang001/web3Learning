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

//	func NewGrpcClient() (grpcClient *grpc.ClientConn, err error) {
//		// 创建grpc连接
//		target := "cosmos-grpc.publicnode.com:443"
//		var opts []grpc.DialOption
//		opts = append(opts, grpc.WithCredentialsBundle(google.NewDefaultCredentials()))
//		grpcConn, err := grpc.NewClient(target, opts...)
//		if err != nil {
//			panic(err)
//		}
//		txClient := typetx.NewServiceClient(grpcConn)
//
//		req := &typetx.GetBlockWithTxsRequest{Height: 20758858}
//		resp, err := txClient.GetBlockWithTxs(context.Background(), req)
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println(resp.Txs[0].Body)
//		return grpcConn, nil
//	}
//
// Import 根据十六进制字符串导入私钥。
// 参数 privKeyHex 是一个十六进制编码的私钥字符串。
// 返回值是一个指向 secp256k1 私钥结构的指针，以及可能的错误。
// 该函数的目的是将十六进制表示的私钥转换为 secp256k1 私钥结构，以供后续加密操作使用。
func Import(privKeyHex string) (*secp256k1.PrivKey, error) {
	// 将十六进制字符串解码为字节序列
	priBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {

		return nil, err
	}
	// 创建并返回一个包含解码后字节序列的 secp256k1 私钥结构
	return &secp256k1.PrivKey{Key: priBytes}, nil
}

// transaction 示例展示了如何在Cosmos SDK中创建和广播一个基本的交易。
// 它包括生成交易消息、构建交易、签名交易以及最后广播交易。
func transaction() {

	// 导入私钥，用于生成交易的签名。私钥是已经派生完的子私钥
	pk, _ := Import("9f904a195e9a0b8dc6f3305928c7654f7ac714370c4a5d7f97a4ae2b7e06e4f8")
	// 从私钥派生账户地址。
	var privKeyAccAddr sdk.AccAddress = pk.PubKey().Address().Bytes()

	// 定义交易的接收方地址和转移的代币数量。
	toAddr := sdk.MustAccAddressFromBech32("cosmos1sphlm2dp2a4v9hy7fzvqznus26g03vgz006ldt")
	coin := sdk.NewInt64Coin("uatom", 1)
	amount := sdk.NewCoins(coin)

	// 创建一个转移代币的消息。定义msg 实例化msgSend
	msg := banktypes.NewMsgSend(privKeyAccAddr, toAddr, amount)

	// 通过gRPC连接到链的查询客户端。
	grpcConn, _ := NewClient()                       // 会返回之前定义的gRPC连接
	authClient := authtypes.NewQueryClient(grpcConn) // 新建客户端

	// 查询发送方账户信息，以获取序列号和账户号，这些都是构建交易必需的。
	fromAddrInfoAny, err := authClient.Account(context.Background(), &authtypes.QueryAccountRequest{Address: privKeyAccAddr.String()})
	// 解析查询结果中的账户信息。
	var f authtypes.BaseAccount
	//  定义f为一个BaseAccount结构体，把fromAddrInfoAny的用户的信息实例化进对象f，这个时候f就有很多属性了
	if err := proto.Unmarshal(fromAddrInfoAny.Account.Value, &f); err != nil {
		panic(err)
	}

	// 定义交易费用和气体限制。定义gas fee 和 gas limit
	fee := sdk.NewInt64Coin("uatom", 2000)
	gasLimit := int64(100000)

	// 构建交易，包括签名。
	// 这里是构建一个完整的签名后的转账消息，后面广播用
	txRaw, err := BuildTxV2(
		"cosmoshub-4",
		f.Sequence,
		f.AccountNumber,
		pk,
		fee,
		gasLimit,
		[]sdk.Msg{msg},
	)
	if err != nil {
		panic(err)
	}

	// 序列化交易，准备广播。
	txBytes, err := proto.Marshal(txRaw)
	if err != nil {
		panic(err)
	}

	// 广播交易到链上。
	// 广播这笔签名后的交易出去
	txClient := typetx.NewServiceClient(grpcConn)
	req := &typetx.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    typetx.BroadcastMode_BROADCAST_MODE_SYNC,
	}
	txResp, err := txClient.BroadcastTx(context.Background(), req)
	if err != nil {
		panic(err)
	}

	// 打印交易响应，包括交易的哈希和状态码。
	fmt.Println(txResp.TxResponse)
}

// BuildTxV2 构建一笔交易
// BuildTxV2 根据给定参数构建交易v2版本。
// 参数:
//
//	chainId - 区块链ID。
//	sequence - 账户序列号。
//	accountNumber - 账户编号。
//	privKey - 私钥，用于交易签名。
//	fee - 交易费用。
//	gaslimit - 交易燃气限制。
//	msgs - 交易消息数组。
//
// 返回:
//
//	*typetx.TxRaw - 构建的交易原始数据。
//	error - 如果构建过程中出现错误，则返回错误。
func BuildTxV2(chainId string, sequence, accountNumber uint64, privKey *secp256k1.PrivKey, fee sdk.Coin, gaslimit int64, msgs []sdk.Msg) (*typetx.TxRaw, error) {
	// 初始化消息数组，用于存储不同类型的消息。
	txBodyMessage := make([]*types.Any, 0)
	// 遍历消息数组，将每个消息封装为Any类型，并添加到消息数组中。
	for i := 0; i < len(msgs); i++ {
		msgAnyValue, err := types.NewAnyWithValue(msgs[i])
		if err != nil {
			return nil, err
		}
		txBodyMessage = append(txBodyMessage, msgAnyValue)
	}
	// 初始化TxBody结构体，设置消息数组和其他必要字段。
	// 实例化TxBody
	txBody := &typetx.TxBody{
		Messages:                    txBodyMessage,
		Memo:                        "",
		TimeoutHeight:               0,
		ExtensionOptions:            nil,
		NonCriticalExtensionOptions: nil,
	}
	// 序列化TxBody为字节流。
	txBodyBytes, err := proto.Marshal(txBody)
	if err != nil {
		return nil, err
	}
	// 将公钥封装为Any类型。
	pubAny, err := types.NewAnyWithValue(privKey.PubKey())
	if err != nil {
		return nil, err
	}
	// 初始化AuthInfo结构体，设置公钥、签名信息和交易费用。定义交易模型
	authInfo := &typetx.AuthInfo{
		SignerInfos: []*typetx.SignerInfo{
			{
				PublicKey: pubAny,
				ModeInfo: &typetx.ModeInfo{
					Sum: &typetx.ModeInfo_Single_{
						Single: &typetx.ModeInfo_Single{Mode: signing.SignMode_SIGN_MODE_DIRECT},
					},
				},
				Sequence: sequence,
			},
		},
		Fee: &typetx.Fee{
			Amount:   sdk.NewCoins(fee),
			GasLimit: uint64(gaslimit),
			Payer:    "",
			Granter:  "",
		},
	}
	// 序列化AuthInfo为字节流。
	// 将账户信息模型序列化
	txAuthInfoBytes, err := proto.Marshal(authInfo)
	if err != nil {
		return nil, err
	}

	// 初始化SignDoc结构体，包含交易体、授权信息、链ID和账户序列号。
	signDoc := &typetx.SignDoc{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: txAuthInfoBytes,
		ChainId:       chainId,
		AccountNumber: accountNumber,
	}
	// 对SignDoc进行签名。
	// 签名序列化
	signatures, err := proto.Marshal(signDoc)
	if err != nil {
		return nil, err
	}
	// 使用私钥对签名数据进行签名。
	sign, err := privKey.Sign(signatures)
	if err != nil {
		return nil, err
	}
	// 构建TxRaw结构体，包含交易体字节流、授权信息字节流和签名。
	// 把签名后的TxRaw返回出去 下一步再是广播
	return &typetx.TxRaw{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: signDoc.AuthInfoBytes,
		Signatures:    [][]byte{sign},
	}, nil
}

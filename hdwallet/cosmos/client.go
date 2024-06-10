package cosmos

import (
	"context"
	"fmt"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/google"
)

// NewClient 新建链接，用于创建一个 grpc 客户端
func NewClient() (grpcClient *grpc.ClientConn, err error) {
	// gRPC 服务器的地址
	serverAddress := "cosmos-grpc.publicnode.com:443"
	var opts []grpc.DialOption // 创建一个空的opts切片，用于存储gRPC客户端连接选项。
	// 这里的append是奖默认凭证添加到gRPC客户端链接选项中，意味着客户端奖使用默认的凭证进行身份认证
	opts = append(opts, grpc.WithCredentialsBundle(google.NewDefaultCredentials()))

	// 创建一个gRPC客户端连接，传入服务器地址和这个opts选项，也即是客户端凭证
	grpcConn, err := grpc.NewClient(serverAddress, opts...)
	if err != nil {
		panic(err)
	}
	//接下来，代码中创建了一个 typetx.ServiceClient 客户端对象 txClient，用于与服务端的 Service 进行交互。这里的 typetx 可能是根据您的 .proto 文件生成的 Go 代码中的命名空间或包名。
	// 这里才是最后的创建里一个客户端连接，生成的一个客户端类型
	//txClient := typetx.NewServiceClient(grpcConn)

	//然后，代码创建了一个 typetx.GetBlockWithTxsRequest 请求对象 req，并设置了其 Height 属性为要查询的块的高度（在这里是 20802601）。
	//req := &typetx.GetBlockWithTxsRequest{Height: 20802601} // 这里应该是定义一个请求对象，就是说你的请求是什么
	//rep2 := &typetx.GetHeight{}
	//最后，代码通过调用 txClient.GetBlockWithTxs() 方法，传入一个 context 对象默认组和请求对象 req，向服务端发送请求并获取响应。响应对象存储在 resp 变量中，错误信息存储在 err 变量中。
	//resp, err := txClient.GetBlockWithTxs(context.Background(), req) // 这里应该就是发送这个请求对象出去，接收响应
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(resp.Txs[0].Body)

	bankClient := banktypes.NewQueryClient(grpcConn)
	bankReq := &banktypes.QueryBalanceRequest{Address: "cosmos167a4tt9k3ue0rxm2qq4a8pzp4t8ccyt5z26r2d", Denom: "uatom"}
	bankResp, err := bankClient.Balance(context.Background(), bankReq)
	if err != nil {
		panic(err)
	}
	fmt.Println(bankResp.String())
	return grpcConn, nil

}

package cosmos

import (
	"context"
	"fmt"
	typetx "github.com/cosmos/cosmos-sdk/types/tx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NewClient 新建链接
func NewClient() (grpcClient *grpc.ClientConn, err error) {
	// gRPC 服务器的地址
	serverAddress := "cosmos-grpc.publicnode.com:443"
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithCredentialsBundle(credentials.NewComputeEngine()))
	grpcConn, err := grpc.NewClient(serverAddress, opts...)
	if err != nil {
		panic(err)
	}
	txClient := typetx.NewServiceClient(grpcConn)

	req := &typetx.GetBlockWithTxsRequest{Height: 2077888}
	resp, err := txClient.GetBlockWithTxs(context.Background(), req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Txs[0].Body)
	return grpcConn, nil

}

package client

import (
	"context"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/chain/chainpb"
	"google.golang.org/grpc"
)

var (
	_ ChainClient = (*chainClient)(nil)
)

type ChainClient interface {
	//Queries
	GetHead() (*chainpb.GetHeadResponse, error)
	GetStatus() (*chainpb.GetStatusResponse, error)
	GetMetaByHeight(int64) (*chainpb.GetMetaByHeightResponse, error)
	GetConstants() (*chainpb.GetConstantsResponse, error)
}

func NewChainClient(conn *grpc.ClientConn) *chainClient {
	return &chainClient{
		client: chainpb.NewChainServiceClient(conn),
	}
}

type chainClient struct {
	client chainpb.ChainServiceClient
}

func (r *chainClient) GetHead() (*chainpb.GetHeadResponse, error) {
	ctx := context.Background()

	return r.client.GetHead(ctx, &chainpb.GetHeadRequest{})
}

func (r *chainClient) GetStatus() (*chainpb.GetStatusResponse, error) {
	ctx := context.Background()

	return r.client.GetStatus(ctx, &chainpb.GetStatusRequest{})
}

func (r *chainClient) GetMetaByHeight(h int64) (*chainpb.GetMetaByHeightResponse, error) {
	ctx := context.Background()

	return r.client.GetMetaByHeight(ctx, &chainpb.GetMetaByHeightRequest{Height: h})
}

func (r *chainClient) GetConstants() (*chainpb.GetConstantsResponse, error) {
	ctx := context.Background()

	return r.client.GetConstants(ctx, &chainpb.GetConstantsRequest{})
}

package client

import (
	"context"

	"google.golang.org/grpc"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/event/eventpb"
)

var (
	_ EventClient = (*eventClient)(nil)
)

type EventClient interface {
	GetEscrowEventsByHeight(int64) (*eventpb.GetEscrowEventsByHeightResponse, error)
}

func NewEventClient(conn *grpc.ClientConn) EventClient {
	return &eventClient{
		client: eventpb.NewEventServiceClient(conn),
	}
}

type eventClient struct {
	client eventpb.EventServiceClient
}

func (r *eventClient) GetEscrowEventsByHeight(h int64) (*eventpb.GetEscrowEventsByHeightResponse, error) {
	ctx := context.Background()

	return r.client.GetEscrowEventsByHeight(ctx, &eventpb.GetEscrowEventsByHeightRequest{Height: h})
}

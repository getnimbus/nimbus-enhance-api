package infra

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
)

func NewRpcClient(ctx context.Context, host string) (*rpc.Client, func(), error) {
	client, err := rpc.DialContext(ctx, host)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		client.Close()
	}

	return client, cleanup, nil
}

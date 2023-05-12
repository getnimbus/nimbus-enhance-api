package infra

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
)

func NewEthClient(ctx context.Context, host string) (*ethclient.Client, func(), error) {
	client, err := ethclient.DialContext(ctx, host)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		client.Close()
	}

	return client, cleanup, nil
}

package service

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"

	"go-nimeth/internal/infra"
	"go-nimeth/internal/setting"
)

type Chain string

const (
	AvalanceC     Chain = "avalanche-c"
	BSC           Chain = "bsc"
	Ethereum      Chain = "ethereum"
	ArbitrumNitro Chain = "arbitrum-nitro"
	Fantom        Chain = "fantom"
)

// TODO: maybe need to add this to postgres database
// TODO: some chains has strange method so maybe we need to customize the result a little bit
var chainInfos = map[Chain]ChainInfo{
	AvalanceC: {
		Endpoint: "https://open-platform.nodereal.io/5acc4f3c88f640b298c8444013d3bf43/avalanche-c/ext/bc/C/rpc",
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
		},
	},
	BSC: {
		Endpoint: "https://bsc-mainnet.nodereal.io/v1/5acc4f3c88f640b298c8444013d3bf43",
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
		},
	},
	Ethereum: {
		Endpoint: "https://eth-mainnet.nodereal.io/v1/5acc4f3c88f640b298c8444013d3bf43",
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
		},
	},
	ArbitrumNitro: {
		Endpoint: "https://open-platform.nodereal.io/5acc4f3c88f640b298c8444013d3bf43/arbitrum-nitro/",
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
		},
	},
	Fantom: {
		Endpoint: "https://open-platform.nodereal.io/5acc4f3c88f640b298c8444013d3bf43/fantom/",
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
		},
	},
}

func NewChainService() ChainService {
	return &chainService{}
}

type ChainService interface {
	GetLatestBlock(ctx context.Context, chain Chain) (*types.Header, error)
	SearchTransactionHash(ctx context.Context, hash string) (map[Chain]*types.Receipt, error)
}

type chainService struct {
}

func (svc *chainService) GetLatestBlock(ctx context.Context, chain Chain) (*types.Header, error) {
	if chain == "" {
		return nil, fmt.Errorf("missing chain")
	}

	chainInfo, ok := svc.getChain(chain)
	if !ok {
		return nil, setting.ErrNotSupportedChain
	}

	client, cleanup, err := infra.NewRpcClient(ctx, chainInfo.Endpoint)
	if err != nil {
		return nil, setting.ErrClientConnectionFailure
	}
	defer cleanup()

	var head *types.Header
	err = client.CallContext(ctx, &head, chainInfo.Methods["getBlockByNumber"], "latest", false)
	if err == nil && head == nil {
		err = ethereum.NotFound
	}
	return head, err
}

func (svc *chainService) SearchTransactionHash(ctx context.Context, hash string) (map[Chain]*types.Receipt, error) {
	if hash == "" {
		return nil, fmt.Errorf("missing tx hash")
	}

	var res = make(map[Chain]*types.Receipt, 0)
	for chain, chainInfo := range chainInfos {
		err := func() error {
			client, cleanup, err := infra.NewRpcClient(ctx, chainInfo.Endpoint)
			if err != nil {
				return nil
			}
			defer cleanup()

			var r *types.Receipt
			err = client.CallContext(ctx, &r, chainInfo.Methods["getTransactionReceipt"], hash)
			if err == nil {
				if r != nil {
					res[chain] = r
				}
			}
			return nil
		}()
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (svc *chainService) getChain(chain Chain) (ChainInfo, bool) {
	val, ok := chainInfos[chain]
	return val, ok
}

type ChainInfo struct {
	Endpoint string            `json:"endpoint"`
	Methods  map[string]string `json:"methods"`
}

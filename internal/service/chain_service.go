package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tikivn/ultrago/u_logger"
	"golang.org/x/sync/errgroup"

	"go-nimeth/internal/infra"
	"go-nimeth/internal/setting"
)

type Chain string

const (
	BSC           Chain = "bsc"
	Ethereum      Chain = "ethereum"
	Polygon       Chain = "polygon"
	Optimism      Chain = "optimism"
	ArbitrumNova  Chain = "arbitrum"
	AvalanceC     Chain = "avalanche-c"
	ArbitrumNitro Chain = "arbitrum-nitro"
	Fantom        Chain = "fantom"

	// not supported yet
	Solana Chain = "solana"
	Aptos  Chain = "aptos"
	Near   Chain = "near"
	Klaytn Chain = "klaytn"
)

// TODO: make this into env later
var NODEREAL_API_KEY string = "5acc4f3c88f640b298c8444013d3bf43"

// TODO: maybe need to add this to postgres database
// TODO: some chains has strange method so maybe we need to customize the result a little bit
var chainInfos = map[Chain]ChainInfo{
	BSC: {
		Endpoint: fmt.Sprintf("https://bsc-mainnet.nodereal.io/v1/%s", NODEREAL_API_KEY),
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
		},
	},
	Ethereum: {
		Endpoint: fmt.Sprintf("https://eth-mainnet.nodereal.io/v1/%s", NODEREAL_API_KEY),
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
		},
	},
	Polygon: {
		Endpoint: fmt.Sprintf("https://polygon-mainnet.nodereal.io/v1/%s", NODEREAL_API_KEY),
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
		},
	},
	Optimism: {
		Endpoint: fmt.Sprintf("https://opt-mainnet.nodereal.io/v1/%s", NODEREAL_API_KEY),
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
		},
	},
	ArbitrumNova: {
		Endpoint: fmt.Sprintf("https://open-platform.nodereal.io/%s/arbitrum/", NODEREAL_API_KEY),
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
		},
	},
	AvalanceC: {
		Endpoint: fmt.Sprintf("https://open-platform.nodereal.io/%s/avalanche-c/ext/bc/C/rpc", NODEREAL_API_KEY),
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
		},
	},
	ArbitrumNitro: {
		Endpoint: fmt.Sprintf("https://open-platform.nodereal.io/%s/arbitrum-nitro/", NODEREAL_API_KEY),
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
		},
	},
	Fantom: {
		Endpoint: fmt.Sprintf("https://open-platform.nodereal.io/%s/fantom/", NODEREAL_API_KEY),
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
	sync.RWMutex
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
	ctx, logger := u_logger.GetLogger(ctx)
	if hash == "" {
		return nil, fmt.Errorf("missing tx hash")
	}

	var (
		eg, childCtx = errgroup.WithContext(ctx)
		res          = make(map[Chain]*types.Receipt, 0)
	)
	for k, v := range chainInfos {
		chain := k
		chainInfo := v
		eg.Go(func() error {
			client, cleanup, err := infra.NewRpcClient(childCtx, chainInfo.Endpoint)
			if err != nil {
				logger.Errorf("failed to connect rpc client of chain %v", chain)
				return nil
			}
			defer cleanup()

			var r *types.Receipt
			err = client.CallContext(childCtx, &r, chainInfo.Methods["getTransactionReceipt"], hash)
			if err == nil {
				if r != nil {
					svc.Lock()
					res[chain] = r
					svc.Unlock()
				}
			} else if err != nil {
				logger.Errorf("failed to get transaction receipt in chain %v: %v", chain, err)
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
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

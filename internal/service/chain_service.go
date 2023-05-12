package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	nearclient "github.com/eteu-technologies/near-api-go/pkg/client"
	"github.com/eteu-technologies/near-api-go/pkg/client/block"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/portto/solana-go-sdk/client"
	solanarpc "github.com/portto/solana-go-sdk/rpc"
	"github.com/tikivn/ultrago/u_logger"
	"golang.org/x/sync/errgroup"

	"nimbus-enhance-api/internal/infra"
	"nimbus-enhance-api/internal/repo"
	"nimbus-enhance-api/internal/repo/redis"
	"nimbus-enhance-api/internal/setting"
	"nimbus-enhance-api/pkg/encoder"
)

type Chain string

const (
	// EVM chain
	BSC           Chain = "bsc"
	Ethereum      Chain = "ethereum"
	Polygon       Chain = "polygon"
	Optimism      Chain = "optimism"
	ArbitrumNova  Chain = "arbitrum"
	AvalanceC     Chain = "avalanche-c"
	ArbitrumNitro Chain = "arbitrum-nitro"
	Fantom        Chain = "fantom"

	// non EVM chain
	Solana Chain = "solana"
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
			"getTransactionByHash":  "eth_getTransactionByHash",
		},
		IsEVM: true,
	},
	Ethereum: {
		Endpoint: fmt.Sprintf("https://eth-mainnet.nodereal.io/v1/%s", NODEREAL_API_KEY),
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
			"getTransactionByHash":  "eth_getTransactionByHash",
		},
		IsEVM: true,
	},
	Polygon: {
		Endpoint: fmt.Sprintf("https://polygon-mainnet.nodereal.io/v1/%s", NODEREAL_API_KEY),
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
			"getTransactionByHash":  "eth_getTransactionByHash",
		},
		IsEVM: true,
	},
	Optimism: {
		Endpoint: fmt.Sprintf("https://opt-mainnet.nodereal.io/v1/%s", NODEREAL_API_KEY),
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
			"getTransactionByHash":  "eth_getTransactionByHash",
		},
		IsEVM: true,
	},
	ArbitrumNova: {
		Endpoint: fmt.Sprintf("https://open-platform.nodereal.io/%s/arbitrum/", NODEREAL_API_KEY),
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
			"getTransactionByHash":  "eth_getTransactionByHash",
		},
		IsEVM: true,
	},
	AvalanceC: {
		Endpoint: fmt.Sprintf("https://open-platform.nodereal.io/%s/avalanche-c/ext/bc/C/rpc", NODEREAL_API_KEY),
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
			"getTransactionByHash":  "eth_getTransactionByHash",
		},
		IsEVM: true,
	},
	ArbitrumNitro: {
		Endpoint: fmt.Sprintf("https://open-platform.nodereal.io/%s/arbitrum-nitro/", NODEREAL_API_KEY),
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
			"getTransactionByHash":  "eth_getTransactionByHash",
		},
		IsEVM: true,
	},
	Fantom: {
		Endpoint: fmt.Sprintf("https://open-platform.nodereal.io/%s/fantom/", NODEREAL_API_KEY),
		Methods: map[string]string{
			"getBlockByNumber":      "eth_getBlockByNumber",
			"getTransactionReceipt": "eth_getTransactionReceipt",
			"getTransactionByHash":  "eth_getTransactionByHash",
		},
		IsEVM: true,
	},
	Solana: {
		Endpoint: fmt.Sprintf("https://open-platform.nodereal.io/%s/solana/", NODEREAL_API_KEY),
		Methods: map[string]string{
			"getBlockHeight":   "getBlockHeight",
			"getBlockByNumber": "getBlock",
		},
		IsEVM: false,
	},
	Near: {
		Endpoint: fmt.Sprintf("https://open-platform.nodereal.io/%s/near/", NODEREAL_API_KEY),
		Methods: map[string]string{
			"getBlockByNumber":      "block",
			"getTransactionReceipt": "EXPERIMENTAL_receipt",
		},
		IsEVM: false,
	},
	Klaytn: {
		Endpoint: fmt.Sprintf("https://open-platform.nodereal.io/%s/klaytn/", NODEREAL_API_KEY),
		Methods: map[string]string{
			"getBlockByNumber":      "klay_getBlockByNumber",
			"getTransactionReceipt": "klay_getTransactionReceipt",
		},
		IsEVM: false,
	},
}

func NewChainService(
	redisClient *infra.RedisClient,
) ChainService {
	return &chainService{
		redisRepo: redis.NewRedisRepo(redisClient, "chain", time.Minute),
	}
}

type ChainService interface {
	GetLatestBlock(ctx context.Context, chain Chain) (interface{}, error)
	SearchTransactionHash(ctx context.Context, hash string) (map[Chain]interface{}, error)
}

type chainService struct {
	sync.RWMutex
	redisRepo repo.RedisRepo
}

func (svc *chainService) GetLatestBlock(ctx context.Context, chain Chain) (interface{}, error) {
	if chain == "" {
		return nil, fmt.Errorf("missing chain")
	}
	chainInfo, ok := svc.getChain(chain)
	if !ok {
		return nil, setting.ErrNotSupportedChain
	}

	// get data from cache
	data, err := svc.redisRepo.Get(ctx, string(chain))
	if err == nil {
		var res interface{}
		err = json.Unmarshal([]byte(data), &res)
		return res, err
	}

	// if no hit cache then call api
	if chainInfo.IsEVM {
		client, cleanup, err := infra.NewRpcClient(ctx, chainInfo.Endpoint)
		if err != nil {
			return nil, setting.ErrClientConnectionFailure
		}
		defer cleanup()

		var res *types.Header
		err = client.CallContext(ctx, &res, chainInfo.Methods["getBlockByNumber"], "latest", false)
		if err != nil {
			return nil, err
		} else if err == nil && res == nil {
			err = ethereum.NotFound
		}
		_ = svc.redisRepo.Set(ctx, string(chain), res)
		return res, nil
	} else if chain == Solana {
		client := client.NewClient(chainInfo.Endpoint)
		latestBlockNumber, err := client.GetSlot(ctx)
		if err != nil {
			return nil, err
		}

		var (
			maxSupportedTransactionVersion uint8 = 0
			rewards                              = false
		)
		res, err := client.GetBlockWithConfig(
			ctx,
			latestBlockNumber,
			solanarpc.GetBlockConfig{
				Encoding:                       solanarpc.GetBlockConfigEncodingBase64,
				TransactionDetails:             solanarpc.GetBlockConfigTransactionDetailsNone,
				MaxSupportedTransactionVersion: &maxSupportedTransactionVersion,
				Rewards:                        &rewards,
			},
		)
		if err != nil {
			return nil, err
		}
		_ = svc.redisRepo.Set(ctx, string(chain), res)
		return res, nil
	} else if chain == Near {
		client, err := nearclient.NewClient(chainInfo.Endpoint)
		if err != nil {
			return nil, setting.ErrClientConnectionFailure
		}
		res, err := client.BlockDetails(ctx, block.FinalityFinal())
		if err != nil {
			return nil, err
		}
		_ = svc.redisRepo.Set(ctx, string(chain), res)
		return res, nil
	} else { // default
		client, cleanup, err := infra.NewRpcClient(ctx, chainInfo.Endpoint)
		if err != nil {
			return nil, setting.ErrClientConnectionFailure
		}
		defer cleanup()

		var res interface{}
		err = client.CallContext(ctx, &res, chainInfo.Methods["getBlockByNumber"], "latest", false)
		if err != nil {
			return nil, err
		} else if err == nil && res == nil {
			err = ethereum.NotFound
		}
		_ = svc.redisRepo.Set(ctx, string(chain), res)
		return res, nil
	}
}

func (svc *chainService) SearchTransactionHash(ctx context.Context, hash string) (map[Chain]interface{}, error) {
	ctx, logger := u_logger.GetLogger(ctx)
	if hash == "" {
		return nil, fmt.Errorf("missing tx hash")
	}

	// get data from cache
	data, err := svc.redisRepo.Get(ctx, hash)
	if err == nil {
		var res map[Chain]interface{}
		err = json.Unmarshal([]byte(data), &res)
		return res, err
	}

	var (
		eg, childCtx = errgroup.WithContext(ctx)
		res          = make(map[Chain]interface{}, 0)
	)

	// chain deploy EVM always has fixed-length 66 in tx hash
	// https://stackoverflow.com/questions/72772567/how-long-ethereum-hash-length-block-transaction-address
	if strings.HasPrefix(hash, "0x") && len(hash) == 66 {
		for k, v := range chainInfos {
			if !v.IsEVM {
				continue
			}

			chain := k
			chainInfo := v
			eg.Go(func() error {
				client, cleanup, err := infra.NewRpcClient(childCtx, chainInfo.Endpoint)
				if err != nil {
					logger.Errorf("failed to connect rpc client of chain %v: %v", chain, err)
					return nil
				}
				defer cleanup()

				var data = struct {
					Receipt *types.Receipt     `json:"receipt"`
					Tx      *types.Transaction `json:"tx"`
					Extra   txExtraInfo        `json:"extra"`
				}{}
				err = client.CallContext(childCtx, &data.Receipt, chainInfo.Methods["getTransactionReceipt"], hash)
				if err != nil {
					logger.Errorf("failed to get transaction receipt in chain %v: %v", chain, err)
					return nil
				}

				if data.Receipt != nil {
					var tx *rpcTransaction
					err = client.CallContext(childCtx, &tx, chainInfo.Methods["getTransactionByHash"], data.Receipt.TxHash)
					if err != nil {
						logger.Errorf("failed to get transaction in chain %v: %v", chain, err)
						return nil
					}
					if tx != nil {
						data.Tx = tx.tx
						data.Extra.BlockNumber = tx.BlockNumber
						data.Extra.BlockHash = tx.BlockHash
						data.Extra.From = tx.From
					}
				}
				if data.Receipt != nil {
					svc.Lock()
					res[chain] = data
					svc.Unlock()
				}
				return nil
			})
		}
	} else if encoder.IsBase58(hash) && len(hash) == 88 {
		eg.Go(func() error {
			chainInfo := chainInfos[Solana]
			client := client.NewClient(chainInfo.Endpoint)
			r, err := client.GetTransaction(ctx, hash)
			if err == nil {
				if r != nil {
					svc.Lock()
					res[Solana] = r
					svc.Unlock()
				}
			} else if err != nil {
				logger.Errorf("failed to get transaction receipt in chain solana: %v", err)
			}
			return nil
		})
	} else if len(hash) >= 42 {
		chainInfo := chainInfos[Near]
		eg.Go(func() error {
			client, cleanup, err := infra.NewRpcClient(childCtx, chainInfo.Endpoint)
			if err != nil {
				logger.Errorf("failed to connect rpc client of chain %v: %v", Near, err)
				return nil
			}
			defer cleanup()

			// TODO: think about make same interface for this type
			var r interface{}
			err = client.CallContext(childCtx, &r, chainInfo.Methods["getTransactionReceipt"], hash)
			if err != nil && err.(rpc.DataError).ErrorData() != nil {
				switch err.(rpc.DataError).ErrorData().(type) {
				case string:
					logger.Errorf("failed to get transaction receipt in chain near: %v", err)
				default:
					svc.Lock()
					res[Near] = err.(rpc.DataError).ErrorData()
					svc.Unlock()
				}
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	_ = svc.redisRepo.Set(ctx, hash, res)
	_ = svc.redisRepo.Expire(ctx, hash, 30*time.Minute)
	return res, nil
}

func (svc *chainService) getChain(chain Chain) (ChainInfo, bool) {
	val, ok := chainInfos[chain]
	return val, ok
}

type ChainInfo struct {
	Endpoint string            `json:"endpoint"`
	Methods  map[string]string `json:"methods"`
	IsEVM    bool              `json:"is_evm"`
}

type JsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type rpcTransaction struct {
	tx *types.Transaction
	txExtraInfo
}

type txExtraInfo struct {
	BlockNumber *string         `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash    `json:"blockHash,omitempty"`
	From        *common.Address `json:"from,omitempty"`
}

func (tx *rpcTransaction) UnmarshalJSON(msg []byte) error {
	if err := json.Unmarshal(msg, &tx.tx); err != nil {
		return err
	}
	return json.Unmarshal(msg, &tx.txExtraInfo)
}

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	nearclient "github.com/eteu-technologies/near-api-go/pkg/client"
	"github.com/eteu-technologies/near-api-go/pkg/client/block"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/portto/solana-go-sdk/client"
	solanarpc "github.com/portto/solana-go-sdk/rpc"
	"github.com/tikivn/ultrago/u_http_client"
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
	Avalance      Chain = "avalanche"
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
	Avalance: {
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
			"getTransactionByHash":  "klay_getTransactionByHash",
		},
		IsEVM: true,
	},
}

func NewChainService(
	redisClient *infra.RedisClient,
	httpExecutor u_http_client.HttpExecutor,
) ChainService {
	return &chainService{
		redisRepo:       redis.NewRedisRepo(redisClient, "chain", time.Minute),
		httpExecutor:    httpExecutor,
		chainBaseApiKey: "2MLSf73Pki3sFxxSx3ytr2TFenm", // TODO: add to config later
	}
}

type ChainService interface {
	GetLatestBlock(ctx context.Context, chain Chain) (interface{}, error)
	SearchTransactionHash(ctx context.Context, hash string) (map[Chain]interface{}, error)
	CountTotalTxLast24h(ctx context.Context, chain Chain) (int64, error)
}

type chainService struct {
	sync.RWMutex
	redisRepo       repo.RedisRepo
	httpExecutor    u_http_client.HttpExecutor
	chainBaseApiKey string
}

func (svc *chainService) GetLatestBlock(ctx context.Context, chain Chain) (interface{}, error) {
	if chain == "" {
		return nil, fmt.Errorf("missing chain")
	}
	chainInfo, ok := svc.getChain(chain)
	if !ok {
		return nil, setting.ErrNotSupportedChain
	}
	cacheKey := fmt.Sprintf("latest_block:%s", string(chain))

	// get data from cache
	data, err := svc.redisRepo.Get(ctx, cacheKey)
	if err == nil {
		var res interface{}
		err = json.Unmarshal([]byte(data), &res)
		return res, err
	}

	// if no hit cache then call api
	if chain == Solana {
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
		_ = svc.redisRepo.Set(ctx, cacheKey, res)
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
		_ = svc.redisRepo.Set(ctx, cacheKey, res)
		return res, nil
	} else if chain == Klaytn {
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
		_ = svc.redisRepo.Set(ctx, cacheKey, res)
		return res, nil
	} else if chainInfo.IsEVM {
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
		_ = svc.redisRepo.Set(ctx, cacheKey, res)
		return res, nil
	} else { // other non EVM chain
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
		_ = svc.redisRepo.Set(ctx, cacheKey, res)
		return res, nil
	}
}

func (svc *chainService) SearchTransactionHash(ctx context.Context, hash string) (map[Chain]interface{}, error) {
	ctx, logger := u_logger.GetLogger(ctx)
	if hash == "" {
		return nil, fmt.Errorf("missing tx hash")
	}
	cacheKey := fmt.Sprintf("tx_hash:%s", hash)

	// get data from cache
	data, err := svc.redisRepo.Get(ctx, cacheKey)
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
					Receipt map[string]interface{} `json:"receipt"`
					Tx      map[string]interface{} `json:"tx"`
				}{}
				err = client.CallContext(childCtx, &data.Receipt, chainInfo.Methods["getTransactionReceipt"], hash)
				if err != nil {
					logger.Errorf("failed to get transaction receipt in chain %v: %v", chain, err)
					return nil
				}

				if data.Receipt == nil {
					return nil
				}
				txHash, ok := data.Receipt["transactionHash"]
				if ok {
					var tx map[string]interface{}
					err = client.CallContext(childCtx, &tx, chainInfo.Methods["getTransactionByHash"], txHash)
					if err != nil {
						logger.Errorf("failed to get transaction in chain %v: %v", chain, err)
						return nil
					}
					if tx != nil {
						data.Tx = tx
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
	} else if len(hash) >= 42 && len(hash) <= 44 {
		eg.Go(func() error {
			chainInfo := chainInfos[Near]
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

	// wait for goroutine
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// cache result into redis
	_ = svc.redisRepo.Set(ctx, cacheKey, res)
	_ = svc.redisRepo.Expire(ctx, cacheKey, 1*time.Hour)
	return res, nil
}

func (svc *chainService) CountTotalTxLast24h(ctx context.Context, chain Chain) (res int64, err error) {
	ctx, logger := u_logger.GetLogger(ctx)
	cacheKey := fmt.Sprintf("total_tx:%s", string(chain))

	// get data from cache
	cacheData, err := svc.redisRepo.Get(ctx, cacheKey)
	if err == nil {
		var res int64
		err = json.Unmarshal([]byte(cacheData), &res)
		return res, err
	}

	// if no hit cache then call api
	headers := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
		"X-API-KEY":    svc.chainBaseApiKey,
	}
	payload := map[string]interface{}{
		"query": fmt.Sprintf("select count(*) as total from %s.transactions where block_timestamp >= now() - interval 24 hour;", string(chain)),
	}
	client := u_http_client.
		NewRetryHttpClient(svc.httpExecutor, 60*time.Second, 3).
		WithUrl("https://api.chainbase.online/v1/dw/query", nil).
		WithHeaders(headers).
		WithPayload(payload)
	resp, err := client.Do(ctx, http.MethodPost)
	if err != nil {
		logger.Errorf("failed to get data from api chainbase: %v", err)
		return 0, err
	}

	var data = struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			TaskID    string  `json:"task_id"`
			Rows      int     `json:"rows"`
			RowsRead  int     `json:"rows_read"`
			BytesRead int     `json:"bytes_read"`
			Elapsed   float64 `json:"elapsed"`
			Meta      []struct {
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"meta"`
			Result []struct {
				Total string `json:"total"`
			} `json:"result"`
			ErrMsg string `json:"err_msg"`
		} `json:"data"`
	}{}
	if err := json.Unmarshal(resp.Payload, &data); err != nil {
		logger.Errorf("failed to unmarshal response payload: %v", err)
		return 0, err
	}
	if len(data.Data.Result) > 0 {
		res, err = strconv.ParseInt(data.Data.Result[0].Total, 10, 64)
	}
	if err == nil {
		// cache result into redis
		_ = svc.redisRepo.Set(ctx, cacheKey, res)
		_ = svc.redisRepo.Expire(ctx, cacheKey, 5*time.Minute)
	}
	return
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

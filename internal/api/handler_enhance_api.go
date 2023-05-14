package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tikivn/ultrago/u_handler"
	"github.com/tikivn/ultrago/u_logger"

	"nimbus-enhance-api/internal/service"
)

func NewEnhanceApiHandler(
	baseHandler *u_handler.BaseHandler,
	chainSvc service.ChainService,
) *EnhanceApiHandler {
	return &EnhanceApiHandler{
		BaseHandler: baseHandler,
		chainSvc:    chainSvc,
	}
}

type EnhanceApiHandler struct {
	*u_handler.BaseHandler
	chainSvc service.ChainService
}

func (h *EnhanceApiHandler) Route() chi.Router {
	mux := chi.NewRouter()
	mux.Get("/blocks/latest/{chain}", h.handlerGetLatestBlockByChain)
	mux.Get("/tx/total/{chain}", h.handlerCountTotalTxByChain)
	mux.Get("/tx/{hash}", h.handlerSearchTxHash)
	return mux
}

func (h *EnhanceApiHandler) handlerGetLatestBlockByChain(w http.ResponseWriter, r *http.Request) {
	var (
		ctx, logger = u_logger.GetLogger(r.Context())
		chain       = chi.URLParam(r, "chain")
	)

	if chain == "" {
		logger.Errorf("missing chain")
		h.BadRequest(w, r, fmt.Errorf("missing chain"))
		return
	}

	res, err := h.chainSvc.GetLatestBlock(ctx, service.Chain(chain))
	if err != nil {
		logger.Errorf("failed to get latest block of chain %s: %v", chain, err)
		h.Internal(w, r, fmt.Errorf("failed to get latest block of chain %s: %v", chain, err))
		return
	}
	h.Success(w, r, res)
}

func (h *EnhanceApiHandler) handlerCountTotalTxByChain(w http.ResponseWriter, r *http.Request) {
	var (
		ctx, logger = u_logger.GetLogger(r.Context())
		chain       = chi.URLParam(r, "chain")
	)

	if chain == "" {
		logger.Errorf("missing chain")
		h.BadRequest(w, r, fmt.Errorf("missing chain"))
		return
	}

	chainName := service.Chain(chain)
	switch chainName {
	case service.ArbitrumNitro,
		service.Near,
		service.Klaytn:
		logger.Errorf("not supported counting tx")
		h.BadRequest(w, r, fmt.Errorf("not supported counting tx"))
		return
	default:
	}

	res, err := h.chainSvc.CountTotalTxLast24h(ctx, chainName)
	if err != nil {
		logger.Errorf("failed to count total tx of chain %s: %v", chain, err)
		h.Internal(w, r, fmt.Errorf("failed to get latest block of chain %s: %v", chain, err))
		return
	}
	h.Success(w, r, res)
}

func (h *EnhanceApiHandler) handlerSearchTxHash(w http.ResponseWriter, r *http.Request) {
	var (
		ctx, logger = u_logger.GetLogger(r.Context())
		hash        = chi.URLParam(r, "hash")
	)

	if hash == "" {
		logger.Errorf("missing tx hash")
		h.BadRequest(w, r, fmt.Errorf("missing tx hash"))
		return
	}

	res, err := h.chainSvc.SearchTransactionHash(ctx, hash)
	if err != nil {
		logger.Errorf("failed to get tx hash %s: %v", hash, err)
		h.Internal(w, r, fmt.Errorf("failed to get tx hash %s: %v", hash, err))
		return
	} else if len(res) == 0 {
		logger.Errorf("not found tx hash %s", hash)
		h.NotFound(w, r, fmt.Errorf("not found tx hash %s", hash))
		return
	}
	h.Success(w, r, res)
}

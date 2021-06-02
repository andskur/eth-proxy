package handlers

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-openapi/runtime/middleware"

	"eth-proxy/api/format"
	ops "eth-proxy/api/server/operations/txs"
	"eth-proxy/lib"
	"eth-proxy/pkg/logger"
)

// NewGetTxByHash create new getTxByHash
// handler instance
func NewGetTxByHash(api lib.IEthProxyAPI) *getTxByHash {
	return &getTxByHash{api: api}
}

// getTxByHash is getTxByHash handler
// structure with grpcLib instance
type getTxByHash struct {
	api lib.IEthProxyAPI
}

// Handle represent getTxByHash endpoint handler function
func (g *getTxByHash) Handle(params ops.GetTxByHashParams) middleware.Responder {
	tx, err := g.api.TxByHash(common.HexToHash(params.Hash))
	if err != nil {
		logger.Log().Error(err)
		return ops.NewGetTxByHashDefault(http.StatusInternalServerError).WithPayload(handleError(errInternalServer))
	}

	return ops.NewGetTxByHashOK().WithPayload(format.Tx(tx))
}

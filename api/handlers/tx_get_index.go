package handlers

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"eth-proxy/api/format"
	ops "eth-proxy/api/server/operations/blocks"
	"eth-proxy/lib"
	"eth-proxy/pkg/logger"
)

// NewGetTxByIndex create new getTxByIndex
// handler instance
func NewGetTxByIndex(api lib.IEthProxyAPI) *getTxByIndex {
	return &getTxByIndex{api: api}
}

// getTxByIndex is getTxByIndex handler
// structure with grpcLib instance
type getTxByIndex struct {
	api lib.IEthProxyAPI
}

// Handle represent getTxByIndex endpoint handler function
func (g *getTxByIndex) Handle(params ops.GetTxByIndexParams) middleware.Responder {
	blockNumber, err := parseBlockNumber(params.Number)
	if err != nil {
		logger.Log().Error(err)
		return ops.NewGetTxByIndexDefault(http.StatusBadRequest).WithPayload(handleError(errbadRequest))
	}

	tx, err := g.api.TxByIndex(blockNumber, int(params.Index))
	if err != nil {
		logger.Log().Error(err)
		return ops.NewGetTxByIndexDefault(http.StatusInternalServerError).WithPayload(handleError(errInternalServer))
	}

	return ops.NewGetTxByIndexOK().WithPayload(format.Tx(tx))
}

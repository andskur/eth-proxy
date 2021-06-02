package handlers

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"eth-proxy/api/format"
	ops "eth-proxy/api/server/operations/blocks"
	"eth-proxy/lib"
	"eth-proxy/pkg/logger"
)

// NewGetBlock create new getBlock
// handler instance
func NewGetBlock(api lib.IEthProxyAPI) *getBlock {
	return &getBlock{api: api}
}

// getBlock is getBlock handler
// structure with grpcLib instance
type getBlock struct {
	api lib.IEthProxyAPI
}

// Handle represent getBlock endpoint handler function
func (g *getBlock) Handle(params ops.GetBlockParams) middleware.Responder {
	blockNumber, err := parseBlockNumber(params.Number)
	if err != nil {
		logger.Log().Error(err)
		return ops.NewGetBlockDefault(http.StatusBadRequest).WithPayload(handleError(errbadRequest))
	}

	block, err := g.api.Block(blockNumber)
	if err != nil {
		logger.Log().Error(err)
		return ops.NewGetBlockDefault(http.StatusInternalServerError).WithPayload(handleError(errInternalServer))
	}

	return ops.NewGetBlockOK().WithPayload(format.Block(block))
}

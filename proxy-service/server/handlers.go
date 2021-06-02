package server

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"eth-proxy/models"
	"eth-proxy/proto"
	"eth-proxy/proxy-service/service"
)

// EthProxyAPI represent all
// GRPC Api methods handlers
type EthProxyAPI struct {
	srv service.IService
}

// NewEthProxyAPI crete new EthProxyAPI instance
func NewEthProxyAPI(srv service.IService) *EthProxyAPI {
	return &EthProxyAPI{srv: srv}
}

// Block is Grpc method implementation
// to fetch Block by given getter params
func (e *EthProxyAPI) Block(context context.Context, getter *proto.BlockGetter) (*proto.Block, error) {
	block, err := e.srv.Block(parseBlockGetter(getter))
	if err != nil {
		return nil, fmt.Errorf("get block: %w", err)
	}

	return block.Proto(), nil
}

// Tx is Grpc method implementation
// to fetch Tx by given getter params
func (e *EthProxyAPI) Tx(context context.Context, getter *proto.TxGetter) (pbTx *proto.Tx, err error) {
	var tx *models.Tx

	switch getBy := getter.GetGetter().(type) {
	case *proto.TxGetter_Hash:
		tx, err = e.srv.TxByHash(common.BytesToHash(getBy.Hash))
		if err != nil {
			return nil, fmt.Errorf("get tx by hash: %w", err)
		}
	case *proto.TxGetter_Index:
		tx, err = e.srv.TxByIndex(parseBlockGetter(getBy.Index.Block), int(getBy.Index.Index))
		if err != nil {
			return nil, fmt.Errorf("get tx by index: %w", err)
		}
	}

	return tx.Proto(), nil
}

// parseBlockGetter parsing oneOf value
// of proto Block getter
func parseBlockGetter(getter *proto.BlockGetter) *int {
	switch getBlock := getter.Getter.(type) {
	case *proto.BlockGetter_Number:
		num := int(getBlock.Number)
		return &num
	default:
		return nil
	}
}

package models

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"eth-proxy/proto"
)

// Tx represent Ethereum Transaction model
// structure with most important data fields
type Tx struct {
	Hash     common.Hash
	From     common.Address
	To       common.Address
	Value    int64
	Gas      uint64
	GasPrice int64
	gethTx   *types.Transaction
}

// GethTx return Geth transaction from model
func (t Tx) GethTx() *types.Transaction {
	return t.gethTx
}

// Proto makes Tx model
// formatting to Proto message
func (t Tx) Proto() *proto.Tx {
	return &proto.Tx{
		Hash:     t.Hash.Bytes(),
		From:     t.From.Bytes(),
		To:       t.To.Bytes(),
		Value:    t.Value,
		Gas:      int64(t.Gas),
		GasPrice: t.GasPrice,
	}
}

// TxFromGeth creates Tx model
// structure from geth package Transaction
func TxFromGeth(gethTx *types.Transaction) *Tx {
	return &Tx{
		Hash:     gethTx.Hash(),
		To:       *gethTx.To(),
		Value:    gethTx.Value().Int64(),
		Gas:      gethTx.Gas(),
		GasPrice: gethTx.GasPrice().Int64(),
		gethTx:   gethTx,
	}
}

// TxsFromGeth creates Tx slice from given geth txs
func TxsFromGeth(gethTxs []*types.Transaction) (txs []*Tx) {
	for _, tx := range gethTxs {
		txs = append(txs, TxFromGeth(tx))
	}
	return
}

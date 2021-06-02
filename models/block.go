package models

import (
	"time"

	"eth-proxy/proto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Block represent Ethereum Chain Block
// structure with mostly important data fields
type Block struct {
	Number       uint64
	Hash         common.Hash
	Parent       common.Hash
	Time         time.Time
	TxCount      int
	Transactions []*Tx
}

// Proto makes Block model
// formatting to Proto message
func (b Block) Proto() *proto.Block {
	return &proto.Block{
		Number:    int64(b.Number),
		Hash:      b.Hash.Bytes(),
		Parent:    b.Parent.Bytes(),
		Timestamp: b.Time.Unix(),
		TxCount:   int64(b.TxCount),
	}
}

// BlockFromProto create Block model
// from protobuf message
func BlockFromProto(pb *proto.Block) *Block {
	return &Block{
		Number:  uint64(pb.Number),
		Hash:    common.BytesToHash(pb.Hash),
		Parent:  common.BytesToHash(pb.Parent),
		Time:    time.Unix(pb.Timestamp, 0),
		TxCount: int(pb.TxCount),
	}
}

// BlockFromGeth creates Block model
// structure from geth package Block
func BlockFromGeth(gethBlock *types.Block) *Block {
	return &Block{
		Number:       gethBlock.NumberU64(),
		Hash:         gethBlock.Hash(),
		Parent:       gethBlock.ParentHash(),
		Time:         time.Unix(int64(gethBlock.Time()), 0),
		TxCount:      len(gethBlock.Transactions()),
		Transactions: TxsFromGeth(gethBlock.Transactions()),
	}
}

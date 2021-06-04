package repository

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"eth-proxy/models"
)

// IRepository represent cache repository patter
// to store, select and delete block from cache
type IRepository interface {
	// SelectLastBlock select latest mined
	// Eth block from cache
	SelectLastBlock() *models.Block

	// StoreLastBlock store latest mined
	// Eth block in cache
	StoreLastBlock(block *models.Block)

	// SelectBlock select Eth block from cache
	// by given number
	SelectBlock(number uint64) *models.Block

	// StoreBlock store Eth block in cache
	StoreBlock(block *models.Block)

	// SelectTx select Eth transaction from
	// cache by given hash
	SelectTx(hash common.Hash) *models.Tx

	// StoreTx store Eth transaction in cache
	StoreTx(tx *models.Tx)
}

// inMemory is IRepository implementation that
// stores cache in-memory
type inMemory struct {
	mu        sync.Mutex
	lastBlock *models.Block
	blocks    IBlocksCache
	txs       ITxsCache
}

// NewInMemory create new inMemory IRepository
// cache implementation
func NewInMemory(cacheSize int) IRepository {
	return &inMemory{
		blocks: NewBlocksCache(cacheSize),
		txs:    NewTxsCache(cacheSize),
	}
}

// SelectLastBlock select latest mined
// Eth block from cache
func (i *inMemory) SelectLastBlock() *models.Block {
	i.mu.Lock()
	defer i.mu.Unlock()

	return i.lastBlock
}

// StoreLastBlock store latest mined
// Eth block in cache
func (i *inMemory) StoreLastBlock(block *models.Block) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.lastBlock = block
}

// SelectBlock select Eth block from cache
// by given number
func (i *inMemory) SelectBlock(number uint64) *models.Block {
	return i.blocks.Get(number)
}

// StoreBlock store Eth block in cache
func (i *inMemory) StoreBlock(block *models.Block) {
	i.blocks.Put(block)
}

// SelectTx select Eth transaction from
// cache by given hash
func (i *inMemory) SelectTx(hash common.Hash) *models.Tx {
	return i.txs.Get(hash)
}

// StoreTx store Eth transaction in cache
func (i *inMemory) StoreTx(tx *models.Tx) {
	i.txs.Put(tx)
}

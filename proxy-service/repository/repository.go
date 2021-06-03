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

	// DeleteBLock delete Eth block from cache
	DeleteBLock(number uint64)

	// SelectTx select Eth transaction from
	// cache by given hash
	SelectTx(hash common.Hash) *models.Tx

	// StoreTx store Eth transaction in cache
	StoreTx(tx *models.Tx)

	// DeleteTx delete Eth Tx from cache
	DeleteTx(hash common.Hash)
}

// inMemory is IRepository implementation that
// stores cache in-memory
type inMemory struct {
	mu        sync.Mutex
	lastBlock *models.Block
	blocks    map[uint64]*models.Block
	txs       map[common.Hash]*models.Tx
}

// NewInMemory create new inMemory IRepository
// cache implementation
func NewInMemory() IRepository {
	return &inMemory{
		blocks: make(map[uint64]*models.Block),
		txs:    make(map[common.Hash]*models.Tx),
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
	i.mu.Lock()
	defer i.mu.Unlock()

	block, ok := i.blocks[number]
	if !ok {
		return nil
	}

	return block
}

// StoreBlock store Eth block in cache
func (i *inMemory) StoreBlock(block *models.Block) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.blocks[block.Number] = block
}

// DeleteBLock delete Eth block from cache
func (i *inMemory) DeleteBLock(number uint64) {
	i.mu.Lock()
	defer i.mu.Unlock()

	delete(i.blocks, number)
}

// SelectTx select Eth transaction from
// cache by given hash
func (i *inMemory) SelectTx(hash common.Hash) *models.Tx {
	i.mu.Lock()
	defer i.mu.Unlock()

	tx, ok := i.txs[hash]
	if !ok {
		return nil
	}

	return tx
}

// StoreTx store Eth transaction in cache
func (i *inMemory) StoreTx(tx *models.Tx) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.txs[tx.Hash] = tx
}

// DeleteTx delete Eth Tx from cache
func (i *inMemory) DeleteTx(hash common.Hash) {
	i.mu.Lock()
	defer i.mu.Unlock()

	delete(i.txs, hash)
}

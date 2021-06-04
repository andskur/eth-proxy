package repository

import (
	"container/heap"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"eth-proxy/models"
	"eth-proxy/pkg/logger"
)

// ITxsCache represent interface for
// Eth Transactions Cache
type ITxsCache interface {
	// Put put Tx to cache and check cache size
	Put(tx *models.Tx)

	// Get get Tx from cache by given hash and
	// update its priority
	Get(hash common.Hash) *models.Tx
}

// TxsCache represent ITxsCache interface
// implementation with built-in priority queue
type TxsCache struct {
	mu           sync.Mutex
	txs          map[common.Hash]*tx
	maxCacheSize int
	queue        txPriorityQueue
}

// NewTxsCache create new TxsCache instance
func NewTxsCache(cacheSize int) ITxsCache {
	cache := &TxsCache{
		mu:           sync.Mutex{},
		txs:          make(map[common.Hash]*tx),
		maxCacheSize: cacheSize,
	}

	heap.Init(&cache.queue)
	return cache
}

// Put put Tx to cache and check cache size
func (b *TxsCache) Put(tx *models.Tx) {
	b.mu.Lock()
	defer b.mu.Unlock()

	cachedTx := wrapTx(tx)

	heap.Push(&b.queue, cachedTx)

	b.txs[tx.Hash] = cachedTx

	logger.Log().Infof("tx %s stored to cache", tx.Hash)

	b.checkCacheSize()
}

// Get get Tx from cache and update its priority
func (b *TxsCache) Get(hash common.Hash) *models.Tx {
	b.mu.Lock()
	defer b.mu.Unlock()

	tx, ok := b.txs[hash]
	if !ok {
		return nil
	}

	b.queue.update(tx)

	logger.Log().Infof("tx %s retrived from cache", tx.Hash.Hex())

	return tx.Tx
}

// checkCacheSize check current cache
// size and free it if needed
func (b *TxsCache) checkCacheSize() {
	if len(b.txs) > b.maxCacheSize {
		logger.Log().Warning("cache is full and need to be free")
		b.freeCache()
	}
}

// freeCache free cache with one element
// with lesser priority
func (b *TxsCache) freeCache() {
	tx := heap.Pop(&b.queue).(*tx)
	delete(b.txs, tx.Hash)
	logger.Log().Infof("tx %s removed from cache", tx.Hash)
}

// txPriorityQueue is
type txPriorityQueue []*tx

// Len is the number of elements in the collection
func (q txPriorityQueue) Len() int {
	return len(q)
}

// Less reports whether the element with index i
func (q txPriorityQueue) Less(i, j int) bool {
	return q[i].priority < q[j].priority
}

// Swap swaps the elements with indexes i and j
func (q txPriorityQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

// Push pushes the element x onto the query
func (q *txPriorityQueue) Push(x interface{}) {
	n := len(*q)
	item := x.(*tx)
	item.index = n
	*q = append(*q, item)
}

// Pop removes and returns the minimum
// element from the queue
func (q *txPriorityQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	item.index = -1
	*q = old[0 : n-1]

	return item
}

// update modifies the priority of
// an Item in the queue.
func (q *txPriorityQueue) update(item *tx) {
	item.priority++
	heap.Fix(q, item.index)
}

// tx represent cached wrapper
// for Eth Txs
type tx struct {
	value    common.Hash
	priority int
	index    int
	*models.Tx
}

// wrapTx is wrap Tx model for cache
func wrapTx(t *models.Tx) *tx {
	return &tx{
		value:    t.Hash,
		Tx:       t,
		priority: 1,
	}
}

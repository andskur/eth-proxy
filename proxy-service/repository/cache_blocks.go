package repository

import (
	"container/heap"
	"sync"

	"eth-proxy/models"
	"eth-proxy/pkg/logger"
)

// IBlocksCache represent interface for
// Eth Blocks Cache
type IBlocksCache interface {
	// Put put Block to cache and check cache size
	Put(block *models.Block)

	// Get get Block from cache by given number
	// and update its priority
	Get(number uint64) *models.Block
}

// BlocksCache represent IBlocksCache interface
// implementation with built-in priority queue
type BlocksCache struct {
	mu           sync.Mutex
	blocks       map[uint64]*block
	maxCacheSize int
	queue        blockPriorityQueue
}

// NewBlocksCache create new BlocksCache instance
func NewBlocksCache(cacheSize int) IBlocksCache {
	cache := &BlocksCache{
		mu:           sync.Mutex{},
		blocks:       make(map[uint64]*block),
		maxCacheSize: cacheSize,
	}

	heap.Init(&cache.queue)
	return cache
}

// Put put Block to cache and check cache size
func (b *BlocksCache) Put(block *models.Block) {
	b.mu.Lock()
	defer b.mu.Unlock()

	cachedBlock := wrapBlock(block)

	heap.Push(&b.queue, cachedBlock)

	b.blocks[block.Number] = cachedBlock

	logger.Log().Infof("block %d stored to cache", block.Number)

	b.checkCacheSize()
}

// Get get Block from cache and update its priority
func (b *BlocksCache) Get(number uint64) *models.Block {
	b.mu.Lock()
	defer b.mu.Unlock()

	block, ok := b.blocks[number]
	if !ok {
		return nil
	}

	b.queue.update(block)

	logger.Log().Infof("block %d retrived from cache", block.Number)

	return block.Block
}

// checkCacheSize check current cache
// size and free it if needed
func (b *BlocksCache) checkCacheSize() {
	if len(b.blocks) > b.maxCacheSize {
		logger.Log().Warning("cache is full and need to be free")
		b.freeCache()
	}
}

// freeCache free cache with one element
// with lesser priority
func (b *BlocksCache) freeCache() {
	block := heap.Pop(&b.queue).(*block)
	delete(b.blocks, block.Number)
	logger.Log().Infof("block %d removed from cache", block.Number)
}

// blockPriorityQueue is
type blockPriorityQueue []*block

// Len is the number of elements in the collection
func (q blockPriorityQueue) Len() int {
	return len(q)
}

// Less reports whether the element with index i
func (q blockPriorityQueue) Less(i, j int) bool {
	return q[i].priority < q[j].priority
}

// Swap swaps the elements with indexes i and j
func (q blockPriorityQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

// Push pushes the element x onto the query
func (q *blockPriorityQueue) Push(x interface{}) {
	n := len(*q)
	item := x.(*block)
	item.index = n
	*q = append(*q, item)
}

// Pop removes and returns the minimum
// element from the queue
func (q *blockPriorityQueue) Pop() interface{} {
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
func (q *blockPriorityQueue) update(item *block) {
	item.priority++
	heap.Fix(q, item.index)
}

// block represent cached wrapper
// for Eth Blocks
type block struct {
	value    uint64
	priority int
	index    int
	*models.Block
}

// wrapBlock is wrap Block model for cache
func wrapBlock(b *models.Block) *block {
	return &block{
		value:    b.Number,
		Block:    b,
		priority: 1,
	}
}

package service

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"eth-proxy/models"
	"eth-proxy/pkg/logger"
	"eth-proxy/proxy-service/repository"
)

// IService is interface that represent Application
// service layer with all business logic methods
type IService interface {
	// Block fetch Ethereum Block by given block number
	// you need to pass nil value as a number argument
	// if you want to fetch latest block
	Block(number *int) (*models.Block, error)

	// TxByIndex fetch Ethereum transaction by related
	// block number and transaction index in the block
	TxByIndex(blockNumber *int, txIndex int) (*models.Tx, error)

	// TxByHash fetch Ethereum transaction by given hash
	TxByHash(hash common.Hash) (*models.Tx, error)

	// ListenNewBlocks listen new mined blocks
	ListenNewBlocks()

	// StopListenNewBlocks stop listen new mined blocks
	StopListenNewBlocks()
}

// Service implements Application service layer interface
// with built-in cache repository
type Service struct {
	ethClient *ethclient.Client
	chainID   *big.Int
	cache     repository.IRepository

	stopChan     chan struct{}
	newBlocks    chan *types.Header
	subscription ethereum.Subscription
}

// NewService create new Service instance
// with given ethereum client
func NewService(ethClient *ethclient.Client, wss bool) (IService, error) {
	chainID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("fetch Ethereum chain id: %w", err)
	}

	srv := &Service{
		ethClient: ethClient,
		chainID:   chainID,
		cache:     repository.NewInMemory(),
	}

	if wss {
		srv.stopChan = make(chan struct{})
		srv.newBlocks = make(chan *types.Header)
		if err := srv.subscribeToBlocks(); err != nil {
			return nil, fmt.Errorf("open new blocks subscription: %w", err)
		}
	}
	return srv, nil
}

// subscribeToBlocks open subscription to new chain blocks in realtime
func (s *Service) subscribeToBlocks() (err error) {
	s.subscription, err = s.ethClient.SubscribeNewHead(context.Background(), s.newBlocks)
	if err != nil {
		return fmt.Errorf("subscribe to new headers: %w", err)
	}

	return nil
}

// Block fetch Ethereum Block by given block number
// you need to pass nil value as a number argument
// if you want to fetch latest block
func (s *Service) Block(number *int) (*models.Block, error) {
	num := bigIntFromPointerInt(number)

	cachedBlock := s.retrieveBlockFromHash(num)
	if cachedBlock != nil {
		logger.Log().Infof("Block %d retrived from cache", cachedBlock.Number)
		return cachedBlock, nil
	}

	block, err := s.ethClient.BlockByNumber(context.Background(), num)
	if err != nil {
		return nil, fmt.Errorf("fetch block from Ethereum node: %w", err)
	}

	formattedBlock := models.BlockFromGeth(block)

	if cachedBlock == nil {
		s.cache.StoreBlock(formattedBlock)
		logger.Log().Infof("Block %d stored to cache", formattedBlock.Number)
	}

	return formattedBlock, nil
}

// TxByIndex fetch Ethereum transaction by related
// block number and transaction index in the block
func (s *Service) TxByIndex(blockNumber *int, txIndex int) (*models.Tx, error) {
	block, err := s.Block(blockNumber)
	if err != nil {
		return nil, fmt.Errorf("get block: %w", err)
	}

	if txIndex > block.TxCount {
		return nil, fmt.Errorf("block %d don't have tx with index %d", block.Number, txIndex)
	}

	tx := block.Transactions[txIndex]

	if err := s.fetchTxSender(tx); err != nil {
		return nil, fmt.Errorf("fetch tx sender: %w", err)
	}

	return tx, nil
}

// TxByHash fetch Ethereum transaction by given hash
func (s *Service) TxByHash(hash common.Hash) (*models.Tx, error) {
	cachedTx := s.cache.SelectTx(hash)
	if cachedTx != nil {
		logger.Log().Infof("Tx %s retrived from cache", hash)
		return cachedTx, nil
	}

	txGeth, _, err := s.ethClient.TransactionByHash(context.Background(), hash)
	if err != nil {
		return nil, fmt.Errorf("fetch transaction by hash from Ethereum node: %w", err)
	}

	tx := models.TxFromGeth(txGeth)

	if err := s.fetchTxSender(tx); err != nil {
		return nil, fmt.Errorf("fetch tx sender: %w", err)
	}

	if cachedTx == nil {
		s.cache.StoreTx(tx)
		logger.Log().Infof("Tx %s stored to cache", hash)
	}

	return tx, nil
}

// ListenNewBlocks listen new mined blocks and store
// it to latestBlock hash key
func (s *Service) ListenNewBlocks() {
	for {
		select {
		case <-s.stopChan:
			return
		case err := <-s.subscription.Err():
			logger.Log().Error(err)
			return
		case header := <-s.newBlocks:
			logger.Log().Infof("new block %d just mined", header.Number.Uint64())
			block, err := s.ethClient.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				logger.Log().Error(err)
				continue
			}

			s.cache.StoreLastBlock(models.BlockFromGeth(block))
			logger.Log().Infof("Block %d stored to cache", block.Number())
		}
	}
}

// StopListenNewBlocks stop listen new mined blocks
func (s *Service) StopListenNewBlocks() {
	close(s.stopChan)
	s.subscription.Unsubscribe()
	close(s.newBlocks)
}

// retrieveBlockFromHash try to retrieve block from cache
func (s *Service) retrieveBlockFromHash(number *big.Int) *models.Block {
	switch number {
	case nil:
		return s.cache.SelectLastBlock()
	default:
		return s.cache.SelectBlock(number.Uint64())
	}
}

// fetchTxSender parse Tx sender
func (s *Service) fetchTxSender(tx *models.Tx) error {
	msg, err := tx.GethTx().AsMessage(types.NewEIP155Signer(s.chainID))
	if err != nil {
		return fmt.Errorf("format tx as a core message: %w", err)
	}

	tx.From = msg.From()
	return nil
}

// bigIntFromPointerInt create big.Int from given int pointer
func bigIntFromPointerInt(number *int) (blockNumBigInt *big.Int) {
	if number != nil {
		blockNumBigInt = big.NewInt(int64(*number))
	}
	return
}

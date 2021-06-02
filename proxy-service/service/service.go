package service

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"eth-proxy/models"
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
}

// Service implements Application service layer interface
type Service struct {
	ethClient *ethclient.Client
	chainId   *big.Int
}

// NewService create new Service instance
// with given ethereum client
func NewService(ethClient *ethclient.Client) (IService, error) {
	chainID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("fetch Ethereum chain id: %w", err)
	}

	return &Service{
		ethClient: ethClient,
		chainId:   chainID,
	}, nil
}

// Block fetch Ethereum Block by given block number
// you need to pass nil value as a number argument
// if you want to fetch latest block
func (s *Service) Block(number *int) (*models.Block, error) {
	block, err := s.ethClient.BlockByNumber(context.Background(), bigIntFromPointerInt(number))
	if err != nil {
		return nil, fmt.Errorf("fetch block from Ethereum node: %w", err)
	}

	return models.BlockFromGeth(block), nil
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
	txGeth, _, err := s.ethClient.TransactionByHash(context.Background(), hash)
	if err != nil {
		return nil, fmt.Errorf("fetch transaction by hash from Ethereum node: %w", err)
	}

	tx := models.TxFromGeth(txGeth)

	if err := s.fetchTxSender(tx); err != nil {
		return nil, fmt.Errorf("fetch tx sender: %w", err)
	}

	return tx, nil
}

// fetchTxSender parse Tx sender
func (s *Service) fetchTxSender(tx *models.Tx) error {
	msg, err := tx.GethTx().AsMessage(types.NewEIP155Signer(s.chainId))
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

package lib

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/keepalive"

	"github.com/ethereum/go-ethereum/common"

	"eth-proxy/models"
	"eth-proxy/proto"
)

// timeOut is hardcoded GRPC requests timeout value
const timeOut = 300

// IEthProxyAPI is interface for eth-proxy GRPC API
type IEthProxyAPI interface {
	// Block fetch Ethereum Block by given block number
	// you need to pass nil value as a number argument
	// if you want to fetch latest block
	Block(number *int) (*models.Block, error)

	// TxByIndex fetch Ethereum transaction by related
	// block number and transaction index in the block
	TxByIndex(blockNumber *int, txIndex int) (*models.Tx, error)

	// TxByHash fetch Ethereum transaction by given hash
	TxByHash(hash common.Hash) (*models.Tx, error)

	// Close GRPC Api connection
	Close() error
}

// Api is eth-proxy GRPC Api
// structure with client Connection
type Api struct {
	timeout time.Duration
	*grpc.ClientConn
	proto.EthProxyServiceClient
}

// New create new EthProxy Api instance
func New(addr string) (IEthProxyAPI, error) {
	api := &Api{timeout: timeOut * time.Second}

	if err := api.initConn(addr); err != nil {
		return nil, fmt.Errorf("create EthProxy Api:  %w", err)
	}

	api.EthProxyServiceClient = proto.NewEthProxyServiceClient(api.ClientConn)
	return api, nil
}

// initConn initialize connection to Grpc servers
func (api *Api) initConn(addr string) (err error) {
	kacp := keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	api.ClientConn, err = grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(kacp),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: 1 * time.Second,
		}),
	)
	return
}

// Block fetch Ethereum Block by given block number
// you need to pass nil value as a number argument
// if you want to fetch latest block
func (api *Api) Block(number *int) (*models.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	block, err := api.EthProxyServiceClient.Block(ctx, parseBlockGetter(number))
	if err != nil {
		return nil, fmt.Errorf("fetch block grpc api request: %w", err)
	}

	return models.BlockFromProto(block), nil
}

// TxByIndex fetch Ethereum transaction by related
// block number and transaction index in the block
func (api *Api) TxByIndex(blockNumber *int, txIndex int) (*models.Tx, error) {
	getter := &proto.TxGetter{
		Getter: &proto.TxGetter_Index{
			Index: &proto.TxGetterIndex{
				Index: int64(txIndex),
				Block: parseBlockGetter(blockNumber),
			}},
	}

	return api.fetchTx(getter)
}

// TxByHash fetch Ethereum transaction by given hash
func (api *Api) TxByHash(hash common.Hash) (*models.Tx, error) {
	getter := &proto.TxGetter{
		Getter: &proto.TxGetter_Hash{
			Hash: hash.Bytes(),
		},
	}

	return api.fetchTx(getter)
}

// fetchTx fetch transaction from Grpc api
func (api *Api) fetchTx(getter *proto.TxGetter) (*models.Tx, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	tx, err := api.EthProxyServiceClient.Tx(ctx, getter)
	if err != nil {
		return nil, fmt.Errorf("fetch tx grpc api request: %w", err)
	}

	return models.TxFromProto(tx), nil
}

// parseBlockGetter parse block getter
// from given number
func parseBlockGetter(number *int) *proto.BlockGetter {
	getter := &proto.BlockGetter{}
	if number != nil {
		getter.Getter = &proto.BlockGetter_Number{Number: int64(*number)}
	} else {
		getter.Getter = &proto.BlockGetter_Latest{Latest: true}
	}

	return getter
}

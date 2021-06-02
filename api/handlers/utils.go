package handlers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-openapi/swag"

	"eth-proxy/api/models"
)

const lastBlock = "latest"

var (
	errInternalServer = errors.New("internal server error")
	errbadRequest     = errors.New("bad request")
	errUnauthorized   = errors.New("authentication failed")
)

// handleError return errors model
func handleError(err error) *models.Error {
	return &models.Error{
		Message: swag.String(err.Error()),
	}
}

// parseBlockNumber parse block number
// or latest block
func parseBlockNumber(number string) (blockNumber *int, err error) {
	if number != lastBlock {
		parsedNum, err := strconv.Atoi(number)
		if err != nil {
			return nil, fmt.Errorf("parse block number: %w", err)
		}
		blockNumber = &parsedNum
	}
	return
}

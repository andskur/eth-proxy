package format

import (
	apiModels "eth-proxy/api/models"
	"eth-proxy/models"
)

// Tx formats Tx model to swagger definition
func Tx(tx *models.Tx) *apiModels.Tx {
	return &apiModels.Tx{
		Hash:     tx.Hash.Hex(),
		From:     tx.From.Hex(),
		To:       tx.To.Hex(),
		Value:    tx.Value,
		Gas:      int64(tx.Gas),
		GasPrice: tx.GasPrice,
	}
}

package format

import (
	"github.com/go-openapi/swag"

	apiModels "eth-proxy/api/models"
	"eth-proxy/models"
)

// Block formats Block model to swagger definition
func Block(block *models.Block) *apiModels.Block {
	return &apiModels.Block{
		Number:    swag.Int64(int64(block.Number)),
		Hash:      swag.String(block.Hash.Hex()),
		Parent:    block.Parent.Hex(),
		Timestamp: swag.Int64(block.Time.Unix()),
		TxCount:   swag.Int64(int64(block.TxCount)),
	}
}

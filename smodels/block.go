package smodels

import "github.com/everstake/cosmoscan-api/dmodels"

type (
	Block struct {
		Height          uint64       `json:"height"`
		Hash            string       `json:"hash"`
		TotalTxs        uint64       `json:"total_txs"`
		ChainID         string       `json:"chain_id"`
		Proposer        string       `json:"proposer"`
		ProposerAddress string       `json:"proposer_address"`
		Txs             []TxItem     `json:"txs"`
		CreatedAt       dmodels.Time `json:"created_at"`
	}
	BlockItem struct {
		Height          uint64       `json:"height"`
		Hash            string       `json:"hash"`
		Proposer        string       `json:"proposer"`
		ProposerAddress string       `json:"proposer_address"`
		CreatedAt       dmodels.Time `json:"created_at"`
	}
)

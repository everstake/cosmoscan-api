package smodels

import (
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/shopspring/decimal"
)

type (
	TxItem struct {
		Hash      string          `json:"hash"`
		Status    bool            `json:"status"`
		Fee       decimal.Decimal `json:"fee"`
		Height    uint64          `json:"height"`
		Messages  uint64          `json:"messages"`
		CreatedAt dmodels.Time    `json:"created_at"`
	}
	Tx struct {
		Hash      string          `json:"hash"`
		Type      string          `json:"type"`
		Status    bool            `json:"status"`
		Fee       decimal.Decimal `json:"fee"`
		Height    uint64          `json:"height"`
		GasUsed   uint64          `json:"gas_used"`
		GasWanted uint64          `json:"gas_wanted"`
		Memo      string          `json:"memo"`
		CreatedAt dmodels.Time    `json:"created_at"`
		Messages  []interface{}   `json:"messages"`
	}
	SendMsg struct {
		From   string          `json:"from"`
		To     string          `json:"to"`
		Amount decimal.Decimal `json:"amount"`
	}
)

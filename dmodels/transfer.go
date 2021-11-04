package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const TransfersTable = "transfers"
const CosmosCurrency = "btsg"

type Transfer struct {
	ID        string          `db:"trf_id"`
	TxHash    string          `db:"trf_tx_hash"`
	From      string          `db:"trf_from"`
	To        string          `db:"trf_to"`
	Amount    decimal.Decimal `db:"trf_amount"`
	Currency  string          `db:"trf_currency"`
	CreatedAt time.Time       `db:"trf_created_at"`
}

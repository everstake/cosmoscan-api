package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dmodels"
)

func (db DB) CreateTransfers(transfers []dmodels.Transfer) error {
	if len(transfers) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.TransfersTable).Columns("trf_id", "trf_tx_hash", "trf_from", "trf_to", "trf_amount", "trf_created_at")
	for _, transfer := range transfers {
		if transfer.ID == "" {
			return fmt.Errorf("field ID can not be empty")
		}
		if transfer.TxHash == "" {
			return fmt.Errorf("field TxHash can not be empty")
		}
		if transfer.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt can not be zero")
		}
		q = q.Values(transfer.ID, transfer.TxHash, transfer.From, transfer.To, transfer.Amount, transfer.CreatedAt)
	}
	return db.Insert(q)
}

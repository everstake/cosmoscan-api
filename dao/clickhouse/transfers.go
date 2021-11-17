package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/config"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/smodels"
	"github.com/shopspring/decimal"
)

func (db DB) CreateTransfers(transfers []dmodels.Transfer) error {
	if len(transfers) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.TransfersTable).Columns("trf_id", "trf_tx_hash", "trf_from", "trf_to", "trf_amount", "trf_created_at", "trf_currency")
	for _, transfer := range transfers {
		if transfer.ID == "" {
			return fmt.Errorf("field ProposalID can not be empty")
		}
		if transfer.TxHash == "" {
			return fmt.Errorf("field TxHash can not be empty")
		}
		if transfer.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt can not be zero")
		}
		q = q.Values(transfer.ID, transfer.TxHash, transfer.From, transfer.To, transfer.Amount, transfer.CreatedAt, transfer.Currency)
	}
	return db.Insert(q)
}

func (db DB) GetAggTransfersVolume(filter filters.Agg) (items []smodels.AggItem, err error) {
	q := squirrel.Select(
		"sum(trf_amount) AS value",
		fmt.Sprintf("toDateTime(%s(trf_created_at)) AS time", filter.AggFunc()),
	).From(dmodels.TransfersTable).
		Where("notEmpty(trf_from)").
		Where(squirrel.Eq{"trf_currency": config.Currency}).
		GroupBy("time").
		OrderBy("time")
	if !filter.From.IsZero() {
		q = q.Where(squirrel.GtOrEq{"trf_created_at": filter.From.Time})
	}
	if !filter.To.IsZero() {
		q = q.Where(squirrel.LtOrEq{"trf_created_at": filter.To.Time})
	}
	err = db.Find(&items, q)
	return items, err
}

func (db DB) GetTransferVolume(filter filters.TimeRange) (total decimal.Decimal, err error) {
	q := squirrel.Select("sum(trf_amount) as total").
		From(dmodels.TransfersTable).
		Where("notEmpty(trf_from)").
		Where(squirrel.Eq{"trf_currency": config.Currency})
	q = filter.Query("trf_created_at", q)
	err = db.FindFirst(&total, q)
	return total, err
}

package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/smodels"
)

func (db DB) CreateTransactions(transactions []dmodels.Transaction) error {
	if len(transactions) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.TransactionsTable).Columns(
		"trn_hash",
		"trn_status",
		"trn_height",
		"trn_messages",
		"trn_fee",
		"trn_gas_used",
		"trn_gas_wanted",
		"trn_created_at",
	)
	for _, tx := range transactions {
		if tx.Hash == "" {
			return fmt.Errorf("field Hash can not be empty")
		}
		if tx.Height == 0 {
			return fmt.Errorf("field Height can not be 0")
		}
		if tx.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt can not be zero")
		}
		q = q.Values(
			tx.Hash,
			tx.Status,
			tx.Height,
			tx.Messages,
			tx.Fee,
			tx.GasUsed,
			tx.GasWanted,
			tx.CreatedAt,
		)
	}
	return db.Insert(q)
}

func (db DB) GetAggTransactionsFee(filter filters.Agg) (items []smodels.AggItem, err error) {
	q := squirrel.Select(
		"sum(trn_fee) AS value",
		fmt.Sprintf("toDateTime(%s(trn_created_at)) AS time", filter.AggFunc()),
	).From(dmodels.TransactionsTable).
		GroupBy("time").
		OrderBy("time")
	if !filter.From.IsZero() {
		q = q.Where(squirrel.GtOrEq{"trn_created_at": filter.From.Time})
	}
	if !filter.To.IsZero() {
		q = q.Where(squirrel.LtOrEq{"trn_created_at": filter.To.Time})
	}
	err = db.Find(&items, q)
	if err != nil {
		return nil, err
	}
	return items, nil
}

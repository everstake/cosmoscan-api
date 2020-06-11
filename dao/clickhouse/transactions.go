package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/smodels"
	"github.com/shopspring/decimal"
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
	q := filter.BuildQuery("sum(trn_fee)", "trn_created_at", dmodels.TransactionsTable)
	err = db.Find(&items, q)
	return items, err
}

func (db DB) GetAggOperationsCount(filter filters.Agg) (items []smodels.AggItem, err error) {
	q := filter.BuildQuery("toDecimal64(sum(trn_messages), 0)", "trn_created_at", dmodels.TransactionsTable)
	err = db.Find(&items, q)
	return items, err
}

func (db DB) GetTransactionsFeeVolume(filter filters.TimeRange) (total decimal.Decimal, err error) {
	q := squirrel.Select("sum(trn_fee) as total").From(dmodels.TransactionsTable)
	q = filter.Query("trn_created_at", q)
	err = db.FindFirst(&total, q)
	return total, err
}

func (db DB) GetTransactionsHighestFee(filter filters.TimeRange) (total decimal.Decimal, err error) {
	q := squirrel.Select("max(trn_fee) as total").From(dmodels.TransactionsTable)
	q = filter.Query("trn_created_at", q)
	err = db.FindFirst(&total, q)
	return total, err
}

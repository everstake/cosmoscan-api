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
		"trn_height",
		"trn_messages",
		"trn_fee",
		"trn_gas_used",
		"trn_gas_wanted",
		"trn_signer",
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
			tx.Height,
			tx.Messages,
			tx.Fee,
			tx.GasUsed,
			tx.GasWanted,
			tx.Signer,
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

func (db DB) GetAvgOperationsPerBlock(filter filters.Agg) (items []smodels.AggItem, err error) {
	// approximate number of blocks by `period`
	blocks := 12000
	switch filter.By {
	case filters.AggByHour:
		blocks = 500
	case filters.AggByWeek:
		blocks = 84000
	case filters.AggByMonth:
		blocks = 360000
	}
	aggValue := fmt.Sprintf("toDecimal64(sum(trn_messages) / %d, 4)", blocks)
	q := filter.BuildQuery(aggValue, "trn_created_at", dmodels.TransactionsTable)
	err = db.Find(&items, q)
	return items, err
}

func (db DB) GetTransactions(filter filters.Transactions) (items []dmodels.Transaction, err error) {
	q := squirrel.Select("*").From(dmodels.TransactionsTable).OrderBy("trn_created_at desc")
	if filter.Height != 0 {
		q = q.Where(squirrel.Eq{"trn_height": filter.Height})
	}
	if filter.Address != "" {
		q = q.Where(squirrel.Eq{"trn_signer": filter.Address})
	}
	if filter.Limit != 0 {
		q = q.Limit(filter.Limit)
	}
	if filter.Offset != 0 {
		q = q.Offset(filter.Offset)
	}
	err = db.Find(&items, q)
	return items, err
}

func (db DB) GetTransactionsCount(filter filters.Transactions) (total uint64, err error) {
	q := squirrel.Select("count(*)").From(dmodels.TransactionsTable)
	if filter.Height != 0 {
		q = q.Where(squirrel.Eq{"trn_height": filter.Height})
	}
	err = db.FindFirst(&total, q)
	return total, err
}

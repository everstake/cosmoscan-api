package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
)

func (db DB) CreateBalanceUpdates(updates []dmodels.BalanceUpdate) error {
	if len(updates) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.BalanceUpdatesTable).Columns("bau_id", "bau_address", "bau_stake", "bau_balance", "bau_unbonding", "bau_created_at")
	for _, update := range updates {
		if update.ID == "" {
			return fmt.Errorf("field ProposalID can not be empty")
		}
		if update.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt can not be 0")
		}
		q = q.Values(update.ID, update.Address, update.Stake, update.Balance, update.Unbonding, update.CreatedAt)
	}
	return db.Insert(q)
}

func (db DB) GetBalanceUpdate(filter filters.BalanceUpdates) (updates []dmodels.BalanceUpdate, err error) {
	q := squirrel.Select("*").From(dmodels.BalanceUpdatesTable).OrderBy("bau_created_at desc")
	if filter.Limit != 0 {
		q = q.Limit(filter.Limit)
	}
	if filter.Offset != 0 {
		q = q.Offset(filter.Offset)
	}
	err = db.Find(&updates, q)
	return updates, err
}

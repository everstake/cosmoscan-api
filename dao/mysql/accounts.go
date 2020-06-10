package mysql

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dmodels"
)

func (m DB) CreateAccounts(accounts []dmodels.Account) error {
	if len(accounts) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.AccountsTable).Columns(
		"acc_address",
		"acc_balance",
		"acc_created_at",
	)
	for _, account := range accounts {
		if account.Address == "" {
			return fmt.Errorf("field Address is empty")
		}
		if account.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt is empty")
		}
		q = q.Values(
			account.Address,
			account.Balance,
			account.CreatedAt,
		)
	}
	_, err := m.insert(q)
	return err
}

func (m DB) UpdateAccount(account dmodels.Account) error {
	q := squirrel.Update(dmodels.AccountsTable).
		Where(squirrel.Eq{"acc_address": account.Address}).
		SetMap(map[string]interface{}{
			"acc_balance": account.Balance,
		})
	return m.update(q)
}

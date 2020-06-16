package mysql

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
)

func (m DB) CreateAccounts(accounts []dmodels.Account) error {
	if len(accounts) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.AccountsTable).Columns(
		"acc_address",
		"acc_balance",
		"acc_stake",
		"acc_unbonding",
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
			account.Stake,
			account.Unbonding,
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
			"acc_balance":   account.Balance,
			"acc_stake":     account.Stake,
			"acc_unbonding": account.Unbonding,
		})
	return m.update(q)
}

func (m DB) GetAccounts(filter filters.Accounts) (accounts []dmodels.Account, err error) {
	q := squirrel.Select("*").From(dmodels.AccountsTable)
	if !filter.GtTotalAmount.IsZero() {
		q = q.Where(squirrel.Gt{"acc_balance + acc_stake": filter.GtTotalAmount})
	}
	if !filter.LtTotalAmount.IsZero() {
		q = q.Where(squirrel.Lt{"acc_balance + acc_stake": filter.LtTotalAmount})
	}
	err = m.find(&accounts, q)
	return accounts, err
}

func (m DB) GetAccountsTotal(filter filters.Accounts) (total uint64, err error) {
	q := squirrel.Select("count(*) as total").From(dmodels.AccountsTable)
	if filter.GtTotalAmount.IsZero() {
		q = q.Where(squirrel.Gt{"acc_balance + acc_stake + acc_unbonding": filter.GtTotalAmount})
	}
	if filter.LtTotalAmount.IsZero() {
		q = q.Where(squirrel.Lt{"acc_balance + acc_stake + acc_unbonding": filter.LtTotalAmount})
	}
	err = m.first(&total, q)
	return total, err
}

func (m DB) GetAccount(address string) (account dmodels.Account, err error) {
	q := squirrel.Select("*").From(dmodels.AccountsTable).Where(squirrel.Eq{"acc_address": address})
	err = m.first(&account, q)
	return account, err
}

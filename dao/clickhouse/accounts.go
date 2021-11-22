package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
)

func (db DB) GetActiveAccounts(filter filters.ActiveAccounts) (addresses []string, err error) {
	var items = []struct {
		field     string
		table     string
		dateField string
	}{
		{field: "dlg_delegator", table: dmodels.DelegationsTable, dateField: "dlg_created_at"},
		{field: "trf_from", table: dmodels.TransfersTable, dateField: "trf_created_at"},
		{field: "trf_to", table: dmodels.TransfersTable, dateField: "trf_created_at"},
		{field: "der_delegator", table: dmodels.DelegatorRewardsTable, dateField: "der_created_at"},
	}

	var qs []squirrel.SelectBuilder
	for _, item := range items {
		q := squirrel.Select(fmt.Sprintf("DISTINCT %s as address", item.field)).
			From(item.table)
		if !filter.From.IsZero() {
			q = q.Where(squirrel.GtOrEq{item.dateField: filter.From})
		}
		if !filter.To.IsZero() {
			q = q.Where(squirrel.LtOrEq{item.dateField: filter.To})
		}
		qs = append(qs, q)
	}

	q := qs[0]

	for i := 1; i < len(qs); i++ {
		sql, args, _ := qs[i].ToSql()
		q = qs[i].Suffix("UNION ALL "+sql, args...)
	}

	query := squirrel.Select("DISTINCT t.address").FromSelect(q, "t")

	err = db.Find(&addresses, query)
	return addresses, err
}

func (db DB) CreateAccountTxs(accountTxs []dmodels.AccountTx) error {
	if len(accountTxs) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.AccountTxsTable).Columns("atx_account", "atx_tx_hash")
	for _, acc := range accountTxs {
		if acc.Account == "" {
			return fmt.Errorf("field Account can not beempty")
		}
		if acc.TxHash == "" {
			return fmt.Errorf("hash can not be empty")
		}
		q = q.Values(acc.Account, acc.TxHash)
	}
	return db.Insert(q)
}

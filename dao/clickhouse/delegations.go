package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/smodels"
	"github.com/shopspring/decimal"
)

func (db DB) CreateDelegations(delegations []dmodels.Delegation) error {
	if len(delegations) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.DelegationsTable).Columns("dlg_id", "dlg_tx_hash", "dlg_delegator", "dlg_validator", "dlg_amount", "dlg_created_at")
	for _, delegation := range delegations {
		if delegation.ID == "" {
			return fmt.Errorf("field ProposalID can not be empty")
		}
		if delegation.TxHash == "" {
			return fmt.Errorf("field TxHash can not be empty")
		}
		if delegation.Delegator == "" {
			return fmt.Errorf("field Delegator can not be empty")
		}
		if delegation.Validator == "" {
			return fmt.Errorf("field Validator can not be empty")
		}
		if delegation.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt can not be zero")
		}
		q = q.Values(delegation.ID, delegation.TxHash, delegation.Delegator, delegation.Validator, delegation.Amount, delegation.CreatedAt)
	}
	return db.Insert(q)
}

func (db DB) GetAggDelegationsVolume(filter filters.DelegationsAgg) (items []smodels.AggItem, err error) {
	q := filter.BuildQuery("sum(dlg_amount)", "dlg_created_at", dmodels.DelegationsTable)
	if len(filter.Validators) != 0 {
		q = q.Where(squirrel.Eq{"dlg_validator": filter.Validators})
	}
	q = q.Where(squirrel.Gt{"dlg_amount": 0})
	err = db.Find(&items, q)
	return items, err
}

func (db DB) GetAggUndelegationsVolume(filter filters.Agg) (items []smodels.AggItem, err error) {
	q := filter.BuildQuery("sum(abs(dlg_amount))", "dlg_created_at", dmodels.DelegationsTable)
	q = q.Where(squirrel.Lt{"dlg_amount": 0})
	err = db.Find(&items, q)
	return items, err
}

func (db DB) GetDelegatorsTotal(filter filters.Delegators) (total uint64, err error) {
	q1 := squirrel.Select("dlg_delegator as delegator", "sum(dlg_amount) as amount").
		From(dmodels.DelegationsTable).GroupBy("dlg_delegator").
		Having(squirrel.Gt{"amount": 0})
	if len(filter.Validators) != 0 {
		q1 = q1.Where(squirrel.Eq{"dlg_validator": filter.Validators})
	}
	q1 = filter.Query("dlg_created_at", q1)
	q := squirrel.Select("count() as total").FromSelect(q1, "t")
	err = db.FindFirst(&total, q)
	return total, err
}

func (db DB) GetMultiDelegatorsTotal(filter filters.TimeRange) (total uint64, err error) {
	q1 := squirrel.Select("dlg_delegator as delegator", "sum(dlg_amount) as amount", "count(DISTINCT dlg_validator) as validators_count").
		From(dmodels.DelegationsTable).GroupBy("dlg_delegator").
		Having(squirrel.Gt{"amount": 0}).Having(squirrel.Gt{"validators_count": 1})
	q1 = filter.Query("dlg_created_at", q1)
	q := squirrel.Select("count() as total").FromSelect(q1, "t")
	err = db.FindFirst(&total, q)
	return total, err
}

func (db DB) GetUndelegationsVolume(filter filters.TimeRange) (total decimal.Decimal, err error) {
	q := squirrel.Select("sum(abs(dlg_amount)) as total").
		From(dmodels.DelegationsTable).
		Where(squirrel.Lt{"dlg_amount": 0})
	q = filter.Query("dlg_created_at", q)
	err = db.FindFirst(&total, q)
	return total, err
}

func (db DB) GetVotingPower(filter filters.VotingPower) (volume decimal.Decimal, err error) {
	q := squirrel.Select("sum(dlg_amount) as volume").From(dmodels.DelegationsTable)
	q = filter.Query("dlg_created_at", q)
	if len(filter.Delegators) != 0 {
		q = q.Where(squirrel.Eq{"dlg_delegator": filter.Delegators})
	}
	if len(filter.Validators) != 0 {
		q = q.Where(squirrel.Eq{"dlg_validator": filter.Validators})
	}
	err = db.FindFirst(&volume, q)
	return volume, err
}

func (db DB) GetValidatorsDelegatorsTotal() (values []dmodels.ValidatorValue, err error) {
	q1 := squirrel.Select("sum(dlg_amount) as volume", "dlg_delegator", "dlg_validator").
		From(dmodels.DelegationsTable).
		GroupBy("dlg_delegator", "dlg_validator").
		Having(squirrel.Gt{"volume": 0})
	q := squirrel.Select("count(dlg_validator) as value, t.dlg_validator as validator").
		FromSelect(q1, "t").GroupBy("dlg_validator").OrderBy("value desc")
	err = db.Find(&values, q)
	return values, err
}

func (db DB) GetValidatorDelegators(filter filters.ValidatorDelegators) (items []dmodels.ValidatorDelegator, err error) {
	query := `SELECT  * FROM
	(SELECT dlg_delegator as delegator, sum(dlg_amount) as amount, min(dlg_created_at) as since
	FROM delegations
	WHERE dlg_validator = ?
	GROUP BY dlg_delegator
	HAVING amount > 0 ORDER BY amount DESC) as t1
	ANY LEFT JOIN (
		SELECT sum(dlg_amount) as delta, dlg_delegator as delegator
		FROM delegations
		WHERE dlg_validator = ? and dlg_created_at > yesterday()
		GROUP BY dlg_delegator
	) as t2 USING (delegator)`
	if filter.Limit != 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, filter.Limit)
	}
	if filter.Offset != 0 {
		query = fmt.Sprintf("%s OFFSET %d", query, filter.Offset)
	}
	q, args, err := squirrel.Expr(query, filter.Validator, filter.Validator).ToSql()
	if err != nil {
		return nil, err
	}
	err = db.conn.Select(&items, q, args...)
	return items, err
}

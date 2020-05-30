package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dmodels"
)

func (db DB) CreateDelegatorRewards(rewards []dmodels.DelegatorReward) error {
	if len(rewards) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.DelegatorRewardsTable).Columns("der_id", "der_tx_hash", "der_delegator", "der_validator", "der_amount", "der_created_at")
	for _, reward := range rewards {
		if reward.ID == "" {
			return fmt.Errorf("field ID can not be empty")
		}
		if reward.TxHash == "" {
			return fmt.Errorf("field TxHash can not be empty")
		}
		if reward.Delegator == "" {
			return fmt.Errorf("field Delegator can not be empty")
		}
		if reward.Validator == "" {
			return fmt.Errorf("field Validator can not be empty")
		}
		if reward.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt can not be zero")
		}
		q = q.Values(reward.ID, reward.TxHash, reward.Delegator, reward.Validator, reward.Amount, reward.CreatedAt)
	}
	return db.Insert(q)
}

func (db DB) CreateValidatorRewards(rewards []dmodels.ValidatorReward) error {
	if len(rewards) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.ValidatorRewardsTable).Columns("var_id", "var_address", "var_amount", "var_created_at")
	for _, reward := range rewards {
		if reward.ID == "" {
			return fmt.Errorf("field ID can not be empty")
		}
		if reward.Address == "" {
			return fmt.Errorf("field TxHash can not be empty")
		}
		if reward.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt can not be zero")
		}
		q = q.Values(reward.ID, reward.Address, reward.Amount, reward.CreatedAt)
	}
	return db.Insert(q)
}

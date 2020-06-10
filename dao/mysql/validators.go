package mysql

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dmodels"
)

func (m DB) CreateValidators(validators []dmodels.Validator) error {
	if len(validators) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.ValidatorsTable).Columns(
		"val_cons_address",
		"val_address",
		"val_operator_address",
		"val_cons_pub_key",
		"val_name",
		"val_description",
		"val_commission",
		"val_min_commission",
		"val_max_commission",
		"val_self_delegations",
		"val_delegations",
		"val_voting_power",
		"val_website",
		"val_jailed",
		"val_created_at",
	)
	for _, validator := range validators {
		if validator.ConsAddress == "" {
			return fmt.Errorf("ConsAddress is empty")
		}
		q = q.Values(
			validator.ConsAddress,
			validator.Address,
			validator.OperatorAddress,
			validator.ConsPubKey,
			validator.Name,
			validator.Description,
			validator.Commission,
			validator.MinCommission,
			validator.MaxCommission,
			validator.SelfDelegations,
			validator.Delegations,
			validator.VotingPower,
			validator.Website,
			validator.Jailed,
			validator.CreatedAt,
		)
	}
	_, err := m.insert(q)
	return err
}

func (m DB) UpdateValidators(validator dmodels.Validator) error {
	q := squirrel.Update(dmodels.ValidatorsTable).
		Where(squirrel.Eq{"val_cons_address": validator.ConsAddress}).
		SetMap(map[string]interface{}{
			"val_address":          validator.Address,
			"val_operator_address": validator.OperatorAddress,
			"val_cons_pub_key":     validator.ConsPubKey,
			"val_name":             validator.Name,
			"val_description":      validator.Description,
			"val_commission":       validator.Commission,
			"val_min_commission":   validator.MinCommission,
			"val_max_commission":   validator.MaxCommission,
			"val_self_delegations": validator.SelfDelegations,
			"val_delegations":      validator.Delegations,
			"val_voting_power":     validator.VotingPower,
			"val_website":          validator.Website,
			"val_jailed":           validator.Jailed,
		})
	return m.update(q)
}

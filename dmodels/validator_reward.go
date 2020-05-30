package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const ValidatorRewardsTable = "validator_rewards"

type ValidatorReward struct {
	ID        string          `db:"var_id"`
	Address   string          `db:"var_address"`
	Amount    decimal.Decimal `db:"var_amount"`
	CreatedAt time.Time       `db:"var_created_at"`
}

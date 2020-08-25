package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const StatsTable = "stats"

const (
	StatsTotalStakingBalance   = "total_staking_balance"
	StatsNumberDelegators      = "number_delegators"
	StatsTotalDelegators       = "total_delegators"
	StatsNumberMultiDelegators = "number_multi_delegators"
	StatsTransfersVolume       = "transfer_volume"
	StatsFeeVolume             = "fee_volume"
	StatsHighestFee            = "highest_fee"
	StatsUndelegationVolume    = "undelegation_volume"
	StatsBlockDelay            = "block_delay"
	StatsNetworkSize           = "network_size"
	StatsTotalAccounts         = "total_accounts"
	StatsTotalWhaleAccounts    = "total_whale_accounts"
	StatsTotalSmallAccounts    = "total_small_accounts"
	StatsTotalJailers          = "total_jailers"
	StatsValidatorsWith33Power = "validators_with_33_power"
)

type Stat struct {
	ID        string          `db:"stt_id"`
	Title     string          `db:"stt_title"`
	Value     decimal.Decimal `db:"stt_value"`
	CreatedAt time.Time       `db:"stt_created_at"`
}

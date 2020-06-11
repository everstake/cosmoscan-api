package dmodels

import "time"

const (
	RangeStateTotalStakingBalance   = "total_staking_balance"
	RangeStateNumberDelegators      = "number_delegators"
	RangeStateNumberMultiDelegators = "number_multi_delegators"
	RangeStateTransfersVolume       = "transfer_volume"
	RangeStateFeeVolume             = "fee_volume"
	RangeStateHighestFee            = "highest_fee"
	RangeStateUndelegationVolume    = "undelegation_volume"
	RangeStateBlockDelay            = "block_delay"
)

const RangeStatesTable = "range_states"

type RangeState struct {
	Title     string    `db:"rst_title"`
	Value1d   string    `db:"rst_value_1d"`
	Value7d   string    `db:"rst_value_7d"`
	Value30d  string    `db:"rst_value_30d"`
	Value90d  string    `db:"rst_value_90d"`
	UpdatedAt time.Time `db:"rst_updated_at"`
}

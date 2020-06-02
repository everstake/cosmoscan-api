package dmodels

import (
	"github.com/shopspring/decimal"
)

const HistoricalStates = "historical_states"

type HistoricalState struct {
	Price             decimal.Decimal `db:"his_price"`
	MarketCap         decimal.Decimal `db:"his_market_cap"`
	CirculationSupply decimal.Decimal `db:"his_circulation_supply"`
	TradingVolume     decimal.Decimal `db:"his_trading_volume"`
	StakedRatio       decimal.Decimal `db:"his_staked_ratio"`
	InflationRate     decimal.Decimal `db:"his_inflation_rate"`
	TransactionsCount uint64          `db:"his_transactions_count"`
	CommunityPool     decimal.Decimal `db:"his_community_pool"`
	Top20Weight       decimal.Decimal `db:"his_top_20_weight"`
	CreatedAt         Time            `db:"his_created_at"`
}

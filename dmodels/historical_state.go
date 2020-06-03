package dmodels

import (
	"github.com/shopspring/decimal"
)

const HistoricalStates = "historical_states"

type HistoricalState struct {
	Price             decimal.Decimal `db:"his_price" json:"price"`
	MarketCap         decimal.Decimal `db:"his_market_cap" json:"market_cap"`
	CirculatingSupply decimal.Decimal `db:"his_circulating_supply" json:"circulating_supply"`
	TradingVolume     decimal.Decimal `db:"his_trading_volume" json:"trading_volume"`
	StakedRatio       decimal.Decimal `db:"his_staked_ratio" json:"staked_ratio"`
	InflationRate     decimal.Decimal `db:"his_inflation_rate" json:"inflation_rate"`
	TransactionsCount uint64          `db:"his_transactions_count" json:"transactions_count"`
	CommunityPool     decimal.Decimal `db:"his_community_pool" json:"community_pool"`
	Top20Weight       decimal.Decimal `db:"his_top_20_weight" json:"top20_weight"`
	CreatedAt         Time            `db:"his_created_at" json:"created_at"`
}

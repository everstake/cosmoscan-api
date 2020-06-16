package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const StatsTable = "stats"

type Stat struct {
	ID        string          `db:"stt_id"`
	Title     string          `db:"stt_title"`
	Value     decimal.Decimal `db:"stt_value"`
	CreatedAt time.Time       `db:"stt_created_at"`
}

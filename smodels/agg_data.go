package smodels

import (
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/shopspring/decimal"
)

type AggItem struct {
	Time  dmodels.Time    `db:"time" json:"time"`
	Value decimal.Decimal `db:"value" json:"value"`
}

package filters

import (
	"github.com/shopspring/decimal"
	"time"
)

type Accounts struct {
	LtTotalAmount decimal.Decimal
	GtTotalAmount decimal.Decimal
}

type ActiveAccounts struct {
	From time.Time
	To   time.Time
}

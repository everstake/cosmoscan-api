package filters

import "github.com/everstake/cosmoscan-api/dmodels"

type Stats struct {
	Titles []string     `schema:"-"`
	To     dmodels.Time `schema:"to"`
	From   dmodels.Time `schema:"-"`
}

package services

import (
	"github.com/everstake/cosmoscan-api/config"
	"github.com/everstake/cosmoscan-api/dao"
	"github.com/everstake/cosmoscan-api/services/cmc"
)

type (
	Services interface {
	}
	CMC interface {
		GetCurrencies() (currencies []cmc.Currency, err error)
	}

	ServiceFacade struct {
		dao dao.DAO
		cfg config.Config
		cmc CMC
	}
)

func NewServices(d dao.DAO, cfg config.Config) (svc Services, err error) {
	return &ServiceFacade{
		dao: d,
		cfg: cfg,
		cmc: cmc.NewCMC(cfg),
	}, nil
}

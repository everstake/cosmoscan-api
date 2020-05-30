package services

import (
	"github.com/everstake/cosmoscan-api/config"
	"github.com/everstake/cosmoscan-api/dao"
)

type (
	Services interface {
	}

	ServiceFacade struct {
		dao dao.DAO
		cfg config.Config
	}
)

func NewServices(d dao.DAO, cfg config.Config) (svc Services, err error) {
	return &ServiceFacade{
		dao: d,
		cfg: cfg,
	}, nil
}

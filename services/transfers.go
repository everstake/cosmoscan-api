package services

import (
	"fmt"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/smodels"
)

func (s *ServiceFacade) GetAggTransfersVolume(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAggTransfersVolume(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAggTransfersVolume: %s", err.Error())
	}
	return items, nil
}


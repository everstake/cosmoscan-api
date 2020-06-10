package services

import (
	"fmt"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/smodels"
)

func (s *ServiceFacade) GetAggDelegationsVolume(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAggDelegationsVolume(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAggDelegationsVolume: %s", err.Error())
	}
	return items, nil
}

func (s *ServiceFacade) GetAggUndelegationsVolume(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAggUndelegationsVolume(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAggUndelegationsVolume: %s", err.Error())
	}
	return items, nil
}

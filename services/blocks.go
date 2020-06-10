package services

import (
	"fmt"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/smodels"
)

func (s *ServiceFacade) GetAggBlocksCount(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAggBlocksCount(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAggBlocksCount: %s", err.Error())
	}
	return items, nil
}

func (s *ServiceFacade) GetAggBlocksDelay(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAggBlocksDelay(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAggBlocksDelay: %s", err.Error())
	}
	return items, nil
}

func (s *ServiceFacade) GetAggUniqBlockValidators(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAggUniqBlockValidators(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAggUniqBlockValidators: %s", err.Error())
	}
	return items, nil
}

package services

import (
	"fmt"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/smodels"
)

func (s *ServiceFacade) GetAggTransactionsFee(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAggTransactionsFee(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAggTransactionsFee: %s", err.Error())
	}
	return items, nil
}

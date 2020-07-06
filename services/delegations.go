package services

import (
	"fmt"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/smodels"
	"github.com/shopspring/decimal"
	"time"
)

const getValidatorDelegatorsAggCacheKey = "GetValidatorDelegatorsAgg"

func (s *ServiceFacade) GetAggDelegationsVolume(filter filters.DelegationsAgg) (items []smodels.AggItem, err error) {
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

func (s *ServiceFacade) GetValidatorDelegationsAgg(validatorAddress string) (items []smodels.AggItem, err error) {
	validator, err := s.GetValidator(validatorAddress)
	if err != nil {
		return nil, fmt.Errorf("GetValidator: %s", err.Error())
	}
	items, err = s.dao.GetAggDelegationsVolume(filters.DelegationsAgg{
		Agg: filters.Agg{
			By:   filters.AggByDay,
			From: dmodels.NewTime(time.Now().Add(-time.Hour * 24 * 30)),
			To:   dmodels.NewTime(time.Now()),
		},
		Validators: []string{validatorAddress},
	})
	if err != nil {
		return nil, fmt.Errorf("dao.GetAggDelegationsVolume: %s", err.Error())
	}
	powerValue := validator.Power
	for i := len(items) - 1; i >= 0; i-- {
		v := items[i].Value
		items[i].Value = powerValue
		powerValue = items[i].Value.Sub(v)
	}
	return items, nil
}

func (s *ServiceFacade) GetValidatorDelegatorsAgg(validatorAddress string) (items []smodels.AggItem, err error) {
	data, found := s.dao.CacheGet(getValidatorDelegatorsAggCacheKey)
	if found {
		return data.([]smodels.AggItem), nil
	}
	for i := 29; i >= 0; i-- {
		y, m, d := time.Now().Add(-time.Hour * 24 * time.Duration(i)).Date()
		date := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
		total, err := s.dao.GetDelegatorsTotal(filters.Delegators{
			TimeRange: filters.TimeRange{
				To: dmodels.NewTime(date),
			},
			Validators: []string{validatorAddress},
		})
		if err != nil {
			return nil, fmt.Errorf("dao.GetDelegatorsTotal: %s", err.Error())
		}
		items = append(items, smodels.AggItem{
			Time:  dmodels.NewTime(date),
			Value: decimal.NewFromInt(int64(total)),
		})
	}
	s.dao.CacheSet(getValidatorDelegatorsAggCacheKey, items, time.Hour)
	return items, nil
}

func (s *ServiceFacade) GetValidatorDelegators(filter filters.ValidatorDelegators) (resp smodels.PaginatableResponse, err error) {
	items, err := s.dao.GetValidatorDelegators(filter)
	if err != nil {
		return resp, fmt.Errorf("dao.GetValidatorDelegators: %s", err.Error())
	}
	total, err := s.dao.GetValidatorDelegatorsTotal(filter)
	if err != nil {
		return resp, fmt.Errorf("dao.GetValidatorDelegatorsTotal: %s", err.Error())
	}
	return smodels.PaginatableResponse{
		Items: items,
		Total: total,
	}, nil
}

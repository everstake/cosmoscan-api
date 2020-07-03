package services

import (
	"fmt"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/smodels"
	"github.com/shopspring/decimal"
)

const topProposedBlocksValidatorsKey = "topProposedBlocksValidatorsKey"
const rewardPerBlock = 4.0

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

func (s *ServiceFacade) GetValidatorBlocksStat(validatorAddress string) (stat smodels.ValidatorBlocksStat, err error) {
	validator, err := s.GetValidator(validatorAddress)
	if err != nil {
		return stat, fmt.Errorf("GetValidator: %s", err.Error())
	}
	stat.Proposed, err = s.dao.GetProposedBlocksTotal(filters.BlocksProposed{
		Proposers: []string{validator.ConsAddress},
	})
	if err != nil {
		return stat, fmt.Errorf("dao.GetProposedBlocksTotal: %s", err.Error())
	}
	stat.MissedValidations, err = s.dao.GetMissedBlocksCount(filters.MissedBlocks{
		Validators: []string{validator.ConsAddress},
	})
	if err != nil {
		return stat, fmt.Errorf("dao.GetMissedBlocksCount: %s", err.Error())
	}
	stat.Revenue = decimal.NewFromFloat(rewardPerBlock)
	return stat, nil
}

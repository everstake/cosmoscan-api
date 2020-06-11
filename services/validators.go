package services

import (
	"fmt"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/everstake/cosmoscan-api/services/node"
	"github.com/everstake/cosmoscan-api/smodels"
	"sort"
	"time"
)

const validatorsCacheKey = "validators"

func (s *ServiceFacade) UpdateValidatorsMap() {
	mp, err := s.makeValidatorMap()
	if err != nil {
		log.Error("UpdateValidatorsMap: makeValidatorMap: %s", err.Error())
		return
	}
	s.dao.CacheSet(validatorsCacheKey, mp, time.Minute*30)
}

func (s *ServiceFacade) GetValidatorMap() (map[string]node.Validator, error) {
	data, found := s.dao.CacheGet(validatorsCacheKey)
	if found {
		return data.(map[string]node.Validator), nil
	}
	mp, err := s.makeValidatorMap()
	if err != nil {
		return nil, fmt.Errorf("makeValidatorMap: %s", err.Error())
	}
	return mp, nil
}

func (s *ServiceFacade) makeValidatorMap() (map[string]node.Validator, error) {
	mp := make(map[string]node.Validator)
	validators, err := s.node.GetValidators()
	if err != nil {
		return nil, fmt.Errorf("node.GetValidators: %s", err.Error())
	}
	for _, validator := range validators {
		mp[validator.OperatorAddress] = validator
	}
	return mp, nil
}

func (s *ServiceFacade) GetStakingPie() (pie smodels.Pie, err error) {
	stakingPool, err := s.node.GetStakingPool()
	if err != nil {
		return pie, fmt.Errorf("node.GetStakingPool: %s", err.Error())
	}
	pie.Total = stakingPool.Result.BondedTokens
	validatorsMap, err := s.GetValidatorMap()
	if err != nil {
		return pie, fmt.Errorf("s.GetValidatorMap: %s", err.Error())
	}
	var validators []node.Validator
	for _, v := range validatorsMap {
		validators = append(validators, v)
	}
	sort.Slice(validators, func(i, j int) bool {
		return validators[i].DelegatorShares.GreaterThan(validators[j].DelegatorShares)
	})
	if len(validators) < 20 {
		return pie, fmt.Errorf("not enought validators")
	}
	parts := make([]smodels.PiePart, 20)
	for i := 0; i < 20; i++ {
		parts[i] = smodels.PiePart{
			Label: validators[i].OperatorAddress,
			Title: validators[i].Description.Moniker,
			Value: validators[i].DelegatorShares.Div(node.PrecisionDiv),
		}
	}
	pie.Parts = parts
	return pie, nil
}

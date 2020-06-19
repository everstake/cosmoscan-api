package services

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/services/node"
	"github.com/everstake/cosmoscan-api/smodels"
	"github.com/shopspring/decimal"
)

func (s *ServiceFacade) GetMetaData() (meta smodels.MetaData, err error) {
	states, err := s.dao.GetHistoricalStates(filters.HistoricalState{Limit: 1})
	if err != nil {
		return meta, fmt.Errorf("dao.GetHistoricalStates: %s", err.Error())
	}
	if len(states) != 0 {
		state := states[0]
		meta.CurrentPrice = state.Price
	}
	blocks, err := s.dao.GetBlocks(filters.Blocks{Limit: 2})
	if err != nil {
		return meta, fmt.Errorf("dao.GetBlocks: %s", err.Error())
	}
	if len(blocks) == 2 {
		meta.BlockTime = blocks[0].CreatedAt.Sub(blocks[1].CreatedAt).Seconds()
		meta.Height = blocks[0].ID
	}
	proposer := blocks[0].Proposer

	data, found := s.dao.CacheGet(validatorsCacheKey)
	if found {
		validators := data.(map[string]node.Validator)
		avgFee := decimal.Zero
		for _, validator := range validators {
			avgFee = avgFee.Add(validator.Commission.CommissionRates.Rate)
		}
		if len(validators) > 0 {
			meta.ValidatorAvgFee = avgFee.Div(decimal.New(int64(len(validators)), 0)).Mul(decimal.New(100, 0))
		}
		for _, validator := range validators {
			pk, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, validator.ConsensusPubkey)
			if err != nil {
				return meta, fmt.Errorf("types.GetPubKeyFromBech32(%s): %s", validator.ConsensusPubkey, err.Error())
			}
			if pk.Address().String() == proposer {
				meta.LatestValidator = validator.Description.Moniker
				break
			}
		}
	}
	proposals, err := s.dao.GetProposals(filters.Proposals{Limit: 1})
	if err != nil {
		return meta, fmt.Errorf("dao.GetProposals: %s", err.Error())
	}
	if len(proposals) != 0 {
		meta.LatestProposal = smodels.MetaDataProposal{
			Name: proposals[0].Title,
			ID:   proposals[0].ID,
		}
	}
	return meta, nil
}

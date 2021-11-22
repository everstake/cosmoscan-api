package services

import (
	"fmt"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/services/helpers"
	"github.com/everstake/cosmoscan-api/smodels"
	"github.com/shopspring/decimal"
	"strings"
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
	stat.Revenue = decimal.NewFromFloat(rewardPerBlock).Mul(decimal.NewFromInt(int64(stat.Proposed)))
	return stat, nil
}

func (s *ServiceFacade) GetBlock(height uint64) (block smodels.Block, err error) {
	dBlock, err := s.node.GetBlock(height)
	if err != nil {
		return block, fmt.Errorf("node.GetBlock: %s", err.Error())
	}
	validators, err := s.getConsensusValidatorMap()
	if err != nil {
		return block, fmt.Errorf("s.getConsensusValidatorMap: %s", err.Error())
	}
	proposerKey, err := helpers.B64ToHex(dBlock.Block.Header.ProposerAddress)
	if err != nil {
		return block, fmt.Errorf("helpers.B64ToHex: %s", err.Error())
	}
	hashHex, err := helpers.B64ToHex(dBlock.BlockID.Hash)
	if err != nil {
		return block, fmt.Errorf("helpers.B64ToHex: %s", err.Error())
	}
	var proposer, proposerAddress string
	validator, ok := validators[strings.ToUpper(proposerKey)]
	if ok {
		proposer = validator.Description.Moniker
		proposerAddress = validator.OperatorAddress
	}
	dTxs, err := s.dao.GetTransactions(filters.Transactions{Height: height})
	if err != nil {
		return block, fmt.Errorf("dao.GetTransactions: %s", err.Error())
	}
	var txs []smodels.TxItem
	for _, tx := range dTxs {
		txs = append(txs, smodels.TxItem{
			Hash:      tx.Hash,
			Status:    tx.Status,
			Fee:       tx.Fee,
			Height:    tx.Height,
			Messages:  tx.Messages,
			CreatedAt: dmodels.NewTime(tx.CreatedAt),
		})
	}
	return smodels.Block{
		Height:          dBlock.Block.Header.Height,
		Hash:            strings.ToUpper(hashHex),
		TotalTxs:        uint64(len(dBlock.Block.Data.Txs)),
		ChainID:         dBlock.Block.Header.ChainID,
		Proposer:        proposer,
		ProposerAddress: proposerAddress,
		Txs:             txs,
		CreatedAt:       dmodels.NewTime(dBlock.Block.Header.Time),
	}, nil
}

func (s *ServiceFacade) GetBlocks(filter filters.Blocks) (resp smodels.PaginatableResponse, err error) {
	dBlocks, err := s.dao.GetBlocks(filter)
	if err != nil {
		return resp, fmt.Errorf("dao.GetBlocks: %s", err.Error())
	}
	total, err := s.dao.GetBlocksCount(filter)
	if err != nil {
		return resp, fmt.Errorf("dao.GetBlocksCount: %s", err.Error())
	}
	validators, err := s.getConsensusValidatorMap()
	if err != nil {
		return resp, fmt.Errorf("s.getConsensusValidatorMap: %s", err.Error())
	}
	var blocks []smodels.BlockItem
	for _, b := range dBlocks {
		var proposer, proposerAddress string
		validator, ok := validators[b.Proposer]
		if ok {
			proposer = validator.Description.Moniker
			proposerAddress = validator.OperatorAddress
		}
		blocks = append(blocks, smodels.BlockItem{
			Height:          b.ID,
			Hash:            b.Hash,
			Proposer:        proposer,
			ProposerAddress: proposerAddress,
			CreatedAt:       dmodels.NewTime(b.CreatedAt),
		})
	}
	return smodels.PaginatableResponse{
		Items: blocks,
		Total: total,
	}, nil
}

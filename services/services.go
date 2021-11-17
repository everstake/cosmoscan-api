package services

import (
	"github.com/everstake/cosmoscan-api/config"
	"github.com/everstake/cosmoscan-api/dao"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/services/cmc"
	"github.com/everstake/cosmoscan-api/services/node"
	"github.com/everstake/cosmoscan-api/smodels"
	"github.com/shopspring/decimal"
)

type (
	Services interface {
		KeepHistoricalState()
		UpdateValidatorsMap()
		GetValidatorMap() (map[string]node.Validator, error)
		GetMetaData() (meta smodels.MetaData, err error)
		GetAggTransactionsFee(filter filters.Agg) (items []smodels.AggItem, err error)
		GetAggOperationsCount(filter filters.Agg) (items []smodels.AggItem, err error)
		GetAggTransfersVolume(filter filters.Agg) (items []smodels.AggItem, err error)
		GetHistoricalState() (state smodels.HistoricalState, err error)
		GetAggBlocksCount(filter filters.Agg) (items []smodels.AggItem, err error)
		GetAggBlocksDelay(filter filters.Agg) (items []smodels.AggItem, err error)
		GetAggUniqBlockValidators(filter filters.Agg) (items []smodels.AggItem, err error)
		GetAggDelegationsVolume(filter filters.DelegationsAgg) (items []smodels.AggItem, err error)
		GetAggUndelegationsVolume(filter filters.Agg) (items []smodels.AggItem, err error)
		GetNetworkStates(filter filters.Stats) (map[string][]decimal.Decimal, error)
		GetStakingPie() (pie smodels.Pie, err error)
		MakeUpdateBalances()
		GetSizeOfNode() (size float64, err error)
		MakeStats()
		UpdateProposals()
		GetProposals(filter filters.Proposals) (proposals []dmodels.Proposal, err error)
		GetProposalVotes(filter filters.ProposalVotes) (items []smodels.ProposalVote, err error)
		GetProposalDeposits(filter filters.ProposalDeposits) (deposits []dmodels.ProposalDeposit, err error)
		GetProposalsChartData() (items []smodels.ProposalChartData, err error)
		GetAggValidators33Power(filter filters.Agg) (items []smodels.AggItem, err error)
		GetValidators() (validators []smodels.Validator, err error)
		UpdateValidators()
		GetAvgOperationsPerBlock(filter filters.Agg) (items []smodels.AggItem, err error)
		GetAggWhaleAccounts(filter filters.Agg) (items []smodels.AggItem, err error)
		GetTopProposedBlocksValidators() (items []dmodels.ValidatorValue, err error)
		GetMostJailedValidators() (items []dmodels.ValidatorValue, err error)
		GetFeeRanges() (items []smodels.FeeRange, err error)
		GetValidatorsDelegatorsTotal() (values []dmodels.ValidatorValue, err error)
		GetValidator(address string) (validator smodels.Validator, err error)
		GetValidatorBalance(valAddress string) (balance smodels.Balance, err error)
		GetValidatorDelegationsAgg(validatorAddress string) (items []smodels.AggItem, err error)
		GetValidatorDelegatorsAgg(validatorAddress string) (items []smodels.AggItem, err error)
		GetValidatorBlocksStat(validatorAddress string) (stat smodels.ValidatorBlocksStat, err error)
		GetValidatorDelegators(filter filters.ValidatorDelegators) (resp smodels.PaginatableResponse, err error)
		GetAggBondedRatio(filter filters.Agg) (items []smodels.AggItem, err error)
		GetAggUnbondingVolume(filter filters.Agg) (items []smodels.AggItem, err error)
		GetBlock(height uint64) (block smodels.Block, err error)
		GetBlocks(filter filters.Blocks) (resp smodels.PaginatableResponse, err error)
		GetTransaction(hash string) (tx smodels.Tx, err error)
		GetTransactions(filter filters.Transactions) (resp smodels.PaginatableResponse, err error)
		GetAccount(address string) (account smodels.Account, err error)
	}
	CMC interface {
		GetCurrencies() (currencies []cmc.Currency, err error)
	}
	Node interface {
		GetCommunityPoolAmount() (amount decimal.Decimal, err error)
		GetValidators() (items []node.Validator, err error)
		GetInflation() (amount decimal.Decimal, err error)
		GetTotalSupply() (amount decimal.Decimal, err error)
		GetStakingPool() (sp node.StakingPool, err error)
		GetBalance(address string) (amount decimal.Decimal, err error)
		GetStake(address string) (amount decimal.Decimal, err error)
		GetUnbonding(address string) (amount decimal.Decimal, err error)
		GetProposals() (proposals node.ProposalsResult, err error)
		GetDelegatorValidatorStake(delegator string, validator string) (amount decimal.Decimal, err error)
		ProposalTallyResult(id uint64) (result node.ProposalTallyResult, err error)
		GetBlock(id uint64) (result node.Block, err error)
		GetTransaction(hash string) (result node.TxResult, err error)
		GetBalances(address string) (result node.AmountResult, err error)
		GetStakeRewards(address string) (amount decimal.Decimal, err error)
	}

	ServiceFacade struct {
		dao  dao.DAO
		cfg  config.Config
		cmc  CMC
		node Node
	}
)

func NewServices(d dao.DAO, cfg config.Config) (svc Services, err error) {
	return &ServiceFacade{
		dao:  d,
		cfg:  cfg,
		cmc:  cmc.NewCMC(cfg),
		node: node.NewAPI(cfg),
	}, nil
}

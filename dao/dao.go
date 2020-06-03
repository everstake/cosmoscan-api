package dao

import (
	"fmt"
	"github.com/everstake/cosmoscan-api/config"
	"github.com/everstake/cosmoscan-api/dao/cache"
	"github.com/everstake/cosmoscan-api/dao/clickhouse"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dao/mysql"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/smodels"
	"time"
)

type (
	DAO interface {
		Mysql
		Clickhouse
		Cache
	}
	Mysql interface {
		GetParsers() (parsers []dmodels.Parser, err error)
		GetParser(title string) (parser dmodels.Parser, err error)
		UpdateParser(parser dmodels.Parser) error
	}
	Clickhouse interface {
		CreateBlocks(blocks []dmodels.Block) error
		GetBlocks(filter filters.Blocks) (blocks []dmodels.Block, err error)
		CreateTransactions(transactions []dmodels.Transaction) error
		GetAggTransactionsFee(filter filters.Agg) (items []smodels.AggItem, err error)
		GetAggTransfersVolume(filter filters.Agg) (items []smodels.AggItem, err error)
		CreateTransfers(transfers []dmodels.Transfer) error
		CreateDelegations(delegations []dmodels.Delegation) error
		CreateDelegatorRewards(rewards []dmodels.DelegatorReward) error
		CreateValidatorRewards(rewards []dmodels.ValidatorReward) error
		CreateProposals(proposals []dmodels.Proposal) error
		CreateProposalDeposits(deposits []dmodels.ProposalDeposit) error
		CreateProposalVotes(votes []dmodels.ProposalVote) error
		CreateHistoricalStates(states []dmodels.HistoricalState) error
		GetHistoricalStates(state filters.HistoricalState) (states []dmodels.HistoricalState, err error)
		GetAggHistoricalStatesByField(filter filters.Agg, field string) (items []smodels.AggItem, err error)
	}

	Cache interface {
		CacheSet(key string, data interface{}, duration time.Duration)
		CacheGet(key string) (data interface{}, found bool)
	}

	daoImpl struct {
		Mysql
		Clickhouse
		Cache
	}
)

func NewDAO(cfg config.Config) (DAO, error) {
	mysqlDB, err := mysql.NewDB(cfg.Mysql)
	if err != nil {
		return nil, fmt.Errorf("mysql.NewDB: %s", err.Error())
	}
	ch, err := clickhouse.NewDB(cfg.Clickhouse)
	if err != nil {
		return nil, fmt.Errorf("clickhouse.NewDB: %s", err.Error())
	}
	return daoImpl{
		Mysql:      mysqlDB,
		Clickhouse: ch,
		Cache:      cache.New(),
	}, nil
}

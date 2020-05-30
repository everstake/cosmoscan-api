package dao

import (
	"fmt"
	"github.com/everstake/cosmoscan-api/config"
	"github.com/everstake/cosmoscan-api/dao/clickhouse"
	"github.com/everstake/cosmoscan-api/dao/mysql"
	"github.com/everstake/cosmoscan-api/dmodels"
)

type (
	DAO interface {
		Mysql
		Clickhouse
	}
	Mysql interface {
		GetParsers() (parsers []dmodels.Parser, err error)
		GetParser(title string) (parser dmodels.Parser, err error)
		UpdateParser(parser dmodels.Parser) error
	}
	Clickhouse interface{
		CreateBlocks(blocks []dmodels.Block) error
		CreateTransactions(transactions []dmodels.Transaction) error
		CreateTransfers(transfers []dmodels.Transfer) error
		CreateDelegations(delegations []dmodels.Delegation) error
		CreateDelegatorRewards(rewards []dmodels.DelegatorReward) error
		CreateValidatorRewards(rewards []dmodels.ValidatorReward) error
		CreateProposal(proposals []dmodels.Proposal) error
		CreateProposalDeposits(deposits []dmodels.ProposalDeposit) error
		CreateProposalVotes(votes []dmodels.ProposalVote) error
	}

	daoImpl struct {
		Mysql
		Clickhouse
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
	}, nil
}

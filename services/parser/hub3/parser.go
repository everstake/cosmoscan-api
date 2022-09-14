package hub3

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/everstake/cosmoscan-api/config"
	"github.com/everstake/cosmoscan-api/dao"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/everstake/cosmoscan-api/services/helpers"
	"github.com/everstake/cosmoscan-api/services/node"
	"github.com/shopspring/decimal"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/bytes"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const repeatDelay = time.Second * 5
const ParserTitle = "hub3"
const AddressLength = 45

const batchTxs = 50
const precision = 6

var precisionDiv = decimal.New(1, precision)

type (
	Parser struct {
		cfg       config.Config
		api       api
		dao       dao.DAO
		fetcherCh chan uint64
		saverCh   chan data
		accounts  map[string]struct{}
		ctx       context.Context
		cancel    context.CancelFunc
		wg        *sync.WaitGroup
	}
	api interface {
		GetLatestBlock() (block Block, err error)
		GetBlock(height uint64) (block Block, err error)
		GetTx(hash string) (txs Tx, err error)
		GetValidatorset(height uint64) (set Validatorsets, err error)
	}
	data struct {
		height           uint64
		blocks           []dmodels.Block
		transactions     []dmodels.Transaction
		transfers        []dmodels.Transfer
		delegations      []dmodels.Delegation
		delegatorRewards []dmodels.DelegatorReward
		validatorRewards []dmodels.ValidatorReward
		proposals        []dmodels.HistoryProposal
		proposalVotes    []dmodels.ProposalVote
		proposalDeposits []dmodels.ProposalDeposit
		jailers          []dmodels.Jailer
		missedBlocks     []dmodels.MissedBlock
		accountTxs       []dmodels.AccountTx
	}
)

func NewParser(cfg config.Config, d dao.DAO) *Parser {
	ctx, cancel := context.WithCancel(context.Background())
	return &Parser{
		cfg:       cfg,
		dao:       d,
		api:       NewAPI(cfg.Parser.Node),
		fetcherCh: make(chan uint64, 5000),
		saverCh:   make(chan data, 5000),
		accounts:  make(map[string]struct{}),
		ctx:       ctx,
		cancel:    cancel,
		wg:        &sync.WaitGroup{},
	}
}

func (p *Parser) Run() error {
	model, err := p.dao.GetParser(ParserTitle)
	if err != nil {
		return fmt.Errorf("parser not found")
	}
	fmt.Printf("model = %v\n", model) ////////////////////////
	for i := uint64(0); i < p.cfg.Parser.Fetchers; i++ {
		go p.runFetcher()
	}
	if model.Height == 0 {
		err = p.parseGenesisState()
		if err != nil {
			return fmt.Errorf("parseGenesisState: %s", err.Error())
		}
		p.setAccounts()
	}
	go p.saving()
	for {
		latestBlock, err := p.api.GetLatestBlock()
		fmt.Printf("latestBlock = %v\n", latestBlock.Block.Header.Height) ////////////////////////////
		if err != nil {
			log.Error("Parser: api.GetLatestBlock: %s", err.Error())
			continue
		}
		if model.Height >= latestBlock.Block.Header.Height {
			<-time.After(time.Second)
			continue
		}
		for ; model.Height < latestBlock.Block.Header.Height; model.Height++ {
			select {
			case <-p.ctx.Done():
				return nil
			case p.fetcherCh <- model.Height + 1:
			}
		}
	}
}

func (p *Parser) Title() string {
	return "Parser"
}

func (p *Parser) Stop() error {
	p.cancel()
	p.wg.Wait()
	return nil
}

func (p *Parser) runFetcher() {
	for {
		select {
		case <-p.ctx.Done():
			return
		default:
		}
		height := <-p.fetcherCh
		for {
			var d data
			d.height = height
			fmt.Printf("before getBlock %d\n", height) ////////////////////
			block, err := p.api.GetBlock(height)
			fmt.Printf("after getBlock %v\n", block) ////////////////////
			if err != nil {
				log.Error("Parser: fetcher: api.GetBlock: %s", err.Error())
				<-time.After(time.Second)
				continue
			}
			validatorsSets, err := p.api.GetValidatorset(height)
			if err != nil {
				log.Error("Parser: fetcher: api.GetValidatorset: %s", err.Error())
				<-time.After(time.Second)
				continue
			}

			d.blocks = append(d.blocks, dmodels.Block{
				ID:        block.Block.Header.Height,
				Hash:      block.BlockID.Hash,
				Proposer:  block.Block.Header.ProposerAddress,
				CreatedAt: block.Block.Header.Time,
			})

			// find missed blocks
			set := make(map[string]struct{})
			for _, s := range validatorsSets.Validators {
				address, err := helpers.GetHexAddressFromBase64PK(s.PubKey.Key)
				if err != nil {
					log.Warn("Parser: helpers.GetHexAddressFromBase64PK: %s", err.Error())
					continue
				}
				set[address] = struct{}{}
			}

			precommits := make(map[string]struct{})
			for _, precommit := range block.Block.LastCommit.Signatures {
				precommits[precommit.ValidatorAddress] = struct{}{}
			}

			for address := range set {
				_, ok := precommits[address]
				if !ok {
					id := makeHash(fmt.Sprintf("%d.%s", block.Block.Header.Height, address))
					d.missedBlocks = append(d.missedBlocks, dmodels.MissedBlock{
						ID:        id,
						Height:    block.Block.Header.Height,
						Validator: address,
						CreatedAt: block.Block.Header.Time,
					})
				}
			}

			fail := false
			for _, txData := range block.Block.Data.Txs {
				decodedTx, _ := base64.StdEncoding.DecodeString(txData)
				sha256Hash := crypto.Sha256(decodedTx)
				hash := bytes.HexBytes(sha256Hash)

				tx, err := p.api.GetTx(hash.String())
				if err != nil {
					log.Error("Parser: fetcher: api.GetTxs: %s", err.Error())
					<-time.After(time.Second)
					fail = true
					break
				}

				success := tx.TxResponse.Code == 0

				fee, err := calculateAtomAmount(tx.Tx.AuthInfo.Fee.Amount)
				if err != nil {
					log.Warn("Parser: height: %d, calculateAtomAmount: %s", tx.TxResponse.Height, err.Error())
				}

				if tx.TxResponse.Hash == "" {
					log.Error("Parser: fetcher: empty tx hash")
					<-time.After(time.Second)
					fail = true
					break
				}

				d.transactions = append(d.transactions, dmodels.Transaction{
					Hash:      tx.TxResponse.Hash,
					Status:    success,
					Height:    tx.TxResponse.Height,
					Messages:  uint64(len(tx.TxResponse.Tx.Body.Messages)),
					Fee:       fee,
					GasUsed:   tx.TxResponse.GasUsed,
					GasWanted: tx.TxResponse.GasWanted,
					CreatedAt: tx.TxResponse.Timestamp,
				})

				// account - transactions relations
				accTxsMap := make(map[string]struct{})
				for _, msg := range tx.Tx.Body.Messages {
					for _, address := range fetchAddressesFromMessage(msg) {
						accTxsMap[address] = struct{}{}
					}
				}
				for address := range accTxsMap {
					d.accountTxs = append(d.accountTxs, dmodels.AccountTx{
						Account: address,
						TxHash:  tx.TxResponse.Hash,
					})
				}

				if success {
					for i, msg := range tx.Tx.Body.Messages {
						var baseMsg BaseMsg
						err = json.Unmarshal(msg, &baseMsg)
						if err != nil {
							log.Error("Parser: BaseMsg: json.Unmarshal: %s", err.Error())
							<-time.After(time.Second)
							fail = true
							break
						}
						switch baseMsg.Type {
						case SendMsg:
							err = d.parseMsgSend(i, tx, msg)
						case MultiSendMsg:
							err = d.parseMultiSendMsg(i, tx, msg)
						case DelegateMsg:
							err = d.parseDelegateMsg(i, tx, msg)
						case UndelegateMsg:
							err = d.parseUndelegateMsg(i, tx, msg)
						case BeginRedelegateMsg:
							err = d.parseBeginRedelegateMsg(i, tx, msg)
						case WithdrawDelegationRewardMsg:
							err = d.parseWithdrawDelegationRewardMsg(i, tx, msg)
						case WithdrawValidatorCommissionMsg:
							err = d.parseWithdrawValidatorCommissionMsg(i, tx, msg)
						case SubmitProposalMsg:
							err = d.parseSubmitProposalMsg(i, tx, msg)
						case DepositMsg:
							err = d.parseDepositMsg(i, tx, msg)
						case VoteMsg:
							err = d.parseVoteMsg(i, tx, msg)
						case UnJailMsg:
							err = d.parseUnjailMsg(i, tx, msg)
						}
						if err != nil {
							log.Error("Parser: (height: %d): %s", tx.TxResponse.Height, err.Error())
							<-time.After(time.Second)
							fail = true
							break
						}
					}
				}
			}
			if fail { // try again
				continue
			}

			p.saverCh <- d
			break
		}

	}
}

func (p *Parser) saving() {
	var model dmodels.Parser
	for {
		var err error
		model, err = p.dao.GetParser(ParserTitle)
		if err != nil {
			log.Error("Parser: saving: dao.GetParser: %s", err.Error())
			<-time.After(time.Second * 5)
			continue
		}
		break
	}
	p.setAccounts()

	ticker := time.After(time.Second)

	var dataset []data

	for {
		select {
		case <-p.ctx.Done():
			return
		case d := <-p.saverCh:
			dataset = append(dataset, d)
			continue
		case <-ticker:
			sort.Slice(dataset, func(i, j int) bool {
				return dataset[i].height < dataset[j].height
			})
			ticker = time.After(time.Second * 2)
		}

		var count int
		for i, item := range dataset {
			if item.height == model.Height+uint64(i+1) {
				count = i + 1
			} else {
				break
			}
		}

		if count == 0 {
			continue
		}

		if count > int(p.cfg.Parser.Batch) {
			count = int(p.cfg.Parser.Batch)
		}

		var singleData data
		for _, item := range dataset[:count] {
			singleData.blocks = append(singleData.blocks, item.blocks...)
			singleData.proposals = append(singleData.proposals, item.proposals...)
			singleData.delegations = append(singleData.delegations, item.delegations...)
			singleData.jailers = append(singleData.jailers, item.jailers...)
			singleData.transactions = append(singleData.transactions, item.transactions...)
			singleData.delegatorRewards = append(singleData.delegatorRewards, item.delegatorRewards...)
			singleData.validatorRewards = append(singleData.validatorRewards, item.validatorRewards...)
			singleData.transfers = append(singleData.transfers, item.transfers...)
			singleData.proposalVotes = append(singleData.proposalVotes, item.proposalVotes...)
			singleData.proposalDeposits = append(singleData.proposalDeposits, item.proposalDeposits...)
			singleData.missedBlocks = append(singleData.missedBlocks, item.missedBlocks...)
			singleData.accountTxs = append(singleData.accountTxs, item.accountTxs...)
		}
		p.wg.Add(1)
		var err error
		for {
			err = p.dao.CreateBlocks(singleData.blocks)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateBlocks: %s", err.Error())
			<-time.After(repeatDelay)
		}
		for {
			err = p.dao.CreateTransactions(singleData.transactions)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateTransactions: %s", err.Error())
			<-time.After(repeatDelay)
		}
		for {
			err = p.dao.CreateTransfers(singleData.transfers)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateTransfers: %s", err.Error())
			<-time.After(repeatDelay)
		}
		for {
			err = p.dao.CreateDelegations(singleData.delegations)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateDelegations: %s", err.Error())
			<-time.After(repeatDelay)
		}
		for {
			err = p.dao.CreateDelegatorRewards(singleData.delegatorRewards)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateDelegatorRewards: %s", err.Error())
			<-time.After(repeatDelay)
		}
		for {
			err = p.dao.CreateValidatorRewards(singleData.validatorRewards)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateValidatorRewards: %s", err.Error())
			<-time.After(repeatDelay)
		}
		for {
			err = p.dao.CreateHistoryProposals(singleData.proposals)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateProposals: %s", err.Error())
			<-time.After(repeatDelay)
		}
		for {
			err = p.dao.CreateProposalDeposits(singleData.proposalDeposits)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateProposalDeposits: %s", err.Error())
			<-time.After(repeatDelay)
		}
		for {
			err = p.dao.CreateProposalVotes(singleData.proposalVotes)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateProposalVotes: %s", err.Error())
			<-time.After(repeatDelay)
		}
		for {
			err = p.dao.CreateJailers(singleData.jailers)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateJailers: %s", err.Error())
			<-time.After(repeatDelay)
		}
		for {
			err = p.dao.CreateMissedBlocks(singleData.missedBlocks)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateMissedBlocks: %s", err.Error())
			<-time.After(repeatDelay)
		}
		for {
			err = p.dao.CreateAccountTxs(singleData.accountTxs)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateAccountTxs: %s", err.Error())
			<-time.After(repeatDelay)
		}
		p.saveNewAccounts(singleData)
		for {
			model.Height += uint64(count)
			err = p.dao.UpdateParser(model)
			if err == nil {
				break
			}
			log.Error("Parser: dao.UpdateParser: %s", err.Error())
			<-time.After(repeatDelay)
		}
		dataset = append(dataset[count:])
		p.wg.Done()
	}
}

func (p *Parser) setAccounts() {
	var accounts []dmodels.Account
	var err error
	for {
		accounts, err = p.dao.GetAccounts(filters.Accounts{})
		if err != nil {
			log.Error("Parser: setAccounts: dao.GetAccounts: %s", err.Error())
			<-time.After(repeatDelay)
			continue
		}
		break
	}
	for _, account := range accounts {
		p.accounts[account.Address] = struct{}{}
	}
}

func (p *Parser) saveNewAccounts(data data) {
	var newAccounts []dmodels.Account
	addAccount := func(acc string, tm time.Time) {
		_, ok := p.accounts[acc]
		if !ok {
			p.accounts[acc] = struct{}{}
			newAccounts = append(newAccounts, dmodels.Account{
				Address:   acc,
				CreatedAt: tm,
			})
		}
	}
	for _, delegation := range data.delegations {
		addAccount(delegation.Delegator, delegation.CreatedAt)
	}
	for _, transfer := range data.transfers {
		if strings.TrimSpace(transfer.From) != "" {
			addAccount(transfer.From, transfer.CreatedAt)
		}
		if strings.TrimSpace(transfer.To) != "" {
			addAccount(transfer.To, transfer.CreatedAt)
		}
	}
	for _, reward := range data.delegatorRewards {
		addAccount(reward.Delegator, reward.CreatedAt)
	}
	for {
		err := p.dao.CreateAccounts(newAccounts)
		if err == nil {
			break
		}
		log.Error("Parser: dao.CreateAccounts: %s", err.Error())
		<-time.After(repeatDelay)
	}
}

func (d *data) parseMsgSend(index int, tx Tx, data []byte) (err error) {
	var m MsgSend
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	currency, amount, err := calculateAmount(m.Amount)
	if err != nil {
		return fmt.Errorf("calculateAtomAmount: %s", err.Error())
	}
	id := makeHash(fmt.Sprintf("%s.%d", tx.TxResponse.Hash, index))
	d.transfers = append(d.transfers, dmodels.Transfer{
		ID:        id,
		TxHash:    tx.TxResponse.Hash,
		From:      m.FromAddress,
		To:        m.ToAddress,
		Amount:    amount,
		Currency:  currency,
		CreatedAt: tx.TxResponse.Timestamp,
	})
	return nil
}

func (d *data) parseMultiSendMsg(index int, tx Tx, data []byte) (err error) {
	var m MsgMultiSendValue
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	for i, input := range m.Inputs {
		id := makeHash(fmt.Sprintf("%s.%d.i.%d", tx.TxResponse.Hash, index, i))
		currency, amount, err := calculateAmount(input.Coins)
		if err != nil {
			return fmt.Errorf("calculateAtomAmount: %s", err.Error())
		}
		d.transfers = append(d.transfers, dmodels.Transfer{
			ID:        id,
			TxHash:    tx.TxResponse.Hash,
			From:      input.Address,
			To:        "",
			Amount:    amount,
			Currency:  currency,
			CreatedAt: tx.TxResponse.Timestamp,
		})
	}
	for i, output := range m.Outputs {
		id := makeHash(fmt.Sprintf("%s.%d.o.%d", tx.TxResponse.Hash, index, i))
		currency, amount, err := calculateAmount(output.Coins)
		if err != nil {
			return fmt.Errorf("calculateAtomAmount: %s", err.Error())
		}
		d.transfers = append(d.transfers, dmodels.Transfer{
			ID:        id,
			TxHash:    tx.TxResponse.Hash,
			From:      "",
			To:        output.Address,
			Amount:    amount,
			Currency:  currency,
			CreatedAt: tx.TxResponse.Timestamp,
		})
	}
	return nil
}

func (d *data) parseDelegateMsg(index int, tx Tx, data []byte) (err error) {
	var m MsgDelegate
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	amount, err := m.Amount.getAmount()
	if err != nil {
		return fmt.Errorf("getAmount: %s", err.Error())
	}
	id := makeHash(fmt.Sprintf("%s.%d", tx.TxResponse.Hash, index))
	d.delegations = append(d.delegations, dmodels.Delegation{
		ID:        id,
		TxHash:    tx.TxResponse.Hash,
		Delegator: m.DelegatorAddress,
		Validator: m.ValidatorAddress,
		Amount:    amount,
		CreatedAt: tx.TxResponse.Timestamp,
	})
	return nil
}

func (d *data) parseUndelegateMsg(index int, tx Tx, data []byte) (err error) {
	var m MsgUndelegate
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	amount, err := m.Amount.getAmount()
	if err != nil {
		return fmt.Errorf("getAmount: %s", err.Error())
	}
	id := makeHash(fmt.Sprintf("%s.%d", tx.TxResponse.Hash, index))
	d.delegations = append(d.delegations, dmodels.Delegation{
		ID:        id,
		TxHash:    tx.TxResponse.Hash,
		Delegator: m.DelegatorAddress,
		Validator: m.ValidatorAddress,
		Amount:    amount.Mul(decimal.NewFromFloat(-1)),
		CreatedAt: tx.TxResponse.Timestamp,
	})
	return nil
}

func (d *data) parseBeginRedelegateMsg(index int, tx Tx, data []byte) (err error) {
	var m MsgBeginRedelegate
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	amount, err := m.Amount.getAmount()
	if err != nil {
		return fmt.Errorf("getAmount: %s", err.Error())
	}
	id := makeHash(fmt.Sprintf("%s.%d.s", tx.TxResponse.Hash, index))
	d.delegations = append(d.delegations, dmodels.Delegation{
		ID:        id,
		TxHash:    tx.TxResponse.Hash,
		Delegator: m.DelegatorAddress,
		Validator: m.ValidatorSrcAddress,
		Amount:    amount.Mul(decimal.NewFromFloat(-1)),
		CreatedAt: tx.TxResponse.Timestamp,
	})
	id = makeHash(fmt.Sprintf("%s.%d.d", tx.TxResponse.Hash, index))
	d.delegations = append(d.delegations, dmodels.Delegation{
		ID:        id,
		TxHash:    tx.TxResponse.Hash,
		Delegator: m.DelegatorAddress,
		Validator: m.ValidatorDstAddress,
		Amount:    amount,
		CreatedAt: tx.TxResponse.Timestamp,
	})
	return nil
}

func (d *data) parseWithdrawDelegationRewardMsg(index int, tx Tx, data []byte) (err error) {
	var m MsgWithdrawDelegationReward
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}

	mp := make(map[string]decimal.Decimal)
	for _, log := range tx.TxResponse.Logs {
		for _, event := range log.Events {
			if event.Type == "withdraw_rewards" {
				for i := 0; i < len(event.Attributes); i += 2 {
					amount, err := strToAmount(event.Attributes[i].Value)
					if err != nil {
						return fmt.Errorf("strToAmount: %s", err.Error())
					}
					if event.Attributes[i+1].Key != "validator" {
						return fmt.Errorf("not found validator in events")
					}
					mp[event.Attributes[i+1].Value] = amount
				}
				break
			}
		}
	}

	amount, ok := mp[m.ValidatorAddress]
	if !ok {
		return fmt.Errorf("not found validator %s in map", m.ValidatorAddress)
	}

	id := makeHash(fmt.Sprintf("%s.%d", tx.TxResponse.Hash, index))
	d.delegatorRewards = append(d.delegatorRewards, dmodels.DelegatorReward{
		ID:        id,
		TxHash:    tx.TxResponse.Hash,
		Delegator: m.DelegatorAddress,
		Validator: m.ValidatorAddress,
		Amount:    amount,
		CreatedAt: tx.TxResponse.Timestamp,
	})
	return nil
}

func (d *data) parseSubmitProposalMsg(index int, tx Tx, data []byte) (err error) {
	var m MsgSubmitProposal
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	var id uint64
	for _, log := range tx.TxResponse.Logs {
		for _, event := range log.Events {
			if event.Type == "submit_proposal" {
				for _, att := range event.Attributes {
					if att.Key == "proposal_id" {
						id, err = strconv.ParseUint(att.Value, 10, 64)
						if err != nil {
							return fmt.Errorf("strconv.ParseUint: %s", err.Error())
						}
					}
				}
			}
		}
	}
	if id == 0 {
		return fmt.Errorf("not found proposal_id")
	}
	amount, err := calculateAtomAmount(m.Content.Value.Amount)
	if err != nil {
		return fmt.Errorf("calculateAtomAmount: %s", err.Error())
	}
	initDeposit, err := calculateAtomAmount(m.InitialDeposit)
	if err != nil {
		return fmt.Errorf("calculateAtomAmount: %s", err.Error())
	}
	d.proposals = append(d.proposals, dmodels.HistoryProposal{
		ID:          id,
		TxHash:      tx.TxResponse.Hash,
		Title:       m.Content.Value.Title,
		Description: m.Content.Value.Description,
		Recipient:   m.Content.Value.Recipient,
		Amount:      amount,
		InitDeposit: initDeposit,
		Proposer:    m.Proposer,
		CreatedAt:   tx.TxResponse.Timestamp,
	})
	return nil
}

func (d *data) parseVoteMsg(index int, tx Tx, data []byte) (err error) {
	var m MsgVote
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	var option string
	switch m.Option {
	case "VOTE_OPTION_YES":
		option = "Yes"
	case "VOTE_OPTION_ABSTAIN":
		option = "Abstain"
	case "VOTE_OPTION_NO":
		option = "No"
	case "VOTE_OPTION_NO_WITH_VETO":
		option = "NoWithVeto"
	default:
		return fmt.Errorf("unknown type of option: %d", m.Option)
	}
	id := makeHash(fmt.Sprintf("%s.%d.s", tx.TxResponse.Hash, index))
	d.proposalVotes = append(d.proposalVotes, dmodels.ProposalVote{
		ID:         id,
		ProposalID: m.ProposalID,
		Voter:      m.Voter,
		TxHash:     tx.TxResponse.Hash,
		Option:     option,
		CreatedAt:  dmodels.NewTime(tx.TxResponse.Timestamp),
	})
	return nil
}

func (d *data) parseDepositMsg(index int, tx Tx, data []byte) (err error) {
	var m MsgDeposit
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	amount := decimal.Zero
	for _, a := range m.Amount {
		amt, err := a.getAmount()
		if err != nil {
			return fmt.Errorf("getAmount: %s", err.Error())
		}
		amount = amount.Add(amt)
	}

	id := makeHash(fmt.Sprintf("%s.%d.s", tx.TxResponse.Hash, index))
	d.proposalDeposits = append(d.proposalDeposits, dmodels.ProposalDeposit{
		ID:         id,
		ProposalID: m.ProposalID,
		Depositor:  m.Depositor,
		Amount:     amount,
		CreatedAt:  dmodels.NewTime(tx.TxResponse.Timestamp),
	})
	return nil
}

func (d *data) parseWithdrawValidatorCommissionMsg(index int, tx Tx, data []byte) (err error) {
	var m MsgWithdrawValidatorCommission
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	var amount decimal.Decimal
	found := false
	for _, log := range tx.TxResponse.Logs {
		for _, event := range log.Events {
			if event.Type == "withdraw_commission" {
				for _, att := range event.Attributes {
					if att.Key == "amount" {
						amount, err = strToAmount(att.Value)
						if err != nil {
							return fmt.Errorf("strToAmount: %s", err.Error())
						}
						found = true
					}
				}
			}
		}
	}
	if !found {
		return fmt.Errorf("amount not found")
	}
	id := makeHash(fmt.Sprintf("%s.%d", tx.TxResponse.Hash, index))
	d.validatorRewards = append(d.validatorRewards, dmodels.ValidatorReward{
		TxHash:    tx.TxResponse.Hash,
		ID:        id,
		Address:   m.ValidatorAddress,
		Amount:    amount,
		CreatedAt: tx.TxResponse.Timestamp,
	})
	return nil
}

func (d *data) parseUnjailMsg(index int, tx Tx, data []byte) (err error) {
	var m MsgUnjail
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	id := makeHash(fmt.Sprintf("%s.%d", tx.TxResponse.Hash, index))
	d.jailers = append(d.jailers, dmodels.Jailer{
		ID:        id,
		Address:   m.ValidatorAddr,
		CreatedAt: tx.TxResponse.Timestamp,
	})
	return nil
}

func calculateAtomAmount(amountItems []Amount) (decimal.Decimal, error) {
	volume := decimal.Zero
	for _, item := range amountItems {
		if item.Denom == "" && item.Amount.IsZero() { // example height=1245781
			break
		}
		if item.Denom != node.MainUnit {
			return volume, fmt.Errorf("unknown demon (currency): %s", item.Denom)
		}
		volume = volume.Add(item.Amount)
	}
	volume = volume.Div(precisionDiv)
	return volume, nil
}

func calculateAmount(amountItems []Amount) (string, decimal.Decimal, error) {
	volume := decimal.Zero
	var lastCurrency string
	for _, item := range amountItems {
		if item.Denom == "" {
			return lastCurrency, volume, errors.New("empty denom")
		}
		if lastCurrency == "" {
			lastCurrency = item.Denom
		} else if item.Denom != lastCurrency {
			return lastCurrency, volume, fmt.Errorf("different currencies: %s, %s", lastCurrency, item.Denom)
		}
		volume = volume.Add(item.Amount)
	}
	if lastCurrency == node.MainUnit {
		volume = volume.Div(precisionDiv)
		lastCurrency = config.Currency
	} else {
		if volume.GreaterThanOrEqual(decimal.New(1, 20)) {
			volume = decimal.Zero
		}
	}
	return lastCurrency, volume, nil
}

func (a Amount) getAmount() (decimal.Decimal, error) {
	if a.Denom == "" && a.Amount.IsZero() {
		return decimal.Zero, nil
	}
	if a.Denom != node.MainUnit {
		return decimal.Zero, fmt.Errorf("unknown demon (currency): %s", a.Denom)
	}
	a.Amount = a.Amount.Div(precisionDiv)
	return a.Amount, nil
}

func strToAmount(str string) (decimal.Decimal, error) {
	if str == "" {
		return decimal.Zero, nil
	}
	val := strings.TrimSuffix(str, node.MainUnit)
	amount, err := decimal.NewFromString(val)
	if err != nil {
		index := strings.LastIndex(val, ",")
		if index == -1 {
			return decimal.Zero, nil
		} else {
			val = val[index+1:]
		}
		amount, err = decimal.NewFromString(val)
		if err != nil {
			return amount, fmt.Errorf("decimal.NewFromString: %s", err.Error())
		}
	}
	amount = amount.Div(precisionDiv)
	return amount, nil
}

func makeHash(str string) string {
	hash := sha1.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

func fetchAddressesFromMessage(msg json.RawMessage) []string {
	var obj map[string]interface{}
	json.Unmarshal(msg, &obj)
	return getAddresses(obj)
}

func getAddresses(v interface{}) []string {
	var addresses []string
	switch val := v.(type) {
	case map[string]interface{}:
		for _, vi := range val {
			addresses = append(addresses, getAddresses(vi)...)
		}
	case string:
		if len(val) == AddressLength && strings.HasPrefix(val, types.Bech32MainPrefix) {
			addresses = append(addresses, val)
		}
	}
	return addresses
}

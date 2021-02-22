package hub3

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/everstake/cosmoscan-api/config"
	"github.com/everstake/cosmoscan-api/dao"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/shopspring/decimal"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const repeatDelay = time.Second * 5
const ParserTitle = "hub3"

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
		GetTxs(filter TxsFilter) (txs TxsBatch, err error)
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
			block, err := p.api.GetBlock(height)
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
			for _, s := range validatorsSets.Result.Validators {
				address, err := types.ConsAddressFromBech32(s.Address)
				if err != nil {
					log.Warn("Parser: types.ConsAddressFromBech32: %s", err.Error())
					continue
				}
				set[strings.ToUpper(hex.EncodeToString(address.Bytes()))] = struct{}{}
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

			pages := int(math.Ceil(float64(len(block.Block.Data.Txs)) / float64(batchTxs)))
			for page := 1; page <= pages; page++ {
				txs, err := p.api.GetTxs(TxsFilter{
					Limit:  batchTxs,
					Page:   uint64(page),
					Height: height,
				})
				if err != nil {
					log.Error("Parser: fetcher: api.GetTxs: %s", err.Error())
					<-time.After(time.Second)
					continue
				}

				for _, tx := range txs.Txs {

					success := tx.Code == 0

					fee, err := calculateAmount(tx.Tx.Value.Fee.Amount)
					if err != nil {
						log.Warn("Parser: height: %d, calculateAmount: %s", tx.Height, err.Error())
					}

					if tx.Hash == "" {
						log.Error("Parser: fetcher: empty tx hash")
						<-time.After(time.Second)
						continue
					}

					d.transactions = append(d.transactions, dmodels.Transaction{
						Hash:      tx.Hash,
						Status:    success,
						Height:    tx.Height,
						Messages:  uint64(len(tx.Tx.Value.Msg)),
						Fee:       fee,
						GasUsed:   tx.GasUsed,
						GasWanted: tx.GasWanted,
						CreatedAt: tx.Timestamp,
					})

					if success {
						for i, msg := range tx.Tx.Value.Msg {
							switch msg.Type {
							case SendMsg:
								err = d.parseMsgSend(i, tx, msg.Value)
							case MultiSendMsg:
								err = d.parseMultiSendMsg(i, tx, msg.Value)
							case DelegateMsg:
								err = d.parseDelegateMsg(i, tx, msg.Value)
							case UndelegateMsg:
								err = d.parseUndelegateMsg(i, tx, msg.Value)
							case BeginRedelegateMsg:
								err = d.parseBeginRedelegateMsg(i, tx, msg.Value)
							case WithdrawDelegationRewardMsg:
								err = d.parseWithdrawDelegationRewardMsg(i, tx, msg.Value)
							case WithdrawValidatorCommissionMsg:
								err = d.parseWithdrawValidatorCommissionMsg(i, tx, msg.Value)
							case SubmitProposalMsg:
								err = d.parseSubmitProposalMsg(i, tx, msg.Value)
							case DepositMsg:
								err = d.parseDepositMsg(i, tx, msg.Value)
							case VoteMsg:
								err = d.parseVoteMsg(i, tx, msg.Value)
							case UnJailMsg:
								err = d.parseUnjailMsg(i, tx, msg.Value)
							}
							if err != nil {
								log.Error("%s, (height: %d): %s", msg.Type, tx.Height, err.Error())
								<-time.After(time.Second)
								continue
							}
						}
					}
				}
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
	amount, err := calculateAmount(tx.Tx.Value.Fee.Amount)
	if err != nil {
		return fmt.Errorf("calculateAmount: %s", err.Error())
	}
	id := makeHash(fmt.Sprintf("%s.%d", tx.Hash, index))
	d.transfers = append(d.transfers, dmodels.Transfer{
		ID:        id,
		TxHash:    tx.Hash,
		From:      m.FromAddress,
		To:        m.ToAddress,
		Amount:    amount,
		CreatedAt: tx.Timestamp,
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
		id := makeHash(fmt.Sprintf("%s.%d.i.%d", tx.Hash, index, i))
		amount, err := calculateAmount(input.Coins)
		if err != nil {
			return fmt.Errorf("calculateAmount: %s", err.Error())
		}
		d.transfers = append(d.transfers, dmodels.Transfer{
			ID:        id,
			TxHash:    tx.Hash,
			From:      input.Address,
			To:        "",
			Amount:    amount,
			CreatedAt: tx.Timestamp,
		})
	}
	for i, output := range m.Outputs {
		id := makeHash(fmt.Sprintf("%s.%d.o.%d", tx.Hash, index, i))
		amount, err := calculateAmount(output.Coins)
		if err != nil {
			return fmt.Errorf("calculateAmount: %s", err.Error())
		}
		d.transfers = append(d.transfers, dmodels.Transfer{
			ID:        id,
			TxHash:    tx.Hash,
			From:      "",
			To:        output.Address,
			Amount:    amount,
			CreatedAt: tx.Timestamp,
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
	id := makeHash(fmt.Sprintf("%s.%d", tx.Hash, index))
	d.delegations = append(d.delegations, dmodels.Delegation{
		ID:        id,
		TxHash:    tx.Hash,
		Delegator: m.DelegatorAddress,
		Validator: m.ValidatorAddress,
		Amount:    amount,
		CreatedAt: tx.Timestamp,
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
	id := makeHash(fmt.Sprintf("%s.%d", tx.Hash, index))
	d.delegations = append(d.delegations, dmodels.Delegation{
		ID:        id,
		TxHash:    tx.Hash,
		Delegator: m.DelegatorAddress,
		Validator: m.ValidatorAddress,
		Amount:    amount.Mul(decimal.NewFromFloat(-1)),
		CreatedAt: tx.Timestamp,
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
	id := makeHash(fmt.Sprintf("%s.%d.s", tx.Hash, index))
	d.delegations = append(d.delegations, dmodels.Delegation{
		ID:        id,
		TxHash:    tx.Hash,
		Delegator: m.DelegatorAddress,
		Validator: m.ValidatorSrcAddress,
		Amount:    amount.Mul(decimal.NewFromFloat(-1)),
		CreatedAt: tx.Timestamp,
	})
	id = makeHash(fmt.Sprintf("%s.%d.d", tx.Hash, index))
	d.delegations = append(d.delegations, dmodels.Delegation{
		ID:        id,
		TxHash:    tx.Hash,
		Delegator: m.DelegatorAddress,
		Validator: m.ValidatorDstAddress,
		Amount:    amount,
		CreatedAt: tx.Timestamp,
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
	for _, log := range tx.Logs {
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

	id := makeHash(fmt.Sprintf("%s.%d", tx.Hash, index))
	d.delegatorRewards = append(d.delegatorRewards, dmodels.DelegatorReward{
		ID:        id,
		TxHash:    tx.Hash,
		Delegator: m.DelegatorAddress,
		Validator: m.ValidatorAddress,
		Amount:    amount,
		CreatedAt: tx.Timestamp,
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
	for _, log := range tx.Logs {
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
	amount, err := calculateAmount(m.Content.Value.Amount)
	if err != nil {
		return fmt.Errorf("calculateAmount: %s", err.Error())
	}
	initDeposit, err := calculateAmount(m.InitialDeposit)
	if err != nil {
		return fmt.Errorf("calculateAmount: %s", err.Error())
	}
	d.proposals = append(d.proposals, dmodels.HistoryProposal{
		ID:          id,
		TxHash:      tx.Hash,
		Title:       m.Content.Value.Title,
		Description: m.Content.Value.Description,
		Recipient:   m.Content.Value.Recipient,
		Amount:      amount,
		InitDeposit: initDeposit,
		Proposer:    m.Proposer,
		CreatedAt:   tx.Timestamp,
	})
	return nil
}

func (d *data) parseVoteMsg(index int, tx Tx, data []byte) (err error) {
	var m MsgVote
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	id := makeHash(fmt.Sprintf("%s.%d.s", tx.Hash, index))
	d.proposalVotes = append(d.proposalVotes, dmodels.ProposalVote{
		ID:         id,
		ProposalID: m.ProposalID,
		Voter:      m.Voter,
		TxHash:     tx.Hash,
		Option:     m.Option,
		CreatedAt:  dmodels.NewTime(tx.Timestamp),
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

	id := makeHash(fmt.Sprintf("%s.%d.s", tx.Hash, index))
	d.proposalDeposits = append(d.proposalDeposits, dmodels.ProposalDeposit{
		ID:         id,
		ProposalID: m.ProposalID,
		Depositor:  m.Depositor,
		Amount:     amount,
		CreatedAt:  dmodels.NewTime(tx.Timestamp),
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
	for _, log := range tx.Logs {
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
	id := makeHash(fmt.Sprintf("%s.%d", tx.Hash, index))
	d.validatorRewards = append(d.validatorRewards, dmodels.ValidatorReward{
		TxHash:    tx.Hash,
		ID:        id,
		Address:   m.ValidatorAddress,
		Amount:    amount,
		CreatedAt: tx.Timestamp,
	})
	return nil
}

func (d *data) parseUnjailMsg(index int, tx Tx, data []byte) (err error) {
	var m MsgUnjail
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	id := makeHash(fmt.Sprintf("%s.%d", tx.Hash, index))
	d.jailers = append(d.jailers, dmodels.Jailer{
		ID:        id,
		Address:   m.Address,
		CreatedAt: tx.Timestamp,
	})
	return nil
}

func calculateAmount(amountItems []Amount) (decimal.Decimal, error) {
	volume := decimal.Zero
	for _, item := range amountItems {
		if item.Denom == "" && item.Amount.IsZero() { // example height=1245781
			break
		}
		if item.Denom != "uatom" {
			return volume, fmt.Errorf("unknown demon (currency): %s", item.Denom)
		}
		volume = volume.Add(item.Amount)
	}
	volume = volume.Div(precisionDiv)
	return volume, nil
}

func (a Amount) getAmount() (decimal.Decimal, error) {
	if a.Denom == "" && a.Amount.IsZero() {
		return decimal.Zero, nil
	}
	if a.Denom != "uatom" {
		return decimal.Zero, fmt.Errorf("unknown demon (currency): %s", a.Denom)
	}
	a.Amount = a.Amount.Div(precisionDiv)
	return a.Amount, nil
}

func strToAmount(str string) (decimal.Decimal, error) {
	if str == "" {
		return decimal.Zero, nil
	}
	val := strings.TrimSuffix(str, "uatom")
	amount, err := decimal.NewFromString(val)
	if err != nil {
		return amount, fmt.Errorf("decimal.NewFromString: %s", err.Error())
	}
	amount = amount.Div(precisionDiv)
	return amount, nil
}

func makeHash(str string) string {
	hash := sha1.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

package hub3

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/everstake/cosmoscan-api/config"
	"github.com/everstake/cosmoscan-api/dao"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/shopspring/decimal"
	"math"
	"sync"
	"time"
)

const repeatDelay = time.Second * 5
const parserTitle = "hub3"

const taskNameBlock = "block"
const taskNameTxs = "txs"
const batchTxs = 50

type (
	Parser struct {
		cfg       config.Config
		api       api
		dao       dao.DAO
		fetcherCh chan task
		workerCh  chan task
		saverCh   chan data
		batchDone chan struct{}
		errCh     chan error
		data      data

		ctx  context.Context
		stop context.CancelFunc
		wg   *sync.WaitGroup
	}
	api interface {
		GetLatestBlock() (block Block, err error)
		GetBlock(height uint64) (block Block, err error)
		GetTxs(filter TxsFilter) (txs TxsBatch, err error)
	}
	data struct {
		blocks           []dmodels.Block
		transactions     []dmodels.Transaction
		transfers        []dmodels.Transfer
		delegations      []dmodels.Delegation
		delegatorRewards []dmodels.DelegatorReward
		validatorRewards []dmodels.ValidatorReward
		proposals        []dmodels.Proposal
		votes            []dmodels.ProposalVote
		proposalDeposits []dmodels.ProposalDeposit
	}
	task struct {
		name   string
		value  interface{}
		height uint64
		page   uint64
		batch  uint64
	}
)

func NewParser(cfg config.Config, d dao.DAO) *Parser {
	ctx, cancel := context.WithCancel(context.Background())
	return &Parser{
		cfg:       cfg,
		dao:       d,
		api:       NewAPI(cfg.Parser.Node),
		fetcherCh: make(chan task, 100000),
		workerCh:  make(chan task, 100000),
		saverCh:   make(chan data),
		errCh:     make(chan error),
		batchDone: make(chan struct{}),
		ctx:       ctx,
		stop:      cancel,
		wg:        &sync.WaitGroup{},
	}
}

func (p *Parser) Stop() error {
	p.stop()
	p.wg.Wait()
	return nil
}

func (p *Parser) Title() string {
	return "Parser"
}

func (p *Parser) Run() error {
	model, err := p.dao.GetParser(parserTitle)
	if err != nil {
		return fmt.Errorf("parser not found")
	}
	p.wg.Add(1)
	defer p.wg.Done()
	go p.runWorker()
	go p.runFetchers()
	go p.runSaver(model)
	for {
		latestBlock, err := p.api.GetLatestBlock()
		if err != nil {
			log.Error("Parser: api.GetLatestBlock: %s", err.Error())
			<-time.After(repeatDelay)
			continue
		}
		if model.Height >= latestBlock.Block.Header.Height {
			<-time.After(repeatDelay)
			continue
		}
		for model.Height < latestBlock.Block.Header.Height {
			batch := p.cfg.Parser.Batch
			if latestBlock.Block.Header.Height-model.Height < batch {
				batch = latestBlock.Block.Header.Height - model.Height
			}

			to := model.Height + batch
			for i := model.Height + 1; i <= to; i++ {
				p.fetcherCh <- task{
					name:   taskNameBlock,
					height: i,
					batch:  batch,
				}
			}
			select {
			case <-p.ctx.Done():
				return nil
			case <-p.batchDone:
				model.Height += batch
			case err := <-p.errCh:
				return err
			}
		}
	}
}

func (p *Parser) runFetchers() {
	for i := uint64(0); i < p.cfg.Parser.Fetchers; i++ {
		go func() {
			p.wg.Add(1)
			defer p.wg.Done()
			for {
				select {
				case <-p.ctx.Done():
					return
				case task := <-p.fetcherCh:
					switch task.name {
					case taskNameBlock:
						var err error
						for {
							task.value, err = p.api.GetBlock(task.height)
							if err == nil {
								p.workerCh <- task
								break
							}
							log.Error("Parser: fetcher: api.GetBlock: %s", err.Error())
							<-time.After(time.Second)
						}
					case taskNameTxs:
						var err error
						for {
							task.value, err = p.api.GetTxs(TxsFilter{
								Limit:  batchTxs,
								Page:   task.page,
								Height: task.height,
							})
							if err == nil {
								p.workerCh <- task
								break
							}
							log.Error("Parser: fetcher: api.GetTxs: %s", err.Error())
							<-time.After(time.Second)
						}
					}
				}
			}
		}()
	}
}

func (p *Parser) runWorker() {
	var batch uint64
	var totalTxs int
	var parsedTxs int
	for {
		select {
		case <-p.ctx.Done():
			return
		case t := <-p.workerCh:
			switch t.name {
			case taskNameBlock:
				block := t.value.(Block)
				batch = t.batch

				p.data.blocks = append(p.data.blocks, dmodels.Block{
					ID:        block.Block.Header.Height,
					Hash:      block.BlockMeta.BlockID.Hash,
					Proposer:  "test",
					CreatedAt: block.BlockMeta.Header.Time,
				})
				totalTxs += block.BlockMeta.Header.NumTxs

				if block.BlockMeta.Header.NumTxs != 0 {
					pages := int(math.Ceil(float64(block.BlockMeta.Header.NumTxs) / float64(batchTxs)))
					for page := 1; page <= pages; page++ {
						p.fetcherCh <- task{
							name:   taskNameTxs,
							height: t.height,
							page:   uint64(page),
							batch:  batchTxs,
						}
					}
				}

			case taskNameTxs:
				txs := t.value.(TxsBatch)
				for _, tx := range txs.Txs {
					fail := false
					for _, l := range tx.Logs {
						if !l.Success {
							fail = true
						}
					}

					fee, err := calculateAmount(tx.Tx.Value.Fee.Amount)
					if err != nil {
						p.errCh <- fmt.Errorf("height: %d, calculateAmount: %s", tx.Height, err.Error())
						return
					}

					p.data.transactions = append(p.data.transactions, dmodels.Transaction{
						Hash:      tx.Hash,
						Status:    !fail,
						Height:    tx.Height,
						Messages:  uint64(len(tx.Tx.Value.Msg)),
						Fee:       fee,
						GasUsed:   tx.Result.GasUsed,
						GasWanted: tx.Result.GasWanted,
						CreatedAt: tx.Timestamp,
					})

					for i, msg := range tx.Tx.Value.Msg {
						switch msg.Type {
						case SendMsg:
							err = p.parseMsgSend(i, tx, msg.Value)
						case MultiSendMsg:
							err = p.parseMultiSendMsg(i, tx, msg.Value)
						case DelegateMsg:
							err = p.parseDelegateMsg(i, tx, msg.Value)
						case UndelegateMsg:
							err = p.parseUndelegateMsg(i, tx, msg.Value)
						case BeginRedelegateMsg:
							err = p.parseBeginRedelegateMsg(i, tx, msg.Value)
						case WithdrawDelegationRewardMsg:
							err = p.parseWithdrawDelegationRewardMsg(i, tx, msg.Value)
						case WithdrawDelegationRewardsAllMsg:
							// todo
							fmt.Println(WithdrawDelegationRewardsAllMsg, tx.Height, string(msg.Value))
							//err = p.parseWithdrawDelegationRewardsAllMsg(i, tx, msg.Value)
						case WithdrawValidatorCommissionMsg:
							// todo
							fmt.Println(WithdrawDelegationRewardsAllMsg, tx.Height, string(msg.Value))
							//err = p.parseWithdrawValidatorCommissionMsg(i, tx, msg.Value)
						case SubmitProposalBaseMsg:
							err = p.parseSubmitProposalBaseMsg(i, tx, msg.Value)
						case DepositMsg:
							err = p.parseDepositMsg(i, tx, msg.Value)
						case VoteMsg:
							err = p.parseVoteMsg(i, tx, msg.Value)
						}
						if err != nil {
							p.errCh <- fmt.Errorf("%s, (height: %d): %s", msg.Type, tx.Height, err.Error())
							return
						}
					}

					parsedTxs++
				}
			}
		}
		if batch == uint64(len(p.data.blocks)) && parsedTxs == totalTxs {

			p.saverCh <- p.data
			p.data = data{}
			batch = 0
			totalTxs = 0
			parsedTxs = 0
			p.batchDone <- struct{}{}
		}
	}
}

func (p *Parser) parseMsgSend(index int, tx Tx, data []byte) (err error) {
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
	p.data.transfers = append(p.data.transfers, dmodels.Transfer{
		ID:        id,
		TxHash:    tx.Hash,
		From:      m.FromAddress,
		To:        m.ToAddress,
		Amount:    amount,
		CreatedAt: tx.Timestamp,
	})
	return nil
}

func (p *Parser) parseMultiSendMsg(index int, tx Tx, data []byte) (err error) {
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
		p.data.transfers = append(p.data.transfers, dmodels.Transfer{
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
		p.data.transfers = append(p.data.transfers, dmodels.Transfer{
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

func (p *Parser) parseDelegateMsg(index int, tx Tx, data []byte) (err error) {
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
	p.data.delegations = append(p.data.delegations, dmodels.Delegation{
		ID:        id,
		TxHash:    tx.Hash,
		Delegator: m.DelegatorAddress,
		Validator: m.ValidatorAddress,
		Amount:    amount,
		CreatedAt: tx.Timestamp,
	})
	return nil
}

func (p *Parser) parseUndelegateMsg(index int, tx Tx, data []byte) (err error) {
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
	p.data.delegations = append(p.data.delegations, dmodels.Delegation{
		ID:        id,
		TxHash:    tx.Hash,
		Delegator: m.DelegatorAddress,
		Validator: m.ValidatorAddress,
		Amount:    amount.Mul(decimal.NewFromFloat(-1)),
		CreatedAt: tx.Timestamp,
	})
	return nil
}

func (p *Parser) parseBeginRedelegateMsg(index int, tx Tx, data []byte) (err error) {
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
	p.data.delegations = append(p.data.delegations, dmodels.Delegation{
		ID:        id,
		TxHash:    tx.Hash,
		Delegator: m.DelegatorAddress,
		Validator: m.ValidatorSrcAddress,
		Amount:    amount.Mul(decimal.NewFromFloat(-1)),
		CreatedAt: tx.Timestamp,
	})
	id = makeHash(fmt.Sprintf("%s.%d.d", tx.Hash, index))
	p.data.delegations = append(p.data.delegations, dmodels.Delegation{
		ID:        id,
		TxHash:    tx.Hash,
		Delegator: m.DelegatorAddress,
		Validator: m.ValidatorDstAddress,
		Amount:    amount,
		CreatedAt: tx.Timestamp,
	})
	return nil
}

func (p *Parser) parseWithdrawDelegationRewardMsg(index int, tx Tx, data []byte) (err error) {
	var m MsgWithdrawDelegationReward
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	// todo amount5
	id := makeHash(fmt.Sprintf("%s.%d.s", tx.Hash, index))
	p.data.delegatorRewards = append(p.data.delegatorRewards, dmodels.DelegatorReward{
		ID:        id,
		TxHash:    "",
		Delegator: "",
		Validator: "",
		Amount:    decimal.Decimal{},
		CreatedAt: time.Time{},
	})
	return nil
}

func (p *Parser) parseSubmitProposalBaseMsg(index int, tx Tx, data []byte) (err error) {
	var m MsgSubmitProposalBase
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	amount, err := m.InitialDeposit.getAmount()
	if err != nil {
		return fmt.Errorf("InitialDeposit.getAmount: %s", err.Error())
	}
	// todo PROPOSAL_ID
	p.data.proposals = append(p.data.proposals, dmodels.Proposal{
		ID:          tx.Hash,
		InitDeposit: amount,
		Proposer:    m.Proposer,
		Content:     string(m.Content),
		CreatedAt:   tx.Timestamp,
	})
	return nil
}

func (p *Parser) parseVoteMsg(index int, tx Tx, data []byte) (err error) {
	var m MsgVote
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	id := makeHash(fmt.Sprintf("%s.%d.s", tx.Hash, index))
	p.data.votes = append(p.data.votes, dmodels.ProposalVote{
		ID:         id,
		ProposalID: m.ProposalID,
		Voter:      m.Voter,
		Option:     m.Option,
		CreatedAt:  tx.Timestamp,
	})
	return nil
}

func (p *Parser) parseDepositMsg(index int, tx Tx, data []byte) (err error) {
	var m MsgDeposit
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	amount, err := m.Amount.getAmount()
	if err != nil {
		return fmt.Errorf("getAmount: %s", err.Error())
	}
	id := makeHash(fmt.Sprintf("%s.%d.s", tx.Hash, index))
	p.data.proposalDeposits = append(p.data.proposalDeposits, dmodels.ProposalDeposit{
		ID:         id,
		ProposalID: m.ProposalID,
		Depositor:  m.Depositor,
		Amount:     amount,
		CreatedAt:  tx.Timestamp,
	})
	return nil
}

//func (p *Parser) parseWithdrawDelegationRewardsAllMsg(index int, tx Tx, data []byte) (err error) {
//	var m MsgWithdrawDelegationRewardsAll
//	err = json.Unmarshal(data, &m)
//	if err != nil {
//		return fmt.Errorf("json.Unmarshal: %s", err.Error())
//	}
//	// TODO
//	id := makeHash(fmt.Sprintf("%s.%d.s", tx.Hash, index))
//	p.data.delegatorRewards = append(p.data.delegatorRewards, dmodels.DelegatorReward{
//		ID:        "",
//		TxHash:    "",
//		Delegator: "",
//		Validator: "",
//		Amount:    decimal.Decimal{},
//		CreatedAt: time.Time{},
//	})
//	return nil
//}

//func (p *Parser) parseWithdrawValidatorCommissionMsg(index int, tx Tx, data []byte) (err error) {
//	var m MsgWithdrawValidatorCommission
//	err = json.Unmarshal(data, &m)
//	if err != nil {
//		return fmt.Errorf("json.Unmarshal: %s", err.Error())
//	}
//	// todo amount
//	id := makeHash(fmt.Sprintf("%s.%d", tx.Hash, index))
//	p.data.validatorRewards = append(p.data.validatorRewards, dmodels.ValidatorReward{
//		ID:        id,
//		Address:   m.ValidatorAddress,
//		Amount:    ?,
//		CreatedAt:  tx.Timestamp,
//	})
//	return nil
//}

func (p *Parser) runSaver(model dmodels.Parser) {
	height := model.Height
	for {
		data := <-p.saverCh
		var err error
		for {
			err = p.dao.CreateBlocks(data.blocks)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateBlocks: %s", err.Error())
			<-time.After(repeatDelay)
		}
		for {
			height += uint64(len(data.blocks))
			err = p.dao.UpdateParser(dmodels.Parser{
				ID:     model.ID,
				Title:  parserTitle,
				Height: height,
			})
			if err == nil {
				break
			}
			log.Error("Parser: dao.UpdateParser: %s", err.Error())
			<-time.After(repeatDelay)
		}
	}
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
	return volume, nil
}

func (a Amount) getAmount() (decimal.Decimal, error) {
	if a.Denom == "" && a.Amount.IsZero() {
		return decimal.Zero, nil
	}
	if a.Denom != "uatom" {
		return decimal.Zero, fmt.Errorf("unknown demon (currency): %s", a.Denom)
	}
	return a.Amount, nil
}

func makeHash(str string) string {
	hash := sha1.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

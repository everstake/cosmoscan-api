package hub3

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	SendMsg                        = "cosmos-sdk/MsgSend"
	MultiSendMsg                   = "cosmos-sdk/MsgMultiSend"
	DelegateMsg                    = "cosmos-sdk/MsgDelegate"
	UndelegateMsg                  = "cosmos-sdk/MsgUndelegate"
	BeginRedelegateMsg             = "cosmos-sdk/MsgBeginRedelegate"
	WithdrawDelegationRewardMsg    = "cosmos-sdk/MsgWithdrawDelegationReward"
	WithdrawValidatorCommissionMsg = "cosmos-sdk/MsgWithdrawValidatorCommission"
	SubmitProposalMsg              = "cosmos-sdk/MsgSubmitProposal"
	DepositMsg                     = "cosmos-sdk/MsgDeposit"
	VoteMsg                        = "cosmos-sdk/MsgVote"
	UnJailMsg                      = "cosmos-sdk/MsgUnjail"
)

type (
	API struct {
		address string
		client  *http.Client
	}

	Block struct {
		BlockID struct {
			Hash  string `json:"hash"`
			Parts struct {
				Total int    `json:"total"`
				Hash  string `json:"hash"`
			} `json:"parts"`
		} `json:"block_id"`
		Block struct {
			Header struct {
				Version struct {
					Block uint64 `json:"block,string"`
				} `json:"version"`
				ChainID     string    `json:"chain_id"`
				Height      uint64    `json:"height,string"`
				Time        time.Time `json:"time"`
				LastBlockID struct {
					Hash  string `json:"hash"`
					Parts struct {
						Total int    `json:"total"`
						Hash  string `json:"hash"`
					} `json:"parts"`
				} `json:"last_block_id"`
				ProposerAddress string `json:"proposer_address"`
			} `json:"header"`
			Data struct {
				Txs []string `json:"txs"`
			} `json:"data"`
			Evidence struct {
				Evidence []interface{} `json:"evidence"`
			} `json:"evidence"`
			LastCommit struct {
				Height  string `json:"height"`
				Round   int    `json:"round"`
				BlockID struct {
					Hash  string `json:"hash"`
					Parts struct {
						Total int    `json:"total"`
						Hash  string `json:"hash"`
					} `json:"parts"`
				} `json:"block_id"`
				Signatures []struct {
					ValidatorAddress string `json:"validator_address"`
				}
			} `json:"last_commit"`
		} `json:"block"`
	}

	TxsBatch struct {
		TotalCount int  `json:"total_count,string"`
		Count      int  `json:"count,string"`
		PageNumber int  `json:"page_number,string"`
		PageTotal  int  `json:"page_total,string"`
		Limit      int  `json:"limit,string"`
		Txs        []Tx `json:"txs"`
	}

	Tx struct {
		Height uint64 `json:"height,string"`
		Hash   string `json:"txhash"`
		Data   string `json:"data"`
		RawLog string `json:"raw_log"`
		Code   int64  `json:"code"`
		Logs   []struct {
			Events []struct {
				Type       string `json:"type"`
				Attributes []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"attributes"`
			} `json:"events"`
		} `json:"logs"`
		GasWanted uint64 `json:"gas_wanted,string"`
		GasUsed   uint64 `json:"gas_used,string"`
		Tx        struct {
			Type  string `json:"type"`
			Value struct {
				Msg []struct {
					Type  string          `json:"type"`
					Value json.RawMessage `json:"value"`
				} `json:"msg"`
				Fee struct {
					Amount []Amount `json:"amount"`
					Gas    uint64   `json:"gas,string"`
				} `json:"fee"`
				Memo string `json:"memo"`
			} `json:"value"`
		} `json:"tx"`
		Timestamp time.Time `json:"timestamp"`
	}
	Amount struct {
		Denom  string          `json:"denom"`
		Amount decimal.Decimal `json:"amount"`
	}

	MsgSend struct {
		FromAddress string   `json:"from_address,omitempty"`
		ToAddress   string   `json:"to_address,omitempty"`
		Amount      []Amount `json:"amount"`
	}
	MsgMultiSendValue struct {
		Inputs []struct {
			Address string   `json:"address"`
			Coins   []Amount `json:"coins"`
		} `json:"inputs"`
		Outputs []struct {
			Address string   `json:"address"`
			Coins   []Amount `json:"coins"`
		} `json:"outputs"`
	}
	MsgDelegate struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Amount           Amount `json:"amount"`
	}
	MsgUndelegate struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Amount           Amount `json:"amount"`
	}
	MsgBeginRedelegate struct {
		DelegatorAddress    string `json:"delegator_address"`
		ValidatorSrcAddress string `json:"validator_src_address"`
		ValidatorDstAddress string `json:"validator_dst_address"`
		Amount              Amount `json:"amount"`
	}
	MsgWithdrawDelegationReward struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
	}
	MsgWithdrawDelegationRewardsAll struct {
		DelegatorAddress string `json:"delegator_address"`
	}
	MsgWithdrawValidatorCommission struct {
		ValidatorAddress string `json:"validator_address"`
	}
	MsgSubmitProposal struct {
		Content struct {
			Type  string `json:"type"`
			Value struct {
				Title       string   `json:"title"`
				Description string   `json:"description"`
				Recipient   string   `json:"recipient"`
				Amount      []Amount `json:"amount"`
			} `json:"value"`
		} `json:"content"`
		InitialDeposit []Amount `json:"initial_deposit"`
		Proposer       string   `json:"proposer"`
	}
	MsgDeposit struct {
		ProposalID uint64   `json:"proposal_id,string"`
		Depositor  string   `json:"depositor" `
		Amount     []Amount `json:"amount" `
	}
	MsgVote struct {
		ProposalID uint64 `json:"proposal_id,string"`
		Voter      string `json:"voter"`
		Option     int    `json:"option"`
	}
	MsgUnjail struct {
		Address string `json:"address"`
	}

	TxsFilter struct {
		Limit     uint64
		Page      uint64
		Height    uint64
		MinHeight uint64
		MaxHeight uint64
	}

	Validatorsets struct {
		Result struct {
			Validators []struct {
				Address     string          `json:"address"`
				VotingPower decimal.Decimal `json:"voting_power"`
			} `json:"validators"`
		} `json:"result"`
	}
)

func NewAPI(address string) *API {
	return &API{
		address: address,
		client: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (api *API) GetBlock(height uint64) (block Block, err error) {
	endpoint := fmt.Sprintf("blocks/%d", height)
	err = api.get(endpoint, nil, &block)
	return block, err
}

func (api *API) GetLatestBlock() (block Block, err error) {
	err = api.get("blocks/latest", nil, &block)
	return block, err
}

func (api *API) GetTxs(filter TxsFilter) (txs TxsBatch, err error) {
	params := make(map[string]string)
	if filter.Limit != 0 {
		params["limit"] = fmt.Sprintf("%d", filter.Limit)
	}
	if filter.Page != 0 {
		params["page"] = fmt.Sprintf("%d", filter.Page)
	}
	if filter.MinHeight != 0 {
		params["tx.minheight"] = fmt.Sprintf("%d", filter.MinHeight)
	}
	if filter.MaxHeight != 0 {
		params["tx.maxheight"] = fmt.Sprintf("%d", filter.MaxHeight)
	}
	if filter.Height != 0 {
		params["tx.height"] = fmt.Sprintf("%d", filter.Height)
	}
	err = api.get("txs", params, &txs)
	return txs, err
}

func (api *API) get(endpoint string, params map[string]string, result interface{}) error {
	fullURL := fmt.Sprintf("%s/%s", api.address, endpoint)
	if len(params) != 0 {
		values := url.Values{}
		for key, value := range params {
			values.Add(key, value)
		}
		fullURL = fmt.Sprintf("%s?%s", fullURL, values.Encode())
	}
	resp, err := api.client.Get(fullURL)
	if err != nil {
		return fmt.Errorf("client.Get: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %s", err.Error())
	}
	err = json.Unmarshal(data, result)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	return nil
}

func (api *API) GetValidatorset(height uint64) (set Validatorsets, err error) {
	err = api.get(fmt.Sprintf("validatorsets/%d", height), nil, &set)
	return set, err
}

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
	SendMsg                        = "/cosmos.bank.v1beta1.MsgSend"
	MultiSendMsg                   = "/cosmos.bank.v1beta1.MsgMultiSend"
	DelegateMsg                    = "/cosmos.staking.v1beta1.MsgDelegate"
	UndelegateMsg                  = "/cosmos.staking.v1beta1.MsgUndelegate"
	BeginRedelegateMsg             = "/cosmos.staking.v1beta1.MsgBeginRedelegate"
	WithdrawDelegationRewardMsg    = "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward"
	WithdrawValidatorCommissionMsg = "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission"
	SubmitProposalMsg              = "/cosmos.gov.v1beta1.MsgSubmitProposal"
	DepositMsg                     = "/cosmos.gov.v1beta1.MsgDeposit"
	VoteMsg                        = "/cosmos.gov.v1beta1.MsgVote"
	UnJailMsg                      = "/cosmos.slashing.v1beta1.MsgUnjail"
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

	Tx struct {
		Tx struct {
			Body struct {
				Messages []json.RawMessage `json:"messages"`
				Memo     string            `json:"memo"`
			} `json:"body"`
			AuthInfo struct {
				SignerInfos []struct {
					PublicKey struct {
						Type string `json:"@type"`
						Key  string `json:"key"`
					} `json:"public_key"`
				} `json:"signer_infos"`
				Fee struct {
					Amount   []Amount `json:"amount"`
					GasLimit uint64   `json:"gas_limit,string"`
					Payer    string   `json:"payer"`
					Granter  string   `json:"granter"`
				} `json:"fee"`
				Signatures []string `json:"signatures"`
			} `json:"auth_info"`
		} `json:"tx"`
		TxResponse struct {
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
				Type string `json:"@type"`
				Body struct {
					Messages []json.RawMessage `json:"messages"`
					Memo     string            `json:"memo"`
				} `json:"body"`
			} `json:"tx"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"tx_response"`
	}

	BaseMsg struct {
		Type string `json:"@type"`
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
		Option     string `json:"option"`
	}
	MsgUnjail struct {
		ValidatorAddr string `json:"validator_addr"`
	}

	TxsFilter struct {
		Limit     uint64
		Page      uint64
		Height    uint64
		MinHeight uint64
		MaxHeight uint64
	}

	Validatorsets struct {
		Validators []struct {
			Address string `json:"address"`
			PubKey  struct {
				Type string `json:"@type"`
				Key  string `json:"key"`
			} `json:"pub_key"`
			VotingPower decimal.Decimal `json:"voting_power"`
		} `json:"validators"`
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

func (api *API) GetTx(hash string) (tx Tx, err error) {
	endpoint := fmt.Sprintf("cosmos/tx/v1beta1/txs/%s", hash)
	err = api.get(endpoint, nil, &tx)
	return tx, err
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
		d, _ := ioutil.ReadAll(resp.Body)
		text := string(d)
		if len(text) > 150 {
			text = text[:150]
		}
		return fmt.Errorf("bad status: %d, %s", resp.StatusCode, text)
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
	err = api.get(fmt.Sprintf("cosmos/base/tendermint/v1beta1/validatorsets/%d", height), nil, &set)
	return set, err
}

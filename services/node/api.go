package node

import (
	"encoding/json"
	"fmt"
	"github.com/everstake/cosmoscan-api/config"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	precision = 6

	DepositPeriodProposalStatus = "PROPOSAL_STATUS_DEPOSIT_PERIOD"
	VotingPeriodProposalStatus  = "PROPOSAL_STATUS_VOTING_PERIOD"
	PassedProposalStatus        = "PROPOSAL_STATUS_PASSED"
	RejectedProposalStatus      = "PROPOSAL_STATUS_REJECTED"
	FailedProposalStatus        = "PROPOSAL_STATUS_FAILED"

	MainUnit = "uatom"
)

var PrecisionDiv = decimal.New(1, precision)

type (
	API struct {
		cfg    config.Config
		client *http.Client
	}
	CommunityPool struct {
		Pool []struct {
			Denom  string          `json:"denom"`
			Amount decimal.Decimal `json:"amount"`
		} `json:"pool"`
	}
	Validators struct {
		Validators []Validator `json:"validators"`
	}
	Validator struct {
		OperatorAddress string `json:"operator_address"`
		ConsensusPubkey struct {
			Type string `json:"@type"`
			Key  string `json:"key"`
		} `json:"consensus_pubkey"`
		Tokens          uint64          `json:"tokens,string"`
		DelegatorShares decimal.Decimal `json:"delegator_shares"`
		Description     struct {
			Moniker  string `json:"moniker"`
			Identity string `json:"identity"`
			Website  string `json:"website"`
			Details  string `json:"details"`
		} `json:"description"`
		UnbondingHeight uint64    `json:"unbonding_height,string"`
		UnbondingTime   time.Time `json:"unbonding_time"`
		Commission      struct {
			CommissionRates struct {
				Rate          decimal.Decimal `json:"rate"`
				MaxRate       decimal.Decimal `json:"max_rate"`
				MaxChangeRate decimal.Decimal `json:"max_change_rate"`
			} `json:"commission_rates"`
		} `json:"commission"`
		MaxChangeRate decimal.Decimal `json:"max_change_rate"`
	}
	Inflation struct {
		Inflation decimal.Decimal `json:"inflation"`
	}
	AmountResult struct {
		Balances []struct {
			Denom  string          `json:"denom"`
			Amount decimal.Decimal `json:"amount"`
		} `json:"balances"`
	}
	DelegatorRewards struct {
		Rewards []struct {
			ValidatorAddress string   `json:"validator_address"`
			Reward           []Amount `json:"reward"`
		} `json:"rewards"`
		Total []Amount `json:"total"`
	}
	StakingPool struct {
		Pool struct {
			NotBondedTokens decimal.Decimal `json:"not_bonded_tokens"`
			BondedTokens    decimal.Decimal `json:"bonded_tokens"`
		} `json:"pool"`
	}
	Supply struct {
		Amount struct {
			Denom  string          `json:"denom"`
			Amount decimal.Decimal `json:"amount"`
		} `json:"amount"`
	}
	StakeResult struct {
		DelegationResponses []struct {
			Delegation struct {
				DelegatorAddress string          `json:"delegator_address"`
				ValidatorAddress string          `json:"validator_address"`
				Shares           decimal.Decimal `json:"shares"`
			} `json:"delegation"`
		} `json:"delegation_responses"`
	}
	UnbondingResult struct {
		UnbondingResponses []struct {
			DelegatorAddress string `json:"delegator_address"`
			ValidatorAddress string `json:"validator_address"`
			Entries          []struct {
				Balance decimal.Decimal `json:"balance"`
			} `json:"entries"`
		} `json:"unbonding_responses"`
	}
	ProposalsResult struct {
		Proposals []struct {
			Content struct {
				Type        string `json:"@type"`
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"content"`
			ProposalID       uint64 `json:"proposal_id,string"`
			Status           string `json:"status"`
			FinalTallyResult struct {
				Yes        int64 `json:"yes,string"`
				Abstain    int64 `json:"abstain,string"`
				No         int64 `json:"no,string"`
				NoWithVeto int64 `json:"no_with_veto,string"`
			} `json:"final_tally_result"`
			SubmitTime     time.Time `json:"submit_time"`
			DepositEndTime time.Time `json:"deposit_end_time"`
			TotalDeposit   []struct {
				Amount decimal.Decimal `json:"amount"`
			} `json:"total_deposit"`
			VotingStartTime time.Time `json:"voting_start_time"`
			VotingEndTime   time.Time `json:"voting_end_time"`
		} `json:"proposals"`
	}
	ProposalProposer struct {
		Proposal struct {
			ProposalID uint64 `json:"proposal_id,string"`
			Proposer   string `json:"proposer"`
		} `json:"proposal"`
	}
	DelegatorValidatorStakeResult struct {
		DelegationResponse struct {
			Delegation struct {
				DelegatorAddress string          `json:"delegator_address"`
				ValidatorAddress string          `json:"validator_address"`
				Shares           decimal.Decimal `json:"shares"`
			} `json:"delegation"`
			Balance struct {
				Denom  string          `json:"denom"`
				Amount decimal.Decimal `json:"amount"`
			} `json:"balance"`
		} `json:"delegation_response"`
	}
	ProposalVotersResult struct {
		Result []struct {
			ProposalID uint64 `json:"proposal_id,string"`
			Voter      string `json:"voter"`
			Option     string `json:"option"`
		} `json:"result"`
	}
	ProposalTallyResult struct {
		Tally struct {
			Yes        int64 `json:"yes,string"`
			Abstain    int64 `json:"abstain,string"`
			No         int64 `json:"no,string"`
			NoWithVeto int64 `json:"no_with_veto,string"`
		} `json:"tally"`
	}
	Block struct {
		BlockID struct {
			Hash          string `json:"hash"`
			PartSetHeader struct {
				Total int    `json:"total"`
				Hash  string `json:"hash"`
			} `json:"part_set_header"`
		} `json:"block_id"`
		Block struct {
			Header struct {
				Version struct {
					Block string `json:"block"`
					App   string `json:"app"`
				} `json:"version"`
				ChainID     string    `json:"chain_id"`
				Height      uint64    `json:"height,string"`
				Time        time.Time `json:"time"`
				LastBlockID struct {
					Hash          string `json:"hash"`
					PartSetHeader struct {
						Total int    `json:"total"`
						Hash  string `json:"hash"`
					} `json:"part_set_header"`
				} `json:"last_block_id"`
				LastCommitHash     string `json:"last_commit_hash"`
				DataHash           string `json:"data_hash"`
				ValidatorsHash     string `json:"validators_hash"`
				NextValidatorsHash string `json:"next_validators_hash"`
				ConsensusHash      string `json:"consensus_hash"`
				AppHash            string `json:"app_hash"`
				LastResultsHash    string `json:"last_results_hash"`
				EvidenceHash       string `json:"evidence_hash"`
				ProposerAddress    string `json:"proposer_address"`
			} `json:"header"`
			Data struct {
				Txs []string `json:"txs"`
			} `json:"data"`
		} `json:"block"`
	}
	TxResult struct {
		Tx struct {
			Type string `json:"@type"`
			Body struct {
				Messages         []json.RawMessage `json:"messages"`
				Memo             string            `json:"memo"`
				TimeoutHeight    string            `json:"timeout_height"`
				ExtensionOptions []struct {
					TypeURL string `json:"type_url"`
					Value   string `json:"value"`
				} `json:"extension_options"`
				NonCriticalExtensionOptions []struct {
					TypeURL string `json:"type_url"`
					Value   string `json:"value"`
				} `json:"non_critical_extension_options"`
			} `json:"body"`
			AuthInfo struct {
				SignerInfos []struct {
					PublicKey struct {
						TypeURL string `json:"type_url"`
						Value   string `json:"value"`
					} `json:"public_key"`
					ModeInfo struct {
						Single struct {
							Mode string `json:"mode"`
						} `json:"single"`
						Multi struct {
							Bitarray struct {
								ExtraBitsStored int    `json:"extra_bits_stored"`
								Elems           string `json:"elems"`
							} `json:"bitarray"`
							ModeInfos []interface{} `json:"mode_infos"`
						} `json:"multi"`
					} `json:"mode_info"`
					Sequence string `json:"sequence"`
				} `json:"signer_infos"`
				Fee struct {
					Amount []struct {
						Denom  string          `json:"denom"`
						Amount decimal.Decimal `json:"amount"`
					} `json:"amount"`
					GasLimit decimal.Decimal `json:"gas_limit"`
					Payer    string          `json:"payer"`
					Granter  string          `json:"granter"`
				} `json:"fee"`
			} `json:"auth_info"`
			Signatures []string `json:"signatures"`
		} `json:"tx"`
		TxResponse struct {
			Height    uint64 `json:"height,string"`
			Txhash    string `json:"txhash"`
			Codespace string `json:"codespace"`
			Code      int    `json:"code"`
			Data      string `json:"data"`
			RawLog    string `json:"raw_log"`
			Logs      []struct {
				MsgIndex int    `json:"msg_index"`
				Log      string `json:"log"`
				Events   []struct {
					Type       string `json:"type"`
					Attributes []struct {
						Key   string `json:"key"`
						Value string `json:"value"`
					} `json:"attributes"`
				} `json:"events"`
			} `json:"logs"`
			Info      string `json:"info"`
			GasWanted uint64 `json:"gas_wanted,string"`
			GasUsed   uint64 `json:"gas_used,string"`
			Tx        struct {
				TypeURL string `json:"type_url"`
				Value   string `json:"value"`
			} `json:"tx"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"tx_response"`
	}
	Amount struct {
		Denom  string          `json:"denom"`
		Amount decimal.Decimal `json:"amount"`
	}
)

func NewAPI(cfg config.Config) *API {
	return &API{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (api API) request(endpoint string, data interface{}) error {
	url := fmt.Sprintf("%s/%s", api.cfg.Parser.Node, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %s", err.Error())
	}
	resp, err := api.client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %s", err.Error())
	}
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %s", err.Error())
	}
	err = json.Unmarshal(d, data)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	return nil
}

func (api API) GetCommunityPoolAmount() (amount decimal.Decimal, err error) {
	var cp CommunityPool
	err = api.request("cosmos/distribution/v1beta1/community_pool", &cp)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	for _, p := range cp.Pool {
		if p.Denom == MainUnit {
			amount = amount.Add(p.Amount)
		}
	}
	return amount.Div(PrecisionDiv), nil
}

func (api API) GetValidators() (items []Validator, err error) {
	var validators Validators
	err = api.request("cosmos/staking/v1beta1/validators?pagination.limit=10000", &validators)
	if err != nil {
		return nil, fmt.Errorf("request: %s", err.Error())
	}
	return validators.Validators, nil
}

func (api API) GetInflation() (amount decimal.Decimal, err error) {
	var inflation Inflation
	err = api.request("cosmos/mint/v1beta1/inflation", &inflation)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	return inflation.Inflation.Mul(decimal.New(100, 0)), nil
}

func (api API) GetTotalSupply() (amount decimal.Decimal, err error) {
	var s Supply
	err = api.request("cosmos/bank/v1beta1/supply/uatom", &s)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	return s.Amount.Amount.Div(PrecisionDiv), nil
}

func (api API) GetStakingPool() (sp StakingPool, err error) {
	err = api.request("cosmos/staking/v1beta1/pool", &sp)
	if err != nil {
		return sp, fmt.Errorf("request: %s", err.Error())
	}
	sp.Pool.BondedTokens = sp.Pool.BondedTokens.Div(PrecisionDiv)
	sp.Pool.NotBondedTokens = sp.Pool.NotBondedTokens.Div(PrecisionDiv)
	return sp, nil
}

func (api API) GetBalance(address string) (amount decimal.Decimal, err error) {
	var result AmountResult
	err = api.request(fmt.Sprintf("cosmos/bank/v1beta1/balances/%s", address), &result)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	for _, b := range result.Balances {
		if b.Denom == MainUnit {
			amount = amount.Add(b.Amount)
		}
	}
	return amount.Div(PrecisionDiv), nil
}

func (api API) GetBalances(address string) (result AmountResult, err error) {
	err = api.request(fmt.Sprintf("cosmos/bank/v1beta1/balances/%s", address), &result)
	if err != nil {
		return result, fmt.Errorf("request: %s", err.Error())
	}
	return result, nil
}

func (api API) GetStakeRewards(address string) (amount decimal.Decimal, err error) {
	var result DelegatorRewards
	err = api.request(fmt.Sprintf("cosmos/distribution/v1beta1/delegators/%s/rewards", address), &result)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	for _, b := range result.Total {
		if b.Denom == MainUnit {
			amount = amount.Add(b.Amount)
		}
	}
	return amount.Div(PrecisionDiv), nil
}

func (api API) GetStake(address string) (amount decimal.Decimal, err error) {
	var result StakeResult
	err = api.request(fmt.Sprintf("cosmos/staking/v1beta1/delegations/%s?pagination.limit=10000", address), &result)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	shares := decimal.Zero
	for _, r := range result.DelegationResponses {
		shares = shares.Add(r.Delegation.Shares)
	}
	return shares.Div(PrecisionDiv), nil
}

func (api API) GetUnbonding(address string) (amount decimal.Decimal, err error) {
	var result UnbondingResult
	err = api.request(fmt.Sprintf("cosmos/staking/v1beta1/delegators/%s/unbonding_delegations?pagination.limit=10000", address), &result)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	for _, r := range result.UnbondingResponses {
		for _, entry := range r.Entries {
			amount = amount.Add(entry.Balance)
		}
	}
	amount = amount.Div(PrecisionDiv)
	return amount, nil
}

func (api API) GetProposals() (proposals ProposalsResult, err error) {
	err = api.request("cosmos/gov/v1beta1/proposals?pagination.limit=10000", &proposals)
	if err != nil {
		return proposals, fmt.Errorf("request: %s", err.Error())
	}
	return proposals, nil
}

func (api API) GetDelegatorValidatorStake(delegator string, validator string) (amount decimal.Decimal, err error) {
	var result DelegatorValidatorStakeResult
	err = api.request(fmt.Sprintf("cosmos/staking/v1beta1/validators/%s/delegations/%s", validator, delegator), &result)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	return result.DelegationResponse.Delegation.Shares.Div(PrecisionDiv), nil
}

func (api API) ProposalTallyResult(id uint64) (result ProposalTallyResult, err error) {
	err = api.request(fmt.Sprintf("/cosmos/gov/v1beta1/proposals/%d/tally", id), &result)
	if err != nil {
		return result, fmt.Errorf("request: %s", err.Error())
	}
	return result, nil
}

func (api API) GetBlock(id uint64) (result Block, err error) {
	err = api.request(fmt.Sprintf("/cosmos/base/tendermint/v1beta1/blocks/%d", id), &result)
	if err != nil {
		return result, fmt.Errorf("request: %s", err.Error())
	}
	return result, nil
}

func (api API) GetTransaction(hash string) (result TxResult, err error) {
	err = api.request(fmt.Sprintf("/cosmos/tx/v1beta1/txs/%s", hash), &result)
	if err != nil {
		return result, fmt.Errorf("request: %s", err.Error())
	}
	return result, nil
}

func Precision(amount decimal.Decimal) decimal.Decimal {
	return amount.Div(PrecisionDiv)
}

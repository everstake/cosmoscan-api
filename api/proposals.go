package api

import (
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/log"
	"net/http"
)

func (api *API) GetProposals(w http.ResponseWriter, r *http.Request) {
	var filter filters.Proposals
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API Decode: %s", err.Error())
		jsonBadRequest(w, "")
		return
	}
	resp, err := api.svc.GetProposals(filter)
	if err != nil {
		log.Error("API GetProposals: svc.GetProposals: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetProposalVotes(w http.ResponseWriter, r *http.Request) {
	var filter filters.ProposalVotes
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API Decode: %s", err.Error())
		jsonBadRequest(w, "")
		return
	}
	resp, err := api.svc.GetProposalVotes(filter)
	if err != nil {
		log.Error("API GetProposalVotes: svc.GetProposalVotes: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetProposalDeposits(w http.ResponseWriter, r *http.Request) {
	var filter filters.ProposalDeposits
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API Decode: %s", err.Error())
		jsonBadRequest(w, "")
		return
	}
	resp, err := api.svc.GetProposalDeposits(filter)
	if err != nil {
		log.Error("API GetProposalDeposits: svc.GetProposalDeposits: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetProposalChartData(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetProposalsChartData()
	if err != nil {
		log.Error("API GetProposalsChartData: svc.GetProposalsChartData: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

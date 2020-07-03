package api

import (
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) GetAggDelegationsVolume(w http.ResponseWriter, r *http.Request) {
	var filter filters.DelegationsAgg
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API Decode: %s", err.Error())
		jsonBadRequest(w, "")
		return
	}
	resp, err := api.svc.GetAggDelegationsVolume(filter)
	if err != nil {
		log.Error("API GetAggDelegationsVolume: svc.GetAggDelegationsVolume: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetAggUndelegationsVolume(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggUndelegationsVolume)
}

func (api *API) GetStakingPie(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetStakingPie()
	if err != nil {
		log.Error("API GetStakingPie: svc.GetStakingPie: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetValidatorDelegationsAgg(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["address"]
	if !ok || address == "" {
		jsonBadRequest(w, "invalid address")
		return
	}
	resp, err := api.svc.GetValidatorDelegationsAgg(address)
	if err != nil {
		log.Error("API GetValidatorDelegationsAgg: svc.GetValidatorDelegationsAgg: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetValidatorDelegatorsAgg(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["address"]
	if !ok || address == "" {
		jsonBadRequest(w, "invalid address")
		return
	}
	resp, err := api.svc.GetValidatorDelegatorsAgg(address)
	if err != nil {
		log.Error("API GetValidatorDelegatorsAgg: svc.GetValidatorDelegatorsAgg: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetValidatorDelegators(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["address"]
	if !ok || address == "" {
		jsonBadRequest(w, "invalid address")
		return
	}
	var filter filters.ValidatorDelegators
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API Decode: %s", err.Error())
		jsonBadRequest(w, "")
		return
	}
	if filter.Limit > 20 {
		filter.Limit = 20
	}
	filter.Validator = address
	resp, err := api.svc.GetValidatorDelegators(filter)
	if err != nil {
		log.Error("API GetValidatorDelegators: svc.GetValidatorDelegators: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

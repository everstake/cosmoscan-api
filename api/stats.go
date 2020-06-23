package api

import (
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/log"
	"net/http"
)

func (api *API) GetNetworkStats(w http.ResponseWriter, r *http.Request) {
	var filter filters.Stats
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API Decode: %s", err.Error())
		jsonBadRequest(w, "")
		return
	}
	resp, err := api.svc.GetNetworkStates(filter)
	if err != nil {
		log.Error("API GetNetworkStats: svc.GetNetworkStates: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetAggValidators33Power(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggValidators33Power)
}

package api

import (
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (api *API) GetAggBlocksCount(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggBlocksCount)
}

func (api *API) GetAggBlocksDelay(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggBlocksDelay)
}

func (api *API) GetAggUniqBlockValidators(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggUniqBlockValidators)
}

func (api *API) GetBlock(w http.ResponseWriter, r *http.Request) {
	heightStr, ok := mux.Vars(r)["height"]
	if !ok || heightStr == "" {
		jsonBadRequest(w, "invalid address")
		return
	}
	height, err := strconv.ParseUint(heightStr, 10, 64)
	if err != nil {
		jsonBadRequest(w, "invalid height")
		return
	}
	resp, err := api.svc.GetBlock(height)
	if err != nil {
		log.Error("API GetValidator: svc.GetBlock: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetBlocks(w http.ResponseWriter, r *http.Request) {
	var filter filters.Blocks
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API Decode: %s", err.Error())
		jsonBadRequest(w, "")
		return
	}
	if filter.Limit == 0 || filter.Limit > 100 {
		filter.Limit = 100
	}
	resp, err := api.svc.GetBlocks(filter)
	if err != nil {
		log.Error("API GetBlocks: svc.GetBlocks: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

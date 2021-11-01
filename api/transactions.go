package api

import (
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) GetAggTransactionsFee(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggTransactionsFee)
}

func (api *API) GetAggOperationsCount(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggOperationsCount)
}

func (api *API) GetAvgOperationsPerBlock(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAvgOperationsPerBlock)
}

func (api *API) GetTransaction(w http.ResponseWriter, r *http.Request) {
	hash, ok := mux.Vars(r)["hash"]
	if !ok || hash == "" {
		jsonBadRequest(w, "invalid hash")
		return
	}
	resp, err := api.svc.GetTransaction(hash)
	if err != nil {
		log.Error("API GetTransaction: svc.GetTransaction: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetTransactions(w http.ResponseWriter, r *http.Request) {
	var filter filters.Transactions
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API Decode: %s", err.Error())
		jsonBadRequest(w, "")
		return
	}
	if filter.Limit == 0 || filter.Limit > 100 {
		filter.Limit = 100
	}
	resp, err := api.svc.GetTransactions(filter)
	if err != nil {
		log.Error("API GetTransactions: svc.GetTransactions: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

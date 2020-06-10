package api

import (
	"net/http"
)

func (api *API) GetAggTransactionsFee(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggTransactionsFee)
}

func (api *API) GetAggOperationsCount(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggOperationsCount)
}

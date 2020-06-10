package api

import (
	"net/http"
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

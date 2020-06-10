package api

import "net/http"

func (api *API) GetAggDelegationsVolume(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggDelegationsVolume)
}

func (api *API) GetAggUndelegationsVolume(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggUndelegationsVolume)
}

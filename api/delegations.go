package api

import (
	"github.com/everstake/cosmoscan-api/log"
	"net/http"
)

func (api *API) GetAggDelegationsVolume(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggDelegationsVolume)
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

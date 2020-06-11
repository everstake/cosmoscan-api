package api

import (
	"github.com/everstake/cosmoscan-api/log"
	"net/http"
)

func (api *API) GetNetworkStats(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetNetworkStates()
	if err != nil {
		log.Error("API GetNetworkStats: svc.GetNetworkStates: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

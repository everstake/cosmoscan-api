package api

import (
	"github.com/everstake/cosmoscan-api/log"
	"net/http"
)

func (api *API) GetHistoricalState(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetHistoricalState()
	if err != nil {
		log.Error("API GetHistoricalState: svc.GetHistoricalState: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

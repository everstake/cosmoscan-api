package api

import (
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/log"
	"net/http"
)

func (api *API) GetProposals(w http.ResponseWriter, r *http.Request) {
	var filter filters.Proposals
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API Decode: %s", err.Error())
		jsonBadRequest(w, "")
		return
	}
	resp, err := api.svc.GetProposals(filter)
	if err != nil {
		log.Error("API GetProposals: svc.GetProposals: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

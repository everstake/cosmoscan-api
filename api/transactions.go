package api

import (
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/log"
	"net/http"
)

func (api *API) GetAggTransactionsFee(w http.ResponseWriter, r *http.Request) {
	var filter filters.Agg
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API GetAggTransactionsFee: Decode: %s", err.Error())
		jsonBadRequest(w, "")
		return
	}
	err = filter.Validate()
	if err != nil {
		log.Debug("API GetAggTransactionsFee: Validate: %s", err.Error())
		jsonBadRequest(w, err.Error())
		return
	}
	resp, err := api.svc.GetAggTransactionsFee(filter)
	if err != nil {
		log.Error("API GetAggTransactionsFee: svc.GetAggTransactionsFee: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}



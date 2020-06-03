package api

import (
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/log"
	"net/http"
)

func (api *API) GetAggTransfersVolume(w http.ResponseWriter, r *http.Request) {
	var filter filters.Agg
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API GetAggTransfersVolume: Decode: %s", err.Error())
		jsonBadRequest(w, "")
		return
	}
	err = filter.Validate()
	if err != nil {
		log.Debug("API GetAggTransfersVolume: Validate: %s", err.Error())
		jsonBadRequest(w, err.Error())
		return
	}
	resp, err := api.svc.GetAggTransfersVolume(filter)
	if err != nil {
		log.Error("API GetAggTransfersVolume: svc.GetAggTransfersVolume: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

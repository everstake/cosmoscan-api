package api

import (
	"github.com/everstake/cosmoscan-api/log"
	"net/http"
)

func (api *API) GetMetaData(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetMetaData()
	if err != nil {
		log.Error("API GetMetaData: svc.GetMetaData: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

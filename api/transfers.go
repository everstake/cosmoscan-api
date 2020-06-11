package api

import (
	"net/http"
)

func (api *API) GetAggTransfersVolume(w http.ResponseWriter, r *http.Request) {
	api.aggHandler(w, r, api.svc.GetAggTransfersVolume)
}

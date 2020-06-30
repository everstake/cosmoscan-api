package api

import (
	"github.com/everstake/cosmoscan-api/log"
	"net/http"
)

func (api *API) GetTopProposedBlocksValidators(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetTopProposedBlocksValidators()
	if err != nil {
		log.Error("API GetTopProposedBlocksValidators: svc.GetTopProposedBlocksValidators: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)

}

func (api *API) GetMostJailedValidators(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetMostJailedValidators()
	if err != nil {
		log.Error("API GetMostJailedValidators: svc.GetMostJailedValidators: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)

}

func (api *API) GetFeeRanges(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetFeeRanges()
	if err != nil {
		log.Error("API GetFeeRanges: svc.GetFeeRanges: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)

}

func (api *API) GetValidators(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetValidators()
	if err != nil {
		log.Error("API GetValidators: svc.GetValidators: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)

}

func (api *API) GetValidatorsDelegatorsTotal(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetValidatorsDelegatorsTotal()
	if err != nil {
		log.Error("API GetValidatorsDelegatorsTotal: svc.GetValidatorsDelegatorsTotal: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)

}

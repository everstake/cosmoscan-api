package api

import (
	"github.com/everstake/cosmoscan-api/config"
	"net/http"
)

func (api *API) Index(w http.ResponseWriter, r *http.Request) {
	jsonData(w, map[string]string{
		"service": config.ServiceName,
	})
}

func (api *API) Health(w http.ResponseWriter, r *http.Request) {
	jsonData(w, map[string]bool{
		"status": true,
	})
}

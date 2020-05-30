package api

import (
	"encoding/json"
	"fmt"
	"github.com/everstake/cosmoscan-api/config"
	"github.com/everstake/cosmoscan-api/dao"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/everstake/cosmoscan-api/services"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"net/http"
)

type API struct {
	dao          dao.DAO
	cfg          config.Config
	router       *mux.Router
	queryDecoder *schema.Decoder
}

type errResponse struct {
	Error string `json:"error"`
	Value string `json:"value"`
}

func NewAPI(cfg config.Config, svc services.Services, dao dao.DAO) *API {
	sd := schema.NewDecoder()
	sd.IgnoreUnknownKeys(true)
	//sd.RegisterConverter(dmodels.Time{}, func(s string) reflect.Value {
	//	timestamp, err := strconv.ParseInt(s, 10, 64)
	//	if err != nil {
	//		return reflect.Value{}
	//	}
	//	t := dmodels.NewTime(time.Unix(timestamp, 0))
	//	return reflect.ValueOf(t)
	//})
	return &API{
		cfg:          cfg,
		dao:          dao,
		queryDecoder: sd,
	}
}

func (api *API) Title() string {
	return "API"
}

func (api *API) Run() error {
	api.router = mux.NewRouter()
	api.loadRoutes()

	http.Handle("/", api.router)
	log.Info("Listen API server on %s port", api.cfg.API.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", api.cfg.API.Port), nil)
	if err != nil {
		return err
	}
	return nil
}

func (api *API) Stop() error {
	return nil
}

func (api *API) loadRoutes() {

	api.router = mux.NewRouter()

	wrapper := negroni.New()

	wrapper.Use(cors.New(cors.Options{
		AllowedOrigins:   api.cfg.API.AllowedHosts,
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "X-User-Env", "Sec-Fetch-Mode"},
	}))

	// public
	HandleActions(api.router, wrapper, "", []*Route{

	})

}

func jsonData(writer http.ResponseWriter, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("can`t marshal json"))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(bytes)
}

//func jsonError(writer http.ResponseWriter, err error) {
//	sErr, ok := err.(tp.Err)
//	if ok {
//		jsonData(writer, errResponse{
//			Error: sErr.Code,
//			Value: sErr.Value,
//		})
//		return
//	}
//	jsonData(writer, errResponse{
//		Error: serrors.ErrService,
//	})
//}

package api

import (
	"encoding/json"
	"fmt"
	"github.com/everstake/cosmoscan-api/config"
	"github.com/everstake/cosmoscan-api/dao"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/everstake/cosmoscan-api/services"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type API struct {
	dao          dao.DAO
	cfg          config.Config
	svc          services.Services
	router       *mux.Router
	queryDecoder *schema.Decoder
}

type errResponse struct {
	Error string `json:"error"`
	Msg   string `json:"msg"`
}

func NewAPI(cfg config.Config, svc services.Services, dao dao.DAO) *API {
	sd := schema.NewDecoder()
	sd.IgnoreUnknownKeys(true)
	sd.RegisterConverter(dmodels.Time{}, func(s string) reflect.Value {
		timestamp, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return reflect.Value{}
		}
		t := dmodels.NewTime(time.Unix(timestamp, 0))
		return reflect.ValueOf(t)
	})
	return &API{
		cfg:          cfg,
		dao:          dao,
		svc:          svc,
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

	api.router.
		PathPrefix("/static").
		Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./resources/static"))))

	wrapper := negroni.New()

	wrapper.Use(cors.New(cors.Options{
		AllowedOrigins:   api.cfg.API.AllowedHosts,
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "X-User-Env", "Sec-Fetch-Mode"},
	}))

	// public
	HandleActions(api.router, wrapper, "", []*Route{
		{Path: "/", Method: http.MethodGet, Func: api.Index},
		{Path: "/health", Method: http.MethodGet, Func: api.Health},
		{Path: "/api", Method: http.MethodGet, Func: api.GetSwaggerAPI},

		{Path: "/meta", Method: http.MethodGet, Func: api.GetMetaData},
		{Path: "/historical-state", Method: http.MethodGet, Func: api.GetHistoricalState},
		{Path: "/transactions/fee/agg", Method: http.MethodGet, Func: api.GetAggTransactionsFee},
		{Path: "/transfers/volume/agg", Method: http.MethodGet, Func: api.GetAggTransfersVolume},
		{Path: "/operations/count/agg", Method: http.MethodGet, Func: api.GetAggOperationsCount},
		{Path: "/blocks/count/agg", Method: http.MethodGet, Func: api.GetAggBlocksCount},
		{Path: "/blocks/delay/agg", Method: http.MethodGet, Func: api.GetAggBlocksDelay},
		{Path: "/blocks/validators/uniq/agg", Method: http.MethodGet, Func: api.GetAggUniqBlockValidators},
		{Path: "/blocks/operations/agg", Method: http.MethodGet, Func: api.GetAvgOperationsPerBlock},
		{Path: "/delegations/volume/agg", Method: http.MethodGet, Func: api.GetAggDelegationsVolume},
		{Path: "/undelegations/volume/agg", Method: http.MethodGet, Func: api.GetAggUndelegationsVolume},
		{Path: "/network/stats", Method: http.MethodGet, Func: api.GetNetworkStats},
		{Path: "/staking/pie", Method: http.MethodGet, Func: api.GetStakingPie},
		{Path: "/proposals", Method: http.MethodGet, Func: api.GetProposals},
		{Path: "/proposals/votes", Method: http.MethodGet, Func: api.GetProposalVotes},
		{Path: "/proposals/deposits", Method: http.MethodGet, Func: api.GetProposalDeposits},
		{Path: "/proposals/chart", Method: http.MethodGet, Func: api.GetProposalChartData},
		{Path: "/validators/33power/agg", Method: http.MethodGet, Func: api.GetAggValidators33Power},
		{Path: "/validators/top/proposed", Method: http.MethodGet, Func: api.GetTopProposedBlocksValidators},
		{Path: "/validators/top/jailed", Method: http.MethodGet, Func: api.GetMostJailedValidators},
		{Path: "/validators/fee/ranges", Method: http.MethodGet, Func: api.GetFeeRanges},
		{Path: "/accounts/whale/agg", Method: http.MethodGet, Func: api.GetAggWhaleAccounts},
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

func jsonError(writer http.ResponseWriter) {
	writer.WriteHeader(500)
	bytes, err := json.Marshal(errResponse{
		Error: "service_error",
	})
	if err != nil {
		writer.Write([]byte("can`t marshal json"))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(bytes)
}

func jsonBadRequest(writer http.ResponseWriter, msg string) {
	bytes, err := json.Marshal(errResponse{
		Error: "bad_request",
		Msg:   msg,
	})
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("can`t marshal json"))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(400)
	writer.Write(bytes)
}

func (api *API) GetSwaggerAPI(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadFile("./resources/templates/swagger.html")
	if err != nil {
		log.Error("GetSwaggerAPI: ReadFile", zap.Error(err))
		return
	}
	_, err = w.Write(body)
	if err != nil {
		log.Error("GetSwaggerAPI: Write", zap.Error(err))
		return
	}
}

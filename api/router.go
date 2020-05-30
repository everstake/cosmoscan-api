package api

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
)

// Route stores an API route data
type Route struct {
	Path       string
	Method     string
	Func       func(http.ResponseWriter, *http.Request)
	Middleware []negroni.HandlerFunc
}

// HandleActions is used to handle all given routes
func HandleActions(router *mux.Router, wrapper *negroni.Negroni, prefix string, routes []*Route) {
	for _, r := range routes {
		w := wrapper.With()
		for _, m := range r.Middleware {
			w.Use(m)
		}

		w.Use(negroni.Wrap(http.HandlerFunc(r.Func)))
		router.Handle(prefix+r.Path, w).Methods(r.Method, "OPTIONS")
	}
}

func getParamsFromVars(r *http.Request) map[string][]string {
	mp := make(map[string][]string, 0)
	for k, v := range mux.Vars(r) {
		mp[k] = []string{v}
	}
	return mp
}

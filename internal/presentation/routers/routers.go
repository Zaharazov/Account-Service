package routers

import (
	"Account-Service/internal/server/logger"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes Routes

func NewRouter() *mux.Router {

	routes = append(routes, s_routes...)
	routes = append(routes, o_routes...)
	routes = append(routes, e_routes...)
	routes = append(routes, u_routes...)

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

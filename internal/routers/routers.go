package routers

import (
	"net/http"
	"restapi/internal/api/users"
	"restapi/internal/logger"

	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
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

var routes = Routes{
	Route{
		"CreateUser",
		strings.ToUpper("Post"),
		"/v1/users",
		users.CreateUser,
	},

	Route{
		"DeleteUserById",
		strings.ToUpper("Delete"),
		"/v1/users/{user_id}",
		users.DeleteUserById,
	},

	Route{
		"GetUserById",
		strings.ToUpper("Get"),
		"/v1/users/{user_id}",
		users.GetUserById,
	},

	Route{
		"GetAllUsers",
		strings.ToUpper("Get"),
		"/v1/users",
		users.GetAllUsers,
	},

	Route{
		"EditUserById",
		strings.ToUpper("Put"),
		"/v1/users/{user_id}",
		users.EditUserById,
	},
}

package routers

import (
	"fmt"
	"net/http"
	"restapi/integer/api/users"
	"restapi/integer/logger"
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

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/v1/",
		Index,
	},

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
		"/v1/users/{users_id}",
		users.EditUserById,
	},
}

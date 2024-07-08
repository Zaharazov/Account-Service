package routers

import (
	"Account-Service/internal/api/organizers"
	"Account-Service/internal/api/students"
	"Account-Service/internal/api/users"
	"Account-Service/internal/logger"
	"net/http"
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
		"GetOrganizersByParams",
		strings.ToUpper("Get"),
		"/v1/users/organizers/search", // Organizers
		organizers.GetOrganizersByParams,
	},

	Route{
		"GetOrganizerById",
		strings.ToUpper("Get"),
		"/v1/users/organizers/{user_id}/",
		organizers.GetOrganizerById,
	},

	Route{
		"GetAllOrganizers",
		strings.ToUpper("Get"),
		"/v1/users/organizers/",
		organizers.GetAllOrganizers,
	},

	Route{
		"EditOrganizersById",
		strings.ToUpper("Put"),
		"/v1/users/organizers/{user_id}/",
		organizers.EditOrganizerById,
	},

	Route{
		"GetStudentsByParams",
		strings.ToUpper("Get"),
		"/v1/users/students/search", // Students
		students.GetStudentsByParams,
	},

	Route{
		"GetStudentById",
		strings.ToUpper("Get"),
		"/v1/users/students/{user_id}/",
		students.GetStudentById,
	},

	Route{
		"GetAllStudents",
		strings.ToUpper("Get"),
		"/v1/users/students/",
		students.GetAllStudents,
	},

	Route{
		"EditStudentById",
		strings.ToUpper("Put"),
		"/v1/users/students/{user_id}/",
		students.EditStudentById,
	},

	Route{
		"CreateUser",
		strings.ToUpper("Post"), // Users
		"/v1/users/",
		users.CreateUser,
	},

	Route{
		"DeleteUserById",
		strings.ToUpper("Delete"),
		"/v1/users/{user_id}/",
		users.DeleteUserById,
	},

	Route{
		"GetUserById",
		strings.ToUpper("Get"),
		"/v1/users/{user_id}/",
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
		"/v1/users/{user_id}/",
		users.EditUserById,
	},
}

package routers

import (
	"Account-Service/internal/services"
	"strings"
)

var e_routes = Routes{

	Route{
		"GetEmployersByParams",
		strings.ToUpper("Get"),
		"/v1/users/employers/search", // Employers
		services.GetEmployersByParams,
	},

	Route{
		"GetEmployerById",
		strings.ToUpper("Get"),
		"/v1/users/employers/{user_id}/",
		services.GetEmployerById,
	},

	Route{
		"GetAllEmployers",
		strings.ToUpper("Get"),
		"/v1/users/employers/",
		services.GetAllEmployers,
	},

	Route{
		"EditEmployersById",
		strings.ToUpper("Put"),
		"/v1/users/employers/{user_id}/",
		services.EditEmployerById,
	},
}

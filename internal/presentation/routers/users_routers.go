package routers

import (
	"Account-Service/internal/services"
	"strings"
)

var u_routes = Routes{

	Route{
		"CreateUser",
		strings.ToUpper("Post"), // Users
		"/v1/users/",
		services.CreateUser,
	},

	Route{
		"DeleteUserById",
		strings.ToUpper("Delete"),
		"/v1/users/{user_id}/",
		services.DeleteUserById,
	},

	Route{
		"GetUserById",
		strings.ToUpper("Get"),
		"/v1/users/{user_id}/",
		services.GetUserById,
	},

	Route{
		"GetAllUsers",
		strings.ToUpper("Get"),
		"/v1/users",
		services.GetAllUsers,
	},

	Route{
		"EditUserById",
		strings.ToUpper("Put"),
		"/v1/users/{user_id}/",
		services.EditUserById,
	},
}

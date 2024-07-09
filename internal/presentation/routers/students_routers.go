package routers

import (
	"Account-Service/internal/services"
	"strings"
)

var s_routes = Routes{

	Route{
		"GetStudentsByParams",
		strings.ToUpper("Get"),
		"/v1/users/students/search", // Students
		services.GetStudentsByParams,
	},

	Route{
		"GetStudentById",
		strings.ToUpper("Get"),
		"/v1/users/students/{user_id}/",
		services.GetStudentById,
	},

	Route{
		"GetAllStudents",
		strings.ToUpper("Get"),
		"/v1/users/students/",
		services.GetAllStudents,
	},

	Route{
		"EditStudentById",
		strings.ToUpper("Put"),
		"/v1/users/students/{user_id}/",
		services.EditStudentById,
	},
}

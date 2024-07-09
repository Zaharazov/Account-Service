package routers

import (
	"Account-Service/internal/services"
	"strings"
)

var o_routes = Routes{

	Route{
		"GetOrganizersByParams",
		strings.ToUpper("Get"),
		"/v1/users/organizers/search", // Organizers
		services.GetOrganizersByParams,
	},

	Route{
		"AddEventToOrganizerById",
		strings.ToUpper("Put"),
		"/v1/users/organizers/{user_id}/events/{event_id}",
		services.AddEventToOrganizerById,
	},

	Route{
		"GetOrganizerById",
		strings.ToUpper("Get"),
		"/v1/users/organizers/{user_id}/",
		services.GetOrganizerById,
	},

	Route{
		"GetAllOrganizers",
		strings.ToUpper("Get"),
		"/v1/users/organizers/",
		services.GetAllOrganizers,
	},

	Route{
		"EditOrganizersById",
		strings.ToUpper("Put"),
		"/v1/users/organizers/{user_id}/",
		services.EditOrganizerById,
	},
}

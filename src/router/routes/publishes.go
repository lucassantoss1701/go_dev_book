package routes

import (
	"api/src/controllers"
	"net/http"
)

var routesPublishes = []Route{
	{
		URI:                    "/publishes",
		Method:                 http.MethodPost,
		Function:               controllers.CreatePublish,
		RequiredAuthentication: true,
	},
	{
		URI:                    "/publishes",
		Method:                 http.MethodGet,
		Function:               controllers.GetPublishes,
		RequiredAuthentication: true,
	},
	{
		URI:                    "/publishes/{publishId}",
		Method:                 http.MethodGet,
		Function:               controllers.GetPublish,
		RequiredAuthentication: true,
	},
	{
		URI:                    "/publishes/{publishId}",
		Method:                 http.MethodPut,
		Function:               controllers.UpdatePublish,
		RequiredAuthentication: true,
	},
	{
		URI:                    "/publishes/{publishId}",
		Method:                 http.MethodDelete,
		Function:               controllers.DeletePublish,
		RequiredAuthentication: true,
	},
	{
		URI:                    "/usuarios/{userId}/publishes",
		Method:                 http.MethodGet,
		Function:               controllers.GetPublishByUser,
		RequiredAuthentication: true,
	},
	{
		URI:                    "/publishes/{publishId}/like",
		Method:                 http.MethodPost,
		Function:               controllers.LikePublish,
		RequiredAuthentication: true,
	},
	{
		URI:                    "/publishes/{publishId}/deslike",
		Method:                 http.MethodPost,
		Function:               controllers.DeslikePublish,
		RequiredAuthentication: true,
	},
}

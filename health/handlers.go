package health

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/movieManagement/gen/restapi/operations"
	"github.com/movieManagement/gen/restapi/operations/health"
	"github.com/movieManagement/swagger"
)

// Configure setups handlers on api with Service
func Configure(api *operations.MovieServiceAPI, service Service) {
	api.HealthGetHealthHandler = health.GetHealthHandlerFunc(func(params health.GetHealthParams) middleware.Responder {

		result, err := service.HealthCheck(params.HTTPRequest.Context(), &params)
		if err != nil {
			return health.NewGetHealthBadRequest().WithPayload(swagger.ErrorResponse(err))
		}
		return health.NewGetHealthOK().WithPayload(result)
	})
}

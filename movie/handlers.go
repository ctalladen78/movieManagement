package movie

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/movieManagement/gen/restapi/operations"
	"github.com/movieManagement/gen/restapi/operations/movie"
	"github.com/movieManagement/swagger"
)

// Configure configures the affiliation service
func Configure(api *operations.MovieServiceAPI, service Service) {
	api.MovieCreateMovieHandler = movie.CreateMovieHandlerFunc(func(params movie.CreateMovieParams) middleware.Responder {
		result, err := service.CreateMovie(params.HTTPRequest.Context(), &params)
		if err != nil {
			return swagger.ErrorHandler("CreateMovie :: ", err)
		}
		return movie.NewCreateMovieCreated().WithPayload(result)
	})

	api.MovieGetmovieHandler = movie.GetmovieHandlerFunc(func(params movie.GetmovieParams) middleware.Responder {
		return middleware.NotImplemented("operation movies.Getmovie has not yet been implemented")
	})

	api.MovieSearchMoviesHandler = movie.SearchMoviesHandlerFunc(func(params movie.SearchMoviesParams) middleware.Responder {
		result, err := service.SearchMovies(params.HTTPRequest.Context(), &params)
		if err != nil {
			return swagger.ErrorHandler("SearchMovies :: ", err)
		}
		return movie.NewSearchMoviesOK().WithPayload(result)
	})
}

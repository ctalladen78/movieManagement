package swagger

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/movieManagement/errs"
	"github.com/movieManagement/gen/models"
	"github.com/movieManagement/gen/restapi/operations/movie"
	log "github.com/sirupsen/logrus"
)

type codedResponse interface {
	Code() string
}

// ErrorResponse wraps the error in the api standard models.ErrorResponse object
func ErrorResponse(err error) *models.ErrorResponse {
	cd := ""
	if e, ok := err.(codedResponse); ok {
		cd = e.Code()
	}

	e := models.ErrorResponse{
		Code:    cd,
		Message: err.Error(),
	}
	return &e
}

// ErrorHandler accepts a string and error and returns the appropriate responder
func ErrorHandler(label string, err error) middleware.Responder {
	//logrus.WithError(err).Error(label)
	orignalErr := err
	errorStr := err.Error()
	errorStrSplit := strings.Split(errorStr, ":")

	if len(errorStrSplit) > 1 {
		err = errors.New(strings.TrimSpace(errorStrSplit[len(errorStrSplit)-1]))
	}

	switch err.Error() {
	case errs.ErrUnauthorized.Error():
		log.Error(err)
		return movie.NewCreateMovieUnauthorized().WithPayload(ErrorResponse(orignalErr))
	case errs.ErrForbidden.Error():
		log.Error(err)
		return movie.NewCreateMovieForbidden().WithPayload(ErrorResponse(orignalErr))
	case errs.ErrNotFound.Error(), sql.ErrNoRows.Error():
		log.Error(err)
		return movie.NewGetmovieNotFound().WithPayload(ErrorResponse(orignalErr))
	case errs.ErrConflict.Error():
		log.Error(err)
		return movie.NewCreateMovieConflict().WithPayload(ErrorResponse(orignalErr))
	default:
		log.Error(err)
		return movie.NewGetmovieBadRequest().WithPayload(ErrorResponse(orignalErr))
	}
}

package movie

import (
	"context"
	"strconv"

	"github.com/labstack/gommon/log"
	"github.com/movieManagement/gen/models"
	"github.com/movieManagement/gen/restapi/operations/movie"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Service interface is a list of services for the affiliation
type Service interface {
	CreateMovie(ctx context.Context, in *movie.CreateMovieParams) (*models.Movie, error)
	SearchMovies(ctx context.Context, in *movie.SearchMoviesParams) (*models.MovieList, error)
}

type service struct {
	repo Repository
}

// New is a simple helper function to create a service instance
func New(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// CreateMovie service definition
func (s *service) CreateMovie(ctx context.Context, in *movie.CreateMovieParams) (*models.Movie, error) {
	logrus.Debugf("entered service CreateAffiliation")
	movie, err := s.repo.CreateMovie(ctx, in)
	if err != nil {
		log.Error(err)
		return nil, errors.Wrap(err, "service.CreateAffiliation")
	}
	return movie, nil
}

// SearchMovies service definition
func (s *service) SearchMovies(ctx context.Context, in *movie.SearchMoviesParams) (*models.MovieList, error) {
	log.Debugf("entered service ListCommunities")

	var movies []*models.Movie
	var meta models.ListMetadata
	var ol models.MovieList
	var err error
	var count int64

	movies, count, err = s.repo.SearchMovies(ctx, in)
	if err != nil {
		log.Error(err)
		return nil, errors.Wrap(err, "service.SearchMovies")
	}

	offset, err := strconv.Atoi(*in.Offset)
	if err != nil {
		log.Error(err)
		return nil, errors.Wrap(err, "service.convertOffset")
	}

	pageSize, err := strconv.Atoi(*in.PageSize)
	if err != nil {
		log.Error(err)
		return nil, errors.Wrap(err, "service.convertPageSize")
	}

	// We should always have movie in our response unless we have an empty database
	// rather than returning 'null' in JSON, we will allocate an empty array so it unmarshalled nicely as [] for the client
	if movies == nil {
		movies = make([]*models.Movie, 0)
	}

	meta.Offset = int64(offset)
	meta.PageSize = int64(pageSize)
	meta.TotalSize = count
	ol.Data = movies
	ol.Metadata = &meta
	return &ol, nil
}

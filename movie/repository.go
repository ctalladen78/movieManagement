package movie

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	gomdb "github.com/eefret/go-imdb"
	"github.com/google/uuid"
	"github.com/ido50/sqlz"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/movieManagement/gen/models"
	"github.com/movieManagement/gen/restapi/operations/movie"
	ini "github.com/movieManagement/init"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	// MovieTable . . .
	MovieTable = "public.moviestbl as mv"
)

// sorted by field alias. Same as movieReturnFields
var movieReturnFields = []string{
	"COALESCE(mv.createddate, '2019-01-01') as CreatedAt",
	"COALESCE(mv.title, '') as Title",
	"COALESCE(mv.rating, '') as Rating",
	"COALESCE(mv.releasedYear, '') as ReleasedYear",
	"COALESCE(mv.genres, '') as Genres",
	"COALESCE(mv.lastmodifieddate, '2019-01-01') as LastModifiedAt",
	"COALESCE(mv.sfid, '') as ID",
}

// Repository interface includes a list of supported repository operations
type Repository interface {
	CreateMovie(ctx context.Context, params *movie.CreateMovieParams) (*models.Movie, error)
	SearchMovies(ctx context.Context, params *movie.SearchMoviesParams) ([]*models.Movie, int64, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new repository from the specified DB reference
func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

// GetDB returns a reference to the underlying database connection
func (repo *repository) GetDB() *sqlx.DB {
	return repo.db
}

// CreateMovie create the affiliation..
func (repo *repository) CreateMovie(ctx context.Context, params *movie.CreateMovieParams) (*models.Movie, error) {
	logrus.Debugf("CreateMovie repo")
	var movies []*models.Movie
	sqlMovies := SQLMovies{}
	uuid := uuid.New().String()
	createMap := insertFields(params, repo)

	createMap["lastmodifieddate"] = sqlz.Indirect("now()::timestamp")
	createMap["createddate"] = sqlz.Indirect("now()::timestamp")
	createMap["sfid"] = uuid

	err := sqlz.Newx(repo.db).
		InsertInto(MovieTable).
		ValueMap(createMap).
		Returning(movieReturnFields...).
		GetRow(&sqlMovies)
	if err != nil {
		logrus.Errorf("error to create membership %v", err)
		//return nil, errors.Wrap(err, "CreateMembership.Exec")
	}
	var movie = sqlMovies.toMovie()

	movies = append(movies, movie)

	return movies[0], nil
}

func insertFields(params *movie.CreateMovieParams, repo *repository) map[string]interface{} {
	var genresDetails *string
	insertMap := make(map[string]interface{})

	addIfNotEmpty(insertMap, "title", params.Movie.Title)
	addIfNotEmpty(insertMap, "releasedYear", params.Movie.ReleasedYear)
	addIfNotEmpty(insertMap, "rating", params.Movie.Rating)

	if len(params.Movie.Genres) > 0 {
		val := strings.Join(params.Movie.Genres[:], ",")
		genresDetails = &val
	}

	addIfNotNil(insertMap, "genres", genresDetails)

	return insertMap
}

// addIfNotEmpty simply adds the key/value pair if the key and value is not empty
func addIfNotEmpty(m map[string]interface{}, key string, value string) {
	if key != "" && value != "" {
		m[key] = value
	}
}

/* // addIfNotNil simply adds the key/value pair if the key and value is not nil
func addIfNotNilBool(m map[string]interface{}, key string, value *bool) {
	if key != "" && value != nil {
		m[key] = *value
	}
} */

// addIfNotNil simply adds the key/value pair if the key and value is not nil
func addIfNotNil(m map[string]interface{}, key string, value *string) {
	if key != "" && value != nil {
		m[key] = *value
	}
}

// SearchMovies returns a list of movies based on the input
// parameters and security permissions
func (repo *repository) SearchMovies(ctx context.Context, params *movie.SearchMoviesParams) ([]*models.Movie, int64, error) {
	log.Debugf("entered function ListCommunities")
	code := "SearchMovies"
	community, count, err := getMovies(ctx, params, repo, code)
	if err != nil {
		log.Error(err)
		return nil, 0, errors.Wrap(err, "ListCommunities.getCommunities")
	}
	return community, count, nil
}

func getMovies(ctx context.Context, params *movie.SearchMoviesParams, repo *repository, code string) ([]*models.Movie, int64, error) {
	pageSize, err := strconv.Atoi(*params.PageSize)
	if err != nil {
		log.Error(err)
		return nil, 0, errors.Wrap(err, fmt.Sprintf("%s.%s", code, "convertPageSize"))
	}

	offset, err := strconv.Atoi(*params.Offset)
	if err != nil {
		log.Error(err)
		return nil, 0, errors.Wrap(err, fmt.Sprintf("%s.%s", code, "convertOffset"))
	}

	var id, title, rating, year string

	if params.ID == nil {
		id = "%"
	} else {
		id = *params.ID
	}

	if params.Title == nil {
		title = "%"
	} else {
		title = *params.Title
	}

	if params.Rating == nil {
		rating = "%"
	} else {
		rating = *params.Rating
	}

	if params.Year == nil {
		year = "%"
	} else {
		year = *params.Year
	}

	var createdMovie *models.Movie
	var movieArray []*models.Movie
	sqlMovies := []SQLMovies{}
	conditions := []sqlz.WhereCondition{}

	if id != "%" {
		conditions = append(conditions, sqlz.Eq("mv.sfid", id))
	}

	if title != "%" {
		conditions = append(conditions, sqlz.Eq("mv.title", title))
	}

	if rating != "%" {
		conditions = append(conditions, sqlz.Eq("mv.rating", rating))
	}

	if year != "%" {
		conditions = append(conditions, sqlz.Eq("mv.releasedYear", year))
	}

	if len(params.Genres) != 0 {
		var genresList []interface{}
		for _, val := range params.Genres {
			genresList = append(genresList, val)
		}
		conditions = append(conditions, sqlz.In("mv.genres", genresList...))
	}

	query := sqlz.Newx(repo.GetDB()).
		Select(movieReturnFields...).
		From(MovieTable).
		Where(conditions...).
		Limit(int64(pageSize)).
		Offset(int64(offset))

	sql, b := query.ToSQL(true)
	log.Info(sql, b)

	count, errCount := query.GetCount()
	if errCount != nil {
		log.Error(err)
		return nil, 0, errors.Wrap(err, fmt.Sprintf("%s.%s", code, "GetCount"))
	}
	err = query.GetAll(&sqlMovies)

	if err != nil {
		log.Error(err)
		return nil, 0, errors.Wrap(err, fmt.Sprintf("%s.%s", code, "SelectQuery"))
	}
	for _, sqlMovie := range sqlMovies {
		var movieData = sqlMovie.toMovie()
		movieArray = append(movieArray, movieData)
	}

	if len(movieArray) == 0 {
		imdb := ini.GetImdbInit()
		movieObject, err := imdb.MovieByTitle(&gomdb.QueryData{Title: *params.Title})
		if err != nil {
			log.Error(err)
			return nil, 0, errors.Wrap(err, fmt.Sprintf("%s.%s", code, "MovieByTitle"))
		}
		log.Debugf("movieObject %s", movieObject)
		if movieObject != nil {
			var in movie.CreateMovieParams
			in.Movie.Title = movieObject.Title
			in.Movie.Rating = movieObject.ImdbRating
			in.Movie.ReleasedYear = movieObject.Released
			in.Movie.Genres = []string{movieObject.Genre}
			createdMovie, err = repo.CreateMovie(ctx, &in)
			if err != nil {
				log.Error(err)
				return nil, 0, errors.Wrap(err, fmt.Sprintf("%s.%s", code, "CreateMovie"))
			}
			if createdMovie != nil {
				count = count + 1
			}
		}
	}
	movieArray = append(movieArray, createdMovie)
	return movieArray, count, nil
}

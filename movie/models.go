package movie

import (
	"database/sql"

	"github.com/go-openapi/strfmt"
	"github.com/movieManagement/gen/models"
)

// SQLMovies . . .
type SQLMovies struct {
	ID             sql.NullString  `json:"ID,omitempty"`
	Title          sql.NullString  `json:"Title,omitempty"`
	LastModifiedAt strfmt.DateTime `json:"LastModifiedAt,omitempty"`
	CreatedAt      strfmt.DateTime `json:"CreatedAt,omitempty"`
	Genres         sql.NullString  `json:"Genres,omitempty"`
	Rating         sql.NullString  `json:"Rating,omitempty"`
	ReleasedYear   sql.NullString  `json:"ReleasedYear,omitempty"`
}

func (sql *SQLMovies) toMovie() *models.Movie {

	movie := models.Movie{
		ID:             sql.ID.String,
		Genres:         []string{sql.Genres.String},
		LastModifiedAt: sql.LastModifiedAt,
		CreatedAt:      sql.CreatedAt,
		Rating:         sql.Rating.String,
		ReleasedYear:   sql.ReleasedYear.String,
		Title:          sql.Title.String,
	}
	return &movie
}

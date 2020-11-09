package health

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/movieManagement/gen/models"
	"github.com/movieManagement/gen/restapi/operations/health"
)

// Service handles async log of audit event
type Service interface {
	HealthCheck(ctx context.Context, in *health.GetHealthParams) (*models.Health, error)
}

type service struct {
	hcdb       *sqlx.DB
	gitHash    string
	buildStamp string
}

// New is a simple helper function to create a service instance
func New(hcdb *sqlx.DB, GitHash, BuildStamp string) Service {
	return &service{
		hcdb:       hcdb,
		gitHash:    GitHash,
		buildStamp: BuildStamp,
	}
}

func dbCheck(db *sqlx.DB) error {
	//var pong string
	//return db.Get(&pong, "select 'pong'")
	return nil
}

func (s service) HealthCheck(ctx context.Context, in *health.GetHealthParams) (*models.Health, error) {
	t := time.Now()
	dbErr := dbCheck(s.hcdb)
	duration := time.Since(t)

	hs := models.HealthStatus{TimeStamp: time.Now().String(), Healthy: true, Name: "HCDB", Duration: duration.String()}
	if dbErr != nil {
		hs.Healthy = false
		hs.Error = dbErr.Error()
	}

	response := models.Health{
		Status:         "healthy",
		TimeStamp:      time.Now().String(),
		Githash:        s.gitHash,
		BuildTimeStamp: s.buildStamp,
		Healths:        []*models.HealthStatus{&hs}}
	return &response, nil
}

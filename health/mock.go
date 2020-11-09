package health

import (
	"context"
	"time"

	"github.com/movieManagement/gen/models"
	"github.com/movieManagement/gen/restapi/operations/health"
)

type mock struct {
}

// NewMock is a simple helper function to create a Mock service instance
func NewMock() Service {
	return &mock{}
}

func (mock) HealthCheck(ctx context.Context, in *health.GetHealthParams) (*models.Health, error) {
	hs := models.HealthStatus{TimeStamp: time.Now().String(), Healthy: true, Name: "TestHealth", Duration: (time.Millisecond * 5).String()}

	response := models.Health{
		Status:    "healthy",
		TimeStamp: time.Now().String(),
		Healths:   []*models.HealthStatus{&hs},
	}
	return &response, nil
}

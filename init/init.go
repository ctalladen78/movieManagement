package init

import (
	"fmt"

	imdb "github.com/eefret/go-imdb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	stage string
	im    *imdb.OmdbApi
)

// CommonInit initializes the common properties
func CommonInit() {
	stage = getProperty("STAGE")

}

// getProperty is a common routine to bind and return the specified environment variable
func getProperty(property string) string {
	err := viper.BindEnv(property)
	if err != nil {
		message := fmt.Sprintf("Unable to load property: %s_%s - value not defined or empty", "Movie", property)
		logrus.Fatal(message)
	}

	value := viper.GetString(property)
	if value == "" {
		err := fmt.Errorf("%s_%s environment variable cannot be empty", "Movie", property)
		logrus.Fatal(err)
	}

	return value
}

// Init initialization logic for all the handlers
func Init() {
	CommonInit()
}

// GetStage returns the deployment stage, e.g. dev, test, stage or prod
func GetStage() string {
	return stage
}

// ImdbVariable loads all the SSM values based on stage.
func ImdbVariable() {
	im = imdb.Init("my-api-key")
}

// GetImdbInit . . .
func GetImdbInit() *imdb.OmdbApi {
	return im
}

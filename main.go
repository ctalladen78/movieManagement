package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/go-openapi/loads"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/movieManagement/cmd"
	"github.com/movieManagement/gen/restapi"
	"github.com/movieManagement/gen/restapi/operations"
	"github.com/movieManagement/health"
	"github.com/movieManagement/movie"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Build and version variables defined and set during the build process
var (
	// Version the application version
	version string
	// Build/Commit the application build number
	commit string
	// Build date
	buildDate string
)

func setupDefaults() map[string]interface{} {
	viper.AutomaticEnv()
	defaults := map[string]interface{}{
		"PORT":               8080,
		"APP_ENV":            "local",
		"USE_MOCK":           "False",
		"HCDB":               "host=localhost port=3306 dbname=pmm sslmode=disable application_name='pmm'",
		"DB_MAX_CONNECTIONS": 20,
	}

	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	return defaults
}

// getProperty is a common routine to bind and return the specified environment variable
func getProperty(property string) string {
	err := viper.BindEnv(property)
	if err != nil {
		message := fmt.Sprintf("Unable to load property: %s - value not defined or empty", property)
		logrus.Fatal(message)
	}

	value := viper.GetString(property)
	if value == "" {
		err := fmt.Errorf("%s environment variable cannot be empty", property)
		logrus.Fatal(err)
	}

	return value
}

func initDB(name string) *sqlx.DB {
	hcDBURL := getProperty(name)
	logrus.Infof("Initializing DB %s connection with URL (prefix) %s...", name, hcDBURL[0:11])
	d, err := sqlx.Connect("postgres", hcDBURL)
	if err != nil {
		logrus.Panic(err)
	}

	d.SetMaxOpenConns(viper.GetInt("DB_MAX_CONNECTIONS"))
	d.SetMaxIdleConns(5)
	d.SetConnMaxLifetime(15 * time.Minute)

	return d
}

func main() {

	_ = setupDefaults()

	host, err := os.Hostname()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("Service Startup")

	var portFlag = flag.Int("port", viper.GetInt("PORT"), "Port to listen for web requests on")

	// Show the version and build info
	logrus.Infof("Version               : %s", version)
	logrus.Infof("Git commit hash       : %s", commit)
	logrus.Infof("Build date            : %s", buildDate)
	logrus.Infof("Golang OS             : %s", runtime.GOOS)
	logrus.Infof("Golang Arch           : %s", runtime.GOARCH)
	logrus.Infof("Service Host          : %s", host)
	logrus.Infof("Service Port          : %d", *portFlag)

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		logrus.Fatal(err)
	}

	// Initialize hcDB connection
	hcDB := initDB("HCDB")

	api := operations.NewMovieServiceAPI(swaggerSpec)

	// Setup the movie service
	movieRepo := movie.NewRepository(hcDB)
	movieService := movie.New(movieRepo)
	movie.Configure(api, movieService)

	// Setup the health service
	var healthService health.Service
	if viper.GetBool("USE_MOCK") {
		healthService = health.NewMock()
	} else {
		healthService = health.New(hcDB, commit, buildDate)
	}
	health.Configure(api, healthService)

	flag.Parse()

	if err := cmd.Start(api, *portFlag); err != nil {
		logrus.Fatal(err)
	}
}

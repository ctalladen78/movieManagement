// +build !aws_lambda

package cmd

import (
	"github.com/movieManagement/gen/restapi"
	"github.com/movieManagement/gen/restapi/operations"
	ini "github.com/movieManagement/init"
	log "github.com/movieManagement/logging"
)

// Start is the local entry point method
func Start(api *operations.MovieServiceAPI, portFlag int) error {
	log.Infof("Running Init()...")
	ini.Init()
	server := restapi.NewServer(api)
	defer server.Shutdown() // nolint
	server.Port = portFlag

	return server.Serve()
}

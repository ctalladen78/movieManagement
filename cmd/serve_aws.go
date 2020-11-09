// +build aws_lambda

package cmd

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/movieManagement/gen/restapi/operations"
	ini "github.com/movieManagement/init"
	log "github.com/movieManagement/logging"
)

// Start is the lambda main entry point
func Start(api *operations.MovieServiceAPI, _ int) error {

	log.Infof("Running Init()...")
	ini.Init()

	adapter := httpadapter.New(api.Serve(nil))

	log.Debugf("Starting Lambda")
	lambda.Start(adapter.Proxy)
	return nil
}

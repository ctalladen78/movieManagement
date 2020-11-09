# Movie Service


This repository contains schema, business logic, service endpoints, CI/CD and
deployment configuration for the Movie Service.

## Service Endpoints

| Endpoint                     | Action    | Description                                                                | Source      |
|:-----------------------------|:----------|:---------------------------------------------------------------------------|:------------|
| /health                      | GET       | Returns a short JSON document indicating the overall health of the service |DB         |
| movies                | GET       | Returns my movies                              |DB         |
| movies                 | POST      | Create a new movie under current user                                            |DB    |
| /api-docs                    | GET       | Returns a fancy HTML page for the swagger documentation                    |Swagger file |


## Data Models

### Service Request/Response Models

The service request and response models are defined in the `swagger`
specification file. `go` code is generated from these definitions
by running the `make swagger` Makefile target.  The result of
running the swagger generator is a set of models and REST API
objects in the `gen` folder. These should NOT be committed to git.

## Building

### Prerequisites

To build the project, the following tools are required:

* Make - a build tool
* Go 1.12+ - Go language compiler
* Node 6.5+ 
* Swagger tool - downloaded automatically on build - a code generator for REST API

### Developer Setup

```bin
cd ~/go/ && \
mkdir -p src/github.com/ && \

git clone git@github.com:movieManagement ~/go/src/github.com/movieManagement

cd ~/go/src/github.com/movieManagement
make setup swagger build lint test
```
> The movie service requires database which uses postgressql DB. Make sure env is having the connection string of postgressql DB. Once Env sourced, we can run the proejct as
```
make run
```

### Execute Build Targets

To build the project, run:

```bash
make
```

To build clean with testing and linting:

```bash
make clean build test lint
```

## Running Locally

When running locally:

* Generate the models from swagger
* Build the codebase
* Setup the Environment
* Run the main executable

### Environment

See deployment environment notes below.

### Command

```bash
make run

# or simply
./bin/movie-service
```

## Database

### AWS Setup

Create a deployment-user user with "Programmatic Access' enabled (or use an
existing deployment user ID) for each environment and set the "Administrator
Access" permissions.  This user will be leveraged by the CircleCI pipeline
to invoke the [serverless](serverless.yml) configuration. This tool will
generate the necessary resources within AWS. Create and save the API keys.

Create `movie-service-user` user for each environment with "Programmatic access"
enabled. Add the following policies:

* `AWSLambdaFullAccess`
* `AmazonSNSFullAccess`
* `AmazonSNSRole` - allows `logs:CreateLogGroup`, `logs:CreateLogStream`, `logs:PutLogEvents`, `logs:PutMetricFilter`, and `logs:PutRetentionPolicy`
    
### Environment Variables

Environment variables need to be set either in the CI/CD system which are
passed through to the deployment run-time or within the deployment
configuration file. The following sections describes each.

#### Deployment Environment Variables

Create a `deployment-user` (see AWS Setup above or use an existing user ID) for
each environment and set the "Administrator Access" permissions (see above in
the AWS Setup section). Replace the ${CURRENT_STAGE} notation below ith the
appropriate deployment environment value, such as `DEV`, `STAGING`, or `PROD`.
For example `AWS_DEFAULT_REGION_${CURRENT_STAGE}` becomes
`AWS_DEFAULT_REGION_DEV` for the DEV environment.

* `AWS_DEFAULT_REGION_${CURRENT_STAGE}` - the AWS region, e.g. us-east-2 for the stage
* `AWS_ACCOUNT_ID_${CURRENT_STAGE}` - the AWS account ID for the stage
* `AWS_ACCESS_KEY_ID_${CURRENT_STAGE}` - the AWS access key ID of the `deployment-user` - used by the serverless deployment to create the required AWS resources
* `AWS_SECRET_ACCESS_KEY_${CURRENT_STAGE}` the AWS secret access key of the `deployment-user` - used by the serverless deployment to create the required AWS resources
* `MOVIE_SERVICE_HCDB` The movie service requires database which uses postgressql DB

The following environment variables are statically or dynamically set in the
[serverless.yml](serverless.yml) file and therefore available to the
application (no need to set these values - they are created/computed and made
available to the runtime environment).

* `STAGE` - set by the CI system configuration based on the deployment stage
* `MOVIE_SERVICE_STAGE`- set by the CI system configuration based on the deployment stage
* `MOVIE_SERVICE_HCDB` - set dynamically by the CI system based on the environment STAGE
* `MOVIE_SERVICE_AWS_ACCESS_KEY_ID`  - set dynamically by the CI deployment script based on the environment AWS ACCESS KEY ID
* `MOVIE_SERVICE_AWS_SECRET_ACCESS_KEY` - set dynamically by the CI deployment script based on the environment AWS SECRET ACCESS KEY
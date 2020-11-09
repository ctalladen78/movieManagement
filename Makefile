BUILD_TIME=`date +%FT%T%z`
VERSION := $(shell sh -c 'git describe --always --tags')
BRANCH := $(shell sh -c 'git rev-parse --abbrev-ref HEAD')
COMMIT := $(shell sh -c 'git rev-parse --short HEAD')
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.branch=$(BRANCH) -X main.buildDate=$(BUILD_TIME)"
LINT_TOOL=$(shell go env GOPATH)/bin/golangci-lint
BUILD_TAGS=-tags aws_lambda
GO_PKGS=$(shell go list ./... | grep -v /vendor/ | grep -v /node_modules/)
GO_FILES=$(shell find . -type f -name '*.go' -not -path './vendor/*')

.PHONY: setup_dev setup_deploy build build-mac swagger fmt clean test lint qc deploy

setup: $(LINT_TOOL) setup_dev setup_deploy

setup_dev:
	./scripts/install-go-swagger.sh
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/golang/dep/cmd/dep	
	go get golang.org/x/tools/cmd/cover
	go get -u github.com/stripe/safesql

setup_deploy:
	npm install serverless

deps:
	dep ensure

build: deps
	env GOOS=linux GOARCH=amd64 go build $(BUILD_TAGS) $(LDFLAGS) -o bin/movie-service main.go
	chmod +x bin/movie-service

build-mac: deps
	env GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/movie-service main.go
	chmod +x bin/movie-service

swagger:
	@rm -Rf gen
	@mkdir -p gen
	swagger -q generate server -t gen -f swagger/movie-service.yaml --exclude-main -A MovieService 

fmt:
	@go fmt $(GO_PKGS)
	@goimports -w -l $(GO_FILES)

test:
	@go test -v $(shell go list ./... | grep -v /vendor/ | grep -v /node_modules/) -coverprofile=cover.out

clean:
	rm -rf ./gen ./bin ./vendor Gopkg.lock

$(LINT_TOOL):
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.28.3

qc: $(LINT_TOOL)
	$(LINT_TOOL) run --config=.golangci.yaml ./...

lint: qc

safesql:
	safesql -v $(GO_FILES)
	
run:
	go run main.go

deploy: clean build
	sls deploy --verbose

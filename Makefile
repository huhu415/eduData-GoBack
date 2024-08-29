# go build -ldflags "-X main.build=`git rev-parse HEAD`
LDFLAGS := -s -w

VERSION := 1.0.0
BUILD_DATE := $(shell date +%Y-%m-%dT%H:%M:%S)
GIT_COMMIT := $(shell git rev-parse --short HEAD)

# build: build
build:
	@env CGO_ENABLED=0 								go build -trimpath -ldflags "-X 'eduData/bootstrap.Version=$(VERSION)' -X 'eduData/bootstrap.BuildDate=$(BUILD_DATE)' -X 'eduData/bootstrap.GitCommit=$(GIT_COMMIT)'"

# cbuild: cross build
cbuild:
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 		go build -trimpath -ldflags "$(LDFLAGS)" .

# vet: 找错误
vet:
	@go vet ./...

.PHONY: build cbuild vet

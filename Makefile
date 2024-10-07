# go build -ldflags "-X main.build=`git rev-parse HEAD`
LDFLAGS := -s -w

VERSION := 1.0.0
BUILD_DATE := $(shell date +%Y-%m-%dT%H:%M:%S)
GIT_COMMIT := $(shell git rev-parse --short HEAD)  $(shell git log -1 --pretty=%s)

# build: build the project
build:
	@env CGO_ENABLED=0 \
	go build -trimpath \
		-ldflags "$(LDFLAGS) \
		-X 'eduData/bootstrap.Version=$(VERSION)' \
		-X 'eduData/bootstrap.BuildDate=$(BUILD_DATE)' \
		-X 'eduData/bootstrap.GitCommit=$(GIT_COMMIT)'" \
		.

# cbuild: cross build for Linux amd64
cbuild:
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -trimpath \
		-ldflags "$(LDFLAGS) \
		-X 'eduData/bootstrap.Version=$(VERSION)' \
		-X 'eduData/bootstrap.BuildDate=$(BUILD_DATE)' \
		-X 'eduData/bootstrap.GitCommit=$(GIT_COMMIT)'" \
		.

# debug: debug
debug:
	@brew services start postgresql@17
	@CompileDaemon -build="make build" -command="./eduData --debug"

.PHONY: build cbuild

# go build -ldflags "-X main.build=`git rev-parse HEAD`
LDFLAGS := -s -w

# build: build
build:
	@env CGO_ENABLED=0 								go build -trimpath -ldflags "$(LDFLAGS)" .

# cbuild: cross build
cbuild:
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 		go build -trimpath -ldflags "$(LDFLAGS)" .

# vet: 找错误
vet:
	go vet ./...

.PHONY: build cbuild vet

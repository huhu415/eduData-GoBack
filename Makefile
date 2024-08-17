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

# dockerbuild: build docker image
dockerbuildamd64:
	rm eduData
	make cbuild
	docker build -t registry.cn-wulanchabu.aliyuncs.com/zzyan/back-go:latest .
	docker push registry.cn-wulanchabu.aliyuncs.com/zzyan/back-go:latest

.PHONY: build cbuild vet dockerbuildamd64

VERSION := $(shell git describe --always --long --dirty)
VERSION_ARGS := -ldflags="-X main.version=${VERSION}"

run\:piximan:
	go run ${VERSION_ARGS} cmd/piximan/main.go $$ARGS

run\:piximanctl:
	go run ${VERSION_ARGS} cmd/piximanctl/main.go $$ARGS

build\:piximan:
	go build ${VERSION_ARGS} -v -o bin/piximan cmd/piximan/main.go

build\:piximanctl:
	go build ${VERSION_ARGS} -v -o bin/piximanctl cmd/piximanctl/main.go

build:
	make build:piximan
	make build:piximanctl

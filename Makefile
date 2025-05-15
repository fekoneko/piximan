VERSION := $(shell git describe --always --tags --dirty)
VERSION_ARGS := -ldflags="-X main.version=${VERSION}"
MAKEFLAGS += --no-print-directory

define BUILD_CMD
	if [ -z "$${GOOS}" ]; then GOOS="$$(go env GOOS)"; fi; \
  if [ -z "$${GOARCH}" ]; then GOARCH="$$(go env GOARCH)"; fi; \
  BIN_DIR="bin/$${GOOS}_$${GOARCH}"; \
	mkdir -p "$${BIN_DIR}"; \
	BIN_NAME="$${CMD_NAME}"; \
	if [ "$${GOOS}" = 'windows' ]; then BIN_NAME="$${BIN_NAME}.exe"; fi; \
	env GOOS="$${GOOS}" GOARCH="$${GOARCH}" \
		go build ${VERSION_ARGS} -v -o "$${BIN_DIR}/$${BIN_NAME}" "cmd/$${CMD_NAME}/main.go" ${ARGS}
endef

run\:piximan:
	@go run ${VERSION_ARGS} 'cmd/piximan/main.go' ${ARGS}

run\:piximanctl:
	@go run ${VERSION_ARGS} 'cmd/piximanctl/main.go' ${ARGS}

build\:piximan\:current:
	@echo "Building piximan for current platform"
	@CMD_NAME=piximan;    GOOS=;        GOARCH=;      $(call BUILD_CMD)

build\:piximanctl\:current:
	@echo "Building piximanctl for current platform"
	@CMD_NAME=piximanctl; GOOS=;        GOARCH=;      $(call BUILD_CMD)

build\:piximan\:linux:
	@echo "Building piximan for Linux"
	@CMD_NAME=piximan;    GOOS=linux;   GOARCH=386;   $(call BUILD_CMD)
	@CMD_NAME=piximan;    GOOS=linux;   GOARCH=amd64; $(call BUILD_CMD)
	@CMD_NAME=piximan;    GOOS=linux;   GOARCH=arm64; $(call BUILD_CMD)

build\:piximanctl\:linux:
	@echo "Building piximanctl for Linux"
	@CMD_NAME=piximanctl; GOOS=linux;   GOARCH=386;   $(call BUILD_CMD)
	@CMD_NAME=piximanctl; GOOS=linux;   GOARCH=amd64; $(call BUILD_CMD)
	@CMD_NAME=piximanctl; GOOS=linux;   GOARCH=arm64; $(call BUILD_CMD)

build\:piximan\:darwin:
	@echo "Building piximan for Darwin"
	@CMD_NAME=piximan;    GOOS=darwin;  GOARCH=amd64; $(call BUILD_CMD)
	@CMD_NAME=piximan;    GOOS=darwin;  GOARCH=arm64; $(call BUILD_CMD)

build\:piximanctl\:darwin:
	@echo "Building piximanctl for Darwin"
	@CMD_NAME=piximanctl; GOOS=darwin;  GOARCH=amd64; $(call BUILD_CMD)
	@CMD_NAME=piximanctl; GOOS=darwin;  GOARCH=arm64; $(call BUILD_CMD)

build\:piximan\:windows:
	@echo "Building piximan for Windows"
	@CMD_NAME=piximan;    GOOS=windows; GOARCH=386;   $(call BUILD_CMD)
	@CMD_NAME=piximan;    GOOS=windows; GOARCH=amd64; $(call BUILD_CMD)
	@CMD_NAME=piximan;    GOOS=windows; GOARCH=arm64; $(call BUILD_CMD)

build\:piximanctl\:windows:
	@echo "Building piximanctl for Windows"
	@CMD_NAME=piximanctl; GOOS=windows; GOARCH=386;   $(call BUILD_CMD)
	@CMD_NAME=piximanctl; GOOS=windows; GOARCH=amd64; $(call BUILD_CMD)
	@CMD_NAME=piximanctl; GOOS=windows; GOARCH=arm64; $(call BUILD_CMD)

build\:piximan:
	@make build:piximan:linux
	@make build:piximan:darwin
	@make build:piximan:windows

build\:piximanctl:
	@make build:piximanctl:linux
	@make build:piximanctl:darwin
	@make build:piximanctl:windows

build:
	@make build:piximan
	@make build:piximanctl

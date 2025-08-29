VERSION := $(shell if [ -n "$$PIXIMAN_VERSION" ]; then \
	echo "$$PIXIMAN_VERSION"; \
	else git describe --always --tags --dirty; \
fi)
GOFLAGS := -trimpath -mod=readonly -ldflags="-X main.version=${VERSION}"
MAKEFLAGS += --no-print-directory

define BUILD_CMD
	if [ -z "$${GOOS}" ]; then GOOS="$$(go env GOOS)"; fi; \
	if [ -z "$${GOARCH}" ]; then GOARCH="$$(go env GOARCH)"; fi; \
	BIN_DIR="bin/$${GOOS}_$${GOARCH}"; \
	mkdir -p "$${BIN_DIR}"; \
	BIN_NAME='piximan'; \
	if [ "$${GOOS}" = 'windows' ]; then BIN_NAME="$${BIN_NAME}.exe"; fi; \
	GOOS="$${GOOS}" GOARCH="$${GOARCH}" \
		go build ${GOFLAGS} -o "$${BIN_DIR}/$${BIN_NAME}" "cmd/piximan/main.go" ${ARGS}
endef

run:
	@go run ${GOFLAGS} 'cmd/piximan/main.go' ${ARGS}

build\:current:
	@echo "Building for current platform"
	@GOOS=;        GOARCH=;      $(call BUILD_CMD)

build\:linux:
	@echo "Building for Linux"
	@GOOS=linux;   GOARCH=386;   $(call BUILD_CMD)
	@GOOS=linux;   GOARCH=amd64; $(call BUILD_CMD)
	@GOOS=linux;   GOARCH=arm64; $(call BUILD_CMD)

build\:darwin:
	@echo "Building for Darwin"
	@GOOS=darwin;  GOARCH=amd64; $(call BUILD_CMD)
	@GOOS=darwin;  GOARCH=arm64; $(call BUILD_CMD)

build\:windows:
	@echo "Building for Windows"
	@GOOS=windows; GOARCH=386;   $(call BUILD_CMD)
	@GOOS=windows; GOARCH=amd64; $(call BUILD_CMD)
	@GOOS=windows; GOARCH=arm64; $(call BUILD_CMD)

build:
	@make build:linux
	@make build:darwin
	@make build:windows

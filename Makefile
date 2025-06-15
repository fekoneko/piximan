VERSION := $(shell git describe --always --tags --dirty)
VERSION_ARGS := -ldflags="-X main.version=${VERSION}"
MAKEFLAGS += --no-print-directory

define BUILD
	if [ -z "$${GOOS}" ]; then GOOS="$$(go env GOOS)"; fi; \
  if [ -z "$${GOARCH}" ]; then GOARCH="$$(go env GOARCH)"; fi; \
  BIN_DIR="bin/$${GOOS}_$${GOARCH}"; \
	mkdir -p "$${BIN_DIR}"; \
	BIN_NAME='piximan'; \
	if [ "$${GOOS}" = 'windows' ]; then BIN_NAME="$${BIN_NAME}.exe"; fi; \
	env GOOS="$${GOOS}" GOARCH="$${GOARCH}" \
		go build ${VERSION_ARGS} -v -o "$${BIN_DIR}/$${BIN_NAME}" "cmd/piximan/main.go" ${ARGS}
endef

run:
	@go run ${VERSION_ARGS} 'cmd/piximan/main.go' ${ARGS}

build\:current:
	@echo "Building for current platform"
	@GOOS=;        GOARCH=;      $(call BUILD)

build\:linux:
	@echo "Building for Linux"
	@GOOS=linux;   GOARCH=386;   $(call BUILD)
	@GOOS=linux;   GOARCH=amd64; $(call BUILD)
	@GOOS=linux;   GOARCH=arm64; $(call BUILD)

build\:darwin:
	@echo "Building for Darwin"
	@GOOS=darwin;  GOARCH=amd64; $(call BUILD)
	@GOOS=darwin;  GOARCH=arm64; $(call BUILD)

build\:windows:
	@echo "Building for Windows"
	@GOOS=windows; GOARCH=386;   $(call BUILD)
	@GOOS=windows; GOARCH=amd64; $(call BUILD)
	@GOOS=windows; GOARCH=arm64; $(call BUILD)

# TODO: test building GTK for all platforms and architectures
# TODO: migrate to meson to include blueprint nicely?
build: build\:linux build\:darwin build\:windows

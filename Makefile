# TODO: rewrite this thing to be more ideomatic to Makefile

VERSION := $(shell git describe --always --tags --dirty)
MAKEFLAGS += --no-print-directory

define COMPILE_RESOURCES
	resources_path='resources'; \
	gresource_xml_path="$$resources_path/piximan.gresource.xml"; \
	gresource_path="cmd/piximan/app/piximan.gresource"; \
	mkdir -p "$$resources_path"; \
	rm "$$resources_path"/*; \
	printf '<?xml version="1.0" encoding="UTF-8"?>\n<gresources>\n  <gresource prefix="/com/fekoneko/piximan">\n' \
		> "$$gresource_xml_path"; \
	for input_path in blueprints/*.blp; do \
		output_name="$$(basename "$$input_path" .blp).ui"; \
		echo "Compiling $$input_path to $$resources_path/$$output_name"; \
		printf '    <file preprocess="xml-stripblanks">%s</file>\n' "$$output_name" >> "$$gresource_xml_path"; \
		blueprint-compiler compile --output "$$resources_path/$$output_name" "$$input_path"; \
	done; \
	printf '  </gresource>\n</gresources>\n' >> "$$gresource_xml_path"; \
	glib-compile-resources \
		--sourcedir="$$resources_path" \
		--target="$$gresource_path" \
		"$$gresource_xml_path"
endef

define BUILD
	if [ -z "$${GOOS}" ]; then GOOS="$$(go env GOOS)"; fi; \
  if [ -z "$${GOARCH}" ]; then GOARCH="$$(go env GOARCH)"; fi; \
  BIN_DIR="bin/$${GOOS}_$${GOARCH}"; \
	mkdir -p "$${BIN_DIR}"; \
	BIN_NAME='piximan'; \
	if [ "$${GOOS}" = 'windows' ]; then BIN_NAME="$${BIN_NAME}.exe"; fi; \
	env GOOS="$${GOOS}" GOARCH="$${GOARCH}" \
		go build -ldflags='-X main.version=${VERSION}' -v -o "$${BIN_DIR}/$${BIN_NAME}" \
		"cmd/piximan/main.go" ${ARGS}
endef

prepare:
	@go mod tidy
	@$(call COMPILE_RESOURCES)

run: prepare
	@go run ${LDFLAGS_ARGS} 'cmd/piximan/main.go' ${ARGS}

build\:current: prepare
	@echo "Building for current platform"
	@GOOS=;        GOARCH=;      $(call BUILD)

build\:linux: prepare
	@echo "Building for Linux"
	@GOOS=linux;   GOARCH=386;   $(call BUILD)
	@GOOS=linux;   GOARCH=amd64; $(call BUILD)
	@GOOS=linux;   GOARCH=arm64; $(call BUILD)

build\:darwin: prepare
	@echo "Building for Darwin"
	@GOOS=darwin;  GOARCH=amd64; $(call BUILD)
	@GOOS=darwin;  GOARCH=arm64; $(call BUILD)

build\:windows: prepare
	@echo "Building for Windows"
	@GOOS=windows; GOARCH=386;   $(call BUILD)
	@GOOS=windows; GOARCH=amd64; $(call BUILD)
	@GOOS=windows; GOARCH=arm64; $(call BUILD)

# TODO: test building GTK for all platforms and architectures
build: build\:linux build\:darwin build\:windows
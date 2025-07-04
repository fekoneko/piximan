#!/usr/bin/env bash
# Usege: ./build.sh [<os>] [<arch>] 

dirname="$(dirname "$0")" || exit 1
version="$(git -C "$dirname" describe --always --tags --dirty)" || exit 1
main_path="$dirname/cmd/piximan/main.go"

# Determine target OS and architecture
os="$1";   if [ -z "$os" ];   then os="$(go env GOOS)"     || exit 1; fi
arch="$2"; if [ -z "$arch" ]; then arch="$(go env GOARCH)" || exit 1; fi

# Compile the resources
"$dirname/compile-resources.sh" || exit 1

# Determine and create directory for the binary
bin_path="$dirname/bin/${os}_$arch"
bin_name='piximan'
if [ "$os" = 'windows' ]; then bin_name="$bin_name.exe"; fi
mkdir -p "$bin_path" || exit 1

# Build the binary
echo "Building piximan $version for $os/$arch"
CGO_ENABLED=1 GOOS="$os" GOARCH="$arch" \
  go build -ldflags="-X main.version=$version" -o "$bin_path/$bin_name" \
  "$main_path" || exit 1

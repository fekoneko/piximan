#!/bin/sh
# Usage: ./build.sh [<version>] [<os>] [<arch>]

dirname="$(dirname "$0")" || exit 1
main_path="$dirname/cmd/piximan/main.go"

# Determine target OS and architecture
os="$1"; if [ -z "$os" ]; then
  os="$(go env GOOS)" || exit 1
fi
arch="$2"; if [ -z "$arch" ]; then
  arch="$(go env GOARCH)" || exit 1
fi
version="$3"; if [ -z "$version" ]; then
  version="$(git -C "$dirname" describe --always --tags --dirty)" || exit 1
fi

# Compile the resources
"$dirname/compile-resources.sh" || exit 1

# Determine and create directory for the binary
bin_path="$dirname/bin/${os}_$arch"
bin_name='piximan'
if [ "$os" = 'windows' ]; then bin_name="$bin_name.exe"; fi
mkdir -p "$bin_path" || exit 1

export GOOS="$os"
export GOARCH="$arch"
export CGO_ENABLED=1
export CGO_CPPFLAGS="$$CGO_CPPFLAGS $$CPPFLAGS"
export CGO_CFLAGS="$$CGO_CFLAGS $$CFLAGS"
export CGO_CXXFLAGS="$$CGO_CXXFLAGS $$CXXFLAGS"
export CGO_LDFLAGS="$$CGO_LDFLAGS $$LDFLAGS"
export GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw -ldflags='-X main.version=$version' $GOFLAGS"

# Build the binary
echo "Building piximan $version for $os/$arch"
go build -o "$bin_path/$bin_name" \
  "$main_path"

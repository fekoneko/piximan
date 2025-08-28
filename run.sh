#!/bin/sh
# Usage: ./run.sh [<args>]

dirname="$(dirname "$0")" || exit 1
version="$(git -C "$dirname" describe --always --tags --dirty)" || exit 1
main_path="$dirname/cmd/piximan/main.go"

# Compile the resources
"$dirname/compile-resources.sh" || exit 1

export CGO_ENABLED=1
export CGO_CPPFLAGS="$$CGO_CPPFLAGS $$CPPFLAGS"
export CGO_CFLAGS="$$CGO_CFLAGS $$CFLAGS"
export CGO_CXXFLAGS="$$CGO_CXXFLAGS $$CXXFLAGS"
export CGO_LDFLAGS="$$CGO_LDFLAGS $$LDFLAGS"
export GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw -ldflags='-X main.version=$version' $GOFLAGS"

# Run the program with Go using provided arguments
echo "Running piximan $version"
go run "$main_path" "$@"

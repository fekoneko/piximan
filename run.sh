#!/bin/sh
# Usage: ./run.sh [<args>]

dirname="$(dirname "$0")" || exit 1
version="$(git -C "$dirname" describe --always --tags --dirty)" || exit 1
main_path="$dirname/cmd/piximan/main.go"

# Compile the resources
"$dirname/compile-resources.sh" || exit 1

# Run the program with Go using provided arguments
echo "Running piximan $version"
# shellcheck disable=SC2068
go run -ldflags="-X main.version=$version" "$main_path" $@

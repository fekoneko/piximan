#!/usr/bin/env bash
set -eu -o pipefail

external_packages_path="external-packages.txt"
module_name=$(awk '/^module / {print $2; exit}' go.mod)

go list -f '{{ join .Imports "\n" }}{{ if .TestImports}}
{{ join .TestImports "\n" }}{{ end }}{{ if .XTestImports}}
{{ join .XTestImports "\n" }}{{ end }}' "./..." \
  | LC_ALL=C sort -u \
  | grep -v "$module_name" \
  | grep -Ev '^C$' \
  > "$external_packages_path"

#!/usr/bin/env bash
set -eu -o pipefail

external_imports_path="external-imports"
module_name=$(awk '/^module / {print $2; exit}' go.mod)

go list -f '{{ join .Imports "\n" }}{{ if .TestImports}}
{{ join .TestImports "\n" }}{{ end }}{{ if .XTestImports}}
{{ join .XTestImports "\n" }}{{ end }}' "./..." \
  | LC_ALL=C sort -u \
  | grep -v "$module_name" \
  | grep -Ev '^C$' \
  > "$external_imports_path"

echo "$external_imports_path"

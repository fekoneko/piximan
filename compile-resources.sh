#!/usr/bin/env bash
# Usege: ./compile-resources.sh

dirname="$(dirname "$0")" || exit 1
blueprints_path="$dirname/blueprints"
resources_path="$dirname/resources"
gresource_xml_path="$resources_path/piximan.gresource.xml"
gresource_path="$dirname/cmd/piximan/app/piximan.gresource"

# Make resources directory and cleanup
mkdir -p "$resources_path"          || exit 1
rm "$resources_path"/* 2> /dev/null || true

# Compile blueprints and produce .gresource.xml
printf '<?xml version="1.0" encoding="UTF-8"?>\n<gresources>\n  <gresource prefix="/com/fekoneko/piximan">\n' \
  > "$gresource_xml_path" || exit 1
for input_path in "$blueprints_path"/*.blp; do
  output_name="$(basename "$input_path" .blp).ui" || exit 1
  echo "Compiling $input_path to $resources_path/$output_name"
  printf '    <file preprocess="xml-stripblanks">%s</file>\n' "$output_name" \
    >> "$gresource_xml_path" || exit 1
  blueprint-compiler compile --output "$resources_path/$output_name" "$input_path" || exit 1
done
printf '  </gresource>\n</gresources>\n' >> "$gresource_xml_path" || exit 1

# Compile .gresource based on .gresource.xml
echo 'Compiling .gresource'
glib-compile-resources \
  --sourcedir="$resources_path" \
  --target="$gresource_path" \
  "$gresource_xml_path" || exit 1

#!/bin/sh
# Usage: ./compile-resources.sh

dirname="$(dirname "$0")" || exit 1
blueprints_path="$dirname/blueprints"
styles_path="$dirname/styles"
resources_path="$dirname/resources"
gresource_xml_path="$resources_path/piximan.gresource.xml"
gresource_path="$dirname/internal/resources/piximan.gresource"

# Make resources directory and cleanup
mkdir -p "$resources_path"          || exit 1
rm "$resources_path"/* 2> /dev/null || true

# Start generating .gresource.xml
printf '<?xml version="1.0" encoding="UTF-8"?>\n<gresources>\n  <gresource prefix="/com/fekoneko/piximan">\n' \
  > "$gresource_xml_path" || exit 1

# Compile blueprints and add to .gresource.xml
for input_path in "$blueprints_path"/*.blp; do
  output_name="$(basename "$input_path" .blp).ui" || exit 1
  echo "Compiling blueprint $input_path to $resources_path/$output_name"
  blueprint-compiler compile --output "$resources_path/$output_name" "$input_path" || exit 1
  printf '    <file preprocess="xml-stripblanks">%s</file>\n' "$output_name" \
    >> "$gresource_xml_path" || exit 1
done

# Copy stylesheets and add to .gresource.xml
for input_path in "$styles_path"/*.css; do
  output_name="$(basename "$input_path" .css).css" || exit 1
  echo "Compying stylesheet $input_path to $resources_path/$output_name"
  cp "$input_path" "$resources_path/$output_name" || exit 1
  printf '    <file>%s</file>\n' "$output_name" \
    >> "$gresource_xml_path" || exit 1
done

# Finish generating .gresource.xml
printf '  </gresource>\n</gresources>\n' >> "$gresource_xml_path" || exit 1

# Compile .gresource based on .gresource.xml
echo 'Compiling .gresource'
glib-compile-resources \
  --sourcedir="$resources_path" \
  --target="$gresource_path" \
  "$gresource_xml_path" || exit 1

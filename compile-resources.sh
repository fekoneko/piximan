#!/usr/bin/env bash

blueprints_path='blueprints'
resources_path='resources'
gresource_xml_path="$resources_path/piximan.gresource.xml"
gresource_path='cmd/piximan/app/piximan.gresource'

# Make resources directory / cleanup
mkdir -p "$resources_path"
rm "$resources_path"/* 2> /dev/null || true

# Compile blueprints and produce .gresource.xml
printf '<?xml version="1.0" encoding="UTF-8"?>\n<gresources>\n  <gresource prefix="/com/fekoneko/piximan">\n' \
  > "$gresource_xml_path"
for input_path in "$blueprints_path"/*.blp; do
  output_name="$(basename "$input_path" .blp).ui"
  echo "Compiling $input_path to $resources_path/$output_name"
  printf '    <file preprocess="xml-stripblanks">%s</file>\n' "$output_name" >> "$gresource_xml_path"
  blueprint-compiler compile --output "$resources_path/$output_name" "$input_path"
done
printf '  </gresource>\n</gresources>\n' >> "$gresource_xml_path"

# Compile .gresource based on .gresource.xml
glib-compile-resources \
  --sourcedir="$resources_path" \
  --target="$gresource_path" \
  "$gresource_xml_path"

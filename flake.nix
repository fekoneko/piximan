{
  description = "piximan";

  inputs = {
    nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.1";
  };

  outputs = { self, nixpkgs }:
    let
      allSystems = [
        "x86_64-linux" # 64-bit Intel/AMD Linux
        "aarch64-linux" # 64-bit ARM Linux
        "x86_64-darwin" # 64-bit Intel macOS
        "aarch64-darwin" # 64-bit ARM macOS
      ];

      forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {
        pkgs = import nixpkgs { inherit system; };
      });
    in
    {
      packages = forAllSystems ({ pkgs }: {
        default = pkgs.buildGoModule {
          name = "piximan";
          src = self;
          goSum = ./go.sum;
          subPackages = [ "cmd/piximan" ];
          vendorHash = "sha256-2aoU7xoumrq+0rQ0aIHoVLTWpJka8Q3XWuELkKAO4fc=";

          buildInputs = with pkgs; [
            gtk4
            gobject-introspection
            libadwaita
          ];
          nativeBuildInputs = with pkgs; [
            pkg-config
            blueprint-compiler
          ];

          # Compile the resources
          preBuild = ''
            resources_path='resources'
            gresource_xml_path="$resources_path/piximan.gresource.xml"
            gresource_path="cmd/piximan/app/piximan.gresource"
            mkdir -p "$resources_path"

            printf '<?xml version="1.0" encoding="UTF-8"?>\n<gresources>\n  <gresource prefix="/com/fekoneko/piximan">\n' \
              > "$gresource_xml_path"
            for input_path in blueprints/*.blp; do
              output_name="$(basename "$input_path" .blp).ui"
              echo "Compiling $input_path to $resources_path/$output_name"
              printf '    <file preprocess="xml-stripblanks">%s</file>\n' "$output_name" >> "$gresource_xml_path"
              blueprint-compiler compile --output "$resources_path/$output_name" "$input_path"
            done
            printf '  </gresource>\n</gresources>\n' >> "$gresource_xml_path"

            glib-compile-resources \
              --sourcedir="$resources_path" \
              --target="$gresource_path" \
              "$gresource_xml_path"
          '';
        };
      });
    };
}

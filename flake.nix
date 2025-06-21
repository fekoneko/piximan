{
  description = "piximan";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      src = self;
      vendorHash = "sha256-2aoU7xoumrq+0rQ0aIHoVLTWpJka8Q3XWuELkKAO4fc=";
      proxyVendor = false;

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
      packages = forAllSystems ({ pkgs }:
        let
          version = ''git describe --always --tags --dirty'';

          resources = pkgs.stdenv.mkDerivation {
            inherit src;
            name = "piximan-resources";

            # TODO: use native system package
            nativeBuildInputs = with pkgs; [
              blueprint-compiler
            ];

            buildPhase = ''bash ./compile-resources.sh'';
            installPhase = ''
              mkdir -p $out/cmd/piximan/app
              cp ./cmd/piximan/app/piximan.gresource $out/
            '';
          };
        in
        {
          default = pkgs.buildGoModule {
            inherit src vendorHash;
            name = "piximan";

            buildInputs = with pkgs; [
              gtk4
              gobject-introspection
              libadwaita
              resources
            ];

            # TODO: use native system package
            nativeBuildInputs = with pkgs; [
              pkg-config
            ];

            ldflags = [ "-X main.version=${version}" ];
            preBuild = ''cp ${resources}/piximan.gresource ./cmd/piximan/app/'';
          };
        }
      );
    };
}

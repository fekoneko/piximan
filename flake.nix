{
  description = "piximan";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    build-go-cache.url = "github:numtide/build-go-cache";
    build-go-cache.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, nixpkgs, build-go-cache }:
    let
      src = self;
      vendorHash = "sha256-2aoU7xoumrq+0rQ0aIHoVLTWpJka8Q3XWuELkKAO4fc=";
      proxyVendor = true;

      allSystems = [
        "x86_64-linux" # 64-bit Intel/AMD Linux
        "aarch64-linux" # 64-bit ARM Linux
        "x86_64-darwin" # 64-bit Intel macOS
        "aarch64-darwin" # 64-bit ARM macOS
      ];

      forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {
        pkgs = import nixpkgs { inherit system; };
        buildGoCache = import build-go-cache { inherit system; };
      });
    in
    {
      packages = forAllSystems ({ pkgs, buildGoCache }:
        let
          goCache = buildGoCache {
            inherit src vendorHash proxyVendor;
            importPackagesFile = ./external-packages.txt;
          };
        in
        {
          default = pkgs.buildGoModule {
            inherit src vendorHash proxyVendor;
            name = "piximan";
            # TODO: version
            goSum = ./go.sum;
            subPackages = [ "cmd/piximan" ];

            buildInputs = with pkgs; [
              goCache
              gtk4
              gobject-introspection
              libadwaita
            ];
            nativeBuildInputs = with pkgs; [
              pkg-config
              blueprint-compiler
            ];

            preBuild = ./compile-resources.sh;
          };
        }
      );
    };
}

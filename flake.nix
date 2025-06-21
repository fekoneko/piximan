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
      version = builtins.readFile ./.version;

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
          inherit src vendorHash version;
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
            libxml2
            blueprint-compiler
          ];

          ldflags = [ "-X main.version=${version}" ];
          preBuild = "bash ./compile-resources.sh";
        };
      }
    );
  };
}

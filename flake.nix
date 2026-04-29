{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = rec {
          default = s2h;

          s2h = pkgs.buildGoModule {
            pname = "s2h";
            version = "1.1.0";

            src = ./.;

            vendorHash = "sha256-7FNoAWa+OaRm9sLajU3WU6zJ5dJpS5NmXpbTnU4G0eE=";

            meta = {
              description = "socks2http proxy";
              mainProgram = "s2h";
            };
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.go
            pkgs.gopls
            pkgs.gotools
          ];
        };
      }
    );
}

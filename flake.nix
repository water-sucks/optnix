{
  description = "A fast options searcher for Nix module systems";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

    flake-parts.url = "github:hercules-ci/flake-parts";

    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
  };

  outputs = {
    nixpkgs,
    flake-parts,
    ...
  } @ inputs: let
    inherit (nixpkgs) lib;
  in
    flake-parts.lib.mkFlake {inherit inputs;} {
      imports = [];

      systems = lib.systems.flakeExposed;

      perSystem = {
        pkgs,
        self',
        ...
      }: let
        inherit (pkgs) callPackage go golangci-lint mkShell;
      in {
        packages = {
          default = self'.packages.optnix;

          optnix = callPackage ./package.nix {};
        };

        devShells.default = mkShell {
          name = "optnix-shell";
          nativeBuildInputs = [
            go
            golangci-lint
          ];
        };
      };
    };
}

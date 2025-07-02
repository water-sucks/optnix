{
  description = "A fast options searcher for Nix module systems";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
  };

  outputs = {
    self,
    nixpkgs,
    ...
  }: let
    inherit (nixpkgs) lib;
    eachSystem = lib.genAttrs lib.systems.flakeExposed;
  in {
    mkLib = import ./nix/lib.nix;

    packages = eachSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
      inherit (pkgs) callPackage;
    in {
      default = self.packages.${system}.optnix;
      optnix = callPackage ./package.nix {};
    });

    devShells = eachSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
      inherit (pkgs) go golangci-lint;
    in {
      default = pkgs.mkShell {
        name = "optnix-shell";
        buildInputs = [
          go
          golangci-lint
        ];
      };
    });

    nixosModules.optnix = import ./nix/modules/nixos.nix self;
    homeModules.optnix = import ./nix/modules/home-manager.nix self;
  };
}

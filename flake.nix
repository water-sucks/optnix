{
  description = "A fast options searcher for Nix module systems";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };

    nix-options-doc.url = "github:Thunderbottom/nix-options-doc/v0.2.0";
  };

  outputs = {
    self,
    nixpkgs,
    ...
  } @ inputs: let
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
      inherit (pkgs) go golangci-lint mdbook prettier scdoc;
      nix-options-doc = inputs.nix-options-doc.packages.${system}.default;
    in {
      default = pkgs.mkShell {
        name = "optnix-shell";
        buildInputs = [
          go
          golangci-lint

          mdbook
          prettier
          scdoc
          nix-options-doc
        ];
      };

      ci = pkgs.mkShell {
        name = "optnix-shell-ci";
        buildInputs = [
          go
          golangci-lint
          mdbook
          prettier
        ];
      };
    });

    nixosModules.optnix = import ./nix/modules/nixos.nix self;
    darwinModules.optnix = import ./nix/modules/nix-darwin.nix self;
    homeModules.optnix = import ./nix/modules/home-manager.nix self;

    flakeModules.flake-parts-doc = import ./flake-module.nix;
  };
}

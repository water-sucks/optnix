{pkgs ? import <nixpkgs> {}}: let
  flakeSelf = (import ./nix/flake-compat.nix).outputs;
  inherit (pkgs.stdenv.hostPlatform) system;
in {
  inherit (flakeSelf.packages.${system}) optnix;

  nixosModules.optnix = flakeSelf.nixosModules.optnix;
  darwinModules.optnix = flakeSelf.darwinModules.optnix;
  homeModules.optnix = flakeSelf.homeModules.optnix;

  mkLib = import ./nix/lib.nix;
}

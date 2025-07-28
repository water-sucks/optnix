{pkgs ? import <nixpkgs> {}}: let
  flakeSelf =
    (
      import
      (
        let
          lock = builtins.fromJSON (builtins.readFile ./flake.lock);
        in
          fetchTarball {
            url = "https://github.com/edolstra/flake-compat/archive/${lock.nodes.flake-compat.locked.rev}.tar.gz";
            sha256 = lock.nodes.flake-compat.locked.narHash;
          }
      )
      {src = ./.;}
    ).outputs;
in {
  inherit (flakeSelf.packages.${pkgs.system}) optnix;

  nixosModules.optnix = flakeSelf.nixosModules.optnix;
  darwinModules.optnix = flakeSelf.darwinModules.optnix;
  homeModules.optnix = flakeSelf.homeModules.optnix;

  mkLib = import ./nix/lib.nix;
}

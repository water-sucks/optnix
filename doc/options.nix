{pkgs ? import <nixpkgs> {}}: let
  self = (import ../nix/flake-compat.nix).outputs;

  optnixLib = self.mkLib pkgs;
in
  # Yes, we are dogfooding optnix to generate
  # its own list of options!!!
  optnixLib.mkOptionsListFromModules {
    modules = [
      self.nixosModules.optnix
    ];
    specialArgs = {
      inherit pkgs;
    };
    excluded = ["_module.args"];
  }

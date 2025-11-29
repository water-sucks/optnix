self: {
  config,
  pkgs,
  lib,
  ...
}: let
  cfg = config.programs.optnix;

  inherit (pkgs.stdenv.hostPlatform) system;

  tomlFormat = pkgs.formats.toml {};
in {
  options.programs.optnix = {
    enable = lib.mkEnableOption "CLI searcher for Nix module system options";

    package = lib.mkOption {
      type = lib.types.package;
      description = "Package that provides optnix";
      default = self.packages.${system}.optnix;
    };

    settings = lib.mkOption {
      type = lib.types.attrs;
      description = "Settings to put into optnix.toml";
      default = {};
    };
  };

  config = lib.mkIf cfg.enable {
    home.packages = [cfg.package];

    xdg.configFile = {
      "optnix/config.toml" = lib.mkIf (cfg.settings != {}) {
        source = tomlFormat.generate "config.toml" cfg.settings;
      };
    };
  };
}

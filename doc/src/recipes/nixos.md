# NixOS Recipes

`optnix` recipes for the [NixOS](https://nixos.org) module system.

## `optnix` Module Examples

A simple `nix-darwin` module showcasing the `programs.optnix` option and using
`optnixLib.mkOptionList` for option list generation:

```nix
{ options, pkgs, ... }: let
  optnixLib = inputs.optnix.mkLib pkgs;
in {
  programs.optnix = {
    enable = true;
    settings = {
      min_score = 3;

      scopes = {
        CharlesWoodson = {
          description = "NixOS configuration for CharlesWoodson";
          options-list-file = optnixLib.mkOptionsList { inherit options; };
          # For flake systems
          # evaluator = "nix eval /path/to/flake#nixosConfigurations.CharlesWoodson.config.{{ .Option }}";
          # For legacy systems
          # evaluator = "nix-instantiate --eval '<nixpkgs/nixos>' -A config.{{ .Option }}";
        };
      };
    };
  };
}
```

## Raw TOML Examples

### Flakes

Inside a flake directory `/path/to/flake` with a NixOS system named
`CharlesWoodson`:

```toml
[scopes.nixos]
description = "NixOS flake configuration for CharlesWoodson"
options-list-cmd = '''
nix eval "/path/to/flake#nixosConfigurations.CharlesWoodson" --json --apply 'input: let
  inherit (input) options pkgs;

  optionsList = builtins.filter
    (v: v.visible && !v.internal)
    (pkgs.lib.optionAttrSetToDocList options);
in
  optionsList'
'''
evaluator = "nix eval /path/to/flake#nixosConfigurations.CharlesWoodson.config.{{ .Option }}"
```

### Legacy

This uses the local NixOS system attributes located in `<nixos/nixpkgs>`.

```toml
[scopes.nixos-legacy]
description = "NixOS flake configuration on local host"
options-list-cmd = '''
nix-instantiate --eval --expr --strict --json 'let
  system = import <nixpkgs/nixos> {};
  pkgs = system.pkgs;
  optionsList = pkgs.lib.optionAttrSetToDocList system.options;
in
  optionsList'
'''
evaluator = "nix-instantiate --eval '<nixpkgs/nixos>' -A config.{{ .Option }}"
```

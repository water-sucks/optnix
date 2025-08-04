# `nix-darwin` Recipes

`optnix` recipes for the [nix-darwin](https://github.com/nix-darwin/nix-darwin)
module system.

## `optnix` Module Examples

A simple `nix-darwin` module showcasing the `programs.optnix` option and using
`optnixLib.mkOptionList` for option list generation:

```nix
{ options, pkgs, ...
}: let
  # Assume `optnix` is correctly instantiated.
  optnixLib = inputs.optnix.mkLib pkgs;
in {
  programs.optnix = {
    enable = true;
    settings = {
      min_score = 3;

      scopes = {
        TimBrown = {
          description = "nix-darwin configuration for TimBrown";
          options-list-file = optnixLib.mkOptionsList { inherit options; };
          # For flake systems
          # evaluator = "nix eval /path/to/flake#darwinConfigurations.TimBrown.config.{{ .Option }}";
          # For legacy systems
          # evaluator = "nix-instantiate --eval '<darwin-config>' -A {{ .Option }}";
        };
      };
    };
  };
}
```

## Raw TOML Examples

### Flakes

Inside a flake directory `/path/to/flake` with a `nix-darwin` system named
`TimBrown`:

```toml
[scopes.nix-darwin]
description = "nix-darwin flake configuration for TimBrown"
options-list-cmd = '''
nix eval "/path/to/flake#darwinConfigurations.TimBrown" --json --apply 'input: let
  inherit (input) options pkgs;

  optionsList = builtins.filter
    (v: v.visible && !v.internal)
    (pkgs.lib.optionAttrSetToDocList options);
in
  optionsList'
'''
evaluator = "nix eval /path/to/flake#darwinConfigurations.TimBrown.config.{{ .Option }}"
```

### Legacy

**TODO: add an example of legacy, non-flake nix-darwin configuration with raw
TOML.**

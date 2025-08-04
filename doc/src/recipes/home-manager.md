# `home-manager` Recipes

`optnix` recipes for
[`home-manager`](https://github.com/nix-community/home-manager) (HM), a popular
module system that manages user configurations.

HM has some quirks that can make it rather weird to use in `optnix`. Absolutely
prefer using the `optnix` modules + the `optnix` Nix library to generate HM
options. If using `optnix` with HM with raw TOML, you are on your own; I have
not been able to create good examples at the time of writing.

## Standalone/Inside `home-manager` Module

Inside of `home-manager`, the `options` attribute is available and can be used
directly.

```nix
{ options, config, pkgs, lib, ... }: let
  # Assume `optnix` is correctly instantiated.
  optnixLib = optnix.mkLib pkgs;
in {
  programs.optnix = {
    enable = true;
    scopes = {
      home-manager = {
        description = "home-manager configuration for all systems";
        options-list-file = optnixLib.mkOptionsList {
          inherit options;
          transform = o:
            o
            // {
              name = lib.removePrefix "home-manager.users.${config.home.username}." o.name;
            };
        };
        evaluator = "";
      };
    }
  };
}
```

**NOTE**: This may create a separate configuration file for ALL users depending
on, so it may not necessarily be what you want on multi-user systems. Look at
the next section for more on an alternative.

## NixOS/`nix-darwin` Module

`home-manager` does not expose a proper `options` attribute set on NixOS and
`nix-darwin` systems, which makes option introspection a little harder than it
should be.

Instead, an unexposed function from the type of the `home-manager.users` option
itself, `getSubOptions`, can be used to obtain an `options` attribute set for
HM.

However, this is impossible to evaluate due to the fact that it relies on other
settings that may not exist, such as usernames and such.

When this is the case, the following scope declaration using the special
function `optnixLib.mkOptionsListFromHMSource` can be used in any module: NixOS,
`nix-darwin`, or even in standalone HM.

This function is adapted from `home-manager`'s documentation facilities, and
inserts some dummy modules that allow for proper option list generation without
evaluation failing.

```nix
{
  inputs,
  config,
  pkgs,
  ...
}: let
  optnixLib = inputs.optnix.mkLib pkgs;
in {
  programs.optnix = {
    enable = true;
    settings = {
      scopes = {
        home-manager = {
          description = "home-manager options for all systems";
          options-list-file = optnixLib.hm.mkOptionsListFromHMSource {
            home-manager = inputs.home;
            modules = with inputs; [
              # Other extra modules that may exist in your source tree
              # optnix.homeModules.optnix
            ];
          };
        };
      };
    };
  };
}
```

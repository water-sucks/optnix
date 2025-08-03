# Module

There are some packaged Nix modules for easy usage that all function in the same
manner, available at the following paths (for both flake and legacy-style
configs):

- `nixosModules.optnix` :: for NixOS systems
- `darwinModules.optnix` :: for `nix-darwin` systems
- `homeModules.optnix` :: for `home-manager` systems

They all contain the same options.

## Library

When using the Nix modules, it is extremely useful to instantiate the Nix
library provided with `optnix`.

This can be done using the exported `optnix.mkLib` function:

```nix
{pkgs, ...}:
let
  # Assume `optnix` is imported already.
  optnixLib = optnix.mkLib pkgs;
in {
  programs.optnix = {
    # whatever options
  };
}
```

The functions creates option lists from Nix code ahead of time.

See the [API Reference](../reference.md) for what functions are available, as
well as the [recipes page](../recipes/index.md) for some real-life examples on
how to use the module and the corresponding functions from an instantiated
`optnix` library.

## Options

{{ #include generated-module.md }}

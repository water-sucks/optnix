# Module

There are some packaged Nix modules for easy usage that all function in the same
manner, available at the following paths (for both flake and legacy-style
configs):

- `nixosModules.optnix` :: for NixOS systems
- `darwinModules.optnix` :: for `nix-darwin` systems
- `homeModules.optnix` :: for `home-manager` systems

They all contain the same options.

{{ #include generated-module.md }}

# Installation

## Executable

The latest version of `optnix` is almost always available in `nixpkgs`.

Otherwise:

### Flakes

Use the provided flake input:

```nix
{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";

    optnix.url = "github:water-sucks/optnix";
  };

  outputs = inputs: {
    nixosConfigurations.jdoe = inputs.nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [
        ({pkgs, ...}: {
          environment.systemPackages = [
            inputs.optnix.packages.${pkgs.system}.optnix
          ];
        })
      ];
    };
  };
}
```

### Legacy

Or import it inside a Nix expression through `fetchTarball`:

```nix
{pkgs, ...}: let
  optnix-url = "https://github.com/water-sucks/optnix/archive/GITREVORBRANCHDEADBEEFDEADBEEF0000.tar.gz";
  optnix = (import "${builtins.fetchTarball optnix}").packages.${pkgs.system}.optnix;
in {
  environment.systemPackages = [
    optnix
    # ...
  ];
}
```

## When will this be in `nixpkgs`/`home-manager`/`nix-darwin`/etc.?

I'm working on it.

Ideally I'll want to get the whole project into `nix-community` too, once it
gains some popularity.

## Cache

There is a Cachix cache available. Add the following to your Nix configuration
to avoid lengthy rebuilds and fetching extra build-time dependencies, if not
using `nixpkgs`:

```nix
{
  nix.settings = {
    substituters = [ "https://watersucks.cachix.org" ];
    trusted-public-keys = [
      "watersucks.cachix.org-1:6gadPC5R8iLWQ3EUtfu3GFrVY7X6I4Fwz/ihW25Jbv8="
    ];
  };
}
```

Or if using the Cachix CLI outside a NixOS environment:

```sh
$ cachix use watersucks
```

## Extras

There are also some packaged Nix modules for easy usage that all function in the
same manner, available at the following paths:

- `nixosModules.optnix` :: for NixOS systems
- `darwinModules.optnix` :: for `nix-darwin` systems
- `homeModules.optnix` :: for `home-manager` systems

Alongside these modules is a function `mkLib`; this function is the entry point
for the `optnix` Nix library for generating lists.

More information on these resources can be found on the
[module page](./usage/module.md), as well as the [API Reference](./reference.md)
for what functions are available in the `optnix` library.

Additionally, some examples configuring `optnix` for different module systems
are available on the [recipes page](./recipes/index.md).

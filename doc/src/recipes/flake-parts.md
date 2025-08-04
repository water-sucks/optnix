# `flake-parts` Recipes

[`flake-parts`](https://flake.parts) is a framework/module system for writing
flakes using Nix module system semantics.

Singe `flake-parts` configurations can wildly vary in module selection, most
users will want to define them on a per-flake basis. This is well-supported
through the usage of a local `optnix.toml` file, relative to the flake.

## Exposing Documentation Through Flake

Use the following flake module code to expose a flake attribute called
`debug.options-doc`:

```nix
{
  lib,
  options,
  ...
}: {
  # Required for evaluating module option values.
  debug = true;
  flake = {
    debug.options-doc = builtins.unsafeDiscardStringContext
      (builtins.toJSON (lib.optionAttrSetToDocList options));
  };
}
```

OR, if you do not want to copy-paste code, use the
`optnix.flakeModules.flake-parts-doc` in an import like so:

```nix
{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";

    flake-parts.url = "github:hercules-ci/flake-parts";

    optnix.url = "github:water-sucks/optnix";
  };

  outputs = { optnix, ... }@inputs:
    flake-parts.lib.mkFlake {inherit inputs;} {
      imports = [
        optnix.flakeModules.flake-parts-doc
      ];

      # ...
    };
}
```

Then, configure an `optnix.toml` (in the same directory as the top-level
`flake.nix`) for this is rather trivial:

```toml
[scopes.flake-parts]
description = "flake-parts config for NixOS configuration"
options-list-cmd = "nix eval --json .#debug.options-doc"
evaluator = "nix eval .#debug.config.{{ .Option }}"
```

Despite the usage of `options-list-cmd`, `flake-parts` evaluates decently fast
at most times.

If using `options-list-file` is a non-negotiable, exposing a package with
`pkgs.writeText` and using the above code as a base is also possible. But you're
on your own. If you really want an example, congrats on reading this! I love
that you're taking the time to read through this application's documentation and
trying to use it, so please file an issue, I was just too lazy to do this right
now at time of writing.

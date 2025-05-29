<h1 align="center">optnix</h1>
<h6 align="center">An options searcher for Nix module systems.</h6>

## Introduction

`optnix` is a fast, terminal-based options searcher for Nix module systems.

There are multiple module systems that Nix users use on a daily basis:

- [NixOS](https://github.com/nixos/nixpkgs) (the most well-known one)
- [Home Manager](https://github.com/nix-community/home-manager)
- [`nix-darwin`](https://github.com/LnL7/nix-darwin)
- [`flake-parts`](https://github.com/hercules-ci/flake-parts)

And their documentation can be hard to look for. Not to mention, any external
options from imported modules can be impossible to find without reading source
code. `optnix` can solve that problem for you, and allows you to inspect their
values if possible; just like `nix repl` in most cases, but better.

## Concepts

### What's a module system, even?

A _module system_ is a Nix library that allows you to configure a set of exposed
_options_. All the systems mentioned above allow you to configure their
respective options with your own values.

While this can be a powerful paradigm for modeling any configuration system,
these options can be rather hard to discover. Some of these options are found
through web interfaces (like https://search.nixos.org), but many options can
remain out of sight without reading source code, such as external module options
or external module systems.

More information on how module systems work can be found on
[nix.dev](https://nix.dev/tutorials/module-system/index.html).

### Scope

A "scope", in the `optnix` context, is a combination of a module system option
list, as well as a module system instantiation.

For example, a "scope" for a NixOS configuration located at the flake ref
`github:water-sucks/nixed#nixosConfigurations.CharlesWoodson` would be:

- The module system's options attrset, converted to a list using
  `lib.optionAttrSetToDocList` (`nixosConfigurations.CharlesWoodson.options`)
- The module system's config attrset
  (`nixosConfigurations.CharlesWoodson.config`)

Scopes need to be defined from the configuration file.

### Option List

An option list is a list of JSON objects; each JSON object describes a single
option, with the following values for an example option:

```json
{
  "name": "services.nginx.enable",
  "description": "Whether to enable Nginx Web Server.",
  "type": "boolean",
  "default": {
    "_type": "literalExpression",
    "text": "false"
  },
  "example": {
    "_type": "literalExpression",
    "text": "true"
  },
  "loc": ["services", "nginx", "enable"],
  "readOnly": false,
  "declarations": [
    "/nix/store/path/nixos/modules/services/web-servers/nginx/default.nix"
  ]
}
```

Given an options attribute set, a list of these options can be generated using
[`lib.optionAttrSetToDocList`](https://noogle.dev/f/lib/optionAttrSetToDocList)
from `nixpkgs`. This will be seen in later examples.

## Installation

Use the provided flake input:

```nix
{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs";
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

## Cache

There is a Cachix cache available. Add the following to your Nix configuration
to avoid lengthy rebuilds and fetching extra build-time dependencies:

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

## Configuration

Configurations are defined in [`TOML`](https://toml.io) format, and are merged
together in order of priority.

There are four possible locations for configurations (in order of highest to
lowest priority, and only if they exist):

- Configuration paths specified on the command line with `--config`
- `optnix.toml` in the current directory
- `$XDG_CONFIG_HOME/optnix/config.toml` or `$HOME/.config/optnix/config.toml`
- `/etc/optnix/config.toml`

### Config Schema

```toml
min_score = 1 # min score required for search matches
debounce_time = 25 # debounce time for search, in ms
default_scope = "" # default scope to use if not specified on the command line

[scopes.<name>]
# Description text, for command-line completion
description = "NixOS configuration for `system`"
# A path to the options list file. Preferred over options-cmd
options-list-file = "/path/to/file"
# A command to run to generate the options list file. The list must be
# printed on stdout.
options-list-cmd = ""
# Go template for what to run in order to evaluate the option. Optional, but
# useful for previewing values.
# The template MUST contain a single placeholder of "{{ .Option }}"
# (without quotes) to fill in the option path to evaluate
# The following is only an example to show how to use the placeholder..
evaluator = "nix eval /path/to/flake#nixosConfigurations.nixos.config.{{ .Option }}"
```

## Integrations

An integration with [`nixos-cli`](https://github.com/nix-community/nixos-cli)
with zero scope configuration required is coming soon. The original code for
this application is sourced from it.

## Recipes

From the four module systems, it can be a little hard to get things up and
running without knowing how those module systems work first.

Here are a set of common scope configurations that you can use in your
`config.toml` file.

⚠️ CAUTION: Do not assume that these will automatically work with your setup.
Tweak as needed.

Feel free to contribute more examples, or request more for different module
systems, as needed.

### NixOS

NixOS, from a local flake `nixosSystem` named `CharlesWoodson`:

```toml
[scopes.nixos]
description = "NixOS flake configuration for CharlesWoodson"
# This path assumes the `nixos-cli` module is being used.
# Do NOT copy verbatim unless using it.
options-list-file = "/run/current-system/etc/nixos-cli/options-cache.json"
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

NixOS, from a legacy-style `configuration.nix` setup:

```toml
[scopes.nixos]
description = "NixOS configuration options"
options-list-file = ""
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

### nix-darwin

`nix-darwin`, from a local flake `darwinSystem` named `TimBrown`:

```toml
[scopes.nix-darwin]
description = "nix-darwin configuration for TimBrown"
options-list-cmd = '''
nix eval /path/to/flake#darwinConfigurations.TimBrown --json --apply 'input: let
  inherit (input) options pkgs;

  optionsList = builtins.filter
    (v: v.visible && !v.internal)
    (pkgs.lib.optionAttrSetToDocList options);
in
  optionsList'
'''
evaluator = "nix eval /path/to/flake#darwinConfigurations.TimBrown.config.{{ .Option }}"
```

### home-manager

Standalone `homeManagerConfiguration` flake named `MarcusAllen`:

```toml
[scopes.home-manager]
description = "Standalone Home Manager system configured through a flake"
options-list-cmd = '''
nix eval /path/to/flake#homeConfigurations.MarcusAllen --json --apply 'input: let
  inherit (input) options pkgs;

  optionsList = builtins.filter
    (v: v.visible && !v.internal)
    (pkgs.lib.optionAttrSetToDocList options);
in
  optionsList'
'''
evaluator = "nix eval /path/to/flake#homeConfigurations.MarcusAllen.config.{{ .Option }}"
```

`home-manager` does not expose a proper `options` attribute set on NixOS and
`nix-darwin` systems, which makes option introspection a little harder than it
should be.

As such, until this is supported, a common way to retrieve home-manager options
is from their own documentation file, slightly modified using `jq`.

The `home-manager` flake has a `docs-json` package that can be installed using
the `environment.systemPackages` option for this purpose. Install the package,
and then use the following configuration as an example: (this requires `jq` to
be installed):

```toml
[scopes.home-manager]
description = "Home Manager system configured as a NixOS module"
options-list-cmd = '''
jq '[to_entries[] | .value.name = .key | .value.declarations = [.value.declarations[].name] | .value]' /run/current-system/sw/share/doc/home-manager/options.json
'''
evaluator = "nix eval /path/to/flake#nixosConfigurations.JaMarcus.config.home-manager.users.varun.{{ .Option }}"
```


### flake-parts

Flake-parts configurations likely need to be defined on a per-flake basis.

It can be useful to configure a repository-specific `optnix.toml` and have
a `options-search` Nix app wrapper around `optnix`, in order to prevent
excessive command-line parameter specifications.

Import the following flake module code to expose an options-doc directly inside
the flake of choice to inspect options for:

```nix

{
  lib,
  options,
  ...
}: {
  # Required for evaluating module option values.
  debug = true;
  flake = {
    options-doc = lib.optionAttrSetToDocList options;
  };
}
```

Assuming this is implemented, an `optnix.toml` (in the same directory as the
top-level `flake.nix`) for this is rather trivial:

```toml
[scopes.flake-parts]
description = "flake-parts config for NixOS configuration"
options-list-cmd = "nix eval --json .#options-doc"
evaluator = "nix eval .#debug.config.{{ .Option }}"
```

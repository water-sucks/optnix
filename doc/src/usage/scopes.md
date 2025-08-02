# Scopes

A "scope", in the `optnix` context, is a combination of a module system option
list, as well as a module system instantiation.

For example, a "scope" for a NixOS configuration located at the flake ref
`github:water-sucks/nixed#nixosConfigurations.CharlesWoodson` would be:

- The module system's `options` set
  (`nixosConfigurations.CharlesWoodson.options`)
- The module system's `config` set (`nixosConfigurations.CharlesWoodson.config`)

`options` would be used to generate the option list, while `config` would be
used to evaluate and preview values of those options.

## Components

In `optnix`, scopes are defined from the configuration file, and have a few
components.

An example of the fields for a scope in `optnix` can be found on the
[configuration page](./configuration.md), and real configuration examples can be
found on the [recipes page](../recipes/index.md).

#### `description`

A small description of what the purpose of this scope is. Optional, but useful
for command-line completion and listing scopes more descriptively.

#### `options-list-{file,cmd}`

The option list can be specified in two different ways:

- A path to a JSON file containing the option list (`options-list-file`)
- A command that prints the option list to `stdout` (`options-list-cmd`)

**Specifying at least one of these two is mandatory for every scope.**

`options-list-file` is preferred over `options-list-cmd`, but both can be
specified; if the file does not exist or cannot be accessed/is incorrect, then
the command is used as a fallback.

The command will usually end up being some invocation of Nix, but this command
is evaluated using a shell (`/bin/sh`), which means it supports POSIX shell
constructs/available commands as long as they are in `$PATH`.

Prefer using `options-list-file` when creating configurations, since this is
almost always faster than running the equivalent `options-list-cmd`, since
`options-list-cmd` is not cached.

Generating options list files can be done using the `optnix` Nix library, and
examples can be seen on the [recipes page](../recipes/index.md).

#### `evaluator`

An **evaluator** is a command template that can be used to evaluate a Nix
configuration to retrieve values.

It is a shell command (always some invocation of a Nix command), but with a
twist: exactly one placeholder of `{{ .Option }}` for the option to evaluate is
required. This will be filled in with the option to evaluate.

An example evaluator for a Nix flake would be:

```sh
nix eval "/path/to/flake#nixosConfigurations.nixos.config.{{ .Option }}"
```

Specifying an evaluator for a scope is optional.

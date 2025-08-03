# Configuration

Configurations are defined in [`TOML`](https://toml.io) format, and are merged
together in order of priority.

There are four possible locations for configurations (in order of highest to
lowest priority, and only if they exist):

- Configuration paths specified on the command line with `--config`
- `optnix.toml` in the current directory
- `$XDG_CONFIG_HOME/optnix/config.toml` or `$HOME/.config/optnix/config.toml`
- `/etc/optnix/config.toml`

### Schema

A more in-depth explanation for scope configuration values is on the
[scopes page](./scopes.md).

```toml
# Minimum score required for search matches; the higher the score,
# the more fuzzy matches are required before it is displayed in the list.
# Higher scores will generally lead to less (but more relevant) results, but
# this has diminishing returns.
min_score = 1
# Debounce time for search, in ms
debounce_time = 25
# Default scope to use if not specified on the command line
default_scope = ""
# Formatter command to use for evaluated values, if available. Takes input on
# stdin and outputs the formatted code back to stdout.
formatter_cmd = "nixfmt"

# <name> is a placeholder for the name of the scope.
# This is not a working scope! See the recipes page.
# for real examples.
[scopes.<name>]
# Description of this scope
description = "NixOS configuration for `nixos` system"
# A path to the options list file. Preferred over options-list-cmd.
options-list-file = "/path/to/file"
# A command to run to generate the options list file. The list must be
# printed on stdout.
# Check the recipes page for some example commands that can generate this.
# This is not always needed for every scope.
# The following is only an example.
options-list-cmd = """
nix eval /path/to/flake#homeConfigurations.nixos --json --apply 'input: let
  inherit (input) options pkgs;

  optionsList = builtins.filter
    (v: v.visible && !v.internal)
    (pkgs.lib.optionAttrSetToDocList options);
in
  optionsList'
"""
# Go template for what to run in order to evaluate the option. Optional, but
# useful for previewing values.
# Check the scopes page for an explanation of this value.
evaluator = "nix eval /path/to/flake#nixosConfigurations.nixos.config.{{ .Option }}"
```

# API Reference

## `mkLib`

`mkLib :: pkgs -> set`

Creates an `optnixLib` instance using an instantiated `pkgs` set.

## `optnixLib`

`optnixLib` represents an instantiated value crated by `mkLib` above.

### `optnixLib.mkOptionsList`

`mkOptionsList :: AttrSet -> Derivation`

Generates an `options.json` derivation that contains a JSON options list from a
provided options attribute set that can be linked into various Nix module
systems.

A base (but common pattern) of usage:

```nix
{options, pkgs, ...}: let
  # Adapt this to wherever `optnix` came from
  optnixLib = optnix.mkLib pkgs;
  optionsList = mkOptionsList {
    inherit options;
    excluded = []; # Add problematic eval'd paths here
  };
in {
  # use here...
}
```

Arguments are passed as a single attribute set in a named manner.

#### Required Attributes

- `options :: AttrSet` :: The options attrset to generate an option list for

#### Optional Attributes

- `transform :: AttrSet -> AttrSet` :: A function to apply to each generated
  option in the list, useful for stripping prefixes, where `AttrSet` is a single
  option in the list
- `excluded :: [String]` :: A list of dot-paths in the options attribute set to
  exclude from the generated options list using `optnixLib.removeNestedAttrs`

### `optnixLib.mkOptionsListFromModules`

`mkOptionsListFromModules :: AttrSet -> Derivation`

Generates an `options.json` derivation that contains a JSON options list from a
provided list of modules that can be linked into various Nix module systems.

This can be useful when generating documentation for external modules that are
not necessarily part of the configuration, such as generating option lists for
usage outside of `optnix`, or for generating a system-agnostic documentation
list.

Arguments are passed as a single attribute set in a named manner.

#### Arguments

- `modules :: List[Module]` :: A list of modules to generate an option list for

### `optnixLib.combineLists`

`combineLists :: [Derivation] -> Derivation`

Combines together multiple JSON file derivations containing option lists (such
as those created by `mkOptionsList`) into a single JSON file.

`combineLists [optionsList1 optionsList2]`

Internally uses `jq --slurp` add to merge JSON arrays.

#### Arguments

- `[Derivation]` :: A list of JSON file derivations, each containing an option
  list, preferably created using `mkOptionsList`

### `optnixLib.removeAtPath`

`removeAtPath :: String -> AttrSet -> AttrSet`

Recursively remove a nested attribute from an attrset, following a string array
containing the path to remove.

This can be useful for removing problematic options (i.e. ones that fail
evaluation) in a custom manner, and is also used internally by the `excluded`
parameter of `mkOptionsList`.

`removeAtPath "services.nginx" options`

This example will return the same attrset config, but with the `services.nginx`
subtree removed. If a path does not exist or is not an attrset, it is left
untouched.

#### Arguments

- `String` :: A dot-path string representing the attribute path to remove from
  the attrset
- `AttrSet` :: The attrset to remove attributes recursively from

### `optnixLib.removeNestedAttrs`

`removeNestedAttrs :: [String] -> AttrSet -> AttrSet`

Recursively remove multiple attributes using `optnixLib.removeAtPath`.

`removeNestedAttrs ["programs.chromium" "services.emacs"] options`

This example will return the same attrset config, but with the `services.emacs`
and `programs.chromium` subtrees removed. If a path does not exist or is not an
attrset, it is left untouched.

#### Arguments

- `[String]` :: A list of dot-path strings representing the attribute paths to
  remove from the attrset
- `AttrSet` :: The attrset to remove attributes recursively from

### `optnixLib.hm`

Functions for `optnixLib` specifically related to
[`home-manager`](https://github.com/nix-community/home-manager).

### `optnixLib.hm.mkOptionsListFromHMSource`

`mkOptionsListFromHMSource :: AttrSet -> Derivation`

Generate a JSON options list given a `home-manager` source tree, alongside
additional `home-manager` modules if desired.

This implementation is adapted from Home Manager's internal documentation
generation pipeline and simplified for `optnix` usage. Prefer using
`mkOptionsList` with an explicit instantiated `options` attribute set if
possible.

#### Required Attributes

- `home-manager :: Derivation` :: The derivation containing `home-manager`,
  usually sourced using `fetchTarball`, a Nix channel, or a flake input

#### Optional Attributes

- `modules :: List[Module]` :: A list of extra modules to additionally evaluate
  when generating the option list

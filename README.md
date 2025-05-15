<h1 align="center">optnix</h1>
<h6 align="center">An options searcher for Nix module systems.</h6>

## Introduction

`optnix` is a tool to search for Nix options in a module system through the
terminal.

There are multiple module systems that Nix users use on a daily basis:

- [NixOS](https://github.com/nixos/nixpkgs) (the most well-known one)
- [Home Manager](https://github.com/nix-community/home-manager)
- [`nix-darwin`](https://github.com/LnL7/nix-darwin)
- [`flake-parts`](https://github.com/hercules-ci/flake-parts)

And their documentation can be hard to look for. Not to mention, any external
options from imported modules can be impossible to find without reading source
code. `optnix` can solve that problem for you, and also allows you to inspect
their values, just like `nix repl`, but better.

## What's a module system, even?

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

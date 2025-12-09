<h1 align="center">optnix</h1>
<h6 align="center">An options searcher for Nix module systems.</h6>

## Introduction

`optnix` is a fast, terminal-based options searcher for Nix module systems.

[![a demo of optnix](./doc/src/demo.gif)](https://asciinema.org/a/732922?autoplay=1)

There are multiple module systems that Nix users use on a daily basis:

- [NixOS](https://github.com/nixos/nixpkgs) (the most well-known one)
- [Home Manager](https://github.com/nix-community/home-manager)
- [`nix-darwin`](https://github.com/LnL7/nix-darwin)
- [`flake-parts`](https://github.com/hercules-ci/flake-parts)

And their documentation can be hard to look for. Not to mention, any external
options from imported modules can be impossible to find without reading source
code. `optnix` can solve that problem for you, and allows you to inspect their
values if possible; just like `nix repl` in most cases, but better.

There is a website for high-level documentation available at
https://water-sucks.github.io/optnix.

This repository is hosted on [sr.ht](https://sr.ht/~watersucks/optnix), with an
official mirror on [GitHub](https://github.com/water-sucks/optnix).

## Install

`optnix` is available in `nixpkgs`.

More installation instructions can be found on the
[website](https://water-sucks.github.io/optnix/installation.html).

## Integrations

`optnix` can be used as a Go library, and is used as such in the following
applications.

### [`nixos-cli`](https://github.com/nix-community/nixos-cli)

`optnix` is used as a library in the `nixos option` subcommand of `nixos-cli`.

`nixos option` requires zero configuration for discovery of NixOS options.

## Contributing

Prefer emailing patch sets to the
[official development mailing list](mailto:~watersucks/optnix-devel@lists.sr.ht).

While the official repository is located on
[sr.ht](https://git.sr.ht/~watersucks/optnix), contributions are also accepted
through GitHub using the
[official mirror](https://github.com/water-sucks/optnix), if desired.

Additionally, filing GitHub issues is fine, but consider using the official
issue tracker on [sr.ht](https://todo.sr.ht/~watersucks/optnix). All issues from
GitHub will be mirrored there by me anyway.

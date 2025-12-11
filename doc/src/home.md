# `optnix`

`optnix` is a fast, terminal-based options searcher for Nix module systems.

<script src="https://asciinema.org/a/761292.js" id="asciicast-761292" data-autoplay async="true"></script>

There are multiple module systems that Nix users use on a daily basis:

- [NixOS](https://github.com/nixos/nixpkgs) (the most well-known one)
- [Home Manager](https://github.com/nix-community/home-manager)
- [`nix-darwin`](https://github.com/nix-darwin/nix-darwin)
- [`flake-parts`](https://github.com/hercules-ci/flake-parts)

These systems can have difficult-to-navigate documentation, especially for
options in external modules.

`optnix` solves that problem, and lets users inspect option values if possible;
just like `nix repl` in most cases, but prettier.

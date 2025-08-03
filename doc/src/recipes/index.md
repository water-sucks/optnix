# Recipes

There are four popular sets of Nix module systems:

- [NixOS](https://nixos.org)
- [`nix-darwin`](https://github.com/nix-darwin/nix-darwin)
- [`home-manager`](https://github.com/nix-community/home-manager)
- [`flake-parts`](https://flake.parts)

Even despite their popularity, it can be a little hard to get things up and
running without knowing how those module systems work first.

The following sections are a set of common scope configurations that you can use
in your configurations, specified in both Nix form using `optnix.mkLib` and raw
TOML, whenever necessary.

Make sure to look at the [API Reference](../reference.md) to check what
arguments can be passed to `optnix` library functions.

⚠️ **CAUTION: Do not assume that these will automatically work with your setup.
Tweak as needed.**

Feel free to contribute more examples, or request more for different module
systems, as needed.

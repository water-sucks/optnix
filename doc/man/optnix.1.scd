OPTNIX(1)

# NAME

optnix - a fast options searcher for NixOS module systems

# SYNOPSIS

*optnix* [options] [OPTION-NAME]

# DESCRIPTION

There are multiple module systems that Nix users use on a daily basis:

- NixOS (the most well-known one)
- Home Manager
- nix-darwin
- flake-parts

These systems can have difficult-to-navigate documentation, especially for
options in external modules.

*optnix* solves that problem, and lets users inspect option values if possible;
just like *nix repl* in most cases, but prettier.

# EXAMPLES

Search for options interactively using the _nixos_ scope:

	*optnix -s nixos*

Search for the option _programs.zsh.enable_ in interactive mode using the
default scope:

	*optnix programs.zsh.enable*

Search for the option _programs.zsh.enable_ in non-interactive mode in the
default scope:

	*optnix -n programs.zsh.enable*

Display the value of _programs.zsh.enable_ in the _home-manager_ scope:

	*optnix -v -s home-manager programs.zsh.enable*

Display the value of _flake.nixosConfigurations_ in the _flake-parts_ scope
using a custom configuration file:

	*optnix -c ./contrib/optnix.toml -v -s flake-parts flake.apps*

# ARGUMENTS

*OPTION-NAME*
	Specify the name of an option to search for.

	In non-interactive mode, this argument is required, and it *MUST* evaluate
	to a valid option.

	In interactive mode, it serves as an initial input for the search bar, and
	is not required.

# OPTIONS

*-c*, *--config <FILES>*
	Path to extra configuration file(s) to load.

	To specify multiple extra configuration files to load, pass this option
	multiple times.

*-j*, *--json*
	Output information in JSON format.

	Implies non-interactive mode.

*-l*, *--list-scopes*
	List available scopes and their origins (aka what configuration file they
	most recently were modified in) and exit.

*-m*, *--min-score <SCORE>*
	Minimum score threshold for deeming a potential candidate a match.

	The lower the number is, the more search results will show up, but they may
	be less relevant.

*-n*, *--non-interactive*
	Do not show search TUI for options.

	If specified, *OPTION-NAME* will become a mandatory parameter.

*-s*, *--scope <NAME>*
	Scope name to use.

	If a default scope is not defined in the configuration, this parameter is
	required.

*-v*, *--value-only*
	Only show the value of the passed option; useful for scripting.

	Implies non-interactive mode.

*--version*
	Display version information.

*-h*, *--help*
	Show the help message for this command.

# AUTHORS

Maintained by Varun Narravula <varun@snare.dev>. Up-to-date sources can be
found at https://github.com/water-suck/optnix, and bugs reports or patches
can be submitted to GitHub's issue tracker.

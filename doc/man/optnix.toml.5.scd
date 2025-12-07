OPTNIX-TOML(5)

# NAME

optnix.toml - configuration options for optnix

# DESCRIPTION

This man page documents available configuration keys for the *optnix* command.

Configuration files use the TOML format, and are merged together in the
following order (if they exist):

- Configuration paths specified on the command line with _--config <FILE>_
- _optnix.toml_ in the current directory
- _$XDG_CONFIG_HOME/optnix/config.toml_ or _$HOME/.config/optnix/config.toml_
- _/etc/optnix/config.toml_

Defaults are noted alongside each option below.

For more information about the *optnix* command itself, check the *optnix(1)*
man page.

## TERMINOLOGY

*Option List*

An _option list_ is a JSON-formatted list that has the following schema:

```
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
		"/nix/store/path/to/nginx/default.nix"
	]
}
```

Generating this list can be done with the _optnix_ flake library, or by using
the _nixpkgs_ library function _lib.optionAttrSetToDocList_.

See the documentation website for information on common configurations and
setups to generate option lists.

*Scope*

A "scope", in its most basic form, consists of four elements:

- Name
- Description
- Option list
- Evaluator

*optnix* operates on one scope at any given time. It will search the option list
for available options, and evaluate options using the "evaluator", which
is a small template for a _nix eval_ or _nix-instantiate_ command to run to
evaluate a value.

The only elements that are mandatory for a scope are a name and a an option
list.

# OPTIONS

All available settings and their descriptions.

*min_score*

Minimum score required for fuzzy search matches; the higher the score, the more
fuzzy matches are required before it is displayed in the list. Higher scores
will generally lead to less (but more relevant) results, but this has
diminishing returns.

This option is not used in regex search mode.

Default: _1_


*debounce_time*

Debounce time for search, in milliseconds.

Default: _25_


*default_scope*

Default scope to use if not specified on the command line.

Default: _(none)_


*formatter_cmd*

Formatter command to use for evaluated values, if available. Takes input on
stdin and outputs the formatted code back to stdout.

Default: _nixfmt_


*scopes.<name>*

Scopes, specified as a map. Each scope will have a unique name.

Default: _{} (none)_


*scopes.<name>.description*

A small description of what the purpose of this scope is. Optional, but useful
for command-line completion and listing scopes more descriptively.

Default: _(none)_


*scopes.<name>.options-list-file*

A JSON file containing an option list. This is preferred over
_scopes.<name>.options-list-cmd_ if it exists.

Default: _(none)_


*scopes.<name>.options-list-cmd*

A command to evaluate that produces a JSON-formatted option list on _stdout_.
This is a fallback if _scopes.<name>.options-list-file_ does not exist or fails.

Default: _(none)_


*scopes.<name>.evaluator*

A command template that can be used to evaluate a Nix configuration to retrieve
values.

Requires a single placeholder called _{{ .Option }}_ to be present; this is
filled in with the option to evaluate automatically.

Default: _(none)_

# SEE ALSO

*optnix(1)*

# AUTHORS

Maintained by the *optnix* authors. See the main man page *optnix(1)* for
details.

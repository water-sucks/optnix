# Usage

The next few sections describe the different `optnix` concepts.

But first:

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

### How `optnix` Works

`optnix` works by ingesting **option lists** that are generated from these
module systems.

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
from `nixpkgs`, or by using wrapper functions provided with `optnix` as a Nix
library. This will be seen in later examples.

### Operation

There are two modes of operation: interactive (the default) and non-interactive.

Interactive mode will display a search UI that allows looking for options using
fuzzy search keywords or regular expressions. Selected options in the list can
also be evaluated in order to preview their values.

Non-interactive mode requires a valid option name as input, and will display the
option and its values (if applicable) without any user interaction. This kind of
output is useful when an option name is known, such as for scripting.

`optnix` is controlled through its configuration file (or files) that define
"**scopes**". For more, look at the following pages:

- [**Scopes**](./scopes.md)
- [**Configuration**](./configuration.md.md)

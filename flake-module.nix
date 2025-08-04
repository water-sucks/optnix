{
  lib,
  options,
  ...
}: {
  # Required for evaluating module option values.
  debug = true;
  flake = {
    options-doc = lib.optionAttrSetToDocList options;
  };
}

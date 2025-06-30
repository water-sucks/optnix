{lib}: let
  parsePath = pathStr: lib.splitString "." pathStr;

  removeNestedAttrs = paths: set: lib.foldl' (s: p: removeAtPath (parsePath p) s) set paths;

  removeAtPath = path: set:
    if path == []
    then set
    else let
      key = builtins.head path;
      rest = builtins.tail path;
      sub = set.${key} or null;
    in
      if rest == []
      then builtins.removeAttrs set [key]
      else if builtins.isAttrs sub
      then
        set
        // {
          ${key} = removeAtPath rest sub;
        }
      else set;

  /*
  Create an options list JSON file from an options attribute set.

  @param  options   options attribute set to generate options list from
  @param  excluded  paths to options or sets of options to skip over
  @return           a derivation that builds an options.json file
  */
  mkOptionsList = {
    options,
    pkgs,
    excluded ? [],
  }: let
    options' = removeAtPath excluded options;

    rawOptions =
      lib.optionAttrSetToDocList options';
    filteredOptions = lib.filter (opt: opt.visible && !opt.internal) rawOptions;

    optionsJSON = builtins.unsafeDiscardStringContext (builtins.toJSON filteredOptions);
  in
    pkgs.writeText "options.json" optionsJSON;

  /*
  Create an options list JSON string from a list of modules.

  @param  modules   list of modules containing options
  @return           a derivation that builds an options.json file
  */
  mkOptionsListFromModules = {
    modules,
    pkgs,
  }: let
    eval'd = lib.evalModules {
      modules =
        modules
        ++ [
          {_module.check = false;}
        ];
    };
  in
    mkOptionsList {
      inherit pkgs;
      inherit (eval'd) options;
    };
in {
  inherit
    removeAtPath
    removeNestedAttrs
    mkOptionsList
    mkOptionsListFromModules
    ;
}

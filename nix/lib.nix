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
  @param  transform function to apply to each option
  @return           a derivation that builds an options.json file
  */
  mkOptionsList = {
    options,
    pkgs,
    transform ? lib.id,
    excluded ? [],
  }: let
    options' = removeAtPath excluded options;

    rawOptions =
      map transform (lib.optionAttrSetToDocList options');
    filteredOptions = lib.filter (opt: opt.visible && !opt.internal) rawOptions;

    optionsJSON = builtins.unsafeDiscardStringContext (builtins.toJSON filteredOptions);
  in
    pkgs.writeText "options.json" optionsJSON;

  /*
  Create an options list JSON file from a list of modules.

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

  hm = {
    /*
    Create an options list JSON file from a Home Manager source list of modules.

    @param  home-manager  path to a home-manager source (like from a flake input or tarball)
    @param  modules       list of extra modules containing options to include
    @return               a derivation that builds an options.json file
    */
    mkOptionsListFromHMSource = {
      pkgs,
      home-manager,
      modules ? [],
    }: let
      hmLib = import "${home-manager}/modules/lib/stdlib-extended.nix" lib;

      hmModules = import "${home-manager}/modules/modules.nix" {
        inherit pkgs;
        lib = hmLib;
        check = false;
      };

      scrubDerivations = prefixPath: attrs: let
        scrubDerivation = name: value: let
          pkgAttrName = prefixPath + "." + name;
        in
          if lib.isAttrs value
          then
            scrubDerivations pkgAttrName value
            // lib.optionalAttrs (lib.isDerivation value) {
              outPath = "\${${pkgAttrName}}";
            }
          else value;
      in
        lib.mapAttrs scrubDerivation attrs;

      # Make sure the used package is scrubbed to avoid actually
      # instantiating derivations.
      scrubbedPkgsModule = {
        imports = [
          {
            _module.args = {
              pkgs = lib.mkForce (scrubDerivations "pkgs" pkgs);
              pkgs_i686 = lib.mkForce {};
            };
          }
        ];
      };

      allModules =
        hmModules
        ++ modules
        ++ [
          scrubbedPkgsModule
        ];

      options =
        (hmLib.evalModules {
          modules = allModules;
          class = "homeManager";
        }).options;
    in
      mkOptionsList {
        inherit options pkgs;
      };
  };
in {
  inherit
    removeAtPath
    removeNestedAttrs
    mkOptionsList
    mkOptionsListFromModules
    hm
    ;
}

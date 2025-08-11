{
  lib,
  stdenv,
  installShellFiles,
  buildGoModule,
  nix-gitignore,
  scdoc,
}:
buildGoModule (finalAttrs: {
  pname = "optnix";
  version = "0.3.0-dev";
  src = nix-gitignore.gitignoreSource [] ./.;

  vendorHash = "sha256-g/H91PiHWSRRQOkaobw2wAYX/07DFxWTCTlKzf7BT1Y=";

  nativeBuildInputs = [installShellFiles scdoc];

  env = {
    CGO_ENABLED = 0;
    VERSION = finalAttrs.version;
  };

  buildPhase = ''
    runHook preBuild
    make all man
    runHook postBuild
  '';

  installPhase = ''
    runHook preInstall

    install -Dm755 ./optnix -t $out/bin

    install -Dm755 ./optnix.1 -t $out/share/man/man1
    install -Dm755 ./optnix.toml.5 -t $out/share/man/man5

    runHook postInstall
  '';

  postInstall = lib.optionalString (stdenv.buildPlatform.canExecute stdenv.hostPlatform) ''
    installShellCompletion --cmd optnix \
      --bash <($out/bin/optnix --completion bash) \
      --fish <($out/bin/optnix --completion fish) \
      --zsh <($out/bin/optnix --completion zsh)
  '';

  meta = {
    homepage = "https://github.com/water-sucks/optnix";
    description = "A fast options searcher for Nix module systems";
    license = lib.licenses.gpl3Only;
    maintainers = with lib.maintainers; [water-sucks];
  };
})

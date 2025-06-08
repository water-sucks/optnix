{
  lib,
  stdenv,
  installShellFiles,
  buildGoModule,
  nix-gitignore,
}:
buildGoModule (finalAttrs: {
  pname = "optnix";
  version = "0.1.2";
  src = nix-gitignore.gitignoreSource [] ./.;

  vendorHash = "sha256-/rV21mX6VrJj39M6dBw4ubp6+O47hxeLn0ZcsG6Ujno=";

  nativeBuildInputs = [installShellFiles];

  env = {
    CGO_ENABLED = 0;
    VERSION = finalAttrs.version;
  };

  buildPhase = ''
    runHook preBuild
    make all
    runHook postBuild
  '';

  installPhase = ''
    runHook preInstall

    install -Dm755 ./optnix -t $out/bin

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

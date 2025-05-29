{
  lib,
  stdenv,
  installShellFiles,
  buildGoModule,
  nix-gitignore,
}:
buildGoModule (finalAttrs: {
  pname = "optnix";
  version = "0.1.0";
  src = nix-gitignore.gitignoreSource [] ./.;

  vendorHash = "sha256-+s+J1vi69riJWX/wf8xMOAihvUlU80aOXqsOfhQkv4c=";

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

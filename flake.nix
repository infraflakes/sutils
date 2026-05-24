{
  description = "Sane Utils is an opinionated CLI suite to streamline many command line work.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs @ {
    self,
    nixpkgs,
    flake-parts,
    ...
  }:
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"];

      perSystem = {
        config,
        pkgs,
        lib,
        system,
        ...
      }: let
        # Read the version safely, stripping trailing newlines common in version files
        version = let
          versionFile = ./. + "/.version";
        in
          if builtins.pathExists versionFile
          then lib.strings.trim (builtins.readFile versionFile)
          else self.shortRev or "dev";

        # Clean source to ensure we include hidden version files but filter out junk
        cleanedSource = lib.cleanSourceWith {
          src = ./.;
          filter = path: type: let
            baseName = baseNameOf path;
          in
            baseName == ".version" || lib.cleanSourceFilter path type;
        };
      in {
        # Define your package build directly inside perSystem
        packages.default = pkgs.buildGoModule {
          pname = "sutils";
          inherit version;
          src = cleanedSource;

          vendorHash = "sha256-G/kRteKbu1TsvEYAvAGBRMLhYLUEY4ham/PV9eJKvLs=";

          env.CGO_ENABLED = "0";

          ldflags = [
            "-s"
            "-w"
            "-X main.version=${version}"
          ];

          nativeBuildInputs = [pkgs.installShellFiles];

          postInstall = ''
            mv $out/bin/sutils $out/bin/sn
          '';

          postFixup = ''
            installShellCompletion --fish ${cleanedSource}/completions/sn.fish
            installShellCompletion --zsh ${cleanedSource}/completions/sn.zsh
            installShellCompletion --bash ${cleanedSource}/completions/sn.bash
          '';
        };

        # Standard interactive shell for developers working on the source
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            go
            golangci-lint
            goreleaser
          ];

          shellHook = ''
            export GOPATH="$HOME/.local/share/go"
            export GOBIN="$GOPATH/bin"
            export PATH="$PATH:$GOBIN"
          '';
        };
      };
    };
}

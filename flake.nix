{
  description = "Sane Utils is an opinionated CLI suite to streamline many command line work.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};

        buildSutils = {
          src,
          version,
        }:
          pkgs.buildGoModule {
            pname = "sutils";
            inherit version src;
            preBuild = ''
              export CGO_ENABLED=0
            '';
            vendorHash = "sha256-gJhMKUmy/wwlQ9uiiab74hdl5O5w5B/O3k6RuQMPDbo=";
            ldflags = [
              "-s"
              "-w"
              "-X main.version=${version}"
            ];
            nativeBuildInputs = [pkgs.installShellFiles];
            postInstall = ''
              mv $out/bin/sutils $out/bin/sn
            '';
            # postFixup = ''
            #   installShellCompletion --fish ${src}/completions/sn.fish
            #   installShellCompletion --zsh ${src}/completions/sn.zsh
            #   installShellCompletion --bash ${src}/completions/sn.bash
            # '';
          };

        cleanedSource = pkgs.lib.cleanSourceWith {
          src = ./.;
          filter = path: type: let
            baseName = baseNameOf path;
          in
            baseName == ".version" || pkgs.lib.cleanSourceFilter path type;
        };
      in {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            golangci-lint
            cmake
            goreleaser
          ];
        };

        packages.default = buildSutils {
          src = cleanedSource;
          version = let
            versionFile = "${cleanedSource}/.version";
          in
            pkgs.lib.escapeShellArg (
              if builtins.pathExists versionFile
              then builtins.readFile versionFile
              else self.shortRev or "dev"
            );
        };

        apps.default = flake-utils.lib.mkApp {drv = self.packages.${system}.default;};
      }
    );
}

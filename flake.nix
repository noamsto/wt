{
  description = "wt - Git worktree manager";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"];

      perSystem = {pkgs, ...}: {
        packages.default = pkgs.buildGoModule {
          pname = "wt";
          version = "0.1.0";
          src = ./.;
          vendorHash = "sha256-omCeza0ltzWu0RulQGcc6hQ4lT5gH0hh404vrfZjda8=";
          ldflags = ["-s" "-w"];
          subPackages = ["cmd/wt"];
          meta = {
            description = "Git worktree manager with optional tmux integration";
            mainProgram = "wt";
          };
        };
      };

      flake = {
        homeManagerModules.default = {
          config,
          lib,
          pkgs,
          ...
        }: let
          cfg = config.programs.wt;
          wtPkg = inputs.self.packages.${pkgs.system}.default;
        in {
          options.programs.wt = {
            enable = lib.mkEnableOption "wt - git worktree manager";
          };

          config = lib.mkIf cfg.enable {
            home.packages = [wtPkg];

            xdg.configFile."fish/completions/wt.fish".source = ./completions/wt.fish;
          };
        };
      };
    };
}

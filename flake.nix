{
  description = "School Lunch Menu Calendar";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in {
        packages.default = pkgs.buildGoModule {
          pname = "calendar-app";
          version = "1.0.0";
          src = ./.;

          vendorHash = "sha256-AFv5Hp+gmRuetMzRXm1m0AGLNduZtqU8lhiC07GQ9dA=";

          nativeBuildInputs = with pkgs; [ gcc ];

          ldflags = [
            "-s"
            "-w"
          ];

          meta = with pkgs.lib; {
            description = "School Lunch Menu Calendar - web app for managing school meal schedules";
            license = licenses.mit;
            mainProgram = "calendar-app";
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gcc
          ];
        };
      }
    );
}

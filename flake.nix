{
  description = "School Lunch Menu Calendar";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" ];
    in
    flake-utils.lib.eachSystem supportedSystems (system:
      let
        pkgs = import nixpkgs { inherit system; };

        src = pkgs.lib.sources.cleanSourceWith {
          src = pkgs.lib.sources.cleanSource ./.;
          filter = name: type:
            let baseName = baseNameOf (toString name); in
            !(builtins.elem baseName [
              "data" "result" ".git" "calendar-app" "hash-password"
            ]);
        };

        calendar-app = pkgs.buildGoModule {
          pname = "calendar-app";
          version = "1.0.0";
          inherit src;

          vendorHash = "sha256-eMhl9Pc7Ip1X8Mltbe8BaPcpTcv3Ud9ccvsS6Tc1gLc=";

          nativeBuildInputs = with pkgs; [ gcc ];

          ldflags = [ "-s" "-w" ];

          meta = with pkgs.lib; {
            description = "School Lunch Menu Calendar - web app for managing school meal schedules";
            license = licenses.mit;
            mainProgram = "calendar-app";
          };
        };
      in {
        packages = {
          default = calendar-app;
          calendar-app = calendar-app;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [ go gcc ];
        };
      }
    ) // {
      overlays.default = final: prev: {
        calendar-app = self.packages.${prev.system}.calendar-app;
      };

      nixosModules.default = { config, lib, pkgs, ... }:
        let
          cfg = config.services.school-lunch-calendar;
          pkg = cfg.package;
        in {
          options.services.school-lunch-calendar = {
            enable = lib.mkEnableOption "School Lunch Menu Calendar";

            package = lib.mkOption {
              type = lib.types.package;
              default = self.packages.${pkgs.system}.calendar-app;
              description = "The calendar-app package to use.";
            };

            port = lib.mkOption {
              type = lib.types.port;
              default = 8080;
              description = "Port to listen on.";
            };

            dataDir = lib.mkOption {
              type = lib.types.path;
              default = "/var/lib/school-lunch-calendar";
              description = "Directory for the SQLite database and uploaded logos.";
            };

            passwordFile = lib.mkOption {
              type = lib.types.path;
              description = ''
                Path to a file containing the plaintext password.
                Typically provided by sops-nix, e.g. config.sops.secrets.school-lunch-password.path.
                The password is bcrypt-hashed at startup; the plaintext is never written to disk.
              '';
            };

            sessionSecretFile = lib.mkOption {
              type = lib.types.nullOr lib.types.path;
              default = null;
              description = ''
                Path to a file containing the session secret string.
                If null, a random secret is generated and persisted in dataDir on first start.
                Typically provided by sops-nix.
              '';
            };

            openFirewall = lib.mkOption {
              type = lib.types.bool;
              default = false;
              description = "Whether to open the configured port in the firewall.";
            };
          };

          config = lib.mkIf cfg.enable {
            systemd.services.school-lunch-calendar = {
              description = "School Lunch Menu Calendar";
              wantedBy = [ "multi-user.target" ];
              after = [ "network.target" ];

              serviceConfig = {
                Type = "simple";
                ExecStart = "${pkg}/bin/calendar-app -config /run/school-lunch-calendar/config.yaml";
                Restart = "on-failure";
                RestartSec = 5;

                StateDirectory = "school-lunch-calendar";
                RuntimeDirectory = "school-lunch-calendar";
                RuntimeDirectoryMode = "0700";

                NoNewPrivileges = true;
                ProtectSystem = "strict";
                ProtectHome = true;
                PrivateTmp = true;
                PrivateDevices = true;
                ReadWritePaths = [ cfg.dataDir ];

                User = "school-lunch-calendar";
                Group = "school-lunch-calendar";
              };

              preStart = let
                sessionSecretCmd =
                  if cfg.sessionSecretFile != null
                  then "session_secret_file=${cfg.sessionSecretFile}"
                  else "session_secret=$(tr -dc A-Za-z0-9 < /dev/urandom | head -c 64)";
              in ''
                mkdir -p ${cfg.dataDir}
                ${sessionSecretCmd}
                printf 'port: ${toString cfg.port}\ndata_dir: ${cfg.dataDir}\npassword_file: ${cfg.passwordFile}\n' > /run/school-lunch-calendar/config.yaml
                if [ -n "$session_secret_file" ]; then
                  printf 'session_secret_file: %s\n' "$session_secret_file" >> /run/school-lunch-calendar/config.yaml
                else
                  printf 'session_secret: %s\n' "$session_secret" >> /run/school-lunch-calendar/config.yaml
                fi
              '';
            };

            users.users.school-lunch-calendar = {
              isSystemUser = true;
              group = "school-lunch-calendar";
            };
            users.groups.school-lunch-calendar = {};

            networking.firewall = lib.mkIf cfg.openFirewall {
              allowedTCPPorts = [ cfg.port ];
            };
          };
        };
    };
}

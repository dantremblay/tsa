# Project Structure

```
.
├── cmd/
│   ├── tsa/             # CLI entry point (client binary)
│   └── tsad/            # daemon entry point (server binary)
│
├── api/                 # Server-side API
│   ├── server/
│   │   ├── server.go    # route registration (all endpoints wired here)
│   │   ├── middleware/  # Https, AdminPassword, CAInit, Authorization, AdminOnly
│   │   └── router/      # HTTP handlers, one package per group
│   │       ├── acme/    # /acme/new-app, /acme/revoke-cert
│   │       ├── ca/      # /ca (public CA cert)
│   │       ├── crl/     # /crl/CRL.crl
│   │       ├── system/  # /system/* (info, auth, cert, ca/init, admin/password)
│   │       ├── usergroup/# /tsa/ug/register (admin)
│   │       └── tsaca/   # /tsa/ca/cert (admin, revoke-by-CN)
│   ├── auth/            # auth drivers (none, ldap)
│   ├── storage/         # server-side config storage (SQLite)
│   └── types/           # shared API types
│
├── cli/                 # Client-side (tsa CLI)
│   ├── command/         # cobra command tree (access, auth, cert, server, system)
│   ├── session/         # local session/token storage
│   └── storage/         # client-side SQLite storage
│
├── client/              # HTTP client used by the CLI to talk to the daemon
├── daemon/              # tsad daemon command + API bootstrap
├── pkg/                 # shared helpers (adf config, api, token, host, validation)
├── version/             # version struct populated via ldflags
├── contrib/init/systemd/# tsad.service unit
├── scripts/nfpm/        # package maintainer scriptlets (systemd install/remove)
├── gen/                 # man page + shell completion generators (per binary)
├── nfpm-tsa.yaml        # client package definition
├── nfpm-tsad.yaml       # daemon package definition
├── Dockerfile           # scratch-based tsad runtime image
└── Makefile
```

## Conventions

- **Two binaries, two packages**: `tsa` (client → `/usr/bin`) and `tsa-daemon`
  (daemon → `/usr/sbin` + systemd unit). Each has its own `nfpm-*.yaml`.
- **Routes** are all registered in `api/server/server.go`, grouped by auth level
  (public / basic-auth / JWT / JWT+AdminOnly).
- **CLI commands**: each group is a package under `cli/command/`; parent exposes
  `NewCommand()`, subcommands are unexported `newXxxCommand()`, logic in `runXxx()`.
- **Auth drivers** register config keys; LDAP keys live in `api/auth/driver/ldap`.
- **JWT claim parsing** (`pkg/token`) uses comma-ok type assertions — do not use
  bare `.(T)` assertions on claim values (they panic on unexpected types).

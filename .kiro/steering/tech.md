# Tech Stack

## Language & Runtime

- Go 1.26+ (module: `github.com/kassisol/tsa`)
- `CGO_ENABLED=0` — pure Go, statically linked

## Key Dependencies

| Library | Purpose |
|---|---|
| `spf13/cobra` | CLI framework |
| `labstack/echo/v4` + `echo-jwt/v4` | HTTP server and JWT middleware |
| `golang-jwt/jwt/v5` | JWT creation/parsing |
| `gorm.io/gorm` + `glebarez/sqlite` | ORM with pure-Go SQLite driver |
| `golang.org/x/crypto` | OCSP / crypto helpers |
| `juliengk/go-cert` | CA, PKIX, CSR helpers |
| `juliengk/go-ldapc` | LDAP client (auth backend) |
| `juliengk/stack` | HTTP client + JSON:API responses |

## Vendored / local dependencies

- `internal/juliengk/*` are project-controlled forks redirected via `replace`
  in `go.mod`. Edit them directly; they are not synced from upstream.
- The repo does **not** commit the `vendor/` directory (it is a local build
  artifact). Only `go.mod` / `go.sum` are tracked. Run `go mod vendor` locally.

## Build System

GNU Make + [nFPM](https://nfpm.goreleaser.com/) (no Dapper). See
`build-and-release.md` for the full workflow.

## Common Commands

```bash
make build       # build bin/tsa and bin/tsad
make cross       # cross-compile both binaries
make generate    # man pages + shell completions into dist/
make packages    # build tsa + tsa-daemon RPM/DEB into dist/
make image       # build the tsad Docker image
make vet / test / tidy / clean / version
```

## Data locations (daemon)

- `/var/lib/tsa/ca/` — CA database, root cert (`certs/ca.crt`), CRL, private key
- `/var/lib/tsa/certs/` — API server TLS cert and key

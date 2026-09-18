# Development Guide

This document describes how to build, package, and release TSA. TSA produces two
binaries: `tsad` (the CA server daemon) and `tsa` (the admin CLI).

## Prerequisites

- Go 1.26+
- [nFPM](https://nfpm.goreleaser.com/) for building packages:
  ```bash
  go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest
  ```
  Make sure `$(go env GOPATH)/bin` (usually `~/go/bin`) is on your `PATH`, or
  `make packages` will fail with `nfpm: command not found`.

## Building

```bash
make build       # build bin/tsa and bin/tsad
make cross       # cross-compile both binaries for linux/darwin, amd64/arm64
make generate    # generate man pages + bash completions into dist/
make packages    # build all four packages into dist/
make image       # build the tsad Docker image (kassisol/tsa:<version>)
make test        # run tests
make vet         # run go vet
make clean       # remove bin/ and dist/
```

`make packages` produces four artifacts — `tsa` and `tsa-daemon`, each as an RPM
and a DEB. For a different architecture:

```bash
TARGETARCH=arm64 make packages
```

## Packaging

Packaging is driven by two nFPM configs (no Dapper):

- `nfpm-tsa.yaml` — the client: installs `/usr/bin/tsa`, man pages, bash completion.
- `nfpm-tsad.yaml` — the daemon: installs `/usr/sbin/tsad`, the systemd unit at
  `/usr/lib/systemd/system/tsad.service`, and maintainer scriptlets
  (`scripts/nfpm/`) that run `systemctl daemon-reload` on install and stop/disable
  the service on removal.

Both binaries are static (`CGO_ENABLED=0`), so the packages install cleanly on
EL 8/9/10, SLES 12/15, and Debian/Ubuntu.

## Dependencies

- The `vendor/` directory is **not** committed. After changing dependencies,
  commit only `go.mod`/`go.sum`, then run `go mod vendor` locally for builds.
- Forks under `internal/juliengk/*` are redirected via `replace` in `go.mod`;
  edit them directly.

## Versioning

The version string is derived from git by the Makefile:

- dirty working tree → `<commit>-dirty`
- a tag on HEAD → the tag (e.g. `0.5.1`)
- otherwise → the short commit hash

To build packages/images that report a clean release version, the tag must be
**on HEAD**, or pass `VERSION=x.y.z` explicitly. `make image` bakes the resolved
version into both the binary and the image tag.

## Release process

1. Merge feature branches into `master` with `--no-ff`.
2. Verify `make build`, `make vet`, and `make packages` succeed.
3. Create an **annotated** tag with release notes covering everything since the
   previous tag (build from `git log <prevtag>..HEAD` — include features, fixes,
   and build changes, not just the last branch merged):
   ```bash
   git tag -a x.y.z -m "Release x.y.z\n\n..."
   ```
4. Push: `git push origin master && git push origin x.y.z`.
5. Rebuild and republish the image if the release changes the binary.

Guidelines:
- Security fixes that change the binary get their own patch release.
- Documentation-only changes do not require a new release.

## Notes

- The server rejects issuing a certificate while a valid certificate with the
  same distinguished name exists, so renewal must revoke-then-issue.
- JWT claim parsing in `pkg/token` must use comma-ok type assertions to avoid
  panics on unexpected claim types.

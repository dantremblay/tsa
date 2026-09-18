# Build & Release Workflow

TSA uses a Makefile + [nFPM](https://nfpm.goreleaser.com/) build system.
**Dapper and the old CentOS 7 spec were removed** (see release 0.5.0).

## Prerequisites

- Go 1.26+
- `nfpm` on PATH for packaging: `go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest`
  (ensure `$(go env GOPATH)/bin`, typically `~/go/bin`, is on PATH — common gotcha)

## Make targets

```bash
make build       # build bin/tsa and bin/tsad with version via ldflags
make cross       # cross-compile both binaries (linux/darwin, amd64/arm64)
make generate    # man pages + bash completions into dist/{tsa,tsad}/
make packages    # generate + build all 4 packages into dist/
make package-rpm # RPMs only (tsa + tsa-daemon)
make package-deb # DEBs only
make image       # build kassisol/tsa:<version> (scratch + tsad)
make vet / test / tidy / clean / version
```

`make packages` produces four artifacts: `tsa` and `tsa-daemon`, each as RPM and
DEB. For another arch: `TARGETARCH=arm64 make packages`.

## Packaging config

- `nfpm-tsa.yaml` — client: `/usr/bin/tsa` + man + completion.
- `nfpm-tsad.yaml` — daemon: `/usr/sbin/tsad`, systemd unit at
  `/usr/lib/systemd/system/tsad.service`, plus maintainer scriptlets in
  `scripts/nfpm/` (daemon-reload on install; stop/disable on remove).
- Man pages use `type: tree` (this nFPM version rejects `glob`).
- Generators output to `dist/{tsa,tsad}/` (override via `TSA_MAN_PATH`,
  `TSAD_MAN_PATH`, `TSA_SC_PATH`, `TSAD_SC_PATH`).

## Versioning

Version comes from git via the Makefile (dirty → `<commit>-dirty`; tag on HEAD →
tag; else bare commit). A clean `<version>` build requires the tag **on HEAD**.
`make image` bakes the resolved version into the binary and image tag.

## Release conventions

- **Annotated** tags only, with release notes.
- Merge feature branches `--no-ff`.
- Build tag notes from the full range since the previous tag
  (`git log <prevtag>..HEAD`) — include ALL changes (features, fixes, build), not
  just the branch you happened to finish. (0.5.0 initially missed the API-endpoint
  commit for this reason.)
- Security fixes that change the binary get their own patch version.
- `vendor/` is not committed — commit only `go.mod`/`go.sum` after dep changes,
  then `go mod vendor` locally.

## Remote

`tsa` is hosted on GitHub. Normal `git push` to master.
(Note: `tsa-web`, a sibling repo, is on Gerrit — different push flow.)

## Known follow-ups

- Two endpoints revoke certs by CN: `DELETE /system/cert/revoke-by-cn/:cn` and
  the newer `DELETE /tsa/ca/cert`. They overlap and could be consolidated.

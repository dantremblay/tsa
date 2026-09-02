# HBM TSA

HBM TSA is a CA (Certification Authority) server that issues certificates for
enabling TLS on the Docker Engine. It works together with
[twic](https://github.com/kassisol/twic), the client used to request and manage
those certificates.

[![Go Report Card](https://goreportcard.com/badge/github.com/kassisol/tsa)](https://goreportcard.com/report/github.com/kassisol/tsa)

TSA ships two binaries:

- **`tsad`** — the CA server daemon (runs as a service, listens for API requests)
- **`tsa`** — the admin CLI used to configure and operate a running daemon

## Installation

Install the prebuilt packages for your distribution. The daemon and the client
are packaged separately.

```bash
# RHEL / Rocky / AlmaLinux / SLES
sudo dnf install ./tsa-daemon-<version>.x86_64.rpm   # the tsad server
sudo dnf install ./tsa-<version>.x86_64.rpm          # the tsa admin CLI

# Debian / Ubuntu
sudo dpkg -i ./tsa-daemon_<version>_amd64.deb
sudo dpkg -i ./tsa_<version>_amd64.deb
```

The `tsa-daemon` package installs a systemd unit. Start the server with:

```bash
sudo systemctl enable --now tsad
```

## Quick start

Once `tsad` is running, use the `tsa` CLI to point at it, authenticate, and
initialize the CA:

```bash
# Register the daemon as a named server
tsa server add local https://tsa.example.com

# Log in (default user is 'admin'); this stores a session token
tsa access login local

# Initialize the root CA (prompts for country, state, org, etc.)
tsa init

# Change the default admin password
tsa passwd
```

The server is now ready to issue certificates. Clients use
[twic](https://github.com/kassisol/twic) to request them.

To issue certificates to real users (rather than only `admin`), enable an
authentication backend:

```bash
# Enable LDAP authentication
tsa auth enable ldap

# Add the required LDAP settings
tsa auth add auth_host ldap.example.com
tsa auth add auth_port 636
tsa auth add auth_group_admin "cn=docker-admins,ou=groups,dc=example,dc=com"
tsa auth add auth_group_user  "cn=docker-users,ou=groups,dc=example,dc=com"
# ... plus auth_bind_username, auth_bind_password, auth_search_base_user,
#     auth_search_filter, auth_attr_members, auth_tls
```

## Command reference

Run `tsa <command> --help` for full details on any command.

| Command | Description |
|---|---|
| `tsa server add [name] [url]` | Register a TSA server |
| `tsa server ls` / `rm` | List / remove registered servers |
| `tsa access login [server]` | Get an access token for a server |
| `tsa access ls` / `use` / `unuse` / `status` / `rm` | Manage sessions |
| `tsa init` | Initialize the root CA |
| `tsa passwd` | Change the admin password |
| `tsa info` | Show server information |
| `tsa version` | Show client and server versions |
| `tsa auth enable [type]` | Enable an auth backend (e.g. `ldap`) |
| `tsa auth disable` | Disable the auth backend |
| `tsa auth add` / `rm [key] [value]` | Manage auth configuration keys |
| `tsa auth ls` | List auth configuration |
| `tsa cert ls` | List issued certificates |
| `tsa cert revoke [serial] \| --cn [name]` | Revoke a certificate by serial number or CN |

## Daemon (`tsad`) options

| Flag | Default | Description |
|---|---|---|
| `-a, --bind-address` | `0.0.0.0` | Address to bind |
| `-p, --bind-port` | `80` (`443` with TLS) | Port to bind |
| `-f, --fqdn` | | API FQDN (required with `--tlsgen`) |
| `-t, --tls` | `false` | Enable TLS |
| `--tlsgen` | `false` | Auto-generate a self-signed TLS certificate |
| `--tls-duration` | `60` | Self-signed certificate duration (days) |
| `--tlscert` / `--tlskey` | under `/var/lib/tsa/certs` | Paths to TLS certificate and key |

Running with TLS is strongly recommended, since clients send credentials to the
API:

```bash
tsad --tls --tlsgen --fqdn tsa.example.com
```

## Admin API endpoints

These endpoints require a JWT from an authenticated admin. They are typically
used by tooling rather than directly.

| Method | Path | Description |
|---|---|---|
| `POST` | `/tsa/ug/register` | Register a user group (body: group value) |
| `DELETE` | `/tsa/ug/register` | Unregister a user group |
| `DELETE` | `/tsa/ca/cert` | Revoke all valid certificates for a CN (body: CN) |

## Data locations

The daemon stores its state under `/var/lib/tsa`:

- `/var/lib/tsa/ca/` — CA database, root certificate (`certs/ca.crt`), CRL, and private key
- `/var/lib/tsa/certs/` — the API server's TLS certificate and key

## Building from source

Requires Go 1.26+ and, for packaging, [nFPM](https://nfpm.goreleaser.com/)
(`go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest`).

```bash
make build        # build both binaries into bin/
make cross        # cross-compile for linux/darwin, amd64/arm64
make generate     # generate man pages and shell completions
make packages     # build tsa and tsa-daemon RPM + DEB packages into dist/
make image        # build the tsad runtime Docker image
```

To build packages for a different architecture:

```bash
TARGETARCH=arm64 make packages
```

## User feedback

If you have any problems with or questions about this application, please open a
[GitHub issue](https://github.com/kassisol/tsa/issues).

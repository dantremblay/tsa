# Product: TSA

TSA (part of the Harbormaster / HBM suite) is a CA (Certification Authority)
server that issues X.509 certificates for enabling TLS on the Docker Engine. It
is the server-side counterpart to [twic](https://github.com/kassisol/twic), the
client used to request and manage those certificates.

## Two binaries

- **`tsad`** — the CA server daemon. Runs as a systemd service, exposes an HTTP
  API, signs CSRs, maintains the CA database and CRL.
- **`tsa`** — the admin CLI. Configures and operates a running daemon (init CA,
  manage auth backends, list/revoke certificates, manage sessions).

## Certificate types issued

- **client** — for Docker CLI users (CN must match the authenticated user; only
  admins may request other CNs)
- **engine** — for Docker daemons (CN must be a valid FQDN; admin only)

## Authentication & authorization

- Bootstrap admin uses a password; `GET /new-authz` (HTTP Basic) returns a JWT.
- Certificate operations require a JWT bearer token.
- Optional LDAP backend maps users to admin/user groups
  (`auth_group_admin` / `auth_group_user` config keys).

## Important domain rule

The server **refuses to issue a certificate while a valid certificate with the
same distinguished name already exists**. Therefore certificate *renewal* must
**revoke the old cert before requesting the new one** (revoke-then-issue). This
is why `twic`'s renew commands are ordered that way.

## Related note on mTLS

The certificates TSA issues are used by Docker with `--tlsverify`, which is
already mutual TLS (both daemon and client present CA-signed certs). The
TSA API itself uses one-way TLS + token auth, not client-cert auth.

## License

GPL-3.0

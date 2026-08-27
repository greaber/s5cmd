# Fork release policy

This fork packages a narrow set of object-storage compatibility improvements
while upstream s5cmd releases are inactive. It is not intended to become an
independent general-purpose distribution.

The `v2.3.0-headers.1` release adds:

- repeatable, signed `--request-header NAME:VALUE` options;
- retries for the exact S3 `SignatureDoesNotMatch` error code using s5cmd's
  existing retry count and backoff;
- confinement of wildcard S3 downloads to their local destination;
- removal of unnecessary HEAD requests when progress reporting is disabled;
- debug and trace logging on stderr; and
- the final AWS SDK for Go v1 release, `v1.55.8`.

Releases contain one statically linked Linux amd64 archive and a checksum file.
Add another target only when a concrete runtime consumer needs it. Release tags
are built and tested by `.github/workflows/goreleaser.yml` with the pinned Go
toolchain declared there.

Do not add credentials, endpoints, or provider policy to the fork. Callers own
those settings and pass provider-specific headers explicitly.

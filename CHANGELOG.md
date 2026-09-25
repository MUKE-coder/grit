# Changelog

The full changelog lives at
**[gritframework.dev/docs/changelog](https://gritframework.dev/docs/changelog)**.

Every release has an entry there, newest first, grouped by week. Entries are
written to be read: what was wrong, why it was wrong, and what changed. Not a
list of commit subjects.

This file is a pointer rather than a copy, because a second changelog is a
changelog that goes stale.

## Other places to look

- **[Releases](https://github.com/MUKE-coder/grit/releases)** on GitHub, with
  the binaries, their checksums and the signing material described in
  [SECURITY.md](SECURITY.md).
- **[What is stable and what is not](https://gritframework.dev/docs/stability)**,
  which says which parts of Grit you can build on and which are still moving.
- `grit update` brings the CLI to the latest release. `grit upgrade` brings an
  existing project's framework-owned files up to the CLI's version, and leaves
  the files you have edited alone.

## Versioning

Grit is versioned `MAJOR.MINOR.PATCH`, and a release is cut whenever something
ships rather than on a schedule. The stability page above is the contract:
version numbers alone do not tell you whether a thing is settled.

# clipshift
WIP - A cross-platform clipboard syncing application written in Go

# Installation
Pre-built binaries are still a work in progress. In the meantime you'll have to use `go install github.com/getclipshift/clipshift@latest`

# Usage

TODO - run `clipshift --help` for now

# Backends
Backends can be configured with the `clipshift config add-backend` and `clipshift config edit-backend` commands.

## [ntfy.sh](https://ntfy.sh) (recommended)
A push notification service. Currently the best option to use as a backend.

[ntfy docs](docs/ntfy.md)

## [nostr](https://github.com/nostr-protocol/nostr) (advanced)
A flexible protocol using websockets with encrypted direct message support.

[nostr docs](docs/nostr.md)

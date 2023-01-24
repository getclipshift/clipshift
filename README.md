# clipshift
WIP - A cross-platform clipboard syncing application written in Go

# Installation
Pre-built binaries are still a work in progress. In the meantime you'll have to use `go install...`

Prerequisites:
- All: go 1.18+
- Linux: `sudo apt install libgtk-3-dev libayatana-appindicator3-dev` (or similar for your distro)

`go install github.com/jhotmann/clipshift@latest`

# Usage

TODO

# Backends

## [ntfy.sh](https://ntfy.sh) (recommended)
Currently the best option to use as a backend. There is a free public instance available or the server can be self-hosted really easily. The Android client also works well alongside Tasker to allow your Android devices to sync their clipboard with the rest of your devices.

[ntfy docs](docs/ntfy.md)

## [nostr](https://github.com/nostr-protocol/nostr) (advanced)
In order to keep strain on public relays down, I can only recommend this backend if you are willing to run your own relay. [Nostream](https://github.com/Cameri/nostream) is a great option.

[nostr docs](todo)

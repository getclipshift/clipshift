# clipshift
A cross-platform clipboard syncing application written in Go

# Installation
There are two flavors: a CLI named `clipshift` and a version that will show an icon in the system tray while syncing name `clipshift-tray`. The `clipshift-tray` version also includes the ability to start sync at login.

- With Go: `go install github.com/getclipshift/clipshift@latest`
- With [Homebrew](https://brew.sh)
  - `clipshift`

    ```sh
    brew tap getclipshift/tap
    brew install clipshift
    ```

  - `clipshift-tray`

    ```sh
    brew tap getclipshift/tap
    brew install --cask --no-quarantine clipshift
    ```

- With [Scoop](https://scoop.sh/) (installs both `clipshift` and `clipshift-tray`)

    ```sh
    scoop bucket add clipshift https://github.com/getclipshift/bucket.git
    scoop install clipshift
    ```

- Bash or zsh
Download and extract the binary for your platform/arch to some place that is in your `$PATH`, example:

    ```sh
    curl -o clipshift.tar.gz "https://github.com/getclipshift/clipshift/releases/latest/download/clipshift_$(uname -s | tr '[:upper:]' '[:lower:]')_$(uname -m).tar.gz"
    sudo tar -C /usr/local/bin -xzf clipshift.tar.gz
    rm clipshift.tar.gz
    ```

# Usage
1. Configure at least one [backend](#backends): `clipshift config add-backend`
1. Run `clipshift sync` or `clipshift-tray` to have clipboard contents synced automatically
1. Run `clipshift send` to manually send the current clipboard
1. Run `clipshift get` to manually receive the latest clipboard

You can use `clipshift help` to see all commands and their options

# Backends
A backend server must be configured for clipshift to connect to and send/receive clipboard contents.

A backend can operate in 4 different ways:
1. Sync - will send and receive clipboard contents automatically during sync or a manual send/get command
1. Push - will send clipboard contents automatically during sync or a manual send command
1. Pull - will receive clipboard contents automatically during sync or a manual get command
1. Manual - will only be used when a manual send/get command is run

## [ntfy.sh](https://ntfy.sh) (recommended)
A push notification service. Currently the best option to use as a backend.

[ntfy docs](docs/ntfy.md)

## [Nostr relay](https://github.com/nostr-protocol/nostr) (advanced)
A flexible protocol using websockets with encrypted direct message support.

[nostr docs](docs/nostr.md)

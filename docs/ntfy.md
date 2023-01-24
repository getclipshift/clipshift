#  ntfy

[ntfy](https://ntfy.sh) is a simple pub/sub server for sending notifications to your devices.

## clipshift config

- Using the public server

    ```yaml
    backends:
      - type: ntfy
        host: https://ntfy.sh
        topic: some-random-string-here
        action: sync
        encryptionkey: some-key-here # optional, but recommended for the public instance
    ```

- Using a [self-hosted](https://docs.ntfy.sh/install/) instance

    ```yaml
    backends:
      - type: ntfy
        host: https://ntfy.example.com
        user: some-user
        pass: some-password
        topic: clipshift
        action: sync
        encryptionkey: # optional, leave blank for no encryption
    ```


## Android

TODO
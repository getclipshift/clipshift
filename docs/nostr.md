# Nostr

[Nostr](https://github.com/nostr-protocol/nostr) is a new protocol that is mainly used as a Twitter-like social network but is flexible enough for a wide array of uses. It also supports encrypted direct messages, which is what clipshift uses for syncing the clipboard between devices. You should generate a new keypair for clipshift using one of the CLI tools (rana, noscl, etc) or nos2x.

In order to avoid taxing public relays, I suggest running a private relay. [Nostream](https://github.com/Cameri/nostream) is a great option.

[Awesome Nostr](https://github.com/aljazceru/awesome-nostr)

## clipshift config
Use `clipshift config add-backend` to be guided through the configuration options or you can manually edit `~/.clipshift/config.yaml` and add the following:

```yaml
backends:
  - type: nostr
    host: wss://nostr-relay.example.com
    user: hex-public-key # technically you don't even need this
    pass: hex-private-key
    action: sync
```

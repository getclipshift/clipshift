FROM goreleaser/goreleaser-cross:v1.19

RUN apt-get update && apt-get install -y libgtk-3-dev libayatana-appindicator3-dev

ENTRYPOINT ["/usr/bin/tini", "--", "/entrypoint.sh"]
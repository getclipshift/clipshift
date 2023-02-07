#!/bin/bash
# CircleCI does not support arm64 macs yet, so this script is to be run on a machine with Apple silicon to finalize a relase

read -r -p 'Tag: ' tag
CIRCLE_TAG=$tag

rm -rf dist
mkdir dist

CGO_ENABLED=1 go build -ldflags "-X 'github.com/getclipshift/clipshift/cmd.Version=${CIRCLE_TAG/v/}'" -o dist/clipshift main.go
CGO_ENABLED=1 go build -ldflags "-X 'github.com/getclipshift/clipshift/cmd.Version=${CIRCLE_TAG/v/}'" -tags tray -o dist/clipshift-tray main.go
tar -czf dist/clipshift_darwin_arm64.tar.gz -C dist clipshift
tar -czf dist/clipshift-tray_darwin_arm64.tar.gz -C dist clipshift-tray

mkdir -p dist/clipshift.app/Contents/MacOS
mkdir -p dist/clipshift.app/Contents/Resources
cp internal/ui/clipboard.icns dist/clipshift.app/Contents/Resources/clipshift.icns
cp dist/clipshift-tray dist/clipshift.app/Contents/MacOS/clipshift-tray
cp dist/clipshift dist/clipshift.app/Contents/MacOS/clipshift
cp build/clipshift.plist dist/clipshift.app/Contents/Info.plist
tar -czf dist/clipshift.app_arm64.tar.gz -C dist clipshift.app

cd dist || exit

curl -H "Circle-Token: $CIRCLE_TOKEN" https://circleci.com/api/v1.1/project/github/getclipshift/clipshift/latest/artifacts \
  | grep -o "https://[^\"]*" \
  | wget --header "Circle-Token: $CIRCLE_TOKEN" --input-file -

sha256sum clipshift_darwin_arm64.tar.gz >> checksums.txt
sha256sum clipshift-tray_darwin_arm64.tar.gz >> checksums.txt
sha256sum clipshift.app_arm64.tar.gz >> checksums.txt

cd ..

awk -v macarm=$(sha256sum dist/clipshift_darwin_arm64.tar.gz | cut -f 1 -d " ") \
    '{gsub("--MAC-ARM-SHA--", macarm); print}' \
    dist/clipshift.rb > dist/clipshift-final.rb

awk -v arm=$(sha256sum dist/clipshift.app_arm64.tar.gz | cut -f 1 -d " ") \
    '{gsub("--ARM-SHA--", arm); print}' \
    dist/clipshift_cask.rb > dist/clipshift_cask-final.rb


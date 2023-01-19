#!/bin/bash

platform="darwin"
arch=$1

rm -rf dist
mkdir dist
echo "Building clipshift"
CGO_ENABLED=1 go build -o "dist/clipshift_${platform}_${arch}" .

root="dist/clipshift_${platform}_${arch}.app"
echo "Creating macos app $root"
rm -rf "$root"
mkdir -p "$root/Contents/MacOS"
mkdir -p "$root/Contents/Resources"
cp internal/ui/clipboard.icns "$root/Contents/Resources/clipshift.icns"
cp "dist/clipshift_${platform}_$arch" "$root/Contents/MacOS/clipshift"
cp build/clipshift.plist "$root/Contents/Info.plist"

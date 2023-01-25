#!/bin/bash

distdir=$1
appname=$2

root="dist/$appname"
echo "Creating macos app $root"

rm -rf "$root"
mkdir -p "$root/Contents/MacOS"
mkdir -p "$root/Contents/Resources"
cp internal/ui/clipboard.icns "$root/Contents/Resources/clipshift.icns"
cp "dist/$distdir/clipshift" "$root/Contents/MacOS/clipshift"
cp build/clipshift.plist "$root/Contents/Info.plist"
#!/bin/bash

set -e

root=$(mktemp -d)

mkdir -p $root/usr/local/bin
cp verbumsrvr verbumctl $root/usr/local/bin

mkdir -p $root/usr/local/share/verbum
cp -R statics templates $root/usr/local/share/verbum

v="$TRAVIS_BUILD_NUMBER"
if [ -z "$v" ]; then
    t=$(date +"%Y%m%d-%H%M%S")
    h=$(git rev-parse --short HEAD)
    v="$t-$h"
fi

mkdir -p $root/DEBIAN
cat > $root/DEBIAN/control <<EOF
Package: verbum
Version: ${v}
Architecture: amd64
Essential: no
Section: web
Priority: optional
Depends:
Maintainer: Vadzim Ramanenka
Installed-Size:
Description: Verbum - Online Platform Dictionary.
EOF

dpkg -b $root deploy/verbum.deb

rm -r $root
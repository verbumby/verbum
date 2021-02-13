#!/bin/bash

echom () {
	echo "==> $(date --iso-8601=ns) $1"
}
export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
cd /home/deploy/debdelivery

test -f verbum.deb.sha256 || exit 0
echom "==> verifying checksum"
/usr/bin/sha256sum -c verbum.deb.sha256 || { echo "sha256 check fails"; exit 1; }

echom "==> installing"
/usr/bin/dpkg -i verbum.deb || { echo "installation has failed"; exit 2; }

echom "==> removing installation files"
rm verbum.deb verbum.deb.sha256 || { echo "failed to rm verbum.deb and verbum.deb.sha256"; exit 3; }

echom "==> done"

#!/bin/bash
set -xe
cd deploy
sha256sum verbum.deb > verbum.deb.sha256
ssh-keyscan -H 165.227.129.101 >> ~/.ssh/known_hosts
scp -i id_ed25519 verbum.deb verbum.deb.sha256 deploy@165.227.129.101:~/

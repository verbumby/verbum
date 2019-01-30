#!/bin/bash
set -xe
cd deploy
openssl aes-256-cbc -K $encrypted_861f9db00f37_key -iv $encrypted_861f9db00f37_iv -in id_ed25519.enc -out id_ed25519 -d
chmod 600 id_ed25519
sha256sum verbum.deb > verbum.deb.sha256
ssh-keyscan -H 165.227.129.101 >> ~/.ssh/known_hosts
scp -i id_ed25519 verbum.deb verbum.deb.sha256 deploy@165.227.129.101:~/

#!/bin/bash
set -e
cd deploy
openssl aes-256-cbc -K $encrypted_a25db77e2a26_key -iv $encrypted_a25db77e2a26_iv -in id_ed25519.enc -out id_ed25519 -d
chmod 600 id_ed25519
sha256sum verbum.deb > verbum.deb.sha256
ssh-keyscan -H 165.227.129.101 >> ~/.ssh/known_hosts
sftp -i id_ed25519 deploy@165.227.129.101 <<EOF
cd debdelivery
put verbum.deb
put verbum.deb.sha256
EOF

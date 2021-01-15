#!/bin/bash
set -e
cd deploy
echo "$ID_ED25519" > id_ed25519
chmod 600 id_ed25519
sha256sum verbum.deb > verbum.deb.sha256
mkdir -p ~/.ssh
ssh-keyscan -H 49.12.79.27 >> ~/.ssh/known_hosts
sftp -i id_ed25519 deploy@49.12.79.27 <<EOF
cd debdelivery
put verbum.deb
put verbum.deb.sha256
EOF
rm id_ed25519

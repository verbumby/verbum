#!/bin/bash

curl -s -XPUT 'http://localhost:9200/_snapshot/backup' \
    -H 'Content-Type: application/json' \
    -d '{
        "type": "fs",
        "settings": {
            "location": "'$1'",
            "compress": true
        }
    }'

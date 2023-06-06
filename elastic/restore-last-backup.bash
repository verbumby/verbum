#!/usr/bin/env bash

EXISTING=$(curl localhost:9200/_cat/indices?format=json 2>/dev/null | jq -r '.|map(select(.index|startswith("dict-")).index)|join(",")')
curl -XDELETE localhost:9200/$EXISTING ; echo
LAST=$(curl localhost:9200/_snapshot/backup/_all 2>/dev/null | jq -r '.snapshots[].snapshot' | sort | tail -n 1) \
&& curl -XPOST localhost:9200/_snapshot/backup/$LAST/_restore?wait_for_completion=true -H 'content-type: application/json' -d '{"indices":"dict-*"}'

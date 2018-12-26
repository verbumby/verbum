#!/bin/bash

esaddr=http://localhost:9200

q=$1
if [ -z "$q" ]; then
    q="ак"
fi

resp=$(curl -XPOST -s -H 'Content-Type: application/json' "$esaddr/dict-*/_search" -d '{
    "_source": false,
    "suggest": {
        "HeadwordSuggest": {
            "prefix": "'$q'",
            "completion": {
                "field": "Suggest",
                "skip_duplicates": true,
                "size": 10
            }
        }
    }
}')

echo "$resp" | jq '.suggest.HeadwordSuggest[].options[].text'

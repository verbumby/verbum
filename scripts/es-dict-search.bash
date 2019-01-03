#!/bin/bash

esaddr=http://localhost:9200

q=$1
if [ -z "$q" ]; then
    q="штармоўка"
fi

resp=$(curl -XPOST -s -H 'Content-Type: application/json' "$esaddr/dict-*/_search" -d '{
    "query": {
        "multi_match": {
            "query": "'$q'",
            "fields": [
                "Headword^4",
                "Headword.Smaller^3",
                "HeadwordAlt^2",
                "HeadwordAlt.Smaller^1"
            ]
        }
    }
}')

echo "$resp" | jq '.'
# echo "$resp" | jq '.hits.hits[] | (._score | tostring )+ " -- " + ._index + " -- " + (._source.Headword | join(", ")) + " -- " + (._source.HeadwordAlt? | join(", ")) '

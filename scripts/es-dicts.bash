#!/bin/bash

esaddr=http://localhost:9200
index=dicts

curl -s -XDELETE $esaddr/$index | jq .
curl -s -XPUT -H 'Content-Type: application/json' $esaddr/$index -d '{
    "settings": {
        "number_of_shards" : 1,
        "number_of_replicas" : 0
    },
    "mappings": {
        "_doc": {
            "properties": {
                "Title": {
                    "type": "keyword"
                }
            }
        }
    }
}' | jq .

curl -s -XPUT -H 'Content-Type: application/json' $esaddr/$index/_doc/rvblr -d '{
    "Title": "Тлумачальны слоўнік беларускай мовы (rv-blr.com)"
}' | jq .

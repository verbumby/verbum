#!/bin/bash

esaddr=http://localhost:9200
index=rvblr2

curl -s -XDELETE $esaddr/$index | jq .
curl -s -XPUT -H 'Content-Type: application/json' $esaddr/$index -d '{
    "settings": {
        "number_of_shards" : 1,
        "number_of_replicas" : 0,
        "analysis": {
            "analyzer": {
                "hw": {
                    "type": "custom",
                    "tokenizer": "keyword",
                    "filter": ["lowercase"]
                },
                "hw_smaller": {
                    "type": "custom",
                    "tokenizer": "hw_smaller",
                    "filter": ["lowercase"]
                }
            },
            "tokenizer": {
                "hw_smaller": {
                    "type": "char_group",
                    "tokenize_on_chars": ["-"]
                }
            }
        }
    },
    "mappings": {
        "_doc": {
            "properties": {
                "headword": {
                    "type": "text",
                    "analyzer": "hw",
                    "fields": {
                        "hw_smaller": {
                            "type": "text",
                            "analyzer": "hw_smaller",
                            "search_analyzer": "hw"
                        }
                    }
                },
                "headword_alt": {
                    "type": "text",
                    "analyzer": "hw",
                    "fields": {
                        "hw_smaller": {
                            "type": "text",
                            "analyzer": "hw_smaller",
                            "search_analyzer": "hw"
                        }
                    }
                },
                "suggest": {
                    "type": "completion",
                    "analyzer": "hw"
                },
                "content": {
                    "type": "text",
                    "index": false
                }
            }
        }
    }
}' | jq .

curl -s -XPOST -H 'Content-Type: application/json' $esaddr/$index/_doc/  -d '{
    "headword": ["штармоўка"],
    "headword_alt": ["штармоўка"],
    "suggest": ["штармоўка"],
    "content": "Куртка з моцнай воданепрымальнай тканіны (у маракоў, геолагаў, спартсменаў і пад.)."
}' | jq .

curl -s -XGET  $esaddr/$index/_refresh | jq .

curl -s -XPOST -H 'Content-Type: application/json' $esaddr/rvblr2,rvblr/_search -d '{
	"query": {
		"multi_match": {
			"query": "штармоўка",
			"fields": [
				"headword^4",
				"headword.hw_smaller^3",
				"headword_alt^2",
				"headword_alt.hw_smaller^1"
			]
		}
	},
    "suggest": {
        "hw-suggest": {
            "prefix": "Шта",
            "completion": {
                "field": "suggest",
                "skip_duplicates": true
            }
        }
    }
}' | jq .

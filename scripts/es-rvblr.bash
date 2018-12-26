#!/bin/bash

esaddr=http://localhost:9200
index=rvblr

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
    "headword": ["штаб"],
    "headword_alt": ["штабны"],
    "suggest": ["штаб", "штабны"],
    "content": "Орган кіравання войскамі, а таксама асобы, якія ўваходзяць у яго. Генеральны ш. Ш. дывізіі. Ш. палка."
}' | jq .

curl -s -XPOST -H 'Content-Type: application/json' $esaddr/$index/_doc/  -d '{
    "headword": ["штаб", "штаб-"],
    "suggest": ["штаб-", "штаб"],
    "content": "Першая частка складаных слоў са знач. які мае адносіны да штаба (у 1 знач.), напр. штаб-афіцэр, штаб-кватэра, штаб-ротмістр."
}' | jq .

curl -s -XPOST -H 'Content-Type: application/json' $esaddr/$index/_doc/  -d '{
    "headword": ["штаб-афіцэр"],
    "headword_alt": ["штаб-афіцэрскі"],
    "suggest": ["штаб-афіцэр", "штаб-афіцэрскі"],
    "content": "У царскай і некаторых замежных арміях: афіцэр у чыне палкоўніка, падпалкоўніка і маёра."
}' | jq .

curl -s -XGET  $esaddr/$index/_refresh | jq .

curl -s -XPOST -H 'Content-Type: application/json' $esaddr/$index/_search -d '{
	"query": {
		"multi_match": {
			"query": "афіцэрскі",
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
            "prefix": "шта",
            "completion": {
                "field": "suggest",
                "skip_duplicates": true
            }
        }
    }
}' | jq .

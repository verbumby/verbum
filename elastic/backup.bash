#!/bin/bash

echo "---------------------$(date --iso-8601=ns)---------------------"
t=$(date +"%Y%m%d-%H%M%S")
output=$(curl -fs -XPUT "http://localhost:9200/_snapshot/backup/${t}?wait_for_completion=true") || {
    echo "failed to backup";
    exit 11;
}
state=$(echo "$output" | jq -r '.snapshot.state')
if [ "$state" != 'SUCCESS' ]; then
    echo "snapshot state is not success: ${output}"
    exit 12
fi

/usr/local/bin/aws s3 sync \
	/var/lib/elasticsearch-backups/backup/ \
	s3://verbumby-backup \
	--only-show-errors \
	--delete	|| {
	echo "sync failed"
	exit 13
}

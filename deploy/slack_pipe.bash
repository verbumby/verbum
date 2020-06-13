#!/bin/bash

header="$1"; shift
webhook_addr="$1"; shift

exec 5>&1
output=$( $@ 2>&1 | tee >(cat - >&5); exit ${PIPESTATUS[0]} )
exit_code=$?
if [ $exit_code -eq 0 ]; then
    if [ -z "$output" ]; then
        exit 0
    fi
    output=$(printf "$header\n"'```'"$output"'```')
else
    output=$(printf "<!channel> :warning: $header\n"'```'"$output\nExit Code: $exit_code"'```')
fi

curl -X POST -H 'Content-type: application/json' \
    --data "$(jq -n --arg output "$output" '{text: $output}')" \
    $webhook_addr

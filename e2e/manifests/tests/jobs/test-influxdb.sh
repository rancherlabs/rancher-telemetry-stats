#!/usr/bin/env bash

set -x

INFLUXDB_PORT=${INFLUXDB_PORT:-8086}

# If run inside Kubernetes, assume Alpine and install dependencies.
if [ -n "${KUBERNETES_SERVICE_HOST}" ]; then
    apt update && \
        apt install jq -y
fi

# Just try.
while true; do
    json=$(
        influx \
            -host influxdb \
            -database telemetry \
            -execute 'SELECT * FROM telemetry LIMIT 1' \
            -format=json
    )
    if [[ "$?" -ne 0 || -z "$json" ]]; then
        echo >&2 "error: failed to fetch data"
        sleep 2
        continue
    fi

    echo "DEBUG: \$json: $json"

    result=$(echo $json | jq -r '.results[0].series | length')
    echo "DEBUG: \$result: $result"
    if [[ "$?" -ne 0 || -z "$result" || "$result" = "null" ]]; then
        echo >&2 "error: couldn't get records, result returned: $result"
        sleep 2
        continue
    fi

    if [ "$result" -ge 1 ]; then
        echo "success"
        exit 0
    fi

    echo >&2 "warning: no records found, result contains: $result"
    pause 1
done

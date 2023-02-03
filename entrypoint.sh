#!/bin/bash
/bin/rancher-telemetry-stats -hours=$HOURS --refresh=$REFRESH -url=$TELEMETRY_URL -accesskey=$TELEMETRY_ACCESS_KEY -secretkey=$TELEMETRY_SECRET_KEY -influxuser=$INFLUX_USER -influxpass=$INFLUX_PASS -influxurl=$INFLUX_URL -influxdb=$INFLUX_DB

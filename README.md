[![](https://images.microbadger.com/badges/image/rawmind/rancher-telemetry-stats.svg)](https://microbadger.com/images/rawmind/rancher-telemetry-stats "Get your own image badge on microbadger.com")

rancher-telemetry-stats
=====================

This image run rancher-telemetry-stats app. It comes from [rawmind/alpine-base][alpine-base].

## Build

```
docker build -t rawmind/rancher-telemetry-stats:<version> .
```

## Versions

- `0.1` [(Dockerfile)](https://github.com/rawmind0/rancher-telemetry-stats/blob/0.1/Dockerfile)


## Usage

This image run rancher-telemetry-stats service. rancher-telemetry-stats get metrics from rancher nginx logs files and send them to a influx in order to be explored by a grafana. 

If you run in daemon mode it will tail files and send metrics every refresh seconds. 

```
Usage of rancher-telemetry-stats:
  -accessKey string
      Rancher access key. Or env TELEMETRY_ACCESS_KEY (default "")
  -filepath string
      Log files to analyze, wildcard allowed between quotes. (default "/var/log/nginx/access.log")
  -format string
      Output format. influx | json (default "influx")
  -geoipdb string
      Geoip db file. (default "GeoLite2-City.mmdb")
  -influxdb string
      Influx db name (default "telemetry")
  -influxpass string
      Influx password
  -influxurl string
      Influx url connection (default "http://localhost:8086")
  -influxuser string
      Influx username
  -limit int
      Limit batch size (default 2000)
  -refresh int
      Send metrics every refresh seconds. daemon mode (default 120)
  -secretKey string
      Rancher secret key. Or env TELEMETRY_SECRET_KEY (default "")
  -url string
      Rancher telemetry url. (default "http://telemetry.rancher.io")
```

NOTE: You need influx already installed and running. The influx db would be created if doesn't exist.

## Metrics

Metrics are on the form.....

```
requests,city=Toronto,country=Canada,host=git.rancher.io,ip=xx.xx.xx.xx,method=GET,path=/rancher-catalog.git/info/refs?service\=git-upload-pack,status=200 ip="xx.xx.xx.xx" 1491289498000000000
```

[alpine-base]: https://github.com/rawmind0/alpine-base

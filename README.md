[![](https://images.microbadger.com/badges/image/rawmind/rancher-telemetry-stats.svg)](https://microbadger.com/images/rawmind/rancher-telemetry-stats "Get your own image badge on microbadger.com")

rancher-telemetry-stats
=====================

This image run rancher-telemetry-stats app. It comes from [rawmind/alpine-base][alpine-base].

## Build

```
docker build -t rawmind/rancher-telemetry-stats:<version> .
```

## Versions

- `0.1-3` [(Dockerfile)](https://github.com/rawmind0/rancher-telemetry-stats/blob/0.1-3/Dockerfile)


## Usage

This image run rancher-telemetry-stats service. rancher-telemetry-stats get metrics from rancher telemetry service and send them to a influx in order to be explored by a grafana. 

It get data every refresh second and send metrics every flush seconds or limit records. 

```
Usage of rancher-telemetry-stats:
  -accessKey string
      Rancher access key. Or env TELEMETRY_ACCESS_KEY (default "")
  -format string
      Output format. influx | json (default "influx")
  -flush int
      Send metrics to inflush every flush seconds. (default 60)
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
      Get metrics every refresh seconds. (default 3600)
  -secretKey string
      Rancher secret key. Or env TELEMETRY_SECRET_KEY (default "")
  -url string
      Rancher telemetry url. (default "http://telemetry.rancher.io")
```

NOTE: You need influx already installed and running. The influx db would be created if doesn't exist.

## Metrics

SELECT id, uid, ts, data FROM record WHERE date_trunc('day',ts) >= (date_trunc('day',now()) - INTERVAL '2 day') ORDER BY id DESC

SELECT r.id, r.uid, r.ts, r.data, i.last_seen, i.last_ip FROM record r JOIN installation i ON (i.last_record = r.id) WHERE date_trunc('day',ts) >= (date_trunc('day',now()) - INTERVAL '7 day') ORDER BY id DESC

SELECT i.id, i.uid, i.first_seen, r.ts, i.last_ip, r.data FROM installation i JOIN record r ON (i.uid = r.uid) WHERE date_trunc('day',ts) >= (date_trunc('day',now()) - INTERVAL '7 day') ORDER BY id DESC

select count(r.uid) FROM installation i JOIN record r ON (i.uid = r.uid) WHERE date_trunc('day',ts) >= (date_trunc('day',now()) - INTERVAL '7 day') ORDER BY id DESC

WHERE i.last_seen >= NOW() - INTERVAL '7 days'

Metrics are on the form.....

```
telemetry,city=city,country=country,country_isocode=country_isocode,firstseen=2017-05-10\ 17:17:53.713117\ +0000\ UTC,id=XXXX,install_image=rancher/server,install_version=v1.6.0,lastseen=2017-05-10\ 17:17:53.713117\ +0000\ UTC,uid=f186b6a5-62dd-4753-b02a-f44c14352e8e container_running=0,container_total=0,environment_total=1,host_active=0,host_cpu_cores_total=0,host_mem_mb_total=0,ip="XX.XX.XX.XX",orch_cattle=1,orch_kubernetes=0,orch_mesos=0,orch_swarm=0,orch_windows=0,service_active=0,service_total=5,stack_active=4,stack_total=4,uid="XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXX" 1494436673713117000
```

[alpine-base]: https://github.com/rawmind0/alpine-base

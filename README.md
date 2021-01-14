[![](https://images.microbadger.com/badges/image/rawmind/rancher-telemetry-stats.svg)](https://microbadger.com/images/rawmind/rancher-telemetry-stats "Get your own image badge on microbadger.com")

rancher-telemetry-stats
=====================

This image run rancher-telemetry-stats app. It comes from [rawmind/alpine-base][alpine-base].

## Build

```
docker build -t rawmind/rancher-telemetry-stats:<version> .
```

## Versions

- `0.2-8` [(Dockerfile)](https://github.com/rawmind0/rancher-telemetry-stats/blob/0.2-8/Dockerfile)
- `0.1-3` [(Dockerfile)](https://github.com/rawmind0/rancher-telemetry-stats/blob/0.1-3/Dockerfile)


## Usage

This image run rancher-telemetry-stats service. rancher-telemetry-stats get metrics from rancher telemetry service and send them to a influx in order to be explored by grafana. 

It get data every refresh second and send metrics every flush seconds or limit records. 

```
NAME:
   telemetry stats - Rancher telemetry stats

USAGE:
   test [global options] command [command options] [arguments...]

VERSION:
   git

AUTHOR:
   Rancher Labs, Inc.

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug             debug logging [$TELEMETRY_DEBUG]
   --log value         path to log to [$TELEMETRY_LOG]
   --pid-file value    path to write PID to [$TELEMETRY_PID_FILE]
   --url value         url to reach telemetry (default: "https://telemetry.rancher.io") [$TELEMETRY_URL]
   --hours value       telemetry hours to get (default: 48)
   --accesskey value   access key for api [$TELEMETRY_ACCESS_KEY]
   --secretkey value   secret key for api [$TELEMETRY_SECRET_KEY]
   --format value      Output format. influx | json (default: "influx")
   --preview           Just print output to stdout
   --influxurl value   Influx url connection (default: "http://localhost:8086")
   --influxdb value    Influx db name (default: "telemetry")
   --influxuser value  Influx username
   --influxpass value  Influx password
   --insecure          Allow insecure connection to telemetry
   --geoipdb value     Geoip db file. (default: "GeoLite2-City.mmdb")
   --file value        Read requests from file.
   --limit value       Limit batch size (default: 2000)
   --refresh value     Get metrics every refresh seconds (default: 3600)
   --flush value       Send metrics to inflush every flush seconds (default: 60)
   --help, -h          show help
   --version, -v       print the version
```

NOTE: You need influx already installed and running. The influx db would be created if doesn't exist.

## Metrics

Metrics are on the form...

Telemetry record version 1
```
telemetry,city=city,country=country,country_isocode=country_isocode,id=XXXX,install_image=rancher/server,install_version=v1.6.0,record_version=1,status=new,uid=XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXX container_running=0,container_total=0,environment_total=1,host_active=0,host_cpu_cores_total=0,host_mem_mb_total=0,ip="XX.XX.XX.XX",orch_cattle=1,orch_kubernetes=0,orch_mesos=0,orch_swarm=0,orch_windows=0,service_active=0,service_total=5,stack_active=4,stack_total=4,stack_from_catalog=4,uid="XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXX" 1494436673713117000
```

Telemetry record version 2
```
telemetry,city=city,country=country,country_isocode=country_isocode,id=XXXX,install_image=rancher/server,install_version=v2.0.0-alpha16,ip=XX.XX.XX.XX,record_version=2,status=new,uid=XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXX cluster_active=1,cluster_cloud_provider_aws=0,cluster_cloud_provider_azure=0,cluster_cloud_provider_custom=0,cluster_cloud_provider_openstack=0,cluster_cloud_provider_vsphere=1,cluster_cpu_total=3,cluster_cpu_util=25,cluster_driver_aks=0,cluster_driver_eks=0,cluster_driver_gke=0,cluster_driver_imported=1,cluster_driver_k3s=0,cluster_driver_k3sBased=0,cluster_mem_mb_total=2676,cluster_mem_util=12,cluster_namespace_from_catalog=1,cluster_namespace_total=5,cluster_total=1,ip="XX.XX.XX.XX",node_active=3,node_from_template=3,node_mem_mb_total=2676,node_mem_util=12,node_role_controlplane=1,node_role_etcd=1,node_role_worker=2,node_total=3,project_namespace_from_catalog=1,project_namespace_total=2,project_pod_total=2,project_total=1,project_workload_total=1,uid="XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXX" 1496359077329645000
telemetry_apps,catalog=system-library,id=33221669,name=rancher-monitoring,uid=XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXX,version=0.1.1 total=4,uid="XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXX" 1496359077329645000
telemetry_drivers,kind=cluster,name=imported,uid=XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXX total=1,uid="XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXX" 1496359077329645000
```

[alpine-base]: https://github.com/rawmind0/alpine-base

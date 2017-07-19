package main

import (
	"os"
	"flag"
	log "github.com/Sirupsen/logrus"
)

func check(e error, m string) {
    if e != nil {
		log.Error("[Error]: ", m , e)
	}
}

type Params struct {
	url string
	accessKey string
	secretKey string
	influxurl string
	influxdb string
	influxuser string
	influxpass string
	geoipdb string
	file string
	format string
	limit int
	refresh int
	flush int
}

func (p *Params) init() {

	flag.StringVar(&p.url, "url", "http://telemetry.rancher.io", "Rancher telemetry url.")
	flag.StringVar(&p.accessKey, "accessKey", os.Getenv("TELEMETRY_ACCESS_KEY"), "Rancher access key. Or env TELEMETRY_ACCESS_KEY")
	flag.StringVar(&p.secretKey, "secretKey", os.Getenv("TELEMETRY_SECRET_KEY"), "Rancher secret key. Or env TELEMETRY_SECRET_KEY")
	flag.StringVar(&p.format, "format", "influx", "Output format. influx | json")
	flag.StringVar(&p.influxurl, "influxurl", "http://localhost:8086", "Influx url connection")
	flag.StringVar(&p.influxdb, "influxdb", "telemetry", "Influx db name")
	flag.StringVar(&p.influxuser, "influxuser", "", "Influx username")
	flag.StringVar(&p.influxpass, "influxpass", "", "Influx password")
	flag.StringVar(&p.geoipdb, "geoipdb", "GeoLite2-City.mmdb", "Geoip db file.")
	flag.StringVar(&p.file, "file", "", "Read requests from file.")
	flag.IntVar(&p.limit, "limit", 2000, "Limit batch size")
	flag.IntVar(&p.refresh, "refresh", 3600, "Get metrics every refresh seconds.")
	flag.IntVar(&p.flush, "flush", 60, "Send metrics to inflush every flush seconds.")

	flag.Parse()

	p.checkParams()
}

func (p *Params) checkParams() {
	if p.format != "influx" && p.format != "json"{
		flag.Usage()
		log.Info("Check your format params. influx | json ")
		os.Exit(1) 
	}
	if p.format == "influx" {
		if ( len(p.influxdb) == 0 || len(p.influxurl) == 0 ) { 
			flag.Usage()
			log.Info("Check your influxdb and/or influxurl params.")
			os.Exit(1) 
		}
	}
}
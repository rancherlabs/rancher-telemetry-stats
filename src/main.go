package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	VERSION string
)

func main() {
	app := cli.NewApp()
	app.Name = "telemetry stats"
	app.Author = "Rancher Labs, Inc."
	app.Usage = "Rancher telemetry stats"

	if VERSION == "" {
		app.Version = "git"
	} else {
		app.Version = VERSION
	}

	flags := []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug logging",
			EnvVar: "TELEMETRY_DEBUG",
		},

		cli.StringFlag{
			Name:   "log",
			Usage:  "path to log to",
			Value:  "",
			EnvVar: "TELEMETRY_LOG",
		},

		cli.StringFlag{
			Name:   "pid-file",
			Usage:  "path to write PID to",
			Value:  "",
			EnvVar: "TELEMETRY_PID_FILE",
		},

		cli.StringFlag{
			Name:   "url",
			Usage:  "url to reach telemetry",
			Value:  "https://telemetry.rancher.io",
			EnvVar: "TELEMETRY_URL",
		},

		cli.IntFlag{
			Name:  "hours",
			Usage: "telemetry hours to get",
			Value: 48,
		},

		cli.StringFlag{
			Name:   "accesskey",
			Usage:  "access key for api",
			Value:  "",
			EnvVar: "TELEMETRY_ACCESS_KEY",
		},

		cli.StringFlag{
			Name:   "secretkey",
			Usage:  "secret key for api",
			Value:  "",
			EnvVar: "TELEMETRY_SECRET_KEY",
		},

		cli.StringFlag{
			Name:  "format",
			Usage: "Output format. influx | json",
			Value: "influx",
		},

		cli.BoolFlag{
			Name:  "preview",
			Usage: "Just print output to stdout",
		},

		cli.StringFlag{
			Name:  "influxurl",
			Usage: "Influx url connection",
			Value: "http://localhost:8086",
		},

		cli.StringFlag{
			Name:  "influxdb",
			Usage: "Influx db name",
			Value: "telemetry",
		},

		cli.StringFlag{
			Name:  "influxuser",
			Usage: "Influx username",
			Value: "",
		},

		cli.StringFlag{
			Name:  "influxpass",
			Usage: "Influx password",
			Value: "",
		},

		cli.BoolFlag{
			Name:  "insecure",
			Usage: "Allow insecure connection to telemetry",
		},

		cli.StringFlag{
			Name:  "geoipdb",
			Usage: "Geoip db file.",
			Value: "GeoLite2-City.mmdb",
		},

		cli.StringFlag{
			Name:  "file",
			Usage: "Read requests from file.",
			Value: "",
		},

		cli.IntFlag{
			Name:  "limit",
			Usage: "Limit batch size",
			Value: 2000,
		},

		cli.IntFlag{
			Name:  "refresh",
			Usage: "Get metrics every refresh seconds",
			Value: 3600,
		},

		cli.IntFlag{
			Name:  "flush",
			Usage: "Send metrics to inflush every flush seconds",
			Value: 60,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "run",
			Description: "continuously migrate recent data",
			Action:      RunRequests,
			Flags:       flags,
		},
		{
			Name:        "restore",
			Description: "restore data from a specific date or range and exit",
			Action: func(c *cli.Context) {
				cli.NewExitError(RunRequests(c), 1)
			},
			Flags: append(flags,
				cli.StringFlag{
					Name:     "dates",
					Usage:    "dates in the format 2022-12-24, single value or separated by minus or comma",
					Required: true,
				},
			),
		},
	}

	app.Flags = flags
	app.Before = before
	app.Action = RunRequests

	app.Run(os.Args)
}

func before(c *cli.Context) error {
	timeFormat := new(log.TextFormatter)
	timeFormat.TimestampFormat = "2006-01-02 15:04:05"
	timeFormat.FullTimestamp = true
	log.SetFormatter(timeFormat)

	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	logFile := c.String("log")
	if logFile != "" {
		output, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			str := fmt.Sprintf("Failed to log to file %s: %v", logFile, err)
			return cli.NewExitError(str, 1)
		}
		log.SetOutput(output)
	}

	pidFile := c.String("pid-file")
	if pidFile != "" {
		log.Infof("Writing pid %d to %s", os.Getpid(), pidFile)
		if err := ioutil.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0644); err != nil {
			str := fmt.Sprintf("Failed to write pid file %s: %v", pidFile, err)
			return cli.NewExitError(str, 1)
		}
	}

	format := c.String("format")
	if format != "influx" && format != "json" {
		return cli.NewExitError("Check your format params. influx | json ", 1)
	}

	influxdb := c.String("influxdb")
	influxurl := c.String("influxurl")
	if format == "influx" {
		if len(influxdb) == 0 || len(influxurl) == 0 {
			return cli.NewExitError("Check your influxdb and/or influxurl params.", 1)
		}
	}

	_, err := os.Stat(c.String("geoipdb"))
	if err != nil {
		str := fmt.Sprintf("Failed to open geoipdb file %s: %v", c.String("geoipdb"), err)
		return cli.NewExitError(str, 1)
	}

	return nil
}

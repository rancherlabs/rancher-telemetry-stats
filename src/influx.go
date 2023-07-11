package main

import (
	"fmt"
	"time"

	_ "github.com/influxdata/influxdb1-client"
	influx "github.com/influxdata/influxdb1-client/v2"
	log "github.com/sirupsen/logrus"
)

func check(e error, m string) {
	if e != nil {
		log.Error("[Error]: ", m, e)
	}
}

type Influx struct {
	url     string
	db      string
	user    string
	pass    string
	client  influx.Client
	batch   influx.BatchPoints
	timeout time.Duration
}

func newInflux(url, db, user, pass string) *Influx {
	var a = &Influx{
		url:  url,
		db:   db,
		user: user,
		pass: pass,
	}

	a.timeout = time.Duration(10)
	return a
}

func (i *Influx) Check(retry int) bool {
	respTime, _, err := i.client.Ping(i.timeout)
	if err != nil {
		connected := i.Init()
		for index := 1; index <= retry && !connected; index++ {
			log.Warn("Influx disconnected. Reconnecting ", index+1, " of ", retry, "...")
			time.Sleep(time.Duration(1) * time.Second)
			connected = i.Init()
			if connected {
				respTime, _, err = i.client.Ping(i.timeout)
			}
		}
		if err != nil {
			log.Error("Failed to connect to influx ", i.url)
			return false
		}
	}
	log.Debug("Influx response time: ", respTime)
	return true
}

// CheckConnect creates a ticker with an interval. On every tick the connection
// to InfluxDB is checked. The check also tries to re-establish the connection.
// If it fails to do so, the returned channel will be closed to signal, that the
// connection has been lost.
func (i *Influx) CheckConnect(interval int) chan bool {
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	disconnected := make(chan bool)

	go func() {
		defer ticker.Stop()
		for range ticker.C {
			if !i.Check(2) {
				close(disconnected)
				return
			}
		}
	}()

	return disconnected
}

// Init creates a new HTTP client and checks the connection. If it can connect,
// it makes sure necessary databases are created.
func (i *Influx) Init() bool {
	var err error
	log.Info("Connecting to Influx...")

	i.client, err = influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     i.url,
		Username: i.user,
		Password: i.pass,
	})

	if err != nil {
		log.Error("[Error]: ", err)
		return false
	}

	if i.Check(0) {
		i.createDb()
		return true
	} else {
		return false
	}
}

func (i *Influx) Close() {
	message := "Closing Influx connection..."
	err := i.client.Close()
	check(err, message)
	log.Info(message)
}

func (i *Influx) createDb() {
	var err error
	log.Info("Creating Influx database if not exists...")

	comm := "CREATE DATABASE " + i.db

	q := influx.NewQuery(comm, "", "")
	_, err = i.client.Query(q)
	if err != nil {
		log.Error("[Error] ", err)
	} else {
		log.Info("Influx database ", i.db, " created.")
	}
}

func (i *Influx) newBatch() {
	var err error
	message := "Creating Influx batch..."
	i.batch, err = influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  i.db,
		Precision: "s",
	})
	check(err, message)
	log.Info(message)
}

func (i *Influx) newPoint(m influx.Point) {
	message := "Adding point to batch..."
	fields, _ := m.Fields()
	pt, err := influx.NewPoint(m.Name(), m.Tags(), fields, m.Time())
	check(err, message)
	i.batch.AddPoint(pt)
}

// addPoints adds new points to the current batch.
func (i *Influx) addPoints(m []influx.Point) {
	log.Info("Adding ", len(m), " points to batch...")
	for _, p := range m {
		i.newPoint(p)
	}
}

func (i *Influx) Write() {
	start := time.Now()
	log.Info("Writing batch points...")

	// Write the batch
	err := i.client.Write(i.batch)
	if err != nil {
		log.Error("[Error]: ", err)

	}

	log.Info("Time to write ", len(i.batch.Points()), " points: ", float64((time.Since(start))/time.Millisecond), "ms")
}

func (i *Influx) deleteForDay(day string) bool {
	_, err := time.Parse("2006-01-02", day)
	if err != nil {
		log.Warnf("Deleting entries for day %s failed: day not formatted correctly", day)
		return false
	}

	for _, measurement := range [3]string{"telemetry", "telemetry_apps", "telemetry_drivers"} {
		influxQL := fmt.Sprintf(
			`delete from %s where time >= '%s' and time < '%s' + 1d;`,
			measurement,
			day,
			day,
		)
		_, err = i.doQuery(influxQL, 3)
		if err != nil {
			log.Warnf(
				"Deleting entries for day %s in measurement %s failed: %s",
				day,
				measurement,
				err,
			)
			return false
		}
	}

	return true
}

func (i *Influx) sendToInflux(m []influx.Point, retry int) bool {
	if i.Check(retry) {
		i.newBatch()
		i.addPoints(m)
		i.Write()
		return true
	}
	return false
}

func (i *Influx) doQuery(queryString string, retry int) (*influx.Response, error) {
	if len(queryString) > 0 && i.Check(retry) {
		return i.client.Query(influx.NewQuery(queryString, i.db, "s"))
	}
	return nil, nil
}

package main

import (
	"fmt"
	//"regexp"
	"bytes"
	"strings"
	"io/ioutil"
	"time"
	"encoding/json"
	"os"
	"net/http"
	"strconv"
	"os/signal"
	"sync"
	"net"
	log "github.com/Sirupsen/logrus"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/oschwald/maxminddb-golang"
)

func getJSON(url string, accessKey string, secretKey string, target interface{}) error {

	start := time.Now()

	log.Info("Connecting to ", url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(accessKey, secretKey)
	resp, err := client.Do(req)

	if err != nil {
		log.Error("Error Collecting JSON from API: ", err)
		panic(err)
	}

	respFormatted := json.NewDecoder(resp.Body).Decode(target)

	// Timings recorded as part of internal metrics
 	log.Info("Time to get json: ", float64((time.Since(start))/ time.Millisecond), " ms")

	// Close the response body, the underlying Transport should then close the connection.
	resp.Body.Close()

	// return formatted JSON
	return respFormatted
}

type reqLocation struct {
	City          	string `json:"city"`
	Country        	struct {
    	Name 		string 
       	ISOCode 	string 
    } `json:"country"`
} 

type RequestData struct {
	container map[string]float64 `json:"container"`
	environment map[string]float64 `json:"environment"`
    host map[string]float64 `json:"host"`
    install map[string]string `json:"install"`
    service map[string]float64 `json:"service"`
    stack map[string]float64 `json:"stack"`
    ts time.Time
}

type Request struct {
	Id        int64       `json:"id"`
	Uid       string      `json:"uid"`
	FirstSeen time.Time   `json:"first_seen"`
	LastSeen  time.Time   `json:"last_seen"`
	LastIp    string      `json:"last_ip"`
	Record    interface{} `json:"record"`
	Data 	  *RequestData
	Location  reqLocation `json:"location"`	
	Status 	  string 	  `json:"status"`
}

func newrequestByString(line, sep string) (*Request) {
	var err error

	s := strings.Split(line, sep)
	if len(s) == 6 {
		req := &Request{}

		req.Id, err = strconv.ParseInt(s[0], 10, 64)
		if err != nil {
			log.Info("Error parsing Id ", s[0], err)
		}
		req.Uid = s[1]
		req.FirstSeen , err = time.Parse("2006-01-02 15:04:05.999999", s[2])
		if err != nil {
			log.Info("Error parsing Firstseen ", s[2], err)
		}
		req.LastSeen , err = time.Parse("2006-01-02 15:04:05.999999", s[3])
		if err != nil {
			log.Info("Error parsing Lastseen ", s[3], err)
		}
		req.LastIp = s[4]
		err = json.Unmarshal([]byte(s[5]), &req.Record)
		if err != nil {
			log.Info("Error unmarshaling Record ", s[5], err)
		}

		return req
	}
	return nil
}

func (r *Request) initData() bool {

	if len(r.Uid) > 2 { 
		r.Data = &RequestData{}

		r.Data.container = make(map[string]float64)
		r.Data.environment = make(map[string]float64)
	    r.Data.host = make(map[string]float64)
	    r.Data.install = make(map[string]string)
	    r.Data.service = make(map[string]float64)
	    r.Data.stack = make(map[string]float64)

	    return true
	} else {
		return false
	}
}

// Produce wire-formatted string for ingestion into influxdb
func (r *Request) printInflux() {
	p := r.getPoint()
	fmt.Println(p.String())
}

func (r *Request) getPoint() *influx.Point {
	var n = "telemetry"
    v := map[string]interface{}{
        "ip": r.LastIp,
        "uid": r.Uid,
        "container_running": r.Data.container["running"],
        "container_total": r.Data.container["total"],
        "environment_total": r.Data.environment["total"],
        "orch_cattle": r.Data.environment["cattle"],
        "orch_kubernetes": r.Data.environment["kubernetes"],
        "orch_mesos": r.Data.environment["mesos"],
        "orch_swarm": r.Data.environment["swarm"],
        "orch_windows": r.Data.environment["windows"],
        "host_active": r.Data.host["active"],
        "host_cpu_cores_total": r.Data.host["cpu_cores_total"],
        "host_mem_mb_total": r.Data.host["mem_mb_total"],
        "service_active": r.Data.service["active"],
        "service_total": r.Data.service["total"],
        "stack_active": r.Data.stack["active"],
        "stack_total": r.Data.stack["total"],
    }
    t := map[string]string{
    	"id":  strconv.FormatInt(r.Id, 10),
        "uid": r.Uid,
        "ip": r.LastIp,
		"install_image": r.Data.install["image"],
		"install_version": r.Data.install["version"],
        "city": r.Location.City,
        "country": r.Location.Country.Name,
        "country_isocode": r.Location.Country.ISOCode,
        "status": r.Status,
    }

	m, err := influx.NewPoint(n,t,v,r.Data.ts)
	if err != nil {
		log.Warn(err)
	}

	return m
} 

func (r *Request) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Request) printJson() {

	j, err := json.Marshal(r)
	if err != nil {
    	log.Error("json", err)
	}

	fmt.Println(string(j))

}

func (r *Request) getLocation(geoipdb string) {
    db, err := maxminddb.Open(geoipdb)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    ip := net.ParseIP(r.LastIp)

    var record struct {
        City    struct {
            Names map[string]string `maxminddb:"names"`
        } `maxminddb:"city"`
        Country struct {
            Names map[string]string `maxminddb:"names"`
            ISOCode string `maxminddb:"iso_code"`
        } `maxminddb:"country"`
    } // Or any appropriate struct

    err = db.Lookup(ip, &record)
    if err != nil {
        log.Fatal(err)
    }
    
    r.Location.City = record.City.Names["en"] 
    r.Location.Country.Name = record.Country.Names["en"]
    r.Location.Country.ISOCode = record.Country.ISOCode

}

func (r *Request) isNew() string{
	diff := r.LastSeen.Sub(r.FirstSeen)

	if diff.Hours() > 24 {
		return "active"
	} else {
		return "new"
	} 
}

func (r *Request) getData(geoipdb string) bool{

	if r.initData() {
		r.getLocation(geoipdb)

		m := r.Record.(map[string]interface{})

		if m != nil {

			if m["container"] != nil {
				container := m["container"].(map[string]interface{})
				if container["running"] != nil {
					r.Data.container["running"] = container["running"].(float64)
				}
				if container["total"] != nil {
					r.Data.container["total"] = container["total"].(float64)
				}
			}

			if m["environment"] != nil {
				environment := m["environment"].(map[string]interface{})
				if environment["total"] != nil {
					r.Data.environment["total"] = environment["total"].(float64)
				}
				if environment["orch"] != nil {
					orch := environment["orch"].(map[string]interface{})
					if orch["cattle"] != nil {
						r.Data.environment["cattle"] = orch["cattle"].(float64)
					}
					if orch["kubernetes"] != nil {
						r.Data.environment["kubernetes"] = orch["kubernetes"].(float64)
					}
					if orch["mesos"] != nil {
						r.Data.environment["mesos"] = orch["mesos"].(float64)
					}
					if orch["swarm"] != nil {
						r.Data.environment["swarm"] = orch["swarm"].(float64)
					}
					if orch["windows"] != nil {
						r.Data.environment["windows"] = orch["windows"].(float64)
					}
				}
			}

			if m["service"] != nil {
				service := m["service"].(map[string]interface{})
				if service["total"] != nil {
					r.Data.service["total"] = service["total"].(float64)
				}
				if service["active"] != nil {
					r.Data.service["active"] = service["active"].(float64)
				}
			}

			if m["stack"] != nil {
				stack := m["stack"].(map[string]interface{})
				if stack["total"] != nil {
					r.Data.stack["total"] = stack["total"].(float64)
				}
				if stack["active"] != nil {
					r.Data.stack["active"] = stack["active"].(float64)
				}
			}

			if m["install"] != nil {
				install := m["install"].(map[string]interface{})
				if install["image"] != nil {
					r.Data.install["image"] = install["image"].(string)
				}
				if install["version"] != nil {
					r.Data.install["version"] = install["version"].(string)
				}
			}

			if m["host"] != nil {
				host := m["host"].(map[string]interface{})
				if host["cpu"] != nil {
					cpu := host["cpu"].(map[string]interface{})
					if cpu["cores_total"] != nil {
						r.Data.host["cpu_cores_total"] = cpu["cores_total"].(float64)
					}
				}
				if host["mem"] != nil {
					mem := host["mem"].(map[string]interface{})
					if mem["mb_total"] != nil {
						r.Data.host["mem_mb_total"] = mem["mb_total"].(float64)
					}
				}
				var host_active, host_count float64
				if host["active"] != nil {
					host_active = host["active"].(float64)
				}
				if host["count"] != nil {
					host_count = host["count"].(float64)
				}

				if host_active < host_count {
					host_active = host_count
				}

				r.Data.host["active"] = host_active
			}

			now := time.Now()
     		diff := now.Sub(r.LastSeen)

     		if diff < 0 {
     			log.Warn("Time in future ", r.LastSeen)
     			r.printJson()
     			return false
			} 
			r.Data.ts = r.LastSeen

			r.Status = r.isNew()

			return true
		}
	} 
	return false

}

type Requests struct {
	Input 			chan *Request
	Output 			chan *influx.Point
	Exit 			chan os.Signal
	Readers			[]chan struct{}
	Config 			Params
	Data 			[]Request `json:"data"`
}

func newRequests(conf Params) *Requests {
	var r = &Requests{
		Readers: []chan struct{}{},
		Config:	conf,
	}

	r.Input = make(chan *Request,1)
	r.Output = make(chan *influx.Point,1)
	r.Exit = make(chan os.Signal, 1)
	signal.Notify(r.Exit, os.Interrupt, os.Kill)

	customFormatter := new(log.TextFormatter)
    customFormatter.TimestampFormat = "2006-01-02 15:04:05"
    log.SetFormatter(customFormatter)
    customFormatter.FullTimestamp = true

	return r
}

func (r *Requests) Close(){
	close(r.Input)
	close(r.Output)
	close(r.Exit)

}

func (r *Requests) sendToInflux() {
	var points []influx.Point
	var index,p_len int
	
	i := newInflux(r.Config.influxurl, r.Config.influxdb, r.Config.influxuser, r.Config.influxpass)

	if i.Connect() {
		connected := i.CheckConnect(r.Config.refresh)
		defer i.Close()

		ticker := time.NewTicker(time.Second * time.Duration(r.Config.flush))

		index = 0
		for {
	        select {
	        case <-connected:
	        	return
	        case <-ticker.C:
	        	if len(points) > 0 {
	            	log.Info("Tick: Sending to influx ", len(points), " points")
	    			if i.sendToInflux(points,1) {
	    				points = []influx.Point{}
	    			} else {
	    				return
	    			}
	    		} 
	        case p := <- r.Output:
	        	if p != nil {
	        		points = append(points, *p)
	        		p_len = len(points)
	        		if p_len == r.Config.limit {
	            		log.Info("Batch: Sending to influx ", p_len, " points")
	            		if i.sendToInflux(points,1) {
	    					points = []influx.Point{}
	    				} else {
	    					return
	    				}
	            	}
	            	index++
	        	} else {
	        		p_len = len(points)
	        		if p_len > 0 {
	        			log.Info("Batch: Sending to influx ", p_len, " points")
	            		if i.sendToInflux(points,1) {
	    					points = []influx.Point{}
	    				}
	        		}
	        		return
	        	}
	        }
	    }
	} 
}

func (r *Requests) addReader() chan struct{} {
	chan_new := make(chan struct{}, 1)
	r.Readers = append(r.Readers, chan_new)

	return chan_new
}

func (r *Requests) closeReaders() {
	for _, r_chan := range r.Readers {
		if r_chan != nil {
			r_chan <- struct{}{}
		}
	}
	r.Readers = nil
}

func (r *Requests) getDataByUrl() {
	var in, out sync.WaitGroup
	indone := make(chan struct{},1)
	outdone := make(chan struct{},1)

	in.Add(1)
	go func() {
		defer in.Done()
		if r.Config.file != "" {
			r.getDataByFile()
		} else {
			r.getDataByChan(r.addReader())
		}
	}()

	out.Add(1)
	go func() {
		defer out.Done()
		r.getOutput()
	}()

	go func() {
		in.Wait()
		close(r.Input)
		close(r.Output)
		close(indone)
	}()

	go func() {
		out.Wait()
		close(outdone)
	}()

	for {
        select {
        case <- indone:
        	<- outdone
        	return
        case <- outdone:
        	log.Error("Aborting...")
        	go r.closeReaders()
        	return
        case <- r.Exit:
        	//close(r.Exit)
        	log.Info("Exit signal detected....Closing...")
        	go r.closeReaders()
        	select {
        	case <- outdone:
        		return
        	}
        }
    }
}

func (r *Requests) getDataByChan(stop chan struct{}) {

	ticker := time.NewTicker(time.Second * time.Duration(r.Config.refresh))

	go r.getData()

	for {
        select {
        case <-ticker.C:
        		log.Info("Tick: Getting data...")
            	go r.getData()
        case <- stop:
            return
        }
    }
}

func (r *Requests) getDataByFile() {
	log.Info("Reading file ", r.Config.file)
	dat, err := ioutil.ReadFile(r.Config.file)

	if err != nil {
		log.Error("Error reading file ", r.Config.file , err)
	}

	lines := bytes.Split(dat, []byte("\n"))

	for _ , line := range lines {
		if line != nil {
			req := newrequestByString(string(line), "|")
			if req != nil {
				r.Data = append(r.Data, *req)
			}
		}
	}

	for _ , req := range r.Data {
		aux := Request{}
		aux = req
		if aux.getData(r.Config.geoipdb) {
			if r.Config.format == "json" {
				r.Input <- &aux
			} else {
				r.Output <- aux.getPoint()
			}
		}
    }
}

func (r *Requests) getData() {
	var uri string

	uri = "/admin/active"

	err := getJSON(r.Config.url+uri, r.Config.accessKey, r.Config.secretKey, r)
	if err != nil {
		log.Error("Error getting JSON from URL ", r.Config.url+uri , err)
	}

	for _ , req := range r.Data {
		aux := Request{}
		aux = req
		if aux.getData(r.Config.geoipdb) {
			if r.Config.format == "json" {
				r.Input <- &aux
			} else {
				r.Output <- aux.getPoint()
			}
		}
    }
}

func (r *Requests) print() {
	if r.Config.format == "json" {
		r.printJson()
	} else {
		r.printInflux()
	}
}

func (r *Requests) getOutput() {
	if r.Config.format == "json" {
		r.printJson()
	} else {
		r.sendToInflux()
		//r.printInflux()
	}
}

func (r *Requests) printJson() {
	for {
        select {
        case req := <- r.Input:
        	if req != nil {
            	req.printJson()
        	} else {
        		return
        	}
        }
    }
}

func (r *Requests) printInflux() {
	for {
        select {
        case p := <- r.Output:
        	if p != nil {
            	fmt.Println(p.String())
        	} else {
        		return
        	}
        }
    }
}



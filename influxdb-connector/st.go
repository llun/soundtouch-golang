package main

import (
	"fmt"
	"net"
	"net/url"
	"reflect"
	"sync"

	"github.com/jpillora/opts"
	log "github.com/sirupsen/logrus"

	soundtouch "github.com/theovassiliou/soundtouch-golang"
)

// var soundtouchNetwork = make(map[string]string)

var conf = config{}

//set this via ldflags (see https://stackoverflow.com/q/11354518)
var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

const shortUsage = "Captures broadcastet information from your Bose Soundtouch systems."

var influxDB = soundtouch.InfluxDB{
	BaseHTTPURL: url.URL{
		Scheme: "http",
		Host:   "localhost:8086",
	},
	Database:          "soundtouch",
	SoundtouchNetwork: make(map[string]string),
}

type config struct {
	Speakers  []string  `opts:"group=Soundtouch" help:"Speakers to listen for, all if not set"`
	Interface string    `opts:"group=Soundtouch" help:"network interface to listen"`
	InfluxURL string    `opts:"group=InfluxDB" help:"URL of the influx database"`
	Database  string    `opts:"group=InfluxDB" help:"InfluxDB database to send the data to"`
	DryRun    bool      `help:"Dump the lineprotocoll in curl format instead sending to influxdb"`
	LogLevel  log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
}

func main() {
	conf = config{
		LogLevel:  log.DebugLevel,
		InfluxURL: "http://influxdb:8086",
		Database:  "soundtouch",
		Interface: "en0",
	}

	//parse config
	opts.New(&conf).
		Summary(shortUsage).
		Repo("github.com/theovassiliou/soundtouch-golang").
		Version(VERSION).
		Parse()

	log.SetLevel(conf.LogLevel)

	v, err := url.Parse(conf.InfluxURL)
	if err != nil {
		log.Fatalf("Not a valid URL: %v", conf.InfluxURL)
	}

	influxDB.BaseHTTPURL = *v
	influxDB.Database = conf.Database

	i, err := net.InterfaceByName(conf.Interface)

	if err != nil {
		log.Fatalf("Error with interface. %s", err)
	}

	log.Debugf("Listening @ %v, supports: %v, HW Address: %v\n", i.Name, i.Flags.String(), i.HardwareAddr)

	speakerCh := soundtouch.Lookup(i)
	var wg sync.WaitGroup
	messageCh := make(chan *soundtouch.Update)

	for speaker := range speakerCh {
		di, _ := speaker.Info()
		speaker.DeviceInfo = di

		spkLogger := log.WithFields(log.Fields{
			"Speaker": speaker.DeviceInfo.Name,
			"ID":      speaker.DeviceInfo.DeviceID,
		})

		if len(conf.Speakers) > 0 && !isIn(di.Name, conf.Speakers) {
			spkLogger.Traceln("Ignoring messages from: ", di.Name)
			continue
		}

		influxDB.SoundtouchNetwork[di.DeviceID] = di.Name
		spkLogger.Infof("Listening\n")
		spkLogger.Debugf(" with IP: %v", speaker.IP)
		wg.Add(1)
		go func(s *soundtouch.Speaker, msgChan chan *soundtouch.Update) {
			defer wg.Done()

			webSocketCh, _ := s.Listen()
			for message := range webSocketCh {
				msgChan <- message
			}
		}(speaker, messageCh)

	}
	for m := range messageCh {
		mLogger := log.WithFields(log.Fields{
			"Speaker": m.DeviceID,
			"Value":   reflect.TypeOf(m.Value).Name(),
		})
		v, _ := m.Lineproto(influxDB, m)
		if !conf.DryRun && v != "" {
			result, err := influxDB.SetData("write", []byte(v))
			if err != nil {
				mLogger.Errorf("failed")
			}
			mLogger.Debugf("succeeded: %v", string(result))

		} else if v != "" {
			fmt.Printf("curl -i -XPOST \"%v\" --data-binary '%v'\n", influxDB.WriteURL("write"), v)
		}
	}
	wg.Wait()
}
func isIn(name string, selected []string) bool {
	for _, s := range selected {
		if name == s {
			return true
		}
	}
	return false
}

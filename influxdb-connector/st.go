package main

import (
	"fmt"
	"net"
	"net/url"
	"reflect"
	"sync"
	"time"

	"github.com/jpillora/opts"
	log "github.com/sirupsen/logrus"

	soundtouch "github.com/theovassiliou/soundtouch-golang"
	"github.com/theovassiliou/soundtouch-golang/magiczone/magicspeaker"
)

var conf = config{}

//set this via ldflags (see https://stackoverflow.com/q/11354518)
var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

const shortUsage = "Captures broadcastet information from your Bose Soundtouch systems."

type speakerMap map[string]bool

var influxDB = soundtouch.InfluxDB{
	BaseHTTPURL: url.URL{
		Scheme: "http",
		Host:   "localhost:8086",
	},
	Database:          "soundtouch",
	SoundtouchNetwork: make(map[string]string),
}

type config struct {
	Speakers            []string  `opts:"group=Soundtouch" help:"Speakers to listen for, all if not set"`
	Interface           string    `opts:"group=Soundtouch" help:"network interface to listen"`
	NoSoundtouchSystems int       `opts:"group=Soundtouch" help:"Number of Soundtouch systems to scan for."`
	InfluxURL           string    `opts:"group=InfluxDB" help:"URL of the influx database"`
	Database            string    `opts:"group=InfluxDB" help:"InfluxDB database to send the data to"`
	DryRun              bool      `help:"Dump the lineprotocoll in curl format instead sending to influxdb"`
	LogLevel            log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
}

var visibleSpeakers = make(magicspeaker.MagicSpeakers)

func main() {
	conf = config{
		NoSoundtouchSystems: -1,
		LogLevel:            log.DebugLevel,
		InfluxURL:           "http://influxdb:8086",
		Database:            "soundtouch",
		Interface:           "en0",
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

	iff, filteredSpeakers, _ := processConfig(conf)

	var wg sync.WaitGroup
	log.Infof("Scanning for Soundtouch systems.")
	messageCh := make(chan *soundtouch.Update)

	for ok := true; ok; ok = (len(visibleSpeakers) < conf.NoSoundtouchSystems) {
		speakerCh := soundtouch.Lookup(iff)

		for speaker := range speakerCh {
			speakerInfo, _ := speaker.Info()
			speaker.DeviceInfo = speakerInfo
			spkLogger := log.WithFields(log.Fields{
				"Speaker": speaker.DeviceInfo.Name,
				"ID":      speaker.DeviceInfo.DeviceID,
			})

			if checkInMap(speaker.DeviceInfo.DeviceID, visibleSpeakers) {
				spkLogger.Debugf("Already included. Ignoring.")
				continue
			}

			ms := magicspeaker.New(speaker)

			// check wether we might have to ignore the speaker
			if len(filteredSpeakers) > 0 && !(filteredSpeakers)[speakerInfo.Name] {
				// spkLogger.Traceln("Seen but ignoring messages from: ", speakerInfo.Name)
				continue
			}

			visibleSpeakers[speaker.DeviceInfo.DeviceID] = ms
			spkLogger.Infof("Listening\n")
			spkLogger.Debugf(" with IP: %v", speaker.IP)
			wg.Add(1)

			// for each speaker we forward the messages to msgChan
			go func(s *soundtouch.Speaker, msgChan chan *soundtouch.Update) {
				defer wg.Done()

				webSocketCh, _ := s.Listen()
				for message := range webSocketCh {
					msgChan <- message
				}
			}(speaker, messageCh)
		}
		time.Sleep(10 * time.Second)
	}

	log.Infof("Found all Soundtouch systems. Normal Operation.")

	// We need only one loop for all messages.
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

// Will create the interface, and the speakerMap
func processConfig(conf config) (*net.Interface, speakerMap, error) {
	filteredSpeakers := make(speakerMap)

	i, err := net.InterfaceByName(conf.Interface)

	if err != nil {
		log.Fatalf("Error with interface. %s", err)
	}

	log.Debugf("Listening @ %v, supports: %v, HW Address: %v\n", i.Name, i.Flags.String(), i.HardwareAddr)

	for _, value := range conf.Speakers {
		filteredSpeakers[value] = true
		log.Debugf("Reacting only speakers %v\n", value)
	}

	return i, filteredSpeakers, nil
}

func checkInMap(deviceID string, list magicspeaker.MagicSpeakers) bool {
	for _, ms := range list {
		if ms.Speaker.DeviceInfo.DeviceID == deviceID {
			return true
		}
	}
	return false
}

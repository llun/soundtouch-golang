package main

import (
	"fmt"
	"net"
	"net/url"
	"reflect"
	"sync"

	"github.com/jpillora/opts"
	log "github.com/sirupsen/logrus"

	"github.com/theovassiliou/soundtouch-golang"
)

var soundtouchNetwork = make(map[string]string)
var strengthMapping = map[string]int{
	"EXCELLENT_SIGNAL": 100, "GOOD_SIGNAL": 70, "POOR_SIGNAL": 30, "MARGINAL_SIGNAL": 10,
}

var playStateMapping = map[soundtouch.Source]int{
	"PLAY_STATE": 1, "PAUSE_STATE": 2, "STOP_STATE": 3, "STANDBY": 5, "BUFFERING_STATE": 8, "INVALID_PLAY_STATUS": 13,
}

var conf = config{}

//set this via ldflags (see https://stackoverflow.com/q/11354518)
var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"
var influxDBinstance = soundtouch.InfluxDB{
	BaseHTTPURL: url.URL{
		Scheme: "http",
		Host:   "localhost:8086",
	},
	Database: "soundtouch",
}

type config struct {
	Speakers []string  `help:"Speakers to listen for, all if not set"`
	LogLevel log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
}

func main() {
	conf = config{
		LogLevel: log.DebugLevel,
	}

	//parse config
	opts.New(&conf).
		Repo("github.com/theovassiliou/soundtouch-golang").
		Version(VERSION).
		Parse()

	log.SetLevel(conf.LogLevel)

	i, _ := net.InterfaceByName("en0")
	log.Infof("Name : %v, supports: %v, HW Address: %v\n", i.Name, i.Flags.String(), i.HardwareAddr)
	speakerCh := soundtouch.Lookup(i)
	var wg sync.WaitGroup
	messageCh := make(chan *soundtouch.Update)

	for speaker := range speakerCh {
		di, _ := speaker.Info()
		speaker.DeviceInfo = di
		soundtouchNetwork[di.DeviceID] = di.Name
		log.Infof("Speaker: %v\n", speaker)
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
		if reflect.TypeOf(m.Value).Name() == "ConnectionStateUpdated" {
			// wifi,name=„Küche“,deviceID=„08DF1F117BB7“ wifiStrength=82,connected=true
			c, _ := m.Value.(soundtouch.ConnectionStateUpdated)
			lineproto := fmt.Sprintf("wifi,name=\"%s\",deviceID=\"%s\" wifiStrength=%v,connected=\"%v\"",
				soundtouchNetwork[m.DeviceId],
				m.DeviceId,
				strengthMapping[c.Signal],
				func() string {
					if c.Up == "true" {
						return "UP"
					}
					return "DOWN"
				}())
			log.WithFields(log.Fields{
				"database": "telegraf"}).Infof(lineproto)
			result, err := influxDBinstance.SetData("write", []byte(lineproto))
			if err != nil {
				log.WithFields(log.Fields{
					"influx": "send"}).Infof("failed")
			}
			log.WithFields(log.Fields{
				"influx": "send"}).Debugf("succeeded: %v", string(result))
		} else if reflect.TypeOf(m.Value).Name() == "NowPlaying" {
			np, _ := m.Value.(soundtouch.NowPlaying)
			lineproto := fmt.Sprintf("playing,name=\"%s\",deviceID=\"%s\" playStatus=%v,album=\"%v\"",
				soundtouchNetwork[m.DeviceId],
				m.DeviceId,
				func() int {
					ps := playStateMapping[np.PlayStatus]
					if ps == 0 && np.Source == "STANDBY" {
						return playStateMapping["STANDBY"]
					}
					return ps
				}(),
				func() string {
					if np.Album == "" {
						return "none"
					}
					return np.Album
				}())
			log.WithFields(log.Fields{
				"database": "telegraf"}).Infof(lineproto)
			result, err := influxDBinstance.SetData("write", []byte(lineproto))
			if err != nil {
				log.WithFields(log.Fields{
					"influx": "send"}).Infof("failed")
			}
			log.WithFields(log.Fields{
				"influx": "send"}).Debugf("succeeded: %v", string(result))
		} else {
			log.Infof("%v -> %v\n", soundtouchNetwork[m.DeviceId], reflect.TypeOf(m.Value).Name())
		}
	}
	wg.Wait()
}

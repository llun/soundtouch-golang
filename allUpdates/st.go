package main

import (
	"net"
<<<<<<< HEAD
	"reflect"
=======
>>>>>>> a4b4ced7bbfa995b9c559c4a6bf8ed454ba25bd7
	"sync"

	"github.com/jpillora/opts"
	log "github.com/sirupsen/logrus"

	"github.com/theovassiliou/soundtouch-golang"
)

<<<<<<< HEAD
// var soundtouchNetwork = make(map[string]string)

=======
const WEBSOCKET_PORT int = 8080
const MESSAGE_BUFFER_SIZE int = 256

var soundtouchNetwork = make(map[string]string)
>>>>>>> a4b4ced7bbfa995b9c559c4a6bf8ed454ba25bd7
var conf = config{}

//set this via ldflags (see https://stackoverflow.com/q/11354518)
var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

<<<<<<< HEAD
const shortUsage = "Captures broadcastet information from your Bose Soundtouch systems."

var influxDB = soundtouch.InfluxDB{
	SoundtouchNetwork: make(map[string]string),
}

type config struct {
	Speakers  []string  `opts:"group=Soundtouch" help:"Speakers to listen for, all if not set"`
	Interface string    `opts:"group=Soundtouch" help:"network interface to listen"`
	LogLevel  log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
=======
type config struct {
	Speakers []string  `help:"Speakers to listen for, all if not set"`
	LogLevel log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
>>>>>>> a4b4ced7bbfa995b9c559c4a6bf8ed454ba25bd7
}

func main() {
	conf = config{
<<<<<<< HEAD
		LogLevel:  log.DebugLevel,
		Interface: "en0",
=======
		LogLevel: log.DebugLevel,
>>>>>>> a4b4ced7bbfa995b9c559c4a6bf8ed454ba25bd7
	}

	//parse config
	opts.New(&conf).
<<<<<<< HEAD
		Summary(shortUsage).
		Repo("github.com/theovassiliou/soundtouch-golang").
=======
		Repo("github.com/theovassiliou/dta").
>>>>>>> a4b4ced7bbfa995b9c559c4a6bf8ed454ba25bd7
		Version(VERSION).
		Parse()

	log.SetLevel(conf.LogLevel)

<<<<<<< HEAD
	i, err := net.InterfaceByName(conf.Interface)

	if err != nil {
		log.Fatalf("Error with interface. %s", err)
	}

	log.Debugf("Listening @ %v, supports: %v, HW Address: %v\n", i.Name, i.Flags.String(), i.HardwareAddr)

=======
	i, _ := net.InterfaceByName("en0")
	log.Infof("Name : %v, supports: %v, HW Address: %v\n", i.Name, i.Flags.String(), i.HardwareAddr)
>>>>>>> a4b4ced7bbfa995b9c559c4a6bf8ed454ba25bd7
	speakerCh := soundtouch.Lookup(i)
	var wg sync.WaitGroup
	messageCh := make(chan *soundtouch.Update)

	for speaker := range speakerCh {
		di, _ := speaker.Info()
		speaker.DeviceInfo = di
<<<<<<< HEAD
		influxDB.SoundtouchNetwork[di.DeviceID] = di.Name
		spkLogger := log.WithFields(log.Fields{
			"Speaker": speaker.DeviceInfo.Name,
			"ID":      speaker.DeviceInfo.DeviceID,
		})
		spkLogger.Infof("Listening\n")
		spkLogger.Debugf(" with IP: %v", speaker.IP)
=======
		soundtouchNetwork[di.DeviceID] = di.Name
		log.Infof("Speaker: %v\n", speaker)
>>>>>>> a4b4ced7bbfa995b9c559c4a6bf8ed454ba25bd7
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
<<<<<<< HEAD
		mLogger := log.WithFields(log.Fields{
			"Speaker": influxDB.SoundtouchNetwork[m.DeviceId],
			"Value":   reflect.TypeOf(m.Value).Name(),
		})
		mLogger.Infof("%v\n", m)

=======
		log.Infof("%v -> %v\n", soundtouchNetwork[m.DeviceId], m)
>>>>>>> a4b4ced7bbfa995b9c559c4a6bf8ed454ba25bd7
	}
	wg.Wait()
}

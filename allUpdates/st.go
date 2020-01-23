package main

import (
	"net"
	"sync"

	"github.com/jpillora/opts"
	log "github.com/sirupsen/logrus"

	"github.com/theovassiliou/soundtouch-golang"
)

const WEBSOCKET_PORT int = 8080
const MESSAGE_BUFFER_SIZE int = 256

var soundtouchNetwork = make(map[string]string)
var conf = config{}

//set this via ldflags (see https://stackoverflow.com/q/11354518)
var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

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
		Repo("github.com/theovassiliou/dta").
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
		log.Infof("%v -> %v\n", soundtouchNetwork[m.DeviceId], m)
	}
	wg.Wait()
}

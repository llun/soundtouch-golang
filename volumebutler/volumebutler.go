package main

import (
	"net"
	"sync"

	"github.com/jpillora/opts"
	scribble "github.com/nanobox-io/golang-scribble"
	log "github.com/sirupsen/logrus"

	"github.com/theovassiliou/soundtouch-golang"
	"github.com/theovassiliou/soundtouch-golang/volumebutler/magicspeaker"
)

//set this via ldflags (see https://stackoverflow.com/q/11354518)
var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

const shortUsage = "Auto adjust volumes for special media on specific Soundtouch speakers."

type config struct {
	Speakers  []string  `opts:"group=Soundtouch" help:"Speakers to listen for, all if not set"`
	Interface string    `opts:"group=Soundtouch" help:"network interface to listen"`
	LogLevel  log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
	Filename  string    `opts:"group=VolumeButler" help:"Where to store the data"`
}
type speakerMap map[string]bool

var conf = config{}
var visibleSpeakers = make(magicspeaker.MagicSpeakers)

func main() {
	conf = config{
		LogLevel:  log.DebugLevel,
		Interface: "en0",
		Filename:  "./Volumes.db",
	}

	//parse config
	opts.New(&conf).
		Summary(shortUsage).
		Repo("github.com/theovassiliou/soundtouch-golang").
		Version(VERSION).
		Parse()

	log.SetLevel(conf.LogLevel)

	iff, filteredSpeakers, scribbleDb, _ := processConfig(conf)

	speakerCh := soundtouch.Lookup(iff)
	var wg sync.WaitGroup
	messageCh := make(chan *soundtouch.Update)

	for speaker := range speakerCh {
		speakerInfo, _ := speaker.Info()
		speaker.DeviceInfo = speakerInfo
		ms := magicspeaker.New(speaker)

		spkLogger := log.WithFields(log.Fields{
			"Speaker": speaker.DeviceInfo.Name,
			"ID":      speaker.DeviceInfo.DeviceID,
		})

		// check wether we might have to ignore the speaker
		if len(filteredSpeakers) > 0 && !(filteredSpeakers)[speakerInfo.Name] {
			// spkLogger.Traceln("Seen but ignoring messages from: ", speakerInfo.Name)
			continue
		}

		visibleSpeakers[speaker.DeviceInfo.DeviceID] = ms
		spkLogger.Infof("Listening\n")
		spkLogger.Debugf(" with IP: %v", speaker.IP)
		wg.Add(1)

		go func(s *soundtouch.Speaker, msgChan chan *soundtouch.Update) {
			defer wg.Done()
			webSocketCh, _ := s.Listen()
			magicSpeaker := magicspeaker.New(s)
			magicSpeaker.ScribbleDb = scribbleDb
			magicSpeaker.WebSocketCh = webSocketCh
			magicSpeaker.SpeakerName = visibleSpeakers[ms.DeviceInfo.DeviceID].DeviceInfo.Name
			magicSpeaker.MessageLoop()
		}(speaker, messageCh)

	}

	wg.Wait()
}

// Will create the interface, the speakerMap, and the scribble database
func processConfig(conf config) (*net.Interface, speakerMap, *scribble.Driver, error) {
	filteredSpeakers := make(speakerMap)
	i, err := net.InterfaceByName(conf.Interface)

	if err != nil {
		log.Fatalf("Error with interface. %s", err)
	}

	log.Debugf("Listening @ %v, supports: %v, HW Address: %v\n", i.Name, i.Flags.String(), i.HardwareAddr)

	for _, value := range conf.Speakers {
		filteredSpeakers[value] = true
		log.Tracef("Reacting only speakers %v\n", value)
	}

	db, _ := scribble.New(conf.Filename, nil)
	if err != nil {
		log.Fatalf("Error with database. %s", err)
	}

	log.Debugf("Created database: %v\n", db)
	return i, filteredSpeakers, db, nil
}

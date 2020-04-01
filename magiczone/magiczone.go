package main

import (
	"net"
	"sync"
	"time"

	"github.com/jpillora/opts"
	log "github.com/sirupsen/logrus"

	"github.com/theovassiliou/soundtouch-golang"
	"github.com/theovassiliou/soundtouch-golang/magiczone/magicspeaker"
)

//set this via ldflags (see https://stackoverflow.com/q/11354518)
var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

const shortUsage = "Auto adjust volumes for special media on specific Soundtouch speakers."

type config struct {
	Speakers            []string  `opts:"group=Soundtouch" help:"Speakers to listen for, all if not set"`
	Interface           string    `opts:"group=Soundtouch" help:"network interface to listen"`
	NoSoundtouchSystems int       `opts:"group=Soundtouch" help:"Number of Soundtouch systems to scan for."`
	LogLevel            log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
}
type speakerMap map[string]bool

var conf = config{}
var visibleSpeakers = make(magicspeaker.MagicSpeakers)

func main() {
	conf = config{
		NoSoundtouchSystems: -1,
		LogLevel:            log.DebugLevel,
		Interface:           "en0",
	}

	//parse config
	opts.New(&conf).
		Summary(shortUsage).
		Repo("github.com/theovassiliou/soundtouch-golang").
		Version(VERSION).
		Parse()

	log.SetLevel(conf.LogLevel)

	iff, _ := processConfig(conf)

	var wg sync.WaitGroup
	log.Infof("Scanning for Soundtouch systems.")
	for ok := true; ok; ok = (len(visibleSpeakers) < conf.NoSoundtouchSystems) {
		speakerCh := soundtouch.Lookup(iff)
		messageCh := make(chan *soundtouch.Update)

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

			visibleSpeakers[speaker.DeviceInfo.DeviceID] = ms
			spkLogger.Infof("Listening\n")
			spkLogger.Debugf(" with IP: %v", speaker.IP)
			soundtouch.DumpZones(spkLogger, *speaker)
			wg.Add(1)

			go func(s *soundtouch.Speaker, msgChan chan *soundtouch.Update) {
				defer wg.Done()
				webSocketCh, _ := s.Listen()
				magicSpeaker := magicspeaker.New(s)
				magicSpeaker.WebSocketCh = webSocketCh
				magicSpeaker.SpeakerName = visibleSpeakers[ms.DeviceInfo.DeviceID].DeviceInfo.Name
				magicSpeaker.KnownSpeakers = &visibleSpeakers
				magicSpeaker.MessageLoop()
			}(speaker, messageCh)
		}
		time.Sleep(10 * time.Second)
	}
	log.Infof("Found all Soundtouch systems. Normal Operation.")
	wg.Wait()
}

// Will create the interface, and the speakerMap
func processConfig(conf config) (*net.Interface, error) {
	i, err := net.InterfaceByName(conf.Interface)

	if err != nil {
		log.Fatalf("Error with interface. %s", err)
	}

	log.Debugf("Listening @ %v, supports: %v, HW Address: %v\n", i.Name, i.Flags.String(), i.HardwareAddr)

	return i, nil
}

func checkInMap(deviceID string, list magicspeaker.MagicSpeakers) bool {
	for _, ms := range list {
		if ms.Speaker.DeviceInfo.DeviceID == deviceID {
			return true
		}
	}
	return false
}

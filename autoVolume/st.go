package main

import (
	"net"
	"reflect"
	"sync"
	"time"

	"github.com/jpillora/opts"
	scribble "github.com/nanobox-io/golang-scribble"
	log "github.com/sirupsen/logrus"

	"github.com/theovassiliou/soundtouch-golang"
)

// var soundtouchNetwork = make(map[string]string)

//set this via ldflags (see https://stackoverflow.com/q/11354518)
var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

const shortUsage = "Auto adjust volumes for special media on specific Soundtouch speakers."

type config struct {
	Speakers  []string  `opts:"group=Soundtouch" help:"Speakers to listen for, all if not set"`
	Interface string    `opts:"group=Soundtouch" help:"network interface to listen"`
	LogLevel  log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
	Filename  string    `opts:"group=AutoAdjust" help:"Where to store the data"`
}
type speakerMap map[string]bool
type dbEntry struct {
	Name        string
	DeviceID    string
	Volume      int
	LastUpdated time.Time
}

var conf = config{}
var visibleSpeakers = make(map[string]*soundtouch.Speaker)

func main() {
	conf = config{
		LogLevel:  log.DebugLevel,
		Interface: "en0",
		Filename:  "./dreiFragezeichenVolumes.db",
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
		speakerDevice, _ := speaker.Info()

		spkLogger := log.WithFields(log.Fields{
			"Speaker": speaker.DeviceInfo.Name,
			"ID":      speaker.DeviceInfo.DeviceID,
		})

		// check wether we might have to ignore the speaker
		if len(filteredSpeakers) > 0 && !(filteredSpeakers)[speakerDevice.Name] {
			spkLogger.Traceln("Seen but ignoring messages from: ", speakerDevice.Name)
			continue
		}

		speaker.DeviceInfo = speakerDevice
		visibleSpeakers[speakerDevice.DeviceID] = speaker
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
	var observable = make(map[string]bool)
	for m := range messageCh {
		mLogger := log.WithFields(log.Fields{
			"Speaker": visibleSpeakers[m.DeviceId].DeviceInfo.Name,
			"Value":   reflect.TypeOf(m.Value).Name(),
		})

		// Check whether a Drei Fragezeichen has been started
		mLogger.Infof("Checking whether this is a ???\n")
		x, album := m.IsDreiFragezeichen()
		mLogger.Infof("Result of check: %v, %v\n", x, album)
		if x && !observable[m.DeviceId] {

			observable[m.DeviceId] = true
			go func(mLogger *log.Entry) {
				mLogger.Infof("IsDreiFragezeichen: %v\n", x)
				// Check whether we are in the observation window
				mLogger.Infof("Observation window: %v\n", soundtouch.IsObservationWindow())
				// Look in your database whether we have a volume

				currentAlbum := &dbEntry{}
				oldVolume := -1

				scribbleDb.Read(visibleSpeakers[m.DeviceId].DeviceInfo.Name, album, &currentAlbum)
				// If yes, set new volume
				if currentAlbum.Name != "" {
					mLogger.Infof("setting volume %v", currentAlbum.Volume)
					visibleSpeakers[m.DeviceId].SetVolume(currentAlbum.Volume)
					oldVolume = currentAlbum.Volume
				}
				// Wait for a minute

				time.Sleep(60 * time.Second)

				np, _ := visibleSpeakers[m.DeviceId].NowPlaying()
				if !(np.PlayStatus == soundtouch.PLAY_STATE) {
					return
				}
				currentAlbum.Name = album
				currentAlbum.DeviceID = m.DeviceId
				vol, _ := visibleSpeakers[m.DeviceId].Volume()
				currentAlbum.Volume = vol.TargetVolume

				mLogger.Infof("Storing track: %#v\n", currentAlbum)
				// Store current volume
				if oldVolume > -1 {
					currentAlbum.Volume = (currentAlbum.Volume + oldVolume) / 2
				}
				scribbleDb.Write(visibleSpeakers[m.DeviceId].DeviceInfo.Name, album, currentAlbum)
				observable[m.DeviceId] = false

			}(mLogger)
		} else {
			mLogger.Infof("Ignoring because of %v or %v:\n", x, observable[m.DeviceId])
		}
		mLogger.Infof("%v\n", m)
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

	return i, filteredSpeakers, db, nil
}

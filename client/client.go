package main

import (
	"time"

	"github.com/jpillora/opts"
	log "github.com/sirupsen/logrus"

	"github.com/theovassiliou/soundtouch-golang"
	"github.com/theovassiliou/soundtouch-golang/plugins/logger"
)

//set this via ldflags (see https://stackoverflow.com/q/11354518)
var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

const shortUsage = "A simple example client to interact with Soundtouch speakers."

var conf = config{}
var soundtouchNetwork = make(map[string]*soundtouch.Speaker)

type config struct {
	Interface           string    `opts:"group=Soundtouch" help:"network interface to listen"`
	NoSoundtouchSystems int       `opts:"group=Soundtouch" help:"Number of Soundtouch systems to scan for."`
	LogLevel            log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
}

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

	nConf := soundtouch.NetworkConfig{
		InterfaceName: conf.Interface,
		NoOfSystems:   conf.NoSoundtouchSystems,
		Plugins: []soundtouch.Plugin{
			&logger.Logger{},
		},
	}

	speakerCh := soundtouch.GetDevices(nConf)
	for speaker := range speakerCh {

		spkLogger := log.WithFields(log.Fields{
			"Speaker": speaker.Name(),
			"ID":      speaker.DeviceID(),
		})
		spkLogger.Infof("Found device\n")
		spkLogger.Debugf(" with IP: %v", speaker.IP)
		soundtouchNetwork[speaker.Name()] = speaker
	}

	data, err := soundtouchNetwork["Office"].Volume()

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("This is the raw data: %s\n", data.Raw)
	log.Printf("The volume is: %d", data.TargetVolume)

	soundtouchNetwork["Office"].SetVolume(20)
	log.Printf("Set volume to 20")
	time.Sleep(15 * time.Second)
	log.Printf("Returning to volume %d", data.TargetVolume)

	soundtouchNetwork["Office"].SetVolume(data.TargetVolume)
	time.Sleep(2 * time.Second)

}

func sliceContains(name string, list []string) bool {
	for _, s := range list {
		if name == s {
			return true
		}
	}
	return false
}

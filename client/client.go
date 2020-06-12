package main

import (
	"reflect"

	"github.com/jpillora/opts"
	log "github.com/sirupsen/logrus"

	"github.com/theovassiliou/soundtouch-golang"
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
		UpdateHandlers: []soundtouch.UpdateHandlerConfig{
			{
				Name:          "",
				UpdateHandler: soundtouch.UpdateHandlerFunc(basicHandler),
				Terminate:     false,
			},
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
	log.Printf("%v\n", data)
	log.Printf("%s\n", data.Raw)
	log.Printf("The volume is: %d", data.TargetVolume)

	soundtouchNetwork["Office"].SetVolume(20)
	log.Printf("Set volume to 20")

}

func basicHandler(hndlName string, update soundtouch.Update, speaker soundtouch.Speaker) {

	mLogger := log.WithFields(log.Fields{
		"Speaker":       speaker.Name(),
		"UpdateMsgType": reflect.TypeOf(update.Value).Name(),
	})
	mLogger.Infof("%v\n", update)

}

func isIn(name string, selected []string) bool {
	for _, s := range selected {
		if name == s {
			return true
		}
	}
	return false
}

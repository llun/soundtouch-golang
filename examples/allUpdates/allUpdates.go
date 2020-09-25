package main

import (
	"io/ioutil"
	"os"

	"github.com/influxdata/toml"
	"github.com/jpillora/opts"
	log "github.com/sirupsen/logrus"

	"github.com/theovassiliou/soundtouch-golang"
	"github.com/theovassiliou/soundtouch-golang/plugins/episodecollector"
	"github.com/theovassiliou/soundtouch-golang/plugins/influxconnector"
	"github.com/theovassiliou/soundtouch-golang/plugins/logger"
	"github.com/theovassiliou/soundtouch-golang/plugins/magiczone"
)

var conf = config{}

type config struct {
	Global

	LogLevel     log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
	SampleConfig bool      `opts:"group=Configuration" help:"If set creates a sample config file that can be used later"`
	Config       string    `opts:"group=Soundtouch" help:"configuration file to load"`
}
type Global struct {
	Interface             string `opts:"group=Soundtouch" help:"network interface to listen"`
	NoOfSoundtouchSystems int    `opts:"group=Soundtouch" help:"Number of Soundtouch systems to scan for."`
}
type tomlConfig struct {
	Title            string
	Global           Global
	Logger           *logger.Config           `toml:"logger"`
	EpisodeCollector *episodecollector.Config `toml:"episodeCollector"`
	MagicZone        *magiczone.Config        `toml:"magicZone"`
	InfluxDB         *influxconnector.Config
}

func main() {
	conf = config{
		Global: Global{
			Interface:             "en0",
			NoOfSoundtouchSystems: -1,
		},
		SampleConfig: false,
		LogLevel:     log.DebugLevel,
		Config:       "config.toml",
	}

	//parse config
	opts.New(&conf).
		Parse()

	log.SetLevel(conf.LogLevel)

	f, err := os.Open(conf.Config)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	var config tomlConfig
	if err := toml.Unmarshal(buf, &config); err != nil {
		panic(err)
	}

	pl := []soundtouch.Plugin{}

	if config.Logger != nil {
		pl = append(pl, logger.NewLogger(*config.Logger))
	}

	if config.EpisodeCollector != nil {
		pl = append(pl, episodecollector.NewCollector(*config.EpisodeCollector))
	}

	if config.MagicZone != nil {
		pl = append(pl, magiczone.NewCollector(*config.MagicZone))
	}

	if config.InfluxDB != nil {
		pl = append(pl, influxconnector.NewLogger(*config.InfluxDB))
	}

	nConf := soundtouch.NetworkConfig{
		InterfaceName: conf.Global.Interface,
		NoOfSystems:   conf.Global.NoOfSoundtouchSystems,
		Plugins:       pl,
	}

	// SearchDevices does not closes the channel
	speakerCh := soundtouch.SearchDevices(nConf)
	for speaker := range speakerCh {
		log.Infof("Found device %s-%s with IP %s\n", speaker.Name(), speaker.DeviceID(), speaker.IP)
	}

}

package main

import (
	"io/ioutil"
	"os"

	"github.com/influxdata/toml"
	"github.com/jpillora/opts"
	log "github.com/sirupsen/logrus"

	"github.com/theovassiliou/soundtouch-golang"
	"github.com/theovassiliou/soundtouch-golang/plugins/episodeCollector"
	"github.com/theovassiliou/soundtouch-golang/plugins/logger"
)

var conf = config{}

type config struct {
	Global

	LogLevel     log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
	SampleConfig bool      `opts:"group=Configuration" help:"If set creates a sampe config file that can be used later"`
	Config       string    `opts:"group=Soundtouch" help:"a configuration file"`
}
type Global struct {
	Interface             string `opts:"group=Soundtouch" help:"network interface to listen"`
	NoOfSoundtouchSystems int    `opts:"group=Soundtouch" help:"Number of Soundtouch systems to scan for."`
}
type tomlConfig struct {
	Title            string
	Global           Global
	Logger           logger.Config           `toml:"logger"`
	EpisodeCollector episodeCollector.Config `toml:"episodeCollector"`
}

func main() {
	conf = config{
		Global: Global{
			Interface:             "",
			NoOfSoundtouchSystems: 0,
		},
		LogLevel:     0,
		SampleConfig: false,
		Config:       "",
	}

	conf = config{
		Global: Global{
			Interface:             "en0",
			NoOfSoundtouchSystems: -1,
		},
		Config:   "config.toml",
		LogLevel: log.DebugLevel,
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

	nConf := soundtouch.NetworkConfig{
		InterfaceName: config.Global.Interface,
		NoOfSystems:   config.Global.NoOfSoundtouchSystems,
		Plugins: []soundtouch.Plugin{
			// logger.NewLogger(config.Logger),
			episodeCollector.NewCollector(config.EpisodeCollector)},
	}

	// SearchDevices does not closes the channel
	speakerCh := soundtouch.SearchDevices(nConf)
	for speaker := range speakerCh {
		log.Infof("Found device %s-%s with IP %s\n", speaker.Name(), speaker.DeviceID(), speaker.IP)
	}

}

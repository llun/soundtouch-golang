package influxconnector

import (
	"fmt"
	"net/url"
	"reflect"

	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/soundtouch-golang"
)

var name = "InfluxConnector"

const sampleConfig = `
  # [influxDB]
  ## speakers for which messages should be logged. If empty, all 
  # speakers = ["Office", "Kitchen"]

  ## log_messages describes the message types to be logged
  ## one or more of "ConnectionStateUpdated", "NowPlaying", "Volume"
  ## all if empty
  # log_messages = ["ConnectionStateUpdated", "NowPlaying", "Volume"] 
  
  ## URL of the InfluxDB
  # influxURL = "http://influxdb:8086"

  ## Database where to store the events
  # database = "soundtouch"
  #
  ## dry_run indicates that the plugin dumps lineporoto for the influxDB conncetion
  ## as curl statement. 
  # dry_run = true
  ## 
  ## terminate indicates whether no further plugin will be called after this plugin has been executed
  # terminate = true
`

const description = "Writes event data to InfluxDB "

// InfluxDB describes the plugin. It has a
// Config to store the configuration
// Plugin the plugin function
// suspended indicates that the plugin is temporarely suspended
type InfluxDB struct {
	Config
	Plugin    soundtouch.PluginFunc
	suspended bool
}

var influxDB = soundtouch.InfluxDB{
	BaseHTTPURL: url.URL{
		Scheme: "http",
		Host:   "localhost:8086",
	},
	Database: "soundtouch",
}

// NewLogger creates a new Logger plugin with the configuration
func NewLogger(config Config) (d *InfluxDB) {
	d = &InfluxDB{}
	d.Config = config

	mLogger := log.WithFields(log.Fields{
		"Plugin": name,
	})

	v, err := url.Parse(config.InfluxURL)
	if err != nil {
		mLogger.Infof("Not a valid URL: %v", config.InfluxURL)
	}

	influxDB.BaseHTTPURL = *v
	influxDB.Database = config.Database

	mLogger.Infof("Initialised\n")
	return d
}

// Config contains the configuration of the plugin
// Speakers list of SpeakerNames the handler is added. All if empty
// Terminate indicates whether this is the last handler to be called
// IgnoreMessages a list of message types to be ignored
type Config struct {
	InfluxURL   string   `toml:"influxURL"`
	Database    string   `toml:"database"`
	Speakers    []string `toml:"speakers"`
	Terminate   bool     `toml:"terminate"`
	LogMessages []string `toml:"log_messages"`
	DryRun      bool     `toml:dry_run"`
}

// Name returns the plugin name
func (d *InfluxDB) Name() string {
	return name
}

// Description returns a string explaining the purpose of this plugin
func (d *InfluxDB) Description() string { return description }

// SampleConfig returns text explaining how plugin should be configured
func (d *InfluxDB) SampleConfig() string { return sampleConfig }

// Terminate indicates that no further plugin will be executed on this speaker
func (d *InfluxDB) Terminate() bool { return d.Config.Terminate }

// Disable temporarely the execution of the plugin
func (d *InfluxDB) Disable() { d.suspended = true }

// Enable temporarely the execution of the plugin
func (d *InfluxDB) Enable() { d.suspended = false }

// Execute runs the plugin with the given parameter
func (d *InfluxDB) Execute(pluginName string, update soundtouch.Update, speaker soundtouch.Speaker) {
	if len(d.LogMessages) > 0 && isIn(reflect.TypeOf(update.Value).Name(), d.LogMessages) {
		return
	}
	if len(d.Speakers) > 0 && !isIn(speaker.Name(), d.Speakers) {
		return
	}

	mLogger := log.WithFields(log.Fields{
		"Plugin":        name,
		"Speaker":       speaker.Name(),
		"UpdateMsgType": reflect.TypeOf(update.Value).Name(),
	})
	v, _ := update.Lineproto(influxDB, &update)

	if !(d.Config.DryRun) && v != "" {
		result, err := influxDB.SetData("write", []byte(v))
		if err != nil {
			mLogger.Errorf("failed")
		}
		mLogger.Debugf("succeeded: %v", string(result))

	} else if v != "" {
		fmt.Printf("curl -i -XPOST \"%v\" --data-binary '%v'\n", influxDB.WriteURL("write"), v)
	}
}

func isIn(name string, selected []string) bool {
	for _, s := range selected {
		if name == s {
			return true
		}
	}
	return false
}

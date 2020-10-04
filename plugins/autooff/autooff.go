package autooff

import (
	"reflect"

	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/soundtouch-golang"
)

var name = "AutoOff"

const sampleConfig = `
  ## Enabling the AutoOff plugin
  # [autoOff]
  
  ## speakers that trigger an autooff 
  # 	[autooff.Wohnzimmer]
  #			thenOff = ["Kueche", "Schrank"]
  #		[autooff.Schlafzimmer]
  #			thenOff = ["Office"]

  ## terminate indicates whether no further plugin will be called after this plugin has been executed
  # terminate = true
`

const description = "Switches speakers off if one is switched on"

// Collector describes the plugin. It has a
// Config to store the configuration
// Plugin the plugin function
// suspended indicates that the plugin is temporarely suspended
type Collector struct {
	Config
	Plugin    soundtouch.PluginFunc
	suspended bool
}

// NewCollector creates a new Collector plugin with the configuration
func NewCollector(config Config) (d *Collector) {
	d = &Collector{}
	d.Config = config

	mLogger := log.WithFields(log.Fields{
		"Plugin": name,
	})

	mLogger.Infof("Initialised\n")

	return d
}

// Config contains the configuration of the plugin
// Groups list of Actions.
type Config map[string]struct {
	ThenOff []string `toml:"thenOff"`
}

// Name returns the plugin name
func (d *Collector) Name() string {
	return name
}

// Description returns a string explaining the purpose of this plugin
func (d *Collector) Description() string { return description }

// SampleConfig returns text explaining how plugin should be configured
func (d *Collector) SampleConfig() string { return sampleConfig }

// Terminate indicates that no further plugin will be executed on this speaker
func (d *Collector) Terminate() bool { return false }

// Disable temporarely the execution of the plugin
func (d *Collector) Disable() { d.suspended = true }

// Enable temporarely the execution of the plugin
func (d *Collector) Enable() { d.suspended = false }

// Execute runs the plugin with the given parameter
func (d *Collector) Execute(pluginName string, update soundtouch.Update, speaker soundtouch.Speaker) {
	if reflect.TypeOf(update.Value).Name() != "NowPlaying" {
		return
	}
	mLogger := log.WithFields(log.Fields{
		"Plugin":        name,
		"Speaker":       speaker.Name(),
		"UpdateMsgType": reflect.TypeOf(update.Value).Name(),
	})
	mLogger.Debugln("Executing", pluginName)

	for observedSpeaker, thenOff := range d.Config {
		if speaker.Name() == observedSpeaker {
			// If speaker is playing and is playing from TV
			if speaker.IsAlive() && update.ContentItem().Source == "PRODUCT" {
				for _, offSpeaker := range thenOff.ThenOff {
					s := soundtouch.GetSpeakerByName(offSpeaker)
					if s != nil {
						s.PowerOff()
					} else {
						mLogger.Errorf("Configured speaker %s not present in soundtouch network. Please check config file.\n", offSpeaker)
					}
				}
			}
		}
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

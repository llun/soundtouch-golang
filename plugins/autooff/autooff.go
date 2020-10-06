package autooff

import (
	"reflect"

	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/soundtouch-golang"
)

var name = "AutoOff"

const description = "Switches speakers off if one is switched on"

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

// Config contains the configuration of the plugin
// Groups list of Actions.
type Config map[string]struct {
	ThenOff []string `toml:"thenOff"`
}

// AutoOff describes the plugin. It has a
// Config to store the configuration
// Plugin the plugin function
// suspended indicates that the plugin is temporarely suspended
type AutoOff struct {
	Config
	Plugin    soundtouch.PluginFunc
	suspended bool
}

// NewObserver creates a new Collector plugin with the configuration
func NewObserver(config Config) (d *AutoOff) {
	d = &AutoOff{}
	d.Config = config

	mLogger := log.WithFields(log.Fields{
		"Plugin": name,
	})

	mLogger.Infof("Initialised\n")

	return d
}

// Name returns the plugin name
func (d *AutoOff) Name() string {
	return name
}

// Description returns a string explaining the purpose of this plugin
func (d *AutoOff) Description() string { return description }

// SampleConfig returns text explaining how plugin should be configured
func (d *AutoOff) SampleConfig() string { return sampleConfig }

// Terminate indicates that no further plugin will be executed on this speaker
func (d *AutoOff) Terminate() bool { return false }

// Disable temporarely the execution of the plugin
func (d *AutoOff) Disable() { d.suspended = true }

// Enable temporarely the execution of the plugin
func (d *AutoOff) Enable() { d.suspended = false }

// Execute runs the plugin with the given parameter
func (d *AutoOff) Execute(pluginName string, update soundtouch.Update, speaker soundtouch.Speaker) {
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

func sliceContains(name string, list []string) bool {
	for _, s := range list {
		if name == s {
			return true
		}
	}
	return false
}

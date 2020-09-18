package episodeCollector

import (
	"reflect"
	"time"

	scribble "github.com/nanobox-io/golang-scribble"
	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/soundtouch-golang"
)

var name = "EpisodeCollector"

const sampleConfig = `
  ## speakers for which episodes should be stores. If empty, all 
  # speakers = ["Office", "Kitchen"]

  ## For which artists to collect the episodes
  ## all if empty
  # artists = ["Drei Frageezeichen","John Sinclair"] 
  
  ## terminate indicates whether no further plugin will be called after this plugin has been executed
  # terminate = true
`

const description = "Collects episodes for specific artists"

// Collector describes the plugin. It has a
// Config to store the configuration
// Plugin the plugin function
// suspended indicates that the plugin is temporarely suspended
type Collector struct {
	Config
	Plugin     soundtouch.PluginFunc
	suspended  bool
	scribbleDb *scribble.Driver
}

// NewCollector creates a new Collector plugin with the configuration
func NewCollector(config Config) (d *Collector) {
	d = &Collector{}
	d.Config = config

	mLogger := log.WithFields(log.Fields{
		"Plugin": name,
	})

	mLogger.Infof("Initialised\nø@")
	mLogger.Infof("Scanning for: %v\n", d.Artists)

	return d
}

// Config contains the configuration of the plugin
// Speakers list of SpeakerNames the handler is added. All if empty
// Terminate indicates whether this is the last handler to be called
// Artists a list of artists for which episodes should be collected
type Config struct {
	Speakers  []string `toml:"-"`
	Terminate bool     `toml:"terminate"`
	Artists   []string `toml:"artists"`
}

// Name returns the plugin name
func (d *Collector) Name() string {
	return name
}

type DbEntry struct {
	ContentItem soundtouch.ContentItem
	AlbumName   string
	Volume      int
	DeviceID    string
	LastUpdated time.Time
}

// Description returns a string explaining the purpose of this plugin
func (d *Collector) Description() string { return description }

// SampleConfig returns text explaining how plugin should be configured
func (d *Collector) SampleConfig() string { return sampleConfig }

// Terminate indicates that no further plugin will be executed on this speaker
func (d *Collector) Terminate() bool { return d.Config.Terminate }

// Disable temporarely the execution of the plugin
func (d *Collector) Disable() { d.suspended = true }

// Enable temporarely the execution of the plugin
func (d *Collector) Enable() { d.suspended = false }

func (d *Collector) Execute(pluginName string, update soundtouch.Update, speaker soundtouch.Speaker) {

	if len(d.Speakers) > 0 && !isIn(speaker.Name(), d.Speakers) {
		return
	}

	mLogger := log.WithFields(log.Fields{
		"Plugin":        name,
		"Speaker":       speaker.Name(),
		"UpdateMsgType": reflect.TypeOf(update.Value).Name(),
	})

	if !(update.Is("NowPlaying") || update.Is("Volume")) {
		if !update.Is("ConnectionStateUpdated") {
			typeName := reflect.TypeOf(update.Value).Name()
			mLogger.Debugf("Ignoring %s\n", typeName)
		}
		return
	}

	artist := update.Artist()
	album := update.Album()

	if !isIn(artist, d.Config.Artists) || !update.HasContentItem() {
		mLogger.Debugf("Ignoring album: %s\n", album)
		return
	}

	mLogger.Infof("Found an album: %v\n", album)

}

func isIn(name string, selected []string) bool {
	for _, s := range selected {
		if name == s {
			return true
		}
	}
	return false
}

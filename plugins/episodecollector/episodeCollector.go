package episodecollector

import (
	"reflect"
	"time"

	scribble "github.com/nanobox-io/golang-scribble"
	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/soundtouch-golang"
)

var name = "EpisodeCollector"

const sampleConfig = `
  ## Enabling the episodeCollector plugin
  # [episodeCollector]
  
  ## speakers for which episodes should be stores. If empty, all 
  # speakers = ["Office", "Kitchen"]

  ## For which artists to collect the episodes
  ## all if empty
  # artists = ["Drei Frageezeichen","John Sinclair"] 
  
  ## terminate indicates whether no further plugin will be called after this plugin has been executed
  # terminate = true

  ## database contains the directory name for the episodes database
  # database = "episode.db"
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

	mLogger.Infof("Initialised\n")
	mLogger.Tracef("Scanning for: %v\n", d.Artists)

	db, err := scribble.New(d.Database, nil)
	if err != nil {
		log.Fatalf("Error with database. %s", err)
	}

	mLogger.Debugf("Initialised database: %v\n", db)
	d.scribbleDb = db

	return d
}

// Config contains the configuration of the plugin
// Speakers list of SpeakerNames the handler is added. All if empty
// Terminate indicates whether this is the last handler to be called
// Artists a list of artists for which episodes should be collected
type Config struct {
	Speakers  []string `toml:"speakers"`
	Terminate bool     `toml:"terminate"`
	Artists   []string `toml:"artists"`
	Database  string   `toml:"database"`
}

// Name returns the plugin name
func (d *Collector) Name() string {
	return name
}

type dbEntry struct {
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

// Execute runs the plugin with the given parameter
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

	mLogger.Infof("Found album: %v\n", album)
	d.readAlbumDB(album, update)
}

func (d *Collector) readDB(album string, currentAlbum *dbEntry) *dbEntry {
	if currentAlbum == nil {
		currentAlbum = &dbEntry{}
	}
	d.scribbleDb.Read("All", album, &currentAlbum)
	return currentAlbum
}

func (d *Collector) writeDB(album string, storedAlbum *dbEntry) {
	storedAlbum.LastUpdated = time.Now()
	d.scribbleDb.Write("All", album, &storedAlbum)
}

func (d *Collector) readAlbumDB(album string, updateMsg soundtouch.Update) *dbEntry {

	storedAlbum := d.readDB(album, &dbEntry{})

	if storedAlbum.AlbumName == "" {
		// no, write this into the database
		storedAlbum.AlbumName = album
		// HYPO: We are in observation window, then the current volume could also
		// be a good measurement
		storedAlbum.DeviceID = updateMsg.DeviceID
		storedAlbum.ContentItem = updateMsg.ContentItem()
		d.writeDB(album, storedAlbum)
	}
	return storedAlbum
}

func isIn(name string, selected []string) bool {
	for _, s := range selected {
		if name == s {
			return true
		}
	}
	return false
}

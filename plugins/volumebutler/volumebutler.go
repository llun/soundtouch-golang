package volumebutler

import (
	"reflect"
	"time"

	scribble "github.com/nanobox-io/golang-scribble"
	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/soundtouch-golang"
)

var name = "volumeButler"

const sampleConfig = `
  ## Enabling the volumeButler plugin
  # [volumeButler]
  
  ## speakers for which volumeButler will handle volumes. None if empty. 
  # speakers = ["Office", "Kitchen"]

  ## For which artists volumes should be handled
  ## all if empty
  # artists = ["Drei Frageezeichen","John Sinclair"] 
  
  ## terminate indicates whether no further plugin will be called after this plugin has been executed
  # terminate = true

  ## database contains the directory name for the episodes database
  # database = "episode.db"
`

const description = "Automatically adjust sets volume based on listening history."

// Collector describes the plugin. It has a
// Config to store the configuration
// Plugin the plugin function
// suspended indicates that the plugin is temporarely suspended
// scribbleDB a link to the volumes database
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

	typeName := reflect.TypeOf(update.Value).Name()
	mLogger := log.WithFields(log.Fields{
		"Plugin":        name,
		"Speaker":       speaker.Name(),
		"UpdateMsgType": reflect.TypeOf(update.Value).Name(),
	})
	mLogger.Debugln("Executing", pluginName)
	mLogger.Debugf("  with %#v", d)

	if len(d.Speakers) == 0 || !isIn(speaker.Name(), d.Speakers) {
		mLogger.Debugln("Speaker not handled. --> Done!")
		return
	}

	if !(update.Is("NowPlaying") || update.Is("Volume")) {
		mLogger.Debugf("Ignoring %s. --> Done!\n", typeName)
		return
	}

	artist := update.Artist()
	album := update.Album()

	if !isIn(artist, d.Config.Artists) || !update.HasContentItem() {
		mLogger.Debugf("Ignoring album %s from %s\n", album, artist)
		return
	}
	mLogger.Infof("Found album %s from %s\n", album, artist)
	// time window independend
	// Do we know this album already?  - read from database
	storedAlbum := d.ReadAlbumDB(album, update)

	// time window and speaker depended
	// 		if available for this album
	//			set the volume
	if storedAlbum.Volume != 0 && time.Now().After(storedAlbum.LastUpdated.Add(20*time.Minute)) {
		// Only setting a volume if it was last update more than 20 minutes ago
		mLogger.Infof("Stored volume was set more than 20minutes ago\n")
		mLogger.Infof("Setting volume to %d\n", storedAlbum.Volume)
		speaker.SetVolume(storedAlbum.Volume)
	}

	// wait for a minute and process last volume observed
	// construct the mean value of current and past volumes
	// store the update value
	mLogger.Infof("Going to sleep for 60s\n")
	time.Sleep(60 * time.Second)

	// clear message and keep last volume update
	mLogger.Infof("Scanning for Volume\n")
	lastVolume := ScanForVolume(&speaker)
	d.ReadDB(speaker.Name(), album, storedAlbum)
	if lastVolume != nil {
		storedAlbum.Volume = storedAlbum.calcNewVolume(lastVolume.TargetVolume)
		mLogger.Infof("writing volume to %v\n", storedAlbum.Volume)
		d.scribbleDb.Write(speaker.Name(), album, &storedAlbum)
	}
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

func getSpeaker(updateMsg soundtouch.Update) *soundtouch.Speaker {
	for _, aKnownDevice := range soundtouch.GetKnownDevices() {
		if aKnownDevice.DeviceID() == updateMsg.DeviceID {
			return aKnownDevice
		}
	}
	return nil
}

func (d *Collector) ReadAlbumDB(album string, updateMsg soundtouch.Update) *dbEntry {

	speaker := getSpeaker(updateMsg)
	if speaker == nil {
		return nil
	}

	mLogger := log.WithFields(log.Fields{
		"Plugin":        name,
		"Speaker":       speaker.Name(),
		"UpdateMsgType": reflect.TypeOf(updateMsg.Value).Name(),
	})

	storedAlbum := d.ReadDB(speaker.Name(), album, &dbEntry{})

	if storedAlbum.AlbumName == "" {
		mLogger.Infof("Album %s not yet known. Reading volume.", storedAlbum.AlbumName)
		// no, write this into the database
		retrievedVol, _ := speaker.Volume()
		mLogger.Infof("Volume is %d", retrievedVol.ActualVolume)
		storedAlbum.AlbumName = album
		// HYPO: We are in observation window, then the current volume could also
		// be a good measurement
		storedAlbum.Volume = retrievedVol.TargetVolume
		storedAlbum.DeviceID = updateMsg.DeviceID
		storedAlbum.LastUpdated = time.Now()
		storedAlbum.ContentItem = updateMsg.ContentItem()
		d.WriteDB(speaker.Name(), album, storedAlbum)
	}
	return storedAlbum
}

func (d *Collector) WriteDB(speakerName, album string, storedAlbum *dbEntry) {
	storedAlbum.LastUpdated = time.Now()
	d.scribbleDb.Write(speakerName, album, &storedAlbum)
	d.scribbleDb.Write("All", album, &storedAlbum)
}

func (d *Collector) ReadDB(speakerName string, album string, currentAlbum *dbEntry) *dbEntry {
	if currentAlbum == nil {
		currentAlbum = &dbEntry{}
	}
	d.scribbleDb.Read(speakerName, album, &currentAlbum)
	return currentAlbum
}

func ScanForVolume(m *soundtouch.Speaker) *soundtouch.Volume {
	var lastVolume *soundtouch.Volume
	var mLogger *log.Entry
	for scanMsg := range m.WebSocketCh {
		typeName := reflect.TypeOf(scanMsg.Value).Name()
		mLogger = log.WithFields(log.Fields{
			"Plugin":        name,
			"Speaker":       m.Name(),
			"UpdateMsgType": typeName,
		})
		if scanMsg.Is("Volume") {
			aVol, _ := scanMsg.Value.(soundtouch.Volume)
			lastVolume = &aVol
			mLogger.Printf("Ignoring! Volume: %#v", lastVolume)
		} else {
			mLogger.Debugf("Ignoring!! %s\n", typeName)
		}
		if len(m.WebSocketCh) == 0 {
			break
		}
	}
	mLogger.Infof("lastVolume was %d\n", lastVolume.ActualVolume)
	return lastVolume
}

func (db *dbEntry) calcNewVolume(currVolume int) int {
	oldVol := db.Volume
	if oldVol == 0 {
		oldVol = currVolume
	}
	return (oldVol + currVolume) / 2
}

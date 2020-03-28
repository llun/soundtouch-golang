package magicspeaker

import (
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"

	scribble "github.com/nanobox-io/golang-scribble"
	soundtouch "github.com/theovassiliou/soundtouch-golang"
)

type MagicSpeaker struct {
	*soundtouch.Speaker
	ScribbleDb  *scribble.Driver
	Vol         int
	SpeakerName string
	WebSocketCh chan *soundtouch.Update
}
type MagicSpeakers map[string]*MagicSpeaker

type DbEntry struct {
	ContentItem soundtouch.ContentItem
	AlbumName   string
	Volume      int
	DeviceID    string
	LastUpdated time.Time
}

func New(s *soundtouch.Speaker) *MagicSpeaker {
	return &MagicSpeaker{s, nil, 0, "", nil}
}

func (m *MagicSpeaker) ReadDB(album string, currentAlbum *DbEntry) *DbEntry {
	if currentAlbum == nil {
		currentAlbum = &DbEntry{}
	}
	m.ScribbleDb.Read(m.SpeakerName, album, &currentAlbum)
	return currentAlbum
}

func (m *MagicSpeaker) WriteDB(album string, storedAlbum *DbEntry) {
	storedAlbum.LastUpdated = time.Now()
	m.ScribbleDb.Write(m.SpeakerName, album, &storedAlbum)
}

func (m *MagicSpeaker) ReadAlbumDB(album string, updateMsg MagicUpdate) *DbEntry {

	storedAlbum := m.ReadDB(album, &DbEntry{})

	if storedAlbum.AlbumName == "" {
		// no, write this into the database
		retrievedVol, _ := m.Volume()
		storedAlbum.AlbumName = album
		// HYPO: We are in observation window, then the current volume could also
		// be a good measurement
		storedAlbum.Volume = retrievedVol.TargetVolume
		storedAlbum.DeviceID = updateMsg.DeviceID
		storedAlbum.LastUpdated = time.Now()
		storedAlbum.ContentItem = updateMsg.ContentItem()
		m.WriteDB(album, storedAlbum)
	}
	return storedAlbum
}

func (m *MagicSpeaker) MessageLoop() {
	for message := range m.WebSocketCh {
		mu := MagicUpdate{Update: *message}
		m.HandleUpdate(mu, m.WebSocketCh)
	}
}

func (m *MagicSpeaker) ScanForVolume() *soundtouch.Volume {
	var lastVolume *soundtouch.Volume
	for scanMsg := range m.WebSocketCh {
		typeName := reflect.TypeOf(scanMsg.Value).Name()
		mLogger2 := log.WithFields(log.Fields{
			"Speaker":     m.SpeakerName,
			"MessageType": typeName,
		})

		if scanMsg.Is("Volume") {
			aVol, _ := scanMsg.Value.(soundtouch.Volume)
			lastVolume = &aVol
			mLogger2.Printf("Ignoring! Volume: %#v", lastVolume)
		} else {
			mLogger2.Debugf("Ignoring!! %s\n", typeName)
		}
		if len(m.WebSocketCh) == 0 {
			break
		}
	}
	return lastVolume
}

// handle message per speaker
func (m *MagicSpeaker) HandleUpdate(msg MagicUpdate, webSocketCh chan *soundtouch.Update) {
	typeName := reflect.TypeOf(msg.Value).Name()
	mLogger := log.WithFields(log.Fields{
		"Speaker":     m.SpeakerName,
		"MessageType": typeName,
	})

	if !(msg.Is("NowPlaying") || msg.Is("Volume")) {
		if !msg.Is("ConnectionStateUpdated") {
			mLogger.Debugf("Ignoring %s\n", typeName)
		}
		return
	}

	// Check whether it matches the artist criteria
	isDreiFragezeichen, album := msg.ArtistMatches("Drei Fragezeichen")

	if !isDreiFragezeichen || !msg.HasContentItem() {
		mLogger.Debugf("Ignoring album: %s\n", album)
		return
	}

	// time window independend
	// Do we know this album already?  - read from database
	storedAlbum := m.ReadAlbumDB(album, msg)

	// time window and speaker depended
	// 		if available for this album
	//			set the volume
	if storedAlbum.Volume != 0 && time.Now().After(storedAlbum.LastUpdated.Add(20*time.Minute)) {
		// Only setting a volume if it was last update more than 20 minutes ago
		m.SetVolume(storedAlbum.Volume)
	}

	// wait for a minute and process last volume observed
	// construct the mean value of current and past volumes
	// store the update value
	time.Sleep(60 * time.Second)

	// clear message and keep last volume update
	lastVolume := m.ScanForVolume()
	m.ReadDB(album, storedAlbum)
	if lastVolume != nil {
		storedAlbum.Volume = storedAlbum.calcNewVolume(lastVolume.TargetVolume)
		m.ScribbleDb.Write(m.SpeakerName, album, &storedAlbum)
	}
}

func (db *DbEntry) calcNewVolume(currVolume int) int {
	oldVol := db.Volume
	if oldVol == 0 {
		oldVol = currVolume
	}
	return (oldVol + currVolume) / 2
}

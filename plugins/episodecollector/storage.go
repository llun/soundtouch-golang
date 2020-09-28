package episodecollector

import (
	"time"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/theovassiliou/soundtouch-golang"
)

type dbEntry struct {
	ContentItem soundtouch.ContentItem
	AlbumName   string
	Volume      int
	DeviceID    string
	LastUpdated time.Time
}

func readDB(sbd *scribble.Driver, album string, currentAlbum *dbEntry) *dbEntry {
	if currentAlbum == nil {
		currentAlbum = &dbEntry{}
	}
	sbd.Read("All", album, &currentAlbum)
	return currentAlbum
}

func writeDB(sbd *scribble.Driver, album string, storedAlbum *dbEntry) {
	storedAlbum.LastUpdated = time.Now()
	sbd.Write("All", album, &storedAlbum)
}

func readAlbumDB(sbd *scribble.Driver, album string, updateMsg soundtouch.Update) *dbEntry {

	storedAlbum := readDB(sbd, album, &dbEntry{})

	if storedAlbum.AlbumName == "" {
		// no, write this into the database
		storedAlbum.AlbumName = album
		// HYPO: We are in observation window, then the current volume could also
		// be a good measurement
		storedAlbum.DeviceID = updateMsg.DeviceID
		storedAlbum.ContentItem = updateMsg.ContentItem()
		writeDB(sbd, album, storedAlbum)
	}
	return storedAlbum
}

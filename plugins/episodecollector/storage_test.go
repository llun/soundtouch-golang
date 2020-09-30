package episodecollector

import (
	"reflect"
	"testing"
	"time"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/theovassiliou/soundtouch-golang"
)

func Test_readDB(t *testing.T) {
	db, _ := scribble.New("episodes.tests.db", nil)
	albumA := &dbEntry{
		ContentItem: soundtouch.ContentItem{
			Type:         "",
			Source:       "STORED_MUSIC",
			Location:     "Loc",
			Name:         "AlbumA",
			IsPresetable: true,
		},
		AlbumName:   "AlbumA",
		Volume:      0,
		DeviceID:    "Dev",
		LastUpdated: time.Time{},
	}
	type args struct {
		sbd          *scribble.Driver
		album        string
		currentAlbum *dbEntry
	}
	tests := []struct {
		name string
		args args
		want *dbEntry
	}{
		{"nil dbEntry",
			args{
				db,
				"",
				nil,
			},
			&dbEntry{},
		},
		{"empty Album",
			args{
				db,
				"",
				&dbEntry{},
			},
			&dbEntry{},
		},
		{"one Album",
			args{
				db,
				"AlbumA",
				&dbEntry{},
			},
			albumA,
		},
		{"wrong Album",
			args{
				db,
				"AlbumB",
				&dbEntry{},
			},
			&dbEntry{},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readDB(tt.args.sbd, tt.args.album, tt.args.currentAlbum); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

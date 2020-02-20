package magicspeaker

import (
	"reflect"
	"strings"

	soundtouch "github.com/theovassiliou/soundtouch-golang"
)

type MagicUpdate struct {
	soundtouch.Update
}

func (u MagicUpdate) ArtistMatches(artist string) (bool, string) {
	switch reflect.TypeOf(u.Value).Name() {
	case "NowPlaying":
		np := u.Value.(soundtouch.NowPlaying)
		if strings.HasPrefix(np.Artist, artist) {
			return true, np.Album
		}
	}
	return false, ""
}

func (u MagicUpdate) ContentItem() soundtouch.ContentItem {
	if u.HasContentItem() {
		return u.Value.(soundtouch.NowPlaying).Content
	}
	return soundtouch.ContentItem{}

}

func (u MagicUpdate) HasContentItem() bool {
	switch reflect.TypeOf(u.Value).Name() {
	case "NowPlaying":
		return true
	}
	return false
}

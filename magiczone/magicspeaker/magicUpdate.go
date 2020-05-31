package magicspeaker

import (
	"reflect"

	soundtouch "github.com/theovassiliou/soundtouch-golang"
)

type MagicUpdate struct {
	soundtouch.Update
}

func (u MagicUpdate) PlayState() soundtouch.PlayStatus {
	switch reflect.TypeOf(u.Value).Name() {
	case "NowPlaying":
		np := u.Value.(soundtouch.NowPlaying)
		return np.PlayStatus
	}
	return ""
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

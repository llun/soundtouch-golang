package soundtouch

import (
	"reflect"
	"strings"
	"time"
)

func (u Update) IsDreiFragezeichen() (bool, string) {
	switch reflect.TypeOf(u.Value).Name() {
	case "NowPlaying":
		np := u.Value.(NowPlaying)
		if strings.HasPrefix(np.Artist, "Drei Fragezeichen") {
			return true, np.Album
		}
	}
	return false, ""
}

func IsObservationWindow() bool {
	now := time.Now()

	return (now.Hour() > 22 || now.Hour() < 8)
}

package soundtouch

import (
	"fmt"
)

// Key typing
type Key string

// All soundtouch key constants
const (
	PLAY      = "PLAY"
	PAUSE     = "PAUSE"
	PLAYPAUSE = "PLAY_PAUSE"
	// STOP       = "STOP" // Deprecated
	PREVTRACK  = "PREV_TRACK"
	NEXTTRACK  = "NEXT_TRACK"
	THUMBSUP   = "THUMBS_UP"
	THUMBSDOWN = "THUMBS_DOWN"
	BOOKMARK   = "BOOKMARK"
	POWER      = "POWER"
	MUTE       = "MUTE"
	// VOLUMEUP       = "VOLUME_UP" // Deprecated
	// VOLUMEDOWN     = "VOLUME_DOWN" // Deprecated
	PRESET1 = "PRESET_1"
	PRESET2 = "PRESET_2"
	PRESET3 = "PRESET_3"
	PRESET4 = "PRESET_4"
	PRESET5 = "PRESET_5"
	PRESET6 = "PRESET_6"
	// AUXINPUT       = "AUX_INPUT" // Deprecated
	SHUFFLEOFF     = "SHUFFLE_OFF"
	SHUFFLEON      = "SHUFFLE_ON"
	REPEATOFF      = "REPEAT_OFF"
	REPEATONE      = "REPEAT_ONE"
	REPEATALL      = "REPEAT_ALL"
	ADDFAVORITE    = "ADD_FAVORITE"
	REMOVEFAVORITE = "REMOVE_FAVORITE"
)

// ALLKEYS contains all KEY constant that can be sent to soundtouch
var ALLKEYS = []string{
	PLAY,
	PAUSE,
	PLAYPAUSE,
	// STOP,
	PREVTRACK,
	NEXTTRACK,
	THUMBSUP,
	THUMBSDOWN,
	BOOKMARK,
	POWER,
	MUTE,
	// VOLUMEUP,
	// VOLUMEDOWN,
	PRESET1,
	PRESET2,
	PRESET3,
	PRESET4,
	PRESET5,
	PRESET6,
	// AUXINPUT,
	SHUFFLEOFF,
	SHUFFLEON,
	REPEATOFF,
	REPEATONE,
	REPEATALL,
	ADDFAVORITE,
	REMOVEFAVORITE,
}

// PressKey sends key press command to soundtouch system. For POWER also release is send immediatly afterwards.
func (s *Speaker) PressKey(key Key) error {
	press := []byte(fmt.Sprintf(`<key state="press" sender="Gabbo">%v</key>`, key))
	_, err := s.SetData("key", press)
	if err != nil {
		return err
	}

	if key == POWER {
		release := []byte(fmt.Sprintf(`<key state="release" sender="Gabbo">%v</key>`, key))
		_, err = s.SetData("key", release)
		if err != nil {
			return err
		}
	}
	return nil
}

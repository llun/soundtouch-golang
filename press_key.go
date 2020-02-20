package soundtouch

import (
	"fmt"
)

// Key typing
type Key string

// All soundtouch key constants
const (
	PLAY           = "PLAY"
	PAUSE          = "PAUSE"
	STOP           = "STOP"
	PREVTRACK      = "PREV_TRACK"
	NEXTTRACK      = "NEXT_TRACK"
	THUMBSUP       = "THUMBS_UP"
	THUMBSDOWN     = "THUMBS_DOWN"
	BOOKMARK       = "BOOKMARK"
	POWER          = "POWER"
	MUTE           = "MUTE"
	VOLUMEUP       = "VOLUME_UP"
	VOLUMEDOWN     = "VOLUME_DOWN"
	PRESET1        = "PRESET_1"
	PRESET2        = "PRESET_2"
	PRESET3        = "PRESET_3"
	PRESET4        = "PRESET_4"
	PRESET5        = "PRESET_5"
	PRESET6        = "PRESET_6"
	AUXINPUT       = "AUX_INPUT"
	SHUFFLEOFF     = "SHUFFLE_OFF"
	SHUFFLEON      = "SHUFFLE_ON"
	REPEATOFF      = "REPEAT_OFF"
	REPEATONE      = "REPEAT_ONE"
	REPEATALL      = "REPEAT_ALL"
	ADDFAVORITE    = "ADD_FAVORITE"
	REMOVEFAVORITE = "REMOVE_FAVORITE"
)

// PressKey send key press/release pair to soundtouch system
func (s *Speaker) PressKey(key Key) error {
	press := []byte(fmt.Sprintf(`<key state="press" sender="Gabbo">%v</key>`, key))
	_, err := s.SetData("key", press)
	if err != nil {
		return err
	}

	release := []byte(fmt.Sprintf(`<key state="release" sender="Gabbo">%v</key>`, key))
	_, err = s.SetData("key", release)
	if err != nil {
		return err
	}
	return nil
}

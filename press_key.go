package soundtouch

import (
	"fmt"
)

type Key string

const (
	PLAY            Key = "PLAY"
	PAUSE               = "PAUSE"
	PLAY_PAUSE          = "PLAY_PAUSE"
	STOP                = "STOP"
	PREV_TRACK          = "PREV_TRACK"
	NEXT_TRACK          = "NEXT_TRACK"
	THUMBS_UP           = "THUMBS_UP"
	THUMBS_DOWN         = "THUMBS_DOWN"
	BOOKMARK            = "BOOKMARK"
	POWER               = "POWER"
	MUTE                = "MUTE"
	VOLUME_UP           = "VOLUME_UP"
	VOLUME_DOWN         = "VOLUME_DOWN"
	PRESET_1            = "PRESET_1"
	PRESET_2            = "PRESET_2"
	PRESET_3            = "PRESET_3"
	PRESET_4            = "PRESET_4"
	PRESET_5            = "PRESET_5"
	PRESET_6            = "PRESET_6"
	AUX_INPUT           = "AUX_INPUT"
	SHUFFLE_OFF         = "SHUFFLE_OFF"
	SHUFFLE_ON          = "SHUFFLE_ON"
	REPEAT_OFF          = "REPEAT_OFF"
	REPEAT_ONE          = "REPEAT_ONE"
	REPEAT_ALL          = "REPEAT_ALL"
	ADD_FAVORITE        = "ADD_FAVORITE"
	REMOVE_FAVORITE     = "REMOVE_FAVORITE"
)

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

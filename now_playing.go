package soundtouch

import (
	"encoding/xml"
)

type PlayStatus string

const (
	PLAY_STATE          = "PLAY_STATE"
	PAUSE_STATE         = "PAUSE_STATE"
	BUFFERING_STATE     = "BUFFERING_STATE"
	INVALID_PLAY_STATUS = "INVALID_PLAY_STATUS"
	STOP_STATE          = "STOP"
	STANDBY             = "STANDBY"
)

type Source string

const (
	SLAVE                = "SLAVE_SOURCE"
	INTERNET_RADIO       = "INTERNET_RADIO"
	LOCAL_INTERNET_RADIO = "LOCAL_INTERNET_RADIO"
	PANDORA              = "PANDORA"
	TUNEIN               = "TUNEIN"
	AIRPLAY              = "AIRPLAY"
	STORED_MUSIC         = "STORED_MUSIC"
	AUX                  = "AUX"
	BLUETOOTH            = "BLUETOOTH"
	PRODUCT              = "PRODUCT"
	OFF_SOURCE           = "OFF_SOURCE"
	CURRATED_RADIO       = "CURRATED_RADIO"
	UPDATE               = "UPDATE"
	DEEZER               = "DEEZER"
	SPOTIFY              = "SPOTIFY"
	IHEART               = "IHEART"
)

type NowPlaying struct {
	PlayStatus    PlayStatus  `xml:"playStatus"`
	Source        string      `xml:"source,attr"`
	SourceAccount string      `xml:"sourceAccount,attr"`
	DeviceId      string      `xml:"deviceID,attr"`
	Content       ContentItem `xml:"ContentItem"`
	Track         string      `xml:"track"`
	Artist        string      `xml:"artist"`
	Album         string      `xml:"album"`
	TrackID       string      `xml:"trackID"`
	Art           string      `xml:"art"`
	StreamType    string      `xml:"streamType"`
	Raw           []byte
}

func (s *Speaker) NowPlaying() (NowPlaying, error) {
	body, err := s.GetData("now_playing")
	if err != nil {
		return NowPlaying{}, err
	}

	nowPlaying := NowPlaying{
		Raw: body,
	}
	err = xml.Unmarshal(body, &nowPlaying)
	if err != nil {
		return nowPlaying, err
	}
	return nowPlaying, nil
}

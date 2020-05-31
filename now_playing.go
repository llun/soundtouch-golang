package soundtouch

import (
	"encoding/xml"
)

// PlayStatus is a string type
type PlayStatus string

// All Playing states of soundtouch speaker
const (
	PlayState         = "PLAY_STATE"
	PauseState        = "PAUSE_STATE"
	BufferingState    = "BUFFERING_STATE"
	InvalidPlayStatus = "INVALID_PLAY_STATUS"
	StopState         = "STOP"
	Standby           = "STANDBY"
)

// Source is a string type
type Source string

// All Sources of a soundtouch speaker
const (
	Slave              = "SLAVE_SOURCE"
	InternetRadio      = "INTERNET_RADIO"
	LocalInternetRadio = "LOCAL_INTERNET_RADIO"
	Pandora            = "PANDORA"
	TuneIn             = "TUNEIN"
	Airplay            = "AIRPLAY"
	StoredMusic        = "STORED_MUSIC"
	Aux                = "AUX"
	Bluetooth          = "BLUETOOTH"
	Product            = "PRODUCT"
	OffSource          = "OFF_SOURCE"
	CurratedRadio      = "CURRATED_RADIO"
	UPDATE             = "UPDATE"
	Deezer             = "DEEZER"
	Spotify            = "SPOTIFY"
	IHeart             = "IHEART"
)

// All StreamTypes
const (
	RadioStreaming = "RADIO_STREAMING"
	TrackOnDemand  = "TRACK_ONDEMAND"
)

// NowPlaying defines the now_playing message to/from soundtouch system
type NowPlaying struct {
	PlayStatus    PlayStatus  `xml:"playStatus"`
	Source        string      `xml:"source,attr"`
	SourceAccount string      `xml:"sourceAccount,attr"`
	DeviceID      string      `xml:"deviceID,attr"`
	Content       ContentItem `xml:"ContentItem"`
	Track         string      `xml:"track"`
	Artist        string      `xml:"artist"`
	Album         string      `xml:"album"`
	TrackID       string      `xml:"trackID"`
	Art           string      `xml:"art"`
	StreamType    string      `xml:"streamType"`
	Raw           []byte
}

// NowPlaying sends the now_playing command to the soundtouch system
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

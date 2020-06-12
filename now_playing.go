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

// ALLSTATES contains all soundtouch state constants
var ALLSTATES = []string{
	PlayState,
	PauseState,
	BufferingState,
	InvalidPlayStatus,
	StopState,
	Standby,
}

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

// ALLSOURCES contains all soundtouch sources
var ALLSOURCES = []string{
	Slave,
	InternetRadio,
	LocalInternetRadio,
	Pandora,
	TuneIn,
	Airplay,
	StoredMusic,
	Aux,
	Bluetooth,
	Product,
	OffSource,
	CurratedRadio,
	UPDATE,
	Deezer,
	Spotify,
	IHeart,
}

// All StreamTypes
const (
	RadioStreaming = "RADIO_STREAMING"
	TrackOnDemand  = "TRACK_ONDEMAND"
)

// ALLSTREAMS contains all soundtouch streamtypes
var ALLSTREAMS = []string{
	RadioStreaming, TrackOnDemand,
}

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
	Raw           []byte      `json:"-"`
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

// IsAlive returns true in case the soundtouch system is in PlayState
func (s *Speaker) IsAlive() bool {
	np, err := s.NowPlaying()
	if err != nil {
		return false
	}

	var isAlive bool
	isAlive = np.PlayStatus != Standby

	if isAlive {
		isAlive = np.PlayStatus == PlayState
	}
	return isAlive
}

// IsPoweredOn returns true in case the soundtouch system is not in standby but powered
func (s *Speaker) IsPoweredOn() bool {
	np, err := s.NowPlaying()
	if err != nil {
		return false
	}
	return np.Source != Standby
}

// PowerOn switches the soundtouch system on. True on success. Returns false in case the system was already powered
func (s *Speaker) PowerOn() bool {
	if !s.IsPoweredOn() {
		s.PressKey(POWER)
		return true
	}
	return false
}

// PowerOnWithVolume switches the soundtouch system on and set's a specific volume. True on success. Returns false in case the system was already powered with no volume set.
func (s *Speaker) PowerOnWithVolume(vol int) bool {
	if !s.IsPoweredOn() {
		s.PressKey(POWER)
		s.SetVolume(vol)
		return true
	}
	return false

}

// PowerOff powers off the soundtouch systems. True on success. Returns false in case it was already powered off.
func (s *Speaker) PowerOff() bool {
	if s.IsPoweredOn() {
		s.PressKey(POWER)
		return true
	}
	return false

}

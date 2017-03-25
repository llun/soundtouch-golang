package soundtouch

import (
  "encoding/xml"
)

type Source string

const (
  SLAVE          Source = "SLAVE_SOURCE"
  INTERNET_RADIO        = "INTERNET_RADIO"
  PANDORA               = "PANDORA"
  AIRPLAY               = "AIRPLAY"
  STORED_MUSIC          = "STORED_MUSIC"
  AUX                   = "AUX"
  OFF_SOURCE            = "OFF_SOURCE"
  CURRATED_RADIO        = "CURRATED_RADIO"
  STANDBY               = "STANDBY"
  UPDATE                = "UPDATE"
  DEEZER                = "DEEZER"
  SPOTIFY               = "SPOTIFY"
  IHEART                = "IHEART"
)

type NowPlaying struct {
  Source        string      `xml:"source,attr"`
  SourceAccount string      `xml:"sourceAccount,attr"`
  DeviceId      string      `xml:"deviceID,attr"`
  Content       ContentItem `xml:"ContentItem"`
  Track         string      `xml:"track"`
  Artist        string      `xml:"artist"`
  Album         string      `xml:"album"`
  TrackID       string      `xml:"trackID"`
  Art           string      `xml:"art"`
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

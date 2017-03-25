package soundtouch

import (
  "encoding/xml"
)

type Preset struct {
  ID      int         `xml:"id,attr"`
  Content ContentItem `xml:"ContentItem"`
}

type Presets struct {
  DeviceID string   `xml:"deviceID,attr"`
  Presets  []Preset `xml:"preset"`
  Raw      []byte
}

func (s *Speaker) Presets() (Presets, error) {
  body, err := s.GetData("presets")
  if err != nil {
    return Presets{}, err
  }

  presets := Presets{
    Raw: body,
  }
  err = xml.Unmarshal(body, &presets)
  if err != nil {
    return presets, err
  }
  return presets, nil
}

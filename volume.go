package soundtouch

import (
  "encoding/xml"
  "fmt"
)

type Volume struct {
  DeviceID     string `xml:"deviceID,attr"`
  TargetVolume int    `xml:"targetvolume"`
  ActualVolume int    `xml:"actualvolume"`
  Muted        bool   `xml:"mutedenabled"`
  Raw          []byte
}

func (s *Speaker) Volume() (Volume, error) {
  body, err := s.GetData("volume")
  if err != nil {
    return Volume{}, err
  }

  volume := Volume{
    Raw: body,
  }
  err = xml.Unmarshal(body, &volume)
  if err != nil {
    return volume, err
  }
  return volume, nil
}

func (s *Speaker) SetVolume(volume int) error {
  data := []byte(fmt.Sprintf("<volume>%v</volume>", volume))
  _, err := s.SetData("volume", data)
  if err != nil {
    return err
  }
  return nil
}

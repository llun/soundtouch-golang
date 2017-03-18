package soundtouch

import (
  "encoding/xml"
)

type Info struct {
  DeviceID string `xml:"deviceID,attr"`
  Name     string `xml:"name"`
  Type     string `xml:"type"`
  Raw      []byte
}

func (s *Speaker) Info() (Info, error) {
  body, err := s.GetData("info")
  if err != nil {
    return Info{}, err
  }

  info := Info{
    Raw: body,
  }
  err = xml.Unmarshal(body, &info)
  if err != nil {
    return info, err
  }
  return info, nil
}

package soundtouch

import (
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
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

func (s Info) String() string {
	if log.GetLevel() >= log.TraceLevel {
		return fmt.Sprintf("%v (%v): %v\n%v", s.Name, s.DeviceID, s.Type, string(s.Raw))
	}
	return fmt.Sprintf("%v (%v): %v", s.Name, s.DeviceID, s.Type)
}

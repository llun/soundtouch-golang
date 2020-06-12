package soundtouch

import (
	"encoding/xml"
	"fmt"

	log "github.com/sirupsen/logrus"
)

//Info defines the Info command for the soundtouch system
type Info struct {
	DeviceID  string   `xml:"deviceID,attr"`
	Name      string   `xml:"name"`
	Type      string   `xml:"type"`
	IPAddress []string `xml:"networkInfo>ipAddress"`
	Raw       []byte
}

type IPAddress string

// Info retrieves speaker information
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

// String creates, depending of the loglevel different string representations on info message
func (s Info) String() string {
	if log.GetLevel() >= log.TraceLevel {
		return fmt.Sprintf("%v (%v): %v\n%v", s.Name, s.DeviceID, s.Type, string(s.Raw))
	}
	return fmt.Sprintf("%v (%v): %v", s.Name, s.DeviceID, s.Type)
}

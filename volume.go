package soundtouch

import (
	"encoding/xml"
	"fmt"
)

// Volume defines the Volume command
type Volume struct {
	DeviceID     string `xml:"deviceID,attr"`
	TargetVolume int    `xml:"targetvolume"`
	ActualVolume int    `xml:"actualvolume"`
	Muted        bool   `xml:"mutedenabled"`
	Raw          []byte
}

// Volume sends the volume command to the soundtouch system to retrieve the volume
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

// SetVolume sends the volume command to the soundtouch system to set the volume
func (s *Speaker) SetVolume(volume int) error {
	data := []byte(fmt.Sprintf("<volume>%v</volume>", volume))
	_, err := s.SetData("volume", data)
	if err != nil {
		return err
	}
	return nil
}

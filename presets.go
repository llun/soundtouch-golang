package soundtouch

import (
	"encoding/xml"
)

// Preset specifies a preset field in soundtouch messages
type Preset struct {
	ID      int         `xml:"id,attr"`
	Content ContentItem `xml:"ContentItem"`
}

// Presets specifies the Presets Update message
type Presets struct {
	DeviceID string   `xml:"deviceID,attr"`
	Presets  []Preset `xml:"preset"`
	Raw      []byte   `json:"-"`
}

// Presets queries the presets of a soundtouch system
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

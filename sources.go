package soundtouch

import (
	"encoding/xml"
)

// SourceItem defines a source within a soundtouch system
type SourceItem struct {
	Source        Source `xml:"source,attr"`
	SourceAccount string `xml:"sourceAccount,attr"`
	Status        string `xml:"status,attr"`
	Local         bool   `xml:"isLocal,attr"`
	Value         string `xml:",innerxml"`
}

// Sources defines the soundtouch sources command
type Sources struct {
	DeviceID    string       `xml:"deviceID,attr"`
	SourceItems []SourceItem `xml:"sourceItem"`
	Raw         []byte
}

// Sources sends the sources command to the soundtouch system
func (s *Speaker) Sources() (Sources, error) {
	body, err := s.GetData("sources")
	if err != nil {
		return Sources{}, err
	}

	sources := Sources{
		Raw: body,
	}
	err = xml.Unmarshal(body, &sources)
	if err != nil {
		return sources, err
	}
	return sources, nil
}

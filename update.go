package soundtouch

import (
	"bytes"
	"encoding/xml"
	"errors"
)

type Update struct {
	Value interface{}
}

func NewUpdate(body []byte) (*Update, error) {
	decoder := xml.NewDecoder(bytes.NewBuffer(body))
	root, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	if root == nil {
		return nil, errors.New("Invalid XML format")
	}

	rootElement := root.(xml.StartElement)
	if rootElement.Name.Local != "updates" {
		return nil, errors.New("Unsupported event")
	}

	updateType, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	value, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	updateTypeElement := updateType.(xml.StartElement)
	switch updateTypeElement.Name.Local {
	case "volumeUpdated":
		valueElement := value.(xml.StartElement)

		var volume Volume
		err = decoder.DecodeElement(&volume, &valueElement)
		if err != nil {
			return nil, err
		}

		return &Update{volume}, nil
	case "nowPlayingUpdated":
		valueElement := value.(xml.StartElement)

		var nowPlaying NowPlaying
		err = decoder.DecodeElement(&nowPlaying, &valueElement)
		if err != nil {
			return nil, err
		}

		return &Update{nowPlaying}, nil
	case "nowSelectionUpdated":
		valueElement := value.(xml.StartElement)

		var preset Preset
		err = decoder.DecodeElement(&preset, &valueElement)
		if err != nil {
			return nil, err
		}
		return &Update{preset}, nil
	default:
		return nil, nil
	}
}

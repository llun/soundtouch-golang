package soundtouch

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
)

// Update carries all update messages communicated to/from soundtouch system
type Update struct {
	DeviceID string
	Value    interface{}
}

// NewUpdate returns an Update decoded from a received message
func NewUpdate(body []byte) (*Update, error) {
	decoder := xml.NewDecoder(bytes.NewBuffer(body))
	root, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	if root == nil {
		return nil, errors.New("Invalid XML format")
	}

	rootElement, ok := root.(xml.StartElement)
	if !ok {
		return nil, errors.New("Invalid XML format")
	}

	if rootElement.Name.Local != "updates" {
		return nil, errors.New("Unsupported event")
	}
	var deviceID string
	for i := 0; i < len(rootElement.Attr); i++ {
		if rootElement.Attr[i].Name.Local == "deviceID" {
			deviceID = rootElement.Attr[i].Value
		}
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
	case "connectionStateUpdated":
		valueElement := updateTypeElement

		var connState ConnectionStateUpdated
		err = decoder.DecodeElement(&connState, &valueElement)
		if err != nil {
			return nil, err
		}

		return &Update{deviceID, connState}, nil
	case "volumeUpdated":
		valueElement := value.(xml.StartElement)

		var volume Volume
		err = decoder.DecodeElement(&volume, &valueElement)
		if err != nil {
			return nil, err
		}

		return &Update{deviceID, volume}, nil
	case "nowPlayingUpdated":
		valueElement := value.(xml.StartElement)

		var nowPlaying NowPlaying
		err = decoder.DecodeElement(&nowPlaying, &valueElement)
		if err != nil {
			return nil, err
		}

		return &Update{deviceID, nowPlaying}, nil
	case "nowSelectionUpdated":
		valueElement := value.(xml.StartElement)

		var preset Preset
		err = decoder.DecodeElement(&preset, &valueElement)
		if err != nil {
			return nil, err
		}
		return &Update{deviceID, preset}, nil
	default:
		return nil, fmt.Errorf("Unhandeld Update Message. %v", string(body))
	}
}

// String returns the specific part of an update message
func (u Update) String() string {
	return fmt.Sprintf("%v", u.Value)
}

// Is returns true in case Update message is of type msgTypeName
// is one of
// - ConnectionStateUpdated
// - Volume
// - NowPlaying
// - Preset
func (u Update) Is(msgTypeName string) bool {
	return reflect.TypeOf(u.Value).Name() == msgTypeName
}

// Artist returns the artist if present, empty else
func (u Update) Artist() string {
	switch reflect.TypeOf(u.Value).Name() {
	case "NowPlaying":
		return u.Value.(NowPlaying).Artist
	}
	return ""
}

// Album returns the Album if present, else empty
func (u Update) Album() string {
	switch reflect.TypeOf(u.Value).Name() {
	case "NowPlaying":
		return u.Value.(NowPlaying).Album
	}
	return ""
}

// HasContentItem returns true if the Update message has contentItem, false else
func (u Update) HasContentItem() bool {
	switch reflect.TypeOf(u.Value).Name() {
	case "NowPlaying":
		return true
	}
	return false
}

// ContentItem returns the ContentItem if present, or an empty one if not present
func (u Update) ContentItem() ContentItem {
	if u.HasContentItem() {
		return u.Value.(NowPlaying).Content
	}
	return ContentItem{}

}

// GetSpeaker returns the Speaker instance the Update has been send from
func GetSpeaker(updateMsg Update) *Speaker {
	for _, aKnownDevice := range GetKnownDevices() {
		if aKnownDevice.DeviceID() == updateMsg.DeviceID {
			return aKnownDevice
		}
	}
	return nil
}

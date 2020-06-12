package soundtouch

import (
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
)

func Test_contains(t *testing.T) {
	s1 := Speaker{
		IP:           nil,
		Port:         0,
		BaseHTTPURL:  url.URL{},
		WebSocketURL: url.URL{},
		DeviceInfo: Info{
			DeviceID: "aabbccspeaker1",
			Name:     "Speaker1",
			Type:     "",
			Raw:      nil,
		},
		conn:           &websocket.Conn{},
		webSocketCh:    make(chan *Update),
		UpdateHandlers: nil,
	}
	s2 := Speaker{
		IP:           nil,
		Port:         0,
		BaseHTTPURL:  url.URL{},
		WebSocketURL: url.URL{},
		DeviceInfo: Info{
			DeviceID: "ccbbaaspeaker2",
			Name:     "Speaker2",
			Type:     "",
			Raw:      nil,
		},
		conn:           &websocket.Conn{},
		webSocketCh:    make(chan *Update),
		UpdateHandlers: nil,
	}
	ss := make(speakers)
	ss["aabbccspeaker1"] = &s1
	ss["ccbbaaspeaker2"] = &s2

	type args struct {
		deviceID string
		list     speakers
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Included",
			args{"aabbccspeaker1", ss},
			true,
		},
		{
			"NotIncluded",
			args{"aabbcc", ss},
			false,
		},
		{
			"Empty DeviceId",
			args{"", ss},
			false,
		},
		{
			"No Map",
			args{"aabbccspeaker1", nil},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.list, tt.args.deviceID); got != tt.want {
				t.Errorf("checkInMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		conn:        &websocket.Conn{},
		webSocketCh: make(chan *Update),
		Plugins:     nil,
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
		conn:        &websocket.Conn{},
		webSocketCh: make(chan *Update),
		Plugins:     nil,
	}
	ss := make(Speakers)
	ss["aabbccspeaker1"] = &s1
	ss["ccbbaaspeaker2"] = &s2

	type args struct {
		deviceID string
		list     Speakers
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

func Test_isIn(t *testing.T) {
	type args struct {
		list     []string
		deviceID string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Contains",
			args{
				[]string{"A", "B", "C"},
				"A",
			},
			true,
		},
		{
			"Doesn't Contains",
			args{
				[]string{"A", "B", "C"},
				"X",
			},
			false,
		},
		{
			"Search in Empty",
			args{
				[]string{},
				"A",
			},
			false,
		},
		{
			"Search in Empty",
			args{
				nil,
				"A",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isIn(tt.args.list, tt.args.deviceID); got != tt.want {
				t.Errorf("isIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

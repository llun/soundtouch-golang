package soundtouch

import (
	"net"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
)

func TestSpeaker_Name(t *testing.T) {
	type fields Speaker
	tests := []struct {
		name     string
		fields   fields
		wantName string
	}{
		{
			"Valid name",
			fields{
				IP:           nil,
				Port:         0,
				BaseHTTPURL:  url.URL{},
				WebSocketURL: url.URL{},
				DeviceInfo: Info{
					DeviceID: "xxee",
					Name:     "Speakers Name",
					Type:     "",
					Raw:      nil,
				},
				conn: &websocket.Conn{},
			},
			"Speakers Name",
		},
		{
			"Empty name",
			fields{
				IP:           nil,
				Port:         0,
				BaseHTTPURL:  url.URL{},
				WebSocketURL: url.URL{},
				DeviceInfo:   Info{},
				conn:         &websocket.Conn{},
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Speaker{
				IP:           tt.fields.IP,
				Port:         tt.fields.Port,
				BaseHTTPURL:  tt.fields.BaseHTTPURL,
				WebSocketURL: tt.fields.WebSocketURL,
				DeviceInfo:   tt.fields.DeviceInfo,
				conn:         tt.fields.conn,
			}
			if gotName := s.Name(); gotName != tt.wantName {
				t.Errorf("Speaker.Name() = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}

func TestSpeaker_DeviceID(t *testing.T) {
	type fields struct {
		IP           net.IP
		Port         int
		BaseHTTPURL  url.URL
		WebSocketURL url.URL
		DeviceInfo   Info
		conn         *websocket.Conn
	}
	tests := []struct {
		name     string
		fields   fields
		wantName string
	}{
		{
			"Valid DeviceId",
			fields{
				IP:           nil,
				Port:         0,
				BaseHTTPURL:  url.URL{},
				WebSocketURL: url.URL{},
				DeviceInfo: Info{
					DeviceID: "xxee",
					Name:     "",
					Type:     "",
					Raw:      nil,
				},
				conn: &websocket.Conn{},
			},
			"xxee",
		},
		{
			"Empty name",
			fields{
				IP:           nil,
				Port:         0,
				BaseHTTPURL:  url.URL{},
				WebSocketURL: url.URL{},
				DeviceInfo:   Info{},
				conn:         &websocket.Conn{},
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Speaker{
				IP:           tt.fields.IP,
				Port:         tt.fields.Port,
				BaseHTTPURL:  tt.fields.BaseHTTPURL,
				WebSocketURL: tt.fields.WebSocketURL,
				DeviceInfo:   tt.fields.DeviceInfo,
				conn:         tt.fields.conn,
			}
			if gotName := s.DeviceID(); gotName != tt.wantName {
				t.Errorf("Speaker.DeviceID() = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}

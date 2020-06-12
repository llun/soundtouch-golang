package soundtouch

import (
	"net"
	"net/url"
	"reflect"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/mdns"
	log "github.com/sirupsen/logrus"
)

func TestSpeaker_RemoveUpdateHandler(t *testing.T) {
	type fields struct {
		IP             net.IP
		Port           int
		BaseHTTPURL    url.URL
		WebSocketURL   url.URL
		DeviceInfo     Info
		conn           *websocket.Conn
		webSocketCh    chan *Update
		UpdateHandlers []UpdateHandlerConfig
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"Remove NotConfigured",
			fields{
				UpdateHandlers: []UpdateHandlerConfig{
					UpdateHandlerConfig{
						Name: "NotConfigured",
						UpdateHandler: UpdateHandlerFunc(func(hndlName string, update Update, speaker Speaker) {
							log.Infof("UpdateHandler not configured.")
						}),
						Terminate: false,
					},
				},
			},
			args{name: "NotConfigured"},
		},
		{
			"No Handler",
			fields{
				UpdateHandlers: []UpdateHandlerConfig{},
			},
			args{name: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Speaker{
				IP:             tt.fields.IP,
				Port:           tt.fields.Port,
				BaseHTTPURL:    tt.fields.BaseHTTPURL,
				WebSocketURL:   tt.fields.WebSocketURL,
				DeviceInfo:     tt.fields.DeviceInfo,
				conn:           tt.fields.conn,
				webSocketCh:    tt.fields.webSocketCh,
				UpdateHandlers: tt.fields.UpdateHandlers,
			}
			s.RemoveUpdateHandler(tt.args.name)
			if s.HasUpdateHandler(tt.args.name) {
				t.Errorf("Speaker.RemoveUpdateHandler failed. %v still included", tt.args.name)
			}
		})
	}
}

func TestSpeaker_AddUpdateHandler(t *testing.T) {
	type fields struct {
		IP             net.IP
		Port           int
		BaseHTTPURL    url.URL
		WebSocketURL   url.URL
		DeviceInfo     Info
		conn           *websocket.Conn
		webSocketCh    chan *Update
		UpdateHandlers []UpdateHandlerConfig
	}
	type args struct {
		uhc UpdateHandlerConfig
	}

	f1 := fields{
		UpdateHandlers: []UpdateHandlerConfig{
			{
				Name: "UpdateHandlerConfig1",
				UpdateHandler: UpdateHandlerFunc(func(hndlName string, update Update, speaker Speaker) {
					log.Infof("UpdateHandler not configured.")
				}),
				Terminate: false,
			},
		},
	}

	f2 := fields{
		UpdateHandlers: []UpdateHandlerConfig{
			{
				Name: "UpdateHandlerConfig1",
				UpdateHandler: UpdateHandlerFunc(func(hndlName string, update Update, speaker Speaker) {
					log.Infof("UpdateHandler not configured.")
				}),
				Terminate: false,
			},
			{
				Name: "UpdateHandlerConfig2",
				UpdateHandler: UpdateHandlerFunc(func(hndlName string, update Update, speaker Speaker) {
					log.Infof("UpdateHandler not configured.")
				}),
				Terminate: false,
			},
		},
	}

	fDefault := fields{
		UpdateHandlers: []UpdateHandlerConfig{
			{
				Name: "NotConfigured",
				UpdateHandler: UpdateHandlerFunc(func(hndlName string, update Update, speaker Speaker) {
					log.Infof("UpdateHandler not configured.")
				}),
				Terminate: false,
			},
		},
	}
	uhc1 := UpdateHandlerConfig{
		Name: "NewUpdateHandler",
		UpdateHandler: UpdateHandlerFunc(func(hndlName string, update Update, speaker Speaker) {
			log.Infof("UpdateHandler not configured.")
		}),
		Terminate: false,
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		lenAfter int
	}{
		{
			"Add 1 to 1",
			f1,
			args{uhc1},
			2,
		},
		{
			"Add 1 to 2",
			f2,
			args{uhc1},
			3,
		},
		{
			"Add 1 to empty",
			fields{
				UpdateHandlers: []UpdateHandlerConfig{},
			},
			args{uhc1},
			1,
		},
		{
			"Add 1 to default",
			fDefault,
			args{uhc1},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Speaker{
				IP:             tt.fields.IP,
				Port:           tt.fields.Port,
				BaseHTTPURL:    tt.fields.BaseHTTPURL,
				WebSocketURL:   tt.fields.WebSocketURL,
				DeviceInfo:     tt.fields.DeviceInfo,
				conn:           tt.fields.conn,
				webSocketCh:    tt.fields.webSocketCh,
				UpdateHandlers: tt.fields.UpdateHandlers,
			}
			s.AddUpdateHandler(tt.args.uhc)
			if len(s.UpdateHandlers) != tt.lenAfter {
				t.Errorf("Speaker.AddUpdateHandler %v failed. length want %d is %v ", tt.name, tt.lenAfter, len(s.UpdateHandlers))
			}
			if !s.HasUpdateHandler(tt.args.uhc.Name) {
				t.Errorf("Speaker.AddUpdateHandler %v failed. %v not included.", tt.name, tt.args.uhc.Name)
			}
		})
	}
}

func TestNewSpeaker(t *testing.T) {
	entriesCh := make(chan *mdns.ServiceEntry, 7)

	params := mdns.DefaultParams("_soundtouch._tcp")
	params.Entries = entriesCh
	params.Interface, _ = net.InterfaceByName("en0")

	mdns.Query(params)
	var entry *mdns.ServiceEntry
	entry = <-entriesCh

	type args struct {
		entry *mdns.ServiceEntry
	}
	tests := []struct {
		name string
		args args
		want *Speaker
	}{
		{
			"New empty ServiceEntry",
			args{entry},
			&Speaker{
				IP:   nil,
				Port: 0,
				BaseHTTPURL: url.URL{
					Scheme:     "",
					Opaque:     "",
					User:       &url.Userinfo{},
					Host:       "",
					Path:       "",
					RawPath:    "",
					ForceQuery: false,
					RawQuery:   "",
					Fragment:   "",
				},
				WebSocketURL: url.URL{
					Scheme:     "",
					Opaque:     "",
					User:       &url.Userinfo{},
					Host:       "",
					Path:       "",
					RawPath:    "",
					ForceQuery: false,
					RawQuery:   "",
					Fragment:   "",
				},
				DeviceInfo: Info{
					DeviceID:  "",
					Name:      "",
					Type:      "",
					IPAddress: nil,
					Raw:       nil,
				},
				conn:           &websocket.Conn{},
				webSocketCh:    make(chan *Update),
				UpdateHandlers: nil,
			},
		},
		{
			"New empty ServiceEntry",
			args{nil},
			&Speaker{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSpeaker(tt.args.entry); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSpeaker() = %v, want %v", got, tt.want)
			}
		})
	}
}

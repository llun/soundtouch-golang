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
		UpdateHandlers []PluginConfig
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
				UpdateHandlers: []PluginConfig{
					{
						Name: "NotConfigured",
						Plugin: PluginFunc(func(pluginName string, update Update, speaker Speaker) {
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
				UpdateHandlers: []PluginConfig{},
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
			s.RemovePlugin(tt.args.name)
			if s.HasPlugin(tt.args.name) {
				t.Errorf("Speaker.RemoveUpdateHandler failed. %v still included", tt.args.name)
			}
		})
	}
}

func TestSpeaker_AddUpdateHandler(t *testing.T) {
	type fields struct {
		IP           net.IP
		Port         int
		BaseHTTPURL  url.URL
		WebSocketURL url.URL
		DeviceInfo   Info
		conn         *websocket.Conn
		webSocketCh  chan *Update
		Plugins      []PluginConfig
	}
	type args struct {
		uhc PluginConfig
	}

	f1 := fields{
		Plugins: []PluginConfig{
			{
				Name: "UpdateHandlerConfig1",
				Plugin: PluginFunc(func(pluginName string, update Update, speaker Speaker) {
					log.Infof("UpdateHandler not configured.")
				}),
				Terminate: false,
			},
		},
	}

	f2 := fields{
		Plugins: []PluginConfig{
			{
				Name: "UpdateHandlerConfig1",
				Plugin: PluginFunc(func(pluginName string, update Update, speaker Speaker) {
					log.Infof("UpdateHandler not configured.")
				}),
				Terminate: false,
			},
			{
				Name: "UpdateHandlerConfig2",
				Plugin: PluginFunc(func(pluginName string, update Update, speaker Speaker) {
					log.Infof("UpdateHandler not configured.")
				}),
				Terminate: false,
			},
		},
	}

	fDefault := fields{
		Plugins: []PluginConfig{
			{
				Name: "NotConfigured",
				Plugin: PluginFunc(func(pluginName string, update Update, speaker Speaker) {
					log.Infof("UpdateHandler not configured.")
				}),
				Terminate: false,
			},
		},
	}
	uhc1 := PluginConfig{
		Name: "NewUpdateHandler",
		Plugin: PluginFunc(func(pluginName string, update Update, speaker Speaker) {
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
				Plugins: []PluginConfig{},
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
				UpdateHandlers: tt.fields.Plugins,
			}
			s.AddPlugin(tt.args.uhc)
			if len(s.UpdateHandlers) != tt.lenAfter {
				t.Errorf("Speaker.AddUpdateHandler %v failed. length want %d is %v ", tt.name, tt.lenAfter, len(s.UpdateHandlers))
			}
			if !s.HasPlugin(tt.args.uhc.Name) {
				t.Errorf("Speaker.AddUpdateHandler %v failed. %v not included.", tt.name, tt.args.uhc.Name)
			}
		})
	}
}

func TestNewSpeaker(t *testing.T) {

	type args struct {
		entry *mdns.ServiceEntry
	}
	tests := []struct {
		name string
		args args
		want *Speaker
	}{

		{
			"New nil ServiceEntry",
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

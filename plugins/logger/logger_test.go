package logger

import (
	"net"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/soundtouch-golang"
)

func TestSpeaker_AddUpdateHandler(t *testing.T) {
	type fields struct {
		IP           net.IP
		Port         int
		BaseHTTPURL  url.URL
		WebSocketURL url.URL
		DeviceInfo   soundtouch.Info
		conn         *websocket.Conn
		webSocketCh  chan *soundtouch.Update
		Plugins      []soundtouch.PluginConfig
	}

	f1 := fields{
		Plugins: []soundtouch.PluginConfig{
			{
				Name: "UpdateHandlerConfig1",
				Pfunction: soundtouch.PluginFunc(func(pluginName string, update soundtouch.Update, speaker soundtouch.Speaker) {
					log.Infof("UpdateHandler not configured.")
				}),
				Terminate: false,
			},
		},
	}

	f2 := fields{
		Plugins: []soundtouch.PluginConfig{
			{
				Name: "UpdateHandlerConfig1",
				Pfunction: soundtouch.PluginFunc(func(pluginName string, update soundtouch.Update, speaker soundtouch.Speaker) {
					log.Infof("UpdateHandler not configured.")
				}),
				Terminate: false,
			},
			{
				Name: "UpdateHandlerConfig2",
				Pfunction: soundtouch.PluginFunc(func(pluginName string, update soundtouch.Update, speaker soundtouch.Speaker) {
					log.Infof("UpdateHandler not configured.")
				}),
				Terminate: false,
			},
		},
	}

	fDefault := fields{
		Plugins: []soundtouch.PluginConfig{
			{
				Name: "NotConfigured",
				Pfunction: soundtouch.PluginFunc(func(pluginName string, update soundtouch.Update, speaker soundtouch.Speaker) {
					log.Infof("UpdateHandler not configured.")
				}),
				Terminate: false,
			},
		},
	}
	uhc1 := NewLogger(Config{})
	tests := []struct {
		name     string
		fields   fields
		lenAfter int
	}{
		{
			"Add 1 to 1",
			f1,
			2,
		},
		{
			"Add 1 to 2",
			f2,
			3,
		},
		{
			"Add 1 to empty",
			fields{
				Plugins: []soundtouch.PluginConfig{},
			},
			1,
		},
		{
			"Add 1 to default",
			fDefault,
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &soundtouch.Speaker{
				IP:           tt.fields.IP,
				Port:         tt.fields.Port,
				BaseHTTPURL:  tt.fields.BaseHTTPURL,
				WebSocketURL: tt.fields.WebSocketURL,
				DeviceInfo:   tt.fields.DeviceInfo,
			}
			s.AddPlugin(uhc1)
			if len(s.Plugins) != tt.lenAfter {
				t.Errorf("Speaker.AddUpdateHandler %v failed. length want %d is %v ", tt.name, tt.lenAfter, len(s.Plugins))
			}
			if !s.HasPlugin(uhc1.Name()) {
				t.Errorf("Speaker.AddUpdateHandler %v failed. %v not included.", tt.name, uhc1.Name())
			}
		})
	}
}

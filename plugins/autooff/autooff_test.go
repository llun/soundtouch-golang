package autooff

import (
	"io/ioutil"
	"net"
	"net/url"
	"testing"

	"github.com/theovassiliou/soundtouch-golang"
)

func TestCollector_Execute(t *testing.T) {
	wohnzimmer := soundtouch.Speaker{
		IP:   net.IP{},
		Port: -1,
		BaseHTTPURL: url.URL{
			Scheme: "http",
			Host:   "192.168.178.52:8090",
		},
		WebSocketURL: url.URL{},
		DeviceInfo: soundtouch.Info{
			DeviceID: "aabbccspeaker1",
			Name:     "Wohnzimmer",
			Type:     "",
			Raw:      nil,
		},
		WebSocketCh: make(chan *soundtouch.Update),
		Plugins:     nil,
	}

	c1 := map[string]struct {
		ThenOff []string "toml:\"thenOff\""
	}{
		"Wohmzimmer": {
			ThenOff: []string{"Schrank", "Küche}"},
		},
	}

	u1, _ := soundtouch.NewUpdate(raw(t, "../../testfiles/updates_nowPlaying_Soundbar_Products.xml"))

	type fields struct {
		Config    Config
		Plugin    soundtouch.PluginFunc
		suspended bool
	}
	type args struct {
		pluginName string
		update     soundtouch.Update
		speaker    soundtouch.Speaker
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "first test",
			fields: fields{
				Config:    c1,
				Plugin:    func(string, soundtouch.Update, soundtouch.Speaker) { panic("not implemented") },
				suspended: false,
			},
			args: args{
				pluginName: "",
				update:     *u1,
				speaker:    wohnzimmer,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Collector{
				Config:    tt.fields.Config,
				Plugin:    tt.fields.Plugin,
				suspended: tt.fields.suspended,
			}
			d.Execute(tt.args.pluginName, tt.args.update, tt.args.speaker)
		})
	}
}

func raw(t *testing.T, filename string) []byte {
	b, err := ioutil.ReadFile(filename) // just pass the file name
	if err != nil {
		t.Errorf("readign raw() error = %v, filename %v", err, filename)
	}

	return b
}

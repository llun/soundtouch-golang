package episodecollector

import (
	"net/url"
	"reflect"
	"testing"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/theovassiliou/soundtouch-golang"
)

func TestCollector_Name(t *testing.T) {
	type fields struct {
		Config     Config
		Plugin     soundtouch.PluginFunc
		suspended  bool
		scribbleDb *scribble.Driver
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{{
		name: "Get Name",
		fields: fields{
			Config: Config{
				Speakers:  nil,
				Terminate: false,
				Artists:   nil,
				Database:  "",
			},
			Plugin:     func(string, soundtouch.Update, soundtouch.Speaker) { panic("not implemented") },
			suspended:  false,
			scribbleDb: &scribble.Driver{},
		},
		want: "EpisodeCollector",
	},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Collector{
				Config:     tt.fields.Config,
				Plugin:     tt.fields.Plugin,
				suspended:  tt.fields.suspended,
				scribbleDb: tt.fields.scribbleDb,
			}
			if got := d.Name(); got != tt.want {
				t.Errorf("Collector.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollector_Terminate(t *testing.T) {
	type fields struct {
		Config     Config
		Plugin     soundtouch.PluginFunc
		suspended  bool
		scribbleDb *scribble.Driver
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{{
		name: "Get Terminate false",
		fields: fields{
			Config: Config{
				Speakers:  nil,
				Terminate: false,
				Artists:   nil,
				Database:  "",
			},
			Plugin:     func(string, soundtouch.Update, soundtouch.Speaker) { panic("not implemented") },
			suspended:  false,
			scribbleDb: &scribble.Driver{},
		},
		want: false,
	},
		{
			name: "Get Terminate true",
			fields: fields{
				Config: Config{
					Speakers:  nil,
					Terminate: true,
					Artists:   nil,
					Database:  "",
				},
				Plugin:     func(string, soundtouch.Update, soundtouch.Speaker) { panic("not implemented") },
				suspended:  false,
				scribbleDb: &scribble.Driver{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Collector{
				Config:     tt.fields.Config,
				Plugin:     tt.fields.Plugin,
				suspended:  tt.fields.suspended,
				scribbleDb: tt.fields.scribbleDb,
			}
			if got := d.Terminate(); got != tt.want {
				t.Errorf("Collector.Terminate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollector_Execute(t *testing.T) {

	u1, _ := soundtouch.NewUpdate([]byte("<updates deviceID=\"08DF1F0E9E36\"><volumeUpdated><volume><targetvolume>27</targetvolume><actualvolume>27</actualvolume><muteenabled>false</muteenabled></volume></volumeUpdated></updates>"))

	type fields struct {
		Config     Config
		Plugin     soundtouch.PluginFunc
		suspended  bool
		scribbleDb *scribble.Driver
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
	}{{
		name: "A name",
		fields: fields{
			Config: Config{
				Speakers:  nil,
				Terminate: false,
				Artists:   []string{"AnArtist"},
				Database:  "",
			},
			Plugin:     func(string, soundtouch.Update, soundtouch.Speaker) { panic("not implemented") },
			suspended:  false,
			scribbleDb: &scribble.Driver{},
		},
		args: args{
			pluginName: "",
			update: soundtouch.Update{
				DeviceID: "08DF1F0E9E36",
				Value:    u1,
			},
			speaker: soundtouch.Speaker{
				IP:   nil,
				Port: 0,
				BaseHTTPURL: url.URL{
					Scheme:      "",
					Opaque:      "",
					User:        &url.Userinfo{},
					Host:        "",
					Path:        "",
					RawPath:     "",
					ForceQuery:  false,
					RawQuery:    "",
					Fragment:    "",
					RawFragment: "",
				},
				WebSocketURL: url.URL{
					Scheme:      "",
					Opaque:      "",
					User:        &url.Userinfo{},
					Host:        "",
					Path:        "",
					RawPath:     "",
					ForceQuery:  false,
					RawQuery:    "",
					Fragment:    "",
					RawFragment: "",
				},
				DeviceInfo: soundtouch.Info{
					DeviceID:  "08DF1F0E9E36",
					Name:      "SpeakerA",
					Type:      "",
					IPAddress: nil,
					Raw:       nil,
				},
				WebSocketCh: make(chan *soundtouch.Update),
				Plugins:     nil,
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Collector{
				Config:     tt.fields.Config,
				Plugin:     tt.fields.Plugin,
				suspended:  tt.fields.suspended,
				scribbleDb: tt.fields.scribbleDb,
			}
			d.Execute(tt.args.pluginName, tt.args.update, tt.args.speaker)
		})
	}
}

func TestNewCollector(t *testing.T) {
	type args struct {
		config Config
	}
	tests := []struct {
		name  string
		args  args
		wantD *Collector
	}{
		{
			name: "Valid constructor",
			args: args{
				config: Config{
					Speakers:  nil,
					Terminate: false,
					Artists:   nil,
					Database:  "",
				},
			},
			wantD: &Collector{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotD := NewCollector(tt.args.config); !reflect.DeepEqual(gotD, tt.wantD) {
				t.Errorf("NewCollector() = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}

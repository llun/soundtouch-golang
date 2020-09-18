package soundtouch

import (
	"io/ioutil"
	"net"
	"net/url"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/gorilla/websocket"
)

func TestSpeaker_NowPlaying(t *testing.T) {

	type fields struct {
		IP           net.IP
		Port         int
		BaseHTTPURL  url.URL
		WebSocketURL url.URL
		DeviceInfo   Info
		conn         *websocket.Conn
	}

	f1 := fields{
		net.IP{},
		-1,
		url.URL{
			Scheme: "http",
			Host:   "192.168.178.52:8090",
		},
		url.URL{},
		Info{},
		nil,
	}

	tests := []struct {
		name         string
		fixtureName  string
		responseCode int
		method       string
		route        string
		fields       fields
		want         NowPlaying
		wantErr      bool
	}{
		{
			name:         "now_playing Buffering InternetRadio",
			fixtureName:  "./testfiles/nowPlaying_InternetRadio_Buffering.xml",
			responseCode: 200,
			method:       "GET",
			route:        "http://192.168.178.52:8090/now_playing",
			fields:       f1,
			want: NowPlaying{
				PlayStatus: BufferingState,
				Source:     LocalInternetRadio,
				DeviceID:   "08DF1F117BB7",
				Content: ContentItem{
					Source:       "LOCAL_INTERNET_RADIO",
					Type:         "stationurl",
					Name:         "Radio Eins",
					Location:     "https://content.api.bose.io/core02/svc-bmx-adapter-orion/prod/orion/station?data=eyJuYW1lIjoiUmFkaW8gRWlucyIsImltYWdlVXJsIjoiIiwic3RyZWFtVXJsIjoiaHR0cDovL3d3dy5yYWRpb2VpbnMuZGUvbGl2ZS5tM3UifQ%3D%3D",
					IsPresetable: true,
				},
				Track:      "Radio Eins",
				StreamType: "RADIO_STREAMING",
				Raw:        raw(t, "./testfiles/nowPlaying_InternetRadio_Buffering.xml"),
			},
		},
		{
			name:         "pausing playing stored music",
			fixtureName:  "./testfiles/nowPlayingStoredMusic_Pause.xml",
			responseCode: 200,
			method:       "GET",
			route:        "http://192.168.178.52:8090/now_playing",
			fields:       f1,
			want: NowPlaying{
				PlayStatus:    PauseState,
				Source:        StoredMusic,
				SourceAccount: "44067f6f-6b79-2d88-b531-11189bea9cd7/0",
				DeviceID:      "0122FF1A1234",
				Content: ContentItem{
					Source:       "STORED_MUSIC",
					Name:         "Ploutarchos 2",
					Location:     "0$1$18$213$12377$12389",
					IsPresetable: true,
				},
				Track:  "Xrwma tis zwis",
				Artist: "Giannis Ploutarxos",
				Album:  "Ploutarchos 2",
				Art:    "http://192.168.178.50:9000/disk/DLNA-PNJPEG_TN-OP01-CI1-FLAGS00d00000/defaa/A/O0$1$8I2084108.jpg?scale=org",
				Raw:    raw(t, "./testfiles/nowPlayingStoredMusic_Pause.xml"),
			},
		},
		{
			name:         "now_playing InternetRadio",
			fixtureName:  "./testfiles/nowPlaying_InternetRadio_Playing.xml",
			responseCode: 200,
			method:       "GET",
			route:        "http://192.168.178.52:8090/now_playing",
			fields:       f1,
			want: NowPlaying{
				PlayStatus: PlayState,
				Source:     LocalInternetRadio,
				DeviceID:   "08DF1F117BB7",
				Content: ContentItem{
					Source:       "LOCAL_INTERNET_RADIO",
					Type:         "stationurl",
					Name:         "Radio Eins",
					Location:     "https://content.api.bose.io/core02/svc-bmx-adapter-orion/prod/orion/station?data=eyJuYW1lIjoiUmFkaW8gRWlucyIsImltYWdlVXJsIjoiIiwic3RyZWFtVXJsIjoiaHR0cDovL3d3dy5yYWRpb2VpbnMuZGUvbGl2ZS5tM3UifQ%3D%3D",
					IsPresetable: true,
				},
				Track:      "Radio Eins",
				StreamType: "RADIO_STREAMING",
				Raw:        raw(t, "./testfiles/nowPlaying_InternetRadio_Playing.xml"),
			},
		},
		{
			name:         "now_playing TuneIn",
			fixtureName:  "./testfiles/nowPlaying_TuneInPlaying.xml",
			responseCode: 200,
			method:       "GET",
			route:        "http://192.168.178.52:8090/now_playing",
			fields:       f1,
			want: NowPlaying{
				PlayStatus: PlayState,
				Source:     TuneIn,
				DeviceID:   "08DF1F117BB7",
				Content: ContentItem{
					Source:       "TUNEIN",
					Type:         "stationurl",
					Name:         "105 5 Spreeradio",
					Location:     "/v1/playback/station/s17211",
					IsPresetable: true,
				},
				Track:      "Spreeradio Live",
				Artist:     "Die besten Songs von den 80ern bis heute.",
				Art:        "http://cdn-profiles.tunein.com/s17211/images/logoq.png?t=1",
				StreamType: "RADIO_STREAMING",
				Raw:        raw(t, "./testfiles/nowPlaying_TuneInPlaying.xml"),
			},
		},
		{
			name:         "Aux playing",
			fixtureName:  "./testfiles/nowPlaying_AuxPlaying.xml",
			responseCode: 200,
			method:       "GET",
			route:        "http://192.168.178.52:8090/now_playing",
			fields:       f1,
			want: NowPlaying{
				PlayStatus: PlayState,
				Source:     Aux,
				DeviceID:   "08DF1F117BB7",
				Content: ContentItem{
					Source:       "AUX",
					Name:         "AUX IN",
					IsPresetable: true,
				},
				Raw: raw(t, "./testfiles/nowPlaying_AuxPlaying.xml"),
			},
		},
		{
			name:         "Soundbar playing TV",
			fixtureName:  "./testfiles/nowPlaying_SoundbarTV_Playing.xml",
			responseCode: 200,
			method:       "GET",
			route:        "http://192.168.178.52:8090/now_playing",
			fields:       f1,
			want: NowPlaying{
				PlayStatus:    PlayState,
				Source:        "PRODUCT",
				SourceAccount: "TV",
				DeviceID:      "9884E39B34BE",
				Content: ContentItem{
					Source:       "PRODUCT",
					IsPresetable: false,
				},
				Raw: raw(t, "./testfiles/nowPlaying_SoundbarTV_Playing.xml"),
			},
		},
		{
			name:         "Soundbar searching Bluetooth",
			fixtureName:  "./testfiles/nowPlaying_Soundbar_BT_Searching.xml",
			responseCode: 200,
			method:       "GET",
			route:        "http://192.168.178.52:8090/now_playing",
			fields:       f1,
			want: NowPlaying{
				PlayStatus: InvalidPlayStatus,
				Source:     "BLUETOOTH",
				DeviceID:   "9884E39B34BE",
				Content: ContentItem{
					Source:       "BLUETOOTH",
					IsPresetable: false,
				},
				Raw: raw(t, "./testfiles/nowPlaying_Soundbar_BT_Searching.xml"),
			},
		},

		{
			name:         "standby",
			fixtureName:  "./testfiles/nowPlayingStoredMusic_Standby.xml",
			responseCode: 200,
			method:       "GET",
			route:        "http://192.168.178.52:8090/now_playing",
			fields:       f1,
			want: NowPlaying{
				Source:   Standby,
				DeviceID: "08DF1F0E9E36",
				Content: ContentItem{
					Source: "STANDBY",
				},
				Raw: raw(t, "./testfiles/nowPlayingStoredMusic_Standby.xml"),
			},
		},
		{
			name:         "playing stored music",
			fixtureName:  "./testfiles/nowPlayingStoredMusic_Play.xml",
			responseCode: 200,
			method:       "GET",
			route:        "http://192.168.178.52:8090/now_playing",
			fields:       f1,
			want: NowPlaying{
				PlayStatus:    PlayState,
				Source:        StoredMusic,
				SourceAccount: "44067f6f-6b79-2d88-b531-11189bea9cd7/0",
				DeviceID:      "0122FF1A1234",
				Content: ContentItem{
					Source:       "STORED_MUSIC",
					Name:         "Platinum Edition",
					Location:     "0$1$18$213$18724$18726",
					IsPresetable: true,
				},
				Track:  "sorry",
				Artist: "Elli Kokkinou",
				Album:  "Platinum Edition",
				Art:    "http://192.168.178.50:9000/disk/DLNA-PNPNG_TN-OP01-CI1-FLAGS00d00000/defaa/A/O0$1$8I3478028.png?scale=org",
				Raw:    raw(t, "./testfiles/nowPlayingStoredMusic_Play.xml"),
			},
		},
		{
			name:         "pausing Bluetooth music",
			fixtureName:  "./testfiles/nowPlaying_Soundbar_BT_Pausing.xml",
			responseCode: 200,
			method:       "GET",
			route:        "http://192.168.178.52:8090/now_playing",
			fields:       f1,
			want: NowPlaying{
				PlayStatus:    PauseState,
				Source:        Bluetooth,
				SourceAccount: "",
				DeviceID:      "9884E39B34BE",
				Content: ContentItem{
					Source:       "BLUETOOTH",
					Name:         "iPhone von Theofanis",
					Location:     "",
					IsPresetable: false,
				},
				Track:  "Once in a Lifetime - 2005 Remaster",
				Artist: "Talking Heads",
				Album:  "Remain in Light (Deluxe Version)",
				Raw:    raw(t, "./testfiles/nowPlaying_Soundbar_BT_Pausing.xml"),
			},
		},

		{
			name:         "playing Bluetooth music",
			fixtureName:  "./testfiles/nowPlaying_Soundbar_BT_Playing.xml",
			responseCode: 200,
			method:       "GET",
			route:        "http://192.168.178.52:8090/now_playing",
			fields:       f1,
			want: NowPlaying{
				PlayStatus:    PlayState,
				Source:        Bluetooth,
				SourceAccount: "",
				DeviceID:      "9884E39B34BE",
				Content: ContentItem{
					Source:       "BLUETOOTH",
					Name:         "iPhone von Theofanis",
					Location:     "",
					IsPresetable: false,
				},
				Track:  "Once in a Lifetime - 2005 Remaster",
				Artist: "Talking Heads",
				Album:  "Remain in Light (Deluxe Version)",
				Raw:    raw(t, "./testfiles/nowPlaying_Soundbar_BT_Playing.xml"),
			},
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fixture, err := ioutil.ReadFile(tt.fixtureName)
			if err != nil {
				t.Errorf("fixture not found %v", err)
			}
			responder := httpmock.NewBytesResponder(tt.responseCode, fixture)
			httpmock.RegisterResponder(tt.method, tt.route, responder)

			s := &Speaker{
				IP:           tt.fields.IP,
				Port:         tt.fields.Port,
				BaseHTTPURL:  tt.fields.BaseHTTPURL,
				WebSocketURL: tt.fields.WebSocketURL,
				DeviceInfo:   tt.fields.DeviceInfo,
				conn:         tt.fields.conn,
			}
			got, err := s.NowPlaying()
			if (err != nil) != tt.wantErr {
				t.Errorf("Speaker.NowPlaying() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			AssertNowPlayingEquals(t, got, tt.want)
		})
	}
}

func AssertEqual(t *testing.T, a interface{}, b interface{}) {

	if a == b {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

func AssertEqualDeep(t *testing.T, got interface{}, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Received %v (type %v), expected %v (type %v)", got, reflect.TypeOf(got), want, reflect.TypeOf(want))
	}
}

func AssertNowPlayingEquals(t *testing.T, got, want NowPlaying) {
	AssertEqual(t, got.PlayStatus, want.PlayStatus)
	AssertEqual(t, got.Source, want.Source)
	AssertEqual(t, got.SourceAccount, want.SourceAccount)
	AssertEqual(t, got.Content.Type, want.Content.Type)
	AssertEqual(t, got.Content.Name, want.Content.Name)
	AssertEqual(t, got.Content.Location, want.Content.Location)
	AssertEqual(t, got.Content.IsPresetable, want.Content.IsPresetable)
	AssertEqual(t, got.Track, want.Track)
	AssertEqual(t, got.Artist, want.Artist)
	AssertEqual(t, got.Album, want.Album)
	AssertEqual(t, got.TrackID, want.TrackID)
	AssertEqual(t, got.Art, want.Art)
	AssertEqual(t, got.StreamType, want.StreamType)
	AssertEqualDeep(t, got.Raw, want.Raw)

	// Last check whether we covered all fields
	AssertEqualDeep(t, got, want)
}

func raw(t *testing.T, filename string) []byte {
	b, err := ioutil.ReadFile(filename) // just pass the file name
	if err != nil {
		t.Errorf("readign raw() error = %v, filename %v", err, filename)
	}

	return b
}

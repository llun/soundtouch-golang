package soundtouch

import (
	"io/ioutil"
	"net"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/gorilla/websocket"
)

func TestSpeaker_Presets(t *testing.T) {

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
		want         Presets
		wantErr      bool
	}{
		{
			name:         "now_playing Buffering InternetRadio",
			fixtureName:  "./testfiles/presets.xml",
			responseCode: 200,
			method:       "GET",
			route:        "http://192.168.178.52:8090/presets",
			fields:       f1,
			want: Presets{
				DeviceID: "",
				Presets: []Preset{
					{
						ID: 1,
						Content: ContentItem{
							Source:       "LOCAL_INTERNET_RADIO",
							Type:         "stationurl",
							Name:         "Radio Eins",
							Location:     "https://content.api.bose.io/core02/svc-bmx-adapter-orion/prod/orion/station?data=eyJuYW1lIjoiUmFkaW8gRWlucyIsImltYWdlVXJsIjoiIiwic3RyZWFtVXJsIjoiaHR0cDovL3d3dy5yYWRpb2VpbnMuZGUvbGl2ZS5tM3UifQ%3D%3D",
							IsPresetable: true,
						},
					},
					{
						ID: 2,
						Content: ContentItem{
							Source:       "AMAZON",
							Type:         "tracklist",
							Name:         "Deutscher Pop",
							Location:     "catalog/stations/A2XJSDTOSE5YW5/#playable",
							IsPresetable: true,
						},
					},
					{
						ID: 3,
						Content: ContentItem{
							Source:       "TUNEIN",
							Type:         "stationurl",
							Name:         "105 5 Spreeradio",
							Location:     "/v1/playback/station/s17211",
							IsPresetable: true,
						},
					},
					{
						ID: 4,
						Content: ContentItem{
							Source:       "STORED_MUSIC",
							Type:         "",
							Name:         "066 Die Schattenmaenner",
							Location:     "0$1$18$1066$19053$20006",
							IsPresetable: true,
						},
					},

					{
						ID: 5,
						Content: ContentItem{
							Source:       "STORED_MUSIC",
							Type:         "",
							Name:         "059 - Giftiges Wasser",
							Location:     "0$1$17$1066$19779",
							IsPresetable: true,
						},
					},
					{
						ID: 6,
						Content: ContentItem{
							Source:       "AMAZON",
							Type:         "tracklist",
							Name:         "Songs",
							Location:     "library/albums/7ec7e321-0067-4f1f-bbf0-0c90bd6e5693/#playable",
							IsPresetable: true,
						},
					},
				},
				Raw: raw(t, "./testfiles/presets.xml"),
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
			got, err := s.Presets()
			if (err != nil) != tt.wantErr {
				t.Errorf("Speaker.NowPlaying() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			AssertPresetsEquals(t, got, tt.want)
		})
	}
}

func AssertPresetsEquals(t *testing.T, got, want Presets) {
	AssertEqual(t, got.DeviceID, want.DeviceID)
	for i, p := range want.Presets {
		AssertEqual(t, got.Presets[i].ID, p.ID)
		AssertEqual(t, got.Presets[i].Content.Type, p.Content.Type)
		AssertEqual(t, got.Presets[i].Content.Name, p.Content.Name)
		AssertEqual(t, got.Presets[i].Content.Location, p.Content.Location)
	}
	AssertEqualDeep(t, got.Raw, want.Raw)

	// Last check whether we covered all fields
	AssertEqualDeep(t, got, want)
}

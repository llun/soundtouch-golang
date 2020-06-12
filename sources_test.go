package soundtouch

import (
	"io/ioutil"
	"net"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/gorilla/websocket"
)

func TestSpeaker_Sources(t *testing.T) {

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
		want         Sources
		wantErr      bool
	}{
		{
			name:         "Sources",
			fixtureName:  "./testfiles/sources.xml",
			responseCode: 200,
			method:       "GET",
			route:        "http://192.168.178.52:8090/sources",
			fields:       f1,
			want: Sources{
				DeviceID: "08DF1F1A065C",
				SourceItems: []SourceItem{
					{
						Source:        "STORED_MUSIC",
						SourceAccount: "55076f6e-6b79-1d65-a471-00089bea8bd7/0",
						Status:        "READY",
						Local:         false,
						Value:         "TwonkyServer [maxi]",
					},
					{
						Source:        "AIRPLAY",
						SourceAccount: "",
						Status:        "READY",
						Local:         false,
						Value:         "",
					},
					{
						Source:        "AMAZON",
						SourceAccount: "nicky.pohl@googlemail.com",
						Status:        "READY",
						Local:         false,
						Value:         "nicky.pohl@googlemail.com",
					},
					{
						Source:        "NOTIFICATION",
						SourceAccount: "",
						Status:        "UNAVAILABLE",
						Local:         false,
						Value:         "",
					},
					{
						Source:        "STORED_MUSIC_MEDIA_RENDERER",
						SourceAccount: "StoredMusicUserName",
						Status:        "UNAVAILABLE",
						Local:         false,
						Value:         "StoredMusicUserName",
					},
					{
						Source:        "QPLAY",
						SourceAccount: "QPlay1UserName",
						Status:        "UNAVAILABLE",
						Local:         true,
						Value:         "QPlay1UserName",
					},
				},
				Raw: raw(t, "./testfiles/sources.xml"),
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
			got, err := s.Sources()
			if (err != nil) != tt.wantErr {
				t.Errorf("Speaker.Sources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			AssertSourcesEquals(t, got, tt.want)
		})
	}
}

func AssertSourcesEquals(t *testing.T, got, want Sources) {
	AssertEqual(t, got.DeviceID, want.DeviceID)
	for i, p := range want.SourceItems {
		AssertEqual(t, got.SourceItems[i].Source, p.Source)
		AssertEqual(t, got.SourceItems[i].SourceAccount, p.SourceAccount)
		AssertEqual(t, got.SourceItems[i].Status, p.Status)
		AssertEqual(t, got.SourceItems[i].Local, p.Local)
		AssertEqual(t, got.SourceItems[i].Value, p.Value)
	}

	AssertEqualDeep(t, got.Raw, want.Raw)

	// Last check whether we covered all fields
	AssertEqualDeep(t, got, want)
}

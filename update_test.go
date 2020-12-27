package soundtouch

import (
	"reflect"
	"testing"
)

func TestNewUpdate(t *testing.T) {
	type args struct {
		body []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Update
		wantErr bool
	}{
		{
			name: "First Test",
			args: args{
				raw(t, "testfiles/updates_nowPlaying_Soundbar_Products.xml"),
			},
			want: &Update{"9884E39B34BE", NowPlaying{
				PlayStatus:    "PLAY_STATE",
				Source:        "PRODUCT",
				SourceAccount: "TV",
				DeviceID:      "9884E39B34BE",
				Content: ContentItem{
					Type:         "",
					Source:       "PRODUCT",
					Location:     "",
					Name:         "",
					IsPresetable: false,
				},
				Track:      "",
				Artist:     "",
				Album:      "",
				TrackID:    "",
				Art:        "",
				StreamType: "",
				Raw:        nil,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUpdate(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

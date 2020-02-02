package soundtouch

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

const lineProtoFmtCSU = "wifi,name=\"%s\",deviceID=\"%s\" wifiStrength=%v,connected=\"%v\""
const lineProtoFmtNP = "playing,name=\"%s\",deviceID=\"%s\" playStatus=%v,album=\"%v\""
const lineProtoFmtVU = "playing,name=\"%s\",deviceID=\"%s\" volume=%v"

var strengthMapping = map[string]int{
	"EXCELLENT_SIGNAL": 100, "GOOD_SIGNAL": 70, "POOR_SIGNAL": 30, "MARGINAL_SIGNAL": 10,
}

var playStateMapping = map[PlayStatus]int{
	"PLAY_STATE": 1, "PAUSE_STATE": 2, "STOP_STATE": 3, "STANDBY": 5, "BUFFERING_STATE": 8, "INVALID_PLAY_STATUS": 13,
}

type InfluxDB struct {
	BaseHTTPURL       url.URL
	Database          string
	SoundtouchNetwork map[string]string
}

func (i *InfluxDB) SetData(action string, input []byte) ([]byte, error) {
	actionURL, _ := url.Parse(i.WriteURL())
	buffer := bytes.NewBuffer(input)
	// log.Debugf("Going to send action: %v, %v", action, string(input))
	resp, err := http.Post(actionURL.String(), "", buffer)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (i *InfluxDB) WriteURL() string {
	return fmt.Sprintf("%v/write?db=%v", i.BaseHTTPURL.String(), i.Database)
}

func (u *Update) Lineproto(i InfluxDB, message *Update) (string, error) {
	typeName := reflect.TypeOf(message.Value).Name()
	switch typeName {
	case "ConnectionStateUpdated":
		c, _ := message.Value.(ConnectionStateUpdated)
		return c.Lineproto(i, message)
	case "NowPlaying":
		np, _ := message.Value.(NowPlaying)
		return np.Lineproto(i, message)
	case "Volume":
		v, _ := message.Value.(Volume)
		return v.Lineproto(i, message)
	default:
		return "", fmt.Errorf("lineproto: no lineproto for this Update-type %v", typeName)

	}
}

func (s *ConnectionStateUpdated) Lineproto(i InfluxDB, message *Update) (string, error) {
	lineproto := fmt.Sprintf(lineProtoFmtCSU,
		i.SoundtouchNetwork[message.DeviceId],
		message.DeviceId,
		strengthMapping[s.Signal],
		func() string {
			if s.Up == "true" {
				return "UP"
			}
			return "DOWN"
		}())
	return lineproto, nil
}

func (v *Volume) Lineproto(i InfluxDB, message *Update) (string, error) {
	lineproto := fmt.Sprintf(lineProtoFmtVU,
		i.SoundtouchNetwork[message.DeviceId],
		message.DeviceId,
		v.TargetVolume,
	)
	return lineproto, nil
}

func (s *NowPlaying) Lineproto(i InfluxDB, message *Update) (string, error) {
	lineproto := fmt.Sprintf(lineProtoFmtNP,
		i.SoundtouchNetwork[message.DeviceId],
		message.DeviceId,
		func() int {
			ps := playStateMapping[s.PlayStatus]
			if ps == 0 && s.Source == "STANDBY" {
				return playStateMapping["STANDBY"]
			}
			return ps
		}(),
		func() string {
			if s.Album == "" {
				return "none"
			}
			return s.Album
		}())
	return lineproto, nil
}
